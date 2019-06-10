package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
	"golang.org/x/net/idna"
)

type inputBundle struct {
	role    *roleEntry
	req     *logical.Request
	apiData *framework.FieldData
}

var (
	// A note on hostnameRegex: although we set the StrictDomainName option
	// when doing the idna conversion, this appears to only affect output, not
	// input, so it will allow e.g. host^123.example.com straight through. So
	// we still need to use this to check the output.
	hostnameRegex = regexp.MustCompile(`^(\*\.)?(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
)

func oidInExtensions(oid asn1.ObjectIdentifier, extensions []pkix.Extension) bool {
	for _, e := range extensions {
		if e.Id.Equal(oid) {
			return true
		}
	}
	return false
}

func getFormat(data *framework.FieldData) string {
	format := data.Get("format").(string)
	switch format {
	case "pem":
	case "der":
	case "pem_bundle":
	default:
		format = ""
	}
	return format
}

// Fetches the CA info. Unlike other certificates, the CA info is stored
// in the backend as a CertBundle, because we are storing its private key
func fetchCAInfo(ctx context.Context, req *logical.Request) (*certutil.CAInfoBundle, error) {
	bundleEntry, err := req.Storage.Get(ctx, "config/ca_bundle")
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch local CA certificate/key: %v", err)}
	}
	if bundleEntry == nil {
		return nil, errutil.UserError{Err: "backend must be configured with a CA certificate/key"}
	}

	var bundle certutil.CertBundle
	if err := bundleEntry.DecodeJSON(&bundle); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode local CA certificate/key: %v", err)}
	}

	parsedBundle, err := bundle.ToParsedCertBundle()
	if err != nil {
		return nil, errutil.InternalError{Err: err.Error()}
	}

	if parsedBundle.Certificate == nil {
		return nil, errutil.InternalError{Err: "stored CA information not able to be parsed"}
	}

	caInfo := &certutil.CAInfoBundle{*parsedBundle, nil}

	entries, err := getURLs(ctx, req)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch URL information: %v", err)}
	}
	if entries == nil {
		entries = &certutil.URLEntries{
			IssuingCertificates:   []string{},
			CRLDistributionPoints: []string{},
			OCSPServers:           []string{},
		}
	}
	caInfo.URLs = entries

	return caInfo, nil
}

// Allows fetching certificates from the backend; it handles the slightly
// separate pathing for CA, CRL, and revoked certificates.
func fetchCertBySerial(ctx context.Context, req *logical.Request, prefix, serial string) (*logical.StorageEntry, error) {
	var path, legacyPath string
	var err error
	var certEntry *logical.StorageEntry

	hyphenSerial := normalizeSerial(serial)
	colonSerial := strings.Replace(strings.ToLower(serial), "-", ":", -1)

	switch {
	// Revoked goes first as otherwise ca/crl get hardcoded paths which fail if
	// we actually want revocation info
	case strings.HasPrefix(prefix, "revoked/"):
		legacyPath = "revoked/" + colonSerial
		path = "revoked/" + hyphenSerial
	case serial == "ca":
		path = "ca"
	case serial == "crl":
		path = "crl"
	default:
		legacyPath = "certs/" + colonSerial
		path = "certs/" + hyphenSerial
	}

	certEntry, err = req.Storage.Get(ctx, path)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error fetching certificate %s: %s", serial, err)}
	}
	if certEntry != nil {
		if certEntry.Value == nil || len(certEntry.Value) == 0 {
			return nil, errutil.InternalError{Err: fmt.Sprintf("returned certificate bytes for serial %s were empty", serial)}
		}
		return certEntry, nil
	}

	// If legacyPath is unset, it's going to be a CA or CRL; return immediately
	if legacyPath == "" {
		return nil, nil
	}

	// Retrieve the old-style path.  We disregard errors here because they
	// always manifest on windows, and thus the initial check for a revoked
	// cert fails would return an error when the cert isn't revoked, preventing
	// the happy path from working.
	certEntry, _ = req.Storage.Get(ctx, legacyPath)
	if certEntry == nil {
		return nil, nil
	}
	if certEntry.Value == nil || len(certEntry.Value) == 0 {
		return nil, errutil.InternalError{Err: fmt.Sprintf("returned certificate bytes for serial %s were empty", serial)}
	}

	// Update old-style paths to new-style paths
	certEntry.Key = path
	if err = req.Storage.Put(ctx, certEntry); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error saving certificate with serial %s to new location", serial)}
	}
	if err = req.Storage.Delete(ctx, legacyPath); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error deleting certificate with serial %s from old location", serial)}
	}

	return certEntry, nil
}

// Given a set of requested names for a certificate, verifies that all of them
// match the various toggles set in the role for controlling issuance.
// If one does not pass, it is returned in the string argument.
func validateNames(data *inputBundle, names []string) string {
	for _, name := range names {
		sanitizedName := name
		emailDomain := name
		isEmail := false
		isWildcard := false

		// If it has an @, assume it is an email address and separate out the
		// user from the hostname portion so that we can act on the hostname.
		// Note that this matches behavior from the alt_names parameter. If it
		// ends up being problematic for users, I guess that could be separated
		// into dns_names and email_names in the future to be explicit, but I
		// don't think this is likely.
		if strings.Contains(name, "@") {
			splitEmail := strings.Split(name, "@")
			if len(splitEmail) != 2 {
				return name
			}
			sanitizedName = splitEmail[1]
			emailDomain = splitEmail[1]
			isEmail = true
		}

		// If we have an asterisk as the first part of the domain name, mark it
		// as wildcard and set the sanitized name to the remainder of the
		// domain
		if strings.HasPrefix(sanitizedName, "*.") {
			sanitizedName = sanitizedName[2:]
			isWildcard = true
		}

		// Email addresses using wildcard domain names do not make sense
		if isEmail && isWildcard {
			return name
		}

		// AllowAnyName is checked after this because EnforceHostnames still
		// applies when allowing any name. Also, we check the sanitized name to
		// ensure that we are not either checking a full email address or a
		// wildcard prefix.
		if data.role.EnforceHostnames {
			p := idna.New(
				idna.StrictDomainName(true),
				idna.VerifyDNSLength(true),
			)
			converted, err := p.ToASCII(sanitizedName)
			if err != nil {
				return name
			}
			if !hostnameRegex.MatchString(converted) {
				return name
			}
		}

		// Self-explanatory
		if data.role.AllowAnyName {
			continue
		}

		// The following blocks all work the same basic way:
		// 1) If a role allows a certain class of base (localhost, token
		// display name, role-configured domains), perform further tests
		//
		// 2) If there is a perfect match on either the name itself or it's an
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

		if data.role.AllowLocalhost {
			if name == "localhost" ||
				name == "localdomain" ||
				(isEmail && emailDomain == "localhost") ||
				(isEmail && emailDomain == "localdomain") {
				continue
			}

			if data.role.AllowSubdomains {
				// It is possible, if unlikely, to have a subdomain of "localhost"
				if strings.HasSuffix(sanitizedName, ".localhost") ||
					(isWildcard && sanitizedName == "localhost") {
					continue
				}

				// A subdomain of "localdomain" is also not entirely uncommon
				if strings.HasSuffix(sanitizedName, ".localdomain") ||
					(isWildcard && sanitizedName == "localdomain") {
					continue
				}
			}
		}

		if data.role.AllowTokenDisplayName {
			if name == data.req.DisplayName {
				continue
			}

			if data.role.AllowSubdomains {
				if isEmail {
					// If it's an email address, we need to parse the token
					// display name in order to do a proper comparison of the
					// subdomain
					if strings.Contains(data.req.DisplayName, "@") {
						splitDisplay := strings.Split(data.req.DisplayName, "@")
						if len(splitDisplay) == 2 {
							// Compare the sanitized name against the hostname
							// portion of the email address in the broken
							// display name
							if strings.HasSuffix(sanitizedName, "."+splitDisplay[1]) {
								continue
							}
						}
					}
				}

				if strings.HasSuffix(sanitizedName, "."+data.req.DisplayName) ||
					(isWildcard && sanitizedName == data.req.DisplayName) {
					continue
				}
			}
		}

		if len(data.role.AllowedDomains) > 0 {
			valid := false
			for _, currDomain := range data.role.AllowedDomains {
				// If there is, say, a trailing comma, ignore it
				if currDomain == "" {
					continue
				}

				// First, allow an exact match of the base domain if that role flag
				// is enabled
				if data.role.AllowBareDomains &&
					(name == currDomain ||
						(isEmail && emailDomain == currDomain)) {
					valid = true
					break
				}

				if data.role.AllowSubdomains {
					if strings.HasSuffix(sanitizedName, "."+currDomain) ||
						(isWildcard && sanitizedName == currDomain) {
						valid = true
						break
					}
				}

				if data.role.AllowGlobDomains &&
					strings.Contains(currDomain, "*") &&
					glob.Glob(currDomain, name) {
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

// validateOtherSANs checks if the values requested are allowed. If an OID
// isn't allowed, it will be returned as the first string. If a value isn't
// allowed, it will be returned as the second string. Empty strings + error
// means everything is okay.
func validateOtherSANs(data *inputBundle, requested map[string][]string) (string, string, error) {
	for _, val := range data.role.AllowedOtherSANs {
		if val == "*" {
			// Anything is allowed
			return "", "", nil
		}
	}
	allowed, err := parseOtherSANs(data.role.AllowedOtherSANs)
	if err != nil {
		return "", "", errwrap.Wrapf("error parsing role's allowed SANs: {{err}}", err)
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

func parseOtherSANs(others []string) (map[string][]string, error) {
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

func validateSerialNumber(data *inputBundle, serialNumber string) string {
	valid := false
	if len(data.role.AllowedSerialNumbers) > 0 {
		for _, currSerialNumber := range data.role.AllowedSerialNumbers {
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

func generateCert(ctx context.Context,
	b *backend,
	input *inputBundle,
	caSign *certutil.CAInfoBundle,
	isCA bool) (*certutil.ParsedCertBundle, error) {

	if input.role == nil {
		return nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	if input.role.KeyType == "rsa" && input.role.KeyBits < 2048 {
		return nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
	}

	data, err := generateCreationBundle(b, input, caSign, nil)
	if err != nil {
		return nil, err
	}
	if data.Params == nil {
		return nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	if isCA {
		data.Params.IsCA = isCA
		data.Params.PermittedDNSDomains = input.apiData.Get("permitted_dns_domains").([]string)

		if data.SigningBundle == nil {
			// Generating a self-signed root certificate
			entries, err := getURLs(ctx, input.req)
			if err != nil {
				return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch URL information: %v", err)}
			}
			if entries == nil {
				entries = &certutil.URLEntries{
					IssuingCertificates:   []string{},
					CRLDistributionPoints: []string{},
					OCSPServers:           []string{},
				}
			}
			data.Params.URLs = entries

			if input.role.MaxPathLength == nil {
				data.Params.MaxPathLength = -1
			} else {
				data.Params.MaxPathLength = *input.role.MaxPathLength
			}
		}
	}

	parsedBundle, err := certutil.CreateCertificate(data)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

// N.B.: This is only meant to be used for generating intermediate CAs.
// It skips some sanity checks.
func generateIntermediateCSR(b *backend, input *inputBundle) (*certutil.ParsedCSRBundle, error) {
	creation, err := generateCreationBundle(b, input, nil, nil)
	if err != nil {
		return nil, err
	}
	if creation.Params == nil {
		return nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	addBasicConstraints := input.apiData != nil && input.apiData.Get("add_basic_constraints").(bool)
	parsedBundle, err := certutil.CreateCSR(creation, addBasicConstraints)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

func signCert(b *backend,
	data *inputBundle,
	caSign *certutil.CAInfoBundle,
	isCA bool,
	useCSRValues bool) (*certutil.ParsedCertBundle, error) {

	if data.role == nil {
		return nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	csrString := data.apiData.Get("csr").(string)
	if csrString == "" {
		return nil, errutil.UserError{Err: fmt.Sprintf("\"csr\" is empty")}
	}

	pemBytes := []byte(csrString)
	pemBlock, pemBytes := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, errutil.UserError{Err: "csr contains no data"}
	}
	csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("certificate request could not be parsed: %v", err)}
	}

	switch data.role.KeyType {
	case "rsa":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.RSA {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				data.role.KeyType)}
		}
		pubKey, ok := csr.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		// Verify that the key is at least 2048 bits
		if pubKey.N.BitLen() < 2048 {
			return nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
		}

		// Verify that the bit size is at least the size specified in the role
		if pubKey.N.BitLen() < data.role.KeyBits {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				data.role.KeyBits,
				pubKey.N.BitLen())}
		}

	case "ec":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.ECDSA {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				data.role.KeyType)}
		}
		pubKey, ok := csr.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		// Verify that the bit size is at least the size specified in the role
		if pubKey.Params().BitSize < data.role.KeyBits {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				data.role.KeyBits,
				pubKey.Params().BitSize)}
		}

	case "any":
		// We only care about running RSA < 2048 bit checks, so if not RSA
		// break out
		if csr.PublicKeyAlgorithm != x509.RSA {
			break
		}

		// Run RSA < 2048 bit checks
		pubKey, ok := csr.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}
		if pubKey.N.BitLen() < 2048 {
			return nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
		}

	}

	creation, err := generateCreationBundle(b, data, caSign, csr)
	if err != nil {
		return nil, err
	}
	if creation.Params == nil {
		return nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	creation.Params.IsCA = isCA
	creation.Params.UseCSRValues = useCSRValues

	if isCA {
		creation.Params.PermittedDNSDomains = data.apiData.Get("permitted_dns_domains").([]string)
	}

	parsedBundle, err := certutil.SignCertificate(creation)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

// generateCreationBundle is a shared function that reads parameters supplied
// from the various endpoints and generates a CreationParameters with the
// parameters that can be used to issue or sign
func generateCreationBundle(b *backend, data *inputBundle, caSign *certutil.CAInfoBundle, csr *x509.CertificateRequest) (*certutil.CreationBundle, error) {
	// Read in names -- CN, DNS and email addresses
	var cn string
	var ridSerialNumber string
	dnsNames := []string{}
	emailAddresses := []string{}
	{
		if csr != nil && data.role.UseCSRCommonName {
			cn = csr.Subject.CommonName
		}
		if cn == "" {
			cn = data.apiData.Get("common_name").(string)
			if cn == "" && data.role.RequireCN {
				return nil, errutil.UserError{Err: `the common_name field is required, or must be provided in a CSR with "use_csr_common_name" set to true, unless "require_cn" is set to false`}
			}
		}

		ridSerialNumber = data.apiData.Get("serial_number").(string)

		// only take serial number from CSR if one was not supplied via API
		if ridSerialNumber == "" && csr != nil {
			ridSerialNumber = csr.Subject.SerialNumber
		}

		if csr != nil && data.role.UseCSRSANs {
			dnsNames = csr.DNSNames
			emailAddresses = csr.EmailAddresses
		}

		if cn != "" && !data.apiData.Get("exclude_cn_from_sans").(bool) {
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
					return nil, errutil.UserError{Err: err.Error()}
				}
				if hostnameRegex.MatchString(converted) {
					dnsNames = append(dnsNames, converted)
				}
			}
		}

		if csr == nil || !data.role.UseCSRSANs {
			cnAltRaw, ok := data.apiData.GetOk("alt_names")
			if ok {
				cnAlt := strutil.ParseDedupLowercaseAndSortStrings(cnAltRaw.(string), ",")
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
							return nil, errutil.UserError{Err: err.Error()}
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
			badName := validateNames(data, []string{cn})
			if len(badName) != 0 {
				return nil, errutil.UserError{Err: fmt.Sprintf(
					"common name %s not allowed by this role", badName)}
			}
		}

		if ridSerialNumber != "" {
			badName := validateSerialNumber(data, ridSerialNumber)
			if len(badName) != 0 {
				return nil, errutil.UserError{Err: fmt.Sprintf(
					"serial_number %s not allowed by this role", badName)}
			}
		}

		// Check for bad email and/or DNS names
		badName := validateNames(data, dnsNames)
		if len(badName) != 0 {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"subject alternate name %s not allowed by this role", badName)}
		}

		badName = validateNames(data, emailAddresses)
		if len(badName) != 0 {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"email address %s not allowed by this role", badName)}
		}
	}

	var otherSANs map[string][]string
	if sans := data.apiData.Get("other_sans").([]string); len(sans) > 0 {
		requested, err := parseOtherSANs(sans)
		if err != nil {
			return nil, errutil.UserError{Err: errwrap.Wrapf("could not parse requested other SAN: {{err}}", err).Error()}
		}
		badOID, badName, err := validateOtherSANs(data, requested)
		switch {
		case err != nil:
			return nil, errutil.UserError{Err: err.Error()}
		case len(badName) > 0:
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"other SAN %s not allowed for OID %s by this role", badName, badOID)}
		case len(badOID) > 0:
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"other SAN OID %s not allowed by this role", badOID)}
		default:
			otherSANs = requested
		}
	}

	// Get and verify any IP SANs
	ipAddresses := []net.IP{}
	{
		if csr != nil && data.role.UseCSRSANs {
			if len(csr.IPAddresses) > 0 {
				if !data.role.AllowIPSANs {
					return nil, errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided some via CSR")}
				}
				ipAddresses = csr.IPAddresses
			}
		} else {
			ipAlt := data.apiData.Get("ip_sans").([]string)
			if len(ipAlt) > 0 {
				if !data.role.AllowIPSANs {
					return nil, errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)}
				}
				for _, v := range ipAlt {
					parsedIP := net.ParseIP(v)
					if parsedIP == nil {
						return nil, errutil.UserError{Err: fmt.Sprintf(
							"the value '%s' is not a valid IP address", v)}
					}
					ipAddresses = append(ipAddresses, parsedIP)
				}
			}
		}
	}

	URIs := []*url.URL{}
	{
		if csr != nil && data.role.UseCSRSANs {
			if len(csr.URIs) > 0 {
				if len(data.role.AllowedURISANs) == 0 {
					return nil, errutil.UserError{Err: fmt.Sprintf(
						"URI Subject Alternative Names are not allowed in this role, but were provided via CSR"),
					}
				}

				// validate uri sans
				for _, uri := range csr.URIs {
					valid := false
					for _, allowed := range data.role.AllowedURISANs {
						validURI := glob.Glob(allowed, uri.String())
						if validURI {
							valid = true
							break
						}
					}

					if !valid {
						return nil, errutil.UserError{Err: fmt.Sprintf(
							"URI Subject Alternative Names were provided via CSR which are not valid for this role"),
						}
					}

					URIs = append(URIs, uri)
				}
			}
		} else {
			uriAlt := data.apiData.Get("uri_sans").([]string)
			if len(uriAlt) > 0 {
				if len(data.role.AllowedURISANs) == 0 {
					return nil, errutil.UserError{Err: fmt.Sprintf(
						"URI Subject Alternative Names are not allowed in this role, but were provided via the API"),
					}
				}

				for _, uri := range uriAlt {
					valid := false
					for _, allowed := range data.role.AllowedURISANs {
						validURI := glob.Glob(allowed, uri)
						if validURI {
							valid = true
							break
						}
					}

					if !valid {
						return nil, errutil.UserError{Err: fmt.Sprintf(
							"URI Subject Alternative Names were provided via CSR which are not valid for this role"),
						}
					}

					parsedURI, err := url.Parse(uri)
					if parsedURI == nil || err != nil {
						return nil, errutil.UserError{Err: fmt.Sprintf(
							"the provided URI Subject Alternative Name '%s' is not a valid URI", uri),
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
		Country:            strutil.RemoveDuplicates(data.role.Country, false),
		Organization:       strutil.RemoveDuplicates(data.role.Organization, false),
		OrganizationalUnit: strutil.RemoveDuplicatesStable(data.role.OU, false),
		Locality:           strutil.RemoveDuplicates(data.role.Locality, false),
		Province:           strutil.RemoveDuplicates(data.role.Province, false),
		StreetAddress:      strutil.RemoveDuplicates(data.role.StreetAddress, false),
		PostalCode:         strutil.RemoveDuplicates(data.role.PostalCode, false),
	}

	// Get the TTL and verify it against the max allowed
	var ttl time.Duration
	var maxTTL time.Duration
	var notAfter time.Time
	{
		ttl = time.Duration(data.apiData.Get("ttl").(int)) * time.Second

		if ttl == 0 && data.role.TTL > 0 {
			ttl = data.role.TTL
		}

		if data.role.MaxTTL > 0 {
			maxTTL = data.role.MaxTTL
		}

		if ttl == 0 {
			ttl = b.System().DefaultLeaseTTL()
		}
		if maxTTL == 0 {
			maxTTL = b.System().MaxLeaseTTL()
		}
		if ttl > maxTTL {
			ttl = maxTTL
		}

		notAfter = time.Now().Add(ttl)

		// If it's not self-signed, verify that the issued certificate won't be
		// valid past the lifetime of the CA certificate
		if caSign != nil &&
			notAfter.After(caSign.Certificate.NotAfter) && !data.role.AllowExpirationPastCA {

			return nil, errutil.UserError{Err: fmt.Sprintf(
				"cannot satisfy request, as TTL would result in notAfter %s that is beyond the expiration of the CA certificate at %s", notAfter.Format(time.RFC3339Nano), caSign.Certificate.NotAfter.Format(time.RFC3339Nano))}
		}
	}

	creation := &certutil.CreationBundle{
		Params: &certutil.CreationParameters{
			Subject:                       subject,
			DNSNames:                      dnsNames,
			EmailAddresses:                emailAddresses,
			IPAddresses:                   ipAddresses,
			URIs:                          URIs,
			OtherSANs:                     otherSANs,
			KeyType:                       data.role.KeyType,
			KeyBits:                       data.role.KeyBits,
			NotAfter:                      notAfter,
			KeyUsage:                      x509.KeyUsage(parseKeyUsages(data.role.KeyUsage)),
			ExtKeyUsage:                   parseExtKeyUsages(data.role),
			ExtKeyUsageOIDs:               data.role.ExtKeyUsageOIDs,
			PolicyIdentifiers:             data.role.PolicyIdentifiers,
			BasicConstraintsValidForNonCA: data.role.BasicConstraintsValidForNonCA,
			NotBeforeDuration:             data.role.NotBeforeDuration,
		},
		SigningBundle: caSign,
		CSR:           csr,
	}

	// Don't deal with URLs or max path length if it's self-signed, as these
	// normally come from the signing bundle
	if caSign == nil {
		return creation, nil
	}

	// This will have been read in from the getURLs function
	creation.Params.URLs = caSign.URLs

	// If the max path length in the role is not nil, it was specified at
	// generation time with the max_path_length parameter; otherwise derive it
	// from the signing certificate
	if data.role.MaxPathLength != nil {
		creation.Params.MaxPathLength = *data.role.MaxPathLength
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

	return creation, nil
}

func convertRespToPKCS8(resp *logical.Response) error {
	privRaw, ok := resp.Data["private_key"]
	if !ok {
		return nil
	}
	priv, ok := privRaw.(string)
	if !ok {
		return fmt.Errorf("error converting response to pkcs8: could not parse original value as string")
	}

	privKeyTypeRaw, ok := resp.Data["private_key_type"]
	if !ok {
		return fmt.Errorf("error converting response to pkcs8: %q not found in response", "private_key_type")
	}
	privKeyType, ok := privKeyTypeRaw.(certutil.PrivateKeyType)
	if !ok {
		return fmt.Errorf("error converting response to pkcs8: could not parse original type value as string")
	}

	var keyData []byte
	var pemUsed bool
	var err error
	var signer crypto.Signer

	block, _ := pem.Decode([]byte(priv))
	if block == nil {
		keyData, err = base64.StdEncoding.DecodeString(priv)
		if err != nil {
			return errwrap.Wrapf("error converting response to pkcs8: error decoding original value: {{err}}", err)
		}
	} else {
		keyData = block.Bytes
		pemUsed = true
	}

	switch privKeyType {
	case certutil.RSAPrivateKey:
		signer, err = x509.ParsePKCS1PrivateKey(keyData)
	case certutil.ECPrivateKey:
		signer, err = x509.ParseECPrivateKey(keyData)
	default:
		return fmt.Errorf("unknown private key type %q", privKeyType)
	}
	if err != nil {
		return errwrap.Wrapf("error converting response to pkcs8: error parsing previous key: {{err}}", err)
	}

	keyData, err = x509.MarshalPKCS8PrivateKey(signer)
	if err != nil {
		return errwrap.Wrapf("error converting response to pkcs8: error marshaling pkcs8 key: {{err}}", err)
	}

	if pemUsed {
		block.Type = "PRIVATE KEY"
		block.Bytes = keyData
		resp.Data["private_key"] = strings.TrimSpace(string(pem.EncodeToMemory(block)))
	} else {
		resp.Data["private_key"] = base64.StdEncoding.EncodeToString(keyData)
	}

	return nil
}
