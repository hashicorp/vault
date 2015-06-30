package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/flag-kv"
	"github.com/hashicorp/vault/helper/flag-slice"
)

// TokenCreateCommand is a Command that mounts a new mount.
type TokenCreateCommand struct {
	Meta
}

func (c *TokenCreateCommand) Run(args []string) int {
	var format string
	var displayName, lease string
	var orphan bool
	var metadata map[string]string
	var numUses int
	var policies []string
	flags := c.Meta.FlagSet("mount", FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&displayName, "display-name", "", "")
	flags.StringVar(&lease, "lease", "", "")
	flags.BoolVar(&orphan, "orphan", false, "")
	flags.IntVar(&numUses, "use-limit", 0, "")
	flags.Var((*kvFlag.Flag)(&metadata), "metadata", "")
	flags.Var((*sliceflag.StringFlag)(&policies), "policy", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 0 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ntoken-create expects no arguments"))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies:    policies,
		Metadata:    metadata,
		Lease:       lease,
		NoParent:    orphan,
		DisplayName: displayName,
		NumUses:     numUses,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error creating token: %s", err))
		return 2
	}

	return OutputSecret(c.Ui, format, secret)
}

func (c *TokenCreateCommand) Synopsis() string {
	return "Create a new auth token"
}

func (c *TokenCreateCommand) Help() string {
	helpText := `
Usage: vault token-create [options]

  Create a new auth token.

  This command creates a new token that can be used for authentication.
  This token will be created as a child of your token. The created token
  will inherit your policies, or can be assigned a subset of your policies.

  A lease can also be associated with the token. If a lease is associated,
  it will expire after that amount of time unless it is renewed.

  Metadata associated with the token (specified with "-metadata") is
  written to the audit log when the token is used.

General Options:

  ` + generalOptionsUsage() + `

Token Options:

  -display-name="name"    A display name to associate with this token. This
                          is a non-security sensitive value used to help
                          identify created secrets, i.e. prefixes.

  -lease="1h"             Lease to associate with the token.

  -metadata="key=value"   Metadata to associate with the token. This shows
                          up in the audit log. This can be specified multiple
                          times.

  -orphan                 If specified, the token will have no parent. Only
                          root tokens can create orphan tokens. This prevents
                          the new token from being revoked with your token.

  -policy="name"          Policy to associate with this token. This can be
                          specified multiple times.

  -use-limit=5            The number of times this token can be used until
                          it is automatically revoked.

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json.

`
	return strings.TrimSpace(helpText)
}
