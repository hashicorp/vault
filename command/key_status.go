package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// KeyStatusCommand is a Command that provides information about the key status
type KeyStatusCommand struct {
	meta.Meta
}

func (c *KeyStatusCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("key-status", meta.FlagSetDefault)
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

	status, err := client.Sys().KeyStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading audits: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf("Key Term: %d", status.Term))
	c.Ui.Output(fmt.Sprintf("Installation Time: %v", status.InstallTime))
	return 0
}

func (c *KeyStatusCommand) Synopsis() string {
	return "Provides information about the active encryption key"
}

func (c *KeyStatusCommand) Help() string {
	helpText := `
Usage: vault key-status [options]

  Provides information about the active encryption key. Specifically,
  the current key term and the key installation time.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
