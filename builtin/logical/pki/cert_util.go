package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
)

type certUsage int

const (
	serverUsage certUsage = 1 << iota
	clientUsage
	codeSigningUsage
)

type certCreationBundle struct {
	SigningBundle *certutil.ParsedCertBundle
	CACert        *x509.Certificate
	CommonNames   []string
	IPSANs        []net.IP
	KeyType       string
	KeyBits       int
	Lease         time.Duration
	Usage         certUsage
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
func validateCommonNames(req *logical.Request, commonNames []string, role *roleEntry) (string, error) {
	hostnameRegex, err := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if err != nil {
		return "", fmt.Errorf("Error compiling hostname regex: %s", err)
	}
	subdomainRegex, err := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]))*$`)
	if err != nil {
		return "", fmt.Errorf("Error compiling subdomain regex: %s", err)
	}
	for _, name := range commonNames {
		if role.AllowLocalhost && name == "localhost" {
			continue
		}

		sanitizedName := name
		isWildcard := false
		if strings.HasPrefix(name, "*.") {
			sanitizedName = name[2:]
			isWildcard = true
		}
		if !hostnameRegex.MatchString(sanitizedName) {
			return name, nil
		}
		if role.AllowAnyName {
			continue
		}

		if role.AllowTokenDisplayName {
			if name == req.DisplayName {
				continue
			}

			if role.AllowSubdomains {
				if strings.HasSuffix(name, "."+req.DisplayName) {
					continue
				}
			}
		}

		if len(role.AllowedBaseDomain) != 0 {
			if strings.HasSuffix(name, "."+role.AllowedBaseDomain) {
				if role.AllowSubdomains {
					continue
				}

				if subdomainRegex.MatchString(strings.TrimSuffix(name, "."+role.AllowedBaseDomain)) {
					continue
				}

				if isWildcard && role.AllowedBaseDomain == sanitizedName {
					continue
				}
			}
		}

		return name, nil
	}

	return "", nil
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
		return nil, certutil.InternalError{Err: fmt.Sprintf("Error getting random serial number")}
	}

	switch creationInfo.KeyType {
	case "rsa":
		result.PrivateKeyType = certutil.RSAPrivateKey
		clientPrivKey, err = rsa.GenerateKey(rand.Reader, creationInfo.KeyBits)
		if err != nil {
			return nil, certutil.InternalError{Err: fmt.Sprintf("Error generating RSA private key")}
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
			return nil, certutil.UserError{Err: fmt.Sprintf("Unsupported bit length for EC key: %d", creationInfo.KeyBits)}
		}
		clientPrivKey, err = ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, certutil.InternalError{Err: fmt.Sprintf("Error generating EC private key")}
		}
		result.PrivateKey = clientPrivKey
		result.PrivateKeyBytes, err = x509.MarshalECPrivateKey(clientPrivKey.(*ecdsa.PrivateKey))
		if err != nil {
			return nil, certutil.InternalError{Err: fmt.Sprintf("Error marshalling EC private key")}
		}
	default:
		return nil, certutil.UserError{Err: fmt.Sprintf("Unknown key type: %s", creationInfo.KeyType)}
	}

	subjKeyID, err := certutil.GetSubjKeyID(result.PrivateKey)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("Error getting subject key ID: %s", err)}
	}

	subject := pkix.Name{
		Country:            creationInfo.CACert.Subject.Country,
		Organization:       creationInfo.CACert.Subject.Organization,
		OrganizationalUnit: creationInfo.CACert.Subject.OrganizationalUnit,
		Locality:           creationInfo.CACert.Subject.Locality,
		Province:           creationInfo.CACert.Subject.Province,
		StreetAddress:      creationInfo.CACert.Subject.StreetAddress,
		PostalCode:         creationInfo.CACert.Subject.PostalCode,
		SerialNumber:       serialNumber.String(),
		CommonName:         creationInfo.CommonNames[0],
	}

	certTemplate := &x509.Certificate{
		SignatureAlgorithm:    x509.SHA256WithRSA,
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(creationInfo.Lease),
		KeyUsage:              x509.KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement),
		BasicConstraintsValid: true,
		IsCA:                        false,
		SubjectKeyId:                subjKeyID,
		DNSNames:                    creationInfo.CommonNames,
		IPAddresses:                 creationInfo.IPSANs,
		PermittedDNSDomainsCritical: false,
		PermittedDNSDomains:         nil,
		CRLDistributionPoints:       creationInfo.CACert.CRLDistributionPoints,
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

	cert, err := x509.CreateCertificate(rand.Reader, certTemplate, creationInfo.CACert, clientPrivKey.Public(), creationInfo.SigningBundle.PrivateKey)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("Unable to create certificate: %s", err)}
	}

	result.CertificateBytes = cert
	result.Certificate, err = x509.ParseCertificate(cert)
	if err != nil {
		return nil, certutil.InternalError{Err: fmt.Sprintf("Unable to parse created certificate: %s", err)}
	}

	result.IssuingCABytes = creationInfo.SigningBundle.CertificateBytes
	result.IssuingCA = creationInfo.SigningBundle.Certificate

	return result, nil
}
