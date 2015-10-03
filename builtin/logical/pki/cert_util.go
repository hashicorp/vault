package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/otto/helper/uuid"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type certUsage int

const (
	serverUsage certUsage = 1 << iota
	clientUsage
	codeSigningUsage
	emailProtectionUsage
	caUsage
)

type certCreationBundle struct {
	CommonName     string
	DNSNames       []string
	EmailAddresses []string
	IPAddresses    []net.IP
	PKIAddress     string
	KeyType        string
	KeyBits        int
	SigningBundle  *certutil.ParsedCertBundle
	TTL            time.Duration
	Usage          certUsage
}

// Fetches the CA info. Unlike other certificates, the CA info is stored
// in the backend as a CertBundle, because we are storing its private key
func fetchCAInfo(req *logical.Request) (*certutil.ParsedCertBundle, error) {
	bundleEntry, err := req.Storage.Get("config/ca_bundle")
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("Unable to fetch local CA certificate/key: %s", err)}
	}
	if bundleEntry == nil {
		return nil, certutil.UserError{Err: fmt.Sprintf("Backend must be configured with a CA certificate/key")}
	}

	var bundle certutil.CertBundle
	if err := bundleEntry.DecodeJSON(&bundle); err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("Unable to decode local CA certificate/key: %s", err)}
	}

	parsedBundle, err := bundle.ToParsedCertBundle()
	if err != nil {
		return nil, certutil.InternalError{Err: err.Error()}
	}

	if parsedBundle.Certificate == nil {
		return nil, certutil.InternalError{Err: "Stored CA information not able to be parsed"}
	}

	return parsedBundle, nil
}

// Allows fetching certificates from the backend; it handles the slightly
// separate pathing for CA, CRL, and revoked certificates.
func fetchCertBySerial(req *logical.Request, prefix, serial string) (*logical.StorageEntry, error) {
	var path string

	switch {
	case serial == "ca":
		path = "ca"
	case serial == "crl":
		path = "crl"
	case strings.HasPrefix(prefix, "revoked/"):
		path = "revoked/" + strings.Replace(strings.ToLower(serial), "-", ":", -1)
	default:
		path = "certs/" + strings.Replace(strings.ToLower(serial), "-", ":", -1)
	}

	certEntry, err := req.Storage.Get(path)
	if err != nil || certEntry == nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("Certificate with serial number %s not found", serial)}
	}

	if certEntry.Value == nil || len(certEntry.Value) == 0 {
		return nil, certutil.InternalError{Err: fmt.Sprintf("Returned certificate bytes for serial %s were empty", serial)}
	}

	return certEntry, nil
}

// Given a set of requested names for a certificate, verifies that all of them
// match the various toggles set in the role for controlling issuance.
// If one does not pass, it is returned in the string argument.
func validateNames(req *logical.Request, names []string, role *roleEntry) (string, error) {
	hostnameRegex, err := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if err != nil {
		return "", fmt.Errorf("error compiling hostname regex: %s", err)
	}
	for _, name := range names {
		sanitizedName := name
		emailDomain := name
		isEmail := false
		isWildcard := false
		var emailUserPart string
		if strings.Contains(name, "@") {
			if !role.EmailProtectionFlag {
				return name, nil
			}
			splitEmail := strings.Split(name, "@")
			if len(splitEmail) != 2 {
				return name, nil
			}
			emailUserPart = splitEmail[0]
			sanitizedName = splitEmail[1]
			emailDomain = splitEmail[1]
			isEmail = true
		}
		if strings.HasPrefix(sanitizedName, "*.") {
			sanitizedName = sanitizedName[2:]
			isWildcard = true
		}

		if role.EnforceHostnames {
			if !hostnameRegex.MatchString(sanitizedName) {
				return name, nil
			}
			if isEmail {
				if !hostnameRegex.MatchString(emailUserPart) {
					return name, nil
				}
			}
		}

		if role.AllowLocalhost {
			if name == "localhost" || (isEmail && emailDomain == "localhost") {
				continue
			}

			if role.AllowSubdomains {
				if strings.HasSuffix(sanitizedName, "."+req.DisplayName) {
					// Regex won't match empty string, so protected against
					// someone trying to get base without it being allowed
					trimmed := strings.TrimSuffix(sanitizedName, "."+req.DisplayName)
					if hostnameRegex.MatchString(trimmed) {
						continue
					}
				}

				// The exception to matching the empty string is if it's *.
				// which is pulled out of the sanitized name, so see if it's
				// a wildcard matching the allowed base domain
				if isWildcard && sanitizedName == role.AllowedBaseDomain {
					continue
				}
			}
		}

		if role.AllowAnyName {
			continue
		}

		if role.AllowTokenDisplayName {
			// Don't check sanitized here because it needs to be exact
			if name == req.DisplayName || (isEmail && emailDomain == req.DisplayName) {
				continue
			}

			if role.AllowSubdomains {
				if strings.HasSuffix(sanitizedName, "."+req.DisplayName) {
					trimmed := strings.TrimSuffix(sanitizedName, "."+req.DisplayName)
					if hostnameRegex.MatchString(trimmed) {
						continue
					}
				}

				if isWildcard && sanitizedName == role.AllowedBaseDomain {
					continue
				}
			}
		}

		if len(role.AllowedBaseDomain) != 0 {
			if role.AllowBaseDomain &&
				(name == role.AllowedBaseDomain ||
					(isEmail && emailDomain == role.AllowedBaseDomain)) {
				continue
			}

			if role.AllowSubdomains {
				if strings.HasSuffix(sanitizedName, "."+role.AllowedBaseDomain) {
					trimmed := strings.TrimSuffix(sanitizedName, "."+role.AllowedBaseDomain)
					if hostnameRegex.MatchString(trimmed) {
						continue
					}
				}

				if isWildcard && sanitizedName == role.AllowedBaseDomain {
					continue
				}
			}
		}

		return name, nil
	}

	return "", nil
}

func generateCert(b *backend,
	role *roleEntry,
	signingBundle *certutil.ParsedCertBundle,
	pkiAddress string,
	req *logical.Request,
	data *framework.FieldData) (*certutil.ParsedCertBundle, error) {

	creationBundle, err := generateCreationBundle(b, role, signingBundle, pkiAddress, req, data)
	if err != nil {
		return nil, err
	}

	parsedBundle, err := createCertificate(creationBundle)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

// N.B.: This is only meant to be used for generating intermediate CAs.
// It skips some sanity checks.
func generateCSR(b *backend,
	role *roleEntry,
	signingBundle *certutil.ParsedCertBundle,
	req *logical.Request,
	data *framework.FieldData) (*certutil.ParsedCSRBundle, error) {

	creationBundle := &certCreationBundle{
		DNSNames: []string{
			fmt.Sprintf("intcsr-%s-%s", strings.TrimSuffix(req.MountPoint, "/"), uuid.GenerateUUID()),
		},
		KeyType: "rsa",
		KeyBits: role.KeyBits,
	}

	parsedBundle, err := createCSR(creationBundle)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

func signCert(b *backend,
	role *roleEntry,
	signingBundle *certutil.ParsedCertBundle,
	pkiAddress string,
	req *logical.Request,
	data *framework.FieldData) (*certutil.ParsedCertBundle, error) {

	csrString := req.Data["csr"].(string)
	if csrString == "" {
		return nil, certutil.UserError{Err: fmt.Sprintf(
			"\"csr\" is empty")}
	}

	pemBytes := []byte(csrString)
	pemBlock, pemBytes := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, certutil.UserError{Err: "csr contains no data"}
	}
	csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		return nil, certutil.UserError{Err: "certificate request could not be parsed"}
	}

	creationBundle, err := generateCreationBundle(b, role, signingBundle, pkiAddress, req, data)
	if err != nil {
		return nil, err
	}

	parsedBundle, err := signCertificate(creationBundle, csr)
	if err != nil {
		return nil, err
	}

	return parsedBundle, nil
}

func generateCreationBundle(b *backend,
	role *roleEntry,
	signingBundle *certutil.ParsedCertBundle,
	pkiAddress string,
	req *logical.Request,
	data *framework.FieldData) (*certCreationBundle, error) {
	var err error

	// Get the common name(s)
	cn := data.Get("common_name").(string)
	if len(cn) == 0 {
		return nil, certutil.UserError{Err: "The common_name field is required"}
	}

	dnsNames := []string{}
	emailAddresses := []string{}
	if strings.Contains(cn, "@") {
		emailAddresses = append(emailAddresses, cn)
	} else {
		dnsNames = append(dnsNames, cn)
	}
	cnAltInt, ok := data.GetOk("alt_names")
	if ok {
		cnAlt := cnAltInt.(string)
		if len(cnAlt) != 0 {
			for _, v := range strings.Split(cnAlt, ",") {
				if strings.Contains(v, "@") {
					emailAddresses = append(emailAddresses, cn)
				} else {
					dnsNames = append(dnsNames, v)
				}
			}
		}
	}

	// Get any IP SANs
	ipAddresses := []net.IP{}
	ipAltInt, ok := data.GetOk("ip_sans")
	if ok {
		ipAlt := ipAltInt.(string)
		if len(ipAlt) != 0 {
			if !role.AllowIPSANs {
				return nil, certutil.UserError{Err: fmt.Sprintf(
					"IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)}
			}
			for _, v := range strings.Split(ipAlt, ",") {
				parsedIP := net.ParseIP(v)
				if parsedIP == nil {
					return nil, certutil.UserError{Err: fmt.Sprintf(
						"the value '%s' is not a valid IP address", v)}
				}
				ipAddresses = append(ipAddresses, parsedIP)
			}
		}
	}

	var ttlField string
	ttlFieldInt, ok := data.GetOk("ttl")
	if !ok {
		ttlField = role.TTL
	} else {
		ttlField = ttlFieldInt.(string)
	}

	var ttl time.Duration
	if len(ttlField) == 0 {
		ttl = b.System().DefaultLeaseTTL()
	} else {
		ttl, err = time.ParseDuration(ttlField)
		if err != nil {
			return nil, certutil.UserError{Err: fmt.Sprintf(
				"invalid requested ttl: %s", err)}
		}
	}

	var maxTTL time.Duration
	if len(role.MaxTTL) == 0 {
		maxTTL = b.System().MaxLeaseTTL()
	} else {
		maxTTL, err = time.ParseDuration(role.MaxTTL)
		if err != nil {
			return nil, certutil.UserError{Err: fmt.Sprintf(
				"invalid ttl: %s", err)}
		}
	}

	if ttl > maxTTL {
		// Don't error if they were using system defaults, only error if
		// they specifically chose a bad TTL
		if len(ttlField) == 0 {
			ttl = maxTTL
		} else {
			return nil, certutil.UserError{Err: fmt.Sprintf(
				"ttl is larger than maximum allowed by this role")}
		}
	}

	if signingBundle != nil &&
		time.Now().Add(ttl).After(signingBundle.Certificate.NotAfter) {
		return nil, certutil.UserError{Err: fmt.Sprintf(
			"cannot satisfy request, as TTL is beyond the expiration of the CA certificate")}
	}

	badName, err := validateNames(req, dnsNames, role)
	if len(badName) != 0 {
		return nil, certutil.UserError{Err: fmt.Sprintf(
			"name %s not allowed by this role", badName)}
	} else if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf(
			"error validating name %s: %s", badName, err)}
	}

	badName, err = validateNames(req, emailAddresses, role)
	if len(badName) != 0 {
		return nil, certutil.UserError{Err: fmt.Sprintf(
			"email %s not allowed by this role", badName)}
	} else if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf(
			"error validating name %s: %s", badName, err)}
	}

	var usage certUsage
	if role.ServerFlag {
		usage = usage | serverUsage
	}
	if role.ClientFlag {
		usage = usage | clientUsage
	}
	if role.CodeSigningFlag {
		usage = usage | codeSigningUsage
	}
	if role.EmailProtectionFlag {
		usage = usage | emailProtectionUsage
	}

	creationBundle := &certCreationBundle{
		CommonName:     cn,
		DNSNames:       dnsNames,
		EmailAddresses: emailAddresses,
		IPAddresses:    ipAddresses,
		KeyType:        role.KeyType,
		KeyBits:        role.KeyBits,
		PKIAddress:     pkiAddress,
		SigningBundle:  signingBundle,
		TTL:            ttl,
		Usage:          usage,
	}

	return creationBundle, nil
}

// Performs the heavy lifting of creating a certificate. Returns
// a fully-filled-in ParsedCertBundle.
func createCertificate(creationInfo *certCreationBundle) (*certutil.ParsedCertBundle, error) {
	var clientPrivKey crypto.Signer
	var err error
	result := &certutil.ParsedCertBundle{}

	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("error getting random serial number")}
	}

	switch creationInfo.KeyType {
	case "rsa":
		result.PrivateKeyType = certutil.RSAPrivateKey
		clientPrivKey, err = rsa.GenerateKey(rand.Reader, creationInfo.KeyBits)
		if err != nil {
			return nil, certutil.InternalError{Err: fmt.Sprintf("error generating RSA private key")}
		}
		result.PrivateKey = clientPrivKey
		result.PrivateKeyBytes = x509.MarshalPKCS1PrivateKey(clientPrivKey.(*rsa.PrivateKey))
	case "ec":
		result.PrivateKeyType = certutil.ECPrivateKey
		var curve elliptic.Curve
		switch creationInfo.KeyBits {
		case 224:
			curve = elliptic.P224()
		case 256:
			curve = elliptic.P256()
		case 384:
			curve = elliptic.P384()
		case 521:
			curve = elliptic.P521()
		default:
			return nil, certutil.UserError{Err: fmt.Sprintf("unsupported bit length for EC key: %d", creationInfo.KeyBits)}
		}
		clientPrivKey, err = ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, certutil.InternalError{Err: fmt.Sprintf("error generating EC private key")}
		}
		result.PrivateKey = clientPrivKey
		result.PrivateKeyBytes, err = x509.MarshalECPrivateKey(clientPrivKey.(*ecdsa.PrivateKey))
		if err != nil {
			return nil, certutil.InternalError{Err: fmt.Sprintf("error marshalling EC private key")}
		}
	default:
		return nil, certutil.UserError{Err: fmt.Sprintf("unknown key type: %s", creationInfo.KeyType)}
	}

	subjKeyID, err := certutil.GetSubjKeyID(result.PrivateKey)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("error getting subject key ID: %s", err)}
	}

	subject := pkix.Name{
		SerialNumber: serialNumber.String(),
		CommonName:   creationInfo.CommonName,
	}

	certTemplate := &x509.Certificate{
		SignatureAlgorithm:    x509.SHA256WithRSA,
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(creationInfo.TTL),
		KeyUsage:              x509.KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement),
		BasicConstraintsValid: true,
		IsCA:           false,
		SubjectKeyId:   subjKeyID,
		DNSNames:       creationInfo.DNSNames,
		EmailAddresses: creationInfo.EmailAddresses,
		IPAddresses:    creationInfo.IPAddresses,
	}

	if creationInfo.Usage&serverUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}
	if creationInfo.Usage&clientUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}
	if creationInfo.Usage&codeSigningUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageCodeSigning)
	}
	if creationInfo.Usage&emailProtectionUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageEmailProtection)
	}

	var certBytes []byte
	if creationInfo.SigningBundle != nil {
		caCert := creationInfo.SigningBundle.Certificate
		subject.Country = caCert.Subject.Country
		subject.Organization = caCert.Subject.Organization
		subject.OrganizationalUnit = caCert.Subject.OrganizationalUnit
		subject.Locality = caCert.Subject.Locality
		subject.Province = caCert.Subject.Province
		subject.StreetAddress = caCert.Subject.StreetAddress
		subject.PostalCode = caCert.Subject.PostalCode
		certTemplate.CRLDistributionPoints = caCert.CRLDistributionPoints
		certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, clientPrivKey.Public(), creationInfo.SigningBundle.PrivateKey)
	} else {
		certTemplate.CRLDistributionPoints = []string{
			creationInfo.PKIAddress + "/crl",
		}
		certTemplate.IssuingCertificateURL = []string{
			creationInfo.PKIAddress + "/ca",
		}
		certTemplate.IsCA = true
		certTemplate.KeyUsage = x509.KeyUsage(certTemplate.KeyUsage | x509.KeyUsageCertSign | x509.KeyUsageCRLSign)
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageOCSPSigning)
		certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, clientPrivKey.Public(), clientPrivKey)
	}

	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CertificateBytes = certBytes
	result.Certificate, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	if creationInfo.SigningBundle != nil {
		result.IssuingCABytes = creationInfo.SigningBundle.CertificateBytes
		result.IssuingCA = creationInfo.SigningBundle.Certificate
	} else {
		result.IssuingCABytes = result.CertificateBytes
		result.IssuingCA = result.Certificate
	}

	return result, nil
}

// Creates a CSR. This is currently only meant for use when
// generating an intermediate certificate.
func createCSR(creationInfo *certCreationBundle) (*certutil.ParsedCSRBundle, error) {
	var clientPrivKey crypto.Signer
	var err error
	result := &certutil.ParsedCSRBundle{}

	switch creationInfo.KeyType {
	case "rsa":
		result.PrivateKeyType = certutil.RSAPrivateKey
		clientPrivKey, err = rsa.GenerateKey(rand.Reader, creationInfo.KeyBits)
		if err != nil {
			return nil, certutil.InternalError{Err: fmt.Sprintf("error generating RSA private key")}
		}
		result.PrivateKey = clientPrivKey
		result.PrivateKeyBytes = x509.MarshalPKCS1PrivateKey(clientPrivKey.(*rsa.PrivateKey))
	default:
		return nil, certutil.UserError{Err: fmt.Sprintf("unsupported key type for CA generation: %s", creationInfo.KeyType)}
	}

	// Like many root CAs, other information is ignored
	subject := pkix.Name{
		CommonName: creationInfo.CommonName,
	}

	csrTemplate := &x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject:            subject,
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, result.PrivateKey)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CSRBytes = csr
	result.CSR, err = x509.ParseCertificateRequest(csr)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	return result, nil
}

// Performs the heavy lifting of generating a certificate from a CSR.
// Returns a ParsedCertBundle sans private keys.
func signCertificate(creationInfo *certCreationBundle,
	csr *x509.CertificateRequest) (*certutil.ParsedCertBundle, error) {
	switch {
	case creationInfo == nil:
		return nil, certutil.UserError{Err: "nil creation info given to signCertificate"}
	case creationInfo.SigningBundle == nil:
		return nil, certutil.UserError{Err: "nil signing bundle given to signCertificate"}
	case csr == nil:
		return nil, certutil.UserError{Err: "nil csr given to signCertificate"}
	}

	err := csr.CheckSignature()
	if err != nil {
		return nil, certutil.UserError{Err: "request signature invalid"}
	}

	result := &certutil.ParsedCertBundle{}

	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("error getting random serial number")}
	}

	subject := pkix.Name{
		SerialNumber: serialNumber.String(),
		CommonName:   creationInfo.CommonName,
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(csr.PublicKey)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("error marshalling public key: %s", err)}
	}
	subjKeyID := sha1.Sum(marshaledKey)

	certTemplate := &x509.Certificate{
		SignatureAlgorithm:    x509.SHA256WithRSA,
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(creationInfo.TTL),
		KeyUsage:              x509.KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement),
		BasicConstraintsValid: true,
		IsCA:           false,
		SubjectKeyId:   subjKeyID[:],
		DNSNames:       creationInfo.DNSNames,
		EmailAddresses: creationInfo.EmailAddresses,
		IPAddresses:    creationInfo.IPAddresses,
	}

	if creationInfo.Usage&serverUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}
	if creationInfo.Usage&clientUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}
	if creationInfo.Usage&codeSigningUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageCodeSigning)
	}
	if creationInfo.Usage&emailProtectionUsage != 0 {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageEmailProtection)
	}

	var certBytes []byte
	caCert := creationInfo.SigningBundle.Certificate
	subject.Country = caCert.Subject.Country
	subject.Organization = caCert.Subject.Organization
	subject.OrganizationalUnit = caCert.Subject.OrganizationalUnit
	subject.Locality = caCert.Subject.Locality
	subject.Province = caCert.Subject.Province
	subject.StreetAddress = caCert.Subject.StreetAddress
	subject.PostalCode = caCert.Subject.PostalCode

	certTemplate.IssuingCertificateURL = caCert.IssuingCertificateURL

	if creationInfo.PKIAddress != "" {
		certTemplate.CRLDistributionPoints = []string{
			creationInfo.PKIAddress + "/crl",
		}
		certTemplate.IsCA = true
		certTemplate.KeyUsage = x509.KeyUsage(certTemplate.KeyUsage | x509.KeyUsageCertSign | x509.KeyUsageCRLSign)
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageOCSPSigning)
	} else {
		certTemplate.CRLDistributionPoints = caCert.CRLDistributionPoints
	}

	certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, csr.PublicKey, creationInfo.SigningBundle.PrivateKey)

	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CertificateBytes = certBytes
	result.Certificate, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	result.IssuingCABytes = creationInfo.SigningBundle.CertificateBytes
	result.IssuingCA = creationInfo.SigningBundle.Certificate

	return result, nil
}
