package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*AuditDisableCommand)(nil)
var _ cli.CommandAutocomplete = (*AuditDisableCommand)(nil)

// AuditDisableCommand is a Command that mounts a new mount.
type AuditDisableCommand struct {
	*BaseCommand
}

func (c *AuditDisableCommand) Synopsis() string {
	return "Disables an audit backend"
}

func (c *AuditDisableCommand) Help() string {
	helpText := `
Usage: vault audit-disable [options] PATH

  Disables an audit backend. Once an audit backend is disabled, no future
  audit logs are dispatched to it. The data associated with the audit backend
  is not affected.

  The argument corresponds to the PATH of the mount, not the TYPE!

  Disable the audit backend at file/:

      $ vault audit-disable file/

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuditDisableCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *AuditDisableCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAudits()
}

func (c *AuditDisableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuditDisableCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	path, kvs, err := extractPath(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	path = ensureTrailingSlash(path)

	if len(kvs) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if err := client.Sys().DisableAudit(path); err != nil {
		c.UI.Error(fmt.Sprintf("Error disabling audit backend: %s", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Disabled audit backend (if it was enabled) at: %s", path))

	return 0
}
