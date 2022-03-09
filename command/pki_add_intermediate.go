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
	_ cli.Command = (*PKIAddIntermediateCommand)(nil)
	_ cli.Command = (*PKIAddIntermediateCommand)(nil)
)

type PKIAddIntermediateCommand struct {
	*BaseCommand

	flagRootMount string
	flagMountName string
}

func (c *PKIAddIntermediateCommand) Synopsis() string {
	return "Creates a new parseIntermediateArgs CA"
}

func (c *PKIAddIntermediateCommand) Help() string {
	helpText := `
Usage: vault pki add-intermediate [ARGS]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIAddIntermediateCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "root-mount",
		Target: &c.flagRootMount,
		Usage:  "The name of the root mount to use for signing the parseIntermediateArgs certificate",
	})

	f.StringVar(&StringVar{
		Name:   "mount",
		Target: &c.flagMountName,
		Usage:  "The name of the mount for the root CA. The name must be unique.",
	})

	return set
}

func (c *PKIAddIntermediateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *PKIAddIntermediateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIAddIntermediateCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) != 1 {
		c.UI.Error(fmt.Sprintf("Wrong number of arguments (expected 1, got %d)", len(args)))
		return 1
	}

	rootMount := sanitizePath(c.flagRootMount)
	intMount := sanitizePath(c.flagMountName)

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting client: %s", err))
		return 1
	}

	var params map[string]interface{}

	if err := jsonutil.DecodeJSONFromReader(strings.NewReader(args[0]), &params); err != nil {
		c.UI.Error(fmt.Sprintf("Error parsing arguments for intermediate CA: %s", err))
		return 1
	}

	ops := pkicli.NewOperations(client)

	resp, err := ops.CreateIntermediate(rootMount, intMount, params)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating intermediate CA: %s", err))
		return 1
	}

	fmt.Println(resp)

	return 0
}
