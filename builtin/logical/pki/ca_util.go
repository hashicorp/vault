package pki

import (
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) getGenerationParams(
	data *framework.FieldData,
) (exported bool, format string, role *roleEntry, errorResp *logical.Response) {
	exportedStr := data.Get("exported").(string)
	switch exportedStr {
	case "exported":
		exported = true
	case "internal":
	case "kms":
	default:
		errorResp = logical.ErrorResponse(
			`the "exported" path parameter must be "internal" or "exported"`)
		return
	}

	format = getFormat(data)
	if format == "" {
		errorResp = logical.ErrorResponse(
			`the "format" path parameter must be "pem", "der", "der_pkcs", or "pem_bundle"`)
		return
	}

	if exportedStr == "kms" {
		_, okKeyType := data.Raw["key_type"]
		_, okKeyBits := data.Raw["key_bits"]

		if okKeyType || okKeyBits {
			errorResp = logical.ErrorResponse(
				`invalid parameter for the kms path parameter, key_type nor key_bits arguments can be set in this mode`)
			return
		}
	}

	role = &roleEntry{
		TTL:                  time.Duration(data.Get("ttl").(int)) * time.Second,
		KeyType:              data.Get("key_type").(string),
		KeyBits:              data.Get("key_bits").(int),
		SignatureBits:        data.Get("signature_bits").(int),
		AllowLocalhost:       true,
		AllowAnyName:         true,
		AllowIPSANs:          true,
		EnforceHostnames:     false,
		AllowedURISANs:       []string{"*"},
		AllowedOtherSANs:     []string{"*"},
		AllowedSerialNumbers: []string{"*"},
		OU:                   data.Get("ou").([]string),
		Organization:         data.Get("organization").([]string),
		Country:              data.Get("country").([]string),
		Locality:             data.Get("locality").([]string),
		Province:             data.Get("province").([]string),
		StreetAddress:        data.Get("street_address").([]string),
		PostalCode:           data.Get("postal_code").([]string),
	}

	var err error
	if role.KeyBits, role.SignatureBits, err = certutil.ValidateDefaultOrValueKeyTypeSignatureLength(role.KeyType, role.KeyBits, role.SignatureBits); err != nil {
		errorResp = logical.ErrorResponse(err.Error())
	}

	return
}
