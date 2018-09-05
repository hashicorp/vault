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
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/mitchellh/mapstructure"
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

// PasrsePublicKeyPEM is used to parse RSA and ECDSA public keys from PEMs
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
