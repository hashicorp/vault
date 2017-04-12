package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// RevokeCommand is a Command that mounts a new mount.
type RevokeCommand struct {
	meta.Meta
}

func (c *RevokeCommand) Run(args []string) int {
	var prefix, force bool
	flags := c.Meta.FlagSet("revoke", meta.FlagSetDefault)
	flags.BoolVar(&prefix, "prefix", false, "")
	flags.BoolVar(&force, "force", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nrevoke expects one argument: the ID to revoke"))
		return 1
	}
	leaseId := args[0]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	switch {
	case force && !prefix:
		c.Ui.Error(fmt.Sprintf(
			"-force requires -prefix"))
		return 1
	case force && prefix:
		err = client.Sys().RevokeForce(leaseId)
	case prefix:
		err = client.Sys().RevokePrefix(leaseId)
	default:
		err = client.Sys().Revoke(leaseId)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Revoke error: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Success! Revoked the secret with ID '%s', if it existed.", leaseId))
	return 0
}

func (c *RevokeCommand) Synopsis() string {
	return "Revoke a secret."
}

func (c *RevokeCommand) Help() string {
	helpText := `
Usage: vault revoke [options] id

  Revoke a secret by its lease ID.

  This command revokes a secret by its lease ID that was returned with it. Once
  the key is revoked, it is no longer valid.

  With the -prefix flag, the revoke is done by prefix: any secret prefixed with
  the given partial ID is revoked. Lease IDs are structured in such a way to
  make revocation of prefixes useful.

  With the -force flag, the lease is removed from Vault even if the revocation
  fails. This is meant for certain recovery scenarios and should not be used
  lightly. This option requires -prefix.

General Options:
` + meta.GeneralOptionsUsage() + `
Revoke Options:

  -prefix=true            Revoke all secrets with the matching prefix. This
                          defaults to false: an exact revocation.

  -force=true             Delete the lease even if the actual revocation
                          operation fails.
`
	return strings.TrimSpace(helpText)
}
