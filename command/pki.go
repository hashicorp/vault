package command

import (
	"github.com/hashicorp/vault/command/pkicli"
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*PKICommand)(nil)

type PKICommand struct {
	*BaseCommand
}

func (c *PKICommand) Synopsis() string {
	return "Interact with PKI Secret Engines"
}

func (c *PKICommand) Help() string {
	helpText := `
Usage: vault pki <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's PKI Secrets
  Engine. Operators can manage PKI mounts and roles.

  To test role based issuance:

       $ vault pki role-test -mount=pki-int server-role example.com

  To add new intermediate:
       $ vault pki add-intermediate pki pki-int example.com ttl=43800h csr=@example.csr format=pem_bundle

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PKICommand) Run(args []string) int {
	//c.testCreateRoot()
	//c.testCreateIntermediate()
	//return 0
	return cli.RunResultHelp
}

func (c *PKICommand) testCreateRoot() int {
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
	}
	ops := pkicli.NewOperations(client)

	vaultAddress := client.Address()
	_, err = ops.CreateRoot("pki-root", map[string]interface{}{
		"max_lease_ttl": "24h",
		"common_name": "example.com",
		"ttl": "87600h",
		"issuing_certificates": vaultAddress + "/v1/pki/ca",
		"crl_distribution_points": vaultAddress + "/v1/pki/crl",
	})
	if err != nil {
		c.UI.Error(err.Error())
	}

	return 0
}

func (c *PKICommand) testCreateIntermediate() int {
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
	}
	ops := pkicli.NewOperations(client)

	_, err = ops.CreateIntermediate("pki-root", "pki_int", map[string]interface{}{
		"max_lease_ttl": "24h",
		"common_name": "example.com Intermediate Authority",
		"ttl": "43800h",
	})
	if err != nil {
		c.UI.Error(err.Error())
	}

	return 0
}