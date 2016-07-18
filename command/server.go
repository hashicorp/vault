package command

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/flag-slice"
	"github.com/hashicorp/vault/helper/gated-writer"
	"github.com/hashicorp/vault/helper/mlock"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

// ServerCommand is a Command that starts the Vault server.
type ServerCommand struct {
	AuditBackends      map[string]audit.Factory
	CredentialBackends map[string]logical.Factory
	LogicalBackends    map[string]logical.Factory

	ShutdownCh chan struct{}
	SighupCh   chan struct{}

	meta.Meta

	logger *log.Logger

	ReloadFuncs map[string][]server.ReloadFunc
}

func (c *ServerCommand) Run(args []string) int {
	var dev, verifyOnly bool
	var configPath []string
	var logLevel, devRootTokenID, devListenAddress string
	flags := c.Meta.FlagSet("server", meta.FlagSetDefault)
	flags.BoolVar(&dev, "dev", false, "")
	flags.StringVar(&devRootTokenID, "dev-root-token-id", "", "")
	flags.StringVar(&devListenAddress, "dev-listen-address", "", "")
	flags.StringVar(&logLevel, "log-level", "info", "")
	flags.BoolVar(&verifyOnly, "verify-only", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	flags.Var((*sliceflag.StringFlag)(&configPath), "config", "config")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if os.Getenv("VAULT_DEV_ROOT_TOKEN_ID") != "" && devRootTokenID == "" {
		devRootTokenID = os.Getenv("VAULT_DEV_ROOT_TOKEN_ID")
	}

	if os.Getenv("VAULT_DEV_LISTEN_ADDRESS") != "" && devListenAddress == "" {
		devListenAddress = os.Getenv("VAULT_DEV_LISTEN_ADDRESS")
	}

	// Validation
	if !dev {
		switch {
		case len(configPath) == 0:
			c.Ui.Error("At least one config path must be specified with -config")
			flags.Usage()
			return 1
		case devRootTokenID != "":
			c.Ui.Error("Root token ID can only be specified with -dev")
			flags.Usage()
			return 1
		case devListenAddress != "":
			c.Ui.Error("Development address can only be specified with -dev")
			flags.Usage()
			return 1
		}
	}

	// Load the configuration
	var config *server.Config
	if dev {
		config = server.DevConfig()
		if devListenAddress != "" {
			config.Listeners[0].Config["address"] = devListenAddress
		}
	}
	for _, path := range configPath {
		current, err := server.LoadConfig(path)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error loading configuration from %s: %s", path, err))
			return 1
		}

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}

	// Ensure at least one config was found.
	if config == nil {
		c.Ui.Error("No configuration files found.")
		return 1
	}

	// Ensure that a backend is provided
	if config.Backend == nil {
		c.Ui.Error("A physical backend must be specified")
		return 1
	}

	// If mlockall(2) isn't supported, show a warning.  We disable this
	// in dev because it is quite scary to see when first using Vault.
	if !dev && !mlock.Supported() {
		c.Ui.Output("==> WARNING: mlock not supported on this system!\n")
		c.Ui.Output("  An `mlockall(2)`-like syscall to prevent memory from being")
		c.Ui.Output("  swapped to disk is not supported on this system. Running")
		c.Ui.Output("  Vault on an mlockall(2) enabled system is much more secure.\n")
	}

	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	logGate := &gatedwriter.Writer{Writer: os.Stderr}
	c.logger = log.New(&logutils.LevelFilter{
		Levels: []logutils.LogLevel{
			"TRACE", "DEBUG", "INFO", "WARN", "ERR"},
		MinLevel: logutils.LogLevel(strings.ToUpper(logLevel)),
		Writer:   logGate,
	}, "", log.LstdFlags)

	if err := c.setupTelemetry(config); err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}

	// Initialize the backend
	backend, err := physical.NewBackend(
		config.Backend.Type, c.logger, config.Backend.Config)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing backend of type %s: %s",
			config.Backend.Type, err))
		return 1
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)

	var seal vault.Seal = &vault.DefaultSeal{}

	// Ensure that the seal finalizer is called, even if using verify-only
	defer func() {
		err = seal.Finalize()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error finalizing seals: %v", err))
		}
	}()

	coreConfig := &vault.CoreConfig{
		Physical:           backend,
		AdvertiseAddr:      config.Backend.AdvertiseAddr,
		HAPhysical:         nil,
		Seal:               seal,
		AuditBackends:      c.AuditBackends,
		CredentialBackends: c.CredentialBackends,
		LogicalBackends:    c.LogicalBackends,
		Logger:             c.logger,
		DisableCache:       config.DisableCache,
		DisableMlock:       config.DisableMlock,
		MaxLeaseTTL:        config.MaxLeaseTTL,
		DefaultLeaseTTL:    config.DefaultLeaseTTL,
	}

	// Initialize the separate HA physical backend, if it exists
	var ok bool
	if config.HABackend != nil {
		habackend, err := physical.NewBackend(
			config.HABackend.Type, c.logger, config.HABackend.Config)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error initializing backend of type %s: %s",
				config.HABackend.Type, err))
			return 1
		}

		if coreConfig.HAPhysical, ok = habackend.(physical.HABackend); !ok {
			c.Ui.Error("Specified HA backend does not support HA")
			return 1
		}

		if !coreConfig.HAPhysical.HAEnabled() {
			c.Ui.Error("Specified HA backend has HA support disabled; please consult documentation")
			return 1
		}

		coreConfig.AdvertiseAddr = config.HABackend.AdvertiseAddr
	} else {
		if coreConfig.HAPhysical, ok = backend.(physical.HABackend); ok {
			coreConfig.AdvertiseAddr = config.Backend.AdvertiseAddr
		}
	}

	if envAA := os.Getenv("VAULT_ADVERTISE_ADDR"); envAA != "" {
		coreConfig.AdvertiseAddr = envAA
	}

	// Attempt to detect the advertise address, if possible
	var detect physical.AdvertiseDetect
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		detect, ok = coreConfig.HAPhysical.(physical.AdvertiseDetect)
	} else {
		detect, ok = coreConfig.Physical.(physical.AdvertiseDetect)
	}
	if ok && coreConfig.AdvertiseAddr == "" {
		advertise, err := c.detectAdvertise(detect, config)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error detecting advertise address: %s", err))
		} else if advertise == "" {
			c.Ui.Error("Failed to detect advertise address.")
		} else {
			coreConfig.AdvertiseAddr = advertise
		}
	}

	// Initialize the core
	core, newCoreError := vault.NewCore(coreConfig)
	if newCoreError != nil {
		if !errwrap.ContainsType(newCoreError, new(vault.NonFatalError)) {
			c.Ui.Error(fmt.Sprintf("Error initializing core: %s", newCoreError))
			return 1
		}
	}

	// If we're in dev mode, then initialize the core
	if dev {
		init, err := c.enableDev(core, devRootTokenID)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error initializing dev mode: %s", err))
			return 1
		}

		export := "export"
		quote := "'"
		if runtime.GOOS == "windows" {
			export = "set"
			quote = ""
		}

		c.Ui.Output(fmt.Sprintf(
			"==> WARNING: Dev mode is enabled!\n\n"+
				"In this mode, Vault is completely in-memory and unsealed.\n"+
				"Vault is configured to only have a single unseal key. The root\n"+
				"token has already been authenticated with the CLI, so you can\n"+
				"immediately begin using the Vault CLI.\n\n"+
				"The only step you need to take is to set the following\n"+
				"environment variables:\n\n"+
				"    "+export+" VAULT_ADDR="+quote+"http://"+config.Listeners[0].Config["address"]+quote+"\n\n"+
				"The unseal key and root token are reproduced below in case you\n"+
				"want to seal/unseal the Vault or play with authentication.\n\n"+
				"Unseal Key: %s\nRoot Token: %s\n",
			hex.EncodeToString(init.SecretShares[0]),
			init.RootToken,
		))
	}

	// Compile server information for output later
	info["backend"] = config.Backend.Type
	info["log level"] = logLevel
	info["mlock"] = fmt.Sprintf(
		"supported: %v, enabled: %v",
		mlock.Supported(), !config.DisableMlock)
	infoKeys = append(infoKeys, "log level", "mlock", "backend")

	if config.HABackend != nil {
		info["HA backend"] = config.HABackend.Type
		info["advertise address"] = coreConfig.AdvertiseAddr
		infoKeys = append(infoKeys, "HA backend", "advertise address")
	} else {
		// If the backend supports HA, then note it
		if coreConfig.HAPhysical != nil {
			if coreConfig.HAPhysical.HAEnabled() {
				info["backend"] += " (HA available)"
				info["advertise address"] = coreConfig.AdvertiseAddr
				infoKeys = append(infoKeys, "advertise address")
			} else {
				info["backend"] += " (HA disabled)"
			}
		}
	}

	// If the backend supports service discovery, run service discovery
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		sd, ok := coreConfig.HAPhysical.(physical.ServiceDiscovery)
		if ok {
			activeFunc := func() bool {
				if isLeader, _, err := core.Leader(); err == nil {
					return isLeader
				}
				return false
			}

			sealedFunc := func() bool {
				if sealed, err := core.Sealed(); err == nil {
					return sealed
				}
				return true
			}

			if err := sd.RunServiceDiscovery(c.ShutdownCh, coreConfig.AdvertiseAddr, activeFunc, sealedFunc); err != nil {
				c.Ui.Error(fmt.Sprintf("Error initializing service discovery: %v", err))
				return 1
			}
		}
	}

	// Initialize the listeners
	lns := make([]net.Listener, 0, len(config.Listeners))
	for i, lnConfig := range config.Listeners {
		ln, props, reloadFunc, err := server.NewListener(lnConfig.Type, lnConfig.Config, logGate)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error initializing listener of type %s: %s",
				lnConfig.Type, err))
			return 1
		}

		// Store the listener props for output later
		key := fmt.Sprintf("listener %d", i+1)
		propsList := make([]string, 0, len(props))
		for k, v := range props {
			propsList = append(propsList, fmt.Sprintf(
				"%s: %q", k, v))
		}
		sort.Strings(propsList)
		infoKeys = append(infoKeys, key)
		info[key] = fmt.Sprintf(
			"%s (%s)", lnConfig.Type, strings.Join(propsList, ", "))

		lns = append(lns, ln)

		if reloadFunc != nil {
			relSlice := c.ReloadFuncs["listener|"+lnConfig.Type]
			relSlice = append(relSlice, reloadFunc)
			c.ReloadFuncs["listener|"+lnConfig.Type] = relSlice
		}
	}

	// Make sure we close all listeners from this point on
	defer func() {
		for _, ln := range lns {
			ln.Close()
		}
	}()

	infoKeys = append(infoKeys, "version")
	info["version"] = version.GetVersion().String()

	// Server configuration output
	padding := 24
	sort.Strings(infoKeys)
	c.Ui.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.Ui.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.Ui.Output("")

	if verifyOnly {
		return 0
	}

	// Initialize the HTTP server
	server := &http.Server{}
	server.Handler = vaulthttp.Handler(core)
	for _, ln := range lns {
		go server.Serve(ln)
	}

	if newCoreError != nil {
		c.Ui.Output("==> Warning:\n\nNon-fatal error during initialization; check the logs for more information.")
		c.Ui.Output("")
	}

	// Output the header that the server has started
	c.Ui.Output("==> Vault server started! Log data will stream in below:\n")

	// Release the log gate.
	logGate.Flush()

	// Wait for shutdown
	shutdownTriggered := false
	for !shutdownTriggered {
		select {
		case <-c.ShutdownCh:
			c.Ui.Output("==> Vault shutdown triggered")
			if err := core.Shutdown(); err != nil {
				c.Ui.Error(fmt.Sprintf("Error with core shutdown: %s", err))
			}
			shutdownTriggered = true
		case <-c.SighupCh:
			c.Ui.Output("==> Vault reload triggered")
			if err := c.Reload(configPath); err != nil {
				c.Ui.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
			}
		}
	}

	return 0
}

func (c *ServerCommand) enableDev(core *vault.Core, rootTokenID string) (*vault.InitResult, error) {
	// Initialize it with a basic single key
	init, err := core.Initialize(&vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Copy the key so that it can be zeroed
	key := make([]byte, len(init.SecretShares[0]))
	copy(key, init.SecretShares[0])

	// Unseal the core
	unsealed, err := core.Unseal(key)
	if err != nil {
		return nil, err
	}
	if !unsealed {
		return nil, fmt.Errorf("failed to unseal Vault for dev mode")
	}

	if rootTokenID != "" {
		req := &logical.Request{
			Operation:   logical.UpdateOperation,
			ClientToken: init.RootToken,
			Path:        "auth/token/create",
			Data: map[string]interface{}{
				"id":                rootTokenID,
				"policies":          []string{"root"},
				"no_parent":         true,
				"no_default_policy": true,
			},
		}
		resp, err := core.HandleRequest(req)
		if err != nil {
			return nil, fmt.Errorf("failed to create root token with ID %s: %s", rootTokenID, err)
		}
		if resp == nil {
			return nil, fmt.Errorf("nil response when creating root token with ID %s", rootTokenID)
		}
		if resp.Auth == nil {
			return nil, fmt.Errorf("nil auth when creating root token with ID %s", rootTokenID)
		}

		init.RootToken = resp.Auth.ClientToken

		req.Path = "auth/token/revoke-self"
		req.Data = nil
		resp, err = core.HandleRequest(req)
		if err != nil {
			return nil, fmt.Errorf("failed to revoke initial root token: %s", err)
		}
	}

	// Set the token
	tokenHelper, err := c.TokenHelper()
	if err != nil {
		return nil, err
	}
	if err := tokenHelper.Store(init.RootToken); err != nil {
		return nil, err
	}

	return init, nil
}

// detectAdvertise is used to attempt advertise address detection
func (c *ServerCommand) detectAdvertise(detect physical.AdvertiseDetect,
	config *server.Config) (string, error) {
	// Get the hostname
	host, err := detect.DetectHostAddr()
	if err != nil {
		return "", err
	}

	// set [] for ipv6 addresses
	if strings.Contains(host, ":") && !strings.Contains(host, "]") {
		host = "[" + host + "]"
	}

	// Default the port and scheme
	scheme := "https"
	port := 8200

	// Attempt to detect overrides
	for _, list := range config.Listeners {
		// Only attempt TCP
		if list.Type != "tcp" {
			continue
		}

		// Check if TLS is disabled
		if val, ok := list.Config["tls_disable"]; ok {
			disable, err := strconv.ParseBool(val)
			if err != nil {
				return "", fmt.Errorf("tls_disable: %s", err)
			}

			if disable {
				scheme = "http"
			}
		}

		// Check for address override
		addr, ok := list.Config["address"]
		if !ok {
			addr = "127.0.0.1:8200"
		}

		// Check for localhost
		hostStr, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			continue
		}
		if hostStr == "127.0.0.1" {
			host = hostStr
		}

		// Check for custom port
		listPort, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		port = listPort
	}

	// Build a URL
	url := &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", host, port),
	}

	// Return the URL string
	return url.String(), nil
}

// setupTelemetry is used to setup the telemetry sub-systems
func (c *ServerCommand) setupTelemetry(config *server.Config) error {
	/* Setup telemetry
	Aggregate on 10 second intervals for 1 minute. Expose the
	metrics over stderr when there is a SIGUSR1 received.
	*/
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	var telConfig *server.Telemetry
	if config.Telemetry == nil {
		telConfig = &server.Telemetry{}
	} else {
		telConfig = config.Telemetry
	}

	metricsConf := metrics.DefaultConfig("vault")
	metricsConf.EnableHostname = !telConfig.DisableHostname

	// Configure the statsite sink
	var fanout metrics.FanoutSink
	if telConfig.StatsiteAddr != "" {
		sink, err := metrics.NewStatsiteSink(telConfig.StatsiteAddr)
		if err != nil {
			return err
		}
		fanout = append(fanout, sink)
	}

	// Configure the statsd sink
	if telConfig.StatsdAddr != "" {
		sink, err := metrics.NewStatsdSink(telConfig.StatsdAddr)
		if err != nil {
			return err
		}
		fanout = append(fanout, sink)
	}

	// Initialize the global sink
	if len(fanout) > 0 {
		fanout = append(fanout, inm)
		metrics.NewGlobal(metricsConf, fanout)
	} else {
		metricsConf.EnableHostname = false
		metrics.NewGlobal(metricsConf, inm)
	}
	return nil
}

func (c *ServerCommand) Reload(configPath []string) error {
	// Read the new config
	var config *server.Config
	for _, path := range configPath {
		current, err := server.LoadConfig(path)
		if err != nil {
			retErr := fmt.Errorf("Error loading configuration from %s: %s", path, err)
			c.Ui.Error(retErr.Error())
			return retErr
		}

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}

	// Ensure at least one config was found.
	if config == nil {
		retErr := fmt.Errorf("No configuration files found")
		c.Ui.Error(retErr.Error())
		return retErr
	}

	var reloadErrors *multierror.Error
	// Call reload on the listeners. This will call each listener with each
	// config block, but they verify the address.
	for _, lnConfig := range config.Listeners {
		for _, relFunc := range c.ReloadFuncs["listener|"+lnConfig.Type] {
			if err := relFunc(lnConfig.Config); err != nil {
				retErr := fmt.Errorf("Error encountered reloading configuration: %s", err)
				reloadErrors = multierror.Append(retErr)
			}
		}
	}

	return reloadErrors.ErrorOrNil()
}

func (c *ServerCommand) Synopsis() string {
	return "Start a Vault server"
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: vault server [options]

  Start a Vault server.

  This command starts a Vault server that responds to API requests.
  Vault will start in a "sealed" state. The Vault must be unsealed
  with "vault unseal" or the API before this server can respond to requests.
  This must be done for every server.

  If the server is being started against a storage backend that has
  brand new (no existing Vault data in it), it must be initialized with
  "vault init" or the API first.


General Options:

  -config=<path>          Path to the configuration file or directory. This can
                          be specified multiple times. If it is a directory,
                          all files with a ".hcl" or ".json" suffix will be
                          loaded.

  -dev                    Enables Dev mode. In this mode, Vault is completely
                          in-memory and unsealed. Do not run the Dev server in
                          production!

  -dev-root-token-id=""   If set, the root token returned in Dev mode will have
                          the given ID. This *only* has an effect when running
                          in Dev mode. Can also be specified with the
                          VAULT_DEV_ROOT_TOKEN_ID environment variable.

  -dev-listen-address=""  If set, this overrides the normal Dev mode listen
                          address of "127.0.0.1:8200". Can also be specified
                          with the VAULT_DEV_LISTEN_ADDRESS environment
                          variable.

  -log-level=info         Log verbosity. Defaults to "info", will be output to
                          stderr. Supported values: "trace", "debug", "info",
                          "warn", "err"
`
	return strings.TrimSpace(helpText)
}

// MakeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// SIGINT or SIGTERM received.
func MakeShutdownCh() chan struct{} {
	resultCh := make(chan struct{})

	shutdownCh := make(chan os.Signal, 4)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-shutdownCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}

// MakeSighupCh returns a channel that can be used for SIGHUP
// reloading. This channel will send a message for every
// SIGHUP received.
func MakeSighupCh() chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGHUP)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}
