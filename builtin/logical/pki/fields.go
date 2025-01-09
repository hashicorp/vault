// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
)

const (
	issuerRefParam = "issuer_ref"
	keyNameParam   = "key_name"
	keyRefParam    = "key_ref"
	keyIdParam     = "key_id"
	keyTypeParam   = "key_type"
	keyBitsParam   = "key_bits"
	skidParam      = "subject_key_id"
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

	fields["remove_roots_from_chain"] = &framework.FieldSchema{
		Type:    framework.TypeBool,
		Default: false,
		Description: `Whether or not to remove self-signed CA certificates in the output
of the ca_chain field.`,
	}

	fields["user_ids"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `The requested user_ids value to place in the subject,
if any, in a comma-delimited list. Restricted by allowed_user_ids.
Any values are added with OID 0.9.2342.19200300.100.1.1.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "User ID(s)",
		},
	}
	fields["cert_metadata"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: `User supplied metadata to store associated with this certificate's serial number, base64 encoded`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Certificate Metadata",
		},
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

	fields["not_before_duration"] = &framework.FieldSchema{
		Type:        framework.TypeDurationSecond,
		Default:     30,
		Description: `The duration before now which the certificate needs to be backdated by.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Value: 30,
		},
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
0 (universal default); with rsa key_type: 2048 (default), 3072, 4096 or 8192;
with ec key_type: 224, 256 (default), 384, or 521; ignored with ed25519.`,
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

	fields["use_pss"] = &framework.FieldSchema{
		Type:    framework.TypeBool,
		Default: false,
		Description: `Whether or not to use PSS signatures when using a
RSA key-type issuer. Defaults to false.`,
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
	fields["excluded_dns_domains"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: `Domains for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Excluded DNS Domains",
		},
	}

	fields["permitted_ip_ranges"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `IP ranges for which this certificate is allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).
Ranges must be specified in the notation of IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Permitted IP ranges",
		},
	}
	fields["excluded_ip_ranges"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `IP ranges for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).
Ranges must be specified in the notation of IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Excluded IP ranges",
		},
	}

	fields["permitted_email_addresses"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: `Email addresses for which this certificate is allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Permitted email addresses",
		},
	}
	fields["excluded_email_addresses"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: `Email addresses for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Excluded email addresses",
		},
	}

	fields["permitted_uri_domains"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: `URI domains for which this certificate is allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Permitted URI domains",
		},
	}
	fields["excluded_uri_domains"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: `URI domains for which this certificate is not allowed to sign or issue child certificates (see https://tools.ietf.org/html/rfc5280#section-4.2.1.10).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Name: "Excluded URI domains",
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

func addTidyFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["tidy_cert_store"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to enable tidying up
the certificate store`,
	}

	fields["tidy_revocation_list"] = &framework.FieldSchema{
		Type:        framework.TypeBool,
		Description: `Deprecated; synonym for 'tidy_revoked_certs`,
	}

	fields["tidy_revoked_certs"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to expire all revoked
and expired certificates, removing them both from the CRL and from storage. The
CRL will be rotated if this causes any values to be removed.`,
	}

	fields["tidy_revoked_cert_issuer_associations"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to validate issuer associations
on revocation entries. This helps increase the performance of CRL building
and OCSP responses.`,
	}

	fields["tidy_expired_issuers"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to automatically remove expired issuers
past the issuer_safety_buffer. No keys will be removed as part of this
operation.`,
	}

	fields["tidy_move_legacy_ca_bundle"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to move the legacy ca_bundle from
/config/ca_bundle to /config/ca_bundle.bak. This prevents downgrades
to pre-Vault 1.11 versions (as older PKI engines do not know about
the new multi-issuer storage layout), but improves the performance
on seal wrapped PKI mounts. This will only occur if at least
issuer_safety_buffer time has occurred after the initial storage
migration.

This backup is saved in case of an issue in future migrations.
Operators may consider removing it via sys/raw if they desire.
The backup will be removed via a DELETE /root call, but note that
this removes ALL issuers within the mount (and is thus not desirable
in most operational scenarios).`,
	}

	fields["tidy_acme"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to enable tidying ACME accounts,
orders and authorizations.  ACME orders are tidied (deleted) 
safety_buffer after the certificate associated with them expires,
or after the order and relevant authorizations have expired if no 
certificate was produced.  Authorizations are tidied with the 
corresponding order.

When a valid ACME Account is at least acme_account_safety_buffer
old, and has no remaining orders associated with it, the account is
marked as revoked.  After another acme_account_safety_buffer has 
passed from the revocation or deactivation date, a revoked or 
deactivated ACME account is deleted.`,
		Default: false,
	}

	fields["safety_buffer"] = &framework.FieldSchema{
		Type: framework.TypeDurationSecond,
		Description: `The amount of extra time that must have passed
beyond certificate expiration before it is removed
from the backend storage and/or revocation list.
Defaults to 72 hours.`,
		Default: int(defaultTidyConfig.SafetyBuffer / time.Second), // TypeDurationSecond currently requires defaults to be int
	}

	fields["issuer_safety_buffer"] = &framework.FieldSchema{
		Type: framework.TypeDurationSecond,
		Description: `The amount of extra time that must have passed
beyond issuer's expiration before it is removed
from the backend storage.
Defaults to 8760 hours (1 year).`,
		Default: int(defaultTidyConfig.IssuerSafetyBuffer / time.Second), // TypeDurationSecond currently requires defaults to be int
	}

	fields["acme_account_safety_buffer"] = &framework.FieldSchema{
		Type: framework.TypeDurationSecond,
		Description: `The amount of time that must pass after creation
that an account with no orders is marked revoked, and the amount of time
after being marked revoked or deactivated.`,
		Default: int(defaultTidyConfig.AcmeAccountSafetyBuffer / time.Second), // TypeDurationSecond currently requires defaults to be int
	}

	fields["pause_duration"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The amount of time to wait between processing
certificates. This allows operators to change the execution profile
of tidy to take consume less resources by slowing down how long it
takes to run. Note that the entire list of certificates will be
stored in memory during the entire tidy operation, but resources to
read/process/update existing entries will be spread out over a
greater period of time. By default this is zero seconds.`,
		Default: "0s",
	}

	fields["tidy_revocation_queue"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to remove stale revocation queue entries
that haven't been confirmed by any active cluster. Only runs on the
active primary node`,
		Default: defaultTidyConfig.RevocationQueue,
	}

	fields["revocation_queue_safety_buffer"] = &framework.FieldSchema{
		Type: framework.TypeDurationSecond,
		Description: `The amount of time that must pass from the
cross-cluster revocation request being initiated to when it will be
slated for removal. Setting this too low may remove valid revocation
requests before the owning cluster has a chance to process them,
especially if the cluster is offline.`,
		Default: int(defaultTidyConfig.QueueSafetyBuffer / time.Second), // TypeDurationSecond currently requires defaults to be int
	}

	fields["tidy_cross_cluster_revoked_certs"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Set to true to enable tidying up
the cross-cluster revoked certificate store. Only runs on the active
primary node.`,
	}

	fields["tidy_cert_metadata"] = &framework.FieldSchema{
		Type:        framework.TypeBool,
		Description: `Set to true to enable tidying up certificate metadata`,
	}

	fields["tidy_cmpv2_nonce_store"] = &framework.FieldSchema{
		Type:        framework.TypeBool,
		Description: `Set to true to enable tidying up the CMPv2 nonce store`,
	}

	return fields
}

// generate the entire list of schema fields we need for CSR sign verbatim, this is also
// leveraged by ACME internally.
func getCsrSignVerbatimSchemaFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{}
	fields = addNonCACommonFields(fields)
	fields = addSignVerbatimRoleFields(fields)

	fields["csr"] = &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "",
		Description: `PEM-format CSR to be signed. Values will be
taken verbatim from the CSR, except for
basic constraints.`,
	}

	return fields
}

// addSignVerbatimRoleFields provides the fields and defaults to be used by anything that is building up the fields
// and their corresponding default values when generating/using a sign-verbatim type role such as buildSignVerbatimRole.
func addSignVerbatimRoleFields(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["key_usage"] = &framework.FieldSchema{
		Type:    framework.TypeCommaStringSlice,
		Default: issuing.DefaultRoleKeyUsages,
		Description: `A comma-separated string or list of key usages (not extended
key usages). Valid values can be found at
https://golang.org/pkg/crypto/x509/#KeyUsage
-- simply drop the "KeyUsage" part of the name.
To remove all key usages from being set, set
this value to an empty list.`,
	}

	fields["ext_key_usage"] = &framework.FieldSchema{
		Type:    framework.TypeCommaStringSlice,
		Default: issuing.DefaultRoleEstKeyUsages,
		Description: `A comma-separated string or list of extended key usages. Valid values can be found at
https://golang.org/pkg/crypto/x509/#ExtKeyUsage
-- simply drop the "ExtKeyUsage" part of the name.
To remove all key usages from being set, set
this value to an empty list.`,
	}

	fields["ext_key_usage_oids"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Default:     issuing.DefaultRoleEstKeyUsageOids,
		Description: `A comma-separated string or list of extended key usage oids.`,
	}

	fields["signature_bits"] = &framework.FieldSchema{
		Type:    framework.TypeInt,
		Default: issuing.DefaultRoleSignatureBits,
		Description: `The number of bits to use in the signature
algorithm; accepts 256 for SHA-2-256, 384 for SHA-2-384, and 512 for
SHA-2-512. Defaults to 0 to automatically detect based on key length
(SHA-2-256 for RSA keys, and matching the curve size for NIST P-Curves).`,
		DisplayAttrs: &framework.DisplayAttributes{
			Value: issuing.DefaultRoleSignatureBits,
		},
	}

	fields["use_pss"] = &framework.FieldSchema{
		Type:    framework.TypeBool,
		Default: issuing.DefaultRoleUsePss,
		Description: `Whether or not to use PSS signatures when using a
RSA key-type issuer. Defaults to false.`,
	}

	return fields
}

func addCACertKeyUsage(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["key_usage"] = &framework.FieldSchema{ // Same Name as Leaf-Cert Field, and CA CSR Field, but Description and Default Differ
		Type:    framework.TypeCommaStringSlice,
		Default: []string{"CertSign", "CRLSign"},
		Description: `This list of key usages (not extended key usages) will be 
added to the existing set of key usages, CRL,CertSign, on 
the generated certificate.  Valid values can be found at 
https://golang.org/pkg/crypto/x509/#KeyUsage -- simply drop 
the "KeyUsage" part of the name.  To use the issuer for 
CMPv2, DigitalSignature must be set.`,
	}

	return fields
}

func addCaCsrKeyUsage(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["key_usage"] = &framework.FieldSchema{ // Same Name as Leaf-Cert, CA-Cert Field, but Description and Default Differ
		Type:    framework.TypeCommaStringSlice,
		Default: []string{},
		Description: `Specifies key_usage to encode in the certificate signing
request.  This is a comma-separated string or list of key
usages (not extended key usages). Valid values can be found
at https://golang.org/pkg/crypto/x509/#KeyUsage -- simply 
drop the "KeyUsage" part of the name.  If not set, key 
usage will not appear on the CSR.`,
	}

	return fields
}
