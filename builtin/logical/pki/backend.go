package pki

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	atomic2 "go.uber.org/atomic"

	"github.com/hashicorp/vault/sdk/helper/consts"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
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
 * If a certificate issue request has a role in which no_store is set to true, that node itself
 * will issue the certificate and not forward the request to the active node, as this does not
 * need to write to storage.
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
				"issuer/+/pem",
				"issuer/+/der",
				"issuer/+/json",
				"issuers/", // LIST operations append a '/' to the requested path
				"ocsp",     // OCSP POST
				"ocsp/*",   // OCSP GET
			},

			LocalStorage: []string{
				revokedPath,
				deltaWALPath,
				legacyCRLPath,
				clusterConfigPath,
				"crls/",
				"certs/",
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
		},

		Secrets: []*framework.Secret{
			secretCerts(&b),
		},

		BackendType:    logical.TypeLogical,
		InitializeFunc: b.initialize,
		Invalidate:     b.invalidate,
		PeriodicFunc:   b.periodicFunc,
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

	// Delay the first tidy until after we've started up.
	b.lastTidy = time.Now()

	// Metrics initialization for count of certificates in storage
	b.certsCounted = atomic2.NewBool(false)
	b.certCount = new(uint32)
	b.revokedCertCount = new(uint32)
	b.possibleDoubleCountedSerials = make([]string, 0, 250)
	b.possibleDoubleCountedRevokedSerials = make([]string, 0, 250)

	return &b
}

type backend struct {
	*framework.Backend

	backendUUID       string
	storage           logical.Storage
	revokeStorageLock sync.RWMutex
	tidyCASGuard      *uint32
	tidyCancelCAS     *uint32

	tidyStatusLock sync.RWMutex
	tidyStatus     *tidyStatus
	lastTidy       time.Time

	certCount                           *uint32
	revokedCertCount                    *uint32
	certsCounted                        *atomic2.Bool
	possibleDoubleCountedSerials        []string
	possibleDoubleCountedRevokedSerials []string

	pkiStorageVersion atomic.Value
	crlBuilder        *crlBuilder

	// Write lock around issuers and keys.
	issuersLock sync.RWMutex
}

type (
	tidyStatusState int
	roleOperation   func(ctx context.Context, req *logical.Request, data *framework.FieldData, role *roleEntry) (*logical.Response, error)
)

const (
	tidyStatusInactive   tidyStatusState = iota
	tidyStatusStarted                    = iota
	tidyStatusFinished                   = iota
	tidyStatusError                      = iota
	tidyStatusCancelling                 = iota
	tidyStatusCancelled                  = iota
)

type tidyStatus struct {
	// Parameters used to initiate the operation
	safetyBuffer       int
	issuerSafetyBuffer int
	tidyCertStore      bool
	tidyRevokedCerts   bool
	tidyRevokedAssocs  bool
	tidyExpiredIssuers bool
	tidyBackupBundle   bool
	pauseDuration      string

	// Status
	state                   tidyStatusState
	err                     error
	timeStarted             time.Time
	timeFinished            time.Time
	message                 string
	certStoreDeletedCount   uint
	revokedCertDeletedCount uint
	missingIssuerCertCount  uint
}

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
		var role *roleEntry
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
			role, err = b.getRole(ctx, req.Storage, roleName)
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
func (b *backend) initialize(ctx context.Context, _ *logical.InitializationRequest) error {
	sc := b.makeStorageContext(ctx, b.storage)
	if err := b.crlBuilder.reloadConfigIfRequired(sc); err != nil {
		return err
	}

	err := b.initializePKIIssuersStorage(ctx)
	if err != nil {
		return err
	}

	// Initialize also needs to populate our certificate and revoked certificate count
	err = b.initializeStoredCertificateCounts(ctx)
	if err != nil {
		return err
	}

	return nil
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

func (b *backend) useLegacyBundleCaStorage() bool {
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

func (b *backend) updatePkiStorageVersion(ctx context.Context, grabIssuersLock bool) {
	info, err := getMigrationInfo(ctx, b.storage)
	if err != nil {
		b.Logger().Error(fmt.Sprintf("Failed loading PKI migration status, staying in legacy mode: %v", err))
		return
	}

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
	switch {
	case strings.HasPrefix(key, legacyMigrationBundleLogKey):
		// This is for a secondary cluster to pick up that the migration has completed
		// and reset its compatibility mode and rebuild the CRL locally. Kick it off
		// as a go routine to not block this call due to the lock grabbing
		// within updatePkiStorageVersion.
		go func() {
			b.Logger().Info("Detected a migration completed, resetting pki storage version")
			b.updatePkiStorageVersion(ctx, true)
			b.crlBuilder.requestRebuildIfActiveNode(b)
		}()
	case strings.HasPrefix(key, issuerPrefix):
		if !b.useLegacyBundleCaStorage() {
			// See note in updateDefaultIssuerId about why this is necessary.
			// We do this ahead of CRL rebuilding just so we know that things
			// are stale.
			b.crlBuilder.invalidateCRLBuildTime()

			// If an issuer has changed on the primary, we need to schedule an update of our CRL,
			// the primary cluster would have done it already, but the CRL is cluster specific so
			// force a rebuild of ours.
			b.crlBuilder.requestRebuildIfActiveNode(b)
		} else {
			b.Logger().Debug("Ignoring invalidation updates for issuer as the PKI migration has yet to complete.")
		}
	case key == "config/crl":
		// We may need to reload our OCSP status flag
		b.crlBuilder.markConfigDirty()
	case key == storageIssuerConfig:
		b.crlBuilder.invalidateCRLBuildTime()
	}
}

func (b *backend) periodicFunc(ctx context.Context, request *logical.Request) error {
	sc := b.makeStorageContext(ctx, request.Storage)

	doCRL := func() error {
		// First attempt to reload the CRL configuration.
		if err := b.crlBuilder.reloadConfigIfRequired(sc); err != nil {
			return err
		}

		// As we're (below) modifying the backing storage, we need to ensure
		// we're not on a standby/secondary node.
		if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) ||
			b.System().ReplicationState().HasState(consts.ReplicationDRSecondary) {
			return nil
		}

		// Check if we're set to auto rebuild and a CRL is set to expire.
		if err := b.crlBuilder.checkForAutoRebuild(sc); err != nil {
			return err
		}

		// Then attempt to rebuild the CRLs if required.
		if err := b.crlBuilder.rebuildIfForced(sc); err != nil {
			return err
		}

		// If a delta CRL was rebuilt above as part of the complete CRL rebuild,
		// this will be a no-op. However, if we do need to rebuild delta CRLs,
		// this would cause us to do so.
		if err := b.crlBuilder.rebuildDeltaCRLsIfForced(sc, false); err != nil {
			return err
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
		b.tidyStatusLock.RLock()
		nextOp := b.lastTidy.Add(config.Interval)
		b.tidyStatusLock.RUnlock()
		if now.Before(nextOp) {
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
		b.tidyStatusLock.Lock()
		b.lastTidy = now
		b.tidyStatusLock.Unlock()

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

	crlErr := doCRL()
	tidyErr := doAutoTidy()

	if crlErr != nil && tidyErr != nil {
		return fmt.Errorf("Error building CRLs:\n - %v\n\nError running auto-tidy:\n - %w\n", crlErr, tidyErr)
	}

	if crlErr != nil {
		return fmt.Errorf("Error building CRLs:\n - %w\n", crlErr)
	}

	if tidyErr != nil {
		return fmt.Errorf("Error running auto-tidy:\n - %w\n", tidyErr)
	}

	// Check if the CRL was invalidated due to issuer swap and update
	// accordingly.
	if err := b.crlBuilder.flushCRLBuildTimeInvalidation(sc); err != nil {
		return err
	}

	// All good!
	return nil
}

func (b *backend) initializeStoredCertificateCounts(ctx context.Context) error {
	b.tidyStatusLock.RLock()
	defer b.tidyStatusLock.RUnlock()
	// For performance reasons, we can't lock on issuance/storage of certs until a list operation completes,
	// but we want to limit possible miscounts / double-counts to over-counting, so we take the tidy lock which
	// prevents (most) deletions - in particular we take a read lock (sufficient to block the write lock in
	// tidyStatusStart while allowing tidy to still acquire a read lock to report via its endpoint)

	entries, err := b.storage.List(ctx, "certs/")
	if err != nil {
		return err
	}
	atomic.AddUint32(b.certCount, uint32(len(entries)))

	revokedEntries, err := b.storage.List(ctx, "revoked/")
	if err != nil {
		return err
	}
	atomic.AddUint32(b.revokedCertCount, uint32(len(revokedEntries)))

	b.certsCounted.Store(true)
	// Now that the metrics are set, we can switch from appending newly-stored certificates to the possible double-count
	// list, and instead have them update the counter directly.  We need to do this so that we are looking at a static
	// slice of possibly double counted serials.  Note that certsCounted is computed before the storage operation, so
	// there may be some delay here.

	// Sort the listed-entries first, to accommodate that delay.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i] < entries[j]
	})

	sort.Slice(revokedEntries, func(i, j int) bool {
		return revokedEntries[i] < revokedEntries[j]
	})

	// We assume here that these lists are now complete.
	sort.Slice(b.possibleDoubleCountedSerials, func(i, j int) bool {
		return b.possibleDoubleCountedSerials[i] < b.possibleDoubleCountedSerials[j]
	})

	listEntriesIndex := 0
	possibleDoubleCountIndex := 0
	for {
		if listEntriesIndex >= len(entries) {
			break
		}
		if possibleDoubleCountIndex >= len(b.possibleDoubleCountedSerials) {
			break
		}
		if entries[listEntriesIndex] == b.possibleDoubleCountedSerials[possibleDoubleCountIndex] {
			// This represents a double-counted entry
			b.decrementTotalCertificatesCountNoReport()
			listEntriesIndex = listEntriesIndex + 1
			possibleDoubleCountIndex = possibleDoubleCountIndex + 1
			continue
		}
		if entries[listEntriesIndex] < b.possibleDoubleCountedSerials[possibleDoubleCountIndex] {
			listEntriesIndex = listEntriesIndex + 1
			continue
		}
		if entries[listEntriesIndex] > b.possibleDoubleCountedSerials[possibleDoubleCountIndex] {
			possibleDoubleCountIndex = possibleDoubleCountIndex + 1
			continue
		}
	}

	sort.Slice(b.possibleDoubleCountedRevokedSerials, func(i, j int) bool {
		return b.possibleDoubleCountedRevokedSerials[i] < b.possibleDoubleCountedRevokedSerials[j]
	})

	listRevokedEntriesIndex := 0
	possibleRevokedDoubleCountIndex := 0
	for {
		if listRevokedEntriesIndex >= len(revokedEntries) {
			break
		}
		if possibleRevokedDoubleCountIndex >= len(b.possibleDoubleCountedRevokedSerials) {
			break
		}
		if revokedEntries[listRevokedEntriesIndex] == b.possibleDoubleCountedRevokedSerials[possibleRevokedDoubleCountIndex] {
			// This represents a double-counted revoked entry
			b.decrementTotalRevokedCertificatesCountNoReport()
			listRevokedEntriesIndex = listRevokedEntriesIndex + 1
			possibleRevokedDoubleCountIndex = possibleRevokedDoubleCountIndex + 1
			continue
		}
		if revokedEntries[listRevokedEntriesIndex] < b.possibleDoubleCountedRevokedSerials[possibleRevokedDoubleCountIndex] {
			listRevokedEntriesIndex = listRevokedEntriesIndex + 1
			continue
		}
		if revokedEntries[listRevokedEntriesIndex] > b.possibleDoubleCountedRevokedSerials[possibleRevokedDoubleCountIndex] {
			possibleRevokedDoubleCountIndex = possibleRevokedDoubleCountIndex + 1
			continue
		}
	}

	b.possibleDoubleCountedRevokedSerials = nil
	b.possibleDoubleCountedSerials = nil

	certCount := atomic.LoadUint32(b.certCount)
	metrics.SetGauge([]string{"secrets", "pki", b.backendUUID, "total_certificates_stored"}, float32(certCount))
	revokedCertCount := atomic.LoadUint32(b.revokedCertCount)
	metrics.SetGauge([]string{"secrets", "pki", b.backendUUID, "total_revoked_certificates_stored"}, float32(revokedCertCount))

	return nil
}

// The "certsCounted" boolean here should be loaded from the backend certsCounted before the corresponding storage call:
// eg. certsCounted := b.certsCounted.Load()
func (b *backend) incrementTotalCertificatesCount(certsCounted bool, newSerial string) {
	certCount := atomic.AddUint32(b.certCount, 1)
	switch {
	case !certsCounted:
		// This is unsafe, but a good best-attempt
		if strings.HasPrefix(newSerial, "certs/") {
			newSerial = newSerial[6:]
		}
		b.possibleDoubleCountedSerials = append(b.possibleDoubleCountedSerials, newSerial)
	default:
		metrics.SetGauge([]string{"secrets", "pki", b.backendUUID, "total_certificates_stored"}, float32(certCount))
	}
}

func (b *backend) decrementTotalCertificatesCountReport() {
	certCount := b.decrementTotalCertificatesCountNoReport()
	metrics.SetGauge([]string{"secrets", "pki", b.backendUUID, "total_certificates_stored"}, float32(certCount))
}

// Called directly only by the initialize function to deduplicate the count, when we don't have a full count yet
func (b *backend) decrementTotalCertificatesCountNoReport() uint32 {
	newCount := atomic.AddUint32(b.certCount, ^uint32(0))
	return newCount
}

// The "certsCounted" boolean here should be loaded from the backend certsCounted before the corresponding storage call:
// eg. certsCounted := b.certsCounted.Load()
func (b *backend) incrementTotalRevokedCertificatesCount(certsCounted bool, newSerial string) {
	newRevokedCertCount := atomic.AddUint32(b.revokedCertCount, 1)
	switch {
	case !certsCounted:
		// This is unsafe, but a good best-attempt
		if strings.HasPrefix(newSerial, "revoked/") { // allow passing in the path (revoked/serial) OR the serial
			newSerial = newSerial[8:]
		}
		b.possibleDoubleCountedRevokedSerials = append(b.possibleDoubleCountedRevokedSerials, newSerial)
	default:
		metrics.SetGauge([]string{"secrets", "pki", b.backendUUID, "total_revoked_certificates_stored"}, float32(newRevokedCertCount))
	}
}

func (b *backend) decrementTotalRevokedCertificatesCountReport() {
	revokedCertCount := b.decrementTotalRevokedCertificatesCountNoReport()
	metrics.SetGauge([]string{"secrets", "pki", b.backendUUID, "total_revoked_certificates_stored"}, float32(revokedCertCount))
}

// Called directly only by the initialize function to deduplicate the count, when we don't have a full count yet
func (b *backend) decrementTotalRevokedCertificatesCountNoReport() uint32 {
	newRevokedCertCount := atomic.AddUint32(b.revokedCertCount, ^uint32(0))
	return newRevokedCertCount
}
