// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVPatchCommand)(nil)
	_ cli.CommandAutocomplete = (*KVPatchCommand)(nil)
)

type KVPatchCommand struct {
	*BaseCommand

	flagCAS        int
	flagMethod     string
	flagMount      string
	testStdin      io.Reader // for tests
	flagRemoveData []string
}

func (c *KVPatchCommand) Synopsis() string {
	return "Sets or updates data in the KV store without overwriting"
}

func (c *KVPatchCommand) Help() string {
	helpText := `
Usage: vault kv patch [options] KEY [DATA]

  *NOTE*: This is only supported for KV v2 engine mounts.

  Writes the data to the corresponding path in the key-value store. The data can be of
  any type.

      $ vault kv patch -mount=secret foo bar=baz

  The deprecated path-like syntax can also be used, but this should be avoided, 
  as the fact that it is not actually the full API path to 
  the secret (secret/data/foo) can cause confusion: 
  
      $ vault kv patch secret/foo bar=baz

  The data can also be consumed from a file on disk by prefixing with the "@"
  symbol. For example:

      $ vault kv patch -mount=secret foo @data.json

  Or it can be read from stdin using the "-" symbol:

      $ echo "abcd1234" | vault kv patch -mount=secret foo bar=-

  To perform a Check-And-Set operation, specify the -cas flag with the
  appropriate version number corresponding to the key you want to perform
  the CAS operation on:

      $ vault kv patch -mount=secret -cas=1 foo bar=baz

  By default, this operation will attempt an HTTP PATCH operation. If your
  policy does not allow that, it will fall back to a read/local update/write approach.
  If you wish to specify which method this command should use, you may do so
  with the -method flag. When -method=patch is specified, only an HTTP PATCH
  operation will be tried. If it fails, the entire command will fail.

      $ vault kv patch -mount=secret -method=patch foo bar=baz

  When -method=rw is specified, only a read/local update/write approach will be tried.
  This was the default behavior previous to Vault 1.9.

      $ vault kv patch -mount=secret -method=rw foo bar=baz

  To remove data from the corresponding path in the key-value store, kv patch can be used.

      $ vault kv patch -mount=secret -remove-data=bar foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVPatchCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	// Patch specific options
	f := set.NewFlagSet("Common Options")

	f.IntVar(&IntVar{
		Name:    "cas",
		Target:  &c.flagCAS,
		Default: 0,
		Usage: `Specifies to use a Check-And-Set operation. If set to 0 or not
		set, the patch will be allowed. If the index is non-zero the patch will
		only be allowed if the key’s current version matches the version
		specified in the cas parameter.`,
	})

	f.StringVar(&StringVar{
		Name:   "method",
		Target: &c.flagMethod,
		Usage: `Specifies which method of patching to use. If set to "patch", then
		an HTTP PATCH request will be issued. If set to "rw", then a read will be
		performed, then a local update, followed by a remote update.`,
	})

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "", // no default, because the handling of the next arg is determined by whether this flag has a value
		Usage: `Specifies the path where the KV backend is mounted. If specified, 
		the next argument will be interpreted as the secret path. If this flag is 
		not specified, the next argument will be interpreted as the combined mount 
		path and secret path, with /data/ automatically appended between KV 
		v2 secrets.`,
	})

	f.StringSliceVar(&StringSliceVar{
		Name:    "remove-data",
		Target:  &c.flagRemoveData,
		Default: []string{},
		Usage:   "Key to remove from data. To specify multiple values, specify this flag multiple times.",
	})

	return set
}

func (c *KVPatchCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *KVPatchCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVPatchCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected >1, got %d)", len(args)))
		return 1
	case len(c.flagRemoveData) == 0 && len(args) == 1:
		c.UI.Error("Must supply data")
		return 1
	}

	var err error

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	newData, err := parseArgsData(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
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
		c.UI.Error("KV engine mount must be version 2 for patch support")
		return 2
	}

	fullPath := addPrefixToKVPath(partialPath, mountPath, "data", false)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// collecting data to be removed
	if newData == nil {
		newData = make(map[string]interface{})
	}

	for _, key := range c.flagRemoveData {
		// A null in a JSON merge patch payload will remove the associated key
		newData[key] = nil
	}

	// Check the method and behave accordingly
	var secret *api.Secret
	var code int

	switch c.flagMethod {
	case "rw":
		secret, code = c.readThenWrite(client, fullPath, newData)
	case "patch":
		secret, code = c.mergePatch(client, fullPath, newData, false)
	case "":
		secret, code = c.mergePatch(client, fullPath, newData, true)
	default:
		c.UI.Error(fmt.Sprintf("Unsupported method provided to -method flag: %s", c.flagMethod))
		return 2
	}

	if code != 0 {
		return code
	}
	if secret == nil {
		// Don't output anything if there's no secret
		return 0
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	// If the secret is wrapped, return the wrapped response.
	if secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.UI, secret)
	}

	if Format(c.UI) == "table" {
		outputPath(c.UI, fullPath, "Secret Path")
		metadata := secret.Data
		c.UI.Info(getHeaderForMap("Metadata", metadata))
		return OutputData(c.UI, metadata)
	}

	return OutputSecret(c.UI, secret)
}

func (c *KVPatchCommand) readThenWrite(client *api.Client, path string, newData map[string]interface{}) (*api.Secret, int) {
	// First, do a read.
	// Note that we don't want to see curl output for the read request.
	curOutputCurl := client.OutputCurlString()
	client.SetOutputCurlString(false)
	outputPolicy := client.OutputPolicy()
	client.SetOutputPolicy(false)
	secret, err := kvReadRequest(client, path, nil)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error doing pre-read at %s: %s", path, err))
		return nil, 2
	}
	client.SetOutputCurlString(curOutputCurl)
	client.SetOutputPolicy(outputPolicy)

	// Make sure a value already exists
	if secret == nil || secret.Data == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return nil, 2
	}

	// Verify metadata found
	rawMeta, ok := secret.Data["metadata"]
	if !ok || rawMeta == nil {
		c.UI.Error(fmt.Sprintf("No metadata found at %s; patch only works on existing data", path))
		return nil, 2
	}
	meta, ok := rawMeta.(map[string]interface{})
	if !ok {
		c.UI.Error(fmt.Sprintf("Metadata found at %s is not the expected type (JSON object)", path))
		return nil, 2
	}
	if meta == nil {
		c.UI.Error(fmt.Sprintf("No metadata found at %s; patch only works on existing data", path))
		return nil, 2
	}

	// Verify old data found
	rawData, ok := secret.Data["data"]
	if !ok || rawData == nil {
		c.UI.Error(fmt.Sprintf("No data found at %s; patch only works on existing data", path))
		return nil, 2
	}
	data, ok := rawData.(map[string]interface{})
	if !ok {
		c.UI.Error(fmt.Sprintf("Data found at %s is not the expected type (JSON object)", path))
		return nil, 2
	}
	if data == nil {
		c.UI.Error(fmt.Sprintf("No data found at %s; patch only works on existing data", path))
		return nil, 2
	}

	// Copy new data over
	for k, v := range newData {
		data[k] = v
	}

	secret, err = client.Logical().Write(path, map[string]interface{}{
		"data": data,
		"options": map[string]interface{}{
			"cas": meta["version"],
		},
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		return nil, 2
	}

	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return nil, 0
	}

	if c.flagField != "" {
		return nil, PrintRawField(c.UI, secret, c.flagField)
	}

	return secret, 0
}

func (c *KVPatchCommand) mergePatch(client *api.Client, path string, newData map[string]interface{}, rwFallback bool) (*api.Secret, int) {
	data := map[string]interface{}{
		"data":    newData,
		"options": map[string]interface{}{},
	}

	if c.flagCAS > 0 {
		data["options"].(map[string]interface{})["cas"] = c.flagCAS
	}

	secret, err := client.Logical().JSONMergePatch(context.Background(), path, data)
	if err != nil {
		// If it's a 405, that probably means the server is running a pre-1.9
		// Vault version that doesn't support the HTTP PATCH method.
		// Fall back to the old way of doing it if the user didn't specify a -method.
		// If they did, and it was "patch", then just error.
		if re, ok := err.(*api.ResponseError); ok && re.StatusCode == 405 && rwFallback {
			return c.readThenWrite(client, path, newData)
		}
		// If it's a 403, that probably means they don't have the patch capability in their policy. Fall back to
		// the old way of doing it if the user didn't specify a -method. If they did, and it was "patch", then just error.
		if re, ok := err.(*api.ResponseError); ok && re.StatusCode == 403 && rwFallback {
			c.UI.Warn(fmt.Sprintf("Data was written to %s but we recommend that you add the \"patch\" capability to your ACL policy in order to use HTTP PATCH in the future.", path))
			return c.readThenWrite(client, path, newData)
		}

		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		return nil, 2
	}

	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return nil, 0
	}

	if c.flagField != "" {
		return nil, PrintRawField(c.UI, secret, c.flagField)
	}

	return secret, 0
}
