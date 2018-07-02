package pki

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathIssue(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "issue/" + framework.GenericNameRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathIssue,
		},

		HelpSynopsis:    pathIssueHelpSyn,
		HelpDescription: pathIssueHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})
	return ret
}

func pathSign(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSign,
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Default:     "",
		Description: `PEM-format CSR to be signed.`,
	}

	return ret
}

func pathSignVerbatim(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "sign-verbatim" + framework.OptionalParamRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSignVerbatim,
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "",
		Description: `PEM-format CSR to be signed. Values will be
taken verbatim from the CSR, except for
basic constraints.`,
	}

	ret.Fields["key_usage"] = &framework.FieldSchema{
		Type:    framework.TypeCommaStringSlice,
		Default: []string{"DigitalSignature", "KeyAgreement", "KeyEncipherment"},
		Description: `A comma-separated string or list of key usages (not extended
key usages). Valid values can be found at
https://golang.org/pkg/crypto/x509/#KeyUsage
-- simply drop the "KeyUsage" part of the name.
To remove all key usages from being set, set
this value to an empty list.`,
	}

	ret.Fields["ext_key_usage"] = &framework.FieldSchema{
		Type:    framework.TypeCommaStringSlice,
		Default: []string{},
		Description: `A comma-separated string or list of extended key usages. Valid values can be found at
https://golang.org/pkg/crypto/x509/#ExtKeyUsage
-- simply drop the "ExtKeyUsage" part of the name.
To remove all key usages from being set, set
this value to an empty list.`,
	}

	ret.Fields["ext_key_usage_oids"] = &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: `A comma-separated string or list of extended key usage oids.`,
	}

	return ret
}

// pathIssue issues a certificate and private key from given parameters,
// subject to role restrictions
func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
	}

	if role.KeyType == "any" {
		return logical.ErrorResponse("role key type \"any\" not allowed for issuing certificates, only signing"), nil
	}

	return b.pathIssueSignCert(ctx, req, data, role, false, false)
}

// pathSign issues a certificate from a submitted CSR, subject to role
// restrictions
func (b *backend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
	}

	return b.pathIssueSignCert(ctx, req, data, role, true, false)
}

// pathSignVerbatim issues a certificate from a submitted CSR, *not* subject to
// role restrictions
func (b *backend) pathSignVerbatim(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := data.Get("role").(string)

	// Get the role if one was specified
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	entry := &roleEntry{
		TTL:                  b.System().DefaultLeaseTTL(),
		MaxTTL:               b.System().MaxLeaseTTL(),
		AllowLocalhost:       true,
		AllowAnyName:         true,
		AllowIPSANs:          true,
		EnforceHostnames:     false,
		KeyType:              "any",
		UseCSRCommonName:     true,
		UseCSRSANs:           true,
		AllowedURISANs:       []string{"*"},
		AllowedSerialNumbers: []string{"*"},
		GenerateLease:        new(bool),
		KeyUsage:             data.Get("key_usage").([]string),
		ExtKeyUsage:          data.Get("ext_key_usage").([]string),
		ExtKeyUsageOIDs:      data.Get("ext_key_usage_oids").([]string),
	}

	*entry.GenerateLease = false

	if role != nil {
		if role.TTL > 0 {
			entry.TTL = role.TTL
		}
		if role.MaxTTL > 0 {
			entry.MaxTTL = role.MaxTTL
		}
		if role.GenerateLease != nil {
			*entry.GenerateLease = *role.GenerateLease
		}
		entry.NoStore = role.NoStore
	}

	if entry.MaxTTL > 0 && entry.TTL > entry.MaxTTL {
		return logical.ErrorResponse(fmt.Sprintf("requested ttl of %s is greater than max ttl of %s", entry.TTL, entry.MaxTTL)), nil
	}

	return b.pathIssueSignCert(ctx, req, data, entry, true, true)
}

func (b *backend) pathIssueSignCert(ctx context.Context, req *logical.Request, data *framework.FieldData, role *roleEntry, useCSR, useCSRValues bool) (*logical.Response, error) {
	format := getFormat(data)
	if format == "" {
		return logical.ErrorResponse(
			`the "format" path parameter must be "pem", "der", or "pem_bundle"`), nil
	}

	var caErr error
	signingBundle, caErr := fetchCAInfo(ctx, req)
	switch caErr.(type) {
	case errutil.UserError:
		return nil, errutil.UserError{Err: fmt.Sprintf(
			"could not fetch the CA certificate (was one set?): %s", caErr)}
	case errutil.InternalError:
		return nil, errutil.InternalError{Err: fmt.Sprintf(
			"error fetching CA certificate: %s", caErr)}
	}

	input := &dataBundle{
		req:           req,
		apiData:       data,
		role:          role,
		signingBundle: signingBundle,
	}
	var parsedBundle *certutil.ParsedCertBundle
	var err error
	if useCSR {
		parsedBundle, err = signCert(b, input, false, useCSRValues)
	} else {
		parsedBundle, err = generateCert(ctx, b, input, false)
	}
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case errutil.InternalError:
			return nil, err
		}
	}

	signingCB, err := signingBundle.ToCertBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw signing bundle to cert bundle: {{err}}", err)
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, errwrap.Wrapf("error converting raw cert bundle to cert bundle: {{err}}", err)
	}

	respData := map[string]interface{}{
		"serial_number": cb.SerialNumber,
	}

	switch format {
	case "pem":
		respData["issuing_ca"] = signingCB.Certificate
		respData["certificate"] = cb.Certificate
		if cb.CAChain != nil && len(cb.CAChain) > 0 {
			respData["ca_chain"] = cb.CAChain
		}
		if !useCSR {
			respData["private_key"] = cb.PrivateKey
			respData["private_key_type"] = cb.PrivateKeyType
		}

	case "pem_bundle":
		respData["issuing_ca"] = signingCB.Certificate
		respData["certificate"] = cb.ToPEMBundle()
		if cb.CAChain != nil && len(cb.CAChain) > 0 {
			respData["ca_chain"] = cb.CAChain
		}
		if !useCSR {
			respData["private_key"] = cb.PrivateKey
			respData["private_key_type"] = cb.PrivateKeyType
		}

	case "der":
		respData["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		respData["issuing_ca"] = base64.StdEncoding.EncodeToString(signingBundle.CertificateBytes)

		var caChain []string
		for _, caCert := range parsedBundle.CAChain {
			caChain = append(caChain, base64.StdEncoding.EncodeToString(caCert.Bytes))
		}
		if caChain != nil && len(caChain) > 0 {
			respData["ca_chain"] = caChain
		}

		if !useCSR {
			respData["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
			respData["private_key_type"] = cb.PrivateKeyType
		}
	}

	var resp *logical.Response
	switch {
	case role.GenerateLease == nil:
		return nil, fmt.Errorf("generate lease in role is nil")
	case *role.GenerateLease == false:
		// If lease generation is disabled do not populate `Secret` field in
		// the response
		resp = &logical.Response{
			Data: respData,
		}
	default:
		resp = b.Secret(SecretCertsType).Response(
			respData,
			map[string]interface{}{
				"serial_number": cb.SerialNumber,
			})
		resp.Secret.TTL = parsedBundle.Certificate.NotAfter.Sub(time.Now())
	}

	if data.Get("private_key_format").(string) == "pkcs8" {
		err = convertRespToPKCS8(resp)
		if err != nil {
			return nil, err
		}
	}

	if !role.NoStore {
		err = req.Storage.Put(ctx, &logical.StorageEntry{
			Key:   "certs/" + normalizeSerial(cb.SerialNumber),
			Value: parsedBundle.CertificateBytes,
		})
		if err != nil {
			return nil, errwrap.Wrapf("unable to store certificate locally: {{err}}", err)
		}
	}

	if useCSR {
		if role.UseCSRCommonName && data.Get("common_name").(string) != "" {
			resp.AddWarning("the common_name field was provided but the role is set with \"use_csr_common_name\" set to true")
		}
		if role.UseCSRSANs && data.Get("alt_names").(string) != "" {
			resp.AddWarning("the alt_names field was provided but the role is set with \"use_csr_sans\" set to true")
		}
	}

	return resp, nil
}

const pathIssueHelpSyn = `
Request a certificate using a certain role with the provided details.
`

const pathIssueHelpDesc = `
This path allows requesting a certificate to be issued according to the
policy of the given role. The certificate will only be issued if the
requested details are allowed by the role policy.

This path returns a certificate and a private key. If you want a workflow
that does not expose a private key, generate a CSR locally and use the
sign path instead.
`

const pathSignHelpSyn = `
Request certificates using a certain role with the provided details.
`

const pathSignHelpDesc = `
This path allows requesting certificates to be issued according to the
policy of the given role. The certificate will only be issued if the
requested common name is allowed by the role policy.

This path requires a CSR; if you want Vault to generate a private key
for you, use the issue path instead.
`
