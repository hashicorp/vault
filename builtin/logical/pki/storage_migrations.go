// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// This allows us to record the version of the migration code within the log entry
// in case we find out in the future that something was horribly wrong with the migration,
// and we need to perform it again...
const (
	latestMigrationVersion = 2
	legacyBundleShimID     = issuing.LegacyBundleShimID
	legacyBundleShimKeyID  = issuing.LegacyBundleShimKeyID
)

type legacyBundleMigrationLog struct {
	Hash             string           `json:"hash"`
	Created          time.Time        `json:"created"`
	CreatedIssuer    issuing.IssuerID `json:"issuer_id"`
	CreatedKey       issuing.KeyID    `json:"key_id"`
	MigrationVersion int              `json:"migrationVersion"`
}

type migrationInfo struct {
	isRequired       bool
	legacyBundle     *certutil.CertBundle
	legacyBundleHash string
	migrationLog     *legacyBundleMigrationLog
}

func getMigrationInfo(ctx context.Context, s logical.Storage) (migrationInfo, error) {
	migrationInfo := migrationInfo{
		isRequired:       false,
		legacyBundle:     nil,
		legacyBundleHash: "",
		migrationLog:     nil,
	}

	var err error
	_, migrationInfo.legacyBundle, err = getLegacyCertBundle(ctx, s)
	if err != nil {
		return migrationInfo, err
	}

	migrationInfo.migrationLog, err = getLegacyBundleMigrationLog(ctx, s)
	if err != nil {
		return migrationInfo, err
	}

	migrationInfo.legacyBundleHash, err = computeHashOfLegacyBundle(migrationInfo.legacyBundle)
	if err != nil {
		return migrationInfo, err
	}

	// Even if there isn't anything to migrate, we always want to write out the log entry
	// as that will trigger the secondary clusters to toggle/wake up
	if (migrationInfo.migrationLog == nil) ||
		(migrationInfo.migrationLog.Hash != migrationInfo.legacyBundleHash) ||
		(migrationInfo.migrationLog.MigrationVersion != latestMigrationVersion) {
		migrationInfo.isRequired = true
	}

	return migrationInfo, nil
}

func migrateStorage(ctx context.Context, b *backend, s logical.Storage) error {
	migrationInfo, err := getMigrationInfo(ctx, s)
	if err != nil {
		return err
	}

	if !migrationInfo.isRequired {
		// No migration was deemed to be required.
		return nil
	}

	var issuerIdentifier issuing.IssuerID
	var keyIdentifier issuing.KeyID
	sc := b.makeStorageContext(ctx, s)
	if migrationInfo.legacyBundle != nil {
		// When the legacy bundle still exists, there's three scenarios we
		// need to worry about:
		//
		// 1. When we have no migration log, we definitely want to migrate.
		haveNoLog := migrationInfo.migrationLog == nil
		// 2. When we have an (empty) log and the version is zero, we want to
		//    migrate.
		haveOldVersion := !haveNoLog && migrationInfo.migrationLog.MigrationVersion == 0
		// 3. When we have a log and the version is at least 1 (where this
		//    migration was introduced), we want to run the migration again
		//    only if the legacy bundle hash has changed.
		isCurrentOrBetterVersion := !haveNoLog && migrationInfo.migrationLog.MigrationVersion >= 1
		haveChange := !haveNoLog && migrationInfo.migrationLog.Hash != migrationInfo.legacyBundleHash
		haveVersionWithChange := isCurrentOrBetterVersion && haveChange

		if haveNoLog || haveOldVersion || haveVersionWithChange {
			// Generate a unique name for the migrated items in case things were to be re-migrated again
			// for some weird reason in the future...
			migrationName := fmt.Sprintf("current-%d", time.Now().Unix())

			b.Logger().Info("performing PKI migration to new keys/issuers layout")
			anIssuer, aKey, err := sc.writeCaBundle(migrationInfo.legacyBundle, migrationName, migrationName)
			if err != nil {
				return err
			}
			b.Logger().Info("Migration generated the following ids and set them as defaults",
				"issuer id", anIssuer.ID, "key id", aKey.ID)
			issuerIdentifier = anIssuer.ID
			keyIdentifier = aKey.ID

			// Since we do not have all the mount information available we must schedule
			// the CRL to be rebuilt at a later time.
			b.CrlBuilder().requestRebuildIfActiveNode(b)
		}
	}

	if migrationInfo.migrationLog != nil && migrationInfo.migrationLog.MigrationVersion == 1 {
		// We've seen a bundle with migration version 1; this means an
		// earlier version of the code ran which didn't have the fix for
		// correct write order in rebuildIssuersChains(...). Rather than
		// having every user read the migrated active issuer and see if
		// their chains need rebuilding, we'll schedule a one-off chain
		// migration here.
		b.Logger().Info(fmt.Sprintf("%v: performing maintenance rebuild of ca_chains", b.backendUUID))
		if err := sc.rebuildIssuersChains(nil); err != nil {
			return err
		}
	}

	// We always want to write out this log entry as the secondary clusters leverage this path to wake up
	// if they were upgraded prior to the primary cluster's migration occurred.
	err = setLegacyBundleMigrationLog(ctx, s, &legacyBundleMigrationLog{
		Hash:             migrationInfo.legacyBundleHash,
		Created:          time.Now(),
		CreatedIssuer:    issuerIdentifier,
		CreatedKey:       keyIdentifier,
		MigrationVersion: latestMigrationVersion,
	})
	if err != nil {
		return err
	}

	b.Logger().Info(fmt.Sprintf("%v: succeeded in migrating to issuer storage version %v", b.backendUUID, latestMigrationVersion))

	return nil
}

func computeHashOfLegacyBundle(bundle *certutil.CertBundle) (string, error) {
	hasher := sha256.New()
	// Generate an empty hash if the bundle does not exist.
	if bundle != nil {
		// We only hash the main certificate and the certs within the CAChain,
		// assuming that any sort of change that occurred would have influenced one of those two fields.
		if _, err := hasher.Write([]byte(bundle.Certificate)); err != nil {
			return "", err
		}
		for _, cert := range bundle.CAChain {
			if _, err := hasher.Write([]byte(cert)); err != nil {
				return "", err
			}
		}
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func getLegacyBundleMigrationLog(ctx context.Context, s logical.Storage) (*legacyBundleMigrationLog, error) {
	entry, err := s.Get(ctx, legacyMigrationBundleLogKey)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	lbm := &legacyBundleMigrationLog{}
	err = entry.DecodeJSON(lbm)
	if err != nil {
		// If we can't decode our bundle, lets scrap it and assume a blank value,
		// re-running the migration will at most bring back an older certificate/private key
		return nil, nil
	}
	return lbm, nil
}

func setLegacyBundleMigrationLog(ctx context.Context, s logical.Storage, lbm *legacyBundleMigrationLog) error {
	json, err := logical.StorageEntryJSON(legacyMigrationBundleLogKey, lbm)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func getLegacyCertBundle(ctx context.Context, s logical.Storage) (*issuing.IssuerEntry, *certutil.CertBundle, error) {
	return issuing.GetLegacyCertBundle(ctx, s)
}
