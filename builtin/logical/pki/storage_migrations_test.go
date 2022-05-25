package pki

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func Test_migrateStorageEmptyStorage(t *testing.T) {
	startTime := time.Now()
	ctx := context.Background()
	b, s := createBackendWithStorage(t)

	// Reset the version the helper above set to 1.
	b.pkiStorageVersion.Store(0)
	require.True(t, b.useLegacyBundleCaStorage(), "pre migration we should have been told to use legacy storage.")

	request := &logical.InitializationRequest{Storage: s}
	err := b.initialize(ctx, request)
	require.NoError(t, err)

	issuerIds, err := listIssuers(ctx, s)
	require.NoError(t, err)
	require.Empty(t, issuerIds)

	keyIds, err := listKeys(ctx, s)
	require.NoError(t, err)
	require.Empty(t, keyIds)

	logEntry, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry)
	require.Equal(t, latestMigrationVersion, logEntry.MigrationVersion)
	require.True(t, len(strings.TrimSpace(logEntry.Hash)) > 0,
		"Hash value (%s) should not have been empty", logEntry.Hash)
	require.True(t, startTime.Before(logEntry.Created),
		"created log entry time (%v) was before our start time(%v)?", logEntry.Created, startTime)
	require.Empty(t, logEntry.CreatedIssuer)
	require.Empty(t, logEntry.CreatedKey)

	require.False(t, b.useLegacyBundleCaStorage(), "post migration we are still told to use legacy storage")

	// Make sure we can re-run the migration without issues
	request = &logical.InitializationRequest{Storage: s}
	err = b.initialize(ctx, request)
	require.NoError(t, err)
	logEntry2, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry2)

	// Make sure the hash and created times have not changed.
	require.Equal(t, logEntry.Created, logEntry2.Created)
	require.Equal(t, logEntry.Hash, logEntry2.Hash)
}

func Test_migrateStorageSimpleBundle(t *testing.T) {
	startTime := time.Now()
	ctx := context.Background()
	b, s := createBackendWithStorage(t)
	// Reset the version the helper above set to 1.
	b.pkiStorageVersion.Store(0)
	require.True(t, b.useLegacyBundleCaStorage(), "pre migration we should have been told to use legacy storage.")

	bundle := genCertBundle(t, b, s)
	json, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
	require.NoError(t, err)
	err = s.Put(ctx, json)
	require.NoError(t, err)

	request := &logical.InitializationRequest{Storage: s}
	err = b.initialize(ctx, request)
	require.NoError(t, err)
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
	require.Equal(t, logEntry.CreatedIssuer, issuerIds[0])
	require.Equal(t, logEntry.CreatedKey, keyIds[0])

	issuerId := issuerIds[0]
	keyId := keyIds[0]
	issuer, err := fetchIssuerById(ctx, s, issuerId)
	require.NoError(t, err)
	require.Equal(t, "current", issuer.Name) // RFC says we should import with Name=current
	require.Equal(t, certutil.ErrNotAfterBehavior, issuer.LeafNotAfterBehavior)

	key, err := fetchKeyById(ctx, s, keyId)
	require.NoError(t, err)
	require.Equal(t, "current", key.Name) // RFC says we should import with Name=current

	require.Equal(t, issuerId, issuer.ID)
	require.Equal(t, bundle.SerialNumber, issuer.SerialNumber)
	require.Equal(t, strings.TrimSpace(bundle.Certificate), strings.TrimSpace(issuer.Certificate))
	require.Equal(t, keyId, issuer.KeyID)
	require.Empty(t, issuer.ManualChain)
	require.Equal(t, []string{bundle.Certificate + "\n"}, issuer.CAChain)
	require.Equal(t, AllIssuerUsages, issuer.Usage)
	require.Equal(t, certutil.ErrNotAfterBehavior, issuer.LeafNotAfterBehavior)

	require.Equal(t, keyId, key.ID)
	require.Equal(t, strings.TrimSpace(bundle.PrivateKey), strings.TrimSpace(key.PrivateKey))
	require.Equal(t, bundle.PrivateKeyType, key.PrivateKeyType)

	// Make sure we kept the old bundle
	_, certBundle, err := getLegacyCertBundle(ctx, s)
	require.NoError(t, err)
	require.Equal(t, bundle, certBundle)

	// Make sure we setup the default values
	keysConfig, err := getKeysConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &keyConfigEntry{DefaultKeyId: keyId}, keysConfig)

	issuersConfig, err := getIssuersConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &issuerConfigEntry{DefaultIssuerId: issuerId}, issuersConfig)

	// Make sure if we attempt to re-run the migration nothing happens...
	err = migrateStorage(ctx, b, s)
	require.NoError(t, err)
	logEntry2, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry2)

	require.Equal(t, logEntry.Created, logEntry2.Created)
	require.Equal(t, logEntry.Hash, logEntry2.Hash)

	require.False(t, b.useLegacyBundleCaStorage(), "post migration we are still told to use legacy storage")
}

func TestExpectedOpsWork_PreMigration(t *testing.T) {
	ctx := context.Background()
	b, s := createBackendWithStorage(t)
	// Reset the version the helper above set to 1.
	b.pkiStorageVersion.Store(0)
	require.True(t, b.useLegacyBundleCaStorage(), "pre migration we should have been told to use legacy storage.")

	bundle := genCertBundle(t, b, s)
	json, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
	require.NoError(t, err)
	err = s.Put(ctx, json)
	require.NoError(t, err)

	// generate role
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/allow-all",
		Storage:   s,
		Data: map[string]interface{}{
			"allow_any_name": "true",
			"no_store":       "false",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error from creating role")
	require.Nil(t, resp, "got non-nil response object from creating role")

	// List roles
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ListOperation,
		Path:       "roles",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error from listing roles")
	require.NotNil(t, resp, "got nil response object from listing roles")
	require.False(t, resp.IsError(), "got error response from listing roles: %#v", resp)
	require.Contains(t, resp.Data["keys"], "allow-all", "failed to list our roles")

	// Read roles
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "roles/allow-all",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error from reading role")
	require.NotNil(t, resp, "got nil response object from reading role")
	require.False(t, resp.IsError(), "got error response from reading role: %#v", resp)
	require.NotEmpty(t, resp.Data, "data map should not have been empty of reading role")

	// Issue a cert from our legacy bundle.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issue/allow-all",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test.com",
			"ttl":         "60s",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error issue on allow-all")
	require.NotNil(t, resp, "got nil response object from issue allow-all")
	require.False(t, resp.IsError(), "got error response from issue on allow-all: %#v", resp)
	serialNum := resp.Data["serial_number"].(string)
	require.NotEmpty(t, serialNum)

	// Make sure we can list
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ListOperation,
		Path:       "certs",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error listing certs")
	require.NotNil(t, resp, "got nil response object from listing certs")
	require.False(t, resp.IsError(), "got error response from listing certs: %#v", resp)
	require.Contains(t, resp.Data["keys"], serialNum, "failed to list our cert")

	// Revoke the cert now.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "revoke",
		Storage:   s,
		Data: map[string]interface{}{
			"serial_number": serialNum,
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error revoking cert")
	require.NotNil(t, resp, "got nil response object from revoke cert")
	require.False(t, resp.IsError(), "got error response from revoke cert: %#v", resp)

	// Check our CRL includes the revoked cert.
	resp = requestCrlFromBackend(t, s, b)
	crl := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))
	requireSerialNumberInCRL(t, crl, serialNum)

	// Set CRL config
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/crl",
		Storage:   s,
		Data: map[string]interface{}{
			"expiry":  "72h",
			"disable": "false",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error setting CRL config")
	require.Nil(t, resp, "got non-nil response setting CRL config")

	// Set URL config
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/urls",
		Storage:   s,
		Data: map[string]interface{}{
			"ocsp_servers": []string{"https://localhost:8080"},
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error setting URL config")
	require.Nil(t, resp, "got non-nil response setting URL config")

	// Make sure we can fetch the old values...
	for _, path := range []string{"ca/pem", "ca_chain", "cert/" + serialNum, "cert/ca", "cert/crl", "cert/ca_chain", "config/crl", "config/urls"} {
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation:  logical.ReadOperation,
			Path:       path,
			Storage:    s,
			MountPoint: "pki/",
		})
		require.NoError(t, err, "error reading cert %s", path)
		require.NotNil(t, resp, "got nil response object from reading cert %s", path)
		require.False(t, resp.IsError(), "got error response from reading cert %s: %#v", path, resp)
	}

	// Sign CSR
	_, csr := generateTestCsr(t, certutil.ECPrivateKey, 224)
	for _, path := range []string{"sign/allow-all", "root/sign-intermediate", "sign-verbatim"} {
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      path,
			Storage:   s,
			Data: map[string]interface{}{
				"csr": csr,
			},
			MountPoint: "pki/",
		})
		require.NoError(t, err, "error signing csr from path %s", path)
		require.NotNil(t, resp, "got nil response object from path %s", path)
		require.NotEmpty(t, resp.Data, "data map response was empty from path %s", path)
	}

	// Sign self-issued
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   s,
		Data: map[string]interface{}{
			"certificate": csr,
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error signing csr from path root/sign-self-issued")
	require.NotNil(t, resp, "got nil response object from path root/sign-self-issued")
	require.NotEmpty(t, resp.Data, "data map response was empty from path root/sign-self-issued")

	// Delete Role
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "roles/allow-all",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error deleting role")
	require.Nil(t, resp, "got non-nil response object from deleting role")

	// Delete Root
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "root",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error deleting root")
	require.NotNil(t, resp, "got nil response object from deleting root")
	require.NotEmpty(t, resp.Warnings, "expected warnings set on delete root")

	///////////////////////////////
	// Legacy calls we expect to fail when in migration mode
	///////////////////////////////
	requireFailInMigration(t, b, s, logical.UpdateOperation, "config/ca")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "intermediate/generate/internal")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "intermediate/set-signed")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "root/generate/internal")

	///////////////////////////////
	// New apis should be unavailable
	///////////////////////////////
	requireFailInMigration(t, b, s, logical.ListOperation, "issuers")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "issuers/generate/root/internal")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "issuers/generate/intermediate/internal")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "issuers/import/cert")
	requireFailInMigration(t, b, s, logical.ReadOperation, "issuer/default/json")
	requireFailInMigration(t, b, s, logical.ReadOperation, "issuer/default/crl/pem")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "issuer/test-role")

	// The following calls work as they are shared handlers with existing paths.
	// requireFailInMigration(t, b, s, logical.UpdateOperation, "issuer/default/issue/test-role")
	// requireFailInMigration(t, b, s, logical.UpdateOperation, "issuer/default/sign/test-role")
	// requireFailInMigration(t, b, s, logical.UpdateOperation, "issuer/default/sign-verbatim")
	// requireFailInMigration(t, b, s, logical.UpdateOperation, "issuer/default/sign-self-issued")

	requireFailInMigration(t, b, s, logical.UpdateOperation, "root/replace")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "root/rotate/internal")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "intermediate/cross-sign")

	requireFailInMigration(t, b, s, logical.UpdateOperation, "config/issuers")
	requireFailInMigration(t, b, s, logical.ReadOperation, "config/issuers")

	requireFailInMigration(t, b, s, logical.ListOperation, "keys")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "keys/generate/internal")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "keys/import")
	requireFailInMigration(t, b, s, logical.ReadOperation, "key/default")
	requireFailInMigration(t, b, s, logical.UpdateOperation, "config/keys")
	requireFailInMigration(t, b, s, logical.ReadOperation, "config/keys")
}

// requireFailInMigration validate that we fail the operation with the appropriate error message to the end-user
func requireFailInMigration(t *testing.T, b *backend, s logical.Storage, operation logical.Operation, path string) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  operation,
		Path:       path,
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error from op:%s path:%s", operation, path)
	require.NotNil(t, resp, "got nil response from op:%s path:%s", operation, path)
	require.True(t, resp.IsError(), "error flag was not set from op:%s path:%s resp: %#v", operation, path, resp)
	require.Contains(t, resp.Error().Error(), "migration has completed",
		"error message did not contain migration test for op:%s path:%s resp: %#v", operation, path, resp)
}
