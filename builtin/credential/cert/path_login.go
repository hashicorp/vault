package cert

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"github.com/hashicorp/vault/helper/cidrutil"
	glob "github.com/ryanuber/go-glob"
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

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var matched *ParsedCert
	if verifyResp, resp, err := b.verifyCredentials(ctx, req, data); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		matched = verifyResp
	}

	if matched == nil {
		return nil, nil
	}

	if err := b.checkCIDR(matched.Entry, req); err != nil {
		return nil, err
	}

	clientCerts := req.Connection.ConnState.PeerCertificates
	if len(clientCerts) == 0 {
		return logical.ErrorResponse("no client certificate found"), nil
	}
	skid := base64.StdEncoding.EncodeToString(clientCerts[0].SubjectKeyId)
	akid := base64.StdEncoding.EncodeToString(clientCerts[0].AuthorityKeyId)

	resp := &logical.Response{
		Auth: &logical.Auth{
			Period: matched.Entry.Period,
			InternalData: map[string]interface{}{
				"subject_key_id":   skid,
				"authority_key_id": akid,
			},
			Policies:    matched.Entry.Policies,
			DisplayName: matched.Entry.DisplayName,
			Metadata: map[string]string{
				"cert_name":        matched.Entry.Name,
				"common_name":      clientCerts[0].Subject.CommonName,
				"serial_number":    clientCerts[0].SerialNumber.String(),
				"subject_key_id":   certutil.GetHexFormatted(clientCerts[0].SubjectKeyId, ":"),
				"authority_key_id": certutil.GetHexFormatted(clientCerts[0].AuthorityKeyId, ":"),
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       matched.Entry.TTL,
				MaxTTL:    matched.Entry.MaxTTL,
			},
			Alias: &logical.Alias{
				Name: clientCerts[0].Subject.CommonName,
			},
			BoundCIDRs: matched.Entry.BoundCIDRs,
		},
	}

	// Generate a response
	return resp, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if !config.DisableBinding {
		var matched *ParsedCert
		if verifyResp, resp, err := b.verifyCredentials(ctx, req, d); err != nil {
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
			return logical.ErrorResponse("no client certificate found"), nil
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
	cert, err := b.Cert(ctx, req.Storage, req.Auth.Metadata["cert_name"])
	if err != nil {
		return nil, err
	}
	if cert == nil {
		// User no longer exists, do not renew
		return nil, nil
	}

	if !policyutil.EquivalentPolicies(cert.Policies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = cert.TTL
	resp.Auth.MaxTTL = cert.MaxTTL
	resp.Auth.Period = cert.Period
	return resp, nil
}

func (b *backend) verifyCredentials(ctx context.Context, req *logical.Request, d *framework.FieldData) (*ParsedCert, *logical.Response, error) {
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
	roots, trusted, trustedNonCAs := b.loadTrustedCerts(ctx, req.Storage, certName)

	// Get the list of full chains matching the connection and validates the
	// certificate itself
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
	// This check happens after checking for a matching configured non-CA certs
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
	return !b.checkForChainInCRLs(trustedChain) &&
		b.matchesNames(clientCert, config) &&
		b.matchesCommonName(clientCert, config) &&
		b.matchesDNSSANs(clientCert, config) &&
		b.matchesEmailSANs(clientCert, config) &&
		b.matchesURISANs(clientCert, config) &&
		b.matchesOrganizationalUnits(clientCert, config) &&
		b.matchesCertificateExtensions(clientCert, config)
}

// matchesNames verifies that the certificate matches at least one configured
// allowed name
func (b *backend) matchesNames(clientCert *x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no names) is to allow all names
	if len(config.Entry.AllowedNames) == 0 {
		return true
	}
	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedName := range config.Entry.AllowedNames {
		if glob.Glob(allowedName, clientCert.Subject.CommonName) {
			return true
		}

		for _, name := range clientCert.DNSNames {
			if glob.Glob(allowedName, name) {
				return true
			}
		}

		for _, name := range clientCert.EmailAddresses {
			if glob.Glob(allowedName, name) {
				return true
			}
		}

	}
	return false
}

// matchesCommonName verifies that the certificate matches at least one configured
// allowed common name
func (b *backend) matchesCommonName(clientCert *x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no names) is to allow all names
	if len(config.Entry.AllowedCommonNames) == 0 {
		return true
	}
	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedCommonName := range config.Entry.AllowedCommonNames {
		if glob.Glob(allowedCommonName, clientCert.Subject.CommonName) {
			return true
		}
	}

	return false
}

// matchesDNSSANs verifies that the certificate matches at least one configured
// allowed dns entry in the subject alternate name extension
func (b *backend) matchesDNSSANs(clientCert *x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no names) is to allow all names
	if len(config.Entry.AllowedDNSSANs) == 0 {
		return true
	}
	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedDNS := range config.Entry.AllowedDNSSANs {
		for _, name := range clientCert.DNSNames {
			if glob.Glob(allowedDNS, name) {
				return true
			}
		}
	}

	return false
}

// matchesEmailSANs verifies that the certificate matches at least one configured
// allowed email in the subject alternate name extension
func (b *backend) matchesEmailSANs(clientCert *x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no names) is to allow all names
	if len(config.Entry.AllowedEmailSANs) == 0 {
		return true
	}
	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedEmail := range config.Entry.AllowedEmailSANs {
		for _, email := range clientCert.EmailAddresses {
			if glob.Glob(allowedEmail, email) {
				return true
			}
		}
	}

	return false
}

// matchesURISANs verifies that the certificate matches at least one configured
// allowed uri in the subject alternate name extension
func (b *backend) matchesURISANs(clientCert *x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no names) is to allow all names
	if len(config.Entry.AllowedURISANs) == 0 {
		return true
	}
	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedURI := range config.Entry.AllowedURISANs {
		for _, name := range clientCert.URIs {
			if glob.Glob(allowedURI, name.String()) {
				return true
			}
		}
	}

	return false
}

// matchesOrganizationalUnits verifies that the certificate matches at least one configurd allowed OU
func (b *backend) matchesOrganizationalUnits(clientCert *x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no OUs) is to allow all OUs
	if len(config.Entry.AllowedOrganizationalUnits) == 0 {
		return true
	}

	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedOrganizationalUnits := range config.Entry.AllowedOrganizationalUnits {
		for _, ou := range clientCert.Subject.OrganizationalUnit {
			if glob.Glob(allowedOrganizationalUnits, ou) {
				return true
			}
		}
	}

	return false
}

// matchesCertificateExtensions verifies that the certificate matches configured
// required extensions
func (b *backend) matchesCertificateExtensions(clientCert *x509.Certificate, config *ParsedCert) bool {
	// If no required extensions, nothing to check here
	if len(config.Entry.RequiredExtensions) == 0 {
		return true
	}
	// Fail fast if we have required extensions but no extensions on the cert
	if len(clientCert.Extensions) == 0 {
		return false
	}

	// Build Client Extensions Map for Constraint Matching
	// x509 Writes Extensions in ASN1 with a bitstring tag, which results in the field
	// including its ASN.1 type tag bytes. For the sake of simplicity, assume string type
	// and drop the tag bytes. And get the number of bytes from the tag.
	clientExtMap := make(map[string]string, len(clientCert.Extensions))
	for _, ext := range clientCert.Extensions {
		var parsedValue string
		asn1.Unmarshal(ext.Value, &parsedValue)
		clientExtMap[ext.Id.String()] = parsedValue
	}
	// If any of the required extensions don't match the constraint fails
	for _, requiredExt := range config.Entry.RequiredExtensions {
		reqExt := strings.SplitN(requiredExt, ":", 2)
		clientExtValue, clientExtValueOk := clientExtMap[reqExt[0]]
		if !clientExtValueOk || !glob.Glob(reqExt[1], clientExtValue) {
			return false
		}
	}
	return true
}

// loadTrustedCerts is used to load all the trusted certificates from the backend
func (b *backend) loadTrustedCerts(ctx context.Context, storage logical.Storage, certName string) (pool *x509.CertPool, trusted []*ParsedCert, trustedNonCAs []*ParsedCert) {
	pool = x509.NewCertPool()
	trusted = make([]*ParsedCert, 0)
	trustedNonCAs = make([]*ParsedCert, 0)
	names, err := storage.List(ctx, "cert/")
	if err != nil {
		b.Logger().Error("failed to list trusted certs", "error", err)
		return
	}
	for _, name := range names {
		// If we are trying to select a single CertEntry and this isn't it
		if certName != "" && name != certName {
			continue
		}
		entry, err := b.Cert(ctx, storage, strings.TrimPrefix(name, "cert/"))
		if err != nil {
			b.Logger().Error("failed to load trusted cert", "name", name, "error", err)
			continue
		}
		parsed := parsePEM([]byte(entry.Certificate))
		if len(parsed) == 0 {
			b.Logger().Error("failed to parse certificate", "name", name)
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

func (b *backend) checkCIDR(cert *CertEntry, req *logical.Request) error {
	if cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, cert.BoundCIDRs) {
		return nil
	}
	return logical.ErrPermissionDenied
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
	certs := cs.PeerCertificates
	if len(certs) == 0 {
		return nil, nil
	}

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
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
