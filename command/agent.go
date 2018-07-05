package command

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/kr/pretty"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/command/agent"
	"github.com/hashicorp/vault/helper/gated-writer"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/reload"
	"github.com/hashicorp/vault/version"
)

var _ cli.Command = (*AgentCommand)(nil)
var _ cli.CommandAutocomplete = (*AgentCommand)(nil)

type AgentCommand struct {
	*BaseCommand

	ShutdownCh chan struct{}
	SighupCh   chan struct{}

	WaitGroup *sync.WaitGroup

	logWriter io.Writer
	logGate   *gatedwriter.Writer
	logger    log.Logger

	cleanupGuard sync.Once

	reloadFuncsLock *sync.RWMutex
	reloadFuncs     *map[string][]reload.ReloadFunc
	startedCh       chan (struct{}) // for tests
	reloadedCh      chan (struct{}) // for tests

	flagConfigs  []string
	flagLogLevel string

	flagTestVerifyOnly bool
	flagCombineLogs    bool
}

func (c *AgentCommand) Synopsis() string {
	return "Start a Vault agent"
}

func (c *AgentCommand) Help() string {
	helpText := `
Usage: vault agent [options]

  This command starts a Vault agent that can perform automatic authentication
  in certain environments.

  Start an agent with a configuration file:

      $ vault agent -config=/etc/vault/config.hcl

  For a full list of examples, please see the documentation.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *AgentCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&StringSliceVar{
		Name:   "config",
		Target: &c.flagConfigs,
		Completion: complete.PredictOr(
			complete.PredictFiles("*.hcl"),
			complete.PredictFiles("*.json"),
		),
		Usage: "Path to a configuration file. This configuration file should " +
			"contain only agent directives.",
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

	// Internal-only flags to follow.
	//
	// Why hello there little source code reader! Welcome to the Vault source
	// code. The remaining options are intentionally undocumented and come with
	// no warranty or backwards-compatability promise. Do not use these flags
	// in production. Do not build automation using these flags. Unless you are
	// developing against Vault, you should not need any of these flags.

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

	// End internal-only flags.

	return set
}

func (c *AgentCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AgentCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AgentCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	c.logGate = &gatedwriter.Writer{Writer: os.Stderr}
	c.logWriter = c.logGate
	if c.flagCombineLogs {
		c.logWriter = os.Stdout
	}
	var level log.Level
	c.flagLogLevel = strings.ToLower(strings.TrimSpace(c.flagLogLevel))
	switch c.flagLogLevel {
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
		c.UI.Error(fmt.Sprintf("Unknown log level: %s", c.flagLogLevel))
		return 1
	}

	c.logger = logging.NewVaultLoggerWithWriter(c.logWriter, level)

	// Validation
	if len(c.flagConfigs) != 1 {
		c.UI.Error("Must specify exactly one config path using -config")
		return 1
	}

	// Load the configuration
	config, err := agent.LoadConfig(c.flagConfigs[0], c.logger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error loading configuration from %s: %s", c.flagConfigs[0], err))
		return 1
	}

	// Ensure at least one config was found.
	if config == nil {
		c.UI.Output(wrapAtLength(
			"No configuration read. Please provide the configuration with the " +
				"-config flag."))
		return 1
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = c.flagLogLevel
	infoKeys = append(infoKeys, "log level")

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
	c.UI.Output("==> Vault agent configuration:\n")
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
		if os.Getenv("VAULT_TEST_VERIFY_ONLY_DUMP_CONFIG") != "" {
			c.UI.Output(fmt.Sprintf(
				"\nConfiguration:\n%s\n",
				pretty.Sprint(*config)))
		}
		return 0
	}

	// Instantiate the wait group
	c.WaitGroup = &sync.WaitGroup{}

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
			c.UI.Output("==> Vault agent shutdown triggered")

			shutdownTriggered = true

		case <-c.SighupCh:
			c.UI.Output("==> Vault agent reload triggered")
			if err := c.Reload(c.reloadFuncsLock, c.reloadFuncs, c.flagConfigs); err != nil {
				c.UI.Error(fmt.Sprintf("Error(s) were encountered during reload: %s", err))
			}
		}
	}

	// Wait for dependent goroutines to complete
	c.WaitGroup.Wait()
	return 0
}

func (c *AgentCommand) Reload(lock *sync.RWMutex, reloadFuncs *map[string][]reload.ReloadFunc, configPath []string) error {
	lock.RLock()
	defer lock.RUnlock()

	var reloadErrors *multierror.Error

	for k, relFuncs := range *reloadFuncs {
		switch {
		case strings.HasPrefix(k, "listener|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(nil); err != nil {
						reloadErrors = multierror.Append(reloadErrors, errwrap.Wrapf("error encountered reloading listener: {{err}}", err))
					}
				}
			}

		case strings.HasPrefix(k, "audit_file|"):
			for _, relFunc := range relFuncs {
				if relFunc != nil {
					if err := relFunc(nil); err != nil {
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
func (c *AgentCommand) storePidFile(pidPath string) error {
	// Quit fast if no pidfile
	if pidPath == "" {
		return nil
	}

	// Open the PID file
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
func (c *AgentCommand) removePidFile(pidPath string) error {
	if pidPath == "" {
		return nil
	}
	return os.Remove(pidPath)
}
