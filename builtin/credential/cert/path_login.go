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

	"github.com/ryanuber/go-glob"
)

// ParsedCert is a certificate that has been configured as trusted
type ParsedCert struct {
	Entry        *CertEntry
	Certificates []*x509.Certificate
}

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the certificate role to authenticate against.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},
	}
}

func (b *backend) pathLoginAliasLookahead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	clientCerts := req.Connection.ConnState.PeerCertificates
	if len(clientCerts) == 0 {
		return nil, fmt.Errorf("no client certificate found")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: clientCerts[0].Subject.CommonName,
			},
		},
	}, nil
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	var matched *ParsedCert
	if verifyResp, resp, err := b.verifyCredentials(req, data); err != nil {
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
				"subject_key_id":   certutil.GetHexFormatted(clientCerts[0].SubjectKeyId, ":"),
				"authority_key_id": certutil.GetHexFormatted(clientCerts[0].AuthorityKeyId, ":"),
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       ttl,
			},
			Alias: &logical.Alias{
				Name: clientCerts[0].SerialNumber.String(),
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
		if verifyResp, resp, err := b.verifyCredentials(req, d); err != nil {
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

func (b *backend) verifyCredentials(req *logical.Request, d *framework.FieldData) (*ParsedCert, *logical.Response, error) {
	// Get the connection state
	if req.Connection == nil || req.Connection.ConnState == nil {
		return nil, logical.ErrorResponse("tls connection required"), nil
	}
	connState := req.Connection.ConnState

	if connState.PeerCertificates == nil || len(connState.PeerCertificates) == 0 {
		return nil, logical.ErrorResponse("client certificate must be supplied"), nil
	}
	clientCert := connState.PeerCertificates[0]

	// Allow constraining the login request to a single CertEntry
	var certName string
	if req.Auth != nil { // It's a renewal, use the saved certName
		certName = req.Auth.Metadata["cert_name"]
	} else {
		certName = d.Get("name").(string)
	}

	// Load the trusted certificates
	roots, trusted, trustedNonCAs := b.loadTrustedCerts(req.Storage, certName)

	// Get the list of full chains matching the connection
	trustedChains, err := validateConnState(roots, connState)
	if err != nil {
		return nil, nil, err
	}

	// If trustedNonCAs is not empty it means that client had registered a non-CA cert
	// with the backend.
	if len(trustedNonCAs) != 0 {
		for _, trustedNonCA := range trustedNonCAs {
			tCert := trustedNonCA.Certificates[0]
			// Check for client cert being explicitly listed in the config (and matching other constraints)
			if tCert.SerialNumber.Cmp(clientCert.SerialNumber) == 0 &&
				bytes.Equal(tCert.AuthorityKeyId, clientCert.AuthorityKeyId) &&
				b.matchesConstraints(clientCert, trustedNonCA.Certificates, trustedNonCA) {
				return trustedNonCA, nil, nil
			}
		}
	}

	// If no trusted chain was found, client is not authenticated
	if len(trustedChains) == 0 {
		return nil, logical.ErrorResponse("invalid certificate or no client certificate supplied"), nil
	}

	// Search for a ParsedCert that intersects with the validated chains and any additional constraints
	matches := make([]*ParsedCert, 0)
	for _, trust := range trusted { // For each ParsedCert in the config
		for _, tCert := range trust.Certificates { // For each certificate in the entry
			for _, chain := range trustedChains { // For each root chain that we matched
				for _, cCert := range chain { // For each cert in the matched chain
					if tCert.Equal(cCert) && // ParsedCert intersects with matched chain
						b.matchesConstraints(clientCert, chain, trust) { // validate client cert + matched chain against the config
						// Add the match to the list
						matches = append(matches, trust)
					}
				}
			}
		}
	}

	// Fail on no matches
	if len(matches) == 0 {
		return nil, logical.ErrorResponse("no chain matching all constraints could be found for this login certificate"), nil
	}

	// Return the first matching entry (for backwards compatibility, we continue to just pick one if multiple match)
	return matches[0], nil, nil
}

func (b *backend) matchesConstraints(clientCert *x509.Certificate, trustedChain []*x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no names) is to allow all names
	nameMatched := len(config.Entry.AllowedNames) == 0
	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedName := range config.Entry.AllowedNames {
		if glob.Glob(allowedName, clientCert.Subject.CommonName) {
			nameMatched = true
		}

		for _, name := range clientCert.DNSNames {
			if glob.Glob(allowedName, name) {
				nameMatched = true
			}
		}

		for _, name := range clientCert.EmailAddresses {
			if glob.Glob(allowedName, name) {
				nameMatched = true
			}
		}
	}

	return !b.checkForChainInCRLs(trustedChain) && nameMatched
}

// loadTrustedCerts is used to load all the trusted certificates from the backend
func (b *backend) loadTrustedCerts(store logical.Storage, certName string) (pool *x509.CertPool, trusted []*ParsedCert, trustedNonCAs []*ParsedCert) {
	pool = x509.NewCertPool()
	trusted = make([]*ParsedCert, 0)
	trustedNonCAs = make([]*ParsedCert, 0)
	names, err := store.List("cert/")
	if err != nil {
		b.Logger().Error("cert: failed to list trusted certs", "error", err)
		return
	}
	for _, name := range names {
		// If we are trying to select a single CertEntry and this isn't it
		if certName != "" && name != certName {
			continue
		}
		entry, err := b.Cert(store, strings.TrimPrefix(name, "cert/"))
		if err != nil {
			b.Logger().Error("cert: failed to load trusted cert", "name", name, "error", err)
			continue
		}
		parsed := parsePEM([]byte(entry.Certificate))
		if len(parsed) == 0 {
			b.Logger().Error("cert: failed to parse certificate", "name", name)
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
