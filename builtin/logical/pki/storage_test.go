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
	_, err := fetchPKIIssuerById(ctx, s, pkiIssuer.ID)
	require.Error(t, err)

	// We get an error when pkiKey id not found
	_, err = fetchPKIKeyById(ctx, s, pkiKey.ID)
	require.Error(t, err)

	// Now write out our issuer and key
	err = writePKIKey(ctx, s, pkiKey)
	require.NoError(t, err)
	err = writePKIIssuer(ctx, s, pkiIssuer)
	require.NoError(t, err)

	pkiKey2, err := fetchPKIKeyById(ctx, s, pkiKey.ID)
	require.NoError(t, err)

	pkiIssuer2, err := fetchPKIIssuerById(ctx, s, pkiIssuer.ID)
	require.NoError(t, err)

	require.Equal(t, pkiKey, pkiKey2)
	require.Equal(t, pkiIssuer, pkiIssuer2)
}

func genIssuerAndKey(t *testing.T, b *backend) (pkiIssuer, pkiKey) {
	certBundle, err := genCertBundle(t, b)

	keyId, err := uuid.GenerateUUID()
	require.NoError(t, err)

	pkiKey := pkiKey{
		ID:             pkiKeyId(keyId),
		PrivateKeyType: certBundle.PrivateKeyType,
		PrivateKey:     certBundle.PrivateKey,
	}

	issuerId, err := uuid.GenerateUUID()
	require.NoError(t, err)

	pkiIssuer := pkiIssuer{
		ID:           pkiIssuerId(issuerId),
		PKIKeyID:     pkiKeyId(keyId),
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
