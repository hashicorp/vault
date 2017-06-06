package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// PolicyListCommand is a Command that enables a new endpoint.
type PolicyListCommand struct {
	meta.Meta
}

func (c *PolicyListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy-list", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) == 1 {
		return c.read(args[0])
	} else if len(args) == 0 {
		return c.list()
	} else {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\npolicies expects zero or one arguments"))
		return 1
	}
}

func (c *PolicyListCommand) list() int {
	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	policies, err := client.Sys().ListPolicies()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error: %s", err))
		return 1
	}

	for _, p := range policies {
		c.Ui.Output(p)
	}

	return 0
}

func (c *PolicyListCommand) read(n string) int {
	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	rules, err := client.Sys().GetPolicy(n)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error: %s", err))
		return 1
	}

	c.Ui.Output(rules)
	return 0
}

func (c *PolicyListCommand) Synopsis() string {
	return "List the policies on the server"
}

func (c *PolicyListCommand) Help() string {
	helpText := `
Usage: vault policy-list [options] [name]

  List the policies that are registered with Vault server.

  If the 'name' argument is specified, if a policy exists by
  that name with Vault server, rules of the same is returned.
  Note: 'name' option will be deprecated, use 'policy-read' instead.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
