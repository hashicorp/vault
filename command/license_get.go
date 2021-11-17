package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*LicenseGetCommand)(nil)
	_ cli.CommandAutocomplete = (*LicenseGetCommand)(nil)
)

type LicenseGetCommand struct {
	*BaseCommand

	signed bool
}

func (c *LicenseGetCommand) Synopsis() string {
	return "Get an existing license"
}

func (c *LicenseGetCommand) Help() string {
	helpText := `
Usage: vault license get [options]

  Get the currently installed license, if any:

      $ vault license get

  Get the currently installed license, if any, and output its contents as a signed blob:

      $ vault license get -signed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *LicenseGetCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("License Options")
	f.BoolVar(&BoolVar{
		Name:   "signed",
		Target: &c.signed,
		Usage:  "Whether to return a signed blob from the API.",
	})

	return set
}

func (c *LicenseGetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LicenseGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LicenseGetCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var path string
	if c.signed {
		path = "sys/license/signed"
	} else {
		path = "sys/license"
	}

	secret, err := client.Logical().Read(path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving license: %s", err))
		return 2
	}

	if secret == nil {
		c.UI.Error("License not found")
		return 2
	}

	if c.signed {
		blob := secret.Data["signed"].(string)
		if blob == "" {
			c.UI.Output("License not found or using a temporary license.")
			return 2
		} else {
			c.UI.Output(blob)
			return 0
		}
	} else {
		return OutputSecret(c.UI, secret)
	}
}
