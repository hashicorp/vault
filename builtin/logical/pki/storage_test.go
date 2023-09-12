// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Create an empty key, issuer for testing.
	key := keyEntry{ID: genKeyId()}
	err := sc.writeKey(key)
	require.NoError(t, err)
	issuer := &issuerEntry{ID: genIssuerId()}
	err = sc.writeIssuer(issuer)
	require.NoError(t, err)

	// Verify we handle nothing stored properly
	keyConfigEmpty, err := sc.getKeysConfig()
	require.NoError(t, err)
	require.Equal(t, &keyConfigEntry{}, keyConfigEmpty)

	issuerConfigEmpty, err := sc.getIssuersConfig()
	require.NoError(t, err)
	require.Equal(t, &issuerConfigEntry{}, issuerConfigEmpty)

	// Now attempt to store and reload properly
	origKeyConfig := &keyConfigEntry{
		DefaultKeyId: key.ID,
	}
	origIssuerConfig := &issuerConfigEntry{
		DefaultIssuerId: issuer.ID,
	}

	err = sc.setKeysConfig(origKeyConfig)
	require.NoError(t, err)
	err = sc.setIssuersConfig(origIssuerConfig)
	require.NoError(t, err)

	keyConfig, err := sc.getKeysConfig()
	require.NoError(t, err)
	require.Equal(t, origKeyConfig, keyConfig)

	issuerConfig, err := sc.getIssuersConfig()
	require.NoError(t, err)
	require.Equal(t, origIssuerConfig.DefaultIssuerId, issuerConfig.DefaultIssuerId)
}

func Test_IssuerRoundTrip(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)
	issuer1, key1 := genIssuerAndKey(t, b, s)
	issuer2, key2 := genIssuerAndKey(t, b, s)

	// We get an error when issuer id not found
	_, err := sc.fetchIssuerById(issuer1.ID)
	require.Error(t, err)

	// We get an error when key id not found
	_, err = sc.fetchKeyById(key1.ID)
	require.Error(t, err)

	// Now write out our issuers and keys
	err = sc.writeKey(key1)
	require.NoError(t, err)
	err = sc.writeIssuer(&issuer1)
	require.NoError(t, err)

	err = sc.writeKey(key2)
	require.NoError(t, err)
	err = sc.writeIssuer(&issuer2)
	require.NoError(t, err)

	fetchedKey1, err := sc.fetchKeyById(key1.ID)
	require.NoError(t, err)

	fetchedIssuer1, err := sc.fetchIssuerById(issuer1.ID)
	require.NoError(t, err)

	require.Equal(t, &key1, fetchedKey1)
	require.Equal(t, &issuer1, fetchedIssuer1)

	keys, err := sc.listKeys()
	require.NoError(t, err)

	require.ElementsMatch(t, []keyID{key1.ID, key2.ID}, keys)

	issuers, err := sc.listIssuers()
	require.NoError(t, err)

	require.ElementsMatch(t, []issuerID{issuer1.ID, issuer2.ID}, issuers)
}

func Test_KeysIssuerImport(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	issuer1, key1 := genIssuerAndKey(t, b, s)
	issuer2, key2 := genIssuerAndKey(t, b, s)

	// Key 1 before Issuer 1; Issuer 2 before Key 2.
	// Remove KeyIDs from non-written entities before beginning.
	key1.ID = ""
	issuer1.ID = ""
	issuer1.KeyID = ""

	key1Ref1, existing, err := sc.importKey(key1.PrivateKey, "key1", key1.PrivateKeyType)
	require.NoError(t, err)
	require.False(t, existing)
	require.Equal(t, strings.TrimSpace(key1.PrivateKey), strings.TrimSpace(key1Ref1.PrivateKey))

	// Make sure if we attempt to re-import the same private key, no import/updates occur.
	// So the existing flag should be set to true, and we do not update the existing Name field.
	key1Ref2, existing, err := sc.importKey(key1.PrivateKey, "ignore-me", key1.PrivateKeyType)
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, key1.PrivateKey, key1Ref1.PrivateKey)
	require.Equal(t, key1Ref1.ID, key1Ref2.ID)
	require.Equal(t, key1Ref1.Name, key1Ref2.Name)

	issuer1Ref1, existing, err := sc.importIssuer(issuer1.Certificate, "issuer1")
	require.NoError(t, err)
	require.False(t, existing)
	require.Equal(t, strings.TrimSpace(issuer1.Certificate), strings.TrimSpace(issuer1Ref1.Certificate))
	require.Equal(t, key1Ref1.ID, issuer1Ref1.KeyID)
	require.Equal(t, "issuer1", issuer1Ref1.Name)

	// Make sure if we attempt to re-import the same issuer, no import/updates occur.
	// So the existing flag should be set to true, and we do not update the existing Name field.
	issuer1Ref2, existing, err := sc.importIssuer(issuer1.Certificate, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, strings.TrimSpace(issuer1.Certificate), strings.TrimSpace(issuer1Ref1.Certificate))
	require.Equal(t, issuer1Ref1.ID, issuer1Ref2.ID)
	require.Equal(t, key1Ref1.ID, issuer1Ref2.KeyID)
	require.Equal(t, issuer1Ref1.Name, issuer1Ref2.Name)

	err = sc.writeIssuer(&issuer2)
	require.NoError(t, err)

	err = sc.writeKey(key2)
	require.NoError(t, err)

	// Same double import tests as above, but make sure if the previous was created through writeIssuer not importIssuer.
	issuer2Ref, existing, err := sc.importIssuer(issuer2.Certificate, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, strings.TrimSpace(issuer2.Certificate), strings.TrimSpace(issuer2Ref.Certificate))
	require.Equal(t, issuer2.ID, issuer2Ref.ID)
	require.Equal(t, "", issuer2Ref.Name)
	require.Equal(t, issuer2.KeyID, issuer2Ref.KeyID)

	// Same double import tests as above, but make sure if the previous was created through writeKey not importKey.
	key2Ref, existing, err := sc.importKey(key2.PrivateKey, "ignore-me", key2.PrivateKeyType)
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, key2.PrivateKey, key2Ref.PrivateKey)
	require.Equal(t, key2.ID, key2Ref.ID)
	require.Equal(t, "", key2Ref.Name)
}

func Test_IssuerUpgrade(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Make sure that we add OCSP signing to v0 issuers if CRLSigning is enabled
	issuer, _ := genIssuerAndKey(t, b, s)
	issuer.Version = 0
	issuer.Usage.ToggleUsage(OCSPSigningUsage)

	err := sc.writeIssuer(&issuer)
	require.NoError(t, err, "failed writing out issuer")

	newIssuer, err := sc.fetchIssuerById(issuer.ID)
	require.NoError(t, err, "failed fetching issuer")

	require.Equal(t, uint(1), newIssuer.Version)
	require.True(t, newIssuer.Usage.HasUsage(OCSPSigningUsage))

	// If CRLSigning is not present on a v0, we should not have OCSP signing after upgrade.
	issuer, _ = genIssuerAndKey(t, b, s)
	issuer.Version = 0
	issuer.Usage.ToggleUsage(OCSPSigningUsage)
	issuer.Usage.ToggleUsage(CRLSigningUsage)

	err = sc.writeIssuer(&issuer)
	require.NoError(t, err, "failed writing out issuer")

	newIssuer, err = sc.fetchIssuerById(issuer.ID)
	require.NoError(t, err, "failed fetching issuer")

	require.Equal(t, uint(1), newIssuer.Version)
	require.False(t, newIssuer.Usage.HasUsage(OCSPSigningUsage))
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
		Usage:        AllIssuerUsages,
		Version:      latestIssuerVersion,
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
	sc := b.makeStorageContext(ctx, s)
	_, _, role, respErr := getGenerationParams(sc, apiData)
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
	parsedCertBundle, _, err := generateCert(sc, input, nil, true, b.GetRandomReader())

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
