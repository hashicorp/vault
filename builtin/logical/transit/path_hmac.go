package transit

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathHMAC() *framework.Path {
	return &framework.Path{
		Pattern: "hmac/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("urlalgorithm"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The key to use for the HMAC function",
			},

			"input": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The base64-encoded input data",
			},

			"algorithm": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "sha2-256",
				Description: `Algorithm to use (POST body parameter). Valid values are:

* sha2-224
* sha2-256
* sha2-384
* sha2-512

Defaults to "sha2-256".`,
			},

			"urlalgorithm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Algorithm to use (POST URL parameter)`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `The version of the key to use for generating the HMAC.
Must be 0 (for latest) or a value greater than or equal
to the min_encryption_version configured on the key.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathHMACWrite,
		},

		HelpSynopsis:    pathHMACHelpSyn,
		HelpDescription: pathHMACHelpDesc,
	}
}

func (b *backend) pathHMACWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	ver := d.Get("key_version").(int)
	inputB64 := d.Get("input").(string)
	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}

	input, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode input as base64: %s", err)), logical.ErrInvalidRequest
	}

	// Get the policy
	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}

	switch {
	case ver == 0:
		// Allowed, will use latest; set explicitly here to ensure the string
		// is generated properly
		ver = p.LatestVersion
	case ver == p.LatestVersion:
		// Allowed
	case p.MinEncryptionVersion > 0 && ver < p.MinEncryptionVersion:
		return logical.ErrorResponse("cannot generate HMAC: version is too old (disallowed by policy)"), logical.ErrInvalidRequest
	}

	key, err := p.HMACKey(ver)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if key == nil {
		return nil, fmt.Errorf("HMAC key value could not be computed")
	}

	var hf hash.Hash
	switch algorithm {
	case "sha2-224":
		hf = hmac.New(sha256.New224, key)
	case "sha2-256":
		hf = hmac.New(sha256.New, key)
	case "sha2-384":
		hf = hmac.New(sha512.New384, key)
	case "sha2-512":
		hf = hmac.New(sha512.New, key)
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}
	hf.Write(input)
	retBytes := hf.Sum(nil)

	retStr := base64.StdEncoding.EncodeToString(retBytes)
	retStr = fmt.Sprintf("vault:v%s:%s", strconv.Itoa(ver), retStr)

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"hmac": retStr,
		},
	}
	return resp, nil
}

func (b *backend) pathHMACVerify(
	req *logical.Request, d *framework.FieldData, verificationHMAC string) (*logical.Response, error) {

	name := d.Get("name").(string)
	inputB64 := d.Get("input").(string)
	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}

	input, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode input as base64: %s", err)), logical.ErrInvalidRequest
	}

	// Verify the prefix
	if !strings.HasPrefix(verificationHMAC, "vault:v") {
		return logical.ErrorResponse("invalid HMAC to verify: no prefix"), logical.ErrInvalidRequest
	}

	splitVerificationHMAC := strings.SplitN(strings.TrimPrefix(verificationHMAC, "vault:v"), ":", 2)
	if len(splitVerificationHMAC) != 2 {
		return logical.ErrorResponse("invalid HMAC: wrong number of fields"), logical.ErrInvalidRequest
	}

	ver, err := strconv.Atoi(splitVerificationHMAC[0])
	if err != nil {
		return logical.ErrorResponse("invalid HMAC: version number could not be decoded"), logical.ErrInvalidRequest
	}

	verBytes, err := base64.StdEncoding.DecodeString(splitVerificationHMAC[1])
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode verification HMAC as base64: %s", err)), logical.ErrInvalidRequest
	}

	// Get the policy
	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}

	if ver > p.LatestVersion {
		return logical.ErrorResponse("invalid HMAC: version is too new"), logical.ErrInvalidRequest
	}

	if p.MinDecryptionVersion > 0 && ver < p.MinDecryptionVersion {
		return logical.ErrorResponse("cannot verify HMAC: version is too old (disallowed by policy)"), logical.ErrInvalidRequest
	}

	key, err := p.HMACKey(ver)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if key == nil {
		return nil, fmt.Errorf("HMAC key value could not be computed")
	}

	var hf hash.Hash
	switch algorithm {
	case "sha2-224":
		hf = hmac.New(sha256.New224, key)
	case "sha2-256":
		hf = hmac.New(sha256.New, key)
	case "sha2-384":
		hf = hmac.New(sha512.New384, key)
	case "sha2-512":
		hf = hmac.New(sha512.New, key)
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}
	hf.Write(input)
	retBytes := hf.Sum(nil)

	return &logical.Response{
		Data: map[string]interface{}{
			"valid": hmac.Equal(retBytes, verBytes),
		},
	}, nil
}

const pathHMACHelpSyn = `Generate an HMAC for input data using the named key`

const pathHMACHelpDesc = `
Generates an HMAC sum of the given algorithm and key against the given input data.
`
