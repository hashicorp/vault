package pki

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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
				"crl/pem",
				"crl",
				"issuer/+/crl/der",
				"issuer/+/crl/pem",
				"issuer/+/crl",
				"issuer/+/pem",
				"issuer/+/der",
				"issuer/+/json",
				"issuers",
			},

			LocalStorage: []string{
				"revoked/",
				legacyCRLPath,
				"crls/",
				"certs/",
			},

			Root: []string{
				"root",
				"root/sign-self-issued",
			},

			SealWrapStorage: []string{
				legacyCertBundlePath,
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
			pathSignVerbatim(&b),
			pathSign(&b),
			pathIssue(&b),
			pathRotateCRL(&b),
			pathRevoke(&b),
			pathTidy(&b),
			pathTidyStatus(&b),

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
		},

		Secrets: []*framework.Secret{
			secretCerts(&b),
		},

		BackendType:    logical.TypeLogical,
		InitializeFunc: b.initialize,
		Invalidate:     b.invalidate,
		PeriodicFunc:   b.periodicFunc,
	}

	b.crlLifetime = time.Hour * 72
	b.tidyCASGuard = new(uint32)
	b.tidyStatus = &tidyStatus{state: tidyStatusInactive}
	b.storage = conf.StorageView
	b.backendUUID = conf.BackendUUID

	b.pkiStorageVersion.Store(0)

	b.crlBuilder = &crlBuilder{}
	return &b
}

type backend struct {
	*framework.Backend

	backendUUID       string
	storage           logical.Storage
	crlLifetime       time.Duration
	revokeStorageLock sync.RWMutex
	tidyCASGuard      *uint32

	tidyStatusLock sync.RWMutex
	tidyStatus     *tidyStatus

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
	tidyStatusInactive tidyStatusState = iota
	tidyStatusStarted
	tidyStatusFinished
	tidyStatusError
)

type tidyStatus struct {
	// Parameters used to initiate the operation
	safetyBuffer     int
	tidyCertStore    bool
	tidyRevokedCerts bool

	// Status
	state                   tidyStatusState
	err                     error
	timeStarted             time.Time
	timeFinished            time.Time
	message                 string
	certStoreDeletedCount   uint
	revokedCertDeletedCount uint
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
		// If an issuer has changed on the primary, we need to schedule an update of our CRL,
		// the primary cluster would have done it already, but the CRL is cluster specific so
		// force a rebuild of ours.
		if !b.useLegacyBundleCaStorage() {
			b.crlBuilder.requestRebuildIfActiveNode(b)
		} else {
			b.Logger().Debug("Ignoring invalidation updates for issuer as the PKI migration has yet to complete.")
		}
	}
}

func (b *backend) periodicFunc(ctx context.Context, request *logical.Request) error {
	return b.crlBuilder.rebuildIfForced(ctx, b, request)
}
