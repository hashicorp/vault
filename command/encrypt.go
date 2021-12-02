package command

import (
	"bytes"
	"encoding/base64"
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
	vaultEncryptVersion = 1
	vaultKeyName        = "my-transit-key"
	defaultPassphrase   = "my password"

	// See https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-argon2-04#section-4 for details on parameter selection
	argon2SaltLen     = 16
	argon2Memory      = 32 * 1024
	argon2Parallelism = 4
	argon2Time        = 20
)

type EncryptCommand struct {
	*BaseCommand

	outfile     string
	passphrase  string
	key         string
	path        string
	ciphertext  string
	flagDecrypt bool
	flagTransit bool
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
		Default: "",
		Usage:   "Specify the name of the output file.",
	})

	f.StringVar(&StringVar{
		Name:    "pass",
		Aliases: []string{"p"},
		Target:  &c.passphrase,
		Default: defaultPassphrase,
		Usage:   "Specify a passphrase to encrypt/decrpyt this file.",
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

	f.BoolVar(&BoolVar{
		Name:    "transit",
		Aliases: []string{"t"},
		Target:  &c.flagTransit,
		Default: false,
		Usage:   "Enables data key generation for encryption using Vault Transit.",
	})

	f.StringVar(&StringVar{
		Name:    "key",
		Aliases: []string{"k"},
		Target:  &c.key,
		Default: "",
		Usage:   "Name of the datakey in Vault Transit.",
	})

	f.StringVar(&StringVar{
		Name:    "cipher",
		Aliases: []string{"c"},
		Target:  &c.ciphertext,
		Default: "",
		Usage:   "Encrypted transit ciphertext used for decryption. Can only be used if both decrypt and transit flags are set.",
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
	rawData, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var processedData []byte
	outfile := c.outfile

	passphrase := c.passphrase
	if passphrase == "" {
		passphrase = defaultPassphrase
	}
	if !c.flagDecrypt {
		// Encryption Mode

		if c.flagTransit {
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

			processedData, err = c.encrypt(rawData, passphrase, key)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
				return 1
			}
			fmt.Println(encryptedKey)
		} else {
			// Create key using passphrase
			processedData, err = c.encrypt(rawData, passphrase, nil)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
				return 1
			}
		}

		if outfile == "" {
			outfile = "output.enc"
		}

	} else {
		// Decryption Mode

		if c.flagTransit {
			// Fetch data key from Vault Transit
			if c.key == "" {
				c.UI.Error(fmt.Sprintf("key name not provided for Vault Transit"))
				return 1
			}

			key, _, err := c.fetchDatakey(c.key, c.ciphertext)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error generating transit data key: %s", err.Error()))
				return 1
			}

			processedData, err = c.decrypt(rawData, passphrase, key)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error encrypting file: %s", err.Error()))
				return 1
			}

		} else {
			// Create key using passphrase
			processedData, err = c.decrypt(rawData, passphrase, nil)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error decrypting file: %s", err.Error()))
				return 1
			}
		}

		if outfile == "" {
			outfile = "decoded.txt"
		}

	}

	// Write processed data to file
	if err := os.WriteFile(outfile, processedData, 0644); err != nil {
		c.UI.Error(fmt.Sprintf("error writing processed data to file: %s", err.Error()))
		return 1
	}

	return 0
}

func (c *EncryptCommand) encrypt(dataToEncrypt []byte, passphrase string, key []byte) ([]byte, error) {
	aad := argon2Params{
		Salt:        make([]byte, argon2SaltLen),
		Memory:      argon2Memory,
		Time:        argon2Time,
		Parallelism: argon2Parallelism,
	}
	if _, err := rand.Read(aad.Salt); err != nil {
		return nil, err
	}

	if key == nil {
		key = argon2.IDKey([]byte(passphrase), aad.Salt, aad.Time, aad.Memory, aad.Parallelism, 32)
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

func (c *EncryptCommand) parseAAD(input []byte) (argon2Params, []byte, error) {
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

func (c *EncryptCommand) decrypt(rawDataToDecrypt []byte, passphrase string, key []byte) ([]byte, error) {
	aad, dataToDecrypt, err := c.parseAAD(rawDataToDecrypt)
	if err != nil {
		return nil, err
	}

	if key == nil {
		key = argon2.IDKey([]byte(passphrase), aad.Salt, aad.Time, aad.Memory, aad.Parallelism, 32)
	}
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
