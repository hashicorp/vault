package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) getGenerationParams(ctx context.Context, storage logical.Storage, data *framework.FieldData, mountPoint string) (exported bool, format string, role *roleEntry, errorResp *logical.Response) {
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

	keyType, keyBits, err := getKeyTypeAndBitsForRole(ctx, b, storage, data, mountPoint)
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
	}
	*role.AllowWildcardCertificates = true

	if role.KeyBits, role.SignatureBits, err = certutil.ValidateDefaultOrValueKeyTypeSignatureLength(role.KeyType, role.KeyBits, role.SignatureBits); err != nil {
		errorResp = logical.ErrorResponse(err.Error())
	}

	return
}

func generateCABundle(ctx context.Context, b *backend, input *inputBundle, data *certutil.CreationBundle, randomSource io.Reader) (*certutil.ParsedCertBundle, error) {
	if kmsRequested(input) {
		return generateManagedKeyCABundle(ctx, b, input, data, randomSource)
	}

	if existingKeyRequested(input) {
		keyRef, err := getKeyRefWithErr(input.apiData)
		if err != nil {
			return nil, err
		}
		return certutil.CreateCertificateWithKeyGenerator(data, randomSource, existingGeneratePrivateKey(ctx, input.req.Storage, keyRef))
	}

	return certutil.CreateCertificateWithRandomSource(data, randomSource)
}

func generateCSRBundle(ctx context.Context, b *backend, input *inputBundle, data *certutil.CreationBundle, addBasicConstraints bool, randomSource io.Reader) (*certutil.ParsedCSRBundle, error) {
	if kmsRequested(input) {
		return generateManagedKeyCSRBundle(ctx, b, input, data, addBasicConstraints, randomSource)
	}

	if existingKeyRequested(input) {
		keyRef, err := getKeyRefWithErr(input.apiData)
		if err != nil {
			return nil, err
		}
		return certutil.CreateCSRWithKeyGenerator(data, addBasicConstraints, randomSource, existingGeneratePrivateKey(ctx, input.req.Storage, keyRef))
	}

	return certutil.CreateCSRWithRandomSource(data, addBasicConstraints, randomSource)
}

func parseCABundle(ctx context.Context, b *backend, req *logical.Request, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	if bundle.PrivateKeyType == certutil.ManagedPrivateKey {
		return parseManagedKeyCABundle(ctx, b, req, bundle)
	}
	return bundle.ToParsedCertBundle()
}

func getKeyTypeAndBitsForRole(ctx context.Context, b *backend, storage logical.Storage, data *framework.FieldData, mountPoint string) (string, int, error) {
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
		pubKeyManagedKey, err := getManagedKeyPublicKey(ctx, b, data, mountPoint)
		if err != nil {
			return "", 0, errors.New("failed to lookup public key from managed key: " + err.Error())
		}
		pubKey = pubKeyManagedKey
	}

	if existingKeyRequestedFromFieldData(data) {
		existingPubKey, err := getExistingPublicKey(ctx, storage, data)
		if err != nil {
			return "", 0, errors.New("failed to lookup public key from existing key: " + err.Error())
		}
		pubKey = existingPubKey
	}

	return getKeyTypeAndBitsFromPublicKeyForRole(pubKey)
}

func getExistingPublicKey(ctx context.Context, s logical.Storage, data *framework.FieldData) (crypto.PublicKey, error) {
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
	signer, err := key.GetSigner()
	if err != nil {
		return nil, err
	}
	return signer.Public(), nil
}

func getKeyTypeAndBitsFromPublicKeyForRole(pubKey crypto.PublicKey) (string, int, error) {
	var keyType string
	var keyBits int

	switch pubKey.(type) {
	case *rsa.PublicKey:
		keyType = "rsa"
		keyBits = certutil.GetPublicKeySize(pubKey)
	case *ecdsa.PublicKey:
		keyType = "ec"
	case *ed25519.PublicKey:
		keyType = "ed25519"
	default:
		return "", 0, fmt.Errorf("unsupported public key: %#v", pubKey)
	}
	return keyType, keyBits, nil
}

func getManagedKeyPublicKey(ctx context.Context, b *backend, data *framework.FieldData, mountPoint string) (crypto.PublicKey, error) {
	keyId, err := getManagedKeyId(data)
	if err != nil {
		return nil, errors.New("unable to determine managed key id")
	}
	// Determine key type and key bits from the managed public key
	var pubKey crypto.PublicKey
	err = withManagedPKIKey(ctx, b, keyId, mountPoint, func(ctx context.Context, key logical.ManagedSigningKey) error {
		pubKey, err = key.GetPublicKey(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.New("failed to lookup public key from managed key: " + err.Error())
	}
	return pubKey, nil
}

func existingGeneratePrivateKey(ctx context.Context, s logical.Storage, keyRef string) certutil.KeyGenerator {
	return func(keyType string, keyBits int, container certutil.ParsedPrivateKeyContainer, _ io.Reader) error {
		keyId, err := resolveKeyReference(ctx, s, keyRef)
		if err != nil {
			return err
		}
		key, err := fetchKeyById(ctx, s, keyId)
		if err != nil {
			return err
		}
		signer, err := key.GetSigner()
		if err != nil {
			return err
		}
		privateKeyType := certutil.GetPrivateKeyTypeFromSigner(signer)
		if privateKeyType == certutil.UnknownPrivateKey {
			return errors.New("unknown private key type loaded from key id: " + keyId.String())
		}
		blk, _ := pem.Decode([]byte(key.PrivateKey))
		container.SetParsedPrivateKey(signer, privateKeyType, blk.Bytes)
		return nil
	}
}
