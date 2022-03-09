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
	_ cli.Command             = (*PKICreateCACommand)(nil)
	_ cli.CommandAutocomplete = (*PKICreateCACommand)(nil)
)

type PKICreateCACommand struct {
	*BaseCommand

	flagRootMountName string
	flagIntMountName  string
}

func (c *PKICreateCACommand) Synopsis() string {
	return "Creates a root CA and corresponding parseIntermediateArgs CA"
}

func (c *PKICreateCACommand) Help() string {
	helpText := `
Usage: vault pki create-ca [ARGS]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKICreateCACommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "root-mount",
		Target: &c.flagRootMountName,
		Usage:  "The name of the mount for the root CA. The name must be unique.",
	})

	f.StringVar(&StringVar{
		Name:   "int-mount",
		Target: &c.flagIntMountName,
		Usage:  "The name of the mount for the parseIntermediateArgs CA. The name must be unique.",
	})

	return set
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
	if len(args) < 2 {
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d", len(args)))
		return 1
	}

	rootMount := sanitizePath(c.flagRootMountName)
	intMount := sanitizePath(c.flagIntMountName)

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting client: %s", err))
		return 1
	}

	data, err := parseArgsData(nil, args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	if data == nil {
		c.UI.Error("Missing required arguments")
		return 1
	}

	if _, ok := data["root_args"]; !ok {
		c.UI.Error("Missing arguments for root CA")
		return 1
	}

	if _, ok := data["intermediate_args"]; !ok {
		c.UI.Error("Missing arguments for intermediate CA")
		return 1
	}

	var rootArgs map[string]interface{}
	var intArgs map[string]interface{}

	rootArgs, err = c.parseArgs(data["root_args"].(string))
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting arguments for root CA: %s", err))
		return 1
	}

	intArgs, err = c.parseArgs(data["intermediate_args"].(string))
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting arguments for intermediate CA: %s", err))
		return 1
	}

	ops := pkicli.NewOperations(client)
	rootResp, err := ops.CreateRoot(rootMount, rootArgs)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating root CA: %s", err))
		return 1
	}

	fmt.Println(rootResp)

	intResp, err := ops.CreateIntermediate(rootMount, intMount, intArgs)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating intermediate CA: %s", err))
		return 1
	}

	fmt.Println(intResp)

	return 0
}

func (c *PKICreateCACommand) parseArgs(args string) (map[string]interface{}, error) {
	argMap := make(map[string]interface{})

	if err := jsonutil.DecodeJSONFromReader(strings.NewReader(args), argMap); err != nil {
		return nil, err
	}

	return argMap, nil
}
