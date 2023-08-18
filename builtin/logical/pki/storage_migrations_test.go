// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	t.Parallel()
	startTime := time.Now()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Reset the version the helper above set to 1.
	b.pkiStorageVersion.Store(0)
	require.True(t, b.useLegacyBundleCaStorage(), "pre migration we should have been told to use legacy storage.")

	request := &logical.InitializationRequest{Storage: s}
	err := b.initialize(ctx, request)
	require.NoError(t, err)

	issuerIds, err := sc.listIssuers()
	require.NoError(t, err)
	require.Empty(t, issuerIds)

	keyIds, err := sc.listKeys()
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

func Test_migrateStorageOnlyKey(t *testing.T) {
	t.Parallel()
	startTime := time.Now()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Reset the version the helper above set to 1.
	b.pkiStorageVersion.Store(0)
	require.True(t, b.useLegacyBundleCaStorage(), "pre migration we should have been told to use legacy storage.")

	bundle := genCertBundle(t, b, s)
	// Clear everything except for the key
	bundle.SerialNumber = ""
	bundle.CAChain = []string{}
	bundle.Certificate = ""
	bundle.IssuingCA = ""

	json, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
	require.NoError(t, err)
	err = s.Put(ctx, json)
	require.NoError(t, err)

	request := &logical.InitializationRequest{Storage: s}
	err = b.initialize(ctx, request)
	require.NoError(t, err)

	issuerIds, err := sc.listIssuers()
	require.NoError(t, err)
	require.Equal(t, 0, len(issuerIds))

	keyIds, err := sc.listKeys()
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
	require.Equal(t, logEntry.CreatedIssuer, issuerID(""))
	require.Equal(t, logEntry.CreatedKey, keyIds[0])

	keyId := keyIds[0]
	key, err := sc.fetchKeyById(keyId)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(key.Name, "current-"),
		"expected key name to start with current- was %s", key.Name)
	require.Equal(t, keyId, key.ID)
	require.Equal(t, strings.TrimSpace(bundle.PrivateKey), strings.TrimSpace(key.PrivateKey))
	require.Equal(t, bundle.PrivateKeyType, key.PrivateKeyType)

	// Make sure we kept the old bundle
	_, certBundle, err := getLegacyCertBundle(ctx, s)
	require.NoError(t, err)
	require.Equal(t, bundle, certBundle)

	// Make sure we setup the default values
	keysConfig, err := sc.getKeysConfig()
	require.NoError(t, err)
	require.Equal(t, &keyConfigEntry{DefaultKeyId: keyId}, keysConfig)

	issuersConfig, err := sc.getIssuersConfig()
	require.NoError(t, err)
	require.Equal(t, issuerID(""), issuersConfig.DefaultIssuerId)

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

func Test_migrateStorageSimpleBundle(t *testing.T) {
	t.Parallel()
	startTime := time.Now()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

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

	issuerIds, err := sc.listIssuers()
	require.NoError(t, err)
	require.Equal(t, 1, len(issuerIds))

	keyIds, err := sc.listKeys()
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
	issuer, err := sc.fetchIssuerById(issuerId)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(issuer.Name, "current-"),
		"expected issuer name to start with current- was %s", issuer.Name)
	require.Equal(t, certutil.ErrNotAfterBehavior, issuer.LeafNotAfterBehavior)

	key, err := sc.fetchKeyById(keyId)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(key.Name, "current-"),
		"expected key name to start with current- was %s", key.Name)

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
	keysConfig, err := sc.getKeysConfig()
	require.NoError(t, err)
	require.Equal(t, &keyConfigEntry{DefaultKeyId: keyId}, keysConfig)

	issuersConfig, err := sc.getIssuersConfig()
	require.NoError(t, err)
	require.Equal(t, issuerId, issuersConfig.DefaultIssuerId)

	// Make sure if we attempt to re-run the migration nothing happens...
	err = migrateStorage(ctx, b, s)
	require.NoError(t, err)
	logEntry2, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry2)

	require.Equal(t, logEntry.Created, logEntry2.Created)
	require.Equal(t, logEntry.Hash, logEntry2.Hash)

	require.False(t, b.useLegacyBundleCaStorage(), "post migration we are still told to use legacy storage")

	// Make sure we can re-process a migration from scratch for whatever reason
	err = s.Delete(ctx, legacyMigrationBundleLogKey)
	require.NoError(t, err)

	err = migrateStorage(ctx, b, s)
	require.NoError(t, err)

	logEntry3, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry3)

	require.NotEqual(t, logEntry.Created, logEntry3.Created)
	require.Equal(t, logEntry.Hash, logEntry3.Hash)
}

func TestMigration_OnceChainRebuild(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Create a legacy CA bundle that we'll migrate to the new layout. We call
	// ToParsedCertBundle just to make sure it works and to populate
	// bundle.SerialNumber for us.
	bundle := &certutil.CertBundle{
		PrivateKeyType: certutil.RSAPrivateKey,
		Certificate:    migIntCA,
		IssuingCA:      migRootCA,
		CAChain:        []string{migRootCA},
		PrivateKey:     migIntPrivKey,
	}
	_, err := bundle.ToParsedCertBundle()
	require.NoError(t, err)
	writeLegacyBundle(t, b, s, bundle)

	// Do an initial migration. Ensure we end up at least on version 2.
	request := &logical.InitializationRequest{Storage: s}
	err = b.initialize(ctx, request)
	require.NoError(t, err)

	issuerIds, err := sc.listIssuers()
	require.NoError(t, err)
	require.Equal(t, 2, len(issuerIds))

	keyIds, err := sc.listKeys()
	require.NoError(t, err)
	require.Equal(t, 1, len(keyIds))

	logEntry, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry)
	require.GreaterOrEqual(t, logEntry.MigrationVersion, 2)
	require.GreaterOrEqual(t, latestMigrationVersion, 2)

	// Verify the chain built correctly: current should have a CA chain of
	// length two.
	//
	// Afterwards, we mutate these issuers to only point at themselves and
	// write back out.
	var rootIssuerId issuerID
	var intIssuerId issuerID
	for _, issuerId := range issuerIds {
		issuer, err := sc.fetchIssuerById(issuerId)
		require.NoError(t, err)
		require.NotNil(t, issuer)

		if strings.HasPrefix(issuer.Name, "current-") {
			require.Equal(t, 2, len(issuer.CAChain))
			require.Equal(t, migIntCA, issuer.CAChain[0])
			require.Equal(t, migRootCA, issuer.CAChain[1])
			intIssuerId = issuerId

			issuer.CAChain = []string{migIntCA}
			err = sc.writeIssuer(issuer)
			require.NoError(t, err)
		} else {
			require.Equal(t, 1, len(issuer.CAChain))
			require.Equal(t, migRootCA, issuer.CAChain[0])
			rootIssuerId = issuerId
		}
	}

	// Reset our migration version back to one, as if this never
	// happened...
	logEntry.MigrationVersion = 1
	err = setLegacyBundleMigrationLog(ctx, s, logEntry)
	require.NoError(t, err)
	b.pkiStorageVersion.Store(1)

	// Re-attempt the migration by reinitializing the mount.
	err = b.initialize(ctx, request)
	require.NoError(t, err)

	newIssuerIds, err := sc.listIssuers()
	require.NoError(t, err)
	require.Equal(t, 2, len(newIssuerIds))
	require.Equal(t, issuerIds, newIssuerIds)

	newKeyIds, err := sc.listKeys()
	require.NoError(t, err)
	require.Equal(t, 1, len(newKeyIds))
	require.Equal(t, keyIds, newKeyIds)

	logEntry, err = getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	require.NotNil(t, logEntry)
	require.Equal(t, logEntry.MigrationVersion, latestMigrationVersion)

	// Ensure the chains are correct on the intermediate. By using the
	// issuerId saved above, this ensures we didn't change any issuerIds,
	// we merely updated the existing issuers.
	intIssuer, err := sc.fetchIssuerById(intIssuerId)
	require.NoError(t, err)
	require.NotNil(t, intIssuer)
	require.Equal(t, 2, len(intIssuer.CAChain))
	require.Equal(t, migIntCA, intIssuer.CAChain[0])
	require.Equal(t, migRootCA, intIssuer.CAChain[1])

	rootIssuer, err := sc.fetchIssuerById(rootIssuerId)
	require.NoError(t, err)
	require.NotNil(t, rootIssuer)
	require.Equal(t, 1, len(rootIssuer.CAChain))
	require.Equal(t, migRootCA, rootIssuer.CAChain[0])
}

func TestExpectedOpsWork_PreMigration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
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
	require.NotNil(t, resp, "got nil response object from creating role")

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
	require.NotNil(t, resp, "got nil response setting CRL config")

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
	requireSuccessNonNilResponse(t, resp, err)

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

func TestBackupBundle(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Reset the version the helper above set to 1.
	b.pkiStorageVersion.Store(0)
	require.True(t, b.useLegacyBundleCaStorage(), "pre migration we should have been told to use legacy storage.")

	// Create an empty request and tidy configuration for us.
	req := &logical.Request{
		Storage:    s,
		MountPoint: "pki/",
	}
	cfg := &tidyConfig{
		BackupBundle:       true,
		IssuerSafetyBuffer: 120 * time.Second,
	}

	// Migration should do nothing if we're on an empty mount.
	err := b.doTidyMoveCABundle(ctx, req, b.Logger(), cfg)
	require.NoError(t, err)
	requireFileNotExists(t, sc, legacyCertBundlePath)
	requireFileNotExists(t, sc, legacyCertBundleBackupPath)
	issuerIds, err := sc.listIssuers()
	require.NoError(t, err)
	require.Empty(t, issuerIds)
	keyIds, err := sc.listKeys()
	require.NoError(t, err)
	require.Empty(t, keyIds)

	// Create a legacy CA bundle and write it out.
	bundle := genCertBundle(t, b, s)
	json, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
	require.NoError(t, err)
	err = s.Put(ctx, json)
	require.NoError(t, err)
	legacyContents := requireFileExists(t, sc, legacyCertBundlePath, nil)

	// Doing another tidy should maintain the status quo since we've
	// still not done our migration.
	err = b.doTidyMoveCABundle(ctx, req, b.Logger(), cfg)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, legacyContents)
	requireFileNotExists(t, sc, legacyCertBundleBackupPath)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.Empty(t, issuerIds)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.Empty(t, keyIds)

	// Do a migration; this should provision an issuer and key.
	initReq := &logical.InitializationRequest{Storage: s}
	err = b.initialize(ctx, initReq)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, legacyContents)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)

	// Doing another tidy should maintain the status quo since we've
	// done our migration too recently relative to the safety buffer.
	err = b.doTidyMoveCABundle(ctx, req, b.Logger(), cfg)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, legacyContents)
	requireFileNotExists(t, sc, legacyCertBundleBackupPath)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)

	// Shortening our buffer should ensure the migration occurs, removing
	// the legacy bundle but creating the backup one.
	time.Sleep(2 * time.Second)
	cfg.IssuerSafetyBuffer = 1 * time.Second
	err = b.doTidyMoveCABundle(ctx, req, b.Logger(), cfg)
	require.NoError(t, err)
	requireFileNotExists(t, sc, legacyCertBundlePath)
	requireFileExists(t, sc, legacyCertBundleBackupPath, legacyContents)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)

	// A new initialization should do nothing.
	err = b.initialize(ctx, initReq)
	require.NoError(t, err)
	requireFileNotExists(t, sc, legacyCertBundlePath)
	requireFileExists(t, sc, legacyCertBundleBackupPath, legacyContents)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	require.Equal(t, len(issuerIds), 1)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)
	require.Equal(t, len(keyIds), 1)

	// Restoring the legacy bundles with new issuers should redo the
	// migration.
	newBundle := genCertBundle(t, b, s)
	json, err = logical.StorageEntryJSON(legacyCertBundlePath, newBundle)
	require.NoError(t, err)
	err = s.Put(ctx, json)
	require.NoError(t, err)
	newLegacyContents := requireFileExists(t, sc, legacyCertBundlePath, nil)

	// -> reinit
	err = b.initialize(ctx, initReq)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, newLegacyContents)
	requireFileExists(t, sc, legacyCertBundleBackupPath, legacyContents)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	require.Equal(t, len(issuerIds), 2)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)
	require.Equal(t, len(keyIds), 2)

	// -> when we tidy again, we'll overwrite the old backup with the new
	// one.
	time.Sleep(2 * time.Second)
	err = b.doTidyMoveCABundle(ctx, req, b.Logger(), cfg)
	require.NoError(t, err)
	requireFileNotExists(t, sc, legacyCertBundlePath)
	requireFileExists(t, sc, legacyCertBundleBackupPath, newLegacyContents)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)

	// Finally, restoring the legacy bundle and re-migrating should redo
	// the migration.
	err = s.Put(ctx, json)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, newLegacyContents)
	requireFileExists(t, sc, legacyCertBundleBackupPath, newLegacyContents)

	// -> overwrite the version and re-migrate
	logEntry, err := getLegacyBundleMigrationLog(ctx, s)
	require.NoError(t, err)
	logEntry.MigrationVersion = 0
	err = setLegacyBundleMigrationLog(ctx, s, logEntry)
	require.NoError(t, err)
	err = b.initialize(ctx, initReq)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, newLegacyContents)
	requireFileExists(t, sc, legacyCertBundleBackupPath, newLegacyContents)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	require.Equal(t, len(issuerIds), 2)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)
	require.Equal(t, len(keyIds), 2)

	// -> Re-tidy should remove the legacy one.
	time.Sleep(2 * time.Second)
	err = b.doTidyMoveCABundle(ctx, req, b.Logger(), cfg)
	require.NoError(t, err)
	requireFileNotExists(t, sc, legacyCertBundlePath)
	requireFileExists(t, sc, legacyCertBundleBackupPath, newLegacyContents)
	issuerIds, err = sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	keyIds, err = sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)
}

func TestDeletedIssuersPostMigration(t *testing.T) {
	// We want to simulate the following scenario:
	//
	// 1.10.x: -> Create a CA.
	// 1.11.0: -> Migrate to new issuer layout but version 1.
	//         -> Delete existing issuers, create new ones.
	// (now):  -> Migrate to version 2 layout, make sure we don't see
	//            re-migration.

	t.Parallel()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Reset the version the helper above set to 1.
	b.pkiStorageVersion.Store(0)
	require.True(t, b.useLegacyBundleCaStorage(), "pre migration we should have been told to use legacy storage.")

	// Create a legacy CA bundle and write it out.
	bundle := genCertBundle(t, b, s)
	json, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
	require.NoError(t, err)
	err = s.Put(ctx, json)
	require.NoError(t, err)
	legacyContents := requireFileExists(t, sc, legacyCertBundlePath, nil)

	// Do a migration; this should provision an issuer and key.
	initReq := &logical.InitializationRequest{Storage: s}
	err = b.initialize(ctx, initReq)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, legacyContents)
	issuerIds, err := sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, issuerIds)
	keyIds, err := sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, keyIds)

	// Hack: reset the version to 1, to simulate a pre-version-2 migration
	// log.
	info, err := getMigrationInfo(sc.Context, sc.Storage)
	require.NoError(t, err, "failed to read migration info")
	info.migrationLog.MigrationVersion = 1
	err = setLegacyBundleMigrationLog(sc.Context, sc.Storage, info.migrationLog)
	require.NoError(t, err, "failed to write migration info")

	// Now delete all issuers and keys and create some new ones.
	for _, issuerId := range issuerIds {
		deleted, err := sc.deleteIssuer(issuerId)
		require.True(t, deleted, "expected it to be deleted")
		require.NoError(t, err, "error removing issuer")
	}
	for _, keyId := range keyIds {
		deleted, err := sc.deleteKey(keyId)
		require.True(t, deleted, "expected it to be deleted")
		require.NoError(t, err, "error removing key")
	}
	emptyIssuers, err := sc.listIssuers()
	require.NoError(t, err)
	require.Empty(t, emptyIssuers)
	emptyKeys, err := sc.listKeys()
	require.NoError(t, err)
	require.Empty(t, emptyKeys)

	// Create a new issuer + key.
	bundle = genCertBundle(t, b, s)
	_, _, err = sc.writeCaBundle(bundle, "", "")
	require.NoError(t, err)

	// List which issuers + keys we currently have.
	postDeletionIssuers, err := sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, postDeletionIssuers)
	postDeletionKeys, err := sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, postDeletionKeys)

	// Now do another migration from 1->2. This should retain the newly
	// created issuers+keys, but not revive any deleted ones.
	err = b.initialize(ctx, initReq)
	require.NoError(t, err)
	requireFileExists(t, sc, legacyCertBundlePath, legacyContents)
	postMigrationIssuers, err := sc.listIssuers()
	require.NoError(t, err)
	require.NotEmpty(t, postMigrationIssuers)
	require.Equal(t, postMigrationIssuers, postDeletionIssuers, "regression failed: expected second migration from v1->v2 to not introduce new issuers")
	postMigrationKeys, err := sc.listKeys()
	require.NoError(t, err)
	require.NotEmpty(t, postMigrationKeys)
	require.Equal(t, postMigrationKeys, postDeletionKeys, "regression failed: expected second migration from v1->v2 to not introduce new keys")
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

func requireFileNotExists(t *testing.T, sc *storageContext, path string) {
	t.Helper()

	entry, err := sc.Storage.Get(sc.Context, path)
	require.NoError(t, err)
	if entry != nil {
		require.Empty(t, entry.Value)
	} else {
		require.Empty(t, entry)
	}
}

func requireFileExists(t *testing.T, sc *storageContext, path string, contents []byte) []byte {
	t.Helper()

	entry, err := sc.Storage.Get(sc.Context, path)
	require.NoError(t, err)
	require.NotNil(t, entry)
	require.NotEmpty(t, entry.Value)
	if contents != nil {
		require.Equal(t, entry.Value, contents)
	}
	return entry.Value
}

// Keys to simulate an intermediate CA mount with also-imported root (parent).
const (
	migIntPrivKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAqu88Jcct/EyT8gDF+jdWuAwFplvanQ7KXAO5at58G6Y39UUz
fwnMS3P3VRBUoV5BDX+13wI2ldskbTKITsl6IXBPXUz0sKrdEKzXRVY4D6P2JR7W
YO1IUytfTgR+3F4sotFNQB++3ivT66AYLW7lOkoa+5lxsPM/oJ82DOlD2uGtDVTU
gQy1zugMBgPDlj+8tB562J9MTIdcKe9JpYrN0eO+aHzhbfvaSpScU4aZBgkS0kDv
8G4FxVfrBSDeD/JjCWaC48rLdgei1YrY0NFuw/8p/nPfA9vf2AtHMsWZRSwukQfq
I5HhQu+0OHQy3NWqXaBPzJNu3HnpKLykPHW7sQIDAQABAoIBAHNJy/2G66MhWx98
Ggt7S4fyw9TCWx5XHXEWKfbEfFyBrXhF5kemqh2x5319+DamRaX/HwF8kqhcF6N2
06ygAzmOcFjzUI3fkB5xFPh1AHa8FYZP2DOjloZR2IPcUFv9QInINRwszSU31kUz
w1rRUtYPqUdM5Pt99Mo219O5eMSlGtPKXm09uDAR8ZPuUx4jwGw90pSgeRB1Bg7X
Dt3YXx3X+OOs3Hbir1VDLSqCuy825l6Kn79h3eB8LAi+FUwCBvnTqyOEWyH2XjgP
z+tbz7lwnhGeKtxUl6Jb3m3SHtXpylot/4fwPisRV/9vaEDhVjKTmySH1WM+TRNR
CQLCJekCgYEA3b67DBhAYsFFdUd/4xh4QhHBanOcepV1CwaRln+UUjw1618ZEsTS
DKb9IS72C+ukUusGhQqxjFJlhOdXeMXpEbnEUY3PlREevWwm3bVAxtoAVRcmkQyK
PM4Oj9ODi2z8Cds0NvEXdX69uVutcbvm/JRZr/dsERWcLsfwdV/QqYcCgYEAxVce
d4ylsqORLm0/gcLnEyB9zhEPwmiJe1Yj5sH7LhGZ6JtLCqbOJO4jXmIzCrkbGyuf
BA/U7klc6jSprkBMgYhgOIuaULuFJvtKzJUzoATGFqX4r8WJm2ZycXgooAwZq6SZ
ySXOuQe9V7hlpI0fJfNhw+/HIjivL1jrnjBoXwcCgYEAtTv6LLx1g0Frv5scj0Ok
pntUlei/8ADPlJ9dxp+nXj8P4rvrBkgPVX/2S3TSbJO/znWA8qP20TVW+/UIrRE0
mOQ37F/3VWKUuUT3zyUhOGVc+C7fupWBNolDpZG+ZepBZNzgJDeQcNuRvTmM3PQy
qiWl2AhlLuF2sVWA1q3lIWkCgYEAnuHWgNA3dE1nDWceE351hxvIzklEU/TQhAHF
o/uYHO5E6VdmoqvMG0W0KkCL8d046rZDMAUDHdrpOROvbcENF9lSBxS26LshqFH4
ViDmULanOgLk57f2Y6ynBZ6Frt4vKNe8jYuoFacale67vzFz251JoHSD8pSKz2cb
ROCal68CgYA51hKqvki4r5rmS7W/Yvc3x3Wc0wpDEHTgLMoH+EV7AffJ8dy0/+po
AHK0nnRU63++1JmhQczBR0yTI6PUyeegEBk/d5CgFlY7UJQMTFPsMsiuM0Xw5nAv
KMPykK01D28UAkUxhwF7CqFrwwEv9GislgjewbdF5Za176+EuMEwIw==
-----END RSA PRIVATE KEY-----
`
	migIntCA = `-----BEGIN CERTIFICATE-----
MIIDHTCCAgWgAwIBAgIUfxlNBmrI7jsgH2Sdle1nVTqn5YQwDQYJKoZIhvcNAQEL
BQAwEjEQMA4GA1UEAxMHUm9vdCBYMTAeFw0yMjExMDIxMjI2MjhaFw0yMjEyMDQx
MjI2NThaMBoxGDAWBgNVBAMTD0ludGVybWVkaWF0ZSBSMTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAKrvPCXHLfxMk/IAxfo3VrgMBaZb2p0OylwDuWre
fBumN/VFM38JzEtz91UQVKFeQQ1/td8CNpXbJG0yiE7JeiFwT11M9LCq3RCs10VW
OA+j9iUe1mDtSFMrX04EftxeLKLRTUAfvt4r0+ugGC1u5TpKGvuZcbDzP6CfNgzp
Q9rhrQ1U1IEMtc7oDAYDw5Y/vLQeetifTEyHXCnvSaWKzdHjvmh84W372kqUnFOG
mQYJEtJA7/BuBcVX6wUg3g/yYwlmguPKy3YHotWK2NDRbsP/Kf5z3wPb39gLRzLF
mUUsLpEH6iOR4ULvtDh0MtzVql2gT8yTbtx56Si8pDx1u7ECAwEAAaNjMGEwDgYD
VR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFFusWj3piAiY
CR7tszR6uNYSMLe2MB8GA1UdIwQYMBaAFMNRNkLozstIhNhXCefi+WnaQApbMA0G
CSqGSIb3DQEBCwUAA4IBAQCmH852E/pDGBhf2VI1JAPZy9VYaRkKoqn4+5R1Gnoq
b90zhdCGueIm/usC1wAa0OOn7+xdQXFNfeI8UUB9w10q0QnM/A/G2v8UkdlLPPQP
zPjIYLalOOIOHf8hU2O5lwj0IA4JwjwDQ4xj69eX/N+x2LEI7SHyVVUZWAx0Y67a
QdyubpIJZlW/PI7kMwGyTx3tdkZxk1nTNtf/0nKvNuXKKcVzBCEMfvXyx4LFEM+U
nc2vdWN7PAoXcjUbxD3ZNGinr7mSBpQg82+nur/8yuSwu6iHomnfGxjUsEHic2GC
ja9siTbR+ONvVb4xUjugN/XmMSSaZnxig2vM9xcV8OMG
-----END CERTIFICATE-----
`
	migRootCA = `-----BEGIN CERTIFICATE-----
MIIDFTCCAf2gAwIBAgIURDTnXp8u78jWMe770Jj6Ac1paxkwDQYJKoZIhvcNAQEL
BQAwEjEQMA4GA1UEAxMHUm9vdCBYMTAeFw0yMjExMDIxMjI0NTVaFw0yMjEyMDQx
MjI1MjRaMBIxEDAOBgNVBAMTB1Jvb3QgWDEwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQC/+dh/o1qKTOua/OkHRMIvHiyBxjjoqrLqFSBYhjYKs+alA0qS
lLVzNqIKU8jm3fT73orx7yk/6acWaEYv/6owMaUn51xwS3gQhTHdFR/fLJwXnu2O
PZNqAs6tjAM3Q08aqR0qfxnjDvcgO7TOWSyOvVT2cTRK+uKYzxJEY52BDMUbp+iC
WJdXca9UwKRzi2wFqGliDycYsBBt/tr8tHSbTSZ5Qx6UpFrKpjZn+sT5KhKUlsdd
BYFmRegc0wXq4/kRjum0oEUigUMlHADIEhRasyXPEKa19sGP8nAZfo/hNOusGhj7
z7UPA0Cbe2uclpYPxsKgvcqQmgKugqKLL305AgMBAAGjYzBhMA4GA1UdDwEB/wQE
AwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTDUTZC6M7LSITYVwnn4vlp
2kAKWzAfBgNVHSMEGDAWgBTDUTZC6M7LSITYVwnn4vlp2kAKWzANBgkqhkiG9w0B
AQsFAAOCAQEAu7qdM1Li6V6iDCPpLg5zZReRtcxhUdwb5Xn4sDa8GJCy35f1voew
n0TQgM3Uph5x/djCR/Sj91MyAJ1/Q1PQQTyKGyUjSHvkcOBg628IAnLthn8Ua1fL
oQC/F/mlT1Yv+/W8eNPtD453/P0z8E0xMT5K3kpEDW/6K9RdHZlDJMW/z3UJ+4LN
6ONjIBmgffmLz9sVMpgCFyL7+w3W01bGP7w5AfKj2duoVG/Ekf2yUwmm6r9NgTQ1
oke0ShbZuMocwO8anq7k0R42FoluH3ipv9Qzzhsy+KdK5/fW5oqy1tKFaZsc67Q6
0UmD9DiDpCtn2Wod3nwxn0zW5HvDAWuDwg==
-----END CERTIFICATE-----
`
)
