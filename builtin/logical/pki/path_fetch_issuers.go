package pki

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

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

	sc := b.makeStorageContext(ctx, req.Storage)
	entries, err := sc.listIssuers()
	if err != nil {
		return nil, err
	}

	config, err := sc.getIssuersConfig()
	if err != nil {
		return nil, err
	}

	// For each issuer, we need not only the identifier (as returned by
	// listIssuers), but also the name of the issuer. This means we have to
	// fetch the actual issuer object as well.
	for _, identifier := range entries {
		issuer, err := sc.fetchIssuerById(identifier)
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
this issuer; valid values are "read-only", "issuing-certificates",
"crl-signing", and "ocsp-signing". Multiple values may be specified. Read-only
is implicit and always set.`,
		Default: []string{"read-only", "issuing-certificates", "crl-signing", "ocsp-signing"},
	}
	fields["revocation_signature_algorithm"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Which x509.SignatureAlgorithm name to use for
signing CRLs. This parameter allows differentiation between PKCS#1v1.5
and PSS keys and choice of signature hash algorithm. The default (empty
string) value is for Go to select the signature algorithm. This can fail
if the underlying key does not support the requested signature algorithm,
which may not be known at modification time (such as with PKCS#11 managed
RSA keys).`,
		Default: "",
	}
	fields["issuing_certificates"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `Comma-separated list of URLs to be used
for the issuing certificate attribute. See also RFC 5280 Section 4.2.2.1.`,
	}
	fields["crl_distribution_points"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `Comma-separated list of URLs to be used
for the CRL distribution points attribute. See also RFC 5280 Section 4.2.1.13.`,
	}
	fields["ocsp_servers"] = &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: `Comma-separated list of URLs to be used
for the OCSP servers attribute. See also RFC 5280 Section 4.2.2.1.`,
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
			logical.PatchOperation: &framework.PathOperation{
				Callback: b.pathPatchIssuer,
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

	return respondReadIssuer(issuer)
}

func respondReadIssuer(issuer *issuerEntry) (*logical.Response, error) {
	var respManualChain []string
	for _, entity := range issuer.ManualChain {
		respManualChain = append(respManualChain, string(entity))
	}

	revSigAlgStr, present := certutil.InvSignatureAlgorithmNames[issuer.RevocationSigAlg]
	if !present {
		revSigAlgStr = issuer.RevocationSigAlg.String()
		if issuer.RevocationSigAlg == x509.UnknownSignatureAlgorithm {
			revSigAlgStr = ""
		}
	}

	data := map[string]interface{}{
		"issuer_id":                      issuer.ID,
		"issuer_name":                    issuer.Name,
		"key_id":                         issuer.KeyID,
		"certificate":                    issuer.Certificate,
		"manual_chain":                   respManualChain,
		"ca_chain":                       issuer.CAChain,
		"leaf_not_after_behavior":        issuer.LeafNotAfterBehavior.String(),
		"usage":                          issuer.Usage.Names(),
		"revocation_signature_algorithm": revSigAlgStr,
		"revoked":                        issuer.Revoked,
		"issuing_certificates":           []string{},
		"crl_distribution_points":        []string{},
		"ocsp_servers":                   []string{},
	}

	if issuer.Revoked {
		data["revocation_time"] = issuer.RevocationTime
		data["revocation_time_rfc3339"] = issuer.RevocationTimeUTC.Format(time.RFC3339Nano)
	}

	if issuer.AIAURIs != nil {
		data["issuing_certificates"] = issuer.AIAURIs.IssuingCertificates
		data["crl_distribution_points"] = issuer.AIAURIs.CRLDistributionPoints
		data["ocsp_servers"] = issuer.AIAURIs.OCSPServers
	}

	response := &logical.Response{
		Data: data,
	}

	if issuer.RevocationSigAlg == x509.SHA256WithRSAPSS || issuer.RevocationSigAlg == x509.SHA384WithRSAPSS || issuer.RevocationSigAlg == x509.SHA512WithRSAPSS {
		response.AddWarning("Issuer uses a PSS Revocation Signature Algorithm. This algorithm will be downgraded to PKCS#1v1.5 signature scheme on OCSP responses, due to limitations in the OCSP library.")
	}

	return response, nil
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

	newName, err := getIssuerName(sc, data)
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
	if len(newName) > 0 && !nameMatcher.MatchString(newName) {
		return logical.ErrorResponse("new key name outside of valid character limits"), nil
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

	// Revocation signature algorithm changes
	revSigAlgStr := data.Get("revocation_signature_algorithm").(string)
	revSigAlg, present := certutil.SignatureAlgorithmNames[strings.ToLower(revSigAlgStr)]
	if !present && revSigAlgStr != "" {
		var knownAlgos []string
		for algoName := range certutil.SignatureAlgorithmNames {
			knownAlgos = append(knownAlgos, algoName)
		}

		return logical.ErrorResponse(fmt.Sprintf("Unknown signature algorithm value: %v - valid values are %v", revSigAlg, strings.Join(knownAlgos, ", "))), nil
	} else if revSigAlgStr == "" {
		revSigAlg = x509.UnknownSignatureAlgorithm
	}
	if err := issuer.CanMaybeSignWithAlgo(revSigAlg); err != nil {
		return nil, err
	}

	// AIA access changes
	issuerCertificates := data.Get("issuing_certificates").([]string)
	if badURL := validateURLs(issuerCertificates); badURL != "" {
		return logical.ErrorResponse(fmt.Sprintf("invalid URL found in Authority Information Access (AIA) parameter issuing_certificates: %s", badURL)), nil
	}
	crlDistributionPoints := data.Get("crl_distribution_points").([]string)
	if badURL := validateURLs(crlDistributionPoints); badURL != "" {
		return logical.ErrorResponse(fmt.Sprintf("invalid URL found in Authority Information Access (AIA) parameter crl_distribution_points: %s", badURL)), nil
	}
	ocspServers := data.Get("ocsp_servers").([]string)
	if badURL := validateURLs(ocspServers); badURL != "" {
		return logical.ErrorResponse(fmt.Sprintf("invalid URL found in Authority Information Access (AIA) parameter ocsp_servers: %s", badURL)), nil
	}

	modified := false

	var oldName string
	if newName != issuer.Name {
		oldName = issuer.Name
		issuer.Name = newName
		issuer.LastModified = time.Now().UTC()
		// See note in updateDefaultIssuerId about why this is necessary.
		b.crlBuilder.invalidateCRLBuildTime()
		b.crlBuilder.flushCRLBuildTimeInvalidation(sc)
		modified = true
	}

	if newLeafBehavior != issuer.LeafNotAfterBehavior {
		issuer.LeafNotAfterBehavior = newLeafBehavior
		modified = true
	}

	if newUsage != issuer.Usage {
		if issuer.Revoked && newUsage.HasUsage(IssuanceUsage) {
			// Forbid allowing cert signing on its usage.
			return logical.ErrorResponse("This issuer was revoked; unable to modify its usage to include certificate signing again. Reissue this certificate (preferably with a new key) and modify that entry instead."), nil
		}

		// Ensure we deny adding CRL usage if the bits are missing from the
		// cert itself.
		cert, err := issuer.GetCertificate()
		if err != nil {
			return nil, fmt.Errorf("unable to parse issuer's certificate: %v", err)
		}
		if (cert.KeyUsage&x509.KeyUsageCRLSign) == 0 && newUsage.HasUsage(CRLSigningUsage) {
			return logical.ErrorResponse("This issuer's underlying certificate lacks the CRLSign KeyUsage value; unable to set CRLSigningUsage on this issuer as a result."), nil
		}

		issuer.Usage = newUsage
		modified = true
	}

	if revSigAlg != issuer.RevocationSigAlg {
		issuer.RevocationSigAlg = revSigAlg
		modified = true
	}

	if issuer.AIAURIs == nil && (len(issuerCertificates) > 0 || len(crlDistributionPoints) > 0 || len(ocspServers) > 0) {
		issuer.AIAURIs = &certutil.URLEntries{}
	}
	if issuer.AIAURIs != nil {
		// Associative mapping from data source to destination on the
		// backing issuer object.
		type aiaPair struct {
			Source *[]string
			Dest   *[]string
		}
		pairs := []aiaPair{
			{
				Source: &issuerCertificates,
				Dest:   &issuer.AIAURIs.IssuingCertificates,
			},
			{
				Source: &crlDistributionPoints,
				Dest:   &issuer.AIAURIs.CRLDistributionPoints,
			},
			{
				Source: &ocspServers,
				Dest:   &issuer.AIAURIs.OCSPServers,
			},
		}

		// For each pair, if it is different on the object, update it.
		for _, pair := range pairs {
			if isStringArrayDifferent(*pair.Source, *pair.Dest) {
				*pair.Dest = *pair.Source
				modified = true
			}
		}

		// If no AIA URLs exist on the issuer, set the AIA URLs entry to nil
		// to ease usage later.
		if len(issuer.AIAURIs.IssuingCertificates) == 0 && len(issuer.AIAURIs.CRLDistributionPoints) == 0 && len(issuer.AIAURIs.OCSPServers) == 0 {
			issuer.AIAURIs = nil
		}
	}

	// Updating the chain should be the last modification as there's a chance
	// it'll write it out to disk for us. We'd hate to then modify the issuer
	// again and write it a second time.
	var updateChain bool
	var constructedChain []issuerID
	for index, newPathRef := range newPath {
		// Allow self for the first entry.
		if index == 0 && newPathRef == "self" {
			newPathRef = string(ref)
		}

		resolvedId, err := sc.resolveIssuerReference(newPathRef)
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
		err := sc.rebuildIssuersChains(issuer)
		if err != nil {
			return nil, err
		}
	}

	if modified {
		err := sc.writeIssuer(issuer)
		if err != nil {
			return nil, err
		}
	}

	response, err := respondReadIssuer(issuer)
	if newName != oldName {
		addWarningOnDereferencing(sc, oldName, response)
	}

	return response, err
}

func (b *backend) pathPatchIssuer(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not patch issuer until migration has completed"), nil
	}

	// First we fetch the issuer
	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

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

	// Now We are Looking at What (Might) Have Changed
	modified := false

	// Name Changes First
	_, ok := data.GetOk("issuer_name") // Don't check for conflicts if we aren't updating the name
	var oldName string
	var newName string
	if ok {
		newName, err = getIssuerName(sc, data)
		if err != nil && err != errIssuerNameInUse && err != errIssuerNameIsEmpty {
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
		if len(newName) > 0 && !nameMatcher.MatchString(newName) {
			return logical.ErrorResponse("new key name outside of valid character limits"), nil
		}
		if newName != issuer.Name {
			oldName = issuer.Name
			issuer.Name = newName
			issuer.LastModified = time.Now().UTC()
			// See note in updateDefaultIssuerId about why this is necessary.
			b.crlBuilder.invalidateCRLBuildTime()
			b.crlBuilder.flushCRLBuildTimeInvalidation(sc)
			modified = true
		}
	}

	// Leaf Not After Changes
	rawLeafBehaviorData, ok := data.GetOk("leaf_not_after_behaivor")
	if ok {
		rawLeafBehavior := rawLeafBehaviorData.(string)
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
		if newLeafBehavior != issuer.LeafNotAfterBehavior {
			issuer.LeafNotAfterBehavior = newLeafBehavior
			modified = true
		}
	}

	// Usage Changes
	rawUsageData, ok := data.GetOk("usage")
	if ok {
		rawUsage := rawUsageData.([]string)
		newUsage, err := NewIssuerUsageFromNames(rawUsage)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Unable to parse specified usages: %v - valid values are %v", rawUsage, AllIssuerUsages.Names())), nil
		}
		if newUsage != issuer.Usage {
			if issuer.Revoked && newUsage.HasUsage(IssuanceUsage) {
				// Forbid allowing cert signing on its usage.
				return logical.ErrorResponse("This issuer was revoked; unable to modify its usage to include certificate signing again. Reissue this certificate (preferably with a new key) and modify that entry instead."), nil
			}

			cert, err := issuer.GetCertificate()
			if err != nil {
				return nil, fmt.Errorf("unable to parse issuer's certificate: %v", err)
			}
			if (cert.KeyUsage&x509.KeyUsageCRLSign) == 0 && newUsage.HasUsage(CRLSigningUsage) {
				return logical.ErrorResponse("This issuer's underlying certificate lacks the CRLSign KeyUsage value; unable to set CRLSigningUsage on this issuer as a result."), nil
			}

			issuer.Usage = newUsage
			modified = true
		}
	}

	// Revocation signature algorithm changes
	rawRevSigAlg, ok := data.GetOk("revocation_signature_algorithm")
	if ok {
		revSigAlgStr := rawRevSigAlg.(string)
		revSigAlg, present := certutil.SignatureAlgorithmNames[strings.ToLower(revSigAlgStr)]
		if !present && revSigAlgStr != "" {
			var knownAlgos []string
			for algoName := range certutil.SignatureAlgorithmNames {
				knownAlgos = append(knownAlgos, algoName)
			}

			return logical.ErrorResponse(fmt.Sprintf("Unknown signature algorithm value: %v - valid values are %v", revSigAlg, strings.Join(knownAlgos, ", "))), nil
		} else if revSigAlgStr == "" {
			revSigAlg = x509.UnknownSignatureAlgorithm
		}

		if err := issuer.CanMaybeSignWithAlgo(revSigAlg); err != nil {
			return nil, err
		}

		if revSigAlg != issuer.RevocationSigAlg {
			issuer.RevocationSigAlg = revSigAlg
			modified = true
		}
	}

	// AIA access changes.
	if issuer.AIAURIs == nil {
		issuer.AIAURIs = &certutil.URLEntries{}
	}

	// Associative mapping from data source to destination on the
	// backing issuer object. For PATCH requests, we use the source
	// data parameter as we still need to validate them and process
	// it into a string list.
	type aiaPair struct {
		Source string
		Dest   *[]string
	}
	pairs := []aiaPair{
		{
			Source: "issuing_certificates",
			Dest:   &issuer.AIAURIs.IssuingCertificates,
		},
		{
			Source: "crl_distribution_points",
			Dest:   &issuer.AIAURIs.CRLDistributionPoints,
		},
		{
			Source: "ocsp_servers",
			Dest:   &issuer.AIAURIs.OCSPServers,
		},
	}

	// For each pair, if it is different on the object, update it.
	for _, pair := range pairs {
		rawURLsValue, ok := data.GetOk(pair.Source)
		if ok {
			urlsValue := rawURLsValue.([]string)
			if badURL := validateURLs(urlsValue); badURL != "" {
				return logical.ErrorResponse(fmt.Sprintf("invalid URL found in Authority Information Access (AIA) parameter %v: %s", pair.Source, badURL)), nil
			}

			if isStringArrayDifferent(urlsValue, *pair.Dest) {
				modified = true
				*pair.Dest = urlsValue
			}
		}
	}

	// If no AIA URLs exist on the issuer, set the AIA URLs entry to nil to
	// ease usage later.
	if len(issuer.AIAURIs.IssuingCertificates) == 0 && len(issuer.AIAURIs.CRLDistributionPoints) == 0 && len(issuer.AIAURIs.OCSPServers) == 0 {
		issuer.AIAURIs = nil
	}

	// Manual Chain Changes
	newPathData, ok := data.GetOk("manual_chain")
	if ok {
		newPath := newPathData.([]string)
		var updateChain bool
		var constructedChain []issuerID
		for index, newPathRef := range newPath {
			// Allow self for the first entry.
			if index == 0 && newPathRef == "self" {
				newPathRef = string(ref)
			}

			resolvedId, err := sc.resolveIssuerReference(newPathRef)
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
			err := sc.rebuildIssuersChains(issuer)
			if err != nil {
				return nil, err
			}
		}
	}

	if modified {
		err := sc.writeIssuer(issuer)
		if err != nil {
			return nil, err
		}
	}

	response, err := respondReadIssuer(issuer)
	if newName != oldName {
		addWarningOnDereferencing(sc, oldName, response)
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

	var contentType string
	var certificate []byte

	response := &logical.Response{}
	ret, err := sendNotModifiedResponseIfNecessary(&IfModifiedSinceHelper{req: req, reqType: ifModifiedCA, issuerRef: ref}, sc, response)
	if err != nil {
		return nil, err
	}
	if ret {
		return response, nil
	}

	certificate = []byte(issuer.Certificate)

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

	sc := b.makeStorageContext(ctx, req.Storage)
	ref, err := sc.resolveIssuerReference(issuerName)
	if err != nil {
		// Return as if we deleted it if we fail to lookup the issuer.
		if ref == IssuerRefNotFound {
			return &logical.Response{}, nil
		}
		return nil, err
	}

	response := &logical.Response{}

	issuer, err := sc.fetchIssuerById(ref)
	if err != nil {
		return nil, err
	}
	if issuer.Name != "" {
		addWarningOnDereferencing(sc, issuer.Name, response)
	}
	addWarningOnDereferencing(sc, string(issuer.ID), response)

	wasDefault, err := sc.deleteIssuer(ref)
	if err != nil {
		return nil, err
	}
	if wasDefault {
		response.AddWarning(fmt.Sprintf("Deleted issuer %v (via issuer_ref %v); this was configured as the default issuer. Operations without an explicit issuer will not work until a new default is configured.", ref, issuerName))
		addWarningOnDereferencing(sc, defaultRef, response)
	}

	// Since we've deleted an issuer, the chains might've changed. Call the
	// rebuild code. We shouldn't technically err (as the issuer was deleted
	// successfully), but log a warning (and to the response) if this fails.
	if err := sc.rebuildIssuersChains(nil); err != nil {
		msg := fmt.Sprintf("Failed to rebuild remaining issuers' chains: %v", err)
		b.Logger().Error(msg)
		response.AddWarning(msg)
	}

	return response, nil
}

func addWarningOnDereferencing(sc *storageContext, name string, resp *logical.Response) {
	timeout, inUseBy, err := sc.checkForRolesReferencing(name)
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
	pattern := "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/crl(/pem|/der|/delta(/pem|/der)?)?"
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

	var certificate []byte
	var contentType string

	sc := b.makeStorageContext(ctx, req.Storage)
	response := &logical.Response{}
	var crlType ifModifiedReqType = ifModifiedCRL
	if strings.Contains(req.Path, "delta") {
		crlType = ifModifiedDeltaCRL
	}
	ret, err := sendNotModifiedResponseIfNecessary(&IfModifiedSinceHelper{req: req, reqType: crlType}, sc, response)
	if err != nil {
		return nil, err
	}
	if ret {
		return response, nil
	}
	crlPath, err := sc.resolveIssuerCRLPath(issuerName)
	if err != nil {
		return nil, err
	}

	if strings.Contains(req.Path, "delta") {
		crlPath += deltaCRLPathSuffix
	}

	crlEntry, err := req.Storage.Get(ctx, crlPath)
	if err != nil {
		return nil, err
	}

	if crlEntry != nil && len(crlEntry.Value) > 0 {
		certificate = []byte(crlEntry.Value)
	}

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
