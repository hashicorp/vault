// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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

	systemd "github.com/coreos/go-systemd/daemon"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/cli"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-kms-wrapping/entropy/v2"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/gatedwriter"
	"github.com/hashicorp/go-secure-stdlib/mlock"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/vault/audit"
	config2 "github.com/hashicorp/vault/command/config"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/experiments"
	loghelper "github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/useragent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/hcp_link"
	"github.com/hashicorp/vault/vault/plugincatalog"
	vaultseal "github.com/hashicorp/vault/vault/seal"
	"github.com/hashicorp/vault/version"
	"github.com/posener/complete"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/atomic"
	"golang.org/x/net/http/httpproxy"
	"google.golang.org/grpc/grpclog"
)

var (
	_ cli.Command             = (*ServerCommand)(nil)
	_ cli.CommandAutocomplete = (*ServerCommand)(nil)
)

var memProfilerEnabled = false

const (
	storageMigrationLock = "core/migration"

	// Even though there are more types than the ones below, the following consts
	// are declared internally for value comparison and reusability.
	storageTypeRaft   = "raft"
	storageTypeConsul = "consul"
)

type ServerCommand struct {
	*BaseCommand
	logFlags logFlags

	AuditBackends      map[string]audit.Factory
	CredentialBackends map[string]logical.Factory
	LogicalBackends    map[string]logical.Factory
	PhysicalBackends   map[string]physical.Factory

	ServiceRegistrations map[string]sr.Factory

	ShutdownCh chan struct{}
	SighupCh   chan struct{}
	SigUSR2Ch  chan struct{}

	WaitGroup *sync.WaitGroup

	logWriter io.Writer
	logGate   *gatedwriter.Writer
	logger    hclog.InterceptLogger

	cleanupGuard sync.Once

	reloadFuncsLock   *sync.RWMutex
	reloadFuncs       *map[string][]reloadutil.ReloadFunc
	startedCh         chan (struct{}) // for tests
	reloadedCh        chan (struct{}) // for tests
	licenseReloadedCh chan (error)    // for tests

	allLoggers []hclog.Logger

	flagConfigs            []string
	flagRecovery           bool
	flagExperiments        []string
	flagCLIDump            string
	flagDev                bool
	flagDevTLS             bool
	flagDevTLSCertDir      string
	flagDevTLSSANs         []string
	flagDevRootTokenID     string
	flagDevListenAddr      string
	flagDevNoStoreToken    bool
	flagDevPluginDir       string
	flagDevPluginInit      bool
	flagDevHA              bool
	flagDevLatency         int
	flagDevLatencyJitter   int
	flagDevLeasedKV        bool
	flagDevNoKV            bool
	flagDevKVV1            bool
	flagDevSkipInit        bool
	flagDevTransactional   bool
	flagDevAutoSeal        bool
	flagDevClusterJson     string
	flagTestVerifyOnly     bool
	flagTestServerConfig   bool
	flagDevConsul          bool
	flagExitOnCoreShutdown bool

	sealsToFinalize []*vault.Seal
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

	// Augment with the log flags
	f.addLogFlags(&c.logFlags)

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

	f.BoolVar(&BoolVar{
		Name:    "exit-on-core-shutdown",
		Target:  &c.flagExitOnCoreShutdown,
		Default: false,
		Usage:   "Exit the vault server if the vault core is shutdown.",
	})

	f.BoolVar(&BoolVar{
		Name:   "recovery",
		Target: &c.flagRecovery,
		Usage: "Enable recovery mode. In this mode, Vault is used to perform recovery actions. " +
			"Using a recovery token, \"sys/raw\" API can be used to manipulate the storage.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:       "experiment",
		Target:     &c.flagExperiments,
		Completion: complete.PredictSet(experiments.ValidExperiments()...),
		Usage: "Name of an experiment to enable. Experiments should NOT be used in production, and " +
			"the associated APIs may have backwards incompatible changes between releases. This " +
			"flag can be specified multiple times to specify multiple experiments. This can also be " +
			fmt.Sprintf("specified via the %s environment variable as a comma-separated list. ", EnvVaultExperiments) +
			"Valid experiments are: " + strings.Join(experiments.ValidExperiments(), ", "),
	})

	f.StringVar(&StringVar{
		Name:       "pprof-dump-dir",
		Target:     &c.flagCLIDump,
		Completion: complete.PredictDirs("*"),
		Usage:      "Directory where generated profiles are created. If left unset, files are not generated.",
	})

	f = set.NewFlagSet("Dev Options")

	f.BoolVar(&BoolVar{
		Name:   "dev",
		Target: &c.flagDev,
		Usage: "Enable development mode. In this mode, Vault runs in-memory and " +
			"starts unsealed. As the name implies, do not run \"dev\" mode in " +
			"production.",
	})

	f.BoolVar(&BoolVar{
		Name:   "dev-tls",
		Target: &c.flagDevTLS,
		Usage: "Enable TLS development mode. In this mode, Vault runs in-memory and " +
			"starts unsealed, with a generated TLS CA, certificate and key. " +
			"As the name implies, do not run \"dev-tls\" mode in " +
			"production.",
	})

	f.StringVar(&StringVar{
		Name:    "dev-tls-cert-dir",
		Target:  &c.flagDevTLSCertDir,
		Default: "",
		Usage: "Directory where generated TLS files are created if `-dev-tls` is " +
			"specified. If left unset, files are generated in a temporary directory.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:    "dev-tls-san",
		Target:  &c.flagDevTLSSANs,
		Default: nil,
		Usage: "Additional Subject Alternative Name (as a DNS name or IP address) " +
			"to generate the certificate with if `-dev-tls` is specified. The " +
			"certificate will always use localhost, localhost4, localhost6, " +
			"localhost.localdomain, and the host name as alternate DNS names, " +
			"and 127.0.0.1 as an alternate IP address. This flag can be specified " +
			"multiple times to specify multiple SANs.",
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
		Name:    "dev-no-kv",
		Target:  &c.flagDevNoKV,
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
		Name:    "dev-consul",
		Target:  &c.flagDevConsul,
		Default: false,
		Hidden:  true,
	})

	f.StringVar(&StringVar{
		Name:   "dev-cluster-json",
		Target: &c.flagDevClusterJson,
		Usage:  "File to write cluster definition to",
	})

	// TODO: should the below flags be public?
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
		Output: c.logWriter,
	}, c.logGate)
}

func (c *ServerCommand) parseConfig() (*server.Config, []configutil.ConfigError, error) {
	var configErrors []configutil.ConfigError
	// Load the configuration
	var config *server.Config
	for _, path := range c.flagConfigs {
		current, err := server.LoadConfig(path)
		if err != nil {
			return nil, nil, fmt.Errorf("error loading configuration from %s: %w", path, err)
		}

		configErrors = append(configErrors, current.Validate(path)...)

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}

	if config != nil && config.Entropy != nil && config.Entropy.Mode == configutil.EntropyAugmentation && constants.IsFIPS() {
		c.UI.Warn("WARNING: Entropy Augmentation is not supported in FIPS 140-2 Inside mode; disabling from server configuration!\n")
		config.Entropy = nil
	}

	entCheckRequestLimiter(c, config)

	return config, configErrors, nil
}

func (c *ServerCommand) runRecoveryMode() int {
	config, configErrors, err := c.parseConfig()
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

	// Update the 'log' related aspects of shared config based on config/env var/cli
	c.flags.applyLogConfigOverrides(config.SharedConfig)
	l, err := c.configureLogging(config)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.logger = l
	c.allLoggers = append(c.allLoggers, l)

	// reporting Errors found in the config
	for _, cErr := range configErrors {
		c.logger.Warn(cErr.String())
	}

	// Ensure logging is flushed if initialization fails
	defer c.flushLog()

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
			config.ClusterAddr = configutil.NormalizeAddr(envCA)
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
	info["log level"] = config.LogLevel
	infoKeys = append(infoKeys, "log level")

	var barrierSeal vault.Seal
	var sealConfigError error

	if len(config.Seals) == 0 {
		config.Seals = append(config.Seals, &configutil.KMS{Type: wrapping.WrapperTypeShamir.String()})
	}

	ctx := context.Background()
	existingSealGenerationInfo, err := vault.PhysicalSealGenInfo(ctx, backend)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting seal generation info: %v", err))
		return 1
	}

	hasPartialPaths, err := hasPartiallyWrappedPaths(ctx, backend)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Cannot determine if there are partially seal wrapped entries in storage: %v", err))
		return 1
	}
	setSealResponse, err := setSeal(c, config, infoKeys, info, existingSealGenerationInfo, hasPartialPaths)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	if setSealResponse.barrierSeal == nil {
		c.UI.Error(fmt.Sprintf("Error setting up seal: %v", setSealResponse.sealConfigError))
		return 1
	}
	barrierSeal = setSealResponse.barrierSeal

	// Ensure that the seal finalizer is called, even if using verify-only
	defer func() {
		err = barrierSeal.Finalize(ctx)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
		}
	}()

	coreConfig := &vault.CoreConfig{
		Physical:     backend,
		StorageType:  config.Storage.Type,
		Seal:         barrierSeal,
		UnwrapSeal:   setSealResponse.unwrapSeal,
		LogLevel:     config.LogLevel,
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

	if err := core.InitializeRecovery(ctx); err != nil {
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
		ln, _, _, err := server.NewListener(lnConfig, c.logGate, c.UI)
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

	fipsStatus := entGetFIPSInfoKey()
	if fipsStatus != "" {
		infoKeys = append(infoKeys, "fips")
		info["fips"] = fipsStatus
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

	for _, ln := range lns {
		handler := vaulthttp.Handler.Handler(&vault.HandlerProperties{
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
		init, err := core.InitializedLocally(ctx)
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

	if !c.logFlags.flagCombineLogs {
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
}

func logProxyEnvironmentVariables(logger hclog.Logger) {
	proxyCfg := httpproxy.FromEnvironment()
	cfgMap := map[string]string{
		"http_proxy":  configutil.NormalizeAddr(proxyCfg.HTTPProxy),
		"https_proxy": configutil.NormalizeAddr(proxyCfg.HTTPSProxy),
		"no_proxy":    configutil.NormalizeAddr(proxyCfg.NoProxy),
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

type quiescenceSink struct {
	t *time.Timer
}

func (q quiescenceSink) Accept(name string, level hclog.Level, msg string, args ...interface{}) {
	q.t.Reset(100 * time.Millisecond)
}

func (c *ServerCommand) setupStorage(config *server.Config) (physical.Backend, error) {
	// Ensure that a backend is provided
	if config.Storage == nil {
		return nil, errors.New("A storage backend must be specified")
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
			return nil, errors.New("Cluster address must be set when using raft storage")
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

func beginServiceRegistration(c *ServerCommand, config *server.Config) (sr.ServiceRegistration, error) {
	sdFactory, ok := c.ServiceRegistrations[config.ServiceRegistration.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown service_registration type %s", config.ServiceRegistration.Type)
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
	var err error
	configSR, err := sdFactory(config.ServiceRegistration.Config, namedSDLogger, state)
	if err != nil {
		return nil, fmt.Errorf("Error initializing service_registration of type %s: %s", config.ServiceRegistration.Type, err)
	}

	return configSR, nil
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
		ln, props, reloadFunc, err := server.NewListener(lnConfig, c.logGate, c.UI)
		if err != nil {
			errMsg = fmt.Errorf("Error initializing listener of type %s: %s", lnConfig.Type, err)
			return 1, nil, nil, errMsg
		}

		if reloadFunc != nil {
			relSlice := (*c.reloadFuncs)[fmt.Sprintf("listener|%s", lnConfig.Type)]
			relSlice = append(relSlice, reloadFunc)
			(*c.reloadFuncs)[fmt.Sprintf("listener|%s", lnConfig.Type)] = relSlice
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

		props["disable_request_limiter"] = strconv.FormatBool(lnConfig.DisableRequestLimiter)

		if lnConfig.ChrootNamespace != "" {
			props["chroot_namespace"] = lnConfig.ChrootNamespace
		}

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

func configureDevTLS(c *ServerCommand) (func(), *server.Config, string, error) {
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

	var certDir string
	var err error
	var config *server.Config
	var f func()

	if c.flagDevTLS {
		if c.flagDevTLSCertDir != "" {
			if _, err = os.Stat(c.flagDevTLSCertDir); err != nil {
				return nil, nil, "", err
			}

			certDir = c.flagDevTLSCertDir
		} else {
			if certDir, err = os.MkdirTemp("", "vault-tls"); err != nil {
				return nil, nil, certDir, err
			}
		}
		extraSANs := c.flagDevTLSSANs
		host, _, err := net.SplitHostPort(c.flagDevListenAddr)
		if err == nil {
			// 127.0.0.1 is the default, and already included in the SANs.
			// Empty host means listen on all interfaces, but users should use the
			// -dev-tls-san flag to get the right SANs in that case.
			if host != "" && host != "127.0.0.1" {
				extraSANs = append(extraSANs, host)
			}
		}
		config, err = server.DevTLSConfig(devStorageType, certDir, extraSANs)

		f = func() {
			if err := os.Remove(fmt.Sprintf("%s/%s", certDir, server.VaultDevCAFilename)); err != nil {
				c.UI.Error(err.Error())
			}

			if err := os.Remove(fmt.Sprintf("%s/%s", certDir, server.VaultDevCertFilename)); err != nil {
				c.UI.Error(err.Error())
			}

			if err := os.Remove(fmt.Sprintf("%s/%s", certDir, server.VaultDevKeyFilename)); err != nil {
				c.UI.Error(err.Error())
			}

			// Only delete temp directories we made.
			if c.flagDevTLSCertDir == "" {
				if err := os.Remove(certDir); err != nil {
					c.UI.Error(err.Error())
				}
			}
		}

	} else {
		config, err = server.DevConfig(devStorageType)
	}

	return f, config, certDir, err
}

func (c *ServerCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Don't exit just because we saw a potential deadlock.
	deadlock.Opts.OnPotentialDeadlock = func() {}

	c.logGate = gatedwriter.NewWriter(os.Stderr)
	c.logWriter = c.logGate

	if c.logFlags.flagCombineLogs {
		c.logWriter = os.Stdout
	}

	if c.flagRecovery {
		return c.runRecoveryMode()
	}

	// Automatically enable dev mode if other dev flags are provided.
	if c.flagDevConsul || c.flagDevHA || c.flagDevTransactional || c.flagDevLeasedKV || c.flagDevAutoSeal || c.flagDevKVV1 || c.flagDevNoKV || c.flagDevTLS {
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
	var certDir string
	if c.flagDev {
		df, cfg, dir, err := configureDevTLS(c)
		if df != nil {
			defer df()
		}

		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		config = cfg
		certDir = dir

		if c.flagDevListenAddr != "" {
			config.Listeners[0].Address = c.flagDevListenAddr
		}
		config.Listeners[0].Telemetry.UnauthenticatedMetricsAccess = true
	}

	parsedConfig, configErrors, err := c.parseConfig()
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

	f.applyLogConfigOverrides(config.SharedConfig)

	l, err := c.configureLogging(config)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.logger = l
	c.allLoggers = append(c.allLoggers, l)

	// flush logs right away if the server is started with the disable-gated-logs flag
	if c.logFlags.flagDisableGatedLogs {
		c.flushLog()
	}

	// reporting Errors found in the config
	for _, cErr := range configErrors {
		c.logger.Warn(cErr.String())
	}

	// Ensure logging is flushed if initialization fails
	defer c.flushLog()

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

	if envLicensePath := os.Getenv(EnvVaultLicensePath); envLicensePath != "" {
		config.LicensePath = envLicensePath
	}
	if envLicense := os.Getenv(EnvVaultLicense); envLicense != "" {
		config.License = envLicense
	}

	if envPluginTmpdir := os.Getenv(EnvVaultPluginTmpdir); envPluginTmpdir != "" {
		config.PluginTmpdir = envPluginTmpdir
	}

	if err := server.ExperimentsFromEnvAndCLI(config, EnvVaultExperiments, c.flagExperiments); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	for _, experiment := range config.Experiments {
		if experiments.IsUnused(experiment) {
			c.UI.Warn(fmt.Sprintf("WARNING! Experiment %s is no longer used", experiment))
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

	// Initialize the storage backend
	var backend physical.Backend
	if !c.flagDev || config.Storage != nil {
		backend, err = c.setupStorage(config)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		// Prevent server startup if migration is active
		// TODO: Use OpenTelemetry to integrate this into Diagnose
		if c.storageMigrationActive(backend) {
			return 1
		}
	}

	clusterName := config.ClusterName

	// Attempt to retrieve cluster name from insecure storage
	if clusterName == "" {
		clusterName, err = c.readClusterNameFromInsecureStorage(backend)
	}

	inmemMetrics, metricSink, prometheusEnabled, err := configutil.SetupTelemetry(&configutil.SetupTelemetryOpts{
		Config:      config.Telemetry,
		Ui:          c.UI,
		ServiceName: "vault",
		DisplayName: "Vault",
		UserAgent:   useragent.String(),
		ClusterName: clusterName,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}
	metricsHelper := metricsutil.NewMetricsHelper(inmemMetrics, prometheusEnabled)

	// Initialize the Service Discovery, if there is one
	var configSR sr.ServiceRegistration
	if config.ServiceRegistration != nil {
		configSR, err = beginServiceRegistration(c, config)
		if err != nil {
			c.UI.Output(err.Error())
			return 1
		}
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = config.LogLevel
	infoKeys = append(infoKeys, "log level")

	// returns a slice of env vars formatted as "key=value"
	envVars := os.Environ()
	var envVarKeys []string
	for _, v := range envVars {
		splitEnvVars := strings.Split(v, "=")
		envVarKeys = append(envVarKeys, splitEnvVars[0])
	}

	sort.Strings(envVarKeys)

	key := "environment variables"
	info[key] = strings.Join(envVarKeys, ", ")
	infoKeys = append(infoKeys, key)

	if len(config.Experiments) != 0 {
		expKey := "experiments"
		info[expKey] = strings.Join(config.Experiments, ", ")
		infoKeys = append(infoKeys, expKey)
	}

	ctx := context.Background()

	setSealResponse, secureRandomReader, err := c.configureSeals(ctx, config, backend, infoKeys, info)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.setSealsToFinalize(setSealResponse.getCreatedSeals())
	defer func() {
		c.finalizeSeals(ctx, c.sealsToFinalize)
	}()

	coreConfig := createCoreConfig(c, config, backend, configSR, setSealResponse.barrierSeal, setSealResponse.unwrapSeal, metricsHelper, metricSink, secureRandomReader)

	if allowPendingRemoval := os.Getenv(consts.EnvVaultAllowPendingRemovalMounts); allowPendingRemoval != "" {
		var err error
		coreConfig.PendingRemovalMountsAllowed, err = strconv.ParseBool(allowPendingRemoval)
		if err != nil {
			c.UI.Warn(wrapAtLength("WARNING! failed to parse " +
				consts.EnvVaultAllowPendingRemovalMounts + " env var: " +
				"defaulting to false."))
		}
	}

	// Initialize the separate HA storage backend, if it exists
	disableClustering, err := initHaBackend(c, config, &coreConfig, backend)
	if err != nil {
		c.UI.Output(err.Error())
		return 1
	}

	// Determine the redirect address from environment variables
	err = determineRedirectAddr(c, &coreConfig, config)
	if err != nil {
		c.UI.Output(err.Error())
	}

	// After the redirect bits are sorted out, if no cluster address was
	// explicitly given, derive one from the redirect addr
	err = findClusterAddress(c, &coreConfig, config, disableClustering)
	if err != nil {
		c.UI.Output(err.Error())
		return 1
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
	entAdjustCoreConfig(config, &coreConfig)

	if !entCheckStorageType(&coreConfig) {
		c.UI.Warn("")
		c.UI.Warn(wrapAtLength(fmt.Sprintf("WARNING: storage configured to use %q which is not supported for Vault Enterprise, must be \"raft\" or \"consul\"", coreConfig.StorageType)))
		c.UI.Warn("")
	}

	if !c.flagDev {
		inMemStorageTypes := []string{
			"inmem", "inmem_ha", "inmem_transactional", "inmem_transactional_ha",
		}

		if strutil.StrListContains(inMemStorageTypes, coreConfig.StorageType) {
			c.UI.Warn("")
			c.UI.Warn(wrapAtLength(fmt.Sprintf("WARNING: storage configured to use %q which should NOT be used in production", coreConfig.StorageType)))
			c.UI.Warn("")
		}
	}

	// Initialize the core
	core, newCoreError := vault.NewCore(&coreConfig)
	if newCoreError != nil {
		if vault.IsFatalError(newCoreError) {
			c.UI.Error(fmt.Sprintf("Error initializing core: %s", newCoreError))
			return 1
		}
		c.UI.Warn(wrapAtLength(
			"WARNING! A non-fatal error occurred during initialization. Please " +
				"check the logs for more information."))
		c.UI.Warn("")

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

	fipsStatus := entGetFIPSInfoKey()
	if fipsStatus != "" {
		infoKeys = append(infoKeys, "fips")
		info["fips"] = fipsStatus
	}

	if config.HCPLinkConf != nil {
		infoKeys = append(infoKeys, "HCP organization")
		info["HCP organization"] = config.HCPLinkConf.Resource.Organization

		infoKeys = append(infoKeys, "HCP project")
		info["HCP project"] = config.HCPLinkConf.Resource.Project

		infoKeys = append(infoKeys, "HCP resource ID")
		info["HCP resource ID"] = config.HCPLinkConf.Resource.ID
	}

	infoKeys = append(infoKeys, "administrative namespace")
	info["administrative namespace"] = config.AdministrativeNamespacePath

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
	core.SetClusterHandler(vaulthttp.Handler.Handler(&vault.HandlerProperties{
		Core:           core,
		ListenerConfig: &configutil.Listener{},
	}))

	// Attempt unsealing in a background goroutine. This is needed for when a
	// Vault cluster with multiple servers is configured with auto-unseal but is
	// uninitialized. Once one server initializes the storage backend, this
	// goroutine will pick up the unseal keys and unseal this instance.
	if !core.IsInSealMigrationMode(true) {
		go runUnseal(c, core, ctx)
	}

	// When the underlying storage is raft, kick off retry join if it was specified
	// in the configuration
	// TODO: Should we also support retry_join for ha_storage?
	if config.Storage.Type == storageTypeRaft {
		if err := core.InitiateRetryJoin(ctx); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to initiate raft retry join, %q", err.Error()))
			return 1
		}
	}

	// Perform initialization of HTTP server after the verifyOnly check.

	// Instantiate the wait group
	c.WaitGroup = &sync.WaitGroup{}

	// If service discovery is available, run service discovery
	err = runListeners(c, &coreConfig, config, configSR)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// If we're in Dev mode, then initialize the core
	clusterJson := &testcluster.ClusterJson{}
	err = initDevCore(c, &coreConfig, config, core, certDir, clusterJson)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Initialize the HTTP servers
	err = startHttpServers(c, core, config, lns)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	hcpLogger := c.logger.Named("hcp-connectivity")
	hcpLink, err := hcp_link.NewHCPLink(config.HCPLinkConf, core, hcpLogger)
	if err != nil {
		c.logger.Error("failed to establish HCP connection", "error", err)
	} else if hcpLink != nil {
		c.logger.Trace("established HCP connection")
	}

	if c.flagTestServerConfig {
		return 0
	}

	if setSealResponse.sealConfigError != nil {
		init, err := core.InitializedLocally(ctx)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error checking if core is initialized: %v", err))
			return 1
		}
		if init {
			c.UI.Error("Vault is initialized but no Seal key could be loaded")
			return 1
		}
	}

	core.SetSealReloadFunc(func(ctx context.Context) error {
		return c.reloadSealsOnLeaderActivation(ctx, core)
	})

	// Output the header that the server has started
	if !c.logFlags.flagCombineLogs {
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

	// Notify systemd that the server is ready (if applicable)
	c.notifySystemd(systemd.SdNotifyReady)

	if c.flagDev {
		protocol := "http://"
		if c.flagDevTLS {
			protocol = "https://"
		}
		clusterJson.Nodes = []testcluster.ClusterNode{
			{
				APIAddress: protocol + config.Listeners[0].Address,
			},
		}
		if c.flagDevTLS {
			clusterJson.CACertPath = fmt.Sprintf("%s/%s", certDir, server.VaultDevCAFilename)
		}

		if c.flagDevClusterJson != "" {
			b, err := jsonutil.EncodeJSON(clusterJson)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error encoding cluster.json: %s", err))
				return 1
			}
			err = os.WriteFile(c.flagDevClusterJson, b, 0o600)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error writing cluster.json %q: %s", c.flagDevClusterJson, err))
				return 1
			}
		}
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

	cliDumpCh := make(chan struct{})
	if c.flagCLIDump != "" {
		go func() { cliDumpCh <- struct{}{} }()
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

			// Notify systemd that the server is reloading config
			c.notifySystemd(systemd.SdNotifyReloading)

			// Check for new log level
			config, configErrors, err := c.reloadConfigFiles()
			if err != nil {
				c.logger.Error("could not reload config", "error", err)
				goto RUNRELOADFUNCS
			}

			// Ensure at least one config was found.
			if config == nil {
				c.logger.Error("no config found at reload time")
				goto RUNRELOADFUNCS
			}

			// reporting Errors found in the config
			for _, cErr := range configErrors {
				c.logger.Warn(cErr.String())
			}

			// Note that seal reloading can also be triggered via Core.TriggerSealReload.
			// See the call to Core.SetSealReloadFunc above.
			if reloaded, err := c.reloadSealsOnSigHup(ctx, core, config); err != nil {
				c.UI.Error(fmt.Errorf("error reloading seal config: %s", err).Error())
				config.Seals = core.GetCoreConfigInternal().Seals
				goto RUNRELOADFUNCS
			} else if !reloaded {
				config.Seals = core.GetCoreConfigInternal().Seals
			}

			core.SetConfig(config)

			// reloading custom response headers to make sure we have
			// the most up to date headers after reloading the config file
			if err = core.ReloadCustomResponseHeaders(); err != nil {
				c.logger.Error(err.Error())
			}

			// Setting log request with the new value in the config after reload
			core.ReloadLogRequestsLevel()

			core.ReloadRequestLimiter()

			core.ReloadOverloadController()

			// reloading HCP link
			hcpLink, err = c.reloadHCPLink(hcpLink, config, core, hcpLogger)
			if err != nil {
				c.logger.Error(err.Error())
			}

			// Reload log level for loggers
			if config.LogLevel != "" {
				level, err := loghelper.ParseLogLevel(config.LogLevel)
				if err != nil {
					c.logger.Error("unknown log level found on reload", "level", config.LogLevel)
				} else {
					core.SetLogLevel(level)
				}
			}

			// notify ServiceRegistration that a configuration reload has occurred
			if sr := coreConfig.GetServiceRegistration(); sr != nil {
				var srConfig *map[string]string
				if config.ServiceRegistration != nil {
					srConfig = &config.ServiceRegistration.Config
				}
				sr.NotifyConfigurationReload(srConfig)
			}

			if err := core.ReloadCensusManager(false); err != nil {
				c.UI.Error(err.Error())
			}

			core.ReloadReplicationCanaryWriteInterval()

		RUNRELOADFUNCS:
			if err := c.Reload(c.reloadFuncsLock, c.reloadFuncs, c.flagConfigs, core); err != nil {
				c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
			}

			// Reload license file
			if err = core.EntReloadLicense(); err != nil {
				c.UI.Error(err.Error())
			}

			select {
			case c.licenseReloadedCh <- err:
			default:
			}

			// Let the managedKeyRegistry react to configuration changes (i.e.
			// changes in kms_libraries)
			core.ReloadManagedKeyRegistryConfig()

			// Notify systemd that the server has completed reloading config
			c.notifySystemd(systemd.SdNotifyReady)
		case <-c.SigUSR2Ch:
			c.logger.Info("Received SIGUSR2, dumping goroutines. This is expected behavior. Vault continues to run normally.")
			logWriter := c.logger.StandardWriter(&hclog.StandardLoggerOptions{})
			pprof.Lookup("goroutine").WriteTo(logWriter, 2)

			if os.Getenv("VAULT_STACKTRACE_WRITE_TO_FILE") != "" {
				c.logger.Info("Writing stacktrace to file")

				dir := ""
				path := os.Getenv("VAULT_STACKTRACE_FILE_PATH")
				if path != "" {
					if _, err := os.Stat(path); err != nil {
						c.logger.Error("Checking stacktrace path failed", "error", err)
						continue
					}
					dir = path
				} else {
					dir, err = os.MkdirTemp("", "vault-stacktrace")
					if err != nil {
						c.logger.Error("Could not create temporary directory for stacktrace", "error", err)
						continue
					}
				}

				f, err := os.CreateTemp(dir, "stacktrace")
				if err != nil {
					c.logger.Error("Could not create stacktrace file", "error", err)
					continue
				}

				if err := pprof.Lookup("goroutine").WriteTo(f, 2); err != nil {
					f.Close()
					c.logger.Error("Could not write stacktrace to file", "error", err)
					continue
				}

				c.logger.Info(fmt.Sprintf("Wrote stacktrace to: %s", f.Name()))
				f.Close()
			}

			// We can only get pprof outputs via the API but sometimes Vault can get
			// into a state where it cannot process requests so we can get pprof outputs
			// via SIGUSR2.
			pprofPath := filepath.Join(os.TempDir(), "vault-pprof")
			cpuProfileDuration := time.Second * 1
			err := WritePprofToFile(pprofPath, cpuProfileDuration)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}

			c.logger.Info(fmt.Sprintf("Wrote pprof files to: %s", pprofPath))
		case <-cliDumpCh:
			path := c.flagCLIDump

			if _, err := os.Stat(path); err != nil && !errors.Is(err, os.ErrNotExist) {
				c.logger.Error("Checking cli dump path failed", "error", err)
				continue
			}

			pprofPath := filepath.Join(path, "vault-pprof")
			err := os.MkdirAll(pprofPath, os.ModePerm)
			if err != nil {
				c.logger.Error("Could not create temporary directory for pprof", "error", err)
				continue
			}

			dumps := []string{"goroutine", "heap", "allocs", "threadcreate", "profile"}
			for _, dump := range dumps {
				pFile, err := os.Create(filepath.Join(pprofPath, dump))
				if err != nil {
					c.logger.Error("error creating pprof file", "name", dump, "error", err)
					break
				}

				if dump != "profile" {
					err = pprof.Lookup(dump).WriteTo(pFile, 0)
					if err != nil {
						c.logger.Error("error generating pprof data", "name", dump, "error", err)
						pFile.Close()
						break
					}
				} else {
					// CPU profiles need to run for a duration so we're going to run it
					// just for one second to avoid blocking here.
					if err := pprof.StartCPUProfile(pFile); err != nil {
						c.logger.Error("could not start CPU profile: ", err)
						pFile.Close()
						break
					}
					time.Sleep(time.Second * 1)
					pprof.StopCPUProfile()
				}
				pFile.Close()
			}

			c.logger.Info(fmt.Sprintf("Wrote startup pprof files to: %s", pprofPath))
		}
	}
	// Notify systemd that the server is shutting down
	c.notifySystemd(systemd.SdNotifyStopping)

	// Stop the listeners so that we don't process further client requests.
	c.cleanupGuard.Do(listenerCloseFunc)

	if hcpLink != nil {
		if err := hcpLink.Shutdown(); err != nil {
			c.UI.Error(fmt.Sprintf("Error with HCP Link shutdown: %v", err.Error()))
		}
	}

	// Finalize will wait until after Vault is sealed, which means the
	// request forwarding listeners will also be closed (and also
	// waited for).
	if err := core.Shutdown(); err != nil {
		c.UI.Error(fmt.Sprintf("Error with core shutdown: %s", err))
	}

	// Wait for dependent goroutines to complete
	c.WaitGroup.Wait()
	return retCode
}

func (c *ServerCommand) reloadConfigFiles() (*server.Config, []configutil.ConfigError, error) {
	var config *server.Config
	var configErrors []configutil.ConfigError
	for _, path := range c.flagConfigs {
		current, err := server.LoadConfig(path)
		if err != nil {
			return nil, nil, err
		}

		configErrors = append(configErrors, current.Validate(path)...)

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}

	return config, configErrors, nil
}

func (c *ServerCommand) configureSeals(ctx context.Context, config *server.Config, backend physical.Backend, infoKeys []string, info map[string]string) (*SetSealResponse, io.Reader, error) {
	existingSealGenerationInfo, err := vault.PhysicalSealGenInfo(ctx, backend)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting seal generation info: %v", err)
	}

	hasPartialPaths, err := hasPartiallyWrappedPaths(ctx, backend)
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot determine if there are partially seal wrapped entries in storage: %v", err)
	}
	setSealResponse, err := setSeal(c, config, infoKeys, info, existingSealGenerationInfo, hasPartialPaths)
	if err != nil {
		return nil, nil, err
	}
	if setSealResponse.sealConfigWarning != nil {
		c.UI.Warn(fmt.Sprintf("Warnings during seal configuration: %v", setSealResponse.sealConfigWarning))
	}

	if setSealResponse.barrierSeal == nil {
		return nil, nil, errors.New("Could not create barrier seal! Most likely proper Seal configuration information was not set, but no error was generated.")
	}

	// prepare a secure random reader for core
	entropyAugLogger := c.logger.Named("entropy-augmentation")
	var entropySources []*configutil.EntropySourcerInfo
	for _, sealWrapper := range setSealResponse.barrierSeal.GetAccess().GetEnabledSealWrappersByPriority() {
		if s, ok := sealWrapper.Wrapper.(entropy.Sourcer); ok {
			entropySources = append(entropySources, &configutil.EntropySourcerInfo{
				Sourcer: s,
				Name:    sealWrapper.Name,
			})
		}
	}
	secureRandomReader, err := configutil.CreateSecureRandomReaderFunc(config.SharedConfig, entropySources, entropyAugLogger)
	if err != nil {
		return nil, nil, err
	}

	return setSealResponse, secureRandomReader, nil
}

func (c *ServerCommand) setSealsToFinalize(seals []*vault.Seal) {
	prev := c.sealsToFinalize
	c.sealsToFinalize = seals

	c.finalizeSeals(context.Background(), prev)
}

func (c *ServerCommand) finalizeSeals(ctx context.Context, seals []*vault.Seal) {
	for _, seal := range seals {
		// Ensure that the seal finalizer is called, even if using verify-only
		err := (*seal).Finalize(ctx)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
		}
	}
}

// configureLogging takes the configuration and attempts to parse config values into 'log' friendly configuration values
// If all goes to plan, a logger is created and setup.
func (c *ServerCommand) configureLogging(config *server.Config) (hclog.InterceptLogger, error) {
	// Parse all the log related config
	logLevel, err := loghelper.ParseLogLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	logFormat, err := loghelper.ParseLogFormat(config.LogFormat)
	if err != nil {
		return nil, err
	}

	logRotateDuration, err := parseutil.ParseDurationSecond(config.LogRotateDuration)
	if err != nil {
		return nil, err
	}

	logCfg, err := loghelper.NewLogConfig("vault")
	if err != nil {
		return nil, err
	}
	logCfg.LogLevel = logLevel
	logCfg.LogFormat = logFormat
	logCfg.LogFilePath = config.LogFile
	logCfg.LogRotateDuration = logRotateDuration
	logCfg.LogRotateBytes = config.LogRotateBytes
	logCfg.LogRotateMaxFiles = config.LogRotateMaxFiles

	return loghelper.Setup(logCfg, c.logWriter)
}

func (c *ServerCommand) reloadHCPLink(hcpLinkVault *hcp_link.HCPLinkVault, conf *server.Config, core *vault.Core, hcpLogger hclog.Logger) (*hcp_link.HCPLinkVault, error) {
	// trigger a shutdown
	if hcpLinkVault != nil {
		err := hcpLinkVault.Shutdown()
		if err != nil {
			return nil, err
		}
	}

	if conf.HCPLinkConf == nil {
		// if cloud stanza is not configured, we should not show anything
		// in the seal-status related to HCP link
		core.SetHCPLinkStatus("", "")
		return nil, nil
	}

	// starting HCP link
	hcpLink, err := hcp_link.NewHCPLink(conf.HCPLinkConf, core, hcpLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to restart HCP Link and it is no longer running, %w", err)
	}

	return hcpLink, nil
}

func (c *ServerCommand) notifySystemd(status string) {
	sent, err := systemd.SdNotify(false, status)
	if err != nil {
		c.logger.Error("error notifying systemd", "error", err)
	} else {
		if sent {
			c.logger.Debug("sent systemd notification", "notification", status)
		} else {
			c.logger.Debug("would have sent systemd notification (systemd not present)", "notification", status)
		}
	}
}

func (c *ServerCommand) enableDev(core *vault.Core, coreConfig *vault.CoreConfig) (*vault.InitResult, error) {
	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	var recoveryConfig *vault.SealConfig
	barrierConfig := &vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
		Name:            "shamir",
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
		return nil, fmt.Errorf("failed to check active status: %w", err)
	}
	if err == nil {
		leaderCount := 5
		for !isLeader {
			if leaderCount == 0 {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				return nil, fmt.Errorf("failed to get active status after five seconds; call stack is\n%s", buf)
			}
			time.Sleep(1 * time.Second)
			isLeader, _, _, err = core.Leader()
			if err != nil {
				return nil, fmt.Errorf("failed to check active status: %w", err)
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
			return nil, fmt.Errorf("failed to create root token with ID %q: %w", coreConfig.DevToken, err)
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
		_, err = core.HandleRequest(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to revoke initial root token: %w", err)
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

	if !c.flagDevNoKV {
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
			return nil, fmt.Errorf("error creating default KV store: %w", err)
		}
		if resp.IsError() {
			return nil, fmt.Errorf("failed to create default KV store: %w", resp.Error())
		}
	}

	return init, nil
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
	config *server.Config,
) (string, error) {
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
	return configutil.NormalizeAddr(url.String()), nil
}

func (c *ServerCommand) Reload(lock *sync.RWMutex, reloadFuncs *map[string][]reloadutil.ReloadFunc, configPath []string, core *vault.Core) error {
	lock.RLock()
	defer lock.RUnlock()

	var reloadErrors *multierror.Error

	for k, relFuncs := range *reloadFuncs {
		switch {
		case strings.HasPrefix(k, "listener|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(); err != nil {
						reloadErrors = multierror.Append(reloadErrors, fmt.Errorf("error encountered reloading listener: %w", err))
					}
				}
			}

		case strings.HasPrefix(k, "audit_file|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(); err != nil {
						reloadErrors = multierror.Append(reloadErrors, fmt.Errorf("error encountered reloading file audit device at path %q: %w", strings.TrimPrefix(k, "audit_file|"), err))
					}
				}
			}
		}
	}

	// Set Introspection Endpoint to enabled with new value in the config after reload
	core.ReloadIntrospectionEndpointEnabled()

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
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("could not open pid file: %w", err)
	}
	defer pidFile.Close()

	// Write out the PID
	pid := os.Getpid()
	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return fmt.Errorf("could not write to pid file: %w", err)
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

		timer := time.NewTimer(2 * time.Second)
		select {
		case <-timer.C:
		case <-c.ShutdownCh:
			timer.Stop()
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

type SetSealResponse struct {
	barrierSeal vault.Seal
	unwrapSeal  vault.Seal

	// sealConfigError is present if there was an error configuring wrappers, other than KeyNotFound.
	sealConfigError          error
	sealConfigWarning        error
	hasPartiallyWrappedPaths bool
}

func (r *SetSealResponse) getCreatedSeals() []*vault.Seal {
	var ret []*vault.Seal
	if r.barrierSeal != nil {
		ret = append(ret, &r.barrierSeal)
	}
	if r.unwrapSeal != nil {
		ret = append(ret, &r.unwrapSeal)
	}
	return ret
}

// setSeal return barrierSeal, barrierWrapper, unwrapSeal, all the created seals, and all the provided seals from the configs so we can close them in Run
// The two errors are the sealConfigError and the regular error
func setSeal(c *ServerCommand, config *server.Config, infoKeys []string, info map[string]string, existingSealGenerationInfo *vaultseal.SealGenerationInfo, hasPartiallyWrappedPaths bool) (*SetSealResponse, error) {
	if c.flagDevAutoSeal {
		access, _ := vaultseal.NewTestSeal(nil)
		barrierSeal := vault.NewAutoSeal(access)

		return &SetSealResponse{barrierSeal: barrierSeal}, nil
	}

	// Handle the case where no seal is provided
	switch len(config.Seals) {
	case 0:
		config.Seals = append(config.Seals, &configutil.KMS{
			Type:     vault.SealConfigTypeShamir.String(),
			Priority: 1,
			Name:     "shamir",
		})
	default:
		allSealsDisabled := true
		for _, c := range config.Seals {
			if !c.Disabled {
				allSealsDisabled = false
			} else if c.Type == vault.SealConfigTypeShamir.String() {
				return nil, errors.New("shamir seals cannot be set disabled (they should simply not be set)")
			}
		}
		// If all seals are disabled assume they want to
		// migrate to a shamir seal and simply didn't provide it
		if allSealsDisabled {
			config.Seals = append(config.Seals, &configutil.KMS{
				Type:     vault.SealConfigTypeShamir.String(),
				Priority: 1,
				Name:     "shamir",
			})
		}
	}

	var sealConfigError error
	var sealConfigWarning error
	recordSealConfigError := func(err error) {
		sealConfigError = errors.Join(sealConfigError, err)
	}
	recordSealConfigWarning := func(err error) {
		sealConfigWarning = errors.Join(sealConfigWarning, err)
	}
	enabledSealWrappers := make([]*vaultseal.SealWrapper, 0)
	disabledSealWrappers := make([]*vaultseal.SealWrapper, 0)
	allSealKmsConfigs := make([]*configutil.KMS, 0)

	type infoKeysAndMap struct {
		keys   []string
		theMap map[string]string
	}
	sealWrapperInfoKeysMap := make(map[string]infoKeysAndMap)

	configuredSeals := 0
	for _, configSeal := range config.Seals {
		sealTypeEnvVarName := "VAULT_SEAL_TYPE"
		if configSeal.Priority > 1 {
			sealTypeEnvVarName = sealTypeEnvVarName + "_" + configSeal.Name
		}

		if !configSeal.Disabled && os.Getenv(sealTypeEnvVarName) != "" {
			sealType := os.Getenv(sealTypeEnvVarName)
			configSeal.Type = sealType
		}

		sealLogger := c.logger.ResetNamed(fmt.Sprintf("seal.%s", configSeal.Type))
		c.allLoggers = append(c.allLoggers, sealLogger)

		allSealKmsConfigs = append(allSealKmsConfigs, configSeal)
		var wrapperInfoKeys []string
		wrapperInfoMap := map[string]string{}
		wrapper, wrapperConfigError := configutil.ConfigureWrapper(configSeal, &wrapperInfoKeys, &wrapperInfoMap, sealLogger)
		if wrapperConfigError == nil {
			// for some reason configureWrapper in kms.go returns nil wrapper and nil error for wrapping.WrapperTypeShamir
			if wrapper == nil {
				wrapper = aeadwrapper.NewShamirWrapper()
			}
			configuredSeals++
		} else if config.IsMultisealEnabled() {
			recordSealConfigWarning(fmt.Errorf("error configuring seal: %v", wrapperConfigError))
		} else {
			// It seems that we are checking for this particular error here is to distinguish between a
			// mis-configured seal vs one that fails for another reason. Apparently the only other reason is
			// a key not found error. It seems the intention is for the key not found error to be returned
			// as a seal specific error later
			if !errwrap.ContainsType(wrapperConfigError, new(logical.KeyNotFoundError)) {
				return nil, fmt.Errorf("error parsing Seal configuration: %s", wrapperConfigError)
			} else {
				sealLogger.Error("error configuring seal", "name", configSeal.Name, "err", wrapperConfigError)
				recordSealConfigError(wrapperConfigError)
			}
		}

		sealWrapper := vaultseal.NewSealWrapper(
			wrapper,
			configSeal.Priority,
			configSeal.Name,
			configSeal.Type,
			configSeal.Disabled,
			wrapperConfigError == nil,
		)

		if configSeal.Disabled {
			disabledSealWrappers = append(disabledSealWrappers, sealWrapper)
		} else {
			enabledSealWrappers = append(enabledSealWrappers, sealWrapper)
		}

		sealWrapperInfoKeysMap[sealWrapper.Name] = infoKeysAndMap{
			keys:   wrapperInfoKeys,
			theMap: wrapperInfoMap,
		}
	}

	if len(enabledSealWrappers) == 0 && len(disabledSealWrappers) == 0 && sealConfigWarning != nil {
		// All of them errored out, so warnings are now errors
		recordSealConfigError(sealConfigWarning)
		sealConfigWarning = nil
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Set the info keys, this modifies the function arguments `info` and `infoKeys`
	// TODO(SEALHA): Why are we doing this? What is its use?
	appendWrapperInfoKeys := func(prefix string, sealWrappers []*vaultseal.SealWrapper) {
		if len(sealWrappers) == 0 {
			return
		}
		useName := false
		if len(sealWrappers) > 1 {
			useName = true
		}
		for _, sealWrapper := range sealWrappers {
			if useName {
				prefix = fmt.Sprintf("%s %s ", prefix, sealWrapper.Name)
			}
			for _, k := range sealWrapperInfoKeysMap[sealWrapper.Name].keys {
				infoKeys = append(infoKeys, prefix+k)
				info[prefix+k] = sealWrapperInfoKeysMap[sealWrapper.Name].theMap[k]
			}
		}
	}
	appendWrapperInfoKeys("", enabledSealWrappers)
	appendWrapperInfoKeys("Old", disabledSealWrappers)

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Compute seal generation
	sealGenerationInfo, err := c.computeSealGenerationInfo(existingSealGenerationInfo, allSealKmsConfigs, hasPartiallyWrappedPaths, config.IsMultisealEnabled())
	if err != nil {
		return nil, err
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Create the Seals

	containsShamir := func(sealWrappers []*vaultseal.SealWrapper) bool {
		for _, si := range sealWrappers {
			if vault.SealConfigTypeShamir.IsSameAs(si.SealConfigType) {
				return true
			}
		}
		return false
	}

	var barrierSeal vault.Seal
	var unwrapSeal vault.Seal

	sealLogger := c.logger
	switch {
	case len(enabledSealWrappers) == 0:
		return nil, errors.Join(sealConfigWarning, errors.New("no enabled Seals in configuration"))
	case configuredSeals == 0:
		return nil, errors.Join(sealConfigWarning, errors.New("no seals were successfully initialized"))
	case len(enabledSealWrappers) == 1 && containsShamir(enabledSealWrappers):
		// The barrier seal is Shamir. If there are any disabled seals, then we put them all in the same
		// autoSeal.
		a, err := vaultseal.NewAccess(sealLogger, sealGenerationInfo, enabledSealWrappers)
		if err != nil {
			return nil, err
		}
		barrierSeal = vault.NewDefaultSeal(a)
		if len(disabledSealWrappers) > 0 {
			a, err = vaultseal.NewAccess(sealLogger, sealGenerationInfo, disabledSealWrappers)
			if err != nil {
				return nil, err
			}
			unwrapSeal = vault.NewAutoSeal(a)
		} else if sealGenerationInfo.Generation == 1 {
			// First generation, and shamir, with no disabled wrapperrs, so there can be no wrapped values
			sealGenerationInfo.SetRewrapped(true)
		}

	case len(disabledSealWrappers) == 1 && containsShamir(disabledSealWrappers):
		// The unwrap seal is Shamir, we are migrating to an autoSeal.
		a, err := vaultseal.NewAccess(sealLogger, sealGenerationInfo, enabledSealWrappers)
		if err != nil {
			return nil, err
		}
		barrierSeal = vault.NewAutoSeal(a)
		a, err = vaultseal.NewAccess(sealLogger, sealGenerationInfo, disabledSealWrappers)
		if err != nil {
			return nil, err
		}
		unwrapSeal = vault.NewDefaultSeal(a)

	case config.IsMultisealEnabled():
		// We know we are not using Shamir seal, that we are not migrating away from one, and multi seal is supported,
		// so just put enabled and disabled wrappers on the same seal Access
		allSealWrappers := append(enabledSealWrappers, disabledSealWrappers...)
		a, err := vaultseal.NewAccess(sealLogger, sealGenerationInfo, allSealWrappers)
		if err != nil {
			return nil, err
		}
		barrierSeal = vault.NewAutoSeal(a)
		if configuredSeals < len(enabledSealWrappers) {
			c.UI.Warn("WARNING: running with fewer than all configured seals during unseal.  Will not be fully highly available until errors are corrected and Vault restarted.")
		}
	case len(enabledSealWrappers) == 1:
		// We may have multiple seals disabled, but we know Shamir is not one of them.
		a, err := vaultseal.NewAccess(sealLogger, sealGenerationInfo, enabledSealWrappers)
		if err != nil {
			return nil, err
		}
		barrierSeal = vault.NewAutoSeal(a)
		if len(disabledSealWrappers) > 0 {
			a, err = vaultseal.NewAccess(sealLogger, sealGenerationInfo, disabledSealWrappers)
			if err != nil {
				return nil, err
			}
			unwrapSeal = vault.NewAutoSeal(a)
		}

	default:
		// We know there are multiple enabled seals but multi seal is not supported.
		return nil, errors.Join(sealConfigWarning, errors.New("error: more than one enabled seal found"))
	}

	return &SetSealResponse{
		barrierSeal:              barrierSeal,
		unwrapSeal:               unwrapSeal,
		sealConfigError:          sealConfigError,
		sealConfigWarning:        sealConfigWarning,
		hasPartiallyWrappedPaths: hasPartiallyWrappedPaths,
	}, nil
}

func (c *ServerCommand) computeSealGenerationInfo(existingSealGenInfo *vaultseal.SealGenerationInfo, sealConfigs []*configutil.KMS, hasPartiallyWrappedPaths, multisealEnabled bool) (*vaultseal.SealGenerationInfo, error) {
	generation := uint64(1)

	if existingSealGenInfo != nil {
		// This forces a seal re-wrap on all seal related config changes, as we can't
		// be sure what effect the config change might do. This is purposefully different
		// from within the Validate call below that just matches on seal configs based
		// on name/type.
		if cmp.Equal(existingSealGenInfo.Seals, sealConfigs) {
			return existingSealGenInfo, nil
		}
		generation = existingSealGenInfo.Generation + 1
	}
	c.logger.Info("incrementing seal generation", "generation", generation)

	// If the stored copy doesn't match the current configuration, we introduce a new generation
	// which keeps track if a rewrap of all CSPs and seal wrapped values has completed (initially false).
	newSealGenInfo := &vaultseal.SealGenerationInfo{
		Generation: generation,
		Seals:      sealConfigs,
		Enabled:    multisealEnabled,
	}

	if multisealEnabled || (existingSealGenInfo != nil && existingSealGenInfo.Enabled) {
		err := newSealGenInfo.Validate(existingSealGenInfo, hasPartiallyWrappedPaths)
		if err != nil {
			return nil, err
		}
	}

	return newSealGenInfo, nil
}

func hasPartiallyWrappedPaths(ctx context.Context, backend physical.Backend) (bool, error) {
	paths, err := vault.GetPartiallySealWrappedPaths(ctx, backend)
	if err != nil {
		return false, err
	}

	return len(paths) > 0, nil
}

func initHaBackend(c *ServerCommand, config *server.Config, coreConfig *vault.CoreConfig, backend physical.Backend) (bool, error) {
	// Initialize the separate HA storage backend, if it exists
	var ok bool
	if config.HAStorage != nil {
		if config.Storage.Type == storageTypeRaft && config.HAStorage.Type == storageTypeRaft {
			return false, fmt.Errorf("Raft cannot be set both as 'storage' and 'ha_storage'. Setting 'storage' to 'raft' will automatically set it up for HA operations as well")
		}

		if config.Storage.Type == storageTypeRaft {
			return false, fmt.Errorf("HA storage cannot be declared when Raft is the storage type")
		}

		factory, exists := c.PhysicalBackends[config.HAStorage.Type]
		if !exists {
			return false, fmt.Errorf("Unknown HA storage type %s", config.HAStorage.Type)
		}

		namedHALogger := c.logger.Named("ha." + config.HAStorage.Type)
		c.allLoggers = append(c.allLoggers, namedHALogger)
		habackend, err := factory(config.HAStorage.Config, namedHALogger)
		if err != nil {
			return false, fmt.Errorf("Error initializing HA storage of type %s: %s", config.HAStorage.Type, err)
		}

		if coreConfig.HAPhysical, ok = habackend.(physical.HABackend); !ok {
			return false, fmt.Errorf("Specified HA storage does not support HA")
		}

		if !coreConfig.HAPhysical.HAEnabled() {
			return false, fmt.Errorf("Specified HA storage has HA support disabled; please consult documentation")
		}

		coreConfig.RedirectAddr = config.HAStorage.RedirectAddr
		disableClustering := config.HAStorage.DisableClustering

		if config.HAStorage.Type == storageTypeRaft && disableClustering {
			return disableClustering, fmt.Errorf("Disable clustering cannot be set to true when Raft is the HA storage type")
		}

		if !disableClustering {
			coreConfig.ClusterAddr = config.HAStorage.ClusterAddr
		}
	} else {
		if coreConfig.HAPhysical, ok = backend.(physical.HABackend); ok {
			coreConfig.RedirectAddr = config.Storage.RedirectAddr
			disableClustering := config.Storage.DisableClustering

			if (config.Storage.Type == storageTypeRaft) && disableClustering {
				return disableClustering, fmt.Errorf("Disable clustering cannot be set to true when Raft is the storage type")
			}

			if !disableClustering {
				coreConfig.ClusterAddr = config.Storage.ClusterAddr
			}
		}
	}
	return config.DisableClustering, nil
}

func determineRedirectAddr(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config) error {
	var retErr error
	if envRA := os.Getenv("VAULT_API_ADDR"); envRA != "" {
		coreConfig.RedirectAddr = configutil.NormalizeAddr(envRA)
	} else if envRA := os.Getenv("VAULT_REDIRECT_ADDR"); envRA != "" {
		coreConfig.RedirectAddr = configutil.NormalizeAddr(envRA)
	} else if envAA := os.Getenv("VAULT_ADVERTISE_ADDR"); envAA != "" {
		coreConfig.RedirectAddr = configutil.NormalizeAddr(envAA)
	}

	// Attempt to detect the redirect address, if possible
	if coreConfig.RedirectAddr == "" {
		c.logger.Warn("no `api_addr` value specified in config or in VAULT_API_ADDR; falling back to detection if possible, but this value should be manually set")
	}

	var ok bool
	var detect physical.RedirectDetect
	if coreConfig.HAPhysical != nil && coreConfig.HAPhysical.HAEnabled() {
		detect, ok = coreConfig.HAPhysical.(physical.RedirectDetect)
	} else {
		detect, ok = coreConfig.Physical.(physical.RedirectDetect)
	}
	if ok && coreConfig.RedirectAddr == "" {
		redirect, err := c.detectRedirect(detect, config)
		// the following errors did not cause Run to return, so I'm not returning these
		// as errors.
		if err != nil {
			retErr = fmt.Errorf("Error detecting api address: %s", err)
		} else if redirect == "" {
			retErr = fmt.Errorf("Failed to detect api address")
		} else {
			coreConfig.RedirectAddr = redirect
		}
	}
	if coreConfig.RedirectAddr == "" && c.flagDev {
		protocol := "http"
		if c.flagDevTLS {
			protocol = "https"
		}
		coreConfig.RedirectAddr = configutil.NormalizeAddr(fmt.Sprintf("%s://%s", protocol, config.Listeners[0].Address))
	}
	return retErr
}

func findClusterAddress(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config, disableClustering bool) error {
	if disableClustering {
		coreConfig.ClusterAddr = ""
	} else if envCA := os.Getenv("VAULT_CLUSTER_ADDR"); envCA != "" {
		coreConfig.ClusterAddr = configutil.NormalizeAddr(envCA)
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
			return fmt.Errorf("Error parsing synthesized cluster address %s: %v", addrToUse, err)
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			// This sucks, as it's a const in the function but not exported in the package
			if strings.Contains(err.Error(), "missing port in address") {
				host = u.Host
				port = "443"
			} else {
				return fmt.Errorf("Error parsing api address: %v", err)
			}
		}
		nPort, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("Error parsing synthesized address; failed to convert %q to a numeric: %v", port, err)
		}
		u.Host = net.JoinHostPort(host, strconv.Itoa(nPort+1))
		// Will always be TLS-secured
		u.Scheme = "https"
		coreConfig.ClusterAddr = configutil.NormalizeAddr(u.String())
	}

CLUSTER_SYNTHESIS_COMPLETE:

	if coreConfig.RedirectAddr == coreConfig.ClusterAddr && len(coreConfig.RedirectAddr) != 0 {
		return fmt.Errorf("Address %q used for both API and cluster addresses", coreConfig.RedirectAddr)
	}

	if coreConfig.ClusterAddr != "" {
		rendered, err := configutil.ParseSingleIPTemplate(coreConfig.ClusterAddr)
		if err != nil {
			return fmt.Errorf("Error parsing cluster address %s: %v", coreConfig.ClusterAddr, err)
		}
		coreConfig.ClusterAddr = rendered
		// Force https as we'll always be TLS-secured
		u, err := url.ParseRequestURI(coreConfig.ClusterAddr)
		if err != nil {
			return fmt.Errorf("Error parsing cluster address %s: %v", coreConfig.ClusterAddr, err)
		}
		u.Scheme = "https"
		coreConfig.ClusterAddr = u.String()
	}
	return nil
}

func runUnseal(c *ServerCommand, core *vault.Core, ctx context.Context) {
	for {
		err := core.UnsealWithStoredKeys(ctx)
		if err == nil {
			return
		}

		if vault.IsFatalError(err) {
			c.logger.Error("error unsealing core", "error", err)
			return
		}
		c.logger.Warn("failed to unseal core", "error", err)

		timer := time.NewTimer(5 * time.Second)
		select {
		case <-c.ShutdownCh:
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func createCoreConfig(c *ServerCommand, config *server.Config, backend physical.Backend, configSR sr.ServiceRegistration, barrierSeal, unwrapSeal vault.Seal,
	metricsHelper *metricsutil.MetricsHelper, metricSink *metricsutil.ClusterMetricSink, secureRandomReader io.Reader,
) vault.CoreConfig {
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
		LogLevel:                       config.LogLevel,
		Logger:                         c.logger,
		DetectDeadlocks:                config.DetectDeadlocks,
		ImpreciseLeaseRoleTracking:     config.ImpreciseLeaseRoleTracking,
		DisableSentinelTrace:           config.DisableSentinelTrace,
		DisableCache:                   config.DisableCache,
		DisableMlock:                   config.DisableMlock,
		MaxLeaseTTL:                    config.MaxLeaseTTL,
		DefaultLeaseTTL:                config.DefaultLeaseTTL,
		ClusterName:                    config.ClusterName,
		CacheSize:                      config.CacheSize,
		PluginDirectory:                config.PluginDirectory,
		PluginTmpdir:                   config.PluginTmpdir,
		PluginFileUid:                  config.PluginFileUid,
		PluginFilePermissions:          config.PluginFilePermissions,
		EnableUI:                       config.EnableUI,
		EnableRaw:                      config.EnableRawEndpoint,
		EnableIntrospection:            config.EnableIntrospectionEndpoint,
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
		License:                        config.License,
		LicensePath:                    config.LicensePath,
		DisableSSCTokens:               config.DisableSSCTokens,
		Experiments:                    config.Experiments,
		AdministrativeNamespacePath:    config.AdministrativeNamespacePath,
	}

	if c.flagDev {
		coreConfig.EnableRaw = true
		coreConfig.EnableIntrospection = true
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
	return *coreConfig
}

func runListeners(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config, configSR sr.ServiceRegistration) error {
	if sd := coreConfig.GetServiceRegistration(); sd != nil {
		if err := configSR.Run(c.ShutdownCh, c.WaitGroup, coreConfig.RedirectAddr); err != nil {
			return fmt.Errorf("Error running service_registration of type %s: %s", config.ServiceRegistration.Type, err)
		}
	}
	return nil
}

func initDevCore(c *ServerCommand, coreConfig *vault.CoreConfig, config *server.Config, core *vault.Core, certDir string, clusterJSON *testcluster.ClusterJson) error {
	if c.flagDev && !c.flagDevSkipInit {

		init, err := c.enableDev(core, coreConfig)
		if err != nil {
			return fmt.Errorf("Error initializing Dev mode: %s", err)
		}

		if clusterJSON != nil {
			clusterJSON.RootToken = init.RootToken
		}

		var plugins, pluginsNotLoaded []string
		if c.flagDevPluginDir != "" && c.flagDevPluginInit {

			f, err := os.Open(c.flagDevPluginDir)
			if err != nil {
				return fmt.Errorf("Error reading plugin dir: %s", err)
			}

			list, err := f.Readdirnames(0)
			f.Close()
			if err != nil {
				return fmt.Errorf("Error listing plugins: %s", err)
			}

			for _, name := range list {
				path := filepath.Join(f.Name(), name)
				if err := c.addPlugin(path, init.RootToken, core); err != nil {
					if !errwrap.Contains(err, plugincatalog.ErrPluginBadType.Error()) {
						return fmt.Errorf("Error enabling plugin %s: %s", name, err)
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
					c.UI.Warn("You may need to set the following environment variables:")
					c.UI.Warn("")

					protocol := "http://"
					if c.flagDevTLS {
						protocol = "https://"
					}

					endpointURL := protocol + config.Listeners[0].Address
					if runtime.GOOS == "windows" {
						c.UI.Warn("PowerShell:")
						c.UI.Warn(fmt.Sprintf("    $env:VAULT_ADDR=\"%s\"", endpointURL))
						c.UI.Warn("cmd.exe:")
						c.UI.Warn(fmt.Sprintf("    set VAULT_ADDR=%s", endpointURL))
					} else {
						c.UI.Warn(fmt.Sprintf("    $ export VAULT_ADDR='%s'", endpointURL))
					}

					if c.flagDevTLS {
						if runtime.GOOS == "windows" {
							c.UI.Warn("PowerShell:")
							c.UI.Warn(fmt.Sprintf("    $env:VAULT_CACERT=\"%s/vault-ca.pem\"", certDir))
							c.UI.Warn("cmd.exe:")
							c.UI.Warn(fmt.Sprintf("    set VAULT_CACERT=%s/vault-ca.pem", certDir))
						} else {
							c.UI.Warn(fmt.Sprintf("    $ export VAULT_CACERT='%s/vault-ca.pem'", certDir))
						}
						c.UI.Warn("")
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
	return nil
}

// Initialize the HTTP servers
func startHttpServers(c *ServerCommand, core *vault.Core, config *server.Config, lns []listenerutil.Listener) error {
	for _, ln := range lns {
		if ln.Config == nil {
			return fmt.Errorf("found nil listener config after parsing")
		}

		if err := config2.IsValidListener(ln.Config); err != nil {
			return err
		}

		handler := vaulthttp.Handler.Handler(&vault.HandlerProperties{
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
	return nil
}

// reloadSealsOnLeaderActivation checks to see if the in-memory seal generation info is stale, and if so,
// reloads the seal configuration.
func (c *ServerCommand) reloadSealsOnLeaderActivation(ctx context.Context, core *vault.Core) error {
	existingSealGenerationInfo, err := vault.PhysicalSealGenInfo(ctx, core.PhysicalAccess())
	if err != nil {
		return fmt.Errorf("error checking for stale seal generation info: %w", err)
	}
	if existingSealGenerationInfo == nil {
		c.logger.Debug("not reloading seals config since there is no seal generation info in storage")
		return nil
	}

	currentSealGenerationInfo := core.SealAccess().GetAccess().GetSealGenerationInfo()
	if currentSealGenerationInfo == nil {
		c.logger.Debug("not reloading seal config since there is no current generation info (the seal has not been initialized)")
		return nil
	}
	if currentSealGenerationInfo.Generation >= existingSealGenerationInfo.Generation {
		c.logger.Debug("seal generation info is up to date, not reloading seal configuration")
		return nil
	}

	// Reload seal configuration

	config, _, err := c.reloadConfigFiles()
	if err != nil {
		return fmt.Errorf("error reading configuration files while reloading seal configuration: %w", err)
	}
	if config == nil {
		return errors.New("no configuration files found while reloading seal configuration")
	}
	reloaded, err := c.reloadSeals(ctx, false, core, config)
	if reloaded {
		core.SetConfig(config)
	}
	return err
}

// reloadSealsOnSigHup will reload seal configurtion as a result of receiving a SIGHUP signal.
func (c *ServerCommand) reloadSealsOnSigHup(ctx context.Context, core *vault.Core, config *server.Config) (bool, error) {
	return c.reloadSeals(ctx, true, core, config)
}

// reloadSeals reloads configuration files and determines whether it needs to re-create the Seal.Access() objects.
// This function needs do detect that core.SealAccess() is no longer using the seal Wrapper that is specified
// in the seal configuration files.
// This function returns true if the newConfig was used to re-create the Seal.Access() objects. In other words,
// if false is returned, there were no changes done to the seals.
func (c *ServerCommand) reloadSeals(ctx context.Context, grabStateLock bool, core *vault.Core, newConfig *server.Config) (bool, error) {
	if core.IsInSealMigrationMode(grabStateLock) {
		c.logger.Debug("not reloading seal configuration since Vault is in migration mode")
		return false, nil
	}

	currentConfig := core.GetCoreConfigInternal()

	// We only want to reload if multiseal is currently enabled, or it is being enabled
	if !(currentConfig.IsMultisealEnabled() || newConfig.IsMultisealEnabled()) {
		c.logger.Debug("not reloading seal configuration since enable_multiseal is not set, nor is it being disabled")
		return false, nil
	}

	if conf, err := core.PhysicalBarrierSealConfig(ctx); err != nil {
		return false, fmt.Errorf("error reading barrier seal configuration from storage while reloading seals: %w", err)
	} else if conf == nil {
		c.logger.Debug("not reloading seal configuration since there is no barrier config in storage (the seal has not been initialized)")
		return false, nil
	}

	if core.SealAccess().BarrierSealConfigType() == vault.SealConfigTypeShamir {
		switch {
		case len(newConfig.Seals) == 0:
			// We are fine, our ServerCommand.reloadConfigFiles() does not do the "automagic" creation
			// of the Shamir seal configuration.
			c.logger.Debug("not reloading seal configuration since the new one has no seal stanzas")
			return false, nil

		case len(newConfig.Seals) == 1 && newConfig.Seals[0].Disabled:
			// If we have only one seal and it is disabled, it means that the newConfig wants to migrate
			// to Shamir, which is not supported by seal reloading.
			c.logger.Debug("not reloading seal configuration since the new one specifies migration to Shamir")
			return false, nil

		case len(newConfig.Seals) == 1 && newConfig.Seals[0].Type == vault.SealConfigTypeShamir.String():
			// Having a single Shamir seal in newConfig is not really possible, since a Shamir seal
			// is specified in configuration by *not* having a seal stanza. If we were to hit this
			// case, though, it is equivalent to trying to migrate to Shamir, which is not supported
			// by seal reloading.
			c.logger.Debug("not reloading seal configuration since the new one has single Shamir stanza")
			return false, nil
		}
	}

	// Verify that the new config we picked up is not trying to migrate from autoseal to shamir
	if len(newConfig.Seals) == 1 && newConfig.Seals[0].Disabled {
		// If we get here, it means the node was not started in migration mode, but the new config says
		// we should go into migration mode. This case should be caught by the core.IsInSealMigrationMode()
		// above.

		return false, errors.New("not reloading seal configuration: moving from autoseal to shamir requires seal migration")
	}

	// Verify that the new config we picked up is not trying to migrate shamir to autoseal
	if core.SealAccess().BarrierSealConfigType() == vault.SealConfigTypeShamir {
		return false, errors.New("not reloading seal configuration: moving from Shamir to autoseal requires seal migration")
	}

	infoKeysReload := make([]string, 0)
	infoReload := make(map[string]string)

	core.SetMultisealEnabled(newConfig.IsMultisealEnabled())
	setSealResponse, secureRandomReader, err := c.configureSeals(ctx, newConfig, core.PhysicalAccess(), infoKeysReload, infoReload)
	if err != nil {
		return false, fmt.Errorf("error reloading seal configuration: %w", err)
	}
	if setSealResponse.sealConfigError != nil {
		return false, fmt.Errorf("error reloading seal configuration: %w", setSealResponse.sealConfigError)
	}

	newGen := setSealResponse.barrierSeal.GetAccess().GetSealGenerationInfo()

	var standby, perf bool
	if grabStateLock {
		// If grabStateLock is false we know we are on a leader activation
		standby, perf = core.StandbyStates()
	}
	switch {
	case !perf && !standby:
		c.logger.Debug("persisting reloaded seals as we are the active node")
		err = core.SetSeals(ctx, grabStateLock, setSealResponse.barrierSeal, secureRandomReader, !newGen.IsRewrapped() || setSealResponse.hasPartiallyWrappedPaths)
		if err != nil {
			return false, fmt.Errorf("error setting seal: %s", err)
		}

		if err := core.SetPhysicalSealGenInfo(ctx, newGen); err != nil {
			c.logger.Warn("could not update seal information in storage", "err", err)
		}
	case perf:
		c.logger.Debug("updating reloaded seals in memory on perf standby")
		err = core.SetSealsOnPerfStandby(ctx, grabStateLock, setSealResponse.barrierSeal, secureRandomReader)
		if err != nil {
			return false, fmt.Errorf("error setting seal on perf standby: %s", err)
		}
	default:
		return false, errors.New("skipping seal reload as we are a standby")
	}

	// finalize the old seals and set the new seals as the current ones
	c.setSealsToFinalize(setSealResponse.getCreatedSeals())

	c.logger.Debug("seal configuration reloaded successfully")

	return true, nil
}

// Attempt to read the cluster name from the insecure storage.
func (c *ServerCommand) readClusterNameFromInsecureStorage(b physical.Backend) (string, error) {
	ctx := context.Background()
	entry, err := b.Get(ctx, "core/cluster/local/name")
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	// Decode JSON data into the map

	if entry != nil {
		if err := jsonutil.DecodeJSON(entry.Value, &result); err != nil {
			return "", fmt.Errorf("failed to decode JSON data: %w", err)
		}
	}

	// Retrieve the value of the "name" field from the map
	name, ok := result["name"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract name field from decoded JSON")
	}

	return name, nil
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
	logger hclog.Logger
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
