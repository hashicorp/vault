// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_                cli.Command             = (*TransformImportCommand)(nil)
	_                cli.CommandAutocomplete = (*TransformImportCommand)(nil)
	transformKeyPath                         = regexp.MustCompile("^(.*)/transformations/(fpe|tokenization)/([^/]*)$")
)

type TransformImportCommand struct {
	*BaseCommand
}

func (c *TransformImportCommand) Synopsis() string {
	return "Import a key into the Transform secrets engines."
}

func (c *TransformImportCommand) Help() string {
	helpText := `
Usage: vault transform import PATH KEY [options...]

  Using the Transform key wrapping system, imports key material from
  the base64 encoded KEY (either directly on the CLI or via @path notation),
  into a new FPE or tokenization transformation whose API path is PATH. 

  To import a new key version into an existing tokenization transformation, 
  use import_version. 
  
  The remaining options after KEY (key=value style) are passed on to 
  Create/Update FPE Transformation or Create/Update Tokenization Transformation 
  API endpoints.

  For example:
  $ vault transform import transform/transformations/tokenization/application-form @path/to/key \        
       allowed_roles=legacy-system 
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransformImportCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *TransformImportCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransformImportCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransformImportCommand) Run(args []string) int {
	return ImportKey(c.BaseCommand, "import", transformImportKeyPath, c.Flags(), args)
}

func transformImportKeyPath(s string, operation string) (path string, apiPath string, err error) {
	parts := transformKeyPath.FindStringSubmatch(s)
	if len(parts) != 4 {
		return "", "", errors.New("expected transform path and key name in the form :path:/transformations/fpe|tokenization/:name:")
	}
	path = parts[1]
	transformation := parts[2]
	keyName := parts[3]
	apiPath = path + "/transformations/" + transformation + "/" + keyName + "/" + operation

	return path, apiPath, nil
}
