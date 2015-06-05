package pki

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"math/big"
	mrand "math/rand"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
)

// The type of Private Key, for storage
const (
	UnknownPrivateKeyType = iota
	RSAPrivateKeyType
	ECPrivateKeyType
)

type certUsage int

const (
	serverUsage certUsage = 1 << iota
	clientUsage
	codeSigningUsage
)

type certBundle struct {
	PrivateKeyType    int    `json:"private_key_type"`
	PrivateKeyString  string `json:"private_key_string"`
	CertificateString string `json:"certificate_string"`
}

type rawCertBundle struct {
	PrivateKeyType   int
	PrivateKeyBytes  []byte
	CertificateBytes []byte
	SerialNumber     *big.Int
}

type certCreationBundle struct {
	RawSigningBundle *rawCertBundle
	CACert           *x509.Certificate
	CommonNames      []string
	IPSANs           []net.IP
	KeyType          string
	KeyBits          int
	Lease            time.Duration
	Usage            certUsage
}

func (c *certBundle) toRawCertBundle() (*rawCertBundle, error) {
	decoder := base64.URLEncoding
	result := &rawCertBundle{
		PrivateKeyType: c.PrivateKeyType,
	}
	var err error
	if result.PrivateKeyBytes, err = decoder.DecodeString(c.PrivateKeyString); err != nil {
		return nil, err
	}
	if result.CertificateBytes, err = decoder.DecodeString(c.CertificateString); err != nil {
		return nil, err
	}

	if err := result.populateSerialNumber(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *rawCertBundle) toCertBundle() *certBundle {
	encoder := base64.URLEncoding
	result := &certBundle{
		PrivateKeyType:    r.PrivateKeyType,
		PrivateKeyString:  encoder.EncodeToString(r.PrivateKeyBytes),
		CertificateString: encoder.EncodeToString(r.CertificateBytes),
	}
	return result
}

func (r *rawCertBundle) populateSerialNumber() error {
	cert, err := x509.ParseCertificate(r.CertificateBytes)
	if err != nil {
		return fmt.Errorf("Error encountered parsing certificate bytes from raw bundle")
	}
	r.SerialNumber = cert.SerialNumber
	return nil
}

// "Signer" corresponds to the Go interface that private keys implement
// that provides a Public() function for getting the corresponding public
// key. It can be type converted to private keys.
func (r *rawCertBundle) getSigner() (crypto.Signer, error) {
	var signer crypto.Signer
	var err error
	switch r.PrivateKeyType {
	case ECPrivateKeyType:
		signer, err = x509.ParseECPrivateKey(r.PrivateKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse CA's private EC key: %s", err)
		}
	case RSAPrivateKeyType:
		signer, err = x509.ParsePKCS1PrivateKey(r.PrivateKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse CA's private RSA key: %s", err)
		}
	default:
		return nil, fmt.Errorf("Unable to determine the type of private key")
	}
	return signer, nil
}

func (r *rawCertBundle) getSubjKeyID() ([]byte, error) {
	privateKey, err := r.getSigner()
	if err != nil {
		return nil, err
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("Error marshalling public key: %s", err)
	}

	subjKeyID := sha1.Sum(marshaledKey)

	return subjKeyID[:], nil
}

func getCertBundle(s logical.Storage, path string) (*certBundle, error) {
	bundle, err := s.Get(path)
	if err != nil {
		return nil, err
	}
	if bundle == nil {
		return nil, nil
	}

	var result certBundle
	if err := bundle.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func getOctalFormatted(buf []byte, sep string) string {
	var ret bytes.Buffer
	for _, cur := range buf {
		if ret.Len() > 0 {
			fmt.Fprintf(&ret, sep)
		}
		fmt.Fprintf(&ret, "%02x", cur)
	}
	return ret.String()
}

func fetchCAInfo(req *logical.Request) (*rawCertBundle, *x509.Certificate, error) {
	bundle, err := getCertBundle(req.Storage, "config/ca_bundle")
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to fetch local CA certificate/key: %s", err)
	}
	if bundle == nil {
		return nil, nil, fmt.Errorf("Backend must be configured with a CA certificate/key")
	}

	rawBundle, err := bundle.toRawCertBundle()
	if err != nil {
		return nil, nil, err
	}

	certificates, err := x509.ParseCertificates(rawBundle.CertificateBytes)
	switch {
	case err != nil:
		return nil, nil, err
	case len(certificates) != 1:
		return nil, nil, fmt.Errorf("Length of CA certificate bundle is wrong")
	}

	return rawBundle, certificates[0], nil
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

func createCertificate(creationInfo *certCreationBundle) (rawBundle *rawCertBundle, userErr, intErr error) {
	rawBundle = &rawCertBundle{
		SerialNumber: (&big.Int{}).Rand(mrand.New(mrand.NewSource(time.Now().UnixNano())), (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil)),
	}

	var clientPrivKey crypto.Signer
	var err error
	switch creationInfo.KeyType {
	case "rsa":
		rawBundle.PrivateKeyType = RSAPrivateKeyType
		clientPrivKey, err = rsa.GenerateKey(crand.Reader, creationInfo.KeyBits)
		if err != nil {
			return nil, nil, fmt.Errorf("Error generating RSA private key")
		}
		rawBundle.PrivateKeyBytes = x509.MarshalPKCS1PrivateKey(clientPrivKey.(*rsa.PrivateKey))
	case "ec":
		rawBundle.PrivateKeyType = ECPrivateKeyType
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
		clientPrivKey, err = ecdsa.GenerateKey(curve, crand.Reader)
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

	subjKeyID, err := rawBundle.getSubjKeyID()
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

	signingPrivKey, err := creationInfo.RawSigningBundle.getSigner()
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to get signing private key: %s", err)
	}

	cert, err := x509.CreateCertificate(crand.Reader, certTemplate, creationInfo.CACert, clientPrivKey.Public(), signingPrivKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to create certificate: %s", err)
	}

	rawBundle.CertificateBytes = cert

	return
}
