// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/mldsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestGetKeyType verifies that GetKeyType maps the names Go uses for public key
// algorithms to Vault key types. Go reports every ML-DSA parameter set as
// "ML-DSA", so that name must map to "ml-dsa". Matching ignores case and
// surrounding whitespace, and unknown algorithms yield an empty string.
func TestGetKeyType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "rsa", input: "RSA", want: "rsa"},
		{name: "ecdsa", input: "ECDSA", want: "ec"},
		{name: "ed25519 as reported by Go", input: "Ed25519", want: "ed25519"},
		{name: "ed25519 upper case", input: "ED25519", want: "ed25519"},
		{name: "ml-dsa as reported by Go", input: "ML-DSA", want: "ml-dsa"},
		{name: "ml-dsa-44", input: "ML-DSA-44", want: "ml-dsa"},
		{name: "ml-dsa-65", input: "ML-DSA-65", want: "ml-dsa"},
		{name: "ml-dsa-87", input: "ML-DSA-87", want: "ml-dsa"},
		{name: "lower case", input: "ml-dsa", want: "ml-dsa"},
		{name: "mixed case", input: "Ecdsa", want: "ec"},
		{name: "surrounding whitespace", input: "  ML-DSA\t", want: "ml-dsa"},
		{name: "unknown algorithm", input: "DSA", want: ""},
		{name: "empty", input: "", want: ""},
		{name: "unknown public key algorithm string", input: x509.UnknownPublicKeyAlgorithm.String(), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, GetKeyType(tt.input))
		})
	}
}

// TestGetKeyType_GoPublicKeyAlgorithmNames verifies GetKeyType against the
// strings Go itself returns from x509.PublicKeyAlgorithm.String(), since that
// is what the CSR and certificate parsers pass in. Go uses a single "ML-DSA"
// name for all three parameter sets.
func TestGetKeyType_GoPublicKeyAlgorithmNames(t *testing.T) {
	t.Parallel()

	require.Equal(t, "rsa", GetKeyType(x509.RSA.String()))
	require.Equal(t, "ec", GetKeyType(x509.ECDSA.String()))
	require.Equal(t, "ed25519", GetKeyType(x509.Ed25519.String()))
	require.Equal(t, "ml-dsa", GetKeyType(x509.MLDSA.String()))
}

// TestGetMLDSAParameterSet verifies that GetMLDSAParameterSet returns the right
// parameter set for ML-DSA public and private keys of every size, and an empty
// value for nil, non-ML-DSA and untyped-nil inputs.
func TestGetMLDSAParameterSet(t *testing.T) {
	t.Parallel()

	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	edPublic, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	for _, tt := range []struct {
		name   string
		params mldsa.Parameters
		want   ParameterSet
	}{
		{name: "ml-dsa-44", params: mldsa.MLDSA44(), want: MLDSA44},
		{name: "ml-dsa-65", params: mldsa.MLDSA65(), want: MLDSA65},
		{name: "ml-dsa-87", params: mldsa.MLDSA87(), want: MLDSA87},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			privateKey, err := mldsa.GenerateKey(tt.params)
			require.NoError(t, err)

			require.Equal(t, tt.want, GetMLDSAParameterSet(privateKey), "private key")
			require.Equal(t, tt.want, GetMLDSAParameterSet(privateKey.PublicKey()), "public key")
			require.Equal(t, tt.want, GetMLDSAParameterSet(privateKey.Public()), "public key via crypto.Signer")
		})
	}

	t.Run("non ml-dsa keys", func(t *testing.T) {
		t.Parallel()

		require.Empty(t, GetMLDSAParameterSet(&ecKey.PublicKey))
		require.Empty(t, GetMLDSAParameterSet(ecKey))
		require.Empty(t, GetMLDSAParameterSet(&rsaKey.PublicKey))
		require.Empty(t, GetMLDSAParameterSet(rsaKey))
		require.Empty(t, GetMLDSAParameterSet(edPublic))
	})

	t.Run("nil inputs", func(t *testing.T) {
		t.Parallel()

		require.Empty(t, GetMLDSAParameterSet(nil))
		require.Empty(t, GetMLDSAParameterSet((*mldsa.PublicKey)(nil)))
		require.Empty(t, GetMLDSAParameterSet((*mldsa.PrivateKey)(nil)))
	})
}

// parseFieldsTestKey describes a key to use when checking the fields that
// ParseCsrToFields and ParseCertificateToFields report.
type parseFieldsTestKey struct {
	name             string
	generate         func(t *testing.T) crypto.Signer
	wantKeyType      string
	wantKeyBits      int
	wantParameterSet string
}

func parseFieldsTestKeys() []parseFieldsTestKey {
	generateMLDSA := func(params mldsa.Parameters) func(t *testing.T) crypto.Signer {
		return func(t *testing.T) crypto.Signer {
			key, err := mldsa.GenerateKey(params)
			require.NoError(t, err)
			return key
		}
	}

	return []parseFieldsTestKey{
		{
			name:             "ml-dsa-44",
			generate:         generateMLDSA(mldsa.MLDSA44()),
			wantKeyType:      "ml-dsa",
			wantParameterSet: MLDSA44,
		},
		{
			name:             "ml-dsa-65",
			generate:         generateMLDSA(mldsa.MLDSA65()),
			wantKeyType:      "ml-dsa",
			wantParameterSet: MLDSA65,
		},
		{
			name:             "ml-dsa-87",
			generate:         generateMLDSA(mldsa.MLDSA87()),
			wantKeyType:      "ml-dsa",
			wantParameterSet: MLDSA87,
		},
		{
			name: "ec",
			generate: func(t *testing.T) crypto.Signer {
				key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				require.NoError(t, err)
				return key
			},
			wantKeyType: "ec",
			wantKeyBits: 256,
		},
		{
			name: "rsa",
			generate: func(t *testing.T) crypto.Signer {
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				require.NoError(t, err)
				return key
			},
			wantKeyType: "rsa",
			wantKeyBits: 2048,
		},
		{
			name: "ed25519",
			generate: func(t *testing.T) crypto.Signer {
				_, key, err := ed25519.GenerateKey(rand.Reader)
				require.NoError(t, err)
				return key
			},
			wantKeyType: "ed25519",
		},
	}
}

// TestParseCsrToFields_KeyFields verifies that key_type, key_bits and
// parameter_set are reported for CSRs of every supported key type. The
// parameter_set key must always be present: it holds the ML-DSA parameter set
// for ML-DSA keys and an empty string for every other key type.
func TestParseCsrToFields_KeyFields(t *testing.T) {
	t.Parallel()

	for _, tt := range parseFieldsTestKeys() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			key := tt.generate(t)
			csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
				Subject: pkix.Name{CommonName: "localhost"},
			}, key)
			require.NoError(t, err)
			csr, err := x509.ParseCertificateRequest(csrBytes)
			require.NoError(t, err)

			fields, err := ParseCsrToFields(*csr)
			require.NoError(t, err)

			require.Equal(t, tt.wantKeyType, fields["key_type"], "key_type")
			require.Equal(t, tt.wantKeyBits, fields["key_bits"], "key_bits")
			require.Contains(t, fields, "parameter_set", "parameter_set must always be present")
			require.Equal(t, string(tt.wantParameterSet), fields["parameter_set"], "parameter_set")
		})
	}
}

// TestParseCertificateToFields_KeyFields verifies that key_type, key_bits and
// parameter_set are reported for certificates of every supported key type. The
// parameter_set key must always be present: it holds the ML-DSA parameter set
// for ML-DSA keys and an empty string for every other key type.
func TestParseCertificateToFields_KeyFields(t *testing.T) {
	t.Parallel()

	for _, tt := range parseFieldsTestKeys() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			key := tt.generate(t)
			template := &x509.Certificate{
				Subject:      pkix.Name{CommonName: "localhost"},
				SerialNumber: big.NewInt(1),
				NotBefore:    time.Now().Add(-time.Minute),
				NotAfter:     time.Now().Add(time.Hour),
			}
			certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
			require.NoError(t, err)
			cert, err := x509.ParseCertificate(certBytes)
			require.NoError(t, err)

			fields, err := ParseCertificateToFields(*cert)
			require.NoError(t, err)

			require.Equal(t, tt.wantKeyType, fields["key_type"], "key_type")
			require.Equal(t, tt.wantKeyBits, fields["key_bits"], "key_bits")
			require.Contains(t, fields, "parameter_set", "parameter_set must always be present")
			require.Equal(t, string(tt.wantParameterSet), fields["parameter_set"], "parameter_set")
		})
	}
}
