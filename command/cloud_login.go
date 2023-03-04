package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*CloudLoginCommand)(nil)

type CloudLoginCommand struct {
	*BaseCommand

	flagClientID string
	flagSecretID string
}

func (c *CloudLoginCommand) Synopsis() string {
	return "Login to HCP"
}

func (c *CloudLoginCommand) Help() string {
	helpText := `
Usage: vault cloud login [options] [args]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *CloudLoginCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "client-id",
		Target:  &c.flagClientID,
		Default: "",
		Usage:   "HCP Service Principal Client ID",
	})

	f.StringVar(&StringVar{
		Name:    "secret-id",
		Target:  &c.flagSecretID,
		Default: "",
		Usage:   "HCP Service Principal Secret ID",
	})
	return set
}

func (c *CloudLoginCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *CloudLoginCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *CloudLoginCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	opts := []config.HCPConfigOption{config.FromEnv()}
	if c.flagClientID != "" || c.flagSecretID != "" {
		opts = append(opts, config.WithClientCredentials(c.flagClientID, c.flagSecretID))
	}

	_, err := config.NewHCPConfig(opts...)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating HCP Config: %s", err))
		return 1
	}

	return 0
}
