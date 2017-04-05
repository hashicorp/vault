package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/meta"
)

type PluginExec struct {
	meta.Meta
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

	runner, ok := builtinplugins.BuiltinPlugins[pluginName]
	if !ok {
		c.Ui.Error(fmt.Sprintf(
			"No plugin with the name %s found", pluginName))
		return 1
	}

	err := runner()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error running plugin: %s", err))
		return 1
	}

	return 0
}

func (c *PluginExec) Synopsis() string {
	return "Runs a builtin plugin. Should only be called by vault."
}

func (c *PluginExec) Help() string {
	helpText := `
Usage: vault plugin-exec type

  Runs a builtin plugin. Should only be called by vault.

  This will execute a plugin for use in a plugable location in vault. If run by
  a cli user it will print a message indicating it can not be executed by anyone
  other than vault. For supported plugin types see the vault documentation.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
