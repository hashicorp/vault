// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

const rsaMinimumSecureKeySize = 2048

// Mapping of key types to default key lengths
var defaultAlgorithmKeyBits = map[string]int{
	"rsa": 2048,
	"ec":  256,
}

// Mapping of NIST P-Curve's key length to expected signature bits.
var expectedNISTPCurveHashBits = map[int]int{
	224: 256,
	256: 256,
	384: 384,
	521: 512,
}

// Mapping of constant names<->constant values for SignatureAlgorithm
var SignatureAlgorithmNames = map[string]x509.SignatureAlgorithm{
	"sha256withrsa":    x509.SHA256WithRSA,
	"sha384withrsa":    x509.SHA384WithRSA,
	"sha512withrsa":    x509.SHA512WithRSA,
	"ecdsawithsha256":  x509.ECDSAWithSHA256,
	"ecdsawithsha384":  x509.ECDSAWithSHA384,
	"ecdsawithsha512":  x509.ECDSAWithSHA512,
	"sha256withrsapss": x509.SHA256WithRSAPSS,
	"sha384withrsapss": x509.SHA384WithRSAPSS,
	"sha512withrsapss": x509.SHA512WithRSAPSS,
	"pureed25519":      x509.PureEd25519,
	"ed25519":          x509.PureEd25519, // Duplicated for clarity; most won't expect the "Pure" prefix.
}

// Mapping of constant values<->constant names for SignatureAlgorithm
var InvSignatureAlgorithmNames = map[x509.SignatureAlgorithm]string{
	x509.SHA256WithRSA:    "SHA256WithRSA",
	x509.SHA384WithRSA:    "SHA384WithRSA",
	x509.SHA512WithRSA:    "SHA512WithRSA",
	x509.ECDSAWithSHA256:  "ECDSAWithSHA256",
	x509.ECDSAWithSHA384:  "ECDSAWithSHA384",
	x509.ECDSAWithSHA512:  "ECDSAWithSHA512",
	x509.SHA256WithRSAPSS: "SHA256WithRSAPSS",
	x509.SHA384WithRSAPSS: "SHA384WithRSAPSS",
	x509.SHA512WithRSAPSS: "SHA512WithRSAPSS",
	x509.PureEd25519:      "Ed25519",
}

// OID for RFC 5280 CRL Number extension.
//
// > id-ce-cRLNumber OBJECT IDENTIFIER ::= { id-ce 20 }
var CRLNumberOID = asn1.ObjectIdentifier([]int{2, 5, 29, 20})

// OID for RFC 5280 Delta CRL Indicator CRL extension.
//
// > id-ce-deltaCRLIndicator OBJECT IDENTIFIER ::= { id-ce 27 }
var DeltaCRLIndicatorOID = asn1.ObjectIdentifier([]int{2, 5, 29, 27})

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
	var inBits uint64
	inBytes := strings.Split(in, sep)
	for _, inByte := range inBytes {
		if inBits, err = strconv.ParseUint(inByte, 16, 8); err != nil {
			return nil
		}
		ret.WriteByte(uint8(inBits))
	}
	return ret.Bytes()
}

// GetSubjKeyID returns the subject key ID. The computed ID is the SHA-1 hash of
// the marshaled public key according to
// https://tools.ietf.org/html/rfc5280#section-4.2.1.2 (1)
func GetSubjKeyID(privateKey crypto.Signer) ([]byte, error) {
	if privateKey == nil {
		return nil, errutil.InternalError{Err: "passed-in private key is nil"}
	}
	return GetSubjectKeyID(privateKey.Public())
}

// Returns the explicit SKID when used for cross-signing, else computes a new
// SKID from the key itself.
func getSubjectKeyIDFromBundle(data *CreationBundle) ([]byte, error) {
	if len(data.Params.SKID) > 0 {
		return data.Params.SKID, nil
	}

	return GetSubjectKeyID(data.CSR.PublicKey)
}

func GetSubjectKeyID(pub interface{}) ([]byte, error) {
	var publicKeyBytes []byte
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		type pkcs1PublicKey struct {
			N *big.Int
			E int
		}

		var err error
		publicKeyBytes, err = asn1.Marshal(pkcs1PublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("error marshalling public key: %s", err)}
		}
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	case ed25519.PublicKey:
		publicKeyBytes = pub
	default:
		return nil, errutil.InternalError{Err: fmt.Sprintf("unsupported public key type: %T", pub)}
	}
	skid := sha1.Sum(publicKeyBytes)
	return skid[:], nil
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

func ParseDERKey(privateKeyBytes []byte) (signer crypto.Signer, format BlockType, err error) {
	var firstError error
	if signer, firstError = x509.ParseECPrivateKey(privateKeyBytes); firstError == nil {
		format = ECBlock
		return
	}

	var secondError error
	if signer, secondError = x509.ParsePKCS1PrivateKey(privateKeyBytes); secondError == nil {
		format = PKCS1Block
		return
	}

	var thirdError error
	var rawKey interface{}
	if rawKey, thirdError = x509.ParsePKCS8PrivateKey(privateKeyBytes); thirdError == nil {
		switch rawSigner := rawKey.(type) {
		case *rsa.PrivateKey:
			signer = rawSigner
		case *ecdsa.PrivateKey:
			signer = rawSigner
		case ed25519.PrivateKey:
			signer = rawSigner
		default:
			return nil, UnknownBlock, errutil.InternalError{Err: "unknown type for parsed PKCS8 Private Key"}
		}

		format = PKCS8Block
		return
	}

	return nil, UnknownBlock, fmt.Errorf("got errors attempting to parse DER private key:\n1. %v\n2. %v\n3. %v", firstError, secondError, thirdError)
}

func ParsePEMKey(keyPem string) (crypto.Signer, BlockType, error) {
	pemBlock, _ := pem.Decode([]byte(keyPem))
	if pemBlock == nil {
		return nil, UnknownBlock, errutil.UserError{Err: "no data found in PEM block"}
	}

	return ParseDERKey(pemBlock.Bytes)
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

		if signer, format, err := ParseDERKey(pemBlock.Bytes); err == nil {
			if parsedBundle.PrivateKeyType != UnknownPrivateKey {
				return nil, errutil.UserError{Err: "more than one private key given; provide only one private key in the bundle"}
			}

			parsedBundle.PrivateKeyFormat = format
			parsedBundle.PrivateKeyType = GetPrivateKeyTypeFromSigner(signer)
			if parsedBundle.PrivateKeyType == UnknownPrivateKey {
				return nil, errutil.UserError{Err: "Unknown type of private key included in the bundle: %v"}
			}

			parsedBundle.PrivateKeyBytes = pemBlock.Bytes
			parsedBundle.PrivateKey = signer
		} else if certificates, err := x509.ParseCertificates(pemBlock.Bytes); err == nil {
			certPath = append(certPath, &CertBlock{
				Certificate: certificates[0],
				Bytes:       pemBlock.Bytes,
			})
		} else if x509.IsEncryptedPEMBlock(pemBlock) {
			return nil, errutil.UserError{Err: "Encrypted private key given; provide only decrypted private key in the bundle"}
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

func (p *ParsedCertBundle) ToTLSCertificate() tls.Certificate {
	var cert tls.Certificate
	cert.Certificate = append(cert.Certificate, p.CertificateBytes)
	cert.Leaf = p.Certificate
	cert.PrivateKey = p.PrivateKey
	for _, ca := range p.CAChain {
		cert.Certificate = append(cert.Certificate, ca.Bytes)
	}

	return cert
}

// GeneratePrivateKey generates a private key with the specified type and key bits.
func GeneratePrivateKey(keyType string, keyBits int, container ParsedPrivateKeyContainer) error {
	return generatePrivateKey(keyType, keyBits, container, nil)
}

// GeneratePrivateKeyWithRandomSource generates a private key with the specified type and key bits.
// GeneratePrivateKeyWithRandomSource uses randomness from the entropyReader to generate the private key.
func GeneratePrivateKeyWithRandomSource(keyType string, keyBits int, container ParsedPrivateKeyContainer, entropyReader io.Reader) error {
	return generatePrivateKey(keyType, keyBits, container, entropyReader)
}

// generatePrivateKey generates a private key with the specified type and key bits.
// generatePrivateKey uses randomness from the entropyReader to generate the private key.
func generatePrivateKey(keyType string, keyBits int, container ParsedPrivateKeyContainer, entropyReader io.Reader) error {
	var err error
	var privateKeyType PrivateKeyType
	var privateKeyBytes []byte
	var privateKey crypto.Signer

	var randReader io.Reader = rand.Reader
	if entropyReader != nil {
		randReader = entropyReader
	}

	switch keyType {
	case "rsa":
		// XXX: there is a false-positive CodeQL path here around keyBits;
		// because of a default zero value in the TypeDurationSecond and
		// TypeSignedDurationSecond cases of schema.DefaultOrZero(), it
		// thinks it is possible to end up with < 2048 bit RSA Key here.
		// While this is true for SSH keys, it isn't true for PKI keys
		// due to ValidateKeyTypeLength(...) below. While we could close
		// the report as a false-positive, enforcing a minimum keyBits size
		// here of 2048 would ensure no other paths exist.
		if keyBits < 2048 {
			return errutil.InternalError{Err: fmt.Sprintf("insecure bit length for RSA private key: %d", keyBits)}
		}
		privateKeyType = RSAPrivateKey
		privateKey, err = rsa.GenerateKey(randReader, keyBits)
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
		privateKey, err = ecdsa.GenerateKey(curve, randReader)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error generating EC private key: %v", err)}
		}
		privateKeyBytes, err = x509.MarshalECPrivateKey(privateKey.(*ecdsa.PrivateKey))
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error marshalling EC private key: %v", err)}
		}
	case "ed25519":
		privateKeyType = Ed25519PrivateKey
		_, privateKey, err = ed25519.GenerateKey(randReader)
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error generating ed25519 private key: %v", err)}
		}
		privateKeyBytes, err = x509.MarshalPKCS8PrivateKey(privateKey.(ed25519.PrivateKey))
		if err != nil {
			return errutil.InternalError{Err: fmt.Sprintf("error marshalling Ed25519 private key: %v", err)}
		}
	default:
		return errutil.UserError{Err: fmt.Sprintf("unknown key type: %s", keyType)}
	}

	container.SetParsedPrivateKey(privateKey, privateKeyType, privateKeyBytes)
	return nil
}

// GenerateSerialNumber generates a serial number suitable for a certificate
func GenerateSerialNumber() (*big.Int, error) {
	return generateSerialNumber(rand.Reader)
}

// GenerateSerialNumberWithRandomSource generates a serial number suitable
// for a certificate with custom entropy.
func GenerateSerialNumberWithRandomSource(randReader io.Reader) (*big.Int, error) {
	return generateSerialNumber(randReader)
}

func generateSerialNumber(randReader io.Reader) (*big.Int, error) {
	serial, err := rand.Int(randReader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error generating serial number: %v", err)}
	}
	return serial, nil
}

// ComparePublicKeysAndType compares two public keys and returns true if they match,
// false if their types or contents differ, and an error on unsupported key types.
func ComparePublicKeysAndType(key1Iface, key2Iface crypto.PublicKey) (bool, error) {
	equal, err := ComparePublicKeys(key1Iface, key2Iface)
	if err != nil {
		if strings.Contains(err.Error(), "key types do not match:") {
			return false, nil
		}
	}

	return equal, err
}

// ComparePublicKeys compares two public keys and returns true if they match,
// returns an error if public key types are mismatched, or they are an unsupported key type.
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
	case ed25519.PublicKey:
		key1 := key1Iface.(ed25519.PublicKey)
		key2, ok := key2Iface.(ed25519.PublicKey)
		if !ok {
			return false, fmt.Errorf("key types do not match: %T and %T", key1Iface, key2Iface)
		}
		if !key1.Equal(key2) {
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
		if len(bytes.TrimSpace(data)) > 0 {
			return nil, errutil.UserError{Err: "unexpected trailing data after parsed PEM block"}
		}
		var rawKey interface{}
		var err error
		if rawKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				rawKey = cert.PublicKey
			} else {
				return nil, err
			}
		}

		switch key := rawKey.(type) {
		case *rsa.PublicKey:
			return key, nil
		case *ecdsa.PublicKey:
			return key, nil
		case ed25519.PublicKey:
			return key, nil
		}
	}
	return nil, errors.New("data does not contain any valid public keys")
}

// AddPolicyIdentifiers adds certificate policies extension, based on CreationBundle
func AddPolicyIdentifiers(data *CreationBundle, certTemplate *x509.Certificate) {
	oidOnly := true
	for _, oidStr := range data.Params.PolicyIdentifiers {
		oid, err := StringToOid(oidStr)
		if err == nil {
			certTemplate.PolicyIdentifiers = append(certTemplate.PolicyIdentifiers, oid)
		}
		if err != nil {
			oidOnly = false
		}
	}
	if !oidOnly { // Because all policy information is held in the same extension, when we use an extra extension to
		// add policy qualifier information, that overwrites any information in the PolicyIdentifiers field on the Cert
		// Template, so we need to reparse all the policy identifiers here
		extension, err := CreatePolicyInformationExtensionFromStorageStrings(data.Params.PolicyIdentifiers)
		if err == nil {
			// If this errors out, don't add it, rely on the OIDs parsed into PolicyIdentifiers above
			certTemplate.ExtraExtensions = append(certTemplate.ExtraExtensions, *extension)
		}
	}
}

// AddExtKeyUsageOids adds custom extended key usage OIDs to certificate
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

// Returns default key bits for the specified key type, or the present value
// if keyBits is non-zero.
func DefaultOrValueKeyBits(keyType string, keyBits int) (int, error) {
	if keyBits == 0 {
		newValue, present := defaultAlgorithmKeyBits[keyType]
		if present {
			keyBits = newValue
		} /* else {
		  // We cannot return an error here as ed25519 (and potentially ed448
		  // in the future) aren't in defaultAlgorithmKeyBits -- the value of
		  // the keyBits parameter is ignored under that algorithm.
		} */
	}

	return keyBits, nil
}

// Returns default signature hash bit length for the specified key type and
// bits, or the present value if hashBits is non-zero. Returns an error under
// certain internal circumstances.
func DefaultOrValueHashBits(keyType string, keyBits int, hashBits int) (int, error) {
	if keyType == "ec" {
		// Enforcement of curve moved to selectSignatureAlgorithmForECDSA. See
		// note there about why.
	} else if keyType == "rsa" && hashBits == 0 {
		// To match previous behavior (and ignoring NIST's recommendations for
		// hash size to align with RSA key sizes), default to SHA-2-256.
		hashBits = 256
	} else if keyType == "ed25519" || keyType == "ed448" || keyType == "any" {
		// No-op; ed25519 and ed448 internally specify their own hash and
		// we do not need to select one. Double hashing isn't supported in
		// certificate signing. Additionally, the any key type can't know
		// what hash algorithm to use yet, so default to zero.
		return 0, nil
	}

	return hashBits, nil
}

// Validates that the combination of keyType, keyBits, and hashBits are
// valid together; replaces individual calls to ValidateSignatureLength and
// ValidateKeyTypeLength. Also updates the value of keyBits and hashBits on
// return.
func ValidateDefaultOrValueKeyTypeSignatureLength(keyType string, keyBits int, hashBits int) (int, int, error) {
	var err error

	if keyBits, err = DefaultOrValueKeyBits(keyType, keyBits); err != nil {
		return keyBits, hashBits, err
	}

	if err = ValidateKeyTypeLength(keyType, keyBits); err != nil {
		return keyBits, hashBits, err
	}

	if hashBits, err = DefaultOrValueHashBits(keyType, keyBits, hashBits); err != nil {
		return keyBits, hashBits, err
	}

	// Note that this check must come after we've selected a value for
	// hashBits above, in the event it was left as the default, but we
	// were allowed to update it.
	if err = ValidateSignatureLength(keyType, hashBits); err != nil {
		return keyBits, hashBits, err
	}

	return keyBits, hashBits, nil
}

// Validates that the length of the hash (in bits) used in the signature
// calculation is a known, approved value.
func ValidateSignatureLength(keyType string, hashBits int) error {
	if keyType == "any" || keyType == "ec" || keyType == "ed25519" || keyType == "ed448" {
		// ed25519 and ed448 include built-in hashing and is not externally
		// configurable. There are three modes for each of these schemes:
		//
		// 1. Built-in hash (default, used in TLS, x509).
		// 2. Double hash (notably used in some block-chain implementations,
		//    but largely regarded as a specialized use case with security
		//    concerns).
		// 3. No hash (bring your own hash function, less commonly used).
		//
		// In all cases, we won't have a hash algorithm to validate here, so
		// return nil.
		//
		// Additionally, when KeyType is any, we can't yet validate the
		// signature algorithm size, so it takes the default zero value.
		//
		// When KeyType is ec, we also can't validate this value as we're
		// forcefully ignoring the users' choice and specifying a value based
		// on issuer type.
		return nil
	}

	switch hashBits {
	case 256:
	case 384:
	case 512:
	default:
		return fmt.Errorf("unsupported hash signature algorithm: %d", hashBits)
	}

	return nil
}

func ValidateKeyTypeLength(keyType string, keyBits int) error {
	switch keyType {
	case "rsa":
		if keyBits < rsaMinimumSecureKeySize {
			return fmt.Errorf("RSA keys < %d bits are unsafe and not supported: got %d", rsaMinimumSecureKeySize, keyBits)
		}

		switch keyBits {
		case 2048:
		case 3072:
		case 4096:
		case 8192:
		default:
			return fmt.Errorf("unsupported bit length for RSA key: %d", keyBits)
		}
	case "ec":
		_, present := expectedNISTPCurveHashBits[keyBits]
		if !present {
			return fmt.Errorf("unsupported bit length for EC key: %d", keyBits)
		}
	case "any", "ed25519":
	default:
		return fmt.Errorf("unknown key type %s", keyType)
	}

	return nil
}

// CreateCertificate uses CreationBundle and the default rand.Reader to
// generate a cert/keypair.
func CreateCertificate(data *CreationBundle) (*ParsedCertBundle, error) {
	return createCertificate(data, rand.Reader, generatePrivateKey)
}

// CreateCertificateWithRandomSource uses CreationBundle and a custom
// io.Reader for randomness to generate a cert/keypair.
func CreateCertificateWithRandomSource(data *CreationBundle, randReader io.Reader) (*ParsedCertBundle, error) {
	return createCertificate(data, randReader, generatePrivateKey)
}

// KeyGenerator Allow us to override how/what generates the private key
type KeyGenerator func(keyType string, keyBits int, container ParsedPrivateKeyContainer, entropyReader io.Reader) error

func CreateCertificateWithKeyGenerator(data *CreationBundle, randReader io.Reader, keyGenerator KeyGenerator) (*ParsedCertBundle, error) {
	return createCertificate(data, randReader, keyGenerator)
}

// Set correct RSA sig algo
func certTemplateSetSigAlgo(certTemplate *x509.Certificate, data *CreationBundle) {
	if data.Params.UsePSS {
		switch data.Params.SignatureBits {
		case 256:
			certTemplate.SignatureAlgorithm = x509.SHA256WithRSAPSS
		case 384:
			certTemplate.SignatureAlgorithm = x509.SHA384WithRSAPSS
		case 512:
			certTemplate.SignatureAlgorithm = x509.SHA512WithRSAPSS
		}
	} else {
		switch data.Params.SignatureBits {
		case 256:
			certTemplate.SignatureAlgorithm = x509.SHA256WithRSA
		case 384:
			certTemplate.SignatureAlgorithm = x509.SHA384WithRSA
		case 512:
			certTemplate.SignatureAlgorithm = x509.SHA512WithRSA
		}
	}
}

// selectSignatureAlgorithmForRSA returns the proper x509.SignatureAlgorithm based on various properties set in the
// Creation Bundle parameter. This method will default to a SHA256 signature algorithm if the requested signature
// bits is not set/unknown.
func selectSignatureAlgorithmForRSA(data *CreationBundle) x509.SignatureAlgorithm {
	if data.Params.UsePSS {
		switch data.Params.SignatureBits {
		case 256:
			return x509.SHA256WithRSAPSS
		case 384:
			return x509.SHA384WithRSAPSS
		case 512:
			return x509.SHA512WithRSAPSS
		default:
			return x509.SHA256WithRSAPSS
		}
	}

	switch data.Params.SignatureBits {
	case 256:
		return x509.SHA256WithRSA
	case 384:
		return x509.SHA384WithRSA
	case 512:
		return x509.SHA512WithRSA
	default:
		return x509.SHA256WithRSA
	}
}

func createCertificate(data *CreationBundle, randReader io.Reader, privateKeyGenerator KeyGenerator) (*ParsedCertBundle, error) {
	var err error
	result := &ParsedCertBundle{}

	serialNumber, err := GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	if err := privateKeyGenerator(data.Params.KeyType,
		data.Params.KeyBits,
		result, randReader); err != nil {
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
		privateKeyType := data.SigningBundle.PrivateKeyType
		if privateKeyType == ManagedPrivateKey {
			privateKeyType = GetPrivateKeyTypeFromSigner(data.SigningBundle.PrivateKey)
		}
		switch privateKeyType {
		case RSAPrivateKey:
			certTemplateSetSigAlgo(certTemplate, data)
		case Ed25519PrivateKey:
			certTemplate.SignatureAlgorithm = x509.PureEd25519
		case ECPrivateKey:
			certTemplate.SignatureAlgorithm = selectSignatureAlgorithmForECDSA(data.SigningBundle.PrivateKey.Public(), data.Params.SignatureBits)
		}

		caCert := data.SigningBundle.Certificate
		certTemplate.AuthorityKeyId = caCert.SubjectKeyId

		certBytes, err = x509.CreateCertificate(randReader, certTemplate, caCert, result.PrivateKey.Public(), data.SigningBundle.PrivateKey)
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
			certTemplateSetSigAlgo(certTemplate, data)
		case "ed25519":
			certTemplate.SignatureAlgorithm = x509.PureEd25519
		case "ec":
			certTemplate.SignatureAlgorithm = selectSignatureAlgorithmForECDSA(result.PrivateKey.Public(), data.Params.SignatureBits)
		}

		certTemplate.AuthorityKeyId = subjKeyID
		certTemplate.BasicConstraintsValid = true
		certBytes, err = x509.CreateCertificate(randReader, certTemplate, certTemplate, result.PrivateKey.Public(), result.PrivateKey)
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
		if (len(data.SigningBundle.Certificate.AuthorityKeyId) > 0 &&
			!bytes.Equal(data.SigningBundle.Certificate.AuthorityKeyId, data.SigningBundle.Certificate.SubjectKeyId)) ||
			data.Params.ForceAppendCaChain {
			var chain []*CertBlock

			signingChain := data.SigningBundle.CAChain
			// Some bundles already include the root included in the chain, so don't include it twice.
			if len(signingChain) == 0 || !bytes.Equal(signingChain[0].Bytes, data.SigningBundle.CertificateBytes) {
				chain = append(chain, &CertBlock{
					Certificate: data.SigningBundle.Certificate,
					Bytes:       data.SigningBundle.CertificateBytes,
				})
			}

			if len(signingChain) > 0 {
				chain = append(chain, signingChain...)
			}

			result.CAChain = chain
		}
	}

	return result, nil
}

func selectSignatureAlgorithmForECDSA(pub crypto.PublicKey, signatureBits int) x509.SignatureAlgorithm {
	// Previously we preferred the user-specified signature bits for ECDSA
	// keys. However, this could result in using a longer hash function than
	// the underlying NIST P-curve will encode (e.g., a SHA-512 hash with a
	// P-256 key). This isn't ideal: the hash is implicitly truncated
	// (effectively turning it into SHA-512/256) and we then need to rely
	// on the prefix security of the hash. Since both NIST and Mozilla guidance
	// suggest instead using the correct hash function, we should prefer that
	// over the operator-specified signatureBits.
	//
	// Lastly, note that pub above needs to be the _signer's_ public key;
	// the issue with DefaultOrValueHashBits is that it is called at role
	// configuration time, which might _precede_ issuer generation. Thus
	// it only has access to the desired key type and not the actual issuer.
	// The reference from that function is reproduced below:
	//
	// > To comply with BSI recommendations Section 4.2 and Mozilla root
	// > store policy section 5.1.2, enforce that NIST P-curves use a hash
	// > length corresponding to curve length. Note that ed25519 does not
	// > implement the "ec" key type.
	key, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return x509.ECDSAWithSHA256
	}
	switch key.Curve {
	case elliptic.P224(), elliptic.P256():
		return x509.ECDSAWithSHA256
	case elliptic.P384():
		return x509.ECDSAWithSHA384
	case elliptic.P521():
		return x509.ECDSAWithSHA512
	default:
		return x509.ECDSAWithSHA256
	}
}

var (
	ExtensionBasicConstraintsOID = []int{2, 5, 29, 19}
	ExtensionSubjectAltNameOID   = []int{2, 5, 29, 17}
)

// CreateCSR creates a CSR with the default rand.Reader to
// generate a cert/keypair. This is currently only meant
// for use when generating an intermediate certificate.
func CreateCSR(data *CreationBundle, addBasicConstraints bool) (*ParsedCSRBundle, error) {
	return createCSR(data, addBasicConstraints, rand.Reader, generatePrivateKey)
}

// CreateCSRWithRandomSource creates a CSR with a custom io.Reader
// for randomness to generate a cert/keypair.
func CreateCSRWithRandomSource(data *CreationBundle, addBasicConstraints bool, randReader io.Reader) (*ParsedCSRBundle, error) {
	return createCSR(data, addBasicConstraints, randReader, generatePrivateKey)
}

// CreateCSRWithKeyGenerator creates a CSR with a custom io.Reader
// for randomness to generate a cert/keypair with the provided private key generator.
func CreateCSRWithKeyGenerator(data *CreationBundle, addBasicConstraints bool, randReader io.Reader, keyGenerator KeyGenerator) (*ParsedCSRBundle, error) {
	return createCSR(data, addBasicConstraints, randReader, keyGenerator)
}

func createCSR(data *CreationBundle, addBasicConstraints bool, randReader io.Reader, keyGenerator KeyGenerator) (*ParsedCSRBundle, error) {
	var err error
	result := &ParsedCSRBundle{}

	if err := keyGenerator(data.Params.KeyType,
		data.Params.KeyBits,
		result, randReader); err != nil {
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
			Id:       ExtensionBasicConstraintsOID,
			Value:    val,
			Critical: true,
		}
		csrTemplate.ExtraExtensions = append(csrTemplate.ExtraExtensions, ext)
	}

	switch data.Params.KeyType {
	case "rsa":
		// use specified RSA algorithm defaulting to the appropriate SHA256 RSA signature type
		csrTemplate.SignatureAlgorithm = selectSignatureAlgorithmForRSA(data)
	case "ec":
		csrTemplate.SignatureAlgorithm = selectSignatureAlgorithmForECDSA(result.PrivateKey.Public(), data.Params.SignatureBits)
	case "ed25519":
		csrTemplate.SignatureAlgorithm = x509.PureEd25519
	}

	csr, err := x509.CreateCertificateRequest(randReader, csrTemplate, result.PrivateKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CSRBytes = csr
	result.CSR, err = x509.ParseCertificateRequest(csr)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %v", err)}
	}

	if err = result.CSR.CheckSignature(); err != nil {
		return nil, errors.New("failed signature validation for CSR")
	}

	return result, nil
}

// SignCertificate performs the heavy lifting
// of generating a certificate from a CSR.
// Returns a ParsedCertBundle sans private keys.
func SignCertificate(data *CreationBundle) (*ParsedCertBundle, error) {
	return signCertificate(data, rand.Reader)
}

// SignCertificateWithRandomSource generates a certificate
// from a CSR, using custom randomness from the randReader.
// Returns a ParsedCertBundle sans private keys.
func SignCertificateWithRandomSource(data *CreationBundle, randReader io.Reader) (*ParsedCertBundle, error) {
	return signCertificate(data, randReader)
}

func signCertificate(data *CreationBundle, randReader io.Reader) (*ParsedCertBundle, error) {
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

	subjKeyID, err := getSubjectKeyIDFromBundle(data)
	if err != nil {
		return nil, err
	}

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

	privateKeyType := data.SigningBundle.PrivateKeyType
	if privateKeyType == ManagedPrivateKey {
		privateKeyType = GetPrivateKeyTypeFromSigner(data.SigningBundle.PrivateKey)
	}

	switch privateKeyType {
	case RSAPrivateKey:
		certTemplateSetSigAlgo(certTemplate, data)
	case ECPrivateKey:
		switch data.Params.SignatureBits {
		case 256:
			certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA256
		case 384:
			certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA384
		case 512:
			certTemplate.SignatureAlgorithm = x509.ECDSAWithSHA512
		}
	}

	if data.Params.UseCSRValues {
		certTemplate.Subject = data.CSR.Subject
		certTemplate.Subject.ExtraNames = certTemplate.Subject.Names

		certTemplate.DNSNames = data.CSR.DNSNames
		certTemplate.EmailAddresses = data.CSR.EmailAddresses
		certTemplate.IPAddresses = data.CSR.IPAddresses
		certTemplate.URIs = data.CSR.URIs

		for _, name := range data.CSR.Extensions {
			if !name.Id.Equal(ExtensionBasicConstraintsOID) && !(len(data.Params.OtherSANs) > 0 && name.Id.Equal(ExtensionSubjectAltNameOID)) {
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

	certBytes, err = x509.CreateCertificate(randReader, certTemplate, caCert, data.CSR.PublicKey, data.SigningBundle.PrivateKey)

	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to create certificate: %s", err)}
	}

	result.CertificateBytes = certBytes
	result.Certificate, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse created certificate: %s", err)}
	}

	result.CAChain = data.SigningBundle.GetFullChain()

	return result, nil
}

func NewCertPool(reader io.Reader) (*x509.CertPool, error) {
	pemBlock, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	certs, err := ParseCertsPEM(pemBlock)
	if err != nil {
		return nil, fmt.Errorf("error reading certs: %s", err)
	}
	pool := x509.NewCertPool()
	for _, cert := range certs {
		pool.AddCert(cert)
	}
	return pool, nil
}

// ParseCertsPEM returns the x509.Certificates contained in the given PEM-encoded byte array
// Returns an error if a certificate could not be parsed, or if the data does not contain any certificates
func ParseCertsPEM(pemCerts []byte) ([]*x509.Certificate, error) {
	ok := false
	certs := []*x509.Certificate{}
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		// Only use PEM "CERTIFICATE" blocks without extra headers
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return certs, err
		}

		certs = append(certs, cert)
		ok = true
	}

	if !ok {
		return certs, errors.New("data does not contain any valid RSA or ECDSA certificates")
	}
	return certs, nil
}

// GetPublicKeySize returns the key size in bits for a given arbitrary crypto.PublicKey
// Returns -1 for an unsupported key type.
func GetPublicKeySize(key crypto.PublicKey) int {
	if key, ok := key.(*rsa.PublicKey); ok {
		return key.Size() * 8
	}
	if key, ok := key.(*ecdsa.PublicKey); ok {
		return key.Params().BitSize
	}
	if key, ok := key.(ed25519.PublicKey); ok {
		return len(key) * 8
	}
	if key, ok := key.(dsa.PublicKey); ok {
		return key.Y.BitLen()
	}

	return -1
}

// CreateKeyBundle create a KeyBundle struct object which includes a generated key
// of keyType with keyBits leveraging the randomness from randReader.
func CreateKeyBundle(keyType string, keyBits int, randReader io.Reader) (KeyBundle, error) {
	return CreateKeyBundleWithKeyGenerator(keyType, keyBits, randReader, generatePrivateKey)
}

// CreateKeyBundleWithKeyGenerator create a KeyBundle struct object which includes
// a generated key of keyType with keyBits leveraging the randomness from randReader and
// delegates the actual key generation to keyGenerator
func CreateKeyBundleWithKeyGenerator(keyType string, keyBits int, randReader io.Reader, keyGenerator KeyGenerator) (KeyBundle, error) {
	result := KeyBundle{}
	if err := keyGenerator(keyType, keyBits, &result, randReader); err != nil {
		return result, err
	}
	return result, nil
}

// CreateDeltaCRLIndicatorExt allows creating correctly formed delta CRLs
// that point back to the last complete CRL that they're based on.
func CreateDeltaCRLIndicatorExt(completeCRLNumber int64) (pkix.Extension, error) {
	bigNum := big.NewInt(completeCRLNumber)
	bigNumValue, err := asn1.Marshal(bigNum)
	if err != nil {
		return pkix.Extension{}, fmt.Errorf("unable to marshal complete CRL number (%v): %v", completeCRLNumber, err)
	}
	return pkix.Extension{
		Id: DeltaCRLIndicatorOID,
		// > When a conforming CRL issuer generates a delta CRL, the delta
		// > CRL MUST include a critical delta CRL indicator extension.
		Critical: true,
		// This extension only includes the complete CRL number:
		//
		// > BaseCRLNumber ::= CRLNumber
		//
		// But, this needs to be encoded as a big number for encoding/asn1
		// to work properly.
		Value: bigNumValue,
	}, nil
}

// ParseBasicConstraintExtension parses a basic constraint pkix.Extension, useful if attempting to validate
// CSRs are requesting CA privileges as Go does not expose its implementation. Values returned are
// IsCA, MaxPathLen or error. If MaxPathLen was not set, a value of -1 will be returned.
func ParseBasicConstraintExtension(ext pkix.Extension) (bool, int, error) {
	if !ext.Id.Equal(ExtensionBasicConstraintsOID) {
		return false, -1, fmt.Errorf("passed in extension was not a basic constraint extension")
	}

	// All elements are set to optional here, as it is possible that we receive a CSR with the extension
	// containing an empty sequence by spec.
	type basicConstraints struct {
		IsCA       bool `asn1:"optional"`
		MaxPathLen int  `asn1:"optional,default:-1"`
	}
	bc := &basicConstraints{}
	leftOver, err := asn1.Unmarshal(ext.Value, bc)
	if err != nil {
		return false, -1, fmt.Errorf("failed unmarshalling extension value: %w", err)
	}

	numLeftOver := len(bytes.TrimSpace(leftOver))
	if numLeftOver > 0 {
		return false, -1, fmt.Errorf("%d extra bytes within basic constraints value extension", numLeftOver)
	}

	return bc.IsCA, bc.MaxPathLen, nil
}

// CreateBasicConstraintExtension create a basic constraint extension based on inputs,
// if isCa is false, an empty value sequence will be returned with maxPath being
// ignored. If isCa is true maxPath can be set to -1 to not set a maxPath value.
func CreateBasicConstraintExtension(isCa bool, maxPath int) (pkix.Extension, error) {
	var asn1Bytes []byte
	var err error

	switch {
	case isCa && maxPath >= 0:
		CaAndMaxPathLen := struct {
			IsCa       bool `asn1:""`
			MaxPathLen int  `asn1:""`
		}{
			IsCa:       isCa,
			MaxPathLen: maxPath,
		}
		asn1Bytes, err = asn1.Marshal(CaAndMaxPathLen)
	case isCa && maxPath < 0:
		justCa := struct {
			IsCa bool `asn1:""`
		}{IsCa: isCa}
		asn1Bytes, err = asn1.Marshal(justCa)
	default:
		asn1Bytes, err = asn1.Marshal(struct{}{})
	}

	if err != nil {
		return pkix.Extension{}, err
	}

	return pkix.Extension{
		Id:       ExtensionBasicConstraintsOID,
		Critical: true,
		Value:    asn1Bytes,
	}, nil
}
