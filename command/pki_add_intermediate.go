package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/command/pkicli"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PKIAddIntermediateCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIAddIntermediateCommand)(nil)
)

type PKIAddIntermediateCommand struct {
	*BaseCommand

	flagRootMount string
}

func (c *PKIAddIntermediateCommand) Synopsis() string {
	return "Generate intermediate certificate"
}

func (c *PKIAddIntermediateCommand) Help() string {
	helpText := `

Usage: vault pki add-intermediate [options] PATH COMMON_NAME [K=V]

  Configures an intermediate mount and generate the intermediate certificate.
  The intermediate certificate is the one from which all leaf certificates will be generated.
  This intermediate will be signed by the root if root-mount specified in the input parameters.
  Other parameters can be specified in K=V format (mirroring vault write).

  Configure an intermediate mount at path pki-int with a specific ttl:
      $ vault pki add-intermediate pki pki-int example.com ttl=48000h

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIAddIntermediateCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "root-mount",
		Target:  &c.flagRootMount,
		Default: "",
		EnvVar:  "",
		Usage:   "PKI root mount",
	})

	return set
}

func (c *PKIAddIntermediateCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here, similar to `vault write`. We
	// don't know what values are valid for the role and/or common names.
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
	if len(args) < 2 {
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2+, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	rootMountPath := c.flagRootMount
	if rootMountPath != "" {
		// Check if root-mount is already configured, if not return error
		_, err = client.Logical().Read(sanitizePath(fmt.Sprintf("sys/mounts/%s", rootMountPath)))
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		// It is assumed that root certificate is generated before making this request to add intermediate
	}

	rootMountPath = sanitizePath(rootMountPath)
	mountPath := sanitizePath(args[0])
	commonName := args[1]

	// Get the remaining parameters
	data, err := parseArgsData(nil, args[2:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	if data == nil {
		data = make(map[string]interface{})
	}
	data["common_name"] = commonName

	ops := pkicli.NewOperations(c.client)
	intResp, err := ops.CreateIntermediate(rootMountPath, mountPath, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating intermediate CA: %s", err))
		return 1
	}
	fmt.Println(*intResp)
	return 0
}