package pkcs7

import (
	"bytes"
	"crypto"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecrypt(t *testing.T) {
	fixture := UnmarshalTestFixture(EncryptedTestFixture)
	p7, err := Parse(fixture.Input)
	if err != nil {
		t.Fatal(err)
	}
	content, err := p7.Decrypt(fixture.Certificate, fixture.PrivateKey)
	if err != nil {
		t.Errorf("Cannot Decrypt with error: %v", err)
	}
	expected := []byte("This is a test")
	if !bytes.Equal(content, expected) {
		t.Errorf("Decrypted result does not match.\n\tExpected:%s\n\tActual:%s", expected, content)
	}
}

// echo -n "This is a test" > test.txt
// openssl cms -encrypt -in test.txt cert.pem
var EncryptedTestFixture = `
-----BEGIN PKCS7-----
MIIBGgYJKoZIhvcNAQcDoIIBCzCCAQcCAQAxgcwwgckCAQAwMjApMRAwDgYDVQQK
EwdBY21lIENvMRUwEwYDVQQDEwxFZGRhcmQgU3RhcmsCBQDL+CvWMA0GCSqGSIb3
DQEBAQUABIGAyFz7bfI2noUs4FpmYfztm1pVjGyB00p9x0H3gGHEYNXdqlq8VG8d
iq36poWtEkatnwsOlURWZYECSi0g5IAL0U9sj82EN0xssZNaK0S5FTGnB3DPvYgt
HJvcKq7YvNLKMh4oqd17C6GB4oXyEBDj0vZnL7SUoCAOAWELPeC8CTUwMwYJKoZI
hvcNAQcBMBQGCCqGSIb3DQMHBAhEowTkot3a7oAQFD//J/IhFnk+JbkH7HZQFA==
-----END PKCS7-----
-----BEGIN CERTIFICATE-----
MIIB1jCCAUGgAwIBAgIFAMv4K9YwCwYJKoZIhvcNAQELMCkxEDAOBgNVBAoTB0Fj
bWUgQ28xFTATBgNVBAMTDEVkZGFyZCBTdGFyazAeFw0xNTA1MDYwMzU2NDBaFw0x
NjA1MDYwMzU2NDBaMCUxEDAOBgNVBAoTB0FjbWUgQ28xETAPBgNVBAMTCEpvbiBT
bm93MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDK6NU0R0eiCYVquU4RcjKc
LzGfx0aa1lMr2TnLQUSeLFZHFxsyyMXXuMPig3HK4A7SGFHupO+/1H/sL4xpH5zg
8+Zg2r8xnnney7abxcuv0uATWSIeKlNnb1ZO1BAxFnESc3GtyOCr2dUwZHX5mRVP
+Zxp2ni5qHNraf3wE2VPIQIDAQABoxIwEDAOBgNVHQ8BAf8EBAMCAKAwCwYJKoZI
hvcNAQELA4GBAIr2F7wsqmEU/J/kLyrCgEVXgaV/sKZq4pPNnzS0tBYk8fkV3V18
sBJyHKRLL/wFZASvzDcVGCplXyMdAOCyfd8jO3F9Ac/xdlz10RrHJT75hNu3a7/n
9KNwKhfN4A1CQv2x372oGjRhCW5bHNCWx4PIVeNzCyq/KZhyY9sxHE6f
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MIICXgIBAAKBgQDK6NU0R0eiCYVquU4RcjKcLzGfx0aa1lMr2TnLQUSeLFZHFxsy
yMXXuMPig3HK4A7SGFHupO+/1H/sL4xpH5zg8+Zg2r8xnnney7abxcuv0uATWSIe
KlNnb1ZO1BAxFnESc3GtyOCr2dUwZHX5mRVP+Zxp2ni5qHNraf3wE2VPIQIDAQAB
AoGBALyvnSt7KUquDen7nXQtvJBudnf9KFPt//OjkdHHxNZNpoF/JCSqfQeoYkeu
MdAVYNLQGMiRifzZz4dDhA9xfUAuy7lcGQcMCxEQ1dwwuFaYkawbS0Tvy2PFlq2d
H5/HeDXU4EDJ3BZg0eYj2Bnkt1sJI35UKQSxblQ0MY2q0uFBAkEA5MMOogkgUx1C
67S1tFqMUSM8D0mZB0O5vOJZC5Gtt2Urju6vywge2ArExWRXlM2qGl8afFy2SgSv
Xk5eybcEiQJBAOMRwwbEoW5NYHuFFbSJyWll4n71CYuWuQOCzehDPyTb80WFZGLV
i91kFIjeERyq88eDE5xVB3ZuRiXqaShO/9kCQQCKOEkpInaDgZSjskZvuJ47kByD
6CYsO4GIXQMMeHML8ncFH7bb6AYq5ybJVb2NTU7QLFJmfeYuhvIm+xdOreRxAkEA
o5FC5Jg2FUfFzZSDmyZ6IONUsdF/i78KDV5nRv1R+hI6/oRlWNCtTNBv/lvBBd6b
dseUE9QoaQZsn5lpILEvmQJAZ0B+Or1rAYjnbjnUhdVZoy9kC4Zov+4UH3N/BtSy
KJRWUR0wTWfZBPZ5hAYZjTBEAFULaYCXlQKsODSp0M1aQA==
-----END PRIVATE KEY-----`

// TestRsaOAEPKeyIdentifier validates that we can parse the RSA-OAEP algo
// identifier's parameters correctly. See RFC 3447 for the definition
func TestRsaOAEPKeyIdentifier(t *testing.T) {
	t.Parallel()

	t.Run("no params", func(t *testing.T) {
		hash, err := parseRSAOaepParam(asn1.NullRawValue)
		require.NoError(t, err)
		require.Equal(t, hash, crypto.SHA1)
	})

	t.Run("only pSource parameter set", func(t *testing.T) {
		// This is the base64 encoded key identifier that only has pSource parameter set.
		// Example pulled from https://github.com/pnwamk/mbedtls/blob/rsaes-oaep/example_tpm_cert.der linked
		// from https://github.com/Mbed-TLS/mbedtls/issues/1015
		base64KeyId := "MCIGCSqGSIb3DQEBBzAVohMwEQYJKoZIhvcNAQEJBARUQ1BB"
		rawKeyId, err := base64.StdEncoding.DecodeString(base64KeyId)
		require.NoError(t, err)

		var keyId pkix.AlgorithmIdentifier
		_, err = asn1.Unmarshal(rawKeyId, &keyId)
		require.NoError(t, err)

		hash, err := parseRSAOaepParam(keyId.Parameters)
		require.NoError(t, err)
		require.Equal(t, hash, crypto.SHA1)
	})

	t.Run("full parameters", func(t *testing.T) {
		// This is the base64 encoded key identifier that has all parameters set, example pulled from
		// https://stackoverflow.com/questions/22194359/rsaes-oaep-certificate-public-key-0-bits
		base64KeyId := "MFIGCSqGSIb3DQEBBzBFoA8wDQYJYIZIAWUDBAIBBQChHDAaBgkqhkiG9w0BAQgwDQYJYIZIAWUDBAIBBQCiFDASBgkqhkiG9w0BAQkEBVRDUEEA"
		rawKeyId, err := base64.StdEncoding.DecodeString(base64KeyId)
		require.NoError(t, err)

		var keyId pkix.AlgorithmIdentifier
		_, err = asn1.Unmarshal(rawKeyId, &keyId)
		require.NoError(t, err)

		hash, err := parseRSAOaepParam(keyId.Parameters)
		require.NoError(t, err)
		require.Equal(t, hash, crypto.SHA256)
	})
}
