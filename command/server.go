package command

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	colorable "github.com/mattn/go-colorable"
	log "github.com/mgutz/logxi/v1"
	testing "github.com/mitchellh/go-testing-interface"
	"github.com/posener/complete"

	"google.golang.org/grpc/grpclog"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/circonus"
	"github.com/armon/go-metrics/datadog"
	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/flag-slice"
	"github.com/hashicorp/vault/helper/gated-writer"
	"github.com/hashicorp/vault/helper/logbridge"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/helper/mlock"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/reload"
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
	PhysicalBackends   map[string]physical.Factory

	ShutdownCh chan struct{}
	SighupCh   chan struct{}

	WaitGroup *sync.WaitGroup

	meta.Meta

	logGate *gatedwriter.Writer
	logger  log.Logger

	cleanupGuard sync.Once

	reloadFuncsLock *sync.RWMutex
	reloadFuncs     *map[string][]reload.ReloadFunc
}

func (c *ServerCommand) Run(args []string) int {
	var dev, verifyOnly, devHA, devTransactional, devLeasedKV, devThreeNode, devSkipInit bool
	var configPath []string
	var logLevel, devRootTokenID, devListenAddress, devPluginDir string
	var devLatency, devLatencyJitter int
	flags := c.Meta.FlagSet("server", meta.FlagSetDefault)
	flags.BoolVar(&dev, "dev", false, "")
	flags.StringVar(&devRootTokenID, "dev-root-token-id", "", "")
	flags.StringVar(&devListenAddress, "dev-listen-address", "", "")
	flags.StringVar(&devPluginDir, "dev-plugin-dir", "", "")
	flags.StringVar(&logLevel, "log-level", "info", "")
	flags.IntVar(&devLatency, "dev-latency", 0, "")
	flags.IntVar(&devLatencyJitter, "dev-latency-jitter", 20, "")
	flags.BoolVar(&verifyOnly, "verify-only", false, "")
	flags.BoolVar(&devHA, "dev-ha", false, "")
	flags.BoolVar(&devTransactional, "dev-transactional", false, "")
	flags.BoolVar(&devLeasedKV, "dev-leased-kv", false, "")
	flags.BoolVar(&devThreeNode, "dev-three-node", false, "")
	flags.BoolVar(&devSkipInit, "dev-skip-init", false, "")
	flags.Usage = func() { c.Ui.Output(c.Help()) }
	flags.Var((*sliceflag.StringFlag)(&configPath), "config", "config")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	c.logGate = &gatedwriter.Writer{Writer: colorable.NewColorable(os.Stderr)}
	var level int
	logLevel = strings.ToLower(strings.TrimSpace(logLevel))
	switch logLevel {
	case "trace":
		level = log.LevelTrace
	case "debug":
		level = log.LevelDebug
	case "info":
		level = log.LevelInfo
	case "notice":
		level = log.LevelNotice
	case "warn":
		level = log.LevelWarn
	case "err":
		level = log.LevelError
	default:
		c.Ui.Output(fmt.Sprintf("Unknown log level %s", logLevel))
		return 1
	}

	logFormat := os.Getenv("VAULT_LOG_FORMAT")
	if logFormat == "" {
		logFormat = os.Getenv("LOGXI_FORMAT")
	}
	switch strings.ToLower(logFormat) {
	case "vault", "vault_json", "vault-json", "vaultjson", "json", "":
		if devThreeNode {
			c.logger = logbridge.NewLogger(hclog.New(&hclog.LoggerOptions{
				Mutex:  &sync.Mutex{},
				Output: c.logGate,
			})).LogxiLogger()
		} else {
			c.logger = logformat.NewVaultLoggerWithWriter(c.logGate, level)
		}
	default:
		c.logger = log.NewLogger(c.logGate, "vault")
		c.logger.SetLevel(level)
	}
	grpclog.SetLogger(&grpclogFaker{
		logger: c.logger,
		log:    os.Getenv("VAULT_GRPC_LOGGING") != "",
	})

	if os.Getenv("VAULT_DEV_ROOT_TOKEN_ID") != "" && devRootTokenID == "" {
		devRootTokenID = os.Getenv("VAULT_DEV_ROOT_TOKEN_ID")
	}

	if os.Getenv("VAULT_DEV_LISTEN_ADDRESS") != "" && devListenAddress == "" {
		devListenAddress = os.Getenv("VAULT_DEV_LISTEN_ADDRESS")
	}

	if devHA || devTransactional || devLeasedKV || devThreeNode {
		dev = true
	}

	// Validation
	if !dev {
		switch {
		case len(configPath) == 0:
			c.Ui.Output("At least one config path must be specified with -config")
			flags.Usage()
			return 1
		case devRootTokenID != "":
			c.Ui.Output("Root token ID can only be specified with -dev")
			flags.Usage()
			return 1
		}
	}

	// Load the configuration
	var config *server.Config
	if dev {
		config = server.DevConfig(devHA, devTransactional)
		if devListenAddress != "" {
			config.Listeners[0].Config["address"] = devListenAddress
		}
	}
	for _, path := range configPath {
		current, err := server.LoadConfig(path, c.logger)
		if err != nil {
			c.Ui.Output(fmt.Sprintf(
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
		c.Ui.Output("No configuration files found.")
		return 1
	}

	// Ensure that a backend is provided
	if config.Storage == nil {
		c.Ui.Output("A storage backend must be specified")
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

	if err := c.setupTelemetry(config); err != nil {
		c.Ui.Output(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}

	// Initialize the backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
	if !exists {
		c.Ui.Output(fmt.Sprintf(
			"Unknown storage type %s",
			config.Storage.Type))
		return 1
	}
	backend, err := factory(config.Storage.Config, c.logger)
	if err != nil {
		c.Ui.Output(fmt.Sprintf(
			"Error initializing storage of type %s: %s",
			config.Storage.Type, err))
		return 1
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = logLevel
	infoKeys = append(infoKeys, "log level")

	var seal vault.Seal = &vault.DefaultSeal{}

	// Ensure that the seal finalizer is called, even if using verify-only
	defer func() {
		if seal != nil {
			err = seal.Finalize()
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Error finalizing seals: %v", err))
			}
		}
	}()

	if seal == nil {
		c.Ui.Error(fmt.Sprintf("Could not create seal; most likely proper Seal configuration information was not set, but no error was generated."))
		return 1
	}

	coreConfig := &vault.CoreConfig{
		Physical:           backend,
		RedirectAddr:       config.Storage.RedirectAddr,
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
		ClusterName:        config.ClusterName,
		CacheSize:          config.CacheSize,
		PluginDirectory:    config.PluginDirectory,
		EnableRaw:          config.EnableRawEndpoint,
	}

	if dev {
		coreConfig.DevToken = devRootTokenID
		if devLeasedKV {
			coreConfig.LogicalBackends["kv"] = vault.LeasedPassthroughBackendFactory
		}
		if devPluginDir != "" {
			coreConfig.PluginDirectory = devPluginDir
		}
		if devLatency > 0 {
			injectLatency := time.Duration(devLatency) * time.Millisecond
			if _, txnOK := backend.(physical.Transactional); txnOK {
				coreConfig.Physical = physical.NewTransactionalLatencyInjector(backend, injectLatency, devLatencyJitter, c.logger)
			} else {
				coreConfig.Physical = physical.NewLatencyInjector(backend, injectLatency, devLatencyJitter, c.logger)
			}
		}
	}

	if devThreeNode {
		return c.enableThreeNodeDevCluster(coreConfig, info, infoKeys, devListenAddress)
	}

	var disableClustering bool

	// Initialize the separate HA storage backend, if it exists
	var ok bool
	if config.HAStorage != nil {
		factory, exists := c.PhysicalBackends[config.HAStorage.Type]
		if !exists {
			c.Ui.Output(fmt.Sprintf(
				"Unknown HA storage type %s",
				config.HAStorage.Type))
			return 1
		}
		habackend, err := factory(config.HAStorage.Config, c.logger)
		if err != nil {
			c.Ui.Output(fmt.Sprintf(
				"Error initializing HA storage of type %s: %s",
				config.HAStorage.Type, err))
			return 1
		}

		if coreConfig.HAPhysical, ok = habackend.(physical.HABackend); !ok {
			c.Ui.Output("Specified HA storage does not support HA")
			return 1
		}

		if !coreConfig.HAPhysical.HAEnabled() {
			c.Ui.Output("Specified HA storage has HA support disabled; please consult documentation")
			return 1
		}

		coreConfig.RedirectAddr = config.HAStorage.RedirectAddr
		disableClustering = config.HAStorage.DisableClustering
		if !disableClustering {
			coreConfig.ClusterAddr = config.HAStorage.ClusterAddr
		}
	} else {
		if coreConfig.HAPhysical, ok = backend.(physical.HABackend); ok {
			coreConfig.RedirectAddr = config.Storage.RedirectAddr
			disableClustering = config.Storage.DisableClustering
			if !disableClustering {
				coreConfig.ClusterAddr = config.Storage.ClusterAddr
			}
		}
	}

	if envRA := os.Getenv("VAULT_API_ADDR"); envRA != "" {
		coreConfig.RedirectAddr = envRA
	} else if envRA := os.Getenv("VAULT_REDIRECT_ADDR"); envRA != "" {
		coreConfig.RedirectAddr = envRA
	} else if envAA := os.Getenv("VAULT_ADVERTISE_ADDR"); envAA != "" {
		coreConfig.RedirectAddr = envAA
	}

	// Attempt to detect the redirect address, if possible
	var detect physical.RedirectDetect
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		detect, ok = coreConfig.HAPhysical.(physical.RedirectDetect)
	} else {
		detect, ok = coreConfig.Physical.(physical.RedirectDetect)
	}
	if ok && coreConfig.RedirectAddr == "" {
		redirect, err := c.detectRedirect(detect, config)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Error detecting redirect address: %s", err))
		} else if redirect == "" {
			c.Ui.Output("Failed to detect redirect address.")
		} else {
			coreConfig.RedirectAddr = redirect
		}
	}
	if coreConfig.RedirectAddr == "" && dev {
		coreConfig.RedirectAddr = fmt.Sprintf("http://%s", config.Listeners[0].Config["address"])
	}

	// After the redirect bits are sorted out, if no cluster address was
	// explicitly given, derive one from the redirect addr
	if disableClustering {
		coreConfig.ClusterAddr = ""
	} else if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
		coreConfig.ClusterAddr = envCA
	} else {
		var addrToUse string
		switch {
		case coreConfig.ClusterAddr == "" && coreConfig.RedirectAddr != "":
			addrToUse = coreConfig.RedirectAddr
		case dev:
			addrToUse = fmt.Sprintf("http://%s", config.Listeners[0].Config["address"])
		default:
			goto CLUSTER_SYNTHESIS_COMPLETE
		}
		u, err := url.ParseRequestURI(addrToUse)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Error parsing synthesized cluster address %s: %v", addrToUse, err))
			return 1
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			// This sucks, as it's a const in the function but not exported in the package
			if strings.Contains(err.Error(), "missing port in address") {
				host = u.Host
				port = "443"
			} else {
				c.Ui.Output(fmt.Sprintf("Error parsing redirect address: %v", err))
				return 1
			}
		}
		nPort, err := strconv.Atoi(port)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Error parsing synthesized address; failed to convert %q to a numeric: %v", port, err))
			return 1
		}
		u.Host = net.JoinHostPort(host, strconv.Itoa(nPort+1))
		// Will always be TLS-secured
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}

CLUSTER_SYNTHESIS_COMPLETE:

	if coreConfig.ClusterAddr != "" {
		// Force https as we'll always be TLS-secured
		u, err := url.ParseRequestURI(coreConfig.ClusterAddr)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("Error parsing cluster address %s: %v", coreConfig.RedirectAddr, err))
			return 1
		}
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}

	// Initialize the core
	core, newCoreError := vault.NewCore(coreConfig)
	if newCoreError != nil {
		if !errwrap.ContainsType(newCoreError, new(vault.NonFatalError)) {
			c.Ui.Output(fmt.Sprintf("Error initializing core: %s", newCoreError))
			return 1
		}
	}

	// Copy the reload funcs pointers back
	c.reloadFuncs = coreConfig.ReloadFuncs
	c.reloadFuncsLock = coreConfig.ReloadFuncsLock

	// Compile server information for output later
	info["storage"] = config.Storage.Type
	info["mlock"] = fmt.Sprintf(
		"supported: %v, enabled: %v",
		mlock.Supported(), !config.DisableMlock && mlock.Supported())
	infoKeys = append(infoKeys, "mlock", "storage")

	if coreConfig.ClusterAddr != "" {
		info["cluster address"] = coreConfig.ClusterAddr
		infoKeys = append(infoKeys, "cluster address")
	}
	if coreConfig.RedirectAddr != "" {
		info["redirect address"] = coreConfig.RedirectAddr
		infoKeys = append(infoKeys, "redirect address")
	}

	if config.HAStorage != nil {
		info["HA storage"] = config.HAStorage.Type
		infoKeys = append(infoKeys, "HA storage")
	} else {
		// If the storage supports HA, then note it
		if coreConfig.HAPhysical != nil {
			if coreConfig.HAPhysical.HAEnabled() {
				info["storage"] += " (HA available)"
			} else {
				info["storage"] += " (HA disabled)"
			}
		}
	}

	clusterAddrs := []*net.TCPAddr{}

	// Initialize the listeners
	c.reloadFuncsLock.Lock()
	lns := make([]net.Listener, 0, len(config.Listeners))
	for i, lnConfig := range config.Listeners {
		ln, props, reloadFunc, err := server.NewListener(lnConfig.Type, lnConfig.Config, c.logGate)
		if err != nil {
			c.Ui.Output(fmt.Sprintf(
				"Error initializing listener of type %s: %s",
				lnConfig.Type, err))
			return 1
		}

		lns = append(lns, ln)

		if reloadFunc != nil {
			relSlice := (*c.reloadFuncs)["listener|"+lnConfig.Type]
			relSlice = append(relSlice, reloadFunc)
			(*c.reloadFuncs)["listener|"+lnConfig.Type] = relSlice
		}

		if !disableClustering && lnConfig.Type == "tcp" {
			var addrRaw interface{}
			var addr string
			var ok bool
			if addrRaw, ok = lnConfig.Config["cluster_address"]; ok {
				addr = addrRaw.(string)
				tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
				if err != nil {
					c.Ui.Output(fmt.Sprintf(
						"Error resolving cluster_address: %s",
						err))
					return 1
				}
				clusterAddrs = append(clusterAddrs, tcpAddr)
			} else {
				tcpAddr, ok := ln.Addr().(*net.TCPAddr)
				if !ok {
					c.Ui.Output("Failed to parse tcp listener")
					return 1
				}
				clusterAddr := &net.TCPAddr{
					IP:   tcpAddr.IP,
					Port: tcpAddr.Port + 1,
				}
				clusterAddrs = append(clusterAddrs, clusterAddr)
				addr = clusterAddr.String()
			}
			props["cluster address"] = addr
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

	}
	c.reloadFuncsLock.Unlock()
	if !disableClustering {
		if c.logger.IsTrace() {
			c.logger.Trace("cluster listener addresses synthesized", "cluster_addresses", clusterAddrs)
		}
	}

	// Make sure we close all listeners from this point on
	listenerCloseFunc := func() {
		for _, ln := range lns {
			ln.Close()
		}
	}

	defer c.cleanupGuard.Do(listenerCloseFunc)

	infoKeys = append(infoKeys, "version")
	verInfo := version.GetVersion()
	info["version"] = verInfo.FullVersionNumber(false)
	if verInfo.Revision != "" {
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	}
	infoKeys = append(infoKeys, "cgo")
	info["cgo"] = "disabled"
	if version.CgoEnabled {
		info["cgo"] = "enabled"
	}

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

	handler := vaulthttp.Handler(core)

	// This needs to happen before we first unseal, so before we trigger dev
	// mode if it's set
	core.SetClusterListenerAddrs(clusterAddrs)
	core.SetClusterHandler(handler)

	err = core.UnsealWithStoredKeys()
	if err != nil {
		if !errwrap.ContainsType(err, new(vault.NonFatalError)) {
			c.Ui.Output(fmt.Sprintf("Error initializing core: %s", err))
			return 1
		}
	}

	// Perform service discovery registrations and initialization of
	// HTTP server after the verifyOnly check.

	// Instantiate the wait group
	c.WaitGroup = &sync.WaitGroup{}

	// If the backend supports service discovery, run service discovery
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		sd, ok := coreConfig.HAPhysical.(physical.ServiceDiscovery)
		if ok {
			activeFunc := func() bool {
				if isLeader, _, _, err := core.Leader(); err == nil {
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

			if err := sd.RunServiceDiscovery(c.WaitGroup, c.ShutdownCh, coreConfig.RedirectAddr, activeFunc, sealedFunc); err != nil {
				c.Ui.Output(fmt.Sprintf("Error initializing service discovery: %v", err))
				return 1
			}
		}
	}

	// If we're in Dev mode, then initialize the core
	if dev && !devSkipInit {
		init, err := c.enableDev(core, coreConfig)
		if err != nil {
			c.Ui.Output(fmt.Sprintf(
				"Error initializing Dev mode: %s", err))
			return 1
		}

		export := "export"
		quote := "'"
		if runtime.GOOS == "windows" {
			export = "set"
			quote = ""
		}

		c.Ui.Output(fmt.Sprint(
			"==> WARNING: Dev mode is enabled!\n\n" +
				"In this mode, Vault is completely in-memory and unsealed.\n" +
				"Vault is configured to only have a single unseal key. The root\n" +
				"token has already been authenticated with the CLI, so you can\n" +
				"immediately begin using the Vault CLI.\n\n" +
				"The only step you need to take is to set the following\n" +
				"environment variables:\n\n" +
				"    " + export + " VAULT_ADDR=" + quote + "http://" + config.Listeners[0].Config["address"].(string) + quote + "\n\n" +
				"The unseal key and root token are reproduced below in case you\n" +
				"want to seal/unseal the Vault or play with authentication.\n",
		))

		// Unseal key is not returned if stored shares is supported
		if len(init.SecretShares) > 0 {
			c.Ui.Output(fmt.Sprintf(
				"Unseal Key: %s",
				base64.StdEncoding.EncodeToString(init.SecretShares[0]),
			))
		}

		if len(init.RecoveryShares) > 0 {
			c.Ui.Output(fmt.Sprintf(
				"Recovery Key: %s",
				base64.StdEncoding.EncodeToString(init.RecoveryShares[0]),
			))
		}

		c.Ui.Output(fmt.Sprintf(
			"Root Token: %s\n",
			init.RootToken,
		))
	}

	// Initialize the HTTP servers
	for _, ln := range lns {
		server := &http.Server{
			Handler: handler,
		}
		go server.Serve(ln)
	}

	if newCoreError != nil {
		c.Ui.Output("==> Warning:\n\nNon-fatal error during initialization; check the logs for more information.")
		c.Ui.Output("")
	}

	// Output the header that the server has started
	c.Ui.Output("==> Vault server started! Log data will stream in below:\n")

	// Release the log gate.
	c.logGate.Flush()

	// Write out the PID to the file now that server has successfully started
	if err := c.storePidFile(config.PidFile); err != nil {
		c.Ui.Output(fmt.Sprintf("Error storing PID: %v", err))
		return 1
	}

	defer func() {
		if err := c.removePidFile(config.PidFile); err != nil {
			c.Ui.Output(fmt.Sprintf("Error deleting the PID file: %v", err))
		}
	}()

	// Wait for shutdown
	shutdownTriggered := false

	for !shutdownTriggered {
		select {
		case <-c.ShutdownCh:
			c.Ui.Output("==> Vault shutdown triggered")

			// Stop the listners so that we don't process further client requests.
			c.cleanupGuard.Do(listenerCloseFunc)

			// Shutdown will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			if err := core.Shutdown(); err != nil {
				c.Ui.Output(fmt.Sprintf("Error with core shutdown: %s", err))
			}

			shutdownTriggered = true

		case <-c.SighupCh:
			c.Ui.Output("==> Vault reload triggered")
			if err := c.Reload(c.reloadFuncsLock, c.reloadFuncs, configPath); err != nil {
				c.Ui.Output(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
			}
		}
	}

	// Wait for dependent goroutines to complete
	c.WaitGroup.Wait()
	return 0
}

func (c *ServerCommand) enableDev(core *vault.Core, coreConfig *vault.CoreConfig) (*vault.InitResult, error) {
	var recoveryConfig *vault.SealConfig
	barrierConfig := &vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}

	if core.SealAccess().RecoveryKeySupported() {
		recoveryConfig = &vault.SealConfig{
			SecretShares:    1,
			SecretThreshold: 1,
		}
	}

	if core.SealAccess().StoredKeysSupported() {
		barrierConfig.StoredShares = 1
	}

	// Initialize it with a basic single key
	init, err := core.Initialize(&vault.InitParams{
		BarrierConfig:  barrierConfig,
		RecoveryConfig: recoveryConfig,
	})
	if err != nil {
		return nil, err
	}

	// Handle unseal with stored keys
	if core.SealAccess().StoredKeysSupported() {
		err := core.UnsealWithStoredKeys()
		if err != nil {
			return nil, err
		}
	} else {
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
	}

	isLeader, _, _, err := core.Leader()
	if err != nil && err != vault.ErrHANotEnabled {
		return nil, fmt.Errorf("failed to check active status: %v", err)
	}
	if err == nil {
		leaderCount := 5
		for !isLeader {
			if leaderCount == 0 {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				return nil, fmt.Errorf("failed to get active status after five seconds; call stack is\n%s\n", buf)
			}
			time.Sleep(1 * time.Second)
			isLeader, _, _, err = core.Leader()
			if err != nil {
				return nil, fmt.Errorf("failed to check active status: %v", err)
			}
			leaderCount--
		}
	}

	// Generate a dev root token if one is provided in the flag
	if coreConfig.DevToken != "" {
		req := &logical.Request{
			ID:          "dev-gen-root",
			Operation:   logical.UpdateOperation,
			ClientToken: init.RootToken,
			Path:        "auth/token/create",
			Data: map[string]interface{}{
				"id":                coreConfig.DevToken,
				"policies":          []string{"root"},
				"no_parent":         true,
				"no_default_policy": true,
			},
		}
		resp, err := core.HandleRequest(req)
		if err != nil {
			return nil, fmt.Errorf("failed to create root token with ID %s: %s", coreConfig.DevToken, err)
		}
		if resp == nil {
			return nil, fmt.Errorf("nil response when creating root token with ID %s", coreConfig.DevToken)
		}
		if resp.Auth == nil {
			return nil, fmt.Errorf("nil auth when creating root token with ID %s", coreConfig.DevToken)
		}

		init.RootToken = resp.Auth.ClientToken

		req.ID = "dev-revoke-init-root"
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

func (c *ServerCommand) enableThreeNodeDevCluster(base *vault.CoreConfig, info map[string]string, infoKeys []string, devListenAddress string) int {
	testCluster := vault.NewTestCluster(&testing.RuntimeT{}, base, &vault.TestClusterOptions{
		HandlerFunc:       vaulthttp.Handler,
		BaseListenAddress: devListenAddress,
		RawLogger:         c.logger,
	})
	defer c.cleanupGuard.Do(testCluster.Cleanup)

	info["cluster parameters path"] = testCluster.TempDir
	infoKeys = append(infoKeys, "cluster parameters path")

	for i, core := range testCluster.Cores {
		info[fmt.Sprintf("node %d redirect address", i)] = fmt.Sprintf("https://%s", core.Listeners[0].Address.String())
		infoKeys = append(infoKeys, fmt.Sprintf("node %d redirect address", i))
	}

	infoKeys = append(infoKeys, "version")
	verInfo := version.GetVersion()
	info["version"] = verInfo.FullVersionNumber(false)
	if verInfo.Revision != "" {
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	}
	infoKeys = append(infoKeys, "cgo")
	info["cgo"] = "disabled"
	if version.CgoEnabled {
		info["cgo"] = "enabled"
	}

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

	for _, core := range testCluster.Cores {
		core.Server.Handler = vaulthttp.Handler(core.Core)
		core.SetClusterHandler(core.Server.Handler)
	}

	testCluster.Start()

	if base.DevToken != "" {
		req := &logical.Request{
			ID:          "dev-gen-root",
			Operation:   logical.UpdateOperation,
			ClientToken: testCluster.RootToken,
			Path:        "auth/token/create",
			Data: map[string]interface{}{
				"id":                base.DevToken,
				"policies":          []string{"root"},
				"no_parent":         true,
				"no_default_policy": true,
			},
		}
		resp, err := testCluster.Cores[0].HandleRequest(req)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("failed to create root token with ID %s: %s", base.DevToken, err))
			return 1
		}
		if resp == nil {
			c.Ui.Output(fmt.Sprintf("nil response when creating root token with ID %s", base.DevToken))
			return 1
		}
		if resp.Auth == nil {
			c.Ui.Output(fmt.Sprintf("nil auth when creating root token with ID %s", base.DevToken))
			return 1
		}

		testCluster.RootToken = resp.Auth.ClientToken

		req.ID = "dev-revoke-init-root"
		req.Path = "auth/token/revoke-self"
		req.Data = nil
		resp, err = testCluster.Cores[0].HandleRequest(req)
		if err != nil {
			c.Ui.Output(fmt.Sprintf("failed to revoke initial root token: %s", err))
			return 1
		}
	}

	// Set the token
	tokenHelper, err := c.TokenHelper()
	if err != nil {
		c.Ui.Output(fmt.Sprintf("%v", err))
		return 1
	}
	if err := tokenHelper.Store(testCluster.RootToken); err != nil {
		c.Ui.Output(fmt.Sprintf("%v", err))
		return 1
	}

	if err := ioutil.WriteFile(filepath.Join(testCluster.TempDir, "root_token"), []byte(testCluster.RootToken), 0755); err != nil {
		c.Ui.Output(fmt.Sprintf("%v", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf(
		"==> Three node dev mode is enabled\n\n" +
			"The unseal key and root token are reproduced below in case you\n" +
			"want to seal/unseal the Vault or play with authentication.\n",
	))

	for i, key := range testCluster.BarrierKeys {
		c.Ui.Output(fmt.Sprintf(
			"Unseal Key %d: %s",
			i+1, base64.StdEncoding.EncodeToString(key),
		))
	}

	c.Ui.Output(fmt.Sprintf(
		"\nRoot Token: %s\n", testCluster.RootToken,
	))

	c.Ui.Output(fmt.Sprintf(
		"\nUseful env vars:\n"+
			"VAULT_TOKEN=%s\n"+
			"VAULT_ADDR=%s\n"+
			"VAULT_CACERT=%s/ca_cert.pem\n",
		testCluster.RootToken,
		testCluster.Cores[0].Client.Address(),
		testCluster.TempDir,
	))

	// Output the header that the server has started
	c.Ui.Output("==> Vault server started! Log data will stream in below:\n")

	// Release the log gate.
	c.logGate.Flush()

	// Wait for shutdown
	shutdownTriggered := false

	for !shutdownTriggered {
		select {
		case <-c.ShutdownCh:
			c.Ui.Output("==> Vault shutdown triggered")

			// Stop the listners so that we don't process further client requests.
			c.cleanupGuard.Do(testCluster.Cleanup)

			// Shutdown will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			for _, core := range testCluster.Cores {
				if err := core.Shutdown(); err != nil {
					c.Ui.Output(fmt.Sprintf("Error with core shutdown: %s", err))
				}
			}

			shutdownTriggered = true

		case <-c.SighupCh:
			c.Ui.Output("==> Vault reload triggered")
			for _, core := range testCluster.Cores {
				if err := c.Reload(core.ReloadFuncsLock, core.ReloadFuncs, nil); err != nil {
					c.Ui.Output(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
				}
			}
		}
	}

	return 0
}

// detectRedirect is used to attempt redirect address detection
func (c *ServerCommand) detectRedirect(detect physical.RedirectDetect,
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
			disable, err := parseutil.ParseBool(val)
			if err != nil {
				return "", fmt.Errorf("tls_disable: %s", err)
			}

			if disable {
				scheme = "http"
			}
		}

		// Check for address override
		var addr string
		addrRaw, ok := list.Config["address"]
		if !ok {
			addr = "127.0.0.1:8200"
		} else {
			addr = addrRaw.(string)
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

	// Configure the Circonus sink
	if telConfig.CirconusAPIToken != "" || telConfig.CirconusCheckSubmissionURL != "" {
		cfg := &circonus.Config{}
		cfg.Interval = telConfig.CirconusSubmissionInterval
		cfg.CheckManager.API.TokenKey = telConfig.CirconusAPIToken
		cfg.CheckManager.API.TokenApp = telConfig.CirconusAPIApp
		cfg.CheckManager.API.URL = telConfig.CirconusAPIURL
		cfg.CheckManager.Check.SubmissionURL = telConfig.CirconusCheckSubmissionURL
		cfg.CheckManager.Check.ID = telConfig.CirconusCheckID
		cfg.CheckManager.Check.ForceMetricActivation = telConfig.CirconusCheckForceMetricActivation
		cfg.CheckManager.Check.InstanceID = telConfig.CirconusCheckInstanceID
		cfg.CheckManager.Check.SearchTag = telConfig.CirconusCheckSearchTag
		cfg.CheckManager.Check.DisplayName = telConfig.CirconusCheckDisplayName
		cfg.CheckManager.Check.Tags = telConfig.CirconusCheckTags
		cfg.CheckManager.Broker.ID = telConfig.CirconusBrokerID
		cfg.CheckManager.Broker.SelectTag = telConfig.CirconusBrokerSelectTag

		if cfg.CheckManager.API.TokenApp == "" {
			cfg.CheckManager.API.TokenApp = "vault"
		}

		if cfg.CheckManager.Check.DisplayName == "" {
			cfg.CheckManager.Check.DisplayName = "Vault"
		}

		if cfg.CheckManager.Check.SearchTag == "" {
			cfg.CheckManager.Check.SearchTag = "service:vault"
		}

		sink, err := circonus.NewCirconusSink(cfg)
		if err != nil {
			return err
		}
		sink.Start()
		fanout = append(fanout, sink)
	}

	if telConfig.DogStatsDAddr != "" {
		var tags []string

		if telConfig.DogStatsDTags != nil {
			tags = telConfig.DogStatsDTags
		}

		sink, err := datadog.NewDogStatsdSink(telConfig.DogStatsDAddr, metricsConf.HostName)
		if err != nil {
			return fmt.Errorf("failed to start DogStatsD sink. Got: %s", err)
		}
		sink.SetTags(tags)
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

func (c *ServerCommand) Reload(lock *sync.RWMutex, reloadFuncs *map[string][]reload.ReloadFunc, configPath []string) error {
	lock.RLock()
	defer lock.RUnlock()

	var reloadErrors *multierror.Error

	for k, relFuncs := range *reloadFuncs {
		switch {
		case strings.HasPrefix(k, "listener|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(nil); err != nil {
						reloadErrors = multierror.Append(reloadErrors, fmt.Errorf("Error encountered reloading listener: %v", err))
					}
				}
			}

		case strings.HasPrefix(k, "audit_file|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(nil); err != nil {
						reloadErrors = multierror.Append(reloadErrors, fmt.Errorf("Error encountered reloading file audit backend at path %s: %v", strings.TrimPrefix(k, "audit_file|"), err))
					}
				}
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

  If the server is being started against a storage backend that is
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

func (c *ServerCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-config":             complete.PredictOr(complete.PredictFiles("*.hcl"), complete.PredictFiles("*.json")),
		"-dev":                complete.PredictNothing,
		"-dev-root-token-id":  complete.PredictNothing,
		"-dev-listen-address": complete.PredictNothing,
		"-log-level":          complete.PredictSet("trace", "debug", "info", "warn", "err"),
	}
}

// storePidFile is used to write out our PID to a file if necessary
func (c *ServerCommand) storePidFile(pidPath string) error {
	// Quit fast if no pidfile
	if pidPath == "" {
		return nil
	}

	// Open the PID file
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open pid file: %v", err)
	}
	defer pidFile.Close()

	// Write out the PID
	pid := os.Getpid()
	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return fmt.Errorf("could not write to pid file: %v", err)
	}
	return nil
}

// removePidFile is used to cleanup the PID file if necessary
func (c *ServerCommand) removePidFile(pidPath string) error {
	if pidPath == "" {
		return nil
	}
	return os.Remove(pidPath)
}

// MakeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// SIGINT or SIGTERM received.
func MakeShutdownCh() chan struct{} {
	resultCh := make(chan struct{})

	shutdownCh := make(chan os.Signal, 4)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownCh
		close(resultCh)
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

type grpclogFaker struct {
	logger log.Logger
	log    bool
}

func (g *grpclogFaker) Fatal(args ...interface{}) {
	g.logger.Error(fmt.Sprint(args...))
	os.Exit(1)
}

func (g *grpclogFaker) Fatalf(format string, args ...interface{}) {
	g.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (g *grpclogFaker) Fatalln(args ...interface{}) {
	g.logger.Error(fmt.Sprintln(args...))
	os.Exit(1)
}

func (g *grpclogFaker) Print(args ...interface{}) {
	if g.log && g.logger.IsTrace() {
		g.logger.Trace(fmt.Sprint(args...))
	}
}

func (g *grpclogFaker) Printf(format string, args ...interface{}) {
	if g.log && g.logger.IsTrace() {
		g.logger.Trace(fmt.Sprintf(format, args...))
	}
}

func (g *grpclogFaker) Println(args ...interface{}) {
	if g.log && g.logger.IsTrace() {
		g.logger.Trace(fmt.Sprintln(args...))
	}
}
