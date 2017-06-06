package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// PolicyReadCommand is a Command that enables a new endpoint.
type PolicyReadCommand struct {
	meta.Meta
}

func (c *PolicyReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy-read", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\npolicy-read expects only one argument"))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 2
	}

	rules, err := client.Sys().GetPolicy(args[0])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err))
		return 1
	}

	c.Ui.Output(rules)
	return 0
}

func (c *PolicyReadCommand) Synopsis() string {
	return "Read a policy from the server"
}

func (c *PolicyReadCommand) Help() string {
	helpText := `
Usage: vault policy-read [options] name

  Read an existing policy with the given name.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
