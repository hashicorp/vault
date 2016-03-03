package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// CapabilitiesCommand is a Command that enables a new endpoint.
type CapabilitiesCommand struct {
	Meta
}

func (c *CapabilitiesCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("capabilities", FlagSetDefault)
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
	switch len(args) {
	case 1:
		// only path is provided
		path = args[0]
	case 2:
		// both token and path are provided
		token = args[0]
		path = args[1]
	default:
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	var resp *api.CapabilitiesResponse
	if token == "" {
		resp, err = client.Sys().CapabilitiesSelf(path)
	} else {
		resp, err = client.Sys().Capabilities(token, path)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error retrieving capabilities: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Capabilities: %s", resp.Capabilities))
	if resp.Message != "" {
		c.Ui.Output(fmt.Sprintf("Message: %s", resp.Message))
	}
	return 0
}

func (c *CapabilitiesCommand) Synopsis() string {
	return "Fetch the capabilities of a given token on a given path"
}

func (c *CapabilitiesCommand) Help() string {
	helpText := `
Usage: vault capabilities [options] [token] path

  Fetch the capabilities of a token on a given path.
  If a token is provided to the command, API '/sys/capabilities' will be invoked
  with the given token; otherwise API '/sys/capabilities-self' will be invoked with
  the client token.

General Options:

  ` + generalOptionsUsage()
	return strings.TrimSpace(helpText)
}
