// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/helper/ocsp"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
)

// ParsedCert is a certificate that has been configured as trusted
type ParsedCert struct {
	Entry        *CertEntry
	Certificates []*x509.Certificate
}

const certAuthFailMsg = "failed to match all constraints for this login certificate"

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixCert,
			OperationVerb:   "login",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The name of the certificate role to authenticate against.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.loginPathWrapper(b.pathLogin),
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
			logical.ResolveRoleOperation:    b.loginPathWrapper(b.pathLoginResolveRole),
		},
	}
}

func (b *backend) loginPathWrapper(wrappedOp func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error)) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		// Make sure that the CRLs have been loaded before processing a login request,
		// they might have been nil'd by an invalidate func call.
		if err := b.populateCrlsIfNil(ctx, req.Storage); err != nil {
			return nil, err
		}
		return wrappedOp(ctx, req, data)
	}
}

func (b *backend) pathLoginResolveRole(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Quota role rule creates send a probe to test if the backend returns
	// ErrUnsupportedOperation for ResolveRole, and there's no req.Storage populated
	// for these.  So just return a non-ErrUnsupportedOperation error.
	if req.Storage == nil {
		return logical.ErrorResponse("no storage"), logical.ErrMissingRequiredState
	}

	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if b.configUpdated.Load() {
		b.updatedConfig(config)
	}

	var matched *ParsedCert

	if verifyResp, resp, err := b.verifyCredentials(ctx, req, data); err != nil {
		return nil, err
	} else if resp != nil {
		return certAuthLoginFailureResponse(config, resp, req), nil
	} else {
		matched = verifyResp
	}

	if matched == nil {
		return logical.ErrorResponse("no certificate was matched by this request"), nil
	}

	return logical.ResolveRoleResponse(matched.Entry.Name)
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Connection == nil || req.Connection.ConnState == nil {
		return nil, fmt.Errorf("tls connection not found")
	}
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
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if b.configUpdated.Load() {
		b.updatedConfig(config)
	}

	var matched *ParsedCert
	if verifyResp, resp, err := b.verifyCredentials(ctx, req, data); err != nil {
		return nil, err
	} else if resp != nil {
		return certAuthLoginFailureResponse(config, resp, req), nil
	} else {
		matched = verifyResp
	}

	if matched == nil {
		return nil, nil
	}

	if len(matched.Entry.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, matched.Entry.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	clientCerts := req.Connection.ConnState.PeerCertificates
	if len(clientCerts) == 0 {
		return logical.ErrorResponse("no client certificate found"), nil
	}
	skid := base64.StdEncoding.EncodeToString(clientCerts[0].SubjectKeyId)
	akid := base64.StdEncoding.EncodeToString(clientCerts[0].AuthorityKeyId)

	metadata := map[string]string{
		"cert_name":        matched.Entry.Name,
		"common_name":      clientCerts[0].Subject.CommonName,
		"serial_number":    clientCerts[0].SerialNumber.String(),
		"subject_key_id":   certutil.GetHexFormatted(clientCerts[0].SubjectKeyId, ":"),
		"authority_key_id": certutil.GetHexFormatted(clientCerts[0].AuthorityKeyId, ":"),
	}

	// Add metadata from allowed_metadata_extensions when present,
	// with sanitized oids (dash-separated instead of dot-separated) as keys.
	for k, v := range b.certificateExtensionsMetadata(clientCerts[0], matched) {
		metadata[k] = v
	}

	auth := &logical.Auth{
		InternalData: map[string]interface{}{
			"subject_key_id":   skid,
			"authority_key_id": akid,
		},
		DisplayName: matched.Entry.DisplayName,
		Metadata:    metadata,
		Alias: &logical.Alias{
			Name: clientCerts[0].Subject.CommonName,
		},
	}

	if config.EnableIdentityAliasMetadata {
		auth.Alias.Metadata = metadata
	}

	matched.Entry.PopulateTokenAuth(auth)

	return &logical.Response{
		Auth: auth,
	}, nil
}

func certAuthLoginFailureResponse(config *config, resp *logical.Response, req *logical.Request) *logical.Response {
	if !config.EnableMetadataOnFailures || !resp.IsError() {
		return resp
	}
	var initialErrMsg string
	if err := resp.Error(); err != nil {
		initialErrMsg = err.Error()
	}

	clientCert, exists := getClientCert(req)
	if !exists {
		return logical.ErrorResponse("no client certificate found\n" + initialErrMsg)
	}

	// Trim these values as they can be anything from any sort of failed certificate
	// and we don't want to expose audit entries to randomly large strings.
	const maxChars = 100
	metadata := map[string]string{
		"common_name":      trimToMaxChars(clientCert.Subject.CommonName, maxChars),
		"serial_number":    trimToMaxChars(clientCert.SerialNumber.String(), maxChars),
		"subject_key_id":   trimToMaxChars(certutil.GetHexFormatted(clientCert.SubjectKeyId, ":"), maxChars),
		"authority_key_id": trimToMaxChars(certutil.GetHexFormatted(clientCert.AuthorityKeyId, ":"), maxChars),
	}

	return logical.ErrorResponseWithData(metadata, initialErrMsg)
}

func getClientCert(req *logical.Request) (*x509.Certificate, bool) {
	if req == nil || req.Connection == nil || req.Connection.ConnState == nil || req.Connection.ConnState.PeerCertificates == nil {
		return nil, false
	}
	clientCerts := req.Connection.ConnState.PeerCertificates
	if len(clientCerts) == 0 {
		return nil, false
	}
	clientCert := clientCerts[0]
	if clientCert == nil || clientCert.IsCA {
		return nil, false
	}
	return clientCert, true
}

func trimToMaxChars(formatted string, maxSize int) string {
	if len(formatted) > maxSize {
		return formatted[:maxSize-3] + "..."
	}

	return formatted
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if b.configUpdated.Load() {
		b.updatedConfig(config)
	}

	if !config.DisableBinding {
		var matched *ParsedCert
		if verifyResp, resp, err := b.verifyCredentials(ctx, req, d); err != nil {
			return nil, err
		} else if resp != nil {
			return certAuthLoginFailureResponse(config, resp, req), nil
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
		// For the identity of the certificate to match, both the subject key ID and authority key Id should match
		if req.Auth.InternalData["subject_key_id"] != skid {
			return nil, fmt.Errorf("client identity during renewal not matching client identity used during login")
		}
		if req.Auth.InternalData["authority_key_id"] != akid {
			return nil, fmt.Errorf("client authority identity during renewal not matching client identity used during login")
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

	if !policyutil.EquivalentPolicies(cert.TokenPolicies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = cert.TokenTTL
	resp.Auth.MaxTTL = cert.TokenMaxTTL
	resp.Auth.Period = cert.TokenPeriod
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
	} else if d != nil { // d is nil if handleAuthRenew call the authRenew
		certName = d.Get("name").(string)
	}

	// Load the trusted certificates and other details
	roots, trusted, trustedNonCAs, verifyConf := b.getTrustedCerts(ctx, req.Storage, certName)

	// Get the list of full chains matching the connection and validates the
	// certificate itself
	trustedChains, err := validateConnState(roots, connState)
	if err != nil {
		return nil, nil, err
	}

	var extraCas []*x509.Certificate
	for _, t := range trusted {
		extraCas = append(extraCas, t.Certificates...)
	}

	// If trustedNonCAs is not empty it means that client had registered a non-CA cert
	// with the backend.
	var retErr error
	if len(trustedNonCAs) != 0 {
		for _, trustedNonCA := range trustedNonCAs {
			tCert := trustedNonCA.Certificates[0]
			// Check for client cert being explicitly listed in the config (and matching other constraints)
			if tCert.SerialNumber.Cmp(clientCert.SerialNumber) == 0 {
				matches, err := b.matchesConstraints(ctx, clientCert, trustedNonCA.Certificates, trustedNonCA, verifyConf)

				if matches {
					if !tCert.Equal(clientCert) {
						// Someone may be trying to pass off a forged certificate as the trusted non-CA cert.  Reject early.
						return nil, logical.ErrorResponse("certificate mismatch of a trusted leaf certificate"), nil
					}
				}

				// matchesConstraints returns an error when OCSP verification fails,
				// but some other path might still give us success. Add to the
				// retErr multierror, but avoid duplicates. This way, if we reach a
				// failure later, we can give additional context.
				//
				// XXX: If matchesConstraints is updated to generate additional,
				// immediately fatal errors, we likely need to extend it to return
				// another boolean (fatality) or other detection scheme.
				if err != nil && (retErr == nil || !errwrap.Contains(retErr, err.Error())) {
					retErr = multierror.Append(retErr, err)
				}

				if matches {
					return trustedNonCA, nil, nil
				}
			}
		}
	}

	// If no trusted chain was found, client is not authenticated
	// This check happens after checking for a matching configured non-CA certs
	if len(trustedChains) == 0 {
		if retErr != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("%s; additionally got errors during verification: %v", certAuthFailMsg, retErr)), nil
		}

		return nil, logical.ErrorResponse(certAuthFailMsg), nil
	}

	// Search for a ParsedCert that intersects with the validated chains and any additional constraints
	for _, trust := range trusted { // For each ParsedCert in the config
		for _, tCert := range trust.Certificates { // For each certificate in the entry
			for _, chain := range trustedChains { // For each root chain that we matched
				for _, cCert := range chain { // For each cert in the matched chain
					if tCert.Equal(cCert) { // ParsedCert intersects with matched chain
						match, err := b.matchesConstraints(ctx, clientCert, chain, trust, verifyConf) // validate client cert + matched chain against the config

						// See note above.
						if err != nil && (retErr == nil || !errwrap.Contains(retErr, err.Error())) {
							retErr = multierror.Append(retErr, err)
						}

						// Return the first matching entry (for backwards
						// compatibility, we continue to just pick the first
						// one if we have multiple matches).
						//
						// Here, we return directly: this means that any
						// future OCSP errors would be ignored; in the future,
						// if these become fatal, we could revisit this
						// choice and choose the first match after evaluating
						// all possible candidates.
						if match && err == nil {
							return trust, nil, nil
						}
					}
				}
			}
		}
	}

	if retErr != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("%s; additionally got errors during verification: %v", certAuthFailMsg, retErr)), nil
	}

	return nil, logical.ErrorResponse(certAuthFailMsg), nil
}

func (b *backend) matchesConstraints(ctx context.Context, clientCert *x509.Certificate, trustedChain []*x509.Certificate,
	config *ParsedCert, conf *ocsp.VerifyConfig,
) (bool, error) {
	soFar := !b.checkForChainInCRLs(trustedChain) &&
		b.matchesNames(clientCert, config) &&
		b.matchesCommonName(clientCert, config) &&
		b.matchesDNSSANs(clientCert, config) &&
		b.matchesEmailSANs(clientCert, config) &&
		b.matchesURISANs(clientCert, config) &&
		b.matchesOrganizationalUnits(clientCert, config) &&
		b.matchesOrganizations(clientCert, config) &&
		b.matchesCertificateExtensions(clientCert, config)
	if config.Entry.OcspEnabled {
		ocspGood, err := b.checkForCertInOCSP(ctx, clientCert, trustedChain, conf)
		if err != nil {
			return false, err
		}
		soFar = soFar && ocspGood
	}
	return soFar, nil
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

// matchesOrganizationalUnits verifies that the certificate matches at least one configured allowed OU
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

// matchesOrganizations verifies that the certificate matches at least one configured allowed organization (O)
func (b *backend) matchesOrganizations(clientCert *x509.Certificate, config *ParsedCert) bool {
	// Default behavior (no Organizations (Os)) is to allow all organizations (Os)
	if len(config.Entry.AllowedOrganizations) == 0 {
		return true
	}

	// At least one pattern must match at least one name if any patterns are specified
	for _, allowedOrganization := range config.Entry.AllowedOrganizations {
		for _, o := range clientCert.Subject.Organization {
			if glob.Glob(allowedOrganization, o) {
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
	hexExtMap := make(map[string]string, len(clientCert.Extensions))

	for _, ext := range clientCert.Extensions {
		var parsedValue string
		_, err := asn1.Unmarshal(ext.Value, &parsedValue)
		if err != nil {
			clientExtMap[ext.Id.String()] = ""
		} else {
			clientExtMap[ext.Id.String()] = parsedValue
		}

		hexExtMap[ext.Id.String()] = hex.EncodeToString(ext.Value)
	}

	// If any of the required extensions don't match the constraint fails
	for _, requiredExt := range config.Entry.RequiredExtensions {
		reqExt := strings.SplitN(requiredExt, ":", 2)
		if len(reqExt) != 2 {
			return false
		}

		if reqExt[0] == "hex" {
			reqHexExt := strings.SplitN(reqExt[1], ":", 2)
			if len(reqHexExt) != 2 {
				return false
			}

			clientExtValue, clientExtValueOk := hexExtMap[reqHexExt[0]]
			if !clientExtValueOk || !glob.Glob(strings.ToLower(reqHexExt[1]), clientExtValue) {
				return false
			}
		} else {
			clientExtValue, clientExtValueOk := clientExtMap[reqExt[0]]
			if !clientExtValueOk || !glob.Glob(reqExt[1], clientExtValue) {
				return false
			}
		}
	}
	return true
}

// certificateExtensionsMetadata returns the metadata from configured
// metadata extensions
func (b *backend) certificateExtensionsMetadata(clientCert *x509.Certificate, config *ParsedCert) map[string]string {
	// If no metadata extensions are configured, return an empty map
	if len(config.Entry.AllowedMetadataExtensions) == 0 {
		return map[string]string{}
	}

	// Build a map with the accepted oid strings as keys, and the metadata keys as values.
	allowedOidMap := make(map[string]string, len(config.Entry.AllowedMetadataExtensions))
	for _, oidString := range config.Entry.AllowedMetadataExtensions {
		// Avoid dots in metadata keys and put dashes instead,
		// to allow use policy templates.
		allowedOidMap[oidString] = strings.ReplaceAll(oidString, ".", "-")
	}

	// Collect the metadata from accepted certificate extensions.
	metadata := make(map[string]string, len(config.Entry.AllowedMetadataExtensions))
	for _, ext := range clientCert.Extensions {
		if metadataKey, ok := allowedOidMap[ext.Id.String()]; ok {
			// x509 Writes Extensions in ASN1 with a bitstring tag, which results in the field
			// including its ASN.1 type tag bytes. For the sake of simplicity, assume string type
			// and drop the tag bytes. And get the number of bytes from the tag.
			var parsedValue string
			asn1.Unmarshal(ext.Value, &parsedValue)
			metadata[metadataKey] = parsedValue
		}
	}

	return metadata
}

// getTrustedCerts is used to load all the trusted certificates from the backend, cached

func (b *backend) getTrustedCerts(ctx context.Context, storage logical.Storage, certName string) (pool *x509.CertPool, trusted []*ParsedCert, trustedNonCAs []*ParsedCert, conf *ocsp.VerifyConfig) {
	if !b.trustedCacheDisabled.Load() {
		trusted, complete := b.getTrustedCertsFromCache(certName)
		if complete {
			return trusted.pool, trusted.trusted, trusted.trustedNonCAs, trusted.ocspConf
		}
	}
	return b.loadTrustedCerts(ctx, storage, certName)
}

func (b *backend) getTrustedCertsFromCache(certName string) (*trusted, bool) {
	var trusted *trusted
	if certName == "" {
		trusted = b.trustedCacheFull.Load()
	} else {
		trusted, _ = b.trustedCache.Get(certName)
	}

	if trusted == nil {
		return nil, false
	}

	// We're complete (for our purposes here) if we're really complete
	// (because retry is nil) or if it's just not time to retry the load yet.
	complete := (trusted.retry == nil) || time.Now().Before(trusted.retry.deadline)
	return trusted, complete
}

// loadTrustedCerts is used to load all the trusted certificates from the backend
func (b *backend) loadTrustedCerts(ctx context.Context, storage logical.Storage, certName string) (pool *x509.CertPool, trustedCerts []*ParsedCert, trustedNonCAs []*ParsedCert, conf *ocsp.VerifyConfig) {
	lock := locksutil.LockForKey(b.trustedCacheLocks, certName)
	lock.Lock()
	defer lock.Unlock()

	var cache *trusted
	if !b.trustedCacheDisabled.Load() {
		var complete bool
		cache, complete = b.getTrustedCertsFromCache(certName)
		if complete {
			return cache.pool, cache.trusted, cache.trustedNonCAs, cache.ocspConf
		}
	}

	if cache == nil {
		cache = &trusted{
			pool:          x509.NewCertPool(),
			trusted:       make([]*ParsedCert, 0),
			trustedNonCAs: make([]*ParsedCert, 0),
			loaded:        make(map[string]struct{}),
			ocspConf:      &ocsp.VerifyConfig{},
		}
	} else {
		cache = cache.clone()
	}

	pool = cache.pool
	trustedCerts = cache.trusted
	trustedNonCAs = cache.trustedNonCAs
	conf = cache.ocspConf

	var names []string
	if certName != "" {
		names = append(names, certName)
	} else {
		var err error
		names, err = storage.List(ctx, trustedCertPath)
		if err != nil {
			b.Logger().Error("failed to list trusted certs", "error", err)
			return
		}
	}

	anyErrors := false
	for _, name := range names {
		if _, found := cache.loaded[name]; found {
			continue
		}

		entry, parsed, ocsp_ca_certs := b.loadTrustedCert(ctx, storage, name)
		if entry == nil {
			anyErrors = true
			continue
		}
		parsed = append(parsed, ocsp_ca_certs...)

		// NOTE: From this point on please finish adding the cert to all the
		// appropriate lists and pools and configuration! Perform any error
		// checking above this line.
		cache.loaded[name] = struct{}{}

		if !parsed[0].IsCA {
			trustedNonCAs = append(trustedNonCAs, &ParsedCert{
				Entry:        entry,
				Certificates: parsed,
			})
		} else {
			for _, p := range parsed {
				pool.AddCert(p)
			}

			trustedCerts = append(trustedCerts, &ParsedCert{
				Entry:        entry,
				Certificates: parsed,
			})
		}

		if entry.OcspEnabled {
			conf.OcspEnabled = true
			conf.OcspServersOverride = append(conf.OcspServersOverride, entry.OcspServersOverride...)
			if entry.OcspFailOpen {
				conf.OcspFailureMode = ocsp.FailOpenTrue
			} else {
				conf.OcspFailureMode = ocsp.FailOpenFalse
			}
			conf.QueryAllServers = conf.QueryAllServers || entry.OcspQueryAllServers
			conf.OcspThisUpdateMaxAge = entry.OcspThisUpdateMaxAge
			conf.OcspMaxRetries = entry.OcspMaxRetries

			if len(ocsp_ca_certs) > 0 {
				conf.ExtraCas = ocsp_ca_certs
			}
		}
	}

	if !b.trustedCacheDisabled.Load() {
		if anyErrors {
			// If something went wrong then we are going to set up for an
			// exponential backoff on reloading the certificates.
			if cache.retry == nil {
				cache.retry = &trustedRetry{}
			}

			// Limits are arbitrary. Max of 2^55 backoff just so that delay
			// and the the jitter fits into a double. Is that silly? Yes.
			// 2^55 seconds is more years than the computer will keep
			// running.
			if cache.retry.attempt < 55 {
				cache.retry.attempt += 1
			}

			d := 1 << cache.retry.attempt
			pct := (rand.Float64() * 2.0) - 1.0      // between -100% and +100%
			d += int(math.Floor(float64(d/4) * pct)) // between -25% and +25%
			cache.retry.deadline = time.Now().Add(time.Duration(d) * time.Second)
		} else {
			// No problems, cache is complete, no need to retry.
			cache.retry = nil
		}

		cache.trustedNonCAs = trustedNonCAs
		cache.trusted = trustedCerts
		if certName == "" {
			b.trustedCacheFull.Store(cache)
		} else {
			b.trustedCache.Add(certName, cache)
		}
	}
	return
}

func (b *backend) loadTrustedCert(ctx context.Context, storage logical.Storage, name string) (*CertEntry, []*x509.Certificate, []*x509.Certificate) {
	entry, err := b.Cert(ctx, storage, strings.TrimPrefix(name, trustedCertPath))
	if err != nil {
		b.Logger().Error("failed to load trusted cert", "name", name, "error", err)
		return nil, nil, nil
	}

	if entry == nil {
		// This could happen when the certName was provided and the cert doesn't exist,
		// or just if between the LIST and the GET the cert was deleted.
		b.Logger().Error("loaded a nil trusted cert", "name", name)
		return nil, nil, nil
	}

	parsed := parsePEM([]byte(entry.Certificate))
	if len(parsed) == 0 {
		b.Logger().Error("failed to parse certificate", "name", name)
		return nil, nil, nil
	}

	var ocsp_ca_certs []*x509.Certificate = nil
	if len(entry.OcspCaCertificates) > 0 {
		ocsp_ca_certs, err = certutil.ParseCertsPEM([]byte(entry.OcspCaCertificates))
		if err != nil {
			// NOTE: For compatibility, failure to parse the OcspCaCertificates
			// is never actually fatal to loading the broader file
			b.Logger().Error("failed to parse ocsp_ca_certificates", "name", name, "error", err)
			ocsp_ca_certs = nil
		}
	}

	return entry, parsed, ocsp_ca_certs
}

func (b *backend) checkForCertInOCSP(ctx context.Context, clientCert *x509.Certificate, chain []*x509.Certificate, conf *ocsp.VerifyConfig) (bool, error) {
	if !conf.OcspEnabled || len(chain) < 2 {
		return true, nil
	}
	b.ocspClientMutex.RLock()
	defer b.ocspClientMutex.RUnlock()
	err := b.ocspClient.VerifyLeafCertificate(ctx, clientCert, chain[1], conf)
	if err != nil {
		if ocsp.IsOcspVerificationError(err) {
			// We don't want anything to override an OCSP verification error
			return false, err
		}
		if conf.OcspFailureMode == ocsp.FailOpenTrue {
			onlyNetworkErrors := b.handleOcspErrorInFailOpen(err)
			if onlyNetworkErrors {
				return true, nil
			}
		}
		// We want to preserve error messages when they have additional,
		// potentially useful information. Just having a revoked cert
		// isn't additionally useful.
		if !strings.Contains(err.Error(), "has been revoked") {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (b *backend) handleOcspErrorInFailOpen(err error) bool {
	urlError := &url.Error{}
	allNetworkErrors := true
	if multiError, ok := err.(*multierror.Error); ok {
		for _, myErr := range multiError.Errors {
			if !errors.As(myErr, &urlError) {
				allNetworkErrors = false
			}
		}
	} else if !errors.As(err, &urlError) {
		allNetworkErrors = false
	}

	if allNetworkErrors {
		b.Logger().Warn("OCSP is set to fail-open, and could not retrieve "+
			"OCSP based revocation but proceeding.", "detail", err)
		return true
	}

	return false
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
