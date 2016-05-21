package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// CapabilitiesCommand is a Command that enables a new endpoint.
type CapabilitiesCommand struct {
	meta.Meta
}

func (c *CapabilitiesCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("capabilities", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 2 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ncapabilities expects at most two arguments"))
		return 1
	}

	var token string
	var path string
	switch {
	case len(args) == 1:
		path = args[0]
	case len(args) == 2:
		token = args[0]
		path = args[1]
	default:
		flags.Usage()
		c.Ui.Error(fmt.Sprintf("\ncapabilities expects at least one argument"))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	var capabilities []string
	if token == "" {
		capabilities, err = client.Sys().CapabilitiesSelf(path)
	} else {
		capabilities, err = client.Sys().Capabilities(token, path)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error retrieving capabilities: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Capabilities: %s", capabilities))
	return 0
}

func (c *CapabilitiesCommand) Synopsis() string {
	return "Fetch the capabilities of a token on a given path"
}

func (c *CapabilitiesCommand) Help() string {
	helpText := `
Usage: vault capabilities [options] [token] path

  Fetch the capabilities of a token on a given path.
  If a token is provided as an argument, the '/sys/capabilities' endpoint will be invoked
  with the given token; otherwise the '/sys/capabilities-self' endpoint will be invoked
  with the client token.

  If a token does not have any capability on a given path, or if any of the policies
  belonging to the token explicitly have ["deny"] capability, or if the argument path
  is invalid, this command will respond with a ["deny"].

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
