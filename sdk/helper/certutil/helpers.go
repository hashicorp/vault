package certutil

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/cryptobyte"
	cbasn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// GetHexFormatted returns the byte buffer formatted in hex with
// the specified separator between bytes.
func GetHexFormatted(buf []byte, sep string) string {
	var ret bytes.Buffer
	for _, cur := range buf {
		if ret.Len() > 0 {
			fmt.Fprintf(&ret, sep)
		}
		fmt.Fprintf(&ret, "%02x", cur)
	}
	return ret.String()
}

// ParseHexFormatted returns the raw bytes from a formatted hex string
func ParseHexFormatted(in, sep string) []byte {
	var ret bytes.Buffer
	var err error
	var inBits int64
	inBytes := strings.Split(in, sep)
	for _, inByte := range inBytes {
		if inBits, err = strconv.ParseInt(inByte, 16, 8); err != nil {
			return nil
		}
		ret.WriteByte(byte(inBits))
	}
	return ret.Bytes()
}

// GetSubjKeyID returns the subject key ID, e.g. the SHA1 sum
// of the marshaled public key
func GetSubjKeyID(privateKey crypto.Signer) ([]byte, error) {
	if privateKey == nil {
		return nil, errutil.InternalError{Err: "passed-in private key is nil"}
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error marshalling public key: %s", err)}
	}

	subjKeyID := sha1.Sum(marshaledKey)

	return subjKeyID[:], nil
}

// ParsePKIMap takes a map (for instance, the Secret.Data
// returned from the PKI backend) and returns a ParsedCertBundle.
func ParsePKIMap(data map[string]interface{}) (*ParsedCertBundle, error) {
	result := &CertBundle{}
	err := mapstructure.Decode(data, result)
	if err != nil {
		return nil, errutil.UserError{Err: err.Error()}
	}

	return result.ToParsedCertBundle()
}

// ParsePKIJSON takes a JSON-encoded string and returns a ParsedCertBundle.
//
// This can be either the output of an
// issue call from the PKI backend or just its data member; or,
// JSON not coming from the PKI backend.
func ParsePKIJSON(input []byte) (*ParsedCertBundle, error) {
	result := &CertBundle{}
	err := jsonutil.DecodeJSON(input, &result)

	if err == nil {
		return result.ToParsedCertBundle()
	}

	var secret Secret
	err = jsonutil.DecodeJSON(input, &secret)

	if err == nil {
		return ParsePKIMap(secret.Data)
	}

	return nil, errutil.UserError{Err: "unable to parse out of either secret data or a secret object"}
}

// ParsePEMBundle takes a string of concatenated PEM-format certificate
// and private key values and decodes/parses them, checking validity along
// the way. The first certificate must be the subject certificate and issuing
// certificates may follow.  There must be at most one private key.
func ParsePEMBundle(pemBundle string) (*ParsedCertBundle, error) {
	if len(pemBundle) == 0 {
		return nil, errutil.UserError{Err: "empty pem bundle"}
	}

	pemBytes := []byte(pemBundle)
	var pemBlock *pem.Block
	parsedBundle := &ParsedCertBundle{}
	var certPath []*CertBlock

	for len(pemBytes) > 0 {
		pemBlock, pemBytes = pem.Decode(pemBytes)
		if pemBlock == nil {
			return nil, errutil.UserError{Err: "no data found in PEM block"}
		}

		if signer, err := x509.ParseECPrivateKey(pemBlock.Bytes); err == nil {
			if parsedBundle.PrivateKeyType != UnknownPrivateKey {
				return nil, errutil.UserError{Err: "more than one private key given; provide only one private key in the bundle"}
			}
			parsedBundle.PrivateKeyFormat = ECBlock
			parsedBundle.PrivateKeyType = ECPrivateKey
			parsedBundle.PrivateKeyBytes = pemBlock.Bytes
			parsedBundle.PrivateKey = signer

		} else if signer, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err == nil {
			if parsedBundle.PrivateKeyType != UnknownPrivateKey {
				return nil, errutil.UserError{Err: "more than one private key given; provide only one private key in the bundle"}
			}
			parsedBundle.PrivateKeyType = RSAPrivateKey
			parsedBundle.PrivateKeyFormat = PKCS1Block
			parsedBundle.PrivateKeyBytes = pemBlock.Bytes
			parsedBundle.PrivateKey = signer
		} else if signer, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes); err == nil {
			parsedBundle.PrivateKeyFormat = PKCS8Block

			if parsedBundle.PrivateKeyType != UnknownPrivateKey {
				return nil, errutil.UserError{Err: "More than one private key given; provide only one private key in the bundle"}
			}
			switch signer := signer.(type) {
			case *rsa.PrivateKey:
				parsedBundle.PrivateKey = signer
				parsedBundle.PrivateKeyType = RSAPrivateKey
				parsedBundle.PrivateKeyBytes = pemBlock.Bytes
			case *ecdsa.PrivateKey:
				parsedBundle.PrivateKey = signer
				parsedBundle.PrivateKeyType = ECPrivateKey
				parsedBundle.PrivateKeyBytes = pemBlock.Bytes
			}
		} else if certificates, err := x509.ParseCertificates(pemBlock.Bytes); err == nil {
			certPath = append(certPath, &CertBlock{
				Certificate: certificates[0],
				Bytes:       pemBlock.Bytes,
			})
		}
	}

	for i, certBlock := range certPath {
		if i == 0 {
			parsedBundle.Certificate = certBlock.Certificate
			parsedBundle.CertificateBytes = certBlock.Bytes
		} else {
			parsedBundle.CAChain = append(parsedBundle.CAChain, certBlock)
		}
	}

	if err := parsedBundle.Verify(); err != nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("verification of parsed bundle failed: %s", err)}
	}

	return parsedBundle, nil
}

// GeneratePrivateKey generates a private key with the specified type and key bits
func GeneratePrivateKey(keyType string, keyBits int, container ParsedPrivateKeyContainer) error {
	var err error
	var privateKeyType PrivateKeyType
	var privateKeyBytes []byte
	var privateKey crypto.Signer

	switch keyType {
	case "rsa":
		privateKeyType = RSAPrivateKey
		privateKey, err = rsa.GenerateKey(rand.Reader, keyBits)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error generating RSA private key: %v", err)}
		}
		privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey))
	case "ec":
		privateKeyType = ECPrivateKey
		var curve elliptic.Curve
		switch keyBits {
		case 224:
			curve = elliptic.P224()
		case 256:
			curve = elliptic.P256()
		case 384:
			curve = elliptic.P384()
		case 521:
			curve = elliptic.P521()
		default:
			return errutil.UserError{Err: fmt.Sprintf("unsupported bit length for EC key: %d", keyBits)}
		}
		privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error generating EC private key: %v", err)}
		}
		privateKeyBytes, err = x509.MarshalECPrivateKey(privateKey.(*ecdsa.PrivateKey))
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error marshalling EC private key: %v", err)}
		}
	default:
		return errutil.UserError{Err: fmt.Sprintf("unknown key type: %s", keyType)}
	}

	container.SetParsedPrivateKey(privateKey, privateKeyType, privateKeyBytes)
	return nil
}

// GenerateSerialNumber generates a serial number suitable for a certificate
func GenerateSerialNumber() (*big.Int, error) {
	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error generating serial number: %v", err)}
	}
	return serial, nil
}

// ComparePublicKeys compares two public keys and returns true if they match
func ComparePublicKeys(key1Iface, key2Iface crypto.PublicKey) (bool, error) {
	switch key1Iface.(type) {
	case *rsa.PublicKey:
		key1 := key1Iface.(*rsa.PublicKey)
		key2, ok := key2Iface.(*rsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("key types do not match: %T and %T", key1Iface, key2Iface)
		}
		if key1.N.Cmp(key2.N) != 0 ||
			key1.E != key2.E {
			return false, nil
		}
		return true, nil

	case *ecdsa.PublicKey:
		key1 := key1Iface.(*ecdsa.PublicKey)
		key2, ok := key2Iface.(*ecdsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("key types do not match: %T and %T", key1Iface, key2Iface)
		}
		if key1.X.Cmp(key2.X) != 0 ||
			key1.Y.Cmp(key2.Y) != 0 {
			return false, nil
		}
		key1Params := key1.Params()
		key2Params := key2.Params()
		if key1Params.P.Cmp(key2Params.P) != 0 ||
			key1Params.N.Cmp(key2Params.N) != 0 ||
			key1Params.B.Cmp(key2Params.B) != 0 ||
			key1Params.Gx.Cmp(key2Params.Gx) != 0 ||
			key1Params.Gy.Cmp(key2Params.Gy) != 0 ||
			key1Params.BitSize != key2Params.BitSize {
			return false, nil
		}
		return true, nil

	default:
		return false, fmt.Errorf("cannot compare key with type %T", key1Iface)
	}
}

// ParsePublicKeyPEM is used to parse RSA and ECDSA public keys from PEMs
func ParsePublicKeyPEM(data []byte) (interface{}, error) {
	block, data := pem.Decode(data)
	if block != nil {
		var rawKey interface{}
		var err error
		if rawKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				rawKey = cert.PublicKey
			} else {
				return nil, err
			}
		}

		if rsaPublicKey, ok := rawKey.(*rsa.PublicKey); ok {
			return rsaPublicKey, nil
		}
		if ecPublicKey, ok := rawKey.(*ecdsa.PublicKey); ok {
			return ecPublicKey, nil
		}
	}

	return nil, errors.New("data does not contain any valid RSA or ECDSA public keys")
}

// addPolicyIdentifiers adds certificate policies extension
//
func AddPolicyIdentifiers(data *CreationBundle, certTemplate *x509.Certificate) {
	for _, oidstr := range data.Params.PolicyIdentifiers {
		oid, err := StringToOid(oidstr)
		if err == nil {
			certTemplate.PolicyIdentifiers = append(certTemplate.PolicyIdentifiers, oid)
		}
	}
}

// addExtKeyUsageOids adds custom extended key usage OIDs to certificate
func AddExtKeyUsageOids(data *CreationBundle, certTemplate *x509.Certificate) {
	for _, oidstr := range data.Params.ExtKeyUsageOIDs {
		oid, err := StringToOid(oidstr)
		if err == nil {
			certTemplate.UnknownExtKeyUsage = append(certTemplate.UnknownExtKeyUsage, oid)
		}
	}
}

func HandleOtherCSRSANs(in *x509.CertificateRequest, sans map[string][]string) error {
	certTemplate := &x509.Certificate{
		DNSNames:       in.DNSNames,
		IPAddresses:    in.IPAddresses,
		EmailAddresses: in.EmailAddresses,
		URIs:           in.URIs,
	}
	if err := HandleOtherSANs(certTemplate, sans); err != nil {
		return err
	}
	if len(certTemplate.ExtraExtensions) > 0 {
		for _, v := range certTemplate.ExtraExtensions {
			in.ExtraExtensions = append(in.ExtraExtensions, v)
		}
	}
	return nil
}

func HandleOtherSANs(in *x509.Certificate, sans map[string][]string) error {
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
			oidStr, err := StringToOid(oid)
			if err != nil {
				return err
			}
			b.AddASN1ObjectIdentifier(oidStr)
			b.AddASN1(cbasn1.Tag(0).ContextSpecific().Constructed(), func(b *cryptobyte.Builder) {
				b.AddASN1(cbasn1.UTF8String, func(b *cryptobyte.Builder) {
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

func StringToOid(in string) (asn1.ObjectIdentifier, error) {
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

func ValidateKeyTypeLength(keyType string, keyBits int) error {
	switch keyType {
	case "rsa":
		switch keyBits {
		case 2048:
		case 4096:
		case 8192:
		default:
			return fmt.Errorf("unsupported bit length for RSA key: %d", keyBits)
		}
	case "ec":
		switch keyBits {
		case 224:
		case 256:
		case 384:
		case 521:
		default:
			return fmt.Errorf("unsupported bit length for EC key: %d", keyBits)
		}
	case "any":
	default:
		return fmt.Errorf("unknown key type %s", keyType)
	}

	return nil
}

// Performs the heavy lifting of creating a certificate. Returns
// a fully-filled-in ParsedCertBundle.
func CreateCertificate(data *CreationBundle) (*ParsedCertBundle, error) {
	var err error
	result := &ParsedCertBundle{}

	serialNumber, err := GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	if err := GeneratePrivateKey(data.Params.KeyType,
		data.Params.KeyBits,
		result); err != nil {
		return nil, err
	}

	subjKeyID, err := GetSubjKeyID(result.PrivateKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error getting subject key ID: %s", err)}
	}

	certTemplate := &x509.Certificate{
		SerialNumber:   serialNumber,
		NotBefore:      time.Now().Add(-30 * time.Second),
		NotAfter:       data.Params.NotAfter,
		IsCA:           false,
		SubjectKeyId:   subjKeyID,
		Subject:        data.Params.Subject,
		DNSNames:       data.Params.DNSNames,
		EmailAddresses: data.Params.EmailAddresses,
		IPAddresses:    data.Params.IPAddresses,
		URIs:           data.Params.URIs,
	}
	if data.Params.NotBeforeDuration > 0 {
		certTemplate.NotBefore = time.Now().Add(-1 * data.Params.NotBeforeDuration)
	}

	if err := HandleOtherSANs(certTemplate, data.Params.OtherSANs); err != nil {
		return nil, errutil.InternalError{Err: errwrap.Wrapf("error marshaling other SANs: {{err}}", err).Error()}
	}

	// Add this before calling addKeyUsages
	if data.SigningBundle == nil {
		certTemplate.IsCA = true
	} else if data.Params.BasicConstraintsValidForNonCA {
		certTemplate.BasicConstraintsValid = true
		certTemplate.IsCA = false
	}

	// This will only be filled in from the generation paths
	if len(data.Params.PermittedDNSDomains) > 0 {
		certTemplate.PermittedDNSDomains = data.Params.PermittedDNSDomains
		certTemplate.PermittedDNSDomainsCritical = true
	}

	AddPolicyIdentifiers(data, certTemplate)

	AddKeyUsages(data, certTemplate)

	AddExtKeyUsageOids(data, certTemplate)

	certTemplate.IssuingCertificateURL = data.Params.URLs.IssuingCertificates
	certTemplate.CRLDistributionPoints = data.Params.URLs.CRLDistributionPoints
	certTemplate.OCSPServer = data.Params.URLs.OCSPServers

	var certBytes []byte
	if data.SigningBundle != nil {
		switch data.SigningBundle.PrivateKeyType {
		case RSAPrivateKey:
			certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
		case ECPrivateKey:
			certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
		}

		caCert := data.SigningBundle.Certificate
		certTemplate.AuthorityKeyId = caCert.SubjectKeyId

		certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, result.PrivateKey.Public(), data.SigningBundle.PrivateKey)
	} else {
		// Creating a self-signed root
		if data.Params.MaxPathLength == 0 {
			certTemplate.MaxPathLen = 0
			certTemplate.MaxPathLenZero = true
		} else {
			certTemplate.MaxPathLen = data.Params.MaxPathLength
		}

		switch data.Params.KeyType {
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

	if data.SigningBundle != nil {
		if len(data.SigningBundle.Certificate.AuthorityKeyId) > 0 &&
			!bytes.Equal(data.SigningBundle.Certificate.AuthorityKeyId, data.SigningBundle.Certificate.SubjectKeyId) {

			result.CAChain = []*CertBlock{
				&CertBlock{
					Certificate: data.SigningBundle.Certificate,
					Bytes:       data.SigningBundle.CertificateBytes,
				},
			}
			result.CAChain = append(result.CAChain, data.SigningBundle.CAChain...)
		}
	}

	return result, nil
}

var oidExtensionBasicConstraints = []int{2, 5, 29, 19}

// Creates a CSR. This is currently only meant for use when
// generating an intermediate certificate.
func CreateCSR(data *CreationBundle, addBasicConstraints bool) (*ParsedCSRBundle, error) {
	var err error
	result := &ParsedCSRBundle{}

	if err := GeneratePrivateKey(data.Params.KeyType,
		data.Params.KeyBits,
		result); err != nil {
		return nil, err
	}

	// Like many root CAs, other information is ignored
	csrTemplate := &x509.CertificateRequest{
		Subject:        data.Params.Subject,
		DNSNames:       data.Params.DNSNames,
		EmailAddresses: data.Params.EmailAddresses,
		IPAddresses:    data.Params.IPAddresses,
		URIs:           data.Params.URIs,
	}

	if err := HandleOtherCSRSANs(csrTemplate, data.Params.OtherSANs); err != nil {
		return nil, errutil.InternalError{Err: errwrap.Wrapf("error marshaling other SANs: {{err}}", err).Error()}
	}

	if addBasicConstraints {
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

	switch data.Params.KeyType {
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
func SignCertificate(data *CreationBundle) (*ParsedCertBundle, error) {
	switch {
	case data == nil:
		return nil, errutil.UserError{Err: "nil data bundle given to signCertificate"}
	case data.Params == nil:
		return nil, errutil.UserError{Err: "nil parameters given to signCertificate"}
	case data.SigningBundle == nil:
		return nil, errutil.UserError{Err: "nil signing bundle given to signCertificate"}
	case data.CSR == nil:
		return nil, errutil.UserError{Err: "nil csr given to signCertificate"}
	}

	err := data.CSR.CheckSignature()
	if err != nil {
		return nil, errutil.UserError{Err: "request signature invalid"}
	}

	result := &ParsedCertBundle{}

	serialNumber, err := GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(data.CSR.PublicKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error marshalling public key: %s", err)}
	}
	subjKeyID := sha1.Sum(marshaledKey)

	caCert := data.SigningBundle.Certificate

	certTemplate := &x509.Certificate{
		SerialNumber:   serialNumber,
		Subject:        data.Params.Subject,
		NotBefore:      time.Now().Add(-30 * time.Second),
		NotAfter:       data.Params.NotAfter,
		SubjectKeyId:   subjKeyID[:],
		AuthorityKeyId: caCert.SubjectKeyId,
	}
	if data.Params.NotBeforeDuration > 0 {
		certTemplate.NotBefore = time.Now().Add(-1 * data.Params.NotBeforeDuration)
	}

	switch data.SigningBundle.PrivateKeyType {
	case RSAPrivateKey:
		certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
	case ECPrivateKey:
		certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
	}

	if data.Params.UseCSRValues {
		certTemplate.Subject = data.CSR.Subject
		certTemplate.Subject.ExtraNames = certTemplate.Subject.Names

		certTemplate.DNSNames = data.CSR.DNSNames
		certTemplate.EmailAddresses = data.CSR.EmailAddresses
		certTemplate.IPAddresses = data.CSR.IPAddresses
		certTemplate.URIs = data.CSR.URIs

		for _, name := range data.CSR.Extensions {
			if !name.Id.Equal(oidExtensionBasicConstraints) {
				certTemplate.ExtraExtensions = append(certTemplate.ExtraExtensions, name)
			}
		}

	} else {
		certTemplate.DNSNames = data.Params.DNSNames
		certTemplate.EmailAddresses = data.Params.EmailAddresses
		certTemplate.IPAddresses = data.Params.IPAddresses
		certTemplate.URIs = data.Params.URIs
	}

	if err := HandleOtherSANs(certTemplate, data.Params.OtherSANs); err != nil {
		return nil, errutil.InternalError{Err: errwrap.Wrapf("error marshaling other SANs: {{err}}", err).Error()}
	}

	AddPolicyIdentifiers(data, certTemplate)

	AddKeyUsages(data, certTemplate)

	AddExtKeyUsageOids(data, certTemplate)

	var certBytes []byte

	certTemplate.IssuingCertificateURL = data.Params.URLs.IssuingCertificates
	certTemplate.CRLDistributionPoints = data.Params.URLs.CRLDistributionPoints
	certTemplate.OCSPServer = data.SigningBundle.URLs.OCSPServers

	if data.Params.IsCA {
		certTemplate.BasicConstraintsValid = true
		certTemplate.IsCA = true

		if data.SigningBundle.Certificate.MaxPathLen == 0 &&
			data.SigningBundle.Certificate.MaxPathLenZero {
			return nil, errutil.UserError{Err: "signing certificate has a max path length of zero, and cannot issue further CA certificates"}
		}

		certTemplate.MaxPathLen = data.Params.MaxPathLength
		if certTemplate.MaxPathLen == 0 {
			certTemplate.MaxPathLenZero = true
		}
	} else if data.Params.BasicConstraintsValidForNonCA {
		certTemplate.BasicConstraintsValid = true
		certTemplate.IsCA = false
	}

	if len(data.Params.PermittedDNSDomains) > 0 {
		certTemplate.PermittedDNSDomains = data.Params.PermittedDNSDomains
		certTemplate.PermittedDNSDomainsCritical = true
	}

	certBytes, err = x509.CreateCertificate(rand.Reader, certTemplate, caCert, data.CSR.PublicKey, data.SigningBundle.PrivateKey)

	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CertificateBytes = certBytes
	result.Certificate, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	result.CAChain = data.SigningBundle.GetCAChain()

	return result, nil
}
