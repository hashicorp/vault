// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/x509"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/managed_key"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	operationPrefixPKI        = "pki"
	operationPrefixPKIIssuer  = "pki-issuer"
	operationPrefixPKIIssuers = "pki-issuers"
	operationPrefixPKIRoot    = "pki-root"

	noRole       = 0
	roleOptional = 1
	roleRequired = 2
)

/*
 * PKI requests are a bit special to keep up with the various failure and load issues.
 *
 * Any requests to write/delete shared data (such as roles, issuers, keys, and configuration)
 * are always forwarded to the Primary cluster's active node to write and send the key
 * material/config globally across all clusters. Reads should be handled locally, to give a
 * sense of where this cluster's replication state is at.
 *
 * CRL/Revocation and Fetch Certificate APIs are handled by the active node within the cluster
 * they originate. This means, if a request comes into a performance secondary cluster, the writes
 * will be forwarded to that cluster's active node and not go all the way up to the performance primary's
 * active node.
 *
 * If a certificate issue request has a role in which no_store and no_store_metadata is set to
 * true, that node itself will issue the certificate and not forward the request to the active
 * node, as this does not need to write to storage.
 *
 * Following the same pattern, if a managed key is involved to sign an issued certificate request
 * and the local node does not have access for some reason to it, the request will be forwarded to
 * the active node within the cluster only.
 *
 * To make sense of what goes where the following bits need to be analyzed within the codebase.
 *
 * 1. The backend LocalStorage paths determine what storage paths will remain within a
 *    cluster and not be forwarded to a performance primary
 * 2. Within each path's OperationHandler definition, check to see if ForwardPerformanceStandby &
 *    ForwardPerformanceSecondary flags are set to short-circuit the request to a given active node
 * 3. Within the managed key util class in pki, an initialization failure could cause the request
 *    to be forwarded to an active node if not already on it.
 */

// Factory creates a new backend implementing the logical.Backend interface
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns a new Backend framework struct
func Backend(conf *logical.BackendConfig) *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"cert/*",
				"ca/pem",
				"ca_chain",
				"ca",
				"crl/delta",
				"crl/delta/pem",
				"crl/pem",
				"crl",
				"issuer/+/crl/der",
				"issuer/+/crl/pem",
				"issuer/+/crl",
				"issuer/+/crl/delta/der",
				"issuer/+/crl/delta/pem",
				"issuer/+/crl/delta",
				"issuer/+/unified-crl/der",
				"issuer/+/unified-crl/pem",
				"issuer/+/unified-crl",
				"issuer/+/unified-crl/delta/der",
				"issuer/+/unified-crl/delta/pem",
				"issuer/+/unified-crl/delta",
				"issuer/+/pem",
				"issuer/+/der",
				"issuer/+/json",
				"issuers/", // LIST operations append a '/' to the requested path
				"ocsp",     // OCSP POST
				"ocsp/*",   // OCSP GET
				"unified-crl/delta",
				"unified-crl/delta/pem",
				"unified-crl/pem",
				"unified-crl",
				"unified-ocsp",   // Unified OCSP POST
				"unified-ocsp/*", // Unified OCSP GET

				// ACME paths are added below
			},

			LocalStorage: []string{
				revokedPath,
				localDeltaWALPath,
				legacyCRLPath,
				clusterConfigPath,
				issuing.PathCrls,
				issuing.PathCerts,
				issuing.PathCertMetadata,
				acmePathPrefix,
				autoTidyLastRunPath,
			},

			Root: []string{
				"root",
				"root/sign-self-issued",
			},

			SealWrapStorage: []string{
				legacyCertBundlePath,
				legacyCertBundleBackupPath,
				keyPrefix,
			},

			WriteForwardedStorage: []string{
				crossRevocationPath,
				revocation.UnifiedRevocationWritePathPrefix,
				unifiedDeltaWALPath,
			},

			Limited: []string{
				"issue",
				"issue/*",
			},

			Binary: []string{
				"ocsp",           // OCSP POST
				"ocsp/*",         // OCSP GET
				"unified-ocsp",   // Unified OCSP POST
				"unified-ocsp/*", // Unified OCSP GET
			},
		},

		Paths: []*framework.Path{
			pathListRoles(&b),
			pathRoles(&b),
			pathGenerateRoot(&b),
			pathSignIntermediate(&b),
			pathSignSelfIssued(&b),
			pathDeleteRoot(&b),
			pathGenerateIntermediate(&b),
			pathSetSignedIntermediate(&b),
			pathConfigCA(&b),
			pathConfigCRL(&b),
			pathConfigURLs(&b),
			pathConfigCluster(&b),
			pathSignVerbatim(&b),
			pathSign(&b),
			pathIssue(&b),
			pathRotateCRL(&b),
			pathRotateDeltaCRL(&b),
			pathRevoke(&b),
			pathRevokeWithKey(&b),
			pathListCertsRevoked(&b),
			pathTidy(&b),
			pathTidyCancel(&b),
			pathTidyStatus(&b),
			pathConfigAutoTidy(&b),

			// Issuer APIs
			pathListIssuers(&b),
			pathGetIssuer(&b),
			pathGetUnauthedIssuer(&b),
			pathGetIssuerCRL(&b),
			pathImportIssuer(&b),
			pathIssuerIssue(&b),
			pathIssuerSign(&b),
			pathIssuerSignIntermediate(&b),
			pathIssuerSignSelfIssued(&b),
			pathIssuerSignVerbatim(&b),
			pathIssuerGenerateRoot(&b),
			pathRotateRoot(&b),
			pathIssuerGenerateIntermediate(&b),
			pathCrossSignIntermediate(&b),
			pathConfigIssuers(&b),
			pathReplaceRoot(&b),
			pathRevokeIssuer(&b),

			// Key APIs
			pathListKeys(&b),
			pathKey(&b),
			pathGenerateKey(&b),
			pathImportKey(&b),
			pathConfigKeys(&b),

			// Fetch APIs have been lowered to favor the newer issuer API endpoints
			pathFetchCA(&b),
			pathFetchCAChain(&b),
			pathFetchCRL(&b),
			pathFetchCRLViaCertPath(&b),
			pathFetchValidRaw(&b),
			pathFetchValid(&b),
			pathFetchListCerts(&b),

			// OCSP APIs
			buildPathOcspGet(&b),
			buildPathOcspPost(&b),

			// CRL Signing
			pathResignCrls(&b),
			pathSignRevocationList(&b),

			// ACME
			pathAcmeConfig(&b),
			pathAcmeEabList(&b),
			pathAcmeEabDelete(&b),
			pathAcmeMgmtAccountList(&b),
			pathAcmeMgmtAccountRead(&b),
		},

		Secrets: []*framework.Secret{
			secretCerts(&b),
		},

		BackendType:    logical.TypeLogical,
		InitializeFunc: b.initialize,
		Invalidate:     b.invalidate,
		PeriodicFunc:   b.periodicFunc,
		Clean:          b.cleanup,
	}

	// Add ACME paths to backend
	for _, prefix := range []struct {
		acmePrefix   string
		unauthPrefix string
		opts         acmeWrapperOpts
	}{
		{
			"acme",
			"acme",
			acmeWrapperOpts{true, false},
		},
		{
			"roles/" + framework.GenericNameRegex("role") + "/acme",
			"roles/+/acme",
			acmeWrapperOpts{},
		},
		{
			"issuer/" + framework.GenericNameRegex(issuerRefParam) + "/acme",
			"issuer/+/acme",
			acmeWrapperOpts{},
		},
		{
			"issuer/" + framework.GenericNameRegex(issuerRefParam) + "/roles/" + framework.GenericNameRegex("role") + "/acme",
			"issuer/+/roles/+/acme",
			acmeWrapperOpts{},
		},
	} {
		setupAcmeDirectory(&b, prefix.acmePrefix, prefix.unauthPrefix, prefix.opts)
	}

	b.tidyCASGuard = new(uint32)
	b.tidyCancelCAS = new(uint32)
	b.tidyStatus = &tidyStatus{state: tidyStatusInactive}
	b.storage = conf.StorageView
	b.backendUUID = conf.BackendUUID

	b.pkiStorageVersion.Store(0)

	// b isn't yet initialized with SystemView state; calling b.System() will
	// result in a nil pointer dereference. Instead query BackendConfig's
	// copy of SystemView.
	cannotRebuildCRLs := conf.System.ReplicationState().HasState(consts.ReplicationPerformanceStandby) ||
		conf.System.ReplicationState().HasState(consts.ReplicationDRSecondary)
	b.crlBuilder = newCRLBuilder(!cannotRebuildCRLs)

	// Delay the first tidy until after we've started up, this will be reset within the initialize function
	now := time.Now()
	b.tidyStatusLock.Lock()
	b.lastAutoTidy = now
	b.tidyStatusLock.Unlock()

	// Keep track of when this mount was started up.
	b.mountStartup = now

	b.unifiedTransferStatus = newUnifiedTransferStatus()

	b.acmeState = NewACMEState()
	b.certificateCounter = NewCertificateCounter(b.backendUUID)

	// It is important that we call SetupEnt at the very end as
	// some ENT backends need access to the member vars initialized above.
	b.SetupEnt()
	return &b
}

type backend struct {
	*framework.Backend
	entBackend

	backendUUID       string
	storage           logical.Storage
	revokeStorageLock sync.RWMutex
	tidyCASGuard      *uint32
	tidyCancelCAS     *uint32

	tidyStatusLock sync.RWMutex
	tidyStatus     *tidyStatus
	// lastAutoTidy should be accessed through the tidyStatusLock,
	// use getAutoTidyLastRun and writeAutoTidyLastRun instead of direct access
	lastAutoTidy time.Time

	// autoTidyBackoff a random time in the future in which auto-tidy can't start
	// for after the system starts up to avoid a thundering herd of tidy operations
	// at startup.
	autoTidyBackoff time.Time

	unifiedTransferStatus *UnifiedTransferStatus

	certificateCounter *CertificateCounter

	pkiStorageVersion atomic.Value
	crlBuilder        *CrlBuilder

	// Write lock around issuers and keys.
	issuersLock sync.RWMutex

	// Context around ACME operations
	acmeState       *acmeState
	acmeAccountLock sync.RWMutex // (Write) Locked on Tidy, (Read) Locked on Account Creation

	// Track when this mount was started.
	mountStartup time.Time
}

// BackendOps a bridge/legacy interface until we can further
// separate out backend things into distinct packages.
type BackendOps interface {
	managed_key.PkiManagedKeyView
	pki_backend.SystemViewGetter
	pki_backend.MountInfo
	pki_backend.Logger
	revocation.RevokerFactory

	UseLegacyBundleCaStorage() bool
	CrlBuilder() *CrlBuilder
	GetRevokeStorageLock() *sync.RWMutex
	GetUnifiedTransferStatus() *UnifiedTransferStatus
	GetAcmeState() *acmeState
	GetRole(ctx context.Context, s logical.Storage, n string) (*issuing.RoleEntry, error)
	GetCertificateCounter() *CertificateCounter
}

var _ BackendOps = &backend{}

type roleOperation func(ctx context.Context, req *logical.Request, data *framework.FieldData, role *issuing.RoleEntry) (*logical.Response, error)

const backendHelp = `
The PKI backend dynamically generates X509 server and client certificates.

After mounting this backend, configure the CA using the "pem_bundle" endpoint within
the "config/" path.
`

func metricsKey(req *logical.Request, extra ...string) []string {
	if req == nil || req.MountPoint == "" {
		return extra
	}
	key := make([]string, len(extra)+1)
	key[0] = req.MountPoint[:len(req.MountPoint)-1]
	copy(key[1:], extra)
	return key
}

func (b *backend) metricsWrap(callType string, roleMode int, ofunc roleOperation) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		key := metricsKey(req, callType)
		var role *issuing.RoleEntry
		var labels []metrics.Label
		var err error

		var roleName string
		switch roleMode {
		case roleRequired:
			roleName = data.Get("role").(string)
		case roleOptional:
			r, ok := data.GetOk("role")
			if ok {
				roleName = r.(string)
			}
		}
		if roleMode > noRole {
			// Get the role
			role, err = b.GetRole(ctx, req.Storage, roleName)
			if err != nil {
				return nil, err
			}
			if role == nil && (roleMode == roleRequired || len(roleName) > 0) {
				return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
			}
			labels = []metrics.Label{{"role", roleName}}
		}

		ns, err := namespace.FromContext(ctx)
		if err == nil {
			labels = append(labels, metricsutil.NamespaceLabel(ns))
		}

		start := time.Now()
		defer metrics.MeasureSinceWithLabels(key, start, labels)
		resp, err := ofunc(ctx, req, data, role)

		if err != nil || resp.IsError() {
			metrics.IncrCounterWithLabels(append(key, "failure"), 1.0, labels)
		} else {
			metrics.IncrCounterWithLabels(key, 1.0, labels)
		}
		return resp, err
	}
}

// initialize is used to perform a possible PKI storage migration if needed
func (b *backend) initialize(ctx context.Context, ir *logical.InitializationRequest) error {
	sc := b.makeStorageContext(ctx, b.storage)
	if err := b.CrlBuilder().reloadConfigIfRequired(sc); err != nil {
		return err
	}

	err := b.initializePKIIssuersStorage(ctx)
	if err != nil {
		return err
	}

	err = b.GetAcmeState().Initialize(b, sc)
	if err != nil {
		return err
	}

	// Initialize also needs to populate our certificate and revoked certificate count
	err = b.initializeStoredCertificateCounts(ctx)
	if err != nil {
		// Don't block/err initialize/startup for metrics.  Context on this call can time out due to number of certificates.
		b.Logger().Error("Could not initialize stored certificate counts", "error", err)
		b.GetCertificateCounter().SetError(err)
	}

	// Initialize lastAutoTidy from disk
	b.initializeLastTidyFromStorage(sc)

	return b.initializeEnt(sc, ir)
}

// initializeLastTidyFromStorage reads the time we last ran auto tidy from storage and initializes
// b.lastAutoTidy with the value. If no previous value existed, we persist time.Now() and initialize
// b.lastAutoTidy with that value.
func (b *backend) initializeLastTidyFromStorage(sc *storageContext) {
	now := time.Now()

	lastTidyTime, err := sc.getAutoTidyLastRun()
	if err != nil {
		lastTidyTime = now
		b.Logger().Error("failed loading previous tidy last run time, using now", "error", err.Error())
	}
	if lastTidyTime.IsZero() {
		// No previous time was set, persist now so we can track a starting point across Vault restarts
		lastTidyTime = now
		if err = b.updateLastAutoTidyTime(sc, now); err != nil {
			b.Logger().Error("failed persisting tidy last run time", "error", err.Error())
		}
	}

	// We bypass using updateLastAutoTidyTime here to avoid the storage write on init
	// that normally isn't required
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()
	b.lastAutoTidy = lastTidyTime
}

func (b *backend) cleanup(ctx context.Context) {
	sc := b.makeStorageContext(ctx, b.storage)

	b.GetAcmeState().Shutdown(b)

	b.cleanupEnt(sc)
}

func (b *backend) initializePKIIssuersStorage(ctx context.Context) error {
	// Grab the lock prior to the updating of the storage lock preventing us flipping
	// the storage flag midway through the request stream of other requests.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	// Load up our current pki storage state, no matter the host type we are on.
	b.updatePkiStorageVersion(ctx, false)

	// Early exit if not a primary cluster or performance secondary with a local mount.
	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		b.Logger().Debug("skipping PKI migration as we are not on primary or secondary with a local mount")
		return nil
	}

	if err := migrateStorage(ctx, b, b.storage); err != nil {
		b.Logger().Error("Error during migration of PKI mount: " + err.Error())
		return err
	}

	b.updatePkiStorageVersion(ctx, false)

	return nil
}

func (b *backend) BackendUUID() string {
	return b.backendUUID
}

func (b *backend) CrlBuilder() *CrlBuilder {
	return b.crlBuilder
}

func (b *backend) GetRevokeStorageLock() *sync.RWMutex {
	return &b.revokeStorageLock
}

func (b *backend) GetUnifiedTransferStatus() *UnifiedTransferStatus {
	return b.unifiedTransferStatus
}

func (b *backend) GetAcmeState() *acmeState {
	return b.acmeState
}

func (b *backend) GetCertificateCounter() *CertificateCounter {
	return b.certificateCounter
}

func (b *backend) UseLegacyBundleCaStorage() bool {
	// This helper function is here to choose whether or not we use the newer
	// issuer/key storage format or the older legacy ca bundle format.
	//
	// This happens because we might've upgraded secondary PR clusters to
	// newer vault code versions. We still want to be able to service requests
	// with the old bundle format (e.g., issuing and revoking certs), until
	// the primary cluster's active node is upgraded to the newer Vault version
	// and the storage is migrated to the new format.
	version := b.pkiStorageVersion.Load()
	return version == nil || version == 0
}

func (b *backend) IsSecondaryNode() bool {
	return b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby)
}

func (b *backend) GetManagedKeyView() (logical.ManagedKeySystemView, error) {
	managedKeyView, ok := b.System().(logical.ManagedKeySystemView)
	if !ok {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unsupported system view")}
	}
	return managedKeyView, nil
}

func (b *backend) updatePkiStorageVersion(ctx context.Context, grabIssuersLock bool) {
	info, err := getMigrationInfo(ctx, b.storage)
	if err != nil {
		b.Logger().Error(fmt.Sprintf("Failed loading PKI migration status, staying in legacy mode: %v", err))
		return
	}

	// If this method is called outside the initialize function, like say an
	// invalidate func on a performance replica cluster, we should be grabbing
	// the issuers lock to offer a consistent view of the storage version while
	// other events are processing things. Its unknown what might happen during
	// a single event if one part thinks we are in legacy mode, and then later
	// on we aren't.
	if grabIssuersLock {
		b.issuersLock.Lock()
		defer b.issuersLock.Unlock()
	}

	if info.isRequired {
		b.pkiStorageVersion.Store(0)
	} else {
		b.pkiStorageVersion.Store(1)
	}
}

func (b *backend) invalidate(ctx context.Context, key string) {
	isNotPerfPrimary := b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary))

	switch {
	case strings.HasPrefix(key, legacyMigrationBundleLogKey):
		// This is for a secondary cluster to pick up that the migration has completed
		// and reset its compatibility mode and rebuild the CRL locally. Kick it off
		// as a go routine to not block this call due to the lock grabbing
		// within updatePkiStorageVersion.
		go func() {
			b.Logger().Info("Detected a migration completed, resetting pki storage version")
			b.updatePkiStorageVersion(ctx, true)
			b.CrlBuilder().requestRebuildIfActiveNode(b)
		}()
	case strings.HasPrefix(key, issuerPrefix):
		if !b.UseLegacyBundleCaStorage() {
			// See note in updateDefaultIssuerId about why this is necessary.
			// We do this ahead of CRL rebuilding just so we know that things
			// are stale.
			b.CrlBuilder().invalidateCRLBuildTime()

			// If an issuer has changed on the primary, we need to schedule an update of our CRL,
			// the primary cluster would have done it already, but the CRL is cluster specific so
			// force a rebuild of ours.
			b.CrlBuilder().requestRebuildIfActiveNode(b)
		} else {
			b.Logger().Debug("Ignoring invalidation updates for issuer as the PKI migration has yet to complete.")
		}
	case key == "config/crl":
		// We may need to reload our OCSP status flag
		b.CrlBuilder().markConfigDirty()
	case key == storageAcmeConfig:
		b.GetAcmeState().markConfigDirty()
	case key == storageIssuerConfig:
		b.CrlBuilder().invalidateCRLBuildTime()
	case strings.HasPrefix(key, crossRevocationPrefix):
		split := strings.Split(key, "/")

		if !strings.HasSuffix(key, "/confirmed") {
			cluster := split[len(split)-2]
			serial := split[len(split)-1]
			b.CrlBuilder().addCertForRevocationCheck(cluster, serial)
		} else {
			if len(split) >= 3 {
				cluster := split[len(split)-3]
				serial := split[len(split)-2]
				// Only process confirmations on the perf primary. The
				// performance secondaries cannot remove other clusters'
				// entries, and so do not need to track them (only to
				// ignore them). On performance primary nodes though,
				// we do want to track them to remove them.
				if !isNotPerfPrimary {
					b.CrlBuilder().addCertForRevocationRemoval(cluster, serial)
				}
			}
		}
	case strings.HasPrefix(key, unifiedRevocationReadPathPrefix):
		// Three parts to this key: prefix, cluster, and serial.
		split := strings.Split(key, "/")
		cluster := split[len(split)-2]
		serial := split[len(split)-1]
		b.CrlBuilder().addCertFromCrossRevocation(cluster, serial)
	}

	b.invalidateEnt(ctx, key)
}

func (b *backend) periodicFunc(ctx context.Context, request *logical.Request) error {
	sc := b.makeStorageContext(ctx, request.Storage)

	doCRL := func() error {
		// First attempt to reload the CRL configuration.
		if err := b.CrlBuilder().reloadConfigIfRequired(sc); err != nil {
			return err
		}

		// As we're (below) modifying the backing storage, we need to ensure
		// we're not on a standby/secondary node.
		if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) ||
			b.System().ReplicationState().HasState(consts.ReplicationDRSecondary) {
			return nil
		}

		// First handle any global revocation queue entries.
		if err := b.CrlBuilder().processRevocationQueue(sc); err != nil {
			return err
		}

		// Then handle any unified cross-cluster revocations.
		if err := b.CrlBuilder().processCrossClusterRevocations(sc); err != nil {
			return err
		}

		// Check if we're set to auto rebuild and a CRL is set to expire.
		if err := b.CrlBuilder().checkForAutoRebuild(sc); err != nil {
			return err
		}

		// Then attempt to rebuild the CRLs if required.
		warnings, err := b.CrlBuilder().RebuildIfForced(sc)
		if err != nil {
			return err
		}
		if len(warnings) > 0 {
			msg := "During rebuild of complete CRL, got the following warnings:"
			for index, warning := range warnings {
				msg = fmt.Sprintf("%v\n %d. %v", msg, index+1, warning)
			}
			b.Logger().Warn(msg)
		}

		// If a delta CRL was rebuilt above as part of the complete CRL rebuild,
		// this will be a no-op. However, if we do need to rebuild delta CRLs,
		// this would cause us to do so.
		warnings, err = b.CrlBuilder().rebuildDeltaCRLsIfForced(sc, false)
		if err != nil {
			return err
		}
		if len(warnings) > 0 {
			msg := "During rebuild of delta CRL, got the following warnings:"
			for index, warning := range warnings {
				msg = fmt.Sprintf("%v\n %d. %v", msg, index+1, warning)
			}
			b.Logger().Warn(msg)
		}

		return nil
	}

	doAutoTidy := func() error {
		// As we're (below) modifying the backing storage, we need to ensure
		// we're not on a standby/secondary node.
		if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) ||
			b.System().ReplicationState().HasState(consts.ReplicationDRSecondary) {
			return nil
		}

		config, err := sc.getAutoTidyConfig()
		if err != nil {
			return err
		}

		if !config.Enabled || config.Interval <= 0*time.Second {
			return nil
		}

		// Check if we should run another tidy...
		now := time.Now()
		nextOp := b.getLastAutoTidyTime().Add(config.Interval)
		if now.Before(nextOp) {
			return nil
		}

		if b.autoTidyBackoff.IsZero() {
			b.autoTidyBackoff = config.CalculateStartupBackoff(b.mountStartup)
		}

		if b.autoTidyBackoff.After(now) {
			b.Logger().Info("Auto tidy will not run as we are still within the random backoff ending at", "backoff_until", b.autoTidyBackoff)
			return nil
		}

		// Ensure a tidy isn't already running... If it is, we'll trigger
		// again when the running one finishes.
		if !atomic.CompareAndSwapUint32(b.tidyCASGuard, 0, 1) {
			return nil
		}

		// Prevent ourselves from starting another tidy operation while
		// this one is still running. This operation runs in the background
		// and has a separate error reporting mechanism.
		err = b.updateLastAutoTidyTime(sc, now)
		if err != nil {
			// We don't really mind if this write fails, we'll re-run in the future
			b.Logger().Warn("failed to persist auto tidy last run time", "error", err.Error())
		}

		// Because the request from the parent storage will be cleared at
		// some point (and potentially reused) -- due to tidy executing in
		// a background goroutine -- we need to copy the storage entry off
		// of the backend instead.
		backendReq := &logical.Request{
			Storage: b.storage,
		}

		b.startTidyOperation(backendReq, config)
		return nil
	}

	// First tidy any ACME nonces to free memory.
	b.GetAcmeState().DoTidyNonces()

	// Then run unified transfer.
	backgroundSc := b.makeStorageContext(context.Background(), b.storage)
	go runUnifiedTransfer(backgroundSc)

	// Then run the CRL rebuild and tidy operation.
	crlErr := doCRL()
	tidyErr := doAutoTidy()

	// Periodically re-emit gauges so that they don't disappear/go stale
	b.GetCertificateCounter().EmitCertStoreMetrics()

	var errors error
	if crlErr != nil {
		errors = multierror.Append(errors, fmt.Errorf("Error building CRLs:\n - %w\n", crlErr))
	}

	if tidyErr != nil {
		errors = multierror.Append(errors, fmt.Errorf("Error running auto-tidy:\n - %w\n", tidyErr))
	}

	if errors != nil {
		return errors
	}

	// Check if the CRL was invalidated due to issuer swap and update
	// accordingly.
	if err := b.CrlBuilder().flushCRLBuildTimeInvalidation(sc); err != nil {
		return err
	}

	// All good!
	return b.periodicFuncEnt(backgroundSc, request)
}

func (b *backend) initializeStoredCertificateCounts(ctx context.Context) error {
	// For performance reasons, we can't lock on issuance/storage of certs until a list operation completes,
	// but we want to limit possible miscounts / double-counts to over-counting, so we take the tidy lock which
	// prevents (most) deletions - in particular we take a read lock (sufficient to block the write lock in
	// tidyStatusStart while allowing tidy to still acquire a read lock to report via its endpoint)
	b.tidyStatusLock.RLock()
	defer b.tidyStatusLock.RUnlock()
	sc := b.makeStorageContext(ctx, b.storage)
	config, err := sc.getAutoTidyConfig()
	if err != nil {
		return err
	}

	certCounter := b.GetCertificateCounter()
	isEnabled := certCounter.ReconfigureWithTidyConfig(config)
	if !isEnabled {
		return nil
	}

	entries, err := b.storage.List(ctx, issuing.PathCerts)
	if err != nil {
		return err
	}

	revokedEntries, err := b.storage.List(ctx, "revoked/")
	if err != nil {
		return err
	}

	certCounter.InitializeCountsFromStorage(entries, revokedEntries)
	return nil
}

var _ revocation.Revoker = &revoker{}

type revoker struct {
	backend        *backend
	storageContext *storageContext
	crlConfig      *pki_backend.CrlConfig
}

func (r *revoker) RevokeCert(cert *x509.Certificate) (revocation.RevokeCertInfo, error) {
	r.backend.GetRevokeStorageLock().Lock()
	defer r.backend.GetRevokeStorageLock().Unlock()
	resp, err := revokeCert(r.storageContext, r.crlConfig, cert)
	return parseRevokeCertOutput(resp, err)
}

func (r *revoker) RevokeCertBySerial(serial string) (revocation.RevokeCertInfo, error) {
	// NOTE: tryRevokeCertBySerial grabs the revoke storage lock for us
	resp, err := tryRevokeCertBySerial(r.storageContext, r.crlConfig, serial)
	return parseRevokeCertOutput(resp, err)
}

// There are a bunch of reasons that a certificate will/won't be revoked. Sadly we will need a further
// refactoring but for now handle the basics of the reasons/response objects back to a usable object
// that doesn't directly reply to the API request
func parseRevokeCertOutput(resp *logical.Response, err error) (revocation.RevokeCertInfo, error) {
	if err != nil {
		return revocation.RevokeCertInfo{}, err
	}

	if resp == nil {
		// nil, nil response, most likely means the certificate was missing,
		// but *might* be other things such as a tainted mount
		return revocation.RevokeCertInfo{}, nil
	}

	if resp.IsError() {
		// There are a few reasons we return a response error but not an error,
		// such as UserError's or they tried to revoke the CA
		return revocation.RevokeCertInfo{}, resp.Error()
	}

	// It is possible we don't return the field for various reasons if just a bunch of warnings are set.
	if revTimeRaw, ok := resp.Data["revocation_time"]; ok {
		revTimeInt, err := parseutil.ParseInt(revTimeRaw)
		if err != nil {
			// Lets me lenient for now
			revTimeInt = 0
		}
		revTime := time.Unix(revTimeInt, 0)
		return revocation.RevokeCertInfo{
			RevocationTime: revTime,
		}, nil
	}

	// Since we don't really know what went wrong if anything, for example the certificate might
	// have been expired or close to expiry, lets punt on it for now
	return revocation.RevokeCertInfo{
		Warnings: resp.Warnings,
	}, nil
}

func (b *backend) GetRevoker(ctx context.Context, s logical.Storage) (revocation.Revoker, error) {
	sc := b.makeStorageContext(ctx, s)
	crlConfig, err := b.CrlBuilder().GetConfigWithUpdate(sc)
	if err != nil {
		return nil, err
	}
	return &revoker{
		backend:        b,
		crlConfig:      crlConfig,
		storageContext: sc,
	}, nil
}
