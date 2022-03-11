package pki

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

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
 * The main ca and intermediate requests are always forwarded to the Primary cluster's active
 * node to write and send the key material/config globally across all clusters.
 *
 * CRL/Revocation and Issued certificate apis are handled by the active node within the cluster
 * they originate. Which means if a request comes into a performance secondary cluster the writes
 * will be forwarded to that cluster's active node and not go all the way up to the performance primary's
 * active node.
 *
 * If a certificate issue request has a role in which no_store is set to true that node itself
 * will issue the certificate and not forward the request to the active node.
 *
 * Following the same pattern if a managed key is involved to sign an issued certificate request
 * and the local node does not have access for some reason to it, the request will be forwarded to
 * the active node within the cluster only.
 *
 * To make sense of what goes where the following bits need to be analyzed within the codebase.
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
			},

			LocalStorage: []string{
				"revoked/",
				"crl",
				"certs/",
			},

			Root: []string{
				"root",
				"root/sign-self-issued",
			},

			SealWrapStorage: []string{
				"config/ca_bundle",
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
			pathFetchCA(&b),
			pathFetchCAChain(&b),
			pathFetchCRL(&b),
			pathFetchCRLViaCertPath(&b),
			pathFetchValidRaw(&b),
			pathFetchValid(&b),
			pathFetchListCerts(&b),
			pathRevoke(&b),
			pathTidy(&b),
			pathTidyStatus(&b),
		},

		Secrets: []*framework.Secret{
			secretCerts(&b),
		},

		BackendType: logical.TypeLogical,
	}

	b.crlLifetime = time.Hour * 72
	b.tidyCASGuard = new(uint32)
	b.tidyStatus = &tidyStatus{state: tidyStatusInactive}
	b.storage = conf.StorageView

	return &b
}

type backend struct {
	*framework.Backend

	storage           logical.Storage
	crlLifetime       time.Duration
	revokeStorageLock sync.RWMutex
	tidyCASGuard      *uint32

	tidyStatusLock sync.RWMutex
	tidyStatus     *tidyStatus
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
			if role == nil && roleMode == roleRequired {
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
