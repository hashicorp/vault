package pki

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// This allows us to record the version of the migration code within the log entry
// in case we find out in the future that something was horribly wrong with the migration,
// and we need to perform it again...
const latestMigrationVersion = 1

func migrateStorage(ctx context.Context, req *logical.InitializationRequest, logger log.Logger) error {
	s := req.Storage
	legacyBundle, err := getLegacyCertBundle(ctx, s)
	if err != nil {
		return err
	}

	if legacyBundle == nil {
		// No legacy certs to migrate, we are done...
		logger.Debug("No legacy certs found, no migration required.")
		return nil
	}

	migrationEntry, err := getLegacyBundleMigrationLog(ctx, s)
	if err != nil {
		return err
	}
	hash, err := computeHashOfLegacyBundle(legacyBundle)
	if err != nil {
		return err
	}

	if migrationEntry != nil {
		// At this point we have already migrated something previously.
		if migrationEntry.hash == hash &&
			migrationEntry.migrationVersion == latestMigrationVersion {
			// The hashes are the same, no need to try and re-import...
			logger.Debug("existing migration hash found and matched legacy bundle, skipping migration.")
			return nil
		}
	}

	logger.Warn("performing PKI migration to new keys/issuers layout")

	anIssuer, aKey, err := writeCaBundle(ctx, s, legacyBundle, "current", "current")
	if err != nil {
		return err
	}
	logger.Info("Migration generated the following ids and set them as defaults",
		"issuer id", anIssuer.ID, "key id", aKey.ID)

	err = setLegacyBundleMigrationLog(ctx, s, &legacyBundleMigration{
		hash:             hash,
		created:          time.Now(),
		migrationVersion: latestMigrationVersion,
	})
	if err != nil {
		return err
	}
	logger.Info("successfully completed migration to new keys/issuers layout")
	return nil
}

func computeHashOfLegacyBundle(bundle *certutil.CertBundle) (string, error) {
	// We only hash the main certificate and the certs within the CAChain,
	// assuming that any sort of change that occurred would have influenced one of those two fields.
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(bundle.Certificate)); err != nil {
		return "", err
	}
	for _, cert := range bundle.CAChain {
		if _, err := hasher.Write([]byte(cert)); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

type legacyBundleMigration struct {
	hash             string
	created          time.Time
	migrationVersion int
}

func getLegacyBundleMigrationLog(ctx context.Context, s logical.Storage) (*legacyBundleMigration, error) {
	entry, err := s.Get(ctx, legacyMigrationBundleLogKey)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	lbm := &legacyBundleMigration{}
	err = entry.DecodeJSON(lbm)
	if err != nil {
		// If we can't decode our bundle, lets scrap it and assume a blank value,
		// re-running the migration will at most bring back an older certificate/private key
		return nil, nil
	}
	return lbm, nil
}

func setLegacyBundleMigrationLog(ctx context.Context, s logical.Storage, lbm *legacyBundleMigration) error {
	json, err := logical.StorageEntryJSON(legacyMigrationBundleLogKey, lbm)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func getLegacyCertBundle(ctx context.Context, s logical.Storage) (*certutil.CertBundle, error) {
	entry, err := s.Get(ctx, legacyCertBundlePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	cb := &certutil.CertBundle{}
	err = entry.DecodeJSON(cb)
	if err != nil {
		return nil, err
	}

	return cb, nil
}
