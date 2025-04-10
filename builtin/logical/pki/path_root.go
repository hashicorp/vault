// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ed25519"
)

const intCaTruncatationWarning = "the signed intermediary CA certificate's notAfter was truncated to the issuer's notAfter"

func pathGenerateRoot(b *backend) *framework.Path {
	pattern := "root/generate/" + framework.GenericNameRegex("exported")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKI,
		OperationVerb:   "generate",
		OperationSuffix: "root",
	}

	return buildPathGenerateRoot(b, pattern, displayAttrs)
}

func pathDeleteRoot(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "root",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationSuffix: "root",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathCADeleteRoot,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
					}},
				},
				// Read more about why these flags are set in backend.go
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathDeleteRootHelpSyn,
		HelpDescription: pathDeleteRootHelpDesc,
	}

	return ret
}

func (b *backend) pathCADeleteRoot(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	sc := b.makeStorageContext(ctx, req.Storage)
	if !b.UseLegacyBundleCaStorage() {
		issuers, err := sc.listIssuers()
		if err != nil {
			return nil, err
		}

		keys, err := sc.listKeys()
		if err != nil {
			return nil, err
		}

		// Delete all issuers and keys. Ignore deleting the default since we're
		// explicitly deleting everything.
		for _, issuer := range issuers {
			if _, err = sc.deleteIssuer(issuer); err != nil {
				return nil, err
			}
		}
		for _, key := range keys {
			if _, err = sc.deleteKey(key); err != nil {
				return nil, err
			}
		}
	}

	// Delete legacy CA bundle and its backup, if any.
	if err := req.Storage.Delete(ctx, legacyCertBundlePath); err != nil {
		return nil, err
	}

	if err := req.Storage.Delete(ctx, legacyCertBundleBackupPath); err != nil {
		return nil, err
	}

	// Delete legacy CRL bundle.
	if err := req.Storage.Delete(ctx, legacyCRLPath); err != nil {
		return nil, err
	}

	// Return a warning about preferring to delete issuers and keys
	// explicitly versus deleting everything.
	resp := &logical.Response{}
	resp.AddWarning("DELETE /root deletes all keys and issuers; prefer the new DELETE /key/:key_ref and DELETE /issuer/:issuer_ref for finer granularity, unless removal of all keys and issuers is desired.")
	return resp, nil
}

func (b *backend) pathCAGenerateRoot(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	var err error

	if b.UseLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not create root CA until migration has completed"), nil
	}

	sc := b.makeStorageContext(ctx, req.Storage)

	exported, format, role, errorResp := getGenerationParams(sc, data)
	if errorResp != nil {
		return errorResp, nil
	}

	maxPathLengthIface, ok := data.GetOk("max_path_length")
	if ok {
		maxPathLength := maxPathLengthIface.(int)
		role.MaxPathLength = &maxPathLength
	}

	issuerName, err := getIssuerName(sc, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	// Handle the aliased path specifying the new issuer name as "next", but
	// only do it if its not in use.
	if strings.HasPrefix(req.Path, "root/rotate/") && len(issuerName) == 0 {
		// err is nil when the issuer name is in use.
		_, err = sc.resolveIssuerReference("next")
		if err != nil {
			issuerName = "next"
		}
	}

	keyName, err := getKeyName(sc, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	input := &inputBundle{
		req:     req,
		apiData: data,
		role:    role,
	}
	parsedBundle, warnings, err := generateCert(sc, input, nil, true, b.Backend.GetRandomReader())
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		default:
			return nil, err
		}
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"expiration":    int64(parsedBundle.Certificate.NotAfter.Unix()),
			"serial_number": cb.SerialNumber,
		},
	}

	if keyUsages, ok := data.GetOk("key_usage"); ok {
		err = validateCaKeyUsages(keyUsages.([]string))
		if err != nil {
			resp.AddWarning(fmt.Sprintf("Invalid key usage will be ignored: %v", err.Error()))
		}
	}

	if len(parsedBundle.Certificate.RawSubject) <= 2 {
		// Strictly a subject is a SEQUENCE of SETs of SEQUENCES.
		//
		// The outer SEQUENCE is preserved, having byte value 30 00.
		//
		// Because of the tag and the length encoding each taking up
		// at least one byte, it is impossible to have a non-empty
		// subject in two or fewer bytes. We're also not here to validate
		// our certificate's ASN.1 content, so let's just assume it holds
		// and move on.
		resp.AddWarning("This issuer certificate was generated without a Subject; this makes it likely that issuing leaf certs with this certificate will cause TLS validation libraries to reject this certificate.")
	}

	if len(parsedBundle.Certificate.OCSPServer) == 0 && len(parsedBundle.Certificate.IssuingCertificateURL) == 0 && len(parsedBundle.Certificate.CRLDistributionPoints) == 0 {
		// If the operator hasn't configured any of the URLs prior to
		// generating this issuer, we should add a warning to the response,
		// informing them they might want to do so prior to issuing leaves.
		resp.AddWarning("This mount hasn't configured any authority information access (AIA) fields; this may make it harder for systems to find missing certificates in the chain or to validate revocation status of certificates. Consider updating /config/urls or the newly generated issuer with this information.")
	}

	switch format {
	case "pem":
		resp.Data["certificate"] = cb.Certificate
		resp.Data["issuing_ca"] = cb.Certificate
		if exported {
			resp.Data["private_key"] = cb.PrivateKey
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}

	case "pem_bundle":
		resp.Data["issuing_ca"] = cb.Certificate

		if exported {
			resp.Data["private_key"] = cb.PrivateKey
			resp.Data["private_key_type"] = cb.PrivateKeyType
			resp.Data["certificate"] = fmt.Sprintf("%s\n%s", cb.PrivateKey, cb.Certificate)
		} else {
			resp.Data["certificate"] = cb.Certificate
		}

	case "der":
		resp.Data["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		resp.Data["issuing_ca"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		if exported {
			resp.Data["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}
	default:
		return nil, fmt.Errorf("unsupported format argument: %s", format)
	}

	if data.Get("private_key_format").(string) == "pkcs8" {
		err = convertRespToPKCS8(resp)
		if err != nil {
			return nil, err
		}
	}

	// Store it as the CA bundle
	myIssuer, myKey, err := sc.writeCaBundle(cb, issuerName, keyName)
	if err != nil {
		return nil, err
	}
	resp.Data["issuer_id"] = myIssuer.ID
	resp.Data["issuer_name"] = myIssuer.Name
	resp.Data["key_id"] = myKey.ID
	resp.Data["key_name"] = myKey.Name

	// The one time that it is safe (and good) to copy the
	// SignatureAlgorithm field off the certificate (for the purposes of
	// detecting PSS support) is when we've freshly generated it AND it
	// is a root (exactly this endpoint).
	//
	// For intermediates, this doesn't hold (not this endpoint) as that
	// reflects the parent key's preferences. For imports, this doesn't
	// hold as the old system might've allowed other signature types that
	// the new system (whether Vault or a managed key) doesn't.
	//
	// Previously we did this conditionally on whether or not PSS was in
	// use. This is insufficient as some cloud KMS providers (namely, GCP)
	// restrict the key to a single signature algorithm! So e.g., a RSA 3072
	// key MUST use SHA-384 as the hash algorithm. Thus we pull in the
	// RevocationSigAlg unconditionally on roots now.
	myIssuer.RevocationSigAlg = parsedBundle.Certificate.SignatureAlgorithm
	if err := sc.writeIssuer(myIssuer); err != nil {
		return nil, fmt.Errorf("unable to store PSS-updated issuer: %w", err)
	}

	// Also store it as just the certificate identified by serial number, so it
	// can be revoked
	err = issuing.StoreCertificate(ctx, req.Storage, b.GetCertificateCounter(), parsedBundle)
	if err != nil {
		return nil, err
	}

	// Build a fresh CRL
	warnings, err = b.CrlBuilder().Rebuild(sc, true)
	if err != nil {
		return nil, err
	}
	for index, warning := range warnings {
		resp.AddWarning(fmt.Sprintf("Warning %d during CRL rebuild: %v", index+1, warning))
	}

	if parsedBundle.Certificate.MaxPathLen == 0 {
		resp.AddWarning("Max path length of the generated certificate is zero. This certificate cannot be used to issue intermediate CA certificates.")
	}

	// Check whether we need to update our default issuer configuration.
	config, err := sc.getIssuersConfig()
	if err != nil {
		resp.AddWarning("Unable to fetch default issuers configuration to update default issuer if necessary: " + err.Error())
	} else if config.DefaultFollowsLatestIssuer {
		if err := sc.updateDefaultIssuerId(myIssuer.ID); err != nil {
			resp.AddWarning("Unable to update this new root as the default issuer: " + err.Error())
		}
	}

	resp = addWarnings(resp, warnings)

	return resp, nil
}

func (b *backend) pathIssuerSignIntermediate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	issuerName := GetIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	format := getFormat(data)
	if format == "" {
		return logical.ErrorResponse(`The "format" path parameter must be "pem", "der" or "pem_bundle"`), nil
	}

	role := &issuing.RoleEntry{
		OU:                        data.Get("ou").([]string),
		Organization:              data.Get("organization").([]string),
		Country:                   data.Get("country").([]string),
		Locality:                  data.Get("locality").([]string),
		Province:                  data.Get("province").([]string),
		StreetAddress:             data.Get("street_address").([]string),
		PostalCode:                data.Get("postal_code").([]string),
		TTL:                       time.Duration(data.Get("ttl").(int)) * time.Second,
		AllowLocalhost:            true,
		AllowAnyName:              true,
		AllowIPSANs:               true,
		AllowWildcardCertificates: new(bool),
		EnforceHostnames:          false,
		KeyType:                   "any",
		SignatureBits:             data.Get("signature_bits").(int),
		UsePSS:                    data.Get("use_pss").(bool),
		AllowedOtherSANs:          []string{"*"},
		AllowedSerialNumbers:      []string{"*"},
		AllowedURISANs:            []string{"*"},
		NotAfter:                  data.Get("not_after").(string),
		NotBeforeDuration:         time.Duration(data.Get("not_before_duration").(int)) * time.Second,
		CNValidations:             []string{"disabled"},
		KeyUsage:                  data.Get("key_usage").([]string),
	}
	*role.AllowWildcardCertificates = true

	if cn := data.Get("common_name").(string); len(cn) == 0 {
		role.UseCSRCommonName = true
	}

	var caErr error
	sc := b.makeStorageContext(ctx, req.Storage)
	signingBundle, issuerId, caErr := sc.fetchCAInfoWithIssuer(issuerName, issuing.IssuanceUsage)
	if caErr != nil {
		switch caErr.(type) {
		case errutil.UserError:
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"could not fetch the CA certificate (was one set?): %s", caErr)}
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf(
				"error fetching CA certificate: %s", caErr)}
		}
	}

	warnAboutTruncate := false
	if enforceLeafNotAfter := data.Get("enforce_leaf_not_after_behavior").(bool); !enforceLeafNotAfter {
		// Since we are signing an intermediate, we will by default truncate the
		// signed intermediary in order to generate a valid intermediary chain. This
		// was changed in 1.17.x as the default prior was PermitNotAfterBehavior
		if signingBundle.LeafNotAfterBehavior != certutil.AlwaysEnforceErr {
			warnAboutTruncate = true
			signingBundle.LeafNotAfterBehavior = certutil.TruncateNotAfterBehavior
		}
	}

	useCSRValues := data.Get("use_csr_values").(bool)

	maxPathLengthIface, ok := data.GetOk("max_path_length")
	if ok {
		maxPathLength := maxPathLengthIface.(int)
		role.MaxPathLength = &maxPathLength
	}

	input := &inputBundle{
		req:     req,
		apiData: data,
		role:    role,
	}
	parsedBundle, warnings, err := signCert(b.System(), input, signingBundle, true, useCSRValues)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf(
				"error signing cert: %s", err)}
		}
	}

	if err := issuing.VerifyCertificate(sc.GetContext(), sc.GetStorage(), issuerId, parsedBundle); err != nil {
		return nil, fmt.Errorf("verification of parsed bundle failed: %w", err)
	}

	resp, err := signIntermediateResponse(signingBundle, parsedBundle, format, warnings)
	if err != nil {
		return nil, err
	}

	err = issuing.StoreCertificate(ctx, req.Storage, b.GetCertificateCounter(), parsedBundle)
	if err != nil {
		return nil, err
	}

	if warnAboutTruncate &&
		signingBundle.Certificate.NotAfter.Equal(parsedBundle.Certificate.NotAfter) {
		resp.AddWarning(intCaTruncatationWarning)
	}

	if keyUsages, ok := data.GetOk("key_usage"); ok {
		err = validateCaKeyUsages(keyUsages.([]string))
		if err != nil {
			resp.AddWarning(fmt.Sprintf("Invalid key usage: %v", err.Error()))
		}
	}

	return resp, nil
}

func signIntermediateResponse(signingBundle *certutil.CAInfoBundle, parsedBundle *certutil.ParsedCertBundle, format string, warnings []string) (*logical.Response, error) {
	signingCB, err := signingBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw signing bundle to cert bundle: %w", err)
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"expiration":    int64(parsedBundle.Certificate.NotAfter.Unix()),
			"serial_number": cb.SerialNumber,
		},
	}

	if signingBundle.Certificate.NotAfter.Before(parsedBundle.Certificate.NotAfter) {
		resp.AddWarning("The expiration time for the signed certificate is after the CA's expiration time. If the new certificate is not treated as a root, validation paths with the certificate past the issuing CA's expiration time will fail.")
	}

	if len(parsedBundle.Certificate.RawSubject) <= 2 {
		// Strictly a subject is a SEQUENCE of SETs of SEQUENCES.
		//
		// The outer SEQUENCE is preserved, having byte value 30 00.
		//
		// Because of the tag and the length encoding each taking up
		// at least one byte, it is impossible to have a non-empty
		// subject in two or fewer bytes. We're also not here to validate
		// our certificate's ASN.1 content, so let's just assume it holds
		// and move on.
		resp.AddWarning("This issuer certificate was generated without a Subject; this makes it likely that issuing leaf certs with this certificate will cause TLS validation libraries to reject this certificate.")
	}

	if len(parsedBundle.Certificate.OCSPServer) == 0 && len(parsedBundle.Certificate.IssuingCertificateURL) == 0 && len(parsedBundle.Certificate.CRLDistributionPoints) == 0 {
		// If the operator hasn't configured any of the URLs prior to
		// generating this issuer, we should add a warning to the response,
		// informing them they might want to do so prior to issuing leaves.
		resp.AddWarning("This mount hasn't configured any authority information access (AIA) fields; this may make it harder for systems to find missing certificates in the chain or to validate revocation status of certificates. Consider updating /config/urls or the newly generated issuer with this information.")
	}

	caChain := append([]string{cb.Certificate}, cb.CAChain...)

	switch format {
	case "pem":
		resp.Data["certificate"] = cb.Certificate
		resp.Data["issuing_ca"] = signingCB.Certificate
		resp.Data["ca_chain"] = caChain

	case "pem_bundle":
		resp.Data["certificate"] = cb.ToPEMBundle()
		resp.Data["issuing_ca"] = signingCB.Certificate
		resp.Data["ca_chain"] = caChain

	case "der":
		resp.Data["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		resp.Data["issuing_ca"] = base64.StdEncoding.EncodeToString(signingBundle.CertificateBytes)

		var derCaChain []string
		derCaChain = append(derCaChain, base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes))
		for _, caCert := range parsedBundle.CAChain {
			derCaChain = append(derCaChain, base64.StdEncoding.EncodeToString(caCert.Bytes))
		}
		resp.Data["ca_chain"] = derCaChain

	default:
		return nil, fmt.Errorf("unsupported format argument: %s", format)
	}

	if parsedBundle.Certificate.MaxPathLen == 0 {
		resp.AddWarning("Max path length of the signed certificate is zero. This certificate cannot be used to issue intermediate CA certificates.")
	}

	resp = addWarnings(resp, warnings)
	return resp, nil
}

func (b *backend) pathIssuerSignSelfIssued(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	issuerName := GetIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	certPem := data.Get("certificate").(string)
	certs, err := parsing.ParseCertificatesFromString(certPem)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error parsing certificate: %s", err)), nil
	}
	if len(certs) != 1 {
		return logical.ErrorResponse(fmt.Sprintf("%d certificates found in PEM file, expected 1", len(certs))), nil
	}

	cert := certs[0]
	if !cert.IsCA {
		return logical.ErrorResponse("given certificate is not a CA certificate"), nil
	}
	if !reflect.DeepEqual(cert.Issuer, cert.Subject) {
		return logical.ErrorResponse("given certificate is not self-issued"), nil
	}

	sc := b.makeStorageContext(ctx, req.Storage)
	signingBundle, caErr := sc.fetchCAInfo(issuerName, issuing.IssuanceUsage)
	if caErr != nil {
		switch caErr.(type) {
		case errutil.UserError:
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"could not fetch the CA certificate (was one set?): %s", caErr)}
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf("error fetching CA certificate: %s", caErr)}
		}
	}

	signingCB, err := signingBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw signing bundle to cert bundle: %w", err)
	}

	urls := &certutil.URLEntries{}
	if signingBundle.URLs != nil {
		urls = signingBundle.URLs
	}
	cert.IssuingCertificateURL = urls.IssuingCertificates
	cert.CRLDistributionPoints = urls.CRLDistributionPoints
	cert.OCSPServer = urls.OCSPServers

	// If the requested signature algorithm isn't the same as the signing certificate, and
	// the user has requested a cross-algorithm signature, reset the template's signing algorithm
	// to that of the signing key
	signingPubType, signingAlgorithm, err := publicKeyType(signingBundle.Certificate.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("error determining signing certificate algorithm type: %e", err)
	}
	certPubType, _, err := publicKeyType(cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("error determining template algorithm type: %e", err)
	}

	if signingPubType != certPubType {
		b, ok := data.GetOk("require_matching_certificate_algorithms")
		if !ok || !b.(bool) {
			cert.SignatureAlgorithm = signingAlgorithm
		} else {
			return nil, fmt.Errorf("signing certificate's public key algorithm (%s) does not match submitted certificate's (%s), and require_matching_certificate_algorithms is true",
				signingPubType.String(), certPubType.String())
		}
	}

	newCert, err := x509.CreateCertificate(rand.Reader, cert, signingBundle.Certificate, cert.PublicKey, signingBundle.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing self-issued certificate: %w", err)
	}
	if len(newCert) == 0 {
		return nil, fmt.Errorf("nil cert was created when signing self-issued certificate")
	}
	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCert,
	})

	return &logical.Response{
		Data: map[string]interface{}{
			"certificate": strings.TrimSpace(string(pemCert)),
			"issuing_ca":  signingCB.Certificate,
		},
	}, nil
}

// Adapted from similar code in https://github.com/golang/go/blob/4a4221e8187189adcc6463d2d96fe2e8da290132/src/crypto/x509/x509.go#L1342,
// may need to be updated in the future.
func publicKeyType(pub crypto.PublicKey) (pubType x509.PublicKeyAlgorithm, sigAlgo x509.SignatureAlgorithm, err error) {
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		pubType = x509.RSA
		sigAlgo = x509.SHA256WithRSA
	case *ecdsa.PublicKey:
		pubType = x509.ECDSA
		switch pub.Curve {
		case elliptic.P224(), elliptic.P256():
			sigAlgo = x509.ECDSAWithSHA256
		case elliptic.P384():
			sigAlgo = x509.ECDSAWithSHA384
		case elliptic.P521():
			sigAlgo = x509.ECDSAWithSHA512
		default:
			err = errors.New("x509: unknown elliptic curve")
		}
	case ed25519.PublicKey:
		pubType = x509.Ed25519
		sigAlgo = x509.PureEd25519
	default:
		err = errors.New("x509: only RSA, ECDSA and Ed25519 keys supported")
	}
	return
}

func validateCaKeyUsages(keyUsages []string) error {
	invalidKeyUsages := []string{}
	for _, usage := range keyUsages {
		cleanUsage := strings.ToLower(strings.TrimSpace(usage))
		switch cleanUsage {
		case "crlsign", "certsign", "digitalsignature":
		case "contentcommitment", "keyencipherment", "dataencipherment", "keyagreement", "encipheronly", "decipheronly":
			invalidKeyUsages = append(invalidKeyUsages, fmt.Sprintf("key usage %s is only valid for non-Ca certs", usage))
		default:
			invalidKeyUsages = append(invalidKeyUsages, fmt.Sprintf("unrecognized key usage %s", usage))
		}
	}
	if invalidKeyUsages != nil {
		return fmt.Errorf(strings.Join(invalidKeyUsages, "; "))
	}
	return nil
}

const pathGenerateRootHelpSyn = `
Generate a new CA certificate and private key used for signing.
`

const pathGenerateRootHelpDesc = `
See the API documentation for more information.
`

const pathDeleteRootHelpSyn = `
Deletes the root CA key to allow a new one to be generated.
`

const pathDeleteRootHelpDesc = `
See the API documentation for more information.
`
