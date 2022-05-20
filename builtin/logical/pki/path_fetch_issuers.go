package pki

import (
	"context"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListIssuers(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issuers/?$",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathListIssuersHandler,
			},
		},

		HelpSynopsis:    pathListIssuersHelpSyn,
		HelpDescription: pathListIssuersHelpDesc,
	}
}

func (b *backend) pathListIssuersHandler(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not list issuers until migration has completed"), nil
	}

	var responseKeys []string
	responseInfo := make(map[string]interface{})

	entries, err := listIssuers(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	config, err := getIssuersConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// For each issuer, we need not only the identifier (as returned by
	// listIssuers), but also the name of the issuer. This means we have to
	// fetch the actual issuer object as well.
	for _, identifier := range entries {
		issuer, err := fetchIssuerById(ctx, req.Storage, identifier)
		if err != nil {
			return nil, err
		}

		responseKeys = append(responseKeys, string(identifier))
		responseInfo[string(identifier)] = map[string]interface{}{
			"issuer_name": issuer.Name,
			"is_default":  identifier == config.DefaultIssuerId,
		}
	}

	return logical.ListResponseWithInfo(responseKeys, responseInfo), nil
}

const (
	pathListIssuersHelpSyn  = `Fetch a list of CA certificates.`
	pathListIssuersHelpDesc = `
This endpoint allows listing of known issuing certificates, returning
their identifier and their name (if set).
`
)

func pathGetIssuer(b *backend) *framework.Path {
	pattern := "issuer/" + framework.GenericNameRegex(issuerRefParam) + "(/der|/pem|/json)?"
	return buildPathGetIssuer(b, pattern)
}

func buildPathGetIssuer(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	fields = addIssuerRefNameFields(fields)

	// Fields for updating issuer.
	fields["manual_chain"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `Chain of issuer references to use to build this
issuer's computed CAChain field, when non-empty.`,
	}
	fields["leaf_not_after_behavior"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Behavior of leaf's NotAfter fields: "err" to error
if the computed NotAfter date exceeds that of this issuer; "truncate" to
silently truncate to that of this issuer; or "permit" to allow this
issuance to succeed (with NotAfter exceeding that of an issuer). Note that
not all values will results in certificates that can be validated through
the entire validity period. It is suggested to use "truncate" for
intermediate CAs and "permit" only for root CAs.`,
		Default: "err",
	}
	fields["usage"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `Comma-separated list (or string slice) of usages for
this issuer; valid values are "read-only", "issuing-certificates", and
"crl-signing". Multiple values may be specified. Read-only is implicit
and always set.`,
		Default: []string{"read-only", "issuing-certificates", "crl-signing"},
	}

	return &framework.Path{
		// Returns a JSON entry.
		Pattern: pattern,
		Fields:  fields,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathGetIssuer,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathUpdateIssuer,
				// Read more about why these flags are set in backend.go.
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathDeleteIssuer,
				// Read more about why these flags are set in backend.go.
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathGetIssuerHelpSyn,
		HelpDescription: pathGetIssuerHelpDesc,
	}
}

func (b *backend) pathGetIssuer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Handle raw issuers first.
	if strings.HasSuffix(req.Path, "/der") || strings.HasSuffix(req.Path, "/pem") || strings.HasSuffix(req.Path, "/json") {
		return b.pathGetRawIssuer(ctx, req, data)
	}

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not get issuer until migration has completed"), nil
	}

	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	ref, err := resolveIssuerReference(ctx, req.Storage, issuerName)
	if err != nil {
		return nil, err
	}
	if ref == "" {
		return logical.ErrorResponse("unable to resolve issuer id for reference: " + issuerName), nil
	}

	issuer, err := fetchIssuerById(ctx, req.Storage, ref)
	if err != nil {
		return nil, err
	}

	return respondReadIssuer(issuer)
}

func respondReadIssuer(issuer *issuerEntry) (*logical.Response, error) {
	var respManualChain []string
	for _, entity := range issuer.ManualChain {
		respManualChain = append(respManualChain, string(entity))
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"issuer_id":               issuer.ID,
			"issuer_name":             issuer.Name,
			"key_id":                  issuer.KeyID,
			"certificate":             issuer.Certificate,
			"manual_chain":            respManualChain,
			"ca_chain":                issuer.CAChain,
			"leaf_not_after_behavior": issuer.LeafNotAfterBehavior.String(),
			"usage":                   issuer.Usage.Names(),
		},
	}, nil
}

func (b *backend) pathUpdateIssuer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not update issuer until migration has completed"), nil
	}

	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	ref, err := resolveIssuerReference(ctx, req.Storage, issuerName)
	if err != nil {
		return nil, err
	}
	if ref == "" {
		return logical.ErrorResponse("unable to resolve issuer id for reference: " + issuerName), nil
	}

	issuer, err := fetchIssuerById(ctx, req.Storage, ref)
	if err != nil {
		return nil, err
	}

	newName, err := getIssuerName(ctx, req.Storage, data)
	if err != nil && err != errIssuerNameInUse {
		// If the error is name already in use, and the new name is the
		// old name for this issuer, we're not actually updating the
		// issuer name (or causing a conflict) -- so don't err out. Other
		// errs should still be surfaced, however.
		return logical.ErrorResponse(err.Error()), nil
	}
	if err == errIssuerNameInUse && issuer.Name != newName {
		// When the new name is in use but isn't this name, throw an error.
		return logical.ErrorResponse(err.Error()), nil
	}

	newPath := data.Get("manual_chain").([]string)
	rawLeafBehavior := data.Get("leaf_not_after_behavior").(string)
	var newLeafBehavior certutil.NotAfterBehavior
	switch rawLeafBehavior {
	case "err":
		newLeafBehavior = certutil.ErrNotAfterBehavior
	case "truncate":
		newLeafBehavior = certutil.TruncateNotAfterBehavior
	case "permit":
		newLeafBehavior = certutil.PermitNotAfterBehavior
	default:
		return logical.ErrorResponse("Unknown value for field `leaf_not_after_behavior`. Possible values are `err`, `truncate`, and `permit`."), nil
	}

	rawUsage := data.Get("usage").([]string)
	newUsage, err := NewIssuerUsageFromNames(rawUsage)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("Unable to parse specified usages: %v - valid values are %v", rawUsage, AllIssuerUsages.Names())), nil
	}

	modified := false

	var oldName string
	if newName != issuer.Name {
		oldName = issuer.Name
		issuer.Name = newName
		modified = true
	}

	if newLeafBehavior != issuer.LeafNotAfterBehavior {
		issuer.LeafNotAfterBehavior = newLeafBehavior
		modified = true
	}

	if newUsage != issuer.Usage {
		issuer.Usage = newUsage
		modified = true
	}

	var updateChain bool
	var constructedChain []issuerID
	for index, newPathRef := range newPath {
		// Allow self for the first entry.
		if index == 0 && newPathRef == "self" {
			newPathRef = string(ref)
		}

		resolvedId, err := resolveIssuerReference(ctx, req.Storage, newPathRef)
		if err != nil {
			return nil, err
		}

		if index == 0 && resolvedId != ref {
			return logical.ErrorResponse(fmt.Sprintf("expected first cert in chain to be a self-reference, but was: %v/%v", newPathRef, resolvedId)), nil
		}

		constructedChain = append(constructedChain, resolvedId)
		if len(issuer.ManualChain) < len(constructedChain) || constructedChain[index] != issuer.ManualChain[index] {
			updateChain = true
		}
	}

	if len(issuer.ManualChain) != len(constructedChain) {
		updateChain = true
	}

	if updateChain {
		issuer.ManualChain = constructedChain

		// Building the chain will write the issuer to disk; no need to do it
		// twice.
		modified = false
		err := rebuildIssuersChains(ctx, req.Storage, issuer)
		if err != nil {
			return nil, err
		}
	}

	if modified {
		err := writeIssuer(ctx, req.Storage, issuer)
		if err != nil {
			return nil, err
		}
	}

	response, err := respondReadIssuer(issuer)
	if newName != oldName {
		addWarningOnDereferencing(oldName, response, ctx, req.Storage)
	}

	return response, err
}

func (b *backend) pathGetRawIssuer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not get issuer until migration has completed"), nil
	}

	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	ref, err := resolveIssuerReference(ctx, req.Storage, issuerName)
	if err != nil {
		return nil, err
	}
	if ref == "" {
		return logical.ErrorResponse("unable to resolve issuer id for reference: " + issuerName), nil
	}

	issuer, err := fetchIssuerById(ctx, req.Storage, ref)
	if err != nil {
		return nil, err
	}

	certificate := []byte(issuer.Certificate)

	var contentType string
	if strings.HasSuffix(req.Path, "/pem") {
		contentType = "application/pem-certificate-chain"
	} else if strings.HasSuffix(req.Path, "/der") {
		contentType = "application/pkix-cert"
	}

	if strings.HasSuffix(req.Path, "/der") {
		pemBlock, _ := pem.Decode(certificate)
		if pemBlock == nil {
			return nil, err
		}

		certificate = pemBlock.Bytes
	}

	statusCode := 200
	if len(certificate) == 0 {
		statusCode = 204
	}

	if strings.HasSuffix(req.Path, "/pem") || strings.HasSuffix(req.Path, "/der") {
		return &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPContentType: contentType,
				logical.HTTPRawBody:     certificate,
				logical.HTTPStatusCode:  statusCode,
			},
		}, nil
	} else {
		return &logical.Response{
			Data: map[string]interface{}{
				"certificate": string(certificate),
				"ca_chain":    issuer.CAChain,
			},
		}, nil
	}
}

func (b *backend) pathDeleteIssuer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not delete issuer until migration has completed"), nil
	}

	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	ref, err := resolveIssuerReference(ctx, req.Storage, issuerName)
	if err != nil {
		// Return as if we deleted it if we fail to lookup the issuer.
		if ref == IssuerRefNotFound {
			return &logical.Response{}, nil
		}
		return nil, err
	}

	response := &logical.Response{}

	issuer, err := fetchIssuerById(ctx, req.Storage, ref)
	if err != nil {
		return nil, err
	}
	if issuer.Name != "" {
		addWarningOnDereferencing(issuer.Name, response, ctx, req.Storage)
	}
	addWarningOnDereferencing(string(issuer.ID), response, ctx, req.Storage)

	wasDefault, err := deleteIssuer(ctx, req.Storage, ref)
	if err != nil {
		return nil, err
	}
	if wasDefault {
		response.AddWarning(fmt.Sprintf("Deleted issuer %v (via issuer_ref %v); this was configured as the default issuer. Operations without an explicit issuer will not work until a new default is configured.", ref, issuerName))
		addWarningOnDereferencing(defaultRef, response, ctx, req.Storage)
	}

	// Since we've deleted an issuer, the chains might've changed. Call the
	// rebuild code. We shouldn't technically err (as the issuer was deleted
	// successfully), but log a warning (and to the response) if this fails.
	if err := rebuildIssuersChains(ctx, req.Storage, nil); err != nil {
		msg := fmt.Sprintf("Failed to rebuild remaining issuers' chains: %v", err)
		b.Logger().Error(msg)
		response.AddWarning(msg)
	}

	return response, nil
}

func addWarningOnDereferencing(name string, resp *logical.Response, ctx context.Context, s logical.Storage) {
	timeout, inUseBy, err := checkForRolesReferencing(name, ctx, s)
	if err != nil || timeout {
		if inUseBy == 0 {
			resp.AddWarning(fmt.Sprint("Unable to check if any roles referenced this issuer by ", name))
		} else {
			resp.AddWarning(fmt.Sprint("The name ", name, " was in use by at least ", inUseBy, " roles"))
		}
	} else {
		if inUseBy > 0 {
			resp.AddWarning(fmt.Sprint(inUseBy, " roles reference ", name))
		}
	}
}

const (
	pathGetIssuerHelpSyn  = `Fetch a single issuer certificate.`
	pathGetIssuerHelpDesc = `
This allows fetching information associated with the underlying issuer
certificate.

:ref can be either the literal value "default", in which case /config/issuers
will be consulted for the present default issuer, an identifier of an issuer,
or its assigned name value.

Use /issuer/:ref/der or /issuer/:ref/pem to return just the certificate in
raw DER or PEM form, without the JSON structure of /issuer/:ref.

Writing to /issuer/:ref allows updating of the name field associated with
the certificate.
`
)

func pathGetIssuerCRL(b *backend) *framework.Path {
	pattern := "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/crl(/pem|/der)?"
	return buildPathGetIssuerCRL(b, pattern)
}

func buildPathGetIssuerCRL(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	fields = addIssuerRefNameFields(fields)

	return &framework.Path{
		// Returns raw values.
		Pattern: pattern,
		Fields:  fields,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathGetIssuerCRL,
			},
		},

		HelpSynopsis:    pathGetIssuerCRLHelpSyn,
		HelpDescription: pathGetIssuerCRLHelpDesc,
	}
}

func (b *backend) pathGetIssuerCRL(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not get issuer's CRL until migration has completed"), nil
	}

	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	if err := b.crlBuilder.rebuildIfForced(ctx, b, req); err != nil {
		return nil, err
	}

	crlPath, err := resolveIssuerCRLPath(ctx, b, req.Storage, issuerName)
	if err != nil {
		return nil, err
	}

	crlEntry, err := req.Storage.Get(ctx, crlPath)
	if err != nil {
		return nil, err
	}

	var certificate []byte
	if crlEntry != nil && len(crlEntry.Value) > 0 {
		certificate = []byte(crlEntry.Value)
	}

	var contentType string
	if strings.HasSuffix(req.Path, "/der") {
		contentType = "application/pkix-crl"
	} else if strings.HasSuffix(req.Path, "/pem") {
		contentType = "application/x-pem-file"
	}

	if !strings.HasSuffix(req.Path, "/der") {
		// Rather return an empty response rather than an empty PEM blob.
		// We build this PEM block for both the JSON and PEM endpoints.
		if len(certificate) > 0 {
			pemBlock := pem.Block{
				Type:  "X509 CRL",
				Bytes: certificate,
			}

			certificate = pem.EncodeToMemory(&pemBlock)
		}
	}

	statusCode := 200
	if len(certificate) == 0 {
		statusCode = 204
	}

	if strings.HasSuffix(req.Path, "/der") || strings.HasSuffix(req.Path, "/pem") {
		return &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPContentType: contentType,
				logical.HTTPRawBody:     certificate,
				logical.HTTPStatusCode:  statusCode,
			},
		}, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"crl": string(certificate),
		},
	}, nil
}

const (
	pathGetIssuerCRLHelpSyn  = `Fetch an issuer's Certificate Revocation Log (CRL).`
	pathGetIssuerCRLHelpDesc = `
This allows fetching the specified issuer's CRL. Note that this is different
than the legacy path (/crl and /certs/crl) in that this is per-issuer and not
just the default issuer's CRL.

Two issuers will have the same CRL if they have the same key material and if
they have the same Subject value.

:ref can be either the literal value "default", in which case /config/issuers
will be consulted for the present default issuer, an identifier of an issuer,
or its assigned name value.

 - /issuer/:ref/crl is JSON encoded and contains a PEM CRL,
 - /issuer/:ref/crl/pem contains the PEM-encoded CRL,
 - /issuer/:ref/crl/DER contains the raw DER-encoded (binary) CRL.
`
)
