package pki

import "github.com/hashicorp/vault/sdk/framework"

const (
	issuerRefParam = "issuer_ref"
	keyNameParam   = "key_name"
	keyRefParam    = "key_ref"
	keyIdParam     = "key_id"
	keyTypeParam   = "key_type"
	keyBitsParam   = "key_bits"
)

// addIssueAndSignCommonFields adds fields common to both CA and non-CA issuing
// and signing
func addIssueAndSignCommonFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["exclude_cn_from_sans"] = &framework.FieldSchema{
		Type:    framework.TypeBool,
		Default: false,
		Description: `If true, the Common Name will not be
included in DNS or Email Subject Alternate Names.
Defaults to false (CN is included).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Exclude Common Name from Subject Alternative Names (SANs)",
		},
	}

	fields["format"] = &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "pem",
		Description: `Format for returned data. Can be "pem", "der",
or "pem_bundle". If "pem_bundle", any private
key and issuing cert will be appended to the
certificate pem. If "der", the value will be
base64 encoded. Defaults to "pem".`,
		AllowedValues: []interface{}{"pem", "der", "pem_bundle"},
		DisplayAttrs: &framework.DisplayAttributes{
			Value: "pem",
		},
	}

	fields["private_key_format"] = &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "der",
		Description: `Format for the returned private key.
Generally the default will be controlled by the "format"
parameter as either base64-encoded DER or PEM-encoded DER.
However, this can be set to "pkcs8" to have the returned
private key contain base64-encoded pkcs8 or PEM-encoded
pkcs8 instead. Defaults to "der".`,
		AllowedValues: []interface{}{"", "der", "pem", "pkcs8"},
		DisplayAttrs: &framework.DisplayAttributes{
			Value: "der",
		},
	}

	fields["ip_sans"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `The requested IP SANs, if any, in a
comma-delimited list`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "IP Subject Alternative Names (SANs)",
		},
	}

	fields["uri_sans"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `The requested URI SANs, if any, in a
comma-delimited list.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "URI Subject Alternative Names (SANs)",
		},
	}

	fields["other_sans"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `Requested other SANs, in an array with the format
<oid>;UTF8:<utf8 string value> for each entry.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Other SANs",
		},
	}

	return fields
}

// addNonCACommonFields adds fields with help text specific to non-CA
// certificate issuing and signing
func addNonCACommonFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields = addIssueAndSignCommonFields(fields)

	fields["role"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The desired role with configuration for this
request`,
	}

	fields["common_name"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested common name; if you want more than
one, specify the alternative names in the
alt_names map. If email protection is enabled
in the role, this may be an email address.`,
	}

	fields["alt_names"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested Subject Alternative Names, if any,
in a comma-delimited list. If email protection
is enabled for the role, this may contain
email addresses.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "DNS/Email Subject Alternative Names (SANs)",
		},
	}

	fields["serial_number"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The Subject's requested serial number, if any.
See RFC 4519 Section 2.31 'serialNumber' for a description of this field.
If you want more than one, specify alternative names in the alt_names
map using OID 2.5.4.5. This has no impact on the final certificate's
Serial Number field.`,
	}

	fields["ttl"] = &framework.FieldSchema{
		Type: framework.TypeDurationSecond,
		Description: `The requested Time To Live for the certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be larger than the role max TTL.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "TTL",
		},
	}

	fields["not_after"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Set the not after field of the certificate with specified date value.
The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ`,
	}

	fields = addIssuerRefField(fields)

	return fields
}

// addCACommonFields adds fields with help text specific to CA
// certificate issuing and signing
func addCACommonFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields = addIssueAndSignCommonFields(fields)

	fields["alt_names"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested Subject Alternative Names, if any,
in a comma-delimited list. May contain both
DNS names and email addresses.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "DNS/Email Subject Alternative Names (SANs)",
		},
	}

	fields["common_name"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested common name; if you want more than
one, specify the alternative names in the alt_names
map. If not specified when signing, the common
name will be taken from the CSR; other names
must still be specified in alt_names or ip_sans.`,
	}

	fields["ttl"] = &framework.FieldSchema{
		Type: framework.TypeDurationSecond,
		Description: `The requested Time To Live for the certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be larger than the mount max TTL. Note:
this only has an effect when generating
a CA cert or signing a CA cert, not when
generating a CSR for an intermediate CA.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "TTL",
		},
	}

	fields["ou"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `If set, OU (OrganizationalUnit) will be set to
this value.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "OU (Organizational Unit)",
		},
	}

	fields["organization"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `If set, O (Organization) will be set to
this value.`,
	}

	fields["country"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `If set, Country will be set to
this value.`,
	}

	fields["locality"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `If set, Locality will be set to
this value.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Locality/City",
		},
	}

	fields["province"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `If set, Province will be set to
this value.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Province/State",
		},
	}

	fields["street_address"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `If set, Street Address will be set to
this value.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Street Address",
		},
	}

	fields["postal_code"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `If set, Postal Code will be set to
this value.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Postal Code",
		},
	}

	fields["serial_number"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The Subject's requested serial number, if any.
See RFC 4519 Section 2.31 'serialNumber' for a description of this field.
If you want more than one, specify alternative names in the alt_names
map using OID 2.5.4.5. This has no impact on the final certificate's
Serial Number field.`,
	}
	fields["not_after"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Set the not after field of the certificate with specified date value.
The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ`,
	}

	return fields
}

// addCAKeyGenerationFields adds fields with help text specific to CA key
// generation and exporting
func addCAKeyGenerationFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["exported"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Must be "internal", "exported" or "kms". If set to
"exported", the generated private key will be
returned. This is your *only* chance to retrieve
the private key!`,
		AllowedValues: []interface{}{"internal", "external", "kms"},
	}

	fields["managed_key_name"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The name of the managed key to use when the exported
type is kms. When kms type is the key type, this field or managed_key_id
is required. Ignored for other types.`,
	}

	fields["managed_key_id"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The name of the managed key to use when the exported
type is kms. When kms type is the key type, this field or managed_key_name
is required. Ignored for other types.`,
	}

	fields["key_bits"] = &framework.FieldSchema{
		Type:    framework.TypeInt,
		Default: 0,
		Description: `The number of bits to use. Allowed values are
0 (universal default); with rsa key_type: 2048 (default), 3072, or
4096; with ec key_type: 224, 256 (default), 384, or 521; ignored with
ed25519.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Value: 0,
		},
	}

	fields["signature_bits"] = &framework.FieldSchema{
		Type:    framework.TypeInt,
		Default: 0,
		Description: `The number of bits to use in the signature
algorithm; accepts 256 for SHA-2-256, 384 for SHA-2-384, and 512 for
SHA-2-512. Defaults to 0 to automatically detect based on key length
(SHA-2-256 for RSA keys, and matching the curve size for NIST P-Curves).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Value: 0,
		},
	}

	fields["key_type"] = &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "rsa",
		Description: `The type of key to use; defaults to RSA. "rsa"
"ec" and "ed25519" are the only valid values.`,
		AllowedValues: []interface{}{"rsa", "ec", "ed25519"},
		DisplayAttrs: &framework.DisplayAttributes{
			Value: "rsa",
		},
	}

	fields = addKeyRefNameFields(fields)

	return fields
}

// addCAIssueFields adds fields common to CA issuing, e.g. when returning
// an actual certificate
func addCAIssueFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["max_path_length"] = &framework.FieldSchema{
		Type:        framework.TypeInt,
		Default:     -1,
		Description: "The maximum allowable path length",
	}

	fields["permitted_dns_domains"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: `Domains for which this certificate is allowed to sign or issue child certificates. If set, all DNS names (subject and alt) on child certs must be exact matches or subsets of the given domains (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Permitted DNS Domains",
		},
	}

	fields = addIssuerNameField(fields)

	return fields
}

func addIssuerRefNameFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields = addIssuerNameField(fields)
	fields = addIssuerRefField(fields)
	return fields
}

func addIssuerNameField(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["issuer_name"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Provide a name to the generated or existing issuer, the name
must be unique across all issuers and not be the reserved value 'default'`,
	}
	return fields
}

func addIssuerRefField(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields[issuerRefParam] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Reference to a existing issuer; either "default"
for the configured default issuer, an identifier or the name assigned
to the issuer.`,
		Default: defaultRef,
	}
	return fields
}

func addKeyRefNameFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields = addKeyNameField(fields)
	fields = addKeyRefField(fields)
	return fields
}

func addKeyNameField(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields[keyNameParam] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Provide a name to the generated or existing key, the name
must be unique across all keys and not be the reserved value 'default'`,
	}

	return fields
}

func addKeyRefField(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields[keyRefParam] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Reference to a existing key; either "default"
for the configured default key, an identifier or the name assigned
to the key.`,
		Default: defaultRef,
	}
	return fields
}
