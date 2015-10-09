package pki

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var rootAndSignSchema = map[string]*framework.FieldSchema{
	"common_name": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested common name; if you want more than
one, specify the alternative names in the alt_names
map. If not specified when signing, the common
name will be taken from the CSR; other names
must still be specified in alt_names or ip_sans.`,
	},

	"alt_names": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested Subject Alternative Names, if any,
in a comma-delimited list`,
	},

	"ip_sans": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested IP SANs, if any, in a
comma-delimited list`,
	},

	"ttl": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The requested Time To Live for the certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be larger than the mount max TTL. Note:
this only has an effect when generating
a CA cert or signing a CA cert, not when
creating a CSR for an intermediate CA.`,
	},

	"format": &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "pem",
		Description: `Format for returned data. Can be "pem" or "der";
defaults to "pem".`,
	},
}

var generateSchema = map[string]*framework.FieldSchema{
	"exported": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `Must be "internal" or "exported".
If set to "exported", the generated private
key will be returned. This is your *only*
chance to retrieve the private key!`,
	},

	"key_type": &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "rsa",
		Description: `The type of key to use; defaults to RSA. "rsa"
and "ec" are the only valid values.`,
	},

	"key_bits": &framework.FieldSchema{
		Type:    framework.TypeInt,
		Default: 2048,
		Description: `The number of bits to use. You will almost
certainly want to change this if you adjust
the key_type.`,
	},
}

func pathConfigCA(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/ca",
		Fields: map[string]*framework.FieldSchema{
			"pem_bundle": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `DEPRECATED: use "config/ca/set" instead.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathCASetWrite,
		},

		HelpSynopsis:    pathConfigCASetHelpSyn,
		HelpDescription: pathConfigCASetHelpDesc,
	}
}

func pathSetCA(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/ca/set",
		Fields: map[string]*framework.FieldSchema{
			"pem_bundle": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `PEM-format, concatenated unencrypted
secret key and certificate, or, if a
CSR was generated with the "generate"
endpoint, just the signed certificate.`,
			},

			"pki_address": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The base URL of the PKI mount, e.g.
"https://vault.example.com/v1/pki".
For HA setups, the given host name
should be the address that can always
be used to contact the leader, as this is
is used for generating the CA/CRL URLs in
the certificate. If empty, no CA/CRL
information will be encoded into
certificates.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathCASetWrite,
		},

		HelpSynopsis:    pathConfigCASetHelpSyn,
		HelpDescription: pathConfigCASetHelpDesc,
	}
}

func pathGenerateRootCA(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "config/ca/generate/root/" + framework.GenericNameRegex("exported"),

		Fields: rootAndSignSchema,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathCAGenerateRoot,
		},

		HelpSynopsis:    pathConfigCAGenerateHelpSyn,
		HelpDescription: pathConfigCAGenerateHelpDesc,
	}

	ret.Fields["pki_address"] = &framework.FieldSchema{
		Type: framework.TypeString,
		Description: `The base URL of the PKI mount, e.g.
"https://vault.example.com/v1/root_pki".
For HA setups, the given host name
should be the address that can always
be used to contact the leader, as this is
is used for generating the CA/CRL URLs in
the certificate. If empty, no CA/CRL
information will be encoded into
certificates.`,
	}

	for k, v := range generateSchema {
		ret.Fields[k] = v
	}

	return ret
}

func pathGenerateIntermediateCA(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "config/ca/generate/intermediate/" + framework.GenericNameRegex("exported"),

		Fields: generateSchema,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathCAGenerateIntermediate,
		},

		HelpSynopsis:    pathConfigCAGenerateHelpSyn,
		HelpDescription: pathConfigCAGenerateHelpDesc,
	}

	ret.Fields["format"] = &framework.FieldSchema{
		Type:    framework.TypeString,
		Default: "pem",
		Description: `Format for returned data. Can be "pem" or "der";
defaults to "pem".`,
	}

	return ret
}

func pathSignIntermediateCA(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "config/ca/sign",

		Fields: rootAndSignSchema,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathCASignIntermediate,
		},

		HelpSynopsis:    pathConfigCASignHelpSyn,
		HelpDescription: pathConfigCASignHelpDesc,
	}

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: `PEM-format CSR to be signed.`,
	}

	return ret
}

func (b *backend) pathCAGenerateRoot(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	var exported bool
	exportedStr := data.Get("exported").(string)
	switch exportedStr {
	case "exported":
		exported = true
	case "internal":
	default:
		return logical.ErrorResponse(
			`The "exported" path parameter must be "internal" or "exported"`,
		), nil
	}

	format := data.Get("format").(string)
	switch format {
	case "pem":
	case "der":
	default:
		return logical.ErrorResponse(
			`The "format" path parameter must be "pem" or "der"`,
		), nil
	}

	pkiAddress := strings.ToLower(data.Get("pki_address").(string))
	switch {
	case len(pkiAddress) == 0:
		break
	case !strings.HasPrefix(pkiAddress, "http"):
		return logical.ErrorResponse(
			`"pki_address" must be a URL`,
		), nil
	case !strings.Contains(pkiAddress, "/v1/"+req.MountPoint[:len(req.MountPoint)-1]):
		return logical.ErrorResponse(
			`"pki_address" needs to be the path to this mount"`,
		), nil
	}
	if strings.HasSuffix(pkiAddress, "/") {
		pkiAddress = pkiAddress[:len(pkiAddress)-1]
	}

	role := &roleEntry{
		TTL:              data.Get("ttl").(string),
		KeyType:          data.Get("key_type").(string),
		KeyBits:          data.Get("key_bits").(int),
		AllowLocalhost:   true,
		AllowAnyName:     true,
		EnforceHostnames: false,
	}

	switch role.KeyType {
	case "rsa":
		switch role.KeyBits {
		case 1024:
		case 2048:
		case 4096:
		case 8192:
		default:
			return logical.ErrorResponse(fmt.Sprintf("unsupported bit length for RSA key: %d", role.KeyBits)), nil
		}
	case "ec":
		switch role.KeyBits {
		case 224:
		case 256:
		case 384:
		case 521:
		default:
			return logical.ErrorResponse(fmt.Sprintf("unsupported bit length for EC key: %d", role.KeyBits)), nil
		}
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown key type %s", role.KeyType)), nil
	}

	var resp *logical.Response
	parsedBundle, err := generateCert(b, role, nil, true, pkiAddress, req, data)
	if err != nil {
		switch err.(type) {
		case certutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case certutil.InternalError:
			return nil, err
		}
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw cert bundle to cert bundle: %s", err)
	}

	resp = &logical.Response{
		Data: map[string]interface{}{
			"serial_number": cb.SerialNumber,
			"expiration":    int64(parsedBundle.Certificate.NotAfter.Unix()),
		},
	}

	switch format {
	case "pem":
		resp.Data["certificate"] = cb.Certificate
		resp.Data["issuing_ca"] = cb.IssuingCA
		if exported {
			resp.Data["private_key"] = cb.PrivateKey
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}
	case "der":
		resp.Data["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		resp.Data["issuing_ca"] = base64.StdEncoding.EncodeToString(parsedBundle.IssuingCABytes)
		if exported {
			resp.Data["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
			resp.Data["private_key_type"] = cb.PrivateKeyType
		}
	}

	entry, err := logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	entry.Key = "config/pki_address"
	entry.Value = []byte(pkiAddress)
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	// For ease of later use, also store just the certificate at a known
	// location, plus a fresh CRL
	entry.Key = "ca"
	entry.Value = parsedBundle.CertificateBytes
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	err = buildCRL(b, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *backend) pathCAGenerateIntermediate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	var exported bool
	exportedStr := data.Get("exported").(string)
	switch exportedStr {
	case "exported":
		exported = true
	case "internal":
	default:
		return logical.ErrorResponse(
			`The "exported" path parameter must be "internal" or "exported"`,
		), nil
	}

	format := data.Get("format").(string)
	switch format {
	case "pem":
	case "der":
	default:
		return logical.ErrorResponse(
			`The "format" path parameter must be "pem" or "der"`,
		), nil
	}

	role := &roleEntry{
		KeyType:          data.Get("key_type").(string),
		KeyBits:          data.Get("key_bits").(int),
		AllowLocalhost:   true,
		AllowAnyName:     true,
		EnforceHostnames: false,
	}

	switch role.KeyType {
	case "rsa":
		switch role.KeyBits {
		case 1024:
		case 2048:
		case 4096:
		case 8192:
		default:
			return logical.ErrorResponse(fmt.Sprintf("unsupported bit length for RSA key: %d", role.KeyBits)), nil
		}
	case "ec":
		switch role.KeyBits {
		case 224:
		case 256:
		case 384:
		case 521:
		default:
			return logical.ErrorResponse(fmt.Sprintf("unsupported bit length for EC key: %d", role.KeyBits)), nil
		}
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown key type %s", role.KeyType)), nil
	}

	var resp *logical.Response
	parsedBundle, err := generateCSR(b, role, nil, req, data)
	if err != nil {
		switch err.(type) {
		case certutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case certutil.InternalError:
			return nil, err
		}
	}

	csrb, err := parsedBundle.ToCSRBundle()
	if err != nil {
		return nil, fmt.Errorf("Error converting raw CSR bundle to CSR bundle: %s", err)
	}

	resp = &logical.Response{
		Data: map[string]interface{}{},
	}

	switch format {
	case "pem":
		resp.Data["csr"] = csrb.CSR
		if exported {
			resp.Data["private_key"] = csrb.PrivateKey
			resp.Data["private_key_type"] = csrb.PrivateKeyType
		}
	case "der":
		resp.Data["csr"] = base64.StdEncoding.EncodeToString(parsedBundle.CSRBytes)
		if exported {
			resp.Data["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
			resp.Data["private_key_type"] = csrb.PrivateKeyType
		}
	}

	cb := &certutil.CertBundle{
		PrivateKey:     csrb.PrivateKey,
		PrivateKeyType: csrb.PrivateKeyType,
	}

	entry, err := logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *backend) pathCASignIntermediate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	format := data.Get("format").(string)
	switch format {
	case "pem":
	case "der":
	default:
		return logical.ErrorResponse(
			`The "format" path parameter must be "pem" or "der"`,
		), nil
	}

	role := &roleEntry{
		TTL:              data.Get("ttl").(string),
		AllowLocalhost:   true,
		AllowAnyName:     true,
		EnforceHostnames: false,
	}

	if cn := data.Get("common_name").(string); len(cn) == 0 {
		role.UseCSRCommonName = true
	}

	var caErr error
	signingBundle, caErr := fetchCAInfo(req)
	switch caErr.(type) {
	case certutil.UserError:
		return nil, certutil.UserError{Err: fmt.Sprintf(
			"could not fetch the CA certificate (was one set?): %s", caErr)}
	case certutil.InternalError:
		return nil, certutil.InternalError{Err: fmt.Sprintf(
			"error fetching CA certificate: %s", caErr)}
	}

	parsedBundle, err := signCert(b, role, signingBundle, true, req, data)
	if err != nil {
		switch err.(type) {
		case certutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case certutil.InternalError:
			return nil, err
		}
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Error converting raw cert bundle to cert bundle: %s", err)
	}

	resp := b.Secret(SecretCertsType).Response(
		map[string]interface{}{
			"expiration":    int64(parsedBundle.Certificate.NotAfter.Unix()),
			"serial_number": cb.SerialNumber,
		},
		map[string]interface{}{
			"serial_number": cb.SerialNumber,
		})

	switch format {
	case "pem":
		resp.Data["certificate"] = cb.Certificate
		resp.Data["issuing_ca"] = cb.IssuingCA
	case "der":
		resp.Data["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		resp.Data["issuing_ca"] = base64.StdEncoding.EncodeToString(parsedBundle.IssuingCABytes)
	}

	resp.Secret.TTL = parsedBundle.Certificate.NotAfter.Sub(time.Now())

	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "certs/" + cb.SerialNumber,
		Value: parsedBundle.CertificateBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to store certificate locally")
	}

	return resp, nil
}

func (b *backend) pathCASetWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	pemBundle := data.Get("pem_bundle").(string)

	parsedBundle, err := certutil.ParsePEMBundle(pemBundle)
	if err != nil {
		switch err.(type) {
		case certutil.InternalError:
			return nil, err
		default:
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	// Handle the case of a self-signed certificate
	if parsedBundle.Certificate == nil && parsedBundle.IssuingCA != nil {
		parsedBundle.Certificate = parsedBundle.IssuingCA
		parsedBundle.CertificateBytes = parsedBundle.IssuingCABytes
	}

	pkiAddress := strings.ToLower(data.Get("pki_address").(string))
	switch {
	case len(pkiAddress) == 0:
		break
	case !strings.HasPrefix(pkiAddress, "http"):
		return logical.ErrorResponse(
			`"pki_address" must be a URL`,
		), nil
	case !strings.Contains(pkiAddress, "/v1/"+req.MountPoint[:len(req.MountPoint)-1]):
		return logical.ErrorResponse(
			`"pki_address" needs to be the path to this mount"`,
		), nil
	}
	if strings.HasSuffix(pkiAddress, "/") {
		pkiAddress = pkiAddress[:len(pkiAddress)-1]
	}

	cb := &certutil.CertBundle{}
	entry, err := req.Storage.Get("config/ca_bundle")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		err = entry.DecodeJSON(cb)
		if err != nil {
			return nil, err
		}
		// If we have a stored private key and did not get one, attempt to
		// correlate the two -- this could be due to a CSR being signed
		// for a generated CA cert and the resulting cert now being uploaded
		if len(cb.PrivateKey) != 0 &&
			cb.PrivateKeyType != "" &&
			parsedBundle.PrivateKeyType == certutil.UnknownPrivateKey &&
			(parsedBundle.PrivateKeyBytes == nil || len(parsedBundle.PrivateKeyBytes) == 0) {
			parsedCB, err := cb.ToParsedCertBundle()
			if err != nil {
				return nil, err
			}
			if parsedCB.PrivateKey == nil {
				return nil, fmt.Errorf("Encountered nil private key from saved key")
			}
			// If true, the stored private key corresponds to the cert's
			// public key, so fill it in
			//panic(fmt.Sprintf("\nparsedCB.PrivateKey.Public().: %#v\nparsedBundle.Certificate.PublicKey"))
			if reflect.DeepEqual(parsedCB.PrivateKey.Public(), parsedBundle.Certificate.PublicKey) {
				parsedBundle.PrivateKey = parsedCB.PrivateKey
				parsedBundle.PrivateKeyType = parsedCB.PrivateKeyType
				parsedBundle.PrivateKeyBytes = parsedCB.PrivateKeyBytes
			}
		}
	}

	if parsedBundle.PrivateKey == nil ||
		parsedBundle.PrivateKeyBytes == nil ||
		len(parsedBundle.PrivateKeyBytes) == 0 {
		return logical.ErrorResponse("No private key given and no matching key stored"), nil
	}

	if !parsedBundle.Certificate.IsCA {
		return logical.ErrorResponse("The given certificate is not marked for CA use and cannot be used with this backend"), nil
	}

	cb, err = parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Error converting raw values into cert bundle: %s", err)
	}

	entry, err = logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	entry.Key = "config/pki_address"
	entry.Value = []byte(pkiAddress)
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	// For ease of later use, also store just the certificate at a known
	// location, plus a fresh CRL
	entry.Key = "ca"
	entry.Value = parsedBundle.CertificateBytes
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	err = buildCRL(b, req)

	return nil, err
}

const pathConfigCASetHelpSyn = `
Set the CA certificate and private key used for generated credentials.
`

const pathConfigCASetHelpDesc = `
This sets the CA information used for credentials generated by this
by this mount. This must be a PEM-format, concatenated unencrypted
secret key and certificate.

For security reasons, the secret key cannot be retrieved later.
`

const pathConfigCAGenerateHelpSyn = `
Generate a new CA certificate and private key used for signing.
`

const pathConfigCAGenerateHelpDesc = `
This path generates a CA certificate and private key to be used for
credentials generated by this mount. The path can either
end in "internal" or "exported"; this controls whether the
unencrypted private key is exported after generation. This will
be your only chance to export the private key; for security reasons
it cannot be read or exported later.

If the "type" option is set to "self-signed", the generated
certificate will be a self-signed root CA. Otherwise, this mount
will act as an intermediate CA; a CSR will be returned, to be signed
by your chosen CA (which could be another mount of this backend).
Note that the CRL path will be set to this mount's CRL path; if you
need further customization it is recommended that you create a CSR
separately and get it signed. Either way, use the "config/ca/set"
endpoint to load the signed certificate into Vault.
`

const pathConfigCASignHelpSyn = `
Generate a signed CA certificate from a CSR.
`

const pathConfigCASignHelpDesc = `
This path generates a CA certificate to be used for credentials
generated by the certificate's destination mount.

Use the "config/ca/set" endpoint to load the signed certificate
into Vault another Vault mount.
`
