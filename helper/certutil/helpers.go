package certutil

import (
	"bytes"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// GetOctalFormatted returns the byte buffer formatted in octal with
// the specified separator between bytes.
func GetOctalFormatted(buf []byte, sep string) string {
	var ret bytes.Buffer
	for _, cur := range buf {
		if ret.Len() > 0 {
			fmt.Fprintf(&ret, sep)
		}
		fmt.Fprintf(&ret, "%02x", cur)
	}
	return ret.String()
}

// GetSubjKeyID returns the subject key ID, e.g. the SHA1 sum
// of the marshaled public key
func GetSubjKeyID(privateKey crypto.Signer) ([]byte, error) {
	if privateKey == nil {
		return nil, InternalError{"Passed-in private key is nil"}
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, InternalError{fmt.Sprintf("Error marshalling public key: %s", err)}
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
		return nil, UserError{err.Error()}
	}

	return result.ToParsedCertBundle()
}

// ParsePKIJSON takes a JSON-encoded string and returns a CertBundle
// ParsedCertBundle.
//
// This can be either the output of an
// issue call from the PKI backend or just its data member; or,
// JSON not coming from the PKI backend.
func ParsePKIJSON(input []byte) (*ParsedCertBundle, error) {
	result := &CertBundle{}
	err := json.Unmarshal(input, &result)

	if err == nil {
		return result.ToParsedCertBundle()
	}

	var secret Secret
	err = json.Unmarshal(input, &secret)

	if err == nil {
		return ParsePKIMap(secret.Data)
	}

	return nil, UserError{"Unable to parse out of either secret data or a secret object"}
}

// ParsePEMBundle takes a string of concatenated PEM-format certificate
// and private key values and decodes/parses them, checking validity along
// the way. There must be at max two certificates (a certificate and its
// issuing certificate) and one private key.
func ParsePEMBundle(pemBundle string) (*ParsedCertBundle, error) {
	if len(pemBundle) == 0 {
		return nil, UserError{"Empty PEM bundle"}
	}

	pemBytes := []byte(pemBundle)
	var pemBlock *pem.Block
	parsedBundle := &ParsedCertBundle{}

	for {
		pemBlock, pemBytes = pem.Decode(pemBytes)
		if pemBlock == nil {
			return nil, UserError{"No data found"}
		}

		if signer, err := x509.ParseECPrivateKey(pemBlock.Bytes); err == nil {
			if parsedBundle.PrivateKeyType != UnknownPrivateKey {
				return nil, UserError{"More than one private key given; provide only one private key in the bundle"}
			}
			parsedBundle.PrivateKeyType = ECPrivateKey
			parsedBundle.PrivateKeyBytes = pemBlock.Bytes
			parsedBundle.PrivateKey = signer

		} else if signer, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err == nil {
			if parsedBundle.PrivateKeyType != UnknownPrivateKey {
				return nil, UserError{"More than one private key given; provide only one private key in the bundle"}
			}
			parsedBundle.PrivateKeyType = RSAPrivateKey
			parsedBundle.PrivateKeyBytes = pemBlock.Bytes
			parsedBundle.PrivateKey = signer

		} else if certificates, err := x509.ParseCertificates(pemBlock.Bytes); err == nil {
			switch len(certificates) {
			case 0:
				return nil, UserError{"PEM block cannot be decoded to a private key or certificate"}

			case 1:
				if parsedBundle.Certificate != nil {
					switch {
					// We just found the issuing CA
					case bytes.Equal(parsedBundle.Certificate.AuthorityKeyId, certificates[0].SubjectKeyId) && certificates[0].IsCA:
						parsedBundle.IssuingCABytes = pemBlock.Bytes
						parsedBundle.IssuingCA = certificates[0]

					// Our saved certificate is actually the issuing CA
					case bytes.Equal(parsedBundle.Certificate.SubjectKeyId, certificates[0].AuthorityKeyId) && parsedBundle.Certificate.IsCA:
						parsedBundle.IssuingCA = parsedBundle.Certificate
						parsedBundle.IssuingCABytes = parsedBundle.CertificateBytes
						parsedBundle.CertificateBytes = pemBlock.Bytes
						parsedBundle.Certificate = certificates[0]
					}
				} else {
					switch {
					// If this case isn't correct, the caller needs to assign
					// the values to Certificate/CertificateBytes; assumptions
					// made here will not be valid for all cases.
					case certificates[0].IsCA:
						parsedBundle.IssuingCABytes = pemBlock.Bytes
						parsedBundle.IssuingCA = certificates[0]

					default:
						parsedBundle.CertificateBytes = pemBlock.Bytes
						parsedBundle.Certificate = certificates[0]
					}
				}

			default:
				return nil, UserError{"Too many certificates given; provide a maximum of two certificates in the bundle"}
			}
		}

		if len(pemBytes) == 0 {
			break
		}
	}

	return parsedBundle, nil
}
