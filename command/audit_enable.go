package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/kv-builder"
	"github.com/hashicorp/vault/meta"
	"github.com/mitchellh/mapstructure"
	"github.com/posener/complete"
)

// AuditEnableCommand is a Command that mounts a new mount.
type AuditEnableCommand struct {
	meta.Meta

	// A test stdin that can be used for tests
	testStdin io.Reader
}

func (c *AuditEnableCommand) Run(args []string) int {
	var desc, path string
	var local bool
	flags := c.Meta.FlagSet("audit-enable", meta.FlagSetDefault)
	flags.StringVar(&desc, "description", "", "")
	flags.StringVar(&path, "path", "", "")
	flags.BoolVar(&local, "local", false, "")
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
	if path == "" {
		path = auditType
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

	err = client.Sys().EnableAuditWithOptions(path, &api.EnableAuditOptions{
		Type:        auditType,
		Description: desc,
		Options:     opts,
		Local:       local,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error enabling audit backend: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf(
		"Successfully enabled audit backend '%s' with path '%s'!", auditType, path))
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

  For example, to configure the file audit backend to write audit logs at
  the path /var/log/audit.log:

      $ vault audit-enable file file_path=/var/log/audit.log

  For information on available configuration options, please see the
  documentation.

General Options:
` + meta.GeneralOptionsUsage() + `
Audit Enable Options:

  -description=<desc>     A human-friendly description for the backend. This
                          shows up only when querying the enabled backends.

  -path=<path>            Specify a unique path for this audit backend. This
                          is purely for referencing this audit backend. By
                          default this will be the backend type.

  -local                  Mark the mount as a local mount. Local mounts
                          are not replicated nor (if a secondary)
                          removed by replication.
`
	return strings.TrimSpace(helpText)
}

func (c *AuditEnableCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictSet(
		"file",
		"syslog",
		"socket",
	)
}

func (c *AuditEnableCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-description": complete.PredictNothing,
		"-path":        complete.PredictNothing,
		"-local":       complete.PredictNothing,
	}
}
