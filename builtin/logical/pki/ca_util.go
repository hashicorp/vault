package pki

import (
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) getGenerationParams(
	data *framework.FieldData,
) (exported bool, format string, role *roleEntry, errorResp *logical.Response) {
	exportedStr := data.Get("exported").(string)
	switch exportedStr {
	case "exported":
		exported = true
	case "internal":
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

	role = &roleEntry{
		TTL:                  time.Duration(data.Get("ttl").(int)) * time.Second,
		KeyType:              data.Get("key_type").(string),
		KeyBits:              data.Get("key_bits").(int),
		AllowLocalhost:       true,
		AllowAnyName:         true,
		AllowIPSANs:          true,
		EnforceHostnames:     false,
		AllowedURISANs:       []string{"*"},
		AllowedSerialNumbers: []string{"*"},
		OU:                   data.Get("ou").([]string),
		Organization:         data.Get("organization").([]string),
		Country:              data.Get("country").([]string),
		Locality:             data.Get("locality").([]string),
		Province:             data.Get("province").([]string),
		StreetAddress:        data.Get("street_address").([]string),
		PostalCode:           data.Get("postal_code").([]string),
	}

	if role.KeyType == "rsa" && role.KeyBits < 2048 {
		errorResp = logical.ErrorResponse("RSA keys < 2048 bits are unsafe and not supported")
		return
	}

	errorResp = validateKeyTypeLength(role.KeyType, role.KeyBits)

	return
}
