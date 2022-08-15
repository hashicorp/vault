package pki

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// This allows us to record the version of the migration code within the log entry
// in case we find out in the future that something was horribly wrong with the migration,
// and we need to perform it again...
const (
	latestMigrationVersion = 2
	legacyBundleShimID     = issuerID("legacy-entry-shim-id")
	legacyBundleShimKeyID  = keyID("legacy-entry-shim-key-id")
)

type legacyBundleMigrationLog struct {
	Hash             string    `json:"hash"`
	Created          time.Time `json:"created"`
	CreatedIssuer    issuerID  `json:"issuer_id"`
	CreatedKey       keyID     `json:"key_id"`
	MigrationVersion int       `json:"migrationVersion"`
}

type migrationInfo struct {
	isRequired       bool
	legacyBundle     *certutil.CertBundle
	legacyBundleHash string
	migrationLog     *legacyBundleMigrationLog
}

func getMigrationInfo(ctx context.Context, s logical.Storage) (migrationInfo, error) {
	info := migrationInfo{
		isRequired:       false,
		legacyBundle:     nil,
		legacyBundleHash: "",
		migrationLog:     nil,
	}

	var err error
	_, info.legacyBundle, err = getLegacyCertBundle(ctx, s)
	if err != nil {
		return info, err
	}

	info.migrationLog, err = getLegacyBundleMigrationLog(ctx, s)
	if err != nil {
		return info, err
	}

	info.legacyBundleHash, err = computeHashOfLegacyBundle(info.legacyBundle)
	if err != nil {
		return info, err
	}

	// Even if there isn't anything to migrate, we always want to write out the log entry
	// as that will trigger the secondary clusters to toggle/wake up
	if info.migrationLog == nil {
		// No migration information at all, set to v0
		info.isRequired = true
		info.migrationLog = &legacyBundleMigrationLog{MigrationVersion: 0}
	}

	if info.migrationLog.Hash != info.legacyBundleHash {
		// There was an existing migration log but the legacy bundle and the logged hash do not match
		// so restart the migration from scratch.
		info.isRequired = true
		info.migrationLog.MigrationVersion = 0
		return info, nil
	}

	if info.migrationLog.MigrationVersion != latestMigrationVersion {
		// We aren't at the latest version of the migration
		info.isRequired = true
		return info, nil
	}

	// no migration required
	return info, nil
}

func migrateStorage(ctx context.Context, b *backend, s logical.Storage) error {
	info, err := getMigrationInfo(ctx, s)
	if err != nil {
		return err
	}

	if !info.isRequired {
		// No migration was deemed to be required.
		return nil
	}

	sc := b.makeStorageContext(ctx, s)
	migrationLog := info.migrationLog

	if info.migrationLog == nil || info.migrationLog.MigrationVersion == 0 {
		issuerIdentifier, keyIdentifier, err := migrateLegacyBundle(sc, info)
		if err != nil {
			return err
		}
		migrationLog = &legacyBundleMigrationLog{
			Hash:             info.legacyBundleHash,
			Created:          time.Now(),
			CreatedIssuer:    issuerIdentifier,
			CreatedKey:       keyIdentifier,
			MigrationVersion: 1,
		}
	}

	if info.migrationLog.MigrationVersion <= 1 {
		err = migrateOcspIssuerUsage(sc)
		if err != nil {
			return err
		}
	}

	migrationLog.MigrationVersion = latestMigrationVersion

	// We always want to write out this log entry as the secondary clusters leverage this path to wake up
	// if they were upgraded prior to the primary cluster's migration occurred.
	return setLegacyBundleMigrationLog(ctx, s, migrationLog)
}

// Upgrade existing issuers to include the new OCSPSigningUsage added in 1.12 if
// it previously contained CRLSigningUsage.
func migrateOcspIssuerUsage(sc *storageContext) error {
	issuerIds, err := sc.listIssuers()
	if err != nil {
		return fmt.Errorf("failed listing issuers in ocsp issuer usage migration: %v", err)
	}
	for _, issuerId := range issuerIds {
		issuer, err := sc.fetchIssuerById(issuerId)
		if err != nil {
			return fmt.Errorf("failed fetching issuer in ocsp issuer usage migration: %v", err)
		}

		// If the existing issuer had CRL signing, grant them OCSP usage privileges if missing.
		if issuer.Usage.HasUsage(CRLSigningUsage) && !issuer.Usage.HasUsage(OCSPSigningUsage) {
			issuer.Usage.ToggleUsage(OCSPSigningUsage)
			err := sc.writeIssuer(issuer)
			if err != nil {
				return fmt.Errorf("failed updating issuer in ocsp issuer usage migration: %v", err)
			}
		}
	}

	return nil
}

// Migrate the existing CA bundle that contains 1 or more certificates along with a private key in a single
// storage entry, into the new separated format of a private key entry and 1 or more issuer entries.
func migrateLegacyBundle(sc *storageContext, migrationInfo migrationInfo) (issuerID, keyID, error) {
	if migrationInfo.legacyBundle != nil {
		// Generate a unique name for the migrated items in case things were to be re-migrated again
		// for some weird reason in the future...
		migrationName := fmt.Sprintf("current-%d", time.Now().Unix())

		sc.Backend.Logger().Info("performing PKI migration to new keys/issuers layout")
		anIssuer, aKey, err := sc.writeCaBundle(migrationInfo.legacyBundle, migrationName, migrationName)
		if err != nil {
			return "", "", err
		}
		sc.Backend.Logger().Info("Migration generated the following ids and set them as defaults",
			"issuer id", anIssuer.ID, "key id", aKey.ID)
		issuerIdentifier := anIssuer.ID
		keyIdentifier := aKey.ID

		// Since we do not have all the mount information available we must schedule
		// the CRL to be rebuilt at a later time.
		sc.Backend.crlBuilder.requestRebuildIfActiveNode(sc.Backend)

		return issuerIdentifier, keyIdentifier, nil
	}
	return "", "", nil
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
	issuer.Usage.ToggleUsage(IssuanceUsage, CRLSigningUsage, OCSPSigningUsage)

	return issuer, cb, nil
}
