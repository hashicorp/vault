package pki

import (
	"bytes"
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
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/ryanuber/go-glob"
)

type certExtKeyUsage int

const (
	serverExtKeyUsage certExtKeyUsage = 1 << iota
	clientExtKeyUsage
	codeSigningExtKeyUsage
	emailProtectionExtKeyUsage
)

type creationBundle struct {
	CommonName     string
	OU             []string
	Organization   []string
	DNSNames       []string
	EmailAddresses []string
	IPAddresses    []net.IP
	IsCA           bool
	KeyType        string
	KeyBits        int
	SigningBundle  *caInfoBundle
	NotAfter       time.Time
	KeyUsage       x509.KeyUsage
	ExtKeyUsage    certExtKeyUsage

	// Only used when signing a CA cert
	UseCSRValues        bool
	PermittedDNSDomains []string

	// URLs to encode into the certificate
	URLs *urlEntries

	// The maximum path length to encode
	MaxPathLength int
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
	hostnameRegex                = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
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
	default:
		return logical.ErrorResponse(fmt.Sprintf(
			"unknown key type %s", keyType))
	}

	return nil
}

// Fetches the CA info. Unlike other certificates, the CA info is stored
// in the backend as a CertBundle, because we are storing its private key
func fetchCAInfo(req *logical.Request) (*caInfoBundle, error) {
	bundleEntry, err := req.Storage.Get("config/ca_bundle")
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

	entries, err := getURLs(req)
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
func fetchCertBySerial(req *logical.Request, prefix, serial string) (*logical.StorageEntry, error) {
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

	certEntry, err = req.Storage.Get(path)
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

	// Retrieve the old-style path
	certEntry, err = req.Storage.Get(legacyPath)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error fetching certificate %s: %s", serial, err)}
	}
	if certEntry == nil {
		return nil, nil
	}
	if certEntry.Value == nil || len(certEntry.Value) == 0 {
		return nil, errutil.InternalError{Err: fmt.Sprintf("returned certificate bytes for serial %s were empty", serial)}
	}

	// Update old-style paths to new-style paths
	certEntry.Key = path
	if err = req.Storage.Put(certEntry); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error saving certificate with serial %s to new location", serial)}
	}
	if err = req.Storage.Delete(legacyPath); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error deleting certificate with serial %s from old location", serial)}
	}

	return certEntry, nil
}

// Given a set of requested names for a certificate, verifies that all of them
// match the various toggles set in the role for controlling issuance.
// If one does not pass, it is returned in the string argument.
func validateNames(req *logical.Request, names []string, role *roleEntry) string {
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
		if role.EnforceHostnames {
			if !hostnameRegex.MatchString(sanitizedName) {
				return name
			}
		}

		// Self-explanatory
		if role.AllowAnyName {
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

		if role.AllowLocalhost {
			if name == "localhost" ||
				name == "localdomain" ||
				(isEmail && emailDomain == "localhost") ||
				(isEmail && emailDomain == "localdomain") {
				continue
			}

			if role.AllowSubdomains {
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

		if role.AllowTokenDisplayName {
			if name == req.DisplayName {
				continue
			}

			if role.AllowSubdomains {
				if isEmail {
					// If it's an email address, we need to parse the token
					// display name in order to do a proper comparison of the
					// subdomain
					if strings.Contains(req.DisplayName, "@") {
						splitDisplay := strings.Split(req.DisplayName, "@")
						if len(splitDisplay) == 2 {
							// Compare the sanitized name against the hostname
							// portion of the email address in the roken
							// display name
							if strings.HasSuffix(sanitizedName, "."+splitDisplay[1]) {
								continue
							}
						}
					}
				}

				if strings.HasSuffix(sanitizedName, "."+req.DisplayName) ||
					(isWildcard && sanitizedName == req.DisplayName) {
					continue
				}
			}
		}

		if role.AllowedDomains != "" {
			valid := false
			for _, currDomain := range strings.Split(role.AllowedDomains, ",") {
				// If there is, say, a trailing comma, ignore it
				if currDomain == "" {
					continue
				}

				// First, allow an exact match of the base domain if that role flag
				// is enabled
				if role.AllowBareDomains &&
					(name == currDomain ||
						(isEmail && emailDomain == currDomain)) {
					valid = true
					break
				}

				if role.AllowSubdomains {
					if strings.HasSuffix(sanitizedName, "."+currDomain) ||
						(isWildcard && sanitizedName == currDomain) {
						valid = true
						break
					}
				}

				if role.AllowGlobDomains &&
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

		//panic(fmt.Sprintf("\nName is %s\nRole is\n%#v\n", name, role))
		return name
	}

	return ""
}

func generateCert(b *backend,
	role *roleEntry,
	signingBundle *caInfoBundle,
	isCA bool,
	req *logical.Request,
	data *framework.FieldData) (*certutil.ParsedCertBundle, error) {

	if role.KeyType == "rsa" && role.KeyBits < 2048 {
		return nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
	}

	creationBundle, err := generateCreationBundle(b, role, signingBundle, nil, req, data)
	if err != nil {
		return nil, err
	}

	if isCA {
		creationBundle.IsCA = isCA

		creationBundle.PermittedDNSDomains = data.Get("permitted_dns_domains").([]string)

		if signingBundle == nil {
			// Generating a self-signed root certificate
			entries, err := getURLs(req)
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
			creationBundle.URLs = entries

			if role.MaxPathLength == nil {
				creationBundle.MaxPathLength = -1
			} else {
				creationBundle.MaxPathLength = *role.MaxPathLength
			}
		}
	}

	parsedBundle, err := createCertificate(creationBundle)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

// N.B.: This is only meant to be used for generating intermediate CAs.
// It skips some sanity checks.
func generateIntermediateCSR(b *backend,
	role *roleEntry,
	signingBundle *caInfoBundle,
	req *logical.Request,
	data *framework.FieldData) (*certutil.ParsedCSRBundle, error) {

	creationBundle, err := generateCreationBundle(b, role, signingBundle, nil, req, data)
	if err != nil {
		return nil, err
	}

	parsedBundle, err := createCSR(creationBundle)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

func signCert(b *backend,
	role *roleEntry,
	signingBundle *caInfoBundle,
	isCA bool,
	useCSRValues bool,
	req *logical.Request,
	data *framework.FieldData) (*certutil.ParsedCertBundle, error) {

	csrString := data.Get("csr").(string)
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

	switch role.KeyType {
	case "rsa":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.RSA {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				role.KeyType)}
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
		if pubKey.N.BitLen() < role.KeyBits {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				role.KeyBits,
				pubKey.N.BitLen())}
		}

	case "ec":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.ECDSA {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				role.KeyType)}
		}
		pubKey, ok := csr.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		// Verify that the bit size is at least the size specified in the role
		if pubKey.Params().BitSize < role.KeyBits {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				role.KeyBits,
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

	creationBundle, err := generateCreationBundle(b, role, signingBundle, csr, req, data)
	if err != nil {
		return nil, err
	}

	creationBundle.IsCA = isCA
	creationBundle.UseCSRValues = useCSRValues

	if isCA {
		creationBundle.PermittedDNSDomains = data.Get("permitted_dns_domains").([]string)
	}

	parsedBundle, err := signCertificate(creationBundle, csr)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

// generateCreationBundle is a shared function that reads parameters supplied
// from the various endpoints and generates a creationBundle with the
// parameters that can be used to issue or sign
func generateCreationBundle(b *backend,
	role *roleEntry,
	signingBundle *caInfoBundle,
	csr *x509.CertificateRequest,
	req *logical.Request,
	data *framework.FieldData) (*creationBundle, error) {
	var err error
	var ok bool

	// Read in names -- CN, DNS and email addresses
	var cn string
	dnsNames := []string{}
	emailAddresses := []string{}
	{
		if csr != nil && role.UseCSRCommonName {
			cn = csr.Subject.CommonName
		}
		if cn == "" {
			cn = data.Get("common_name").(string)
			if cn == "" {
				return nil, errutil.UserError{Err: `the common_name field is required, or must be provided in a CSR with "use_csr_common_name" set to true`}
			}
		}

		if csr != nil && role.UseCSRSANs {
			dnsNames = csr.DNSNames
			emailAddresses = csr.EmailAddresses
		}

		if !data.Get("exclude_cn_from_sans").(bool) {
			if strings.Contains(cn, "@") {
				// Note: emails are not disallowed if the role's email protection
				// flag is false, because they may well be included for
				// informational purposes; it is up to the verifying party to
				// ensure that email addresses in a subject alternate name can be
				// used for the purpose for which they are presented
				emailAddresses = append(emailAddresses, cn)
			} else {
				dnsNames = append(dnsNames, cn)
			}
		}

		if csr == nil || !role.UseCSRSANs {
			cnAltRaw, ok := data.GetOk("alt_names")
			if ok {
				cnAlt := strutil.ParseDedupLowercaseAndSortStrings(cnAltRaw.(string), ",")
				for _, v := range cnAlt {
					if strings.Contains(v, "@") {
						emailAddresses = append(emailAddresses, v)
					} else {
						dnsNames = append(dnsNames, v)
					}
				}
			}
		}

		// Check the CN. This ensures that the CN is checked even if it's
		// excluded from SANs.
		badName := validateNames(req, []string{cn}, role)
		if len(badName) != 0 {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"common name %s not allowed by this role", badName)}
		}

		// Check for bad email and/or DNS names
		badName = validateNames(req, dnsNames, role)
		if len(badName) != 0 {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"subject alternate name %s not allowed by this role", badName)}
		}

		badName = validateNames(req, emailAddresses, role)
		if len(badName) != 0 {
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"email address %s not allowed by this role", badName)}
		}
	}

	// Get and verify any IP SANs
	ipAddresses := []net.IP{}
	var ipAltInt interface{}
	{
		if csr != nil && role.UseCSRSANs {
			if len(csr.IPAddresses) > 0 {
				if !role.AllowIPSANs {
					return nil, errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided some via CSR")}
				}
				ipAddresses = csr.IPAddresses
			}
		} else {
			ipAltInt, ok = data.GetOk("ip_sans")
			if ok {
				ipAlt := ipAltInt.(string)
				if len(ipAlt) != 0 {
					if !role.AllowIPSANs {
						return nil, errutil.UserError{Err: fmt.Sprintf(
							"IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)}
					}
					for _, v := range strings.Split(ipAlt, ",") {
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
	}

	// Set OU (organizationalUnit) values if specified in the role
	ou := []string{}
	{
		if role.OU != "" {
			ou = strutil.RemoveDuplicates(strutil.ParseStringSlice(role.OU, ","), false)
		}
	}

	// Set O (organization) values if specified in the role
	organization := []string{}
	{
		if role.Organization != "" {
			organization = strutil.RemoveDuplicates(strutil.ParseStringSlice(role.Organization, ","), false)
		}
	}

	// Get the TTL and verify it against the max allowed
	var ttl time.Duration
	var maxTTL time.Duration
	var notAfter time.Time
	{
		ttl = time.Duration(data.Get("ttl").(int)) * time.Second

		if ttl == 0 {
			if role.TTL != "" {
				ttl, err = parseutil.ParseDurationSecond(role.TTL)
				if err != nil {
					return nil, errutil.UserError{Err: fmt.Sprintf(
						"invalid role ttl: %s", err)}
				}
			}
		}

		if role.MaxTTL != "" {
			maxTTL, err = parseutil.ParseDurationSecond(role.MaxTTL)
			if err != nil {
				return nil, errutil.UserError{Err: fmt.Sprintf(
					"invalid role max_ttl: %s", err)}
			}
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
		if signingBundle != nil &&
			notAfter.After(signingBundle.Certificate.NotAfter) && !role.AllowExpirationPastCA {

			return nil, errutil.UserError{Err: fmt.Sprintf(
				"cannot satisfy request, as TTL is beyond the expiration of the CA certificate")}
		}
	}

	// Build up usages
	var extUsage certExtKeyUsage
	{
		if role.ServerFlag {
			extUsage = extUsage | serverExtKeyUsage
		}
		if role.ClientFlag {
			extUsage = extUsage | clientExtKeyUsage
		}
		if role.CodeSigningFlag {
			extUsage = extUsage | codeSigningExtKeyUsage
		}
		if role.EmailProtectionFlag {
			extUsage = extUsage | emailProtectionExtKeyUsage
		}
	}

	creationBundle := &creationBundle{
		CommonName:     cn,
		OU:             ou,
		Organization:   organization,
		DNSNames:       dnsNames,
		EmailAddresses: emailAddresses,
		IPAddresses:    ipAddresses,
		KeyType:        role.KeyType,
		KeyBits:        role.KeyBits,
		SigningBundle:  signingBundle,
		NotAfter:       notAfter,
		KeyUsage:       x509.KeyUsage(parseKeyUsages(role.KeyUsage)),
		ExtKeyUsage:    extUsage,
	}

	// Don't deal with URLs or max path length if it's self-signed, as these
	// normally come from the signing bundle
	if signingBundle == nil {
		return creationBundle, nil
	}

	// This will have been read in from the getURLs function
	creationBundle.URLs = signingBundle.URLs

	// If the max path length in the role is not nil, it was specified at
	// generation time with the max_path_length parameter; otherwise derive it
	// from the signing certificate
	if role.MaxPathLength != nil {
		creationBundle.MaxPathLength = *role.MaxPathLength
	} else {
		switch {
		case signingBundle.Certificate.MaxPathLen < 0:
			creationBundle.MaxPathLength = -1
		case signingBundle.Certificate.MaxPathLen == 0 &&
			signingBundle.Certificate.MaxPathLenZero:
			// The signing function will ensure that we do not issue a CA cert
			creationBundle.MaxPathLength = 0
		default:
			// If this takes it to zero, we handle this case later if
			// necessary
			creationBundle.MaxPathLength = signingBundle.Certificate.MaxPathLen - 1
		}
	}

	return creationBundle, nil
}

// addKeyUsages adds approrpiate key usages to the template given the creation
// information
func addKeyUsages(creationInfo *creationBundle, certTemplate *x509.Certificate) {
	if creationInfo.IsCA {
		certTemplate.KeyUsage = x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign)
		return
	}

	certTemplate.KeyUsage = creationInfo.KeyUsage

	if creationInfo.ExtKeyUsage&serverExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}
	if creationInfo.ExtKeyUsage&clientExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}
	if creationInfo.ExtKeyUsage&codeSigningExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageCodeSigning)
	}
	if creationInfo.ExtKeyUsage&emailProtectionExtKeyUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageEmailProtection)
	}
}

// Performs the heavy lifting of creating a certificate. Returns
// a fully-filled-in ParsedCertBundle.
func createCertificate(creationInfo *creationBundle) (*certutil.ParsedCertBundle, error) {
	var err error
	result := &certutil.ParsedCertBundle{}

	serialNumber, err := certutil.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	if err := certutil.GeneratePrivateKey(creationInfo.KeyType,
		creationInfo.KeyBits,
		result); err != nil {
		return nil, err
	}

	subjKeyID, err := certutil.GetSubjKeyID(result.PrivateKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error getting subject key ID: %s", err)}
	}

	subject := pkix.Name{
		CommonName:         creationInfo.CommonName,
		OrganizationalUnit: creationInfo.OU,
		Organization:       creationInfo.Organization,
	}

	certTemplate := &x509.Certificate{
		SerialNumber:   serialNumber,
		Subject:        subject,
		NotBefore:      time.Now().Add(-30 * time.Second),
		NotAfter:       creationInfo.NotAfter,
		IsCA:           false,
		SubjectKeyId:   subjKeyID,
		DNSNames:       creationInfo.DNSNames,
		EmailAddresses: creationInfo.EmailAddresses,
		IPAddresses:    creationInfo.IPAddresses,
	}

	// Add this before calling addKeyUsages
	if creationInfo.SigningBundle == nil {
		certTemplate.IsCA = true
	}

	// This will only be filled in from the generation paths
	if len(creationInfo.PermittedDNSDomains) > 0 {
		certTemplate.PermittedDNSDomains = creationInfo.PermittedDNSDomains
		certTemplate.PermittedDNSDomainsCritical = true
	}

	addKeyUsages(creationInfo, certTemplate)

	certTemplate.IssuingCertificateURL = creationInfo.URLs.IssuingCertificates
	certTemplate.CRLDistributionPoints = creationInfo.URLs.CRLDistributionPoints
	certTemplate.OCSPServer = creationInfo.URLs.OCSPServers

	var certBytes []byte
	if creationInfo.SigningBundle != nil {
		switch creationInfo.SigningBundle.PrivateKeyType {
		case certutil.RSAPrivateKey:
			certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
		case certutil.ECPrivateKey:
			certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
		}

		caCert := creationInfo.SigningBundle.Certificate
		certTemplate.AuthorityKeyId = caCert.SubjectKeyId

		err = checkPermittedDNSDomains(certTemplate, caCert)
		if err != nil {
			return nil, errutil.UserError{Err: err.Error()}
		}

		certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, result.PrivateKey.Public(), creationInfo.SigningBundle.PrivateKey)
	} else {
		// Creating a self-signed root
		if creationInfo.MaxPathLength == 0 {
			certTemplate.MaxPathLen = 0
			certTemplate.MaxPathLenZero = true
		} else {
			certTemplate.MaxPathLen = creationInfo.MaxPathLength
		}

		switch creationInfo.KeyType {
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

	if creationInfo.SigningBundle != nil {
		if len(creationInfo.SigningBundle.Certificate.AuthorityKeyId) > 0 &&
			!bytes.Equal(creationInfo.SigningBundle.Certificate.AuthorityKeyId, creationInfo.SigningBundle.Certificate.SubjectKeyId) {

			result.CAChain = []*certutil.CertBlock{
				&certutil.CertBlock{
					Certificate: creationInfo.SigningBundle.Certificate,
					Bytes:       creationInfo.SigningBundle.CertificateBytes,
				},
			}
			result.CAChain = append(result.CAChain, creationInfo.SigningBundle.CAChain...)
		}
	}

	return result, nil
}

// Creates a CSR. This is currently only meant for use when
// generating an intermediate certificate.
func createCSR(creationInfo *creationBundle) (*certutil.ParsedCSRBundle, error) {
	var err error
	result := &certutil.ParsedCSRBundle{}

	if err := certutil.GeneratePrivateKey(creationInfo.KeyType,
		creationInfo.KeyBits,
		result); err != nil {
		return nil, err
	}

	// Like many root CAs, other information is ignored
	subject := pkix.Name{
		CommonName: creationInfo.CommonName,
	}

	csrTemplate := &x509.CertificateRequest{
		Subject:        subject,
		DNSNames:       creationInfo.DNSNames,
		EmailAddresses: creationInfo.EmailAddresses,
		IPAddresses:    creationInfo.IPAddresses,
	}

	switch creationInfo.KeyType {
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
func signCertificate(creationInfo *creationBundle,
	csr *x509.CertificateRequest) (*certutil.ParsedCertBundle, error) {
	switch {
	case creationInfo == nil:
		return nil, errutil.UserError{Err: "nil creation info given to signCertificate"}
	case creationInfo.SigningBundle == nil:
		return nil, errutil.UserError{Err: "nil signing bundle given to signCertificate"}
	case csr == nil:
		return nil, errutil.UserError{Err: "nil csr given to signCertificate"}
	}

	err := csr.CheckSignature()
	if err != nil {
		return nil, errutil.UserError{Err: "request signature invalid"}
	}

	result := &certutil.ParsedCertBundle{}

	serialNumber, err := certutil.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(csr.PublicKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error marshalling public key: %s", err)}
	}
	subjKeyID := sha1.Sum(marshaledKey)

	caCert := creationInfo.SigningBundle.Certificate

	subject := pkix.Name{
		CommonName:         creationInfo.CommonName,
		OrganizationalUnit: creationInfo.OU,
		Organization:       creationInfo.Organization,
	}

	certTemplate := &x509.Certificate{
		SerialNumber:   serialNumber,
		Subject:        subject,
		NotBefore:      time.Now().Add(-30 * time.Second),
		NotAfter:       creationInfo.NotAfter,
		SubjectKeyId:   subjKeyID[:],
		AuthorityKeyId: caCert.SubjectKeyId,
	}

	switch creationInfo.SigningBundle.PrivateKeyType {
	case certutil.RSAPrivateKey:
		certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
	case certutil.ECPrivateKey:
		certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
	}

	if creationInfo.UseCSRValues {
		certTemplate.Subject = csr.Subject

		certTemplate.DNSNames = csr.DNSNames
		certTemplate.EmailAddresses = csr.EmailAddresses
		certTemplate.IPAddresses = csr.IPAddresses

		certTemplate.ExtraExtensions = csr.Extensions
	} else {
		certTemplate.DNSNames = creationInfo.DNSNames
		certTemplate.EmailAddresses = creationInfo.EmailAddresses
		certTemplate.IPAddresses = creationInfo.IPAddresses
	}

	addKeyUsages(creationInfo, certTemplate)

	var certBytes []byte

	certTemplate.IssuingCertificateURL = creationInfo.URLs.IssuingCertificates
	certTemplate.CRLDistributionPoints = creationInfo.URLs.CRLDistributionPoints
	certTemplate.OCSPServer = creationInfo.SigningBundle.URLs.OCSPServers

	if creationInfo.IsCA {
		certTemplate.BasicConstraintsValid = true
		certTemplate.IsCA = true

		if creationInfo.SigningBundle.Certificate.MaxPathLen == 0 &&
			creationInfo.SigningBundle.Certificate.MaxPathLenZero {
			return nil, errutil.UserError{Err: "signing certificate has a max path length of zero, and cannot issue further CA certificates"}
		}

		certTemplate.MaxPathLen = creationInfo.MaxPathLength
		if certTemplate.MaxPathLen == 0 {
			certTemplate.MaxPathLenZero = true
		}
	}

	if len(creationInfo.PermittedDNSDomains) > 0 {
		certTemplate.PermittedDNSDomains = creationInfo.PermittedDNSDomains
		certTemplate.PermittedDNSDomainsCritical = true
	}
	err = checkPermittedDNSDomains(certTemplate, caCert)
	if err != nil {
		return nil, errutil.UserError{Err: err.Error()}
	}

	certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, csr.PublicKey, creationInfo.SigningBundle.PrivateKey)

	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CertificateBytes = certBytes
	result.Certificate, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	result.CAChain = creationInfo.SigningBundle.GetCAChain()

	return result, nil
}

func checkPermittedDNSDomains(template, ca *x509.Certificate) error {
	if len(ca.PermittedDNSDomains) == 0 {
		return nil
	}

	namesToCheck := map[string]struct{}{
		template.Subject.CommonName: struct{}{},
	}
	for _, name := range template.DNSNames {
		namesToCheck[name] = struct{}{}
	}

	var badName string
NameCheck:
	for name := range namesToCheck {
		for _, perm := range ca.PermittedDNSDomains {
			switch {
			case strings.HasPrefix(perm, ".") && strings.HasSuffix(name, perm):
				// .example.com matches my.host.example.com and
				// host.example.com but does not match example.com
				break NameCheck
			case perm == name:
				break NameCheck
			}
		}
		badName = name
		break
	}

	if badName == "" {
		return nil
	}

	return fmt.Errorf("name %q disallowed by CA's permitted DNS domains", badName)
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

	keyData, err = certutil.MarshalPKCS8PrivateKey(signer)
	if err != nil {
		return errwrap.Wrapf("error converting response to pkcs8: error marshaling pkcs8 key: {{err}}", err)
	}

	if pemUsed {
		block.Type = "PRIVATE KEY"
		block.Bytes = keyData
		resp.Data["private_key"] = string(pem.EncodeToMemory(block))
	} else {
		resp.Data["private_key"] = base64.StdEncoding.EncodeToString(keyData)
	}

	return nil
}
