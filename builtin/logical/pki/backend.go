package pki

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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
		PeriodicFunc: b.periodicFunc,
		Help:         strings.TrimSpace(backendHelp),

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
			pathFetchValid(&b),
			pathFetchListCerts(&b),
			pathRevoke(&b),
			pathRevokeByRole(&b),
			pathTidy(&b),
		},

		Secrets: []*framework.Secret{
			secretCerts(&b),
		},

		BackendType: logical.TypeLogical,
	}

	b.crlLifetime = time.Hour * 72
	b.tidyCASGuard = new(uint32)
	b.storage = conf.StorageView

	return &b
}

func (b *backend) periodicFunc(ctx context.Context, req *logical.Request) error {
	fmt.Printf("period function has been called")
	// skip the function if autoTidy isn't being used
	zeroDuration, _ := time.ParseDuration("0s")
	if b.crlAutoTidy == zeroDuration {
		fmt.Printf("zero duration")
		return nil
	}
	// If autoTidy hasn't been run and the periodic function goes off, then simply
	// set the next time for tidying to the appropriate value
	if b.nextCrlAutoTidy.IsZero() {
		// NOTE:Commented out for demo purposes
		// fmt.Printf("first tidy push out")
		// b.nextCrlAutoTidy = time.Now().Add(b.crlAutoTidy)
		// return nil
		b.nextCrlAutoTidy = time.Now()
	}

	// We're ready to autoTidy
	if !time.Now().Before(b.nextCrlAutoTidy) && !b.nextCrlAutoTidy.IsZero() {
		fmt.Printf("ready to tidy")
		// tidy_revoked_certs is all we need, I think
		rawFieldData := map[string]interface{}{}
		rawFieldData["tidy_revoked_certs"] = true
		rawFieldData["tidy_cert_store"] = true
		rawFieldData["safety_buffer"] = 1 // this is for demo purposes. It is unsafe.
		tidyFd := &framework.FieldData{
			Raw:    rawFieldData,
			Schema: pathTidy(b).Fields,
		}
		_, err := b.pathTidyWrite(ctx, req, tidyFd)
		if err != nil {
			fmt.Printf("couldn't tidy")
			return err
		}
	}
	fmt.Printf("returning from tidy")
	return nil
}

type backend struct {
	*framework.Backend

	storage           logical.Storage
	crlLifetime       time.Duration
	revokeStorageLock sync.RWMutex
	tidyCASGuard      *uint32
	crlAutoTidy       time.Duration
	nextCrlAutoTidy   time.Time
}

const backendHelp = `
The PKI backend dynamically generates X509 server and client certificates.

After mounting this backend, configure the CA using the "pem_bundle" endpoint within
the "config/" path.
`
