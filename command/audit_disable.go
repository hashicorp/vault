package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// AuditDisableCommand is a Command that mounts a new mount.
type AuditDisableCommand struct {
	meta.Meta
}

func (c *AuditDisableCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("mount", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\naudit-disable expects one argument: the id to disable"))
		return 1
	}

	id := args[0]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().DisableAudit(id); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error disabling audit backend: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully disabled audit backend '%s' if it was enabled", id))

	return 0
}

func (c *AuditDisableCommand) Synopsis() string {
	return "Disable an audit backend"
}

func (c *AuditDisableCommand) Help() string {
	helpText := `
Usage: vault audit-disable [options] id

  Disable an audit backend.

  Once the audit backend is disabled no more audit logs will be sent to
  it. The data associated with the audit backend isn't affected.

  The "id" parameter should map to the "path" used in "audit-enable". If
  no path was provided to "audit-enable" you should use the backend
  type (e.g. "file").

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
