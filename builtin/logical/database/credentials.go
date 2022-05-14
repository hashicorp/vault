package database

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/vault/helper/random"
	"github.com/mitchellh/mapstructure"
)

// passwordGenerator generates password credentials.
// A zero value passwordGenerator is usable.
type passwordGenerator struct {
	// PasswordPolicy is the named password policy used to generate passwords.
	// If empty (default), a random string of 20 characters will be generated.
	PasswordPolicy string `mapstructure:"password_policy,omitempty"`
}

// newPasswordGenerator returns a new passwordGenerator using the given config.
// Default values will be set on the returned passwordGenerator if not provided
// in the config.
func newPasswordGenerator(config map[string]interface{}) (passwordGenerator, error) {
	var pg passwordGenerator
	if err := mapstructure.WeakDecode(config, &pg); err != nil {
		return pg, err
	}

	return pg, nil
}

// Generate generates a password credential using the configured password policy.
// Returns the generated password or an error.
func (pg passwordGenerator) generate(ctx context.Context, b *databaseBackend, wrapper databaseVersionWrapper) (string, error) {
	if !wrapper.isV5() && !wrapper.isV4() {
		return "", fmt.Errorf("no underlying database specified")
	}

	// The database plugin generates the password if its interface is v4
	if wrapper.isV4() {
		password, err := wrapper.v4.GenerateCredentials(ctx)
		if err != nil {
			return "", err
		}
		return password, nil
	}

	if pg.PasswordPolicy == "" {
		return random.DefaultStringGenerator.Generate(ctx, b.GetRandomReader())
	}
	return b.System().GeneratePasswordFromPolicy(ctx, pg.PasswordPolicy)
}

// configMap returns the configuration of the passwordGenerator
// as a map from string to string.
func (pg passwordGenerator) configMap() (map[string]interface{}, error) {
	config := make(map[string]interface{})
	if err := mapstructure.WeakDecode(pg, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// rsaKeyGenerator generates RSA key pair credentials.
// A zero value rsaKeyGenerator is usable.
type rsaKeyGenerator struct {
	// Format is the output format of the generated private key.
	// Options include: 'pkcs8' (default)
	Format string `mapstructure:"format,omitempty"`

	// KeyBits is the bit size of the RSA key to generate.
	// Options include: 2048 (default), 3072, and 4096
	KeyBits int `mapstructure:"key_bits,omitempty"`
}

// newRSAKeyGenerator returns a new rsaKeyGenerator using the given config.
// Default values will be set on the returned rsaKeyGenerator if not provided
// in the given config.
func newRSAKeyGenerator(config map[string]interface{}) (rsaKeyGenerator, error) {
	var kg rsaKeyGenerator
	if err := mapstructure.WeakDecode(config, &kg); err != nil {
		return kg, err
	}

	switch strings.ToLower(kg.Format) {
	case "":
		kg.Format = "pkcs8"
	case "pkcs8":
	default:
		return kg, fmt.Errorf("invalid format: %v", kg.Format)
	}

	switch kg.KeyBits {
	case 0:
		kg.KeyBits = 2048
	case 2048, 3072, 4096:
	default:
		return kg, fmt.Errorf("invalid key_bits: %v", kg.KeyBits)
	}

	return kg, nil
}

// Generate generates an RSA key pair. Returns a PEM-encoded, PKIX marshaled
// public key and a PEM-encoded private key marshaled into the configuration
// format (in that order) or an error.
func (kg *rsaKeyGenerator) generate(r io.Reader) ([]byte, []byte, error) {
	reader := rand.Reader
	if r != nil {
		reader = r
	}

	var keyBits int
	switch kg.KeyBits {
	case 0:
		keyBits = 2048
	case 2048, 3072, 4096:
		keyBits = kg.KeyBits
	default:
		return nil, nil, fmt.Errorf("invalid key_bits: %v", kg.KeyBits)
	}

	key, err := rsa.GenerateKey(reader, keyBits)
	if err != nil {
		return nil, nil, err
	}

	public, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return nil, nil, err
	}

	var private []byte
	switch strings.ToLower(kg.Format) {
	case "", "pkcs8":
		private, err = x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("invalid format: %v", kg.Format)
	}

	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: public,
	}
	privateBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: private,
	}

	return pem.EncodeToMemory(publicBlock), pem.EncodeToMemory(privateBlock), nil
}

// configMap returns the configuration of the rsaKeyGenerator
// as a map from string to string.
func (kg rsaKeyGenerator) configMap() (map[string]interface{}, error) {
	config := make(map[string]interface{})
	if err := mapstructure.WeakDecode(kg, &config); err != nil {
		return nil, err
	}
	return config, nil
}
