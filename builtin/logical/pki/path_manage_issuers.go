// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathIssuerGenerateRoot(b *backend) *framework.Path {
	pattern := "issuers/generate/root/" + framework.GenericNameRegex("exported")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKIIssuers,
		OperationVerb:   "generate",
		OperationSuffix: "root",
	}

	return buildPathGenerateRoot(b, pattern, displayAttrs)
}

func pathRotateRoot(b *backend) *framework.Path {
	pattern := "root/rotate/" + framework.GenericNameRegex("exported")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKI,
		OperationVerb:   "rotate",
		OperationSuffix: "root",
	}

	return buildPathGenerateRoot(b, pattern, displayAttrs)
}

func buildPathGenerateRoot(b *backend, pattern string, displayAttrs *framework.DisplayAttributes) *framework.Path {
	ret := &framework.Path{
		Pattern:      pattern,
		DisplayAttrs: displayAttrs,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCAGenerateRoot,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"expiration": {
								Type:        framework.TypeInt64,
								Description: `The expiration of the given issuer.`,
								Required:    true,
							},
							"serial_number": {
								Type:        framework.TypeString,
								Description: `The requested Subject's named serial number.`,
								Required:    true,
							},
							"certificate": {
								Type:        framework.TypeString,
								Description: `The generated self-signed CA certificate.`,
								Required:    true,
							},
							"issuing_ca": {
								Type:        framework.TypeString,
								Description: `The issuing certificate authority.`,
								Required:    true,
							},
							"issuer_id": {
								Type:        framework.TypeString,
								Description: `The ID of the issuer`,
								Required:    true,
							},
							"issuer_name": {
								Type:        framework.TypeString,
								Description: `The name of the issuer.`,
								Required:    true,
							},
							"key_id": {
								Type:        framework.TypeString,
								Description: `The ID of the key.`,
								Required:    true,
							},
							"key_name": {
								Type:        framework.TypeString,
								Description: `The key name if given.`,
								Required:    true,
							},
							"private_key": {
								Type:        framework.TypeString,
								Description: `The private key if exported was specified.`,
								Required:    false,
							},
						},
					}},
				},
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathGenerateRootHelpSyn,
		HelpDescription: pathGenerateRootHelpDesc,
	}

	ret.Fields = addCACommonFields(map[string]*framework.FieldSchema{})
	ret.Fields = addCAKeyGenerationFields(ret.Fields)
	ret.Fields = addCAIssueFields(ret.Fields)
	ret.Fields = addCACertKeyUsage(ret.Fields)
	return ret
}

func pathIssuerGenerateIntermediate(b *backend) *framework.Path {
	pattern := "issuers/generate/intermediate/" + framework.GenericNameRegex("exported")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKIIssuers,
		OperationVerb:   "generate",
		OperationSuffix: "intermediate",
	}

	return buildPathGenerateIntermediate(b, pattern, displayAttrs)
}

func pathCrossSignIntermediate(b *backend) *framework.Path {
	pattern := "intermediate/cross-sign"

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKI,
		OperationVerb:   "cross-sign",
		OperationSuffix: "intermediate",
	}

	return buildPathGenerateIntermediate(b, pattern, displayAttrs)
}

func buildPathGenerateIntermediate(b *backend, pattern string, displayAttrs *framework.DisplayAttributes) *framework.Path {
	ret := &framework.Path{
		Pattern:      pattern,
		DisplayAttrs: displayAttrs,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathGenerateIntermediate,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"csr": {
								Type:        framework.TypeString,
								Description: `Certificate signing request.`,
								Required:    true,
							},
							"key_id": {
								Type:        framework.TypeString,
								Description: `Id of the key.`,
								Required:    true,
							},
							"private_key": {
								Type:        framework.TypeString,
								Description: `Generated private key.`,
								Required:    false,
							},
							"private_key_type": {
								Type:        framework.TypeString,
								Description: `Specifies the format used for marshaling the private key.`,
								Required:    false,
							},
						},
					}},
				},
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathGenerateIntermediateHelpSyn,
		HelpDescription: pathGenerateIntermediateHelpDesc,
	}

	ret.Fields = addCACommonFields(map[string]*framework.FieldSchema{})
	ret.Fields = addCAKeyGenerationFields(ret.Fields)
	ret.Fields["add_basic_constraints"] = &framework.FieldSchema{
		Type: framework.TypeBool,
		Description: `Whether to add a Basic Constraints
extension with CA: true. Only needed as a
workaround in some compatibility scenarios
with Active Directory Certificate Services.`,
	}
	ret.Fields = addCaCsrKeyUsage(ret.Fields)

	// At this time Go does not support signing CSRs using PSS signatures, see
	// https://github.com/golang/go/issues/45990
	delete(ret.Fields, "use_pss")

	return ret
}

func pathImportIssuer(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issuers/import/(cert|bundle)",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKIIssuers,
			OperationVerb:   "import",
			OperationSuffix: "cert|bundle",
		},

		Fields: map[string]*framework.FieldSchema{
			"pem_bundle": {
				Type: framework.TypeString,
				Description: `PEM-format, concatenated unencrypted
secret-key (optional) and certificates.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathImportIssuers,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"mapping": {
								Type:        framework.TypeMap,
								Description: "A mapping of issuer_id to key_id for all issuers included in this request",
								Required:    true,
							},
							"imported_keys": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Net-new keys imported as a part of this request",
								Required:    true,
							},
							"imported_issuers": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Net-new issuers imported as a part of this request",
								Required:    true,
							},
							"existing_keys": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Existing keys specified as part of the import bundle of this request",
								Required:    true,
							},
							"existing_issuers": {
								Type:        framework.TypeCommaStringSlice,
								Description: "Existing issuers specified as part of the import bundle of this request",
								Required:    true,
							},
						},
					}},
				},
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathImportIssuersHelpSyn,
		HelpDescription: pathImportIssuersHelpDesc,
	}
}

func (b *backend) pathImportIssuers(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	keysAllowed := strings.HasSuffix(req.Path, "bundle") || req.Path == "config/ca"

	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not import issuers until migration has completed"), nil
	}

	var pemBundle string
	var certificate string
	rawPemBundle, bundleOk := data.GetOk("pem_bundle")
	rawCertificate, certOk := data.GetOk("certificate")
	if bundleOk {
		pemBundle = rawPemBundle.(string)
	}
	if certOk {
		certificate = rawCertificate.(string)
	}

	if len(pemBundle) == 0 && len(certificate) == 0 {
		return logical.ErrorResponse("'pem_bundle' and 'certificate' parameters were empty"), nil
	}
	if len(pemBundle) > 0 && len(certificate) > 0 {
		return logical.ErrorResponse("'pem_bundle' and 'certificate' parameters were both provided"), nil
	}
	if len(certificate) > 0 {
		keysAllowed = false
		pemBundle = certificate
	}
	if len(pemBundle) < 75 {
		// It is almost nearly impossible to store a complete certificate in
		// less than 75 bytes. It is definitely impossible to do so when PEM
		// encoding has been applied. Detect this and give a better warning
		// than "provided PEM block contained no data" in this case. This is
		// because the PEM headers contain 5*4 + 6 + 4 + 2 + 2 = 34 characters
		// minimum (five dashes, "BEGIN" + space + at least one character
		// identifier, "END" + space + at least one character identifier, and
		// a pair of new lines). That would leave 41 bytes for Base64 data,
		// meaning at most a 30-byte DER certificate.
		//
		// However, < 75 bytes is probably a good length for a file path so
		// suggest that is the case.
		return logical.ErrorResponse("provided data for import was too short; perhaps a path was passed to the API rather than the contents of a PEM file"), nil
	}

	var createdKeys []string
	var createdIssuers []string
	var existingKeys []string
	var existingIssuers []string
	issuerKeyMap := make(map[string]string)

	// Rather than using certutil.ParsePEMBundle (which restricts the
	// construction of the PEM bundle), we manually parse the bundle instead.
	pemBytes := []byte(pemBundle)
	var pemBlock *pem.Block

	var issuers []string
	var keys []string

	// By decoding and re-encoding PEM blobs, we can pass strict PEM blobs
	// to the import functionality (importKeys, importIssuers). This allows
	// them to validate no duplicate issuers exist (and place greater
	// restrictions during parsing) but allows this code to accept OpenSSL
	// parsed chains (with full textual output between PEM entries).
	for len(bytes.TrimSpace(pemBytes)) > 0 {
		pemBlock, pemBytes = pem.Decode(pemBytes)
		if pemBlock == nil {
			return logical.ErrorResponse("provided PEM block contained no data"), nil
		}

		pemBlockString := string(pem.EncodeToMemory(pemBlock))

		switch pemBlock.Type {
		case "CERTIFICATE", "X509 CERTIFICATE":
			// Must be a certificate
			issuers = append(issuers, pemBlockString)
		case "CRL", "X509 CRL":
			// Ignore any CRL entries.
		case "EC PARAMS", "EC PARAMETERS":
			// Ignore any EC parameter entries. This is an optional block
			// that some implementations send, to ensure some semblance of
			// compatibility with weird curves. Go doesn't support custom
			// curves and 99% of software doesn't either, so discard them
			// without parsing them.
		default:
			// Otherwise, treat them as keys.
			keys = append(keys, pemBlockString)
		}
	}

	if len(keys) > 0 && !keysAllowed {
		return logical.ErrorResponse("private keys found in the PEM bundle but not allowed by the path; use /issuers/import/bundle"), nil
	}

	sc := b.makeStorageContext(ctx, req.Storage)

	for keyIndex, keyPem := range keys {
		// Handle import of private key.
		key, existing, err := importKeyFromBytes(sc, keyPem, "")
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing key %v: %v", keyIndex, err)), nil
		}

		if !existing {
			createdKeys = append(createdKeys, key.ID.String())
		} else {
			existingKeys = append(existingKeys, key.ID.String())
		}
	}

	for certIndex, certPem := range issuers {
		cert, existing, err := sc.importIssuer(certPem, "")
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing issuer %v: %v\n%v", certIndex, err, certPem)), nil
		}

		issuerKeyMap[cert.ID.String()] = cert.KeyID.String()
		if !existing {
			createdIssuers = append(createdIssuers, cert.ID.String())
		} else {
			existingIssuers = append(existingIssuers, cert.ID.String())
		}
	}

	response := &logical.Response{
		Data: map[string]interface{}{
			"mapping":          issuerKeyMap,
			"imported_keys":    createdKeys,
			"imported_issuers": createdIssuers,
			"existing_keys":    existingKeys,
			"existing_issuers": existingIssuers,
		},
	}

	if len(createdIssuers) > 0 {
		warnings, err := b.CrlBuilder().Rebuild(sc, true)
		if err != nil {
			// Before returning, check if the error message includes the
			// string "PSS". If so, it indicates we might've wanted to modify
			// this issuer, so convert the error to a warning.
			if strings.Contains(err.Error(), "PSS") || strings.Contains(err.Error(), "pss") {
				err = fmt.Errorf("Rebuilding the CRL failed with a message relating to the PSS signature algorithm. This likely means the revocation_signature_algorithm needs to be set on the newly imported issuer(s) because a managed key supports only the PSS algorithm; by default PKCS#1v1.5 was used to build the CRLs. CRLs will not be generated until this has been addressed, however the import was successful. The original error is reproduced below:\n\n\t%w", err)
			} else {
				// Note to the caller that while this is an error, we did
				// successfully import the issuers.
				err = fmt.Errorf("Rebuilding the CRL failed. While this is indicative of a problem with the imported issuers (perhaps because of their revocation_signature_algorithm), they did import successfully and are now usable. It is strongly suggested to fix the CRL building errors before continuing. The original error is reproduced below:\n\n\t%w", err)
			}

			return nil, err
		}
		for index, warning := range warnings {
			response.AddWarning(fmt.Sprintf("Warning %d during CRL rebuild: %v", index+1, warning))
		}

		var issuersWithKeys []string
		for _, issuer := range createdIssuers {
			if issuerKeyMap[issuer] != "" {
				issuersWithKeys = append(issuersWithKeys, issuer)
			}
		}

		// Check whether we need to update our default issuer configuration.
		config, err := sc.getIssuersConfig()
		if err != nil {
			response.AddWarning("Unable to fetch default issuers configuration to update default issuer if necessary: " + err.Error())
		} else if config.DefaultFollowsLatestIssuer {
			if len(issuersWithKeys) == 1 {
				if err := sc.updateDefaultIssuerId(issuing.IssuerID(issuersWithKeys[0])); err != nil {
					response.AddWarning("Unable to update this new root as the default issuer: " + err.Error())
				}
			} else if len(issuersWithKeys) > 1 {
				response.AddWarning("Default issuer left unchanged: could not select new issuer automatically as multiple imported issuers had key material in Vault.")
			}
		}
	}

	// While we're here, check if we should warn about a bad default key. We
	// do this unconditionally if the issuer or key was modified, so the admin
	// is always warned. But if unrelated key material was imported, we do
	// not warn.
	config, err := sc.getIssuersConfig()
	if err == nil && len(config.DefaultIssuerId) > 0 {
		// We can use the mapping above to check the issuer mapping.
		if keyId, ok := issuerKeyMap[string(config.DefaultIssuerId)]; ok && len(keyId) == 0 {
			msg := "The default issuer has no key associated with it. Some operations like issuing certificates and signing CRLs will be unavailable with the requested default issuer until a key is imported or the default issuer is changed."
			response.AddWarning(msg)
			b.Logger().Error(msg)
		}

		// If we imported multiple issuers with keys (or matched existing
		// keys), and we set one of those as a default, warn the end-user we
		// might have selected the wrong one.
		if len(createdIssuers) > 1 {
			numCreatedIssuersWithKeys := 0
			defaultIssuerWasCreated := false
			for _, issuerId := range createdIssuers {
				if keyId, ok := issuerKeyMap[issuerId]; ok && len(keyId) != 0 {
					numCreatedIssuersWithKeys++
				}

				if config.DefaultIssuerId.String() == issuerId {
					defaultIssuerWasCreated = true
				}
			}

			if numCreatedIssuersWithKeys > 1 && defaultIssuerWasCreated {
				msg := "The imported bundle contained multiple certs matching keys, " +
					"the default issuer that was selected should be verified and manually changed if incorrect."
				response.AddWarning(msg)
				b.Logger().Error(msg)
			}
		}
	}

	// Also while we're here, we should let the user know the next steps.
	// In particular, if there's no default AIA URLs configuration, we should
	// tell the user that's probably next.
	if entries, err := getGlobalAIAURLs(ctx, req.Storage); err == nil && len(entries.IssuingCertificates) == 0 && len(entries.CRLDistributionPoints) == 0 && len(entries.OCSPServers) == 0 {
		response.AddWarning("This mount hasn't configured any authority information access (AIA) fields; this may make it harder for systems to find missing certificates in the chain or to validate revocation status of certificates. Consider updating /config/urls or the newly generated issuer with this information.")
	}

	return response, nil
}

const (
	pathImportIssuersHelpSyn  = `Import the specified issuing certificates.`
	pathImportIssuersHelpDesc = `
This endpoint allows importing the specified issuer certificates.

:type is either the literal value "cert", to only allow importing
certificates, else "bundle" to allow importing keys as well as
certificates.

Depending on the value of :type, the pem_bundle request parameter can
either take PEM-formatted certificates, and, if :type="bundle", unencrypted
secret-keys.
`
)

func pathRevokeIssuer(b *backend) *framework.Path {
	fields := addIssuerRefField(map[string]*framework.FieldSchema{})
	responseFields := issuerResponseFields(true)

	return &framework.Path{
		Pattern: "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/revoke",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "revoke",
			OperationSuffix: "issuer",
		},

		Fields: fields,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathRevokeIssuer,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      responseFields,
					}},
				},
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathRevokeIssuerHelpSyn,
		HelpDescription: pathRevokeIssuerHelpDesc,
	}
}

func (b *backend) pathRevokeIssuer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	// Issuer revocation can't work on the legacy cert bundle.
	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("cannot revoke issuer until migration has completed"), nil
	}

	issuerName := GetIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	// Fetch the issuer.
	sc := b.makeStorageContext(ctx, req.Storage)
	ref, err := sc.resolveIssuerReference(issuerName)
	if err != nil {
		return nil, err
	}
	if ref == "" {
		return logical.ErrorResponse("unable to resolve issuer id for reference: " + issuerName), nil
	}

	issuer, err := sc.fetchIssuerById(ref)
	if err != nil {
		return nil, err
	}

	// If its already been revoked, just return the read results sans warnings
	// like we would otherwise.
	if issuer.Revoked {
		return respondReadIssuer(issuer)
	}

	// When revoking, we want to forbid new certificate issuance. We allow
	// new revocations of leaves issued by this issuer to trigger a CRL
	// rebuild still.
	issuer.Revoked = true
	if issuer.Usage.HasUsage(issuing.IssuanceUsage) {
		issuer.Usage.ToggleUsage(issuing.IssuanceUsage)
	}

	currTime := time.Now()
	issuer.RevocationTime = currTime.Unix()
	issuer.RevocationTimeUTC = currTime.UTC()

	err = sc.writeIssuer(issuer)
	if err != nil {
		return nil, err
	}

	// Now, if the parent issuer exists within this mount, we'd have written
	// a storage entry for this certificate, making it appear as any other
	// leaf. We need to add a RevocationInfo entry for this into storage,
	// so that it appears as if it was revoked.
	//
	// This is a _necessary_ but not necessarily _sufficient_ step to
	// consider an arbitrary issuer revoked and the former step (setting
	// issuer.Revoked = true) is more correct: if two intermediates have the
	// same serial number, and one appears somehow in the storage but from a
	// different issuer, we'd only include one in the CRLs, but we'd want to
	// include both in two separate CRLs. Hence, the former is the condition
	// we check in CRL building, but this step satisfies other guarantees
	// within Vault.
	certEntry, err := fetchCertBySerial(sc, issuing.PathCerts, issuer.SerialNumber)
	if err == nil && certEntry != nil {
		// We've inverted this error check as it doesn't matter; we already
		// consider this certificate revoked.
		storageCert, err := x509.ParseCertificate(certEntry.Value)
		if err != nil {
			return nil, fmt.Errorf("error parsing stored certificate value: %w", err)
		}

		issuerCert, err := issuer.GetCertificate()
		if err != nil {
			return nil, fmt.Errorf("error parsing issuer certificate value: %w", err)
		}

		if bytes.Equal(issuerCert.Raw, storageCert.Raw) {
			// If the issuer is on disk at its serial number is the same as
			// our issuer, we know we can write the revocation entry. Since
			// Vault has historically forbid revocation of non-stored certs
			// and issuers, we're the only ones to write this entry, so we
			// don't need the write guard that exists in crl_util.go for the
			// general case (forbidding a newer revocation time).
			//
			// We'll let a cleanup pass or CRL build identify the issuer for
			// us.
			revInfo := revocation.RevocationInfo{
				CertificateBytes:  issuerCert.Raw,
				RevocationTime:    issuer.RevocationTime,
				RevocationTimeUTC: issuer.RevocationTimeUTC,
			}

			revEntry, err := logical.StorageEntryJSON(revokedPath+normalizeSerial(issuer.SerialNumber), revInfo)
			if err != nil {
				return nil, fmt.Errorf("error creating revocation entry for issuer: %w", err)
			}

			err = req.Storage.Put(ctx, revEntry)
			if err != nil {
				return nil, fmt.Errorf("error saving revoked issuer to new location: %w", err)
			}
		}
	}

	// Rebuild the CRL to include the newly revoked issuer.
	warnings, crlErr := b.CrlBuilder().Rebuild(sc, false)
	if crlErr != nil {
		switch crlErr.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
		default:
			return nil, fmt.Errorf("error encountered during CRL building: %w", crlErr)
		}
	}

	// Finally, respond with the issuer's updated data.
	response, err := respondReadIssuer(issuer)
	if err != nil {
		// Impossible.
		return nil, err
	}
	for index, warning := range warnings {
		response.AddWarning(fmt.Sprintf("Warning %d during CRL rebuild: %v", index+1, warning))
	}

	// For sanity, we'll add a warning message here if there's no other
	// issuer which verifies this issuer.
	ourCert, err := issuer.GetCertificate()
	if err != nil {
		return nil, err
	}

	allIssuers, err := sc.listIssuers()
	if err != nil {
		return nil, err
	}

	isSelfSigned := false
	haveOtherIssuer := false
	for _, candidateID := range allIssuers {
		candidate, err := sc.fetchIssuerById(candidateID)
		if err != nil {
			return nil, err
		}

		candidateCert, err := candidate.GetCertificate()
		if err != nil {
			// Returning this error is fine because more things will fail
			// if this issuer can't parse.
			return nil, err
		}

		if err := ourCert.CheckSignatureFrom(candidateCert); err == nil {
			// Signature verification is a success. This means we have a
			// parent for this cert. But notice above we didn't filter out
			// ourselves: we want to see if this is a self-signed cert. So
			// check that now.
			if candidate.ID == issuer.ID {
				isSelfSigned = true
			} else {
				haveOtherIssuer = true
			}
		}

		// If we have both possible warning candidates, no sense continuing
		// to check signatures; exit.
		if isSelfSigned && haveOtherIssuer {
			break
		}
	}

	if isSelfSigned {
		response.AddWarning("This issuer is a self-signed (potentially root) certificate. This means it may not be considered revoked if there is not an external, cross-signed variant of this certificate. This issuer's serial number will not appear on its own CRL.")
	}

	if !haveOtherIssuer {
		response.AddWarning("This issuer lacks another parent issuer within the mount. This means it will not appear on any other CRLs and may not be considered revoked by clients. Consider adding this issuer to its issuer's CRL as well if it is not self-signed.")
	}

	config, err := sc.getIssuersConfig()
	if err == nil && config != nil && config.DefaultIssuerId == issuer.ID {
		response.AddWarning("This issuer is currently configured as the default issuer for this mount; operations such as certificate issuance may not work until a new default issuer is selected.")
	}

	return response, nil
}

const (
	pathRevokeIssuerHelpSyn  = `Revoke the specified issuer certificate.`
	pathRevokeIssuerHelpDesc = `
This endpoint allows revoking the specified issuer certificates.

This is useful when the issuer and its parent exist within the same PKI
mount point (utilizing the multi-issuer functionality). If no suitable
parent is found, this revocation may not appear on any CRL in this mount.

Once revoked, issuers cannot be unrevoked and may not be used to sign any
more certificates.
`
)
