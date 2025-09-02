// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	DefaultRoleKeyUsages       = []string{"DigitalSignature", "KeyAgreement", "KeyEncipherment"}
	DefaultRoleEstKeyUsages    = []string{}
	DefaultRoleEstKeyUsageOids = []string{}
)

const (
	DefaultRoleSignatureBits = 0
	DefaultRoleUsePss        = false
)

type RoleEntry struct {
	LeaseMax                      string        `json:"lease_max"`
	Lease                         string        `json:"lease"`
	DeprecatedMaxTTL              string        `json:"max_ttl"`
	DeprecatedTTL                 string        `json:"ttl"`
	TTL                           time.Duration `json:"ttl_duration"`
	MaxTTL                        time.Duration `json:"max_ttl_duration"`
	AllowLocalhost                bool          `json:"allow_localhost"`
	AllowedBaseDomain             string        `json:"allowed_base_domain"`
	AllowedDomainsOld             string        `json:"allowed_domains,omitempty"`
	AllowedDomains                []string      `json:"allowed_domains_list"`
	AllowedDomainsTemplate        bool          `json:"allowed_domains_template"`
	AllowBaseDomain               bool          `json:"allow_base_domain"`
	AllowBareDomains              bool          `json:"allow_bare_domains"`
	AllowTokenDisplayName         bool          `json:"allow_token_displayname"`
	AllowSubdomains               bool          `json:"allow_subdomains"`
	AllowGlobDomains              bool          `json:"allow_glob_domains"`
	AllowWildcardCertificates     *bool         `json:"allow_wildcard_certificates,omitempty"`
	AllowAnyName                  bool          `json:"allow_any_name"`
	EnforceHostnames              bool          `json:"enforce_hostnames"`
	AllowIPSANs                   bool          `json:"allow_ip_sans"`
	ServerFlag                    bool          `json:"server_flag"`
	ClientFlag                    bool          `json:"client_flag"`
	CodeSigningFlag               bool          `json:"code_signing_flag"`
	EmailProtectionFlag           bool          `json:"email_protection_flag"`
	UseCSRCommonName              bool          `json:"use_csr_common_name"`
	UseCSRSANs                    bool          `json:"use_csr_sans"`
	SerialNumberSource            string        `json:"serial_number_source"`
	KeyType                       string        `json:"key_type"`
	KeyBits                       int           `json:"key_bits"`
	UsePSS                        bool          `json:"use_pss"`
	SignatureBits                 int           `json:"signature_bits"`
	MaxPathLength                 *int          `json:",omitempty"`
	KeyUsageOld                   string        `json:"key_usage,omitempty"`
	KeyUsage                      []string      `json:"key_usage_list"`
	ExtKeyUsage                   []string      `json:"extended_key_usage_list"`
	OUOld                         string        `json:"ou,omitempty"`
	OU                            []string      `json:"ou_list"`
	OrganizationOld               string        `json:"organization,omitempty"`
	Organization                  []string      `json:"organization_list"`
	Country                       []string      `json:"country"`
	Locality                      []string      `json:"locality"`
	Province                      []string      `json:"province"`
	StreetAddress                 []string      `json:"street_address"`
	PostalCode                    []string      `json:"postal_code"`
	GenerateLease                 *bool         `json:"generate_lease,omitempty"`
	NoStore                       bool          `json:"no_store"`
	NoStoreMetadata               bool          `json:"no_store_metadata"`
	RequireCN                     bool          `json:"require_cn"`
	CNValidations                 []string      `json:"cn_validations"`
	AllowedOtherSANs              []string      `json:"allowed_other_sans"`
	AllowedSerialNumbers          []string      `json:"allowed_serial_numbers"`
	AllowedUserIDs                []string      `json:"allowed_user_ids"`
	AllowedURISANs                []string      `json:"allowed_uri_sans"`
	AllowedURISANsTemplate        bool          `json:"allowed_uri_sans_template"`
	PolicyIdentifiers             []string      `json:"policy_identifiers"`
	ExtKeyUsageOIDs               []string      `json:"ext_key_usage_oids"`
	BasicConstraintsValidForNonCA bool          `json:"basic_constraints_valid_for_non_ca"`
	NotBeforeDuration             time.Duration `json:"not_before_duration"`
	NotAfter                      string        `json:"not_after"`
	Issuer                        string        `json:"issuer"`
	// Name is only set when the role has been stored, on the fly roles have a blank name
	Name string `json:"-"`
	// WasModified indicates to callers if the returned entry is different than the persisted version
	WasModified bool `json:"-"`
}

func (r *RoleEntry) ToResponseData() map[string]interface{} {
	responseData := map[string]interface{}{
		"ttl":                                int64(r.TTL.Seconds()),
		"max_ttl":                            int64(r.MaxTTL.Seconds()),
		"allow_localhost":                    r.AllowLocalhost,
		"allowed_domains":                    r.AllowedDomains,
		"allowed_domains_template":           r.AllowedDomainsTemplate,
		"allow_bare_domains":                 r.AllowBareDomains,
		"allow_token_displayname":            r.AllowTokenDisplayName,
		"allow_subdomains":                   r.AllowSubdomains,
		"allow_glob_domains":                 r.AllowGlobDomains,
		"allow_wildcard_certificates":        r.AllowWildcardCertificates,
		"allow_any_name":                     r.AllowAnyName,
		"allowed_uri_sans_template":          r.AllowedURISANsTemplate,
		"enforce_hostnames":                  r.EnforceHostnames,
		"allow_ip_sans":                      r.AllowIPSANs,
		"server_flag":                        r.ServerFlag,
		"client_flag":                        r.ClientFlag,
		"code_signing_flag":                  r.CodeSigningFlag,
		"email_protection_flag":              r.EmailProtectionFlag,
		"use_csr_common_name":                r.UseCSRCommonName,
		"use_csr_sans":                       r.UseCSRSANs,
		"serial_number_source":               r.SerialNumberSource,
		"key_type":                           r.KeyType,
		"key_bits":                           r.KeyBits,
		"signature_bits":                     r.SignatureBits,
		"use_pss":                            r.UsePSS,
		"key_usage":                          r.KeyUsage,
		"ext_key_usage":                      r.ExtKeyUsage,
		"ext_key_usage_oids":                 r.ExtKeyUsageOIDs,
		"ou":                                 r.OU,
		"organization":                       r.Organization,
		"country":                            r.Country,
		"locality":                           r.Locality,
		"province":                           r.Province,
		"street_address":                     r.StreetAddress,
		"postal_code":                        r.PostalCode,
		"no_store":                           r.NoStore,
		"allowed_other_sans":                 r.AllowedOtherSANs,
		"allowed_serial_numbers":             r.AllowedSerialNumbers,
		"allowed_user_ids":                   r.AllowedUserIDs,
		"allowed_uri_sans":                   r.AllowedURISANs,
		"require_cn":                         r.RequireCN,
		"cn_validations":                     r.CNValidations,
		"policy_identifiers":                 r.PolicyIdentifiers,
		"basic_constraints_valid_for_non_ca": r.BasicConstraintsValidForNonCA,
		"not_before_duration":                int64(r.NotBeforeDuration.Seconds()),
		"not_after":                          r.NotAfter,
		"issuer_ref":                         r.Issuer,
	}
	if r.MaxPathLength != nil {
		responseData["max_path_length"] = r.MaxPathLength
	}
	if r.GenerateLease != nil {
		responseData["generate_lease"] = r.GenerateLease
	}
	AddNoStoreMetadata(responseData, r)
	return responseData
}

var ErrRoleNotFound = errors.New("role not found")

// GetRole will load a role from storage based on the provided name and
// update its contents to the latest version if out of date. The WasUpdated field
// will be set to true if modifications were made indicating the caller should if
// possible write them back to disk. If the role is not found an ErrRoleNotFound
// will be returned as an error.
func GetRole(ctx context.Context, s logical.Storage, n string) (*RoleEntry, error) {
	entry, err := s.Get(ctx, "role/"+n)
	if err != nil {
		return nil, fmt.Errorf("failed to load role %s: %w", n, err)
	}
	if entry == nil {
		return nil, fmt.Errorf("%w: with name %s", ErrRoleNotFound, n)
	}

	var result RoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("failed decoding role %s: %w", n, err)
	}

	// Migrate existing saved entries and save back if changed
	modified := false
	if len(result.DeprecatedTTL) == 0 && len(result.Lease) != 0 {
		result.DeprecatedTTL = result.Lease
		result.Lease = ""
		modified = true
	}
	if result.TTL == 0 && len(result.DeprecatedTTL) != 0 {
		parsed, err := parseutil.ParseDurationSecond(result.DeprecatedTTL)
		if err != nil {
			return nil, err
		}
		result.TTL = parsed
		result.DeprecatedTTL = ""
		modified = true
	}
	if len(result.DeprecatedMaxTTL) == 0 && len(result.LeaseMax) != 0 {
		result.DeprecatedMaxTTL = result.LeaseMax
		result.LeaseMax = ""
		modified = true
	}
	if result.MaxTTL == 0 && len(result.DeprecatedMaxTTL) != 0 {
		parsed, err := parseutil.ParseDurationSecond(result.DeprecatedMaxTTL)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_ttl field in %s: %w", n, err)
		}
		result.MaxTTL = parsed
		result.DeprecatedMaxTTL = ""
		modified = true
	}
	if result.AllowBaseDomain {
		result.AllowBaseDomain = false
		result.AllowBareDomains = true
		modified = true
	}
	if result.AllowedDomainsOld != "" {
		result.AllowedDomains = strings.Split(result.AllowedDomainsOld, ",")
		result.AllowedDomainsOld = ""
		modified = true
	}
	if result.AllowedBaseDomain != "" {
		found := false
		for _, v := range result.AllowedDomains {
			if v == result.AllowedBaseDomain {
				found = true
				break
			}
		}
		if !found {
			result.AllowedDomains = append(result.AllowedDomains, result.AllowedBaseDomain)
		}
		result.AllowedBaseDomain = ""
		modified = true
	}
	if result.AllowWildcardCertificates == nil {
		// While not the most secure default, when AllowWildcardCertificates isn't
		// explicitly specified in the stored Role, we automatically upgrade it to
		// true to preserve compatibility with previous versions of Vault. Once this
		// field is set, this logic will not be triggered any more.
		result.AllowWildcardCertificates = new(bool)
		*result.AllowWildcardCertificates = true
		modified = true
	}

	// Upgrade generate_lease in role
	if result.GenerateLease == nil {
		// All the new roles will have GenerateLease always set to a Value. A
		// nil Value indicates that this role needs an upgrade. Set it to
		// `true` to not alter its current behavior.
		result.GenerateLease = new(bool)
		*result.GenerateLease = true
		modified = true
	}

	// Upgrade key usages
	if result.KeyUsageOld != "" {
		result.KeyUsage = strings.Split(result.KeyUsageOld, ",")
		result.KeyUsageOld = ""
		modified = true
	}

	// Upgrade OU
	if result.OUOld != "" {
		result.OU = strings.Split(result.OUOld, ",")
		result.OUOld = ""
		modified = true
	}

	// Upgrade Organization
	if result.OrganizationOld != "" {
		result.Organization = strings.Split(result.OrganizationOld, ",")
		result.OrganizationOld = ""
		modified = true
	}

	// Set the issuer field to default if not set. We want to do this
	// unconditionally as we should probably never have an empty issuer
	// on a stored roles.
	if len(result.Issuer) == 0 {
		result.Issuer = DefaultRef
		modified = true
	}

	// Update CN Validations to be the present default, "email,hostname"
	if len(result.CNValidations) == 0 {
		result.CNValidations = []string{"email", "hostname"}
		modified = true
	}

	result.Name = n
	result.WasModified = modified

	return &result, nil
}

type RoleModifier func(r *RoleEntry)

func WithKeyUsage(keyUsages []string) RoleModifier {
	return func(r *RoleEntry) {
		r.KeyUsage = keyUsages
	}
}

func WithExtKeyUsage(extKeyUsages []string) RoleModifier {
	return func(r *RoleEntry) {
		r.ExtKeyUsage = extKeyUsages
	}
}

func WithExtKeyUsageOIDs(extKeyUsageOids []string) RoleModifier {
	return func(r *RoleEntry) {
		r.ExtKeyUsageOIDs = extKeyUsageOids
	}
}

func WithSignatureBits(signatureBits int) RoleModifier {
	return func(r *RoleEntry) {
		r.SignatureBits = signatureBits
	}
}

func WithUsePSS(usePss bool) RoleModifier {
	return func(r *RoleEntry) {
		r.UsePSS = usePss
	}
}

func WithTTL(ttl time.Duration) RoleModifier {
	return func(r *RoleEntry) {
		r.TTL = ttl
	}
}

func WithMaxTTL(ttl time.Duration) RoleModifier {
	return func(r *RoleEntry) {
		r.MaxTTL = ttl
	}
}

func WithGenerateLease(genLease bool) RoleModifier {
	return func(r *RoleEntry) {
		*r.GenerateLease = genLease
	}
}

func WithNotBeforeDuration(ttl time.Duration) RoleModifier {
	return func(r *RoleEntry) {
		r.NotBeforeDuration = ttl
	}
}

func WithNoStore(noStore bool) RoleModifier {
	return func(r *RoleEntry) {
		r.NoStore = noStore
	}
}

func WithIssuer(issuer string) RoleModifier {
	return func(r *RoleEntry) {
		if issuer == "" {
			issuer = DefaultRef
		}
		r.Issuer = issuer
	}
}

// SignVerbatimRole create a sign-verbatim role with no overrides. This will store
// the signed certificate, allowing any key type and Value from a role restriction.
func SignVerbatimRole() *RoleEntry {
	return SignVerbatimRoleWithOpts()
}

// SignVerbatimRoleWithOpts create a sign-verbatim role with the normal defaults,
// but allowing any field to be tweaked based on the consumers needs.
func SignVerbatimRoleWithOpts(opts ...RoleModifier) *RoleEntry {
	entry := &RoleEntry{
		AllowLocalhost:            true,
		AllowAnyName:              true,
		AllowIPSANs:               true,
		AllowWildcardCertificates: new(bool),
		EnforceHostnames:          false,
		KeyType:                   "any",
		UseCSRCommonName:          true,
		UseCSRSANs:                true,
		SerialNumberSource:        "json-csr",
		AllowedOtherSANs:          []string{"*"},
		AllowedSerialNumbers:      []string{"*"},
		AllowedURISANs:            []string{"*"},
		AllowedUserIDs:            []string{"*"},
		CNValidations:             []string{"disabled"},
		GenerateLease:             new(bool),
		KeyUsage:                  DefaultRoleKeyUsages,
		ExtKeyUsage:               DefaultRoleEstKeyUsages,
		ExtKeyUsageOIDs:           DefaultRoleEstKeyUsageOids,
		SignatureBits:             DefaultRoleSignatureBits,
		UsePSS:                    DefaultRoleUsePss,
	}
	*entry.AllowWildcardCertificates = true
	*entry.GenerateLease = false

	if opts != nil {
		for _, opt := range opts {
			if opt != nil {
				opt(entry)
			}
		}
	}

	return entry
}

func ParseExtKeyUsagesFromRole(role *RoleEntry) certutil.CertExtKeyUsage {
	var parsedKeyUsages certutil.CertExtKeyUsage

	if role.ServerFlag {
		parsedKeyUsages |= certutil.ServerAuthExtKeyUsage
	}

	if role.ClientFlag {
		parsedKeyUsages |= certutil.ClientAuthExtKeyUsage
	}

	if role.CodeSigningFlag {
		parsedKeyUsages |= certutil.CodeSigningExtKeyUsage
	}

	if role.EmailProtectionFlag {
		parsedKeyUsages |= certutil.EmailProtectionExtKeyUsage
	}

	for _, k := range role.ExtKeyUsage {
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "any":
			parsedKeyUsages |= certutil.AnyExtKeyUsage
		case "serverauth":
			parsedKeyUsages |= certutil.ServerAuthExtKeyUsage
		case "clientauth":
			parsedKeyUsages |= certutil.ClientAuthExtKeyUsage
		case "codesigning":
			parsedKeyUsages |= certutil.CodeSigningExtKeyUsage
		case "emailprotection":
			parsedKeyUsages |= certutil.EmailProtectionExtKeyUsage
		case "ipsecendsystem":
			parsedKeyUsages |= certutil.IpsecEndSystemExtKeyUsage
		case "ipsectunnel":
			parsedKeyUsages |= certutil.IpsecTunnelExtKeyUsage
		case "ipsecuser":
			parsedKeyUsages |= certutil.IpsecUserExtKeyUsage
		case "timestamping":
			parsedKeyUsages |= certutil.TimeStampingExtKeyUsage
		case "ocspsigning":
			parsedKeyUsages |= certutil.OcspSigningExtKeyUsage
		case "microsoftservergatedcrypto":
			parsedKeyUsages |= certutil.MicrosoftServerGatedCryptoExtKeyUsage
		case "netscapeservergatedcrypto":
			parsedKeyUsages |= certutil.NetscapeServerGatedCryptoExtKeyUsage
		}
	}

	return parsedKeyUsages
}
