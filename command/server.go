package command

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	colorable "github.com/mattn/go-colorable"
	log "github.com/mgutz/logxi/v1"
	"github.com/mitchellh/cli"
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
	"github.com/hashicorp/vault/helper/gated-writer"
	"github.com/hashicorp/vault/helper/logbridge"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/helper/mlock"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/reload"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

var _ cli.Command = (*ServerCommand)(nil)
var _ cli.CommandAutocomplete = (*ServerCommand)(nil)

type ServerCommand struct {
	*BaseCommand

	AuditBackends      map[string]audit.Factory
	CredentialBackends map[string]logical.Factory
	LogicalBackends    map[string]logical.Factory
	PhysicalBackends   map[string]physical.Factory

	ShutdownCh chan struct{}
	SighupCh   chan struct{}

	WaitGroup *sync.WaitGroup

	logGate *gatedwriter.Writer
	logger  log.Logger

	cleanupGuard sync.Once

	reloadFuncsLock *sync.RWMutex
	reloadFuncs     *map[string][]reload.ReloadFunc
	startedCh       chan (struct{}) // for tests
	reloadedCh      chan (struct{}) // for tests

	// new stuff
	flagConfigs        []string
	flagLogLevel       string
	flagDev            bool
	flagDevRootTokenID string
	flagDevListenAddr  string

	flagDevPluginDir     string
	flagDevHA            bool
	flagDevLatency       int
	flagDevLatencyJitter int
	flagDevLeasedKV      bool
	flagDevSkipInit      bool
	flagDevThreeNode     bool
	flagDevFourCluster   bool
	flagDevTransactional bool
	flagTestVerifyOnly   bool
}

func (c *ServerCommand) Synopsis() string {
	return "Start a Vault server"
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: vault server [options]

  This command starts a Vault server that responds to API requests. By default,
  Vault will start in a "sealed" state. The Vault cluster must be initialized
  before use, usually by the "vault init" command. Each Vault server must also
  be unsealed using the "vault unseal" command or the API before the server can
  respond to requests.

  Start a server with a configuration file:

      $ vault server -config=/etc/vault/config.hcl

  Run in "dev" mode:

      $ vault server -dev -dev-root-token-id="root"

  For a full list of examples, please see the documentation.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *ServerCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&StringSliceVar{
		Name:   "config",
		Target: &c.flagConfigs,
		Completion: complete.PredictOr(
			complete.PredictFiles("*.hcl"),
			complete.PredictFiles("*.json"),
			complete.PredictDirs("*"),
		),
		Usage: "Path to a configuration file or directory of configuration " +
			"files. This flag can be specified multiple times to load multiple " +
			"configurations. If the path is a directory, all files which end in " +
			".hcl or .json are loaded.",
	})

	f.StringVar(&StringVar{
		Name:       "log-level",
		Target:     &c.flagLogLevel,
		Default:    "info",
		EnvVar:     "VAULT_LOG_LEVEL",
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "err"),
		Usage: "Log verbosity level. Supported values (in order of detail) are " +
			"\"trace\", \"debug\", \"info\", \"warn\", and \"err\".",
	})

	f = set.NewFlagSet("Dev Options")

	f.BoolVar(&BoolVar{
		Name:   "dev",
		Target: &c.flagDev,
		Usage: "Enable development mode. In this mode, Vault runs in-memory and " +
			"starts unsealed. As the name implies, do not run \"dev\" mode in " +
			"production.",
	})

	f.StringVar(&StringVar{
		Name:    "dev-root-token-id",
		Target:  &c.flagDevRootTokenID,
		Default: "",
		EnvVar:  "VAULT_DEV_ROOT_TOKEN_ID",
		Usage: "Initial root token. This only applies when running in \"dev\" " +
			"mode.",
	})

	f.StringVar(&StringVar{
		Name:    "dev-listen-address",
		Target:  &c.flagDevListenAddr,
		Default: "127.0.0.1:8200",
		EnvVar:  "VAULT_DEV_LISTEN_ADDRESS",
		Usage:   "Address to bind to in \"dev\" mode.",
	})

	// Internal-only flags to follow.
	//
	// Why hello there little source code reader! Welcome to the Vault source
	// code. The remaining options are intentionally undocumented and come with
	// no warranty or backwards-compatability promise. Do not use these flags
	// in production. Do not build automation using these flags. Unless you are
	// developing against Vault, you should not need any of these flags.

	f.StringVar(&StringVar{
		Name:       "dev-plugin-dir",
		Target:     &c.flagDevPluginDir,
		Default:    "",
		Completion: complete.PredictDirs("*"),
		Hidden:     true,
	})

	f.BoolVar(&BoolVar{
		Name:    "dev-ha",
		Target:  &c.flagDevHA,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "dev-transactional",
		Target:  &c.flagDevTransactional,
		Default: false,
		Hidden:  true,
	})

	f.IntVar(&IntVar{
		Name:   "dev-latency",
		Target: &c.flagDevLatency,
		Hidden: true,
	})

	f.IntVar(&IntVar{
		Name:   "dev-latency-jitter",
		Target: &c.flagDevLatencyJitter,
		Hidden: true,
	})

	f.BoolVar(&BoolVar{
		Name:    "dev-leased-kv",
		Target:  &c.flagDevLeasedKV,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "dev-skip-init",
		Target:  &c.flagDevSkipInit,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "dev-three-node",
		Target:  &c.flagDevThreeNode,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "dev-four-cluster",
		Target:  &c.flagDevFourCluster,
		Default: false,
		Hidden:  true,
	})

	// TODO: should this be a public flag?
	f.BoolVar(&BoolVar{
		Name:    "test-verify-only",
		Target:  &c.flagTestVerifyOnly,
		Default: false,
		Hidden:  true,
	})

	// End internal-only flags.

	return set
}

func (c *ServerCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	c.logGate = &gatedwriter.Writer{Writer: colorable.NewColorable(os.Stderr)}
	var level int
	c.flagLogLevel = strings.ToLower(strings.TrimSpace(c.flagLogLevel))
	switch c.flagLogLevel {
	case "trace":
		level = log.LevelTrace
	case "debug":
		level = log.LevelDebug
	case "info", "":
		level = log.LevelInfo
	case "notice":
		level = log.LevelNotice
	case "warn", "warning":
		level = log.LevelWarn
	case "err", "error":
		level = log.LevelError
	default:
		c.UI.Error(fmt.Sprintf("Unknown log level: %s", c.flagLogLevel))
		return 1
	}

	logFormat := os.Getenv("VAULT_LOG_FORMAT")
	if logFormat == "" {
		logFormat = os.Getenv("LOGXI_FORMAT")
	}
	switch strings.ToLower(logFormat) {
	case "vault", "vault_json", "vault-json", "vaultjson", "json", "":
		if c.flagDevThreeNode || c.flagDevFourCluster {
			c.logger = logbridge.NewLogger(hclog.New(&hclog.LoggerOptions{
				Mutex:  &sync.Mutex{},
				Output: c.logGate,
				Level:  hclog.Trace,
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

	// Automatically enable dev mode if other dev flags are provided.
	if c.flagDevHA || c.flagDevTransactional || c.flagDevLeasedKV || c.flagDevThreeNode || c.flagDevFourCluster {
		c.flagDev = true
	}

	// Validation
	if !c.flagDev {
		switch {
		case len(c.flagConfigs) == 0:
			c.UI.Error("Must specify at least one config path using -config")
			return 1
		case c.flagDevRootTokenID != "":
			c.UI.Warn(wrapAtLength(
				"You cannot specify a custom root token ID outside of \"dev\" mode. " +
					"Your request has been ignored."))
			c.flagDevRootTokenID = ""
		}
	}

	// Load the configuration
	var config *server.Config
	if c.flagDev {
		config = server.DevConfig(c.flagDevHA, c.flagDevTransactional)
		if c.flagDevListenAddr != "" {
			config.Listeners[0].Config["address"] = c.flagDevListenAddr
		}
	}
	for _, path := range c.flagConfigs {
		current, err := server.LoadConfig(path, c.logger)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error loading configuration from %s: %s", path, err))
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
		c.UI.Output(wrapAtLength(
			"No configuration files found. Please provide configurations with the " +
				"-config flag. If you are supply the path to a directory, please " +
				"ensure the directory contains files with the .hcl or .json " +
				"extension."))
		return 1
	}

	// Ensure that a backend is provided
	if config.Storage == nil {
		c.UI.Output("A storage backend must be specified")
		return 1
	}

	// If mlockall(2) isn't supported, show a warning. We disable this in dev
	// because it is quite scary to see when first using Vault. We also disable
	// this if the user has explicitly disabled mlock in configuration.
	if !c.flagDev && !config.DisableMlock && !mlock.Supported() {
		c.UI.Warn(wrapAtLength(
			"WARNING! mlock is not supported on this system! An mlockall(2)-like " +
				"syscall to prevent memory from being swapped to disk is not " +
				"supported on this system. For better security, only run Vault on " +
				"systems where this call is supported. If you are running Vault " +
				"in a Docker container, provide the IPC_LOCK cap to the container."))
	}

	if err := c.setupTelemetry(config); err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}

	// Initialize the backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
	if !exists {
		c.UI.Error(fmt.Sprintf("Unknown storage type %s", config.Storage.Type))
		return 1
	}
	backend, err := factory(config.Storage.Config, c.logger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing storage of type %s: %s", config.Storage.Type, err))
		return 1
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = c.flagLogLevel
	infoKeys = append(infoKeys, "log level")

	var seal vault.Seal = vault.NewDefaultSeal()

	// Ensure that the seal finalizer is called, even if using verify-only
	defer func() {
		if seal != nil {
			err = seal.Finalize(context.Background())
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
			}
		}
	}()

	if seal == nil {
		c.UI.Error(fmt.Sprintf("Could not create seal! Most likely proper Seal configuration information was not set, but no error was generated."))
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
	if c.flagDev {
		coreConfig.DevToken = c.flagDevRootTokenID
		if c.flagDevLeasedKV {
			coreConfig.LogicalBackends["kv"] = vault.LeasedPassthroughBackendFactory
		}
		if c.flagDevPluginDir != "" {
			coreConfig.PluginDirectory = c.flagDevPluginDir
		}
		if c.flagDevLatency > 0 {
			injectLatency := time.Duration(c.flagDevLatency) * time.Millisecond
			if _, txnOK := backend.(physical.Transactional); txnOK {
				coreConfig.Physical = physical.NewTransactionalLatencyInjector(backend, injectLatency, c.flagDevLatencyJitter, c.logger)
			} else {
				coreConfig.Physical = physical.NewLatencyInjector(backend, injectLatency, c.flagDevLatencyJitter, c.logger)
			}
		}
	}

	if c.flagDevThreeNode {
		return c.enableThreeNodeDevCluster(coreConfig, info, infoKeys, c.flagDevListenAddr, os.Getenv("VAULT_DEV_TEMP_DIR"))
	}

	var disableClustering bool

	// Initialize the separate HA storage backend, if it exists
	var ok bool
	if config.HAStorage != nil {
		factory, exists := c.PhysicalBackends[config.HAStorage.Type]
		if !exists {
			c.UI.Error(fmt.Sprintf("Unknown HA storage type %s", config.HAStorage.Type))
			return 1

		}
		habackend, err := factory(config.HAStorage.Config, c.logger)
		if err != nil {
			c.UI.Error(fmt.Sprintf(
				"Error initializing HA storage of type %s: %s", config.HAStorage.Type, err))
			return 1

		}

		if coreConfig.HAPhysical, ok = habackend.(physical.HABackend); !ok {
			c.UI.Error("Specified HA storage does not support HA")
			return 1
		}

		if !coreConfig.HAPhysical.HAEnabled() {
			c.UI.Error("Specified HA storage has HA support disabled; please consult documentation")
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
			c.UI.Error(fmt.Sprintf("Error detecting redirect address: %s", err))
		} else if redirect == "" {
			c.UI.Error("Failed to detect redirect address.")
		} else {
			coreConfig.RedirectAddr = redirect
		}
	}
	if coreConfig.RedirectAddr == "" && c.flagDev {
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
		case c.flagDev:
			addrToUse = fmt.Sprintf("http://%s", config.Listeners[0].Config["address"])
		default:
			goto CLUSTER_SYNTHESIS_COMPLETE
		}
		u, err := url.ParseRequestURI(addrToUse)
		if err != nil {
			c.UI.Error(fmt.Sprintf(
				"Error parsing synthesized cluster address %s: %v", addrToUse, err))
			return 1
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			// This sucks, as it's a const in the function but not exported in the package
			if strings.Contains(err.Error(), "missing port in address") {
				host = u.Host
				port = "443"
			} else {
				c.UI.Error(fmt.Sprintf("Error parsing redirect address: %v", err))
				return 1
			}
		}
		nPort, err := strconv.Atoi(port)
		if err != nil {
			c.UI.Error(fmt.Sprintf(
				"Error parsing synthesized address; failed to convert %q to a numeric: %v", port, err))
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
			c.UI.Error(fmt.Sprintf("Error parsing cluster address %s: %v", coreConfig.RedirectAddr, err))
			return 11
		}
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}

	// Initialize the core
	core, newCoreError := vault.NewCore(coreConfig)
	if newCoreError != nil {
		if !errwrap.ContainsType(newCoreError, new(vault.NonFatalError)) {
			c.UI.Error(fmt.Sprintf("Error initializing core: %s", newCoreError))
			return 1
		}
	}

	// Copy the reload funcs pointers back
	c.reloadFuncs = coreConfig.ReloadFuncs
	c.reloadFuncsLock = coreConfig.ReloadFuncsLock

	// Compile server information for output later
	info["storage"] = config.Storage.Type
	info["log level"] = c.flagLogLevel
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
		ln, props, reloadFunc, err := server.NewListener(lnConfig.Type, lnConfig.Config, c.logGate, c.UI)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing listener of type %s: %s", lnConfig.Type, err))
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
					c.UI.Error(fmt.Sprintf("Error resolving cluster_address: %s", err))
					return 1
				}
				clusterAddrs = append(clusterAddrs, tcpAddr)
			} else {
				tcpAddr, ok := ln.Addr().(*net.TCPAddr)
				if !ok {
					c.UI.Error("Failed to parse tcp listener")
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
	c.UI.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.UI.Output("")

	// Tests might not want to start a vault server and just want to verify
	// the configuration.
	if c.flagTestVerifyOnly {
		return 0
	}

	handler := vaulthttp.Handler(core)

	// This needs to happen before we first unseal, so before we trigger dev
	// mode if it's set
	core.SetClusterListenerAddrs(clusterAddrs)
	core.SetClusterHandler(handler)

	err = core.UnsealWithStoredKeys(context.Background())
	if err != nil {
		if !errwrap.ContainsType(err, new(vault.NonFatalError)) {
			c.UI.Error(fmt.Sprintf("Error initializing core: %s", err))
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
				c.UI.Error(fmt.Sprintf("Error initializing service discovery: %v", err))
				return 1
			}
		}
	}

	// If we're in Dev mode, then initialize the core
	if c.flagDev && !c.flagDevSkipInit {
		init, err := c.enableDev(core, coreConfig)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing Dev mode: %s", err))
			return 1
		}

		export := "export"
		quote := "'"
		if runtime.GOOS == "windows" {
			export = "set"
			quote = ""
		}

		// Print the big dev mode warning!
		c.UI.Warn(wrapAtLength(
			"WARNING! dev mode is enabled! In this mode, Vault runs entirely " +
				"in-memory and starts unsealed with a single unseal key. The root " +
				"token is already authenticated to the CLI, so you can immediately " +
				"begin using Vault."))
		c.UI.Warn("")
		c.UI.Warn("You may need to set the following environment variable:")
		c.UI.Warn("")
		c.UI.Warn(fmt.Sprintf("    $ %s VAULT_ADDR=%s%s%s",
			export, quote, "http://"+config.Listeners[0].Config["address"].(string), quote))

		// Unseal key is not returned if stored shares is supported
		if len(init.SecretShares) > 0 {
			c.UI.Warn("")
			c.UI.Warn(wrapAtLength(
				"The unseal key and root token are displayed below in case you want " +
					"to seal/unseal the Vault or re-authenticate."))
			c.UI.Warn("")
			c.UI.Warn(fmt.Sprintf("Unseal Key: %s", base64.StdEncoding.EncodeToString(init.SecretShares[0])))
		}

		if len(init.RecoveryShares) > 0 {
			c.UI.Warn("")
			c.UI.Warn(wrapAtLength(
				"The recovery key and root token are displayed below in case you want " +
					"to seal/unseal the Vault or re-authenticate."))
			c.UI.Warn("")
			c.UI.Warn(fmt.Sprintf("Unseal Key: %s", base64.StdEncoding.EncodeToString(init.RecoveryShares[0])))
		}

		c.UI.Warn(fmt.Sprintf("Root Token: %s", init.RootToken))

		c.UI.Warn("")
		c.UI.Warn(wrapAtLength(
			"Development mode should NOT be used in production installations!"))
		c.UI.Warn("")
	}

	// Initialize the HTTP servers
	for _, ln := range lns {
		server := &http.Server{
			Handler: handler,
		}
		go server.Serve(ln)
	}

	if newCoreError != nil {
		c.UI.Warn(wrapAtLength(
			"WARNING! A non-fatal error occurred during initialization. Please " +
				"check the logs for more information."))
		c.UI.Warn("")
	}

	// Output the header that the server has started
	c.UI.Output("==> Vault server started! Log data will stream in below:\n")

	// Inform any tests that the server is ready
	select {
	case c.startedCh <- struct{}{}:
	default:
	}

	// Release the log gate.
	c.logGate.Flush()

	// Write out the PID to the file now that server has successfully started
	if err := c.storePidFile(config.PidFile); err != nil {
		c.UI.Error(fmt.Sprintf("Error storing PID: %s", err))
		return 1
	}

	defer func() {
		if err := c.removePidFile(config.PidFile); err != nil {
			c.UI.Error(fmt.Sprintf("Error deleting the PID file: %s", err))
		}
	}()

	// Wait for shutdown
	shutdownTriggered := false

	for !shutdownTriggered {
		select {
		case <-c.ShutdownCh:
			c.UI.Output("==> Vault shutdown triggered")

			// Stop the listeners so that we don't process further client requests.
			c.cleanupGuard.Do(listenerCloseFunc)

			// Shutdown will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			if err := core.Shutdown(); err != nil {
				c.UI.Error(fmt.Sprintf("Error with core shutdown: %s", err))
			}

			shutdownTriggered = true

		case <-c.SighupCh:
			c.UI.Output("==> Vault reload triggered")
			if err := c.Reload(c.reloadFuncsLock, c.reloadFuncs, c.flagConfigs); err != nil {
				c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
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

	ctx := context.Background()

	// Initialize it with a basic single key
	init, err := core.Initialize(ctx, &vault.InitParams{
		BarrierConfig:  barrierConfig,
		RecoveryConfig: recoveryConfig,
	})
	if err != nil {
		return nil, err
	}

	// Handle unseal with stored keys
	if core.SealAccess().StoredKeysSupported() {
		err := core.UnsealWithStoredKeys(ctx)
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

func (c *ServerCommand) enableThreeNodeDevCluster(base *vault.CoreConfig, info map[string]string, infoKeys []string, devListenAddress, tempDir string) int {
	testCluster := vault.NewTestCluster(&testing.RuntimeT{}, base, &vault.TestClusterOptions{
		HandlerFunc:       vaulthttp.Handler,
		BaseListenAddress: c.flagDevListenAddr,
		RawLogger:         c.logger,
		TempDir:           tempDir,
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
	c.UI.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.UI.Output("")

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
			c.UI.Error(fmt.Sprintf("failed to create root token with ID %s: %s", base.DevToken, err))
			return 1
		}
		if resp == nil {
			c.UI.Error(fmt.Sprintf("nil response when creating root token with ID %s", base.DevToken))
			return 1
		}
		if resp.Auth == nil {
			c.UI.Error(fmt.Sprintf("nil auth when creating root token with ID %s", base.DevToken))
			return 1
		}

		testCluster.RootToken = resp.Auth.ClientToken

		req.ID = "dev-revoke-init-root"
		req.Path = "auth/token/revoke-self"
		req.Data = nil
		resp, err = testCluster.Cores[0].HandleRequest(req)
		if err != nil {
			c.UI.Output(fmt.Sprintf("failed to revoke initial root token: %s", err))
			return 1
		}
	}

	// Set the token
	tokenHelper, err := c.TokenHelper()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting token helper: %s", err))
		return 1
	}
	if err := tokenHelper.Store(testCluster.RootToken); err != nil {
		c.UI.Error(fmt.Sprintf("Error storing in token helper: %s", err))
		return 1
	}

	if err := ioutil.WriteFile(filepath.Join(testCluster.TempDir, "root_token"), []byte(testCluster.RootToken), 0755); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing token to tempfile: %s", err))
		return 1
	}

	c.UI.Output(fmt.Sprintf(
		"==> Three node dev mode is enabled\n\n" +
			"The unseal key and root token are reproduced below in case you\n" +
			"want to seal/unseal the Vault or play with authentication.\n",
	))

	for i, key := range testCluster.BarrierKeys {
		c.UI.Output(fmt.Sprintf(
			"Unseal Key %d: %s",
			i+1, base64.StdEncoding.EncodeToString(key),
		))
	}

	c.UI.Output(fmt.Sprintf(
		"\nRoot Token: %s\n", testCluster.RootToken,
	))

	c.UI.Output(fmt.Sprintf(
		"\nUseful env vars:\n"+
			"VAULT_TOKEN=%s\n"+
			"VAULT_ADDR=%s\n"+
			"VAULT_CACERT=%s/ca_cert.pem\n",
		testCluster.RootToken,
		testCluster.Cores[0].Client.Address(),
		testCluster.TempDir,
	))

	// Output the header that the server has started
	c.UI.Output("==> Vault server started! Log data will stream in below:\n")

	// Inform any tests that the server is ready
	select {
	case c.startedCh <- struct{}{}:
	default:
	}

	// Release the log gate.
	c.logGate.Flush()

	// Wait for shutdown
	shutdownTriggered := false

	for !shutdownTriggered {
		select {
		case <-c.ShutdownCh:
			c.UI.Output("==> Vault shutdown triggered")

			// Stop the listeners so that we don't process further client requests.
			c.cleanupGuard.Do(testCluster.Cleanup)

			// Shutdown will wait until after Vault is sealed, which means the
			// request forwarding listeners will also be closed (and also
			// waited for).
			for _, core := range testCluster.Cores {
				if err := core.Shutdown(); err != nil {
					c.UI.Error(fmt.Sprintf("Error with core shutdown: %s", err))
				}
			}

			shutdownTriggered = true

		case <-c.SighupCh:
			c.UI.Output("==> Vault reload triggered")
			for _, core := range testCluster.Cores {
				if err := c.Reload(core.ReloadFuncsLock, core.ReloadFuncs, nil); err != nil {
					c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
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
						reloadErrors = multierror.Append(reloadErrors, fmt.Errorf("Error encountered reloading file audit device at path %s: %v", strings.TrimPrefix(k, "audit_file|"), err))
					}
				}
			}
		}
	}

	// Send a message that we reloaded. This prevents "guessing" sleep times
	// in tests.
	select {
	case c.reloadedCh <- struct{}{}:
	default:
	}

	return reloadErrors.ErrorOrNil()
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
