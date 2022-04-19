package pki

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func Test_migrateStorageEmptyStorage(t *testing.T) {
	ctx := context.Background()
	b, s := createBackendWithStorage(t)
	request := &logical.InitializationRequest{Storage: s}

	err := migrateStorage(ctx, request, b.Logger())
	require.NoError(t, err)

	issuerIds, err := listIssuers(ctx, s)
	require.NoError(t, err)
	require.Empty(t, issuerIds)

	keyIds, err := listKeys(ctx, s)
	require.NoError(t, err)
	require.Empty(t, keyIds)

	logEntry, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.Nil(t, logEntry)
}

func Test_migrateStorageSimpleBundle(t *testing.T) {
	startTime := time.Now()
	ctx := context.Background()
	b, s := createBackendWithStorage(t)

	bundle := genCertBundle(t, b)
	json, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
	require.NoError(t, err)
	err = s.Put(ctx, json)
	require.NoError(t, err)

	request := &logical.InitializationRequest{Storage: s}

	err = migrateStorage(ctx, request, b.Logger())
	require.NoError(t, err)

	issuerIds, err := listIssuers(ctx, s)
	require.NoError(t, err)
	require.Equal(t, 1, len(issuerIds))

	keyIds, err := listKeys(ctx, s)
	require.NoError(t, err)
	require.Equal(t, 1, len(keyIds))

	logEntry, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry)
	require.Equal(t, latestMigrationVersion, logEntry.MigrationVersion)
	require.True(t, len(strings.TrimSpace(logEntry.Hash)) > 0,
		"Hash value (%s) should not have been empty", logEntry.Hash)
	require.True(t, startTime.Before(logEntry.Created),
		"created log entry time (%v) was before our start time(%v)?", logEntry.Created, startTime)

	issuerId := issuerIds[0]
	keyId := keyIds[0]
	issuer, err := fetchIssuerById(ctx, s, issuerId)
	require.NoError(t, err)
	require.Equal(t, "current", issuer.Name) // RFC says we should import with Name=current

	key, err := fetchKeyById(ctx, s, keyId)
	require.NoError(t, err)
	require.Equal(t, "current", key.Name) // RFC says we should import with Name=current

	require.Equal(t, issuerId, issuer.ID)
	require.Equal(t, bundle.SerialNumber, issuer.SerialNumber)
	require.Equal(t, bundle.Certificate, issuer.Certificate)
	require.Equal(t, keyId, issuer.KeyID)
	// FIXME: Add tests for CAChain...

	require.Equal(t, keyId, key.ID)
	require.Equal(t, bundle.PrivateKey, key.PrivateKey)
	require.Equal(t, bundle.PrivateKeyType, key.PrivateKeyType)

	// Make sure we kept the old bundle
	certBundle, err := getLegacyCertBundle(ctx, s)
	require.NoError(t, err)
	require.Equal(t, bundle, certBundle)

	// Make sure we setup the default values
	keysConfig, err := getKeysConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &keyConfig{DefaultKeyId: keyId}, keysConfig)

	issuersConfig, err := getIssuersConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &issuerConfig{DefaultIssuerId: issuerId}, issuersConfig)

	// Make sure if we attempt to re-run the migration nothing happens...
	err = migrateStorage(ctx, request, b.Logger())
	require.NoError(t, err)
	logEntry2, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry2)

	require.Equal(t, logEntry.Created, logEntry2.Created)
	require.Equal(t, logEntry.Hash, logEntry2.Hash)
}
