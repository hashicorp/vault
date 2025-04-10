// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
	"golang.org/x/net/idna"
)

const (
	PathCerts        = "certs/"
	PathCertMetadata = "cert-metadata/"
	PathCrls         = "crls/"
)

var (
	// labelRegex is a single label from a valid domain name and was extracted
	// from hostnameRegex below for use in leftWildLabelRegex, without any
	// label separators (`.`).
	labelRegex = `([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])`

	// A note on hostnameRegex: although we set the StrictDomainName option
	// when doing the idna conversion, this appears to only affect output, not
	// input, so it will allow e.g. host^123.example.com straight through. So
	// we still need to use this to check the output.
	hostnameRegex = regexp.MustCompile(`^(\*\.)?(` + labelRegex + `\.)*` + labelRegex + `\.?$`)

	// Left Wildcard Label Regex is equivalent to a single domain label
	// component from hostnameRegex above, but with additional wildcard
	// characters added. There are four possibilities here:
	//
	//  1. Entire label is a wildcard,
	//  2. Wildcard exists at the start,
	//  3. Wildcard exists at the end,
	//  4. Wildcard exists in the middle.
	allWildRegex       = `\*`
	startWildRegex     = `\*` + labelRegex
	endWildRegex       = labelRegex + `\*`
	middleWildRegex    = labelRegex + `\*` + labelRegex
	leftWildLabelRegex = regexp.MustCompile(`^(` + allWildRegex + `|` + startWildRegex + `|` + endWildRegex + `|` + middleWildRegex + `)$`)
)

type EntityInfo struct {
	DisplayName string
	EntityID    string
}

type CertificateCounter interface {
	IsInitialized() bool
	IncrementTotalCertificatesCount(certsCounted bool, newSerial string)
	IncrementTotalRevokedCertificatesCount(certsCounted bool, newSerial string)
}

func NewEntityInfoFromReq(req *logical.Request) EntityInfo {
	if req == nil {
		return EntityInfo{}
	}
	return EntityInfo{
		DisplayName: req.DisplayName,
		EntityID:    req.EntityID,
	}
}

type CreationBundleInput interface {
	CertNotAfterInput
	GetCommonName() string
	GetSerialNumber() string
	GetExcludeCnFromSans() bool
	GetOptionalAltNames() (interface{}, bool)
	GetOtherSans() []string
	GetIpSans() []string
	GetURISans() []string
	GetOptionalSkid() (interface{}, bool)
	IsUserIdInSchema() (interface{}, bool)
	GetUserIds() []string
	IgnoreCSRSignature() bool
}

// GenerateCreationBundle is a shared function that reads parameters supplied
// from the various endpoints and generates a CreationParameters with the
// parameters that can be used to issue or sign
func GenerateCreationBundle(b logical.SystemView, role *RoleEntry, entityInfo EntityInfo, cb CreationBundleInput, caSign *certutil.CAInfoBundle, csr *x509.CertificateRequest) (*certutil.CreationBundle, []string, error) {
	// Read in names -- CN, DNS and email addresses
	var cn string
	var ridSerialNumber string
	var warnings []string
	dnsNames := []string{}
	emailAddresses := []string{}
	{
		if csr != nil && role.UseCSRCommonName {
			cn = csr.Subject.CommonName
		}
		if cn == "" {
			cn = cb.GetCommonName()
			if cn == "" && role.RequireCN {
				return nil, nil, errutil.UserError{Err: `the common_name field is required, or must be provided in a CSR with "use_csr_common_name" set to true, unless "require_cn" is set to false`}
			}
		}

		ridSerialNumber = cb.GetSerialNumber()

		// only take serial number from CSR if one was not supplied via API
		switch role.SerialNumberSource {
		case "", "json-csr":
			if ridSerialNumber == "" && csr != nil {
				ridSerialNumber = csr.Subject.SerialNumber
			}
		case "json":
			// use the value from cb set above
		default:
			return nil, nil, errutil.UserError{Err: "invalid value for serial_number_source"}
		}

		if csr != nil && role.UseCSRSANs {
			dnsNames = csr.DNSNames
			emailAddresses = csr.EmailAddresses
		}

		if cn != "" && !cb.GetExcludeCnFromSans() {
			if strings.Contains(cn, "@") {
				// Note: emails are not disallowed if the role's email protection
				// flag is false, because they may well be included for
				// informational purposes; it is up to the verifying party to
				// ensure that email addresses in a subject alternate name can be
				// used for the purpose for which they are presented
				emailAddresses = append(emailAddresses, cn)
			} else {
				// Only add to dnsNames if it's actually a DNS name but convert
				// idn first
				p := idna.New(
					idna.StrictDomainName(true),
					idna.VerifyDNSLength(true),
				)
				converted, err := p.ToASCII(cn)
				if err != nil {
					return nil, nil, errutil.UserError{Err: err.Error()}
				}
				if hostnameRegex.MatchString(converted) {
					dnsNames = append(dnsNames, converted)
				}
			}
		}

		if csr == nil || !role.UseCSRSANs {
			cnAltRaw, ok := cb.GetOptionalAltNames()
			if ok {
				cnAlt := strutil.ParseDedupAndSortStrings(cnAltRaw.(string), ",")
				for _, v := range cnAlt {
					if strings.Contains(v, "@") {
						emailAddresses = append(emailAddresses, v)
					} else {
						// Only add to dnsNames if it's actually a DNS name but
						// convert idn first
						p := idna.New(
							idna.StrictDomainName(true),
							idna.VerifyDNSLength(true),
						)
						converted, err := p.ToASCII(v)
						if err != nil {
							return nil, nil, errutil.UserError{Err: err.Error()}
						}
						if hostnameRegex.MatchString(converted) {
							dnsNames = append(dnsNames, converted)
						}
					}
				}
			}
		}

		// Check the CN. This ensures that the CN is checked even if it's
		// excluded from SANs.
		if cn != "" {
			badName := ValidateCommonName(b, role, entityInfo, cn)
			if len(badName) != 0 {
				return nil, nil, errutil.UserError{Err: fmt.Sprintf(
					"common name %s not allowed by this role", badName)}
			}
		}

		if ridSerialNumber != "" {
			badName := ValidateSerialNumber(role, ridSerialNumber)
			if len(badName) != 0 {
				return nil, nil, errutil.UserError{Err: fmt.Sprintf(
					"serial_number %s not allowed by this role", badName)}
			}
		}

		// Check for bad email and/or DNS names
		badName := ValidateNames(b, role, entityInfo, dnsNames)
		if len(badName) != 0 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"subject alternate name %s not allowed by this role", badName)}
		}

		badName = ValidateNames(b, role, entityInfo, emailAddresses)
		if len(badName) != 0 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"email address %s not allowed by this role", badName)}
		}
	}

	// otherSANsInput has the same format as the other_sans HTTP param in the
	// Vault PKI API: it is a list of strings of the form <oid>;<type>:<value>
	// where <type> must be UTF8/UTF-8.
	var otherSANsInput []string
	// otherSANs is the output of parseOtherSANs(otherSANsInput): its keys are
	// the <oid> value, its values are of the form [<type>, <value>]
	var otherSANs map[string][]string
	if sans := cb.GetOtherSans(); len(sans) > 0 {
		otherSANsInput = sans
	}
	if role.UseCSRSANs && csr != nil && len(csr.Extensions) > 0 {
		others, err := certutil.GetOtherSANsFromX509Extensions(csr.Extensions)
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Errorf("could not parse requested other SAN: %w", err).Error()}
		}
		for _, other := range others {
			otherSANsInput = append(otherSANsInput, other.String())
		}
	}
	if len(otherSANsInput) > 0 {
		requested, err := ParseOtherSANs(otherSANsInput)
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Errorf("could not parse requested other SAN: %w", err).Error()}
		}
		badOID, badName, err := ValidateOtherSANs(role, requested)
		switch {
		case err != nil:
			return nil, nil, errutil.UserError{Err: err.Error()}
		case len(badName) > 0:
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"other SAN %s not allowed for OID %s by this role", badName, badOID)}
		case len(badOID) > 0:
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"other SAN OID %s not allowed by this role", badOID)}
		default:
			otherSANs = requested
		}
	}

	// Get and verify any IP SANs
	ipAddresses := []net.IP{}
	{
		if csr != nil && role.UseCSRSANs {
			if len(csr.IPAddresses) > 0 {
				if !role.AllowIPSANs {
					return nil, nil, errutil.UserError{Err: "IP Subject Alternative Names are not allowed in this role, but was provided some via CSR"}
				}
				ipAddresses = csr.IPAddresses
			}
		} else {
			ipAlt := cb.GetIpSans()
			if len(ipAlt) > 0 {
				if !role.AllowIPSANs {
					return nil, nil, errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)}
				}
				for _, v := range ipAlt {
					parsedIP := net.ParseIP(v)
					if parsedIP == nil {
						return nil, nil, errutil.UserError{Err: fmt.Sprintf(
							"the value %q is not a valid IP address", v)}
					}
					ipAddresses = append(ipAddresses, parsedIP)
				}
			}
		}
	}

	URIs := []*url.URL{}
	{
		if csr != nil && role.UseCSRSANs {
			if len(csr.URIs) > 0 {
				if len(role.AllowedURISANs) == 0 {
					return nil, nil, errutil.UserError{
						Err: "URI Subject Alternative Names are not allowed in this role, but were provided via CSR",
					}
				}

				// validate uri sans
				for _, uri := range csr.URIs {
					valid := ValidateURISAN(b, role, entityInfo, uri.String())
					if !valid {
						return nil, nil, errutil.UserError{
							Err: "URI Subject Alternative Names were provided via CSR which are not valid for this role",
						}
					}

					URIs = append(URIs, uri)
				}
			}
		} else {
			uriAlt := cb.GetURISans()
			if len(uriAlt) > 0 {
				if len(role.AllowedURISANs) == 0 {
					return nil, nil, errutil.UserError{
						Err: "URI Subject Alternative Names are not allowed in this role, but were provided via the API",
					}
				}

				for _, uri := range uriAlt {
					valid := ValidateURISAN(b, role, entityInfo, uri)
					if !valid {
						return nil, nil, errutil.UserError{
							Err: "URI Subject Alternative Names were provided via the API which are not valid for this role",
						}
					}

					parsedURI, err := url.Parse(uri)
					if parsedURI == nil || err != nil {
						return nil, nil, errutil.UserError{
							Err: fmt.Sprintf(
								"the provided URI Subject Alternative Name %q is not a valid URI", uri),
						}
					}

					URIs = append(URIs, parsedURI)
				}
			}
		}
	}

	// Most of these could also be RemoveDuplicateStable, or even
	// leave duplicates in, but OU is the one most likely to be duplicated.
	subject := pkix.Name{
		CommonName:         cn,
		SerialNumber:       ridSerialNumber,
		Country:            strutil.RemoveDuplicatesStable(role.Country, false),
		Organization:       strutil.RemoveDuplicatesStable(role.Organization, false),
		OrganizationalUnit: strutil.RemoveDuplicatesStable(role.OU, false),
		Locality:           strutil.RemoveDuplicatesStable(role.Locality, false),
		Province:           strutil.RemoveDuplicatesStable(role.Province, false),
		StreetAddress:      strutil.RemoveDuplicatesStable(role.StreetAddress, false),
		PostalCode:         strutil.RemoveDuplicatesStable(role.PostalCode, false),
	}

	// Get the TTL and verify it against the max allowed
	notAfter, ttlWarnings, err := GetCertificateNotAfter(b, role, cb, caSign)
	if err != nil {
		return nil, warnings, err
	}
	warnings = append(warnings, ttlWarnings...)

	// Parse SKID from the request for cross-signing.
	var skid []byte
	{
		if rawSKIDValue, ok := cb.GetOptionalSkid(); ok {
			// Handle removing common separators to make copy/paste from tool
			// output easier. Chromium uses space, OpenSSL uses colons, and at
			// one point, Vault had preferred dash as a separator for hex
			// strings.
			var err error
			skidValue := rawSKIDValue.(string)
			for _, separator := range []string{":", "-", " "} {
				skidValue = strings.ReplaceAll(skidValue, separator, "")
			}

			skid, err = hex.DecodeString(skidValue)
			if err != nil {
				return nil, nil, errutil.UserError{Err: fmt.Sprintf("cannot parse requested SKID value as hex: %v", err)}
			}
		}
	}

	// Add UserIDs into the Subject, if the request type supports it.
	if _, present := cb.IsUserIdInSchema(); present {
		rawUserIDs := cb.GetUserIds()

		// Only take UserIDs from CSR if one was not supplied via API.
		if len(rawUserIDs) == 0 && csr != nil {
			for _, attr := range csr.Subject.Names {
				if attr.Type.Equal(certutil.SubjectPilotUserIDAttributeOID) {
					switch aValue := attr.Value.(type) {
					case string:
						rawUserIDs = append(rawUserIDs, aValue)
					case []byte:
						rawUserIDs = append(rawUserIDs, string(aValue))
					default:
						return nil, nil, errutil.UserError{Err: "unknown type for user_id attribute in CSR's Subject"}
					}
				}
			}
		}

		// Check for bad userIDs and add to the subject.
		if len(rawUserIDs) > 0 {
			for _, value := range rawUserIDs {
				if !ValidateUserId(role, value) {
					return nil, nil, errutil.UserError{Err: fmt.Sprintf("user_id %v is not allowed by this role", value)}
				}

				subject.ExtraNames = append(subject.ExtraNames, pkix.AttributeTypeAndValue{
					Type:  certutil.SubjectPilotUserIDAttributeOID,
					Value: value,
				})
			}
		}
	}

	creation := &certutil.CreationBundle{
		Params: &certutil.CreationParameters{
			Subject:                       subject,
			DNSNames:                      strutil.RemoveDuplicates(dnsNames, false),
			EmailAddresses:                strutil.RemoveDuplicates(emailAddresses, false),
			IPAddresses:                   ipAddresses,
			URIs:                          URIs,
			OtherSANs:                     otherSANs,
			KeyType:                       role.KeyType,
			KeyBits:                       role.KeyBits,
			SignatureBits:                 role.SignatureBits,
			UsePSS:                        role.UsePSS,
			NotAfter:                      notAfter,
			KeyUsage:                      x509.KeyUsage(parsing.ParseKeyUsages(role.KeyUsage)),
			ExtKeyUsage:                   ParseExtKeyUsagesFromRole(role),
			ExtKeyUsageOIDs:               role.ExtKeyUsageOIDs,
			PolicyIdentifiers:             role.PolicyIdentifiers,
			BasicConstraintsValidForNonCA: role.BasicConstraintsValidForNonCA,
			NotBeforeDuration:             role.NotBeforeDuration,
			ForceAppendCaChain:            caSign != nil,
			SKID:                          skid,
			IgnoreCSRSignature:            cb.IgnoreCSRSignature(),
		},
		SigningBundle: caSign,
		CSR:           csr,
	}

	// Don't deal with URLs or max path length if it's self-signed, as these
	// normally come from the signing bundle
	if caSign == nil {
		return creation, warnings, nil
	}

	// This will have been read in from the getGlobalAIAURLs function
	creation.Params.URLs = caSign.URLs

	// If the max path length in the role is not nil, it was specified at
	// generation time with the max_path_length parameter; otherwise derive it
	// from the signing certificate
	if role.MaxPathLength != nil {
		creation.Params.MaxPathLength = *role.MaxPathLength
	} else {
		switch {
		case caSign.Certificate.MaxPathLen < 0:
			creation.Params.MaxPathLength = -1
		case caSign.Certificate.MaxPathLen == 0 &&
			caSign.Certificate.MaxPathLenZero:
			// The signing function will ensure that we do not issue a CA cert
			creation.Params.MaxPathLength = 0
		default:
			// If this takes it to zero, we handle this case later if
			// necessary
			creation.Params.MaxPathLength = caSign.Certificate.MaxPathLen - 1
		}
	}

	return creation, warnings, nil
}

// Given a set of requested names for a certificate, verifies that all of them
// match the various toggles set in the role for controlling issuance.
// If one does not pass, it is returned in the string argument.
func ValidateNames(b logical.SystemView, role *RoleEntry, entityInfo EntityInfo, names []string) string {
	for _, name := range names {
		// Previously, reducedName was called sanitizedName but this made
		// little sense under the previous interpretation of wildcards,
		// leading to two bugs in this implementation. We presently call it
		// "reduced" to indicate that it is still untrusted input (potentially
		// different from the bare Common Name entry we're validating), it
		// might have been modified such as by the removal of wildcard labels
		// or the email prefix.
		reducedName := name
		emailDomain := reducedName
		wildcardLabel := ""
		isEmail := false
		isWildcard := false

		// If it has an @, assume it is an email address and separate out the
		// user from the hostname portion so that we can act on the hostname.
		// Note that this matches behavior from the alt_names parameter. If it
		// ends up being problematic for users, I guess that could be separated
		// into dns_names and email_names in the future to be explicit, but I
		// don't think this is likely.
		if strings.Contains(reducedName, "@") {
			splitEmail := strings.Split(reducedName, "@")
			if len(splitEmail) != 2 {
				return name
			}
			reducedName = splitEmail[1]
			emailDomain = splitEmail[1]
			isEmail = true
		}

		if IsWildcardDomain(reducedName) {
			// Regardless of later rejections below, this common name contains
			// a wildcard character and is thus technically a wildcard name.
			isWildcard = true

			// Additionally, if AllowWildcardCertificates is explicitly
			// forbidden, it takes precedence over AllowAnyName, thus we should
			// reject the name now.
			//
			// We expect the role to have been correctly migrated but guard for
			// safety.
			if role.AllowWildcardCertificates != nil && !*role.AllowWildcardCertificates {
				return name
			}

			// Check that this domain is well-formatted per RFC 6125.
			var err error
			wildcardLabel, reducedName, err = ValidateWildcardDomain(reducedName)
			if err != nil {
				return name
			}
		}

		// Email addresses using wildcard domain names do not make sense
		// in a Common Name field.
		if isEmail && isWildcard {
			return name
		}

		// AllowAnyName is checked after this because EnforceHostnames still
		// applies when allowing any name. Also, we check the reduced name to
		// ensure that we are not either checking a full email address or a
		// wildcard prefix.
		if role.EnforceHostnames {
			if reducedName != "" {
				// See note above about splitLabels having only one segment
				// and setting reducedName to the empty string.
				p := idna.New(
					idna.StrictDomainName(true),
					idna.VerifyDNSLength(true),
				)
				converted, err := p.ToASCII(reducedName)
				if err != nil {
					return name
				}
				if !hostnameRegex.MatchString(converted) {
					return name
				}
			}

			// When a wildcard is specified, we additionally need to validate
			// the label with the wildcard is correctly formed.
			if isWildcard && !leftWildLabelRegex.MatchString(wildcardLabel) {
				return name
			}
		}

		// Self-explanatory, but validations from EnforceHostnames and
		// AllowWildcardCertificates take precedence.
		if role.AllowAnyName {
			continue
		}

		// The following blocks all work the same basic way:
		// 1) If a role allows a certain class of base (localhost, token
		// display name, role-configured domains), perform further tests
		//
		// 2) If there is a perfect match on either the sanitized name or it's an
		// email address with a perfect match on the hostname portion, allow it
		//
		// 3) If subdomains are allowed, we check based on the sanitized name;
		// note that if not a wildcard, will be equivalent to the email domain
		// for email checks, and we already checked above for both a wildcard
		// and email address being present in the same name
		// 3a) First we check for a non-wildcard subdomain, as in <name>.<base>
		// 3b) Then we check if it's a wildcard and the base domain is a match
		//
		// Variances are noted in-line

		if role.AllowLocalhost {
			if reducedName == "localhost" ||
				reducedName == "localdomain" ||
				(isEmail && emailDomain == "localhost") ||
				(isEmail && emailDomain == "localdomain") {
				continue
			}

			if role.AllowSubdomains {
				// It is possible, if unlikely, to have a subdomain of "localhost"
				if strings.HasSuffix(reducedName, ".localhost") ||
					(isWildcard && reducedName == "localhost") {
					continue
				}

				// A subdomain of "localdomain" is also not entirely uncommon
				if strings.HasSuffix(reducedName, ".localdomain") ||
					(isWildcard && reducedName == "localdomain") {
					continue
				}
			}
		}

		if role.AllowTokenDisplayName {
			if name == entityInfo.DisplayName {
				continue
			}

			if role.AllowSubdomains {
				if isEmail {
					// If it's an email address, we need to parse the token
					// display name in order to do a proper comparison of the
					// subdomain
					if strings.Contains(entityInfo.DisplayName, "@") {
						splitDisplay := strings.Split(entityInfo.DisplayName, "@")
						if len(splitDisplay) == 2 {
							// Compare the sanitized name against the hostname
							// portion of the email address in the broken
							// display name
							if strings.HasSuffix(reducedName, "."+splitDisplay[1]) {
								continue
							}
						}
					}
				}

				if strings.HasSuffix(reducedName, "."+entityInfo.DisplayName) ||
					(isWildcard && reducedName == entityInfo.DisplayName) {
					continue
				}
			}
		}

		if len(role.AllowedDomains) > 0 {
			valid := false
			for _, currDomain := range role.AllowedDomains {
				// If there is, say, a trailing comma, ignore it
				if currDomain == "" {
					continue
				}

				if role.AllowedDomainsTemplate {
					isTemplate, _ := framework.ValidateIdentityTemplate(currDomain)
					if isTemplate && entityInfo.EntityID != "" {
						tmpCurrDomain, err := framework.PopulateIdentityTemplate(currDomain, entityInfo.EntityID, b)
						if err != nil {
							continue
						}

						currDomain = tmpCurrDomain
					}
				}

				// First, allow an exact match of the base domain if that role flag
				// is enabled
				if role.AllowBareDomains &&
					(strings.EqualFold(name, currDomain) ||
						(isEmail && strings.EqualFold(emailDomain, currDomain))) {
					valid = true
					break
				}

				if role.AllowSubdomains {
					if strings.HasSuffix(reducedName, "."+currDomain) ||
						(isWildcard && strings.EqualFold(reducedName, currDomain)) {
						valid = true
						break
					}
				}

				if role.AllowGlobDomains &&
					strings.Contains(currDomain, "*") &&
					glob.Glob(strings.ToLower(currDomain), strings.ToLower(name)) {
					valid = true
					break
				}
			}

			if valid {
				continue
			}
		}

		return name
	}

	return ""
}

func IsWildcardDomain(name string) bool {
	// Per RFC 6125 Section 6.4.3, and explicitly contradicting the earlier
	// RFC 2818 which no modern client will validate against, there are two
	// main types of wildcards, each with a single wildcard specifier (`*`,
	// functionally different from the `*` used as a glob from the
	// AllowGlobDomains parsing path) in the left-most label:
	//
	//  1. Entire label is a single wildcard character (most common and
	//     well-supported),
	//  2. Part of the label contains a single wildcard character (e.g. per
	//     RFC 6125: baz*.example.net, *baz.example.net, or b*z.example.net).
	//
	// We permit issuance of both but not the older RFC 2818 style under
	// the new AllowWildcardCertificates option. However, anything with a
	// glob character is technically a wildcard, though not a valid one.

	return strings.Contains(name, "*")
}

func ValidateWildcardDomain(name string) (string, string, error) {
	// See note in isWildcardDomain(...) about the definition of a wildcard
	// domain.
	var wildcardLabel string
	var reducedName string

	if strings.Count(name, "*") > 1 {
		// As mentioned above, only one wildcard character is permitted
		// under RFC 6125 semantics.
		return wildcardLabel, reducedName, fmt.Errorf("expected only one wildcard identifier in the given domain name")
	}

	// Split the Common Name into two parts: a left-most label and the
	// remaining segments (if present).
	splitLabels := strings.SplitN(name, ".", 2)
	if len(splitLabels) != 2 {
		// We've been given a single-part domain name that consists
		// entirely of a wildcard. This is a little tricky to handle,
		// but EnforceHostnames validates both the wildcard-containing
		// label and the reduced name, but _only_ the latter if it is
		// non-empty. This allows us to still validate the only label
		// component matches hostname expectations still.
		wildcardLabel = splitLabels[0]
		reducedName = ""
	} else {
		// We have a (at least) two label domain name. But before we can
		// update our names, we need to validate the wildcard ended up
		// in the segment we expected it to. While this is (kinda)
		// validated under EnforceHostnames's leftWildLabelRegex, we
		// still need to validate it in the non-enforced mode.
		//
		// By validated assumption above, we know there's strictly one
		// wildcard in this domain so we only need to check the wildcard
		// label or the reduced name (as one is equivalent to the other).
		// Because we later assume reducedName _lacks_ wildcard segments,
		// we validate that.
		wildcardLabel = splitLabels[0]
		reducedName = splitLabels[1]
		if strings.Contains(reducedName, "*") {
			return wildcardLabel, reducedName, fmt.Errorf("expected wildcard to only be present in left-most domain label")
		}
	}

	return wildcardLabel, reducedName, nil
}

// ValidateCommonName Validates a given common name, ensuring it's either an email or a hostname
// after validating it according to the role parameters, or disables
// validation altogether.
func ValidateCommonName(b logical.SystemView, role *RoleEntry, entityInfo EntityInfo, name string) string {
	isDisabled := len(role.CNValidations) == 1 && role.CNValidations[0] == "disabled"
	if isDisabled {
		return ""
	}

	if ValidateNames(b, role, entityInfo, []string{name}) != "" {
		return name
	}

	// Validations weren't disabled, but the role lacked CN Validations, so
	// don't restrict types. This case is hit in certain existing tests.
	if len(role.CNValidations) == 0 {
		return ""
	}

	// If there's an at in the data, ensure email type validation is allowed.
	// Otherwise, ensure hostname is allowed.
	if strings.Contains(name, "@") {
		var allowsEmails bool
		for _, validation := range role.CNValidations {
			if validation == "email" {
				allowsEmails = true
				break
			}
		}
		if !allowsEmails {
			return name
		}
	} else {
		var allowsHostnames bool
		for _, validation := range role.CNValidations {
			if validation == "hostname" {
				allowsHostnames = true
				break
			}
		}
		if !allowsHostnames {
			return name
		}
	}

	return ""
}

// ValidateOtherSANs checks if the values requested are allowed. If an OID
// isn't allowed, it will be returned as the first string. If a Value isn't
// allowed, it will be returned as the second string. Empty strings + error
// means everything is okay.
func ValidateOtherSANs(role *RoleEntry, requested map[string][]string) (string, string, error) {
	if len(role.AllowedOtherSANs) == 1 && role.AllowedOtherSANs[0] == "*" {
		// Anything is allowed
		return "", "", nil
	}

	allowed, err := ParseOtherSANs(role.AllowedOtherSANs)
	if err != nil {
		return "", "", fmt.Errorf("error parsing role's allowed SANs: %w", err)
	}
	for oid, names := range requested {
		for _, name := range names {
			allowedNames, ok := allowed[oid]
			if !ok {
				return oid, "", nil
			}

			valid := false
			for _, allowedName := range allowedNames {
				if glob.Glob(allowedName, name) {
					valid = true
					break
				}
			}

			if !valid {
				return oid, name, nil
			}
		}
	}

	return "", "", nil
}

func ParseOtherSANs(others []string) (map[string][]string, error) {
	result := map[string][]string{}
	for _, other := range others {
		splitOther := strings.SplitN(other, ";", 2)
		if len(splitOther) != 2 {
			return nil, fmt.Errorf("expected a semicolon in other SAN %q", other)
		}
		splitType := strings.SplitN(splitOther[1], ":", 2)
		if len(splitType) != 2 {
			return nil, fmt.Errorf("expected a colon in other SAN %q", other)
		}
		switch {
		case strings.EqualFold(splitType[0], "utf8"):
		case strings.EqualFold(splitType[0], "utf-8"):
		default:
			return nil, fmt.Errorf("only utf8 other SANs are supported; found non-supported type in other SAN %q", other)
		}
		result[splitOther[0]] = append(result[splitOther[0]], splitType[1])
	}

	return result, nil
}

// Given a URI SAN, verify that it is allowed.
func ValidateURISAN(b logical.SystemView, role *RoleEntry, entityInfo EntityInfo, uri string) bool {
	valid := false
	for _, allowed := range role.AllowedURISANs {
		if role.AllowedURISANsTemplate {
			isTemplate, _ := framework.ValidateIdentityTemplate(allowed)
			if isTemplate && entityInfo.EntityID != "" {
				tmpAllowed, err := framework.PopulateIdentityTemplate(allowed, entityInfo.EntityID, b)
				if err != nil {
					continue
				}
				allowed = tmpAllowed
			}
		}
		validURI := glob.Glob(allowed, uri)
		if validURI {
			valid = true
			break
		}
	}
	return valid
}

// ValidateUserId Returns bool stating whether the given UserId is Valid
func ValidateUserId(role *RoleEntry, userId string) bool {
	allowedList := role.AllowedUserIDs

	if len(allowedList) == 0 {
		// Nothing is allowed.
		return false
	}

	if strutil.StrListContainsCaseInsensitive(allowedList, userId) {
		return true
	}

	for _, rolePattern := range allowedList {
		if rolePattern == "" {
			continue
		}

		if strings.Contains(rolePattern, "*") && glob.Glob(rolePattern, userId) {
			return true
		}
	}

	// No matches.
	return false
}

func ValidateSerialNumber(role *RoleEntry, serialNumber string) string {
	valid := false
	if len(role.AllowedSerialNumbers) > 0 {
		for _, currSerialNumber := range role.AllowedSerialNumbers {
			if currSerialNumber == "" {
				continue
			}

			if (strings.Contains(currSerialNumber, "*") &&
				glob.Glob(currSerialNumber, serialNumber)) ||
				currSerialNumber == serialNumber {
				valid = true
				break
			}
		}
	}
	if !valid {
		return serialNumber
	} else {
		return ""
	}
}

type CertNotAfterInput interface {
	GetTTL() int
	GetOptionalNotAfter() (interface{}, bool)
}

// GetCertificateNotAfter compute a certificate's NotAfter date based on the mount ttl, role, signing bundle and input
// api data being sent. Returns a NotAfter time, a set of warnings or an error.
func GetCertificateNotAfter(b logical.SystemView, role *RoleEntry, input CertNotAfterInput, caSign *certutil.CAInfoBundle) (time.Time, []string, error) {
	var warnings []string
	var maxTTL time.Duration
	var notAfter time.Time
	var err error

	ttl := time.Duration(input.GetTTL()) * time.Second
	notAfterAlt := role.NotAfter
	if notAfterAlt == "" {
		notAfterAltRaw, ok := input.GetOptionalNotAfter()
		if ok {
			notAfterAlt = notAfterAltRaw.(string)
		}
	}
	if ttl > 0 && notAfterAlt != "" {
		return time.Time{}, warnings, errutil.UserError{Err: "Either ttl or not_after should be provided. Both should not be provided in the same request."}
	}

	if ttl == 0 && role.TTL > 0 {
		ttl = role.TTL
	}

	if role.MaxTTL > 0 {
		maxTTL = role.MaxTTL
	}

	if ttl == 0 {
		ttl = b.DefaultLeaseTTL()
	}
	if maxTTL == 0 {
		maxTTL = b.MaxLeaseTTL()
	}
	if ttl > maxTTL {
		warnings = append(warnings, fmt.Sprintf("TTL %q is longer than permitted maxTTL %q, so maxTTL is being used", ttl, maxTTL))
		ttl = maxTTL
	}

	if notAfterAlt != "" {
		notAfter, err = time.Parse(time.RFC3339, notAfterAlt)
		if err != nil {
			return notAfter, warnings, errutil.UserError{Err: err.Error()}
		}
	} else {
		notAfter = time.Now().Add(ttl)
	}
	notAfter, err = ApplyIssuerLeafNotAfterBehavior(caSign, notAfter)
	if err != nil {
		return time.Time{}, warnings, err
	}
	return notAfter, warnings, nil
}

// ApplyIssuerLeafNotAfterBehavior resets a certificate's notAfter time or errors out based on the
// issuer's notAfter date along with the LeafNotAfterBehavior configuration
func ApplyIssuerLeafNotAfterBehavior(caSign *certutil.CAInfoBundle, notAfter time.Time) (time.Time, error) {
	if caSign != nil && notAfter.After(caSign.Certificate.NotAfter) {
		// If it's not self-signed, verify that the issued certificate
		// won't be valid past the lifetime of the CA certificate, and
		// act accordingly. This is dependent based on the issuer's
		// LeafNotAfterBehavior argument.
		switch caSign.LeafNotAfterBehavior {
		case certutil.PermitNotAfterBehavior:
			// Explicitly do nothing.
		case certutil.TruncateNotAfterBehavior:
			notAfter = caSign.Certificate.NotAfter
		case certutil.ErrNotAfterBehavior, certutil.AlwaysEnforceErr:
			fallthrough
		default:
			return time.Time{}, errutil.UserError{Err: fmt.Sprintf(
				"cannot satisfy request, as TTL would result in notAfter of %s that is beyond the expiration of the CA certificate at %s", notAfter.UTC().Format(time.RFC3339Nano), caSign.Certificate.NotAfter.UTC().Format(time.RFC3339Nano))}
		}
	}
	return notAfter, nil
}

// StoreCertificate given a certificate bundle that was signed, persist the certificate to storage
func StoreCertificate(ctx context.Context, s logical.Storage, certCounter CertificateCounter, certBundle *certutil.ParsedCertBundle) error {
	hyphenSerialNumber := parsing.NormalizeSerialForStorageFromBigInt(certBundle.Certificate.SerialNumber)
	key := PathCerts + hyphenSerialNumber
	certsCounted := certCounter.IsInitialized()
	err := s.Put(ctx, &logical.StorageEntry{
		Key:   key,
		Value: certBundle.CertificateBytes,
	})
	if err != nil {
		return fmt.Errorf("unable to store certificate locally: %w", err)
	}
	certCounter.IncrementTotalCertificatesCount(certsCounted, key)
	return nil
}
