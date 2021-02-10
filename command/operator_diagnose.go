package command

import (
	"fmt"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const OperatorDiagnoseEnableEnv = "VAULT_DIAGNOSE"

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
}

func (n *StartupObserver) Exited(status int) {
}

func (n *StartupObserver) Success(key string) {
	n.Status[key] = status_ok
}

func (n *StartupObserver) Error(key string, err error) {
	switch key {
	case "config-absent",
		"config-cluster-address",
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
		n.Status["config"] = status_failed
		n.Cause["config"] = err.Error()
	case "config-core-nonfatal":
		n.Status["config"] = status_warn
		n.Cause["config"] = err.Error()
	case "listener",
		"listener-cluster-address":
		n.Status["listener"] = status_failed
		n.Cause["listener"] = err.Error()
	case "service-registration":
		n.Status["listener"] = status_failed
		n.Cause["listener"] = fmt.Sprintf("Error running service_registration: %s", err.Error())
	case "storage",
		"storage-ha",
		"storage-ha-disabled",
		"storage-ha-unsupported",
		"storage-migration",
		"storage-raft-retry-join":
		n.Status["storage"] = status_failed
		n.Cause["storage"] = err.Error()
	case "unseal",
		"unseal-barrier",
		"unseal-config",
		"unseal-load",
		"unseal-uninitialized":
		n.Status["unseal"] = status_failed
		n.Cause["unseal"] = err.Error()
	}
}

func (n *StartupObserver) IsEnabled() bool {
	return true
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

	c.UI.Output(version.GetVersion().FullVersionNumber(true))

	shutdownCh := make(chan struct{})
	startedCh := make(chan struct{})

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

	// We'll run the command in its own goroutine.
	// If it returns we're done.  If we get startedCh then it was successful.
	// in that case we may want to delay until UnsealWithStoredKeys is called.

	phase := "Parse configuration"
	c.UI.Output(status_unknown + phase)
	server.flagConfigs = c.flagConfigs
	config, err := server.parseConfig()
	if err != nil {
		// TODO: show every file one by one?
		c.UI.Output(same_line + status_failed + phase)
		c.UI.Output("Error while reading configuration files:")
		c.UI.Output(err.Error())
		return 1
	}

	// Errors in these items could stop Vault from starting but are not yet covered:
	// TODO: logging configuration
	// TODO: SetupTelemetry
	// TODO: check for storage backend
	c.UI.Output(same_line + status_ok + phase)

	phase = "Access storage"
	c.UI.Output(status_unknown + phase)
	_, err = server.setupStorage(config)
	if err != nil {
		c.UI.Output(same_line + status_failed + phase)
		c.UI.Output(err.Error())
		return 1
	}
	c.UI.Output(same_line + status_ok + phase)

	return 0
}
