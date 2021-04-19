package command

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/gatedwriter"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/internalshared/reloadutil"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/version"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/vault"
	vaultseal "github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-testing-interface"
	"github.com/posener/complete"
	"go.uber.org/atomic"
	"golang.org/x/net/http/httpproxy"
	"google.golang.org/grpc/grpclog"
)

var (
	_ cli.Command             = (*ServerCommand)(nil)
	_ cli.CommandAutocomplete = (*ServerCommand)(nil)
)

var memProfilerEnabled = false

var enableFourClusterDev = func(c *ServerCommand, base *vault.CoreConfig, info map[string]string, infoKeys []string, devListenAddress, tempDir string) int {
	c.logger.Error("-dev-four-cluster only supported in enterprise Vault")
	return 1
}

const (
	storageMigrationLock = "core/migration"

	// Even though there are more types than the ones below, the following consts
	// are declared internally for value comparison and reusability.
	storageTypeRaft   = "raft"
	storageTypeConsul = "consul"
)

type ServerCommand struct {
	*BaseCommand

	AuditBackends      map[string]audit.Factory
	CredentialBackends map[string]logical.Factory
	LogicalBackends    map[string]logical.Factory
	PhysicalBackends   map[string]physical.Factory

	ServiceRegistrations map[string]sr.Factory

	ShutdownCh chan struct{}
	SighupCh   chan struct{}
	SigUSR2Ch  chan struct{}

	WaitGroup *sync.WaitGroup

	logOutput   io.Writer
	gatedWriter *gatedwriter.Writer
	logger      log.InterceptLogger

	cleanupGuard sync.Once

	reloadFuncsLock *sync.RWMutex
	reloadFuncs     *map[string][]reloadutil.ReloadFunc
	startedCh       chan (struct{}) // for tests
	reloadedCh      chan (struct{}) // for tests

	allLoggers []log.Logger

	// new stuff
	flagConfigs            []string
	flagLogLevel           string
	flagLogFormat          string
	flagRecovery           bool
	flagDev                bool
	flagDevRootTokenID     string
	flagDevListenAddr      string
	flagDevNoStoreToken    bool
	flagDevPluginDir       string
	flagDevPluginInit      bool
	flagDevHA              bool
	flagDevLatency         int
	flagDevLatencyJitter   int
	flagDevLeasedKV        bool
	flagDevKVV1            bool
	flagDevSkipInit        bool
	flagDevThreeNode       bool
	flagDevFourCluster     bool
	flagDevTransactional   bool
	flagDevAutoSeal        bool
	flagTestVerifyOnly     bool
	flagCombineLogs        bool
	flagTestServerConfig   bool
	flagDevConsul          bool
	flagExitOnCoreShutdown bool
	flagDiagnose           string
}

func (c *ServerCommand) Synopsis() string {
	return "Start a Vault server"
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: vault server [options]

  This command starts a Vault server that responds to API requests. By default,
  Vault will start in a "sealed" state. The Vault cluster must be initialized
  before use, usually by the "vault operator init" command. Each Vault server must
  also be unsealed using the "vault operator unseal" command or the API before the
  server can respond to requests.

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
		Default:    notSetValue,
		EnvVar:     "VAULT_LOG_LEVEL",
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "err"),
		Usage: "Log verbosity level. Supported values (in order of detail) are " +
			"\"trace\", \"debug\", \"info\", \"warn\", and \"err\".",
	})

	f.StringVar(&StringVar{
		Name:    "log-format",
		Target:  &c.flagLogFormat,
		Default: notSetValue,
		// EnvVar can't be just "VAULT_LOG_FORMAT", because more than one env var name is supported
		// for backwards compatibility reasons.
		// See github.com/hashicorp/vault/sdk/helper/logging.ParseEnvLogFormat()
		Completion: complete.PredictSet("standard", "json"),
		Usage:      `Log format. Supported values are "standard" and "json".`,
	})

	f.BoolVar(&BoolVar{
		Name:    "exit-on-core-shutdown",
		Target:  &c.flagExitOnCoreShutdown,
		Default: false,
		Usage:   "Exit the vault server if the vault core is shutdown.",
	})

	f.BoolVar(&BoolVar{
		Name:   "recovery",
		Target: &c.flagRecovery,
		Usage: "Enable recovery mode. In this mode, Vault is used to perform recovery actions." +
			"Using a recovery operation token, \"sys/raw\" API can be used to manipulate the storage.",
	})

	// Disabled by default until functional
	if os.Getenv(OperatorDiagnoseEnableEnv) != "" {
		f.StringVar(&StringVar{
			Name:    "diagnose",
			Target:  &c.flagDiagnose,
			Default: notSetValue,
			Usage:   "Run diagnostics before starting Vault. Specify a filename to direct output to that file.",
		})
	} else {
		// Ensure diagnose is *not* run when feature flag is off.
		c.flagDiagnose = notSetValue
	}

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
	f.BoolVar(&BoolVar{
		Name:    "dev-no-store-token",
		Target:  &c.flagDevNoStoreToken,
		Default: false,
		Usage: "Do not persist the dev root token to the token helper " +
			"(usually the local filesystem) for use in future requests. " +
			"The token will only be displayed in the command output.",
	})

	// Internal-only flags to follow.
	//
	// Why hello there little source code reader! Welcome to the Vault source
	// code. The remaining options are intentionally undocumented and come with
	// no warranty or backwards-compatibility promise. Do not use these flags
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
		Name:    "dev-plugin-init",
		Target:  &c.flagDevPluginInit,
		Default: true,
		Hidden:  true,
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
		Name:    "dev-kv-v1",
		Target:  &c.flagDevKVV1,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "dev-auto-seal",
		Target:  &c.flagDevAutoSeal,
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

	f.BoolVar(&BoolVar{
		Name:    "dev-consul",
		Target:  &c.flagDevConsul,
		Default: false,
		Hidden:  true,
	})

	// TODO: should the below flags be public?
	f.BoolVar(&BoolVar{
		Name:    "combine-logs",
		Target:  &c.flagCombineLogs,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "test-verify-only",
		Target:  &c.flagTestVerifyOnly,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "test-server-config",
		Target:  &c.flagTestServerConfig,
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

func (c *ServerCommand) flushLog() {
	c.logger.(hclog.OutputResettable).ResetOutputWithFlush(&hclog.LoggerOptions{
		Output: c.logOutput,
	}, c.gatedWriter)
}

func (c *ServerCommand) parseConfig() (*server.Config, error) {
	// Load the configuration
	var config *server.Config
	for _, path := range c.flagConfigs {
		current, err := server.LoadConfig(path)
		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("error loading configuration from %s: {{err}}", path), err)
		}

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}
	return config, nil
}

func (c *ServerCommand) runRecoveryMode() int {
	config, err := c.parseConfig()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Ensure at least one config was found.
	if config == nil {
		c.UI.Output(wrapAtLength(
			"No configuration files found. Please provide configurations with the " +
				"-config flag. If you are supplying the path to a directory, please " +
				"ensure the directory contains files with the .hcl or .json " +
				"extension."))
		return 1
	}

	level, logLevelString, logLevelWasNotSet, logFormat, err := c.processLogLevelAndFormat(config)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.logger = log.NewInterceptLogger(&log.LoggerOptions{
		Output: c.gatedWriter,
		Level:  level,
		// Note that if logFormat is either unspecified or standard, then
		// the resulting logger's format will be standard.
		JSONFormat: logFormat == logging.JSONFormat,
	})

	// Ensure logging is flushed if initialization fails
	defer c.flushLog()

	logLevelStr, err := c.adjustLogLevel(config, logLevelWasNotSet)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	if logLevelStr != "" {
		logLevelString = logLevelStr
	}

	// create GRPC logger
	namedGRPCLogFaker := c.logger.Named("grpclogfaker")
	grpclog.SetLogger(&grpclogFaker{
		logger: namedGRPCLogFaker,
		log:    os.Getenv("VAULT_GRPC_LOGGING") != "",
	})

	if config.Storage == nil {
		c.UI.Output("A storage backend must be specified")
		return 1
	}

	if config.DefaultMaxRequestDuration != 0 {
		vault.DefaultMaxRequestDuration = config.DefaultMaxRequestDuration
	}

	logProxyEnvironmentVariables(c.logger)

	// Initialize the storage backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
	if !exists {
		c.UI.Error(fmt.Sprintf("Unknown storage type %s", config.Storage.Type))
		return 1
	}
	if config.Storage.Type == storageTypeRaft || (config.HAStorage != nil && config.HAStorage.Type == storageTypeRaft) {
		if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
			config.ClusterAddr = envCA
		}

		if len(config.ClusterAddr) == 0 {
			c.UI.Error("Cluster address must be set when using raft storage")
			return 1
		}
	}

	namedStorageLogger := c.logger.Named("storage." + config.Storage.Type)
	backend, err := factory(config.Storage.Config, namedStorageLogger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing storage of type %s: %s", config.Storage.Type, err))
		return 1
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = logLevelString
	infoKeys = append(infoKeys, "log level")

	var barrierSeal vault.Seal
	var sealConfigError error
	var wrapper wrapping.Wrapper

	if len(config.Seals) == 0 {
		config.Seals = append(config.Seals, &configutil.KMS{Type: wrapping.Shamir})
	}

	if len(config.Seals) > 1 {
		c.UI.Error("Only one seal block is accepted in recovery mode")
		return 1
	}

	configSeal := config.Seals[0]
	sealType := wrapping.Shamir
	if !configSeal.Disabled && os.Getenv("VAULT_SEAL_TYPE") != "" {
		sealType = os.Getenv("VAULT_SEAL_TYPE")
		configSeal.Type = sealType
	} else {
		sealType = configSeal.Type
	}

	infoKeys = append(infoKeys, "Seal Type")
	info["Seal Type"] = sealType

	var seal vault.Seal
	defaultSeal := vault.NewDefaultSeal(&vaultseal.Access{
		Wrapper: aeadwrapper.NewShamirWrapper(&wrapping.WrapperOptions{
			Logger: c.logger.Named("shamir"),
		}),
	})
	sealLogger := c.logger.ResetNamed(fmt.Sprintf("seal.%s", sealType))
	wrapper, sealConfigError = configutil.ConfigureWrapper(configSeal, &infoKeys, &info, sealLogger)
	if sealConfigError != nil {
		if !errwrap.ContainsType(sealConfigError, new(logical.KeyNotFoundError)) {
			c.UI.Error(fmt.Sprintf(
				"Error parsing Seal configuration: %s", sealConfigError))
			return 1
		}
	}
	if wrapper == nil {
		seal = defaultSeal
	} else {
		seal = vault.NewAutoSeal(&vaultseal.Access{
			Wrapper: wrapper,
		})
	}
	barrierSeal = seal

	// Ensure that the seal finalizer is called, even if using verify-only
	defer func() {
		err = seal.Finalize(context.Background())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
		}
	}()

	coreConfig := &vault.CoreConfig{
		Physical:     backend,
		StorageType:  config.Storage.Type,
		Seal:         barrierSeal,
		Logger:       c.logger,
		DisableMlock: config.DisableMlock,
		RecoveryMode: c.flagRecovery,
		ClusterAddr:  config.ClusterAddr,
	}

	core, newCoreError := vault.NewCore(coreConfig)
	if newCoreError != nil {
		if vault.IsFatalError(newCoreError) {
			c.UI.Error(fmt.Sprintf("Error initializing core: %s", newCoreError))
			return 1
		}
	}

	if err := core.InitializeRecovery(context.Background()); err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing core in recovery mode: %s", err))
		return 1
	}

	// Compile server information for output later
	infoKeys = append(infoKeys, "storage")
	info["storage"] = config.Storage.Type

	if coreConfig.ClusterAddr != "" {
		info["cluster address"] = coreConfig.ClusterAddr
		infoKeys = append(infoKeys, "cluster address")
	}

	// Initialize the listeners
	lns := make([]listenerutil.Listener, 0, len(config.Listeners))
	for _, lnConfig := range config.Listeners {
		ln, _, _, err := server.NewListener(lnConfig, c.gatedWriter, c.UI)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing listener of type %s: %s", lnConfig.Type, err))
			return 1
		}

		lns = append(lns, listenerutil.Listener{
			Listener: ln,
			Config:   lnConfig,
		})
	}

	listenerCloseFunc := func() {
		for _, ln := range lns {
			ln.Listener.Close()
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

	infoKeys = append(infoKeys, "recovery mode")
	info["recovery mode"] = "true"

	infoKeys = append(infoKeys, "go version")
	info["go version"] = runtime.Version()

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

	for _, ln := range lns {
		handler := vaulthttp.Handler(&vault.HandlerProperties{
			Core:                  core,
			ListenerConfig:        ln.Config,
			DisablePrintableCheck: config.DisablePrintableCheck,
			RecoveryMode:          c.flagRecovery,
			RecoveryToken:         atomic.NewString(""),
		})

		server := &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			IdleTimeout:       5 * time.Minute,
			ErrorLog:          c.logger.StandardLogger(nil),
		}

		go server.Serve(ln.Listener)
	}

	if sealConfigError != nil {
		init, err := core.InitializedLocally(context.Background())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking if core is initialized: %v", err))
			return 1
		}
		if init {
			c.UI.Error("Vault is initialized but no Seal key could be loaded")
			return 1
		}
	}

	if newCoreError != nil {
		c.UI.Warn(wrapAtLength(
			"WARNING! A non-fatal error occurred during initialization. Please " +
				"check the logs for more information."))
		c.UI.Warn("")
	}

	if !c.flagCombineLogs {
		c.UI.Output("==> Vault server started! Log data will stream in below:\n")
	}

	c.flushLog()

	for {
		select {
		case <-c.ShutdownCh:
			c.UI.Output("==> Vault shutdown triggered")

			c.cleanupGuard.Do(listenerCloseFunc)

			if err := core.Shutdown(); err != nil {
				c.UI.Error(fmt.Sprintf("Error with core shutdown: %s", err))
			}

			return 0

		case <-c.SigUSR2Ch:
			buf := make([]byte, 32*1024*1024)
			n := runtime.Stack(buf[:], true)
			c.logger.Info("goroutine trace", "stack", string(buf[:n]))
		}
	}

	return 0
}

func logProxyEnvironmentVariables(logger hclog.Logger) {
	proxyCfg := httpproxy.FromEnvironment()
	cfgMap := map[string]string{
		"http_proxy":  proxyCfg.HTTPProxy,
		"https_proxy": proxyCfg.HTTPSProxy,
		"no_proxy":    proxyCfg.NoProxy,
	}
	for k, v := range cfgMap {
		u, err := url.Parse(v)
		if err != nil {
			// Env vars may contain URLs or host:port values.  We only care
			// about the former.
			continue
		}
		if _, ok := u.User.Password(); ok {
			u.User = url.UserPassword("redacted-username", "redacted-password")
		} else if user := u.User.Username(); user != "" {
			u.User = url.User("redacted-username")
		}
		cfgMap[k] = u.String()
	}
	logger.Info("proxy environment", "http_proxy", cfgMap["http_proxy"],
		"https_proxy", cfgMap["https_proxy"], "no_proxy", cfgMap["no_proxy"])
}

func (c *ServerCommand) adjustLogLevel(config *server.Config, logLevelWasNotSet bool) (string, error) {
	var logLevelString string
	if config.LogLevel != "" && logLevelWasNotSet {
		configLogLevel := strings.ToLower(strings.TrimSpace(config.LogLevel))
		logLevelString = configLogLevel
		switch configLogLevel {
		case "trace":
			c.logger.SetLevel(log.Trace)
		case "debug":
			c.logger.SetLevel(log.Debug)
		case "notice", "info", "":
			c.logger.SetLevel(log.Info)
		case "warn", "warning":
			c.logger.SetLevel(log.Warn)
		case "err", "error":
			c.logger.SetLevel(log.Error)
		default:
			return "", fmt.Errorf("unknown log level: %s", config.LogLevel)
		}
	}
	return logLevelString, nil
}

func (c *ServerCommand) processLogLevelAndFormat(config *server.Config) (log.Level, string, bool, logging.LogFormat, error) {
	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	c.logOutput = os.Stderr
	if c.flagCombineLogs {
		c.logOutput = os.Stdout
	}
	c.gatedWriter = gatedwriter.NewWriter(c.logOutput)
	var level log.Level
	var logLevelWasNotSet bool
	logFormat := logging.UnspecifiedFormat
	logLevelString := c.flagLogLevel
	c.flagLogLevel = strings.ToLower(strings.TrimSpace(c.flagLogLevel))
	switch c.flagLogLevel {
	case notSetValue, "":
		logLevelWasNotSet = true
		logLevelString = "info"
		level = log.Info
	case "trace":
		level = log.Trace
	case "debug":
		level = log.Debug
	case "notice", "info":
		level = log.Info
	case "warn", "warning":
		level = log.Warn
	case "err", "error":
		level = log.Error
	default:
		return level, logLevelString, logLevelWasNotSet, logFormat, fmt.Errorf("unknown log level: %s", c.flagLogLevel)
	}

	if c.flagLogFormat != notSetValue {
		var err error
		logFormat, err = logging.ParseLogFormat(c.flagLogFormat)
		if err != nil {
			return level, logLevelString, logLevelWasNotSet, logFormat, err
		}
	}
	if logFormat == logging.UnspecifiedFormat {
		logFormat = logging.ParseEnvLogFormat()
	}
	if logFormat == logging.UnspecifiedFormat {
		var err error
		logFormat, err = logging.ParseLogFormat(config.LogFormat)
		if err != nil {
			return level, logLevelString, logLevelWasNotSet, logFormat, err
		}
	}

	return level, logLevelString, logLevelWasNotSet, logFormat, nil
}

type quiescenceSink struct {
	t *time.Timer
}

func (q quiescenceSink) Accept(name string, level log.Level, msg string, args ...interface{}) {
	q.t.Reset(100 * time.Millisecond)
}

func (c *ServerCommand) setupStorage(config *server.Config) (physical.Backend, error) {
	// Ensure that a backend is provided
	if config.Storage == nil {
		return nil, fmt.Errorf("A storage backend must be specified")
	}

	// Initialize the backend
	factory, exists := c.PhysicalBackends[config.Storage.Type]
	if !exists {
		return nil, fmt.Errorf("Unknown storage type %s", config.Storage.Type)
	}

	// Do any custom configuration needed per backend
	switch config.Storage.Type {
	case storageTypeConsul:
		if config.ServiceRegistration == nil {
			// If Consul is configured for storage and service registration is unconfigured,
			// use Consul for service registration without requiring additional configuration.
			// This maintains backward-compatibility.
			config.ServiceRegistration = &server.ServiceRegistration{
				Type:   "consul",
				Config: config.Storage.Config,
			}
		}
	case storageTypeRaft:
		if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
			config.ClusterAddr = envCA
		}
		if len(config.ClusterAddr) == 0 {
			return nil, fmt.Errorf("Cluster address must be set when using raft storage")
		}
	}

	namedStorageLogger := c.logger.Named("storage." + config.Storage.Type)
	c.allLoggers = append(c.allLoggers, namedStorageLogger)
	backend, err := factory(config.Storage.Config, namedStorageLogger)
	if err != nil {
		return nil, fmt.Errorf("Error initializing storage of type %s: %w", config.Storage.Type, err)
	}

	return backend, nil
}

// InitListeners returns a response code, error message, Listeners, and a TCP Address list.
func (c *ServerCommand) InitListeners(config *server.Config, disableClustering bool, infoKeys *[]string, info *map[string]string) (int, []listenerutil.Listener, []*net.TCPAddr, error) {
	clusterAddrs := []*net.TCPAddr{}

	// Initialize the listeners
	lns := make([]listenerutil.Listener, 0, len(config.Listeners))

	c.reloadFuncsLock.Lock()

	defer c.reloadFuncsLock.Unlock()

	var errMsg error
	for i, lnConfig := range config.Listeners {
		ln, props, reloadFunc, err := server.NewListener(lnConfig, c.gatedWriter, c.UI)
		if err != nil {
			errMsg = fmt.Errorf("Error initializing listener of type %s: %s", lnConfig.Type, err)
			return 1, nil, nil, errMsg
		}

		if reloadFunc != nil {
			relSlice := (*c.reloadFuncs)["listener|"+lnConfig.Type]
			relSlice = append(relSlice, reloadFunc)
			(*c.reloadFuncs)["listener|"+lnConfig.Type] = relSlice
		}

		if !disableClustering && lnConfig.Type == "tcp" {
			addr := lnConfig.ClusterAddress
			if addr != "" {
				tcpAddr, err := net.ResolveTCPAddr("tcp", lnConfig.ClusterAddress)
				if err != nil {
					errMsg = fmt.Errorf("Error resolving cluster_address: %s", err)
					return 1, nil, nil, errMsg
				}
				clusterAddrs = append(clusterAddrs, tcpAddr)
			} else {
				tcpAddr, ok := ln.Addr().(*net.TCPAddr)
				if !ok {
					errMsg = fmt.Errorf("Failed to parse tcp listener")
					return 1, nil, nil, errMsg
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

		if lnConfig.MaxRequestSize == 0 {
			lnConfig.MaxRequestSize = vaulthttp.DefaultMaxRequestSize
		}
		props["max_request_size"] = fmt.Sprintf("%d", lnConfig.MaxRequestSize)

		if lnConfig.MaxRequestDuration == 0 {
			lnConfig.MaxRequestDuration = vault.DefaultMaxRequestDuration
		}
		props["max_request_duration"] = lnConfig.MaxRequestDuration.String()

		lns = append(lns, listenerutil.Listener{
			Listener: ln,
			Config:   lnConfig,
		})

		// Store the listener props for output later
		key := fmt.Sprintf("listener %d", i+1)
		propsList := make([]string, 0, len(props))
		for k, v := range props {
			propsList = append(propsList, fmt.Sprintf(
				"%s: %q", k, v))
		}
		sort.Strings(propsList)
		*infoKeys = append(*infoKeys, key)
		(*info)[key] = fmt.Sprintf(
			"%s (%s)", lnConfig.Type, strings.Join(propsList, ", "))

	}
	if !disableClustering {
		if c.logger.IsDebug() {
			c.logger.Debug("cluster listener addresses synthesized", "cluster_addresses", clusterAddrs)
		}
	}
	return 0, lns, clusterAddrs, nil
}

func (c *ServerCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.flagRecovery {
		return c.runRecoveryMode()
	}

	// Automatically enable dev mode if other dev flags are provided.
	if c.flagDevConsul || c.flagDevHA || c.flagDevTransactional || c.flagDevLeasedKV || c.flagDevThreeNode || c.flagDevFourCluster || c.flagDevAutoSeal || c.flagDevKVV1 {
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

	if c.flagDiagnose != notSetValue {
		if c.flagDev {
			c.UI.Error("Cannot run diagnose on Vault in dev mode.")
			return 1
		}
		// TODO: add a file output flag to Diagnose
		diagnose := &OperatorDiagnoseCommand{
			BaseCommand: c.BaseCommand,
			flagDebug:   false,
			flagSkips:   []string{},
			flagConfigs: c.flagConfigs,
		}
		diagnose.RunWithParsedFlags()
	}

	// Load the configuration
	var config *server.Config
	var err error
	if c.flagDev {
		var devStorageType string
		switch {
		case c.flagDevConsul:
			devStorageType = "consul"
		case c.flagDevHA && c.flagDevTransactional:
			devStorageType = "inmem_transactional_ha"
		case !c.flagDevHA && c.flagDevTransactional:
			devStorageType = "inmem_transactional"
		case c.flagDevHA && !c.flagDevTransactional:
			devStorageType = "inmem_ha"
		default:
			devStorageType = "inmem"
		}
		config, err = server.DevConfig(devStorageType)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		if c.flagDevListenAddr != "" {
			config.Listeners[0].Address = c.flagDevListenAddr
		}
		config.Listeners[0].Telemetry.UnauthenticatedMetricsAccess = true
	}

	parsedConfig, err := c.parseConfig()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	if config == nil {
		config = parsedConfig
	} else {
		config = config.Merge(parsedConfig)
	}

	// Ensure at least one config was found.
	if config == nil {
		c.UI.Output(wrapAtLength(
			"No configuration files found. Please provide configurations with the " +
				"-config flag. If you are supplying the path to a directory, please " +
				"ensure the directory contains files with the .hcl or .json " +
				"extension."))
		return 1
	}

	level, logLevelString, logLevelWasNotSet, logFormat, err := c.processLogLevelAndFormat(config)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	config.LogFormat = logFormat.String()

	if c.flagDevThreeNode || c.flagDevFourCluster {
		c.logger = log.NewInterceptLogger(&log.LoggerOptions{
			Mutex:  &sync.Mutex{},
			Output: c.gatedWriter,
			Level:  log.Trace,
		})
	} else {
		c.logger = log.NewInterceptLogger(&log.LoggerOptions{
			Output: c.gatedWriter,
			Level:  level,
			// Note that if logFormat is either unspecified or standard, then
			// the resulting logger's format will be standard.
			JSONFormat: logFormat == logging.JSONFormat,
		})
	}

	// Ensure logging is flushed if initialization fails
	defer c.flushLog()

	c.allLoggers = []log.Logger{c.logger}

	logLevelStr, err := c.adjustLogLevel(config, logLevelWasNotSet)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	if logLevelStr != "" {
		logLevelString = logLevelStr
	}

	// create GRPC logger
	namedGRPCLogFaker := c.logger.Named("grpclogfaker")
	c.allLoggers = append(c.allLoggers, namedGRPCLogFaker)
	grpclog.SetLogger(&grpclogFaker{
		logger: namedGRPCLogFaker,
		log:    os.Getenv("VAULT_GRPC_LOGGING") != "",
	})

	if memProfilerEnabled {
		c.startMemProfiler()
	}

	if config.DefaultMaxRequestDuration != 0 {
		vault.DefaultMaxRequestDuration = config.DefaultMaxRequestDuration
	}

	logProxyEnvironmentVariables(c.logger)

	if envMlock := os.Getenv("VAULT_DISABLE_MLOCK"); envMlock != "" {
		var err error
		config.DisableMlock, err = strconv.ParseBool(envMlock)
		if err != nil {
			c.UI.Output("Error parsing the environment variable VAULT_DISABLE_MLOCK")
			return 1
		}
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

	inmemMetrics, metricSink, prometheusEnabled, err := configutil.SetupTelemetry(&configutil.SetupTelemetryOpts{
		Config:      config.Telemetry,
		Ui:          c.UI,
		ServiceName: "vault",
		DisplayName: "Vault",
		UserAgent:   useragent.String(),
		ClusterName: config.ClusterName,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}
	metricsHelper := metricsutil.NewMetricsHelper(inmemMetrics, prometheusEnabled)

	// Initialize the storage backend
	backend, err := c.setupStorage(config)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// If this is the Raft backend, add the metricSink to it
	if rb, ok := backend.(*raft.RaftBackend); ok {
		rb.SetMetricsSink(metricSink)
	}

	// Prevent server startup if migration is active
	// TODO: how to incorporate this check into Diagnose?
	if c.storageMigrationActive(backend) {
		return 1
	}

	// Initialize the Service Discovery, if there is one
	var configSR sr.ServiceRegistration
	if config.ServiceRegistration != nil {
		sdFactory, ok := c.ServiceRegistrations[config.ServiceRegistration.Type]
		if !ok {
			c.UI.Error(fmt.Sprintf("Unknown service_registration type %s", config.ServiceRegistration.Type))
			return 1
		}

		namedSDLogger := c.logger.Named("service_registration." + config.ServiceRegistration.Type)
		c.allLoggers = append(c.allLoggers, namedSDLogger)

		// Since we haven't even begun starting Vault's core yet,
		// we know that Vault is in its pre-running state.
		state := sr.State{
			VaultVersion:         version.GetVersion().VersionNumber(),
			IsInitialized:        false,
			IsSealed:             true,
			IsActive:             false,
			IsPerformanceStandby: false,
		}
		configSR, err = sdFactory(config.ServiceRegistration.Config, namedSDLogger, state)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing service_registration of type %s: %s", config.ServiceRegistration.Type, err))
			return 1
		}
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = logLevelString
	infoKeys = append(infoKeys, "log level")

	var barrierSeal vault.Seal
	var unwrapSeal vault.Seal

	var sealConfigError error
	var wrapper wrapping.Wrapper
	var barrierWrapper wrapping.Wrapper
	if c.flagDevAutoSeal {
		barrierSeal = vault.NewAutoSeal(vaultseal.NewTestSeal(nil))
	} else {
		// Handle the case where no seal is provided
		switch len(config.Seals) {
		case 0:
			config.Seals = append(config.Seals, &configutil.KMS{Type: wrapping.Shamir})
		case 1:
			// If there's only one seal and it's disabled assume they want to
			// migrate to a shamir seal and simply didn't provide it
			if config.Seals[0].Disabled {
				config.Seals = append(config.Seals, &configutil.KMS{Type: wrapping.Shamir})
			}
		}
		for _, configSeal := range config.Seals {
			sealType := wrapping.Shamir
			if !configSeal.Disabled && os.Getenv("VAULT_SEAL_TYPE") != "" {
				sealType = os.Getenv("VAULT_SEAL_TYPE")
				configSeal.Type = sealType
			} else {
				sealType = configSeal.Type
			}

			var seal vault.Seal
			sealLogger := c.logger.ResetNamed(fmt.Sprintf("seal.%s", sealType))
			c.allLoggers = append(c.allLoggers, sealLogger)
			defaultSeal := vault.NewDefaultSeal(&vaultseal.Access{
				Wrapper: aeadwrapper.NewShamirWrapper(&wrapping.WrapperOptions{
					Logger: c.logger.Named("shamir"),
				}),
			})
			var sealInfoKeys []string
			sealInfoMap := map[string]string{}
			wrapper, sealConfigError = configutil.ConfigureWrapper(configSeal, &sealInfoKeys, &sealInfoMap, sealLogger)
			if sealConfigError != nil {
				if !errwrap.ContainsType(sealConfigError, new(logical.KeyNotFoundError)) {
					c.UI.Error(fmt.Sprintf(
						"Error parsing Seal configuration: %s", sealConfigError))
					return 1
				}
			}
			if wrapper == nil {
				seal = defaultSeal
			} else {
				seal = vault.NewAutoSeal(&vaultseal.Access{
					Wrapper: wrapper,
				})
			}

			infoPrefix := ""
			if configSeal.Disabled {
				unwrapSeal = seal
				infoPrefix = "Old "
			} else {
				barrierSeal = seal
				barrierWrapper = wrapper
			}
			for _, k := range sealInfoKeys {
				infoKeys = append(infoKeys, infoPrefix+k)
				info[infoPrefix+k] = sealInfoMap[k]
			}

			// Ensure that the seal finalizer is called, even if using verify-only
			defer func() {
				err = seal.Finalize(context.Background())
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
				}
			}()

		}
	}

	if barrierSeal == nil {
		c.UI.Error(fmt.Sprintf("Could not create barrier seal! Most likely proper Seal configuration information was not set, but no error was generated."))
		return 1
	}

	// prepare a secure random reader for core
	secureRandomReader, err := configutil.CreateSecureRandomReaderFunc(config.SharedConfig, barrierWrapper)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	coreConfig := &vault.CoreConfig{
		RawConfig:                      config,
		Physical:                       backend,
		RedirectAddr:                   config.Storage.RedirectAddr,
		StorageType:                    config.Storage.Type,
		HAPhysical:                     nil,
		ServiceRegistration:            configSR,
		Seal:                           barrierSeal,
		UnwrapSeal:                     unwrapSeal,
		AuditBackends:                  c.AuditBackends,
		CredentialBackends:             c.CredentialBackends,
		LogicalBackends:                c.LogicalBackends,
		Logger:                         c.logger,
		DisableSentinelTrace:           config.DisableSentinelTrace,
		DisableCache:                   config.DisableCache,
		DisableMlock:                   config.DisableMlock,
		MaxLeaseTTL:                    config.MaxLeaseTTL,
		DefaultLeaseTTL:                config.DefaultLeaseTTL,
		ClusterName:                    config.ClusterName,
		CacheSize:                      config.CacheSize,
		PluginDirectory:                config.PluginDirectory,
		EnableUI:                       config.EnableUI,
		EnableRaw:                      config.EnableRawEndpoint,
		DisableSealWrap:                config.DisableSealWrap,
		DisablePerformanceStandby:      config.DisablePerformanceStandby,
		DisableIndexing:                config.DisableIndexing,
		AllLoggers:                     c.allLoggers,
		BuiltinRegistry:                builtinplugins.Registry,
		DisableKeyEncodingChecks:       config.DisablePrintableCheck,
		MetricsHelper:                  metricsHelper,
		MetricSink:                     metricSink,
		SecureRandomReader:             secureRandomReader,
		EnableResponseHeaderHostname:   config.EnableResponseHeaderHostname,
		EnableResponseHeaderRaftNodeID: config.EnableResponseHeaderRaftNodeID,
	}
	if c.flagDev {
		coreConfig.EnableRaw = true
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

	if c.flagDevFourCluster {
		return enableFourClusterDev(c, coreConfig, info, infoKeys, c.flagDevListenAddr, os.Getenv("VAULT_DEV_TEMP_DIR"))
	}

	var disableClustering bool

	// Initialize the separate HA storage backend, if it exists
	var ok bool
	if config.HAStorage != nil {
		if config.Storage.Type == storageTypeRaft && config.HAStorage.Type == storageTypeRaft {
			c.UI.Error("Raft cannot be set both as 'storage' and 'ha_storage'. Setting 'storage' to 'raft' will automatically set it up for HA operations as well")
			return 1
		}

		if config.Storage.Type == storageTypeRaft {
			c.UI.Error("HA storage cannot be declared when Raft is the storage type")
			return 1
		}

		factory, exists := c.PhysicalBackends[config.HAStorage.Type]
		if !exists {
			c.UI.Error(fmt.Sprintf("Unknown HA storage type %s", config.HAStorage.Type))
			return 1

		}

		namedHALogger := c.logger.Named("ha." + config.HAStorage.Type)
		c.allLoggers = append(c.allLoggers, namedHALogger)
		habackend, err := factory(config.HAStorage.Config, namedHALogger)
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

		if config.HAStorage.Type == storageTypeRaft && disableClustering {
			c.UI.Error("Disable clustering cannot be set to true when Raft is the HA storage type")
			return 1
		}

		if !disableClustering {
			coreConfig.ClusterAddr = config.HAStorage.ClusterAddr
		}
	} else {
		if coreConfig.HAPhysical, ok = backend.(physical.HABackend); ok {
			coreConfig.RedirectAddr = config.Storage.RedirectAddr
			disableClustering = config.Storage.DisableClustering

			if (config.Storage.Type == storageTypeRaft) && disableClustering {
				c.UI.Error("Disable clustering cannot be set to true when Raft is the storage type")
				return 1
			}

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
	if coreConfig.RedirectAddr == "" {
		c.logger.Warn("no `api_addr` value specified in config or in VAULT_API_ADDR; falling back to detection if possible, but this value should be manually set")
	}
	var detect physical.RedirectDetect
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		detect, ok = coreConfig.HAPhysical.(physical.RedirectDetect)
	} else {
		detect, ok = coreConfig.Physical.(physical.RedirectDetect)
	}
	if ok && coreConfig.RedirectAddr == "" {
		redirect, err := c.detectRedirect(detect, config)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error detecting api address: %s", err))
		} else if redirect == "" {
			c.UI.Error("Failed to detect api address")
		} else {
			coreConfig.RedirectAddr = redirect
		}
	}
	if coreConfig.RedirectAddr == "" && c.flagDev {
		coreConfig.RedirectAddr = fmt.Sprintf("http://%s", config.Listeners[0].Address)
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
			addrToUse = fmt.Sprintf("http://%s", config.Listeners[0].Address)
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
				c.UI.Error(fmt.Sprintf("Error parsing api address: %v", err))
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

	if coreConfig.RedirectAddr == coreConfig.ClusterAddr && len(coreConfig.RedirectAddr) != 0 {
		c.UI.Error(fmt.Sprintf(
			"Address %q used for both API and cluster addresses", coreConfig.RedirectAddr))
		return 1
	}

	if coreConfig.ClusterAddr != "" {
		// Force https as we'll always be TLS-secured
		u, err := url.ParseRequestURI(coreConfig.ClusterAddr)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing cluster address %s: %v", coreConfig.ClusterAddr, err))
			return 11
		}
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}

	// Override the UI enabling config by the environment variable
	if enableUI := os.Getenv("VAULT_UI"); enableUI != "" {
		var err error
		coreConfig.EnableUI, err = strconv.ParseBool(enableUI)
		if err != nil {
			c.UI.Output("Error parsing the environment variable VAULT_UI")
			return 1
		}
	}

	// If ServiceRegistration is configured, then the backend must support HA
	isBackendHA := coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled()
	if !c.flagDev && (coreConfig.GetServiceRegistration() != nil) && !isBackendHA {
		c.UI.Output("service_registration is configured, but storage does not support HA")
		return 1
	}

	// Apply any enterprise configuration onto the coreConfig.
	adjustCoreConfigForEnt(config, coreConfig)

	// Initialize the core
	core, newCoreError := vault.NewCore(coreConfig)
	if newCoreError != nil {
		if vault.IsFatalError(newCoreError) {
			c.UI.Error(fmt.Sprintf("Error initializing core: %s", newCoreError))
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
		info["api address"] = coreConfig.RedirectAddr
		infoKeys = append(infoKeys, "api address")
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

	status, lns, clusterAddrs, errMsg := c.InitListeners(config, disableClustering, &infoKeys, &info)

	if status != 0 {
		c.UI.Output("Error parsing listener configuration.")
		c.UI.Error(errMsg.Error())
		return 1
	}

	// Make sure we close all listeners from this point on
	listenerCloseFunc := func() {
		for _, ln := range lns {
			ln.Listener.Close()
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

	infoKeys = append(infoKeys, "recovery mode")
	info["recovery mode"] = "false"

	infoKeys = append(infoKeys, "go version")
	info["go version"] = runtime.Version()

	sort.Strings(infoKeys)
	c.UI.Output("==> Vault server configuration:\n")

	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%24s: %s",
			strings.Title(k),
			info[k]))
	}

	c.UI.Output("")

	// Tests might not want to start a vault server and just want to verify
	// the configuration.
	if c.flagTestVerifyOnly {
		return 0
	}

	// This needs to happen before we first unseal, so before we trigger dev
	// mode if it's set
	core.SetClusterListenerAddrs(clusterAddrs)
	core.SetClusterHandler(vaulthttp.Handler(&vault.HandlerProperties{
		Core: core,
	}))

	// Attempt unsealing in a background goroutine. This is needed for when a
	// Vault cluster with multiple servers is configured with auto-unseal but is
	// uninitialized. Once one server initializes the storage backend, this
	// goroutine will pick up the unseal keys and unseal this instance.
	if !core.IsInSealMigrationMode() {
		go func() {
			for {
				err := core.UnsealWithStoredKeys(context.Background())
				if err == nil {
					return
				}

				if vault.IsFatalError(err) {
					c.logger.Error("error unsealing core", "error", err)
					return
				} else {
					c.logger.Warn("failed to unseal core", "error", err)
				}

				select {
				case <-c.ShutdownCh:
					return
				case <-time.After(5 * time.Second):
				}
			}
		}()
	}

	// When the underlying storage is raft, kick off retry join if it was specified
	// in the configuration
	// TODO: Should we also support retry_join for ha_storage?
	if config.Storage.Type == storageTypeRaft {
		if err := core.InitiateRetryJoin(context.Background()); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to initiate raft retry join, %q", err.Error()))
			return 1
		}
	}

	// Perform initialization of HTTP server after the verifyOnly check.

	// Instantiate the wait group
	c.WaitGroup = &sync.WaitGroup{}

	// If service discovery is available, run service discovery
	if sd := coreConfig.GetServiceRegistration(); sd != nil {
		if err := configSR.Run(c.ShutdownCh, c.WaitGroup, coreConfig.RedirectAddr); err != nil {
			c.UI.Error(fmt.Sprintf("Error running service_registration of type %s: %s", config.ServiceRegistration.Type, err))
			return 1
		}
	}

	// If we're in Dev mode, then initialize the core
	if c.flagDev && !c.flagDevSkipInit {

		init, err := c.enableDev(core, coreConfig)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error initializing Dev mode: %s", err))
			return 1
		}

		var plugins, pluginsNotLoaded []string
		if c.flagDevPluginDir != "" && c.flagDevPluginInit {

			f, err := os.Open(c.flagDevPluginDir)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error reading plugin dir: %s", err))
				return 1
			}

			list, err := f.Readdirnames(0)
			f.Close()
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error listing plugins: %s", err))
				return 1
			}

			for _, name := range list {
				path := filepath.Join(f.Name(), name)
				if err := c.addPlugin(path, init.RootToken, core); err != nil {
					if !errwrap.Contains(err, vault.ErrPluginBadType.Error()) {
						c.UI.Error(fmt.Sprintf("Error enabling plugin %s: %s", name, err))
						return 1
					}
					pluginsNotLoaded = append(pluginsNotLoaded, name)
					continue
				}
				plugins = append(plugins, name)
			}

			sort.Strings(plugins)
		}

		var qw *quiescenceSink
		var qwo sync.Once
		qw = &quiescenceSink{
			t: time.AfterFunc(100*time.Millisecond, func() {
				qwo.Do(func() {
					c.logger.DeregisterSink(qw)

					// Print the big dev mode warning!
					c.UI.Warn(wrapAtLength(
						"WARNING! dev mode is enabled! In this mode, Vault runs entirely " +
							"in-memory and starts unsealed with a single unseal key. The root " +
							"token is already authenticated to the CLI, so you can immediately " +
							"begin using Vault."))
					c.UI.Warn("")
					c.UI.Warn("You may need to set the following environment variable:")
					c.UI.Warn("")

					endpointURL := "http://" + config.Listeners[0].Address
					if runtime.GOOS == "windows" {
						c.UI.Warn("PowerShell:")
						c.UI.Warn(fmt.Sprintf("    $env:VAULT_ADDR=\"%s\"", endpointURL))
						c.UI.Warn("cmd.exe:")
						c.UI.Warn(fmt.Sprintf("    set VAULT_ADDR=%s", endpointURL))
					} else {
						c.UI.Warn(fmt.Sprintf("    $ export VAULT_ADDR='%s'", endpointURL))
					}

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
						c.UI.Warn(fmt.Sprintf("Recovery Key: %s", base64.StdEncoding.EncodeToString(init.RecoveryShares[0])))
					}

					c.UI.Warn(fmt.Sprintf("Root Token: %s", init.RootToken))

					if len(plugins) > 0 {
						c.UI.Warn("")
						c.UI.Warn(wrapAtLength(
							"The following dev plugins are registered in the catalog:"))
						for _, p := range plugins {
							c.UI.Warn(fmt.Sprintf("    - %s", p))
						}
					}

					if len(pluginsNotLoaded) > 0 {
						c.UI.Warn("")
						c.UI.Warn(wrapAtLength(
							"The following dev plugins FAILED to be registered in the catalog due to unknown type:"))
						for _, p := range pluginsNotLoaded {
							c.UI.Warn(fmt.Sprintf("    - %s", p))
						}
					}

					c.UI.Warn("")
					c.UI.Warn(wrapAtLength(
						"Development mode should NOT be used in production installations!"))
					c.UI.Warn("")
				})
			}),
		}
		c.logger.RegisterSink(qw)
	}

	// Initialize the HTTP servers
	for _, ln := range lns {
		if ln.Config == nil {
			c.UI.Error("Found nil listener config after parsing")
			return 1
		}
		handler := vaulthttp.Handler(&vault.HandlerProperties{
			Core:                  core,
			ListenerConfig:        ln.Config,
			DisablePrintableCheck: config.DisablePrintableCheck,
			RecoveryMode:          c.flagRecovery,
		})

		if len(ln.Config.XForwardedForAuthorizedAddrs) > 0 {
			handler = vaulthttp.WrapForwardedForHandler(handler, ln.Config)
		}

		// server defaults
		server := &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			IdleTimeout:       5 * time.Minute,
			ErrorLog:          c.logger.StandardLogger(nil),
		}

		// override server defaults with config values for read/write/idle timeouts if configured
		if ln.Config.HTTPReadHeaderTimeout > 0 {
			server.ReadHeaderTimeout = ln.Config.HTTPReadHeaderTimeout
		}
		if ln.Config.HTTPReadTimeout > 0 {
			server.ReadTimeout = ln.Config.HTTPReadTimeout
		}
		if ln.Config.HTTPWriteTimeout > 0 {
			server.WriteTimeout = ln.Config.HTTPWriteTimeout
		}
		if ln.Config.HTTPIdleTimeout > 0 {
			server.IdleTimeout = ln.Config.HTTPIdleTimeout
		}

		// server config tests can exit now
		if c.flagTestServerConfig {
			continue
		}

		go server.Serve(ln.Listener)
	}

	if c.flagTestServerConfig {
		return 0
	}

	if sealConfigError != nil {
		init, err := core.InitializedLocally(context.Background())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking if core is initialized: %v", err))
			return 1
		}
		if init {
			c.UI.Error("Vault is initialized but no Seal key could be loaded")
			return 1
		}
	}

	if newCoreError != nil {
		c.UI.Warn(wrapAtLength(
			"WARNING! A non-fatal error occurred during initialization. Please " +
				"check the logs for more information."))
		c.UI.Warn("")
	}

	// Output the header that the server has started
	if !c.flagCombineLogs {
		c.UI.Output("==> Vault server started! Log data will stream in below:\n")
	}

	// Inform any tests that the server is ready
	select {
	case c.startedCh <- struct{}{}:
	default:
	}

	// Release the log gate.
	c.flushLog()

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

	var coreShutdownDoneCh <-chan struct{}
	if c.flagExitOnCoreShutdown {
		coreShutdownDoneCh = core.ShutdownDone()
	}

	// Wait for shutdown
	shutdownTriggered := false
	retCode := 0

	for !shutdownTriggered {
		select {
		case <-coreShutdownDoneCh:
			c.UI.Output("==> Vault core was shut down")
			retCode = 1
			shutdownTriggered = true
		case <-c.ShutdownCh:
			c.UI.Output("==> Vault shutdown triggered")
			shutdownTriggered = true
		case <-c.SighupCh:
			c.UI.Output("==> Vault reload triggered")

			// Check for new log level
			var config *server.Config
			var level log.Level
			for _, path := range c.flagConfigs {
				current, err := server.LoadConfig(path)
				if err != nil {
					c.logger.Error("could not reload config", "path", path, "error", err)
					goto RUNRELOADFUNCS
				}

				if config == nil {
					config = current
				} else {
					config = config.Merge(current)
				}
			}

			// Ensure at least one config was found.
			if config == nil {
				c.logger.Error("no config found at reload time")
				goto RUNRELOADFUNCS
			}

			core.SetConfig(config)

			if config.LogLevel != "" {
				configLogLevel := strings.ToLower(strings.TrimSpace(config.LogLevel))
				switch configLogLevel {
				case "trace":
					level = log.Trace
				case "debug":
					level = log.Debug
				case "notice", "info", "":
					level = log.Info
				case "warn", "warning":
					level = log.Warn
				case "err", "error":
					level = log.Error
				default:
					c.logger.Error("unknown log level found on reload", "level", config.LogLevel)
					goto RUNRELOADFUNCS
				}
				core.SetLogLevel(level)
			}

		RUNRELOADFUNCS:
			if err := c.Reload(c.reloadFuncsLock, c.reloadFuncs, c.flagConfigs); err != nil {
				c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
			}

		case <-c.SigUSR2Ch:
			logWriter := c.logger.StandardWriter(&hclog.StandardLoggerOptions{})
			pprof.Lookup("goroutine").WriteTo(logWriter, 2)
		}
	}

	// Stop the listeners so that we don't process further client requests.
	c.cleanupGuard.Do(listenerCloseFunc)

	// Shutdown will wait until after Vault is sealed, which means the
	// request forwarding listeners will also be closed (and also
	// waited for).
	if err := core.Shutdown(); err != nil {
		c.UI.Error(fmt.Sprintf("Error with core shutdown: %s", err))
	}

	// Wait for dependent goroutines to complete
	c.WaitGroup.Wait()
	return retCode
}

func (c *ServerCommand) enableDev(core *vault.Core, coreConfig *vault.CoreConfig) (*vault.InitResult, error) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

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

	if core.SealAccess().StoredKeysSupported() != vaultseal.StoredKeysNotSupported {
		barrierConfig.StoredShares = 1
	}

	// Initialize it with a basic single key
	init, err := core.Initialize(ctx, &vault.InitParams{
		BarrierConfig:  barrierConfig,
		RecoveryConfig: recoveryConfig,
	})
	if err != nil {
		return nil, err
	}

	// Handle unseal with stored keys
	if core.SealAccess().StoredKeysSupported() == vaultseal.StoredKeysSupportedGeneric {
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
		return nil, errwrap.Wrapf("failed to check active status: {{err}}", err)
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
				return nil, errwrap.Wrapf("failed to check active status: {{err}}", err)
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
		resp, err := core.HandleRequest(ctx, req)
		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("failed to create root token with ID %q: {{err}}", coreConfig.DevToken), err)
		}
		if resp == nil {
			return nil, fmt.Errorf("nil response when creating root token with ID %q", coreConfig.DevToken)
		}
		if resp.Auth == nil {
			return nil, fmt.Errorf("nil auth when creating root token with ID %q", coreConfig.DevToken)
		}

		init.RootToken = resp.Auth.ClientToken

		req.ID = "dev-revoke-init-root"
		req.Path = "auth/token/revoke-self"
		req.Data = nil
		resp, err = core.HandleRequest(ctx, req)
		if err != nil {
			return nil, errwrap.Wrapf("failed to revoke initial root token: {{err}}", err)
		}
	}

	// Set the token
	if !c.flagDevNoStoreToken {
		tokenHelper, err := c.TokenHelper()
		if err != nil {
			return nil, err
		}
		if err := tokenHelper.Store(init.RootToken); err != nil {
			return nil, err
		}
	}

	kvVer := "2"
	if c.flagDevKVV1 || c.flagDevLeasedKV {
		kvVer = "1"
	}
	req := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: init.RootToken,
		Path:        "sys/mounts/secret",
		Data: map[string]interface{}{
			"type":        "kv",
			"path":        "secret/",
			"description": "key/value secret storage",
			"options": map[string]string{
				"version": kvVer,
			},
		},
	}
	resp, err := core.HandleRequest(ctx, req)
	if err != nil {
		return nil, errwrap.Wrapf("error creating default K/V store: {{err}}", err)
	}
	if resp.IsError() {
		return nil, errwrap.Wrapf("failed to create default K/V store: {{err}}", resp.Error())
	}

	return init, nil
}

func (c *ServerCommand) enableThreeNodeDevCluster(base *vault.CoreConfig, info map[string]string, infoKeys []string, devListenAddress, tempDir string) int {
	testCluster := vault.NewTestCluster(&testing.RuntimeT{}, base, &vault.TestClusterOptions{
		HandlerFunc:       vaulthttp.Handler,
		BaseListenAddress: c.flagDevListenAddr,
		Logger:            c.logger,
		TempDir:           tempDir,
	})
	defer c.cleanupGuard.Do(testCluster.Cleanup)

	info["cluster parameters path"] = testCluster.TempDir
	infoKeys = append(infoKeys, "cluster parameters path")

	for i, core := range testCluster.Cores {
		info[fmt.Sprintf("node %d api address", i)] = fmt.Sprintf("https://%s", core.Listeners[0].Address.String())
		infoKeys = append(infoKeys, fmt.Sprintf("node %d api address", i))
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

	infoKeys = append(infoKeys, "go version")
	info["go version"] = runtime.Version()

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
		core.Server.Handler = vaulthttp.Handler(&vault.HandlerProperties{
			Core: core.Core,
		})
		core.SetClusterHandler(core.Server.Handler)
	}

	testCluster.Start()

	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

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
		resp, err := testCluster.Cores[0].HandleRequest(ctx, req)
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
		resp, err = testCluster.Cores[0].HandleRequest(ctx, req)
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

	if err := ioutil.WriteFile(filepath.Join(testCluster.TempDir, "root_token"), []byte(testCluster.RootToken), 0o755); err != nil {
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
	c.flushLog()

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

// addPlugin adds any plugins to the catalog
func (c *ServerCommand) addPlugin(path, token string, core *vault.Core) error {
	// Get the sha256 of the file at the given path.
	pluginSum := func(p string) (string, error) {
		hasher := sha256.New()
		f, err := os.Open(p)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := io.Copy(hasher, f); err != nil {
			return "", err
		}
		return hex.EncodeToString(hasher.Sum(nil)), nil
	}

	// Mount any test plugins. We do this explicitly before we inform tests of
	// a completely booted server intentionally.
	sha256sum, err := pluginSum(path)
	if err != nil {
		return err
	}

	// Default the name to the basename of the binary
	name := filepath.Base(path)

	// File a request against core to enable the plugin
	req := &logical.Request{
		Operation:   logical.UpdateOperation,
		ClientToken: token,
		Path:        fmt.Sprintf("sys/plugins/catalog/%s", name),
		Data: map[string]interface{}{
			"sha256":  sha256sum,
			"command": name,
		},
	}
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	if _, err := core.HandleRequest(ctx, req); err != nil {
		return err
	}

	return nil
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
		if list.TLSDisable {
			scheme = "http"
		}

		// Check for address override
		addr := list.Address
		if addr == "" {
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

func (c *ServerCommand) Reload(lock *sync.RWMutex, reloadFuncs *map[string][]reloadutil.ReloadFunc, configPath []string) error {
	lock.RLock()
	defer lock.RUnlock()

	var reloadErrors *multierror.Error

	for k, relFuncs := range *reloadFuncs {
		switch {
		case strings.HasPrefix(k, "listener|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(); err != nil {
						reloadErrors = multierror.Append(reloadErrors, errwrap.Wrapf("error encountered reloading listener: {{err}}", err))
					}
				}
			}

		case strings.HasPrefix(k, "audit_file|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(); err != nil {
						reloadErrors = multierror.Append(reloadErrors, errwrap.Wrapf(fmt.Sprintf("error encountered reloading file audit device at path %q: {{err}}", strings.TrimPrefix(k, "audit_file|")), err))
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
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return errwrap.Wrapf("could not open pid file: {{err}}", err)
	}
	defer pidFile.Close()

	// Write out the PID
	pid := os.Getpid()
	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return errwrap.Wrapf("could not write to pid file: {{err}}", err)
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

// storageMigrationActive checks and warns against in-progress storage migrations.
// This function will block until storage is available.
func (c *ServerCommand) storageMigrationActive(backend physical.Backend) bool {
	first := true

	for {
		migrationStatus, err := CheckStorageMigration(backend)
		if err == nil {
			if migrationStatus != nil {
				startTime := migrationStatus.Start.Format(time.RFC3339)
				c.UI.Error(wrapAtLength(fmt.Sprintf("ERROR! Storage migration in progress (started: %s). "+
					"Server startup is prevented until the migration completes. Use 'vault operator migrate -reset' "+
					"to force clear the migration lock.", startTime)))
				return true
			}
			return false
		}
		if first {
			first = false
			c.UI.Warn("\nWARNING! Unable to read storage migration status.")

			// unexpected state, so stop buffering log messages
			c.flushLog()
		}
		c.logger.Warn("storage migration check error", "error", err.Error())

		select {
		case <-time.After(2 * time.Second):
		case <-c.ShutdownCh:
			return true
		}
	}
}

type StorageMigrationStatus struct {
	Start time.Time `json:"start"`
}

func CheckStorageMigration(b physical.Backend) (*StorageMigrationStatus, error) {
	entry, err := b.Get(context.Background(), storageMigrationLock)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var status StorageMigrationStatus
	if err := jsonutil.DecodeJSON(entry.Value, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func SetStorageMigration(b physical.Backend, active bool) error {
	if !active {
		return b.Delete(context.Background(), storageMigrationLock)
	}

	status := StorageMigrationStatus{
		Start: time.Now(),
	}

	enc, err := jsonutil.EncodeJSON(status)
	if err != nil {
		return err
	}

	entry := &physical.Entry{
		Key:   storageMigrationLock,
		Value: enc,
	}

	return b.Put(context.Background(), entry)
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
	if g.log && g.logger.IsDebug() {
		g.logger.Debug(fmt.Sprint(args...))
	}
}

func (g *grpclogFaker) Printf(format string, args ...interface{}) {
	if g.log && g.logger.IsDebug() {
		g.logger.Debug(fmt.Sprintf(format, args...))
	}
}

func (g *grpclogFaker) Println(args ...interface{}) {
	if g.log && g.logger.IsDebug() {
		g.logger.Debug(fmt.Sprintln(args...))
	}
}
