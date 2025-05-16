package pkcs7

import (
	_ "embed"
	"encoding/asn1"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed "testdata/intune-pkcs7-sample.msg"
var base64msg string

// TestCustomParseX509Certificates tests the custom parsing of X.509 certificates from a PKCS7 message,
// around the intune PKCS7 message that contains a certificate with a key authority extension
// that is marked as critical. This is a workaround for the fact that the standard x509 library
// has a specific check forbidding this
func TestCustomParseX509Certificates(t *testing.T) {
	msg, err := base64.StdEncoding.DecodeString(base64msg)
	require.NoError(t, err)

	myPkcs7, err := Parse(msg)
	require.NoErrorf(t, err, "CustomParseX509Certificates: %v", err)
	require.NotNil(t, myPkcs7)
	require.NotNil(t, myPkcs7.Certificates)
	require.Equal(t, 1, len(myPkcs7.Certificates))

	OIDExtensionAuthorityKeyId := asn1.ObjectIdentifier{2, 5, 29, 35}

	foundExt := false
	for _, ext := range myPkcs7.Certificates[0].Extensions {
		if ext.Id.Equal(OIDExtensionAuthorityKeyId) {
			require.Equal(t, true, ext.Critical)
			foundExt = true
		}
	}
	require.True(t, foundExt, "did not find the expected extension")
}
