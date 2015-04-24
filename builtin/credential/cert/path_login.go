package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// TrustedCertificate is a certificate that has been configured as trusted
type TrustedCertificate struct {
	Certificates []*x509.Certificate
	Policies     []string
	DisplayName  string
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
		return nil, nil
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
		return nil, nil
	}

	// Match the trusted chain with the policy
	matched := b.matchPolicy(trustedChains, trusted)
	if matched == nil {
		return nil, nil
	}

	// Generate a response
	resp := &logical.Response{
		Auth: &logical.Auth{
			Policies:    matched.Policies,
			DisplayName: matched.DisplayName,
		},
	}
	return resp, nil
}

// matchPolicy is used to match the associated policy with the certificate that
// was used to establish the client identity.
func (b *backend) matchPolicy(chains [][]*x509.Certificate, trusted []*TrustedCertificate) *TrustedCertificate {
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
func (b *backend) loadTrustedCerts(store logical.Storage) (pool *x509.CertPool, trusted []*TrustedCertificate) {
	pool = x509.NewCertPool()
	names, err := b.MapCertId.List(store, "")
	if err != nil {
		b.Logger().Printf("[ERR] cert: failed to list trusted certs: %v", err)
		return
	}
	for _, name := range names {
		data, err := b.MapCertId.Get(store, name)
		if err != nil {
			b.Logger().Printf("[ERR] cert: failed to load trusted certs '%s': %v", name, err)
			continue
		}
		certRaw, ok := data["certificate"]
		if !ok {
			b.Logger().Printf("[ERR] cert: no certificate for '%s'", name)
			continue
		}
		cert, ok := certRaw.(string)
		if !ok {
			b.Logger().Printf("[ERR] cert: certificate for '%s' is not a string", name)
			continue
		}
		parsed := parsePEM([]byte(cert))
		if len(parsed) == 0 {
			b.Logger().Printf("[ERR] cert: failed to parse certificate for '%s'", name)
			continue
		}
		for _, p := range parsed {
			pool.AddCert(p)
		}

		// Extract the relevant policy
		var policyString string
		raw, ok := data["value"]
		if ok {
			rawS, ok := raw.(string)
			if ok {
				policyString = rawS
			}
		}

		// Extract the display name if any
		var displayName string
		raw, ok = data["display_name"]
		if ok {
			rawS, ok := raw.(string)
			if ok {
				displayName = rawS
			}
		}

		// Create a TrustedCertificate entry
		trusted = append(trusted, &TrustedCertificate{
			Certificates: parsed,
			Policies:     policyStringToList(policyString),
			DisplayName:  displayName,
		})
	}
	return
}

// policyStringToList turns a string with comma seperated
// policies into a sorted, de-duplicated list of policies.
func policyStringToList(s string) []string {
	set := make(map[string]struct{})
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			set[p] = struct{}{}
		}
	}

	list := make([]string, 0, len(set))
	for k, _ := range set {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
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

	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
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
