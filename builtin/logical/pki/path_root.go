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
	"reflect"
	"strings"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/hashicorp/vault/sdk/helper/certutil"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerateRoot(b *backend) *framework.Path {
	return buildPathGenerateRoot(b, "root/generate/"+framework.GenericNameRegex("exported"))
}

func pathDeleteRoot(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "root",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathCADeleteRoot,
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

	if !b.useLegacyBundleCaStorage() {
		issuers, err := listIssuers(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		keys, err := listKeys(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		// Delete all issuers and keys. Ignore deleting the default since we're
		// explicitly deleting everything.
		for _, issuer := range issuers {
			if _, err = deleteIssuer(ctx, req.Storage, issuer); err != nil {
				return nil, err
			}
		}
		for _, key := range keys {
			if _, err = deleteKey(ctx, req.Storage, key); err != nil {
				return nil, err
			}
		}
	}

	// Delete legacy CA bundle.
	if err := req.Storage.Delete(ctx, legacyCertBundlePath); err != nil {
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

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not create root CA until migration has completed"), nil
	}

	exported, format, role, errorResp := b.getGenerationParams(ctx, req.Storage, data)
	if errorResp != nil {
		return errorResp, nil
	}

	maxPathLengthIface, ok := data.GetOk("max_path_length")
	if ok {
		maxPathLength := maxPathLengthIface.(int)
		role.MaxPathLength = &maxPathLength
	}

	issuerName, err := getIssuerName(ctx, req.Storage, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	// Handle the aliased path specifying the new issuer name as "next", but
	// only do it if its not in use.
	if strings.HasPrefix(req.Path, "root/rotate/") && len(issuerName) == 0 {
		// err is nil when the issuer name is in use.
		_, err = resolveIssuerReference(ctx, req.Storage, "next")
		if err != nil {
			issuerName = "next"
		}
	}

	keyName, err := getKeyName(ctx, req.Storage, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	input := &inputBundle{
		req:     req,
		apiData: data,
		role:    role,
	}
	parsedBundle, err := generateCert(ctx, b, input, nil, true, b.Backend.GetRandomReader())
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
		// informing them they might want to do so and re-generate the issuer.
		resp.AddWarning("This mount hasn't configured any authority access information fields; this may make it harder for systems to find missing certificates in the chain or to validate revocation status of certificates. Consider updating /config/urls with this information.")
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
	myIssuer, myKey, err := writeCaBundle(ctx, b, req.Storage, cb, issuerName, keyName)
	if err != nil {
		return nil, err
	}
	resp.Data["issuer_id"] = myIssuer.ID
	resp.Data["key_id"] = myKey.ID

	// Also store it as just the certificate identified by serial number, so it
	// can be revoked
	err = req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   "certs/" + normalizeSerial(cb.SerialNumber),
		Value: parsedBundle.CertificateBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to store certificate locally: %w", err)
	}

	// Build a fresh CRL
	err = b.crlBuilder.rebuild(ctx, b, req, true)
	if err != nil {
		return nil, err
	}

	if parsedBundle.Certificate.MaxPathLen == 0 {
		resp.AddWarning("Max path length of the generated certificate is zero. This certificate cannot be used to issue intermediate CA certificates.")
	}

	return resp, nil
}

func (b *backend) pathIssuerSignIntermediate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	format := getFormat(data)
	if format == "" {
		return logical.ErrorResponse(
			`The "format" path parameter must be "pem" or "der"`,
		), nil
	}

	role := &roleEntry{
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
		AllowedOtherSANs:          []string{"*"},
		AllowedSerialNumbers:      []string{"*"},
		AllowedURISANs:            []string{"*"},
		NotAfter:                  data.Get("not_after").(string),
		NotBeforeDuration:         time.Duration(data.Get("not_before_duration").(int)) * time.Second,
	}
	*role.AllowWildcardCertificates = true

	if cn := data.Get("common_name").(string); len(cn) == 0 {
		role.UseCSRCommonName = true
	}

	var caErr error
	signingBundle, caErr := fetchCAInfo(ctx, b, req, issuerName, IssuanceUsage)
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

	// Since we are signing an intermediate, we explicitly want to override
	// the leaf NotAfterBehavior to permit issuing intermediates longer than
	// the life of this issuer.
	signingBundle.LeafNotAfterBehavior = certutil.PermitNotAfterBehavior

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
	parsedBundle, err := signCert(b, input, signingBundle, true, useCSRValues)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf(
				"error signing cert: %s", err)}
		}
	}

	if err := parsedBundle.Verify(); err != nil {
		return nil, fmt.Errorf("verification of parsed bundle failed: %w", err)
	}

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
		// informing them they might want to do so and re-generate the issuer.
		resp.AddWarning("This mount hasn't configured any authority access information fields; this may make it harder for systems to find missing certificates in the chain or to validate revocation status of certificates. Consider updating /config/urls with this information.")
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

	err = req.Storage.Put(ctx, &logical.StorageEntry{
		Key:   "certs/" + normalizeSerial(cb.SerialNumber),
		Value: parsedBundle.CertificateBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to store certificate locally: %w", err)
	}

	if parsedBundle.Certificate.MaxPathLen == 0 {
		resp.AddWarning("Max path length of the signed certificate is zero. This certificate cannot be used to issue intermediate CA certificates.")
	}

	return resp, nil
}

func (b *backend) pathIssuerSignSelfIssued(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	issuerName := getIssuerRef(data)
	if len(issuerName) == 0 {
		return logical.ErrorResponse("missing issuer reference"), nil
	}

	certPem := data.Get("certificate").(string)
	block, _ := pem.Decode([]byte(certPem))
	if block == nil || len(block.Bytes) == 0 {
		return logical.ErrorResponse("certificate could not be PEM-decoded"), nil
	}
	certs, err := x509.ParseCertificates(block.Bytes)
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

	var caErr error
	signingBundle, caErr := fetchCAInfo(ctx, b, req, issuerName, IssuanceUsage)
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
