package command

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/argon2"
)

var _ cli.Command = (*EncryptCommand)(nil)

const (
	vaultEncryptionHeader = "vault-v1"
	encryptedFileSuffix   = ".enc"

	// See https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-argon2-04#section-4 for details on parameter selection
	argon2SaltLen     = 16
	argon2Memory      = 32 * 1024
	argon2Parallelism = 4
	argon2Time        = 20
)

type EncryptCommand struct {
	*BaseCommand

	outfile     string
	key         string
	path        string
	flagDecrypt bool
}

type argon2Params struct {
	Salt        []byte
	Memory      uint32
	Time        uint32
	Parallelism uint8
}

func (aad argon2Params) encode() string {
	return fmt.Sprintf("%x:t=%d:m=%d:p=%d", aad.Salt, aad.Time, aad.Memory, aad.Parallelism)
}

func (c *EncryptCommand) Synopsis() string {
	return "Encrypts a file using AES-256 encryption."
}

func (c *EncryptCommand) Help() string {
	helpText := `
Usage: vault encrypt [options] [filename]

  Encrypts a file using AES encryption. File can be encrypted using a derived key 
  using a passphrase or a datakey generated using Vault Transit.

  Encrypt a single file using a passphrase:

	  $ vault encrypt -o secrets.enc foo.txt
	
  Decrypt a single file using a passphrase:

	  $ vault encrypt -d -o decoded.txt secrets.enc

  Encrypt a single file using Transit datakey:

	  $ vault encrypt -k my-key -o secrets.enc foo.txt

  Decrypt a single file using a passphrase:

	  $ vault encrypt -d -k my-key -o decoded.txt secrets.enc

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *EncryptCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "out",
		Aliases: []string{"o"},
		Target:  &c.outfile,
		Usage:   "Specify the name of the output file.",
	})

	f.StringVar(&StringVar{
		Name:    "path",
		Target:  &c.path,
		Default: "transit",
		Usage:   "Path to mount used for transit.",
	})

	f.BoolVar(&BoolVar{
		Name:    "decrypt",
		Aliases: []string{"d"},
		Target:  &c.flagDecrypt,
		Default: false,
		Usage:   "Enables AES-256 decrypt mode.",
	})

	f.StringVar(&StringVar{
		Name:    "key",
		Aliases: []string{"k"},
		Target:  &c.key,
		Usage:   "Name of the data key in Vault Transit.",
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

	// Try to intelligently determine output file if not provided
	outfile := c.outfile
	if outfile == "" {
		if c.flagDecrypt {
			if strings.HasSuffix(filename, encryptedFileSuffix) {
				outfile = strings.TrimSuffix(filename, encryptedFileSuffix)
			} else {
				c.UI.Error(fmt.Sprintf("Unknown suffix. Provide the output filename with -out."))
				return 1
			}
		} else {
			outfile = filename + encryptedFileSuffix
		}
	}

	rawData, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var processedData []byte

	// Request either a passphrase or the encrypted datakey
	var passphrase string
	var ciphertext string

	datakeyMode := c.key != ""

	// No provided key name means passphrase mode
	if !datakeyMode {
		passphrase, err = c.UI.AskSecret("Passphrase (will be hidden):")

		if err != nil {
			c.UI.Error(fmt.Sprintf("Error reading input: %s", err.Error()))
			return 2
		}
	} else {
		if c.flagDecrypt {
			ciphertext, err = c.UI.AskSecret("Encrypted data key (will be hidden):")

			if err != nil {
				c.UI.Error(fmt.Sprintf("Error reading input: %s", err.Error()))
				return 2
			}
		}
	}

	if !c.flagDecrypt {
		// Encryption Mode
		if datakeyMode {
			// Fetch data key from Vault Transit
			if c.key == "" {
				c.UI.Error(fmt.Sprintf("key name not provided for Vault Transit"))
				return 1
			}

			key, encryptedKey, err := c.fetchDatakey(c.key, "")
			if err != nil {
				c.UI.Error(fmt.Sprintf("error generating transit data key: %s", err.Error()))
				return 1
			}

			processedData, err = c.encrypt(rawData, key, false)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
				return 1
			}
			c.UI.Output("Encrypted data key:\n\n    " + encryptedKey)
		} else {
			// Create key using passphrase
			processedData, err = c.encrypt(rawData, []byte(passphrase), true)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
				return 1
			}
		}
	} else {
		// Decryption Mode

		if datakeyMode {
			key, _, err := c.fetchDatakey(c.key, ciphertext)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error generating transit data key: %s", err.Error()))
				return 1
			}

			processedData, err = c.decrypt(rawData, key, false)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
				return 1
			}

		} else {
			// Create key using passphrase
			processedData, err = c.decrypt(rawData, []byte(passphrase), true)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error decrypting file: %s", err.Error()))
				return 1
			}
		}
	}

	// Write processed data to file
	if err := os.WriteFile(outfile, processedData, 0644); err != nil {
		c.UI.Error(fmt.Sprintf("error writing processed data to file: %s", err.Error()))
		return 1
	}

	c.UI.Output("\nOutput written to: " + outfile)

	return 0
}

func (c *EncryptCommand) encrypt(dataToEncrypt []byte, key []byte, passphrase bool) ([]byte, error) {
	var aad string
	var out bytes.Buffer

	out.WriteString(vaultEncryptionHeader + "$")

	if passphrase {
		params := argon2Params{
			Salt:        make([]byte, argon2SaltLen),
			Memory:      argon2Memory,
			Time:        argon2Time,
			Parallelism: argon2Parallelism,
		}
		if _, err := rand.Read(params.Salt); err != nil {
			return nil, err
		}
		key = argon2.IDKey(key, params.Salt, params.Time, params.Memory, params.Parallelism, 32)
		aad = params.encode()
		out.WriteString(aad + "$")
	}

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

	ciphertext := aesGCM.Seal(nonce, nonce, dataToEncrypt, []byte(aad))
	out.Write(ciphertext)

	return out.Bytes(), nil
}

func (c *EncryptCommand) parseAAD(input []byte) (argon2Params, []byte, error) {
	var aad argon2Params

	parts := bytes.SplitN(input, []byte("$"), 2)
	if len(parts) != 2 {
		return aad, nil, fmt.Errorf("unexpected number of parts")
	}

	if _, err := fmt.Sscanf(string(parts[0]), "%x:t=%d:m=%d:p=%d", &aad.Salt, &aad.Time, &aad.Memory, &aad.Parallelism); err != nil {
		return aad, nil, err
	}

	return aad, parts[1], nil
}

func (c *EncryptCommand) decrypt(rawDataToDecrypt []byte, key []byte, passphrase bool) ([]byte, error) {
	if !bytes.HasPrefix(rawDataToDecrypt, []byte(vaultEncryptionHeader+"$")) {
		return nil, errors.New("unrecognized file header")
	}

	rawDataToDecrypt = bytes.TrimPrefix(rawDataToDecrypt, []byte(vaultEncryptionHeader+"$"))

	var dataToDecrypt = rawDataToDecrypt
	var aad string

	if passphrase {
		aad2, data, err := c.parseAAD(rawDataToDecrypt)
		if err != nil {
			return nil, fmt.Errorf("error creating key: %w", err)
		}
		dataToDecrypt = data
		aad = aad2.encode()
		key = argon2.IDKey(key, aad2.Salt, aad2.Time, aad2.Memory, aad2.Parallelism, 32)
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
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, []byte(aad))
	if err != nil {
		return nil, fmt.Errorf("error decrypting cipher text: %w", err)

	}

	return plaintext, nil
}

func (c *EncryptCommand) fetchDatakey(name string, ciphertext string) ([]byte, string, error) {
	client, err := c.Client()
	if err != nil {
		return nil, "", fmt.Errorf("Error initializing client: %s", err)
	}
	var key []byte

	if ciphertext == "" {
		// Encryption mode
		data := map[string]interface{}{
			"bits": 256,
		}

		secret, err := client.Logical().Write(fmt.Sprintf("%s/datakey/plaintext/%s", c.path, name), data)
		if err != nil {
			return nil, "", fmt.Errorf("Error making request for datakey: %s", err)
		}

		key, err = base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string)) // TODO: check assertion
		if err != nil {
			return nil, "", fmt.Errorf("Error b64 encoding plaintext: %s", err)
		}

		return key, secret.Data["ciphertext"].(string), nil
	} else {
		// Decrypt ciphertext
		data := map[string]interface{}{
			"ciphertext": ciphertext,
		}

		secret, err := client.Logical().Write(fmt.Sprintf("%s/decrypt/%s", c.path, name), data)
		if err != nil {
			return nil, "", fmt.Errorf("Error decrypting transit key: %s", err)
		}

		key, err = base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string)) // TODO: check assertion
		if err != nil {
			return nil, "", fmt.Errorf("Error b64 encoding plaintext: %s", err)
		}

		return key, "", nil
	}

}
