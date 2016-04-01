package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// StepDownCommand is a Command that seals the vault.
type StepDownCommand struct {
	meta.Meta
}

func (c *StepDownCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("step-down", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	if err := client.Sys().StepDown(); err != nil {
		c.Ui.Error(fmt.Sprintf("Error stepping down: %s", err))
		return 1
	}

	return 0
}

func (c *StepDownCommand) Synopsis() string {
	return "Force the Vault node to give up active duty"
}

func (c *StepDownCommand) Help() string {
	helpText := `
Usage: vault step-down [options]

  Force the Vault node to step down from active duty.

  This causes the indicated node to give up active status. Note that while the
  affected node will have a short delay before attempting to grab the lock
  again, if no other node grabs the lock beforehand, it is possible for the
  same node to re-grab the lock and become active again.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
