package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

type PluginExec struct {
	meta.Meta
}

var builtinFactories = map[string]func() error{
//	"mysql-database-plugin":    mysql.Factory,
//	"postgres-database-plugin": postgres.Factory,
}

func (c *PluginExec) Run(args []string) int {
	flags := c.Meta.FlagSet("plugin-exec", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nplugin-exec expects one argument: the plugin to execute."))
		return 1
	}

	pluginName := args[0]

	factory, ok := builtinFactories[pluginName]
	if !ok {
		c.Ui.Error(fmt.Sprintf(
			"No plugin with the name %s found", pluginName))
		return 1
	}

	err := factory()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error running plugin: %s", err))
		return 1
	}

	return 0
}

func (c *PluginExec) Synopsis() string {
	return "Force the Vault node to give up active duty"
}

func (c *PluginExec) Help() string {
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
