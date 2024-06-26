// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/tink/go/kwp/subtle"
	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
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
	return "Import a key into the Transit secrets engines."
}

func (c *TransitImportCommand) Help() string {
	helpText := `
Usage: vault transit import PATH KEY [options...]

  Using the Transit key wrapping system, imports key material from
  the base64 encoded KEY (either directly on the CLI or via @path notation),
  into a new key whose API path is PATH.  To import a new version into an
  existing key, use import_version.  The remaining options after KEY (key=value
  style) are passed on to the Transit create key endpoint.  If your
  system or device natively supports the RSA AES key wrap mechanism (such as
  the PKCS#11 mechanism CKM_RSA_AES_KEY_WRAP), you should use it directly
  rather than this command.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransitImportCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *TransitImportCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitImportCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransitImportCommand) Run(args []string) int {
	return ImportKey(c.BaseCommand, "import", transitImportKeyPath, c.Flags(), args)
}

func transitImportKeyPath(s string, operation string) (path string, apiPath string, err error) {
	parts := keyPath.FindStringSubmatch(s)
	if len(parts) != 3 {
		return "", "", errors.New("expected transit path and key name in the form :path:/keys/:name:")
	}
	path = parts[1]
	keyName := parts[2]
	apiPath = path + "/keys/" + keyName + "/" + operation

	return path, apiPath, nil
}

type ImportKeyFunc func(s string, operation string) (path string, apiPath string, err error)

// error codes: 1: user error, 2: internal computation error, 3: remote api call error
func ImportKey(c *BaseCommand, operation string, pathFunc ImportKeyFunc, flags *FlagSets, args []string) int {
	// Parse and validate the arguments.
	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = flags.Args()
	if len(args) < 2 {
		c.UI.Error(fmt.Sprintf("Incorrect argument count (expected 2+, got %d). Wanted PATH to import into and KEY material.", len(args)))
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
	path, apiPath, err := pathFunc(args[0], operation)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	keyMaterial := args[1]
	if keyMaterial[0] == '@' {
		keyMaterialBytes, err := os.ReadFile(keyMaterial[1:])
		if err != nil {
			c.UI.Error(fmt.Sprintf("error reading key material file: %v", err))
			return 1
		}

		keyMaterial = string(keyMaterialBytes)
	}

	key, err := base64.StdEncoding.DecodeString(keyMaterial)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error base64 decoding source key material: %v", err))
		return 1
	}
	// Fetch the wrapping key
	c.UI.Output("Retrieving wrapping key.")
	wrappingKey, err := fetchWrappingKey(client, path)
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
	c.UI.Output("Encrypting ephemeral key with wrapping key.")
	wrappedAESKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		wrappingKey,
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
	data, err := parseArgsData(os.Stdin, args[2:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse extra K=V data: %s", err))
		return 1
	}
	if data == nil {
		data = make(map[string]interface{}, 1)
	}

	data["ciphertext"] = importCiphertext

	c.UI.Output("Submitting wrapped key.")
	// Finally, call import

	_, err = client.Logical().Write(apiPath, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to call import:%v", err))
		return 3
	} else {
		c.UI.Output("Success!")
		return 0
	}
}

func fetchWrappingKey(client *api.Client, path string) (*rsa.PublicKey, error) {
	resp, err := client.Logical().Read(path + "/wrapping_key")
	if err != nil {
		return nil, fmt.Errorf("error fetching wrapping key: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("no mount found at %s: %v", path, err)
	}
	key, ok := resp.Data["public_key"]
	if !ok {
		return nil, fmt.Errorf("missing public_key field in response")
	}
	keyBlock, _ := pem.Decode([]byte(key.(string)))
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM information from public_key response field")
	}
	parsedKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing wrapping key: %w", err)
	}
	rsaKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("returned value was not an RSA public key but a %T", rsaKey)
	}
	return rsaKey, nil
}
