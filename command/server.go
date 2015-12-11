package command

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/flag-slice"
	"github.com/hashicorp/vault/helper/gated-writer"
	"github.com/hashicorp/vault/helper/mlock"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

// ServerCommand is a Command that starts the Vault server.
type ServerCommand struct {
	AuditBackends      map[string]audit.Factory
	CredentialBackends map[string]logical.Factory
	LogicalBackends    map[string]logical.Factory

	ShutdownCh <-chan struct{}
	Meta
}

func (c *ServerCommand) Run(args []string) int {
	var dev, verifyOnly bool
	var configPath []string
	var logLevel string
	flags := c.Meta.FlagSet("server", FlagSetDefault)
	flags.BoolVar(&dev, "dev", false, "")
	flags.StringVar(&logLevel, "log-level", "info", "")
	flags.BoolVar(&verifyOnly, "verify-only", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	flags.Var((*sliceflag.StringFlag)(&configPath), "config", "config")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validation
	if !dev && len(configPath) == 0 {
		c.Ui.Error("At least one config path must be specified with -config")
		flags.Usage()
		return 1
	}

	// Load the configuration
	var config *server.Config
	if dev {
		config = server.DevConfig()
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

	// Ensure that a backend is provided
	if config.Backend == nil {
		c.Ui.Error("A physical backend must be specified")
		return 1
	}

	// If mlock isn't supported, show a warning. We disable this in
	// dev because it is quite scary to see when first using Vault.
	if !dev && !mlock.Supported() {
		c.Ui.Output("==> WARNING: mlock not supported on this system!\n")
		c.Ui.Output("  The `mlock` syscall to prevent memory from being swapped to")
		c.Ui.Output("  disk is not supported on this system. Enabling mlock or")
		c.Ui.Output("  running Vault on a system with mlock is much more secure.\n")
	}

	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	logGate := &gatedwriter.Writer{Writer: os.Stderr}
	logger := log.New(&logutils.LevelFilter{
		Levels: []logutils.LogLevel{
			"TRACE", "DEBUG", "INFO", "WARN", "ERR"},
		MinLevel: logutils.LogLevel(strings.ToUpper(logLevel)),
		Writer:   logGate,
	}, "", log.LstdFlags)

	// Initialize the backend
	backend, err := physical.NewBackend(
		config.Backend.Type, config.Backend.Config)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing backend of type %s: %s",
			config.Backend.Type, err))
		return 1
	}

	var advertiseAddr string = config.Backend.AdvertiseAddr

	// Note that "habackend" is a backend that *may* support HA;
	// it defaults to the same backend as normal operations
	var habackend physical.Backend = backend

	// Initialize the separate HA physical backend, if it exists
	if config.HABackend != nil {
		habackend, err = physical.NewBackend(
			config.HABackend.Type, config.HABackend.Config)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error initializing backend of type %s: %s",
				config.HABackend.Type, err))
			return 1
		}
		if _, ok := habackend.(physical.HABackend); !ok {
			c.Ui.Error("Specified HA backend does not support HA")
			return 1
		}

		advertiseAddr = config.HABackend.AdvertiseAddr
	}

	// Attempt to detect the advertise address possible
	if detect, ok := habackend.(physical.AdvertiseDetect); ok && advertiseAddr == "" {
		advertise, err := c.detectAdvertise(detect, config)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error detecting advertise address: %s", err))
		} else if advertise == "" {
			c.Ui.Error("Failed to detect advertise address.")
		} else {
			advertiseAddr = advertise
		}
	}

	// Initialize the core
	core, err := vault.NewCore(&vault.CoreConfig{
		AdvertiseAddr:      advertiseAddr,
		Physical:           backend,
		HAPhysical:         habackend,
		AuditBackends:      c.AuditBackends,
		CredentialBackends: c.CredentialBackends,
		LogicalBackends:    c.LogicalBackends,
		Logger:             logger,
		DisableCache:       config.DisableCache,
		DisableMlock:       config.DisableMlock,
		MaxLeaseTTL:        config.MaxLeaseTTL,
		DefaultLeaseTTL:    config.DefaultLeaseTTL,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing core: %s", err))
		return 1
	}

	// If we're in dev mode, then initialize the core
	if dev {
		init, err := c.enableDev(core)
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
				"    "+export+" VAULT_ADDR="+quote+"http://127.0.0.1:8200"+quote+"\n\n"+
				"The unseal key and root token are reproduced below in case you\n"+
				"want to seal/unseal the Vault or play with authentication.\n\n"+
				"Unseal Key: %s\nRoot Token: %s\n",
			hex.EncodeToString(init.SecretShares[0]),
			init.RootToken,
		))
	}

	// Compile server information for output later
	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["backend"] = config.Backend.Type
	info["log level"] = logLevel
	info["mlock"] = fmt.Sprintf(
		"supported: %v, enabled: %v",
		mlock.Supported(), !config.DisableMlock)
	infoKeys = append(infoKeys, "log level", "mlock", "backend")

	if config.HABackend != nil {
		info["HA backend"] = config.HABackend.Type
		info["advertise address"] = advertiseAddr
		infoKeys = append(infoKeys, "HA backend", "advertise address")
	} else {
		// If the backend supports HA, then note it
		if _, ok := habackend.(physical.HABackend); ok {
			info["backend"] += " (HA available)"
			info["advertise address"] = advertiseAddr
			infoKeys = append(infoKeys, "advertise address")
		}
	}

	// Initialize the telemetry
	if err := c.setupTelementry(config); err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}

	// Initialize the listeners
	lns := make([]net.Listener, 0, len(config.Listeners))
	for i, lnConfig := range config.Listeners {
		ln, props, err := server.NewListener(lnConfig.Type, lnConfig.Config)
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
	}

	if verifyOnly {
		return 0
	}

	// Initialize the HTTP server
	server := &http.Server{}
	server.Handler = vaulthttp.Handler(core)
	for _, ln := range lns {
		go server.Serve(ln)
	}

	infoKeys = append(infoKeys, "version")
	info["version"] = version.GetVersion().String()

	// Server configuration output
	padding := 18
	c.Ui.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.Ui.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.Ui.Output("")

	// Output the header that the server has started
	c.Ui.Output("==> Vault server started! Log data will stream in below:\n")

	// Release the log gate.
	logGate.Flush()

	// Wait for shutdown
	select {
	case <-c.ShutdownCh:
		c.Ui.Output("==> Vault shutdown triggered")
		if err := core.Shutdown(); err != nil {
			c.Ui.Error(fmt.Sprintf("Error with core shutdown: %s", err))
		}
	}
	return 0
}

func (c *ServerCommand) enableDev(core *vault.Core) (*vault.InitResult, error) {
	// Initialize it with a basic single key
	init, err := core.Initialize(&vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
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

// setupTelementry is used ot setup the telemetry sub-systems
func (c *ServerCommand) setupTelementry(config *server.Config) error {
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

  -config=<path>      Path to the configuration file or directory. This can be
                      specified multiple times. If it is a directory, all
                      files with a ".hcl" or ".json" suffix will be loaded.

  -dev                Enables Dev mode. In this mode, Vault is completely
                      in-memory and unsealed. Do not run the Dev server in
                      production!

  -log-level=info     Log verbosity. Defaults to "info", will be outputted
                      to stderr. Supported values: "trace", "debug", "info",
                      "warn", "err"

`
	return strings.TrimSpace(helpText)
}
