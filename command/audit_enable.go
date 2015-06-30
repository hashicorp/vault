package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/helper/kv-builder"
	"github.com/mitchellh/mapstructure"
)

// AuditEnableCommand is a Command that mounts a new mount.
type AuditEnableCommand struct {
	Meta

	// A test stdin that can be used for tests
	testStdin io.Reader
}

func (c *AuditEnableCommand) Run(args []string) int {
	var desc, id string
	flags := c.Meta.FlagSet("audit-enable", FlagSetDefault)
	flags.StringVar(&desc, "description", "", "")
	flags.StringVar(&id, "id", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\naudit-enable expects at least one argument: the type to enable"))
		return 1
	}

	auditType := args[0]
	if id == "" {
		id = auditType
	}

	// Build the options
	var stdin io.Reader = os.Stdin
	if c.testStdin != nil {
		stdin = c.testStdin
	}
	builder := &kvbuilder.Builder{Stdin: stdin}
	if err := builder.Add(args[1:]...); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error parsing options: %s", err))
		return 1
	}

	var opts map[string]string
	if err := mapstructure.WeakDecode(builder.Map(), &opts); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error parsing options: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 1
	}

	err = client.Sys().EnableAudit(id, auditType, desc, opts)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error enabling audit backend: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully enabled audit backend '%s'!", auditType))
	return 0
}

func (c *AuditEnableCommand) Synopsis() string {
	return "Enable an audit backend"
}

func (c *AuditEnableCommand) Help() string {
	helpText := `
Usage: vault audit-enable [options] type [config...]

  Enable an audit backend.

  This command enables an audit backend of type "type". Additional
  options for configuring the audit backend can be specified after the
  type in the same format as the "vault write" command in key/value pairs.
  Example: vault audit-enable file path=audit.log

General Options:

  ` + generalOptionsUsage() + `

Audit Enable Options:

  -description=<desc>     A human-friendly description for the backend. This
                          shows up only when querying the enabled backends.

  -id=<id>                Specify a unique ID for this audit backend. This
                          is purely for referencing this audit backend. By
                          default this will be the backend type.

`
	return strings.TrimSpace(helpText)
}
