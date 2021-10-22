package pki

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

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
				"connfig/ca_bundle",
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
	revokeStorageLock vault.DeadlockRWMutex
	tidyCASGuard      *uint32

	tidyStatusLock vault.DeadlockRWMutex
	tidyStatus     *tidyStatus
}

type tidyStatusState int

const (
	tidyStatusInactive tidyStatusState = iota
	tidyStatusStarted
	tidyStatusFinished
)

type tidyStatus struct {
	// Parameters used to initiate the operation
	safetyBuffer     int
	tidyCertStore    bool
	tidyRevokedCerts bool

	// Status
	state            tidyStatusState
	err              error
	timeStarted      time.Time
	timeFinished     time.Time
	message          string
	certStoreCount   uint
	revokedCertCount uint
}

const backendHelp = `
The PKI backend dynamically generates X509 server and client certificates.

After mounting this backend, configure the CA using the "pem_bundle" endpoint within
the "config/" path.
`
