package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/flag-kv"
	"github.com/hashicorp/vault/helper/flag-slice"
	"github.com/hashicorp/vault/meta"
)

// TokenCreateCommand is a Command that mounts a new mount.
type TokenCreateCommand struct {
	meta.Meta
}

func (c *TokenCreateCommand) Run(args []string) int {
	var format string
	var id, displayName, lease, ttl, role string
	var orphan, noDefaultPolicy bool
	var metadata map[string]string
	var numUses int
	var policies []string
	flags := c.Meta.FlagSet("mount", meta.FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&displayName, "display-name", "", "")
	flags.StringVar(&id, "id", "", "")
	flags.StringVar(&lease, "lease", "", "")
	flags.StringVar(&ttl, "ttl", "", "")
	flags.StringVar(&role, "role", "", "")
	flags.BoolVar(&orphan, "orphan", false, "")
	flags.BoolVar(&noDefaultPolicy, "no-default-policy", false, "")
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

	if ttl == "" {
		ttl = lease
	}

	tcr := &api.TokenCreateRequest{
		ID:              id,
		Policies:        policies,
		Metadata:        metadata,
		TTL:             ttl,
		NoParent:        orphan,
		NoDefaultPolicy: noDefaultPolicy,
		DisplayName:     displayName,
		NumUses:         numUses,
	}

	var secret *api.Secret
	if role != "" {
		secret, err = client.Auth().Token().CreateWithRole(tcr, role)
	} else {
		secret, err = client.Auth().Token().Create(tcr)
	}

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

  A lease can also be associated with the token. If a lease is not associated
  with the token, then it cannot be renewed. If a lease is associated with
  the token, it will expire after that amount of time unless it is renewed.

  Metadata associated with the token (specified with "-metadata") is
  written to the audit log when the token is used.

  If a role is specified, the role may override parameters specified here.

General Options:
` + meta.GeneralOptionsUsage() + `
Token Options:

  -id="7699125c-d8...."   The token value that clients will use to authenticate
                          with vault. If not provided this defaults to a 36
                          character UUID. A root token is required to specify
                          the ID of a token.

  -display-name="name"    A display name to associate with this token. This
                          is a non-security sensitive value used to help
                          identify created secrets, i.e. prefixes.

  -lease="1h"             Deprecated; use "-ttl" instead.

  -ttl="1h"               TTL to associate with the token. This option enables
                          the tokens to be renewable.

  -metadata="key=value"   Metadata to associate with the token. This shows
                          up in the audit log. This can be specified multiple
                          times.

  -orphan                 If specified, the token will have no parent. Only
                          root tokens can create orphan tokens. This prevents
                          the new token from being revoked with your token.

  -no-default-policy      If specified, the token will not have the "default"
                          policy included in its policy set.

  -policy="name"          Policy to associate with this token. This can be
                          specified multiple times.

  -use-limit=5            The number of times this token can be used until
                          it is automatically revoked.

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.

  -role=name              If set, the token will be created against the named
                          role. The role may override other parameters. This
                          requires the client to have permissions on the
                          appropriate endpoint (auth/token/create/<name>).
`
	return strings.TrimSpace(helpText)
}
