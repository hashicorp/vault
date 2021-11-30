package command

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"io/ioutil"

	"github.com/mitchellh/cli"
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

  Decrpyts an  AES-256 encrypted file.

  Decrpyt a single file:

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

	aad, encryptedData, err := parseAAD(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(aad)

	//encryptedData, err := io.ReadAll(input)
	//if err != nil {
	//	log.Fatal(err)
	//}
	fmt.Println(len(encryptedData))

	key, _, err := generateKeyAAD("my password", aad)
	fmt.Printf("decoded key: %x\n", key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error creating key: %s", err.Error()))
		return 1
	}

	err = decrypt(encryptedData, key, aad, c.outfile)
	if err != nil {
		// Confirm if this is the right code
		c.UI.Error(fmt.Sprintf("error decrypting file: %s", err.Error()))
		return 1
	}

	return 0
}

func parseAAD(input []byte) (additionalData, []byte, error) {
	var aad additionalData

	parts := bytes.SplitN(input, []byte("$"), 2)
	if len(parts) != 2 {
		return aad, nil, fmt.Errorf("unexpected number of parts")
	}

	if _, err := fmt.Sscanf(string(parts[0]), "vault:v=%d:%x:t=%d:m=%d:p=%d", &aad.Version, &aad.Salt, &aad.Time, &aad.Memory, &aad.Parallelism); err != nil {
		// TODO: check n
		return aad, nil, err
	}

	return aad, parts[1], nil
}

func decrypt(dataToDecrypt, key []byte, aad additionalData, outfile string) error {

	// Create a new Cipher Block using key
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("error creating new cipher block: %s", err.Error())
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("error creating GCM: %s", err.Error())
	}

	// Create a nonce from GCM
	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := dataToDecrypt[:nonceSize], dataToDecrypt[nonceSize:]

	// Decrypt and write to file
	e := encodeAAD(aad)
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, e)
	if err != nil {
		return fmt.Errorf("error decrypting cipher text: %s", err.Error())

	}

	// TODO ensure output is not written to stdout

	if outfile == "" {
		outfile = "input.txt"
	}
	err = ioutil.WriteFile(outfile, plaintext, 0644)
	if err != nil {
		return fmt.Errorf("error writing decrypted data to file: %s", err.Error())
	}

	return nil
}
