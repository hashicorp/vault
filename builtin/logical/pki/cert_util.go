package pki

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	glob "github.com/ryanuber/go-glob"
	"golang.org/x/crypto/cryptobyte"
	cbbasn1 "golang.org/x/crypto/cryptobyte/asn1"
	"golang.org/x/net/idna"
)

type certExtKeyUsage int

const (
	anyExtKeyUsage certExtKeyUsage = 1 << iota
	serverAuthExtKeyUsage
	clientAuthExtKeyUsage
	codeSigningExtKeyUsage
	emailProtectionExtKeyUsage
	ipsecEndSystemExtKeyUsage
	ipsecTunnelExtKeyUsage
	ipsecUserExtKeyUsage
	timeStampingExtKeyUsage
	ocspSigningExtKeyUsage
	microsoftServerGatedCryptoExtKeyUsage
	netscapeServerGatedCryptoExtKeyUsage
	microsoftCommercialCodeSigningExtKeyUsage
	microsoftKernelCodeSigningExtKeyUsage
)

type dataBundle struct {
	params        *creationParameters
	signingBundle *caInfoBundle
	csr           *x509.CertificateRequest
	role          *roleEntry
	req           *logical.Request
	apiData       *framework.FieldData
}

type creationParameters struct {
	Subject                       pkix.Name
	DNSNames                      []string
	EmailAddresses                []string
	IPAddresses                   []net.IP
	URIs                          []*url.URL
	OtherSANs                     map[string][]string
	IsCA                          bool
	KeyType                       string
	KeyBits                       int
	NotAfter                      time.Time
	KeyUsage                      x509.KeyUsage
	ExtKeyUsage                   certExtKeyUsage
	ExtKeyUsageOIDs               []string
	PolicyIdentifiers             []string
	BasicConstraintsValidForNonCA bool

	// Only used when signing a CA cert
	UseCSRValues        bool
	PermittedDNSDomains []string

	// URLs to encode into the certificate
	URLs *urlEntries

	// The maximum path length to encode
	MaxPathLength int

	// The duration the certificate will use NotBefore
	NotBeforeDuration time.Duration
}

type caInfoBundle struct {
	certutil.ParsedCertBundle
	URLs *urlEntries
}

func (b *caInfoBundle) GetCAChain() []*certutil.CertBlock {
	chain := []*certutil.CertBlock{}

	// Include issuing CA in Chain, not including Root Authority
	if (len(b.Certificate.AuthorityKeyId) > 0 &&
		!bytes.Equal(b.Certificate.AuthorityKeyId, b.Certificate.SubjectKeyId)) ||
		(len(b.Certificate.AuthorityKeyId) == 0 &&
			!bytes.Equal(b.Certificate.RawIssuer, b.Certificate.RawSubject)) {

		chain = append(chain, &certutil.CertBlock{
			Certificate: b.Certificate,
			Bytes:       b.CertificateBytes,
		})
		if b.CAChain != nil && len(b.CAChain) > 0 {
			chain = append(chain, b.CAChain...)
		}
	}

	return chain
}

var (
	// A note on hostnameRegex: although we set the StrictDomainName option
	// when doing the idna conversion, this appears to only affect output, not
	// input, so it will allow e.g. host^123.example.com straight through. So
	// we still need to use this to check the output.
	hostnameRegex                = regexp.MustCompile(`^(\*\.)?(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	oidExtensionBasicConstraints = []int{2, 5, 29, 19}
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

func validateKeyTypeLength(keyType string, keyBits int) *logical.Response {
	switch keyType {
	case "rsa":
		switch keyBits {
		case 2048:
		case 4096:
		case 8192:
		default:
			return logical.ErrorResponse(fmt.Sprintf(
				"unsupported bit length for RSA key: %d", keyBits))
		}
	case "ec":
		switch keyBits {
		case 224:
		case 256:
		case 384:
		case 521:
		default:
			return logical.ErrorResponse(fmt.Sprintf(
				"unsupported bit length for EC key: %d", keyBits))
		}
	case "any":
	default:
		return logical.ErrorResponse(fmt.Sprintf(
			"unknown key type %s", keyType))
	}

	return nil
}

// Fetches the CA info. Unlike other certificates, the CA info is stored
// in the backend as a CertBundle, because we are storing its private key
func fetchCAInfo(ctx context.Context, req *logical.Request) (*caInfoBundle, error) {
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

	caInfo := &caInfoBundle{*parsedBundle, nil}

	entries, err := getURLs(ctx, req)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch URL information: %v", err)}
	}
	if entries == nil {
		entries = &urlEntries{
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
func validateNames(data *dataBundle, names []string) string {
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
func validateOtherSANs(data *dataBundle, requested map[string][]string) (string, string, error) {
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

func validateSerialNumber(data *dataBundle, serialNumber string) string {
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
	data *dataBundle,
	isCA bool) (*certutil.ParsedCertBundle, error) {

	if data.role == nil {
		return nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	if data.role.KeyType == "rsa" && data.role.KeyBits < 2048 {
		return nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
	}

	err := generateCreationBundle(b, data)
	if err != nil {
		return nil, err
	}
	if data.params == nil {
		return nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	if isCA {
		data.params.IsCA = isCA
		data.params.PermittedDNSDomains = data.apiData.Get("permitted_dns_domains").([]string)

		if data.signingBundle == nil {
			// Generating a self-signed root certificate
			entries, err := getURLs(ctx, data.req)
			if err != nil {
				return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch URL information: %v", err)}
			}
			if entries == nil {
				entries = &urlEntries{
					IssuingCertificates:   []string{},
					CRLDistributionPoints: []string{},
					OCSPServers:           []string{},
				}
			}
			data.params.URLs = entries

			if data.role.MaxPathLength == nil {
				data.params.MaxPathLength = -1
			} else {
				data.params.MaxPathLength = *data.role.MaxPathLength
			}
		}
	}

	parsedBundle, err := createCertificate(data)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

// N.B.: This is only meant to be used for generating intermediate CAs.
// It skips some sanity checks.
func generateIntermediateCSR(b *backend, data *dataBundle) (*certutil.ParsedCSRBundle, error) {
	err := generateCreationBundle(b, data)
	if err != nil {
		return nil, err
	}
	if data.params == nil {
		return nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	parsedBundle, err := createCSR(data)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

func signCert(b *backend,
	data *dataBundle,
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

	data.csr = csr

	err = generateCreationBundle(b, data)
	if err != nil {
		return nil, err
	}
	if data.params == nil {
		return nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	data.params.IsCA = isCA
	data.params.UseCSRValues = useCSRValues

	if isCA {
		data.params.PermittedDNSDomains = data.apiData.Get("permitted_dns_domains").([]string)
	}

	parsedBundle, err := signCertificate(data)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

// generateCreationBundle is a shared function that reads parameters supplied
// from the various endpoints and generates a creationParameters with the
// parameters that can be used to issue or sign
func generateCreationBundle(b *backend, data *dataBundle) error {
	// Read in names -- CN, DNS and email addresses
	var cn string
	var ridSerialNumber string
	dnsNames := []string{}
	emailAddresses := []string{}
	{
		if data.csr != nil && data.role.UseCSRCommonName {
			cn = data.csr.Subject.CommonName
		}
		if cn == "" {
			cn = data.apiData.Get("common_name").(string)
			if cn == "" && data.role.RequireCN {
				return errutil.UserError{Err: `the common_name field is required, or must be provided in a CSR with "use_csr_common_name" set to true, unless "require_cn" is set to false`}
			}
		}

		ridSerialNumber = data.apiData.Get("serial_number").(string)

		// only take serial number from CSR if one was not supplied via API
		if ridSerialNumber == "" && data.csr != nil {
			ridSerialNumber = data.csr.Subject.SerialNumber
		}

		if data.csr != nil && data.role.UseCSRSANs {
			dnsNames = data.csr.DNSNames
			emailAddresses = data.csr.EmailAddresses
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
					return errutil.UserError{Err: err.Error()}
				}
				if hostnameRegex.MatchString(converted) {
					dnsNames = append(dnsNames, converted)
				}
			}
		}

		if data.csr == nil || !data.role.UseCSRSANs {
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
							return errutil.UserError{Err: err.Error()}
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
				return errutil.UserError{Err: fmt.Sprintf(
					"common name %s not allowed by this role", badName)}
			}
		}

		if ridSerialNumber != "" {
			badName := validateSerialNumber(data, ridSerialNumber)
			if len(badName) != 0 {
				return errutil.UserError{Err: fmt.Sprintf(
					"serial_number %s not allowed by this role", badName)}
			}
		}

		// Check for bad email and/or DNS names
		badName := validateNames(data, dnsNames)
		if len(badName) != 0 {
			return errutil.UserError{Err: fmt.Sprintf(
				"subject alternate name %s not allowed by this role", badName)}
		}

		badName = validateNames(data, emailAddresses)
		if len(badName) != 0 {
			return errutil.UserError{Err: fmt.Sprintf(
				"email address %s not allowed by this role", badName)}
		}
	}

	var otherSANs map[string][]string
	if sans := data.apiData.Get("other_sans").([]string); len(sans) > 0 {
		requested, err := parseOtherSANs(sans)
		if err != nil {
			return errutil.UserError{Err: errwrap.Wrapf("could not parse requested other SAN: {{err}}", err).Error()}
		}
		badOID, badName, err := validateOtherSANs(data, requested)
		switch {
		case err != nil:
			return errutil.UserError{Err: err.Error()}
		case len(badName) > 0:
			return errutil.UserError{Err: fmt.Sprintf(
				"other SAN %s not allowed for OID %s by this role", badName, badOID)}
		case len(badOID) > 0:
			return errutil.UserError{Err: fmt.Sprintf(
				"other SAN OID %s not allowed by this role", badOID)}
		default:
			otherSANs = requested
		}
	}

	// Get and verify any IP SANs
	ipAddresses := []net.IP{}
	{
		if data.csr != nil && data.role.UseCSRSANs {
			if len(data.csr.IPAddresses) > 0 {
				if !data.role.AllowIPSANs {
					return errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided some via CSR")}
				}
				ipAddresses = data.csr.IPAddresses
			}
		} else {
			ipAlt := data.apiData.Get("ip_sans").([]string)
			if len(ipAlt) > 0 {
				if !data.role.AllowIPSANs {
					return errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)}
				}
				for _, v := range ipAlt {
					parsedIP := net.ParseIP(v)
					if parsedIP == nil {
						return errutil.UserError{Err: fmt.Sprintf(
							"the value '%s' is not a valid IP address", v)}
					}
					ipAddresses = append(ipAddresses, parsedIP)
				}
			}
		}
	}

	URIs := []*url.URL{}
	{
		if data.csr != nil && data.role.UseCSRSANs {
			if len(data.csr.URIs) > 0 {
				if len(data.role.AllowedURISANs) == 0 {
					return errutil.UserError{Err: fmt.Sprintf(
						"URI Subject Alternative Names are not allowed in this role, but were provided via CSR"),
					}
				}

				// validate uri sans
				for _, uri := range data.csr.URIs {
					valid := false
					for _, allowed := range data.role.AllowedURISANs {
						validURI := glob.Glob(allowed, uri.String())
						if validURI {
							valid = true
							break
						}
					}

					if !valid {
						return errutil.UserError{Err: fmt.Sprintf(
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
					return errutil.UserError{Err: fmt.Sprintf(
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
						return errutil.UserError{Err: fmt.Sprintf(
							"URI Subject Alternative Names were provided via CSR which are not valid for this role"),
						}
					}

					parsedURI, err := url.Parse(uri)
					if parsedURI == nil || err != nil {
						return errutil.UserError{Err: fmt.Sprintf(
							"the provided URI Subject Alternative Name '%s' is not a valid URI", uri),
						}
					}

					URIs = append(URIs, parsedURI)
				}
			}
		}
	}

	subject := pkix.Name{
		CommonName:         cn,
		SerialNumber:       ridSerialNumber,
		Country:            strutil.RemoveDuplicates(data.role.Country, false),
		Organization:       strutil.RemoveDuplicates(data.role.Organization, false),
		OrganizationalUnit: strutil.RemoveDuplicates(data.role.OU, false),
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
		if data.signingBundle != nil &&
			notAfter.After(data.signingBundle.Certificate.NotAfter) && !data.role.AllowExpirationPastCA {

			return errutil.UserError{Err: fmt.Sprintf(
				"cannot satisfy request, as TTL would result in notAfter %s that is beyond the expiration of the CA certificate at %s", notAfter.Format(time.RFC3339Nano), data.signingBundle.Certificate.NotAfter.Format(time.RFC3339Nano))}
		}
	}

	data.params = &creationParameters{
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
	}

	// Don't deal with URLs or max path length if it's self-signed, as these
	// normally come from the signing bundle
	if data.signingBundle == nil {
		return nil
	}

	// This will have been read in from the getURLs function
	data.params.URLs = data.signingBundle.URLs

	// If the max path length in the role is not nil, it was specified at
	// generation time with the max_path_length parameter; otherwise derive it
	// from the signing certificate
	if data.role.MaxPathLength != nil {
		data.params.MaxPathLength = *data.role.MaxPathLength
	} else {
		switch {
		case data.signingBundle.Certificate.MaxPathLen < 0:
			data.params.MaxPathLength = -1
		case data.signingBundle.Certificate.MaxPathLen == 0 &&
			data.signingBundle.Certificate.MaxPathLenZero:
			// The signing function will ensure that we do not issue a CA cert
			data.params.MaxPathLength = 0
		default:
			// If this takes it to zero, we handle this case later if
			// necessary
			data.params.MaxPathLength = data.signingBundle.Certificate.MaxPathLen - 1
		}
	}

	return nil
}

// addKeyUsages adds appropriate key usages to the template given the creation
// information
func addKeyUsages(data *dataBundle, certTemplate *x509.Certificate) {
	if data.params.IsCA {
		certTemplate.KeyUsage = x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign)
		return
	}

	certTemplate.KeyUsage = data.params.KeyUsage

	if data.params.ExtKeyUsage&anyExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageAny)
	}

	if data.params.ExtKeyUsage&serverAuthExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}

	if data.params.ExtKeyUsage&clientAuthExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}

	if data.params.ExtKeyUsage&codeSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageCodeSigning)
	}

	if data.params.ExtKeyUsage&emailProtectionExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageEmailProtection)
	}

	if data.params.ExtKeyUsage&ipsecEndSystemExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageIPSECEndSystem)
	}

	if data.params.ExtKeyUsage&ipsecTunnelExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageIPSECTunnel)
	}

	if data.params.ExtKeyUsage&ipsecUserExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageIPSECUser)
	}

	if data.params.ExtKeyUsage&timeStampingExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageTimeStamping)
	}

	if data.params.ExtKeyUsage&ocspSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageOCSPSigning)
	}

	if data.params.ExtKeyUsage&microsoftServerGatedCryptoExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageMicrosoftServerGatedCrypto)
	}

	if data.params.ExtKeyUsage&netscapeServerGatedCryptoExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageNetscapeServerGatedCrypto)
	}

	if data.params.ExtKeyUsage&microsoftCommercialCodeSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageMicrosoftCommercialCodeSigning)
	}

	if data.params.ExtKeyUsage&microsoftKernelCodeSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageMicrosoftKernelCodeSigning)
	}
}

// addPolicyIdentifiers adds certificate policies extension
//
func addPolicyIdentifiers(data *dataBundle, certTemplate *x509.Certificate) {
	for _, oidstr := range data.params.PolicyIdentifiers {
		oid, err := stringToOid(oidstr)
		if err == nil {
			certTemplate.PolicyIdentifiers = append(certTemplate.PolicyIdentifiers, oid)
		}
	}
}

// addExtKeyUsageOids adds custom extended key usage OIDs to certificate
func addExtKeyUsageOids(data *dataBundle, certTemplate *x509.Certificate) {
	for _, oidstr := range data.params.ExtKeyUsageOIDs {
		oid, err := stringToOid(oidstr)
		if err == nil {
			certTemplate.UnknownExtKeyUsage = append(certTemplate.UnknownExtKeyUsage, oid)
		}
	}
}

// Performs the heavy lifting of creating a certificate. Returns
// a fully-filled-in ParsedCertBundle.
func createCertificate(data *dataBundle) (*certutil.ParsedCertBundle, error) {
	var err error
	result := &certutil.ParsedCertBundle{}

	serialNumber, err := certutil.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	if err := certutil.GeneratePrivateKey(data.params.KeyType,
		data.params.KeyBits,
		result); err != nil {
		return nil, err
	}

	subjKeyID, err := certutil.GetSubjKeyID(result.PrivateKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error getting subject key ID: %s", err)}
	}

	certTemplate := &x509.Certificate{
		SerialNumber:   serialNumber,
		NotBefore:      time.Now().Add(-30 * time.Second),
		NotAfter:       data.params.NotAfter,
		IsCA:           false,
		SubjectKeyId:   subjKeyID,
		Subject:        data.params.Subject,
		DNSNames:       data.params.DNSNames,
		EmailAddresses: data.params.EmailAddresses,
		IPAddresses:    data.params.IPAddresses,
		URIs:           data.params.URIs,
	}
	if data.params.NotBeforeDuration > 0 {
		certTemplate.NotBefore = time.Now().Add(-1 * data.params.NotBeforeDuration)
	}

	if err := handleOtherSANs(certTemplate, data.params.OtherSANs); err != nil {
		return nil, errutil.InternalError{Err: errwrap.Wrapf("error marshaling other SANs: {{err}}", err).Error()}
	}

	// Add this before calling addKeyUsages
	if data.signingBundle == nil {
		certTemplate.IsCA = true
	} else if data.params.BasicConstraintsValidForNonCA {
		certTemplate.BasicConstraintsValid = true
		certTemplate.IsCA = false
	}

	// This will only be filled in from the generation paths
	if len(data.params.PermittedDNSDomains) > 0 {
		certTemplate.PermittedDNSDomains = data.params.PermittedDNSDomains
		certTemplate.PermittedDNSDomainsCritical = true
	}

	addPolicyIdentifiers(data, certTemplate)

	addKeyUsages(data, certTemplate)

	addExtKeyUsageOids(data, certTemplate)

	certTemplate.IssuingCertificateURL = data.params.URLs.IssuingCertificates
	certTemplate.CRLDistributionPoints = data.params.URLs.CRLDistributionPoints
	certTemplate.OCSPServer = data.params.URLs.OCSPServers

	var certBytes []byte
	if data.signingBundle != nil {
		switch data.signingBundle.PrivateKeyType {
		case certutil.RSAPrivateKey:
			certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
		case certutil.ECPrivateKey:
			certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
		}

		caCert := data.signingBundle.Certificate
		certTemplate.AuthorityKeyId = caCert.SubjectKeyId

		certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, result.PrivateKey.Public(), data.signingBundle.PrivateKey)
	} else {
		// Creating a self-signed root
		if data.params.MaxPathLength == 0 {
			certTemplate.MaxPathLen = 0
			certTemplate.MaxPathLenZero = true
		} else {
			certTemplate.MaxPathLen = data.params.MaxPathLength
		}

		switch data.params.KeyType {
		case "rsa":
			certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
		case "ec":
			certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
		}

		certTemplate.AuthorityKeyId = subjKeyID
		certTemplate.BasicConstraintsValid = true
		certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, result.PrivateKey.Public(), result.PrivateKey)
	}

	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CertificateBytes = certBytes
	result.Certificate, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	if data.signingBundle != nil {
		if len(data.signingBundle.Certificate.AuthorityKeyId) > 0 &&
			!bytes.Equal(data.signingBundle.Certificate.AuthorityKeyId, data.signingBundle.Certificate.SubjectKeyId) {

			result.CAChain = []*certutil.CertBlock{
				&certutil.CertBlock{
					Certificate: data.signingBundle.Certificate,
					Bytes:       data.signingBundle.CertificateBytes,
				},
			}
			result.CAChain = append(result.CAChain, data.signingBundle.CAChain...)
		}
	}

	return result, nil
}

// Creates a CSR. This is currently only meant for use when
// generating an intermediate certificate.
func createCSR(data *dataBundle) (*certutil.ParsedCSRBundle, error) {
	var err error
	result := &certutil.ParsedCSRBundle{}

	if err := certutil.GeneratePrivateKey(data.params.KeyType,
		data.params.KeyBits,
		result); err != nil {
		return nil, err
	}

	// Like many root CAs, other information is ignored
	csrTemplate := &x509.CertificateRequest{
		Subject:        data.params.Subject,
		DNSNames:       data.params.DNSNames,
		EmailAddresses: data.params.EmailAddresses,
		IPAddresses:    data.params.IPAddresses,
		URIs:           data.params.URIs,
	}

	if err := handleOtherCSRSANs(csrTemplate, data.params.OtherSANs); err != nil {
		return nil, errutil.InternalError{Err: errwrap.Wrapf("error marshaling other SANs: {{err}}", err).Error()}
	}

	if data.apiData != nil && data.apiData.Get("add_basic_constraints").(bool) {
		type basicConstraints struct {
			IsCA       bool `asn1:"optional"`
			MaxPathLen int  `asn1:"optional,default:-1"`
		}
		val, err := asn1.Marshal(basicConstraints{IsCA: true, MaxPathLen: -1})
		if err != nil {
			return nil, errutil.InternalError{Err: errwrap.Wrapf("error marshaling basic constraints: {{err}}", err).Error()}
		}
		ext := pkix.Extension{
			Id:       oidExtensionBasicConstraints,
			Value:    val,
			Critical: true,
		}
		csrTemplate.ExtraExtensions = append(csrTemplate.ExtraExtensions, ext)
	}

	switch data.params.KeyType {
	case "rsa":
		csrTemplate.SignatureAlgorithm = x509.SHA256WithRSA
	case "ec":
		csrTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, result.PrivateKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CSRBytes = csr
	result.CSR, err = x509.ParseCertificateRequest(csr)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %v", err)}
	}

	return result, nil
}

// Performs the heavy lifting of generating a certificate from a CSR.
// Returns a ParsedCertBundle sans private keys.
func signCertificate(data *dataBundle) (*certutil.ParsedCertBundle, error) {
	switch {
	case data == nil:
		return nil, errutil.UserError{Err: "nil data bundle given to signCertificate"}
	case data.params == nil:
		return nil, errutil.UserError{Err: "nil parameters given to signCertificate"}
	case data.signingBundle == nil:
		return nil, errutil.UserError{Err: "nil signing bundle given to signCertificate"}
	case data.csr == nil:
		return nil, errutil.UserError{Err: "nil csr given to signCertificate"}
	}

	err := data.csr.CheckSignature()
	if err != nil {
		return nil, errutil.UserError{Err: "request signature invalid"}
	}

	result := &certutil.ParsedCertBundle{}

	serialNumber, err := certutil.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(data.csr.PublicKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error marshalling public key: %s", err)}
	}
	subjKeyID := sha1.Sum(marshaledKey)

	caCert := data.signingBundle.Certificate

	certTemplate := &x509.Certificate{
		SerialNumber:   serialNumber,
		Subject:        data.params.Subject,
		NotBefore:      time.Now().Add(-30 * time.Second),
		NotAfter:       data.params.NotAfter,
		SubjectKeyId:   subjKeyID[:],
		AuthorityKeyId: caCert.SubjectKeyId,
	}
	if data.params.NotBeforeDuration > 0 {
		certTemplate.NotBefore = time.Now().Add(-1 * data.params.NotBeforeDuration)
	}

	switch data.signingBundle.PrivateKeyType {
	case certutil.RSAPrivateKey:
		certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
	case certutil.ECPrivateKey:
		certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
	}

	if data.params.UseCSRValues {
		certTemplate.Subject = data.csr.Subject
		certTemplate.Subject.ExtraNames = certTemplate.Subject.Names

		certTemplate.DNSNames = data.csr.DNSNames
		certTemplate.EmailAddresses = data.csr.EmailAddresses
		certTemplate.IPAddresses = data.csr.IPAddresses
		certTemplate.URIs = data.csr.URIs

		for _, name := range data.csr.Extensions {
			if !name.Id.Equal(oidExtensionBasicConstraints) {
				certTemplate.ExtraExtensions = append(certTemplate.ExtraExtensions, name)
			}
		}

	} else {
		certTemplate.DNSNames = data.params.DNSNames
		certTemplate.EmailAddresses = data.params.EmailAddresses
		certTemplate.IPAddresses = data.params.IPAddresses
		certTemplate.URIs = data.params.URIs
	}

	if err := handleOtherSANs(certTemplate, data.params.OtherSANs); err != nil {
		return nil, errutil.InternalError{Err: errwrap.Wrapf("error marshaling other SANs: {{err}}", err).Error()}
	}

	addPolicyIdentifiers(data, certTemplate)

	addKeyUsages(data, certTemplate)

	addExtKeyUsageOids(data, certTemplate)

	var certBytes []byte

	certTemplate.IssuingCertificateURL = data.params.URLs.IssuingCertificates
	certTemplate.CRLDistributionPoints = data.params.URLs.CRLDistributionPoints
	certTemplate.OCSPServer = data.signingBundle.URLs.OCSPServers

	if data.params.IsCA {
		certTemplate.BasicConstraintsValid = true
		certTemplate.IsCA = true

		if data.signingBundle.Certificate.MaxPathLen == 0 &&
			data.signingBundle.Certificate.MaxPathLenZero {
			return nil, errutil.UserError{Err: "signing certificate has a max path length of zero, and cannot issue further CA certificates"}
		}

		certTemplate.MaxPathLen = data.params.MaxPathLength
		if certTemplate.MaxPathLen == 0 {
			certTemplate.MaxPathLenZero = true
		}
	} else if data.params.BasicConstraintsValidForNonCA {
		certTemplate.BasicConstraintsValid = true
		certTemplate.IsCA = false
	}

	if len(data.params.PermittedDNSDomains) > 0 {
		certTemplate.PermittedDNSDomains = data.params.PermittedDNSDomains
		certTemplate.PermittedDNSDomainsCritical = true
	}

	certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, data.csr.PublicKey, data.signingBundle.PrivateKey)

	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CertificateBytes = certBytes
	result.Certificate, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	result.CAChain = data.signingBundle.GetCAChain()

	return result, nil
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

func handleOtherCSRSANs(in *x509.CertificateRequest, sans map[string][]string) error {
	certTemplate := &x509.Certificate{
		DNSNames:       in.DNSNames,
		IPAddresses:    in.IPAddresses,
		EmailAddresses: in.EmailAddresses,
		URIs:           in.URIs,
	}
	if err := handleOtherSANs(certTemplate, sans); err != nil {
		return err
	}
	if len(certTemplate.ExtraExtensions) > 0 {
		for _, v := range certTemplate.ExtraExtensions {
			in.ExtraExtensions = append(in.ExtraExtensions, v)
		}
	}
	return nil
}

func handleOtherSANs(in *x509.Certificate, sans map[string][]string) error {
	// If other SANs is empty we return which causes normal Go stdlib parsing
	// of the other SAN types
	if len(sans) == 0 {
		return nil
	}

	var rawValues []asn1.RawValue

	// We need to generate an IMPLICIT sequence for compatibility with OpenSSL
	// -- it's an open question what the default for RFC 5280 actually is, see
	// https://github.com/openssl/openssl/issues/5091 -- so we have to use
	// cryptobyte because using the asn1 package's marshaling always produces
	// an EXPLICIT sequence. Note that asn1 is way too magical according to
	// agl, and cryptobyte is modeled after the CBB/CBS bits that agl put into
	// boringssl.
	for oid, vals := range sans {
		for _, val := range vals {
			var b cryptobyte.Builder
			oidStr, err := stringToOid(oid)
			if err != nil {
				return err
			}
			b.AddASN1ObjectIdentifier(oidStr)
			b.AddASN1(cbbasn1.Tag(0).ContextSpecific().Constructed(), func(b *cryptobyte.Builder) {
				b.AddASN1(cbbasn1.UTF8String, func(b *cryptobyte.Builder) {
					b.AddBytes([]byte(val))
				})
			})
			m, err := b.Bytes()
			if err != nil {
				return err
			}
			rawValues = append(rawValues, asn1.RawValue{Tag: 0, Class: 2, IsCompound: true, Bytes: m})
		}
	}

	// If other SANs is empty we return which causes normal Go stdlib parsing
	// of the other SAN types
	if len(rawValues) == 0 {
		return nil
	}

	// Append any existing SANs, sans marshalling
	rawValues = append(rawValues, marshalSANs(in.DNSNames, in.EmailAddresses, in.IPAddresses, in.URIs)...)

	// Marshal and add to ExtraExtensions
	ext := pkix.Extension{
		// This is the defined OID for subjectAltName
		Id: asn1.ObjectIdentifier{2, 5, 29, 17},
	}
	var err error
	ext.Value, err = asn1.Marshal(rawValues)
	if err != nil {
		return err
	}
	in.ExtraExtensions = append(in.ExtraExtensions, ext)

	return nil
}

// Note: Taken from the Go source code since it's not public, and used in the
// modified function below (which also uses these consts upstream)
const (
	nameTypeEmail = 1
	nameTypeDNS   = 2
	nameTypeURI   = 6
	nameTypeIP    = 7
)

// Note: Taken from the Go source code since it's not public, plus changed to not marshal
// marshalSANs marshals a list of addresses into a the contents of an X.509
// SubjectAlternativeName extension.
func marshalSANs(dnsNames, emailAddresses []string, ipAddresses []net.IP, uris []*url.URL) []asn1.RawValue {
	var rawValues []asn1.RawValue
	for _, name := range dnsNames {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeDNS, Class: 2, Bytes: []byte(name)})
	}
	for _, email := range emailAddresses {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeEmail, Class: 2, Bytes: []byte(email)})
	}
	for _, rawIP := range ipAddresses {
		// If possible, we always want to encode IPv4 addresses in 4 bytes.
		ip := rawIP.To4()
		if ip == nil {
			ip = rawIP
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeIP, Class: 2, Bytes: ip})
	}
	for _, uri := range uris {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeURI, Class: 2, Bytes: []byte(uri.String())})
	}
	return rawValues
}

func stringToOid(in string) (asn1.ObjectIdentifier, error) {
	split := strings.Split(in, ".")
	ret := make(asn1.ObjectIdentifier, 0, len(split))
	for _, v := range split {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, i)
	}
	return asn1.ObjectIdentifier(ret), nil
}
