package command

import (
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/hashicorp/vault/vault/diagnose"
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

func (c *OperatorDiagnoseCommand) RunWithParsedFlags() int {
	if len(c.flagConfigs) == 0 {
		c.UI.Error("Must specify a configuration file using -config.")
		return 1
	}

	c.UI.Output(version.GetVersion().FullVersionNumber(true))

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

		logger:     log.NewInterceptLogger(nil),
		allLoggers: []log.Logger{},
	}

	phase := "Parse configuration"
	c.UI.Output(status_unknown + phase)
	server.flagConfigs = c.flagConfigs
	config, err := server.parseConfig()
	if err != nil {
		c.UI.Output(same_line + status_failed + phase)
		c.UI.Output("Error while reading configuration files:")
		c.UI.Output(err.Error())
		return 1
	}

	// Check Listener Information
	// TODO: Run Diagnose checks on the actual net.Listeners

	disableClustering := config.HAStorage.DisableClustering
	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	status, lns, _ := server.InitListeners(config, disableClustering, &infoKeys, &info)

	if status != 0 {
		// The error should be displayed in InitListeners. We don't need the error message here.
		c.UI.Output("Error parsing listener configuration.")
		return 1
	}

	sanitizedListeners := make([]listenerutil.Listener, 0, len(config.Listeners))
	for _, ln := range lns {
		if ln.Config.TLSDisable {
			c.UI.Warn("WARNING! TLS is disabled in a Listener config stanza.")
			continue
		}
		if ln.Config.TLSDisableClientCerts {
			c.UI.Warn("WARNING! TLS for a listener is turned on without requiring client certs.")
		}

		// Check ciphersuite and load ca/cert/key files
		// TODO: TLSConfig returns a reloadFunc and a TLSConfig. We can use this to
		// perform an active probe.
		_, _, err := listenerutil.TLSConfig(ln.Config, make(map[string]string), c.UI)
		if err != nil {
			c.UI.Output("Error creating TLS Configuration out of config file: ")
			c.UI.Output(err.Error())
			return 1
		}

		sanitizedListeners = append(sanitizedListeners, listenerutil.Listener{
			Listener: ln.Listener,
			Config:   ln.Config,
		})
	}
	err = diagnose.TLSConfigChecks(sanitizedListeners)
	if err != nil {
		c.UI.Output("Diagnose caught configuration errors: ")
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
