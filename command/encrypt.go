package command

import (
	"fmt"
	"log"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*EncryptCommand)(nil)

type EncryptCommand struct {
	*BaseCommand

	outfile string
}

func (c *EncryptCommand) Synopsis() string {
	return "Encrypts a file using AES 128bit or 256bit encryption."
}

func (c *EncryptCommand) Help() string {
	helpText := `
Usage: vault encrypt [options] [filename]

  Encrypts a file using AES encryption.

  Encrypt a single file:

      $ vault encrypt -o foo.enc foo.txt

`
	return strings.TrimSpace(helpText)
}

func (c *EncryptCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "out",
		Aliases: []string{"o"},
		Target:  &c.outfile,
		Default: "output.enc",
		Usage:   "Specify the name of the output file.",
	})

	return set
}

func (c *EncryptCommand) Run(args []string) int {
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
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// TODO Prompt w/ passphrase
	// TODO Create key w/ passphrase
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		fmt.Errorf(err.Error())
	}

	encrypt(data, key)

	return 0
}

func encrypt(dataToEncrypt []byte, key []byte) {

	// Create a new Cipher Block using key
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	// Create a nonce from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Errorf(err.Error())
	}

	// Encrypt and write to file
	ciphertext := aesGCM.Seal(nonce, nonce, dataToEncrypt, nil)

	// TODO ensure output is not written to stdout
	err = ioutil.WriteFile("output.enc", ciphertext, 0644)
	if err != nil {
		fmt.Errorf(err.Error())
	}

}
