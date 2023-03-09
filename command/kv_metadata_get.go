package command

import (
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVMetadataGetCommand)(nil)
	_ cli.CommandAutocomplete = (*KVMetadataGetCommand)(nil)
)

type KVMetadataGetCommand struct {
	*BaseCommand
	flagMount string
}

func (c *KVMetadataGetCommand) Synopsis() string {
	return "Retrieves key metadata from the KV store"
}

func (c *KVMetadataGetCommand) Help() string {
	helpText := `
Usage: vault kv metadata get [options] KEY

  Retrieves the metadata from Vault's key-value store at the given key name. If no
  key exists with that name, an error is returned.

      $ vault kv metadata get -mount=secret foo

  The deprecated path-like syntax can also be used, but this should be avoided 
  for KV v2, as the fact that it is not actually the full API path to 
  the secret (secret/metadata/foo) can cause confusion: 
  
      $ vault kv metadata get secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVMetadataGetCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat | FlagSetClipboard)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "", // no default, because the handling of the next arg is determined by whether this flag has a value
		Usage: `Specifies the path where the KV backend is mounted. If specified, 
		the next argument will be interpreted as the secret path. If this flag is 
		not specified, the next argument will be interpreted as the combined mount 
		path and secret path, with /metadata/ automatically appended between KV 
		v2 secrets.`,
	})

	return set
}

func (c *KVMetadataGetCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVMetadataGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVMetadataGetCommand) Run(args []string) int {
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

	if !validateClipboardFlag(c.BaseCommand) {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// If true, we're working with "-mount=secret foo" syntax.
	// If false, we're using "secret/foo" syntax.
	mountFlagSyntax := c.flagMount != ""

	var (
		mountPath   string
		partialPath string
		v2          bool
	)

	// Parse the paths and grab the KV version
	if mountFlagSyntax {
		// In this case, this arg is the secret path (e.g. "foo").
		partialPath = sanitizePath(args[0])
		mountPath, v2, err = isKVv2(sanitizePath(c.flagMount), client)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}

		if v2 {
			partialPath = path.Join(mountPath, partialPath)
		}
	} else {
		// In this case, this arg is a path-like combination of mountPath/secretPath.
		// (e.g. "secret/foo")
		partialPath = sanitizePath(args[0])
		mountPath, v2, err = isKVv2(partialPath, client)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
	}

	if !v2 {
		c.UI.Error("Metadata not supported on KV Version 1")
		return 1
	}

	fullPath := addPrefixToKVPath(partialPath, mountPath, "metadata")
	secret, err := client.Logical().Read(fullPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading %s: %s", fullPath, err))
		return 2
	}
	if secret == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", fullPath))
		return 2
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField, c.flagClipboard, c.flagClipboardTTL)
	}

	// If we have wrap info print the secret normally.
	if secret.WrapInfo != nil || c.flagFormat != "table" {
		return OutputSecret(c.UI, secret)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok || versionsRaw == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", fullPath))
		OutputSecret(c.UI, secret)
		return 2
	}
	versions := versionsRaw.(map[string]interface{})

	delete(secret.Data, "versions")

	outputPath(c.UI, fullPath, "Metadata Path")

	c.UI.Info(getHeaderForMap("Metadata", secret.Data))
	OutputSecret(c.UI, secret)

	versionKeys := []int{}
	for k := range versions {
		i, err := strconv.Atoi(k)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing version %s", k))
			return 2
		}

		versionKeys = append(versionKeys, i)
	}

	sort.Ints(versionKeys)

	for _, v := range versionKeys {
		c.UI.Info("\n" + getHeaderForMap(fmt.Sprintf("Version %d", v), versions[strconv.Itoa(v)].(map[string]interface{})))
		OutputData(c.UI, versions[strconv.Itoa(v)])
	}

	return 0
}
