package command

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/helper/metricsutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/internalshared/reloadutil"
	physconsul "github.com/hashicorp/vault/physical/consul"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/version"
	sr "github.com/hashicorp/vault/serviceregistration"
	srconsul "github.com/hashicorp/vault/serviceregistration/consul"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/diagnose"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const OperatorDiagnoseEnableEnv = "VAULT_DIAGNOSE"

const CoreUninitializedErr = "diagnose cannot attempt this step because core could not be initialized"
const BackendUninitializedErr = "diagnose cannot attempt this step because backend could not be initialized"
const CoreConfigUninitializedErr = "diagnose cannot attempt this step because core config could not be set"

var (
	_ cli.Command             = (*OperatorDiagnoseCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorDiagnoseCommand)(nil)
)

type OperatorDiagnoseCommand struct {
	*BaseCommand
	diagnose *diagnose.Session

	flagDebug    bool
	flagSkips    []string
	flagConfigs  []string
	cleanupGuard sync.Once

	reloadFuncsLock      *sync.RWMutex
	reloadFuncs          *map[string][]reloadutil.ReloadFunc
	ServiceRegistrations map[string]sr.Factory
	startedCh            chan struct{} // for tests
	reloadedCh           chan struct{} // for tests
	skipEndEnd           bool          // for tests
}

func (c *OperatorDiagnoseCommand) Synopsis() string {
	return "Troubleshoot problems starting Vault"
}

func (c *OperatorDiagnoseCommand) Help() string {
	helpText := `
Usage: vault operator diagnose 

  This command troubleshoots Vault startup issues, such as TLS configuration or
  auto-unseal. It should be run using the same environment variables and configuration
  files as the "vault server" command, so that startup problems can be accurately
  reproduced.

  Start diagnose with a configuration file:
    
     $ vault operator diagnose -config=/etc/vault/config.hcl

  Perform a diagnostic check while Vault is still running:

     $ vault operator diagnose -config=/etc/vault/config.hcl -skip=listener

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *OperatorDiagnoseCommand) Flags() *FlagSets {
	set := NewFlagSets(c.UI)
	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&StringSliceVar{
		Name:   "config",
		Target: &c.flagConfigs,
		Completion: complete.PredictOr(
			complete.PredictFiles("*.hcl"),
			complete.PredictFiles("*.json"),
			complete.PredictDirs("*"),
		),
		Usage: "Path to a Vault configuration file or directory of configuration " +
			"files. This flag can be specified multiple times to load multiple " +
			"configurations. If the path is a directory, all files which end in " +
			".hcl or .json are loaded.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   "skip",
		Target: &c.flagSkips,
		Usage:  "Skip the health checks named as arguments. May be 'listener', 'storage', or 'autounseal'.",
	})

	f.BoolVar(&BoolVar{
		Name:    "debug",
		Target:  &c.flagDebug,
		Default: false,
		Usage:   "Dump all information collected by Diagnose.",
	})
	return set
}

func (c *OperatorDiagnoseCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *OperatorDiagnoseCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

const (
	status_unknown = "[      ] "
	status_ok      = "\u001b[32m[  ok  ]\u001b[0m "
	status_failed  = "\u001b[31m[failed]\u001b[0m "
	status_warn    = "\u001b[33m[ warn ]\u001b[0m "
	same_line      = "\u001b[F"
)

func (c *OperatorDiagnoseCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	return c.RunWithParsedFlags()
}

func (c *OperatorDiagnoseCommand) RunWithParsedFlags() int {

	if len(c.flagConfigs) == 0 {
		c.UI.Error("Must specify a configuration file using -config.")
		return 1
	}

	c.UI.Output(version.GetVersion().FullVersionNumber(true))
	ctx := diagnose.Context(context.Background(), c.diagnose)
	err := c.offlineDiagnostics(ctx)

	if err != nil {
		return 1
	}
	return 0
}

func (c *OperatorDiagnoseCommand) offlineDiagnostics(ctx context.Context) error {
	rloadFuncs := make(map[string][]reloadutil.ReloadFunc)
	server := &ServerCommand{
		// TODO: set up a different one?
		// In particular, a UI instance that won't output?
		BaseCommand: c.BaseCommand,

		// TODO: refactor to a common place?
		AuditBackends:        auditBackends,
		CredentialBackends:   credentialBackends,
		LogicalBackends:      logicalBackends,
		PhysicalBackends:     physicalBackends,
		ServiceRegistrations: serviceRegistrations,

		// TODO: other ServerCommand options?

		logger:          log.NewInterceptLogger(nil),
		allLoggers:      []log.Logger{},
		reloadFuncs:     &rloadFuncs,
		reloadFuncsLock: new(sync.RWMutex),
	}

	ctx, span := diagnose.StartSpan(ctx, "initialization")
	defer span.End()

	server.flagConfigs = c.flagConfigs
	config, err := server.parseConfig()
	if err != nil {
		return diagnose.SpotError(ctx, "parse-config", err)
	} else {
		diagnose.SpotOk(ctx, "parse-config", "")
	}

	var metricSink *metricsutil.ClusterMetricSink
	var metricsHelper *metricsutil.MetricsHelper
	if err := diagnose.Test(ctx, "setup-telemetry", func(ctx context.Context) error {
		var prometheusEnabled bool
		var inmemMetrics *metrics.InmemSink
		inmemMetrics, metricSink, prometheusEnabled, err = configutil.SetupTelemetry(&configutil.SetupTelemetryOpts{
			Config:      config.Telemetry,
			Ui:          c.UI,
			ServiceName: "vault",
			DisplayName: "Vault",
			UserAgent:   useragent.String(),
			ClusterName: config.ClusterName,
		})
		metricsHelper = metricsutil.NewMetricsHelper(inmemMetrics, prometheusEnabled)

		// TODO: Comment this error back in. Currently commenting this error in yields
		// indeterministic test behavior, where some subset of operator_diagnose tests have
		// AlreadyRegisteredError errors. We need these metrics values to be initialized for
		// to add them to the coreConfig in later steps, so we have to keep this step in as a
		// no-op for now.

		// if err != nil {
		// 	return fmt.Errorf("Error initializing telemetry: %s", err)
		// }

		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	var backend *physical.Backend
	if err := diagnose.Test(ctx, "storage", func(ctx context.Context) error {
		b, err := server.setupStorage(config)
		if err != nil {
			return err
		}
		backend = &b

		dirAccess := diagnose.ConsulDirectAccess(config.HAStorage.Config)
		if dirAccess != "" {
			diagnose.Warn(ctx, dirAccess)
		}

		if config.Storage != nil && config.Storage.Type == storageTypeConsul {
			err = physconsul.SetupSecureTLS(api.DefaultConfig(), config.Storage.Config, server.logger, true)
			if err != nil {
				return err
			}

			dirAccess := diagnose.ConsulDirectAccess(config.Storage.Config)
			if dirAccess != "" {
				diagnose.Warn(ctx, dirAccess)
			}
		}

		if config.HAStorage != nil && config.HAStorage.Type == storageTypeConsul {
			err = physconsul.SetupSecureTLS(api.DefaultConfig(), config.HAStorage.Config, server.logger, true)
			if err != nil {
				return err
			}
		}

		// Attempt to use storage backend
		if !c.skipEndEnd {
			err = diagnose.StorageEndToEndLatencyCheck(ctx, *backend)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	var configSR sr.ServiceRegistration
	if err := diagnose.Test(ctx, "service-discovery", func(ctx context.Context) error {
		if config.ServiceRegistration == nil {
			return fmt.Errorf("No service registration config")
		}
		srConfig := config.ServiceRegistration.Config
		if config.ServiceRegistration != nil && config.ServiceRegistration.Type == "consul" {
			dirAccess := diagnose.ConsulDirectAccess(config.ServiceRegistration.Config)
			if dirAccess != "" {
				diagnose.Warn(ctx, dirAccess)
			}

			// SetupSecureTLS for service discovery uses the same cert and key to set up physical
			// storage. See the consul package in physical for details.
			err = srconsul.SetupSecureTLS(api.DefaultConfig(), srConfig, server.logger, true)
			if err != nil {
				return err
			}

			// Initialize the Service Discovery, if there is one
			configSR, err = beginServiceRegistration(server, config)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	var barrierWrapper *wrapping.Wrapper
	var barrierSeal *vault.Seal
	var unwrapSeal *vault.Seal
	if err := diagnose.Test(ctx, "create-seal", func(context.Context) error {
		var seals []vault.Seal
		var sealConfigError error
		bS, bW, uS, seals, sealConfigError, err := setSeal(server, config, make([]string, 0), make(map[string]string))
		// Check error here
		if err != nil {
			return err
		}
		if sealConfigError != nil {
			return fmt.Errorf("seal could not be configured: seals may already be initialized")
		}

		barrierSeal = &bS
		barrierWrapper = &bW
		unwrapSeal = &uS
		if seals != nil {
			for _, seal := range seals {
				// Ensure that the seal finalizer is called, even if using verify-only
				defer func(seal *vault.Seal) {
					err = (*seal).Finalize(context.Background())
					if err != nil {
						c.UI.Error(fmt.Sprintf("Error finalizing seals: %v", err))
					}
				}(&seal)
			}
		}

		if barrierSeal == nil {
			return fmt.Errorf("could not create barrier seal! Most likely proper Seal configuration information was not set, but no error was generated")
		}
		return nil
	}); err != nil {
		return err
	}

	var coreConfig vault.CoreConfig
	if err := diagnose.Test(ctx, "setup-core", func(ctx context.Context) error {
		// prepare a secure random reader for core
		secureRandomReader, err := configutil.CreateSecureRandomReaderFunc(config.SharedConfig, *barrierWrapper)
		if err != nil {
			return err
		}
		if backend == nil {
			return fmt.Errorf(BackendUninitializedErr)
		}
		coreConfig = createCoreConfig(server, config, *backend, configSR, *barrierSeal, *unwrapSeal, metricsHelper, metricSink, secureRandomReader)
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	var disableClustering bool
	if err := diagnose.Test(ctx, "setup-ha-storage", func(ctx context.Context) error {
		if backend == nil {
			return fmt.Errorf(BackendUninitializedErr)
		}
		// Initialize the separate HA storage backend, if it exists
		disableClustering, err = initHaBackend(server, config, &coreConfig, *backend)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	// // Determine the redirect address from environment variables
	if err := diagnose.Test(ctx, "determine-redirect", func(ctx context.Context) error {

		err = determineRedirectAddr(server, &coreConfig, config)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return diagnose.Error(ctx, err)
	}

	if err := diagnose.Test(ctx, "find-cluster-addr", func(ctx context.Context) error {
		err = findClusterAddress(server, &coreConfig, config, disableClustering)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	// Initialize the core
	var core *vault.Core
	if err := diagnose.Test(ctx, "init-core", func(ctx context.Context) error {
		// Initialize the core
		var newCoreError error
		if coreConfig.RawConfig == nil {
			return fmt.Errorf(CoreConfigUninitializedErr)
		}
		core, newCoreError = vault.NewCore(&coreConfig)
		if newCoreError != nil {
			if vault.IsFatalError(newCoreError) {
				return fmt.Errorf("Error initializing core: %s", newCoreError)
			}
			diagnose.Warn(ctx, wrapAtLength(
				"WARNING! A non-fatal error occurred during initialization. Please "+
					"check the logs for more information."))
		}
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	if coreConfig.ReloadFuncs != nil && coreConfig.ReloadFuncsLock != nil {
		// Copy the reload funcs pointers back
		server.reloadFuncs = coreConfig.ReloadFuncs
		server.reloadFuncsLock = coreConfig.ReloadFuncsLock
	}

	var clusterAddrs []*net.TCPAddr
	var lns []listenerutil.Listener
	if err := diagnose.Test(ctx, "init-listeners", func(ctx context.Context) error {
		disableClustering := config.HAStorage.DisableClustering
		infoKeys := make([]string, 0, 10)
		info := make(map[string]string)
		status, listeners, clAddrs, errMsg := server.InitListeners(config, disableClustering, &infoKeys, &info)
		if status != 0 {
			return errMsg
		}

		lns = listeners
		clusterAddrs = clAddrs

		// Make sure we close all listeners from this point on
		listenerCloseFunc := func() {
			for _, ln := range lns {
				ln.Listener.Close()
			}
		}

		defer c.cleanupGuard.Do(listenerCloseFunc)

		sanitizedListeners := make([]listenerutil.Listener, 0, len(config.Listeners))
		for _, ln := range lns {
			if ln.Config.TLSDisable {
				diagnose.Warn(ctx, "TLS is disabled in a Listener config stanza.")
				continue
			}
			if ln.Config.TLSDisableClientCerts {
				diagnose.Warn(ctx, "TLS for a listener is turned on without requiring client certs.")
			}

			// Check ciphersuite and load ca/cert/key files
			// TODO: TLSConfig returns a reloadFunc and a TLSConfig. We can use this to
			// perform an active probe.
			_, _, err := listenerutil.TLSConfig(ln.Config, make(map[string]string), c.UI)
			if err != nil {
				return err
			}

			sanitizedListeners = append(sanitizedListeners, listenerutil.Listener{
				Listener: ln.Listener,
				Config:   ln.Config,
			})
		}
		return diagnose.ListenerChecks(sanitizedListeners)
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	if core != nil {
		// This needs to happen before we first unseal, so before we trigger dev
		// mode if it's set
		core.SetClusterListenerAddrs(clusterAddrs)
		core.SetClusterHandler(vaulthttp.Handler(&vault.HandlerProperties{
			Core: core,
		}))
	}

	// // TODO: Diagnose logging configuration

	if err := diagnose.Test(ctx, "unseal", func(ctx context.Context) error {
		if core != nil {
			runUnseal(server, core)
		} else {
			return fmt.Errorf(CoreUninitializedErr)
		}
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	// If service discovery is available, run service discovery
	if err := diagnose.Test(ctx, "run-listeners", func(ctx context.Context) error {

		// Instantiate the wait group
		server.WaitGroup = &sync.WaitGroup{}

		err = runListeners(server, &coreConfig, config, configSR)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}

	if err := diagnose.Test(ctx, "start-servers", func(ctx context.Context) error {
		// Initialize the HTTP servers
		if core != nil {
			err = startHttpServers(server, core, config, lns)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf(CoreUninitializedErr)
		}
		return nil
	}); err != nil {
		diagnose.Error(ctx, err)
	}
	return nil
}
