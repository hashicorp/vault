package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*AuditListCommand)(nil)
var _ cli.CommandAutocomplete = (*AuditListCommand)(nil)

type AuditListCommand struct {
	*BaseCommand

	flagDetailed bool
}

func (c *AuditListCommand) Synopsis() string {
	return "Lists enabled audit devices"
}

func (c *AuditListCommand) Help() string {
	helpText := `
Usage: vault audit list [options]

  Lists the enabled audit devices in the Vault server. The output lists the
  enabled audit devices and the options for those devices.

  List all audit devices:

      $ vault audit list

  List detailed output about the audit devices:

      $ vault audit list -detailed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuditListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "detailed",
		Target:  &c.flagDetailed,
		Default: false,
		EnvVar:  "",
		Usage: "Print detailed information such as options and replication " +
			"status about each auth device.",
	})

	return set
}

func (c *AuditListCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *AuditListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuditListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	audits, err := client.Sys().ListAudit()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing audits: %s", err))
		return 2
	}

	if len(audits) == 0 {
		c.UI.Output(fmt.Sprintf("No audit devices are enabled."))
		return 0
	}

	switch Format(c.UI) {
	case "table":
		if c.flagDetailed {
			c.UI.Output(tableOutput(c.detailedAudits(audits), nil))
			return 0
		}
		c.UI.Output(tableOutput(c.simpleAudits(audits), nil))
		return 0
	default:
		return OutputData(c.UI, audits)
	}
}

func (c *AuditListCommand) simpleAudits(audits map[string]*api.Audit) []string {
	paths := make([]string, 0, len(audits))
	for path, _ := range audits {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	columns := []string{"Path | Type | Description"}
	for _, path := range paths {
		audit := audits[path]
		columns = append(columns, fmt.Sprintf("%s | %s | %s",
			audit.Path,
			audit.Type,
			audit.Description,
		))
	}

	return columns
}

func (c *AuditListCommand) detailedAudits(audits map[string]*api.Audit) []string {
	paths := make([]string, 0, len(audits))
	for path, _ := range audits {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	columns := []string{"Path | Type | Description | Replication | Options"}
	for _, path := range paths {
		audit := audits[path]

		opts := make([]string, 0, len(audit.Options))
		for k, v := range audit.Options {
			opts = append(opts, k+"="+v)
		}

		replication := "replicated"
		if audit.Local {
			replication = "local"
		}

		columns = append(columns, fmt.Sprintf("%s | %s | %s | %s | %s",
			path,
			audit.Type,
			audit.Description,
			replication,
			strings.Join(opts, " "),
		))
	}

	return columns
}
