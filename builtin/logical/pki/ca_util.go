package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) getGenerationParams(ctx context.Context,
	data *framework.FieldData, mountPoint string,
) (exported bool, format string, role *roleEntry, errorResp *logical.Response) {
	exportedStr := data.Get("exported").(string)
	switch exportedStr {
	case "exported":
		exported = true
	case "internal":
	case "kms":
	default:
		errorResp = logical.ErrorResponse(
			`the "exported" path parameter must be "internal", "exported" or "kms"`)
		return
	}

	format = getFormat(data)
	if format == "" {
		errorResp = logical.ErrorResponse(
			`the "format" path parameter must be "pem", "der", or "pem_bundle"`)
		return
	}

	keyType := data.Get("key_type").(string)
	keyBits := data.Get("key_bits").(int)
	if exportedStr == "kms" {
		_, okKeyType := data.Raw["key_type"]
		_, okKeyBits := data.Raw["key_bits"]

		if okKeyType || okKeyBits {
			errorResp = logical.ErrorResponse(
				`invalid parameter for the kms path parameter, key_type nor key_bits arguments can be set in this mode`)
			return
		}

		keyId, err := getManagedKeyId(data)
		if err != nil {
			errorResp = logical.ErrorResponse("unable to determine managed key id")
			return
		}
		// Determine key type and key bits from the managed public key
		err = withManagedPKIKey(ctx, b, keyId, mountPoint, func(ctx context.Context, key logical.ManagedSigningKey) error {
			pubKey, err := key.GetPublicKey(ctx)
			if err != nil {
				return err
			}
			switch pubKey.(type) {
			case *rsa.PublicKey:
				keyType = "rsa"
				keyBits = pubKey.(*rsa.PublicKey).Size() * 8
			case *ecdsa.PublicKey:
				keyType = "ec"
			case *ed25519.PublicKey:
				keyType = "ed25519"
			default:
				return fmt.Errorf("unsupported public key: %#v", pubKey)
			}
			return nil
		})
		if err != nil {
			errorResp = logical.ErrorResponse("failed to lookup public key from managed key: %s", err.Error())
			return
		}
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

	var err error
	if role.KeyBits, role.SignatureBits, err = certutil.ValidateDefaultOrValueKeyTypeSignatureLength(role.KeyType, role.KeyBits, role.SignatureBits); err != nil {
		errorResp = logical.ErrorResponse(err.Error())
	}

	return
}

func generateCABundle(ctx context.Context, b *backend, input *inputBundle, data *certutil.CreationBundle, randomSource io.Reader) (*certutil.ParsedCertBundle, error) {
	if kmsRequested(input) {
		return generateManagedKeyCABundle(ctx, b, input, data, randomSource)
	}

	return certutil.CreateCertificateWithRandomSource(data, randomSource)
}

func generateCSRBundle(ctx context.Context, b *backend, input *inputBundle, data *certutil.CreationBundle, addBasicConstraints bool, randomSource io.Reader) (*certutil.ParsedCSRBundle, error) {
	if kmsRequested(input) {
		return generateManagedKeyCSRBundle(ctx, b, input, data, addBasicConstraints, randomSource)
	}

	return certutil.CreateCSRWithRandomSource(data, addBasicConstraints, randomSource)
}

func parseCABundle(ctx context.Context, b *backend, req *logical.Request, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	if bundle.PrivateKeyType == certutil.ManagedPrivateKey {
		return parseManagedKeyCABundle(ctx, b, req, bundle)
	}
	return bundle.ToParsedCertBundle()
}
