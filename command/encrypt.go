package command

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/argon2"
)

var _ cli.Command = (*EncryptCommand)(nil)

const vaultEncryptVersion = 1

type EncryptCommand struct {
	*BaseCommand

	outfile string
}

type argon2Params struct {
	Salt        []byte
	Memory      uint32
	Time        uint32
	Parallelism uint8
}

func (c *EncryptCommand) Synopsis() string {
	return "Encrypts a file using AES-256 encryption."
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

	passphrase := "my password"

	encrypted, err := encrypt(data, passphrase)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
		return 1
	}

	outfile := c.outfile
	if outfile == "" {
		outfile = "output.enc"
	}

	if err := os.WriteFile(outfile, encrypted, 0644); err != nil {
		c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
		return 1
	}

	return 0
}

func encrypt(dataToEncrypt []byte, passphrase string) ([]byte, error) {
	aad := argon2Params{
		Salt:        make([]byte, 16),
		Memory:      64 * 1024,
		Time:        3,
		Parallelism: 4,
	}
	if _, err := rand.Read(aad.Salt); err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(passphrase), aad.Salt, aad.Time, aad.Memory, aad.Parallelism, 32)

	// Create a new Cipher Block using key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating new cipher block: %s", err.Error())
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %s", err.Error())
	}

	// Create a nonce for GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error creating nonce: %s", err.Error())
	}

	var out bytes.Buffer
	if _, err := out.Write(aad.encode()); err != nil {
		return nil, err
	}

	// Encrypt and write to file
	ciphertext := aesGCM.Seal(nonce, nonce, dataToEncrypt, aad.encode())

	// TODO ensure output is not written to stdout

	if _, err := out.Write(ciphertext); err != nil {
		return nil, fmt.Errorf("error writing encrypted data to file: %w", err)
	}

	return out.Bytes(), nil
}

func (aad argon2Params) encode() []byte {
	// TODO split out version
	e := fmt.Sprintf("vault:v=%d:%x:t=%d:m=%d:p=%d$", vaultEncryptVersion, aad.Salt, aad.Time, aad.Memory, aad.Parallelism)

	return []byte(e)
}
