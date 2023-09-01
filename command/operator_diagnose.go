// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-kms-wrapping/entropy/v2"

	"golang.org/x/term"

	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	uuid "github.com/hashicorp/go-uuid"
	cserver "github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	physconsul "github.com/hashicorp/vault/physical/consul"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/physical"
	sr "github.com/hashicorp/vault/serviceregistration"
	srconsul "github.com/hashicorp/vault/serviceregistration/consul"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/diagnose"
	"github.com/hashicorp/vault/vault/hcp_link"
	"github.com/hashicorp/vault/version"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const CoreConfigUninitializedErr = "Diagnose cannot attempt this step because core config could not be set."

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

	f.StringVar(&StringVar{
		Name:   "format",
		Target: &c.flagFormat,
		Usage:  "The output format",
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
		return 3
	}
	return c.RunWithParsedFlags()
}

func (c *OperatorDiagnoseCommand) RunWithParsedFlags() int {
	if len(c.flagConfigs) == 0 {
		c.UI.Error("Must specify a configuration file using -config.")
		return 3
	}

	if c.diagnose == nil {
		if c.flagFormat == "json" {
			c.diagnose = diagnose.New(io.Discard)
		} else {
			c.UI.Output(version.GetVersion().FullVersionNumber(true))
			c.diagnose = diagnose.New(os.Stdout)
		}
	}
	ctx := diagnose.Context(context.Background(), c.diagnose)
	c.diagnose.SkipFilters = c.flagSkips
	err := c.offlineDiagnostics(ctx)

	results := c.diagnose.Finalize(ctx)
	if c.flagFormat == "json" {
		resultsJS, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling results: %v.", err)
			return 4
		}
		c.UI.Output(string(resultsJS))
	} else {
		c.UI.Output("\nResults:")
		w, _, err := term.GetSize(0)
		if err == nil {
			results.Write(os.Stdout, w)
		} else {
			results.Write(os.Stdout, 0)
		}
	}

	if err != nil {
		return 4
	}
	// Use a different return code
	switch results.Status {
	case diagnose.WarningStatus:
		return 2
	case diagnose.ErrorStatus:
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

		logger: log.NewInterceptLogger(&log.LoggerOptions{
			Level: log.Off,
		}),
		allLoggers:      []log.Logger{},
		reloadFuncs:     &rloadFuncs,
		reloadFuncsLock: new(sync.RWMutex),
	}

	ctx, span := diagnose.StartSpan(ctx, "Vault Diagnose")
	defer span.End()

	// OS Specific checks
	diagnose.OSChecks(ctx)

	var config *cserver.Config

	diagnose.Test(ctx, "Parse Configuration", func(ctx context.Context) (err error) {
		server.flagConfigs = c.flagConfigs
		var configErrors []configutil.ConfigError
		config, configErrors, err = server.parseConfig()
		if err != nil {
			return fmt.Errorf("Could not parse configuration: %w.", err)
		}
		for _, ce := range configErrors {
			diagnose.Warn(ctx, diagnose.CapitalizeFirstLetter(ce.String())+".")
		}
		diagnose.Success(ctx, "Vault configuration syntax is ok.")
		return nil
	})
	if config == nil {
		return fmt.Errorf("No vault server configuration found.")
	}

	diagnose.Test(ctx, "Check Telemetry", func(ctx context.Context) (err error) {
		if config.Telemetry == nil {
			diagnose.Warn(ctx, "Telemetry is using default configuration")
			diagnose.Advise(ctx, "By default only Prometheus and JSON metrics are available.  Ignore this warning if you are using telemetry or are using these metrics and are satisfied with the default retention time and gauge period.")
		} else {
			t := config.Telemetry
			// If any Circonus setting is present but we're missing the basic fields...
			if coalesce(t.CirconusAPIURL, t.CirconusAPIToken, t.CirconusCheckID, t.CirconusCheckTags, t.CirconusCheckSearchTag,
				t.CirconusBrokerID, t.CirconusBrokerSelectTag, t.CirconusCheckForceMetricActivation, t.CirconusCheckInstanceID,
				t.CirconusCheckSubmissionURL, t.CirconusCheckDisplayName) != nil {
				if t.CirconusAPIURL == "" {
					return errors.New("incomplete Circonus telemetry configuration, missing circonus_api_url")
				} else if t.CirconusAPIToken != "" {
					return errors.New("incomplete Circonus telemetry configuration, missing circonus_api_token")
				}
			}
			if len(t.DogStatsDTags) > 0 && t.DogStatsDAddr == "" {
				return errors.New("incomplete DogStatsD telemetry configuration, missing dogstatsd_addr, while dogstatsd_tags specified")
			}

			// If any Stackdriver setting is present but we're missing the basic fields...
			if coalesce(t.StackdriverNamespace, t.StackdriverLocation, t.StackdriverDebugLogs, t.StackdriverNamespace) != nil {
				if t.StackdriverProjectID == "" {
					return errors.New("incomplete Stackdriver telemetry configuration, missing stackdriver_project_id")
				}
				if t.StackdriverLocation == "" {
					return errors.New("incomplete Stackdriver telemetry configuration, missing stackdriver_location")
				}
				if t.StackdriverNamespace == "" {
					return errors.New("incomplete Stackdriver telemetry configuration, missing stackdriver_namespace")
				}
			}
		}
		return nil
	})

	var metricSink *metricsutil.ClusterMetricSink
	var metricsHelper *metricsutil.MetricsHelper

	var backend *physical.Backend
	diagnose.Test(ctx, "Check Storage", func(ctx context.Context) error {
		// Ensure that there is a storage stanza
		if config.Storage == nil {
			diagnose.Advise(ctx, "To learn how to specify a storage backend, see the Vault server configuration documentation.")
			return fmt.Errorf("No storage stanza in Vault server configuration.")
		}

		diagnose.Test(ctx, "Create Storage Backend", func(ctx context.Context) error {
			b, err := server.setupStorage(config)
			if err != nil {
				return err
			}
			if b == nil {
				diagnose.Advise(ctx, "To learn how to specify a storage backend, see the Vault server configuration documentation.")
				return fmt.Errorf("Storage backend could not be initialized.")
			}
			backend = &b
			return nil
		})

		if backend == nil {
			diagnose.Fail(ctx, "Diagnose could not initialize storage backend.")
			span.End()
			return fmt.Errorf("Diagnose could not initialize storage backend.")
		}

		// Check for raft quorum status
		if config.Storage.Type == storageTypeRaft {
			path := os.Getenv(raft.EnvVaultRaftPath)
			if path == "" {
				path, ok := config.Storage.Config["path"]
				if !ok {
					diagnose.SpotError(ctx, "Check Raft Folder Permissions", fmt.Errorf("Storage folder path is required."))
				}
				diagnose.RaftFileChecks(ctx, path)
			}
			diagnose.RaftStorageQuorum(ctx, (*backend).(*raft.RaftBackend))
		}

		// Consul storage checks
		if config.Storage != nil && config.Storage.Type == storageTypeConsul {
			diagnose.Test(ctx, "Check Consul TLS", func(ctx context.Context) error {
				err := physconsul.SetupSecureTLS(ctx, api.DefaultConfig(), config.Storage.Config, server.logger, true)
				if err != nil {
					return err
				}
				return nil
			})

			diagnose.Test(ctx, "Check Consul Direct Storage Access", func(ctx context.Context) error {
				dirAccess := diagnose.ConsulDirectAccess(config.Storage.Config)
				if dirAccess != "" {
					diagnose.Warn(ctx, dirAccess)
				}
				if dirAccess == diagnose.DirAccessErr {
					diagnose.Advise(ctx, diagnose.DirAccessAdvice)
				}
				return nil
			})
		}

		// Attempt to use storage backend
		if !c.skipEndEnd && config.Storage.Type != storageTypeRaft {
			diagnose.Test(ctx, "Check Storage Access", diagnose.WithTimeout(30*time.Second, func(ctx context.Context) error {
				maxDurationCrudOperation := "write"
				maxDuration := time.Duration(0)
				uuidSuffix, err := uuid.GenerateUUID()
				if err != nil {
					return err
				}
				uuid := "diagnose/latency/" + uuidSuffix
				dur, err := diagnose.EndToEndLatencyCheckWrite(ctx, uuid, *backend)
				if err != nil {
					return err
				}
				maxDuration = dur
				dur, err = diagnose.EndToEndLatencyCheckRead(ctx, uuid, *backend)
				if err != nil {
					return err
				}
				if dur > maxDuration {
					maxDuration = dur
					maxDurationCrudOperation = "read"
				}
				dur, err = diagnose.EndToEndLatencyCheckDelete(ctx, uuid, *backend)
				if err != nil {
					return err
				}
				if dur > maxDuration {
					maxDuration = dur
					maxDurationCrudOperation = "delete"
				}

				if maxDuration > time.Duration(0) {
					diagnose.Warn(ctx, diagnose.LatencyWarning+fmt.Sprintf("duration: %s, operation: %s", maxDuration, maxDurationCrudOperation))
				}
				return nil
			}))
		}
		return nil
	})

	// Return from top-level span when backend is nil
	if backend == nil {
		return fmt.Errorf("Diagnose could not initialize storage backend.")
	}

	var configSR sr.ServiceRegistration
	diagnose.Test(ctx, "Check Service Discovery", func(ctx context.Context) error {
		if config.ServiceRegistration == nil || config.ServiceRegistration.Config == nil {
			diagnose.Skipped(ctx, "No service registration configured.")
			return nil
		}
		srConfig := config.ServiceRegistration.Config

		diagnose.Test(ctx, "Check Consul Service Discovery TLS", func(ctx context.Context) error {
			// SetupSecureTLS for service discovery uses the same cert and key to set up physical
			// storage. See the consul package in physical for details.
			err := srconsul.SetupSecureTLS(ctx, api.DefaultConfig(), srConfig, server.logger, true)
			if err != nil {
				return err
			}
			return nil
		})

		if config.ServiceRegistration != nil && config.ServiceRegistration.Type == "consul" {
			diagnose.Test(ctx, "Check Consul Direct Service Discovery", func(ctx context.Context) error {
				dirAccess := diagnose.ConsulDirectAccess(config.ServiceRegistration.Config)
				if dirAccess != "" {
					diagnose.Warn(ctx, dirAccess)
				}
				if dirAccess == diagnose.DirAccessErr {
					diagnose.Advise(ctx, diagnose.DirAccessAdvice)
				}
				return nil
			})
		}
		return nil
	})

	sealcontext, sealspan := diagnose.StartSpan(ctx, "Create Vault Server Configuration Seals")

	var setSealResponse *SetSealResponse
	existingSealGenerationInfo, err := vault.PhysicalSealGenInfo(sealcontext, *backend)
	if err != nil {
		diagnose.Fail(sealcontext, fmt.Sprintf("Unable to get Seal genration information from storage: %s.", err.Error()))
		goto SEALFAIL
	}

	setSealResponse, err = setSeal(server, config, make([]string, 0), make(map[string]string), existingSealGenerationInfo, false /* unsealed vault has no partially wrapped paths */)
	if err != nil {
		diagnose.Advise(ctx, "For assistance with the seal stanza, see the Vault configuration documentation.")
		diagnose.Fail(sealcontext, fmt.Sprintf("Seal creation resulted in the following error: %s.", err.Error()))
		goto SEALFAIL
	}

	for _, seal := range setSealResponse.getCreatedSeals() {
		seal := seal // capture range variable
		// Ensure that the seal finalizer is called, even if using verify-only
		defer func(seal *vault.Seal) {
			sealType := diagnose.CapitalizeFirstLetter((*seal).BarrierSealConfigType().String())
			finalizeSealContext, finalizeSealSpan := diagnose.StartSpan(ctx, "Finalize "+sealType+" Seal")
			err = (*seal).Finalize(finalizeSealContext)
			if err != nil {
				diagnose.Fail(finalizeSealContext, "Error finalizing seal.")
				diagnose.Advise(finalizeSealContext, "This likely means that the barrier is still in use; therefore, finalizing the seal timed out.")
				finalizeSealSpan.End()
			}
			finalizeSealSpan.End()
		}(seal)
	}

	if setSealResponse.sealConfigError != nil {
		diagnose.Fail(sealcontext, "Seal could not be configured: seals may already be initialized.")
	} else if setSealResponse.barrierSeal == nil {
		diagnose.Fail(sealcontext, "Could not create barrier seal. No error was generated, but it is likely that the seal stanza is misconfigured. For guidance, see Vault's configuration documentation on the seal stanza.")
	}

SEALFAIL:
	sealspan.End()

	var barrierSeal vault.Seal
	var unwrapSeal vault.Seal

	if setSealResponse != nil {
		barrierSeal = setSealResponse.barrierSeal
		unwrapSeal = setSealResponse.unwrapSeal
	}

	diagnose.Test(ctx, "Check Transit Seal TLS", func(ctx context.Context) error {
		var checkSealTransit bool
		for _, seal := range config.Seals {
			if seal.Type == "transit" {
				checkSealTransit = true

				tlsSkipVerify, _ := seal.Config["tls_skip_verify"]
				if tlsSkipVerify == "true" {
					diagnose.Warn(ctx, "TLS verification is skipped. This is highly discouraged and decreases the security of data transmissions to and from the Vault server.")
					return nil
				}

				// Checking tls_client_cert and tls_client_key
				tlsClientCert, ok := seal.Config["tls_client_cert"]
				if !ok {
					diagnose.Warn(ctx, "Missing tls_client_cert in the seal configuration.")
					return nil
				}
				tlsClientKey, ok := seal.Config["tls_client_key"]
				if !ok {
					diagnose.Warn(ctx, "Missing tls_client_key in the seal configuration.")
					return nil
				}
				_, err := diagnose.TLSFileChecks(tlsClientCert, tlsClientKey)
				if err != nil {
					return fmt.Errorf("The TLS certificate and key configured through the tls_client_cert and tls_client_key fields of the transit seal configuration are invalid: %w.", err)
				}

				// checking tls_ca_cert
				tlsCACert, ok := seal.Config["tls_ca_cert"]
				if !ok {
					diagnose.Warn(ctx, "Missing tls_ca_cert in the seal configuration.")
					return nil
				}
				warnings, err := diagnose.TLSCAFileCheck(tlsCACert)
				if len(warnings) != 0 {
					for _, warning := range warnings {
						diagnose.Warn(ctx, warning)
					}
				}
				if err != nil {
					return fmt.Errorf("The TLS CA certificate configured through the tls_ca_cert field of the transit seal configuration is invalid: %w.", err)
				}
			}
		}
		if !checkSealTransit {
			diagnose.Skipped(ctx, "No transit seal found in seal configuration.")
		}
		return nil
	})

	var coreConfig vault.CoreConfig
	diagnose.Test(ctx, "Create Core Configuration", func(ctx context.Context) error {
		var secureRandomReader io.Reader
		// prepare a secure random reader for core
		randReaderTestName := "Initialize Randomness for Core"
		var sources []*configutil.EntropySourcerInfo
		if barrierSeal != nil {
			for _, sealWrapper := range barrierSeal.GetAccess().GetEnabledSealWrappersByPriority() {
				if s, ok := sealWrapper.Wrapper.(entropy.Sourcer); ok {
					sources = append(sources, &configutil.EntropySourcerInfo{
						Sourcer: s,
						Name:    sealWrapper.Name,
					})
				}
			}
		}
		secureRandomReader, err = configutil.CreateSecureRandomReaderFunc(config.SharedConfig, sources, server.logger)
		if err != nil {
			return diagnose.SpotError(ctx, randReaderTestName, fmt.Errorf("could not initialize randomness for core: %w", err))
		}
		diagnose.SpotOk(ctx, randReaderTestName, "")
		coreConfig = createCoreConfig(server, config, *backend, configSR, barrierSeal, unwrapSeal, metricsHelper, metricSink, secureRandomReader)
		return nil
	})

	var disableClustering bool
	diagnose.Test(ctx, "HA Storage", func(ctx context.Context) error {
		diagnose.Test(ctx, "Create HA Storage Backend", func(ctx context.Context) error {
			// Initialize the separate HA storage backend, if it exists
			disableClustering, err = initHaBackend(server, config, &coreConfig, *backend)
			if err != nil {
				return err
			}
			return nil
		})

		diagnose.Test(ctx, "Check HA Consul Direct Storage Access", func(ctx context.Context) error {
			if config.HAStorage == nil {
				diagnose.Skipped(ctx, "No HA storage stanza is configured.")
			} else {
				dirAccess := diagnose.ConsulDirectAccess(config.HAStorage.Config)
				if dirAccess != "" {
					diagnose.Warn(ctx, dirAccess)
				}
				if dirAccess == diagnose.DirAccessErr {
					diagnose.Advise(ctx, diagnose.DirAccessAdvice)
				}
			}
			return nil
		})
		if config.HAStorage != nil && config.HAStorage.Type == storageTypeConsul {
			diagnose.Test(ctx, "Check Consul TLS", func(ctx context.Context) error {
				err = physconsul.SetupSecureTLS(ctx, api.DefaultConfig(), config.HAStorage.Config, server.logger, true)
				if err != nil {
					return err
				}
				return nil
			})
		}
		return nil
	})

	// Determine the redirect address from environment variables
	err = determineRedirectAddr(server, &coreConfig, config)
	if err != nil {
		return diagnose.SpotError(ctx, "Determine Redirect Address", fmt.Errorf("Redirect Address could not be determined: %w.", err))
	}
	diagnose.SpotOk(ctx, "Determine Redirect Address", "")

	err = findClusterAddress(server, &coreConfig, config, disableClustering)
	if err != nil {
		return diagnose.SpotError(ctx, "Check Cluster Address", fmt.Errorf("Cluster Address could not be determined or was invalid: %w.", err),
			diagnose.Advice("Please check that the API and Cluster addresses are different, and that the API, Cluster and Redirect addresses have both a host and port."))
	}
	diagnose.SpotOk(ctx, "Check Cluster Address", "Cluster address is logically valid and can be found.")

	var vaultCore *vault.Core

	// Run all the checks that are utilized when initializing a core object
	// without actually calling core.Init. These are in the init-core section
	// as they are runtime checks.
	diagnose.Test(ctx, "Check Core Creation", func(ctx context.Context) error {
		var newCoreError error
		if coreConfig.RawConfig == nil {
			return fmt.Errorf(CoreConfigUninitializedErr)
		}
		core, newCoreError := vault.CreateCore(&coreConfig)
		if newCoreError != nil {
			if vault.IsFatalError(newCoreError) {
				return fmt.Errorf("Error initializing core: %s.", newCoreError)
			}
			diagnose.Warn(ctx, wrapAtLength(
				"A non-fatal error occurred during initialization. Please check the logs for more information."))
		} else {
			vaultCore = core
		}
		return nil
	})

	if vaultCore == nil {
		return fmt.Errorf("Diagnose could not initialize the Vault core from the Vault server configuration.")
	}

	licenseCtx, licenseSpan := diagnose.StartSpan(ctx, "Check For Autoloaded License")
	// If we are not in enterprise, return from the check
	if !constants.IsEnterprise {
		diagnose.Skipped(licenseCtx, "License check will not run on OSS Vault.")
	} else {
		// Load License from environment variables. These take precedence over the
		// configured license.
		if envLicensePath := os.Getenv(EnvVaultLicensePath); envLicensePath != "" {
			coreConfig.LicensePath = envLicensePath
		}
		if envLicense := os.Getenv(EnvVaultLicense); envLicense != "" {
			coreConfig.License = envLicense
		}
		vault.DiagnoseCheckLicense(licenseCtx, vaultCore, coreConfig, false)
	}
	licenseSpan.End()

	var lns []listenerutil.Listener
	diagnose.Test(ctx, "Start Listeners", func(ctx context.Context) error {
		disableClustering := config.HAStorage != nil && config.HAStorage.DisableClustering
		infoKeys := make([]string, 0, 10)
		info := make(map[string]string)
		var listeners []listenerutil.Listener
		var status int

		diagnose.ListenerChecks(ctx, config.Listeners)

		diagnose.Test(ctx, "Create Listeners", func(ctx context.Context) error {
			status, listeners, _, err = server.InitListeners(config, disableClustering, &infoKeys, &info)
			if status != 0 {
				return err
			}
			return nil
		})

		lns = listeners

		// Make sure we close all listeners from this point on
		listenerCloseFunc := func() {
			for _, ln := range lns {
				ln.Listener.Close()
			}
		}

		c.cleanupGuard.Do(listenerCloseFunc)

		return nil
	})

	// TODO: Diagnose logging configuration

	// The unseal diagnose check will simply attempt to use the barrier to encrypt and
	// decrypt a mock value. It will not call runUnseal.
	diagnose.Test(ctx, "Check Autounseal Encryption", diagnose.WithTimeout(30*time.Second, func(ctx context.Context) error {
		if barrierSeal == nil {
			return fmt.Errorf("Diagnose could not create a barrier seal object.")
		}
		if barrierSeal.BarrierSealConfigType() == vault.SealConfigTypeShamir {
			diagnose.Skipped(ctx, "Skipping barrier encryption test. Only supported for auto-unseal.")
			return nil
		}
		barrierUUID, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Diagnose could not create unique UUID for unsealing.")
		}
		barrierEncValue := "diagnose-" + barrierUUID
		ciphertext, errMap := barrierSeal.GetAccess().Encrypt(ctx, []byte(barrierEncValue), nil)
		if len(errMap) > 0 {
			var sealErrors []error
			for name, err := range errMap {
				sealErrors = append(sealErrors, fmt.Errorf("error encrypting with seal %q: %w", name, err))
			}
			if ciphertext == nil {
				// Full failure
				if len(sealErrors) == 1 {
					return sealErrors[0]
				} else {
					return fmt.Errorf("complete seal encryption failure: %w", errors.Join())
				}
			} else {
				// Partial failure
				return fmt.Errorf("partial seal encryption failure: %w", errors.Join())
			}
		}
		plaintext, _, err := barrierSeal.GetAccess().Decrypt(ctx, ciphertext, nil)
		if err != nil {
			return fmt.Errorf("Error decrypting with seal barrier: %w", err)
		}
		if string(plaintext) != barrierEncValue {
			return fmt.Errorf("Barrier returned incorrect decrypted value for mock data.")
		}
		return nil
	}))

	// The following block contains static checks that are run during the
	// startHttpServers portion of server run. In other words, they are static
	// checks during resource creation. Currently there is nothing important in this
	// diagnose check. For now it is a placeholder for any checks that will be done
	// before server run.
	diagnose.Test(ctx, "Check Server Before Runtime", func(ctx context.Context) error {
		for _, ln := range lns {
			if ln.Config == nil {
				return fmt.Errorf("Found no listener config after parsing the Vault configuration.")
			}
		}
		return nil
	})

	// Checking HCP link to make sure Vault could connect to SCADA.
	// If it could not connect to SCADA in 5 seconds, diagnose reports an issue
	if !constants.IsEnterprise {
		diagnose.Skipped(ctx, "HCP link check will not run on OSS Vault.")
	} else {
		if config.HCPLinkConf != nil {
			// we need to override API and Passthrough capabilities
			// as they could not be initialized when Vault http handler
			// is not fully initialized
			config.HCPLinkConf.EnablePassThroughCapability = false
			config.HCPLinkConf.EnableAPICapability = false

			diagnose.Test(ctx, "Check HCP Connection", func(ctx context.Context) error {
				hcpLink, err := hcp_link.NewHCPLink(config.HCPLinkConf, vaultCore, server.logger)
				if err != nil || hcpLink == nil {
					return fmt.Errorf("failed to start HCP link, %w", err)
				}

				// check if a SCADA session is established successfully
				deadline := time.Now().Add(5 * time.Second)
				linkSessionStatus := "disconnected"
				for time.Now().Before(deadline) {
					linkSessionStatus = hcpLink.GetConnectionStatusMessage(hcpLink.GetScadaSessionStatus())
					if linkSessionStatus == "connected" {
						break
					}
					time.Sleep(500 * time.Millisecond)
				}
				if linkSessionStatus != "connected" {
					return fmt.Errorf("failed to connect to HCP in 5 seconds. HCP session status is: %s", linkSessionStatus)
				}

				err = hcpLink.Shutdown()
				if err != nil {
					return fmt.Errorf("failed to shutdown HCP link: %w", err)
				}

				return nil
			})
		}
	}

	return nil
}

func coalesce(values ...interface{}) interface{} {
	for _, val := range values {
		if val != nil && val != "" {
			return val
		}
	}
	return nil
}
