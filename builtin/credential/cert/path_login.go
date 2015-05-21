package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// ParsedCert is a certificate that has been configured as trusted
type ParsedCert struct {
	Entry        *CertEntry
	Certificates []*x509.Certificate
}

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields:  map[string]*framework.FieldSchema{},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the connection state
	if req.Connection == nil || req.Connection.ConnState == nil {
		return logical.ErrorResponse("tls connection required"), nil
	}
	connState := req.Connection.ConnState

	// Load the trusted certificates
	roots, trusted := b.loadTrustedCerts(req.Storage)

	// Validate the connection state is trusted
	trustedChains, err := validateConnState(roots, connState)
	if err != nil {
		return nil, err
	}

	// If no trusted chain was found, client is not authenticated
	if len(trustedChains) == 0 {
		return logical.ErrorResponse("invalid certificate"), nil
	}

	// Match the trusted chain with the policy
	matched := b.matchPolicy(trustedChains, trusted)
	if matched == nil {
		return nil, nil
	}

	// Generate a response
	resp := &logical.Response{
		Auth: &logical.Auth{
			Policies:    matched.Entry.Policies,
			DisplayName: matched.Entry.DisplayName,
			Metadata: map[string]string{
				"cert_name": matched.Entry.Name,
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				Lease:     matched.Entry.Lease,
			},
		},
	}
	return resp, nil
}

// matchPolicy is used to match the associated policy with the certificate that
// was used to establish the client identity.
func (b *backend) matchPolicy(chains [][]*x509.Certificate, trusted []*ParsedCert) *ParsedCert {
	// There is probably a better way to do this...
	for _, chain := range chains {
		for _, trust := range trusted {
			for _, tCert := range trust.Certificates {
				for _, cCert := range chain {
					if tCert.Equal(cCert) {
						return trust
					}
				}
			}
		}
	}
	return nil
}

// loadTrustedCerts is used to load all the trusted certificates from the backend
func (b *backend) loadTrustedCerts(store logical.Storage) (pool *x509.CertPool, trusted []*ParsedCert) {
	pool = x509.NewCertPool()
	names, err := store.List("cert/")
	if err != nil {
		b.Logger().Printf("[ERR] cert: failed to list trusted certs: %v", err)
		return
	}
	for _, name := range names {
		entry, err := b.Cert(store, strings.TrimPrefix(name, "cert/"))
		if err != nil {
			b.Logger().Printf("[ERR] cert: failed to load trusted certs '%s': %v", name, err)
			continue
		}
		parsed := parsePEM([]byte(entry.Certificate))
		if len(parsed) == 0 {
			b.Logger().Printf("[ERR] cert: failed to parse certificate for '%s'", name)
			continue
		}
		for _, p := range parsed {
			pool.AddCert(p)
		}

		// Create a ParsedCert entry
		trusted = append(trusted, &ParsedCert{
			Entry:        entry,
			Certificates: parsed,
		})
	}
	return
}

// parsePEM parses a PEM encoded x509 certificate
func parsePEM(raw []byte) (certs []*x509.Certificate) {
	for len(raw) > 0 {
		var block *pem.Block
		block, raw = pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		certs = append(certs, cert)
	}
	return
}

// validateConnState is used to validate that the TLS client is authorized
// by at trusted certificate. Most of this logic is lifted from the client
// verification logic here:  http://golang.org/src/crypto/tls/handshake_server.go
// The trusted chains are returned.
func validateConnState(roots *x509.CertPool, cs *tls.ConnectionState) ([][]*x509.Certificate, error) {
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	certs := cs.PeerCertificates
	if len(certs) == 0 {
		return nil, nil
	}

	if len(certs) > 1 {
		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}
	}

	chains, err := certs[0].Verify(opts)
	if err != nil {
		if _, ok := err.(x509.UnknownAuthorityError); ok {
			return nil, nil
		}
		return nil, errors.New("failed to verify client's certificate: " + err.Error())
	}
	return chains, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the cert and validate auth
	cert, err := b.Cert(req.Storage, req.Auth.Metadata["cert_name"])
	if err != nil {
		return nil, err
	}
	if cert == nil {
		// User no longer exists, do not renew
		return nil, nil
	}

	return framework.LeaseExtend(cert.Lease, 0)(req, d)
}
