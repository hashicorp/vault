package database

import (
	"context"
	"crypto/x509/pkix"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCA(b *databaseBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "ca",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.getCA(b.getOrGenerateCABundle),
			},
		},
		{
			Pattern: "ca/rotate",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.getCA(b.generateCABundle),
			},
		},
	}
}

func (b *databaseBackend) getCA(f func(context.Context, logical.Storage) (*certutil.CertBundle, *certutil.ParsedCertBundle, error)) func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		cb, parsedBundle, err := f(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		resp := &logical.Response{
			Data: map[string]interface{}{
				"expiration":    int64(parsedBundle.Certificate.NotAfter.Unix()),
				"serial_number": cb.SerialNumber,
				"certificate":   cb.Certificate,
			},
		}
		return resp, nil
	}
}

func (b *databaseBackend) getOrGenerateCABundle(ctx context.Context, s logical.Storage) (*certutil.CertBundle, *certutil.ParsedCertBundle, error) {
	cb, parsedBundle, err := b.getCABundle(ctx, s)
	if cb != nil || err != nil {
		return cb, parsedBundle, err
	}

	return b.generateCABundle(ctx, s)
}

func (b *databaseBackend) getCABundle(ctx context.Context, s logical.Storage) (*certutil.CertBundle, *certutil.ParsedCertBundle, error) {
	entry, err := s.Get(ctx, databaseCAPath+"root")
	if entry == nil || err != nil {
		return nil, nil, err
	}

	var cb certutil.CertBundle
	if err := entry.DecodeJSON(&cb); err != nil {
		return nil, nil, fmt.Errorf("unable to decode local CA certificate/key: %w", err)
	}
	parsedBundle, err := cb.ToParsedCertBundle()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse certificate bundle: %w", err)
	}

	if time.Now().Add(1 * time.Hour).After(parsedBundle.Certificate.NotAfter) {
		return nil, nil, nil
	}

	return &cb, parsedBundle, nil
}

func (b *databaseBackend) generateCABundle(ctx context.Context, s logical.Storage) (*certutil.CertBundle, *certutil.ParsedCertBundle, error) {
	parsedBundle, err := certutil.CreateCertificate(&certutil.CreationBundle{
		Params: &certutil.CreationParameters{
			NotAfter: time.Now().Add(24 * 365 * time.Hour),
			KeyType:  "rsa",
			KeyBits:  2048,
			URLs:     &certutil.URLEntries{},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, nil, fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}

	certs, err := s.List(ctx, databaseCAPath+"cert/")
	if err != nil {
		return nil, nil, fmt.Errorf("error while listing previously generated certs: %w", err)
	}
	for _, cert := range certs {
		if err := s.Delete(ctx, cert); err != nil {
			return nil, nil, fmt.Errorf("error while removing previous cert: %w", err)
		}
	}

	entry, err := logical.StorageEntryJSON(databaseCAPath+"root", cb)
	if err != nil {
		return nil, nil, err
	}
	err = s.Put(ctx, entry)
	if err != nil {
		return nil, nil, err
	}

	// Reset all the connections that depend on the CA
	for name := range b.connections {
		config, err := b.DatabaseConfig(ctx, s, name)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get the configuration for %s: %w", name, err)
		}
		if strings.Contains(config.ConnectionDetails["connection_url"].(string), "{{sslcert}}") ||
			strings.Contains(config.ConnectionDetails["connection_url"].(string), "{{sslkey}}") {
			b.resetConnection(ctx, s, name)
		}
	}

	return cb, parsedBundle, nil
}

func (b *databaseBackend) getOrGenerateCert(username string, ctx context.Context, s logical.Storage) (*certutil.CertBundle, error) {
	cb, err := b.getCert(username, ctx, s)
	if cb != nil || err != nil {
		return cb, err
	}

	return b.generateCert(username, ctx, s)
}

func (b *databaseBackend) getCert(username string, ctx context.Context, s logical.Storage) (*certutil.CertBundle, error) {
	entry, err := s.Get(ctx, databaseCAPath+"cert/"+username)
	if entry == nil || err != nil {
		return nil, err
	}

	var cb certutil.CertBundle
	if err := entry.DecodeJSON(&cb); err != nil {
		return nil, fmt.Errorf("unable to decode cert: %w", err)
	}

	parsedBundle, err := cb.ToParsedCertBundle()
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate bundle: %w", err)
	}
	if time.Now().Add(1 * time.Hour).After(parsedBundle.Certificate.NotAfter) {
		return nil, nil
	}

	return &cb, nil
}

func (b *databaseBackend) generateCert(username string, ctx context.Context, s logical.Storage) (*certutil.CertBundle, error) {
	_, caParsedBundle, err := b.getOrGenerateCABundle(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("failed to get CA bundle: %w", err)
	}

	parsedBundle, err := certutil.CreateCertificate(&certutil.CreationBundle{
		Params: &certutil.CreationParameters{
			Subject: pkix.Name{
				CommonName: username,
			},
			NotAfter: time.Now().Add(24 * 30 * time.Hour),
			KeyType:  "rsa",
			KeyBits:  2048,
			URLs:     &certutil.URLEntries{},
		},
		SigningBundle: &certutil.CAInfoBundle{
			ParsedCertBundle: *caParsedBundle,
			URLs:             &certutil.URLEntries{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate cert: %w", err)
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}
	entry, err := logical.StorageEntryJSON(databaseCAPath+"cert/"+username, cb)
	if err != nil {
		return nil, err
	}
	err = s.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	// Reset all the connections that depend on the cert
	for name := range b.connections {
		config, err := b.DatabaseConfig(ctx, s, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get the configuration for %s: %w", name, err)
		}
		if strings.Contains(config.ConnectionDetails["connection_url"].(string), "{{sslcert}}") ||
			strings.Contains(config.ConnectionDetails["connection_url"].(string), "{{sslkey}}") {
			if u, ok := config.ConnectionDetails["username"]; ok && u.(string) == username {
				b.resetConnection(ctx, s, name)
			}
		}
	}

	return cb, nil
}
