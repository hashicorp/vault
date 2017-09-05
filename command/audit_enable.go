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

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*AuditEnableCommand)(nil)
var _ cli.CommandAutocomplete = (*AuditEnableCommand)(nil)

// AuditEnableCommand is a Command that mounts a new mount.
type AuditEnableCommand struct {
	*BaseCommand

	flagDescription string
	flagPath        string
	flagLocal       bool

	testStdin io.Reader // For tests
}

func (c *AuditEnableCommand) Synopsis() string {
	return "Enables an audit backend"
}

func (c *AuditEnableCommand) Help() string {
	helpText := `
Usage: vault audit-enable [options] TYPE [CONFIG K=V...]

  Enables an audit backend at a given path.

  This command enables an audit backend of type "type". Additional
  options for configuring the audit backend can be specified after the
  type in the same format as the "vault write" command in key/value pairs.

  For example, to configure the file audit backend to write audit logs at
  the path /var/log/audit.log:

      $ vault audit-enable file file_path=/var/log/audit.log

  For information on available configuration options, please see the
  documentation.

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
			"backend.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "", // The default is complex, so we have to manually document
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Place where the audit backend will be accessible. This must be " +
			"unique across all audit backends. This defaults to the \"type\" of the " +
			"audit backend.",
	})

	f.BoolVar(&BoolVar{
		Name:    "local",
		Target:  &c.flagLocal,
		Default: false,
		EnvVar:  "",
		Usage: "Mark the audit backend as a local-only backned. Local backends " +
			"are not replicated nor removed by replication.",
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
		c.UI.Error(fmt.Sprintf("Error enabling audit backend: %s", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Enabled the %s audit backend at: %s", auditType, auditPath))
	return 0
}
