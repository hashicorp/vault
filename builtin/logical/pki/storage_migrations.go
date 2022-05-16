package pki

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// This allows us to record the version of the migration code within the log entry
// in case we find out in the future that something was horribly wrong with the migration,
// and we need to perform it again...
const (
	latestMigrationVersion = 1
	legacyBundleShimID     = issuerID("legacy-entry-shim-id")
	legacyBundleShimKeyID  = keyID("legacy-entry-shim-key-id")
)

type legacyBundleMigrationLog struct {
	Hash             string    `json:"hash" structs:"hash" mapstructure:"hash"`
	Created          time.Time `json:"created" structs:"created" mapstructure:"created"`
	CreatedIssuer    issuerID  `json:"issuer_id" structs:"issuer_id" mapstructure:"issuer_id"`
	CreatedKey       keyID     `json:"key_id" structs:"key_id" mapstructure:"key_id"`
	MigrationVersion int       `json:"migrationVersion" structs:"migrationVersion" mapstructure:"migrationVersion"`
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

	var issuerIdentifier issuerID
	var keyIdentifier keyID
	if migrationInfo.legacyBundle != nil {
		b.Logger().Info("performing PKI migration to new keys/issuers layout")
		anIssuer, aKey, err := writeCaBundle(ctx, b, s, migrationInfo.legacyBundle, "current", "current")
		if err != nil {
			return err
		}
		b.Logger().Info("Migration generated the following ids and set them as defaults",
			"issuer id", anIssuer.ID, "key id", aKey.ID)
		issuerIdentifier = anIssuer.ID
		keyIdentifier = aKey.ID

		// Since we do not have all the mount information available we must schedule
		// the CRL to be rebuilt at a later time.
		b.crlBuilder.requestRebuildIfActiveNode(b)
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

func getLegacyCertBundle(ctx context.Context, s logical.Storage) (*issuerEntry, *certutil.CertBundle, error) {
	entry, err := s.Get(ctx, legacyCertBundlePath)
	if err != nil {
		return nil, nil, err
	}

	if entry == nil {
		return nil, nil, nil
	}

	cb := &certutil.CertBundle{}
	err = entry.DecodeJSON(cb)
	if err != nil {
		return nil, nil, err
	}

	// Fake a storage entry with backwards compatibility in mind.
	issuer := &issuerEntry{
		ID:                   legacyBundleShimID,
		KeyID:                legacyBundleShimKeyID,
		Name:                 "legacy-entry-shim",
		Certificate:          cb.Certificate,
		CAChain:              cb.CAChain,
		SerialNumber:         cb.SerialNumber,
		LeafNotAfterBehavior: certutil.ErrNotAfterBehavior,
	}
	issuer.Usage.ToggleUsage(IssuanceUsage, CRLSigningUsage)

	return issuer, cb, nil
}
