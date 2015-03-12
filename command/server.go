package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/flag-slice"
)

// ServerCommand is a Command that starts the Vault server.
type ServerCommand struct {
	Meta
}

func (c *ServerCommand) Run(args []string) int {
	var configPath []string
	flags := c.Meta.FlagSet("server", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	flags.Var((*sliceflag.StringFlag)(&configPath), "config", "config")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validation
	if len(configPath) == 0 {
		c.Ui.Error("At least one config path must be specified with -config")
		flags.Usage()
		return 1
	}

	// Load the configuration
	var config *server.Config
	for _, path := range configPath {
		current, err := server.LoadConfig(path)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error loading configuration from %s: %s", path, err))
			return 1
		}

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}

	// Initialize the listeners

	return 0
}

func (c *ServerCommand) Synopsis() string {
	return "Start a Vault server"
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: vault server [options]

  Start a Vault server.

  This command starts a Vault server that responds to API requests.
  Vault will start in a "sealed" state. The Vault must be unsealed
  with "vault unseal" or the API before this server can respond to requests.
  This must be done for every server.

  If the server is being started against a storage backend that has
  brand new (no existing Vault data in it), it must be initialized with
  "vault init" or the API first.


General Options:

  -config=<path>      Path to the configuration file or directory. This can be
                      specified multiple times. If it is a directory, all
                      files with a ".hcl" or ".json" suffix will be loaded.

`
	return strings.TrimSpace(helpText)
}
