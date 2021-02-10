package command

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const OperatorDiagnoseEnableEnv = "VAULT_DIAGNOSE"
const OperatorDiagnoseTraceEnv = "VAULT_DIAGNOSE_TRACE"

var traceDiagnose = false

var _ cli.Command = (*OperatorDiagnoseCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorDiagnoseCommand)(nil)

type OperatorDiagnoseCommand struct {
	*BaseCommand

	flagDebug   bool
	flagSkips   []string
	flagConfigs []string
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

const status_unknown = "[      ] "
const status_ok = "\u001b[32m[  ok  ]\u001b[0m "
const status_failed = "\u001b[31m[failed]\u001b[0m "
const status_warn = "\u001b[33m[ warn ]\u001b[0m "
const same_line = "\u001b[F"

func (c *OperatorDiagnoseCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	return c.RunWithParsedFlags()
}

// This class records any errors issued during execution of the
// server command.
type StartupObserver struct {
	// Placeholder for the whole tree structure
	// Status reports the step status, Errors contains the original error
	Status map[string]string
	Cause  map[string]string
	Config *server.Config
	lock   sync.RWMutex
}

func (o *StartupObserver) Success(key string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Status[key] = status_ok

	if traceDiagnose {
		// TODO: use the UI output
		fmt.Printf("key %v success\n", key)
	}
}

func (o *StartupObserver) Error(key string, err error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if traceDiagnose {
		fmt.Printf("key %v error %v\n", key, err)
	}

	switch key {
	case "config-absent":
		// Failure but no error message provided
		o.Status["config"] = status_failed
		o.Cause["config"] = "No configuration files found."
		// TODO: copy the good explanation already present!

	case "config-cluster-address",
		"config-clustering",
		"config-core",
		"config-ha-storage",
		"config-logging-1",
		"config-logging-2",
		"config-mlock",
		"config-parse-error",
		"config-pid",
		"config-random-reader",
		"config-registration-no-ha",
		"config-service-registration",
		"config-telemetry",
		"config-vault-ui":
		o.Status["config"] = status_failed
		o.Cause["config"] = err.Error()
	case "config-core-nonfatal":
		o.Status["config"] = status_warn
		o.Cause["config"] = err.Error()
	case "listener",
		"listener-cluster-address":
		o.Status["listener"] = status_failed
		o.Cause["listener"] = err.Error()
	case "service-registration":
		o.Status["listener"] = status_failed
		o.Cause["listener"] = fmt.Sprintf("Error running service_registration: %s", err.Error())
	case "storage",
		"storage-ha",
		"storage-ha-disabled",
		"storage-ha-unsupported",
		"storage-migration",
		"storage-raft-retry-join":
		o.Status["storage"] = status_failed
		o.Cause["storage"] = err.Error()
	case "unseal",
		"unseal-barrier",
		"unseal-config",
		"unseal-load",
		"unseal-uninitialized":
		o.Status["unseal"] = status_failed
		o.Cause["unseal"] = err.Error()
	}
}

func (o *StartupObserver) ConfigCreated(config *server.Config) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.Config = config
}

func (o *StartupObserver) IsEnabled() bool {
	return true
}

// Should we delay and give auto-unseal a chance to run?
func (o *StartupObserver) WaitForUnseal() bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if o.Config == nil {
		return true
	}

	_, ok := o.Status["unseal"]
	if ok {
		// Got a status already, we're done
		return false
	}

	// Is there an auto-unseal configured, or only shamir?
	for _, seal := range o.Config.Seals {
		if !seal.Disabled && seal.Type != wrapping.Shamir {
			return true
		}
	}
	return false
}

type NullUI struct {
}

func (n *NullUI) Ask(_ string) (string, error) {
	return "", fmt.Errorf("Ask is unimplemented")
}

func (n *NullUI) AskSecret(_ string) (string, error) {
	return "", fmt.Errorf("AskSecret is unimplemented")
}

func (n *NullUI) Output(_ string) {
}

func (n *NullUI) Info(_ string) {
}

func (n *NullUI) Error(_ string) {
}

func (n *NullUI) Warn(_ string) {
}

func (c *OperatorDiagnoseCommand) RunWithParsedFlags() int {
	if len(c.flagConfigs) == 0 {
		c.UI.Error("Must specify a configuration file using -config.")
		return 1
	}

	if os.Getenv(OperatorDiagnoseTraceEnv) != "" {
		traceDiagnose = true
	}

	c.UI.Output(version.GetVersion().FullVersionNumber(true))

	shutdownCh := make(chan struct{})
	startedCh := make(chan struct{})
	exitedCh := make(chan int)

	nullUI := &VaultUI{
		Ui:     &NullUI{},
		format: "json",
	}
	server := &ServerCommand{
		// Copy of the base command, but with a UI that won't output
		BaseCommand: &BaseCommand{
			UI:          nullUI,
			tokenHelper: c.tokenHelper,
			flagAddress: c.flagAddress,
			client:      c.client,
		},

		// TODO: refactor to a common place?
		AuditBackends:        auditBackends,
		CredentialBackends:   credentialBackends,
		LogicalBackends:      logicalBackends,
		PhysicalBackends:     physicalBackends,
		ServiceRegistrations: serviceRegistrations,

		ShutdownCh: shutdownCh,
		startedCh:  startedCh,

		logger:     log.NewInterceptLogger(nil),
		allLoggers: []log.Logger{},
	}

	// Ensure commands get all their default values
	server.Run([]string{})

	// Override config flags
	server.flagConfigs = c.flagConfigs

	// Supress log messages (TODO: can we do this entirely?)
	server.flagLogLevel = "error"

	observer := &StartupObserver{
		Status: make(map[string]string),
		Cause:  make(map[string]string),
	}

	// We'll run the command in its own goroutine.
	// If it returns we're done.  If we get startedCh then it was successful.
	// in that case we may want to delay until UnsealWithStoredKeys is called.

	go func() {
		ret := server.RunWithObserver(observer)
		exitedCh <- ret
	}()

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	vaultExitCode := 0
	diagnoseExitCode := 0
	started := false
	monitoringServer := true

	for monitoringServer {
		select {
		case <-timeout:
			if started {
				c.UI.Error("Vault shutdown timed out.")
			} else {
				close(shutdownCh)
				c.UI.Error("Vault startup timed out.")
			}
			diagnoseExitCode = 1
			monitoringServer = false
		case vaultExitCode = <-exitedCh:
			monitoringServer = false
		case <-startedCh:
			// Successful start, keep waiting?
			started = true
			if !observer.WaitForUnseal() {
				close(shutdownCh)
				// Wait for exit code
			}
		case <-ticker.C:
			if started && !observer.WaitForUnseal() {
				close(shutdownCh)
				// Wait for exit code
			}
		}
	}
	ticker.Stop()

	// Report status
	// TODO: move to another function, replace by tree structure
	// that generates templated output
	observer.lock.RLock()
	defer observer.lock.RUnlock()

	switch observer.Status["config"] {
	case "":
		c.UI.Output(status_failed + "Unknown error caused configuration failure")
		return 1
	case status_warn:
		c.UI.Output(status_warn + "Parsed configuration")
		c.UI.Output(observer.Cause["config"])
	case status_ok:
		c.UI.Output(status_ok + "Parsed configuration")
	case status_failed:
		c.UI.Output(status_failed + "Configuration error")
		c.UI.Output(observer.Cause["config"])
		return 1
	}

	storage_type := observer.Config.Storage.Type
	switch observer.Status["storage"] {
	case "":
		c.UI.Output(status_failed + "Unknown error caused storage initialization failure")
		return 1
	case status_warn:
		c.UI.Output(status_warn + "Started " + storage_type + " storage")
		c.UI.Output(observer.Cause["storage"])
	case status_ok:
		c.UI.Output(status_ok + "Started " + storage_type + " storage")
	case status_failed:
		c.UI.Output(status_failed + "Storage error initializing " + storage_type)
		c.UI.Output(observer.Cause["storage"])
		return 1
	}

	switch observer.Status["listener"] {
	case "":
		// If we had an unseal failure first, go on to report that
		if observer.Status["unseal"] == "" {
			c.UI.Output(status_failed + "Unknown error caused listener initialization failure")
			return 1
		}
	case status_warn:
		c.UI.Output(status_warn + "Started listeners")
		c.UI.Output(observer.Cause["listener"])
	case status_ok:
		c.UI.Output(status_ok + "Started listeners")
	case status_failed:
		c.UI.Output(status_failed + "Error initializing listeners")
		c.UI.Output(observer.Cause["listener"])
		return 1
	}

	seal_type := wrapping.Shamir
	for _, seal := range observer.Config.Seals {
		if !seal.Disabled && seal.Type != wrapping.Shamir {
			seal_type = seal.Type
			break
		}
	}

	if seal_type == wrapping.Shamir {
		c.UI.Output(status_unknown + "Shamir unseal requires user input")
	} else {
		switch observer.Status["unseal"] {
		case "":
			c.UI.Output(status_failed + "Unknown error caused unseal failure")
			return 1
		case status_warn:
			c.UI.Output(status_warn + "Auto-unseal using " + seal_type + " succeeded with warnings")
			c.UI.Output(observer.Cause["listener"])
		case status_ok:
			c.UI.Output(status_ok + "Auto-unseal using " + seal_type + " succeeded")
		case status_failed:
			c.UI.Output(status_failed + "Error during " + seal_type + " unseal")
			c.UI.Output(observer.Cause["unseal"])
			return 1
		}
	}

	if vaultExitCode != 0 {
		c.UI.Output(fmt.Sprintf("%vVault exited with status code %v due to an undiagnosed error", status_failed, vaultExitCode))
	} else {
		c.UI.Output("Vault Diagnose could not detect any problems during startup.")
	}
	return diagnoseExitCode
}
