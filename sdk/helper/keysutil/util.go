package keysutil

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"

	"golang.org/x/crypto/ed25519"
)

// pkcs8 reflects an ASN.1, PKCS #8 PrivateKey. See
// ftp://ftp.rsasecurity.com/pub/pkcs/pkcs-8/pkcs-8v1_2.asn
// and RFC 5208.
//
// Copied from Go: https://github.com/golang/go/blob/master/src/crypto/x509/pkcs8.go#L17-L80
type pkcs8 struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
	// optional attributes omitted.
}

// ecPrivateKey reflects an ASN.1 Elliptic Curve Private Key Structure.
// References:
//
//	RFC 5915
//	SEC1 - http://www.secg.org/sec1-v2.pdf
//
// Per RFC 5915 the NamedCurveOID is marked as ASN.1 OPTIONAL, however in
// most cases it is not.
//
// Copied from Go: 	https://github.com/golang/go/blob/master/src/crypto/x509/sec1.go#L18-L31
type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`

	// Because the PKCS8/RFC 5915 encoding of the Ed25519 key uses the
	// RFC 8032 Ed25519 seed format, we can ignore the public key parameter
	// and infer it later.
}

var (
	// See crypto/x509/x509.go in the Go toolchain source distribution.
	oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}

	// NSS encodes Ed25519 private keys with the OID 1.3.6.1.4.1.11591.15.1
	// from https://tools.ietf.org/html/draft-josefsson-pkix-newcurves-01.
	// See https://github.com/nss-dev/nss/blob/NSS_3_79_BRANCH/lib/util/secoid.c#L600-L603.
	oidNSSPKIXEd25519 = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 11591, 15, 1}

	// Other implementations may use the OID 1.3.101.110 from
	// https://datatracker.ietf.org/doc/html/rfc8410.
	oidRFC8410Ed25519 = asn1.ObjectIdentifier{1, 3, 101, 110}
)

func isEd25519OID(oid asn1.ObjectIdentifier) bool {
	return oidNSSPKIXEd25519.Equal(oid) || oidRFC8410Ed25519.Equal(oid)
}

// ParsePKCS8PrivateKey parses an unencrypted private key in PKCS #8, ASN.1 DER form.
//
// It returns a *rsa.PrivateKey, a *ecdsa.PrivateKey, or a ed25519.PrivateKey.
// More types might be supported in the future.
//
// This kind of key is commonly encoded in PEM blocks of type "PRIVATE KEY".
func ParsePKCS8Ed25519PrivateKey(der []byte) (key interface{}, err error) {
	var privKey pkcs8
	var ed25519Key ecPrivateKey

	var checkedOID bool

	// If this err is nil, we assume we directly have a ECPrivateKey structure
	// with explicit OID; ignore this error for now and return the latter err
	// instead if neither parse correctly.
	if _, err := asn1.Unmarshal(der, &privKey); err == nil {
		switch {
		case privKey.Algo.Algorithm.Equal(oidPublicKeyECDSA):
			bytes := privKey.Algo.Parameters.FullBytes
			namedCurveOID := new(asn1.ObjectIdentifier)
			if _, err := asn1.Unmarshal(bytes, namedCurveOID); err != nil {
				namedCurveOID = nil
			}

			if namedCurveOID == nil || !isEd25519OID(*namedCurveOID) {
				return nil, errors.New("keysutil: failed to parse private key (invalid, non-ed25519 curve parameter OID)")
			}

			der = privKey.PrivateKey
			checkedOID = true
		default:
			// The Go standard library already parses RFC 8410 keys; the
			// inclusion of the OID here is in case it is used with the
			// regular ECDSA PrivateKey structure, rather than the struct
			// recognized by the Go standard library.
			return nil, errors.New("keysutil: failed to parse key as ed25519 private key")
		}
	}

	_, err = asn1.Unmarshal(der, &ed25519Key)
	if err != nil {
		return nil, fmt.Errorf("keysutil: failed to parse private key (inner Ed25519 ECPrivateKey format was incorrect): %v", err)
	}

	if !checkedOID && !isEd25519OID(ed25519Key.NamedCurveOID) {
		return nil, errors.New("keysutil: failed to parse private key (invalid, non-ed25519 curve parameter OID)")
	}

	if len(ed25519Key.PrivateKey) != 32 {
		return nil, fmt.Errorf("keysutil: failed to parse private key as ed25519 private key: got %v bytes but expected %v byte RFC 8032 seed", len(ed25519Key.PrivateKey), ed25519.SeedSize)
	}

	return ed25519.NewKeyFromSeed(ed25519Key.PrivateKey), nil
}
