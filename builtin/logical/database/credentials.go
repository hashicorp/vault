// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/template"
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

	key, err := cryptoutil.GenerateRSAKey(reader, keyBits)
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

type ClientCertificateGenerator struct {
	// CommonNameTemplate is username template to be used for the client certificate common name.
	CommonNameTemplate string `mapstructure:"common_name_template,omitempty"`

	// CAPrivateKey is the PEM-encoded private key for the given ca_cert.
	CAPrivateKey string `mapstructure:"ca_private_key,omitempty"`

	// CACert is the PEM-encoded CA certificate.
	CACert string `mapstructure:"ca_cert,omitempty"`

	// KeyType specifies the desired key type.
	// Options include: 'rsa', 'ed25519', 'ec'.
	KeyType string `mapstructure:"key_type,omitempty"`

	// KeyBits is the number of bits to use for the generated keys.
	// Options include: with key_type=rsa, 2048 (default), 3072, 4096;
	// With key_type=ec, allowed values are: 224, 256 (default), 384, 521;
	// Ignored with key_type=ed25519.
	KeyBits int `mapstructure:"key_bits,omitempty"`

	// SignatureBits is the number of bits to use in the signature algorithm.
	// Options include: 256 (default), 384, 512.
	SignatureBits int `mapstructure:"signature_bits,omitempty"`

	parsedCABundle *certutil.ParsedCertBundle
	cnProducer     template.StringTemplate
}

// newClientCertificateGenerator returns a new ClientCertificateGenerator
// using the given config. Default values will be set on the returned
// ClientCertificateGenerator if not provided in the config.
func newClientCertificateGenerator(config map[string]interface{}) (ClientCertificateGenerator, error) {
	var cg ClientCertificateGenerator
	if err := mapstructure.WeakDecode(config, &cg); err != nil {
		return cg, err
	}

	switch cg.KeyType {
	case "rsa":
		switch cg.KeyBits {
		case 0:
			cg.KeyBits = 2048
		case 2048, 3072, 4096:
		default:
			return cg, fmt.Errorf("invalid key_bits")
		}
	case "ec":
		switch cg.KeyBits {
		case 0:
			cg.KeyBits = 256
		case 224, 256, 384, 521:
		default:
			return cg, fmt.Errorf("invalid key_bits")
		}
	case "ed25519":
	// key_bits ignored
	default:
		return cg, fmt.Errorf("invalid key_type")
	}

	switch cg.SignatureBits {
	case 0:
		cg.SignatureBits = 256
	case 256, 384, 512:
	default:
		return cg, fmt.Errorf("invalid signature_bits")
	}

	if cg.CommonNameTemplate == "" {
		return cg, fmt.Errorf("missing required common_name_template")
	}

	// Validate the common name template
	t, err := template.NewTemplate(template.Template(cg.CommonNameTemplate))
	if err != nil {
		return cg, fmt.Errorf("failed to create template: %w", err)
	}

	_, err = t.Generate(dbplugin.UsernameMetadata{})
	if err != nil {
		return cg, fmt.Errorf("invalid common_name_template: %w", err)
	}
	cg.cnProducer = t

	if cg.CACert == "" {
		return cg, fmt.Errorf("missing required ca_cert")
	}
	if cg.CAPrivateKey == "" {
		return cg, fmt.Errorf("missing required ca_private_key")
	}
	parsedBundle, err := certutil.ParsePEMBundle(strings.Join([]string{cg.CACert, cg.CAPrivateKey}, "\n"))
	if err != nil {
		return cg, err
	}
	if parsedBundle.PrivateKey == nil {
		return cg, fmt.Errorf("private key not found in the PEM bundle")
	}
	if parsedBundle.PrivateKeyType == certutil.UnknownPrivateKey {
		return cg, fmt.Errorf("unknown private key found in the PEM bundle")
	}
	if parsedBundle.Certificate == nil {
		return cg, fmt.Errorf("certificate not found in the PEM bundle")
	}
	if !parsedBundle.Certificate.IsCA {
		return cg, fmt.Errorf("the given certificate is not marked for CA use")
	}
	if !parsedBundle.Certificate.BasicConstraintsValid {
		return cg, fmt.Errorf("the given certificate does not meet basic constraints for CA use")
	}

	certBundle, err := parsedBundle.ToCertBundle()
	if err != nil {
		return cg, fmt.Errorf("error converting raw values into cert bundle: %w", err)
	}

	parsedCABundle, err := certBundle.ToParsedCertBundle()
	if err != nil {
		return cg, fmt.Errorf("failed to parse cert bundle: %w", err)
	}
	cg.parsedCABundle = parsedCABundle

	return cg, nil
}

func (cg *ClientCertificateGenerator) generate(r io.Reader, expiration time.Time, userMeta dbplugin.UsernameMetadata) (*certutil.CertBundle, string, error) {
	commonName, err := cg.cnProducer.Generate(userMeta)
	if err != nil {
		return nil, "", err
	}

	// Set defaults
	keyBits := cg.KeyBits
	signatureBits := cg.SignatureBits
	switch cg.KeyType {
	case "rsa":
		if keyBits == 0 {
			keyBits = 2048
		}
		if signatureBits == 0 {
			signatureBits = 256
		}
	case "ec":
		if keyBits == 0 {
			keyBits = 256
		}
		if signatureBits == 0 {
			if keyBits == 224 {
				signatureBits = 256
			} else {
				signatureBits = keyBits
			}
		}
	case "ed25519":
		// key_bits ignored
		if signatureBits == 0 {
			signatureBits = 256
		}
	}

	subject := pkix.Name{
		CommonName: commonName,
		// Additional subject DN options intentionally omitted for now
	}

	creation := &certutil.CreationBundle{
		Params: &certutil.CreationParameters{
			Subject:                       subject,
			KeyType:                       cg.KeyType,
			KeyBits:                       cg.KeyBits,
			SignatureBits:                 cg.SignatureBits,
			NotAfter:                      expiration,
			KeyUsage:                      x509.KeyUsageDigitalSignature,
			ExtKeyUsage:                   certutil.ClientAuthExtKeyUsage,
			BasicConstraintsValidForNonCA: false,
			NotBeforeDuration:             30 * time.Second,
			URLs: &certutil.URLEntries{
				IssuingCertificates:   []string{},
				CRLDistributionPoints: []string{},
				OCSPServers:           []string{},
			},
		},
		SigningBundle: &certutil.CAInfoBundle{
			ParsedCertBundle: *cg.parsedCABundle,
			URLs: &certutil.URLEntries{
				IssuingCertificates:   []string{},
				CRLDistributionPoints: []string{},
				OCSPServers:           []string{},
			},
		},
	}

	parsedClientBundle, err := certutil.CreateCertificateWithRandomSource(creation, r)
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate client certificate: %w", err)
	}

	cb, err := parsedClientBundle.ToCertBundle()
	if err != nil {
		return nil, "", fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}

	return cb, subject.String(), nil
}

// configMap returns the configuration of the ClientCertificateGenerator
// as a map from string to string.
func (cg ClientCertificateGenerator) configMap() (map[string]interface{}, error) {
	config := make(map[string]interface{})
	if err := mapstructure.WeakDecode(cg, &config); err != nil {
		return nil, err
	}
	return config, nil
}
