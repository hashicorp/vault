package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*AuditEnableCommand)(nil)
var _ cli.CommandAutocomplete = (*AuditEnableCommand)(nil)

type AuditEnableCommand struct {
	*BaseCommand

	flagDescription string
	flagPath        string
	flagLocal       bool

	testStdin io.Reader // For tests
}

func (c *AuditEnableCommand) Synopsis() string {
	return "Enables an audit device"
}

func (c *AuditEnableCommand) Help() string {
	helpText := `
Usage: vault audit enable [options] TYPE [CONFIG K=V...]

  Enables an audit device at a given path.

  This command enables an audit device of TYPE. Additional options for
  configuring the audit device can be specified after the type in the same
  format as the "vault write" command in key/value pairs.

  For example, to configure the file audit device to write audit logs at the
  path "/var/log/audit.log":

      $ vault audit enable file file_path=/var/log/audit.log

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *AuditEnableCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "description",
		Target:     &c.flagDescription,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Human-friendly description for the purpose of this audit " +
			"device.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "", // The default is complex, so we have to manually document
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Place where the audit device will be accessible. This must be " +
			"unique across all audit devices. This defaults to the \"type\" of the " +
			"audit device.",
	})

	f.BoolVar(&BoolVar{
		Name:    "local",
		Target:  &c.flagLocal,
		Default: false,
		EnvVar:  "",
		Usage: "Mark the audit device as a local-only device. Local devices " +
			"are not replicated or removed by replication.",
	})

	return set
}

func (c *AuditEnableCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictSet(
		"file",
		"syslog",
		"socket",
	)
}

func (c *AuditEnableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuditEnableCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 1 {
		c.UI.Error("Missing TYPE!")
		return 1
	}

	// Grab the type
	auditType := strings.TrimSpace(args[0])

	auditPath := c.flagPath
	if auditPath == "" {
		auditPath = auditType
	}
	auditPath = ensureTrailingSlash(auditPath)

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	options, err := parseArgsDataString(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if err := client.Sys().EnableAuditWithOptions(auditPath, &api.EnableAuditOptions{
		Type:        auditType,
		Description: c.flagDescription,
		Options:     options,
		Local:       c.flagLocal,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error enabling audit device: %s", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Enabled the %s audit device at: %s", auditType, auditPath))
	return 0
}
