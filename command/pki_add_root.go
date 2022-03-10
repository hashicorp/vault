package command

import (
	"fmt"
	"github.com/hashicorp/vault/command/pkicli"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"strings"
)

var (
	_ cli.Command             = (*PKIAddRootCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIAddRootCommand)(nil)
)

type PKIAddRootCommand struct {
	*BaseCommand
}

func (c *PKIAddRootCommand) Synopsis() string {
	return "Creates a new root CA"
}

func (c *PKIAddRootCommand) Help() string {
	helpText := `
Usage: vault pki add-root [mount] [K=V]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIAddRootCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PKIAddRootCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *PKIAddRootCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIAddRootCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 1 {
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1+, got %d)", len(args)))
		return 1
	}

	data, err := parseArgsData(nil, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating client: %s", err))
		return 1
	}

	mount := sanitizePath(args[0])

	ops := pkicli.NewOperations(client)
	resp, err := ops.CreateRoot(mount, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating root CA: %s", err))
		return 1
	}

	fmt.Println(*resp)

	return 0
}
