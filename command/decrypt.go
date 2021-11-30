package command

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"crypto/aes"
	"crypto/cipher"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/argon2"
)

var _ cli.Command = (*DecryptCommand)(nil)

type DecryptCommand struct {
	*BaseCommand

	outfile string
}

func (c *DecryptCommand) Synopsis() string {
	return "Decrypts an  AES-256 encrypted file."
}

func (c *DecryptCommand) Help() string {
	helpText := `
Usage: vault decrypt [options] [filename]

  Decrypts an  AES-256 encrypted file.

  Decrypt a single file:

      $ vault decrypt -o foo.txt foo.enc

`
	return strings.TrimSpace(helpText)
}

func (c *DecryptCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "out",
		Aliases: []string{"o"},
		Target:  &c.outfile,
		Default: "input.txt",
		Usage:   "Specify the name of the output decrypted file.",
	})

	return set
}

func (c *DecryptCommand) Run(args []string) int {
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

	filename := strings.TrimSpace(args[0])

	input, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	passphrase := "my password"

	plaintext, err := c.decrypt(input, passphrase)
	if err != nil {
		// Confirm if this is the right code
		c.UI.Error(fmt.Sprintf("error decrypting file: %s", err.Error()))
		return 1
	}

	// TODO ensure output is not written to stdout
	outfile := c.outfile
	if outfile == "" {
		outfile = "input.txt"
	}
	err = os.WriteFile(outfile, plaintext, 0644)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error writing decrypted data to file: %s", err.Error()))
		return 1
	}

	return 0
}

func (c *DecryptCommand) parseAAD(input []byte) (argon2Params, []byte, error) {
	var aad argon2Params

	parts := bytes.SplitN(input, []byte("$"), 2)
	if len(parts) != 2 {
		return aad, nil, fmt.Errorf("unexpected number of parts")
	}

	var version int
	if _, err := fmt.Sscanf(string(parts[0]), "vault:v=%d:%x:t=%d:m=%d:p=%d", &version, &aad.Salt, &aad.Time, &aad.Memory, &aad.Parallelism); err != nil {
		return aad, nil, err
	}

	if version != vaultEncryptVersion {
		return aad, nil, fmt.Errorf("unknown version: %d", version)
	}

	return aad, parts[1], nil
}

func (c *DecryptCommand) decrypt(rawDataToDecrypt []byte, passphrase string) ([]byte, error) {
	aad, dataToDecrypt, err := c.parseAAD(rawDataToDecrypt)
	if err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(passphrase), aad.Salt, aad.Time, aad.Memory, aad.Parallelism, 32)
	if err != nil {
		return nil, fmt.Errorf("error creating key: %w", err)
	}

	// Create a new Cipher Block using key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating new cipher block: %w", err)
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %w", err)
	}

	// Create a nonce from GCM
	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := dataToDecrypt[:nonceSize], dataToDecrypt[nonceSize:]

	// Decrypt and write to file
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, aad.encode())
	if err != nil {
		return nil, fmt.Errorf("error decrypting cipher text: %w", err)

	}

	return plaintext, nil
}
