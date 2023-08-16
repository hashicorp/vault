package command

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*NamespacePatchCommand)(nil)
	_ cli.CommandAutocomplete = (*NamespacePatchCommand)(nil)
)

type NamespacePatchCommand struct {
	*BaseCommand

	flagCustomMetadata       map[string]string
	flagRemoveCustomMetadata []string
}

func (c *NamespacePatchCommand) Synopsis() string {
	return "Patch an existing namespace"
}

func (c *NamespacePatchCommand) Help() string {
	helpText := `
Usage: vault namespace patch [options] PATH

  Patch an existing namespace. The namespace patched will be relative to the
  namespace provided in either the VAULT_NAMESPACE environment variable or
  -namespace CLI flag.

  Patch an existing child namespace by adding and removing custom-metadata (e.g. ns1/):

      $ vault namespace patch -custom-metadata=foo=abc -remove-custom-metadata=bar ns1

  Patch an existing child namespace from a parent namespace (e.g. ns1/ns2/):

      $ vault namespace patch -namespace=ns1 -custom-metadata=foo=abc ns2

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespacePatchCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")
	f.StringMapVar(&StringMapVar{
		Name:    "custom-metadata",
		Target:  &c.flagCustomMetadata,
		Default: map[string]string{},
		Usage: "Specifies arbitrary key=value metadata meant to describe a namespace." +
			"This can be specified multiple times to add multiple pieces of metadata.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:    "remove-custom-metadata",
		Target:  &c.flagRemoveCustomMetadata,
		Default: []string{},
		Usage:   "Key to remove from custom metadata. To specify multiple values, specify this flag multiple times.",
	})

	return set
}

func (c *NamespacePatchCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *NamespacePatchCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespacePatchCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	namespacePath := strings.TrimSpace(args[0])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data := make(map[string]interface{})
	customMetadata := make(map[string]interface{})

	for key, value := range c.flagCustomMetadata {
		customMetadata[key] = value
	}

	for _, key := range c.flagRemoveCustomMetadata {
		// A null in a JSON merge patch payload will remove the associated key
		customMetadata[key] = nil
	}

	data["custom_metadata"] = customMetadata

	secret, err := client.Logical().JSONMergePatch(context.Background(), "sys/namespaces/"+namespacePath, data)
	if err != nil {
		if re, ok := err.(*api.ResponseError); ok && re.StatusCode == http.StatusNotFound {
			c.UI.Error("Namespace not found")
			return 2
		}

		c.UI.Error(fmt.Sprintf("Error patching namespace: %s", err))
		return 2
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
