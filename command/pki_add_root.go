package command

import (
	"fmt"
	"github.com/hashicorp/vault/command/pkicli"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
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

	flagMountName string
}

func (c *PKIAddRootCommand) Synopsis() string {
	return "Creates a new root CA"
}

func (c *PKIAddRootCommand) Help() string {
	helpText := `
Usage: vault pki add-root [ARGS]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIAddRootCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "mount",
		Target: &c.flagMountName,
		Usage:  "The name of the mount for the root CA. The name must be unique.",
	})

	return set
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

	data, err := parseArgsData(nil, args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	var params map[string]interface{}

	if _, ok := data["config"]; ok {
		if err := jsonutil.DecodeJSONFromReader(strings.NewReader(data["config"].(string)), &params); err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing arguments for root CA: %s", err))
			return 1
		}
	} else {
		params = data
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating client: %s", err))
		return 1
	}

	mount := sanitizePath(c.flagMountName)

	ops := pkicli.NewOperations(client)
	resp, err := ops.CreateRoot(mount, params)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating root CA: %s", err))
		return 1
	}

	fmt.Println(*resp)

	return 0
}
