package pki

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestIntegration_RotateRootUsesNext(t *testing.T) {
	b, s := createBackendWithStorage(t)
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/rotate/internal",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed rotate root")
	require.NotNil(t, resp, "got nil response from rotate root")
	require.False(t, resp.IsError(), "got an error from rotate root: %#v", resp)

	issuerId1 := resp.Data["issuer_id"].(issuerID)
	issuerName1 := resp.Data["issuer_name"]

	require.NotEmpty(t, issuerId1, "issuer id was empty on initial rotate root command")
	require.Equal(t, "next", issuerName1, "expected an issuer name of next on initial rotate root command")

	// Call it again, we should get a new issuer id, but since next issuer_name is used we should get a blank value.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/rotate/internal",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed rotate root")
	require.NotNil(t, resp, "got nil response from rotate root")
	require.False(t, resp.IsError(), "got an error from rotate root: %#v", resp)

	issuerId2 := resp.Data["issuer_id"].(issuerID)
	issuerName2 := resp.Data["issuer_name"]

	require.NotEmpty(t, issuerId2, "issuer id was empty on second rotate root command")
	require.NotEqual(t, issuerId1, issuerId2, "should have been different issuer ids")
	require.Empty(t, issuerName2, "expected a blank issuer name on the second rotate root command")

	// Call it again, making sure we can use our own name if desired.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/rotate/internal",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test.com",
			"issuer_name": "next-cert",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed rotate root")
	require.NotNil(t, resp, "got nil response from rotate root")
	require.False(t, resp.IsError(), "got an error from rotate root: %#v", resp)

	issuerId3 := resp.Data["issuer_id"].(issuerID)
	issuerName3 := resp.Data["issuer_name"]

	require.NotEmpty(t, issuerId3, "issuer id was empty on third rotate root command")
	require.NotEqual(t, issuerId3, issuerId1, "should have been different issuer id from initial")
	require.NotEqual(t, issuerId3, issuerId2, "should have been different issuer id from second call")
	require.Equal(t, "next-cert", issuerName3, "expected an issuer name that we specified on third rotate root command")
}

func TestIntegration_ReplaceRootNormal(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// generate roots
	genTestRootCa(t, b, s)
	issuerId2, _ := genTestRootCa(t, b, s)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": issuerId2.String(),
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	replacedIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, replacedIssuer, "expected return value to match issuer we set")

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "config/issuers",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	defaultIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, defaultIssuer, "expected default issuer to be updated")
}

func TestIntegration_ReplaceRootDefaultsToNext(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// generate roots
	genTestRootCa(t, b, s)
	issuerId2, _ := genTestRootCaWithIssuerName(t, b, s, "next")

	// Do not specify the default value to replace.
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "root/replace",
		Storage:    s,
		Data:       map[string]interface{}{},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	replacedIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, replacedIssuer, "expected return value to match the 'next' issuer we set")

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "config/issuers",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	defaultIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, defaultIssuer, "expected default issuer to be updated")
}

func TestIntegration_ReplaceRootBadIssuer(t *testing.T) {
	b, s := createBackendWithStorage(t)

	// generate roots
	genTestRootCa(t, b, s)
	genTestRootCa(t, b, s)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": "a-bad-issuer-id",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root, should have been an error in the response.")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.True(t, resp.IsError(), "did not get an error from replacing root: %#v", resp)

	// Make sure we trap replacing with default.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": "default",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root, should have been an error in the response.")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.True(t, resp.IsError(), "did not get an error from replacing root: %#v", resp)

	// Make sure we trap replacing with blank string.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": "",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root, should have been an error in the response.")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.True(t, resp.IsError(), "did not get an error from replacing root: %#v", resp)
}

func TestIntegration_SetSignedWithBackwardsPemBundles(t *testing.T) {
	rootBackend, rootStorage := createBackendWithStorage(t)
	intBackend, intStorage := createBackendWithStorage(t)

	// generate root
	resp, err := rootBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issuers/generate/root/internal",
		Storage:   rootStorage,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)
	rootCert := resp.Data["certificate"].(string)

	// generate intermediate
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issuers/generate/intermediate/internal",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed generating int ca")
	require.NotNil(t, resp, "got nil response from generating int ca")
	require.False(t, resp.IsError(), "got an error from generating int ca: %#v", resp)
	intCsr := resp.Data["csr"].(string)

	// sign csr
	resp, err = rootBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-intermediate",
		Storage:   rootStorage,
		Data: map[string]interface{}{
			"csr":    intCsr,
			"format": "pem_bundle",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)

	intCert := resp.Data["certificate"].(string)

	// Send in the chain backwards now and make sure we link intCert as default.
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "intermediate/set-signed",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"certificate": rootCert + "\n" + intCert + "\n",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)

	// setup role
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/example",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"allowed_domains":  "example.com",
			"allow_subdomains": "true",
			"max_ttl":          "1h",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed setting up role example")
	require.Nil(t, resp, "got non-nil response from setting up role example: %#v", resp)

	// Issue cert
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issue/example",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"common_name": "test.example.com",
			"ttl":         "5m",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed issuing a leaf cert from int ca")
	require.NotNil(t, resp, "got nil response issuing a leaf cert from int ca")
	require.False(t, resp.IsError(), "got an error issuing a leaf cert from int ca: %#v", resp)
}

func genTestRootCa(t *testing.T, b *backend, s logical.Storage) (issuerID, keyID) {
	return genTestRootCaWithIssuerName(t, b, s, "")
}

func genTestRootCaWithIssuerName(t *testing.T, b *backend, s logical.Storage, issuerName string) (issuerID, keyID) {
	data := map[string]interface{}{
		"common_name": "test.com",
	}
	if len(issuerName) > 0 {
		data["issuer_name"] = issuerName
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "issuers/generate/root/internal",
		Storage:    s,
		Data:       data,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)

	issuerId := resp.Data["issuer_id"].(issuerID)
	keyId := resp.Data["key_id"].(keyID)

	require.NotEmpty(t, issuerId, "returned issuer id was empty")
	require.NotEmpty(t, keyId, "returned key id was empty")

	return issuerId, keyId
}
