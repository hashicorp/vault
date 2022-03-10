package command

import (
	"fmt"
	"github.com/hashicorp/vault/command/pkicli"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"strings"
)

var (
	_ cli.Command             = (*PKICreateCACommand)(nil)
	_ cli.CommandAutocomplete = (*PKICreateCACommand)(nil)
)

type PKICreateCACommand struct {
	*BaseCommand
}

func (c *PKICreateCACommand) Synopsis() string {
	return "Creates a root CA and corresponding parseIntermediateArgs CA"
}

func (c *PKICreateCACommand) Help() string {
	helpText := `
Usage: vault pki initialize-topology [root mount] [root CN] [intermediate mount] [intermediate CN] [K=V]

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKICreateCACommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *PKICreateCACommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *PKICreateCACommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKICreateCACommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 4 {
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 4+, got %d)", len(args)))
		return 1
	}

	rootMount := args[0]
	rootCommonName := args[1]
	intMount := args[2]
	intCommonName := args[3]

	data, err := parseArgsData(nil, args[4:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating client: %s", err))
		return 1
	}

	data["common_name"] = rootCommonName

	ops := pkicli.NewOperations(client)
	rootResp, err := ops.CreateRoot(rootMount, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating root CA: %s", err))
		return 1
	}

	fmt.Println(*rootResp)

	data["common_name"] = intCommonName
	intResp, err := ops.CreateIntermediate(rootMount, intMount, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating intermediate CA: %s", err))
		return 1
	}

	fmt.Println(*intResp)

	return 0
}
