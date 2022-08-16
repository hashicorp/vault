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
	"github.com/google/tink/go/kwp/subtle"
	"github.com/hashicorp/errwrap"
	"path"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TransitImportCommand)(nil)
	_ cli.CommandAutocomplete = (*TransitImportCommand)(nil)
)

type TransitImportCommand struct {
	*BaseCommand
}

func (c *TransitImportCommand) Synopsis() string {
	return "Imports key material into a new Transit key"
}

func (c *TransitImportCommand) Help() string {
	helpText := `
Usage: vault transit import [options]

  

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *TransitImportCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *TransitImportCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitImportCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

var keyPath = regexp.MustCompile("^(.*)/keys/([^/]*)$")

func (c *TransitImportCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 2 {
		c.UI.Error(fmt.Sprintf("Too few arguments (expected 3, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	return doImport(c.UI, "import", args, client)
}

func doImport(ui cli.Ui, operation string, args []string, client *api.Client) int {
	ephemeralAESKey := make([]byte, 32)
	_, err := rand.Read(ephemeralAESKey)
	if err != nil {
		ui.Error(fmt.Sprintf("failed to generate ephemeral key: %s", err.Error()))
		return 3
	}
	parts := keyPath.FindStringSubmatch(args[0])
	if len(parts) != 3 {
		ui.Error("expected transit path and key name in the form :path:/keys/:name:")
		return 1
	}
	mountPath := parts[1]
	keyName := parts[2]
	key, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		ui.Error(fmt.Sprintf("error base64 decoding source key material: %s", err.Error()))
		return 1
	}
	// Fetch the wrapping key
	ui.Info("Retrieving transit wrapping key.")
	wrappingKey, err := fetchWrappingKey(client, mountPath)
	if err != nil {
		ui.Error(fmt.Sprintf("error fetching wrapping key: %s", err))
		return 2
	}
	ui.Info("Wrapping source key with ephemeral key.")
	wrapKWP, err := subtle.NewKWP(ephemeralAESKey)
	if err != nil {
		ui.Error(fmt.Sprintf("failure building key wrapping key: %s", err.Error()))
	}
	wrappedTargetKey, err := wrapKWP.Wrap(key)
	if err != nil {
		ui.Error(fmt.Sprintf("failure wrapping source key: %s", err.Error()))
	}
	ui.Info("Encrypting ephemeral key with transit wrapping key.")
	wrappedAESKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		wrappingKey.(*rsa.PublicKey),
		ephemeralAESKey,
		[]byte{},
	)
	if err != nil {
		ui.Error(fmt.Sprintf("failure encrypting wrapped key: %s", err.Error()))
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

	ui.Info("Submitting wrapped key to Vault transit.")
	// Finally, call import
	_, err = client.Logical().Write(path.Join(mountPath, "keys", keyName, operation), data)
	if err != nil {
		ui.Error(fmt.Sprintf("failed to call import: %s", err.Error()))
	} else {
		ui.Info("Success!")
	}
	return 0
}

func fetchWrappingKey(client *api.Client, path string) (any, error) {
	resp, err := client.Logical().Read(path + "/wrapping_key")
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("transit not mounted at %s", path)
	}
	key, ok := resp.Data["public_key"]
	if !ok {
		errors.New("could not find wrapping key")
	}
	keyBlock, _ := pem.Decode([]byte(key.(string)))
	parsedKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		return nil, errwrap.Wrap(errors.New("error parsing wrapping key"), err)
	}
	return parsedKey, nil
}
