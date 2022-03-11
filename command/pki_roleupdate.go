package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PKIRoleUpdateCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIRoleUpdateCommand)(nil)
)

type PKIRoleUpdateCommand struct {
	*BaseCommand
}

func (c *PKIRoleUpdateCommand) Synopsis() string {
	return "Update existing PKI Secrets Engine role"
}

func (c *PKIRoleUpdateCommand) Help() string {
	helpText := `
Usage: vault pki role-update ROLE_PATH [K=V...]

  Update the specific role, fetching all present values and then updating
  only the newly specified fields.

  To disallow localhost issuance without modifying the rest of the role:

      $ vault pki role-update /pki/roles/example-com allow_localhost=false

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIRoleUpdateCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	return set
}

func (c *PKIRoleUpdateCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here, similar to `vault write`. We
	// don't know what values are valid for the role and/or common names.
	return complete.PredictAnything
}

func (c *PKIRoleUpdateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIRoleUpdateCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 2 {
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2+, got %d)", len(args)))
		return 1
	}

	rolePath := sanitizePath(args[0])
	updateData, err := parseArgsData(nil, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	roleData, err := client.Logical().Read(rolePath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading role: %v", err))
		return 2
	}

	if roleData == nil || roleData.Data == nil {
		c.UI.Error(fmt.Sprintf("Fetch succeeded but got empty role data: %v", roleData))
		return 2
	}

	for key, value := range updateData {
		roleData.Data[key] = value
	}

	_, err = client.Logical().Write(rolePath, roleData.Data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error updating role: %v", err))
		return 2
	}

	return 0
}
