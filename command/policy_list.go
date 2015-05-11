package command

import (
	"fmt"
	"strings"
)

// PolicyListCommand is a Command that enables a new endpoint.
type PolicyListCommand struct {
	Meta
}

func (c *PolicyListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy-list", FlagSetDefault)
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
Usage: vault policies [options] [name]

  List the policies that are available or read a single policy.

  This command lists the policies that are written to the Vault server.
  If a name of a policy is specified, that policy is outputted.

General Options:

  -address=addr           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -tls-skip-verify        Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

`
	return strings.TrimSpace(helpText)
}
