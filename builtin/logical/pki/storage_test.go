package pki

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func Test_PKIIssuerRoundTrip(t *testing.T) {
	b, s := createBackendWithStorage(t)
	pkiIssuer, pkiKey := genIssuerAndKey(t, b)

	// We get an error when issuer id not found
	_, err := fetchIssuerById(ctx, s, pkiIssuer.ID)
	require.Error(t, err)

	// We get an error when key id not found
	_, err = fetchKeyById(ctx, s, pkiKey.ID)
	require.Error(t, err)

	// Now write out our issuer and key
	err = writeKey(ctx, s, pkiKey)
	require.NoError(t, err)
	err = writeIssuer(ctx, s, pkiIssuer)
	require.NoError(t, err)

	pkiKey2, err := fetchKeyById(ctx, s, pkiKey.ID)
	require.NoError(t, err)

	pkiIssuer2, err := fetchIssuerById(ctx, s, pkiIssuer.ID)
	require.NoError(t, err)

	require.Equal(t, pkiKey, pkiKey2)
	require.Equal(t, pkiIssuer, pkiIssuer2)
}

func genIssuerAndKey(t *testing.T, b *backend) (issuer, key) {
	certBundle, err := genCertBundle(t, b)

	keyIdStr, err := uuid.GenerateUUID()
	require.NoError(t, err)

	pkiKey := key{
		ID:             keyId(keyIdStr),
		PrivateKeyType: certBundle.PrivateKeyType,
		PrivateKey:     certBundle.PrivateKey,
	}

	issuerIdStr, err := uuid.GenerateUUID()
	require.NoError(t, err)

	pkiIssuer := issuer{
		ID:           issuerId(issuerIdStr),
		KeyID:        keyId(keyIdStr),
		Certificate:  certBundle.Certificate,
		CAChain:      certBundle.CAChain,
		SerialNumber: certBundle.SerialNumber,
	}

	return pkiIssuer, pkiKey
}

func genCertBundle(t *testing.T, b *backend) (*certutil.CertBundle, error) {
	// Pretty gross just to generate a cert bundle, but
	fields := addCACommonFields(map[string]*framework.FieldSchema{})
	fields = addCAKeyGenerationFields(fields)
	fields = addCAIssueFields(fields)
	apiData := &framework.FieldData{
		Schema: fields,
		Raw: map[string]interface{}{
			"exported": "internal",
			"cn":       "example.com",
			"ttl":      3600,
		},
	}
	_, _, role, respErr := b.getGenerationParams(ctx, apiData, "/pki")
	require.Nil(t, respErr)

	input := &inputBundle{
		req: &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "issue/testrole",
			Storage:   b.storage,
		},
		apiData: apiData,
		role:    role,
	}
	parsedCertBundle, err := generateCert(ctx, b, input, nil, true, rand.Reader)

	require.NoError(t, err)
	certBundle, err := parsedCertBundle.ToCertBundle()
	require.NoError(t, err)
	return certBundle, err
}
