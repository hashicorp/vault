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

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/certutil"
)

type certUsage int

const (
	serverUsage certUsage = 1 << iota
	clientUsage
	codeSigningUsage
)

type certCreationBundle struct {
	RawSigningBundle *certutil.RawCertBundle
	CACert           *x509.Certificate
	CommonNames      []string
	IPSANs           []net.IP
	KeyType          string
	KeyBits          int
	Lease            time.Duration
	Usage            certUsage
}

func getCertBundle(s logical.Storage, path string) (*certutil.CertBundle, error) {
	bundle, err := s.Get(path)
	if err != nil {
		return nil, err
	}
	if bundle == nil {
		return nil, nil
	}

	var result certutil.CertBundle
	if err := bundle.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func fetchCAInfo(req *logical.Request) (*certutil.RawCertBundle, *x509.Certificate, error, error) {
	bundle, err := getCertBundle(req.Storage, "config/ca_bundle")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Unable to fetch local CA certificate/key: %s", err)
	}
	if bundle == nil {
		return nil, nil, fmt.Errorf("Backend must be configured with a CA certificate/key"), nil
	}

	rawBundle, err := bundle.ToRawCertBundle()
	if err != nil {
		return nil, nil, nil, err
	}

	certificates, err := x509.ParseCertificates(rawBundle.CertificateBytes)
	switch {
	case err != nil:
		return nil, nil, nil, err
	case len(certificates) != 1:
		return nil, nil, nil, fmt.Errorf("Length of CA certificate bundle is wrong")
	}

	return rawBundle, certificates[0], nil, nil
}

func fetchCertBySerial(req *logical.Request, prefix, serial string) (certEntry *logical.StorageEntry, userError, internalError error) {
	var path string
	var err error
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

	certEntry, err = req.Storage.Get(path)
	if err != nil || certEntry == nil {
		return nil, fmt.Errorf("Certificate with serial number %s not found (if it has been revoked, the revoked/ endpoint must be used)", serial), nil
	}

	if len(certEntry.Value) == 0 {
		return nil, nil, fmt.Errorf("Returned certificate bytes for serial %s were empty", serial)
	}

	return
}

func validateCommonNames(req *logical.Request, commonNames []string, role *roleEntry) (string, error) {
	// TODO: handle wildcards
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

func createCertificate(creationInfo *certCreationBundle) (rawBundle *certutil.RawCertBundle, userErr, intErr error) {
	var clientPrivKey crypto.Signer
	var err error
	rawBundle = &certutil.RawCertBundle{}

	rawBundle.SerialNumber, err = rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting random serial number")
	}

	switch creationInfo.KeyType {
	case "rsa":
		rawBundle.PrivateKeyType = certutil.RSAPrivateKeyType
		clientPrivKey, err = rsa.GenerateKey(rand.Reader, creationInfo.KeyBits)
		if err != nil {
			return nil, nil, fmt.Errorf("Error generating RSA private key")
		}
		rawBundle.PrivateKeyBytes = x509.MarshalPKCS1PrivateKey(clientPrivKey.(*rsa.PrivateKey))
	case "ec":
		rawBundle.PrivateKeyType = certutil.ECPrivateKeyType
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
			return nil, fmt.Errorf("Unsupported bit length for EC key: %d", creationInfo.KeyBits), nil
		}
		clientPrivKey, err = ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("Error generating EC private key")
		}
		rawBundle.PrivateKeyBytes, err = x509.MarshalECPrivateKey(clientPrivKey.(*ecdsa.PrivateKey))
		if err != nil {
			return nil, nil, fmt.Errorf("Error marshalling EC private key")
		}
	default:
		return nil, fmt.Errorf("Unknown key type: %s", creationInfo.KeyType), nil
	}

	subjKeyID, err := rawBundle.GetSubjKeyID()
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting subject key ID: %s", err)
	}

	subject := pkix.Name{
		Country:            creationInfo.CACert.Subject.Country,
		Organization:       creationInfo.CACert.Subject.Organization,
		OrganizationalUnit: creationInfo.CACert.Subject.OrganizationalUnit,
		Locality:           creationInfo.CACert.Subject.Locality,
		Province:           creationInfo.CACert.Subject.Province,
		StreetAddress:      creationInfo.CACert.Subject.StreetAddress,
		PostalCode:         creationInfo.CACert.Subject.PostalCode,
		SerialNumber:       rawBundle.SerialNumber.String(),
		CommonName:         creationInfo.CommonNames[0],
	}

	certTemplate := &x509.Certificate{
		SignatureAlgorithm:    x509.SHA256WithRSA,
		SerialNumber:          rawBundle.SerialNumber,
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

	signingPrivKey, err := creationInfo.RawSigningBundle.GetSigner()
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to get signing private key: %s", err)
	}

	cert, err := x509.CreateCertificate(rand.Reader, certTemplate, creationInfo.CACert, clientPrivKey.Public(), signingPrivKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to create certificate: %s", err)
	}

	rawBundle.CertificateBytes = cert
	rawBundle.IssuingCABytes = creationInfo.RawSigningBundle.CertificateBytes

	return
}
