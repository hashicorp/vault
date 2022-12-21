package pki

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

var tidyCancelledError = errors.New("tidy operation cancelled")

type tidyConfig struct {
	Enabled            bool          `json:"enabled"`
	Interval           time.Duration `json:"interval_duration"`
	CertStore          bool          `json:"tidy_cert_store"`
	RevokedCerts       bool          `json:"tidy_revoked_certs"`
	IssuerAssocs       bool          `json:"tidy_revoked_cert_issuer_associations"`
	ExpiredIssuers     bool          `json:"tidy_expired_issuers"`
	BackupBundle       bool          `json:"tidy_move_legacy_ca_bundle"`
	SafetyBuffer       time.Duration `json:"safety_buffer"`
	IssuerSafetyBuffer time.Duration `json:"issuer_safety_buffer"`
	PauseDuration      time.Duration `json:"pause_duration"`
}

var defaultTidyConfig = tidyConfig{
	Enabled:            false,
	Interval:           12 * time.Hour,
	CertStore:          false,
	RevokedCerts:       false,
	IssuerAssocs:       false,
	ExpiredIssuers:     false,
	BackupBundle:       false,
	SafetyBuffer:       72 * time.Hour,
	IssuerSafetyBuffer: 365 * 24 * time.Hour,
	PauseDuration:      0 * time.Second,
}

func pathTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy$",
		Fields:  addTidyFields(map[string]*framework.FieldSchema{}),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                  b.pathTidyWrite,
				ForwardPerformanceStandby: true,
			},
		},
		HelpSynopsis:    pathTidyHelpSyn,
		HelpDescription: pathTidyHelpDesc,
	}
}

func pathTidyCancel(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy-cancel$",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                  b.pathTidyCancelWrite,
				ForwardPerformanceStandby: true,
			},
		},
		HelpSynopsis:    pathTidyCancelHelpSyn,
		HelpDescription: pathTidyCancelHelpDesc,
	}
}

func pathTidyStatus(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy-status$",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:                  b.pathTidyStatusRead,
				ForwardPerformanceStandby: true,
			},
		},
		HelpSynopsis:    pathTidyStatusHelpSyn,
		HelpDescription: pathTidyStatusHelpDesc,
	}
}

func pathConfigAutoTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/auto-tidy",
		Fields: addTidyFields(map[string]*framework.FieldSchema{
			"enabled": {
				Type:        framework.TypeBool,
				Description: `Set to true to enable automatic tidy operations.`,
			},
			"interval_duration": {
				Type:        framework.TypeDurationSecond,
				Description: `Interval at which to run an auto-tidy operation. This is the time between tidy invocations (after one finishes to the start of the next). Running a manual tidy will reset this duration.`,
				Default:     int(defaultTidyConfig.Interval / time.Second), // TypeDurationSecond currently requires the default to be an int.
			},
		}),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigAutoTidyRead,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigAutoTidyWrite,
				// Read more about why these flags are set in backend.go.
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},
		HelpSynopsis:    pathConfigAutoTidySyn,
		HelpDescription: pathConfigAutoTidyDesc,
	}
}

func (b *backend) pathTidyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	safetyBuffer := d.Get("safety_buffer").(int)
	tidyCertStore := d.Get("tidy_cert_store").(bool)
	tidyRevokedCerts := d.Get("tidy_revoked_certs").(bool) || d.Get("tidy_revocation_list").(bool)
	tidyRevokedAssocs := d.Get("tidy_revoked_cert_issuer_associations").(bool)
	tidyExpiredIssuers := d.Get("tidy_expired_issuers").(bool)
	tidyBackupBundle := d.Get("tidy_move_legacy_ca_bundle").(bool)
	issuerSafetyBuffer := d.Get("issuer_safety_buffer").(int)
	pauseDurationStr := d.Get("pause_duration").(string)
	pauseDuration := 0 * time.Second

	if safetyBuffer < 1 {
		return logical.ErrorResponse("safety_buffer must be greater than zero"), nil
	}

	if issuerSafetyBuffer < 1 {
		return logical.ErrorResponse("issuer_safety_buffer must be greater than zero"), nil
	}

	if pauseDurationStr != "" {
		var err error
		pauseDuration, err = time.ParseDuration(pauseDurationStr)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing pause_duration: %v", err)), nil
		}

		if pauseDuration < (0 * time.Second) {
			return logical.ErrorResponse("received invalid, negative pause_duration"), nil
		}
	}

	bufferDuration := time.Duration(safetyBuffer) * time.Second
	issuerBufferDuration := time.Duration(issuerSafetyBuffer) * time.Second

	// Manual run with constructed configuration.
	config := &tidyConfig{
		Enabled:            true,
		Interval:           0 * time.Second,
		CertStore:          tidyCertStore,
		RevokedCerts:       tidyRevokedCerts,
		IssuerAssocs:       tidyRevokedAssocs,
		ExpiredIssuers:     tidyExpiredIssuers,
		BackupBundle:       tidyBackupBundle,
		SafetyBuffer:       bufferDuration,
		IssuerSafetyBuffer: issuerBufferDuration,
		PauseDuration:      pauseDuration,
	}

	if !atomic.CompareAndSwapUint32(b.tidyCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	// Tests using framework will screw up the storage so make a locally
	// scoped req to hold a reference
	req = &logical.Request{
		Storage: req.Storage,
	}

	// Mark the last tidy operation as relatively recent, to ensure we don't
	// try to trigger the periodic function.
	b.tidyStatusLock.Lock()
	b.lastTidy = time.Now()
	b.tidyStatusLock.Unlock()

	// Kick off the actual tidy.
	b.startTidyOperation(req, config)

	resp := &logical.Response{}
	if !tidyCertStore && !tidyRevokedCerts && !tidyRevokedAssocs && !tidyExpiredIssuers && !tidyBackupBundle {
		resp.AddWarning("No targets to tidy; specify tidy_cert_store=true or tidy_revoked_certs=true or tidy_revoked_cert_issuer_associations=true or tidy_expired_issuers=true or tidy_move_legacy_ca_bundle=true to start a tidy operation.")
	} else {
		resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	}

	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

func (b *backend) startTidyOperation(req *logical.Request, config *tidyConfig) {
	go func() {
		atomic.StoreUint32(b.tidyCancelCAS, 0)
		defer atomic.StoreUint32(b.tidyCASGuard, 0)

		b.tidyStatusStart(config)

		// Don't cancel when the original client request goes away.
		ctx := context.Background()

		logger := b.Logger().Named("tidy")

		doTidy := func() error {
			if config.CertStore {
				if err := b.doTidyCertStore(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.RevokedCerts || config.IssuerAssocs {
				if err := b.doTidyRevocationStore(ctx, req, logger, config); err != nil {
					return err
				}
			}

			if config.ExpiredIssuers {
				if err := b.doTidyExpiredIssuers(ctx, req, logger, config); err != nil {
					return err
				}
			}

			if config.BackupBundle {
				if err := b.doTidyMoveCABundle(ctx, req, logger, config); err != nil {
					return err
				}
			}

			return nil
		}

		if err := doTidy(); err != nil {
			logger.Error("error running tidy", "error", err)
			b.tidyStatusStop(err)
		} else {
			b.tidyStatusStop(nil)

			// Since the tidy operation finished without an error, we don't
			// really want to start another tidy right away (if the interval
			// is too short). So mark the last tidy as now.
			b.tidyStatusLock.Lock()
			b.lastTidy = time.Now()
			b.tidyStatusLock.Unlock()
		}
	}()
}

func (b *backend) doTidyCertStore(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	serials, err := req.Storage.List(ctx, "certs/")
	if err != nil {
		return fmt.Errorf("error fetching list of certs: %w", err)
	}

	serialCount := len(serials)
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "cert_store_total_entries"}, float32(serialCount))
	for i, serial := range serials {
		b.tidyStatusMessage(fmt.Sprintf("Tidying certificate store: checking entry %d of %d", i, serialCount))
		metrics.SetGauge([]string{"secrets", "pki", "tidy", "cert_store_current_entry"}, float32(i))

		// Check for cancel before continuing.
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
			return tidyCancelledError
		}

		// Check for pause duration to reduce resource consumption.
		if config.PauseDuration > (0 * time.Second) {
			time.Sleep(config.PauseDuration)
		}

		certEntry, err := req.Storage.Get(ctx, "certs/"+serial)
		if err != nil {
			return fmt.Errorf("error fetching certificate %q: %w", serial, err)
		}

		if certEntry == nil {
			logger.Warn("certificate entry is nil; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
				return fmt.Errorf("error deleting nil entry with serial %s: %w", serial, err)
			}
			b.tidyStatusIncCertStoreCount()
			continue
		}

		if certEntry.Value == nil || len(certEntry.Value) == 0 {
			logger.Warn("certificate entry has no value; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
				return fmt.Errorf("error deleting entry with nil value with serial %s: %w", serial, err)
			}
			b.tidyStatusIncCertStoreCount()
			continue
		}

		cert, err := x509.ParseCertificate(certEntry.Value)
		if err != nil {
			return fmt.Errorf("unable to parse stored certificate with serial %q: %w", serial, err)
		}

		if time.Now().After(cert.NotAfter.Add(config.SafetyBuffer)) {
			if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
				return fmt.Errorf("error deleting serial %q from storage: %w", serial, err)
			}
			b.tidyStatusIncCertStoreCount()
		}
	}

	b.tidyStatusLock.RLock()
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "cert_store_total_entries_remaining"}, float32(uint(serialCount)-b.tidyStatus.certStoreDeletedCount))
	b.tidyStatusLock.RUnlock()

	return nil
}

func (b *backend) doTidyRevocationStore(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	b.revokeStorageLock.Lock()
	defer b.revokeStorageLock.Unlock()

	// Fetch and parse our issuers so we can associate them if necessary.
	sc := b.makeStorageContext(ctx, req.Storage)
	issuerIDCertMap, err := fetchIssuerMapForRevocationChecking(sc)
	if err != nil {
		return err
	}

	rebuildCRL := false

	revokedSerials, err := req.Storage.List(ctx, "revoked/")
	if err != nil {
		return fmt.Errorf("error fetching list of revoked certs: %w", err)
	}

	revokedSerialsCount := len(revokedSerials)
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_total_entries"}, float32(revokedSerialsCount))

	fixedIssuers := 0

	var revInfo revocationInfo
	for i, serial := range revokedSerials {
		b.tidyStatusMessage(fmt.Sprintf("Tidying revoked certificates: checking certificate %d of %d", i, len(revokedSerials)))
		metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_current_entry"}, float32(i))

		// Check for cancel before continuing.
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
			return tidyCancelledError
		}

		// Check for pause duration to reduce resource consumption.
		if config.PauseDuration > (0 * time.Second) {
			b.revokeStorageLock.Unlock()
			time.Sleep(config.PauseDuration)
			b.revokeStorageLock.Lock()
		}

		revokedEntry, err := req.Storage.Get(ctx, "revoked/"+serial)
		if err != nil {
			return fmt.Errorf("unable to fetch revoked cert with serial %q: %w", serial, err)
		}

		if revokedEntry == nil {
			logger.Warn("revoked entry is nil; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
				return fmt.Errorf("error deleting nil revoked entry with serial %s: %w", serial, err)
			}
			b.tidyStatusIncRevokedCertCount()
			continue
		}

		if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
			logger.Warn("revoked entry has nil value; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
				return fmt.Errorf("error deleting revoked entry with nil value with serial %s: %w", serial, err)
			}
			b.tidyStatusIncRevokedCertCount()
			continue
		}

		err = revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return fmt.Errorf("error decoding revocation entry for serial %q: %w", serial, err)
		}

		revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
		if err != nil {
			return fmt.Errorf("unable to parse stored revoked certificate with serial %q: %w", serial, err)
		}

		// Tidy operations over revoked certs should execute prior to
		// tidyRevokedCerts as that may remove the entry. If that happens,
		// we won't persist the revInfo changes (as it was deleted instead).
		var storeCert bool
		if config.IssuerAssocs {
			if !isRevInfoIssuerValid(&revInfo, issuerIDCertMap) {
				b.tidyStatusIncMissingIssuerCertCount()
				revInfo.CertificateIssuer = issuerID("")
				storeCert = true
				if associateRevokedCertWithIsssuer(&revInfo, revokedCert, issuerIDCertMap) {
					fixedIssuers += 1
				}
			}
		}

		if config.RevokedCerts {
			// Only remove the entries from revoked/ and certs/ if we're
			// past its NotAfter value. This is because we use the
			// information on revoked/ to build the CRL and the
			// information on certs/ for lookup.
			if time.Now().After(revokedCert.NotAfter.Add(config.SafetyBuffer)) {
				if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
					return fmt.Errorf("error deleting serial %q from revoked list: %w", serial, err)
				}
				if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
					return fmt.Errorf("error deleting serial %q from store when tidying revoked: %w", serial, err)
				}
				rebuildCRL = true
				storeCert = false
				b.tidyStatusIncRevokedCertCount()
			}
		}

		// If the entry wasn't removed but was otherwise modified,
		// go ahead and write it back out.
		if storeCert {
			revokedEntry, err = logical.StorageEntryJSON("revoked/"+serial, revInfo)
			if err != nil {
				return fmt.Errorf("error building entry to persist changes to serial %v from revoked list: %w", serial, err)
			}

			err = req.Storage.Put(ctx, revokedEntry)
			if err != nil {
				return fmt.Errorf("error persisting changes to serial %v from revoked list: %w", serial, err)
			}
		}
	}

	b.tidyStatusLock.RLock()
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_total_entries_remaining"}, float32(uint(revokedSerialsCount)-b.tidyStatus.revokedCertDeletedCount))
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_entries_incorrect_issuers"}, float32(b.tidyStatus.missingIssuerCertCount))
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_entries_fixed_issuers"}, float32(fixedIssuers))
	b.tidyStatusLock.RUnlock()

	if rebuildCRL {
		// Expired certificates isn't generally an important
		// reason to trigger a CRL rebuild for. Check if
		// automatic CRL rebuilds have been enabled and defer
		// the rebuild if so.
		config, err := sc.getRevocationConfig()
		if err != nil {
			return err
		}

		if !config.AutoRebuild {
			if err := b.crlBuilder.rebuild(sc, false); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *backend) doTidyExpiredIssuers(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		b.Logger().Debug("skipping expired issuer tidy as we're not on the primary or secondary with a local mount")
		return nil
	}

	// Short-circuit to avoid having to deal with the legacy mounts. While we
	// could handle this case and remove these issuers, its somewhat
	// unexpected behavior and we'd prefer to finish the migration first.
	if b.useLegacyBundleCaStorage() {
		return nil
	}

	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	// Fetch and parse our issuers so we have their expiration date.
	sc := b.makeStorageContext(ctx, req.Storage)
	issuerIDCertMap, err := fetchIssuerMapForRevocationChecking(sc)
	if err != nil {
		return err
	}

	// Fetch the issuer config to find the default; we don't want to remove
	// the current active issuer automatically.
	iConfig, err := sc.getIssuersConfig()
	if err != nil {
		return err
	}

	// We want certificates which have expired before this date by a given
	// safety buffer. So we subtract the buffer from now, and anything which
	// has expired before our after buffer can be tidied, and anything that
	// expired after this buffer must be kept.
	now := time.Now()
	afterBuffer := now.Add(-1 * config.IssuerSafetyBuffer)

	rebuildChainsAndCRL := false

	for issuer, cert := range issuerIDCertMap {
		if cert.NotAfter.After(afterBuffer) {
			continue
		}

		entry, err := sc.fetchIssuerById(issuer)
		if err != nil {
			return nil
		}

		// This issuer's certificate has expired. We explicitly persist the
		// key, but log both the certificate and the keyId to the
		// informational logs so an admin can recover the removed cert if
		// necessary or remove the key (and know which cert it belonged to),
		// if desired.
		msg := "[Tidy on mount: %v] Issuer %v has expired by %v and is being removed."
		idAndName := fmt.Sprintf("[id:%v/name:%v]", entry.ID, entry.Name)
		msg = fmt.Sprintf(msg, b.backendUUID, idAndName, config.IssuerSafetyBuffer)

		// Before we log, check if we're the default. While this is late, and
		// after we read it from storage, we have more info here to tell the
		// user that their default has expired AND has passed the safety
		// buffer.
		if iConfig.DefaultIssuerId == issuer {
			msg = "[Tidy on mount: %v] Issuer %v has expired and would be removed via tidy, but won't be, as it is currently the default issuer."
			msg = fmt.Sprintf(msg, b.backendUUID, idAndName)
			b.Logger().Warn(msg)
			continue
		}

		// Log the above message..
		b.Logger().Info(msg, "serial_number", entry.SerialNumber, "key_id", entry.KeyID, "certificate", entry.Certificate)

		wasDefault, err := sc.deleteIssuer(issuer)
		if err != nil {
			b.Logger().Error(fmt.Sprintf("failed to remove %v: %v", idAndName, err))
			return err
		}
		if wasDefault {
			b.Logger().Warn(fmt.Sprintf("expired issuer %v was default; it is strongly encouraged to choose a new default issuer for backwards compatibility", idAndName))
		}

		rebuildChainsAndCRL = true
	}

	if rebuildChainsAndCRL {
		// When issuers are removed, there's a chance chains change as a
		// result; remove them.
		if err := sc.rebuildIssuersChains(nil); err != nil {
			return err
		}

		// Removal of issuers is generally a good reason to rebuild the CRL,
		// even if auto-rebuild is enabled.
		b.revokeStorageLock.Lock()
		defer b.revokeStorageLock.Unlock()

		if err := b.crlBuilder.rebuild(sc, false); err != nil {
			return err
		}
	}

	return nil
}

func (b *backend) doTidyMoveCABundle(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		b.Logger().Debug("skipping moving the legacy CA bundle as we're not on the primary or secondary with a local mount")
		return nil
	}

	// Short-circuit to avoid moving the legacy bundle from under a legacy
	// mount.
	if b.useLegacyBundleCaStorage() {
		return nil
	}

	// If we've already run, exit.
	_, bundle, err := getLegacyCertBundle(ctx, req.Storage)
	if err != nil {
		return fmt.Errorf("failed to fetch the legacy CA bundle: %w", err)
	}

	if bundle == nil {
		b.Logger().Debug("No legacy CA bundle available; nothing to do.")
		return nil
	}

	log, err := getLegacyBundleMigrationLog(ctx, req.Storage)
	if err != nil {
		return fmt.Errorf("failed to fetch the legacy bundle migration log: %w", err)
	}

	if log == nil {
		return fmt.Errorf("refusing to tidy with an empty legacy migration log but present CA bundle: %w", err)
	}

	now := time.Now()
	afterBuffer := now.Add(-1 * config.IssuerSafetyBuffer)

	if log.Created.After(afterBuffer) {
		b.Logger().Debug("Migration was created too recently to remove the legacy bundle; refusing to move legacy CA bundle to backup location.")
		return nil
	}

	// Do the write before the delete.
	entry, err := logical.StorageEntryJSON(legacyCertBundleBackupPath, bundle)
	if err != nil {
		return fmt.Errorf("failed to create new backup storage entry: %w", err)
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to write new backup legacy CA bundle: %w", err)
	}

	err = req.Storage.Delete(ctx, legacyCertBundlePath)
	if err != nil {
		return fmt.Errorf("failed to remove old legacy CA bundle path: %w", err)
	}

	b.Logger().Info("legacy CA bundle successfully moved to backup location")

	return nil
}

func (b *backend) pathTidyCancelWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if atomic.LoadUint32(b.tidyCASGuard) == 0 {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation cannot be cancelled as none is currently running.")
		return resp, nil
	}

	// Grab the status lock before writing the cancel atomic. This lets us
	// update the status correctly as well, avoiding writing it if we're not
	// presently running.
	//
	// Unlock needs to occur prior to calling read.
	b.tidyStatusLock.Lock()
	if b.tidyStatus.state == tidyStatusStarted || atomic.LoadUint32(b.tidyCASGuard) == 1 {
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 0, 1) {
			b.tidyStatus.state = tidyStatusCancelling
		}
	}
	b.tidyStatusLock.Unlock()

	return b.pathTidyStatusRead(ctx, req, d)
}

func (b *backend) pathTidyStatusRead(_ context.Context, _ *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	b.tidyStatusLock.RLock()
	defer b.tidyStatusLock.RUnlock()

	resp := &logical.Response{
		Data: map[string]interface{}{
			"safety_buffer":                         nil,
			"issuer_safety_buffer":                  nil,
			"tidy_cert_store":                       nil,
			"tidy_revoked_certs":                    nil,
			"tidy_revoked_cert_issuer_associations": nil,
			"tidy_expired_issuers":                  nil,
			"tidy_move_legacy_ca_bundle":            nil,
			"pause_duration":                        nil,
			"state":                                 "Inactive",
			"error":                                 nil,
			"time_started":                          nil,
			"time_finished":                         nil,
			"message":                               nil,
			"cert_store_deleted_count":              nil,
			"revoked_cert_deleted_count":            nil,
			"missing_issuer_cert_count":             nil,
			"current_cert_store_count":              nil,
			"current_revoked_cert_count":            nil,
		},
	}

	if b.tidyStatus.state == tidyStatusInactive {
		return resp, nil
	}

	resp.Data["safety_buffer"] = b.tidyStatus.safetyBuffer
	resp.Data["issuer_safety_buffer"] = b.tidyStatus.issuerSafetyBuffer
	resp.Data["tidy_cert_store"] = b.tidyStatus.tidyCertStore
	resp.Data["tidy_revoked_certs"] = b.tidyStatus.tidyRevokedCerts
	resp.Data["tidy_revoked_cert_issuer_associations"] = b.tidyStatus.tidyRevokedAssocs
	resp.Data["tidy_expired_issuers"] = b.tidyStatus.tidyExpiredIssuers
	resp.Data["tidy_move_legacy_ca_bundle"] = b.tidyStatus.tidyBackupBundle
	resp.Data["pause_duration"] = b.tidyStatus.pauseDuration
	resp.Data["time_started"] = b.tidyStatus.timeStarted
	resp.Data["message"] = b.tidyStatus.message
	resp.Data["cert_store_deleted_count"] = b.tidyStatus.certStoreDeletedCount
	resp.Data["revoked_cert_deleted_count"] = b.tidyStatus.revokedCertDeletedCount
	resp.Data["missing_issuer_cert_count"] = b.tidyStatus.missingIssuerCertCount

	switch b.tidyStatus.state {
	case tidyStatusStarted:
		resp.Data["state"] = "Running"
	case tidyStatusFinished:
		resp.Data["state"] = "Finished"
		resp.Data["time_finished"] = b.tidyStatus.timeFinished
		resp.Data["message"] = nil
	case tidyStatusError:
		resp.Data["state"] = "Error"
		resp.Data["time_finished"] = b.tidyStatus.timeFinished
		resp.Data["error"] = b.tidyStatus.err.Error()
		// Don't clear the message so that it serves as a hint about when
		// the error occurred.
	case tidyStatusCancelling:
		resp.Data["state"] = "Cancelling"
	case tidyStatusCancelled:
		resp.Data["state"] = "Cancelled"
		resp.Data["time_finished"] = b.tidyStatus.timeFinished
	}

	resp.Data["current_cert_store_count"] = b.certCount
	resp.Data["current_revoked_cert_count"] = b.revokedCertCount

	if !b.certsCounted.Load() {
		resp.AddWarning("Certificates in storage are still being counted, current counts provided may be " +
			"inaccurate")
	}

	return resp, nil
}

func (b *backend) pathConfigAutoTidyRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := sc.getAutoTidyConfig()
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"enabled":                               config.Enabled,
			"interval_duration":                     int(config.Interval / time.Second),
			"tidy_cert_store":                       config.CertStore,
			"tidy_revoked_certs":                    config.RevokedCerts,
			"tidy_revoked_cert_issuer_associations": config.IssuerAssocs,
			"tidy_expired_issuers":                  config.ExpiredIssuers,
			"tidy_move_legacy_ca_bundle":            config.BackupBundle,
			"safety_buffer":                         int(config.SafetyBuffer / time.Second),
			"issuer_safety_buffer":                  int(config.IssuerSafetyBuffer / time.Second),
			"pause_duration":                        config.PauseDuration.String(),
		},
	}, nil
}

func (b *backend) pathConfigAutoTidyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := sc.getAutoTidyConfig()
	if err != nil {
		return nil, err
	}

	if enabledRaw, ok := d.GetOk("enabled"); ok {
		config.Enabled = enabledRaw.(bool)
	}

	if intervalRaw, ok := d.GetOk("interval_duration"); ok {
		config.Interval = time.Duration(intervalRaw.(int)) * time.Second
		if config.Interval < 0 {
			return logical.ErrorResponse(fmt.Sprintf("given interval_duration must be greater than or equal to zero seconds; got: %v", intervalRaw)), nil
		}
	}

	if certStoreRaw, ok := d.GetOk("tidy_cert_store"); ok {
		config.CertStore = certStoreRaw.(bool)
	}

	if revokedCertsRaw, ok := d.GetOk("tidy_revoked_certs"); ok {
		config.RevokedCerts = revokedCertsRaw.(bool)
	}

	if issuerAssocRaw, ok := d.GetOk("tidy_revoked_cert_issuer_associations"); ok {
		config.IssuerAssocs = issuerAssocRaw.(bool)
	}

	if safetyBufferRaw, ok := d.GetOk("safety_buffer"); ok {
		config.SafetyBuffer = time.Duration(safetyBufferRaw.(int)) * time.Second
		if config.SafetyBuffer < 1*time.Second {
			return logical.ErrorResponse(fmt.Sprintf("given safety_buffer must be at least one second; got: %v", safetyBufferRaw)), nil
		}
	}

	if pauseDurationRaw, ok := d.GetOk("pause_duration"); ok {
		config.PauseDuration, err = time.ParseDuration(pauseDurationRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("unable to parse given pause_duration: %v", err)), nil
		}

		if config.PauseDuration < (0 * time.Second) {
			return logical.ErrorResponse("received invalid, negative pause_duration"), nil
		}
	}

	if expiredIssuers, ok := d.GetOk("tidy_expired_issuers"); ok {
		config.ExpiredIssuers = expiredIssuers.(bool)
	}

	if issuerSafetyBufferRaw, ok := d.GetOk("issuer_safety_buffer"); ok {
		config.IssuerSafetyBuffer = time.Duration(issuerSafetyBufferRaw.(int)) * time.Second
		if config.IssuerSafetyBuffer < 1*time.Second {
			return logical.ErrorResponse(fmt.Sprintf("given safety_buffer must be at least one second; got: %v", issuerSafetyBufferRaw)), nil
		}
	}

	if backupBundle, ok := d.GetOk("tidy_move_legacy_ca_bundle"); ok {
		config.BackupBundle = backupBundle.(bool)
	}

	if config.Enabled && !(config.CertStore || config.RevokedCerts || config.IssuerAssocs || config.ExpiredIssuers || config.BackupBundle) {
		return logical.ErrorResponse("Auto-tidy enabled but no tidy operations were requested. Enable at least one tidy operation to be run (tidy_cert_store / tidy_revoked_certs / tidy_revoked_cert_issuer_associations / tidy_move_legacy_ca_bundle)."), nil
	}

	if err := sc.writeAutoTidyConfig(config); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"enabled":                               config.Enabled,
			"interval_duration":                     int(config.Interval / time.Second),
			"tidy_cert_store":                       config.CertStore,
			"tidy_revoked_certs":                    config.RevokedCerts,
			"tidy_revoked_cert_issuer_associations": config.IssuerAssocs,
			"tidy_expired_issuers":                  config.ExpiredIssuers,
			"tidy_move_legacy_ca_bundle":            config.BackupBundle,
			"safety_buffer":                         int(config.SafetyBuffer / time.Second),
			"issuer_safety_buffer":                  int(config.IssuerSafetyBuffer / time.Second),
			"pause_duration":                        config.PauseDuration.String(),
		},
	}, nil
}

func (b *backend) tidyStatusStart(config *tidyConfig) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus = &tidyStatus{
		safetyBuffer:       int(config.SafetyBuffer / time.Second),
		issuerSafetyBuffer: int(config.IssuerSafetyBuffer / time.Second),
		tidyCertStore:      config.CertStore,
		tidyRevokedCerts:   config.RevokedCerts,
		tidyRevokedAssocs:  config.IssuerAssocs,
		tidyExpiredIssuers: config.ExpiredIssuers,
		tidyBackupBundle:   config.BackupBundle,
		pauseDuration:      config.PauseDuration.String(),

		state:       tidyStatusStarted,
		timeStarted: time.Now(),
	}

	metrics.SetGauge([]string{"secrets", "pki", "tidy", "start_time_epoch"}, float32(b.tidyStatus.timeStarted.Unix()))
}

func (b *backend) tidyStatusStop(err error) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.timeFinished = time.Now()
	b.tidyStatus.err = err
	if err == nil {
		b.tidyStatus.state = tidyStatusFinished
	} else if err == tidyCancelledError {
		b.tidyStatus.state = tidyStatusCancelled
	} else {
		b.tidyStatus.state = tidyStatusError
	}

	metrics.MeasureSince([]string{"secrets", "pki", "tidy", "duration"}, b.tidyStatus.timeStarted)
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "start_time_epoch"}, 0)
	metrics.IncrCounter([]string{"secrets", "pki", "tidy", "cert_store_deleted_count"}, float32(b.tidyStatus.certStoreDeletedCount))
	metrics.IncrCounter([]string{"secrets", "pki", "tidy", "revoked_cert_deleted_count"}, float32(b.tidyStatus.revokedCertDeletedCount))

	if err != nil {
		metrics.IncrCounter([]string{"secrets", "pki", "tidy", "failure"}, 1)
	} else {
		metrics.IncrCounter([]string{"secrets", "pki", "tidy", "success"}, 1)
	}
}

func (b *backend) tidyStatusMessage(msg string) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.message = msg
}

func (b *backend) tidyStatusIncCertStoreCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.certStoreDeletedCount++

	b.decrementTotalCertificatesCountReport()
}

func (b *backend) tidyStatusIncRevokedCertCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.revokedCertDeletedCount++

	b.decrementTotalRevokedCertificatesCountReport()
}

func (b *backend) tidyStatusIncMissingIssuerCertCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.missingIssuerCertCount++
}

const pathTidyHelpSyn = `
Tidy up the backend by removing expired certificates, revocation information,
or both.
`

const pathTidyHelpDesc = `
This endpoint allows expired certificates and/or revocation information to be
removed from the backend, freeing up storage and shortening CRLs.

For safety, this function is a noop if called without parameters; cleanup from
normal certificate storage must be enabled with 'tidy_cert_store' and cleanup
from revocation information must be enabled with 'tidy_revocation_list'.

The 'safety_buffer' parameter is useful to ensure that clock skew amongst your
hosts cannot lead to a certificate being removed from the CRL while it is still
considered valid by other hosts (for instance, if their clocks are a few
minutes behind). The 'safety_buffer' parameter can be an integer number of
seconds or a string duration like "72h".

All certificates and/or revocation information currently stored in the backend
will be checked when this endpoint is hit. The expiration of the
certificate/revocation information of each certificate being held in
certificate storage or in revocation information will then be checked. If the
current time, minus the value of 'safety_buffer', is greater than the
expiration, it will be removed.
`

const pathTidyCancelHelpSyn = `
Cancels a currently running tidy operation.
`

const pathTidyCancelHelpDesc = `
This endpoint allows cancelling a currently running tidy operation.

Periodically throughout the invocation of tidy, we'll check if the operation
has been requested to be cancelled. If so, we'll stop the currently running
tidy operation.
`

const pathTidyStatusHelpSyn = `
Returns the status of the tidy operation.
`

const pathTidyStatusHelpDesc = `
This is a read only endpoint that returns information about the current tidy
operation, or the most recent if none is currently running.

The result includes the following fields:
* 'safety_buffer': the value of this parameter when initiating the tidy operation
* 'tidy_cert_store': the value of this parameter when initiating the tidy operation
* 'tidy_revoked_certs': the value of this parameter when initiating the tidy operation
* 'tidy_revoked_cert_issuer_associations': the value of this parameter when initiating the tidy operation
* 'state': one of "Inactive", "Running", "Finished", "Error"
* 'error': the error message, if the operation ran into an error
* 'time_started': the time the operation started
* 'time_finished': the time the operation finished
* 'message': One of "Tidying certificate store: checking entry N of TOTAL" or
  "Tidying revoked certificates: checking certificate N of TOTAL"
* 'cert_store_deleted_count': The number of certificate storage entries deleted
* 'revoked_cert_deleted_count': The number of revoked certificate entries deleted
* 'missing_issuer_cert_count': The number of revoked certificates which were missing a valid issuer reference
* 'tidy_expired_issuers': the value of this parameter when initiating the tidy operation
* 'issuer_safety_buffer': the value of this parameter when initiating the tidy operation
* 'tidy_move_legacy_ca_bundle': the value of this parameter when initiating the tidy operation
`

const pathConfigAutoTidySyn = `
Modifies the current configuration for automatic tidy execution.
`

const pathConfigAutoTidyDesc = `
This endpoint accepts parameters to a tidy operation (see /tidy) that
will be used for automatic tidy execution. This takes two extra parameters,
enabled (to enable or disable auto-tidy) and interval_duration (which
controls the frequency of auto-tidy execution).

Once enabled, a tidy operation will be kicked off automatically, as if it
were executed with the posted configuration.
`
