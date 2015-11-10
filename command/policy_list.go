package command

import (
	"fmt"
	"github.com/hashicorp/vault/vault"
	"strings"
)

const (
	PathPolicyDeny  = "deny"
	PathPolicyRead  = "read"
	PathPolicyWrite = "write"
	PathPolicySudo  = "sudo"
)

// PolicyListCommand is a Command that enables a new endpoint.
type PolicyListCommand struct {
	Meta
}

func (c *PolicyListCommand) Run(args []string) int {
	var detailed bool
	flags := c.Meta.FlagSet("policy-list", FlagSetDefault)
	flags.BoolVar(&detailed, "detailed", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if detailed {
		return c.detailedOutput()
	}
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

//sorted policy from higher to lower
func (c *PolicyListCommand) detailedOutput() int {
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

	for i, policy := range policies {
		rules, err := client.Sys().GetPolicy(policy)
		if err != nil {

		}
		pol, err := vault.Parse(rules)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error: %s", err))
		}

		if policy == "root" {
			c.Ui.Output(fmt.Sprintf("%d. Name:'%s'", i+1, policy))
			continue
		}

		var deny, write, read, total, sudo int
		deny, write, read, total, sudo = 0, 0, 0, 0, 0
		pols := make(map[string]int)
		for _, path := range pol.Paths {
			if len(path.Prefix) > 0 {
				splitter := strings.Split(path.Prefix, "/")[0]
				pols[splitter+"/*"]++
			}
			switch path.Policy {
			case PathPolicyDeny:
				deny++
			case PathPolicyRead:
				read++
			case PathPolicyWrite:
				write++
			case PathPolicySudo:
				sudo++
			}
			total++
		}
		c.Ui.Output(fmt.Sprintf(
			"%d. Name: '%s'. Total policies: %d", i+1, policy, total))
		c.Ui.Output(fmt.Sprintf(
			"deny: %d, write: %d, read: %d, sudo: %d\n", deny, write, read, sudo))

		output := ""
		for pol, value := range pols {
			output += fmt.Sprintf("%s: %d\n", pol, value)
		}

		c.Ui.Output(output)

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
  -detailed               Output detailed information about policies.

  ` + generalOptionsUsage()
	return strings.TrimSpace(helpText)
}
