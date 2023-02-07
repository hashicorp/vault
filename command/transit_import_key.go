// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"

	"github.com/google/tink/go/kwp/subtle"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_       cli.Command             = (*TransitImportCommand)(nil)
	_       cli.CommandAutocomplete = (*TransitImportCommand)(nil)
	keyPath                         = regexp.MustCompile("^(.*)/keys/([^/]*)$")
)

type TransitImportCommand struct {
	*BaseCommand
}

func (c *TransitImportCommand) Synopsis() string {
	return "Import a key into the Transit or Transform secrets engines."
}

func (c *TransitImportCommand) Help() string {
	helpText := `
Usage: vault transit import PATH KEY [options...]

  Using the Transit or Transform key wrapping system, imports key material from
  the base64 encoded KEY, into a new key whose API path is PATH.  To import a new version
  into an existing key, use import_version.  The remaining options after KEY (key=value style) are passed
  on to the transit/transform create key endpoint.  If your system or device natively supports
  the RSA AES key wrap mechanism, you should use it directly rather than this command. 
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransitImportCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
}

func (c *TransitImportCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitImportCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransitImportCommand) Run(args []string) int {
	return importKey(c.BaseCommand, "import", args)
}

// error codes: 1: user error, 2: internal computation error, 3: remote api call error
func importKey(c *BaseCommand, operation string, args []string) int {
	if len(args) != 2 {
		c.UI.Error(fmt.Sprintf("Incorrect argument count (expected 2, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	ephemeralAESKey := make([]byte, 32)
	_, err = rand.Read(ephemeralAESKey)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to generate ephemeral key: %v", err))
	}
	parts := keyPath.FindStringSubmatch(args[0])
	if len(parts) != 3 {
		c.UI.Error("expected transit path and key name in the form :path:/keys/:name:")
		return 1
	}
	path := parts[1]
	keyName := parts[2]

	key, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		c.UI.Error(fmt.Sprintf("error base64 decoding source key material: %v", err))
		return 1
	}
	// Fetch the wrapping key
	c.UI.Output("Retrieving transit wrapping key.")
	wrappingKey, err := fetchWrappingKey(c, client, path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to fetch wrapping key: %v", err))
		return 3
	}
	c.UI.Output("Wrapping source key with ephemeral key.")
	wrapKWP, err := subtle.NewKWP(ephemeralAESKey)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failure building key wrapping key: %v", err))
		return 2
	}
	wrappedTargetKey, err := wrapKWP.Wrap(key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failure wrapping source key: %v", err))
		return 2
	}
	c.UI.Output("Encrypting ephemeral key with transit wrapping key.")
	wrappedAESKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		wrappingKey.(*rsa.PublicKey),
		ephemeralAESKey,
		[]byte{},
	)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failure encrypting wrapped key: %v", err))
		return 2
	}
	combinedCiphertext := append(wrappedAESKey, wrappedTargetKey...)
	importCiphertext := base64.StdEncoding.EncodeToString(combinedCiphertext)
	// Parse all the key options
	data := map[string]interface{}{
		"ciphertext": importCiphertext,
	}
	for _, v := range args[2:] {
		parts := strings.Split(v, "=")
		data[parts[0]] = parts[1]
	}

	c.UI.Output("Submitting wrapped key to Vault transit.")
	// Finally, call import
	_, err = client.Logical().Write(path+"/keys/"+keyName+"/"+operation, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to call import:%v", err))
		return 3
	} else {
		c.UI.Output("Success!")
		return 0
	}
}

func fetchWrappingKey(c *BaseCommand, client *api.Client, path string) (any, error) {
	resp, err := client.Logical().Read(path + "/wrapping_key")
	if err != nil {
		return nil, fmt.Errorf("error fetching wrapping key: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("transit not mounted at %s: %v", path, err)
	}
	key, ok := resp.Data["public_key"]
	if !ok {
		c.UI.Error("could not find wrapping key")
	}
	keyBlock, _ := pem.Decode([]byte(key.(string)))
	parsedKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing wrapping key: %w", err)
	}
	return parsedKey, nil
}
