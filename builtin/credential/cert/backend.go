package cert

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	if err := b.lockThenpopulateCRLs(ctx, conf.StorageView); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: []*framework.Path{
			pathConfig(&b),
			pathLogin(&b),
			pathListCerts(&b),
			pathCerts(&b),
			pathCRLs(&b),
		},
		AuthRenew:    b.loginPathWrapper(b.pathLoginRenew),
		Invalidate:   b.invalidate,
		BackendType:  logical.TypeCredential,
		PeriodicFunc: b.updateCRLs,
	}

	b.crlUpdateMutex = &sync.RWMutex{}

	return &b
}

type backend struct {
	*framework.Backend
	MapCertId *framework.PathMap

	crls           map[string]CRLInfo
	crlUpdateMutex *sync.RWMutex
}

func (b *backend) invalidate(_ context.Context, key string) {
	switch {
	case strings.HasPrefix(key, "crls/"):
		b.crlUpdateMutex.Lock()
		defer b.crlUpdateMutex.Unlock()
		b.crls = nil
	}
}

func (b *backend) fetchCRL(ctx context.Context, storage logical.Storage, name string, crl *CRLInfo) error {
	response, err := http.Get(crl.CDP.Url)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		certList, err := x509.ParseCRL(body)
		if err != nil {
			return err
		}
		crl.CDP.ValidUntil = certList.TBSCertList.NextUpdate
		return b.setCRL(ctx, storage, certList, name, crl.CDP)
	}
	return fmt.Errorf("unexpected response code %d fetching CRL from %s", response.StatusCode, crl.CDP.Url)
}

func (b *backend) updateCRLs(ctx context.Context, req *logical.Request) error {
	b.crlUpdateMutex.Lock()
	defer b.crlUpdateMutex.Unlock()
	var errs *multierror.Error
	for name, crl := range b.crls {
		if crl.CDP != nil && time.Now().After(crl.CDP.ValidUntil) {
			if err := b.fetchCRL(ctx, req.Storage, name, &crl); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}
	return errs.ErrorOrNil()
}

const backendHelp = `
The "cert" credential provider allows authentication using
TLS client certificates. A client connects to Vault and uses
the "login" endpoint to generate a client token.

Trusted certificates are configured using the "certs/" endpoint
by a user with root access. A certificate authority can be trusted,
which permits all keys signed by it. Alternatively, self-signed
certificates can be trusted avoiding the need for a CA.
`
