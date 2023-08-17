package pki

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func Test_ConfigsRoundTrip(t *testing.T) {
	_, s := createBackendWithStorage(t)

	// Create an empty key, issuer for testing.
	key := keyEntry{ID: genKeyId()}
	err := writeKey(ctx, s, key)
	require.NoError(t, err)
	issuer := &issuerEntry{ID: genIssuerId()}
	err = writeIssuer(ctx, s, issuer)
	require.NoError(t, err)

	// Verify we handle nothing stored properly
	keyConfigEmpty, err := getKeysConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &keyConfigEntry{}, keyConfigEmpty)

	issuerConfigEmpty, err := getIssuersConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &issuerConfigEntry{}, issuerConfigEmpty)

	// Now attempt to store and reload properly
	origKeyConfig := &keyConfigEntry{
		DefaultKeyId: key.ID,
	}
	origIssuerConfig := &issuerConfigEntry{
		DefaultIssuerId: issuer.ID,
	}

	err = setKeysConfig(ctx, s, origKeyConfig)
	require.NoError(t, err)
	err = setIssuersConfig(ctx, s, origIssuerConfig)
	require.NoError(t, err)

	keyConfig, err := getKeysConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, origKeyConfig, keyConfig)

	issuerConfig, err := getIssuersConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, origIssuerConfig.DefaultIssuerId, issuerConfig.DefaultIssuerId)
}

func Test_IssuerRoundTrip(t *testing.T) {
	b, s := createBackendWithStorage(t)
	issuer1, key1 := genIssuerAndKey(t, b, s)
	issuer2, key2 := genIssuerAndKey(t, b, s)

	// We get an error when issuer id not found
	_, err := fetchIssuerById(ctx, s, issuer1.ID)
	require.Error(t, err)

	// We get an error when key id not found
	_, err = fetchKeyById(ctx, s, key1.ID)
	require.Error(t, err)

	// Now write out our issuers and keys
	err = writeKey(ctx, s, key1)
	require.NoError(t, err)
	err = writeIssuer(ctx, s, &issuer1)
	require.NoError(t, err)

	err = writeKey(ctx, s, key2)
	require.NoError(t, err)
	err = writeIssuer(ctx, s, &issuer2)
	require.NoError(t, err)

	fetchedKey1, err := fetchKeyById(ctx, s, key1.ID)
	require.NoError(t, err)

	fetchedIssuer1, err := fetchIssuerById(ctx, s, issuer1.ID)
	require.NoError(t, err)

	require.Equal(t, &key1, fetchedKey1)
	require.Equal(t, &issuer1, fetchedIssuer1)

	keys, err := listKeys(ctx, s)
	require.NoError(t, err)

	require.ElementsMatch(t, []keyID{key1.ID, key2.ID}, keys)

	issuers, err := listIssuers(ctx, s)
	require.NoError(t, err)

	require.ElementsMatch(t, []issuerID{issuer1.ID, issuer2.ID}, issuers)
}

func Test_KeysIssuerImport(t *testing.T) {
	b, s := createBackendWithStorage(t)

	issuer1, key1 := genIssuerAndKey(t, b, s)
	issuer2, key2 := genIssuerAndKey(t, b, s)

	// Key 1 before Issuer 1; Issuer 2 before Key 2.
	// Remove KeyIDs from non-written entities before beginning.
	key1.ID = ""
	issuer1.ID = ""
	issuer1.KeyID = ""

	key1Ref1, existing, err := importKey(ctx, b, s, key1.PrivateKey, "key1", key1.PrivateKeyType)
	require.NoError(t, err)
	require.False(t, existing)
	require.Equal(t, strings.TrimSpace(key1.PrivateKey), strings.TrimSpace(key1Ref1.PrivateKey))

	// Make sure if we attempt to re-import the same private key, no import/updates occur.
	// So the existing flag should be set to true, and we do not update the existing Name field.
	key1Ref2, existing, err := importKey(ctx, b, s, key1.PrivateKey, "ignore-me", key1.PrivateKeyType)
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, key1.PrivateKey, key1Ref1.PrivateKey)
	require.Equal(t, key1Ref1.ID, key1Ref2.ID)
	require.Equal(t, key1Ref1.Name, key1Ref2.Name)

	issuer1Ref1, existing, err := importIssuer(ctx, b, s, issuer1.Certificate, "issuer1")
	require.NoError(t, err)
	require.False(t, existing)
	require.Equal(t, strings.TrimSpace(issuer1.Certificate), strings.TrimSpace(issuer1Ref1.Certificate))
	require.Equal(t, key1Ref1.ID, issuer1Ref1.KeyID)
	require.Equal(t, "issuer1", issuer1Ref1.Name)

	// Make sure if we attempt to re-import the same issuer, no import/updates occur.
	// So the existing flag should be set to true, and we do not update the existing Name field.
	issuer1Ref2, existing, err := importIssuer(ctx, b, s, issuer1.Certificate, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, strings.TrimSpace(issuer1.Certificate), strings.TrimSpace(issuer1Ref1.Certificate))
	require.Equal(t, issuer1Ref1.ID, issuer1Ref2.ID)
	require.Equal(t, key1Ref1.ID, issuer1Ref2.KeyID)
	require.Equal(t, issuer1Ref1.Name, issuer1Ref2.Name)

	err = writeIssuer(ctx, s, &issuer2)
	require.NoError(t, err)

	err = writeKey(ctx, s, key2)
	require.NoError(t, err)

	// Same double import tests as above, but make sure if the previous was created through writeIssuer not importIssuer.
	issuer2Ref, existing, err := importIssuer(ctx, b, s, issuer2.Certificate, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, strings.TrimSpace(issuer2.Certificate), strings.TrimSpace(issuer2Ref.Certificate))
	require.Equal(t, issuer2.ID, issuer2Ref.ID)
	require.Equal(t, "", issuer2Ref.Name)
	require.Equal(t, issuer2.KeyID, issuer2Ref.KeyID)

	// Same double import tests as above, but make sure if the previous was created through writeKey not importKey.
	key2Ref, existing, err := importKey(ctx, b, s, key2.PrivateKey, "ignore-me", key2.PrivateKeyType)
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, key2.PrivateKey, key2Ref.PrivateKey)
	require.Equal(t, key2.ID, key2Ref.ID)
	require.Equal(t, "", key2Ref.Name)
}

func genIssuerAndKey(t *testing.T, b *backend, s logical.Storage) (issuerEntry, keyEntry) {
	certBundle := genCertBundle(t, b, s)

	keyId := genKeyId()

	pkiKey := keyEntry{
		ID:             keyId,
		PrivateKeyType: certBundle.PrivateKeyType,
		PrivateKey:     strings.TrimSpace(certBundle.PrivateKey) + "\n",
	}

	issuerId := genIssuerId()

	pkiIssuer := issuerEntry{
		ID:           issuerId,
		KeyID:        keyId,
		Certificate:  strings.TrimSpace(certBundle.Certificate) + "\n",
		CAChain:      certBundle.CAChain,
		SerialNumber: certBundle.SerialNumber,
	}

	return pkiIssuer, pkiKey
}

func genCertBundle(t *testing.T, b *backend, s logical.Storage) *certutil.CertBundle {
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
	_, _, role, respErr := b.getGenerationParams(ctx, s, apiData)
	require.Nil(t, respErr)

	input := &inputBundle{
		req: &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "issue/testrole",
			Storage:   s,
		},
		apiData: apiData,
		role:    role,
	}
	parsedCertBundle, err := generateCert(ctx, b, input, nil, true, b.GetRandomReader())

	require.NoError(t, err)
	certBundle, err := parsedCertBundle.ToCertBundle()
	require.NoError(t, err)
	return certBundle
}

func writeLegacyBundle(t *testing.T, b *backend, s logical.Storage, bundle *certutil.CertBundle) {
	entry, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
	require.NoError(t, err)

	err = s.Put(context.Background(), entry)
	require.NoError(t, err)
}
