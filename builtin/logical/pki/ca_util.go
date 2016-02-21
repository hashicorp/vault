package pki

import (
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
			`The "exported" path parameter must be "internal" or "exported"`)
		return
	}

	format = getFormat(data)
	if format == "" {
		errorResp = logical.ErrorResponse(
			`The "format" path parameter must be "pem", "der", or "pem_bundle"`)
		return
	}

	role = &roleEntry{
		TTL:              data.Get("ttl").(string),
		KeyType:          data.Get("key_type").(string),
		KeyBits:          data.Get("key_bits").(int),
		AllowLocalhost:   true,
		AllowAnyName:     true,
		AllowIPSANs:      true,
		EnforceHostnames: false,
	}

	if role.KeyType == "rsa" && role.KeyBits < 2048 {
		errorResp = logical.ErrorResponse("RSA keys < 2048 bits are unsafe and not supported")
		return
	}

	errorResp = validateKeyTypeLength(role.KeyType, role.KeyBits)

	return
}
