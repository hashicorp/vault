package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ryanuber/columnize"
)

// AuditListCommand is a Command that lists the enabled audits.
type AuditListCommand struct {
	Meta
}

func (c *AuditListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("audit-list", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	audits, err := client.Sys().ListAudit()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading audits: %s", err))
		return 2
	}

	if len(audits) == 0 {
		c.Ui.Error(fmt.Sprintf(
			"No audit backends are enabled. Use `vault audit-enable` to\n" +
				"enable an audit backend."))
		return 1
	}

	paths := make([]string, 0, len(audits))
	for path, _ := range audits {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	columns := []string{"Type | Description | Options"}
	for _, path := range paths {
		audit := audits[path]
		opts := make([]string, 0, len(audit.Options))
		for k, v := range audit.Options {
			opts = append(opts, k+"="+v)
		}

		columns = append(columns, fmt.Sprintf(
			"%s | %s | %s", audit.Type, audit.Description, strings.Join(opts, " ")))
	}

	c.Ui.Output(columnize.SimpleFormat(columns))
	return 0
}

func (c *AuditListCommand) Synopsis() string {
	return "Lists enabled audit backends in Vault"
}

func (c *AuditListCommand) Help() string {
	helpText := `
Usage: vault audit-list [options]

  List the enabled audit backends.

  The output lists the enabled audit backends and the options for those
  backends. The options may contain sensitive information, and therefore
  only a root Vault user can view this.

General Options:

  ` + generalOptionsUsage()
	return strings.TrimSpace(helpText)
}
