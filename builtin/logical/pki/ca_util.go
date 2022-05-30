package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) getGenerationParams(ctx context.Context, storage logical.Storage, data *framework.FieldData) (exported bool, format string, role *roleEntry, errorResp *logical.Response) {
	exportedStr := data.Get("exported").(string)
	switch exportedStr {
	case "exported":
		exported = true
	case "internal":
	case "existing":
	case "kms":
	default:
		errorResp = logical.ErrorResponse(
			`the "exported" path parameter must be "internal", "existing", exported" or "kms"`)
		return
	}

	format = getFormat(data)
	if format == "" {
		errorResp = logical.ErrorResponse(
			`the "format" path parameter must be "pem", "der", or "pem_bundle"`)
		return
	}
	keyType, keyBits, err := getKeyTypeAndBitsForRole(ctx, b, storage, data)
	if err != nil {
		errorResp = logical.ErrorResponse(err.Error())
		return
	}

	role = &roleEntry{
		TTL:                       time.Duration(data.Get("ttl").(int)) * time.Second,
		KeyType:                   keyType,
		KeyBits:                   keyBits,
		SignatureBits:             data.Get("signature_bits").(int),
		AllowLocalhost:            true,
		AllowAnyName:              true,
		AllowIPSANs:               true,
		AllowWildcardCertificates: new(bool),
		EnforceHostnames:          false,
		AllowedURISANs:            []string{"*"},
		AllowedOtherSANs:          []string{"*"},
		AllowedSerialNumbers:      []string{"*"},
		OU:                        data.Get("ou").([]string),
		Organization:              data.Get("organization").([]string),
		Country:                   data.Get("country").([]string),
		Locality:                  data.Get("locality").([]string),
		Province:                  data.Get("province").([]string),
		StreetAddress:             data.Get("street_address").([]string),
		PostalCode:                data.Get("postal_code").([]string),
		NotBeforeDuration:         time.Duration(data.Get("not_before_duration").(int)) * time.Second,
	}
	*role.AllowWildcardCertificates = true

	if role.KeyBits, role.SignatureBits, err = certutil.ValidateDefaultOrValueKeyTypeSignatureLength(role.KeyType, role.KeyBits, role.SignatureBits); err != nil {
		errorResp = logical.ErrorResponse(err.Error())
	}

	return
}

func generateCABundle(ctx context.Context, b *backend, input *inputBundle, data *certutil.CreationBundle, randomSource io.Reader) (*certutil.ParsedCertBundle, error) {
	if kmsRequested(input) {
		keyId, err := getManagedKeyId(input.apiData)
		if err != nil {
			return nil, err
		}
		return generateManagedKeyCABundle(ctx, b, keyId, data, randomSource)
	}

	if existingKeyRequested(input) {
		keyRef, err := getKeyRefWithErr(input.apiData)
		if err != nil {
			return nil, err
		}

		keyEntry, err := getExistingKeyFromRef(ctx, input.req.Storage, keyRef)
		if err != nil {
			return nil, err
		}

		if keyEntry.isManagedPrivateKey() {
			keyId, err := keyEntry.getManagedKeyUUID()
			if err != nil {
				return nil, err
			}
			return generateManagedKeyCABundle(ctx, b, keyId, data, randomSource)
		}

		return certutil.CreateCertificateWithKeyGenerator(data, randomSource, existingKeyGeneratorFromBytes(keyEntry))
	}

	return certutil.CreateCertificateWithRandomSource(data, randomSource)
}

func generateCSRBundle(ctx context.Context, b *backend, input *inputBundle, data *certutil.CreationBundle, addBasicConstraints bool, randomSource io.Reader) (*certutil.ParsedCSRBundle, error) {
	if kmsRequested(input) {
		keyId, err := getManagedKeyId(input.apiData)
		if err != nil {
			return nil, err
		}

		return generateManagedKeyCSRBundle(ctx, b, keyId, data, addBasicConstraints, randomSource)
	}

	if existingKeyRequested(input) {
		keyRef, err := getKeyRefWithErr(input.apiData)
		if err != nil {
			return nil, err
		}

		key, err := getExistingKeyFromRef(ctx, input.req.Storage, keyRef)
		if err != nil {
			return nil, err
		}

		if key.isManagedPrivateKey() {
			keyId, err := key.getManagedKeyUUID()
			if err != nil {
				return nil, err
			}
			return generateManagedKeyCSRBundle(ctx, b, keyId, data, addBasicConstraints, randomSource)
		}

		return certutil.CreateCSRWithKeyGenerator(data, addBasicConstraints, randomSource, existingKeyGeneratorFromBytes(key))
	}

	return certutil.CreateCSRWithRandomSource(data, addBasicConstraints, randomSource)
}

func parseCABundle(ctx context.Context, b *backend, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	if bundle.PrivateKeyType == certutil.ManagedPrivateKey {
		return parseManagedKeyCABundle(ctx, b, bundle)
	}
	return bundle.ToParsedCertBundle()
}

func getKeyTypeAndBitsForRole(ctx context.Context, b *backend, storage logical.Storage, data *framework.FieldData) (string, int, error) {
	exportedStr := data.Get("exported").(string)
	var keyType string
	var keyBits int

	switch exportedStr {
	case "internal":
		fallthrough
	case "exported":
		keyType = data.Get("key_type").(string)
		keyBits = data.Get("key_bits").(int)
		return keyType, keyBits, nil
	}

	// existing and kms types don't support providing the key_type and key_bits args.
	_, okKeyType := data.Raw["key_type"]
	_, okKeyBits := data.Raw["key_bits"]

	if okKeyType || okKeyBits {
		return "", 0, errors.New("invalid parameter for the kms/existing path parameter, key_type nor key_bits arguments can be set in this mode")
	}

	var pubKey crypto.PublicKey
	if kmsRequestedFromFieldData(data) {
		keyId, err := getManagedKeyId(data)
		if err != nil {
			return "", 0, errors.New("unable to determine managed key id" + err.Error())
		}

		pubKeyManagedKey, err := getManagedKeyPublicKey(ctx, b, keyId)
		if err != nil {
			return "", 0, errors.New("failed to lookup public key from managed key: " + err.Error())
		}
		pubKey = pubKeyManagedKey
	}

	if existingKeyRequestedFromFieldData(data) {
		existingPubKey, err := getExistingPublicKey(ctx, b, storage, data)
		if err != nil {
			return "", 0, errors.New("failed to lookup public key from existing key: " + err.Error())
		}
		pubKey = existingPubKey
	}

	privateKeyType, keyBits, err := getKeyTypeAndBitsFromPublicKeyForRole(pubKey)
	return string(privateKeyType), keyBits, err
}

func getExistingPublicKey(ctx context.Context, b *backend, s logical.Storage, data *framework.FieldData) (crypto.PublicKey, error) {
	keyRef, err := getKeyRefWithErr(data)
	if err != nil {
		return nil, err
	}
	id, err := resolveKeyReference(ctx, s, keyRef)
	if err != nil {
		return nil, err
	}
	key, err := fetchKeyById(ctx, s, id)
	if err != nil {
		return nil, err
	}
	return getPublicKey(ctx, b, key)
}

func getKeyTypeAndBitsFromPublicKeyForRole(pubKey crypto.PublicKey) (certutil.PrivateKeyType, int, error) {
	var keyType certutil.PrivateKeyType
	var keyBits int

	switch pubKey.(type) {
	case *rsa.PublicKey:
		keyType = certutil.RSAPrivateKey
		keyBits = certutil.GetPublicKeySize(pubKey)
	case *ecdsa.PublicKey:
		keyType = certutil.ECPrivateKey
	case *ed25519.PublicKey:
		keyType = certutil.Ed25519PrivateKey
	default:
		return certutil.UnknownPrivateKey, 0, fmt.Errorf("unsupported public key: %#v", pubKey)
	}
	return keyType, keyBits, nil
}

func getExistingKeyFromRef(ctx context.Context, s logical.Storage, keyRef string) (*keyEntry, error) {
	keyId, err := resolveKeyReference(ctx, s, keyRef)
	if err != nil {
		return nil, err
	}
	return fetchKeyById(ctx, s, keyId)
}

func existingKeyGeneratorFromBytes(key *keyEntry) certutil.KeyGenerator {
	return func(_ string, _ int, container certutil.ParsedPrivateKeyContainer, _ io.Reader) error {
		signer, _, pemBytes, err := getSignerFromKeyEntryBytes(key)
		if err != nil {
			return err
		}

		container.SetParsedPrivateKey(signer, key.PrivateKeyType, pemBytes.Bytes)
		return nil
	}
}
