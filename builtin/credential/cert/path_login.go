package cert

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/policyutil"
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
			logical.UpdateOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	var matched *ParsedCert
	if verifyResp, resp, err := b.verifyCredentials(req); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		matched = verifyResp
	}

	if matched == nil {
		return nil, nil
	}

	ttl := matched.Entry.TTL
	if ttl == 0 {
		ttl = b.System().DefaultLeaseTTL()
	}

	clientCerts := req.Connection.ConnState.PeerCertificates
	if len(clientCerts) == 0 {
		return logical.ErrorResponse("no client certificate found"), nil
	}
	skid := base64.StdEncoding.EncodeToString(clientCerts[0].SubjectKeyId)
	akid := base64.StdEncoding.EncodeToString(clientCerts[0].AuthorityKeyId)

	// Generate a response
	resp := &logical.Response{
		Auth: &logical.Auth{
			InternalData: map[string]interface{}{
				"subject_key_id":   skid,
				"authority_key_id": akid,
			},
			Policies:    matched.Entry.Policies,
			DisplayName: matched.Entry.DisplayName,
			Metadata: map[string]string{
				"cert_name":        matched.Entry.Name,
				"common_name":      clientCerts[0].Subject.CommonName,
				"subject_key_id":   certutil.GetOctalFormatted(clientCerts[0].SubjectKeyId, ":"),
				"authority_key_id": certutil.GetOctalFormatted(clientCerts[0].AuthorityKeyId, ":"),
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       ttl,
			},
		},
	}
	return resp, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	if !config.DisableBinding {
		var matched *ParsedCert
		if verifyResp, resp, err := b.verifyCredentials(req); err != nil {
			return nil, err
		} else if resp != nil {
			return resp, nil
		} else {
			matched = verifyResp
		}

		if matched == nil {
			return nil, nil
		}

		clientCerts := req.Connection.ConnState.PeerCertificates
		if len(clientCerts) == 0 {
			return nil, fmt.Errorf("no client certificate found")
		}
		skid := base64.StdEncoding.EncodeToString(clientCerts[0].SubjectKeyId)
		akid := base64.StdEncoding.EncodeToString(clientCerts[0].AuthorityKeyId)

		// Certificate should not only match a registered certificate policy.
		// Also, the identity of the certificate presented should match the identity of the certificate used during login
		if req.Auth.InternalData["subject_key_id"] != skid && req.Auth.InternalData["authority_key_id"] != akid {
			return nil, fmt.Errorf("client identity during renewal not matching client identity used during login")
		}

	}
	// Get the cert and use its TTL
	cert, err := b.Cert(req.Storage, req.Auth.Metadata["cert_name"])
	if err != nil {
		return nil, err
	}
	if cert == nil {
		// User no longer exists, do not renew
		return nil, nil
	}

	if !policyutil.EquivalentPolicies(cert.Policies, req.Auth.Policies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	return framework.LeaseExtend(cert.TTL, 0, b.System())(req, d)
}

func (b *backend) verifyCredentials(req *logical.Request) (*ParsedCert, *logical.Response, error) {
	// Get the connection state
	if req.Connection == nil || req.Connection.ConnState == nil {
		return nil, logical.ErrorResponse("tls connection required"), nil
	}
	connState := req.Connection.ConnState

	if connState.PeerCertificates == nil || len(connState.PeerCertificates) == 0 {
		return nil, logical.ErrorResponse("client certificate must be supplied"), nil
	}

	// Load the trusted certificates
	roots, trusted, trustedNonCAs := b.loadTrustedCerts(req.Storage)

	// If trustedNonCAs is not empty it means that client had registered a non-CA cert
	// with the backend.
	if len(trustedNonCAs) != 0 {
		policy := b.matchNonCAPolicy(connState.PeerCertificates[0], trustedNonCAs)
		if policy != nil && !b.checkForChainInCRLs(policy.Certificates) {
			return policy, nil, nil
		}
	}

	// Validate the connection state is trusted
	trustedChains, err := validateConnState(roots, connState)
	if err != nil {
		return nil, nil, err
	}

	// If no trusted chain was found, client is not authenticated
	if len(trustedChains) == 0 {
		return nil, logical.ErrorResponse("invalid certificate or no client certificate supplied"), nil
	}

	validChain := b.checkForValidChain(trustedChains)
	if !validChain {
		return nil, logical.ErrorResponse(
			"no chain containing non-revoked certificates could be found for this login certificate",
		), nil
	}

	// Match the trusted chain with the policy
	return b.matchPolicy(trustedChains, trusted), nil, nil
}

// matchNonCAPolicy is used to match the client cert with the registered non-CA
// policies to establish client identity.
func (b *backend) matchNonCAPolicy(clientCert *x509.Certificate, trustedNonCAs []*ParsedCert) *ParsedCert {
	for _, trustedNonCA := range trustedNonCAs {
		tCert := trustedNonCA.Certificates[0]
		if tCert.SerialNumber.Cmp(clientCert.SerialNumber) == 0 && bytes.Equal(tCert.AuthorityKeyId, clientCert.AuthorityKeyId) {
			return trustedNonCA
		}
	}
	return nil
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
func (b *backend) loadTrustedCerts(store logical.Storage) (pool *x509.CertPool, trusted []*ParsedCert, trustedNonCAs []*ParsedCert) {
	pool = x509.NewCertPool()
	trusted = make([]*ParsedCert, 0)
	trustedNonCAs = make([]*ParsedCert, 0)
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
		if !parsed[0].IsCA {
			trustedNonCAs = append(trustedNonCAs, &ParsedCert{
				Entry:        entry,
				Certificates: parsed,
			})
		} else {
			for _, p := range parsed {
				pool.AddCert(p)
			}

			// Create a ParsedCert entry
			trusted = append(trusted, &ParsedCert{
				Entry:        entry,
				Certificates: parsed,
			})
		}
	}
	return
}

func (b *backend) checkForChainInCRLs(chain []*x509.Certificate) bool {
	badChain := false
	for _, cert := range chain {
		badCRLs := b.findSerialInCRLs(cert.SerialNumber)
		if len(badCRLs) != 0 {
			badChain = true
			break
		}
	}
	return badChain
}

func (b *backend) checkForValidChain(chains [][]*x509.Certificate) bool {
	for _, chain := range chains {
		if !b.checkForChainInCRLs(chain) {
			return true
		}
	}
	return false
}

// parsePEM parses a PEM encoded x509 certificate
func parsePEM(raw []byte) (certs []*x509.Certificate) {
	for len(raw) > 0 {
		var block *pem.Block
		block, raw = pem.Decode(raw)
		if block == nil {
			break
		}
		if (block.Type != "CERTIFICATE" && block.Type != "TRUSTED CERTIFICATE") || len(block.Headers) != 0 {
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
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
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
