package transit

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathSign() *framework.Path {
	return &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("urlalgorithm"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The key to use",
			},

			"input": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The base64-encoded input data",
			},

			"algorithm": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "sha2-256",
				Description: `Hash algorithm to use (POST body parameter). Valid values are:

* sha2-224
* sha2-256
* sha2-384
* sha2-512

Defaults to "sha2-256".`,
			},

			"urlalgorithm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Hash algorithm to use (POST URL parameter)`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSignWrite,
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}
}

func (b *backend) pathVerify() *framework.Path {
	return &framework.Path{
		Pattern: "verify/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("urlalgorithm"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The key to use",
			},

			"signature": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The signature, including vault header/key version",
			},

			"hmac": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The HMAC, including vault header/key version",
			},

			"input": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The base64-encoded input data to verify",
			},

			"urlalgorithm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Hash algorithm to use (POST URL parameter)`,
			},

			"algorithm": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "sha2-256",
				Description: `Hash algorithm to use (POST body parameter). Valid values are:

* sha2-224
* sha2-256
* sha2-384
* sha2-512

Defaults to "sha2-256".`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathVerifyWrite,
		},

		HelpSynopsis:    pathVerifyHelpSyn,
		HelpDescription: pathVerifyHelpDesc,
	}
}

func (b *backend) pathSignWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

	var hf hash.Hash
	switch algorithm {
	case "sha2-224":
		hf = sha256.New224()
	case "sha2-256":
		hf = sha256.New()
	case "sha2-384":
		hf = sha512.New384()
	case "sha2-512":
		hf = sha512.New()
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}
	hf.Write(input)
	hashedInput := hf.Sum(nil)

	// Get the policy
	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
	}

	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type %v does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	sig, err := p.Sign(hashedInput)
	if err != nil {
		return nil, err
	}
	if sig == "" {
		return nil, fmt.Errorf("signature could not be computed")
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"signature": sig,
		},
	}
	return resp, nil
}

func (b *backend) pathVerifyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	sig := d.Get("signature").(string)
	hmac := d.Get("hmac").(string)
	switch {
	case sig != "" && hmac != "":
		return logical.ErrorResponse("provide one of 'signature' or 'hmac'"), logical.ErrInvalidRequest

	case sig == "" && hmac == "":
		return logical.ErrorResponse("neither a 'signature' nor an 'hmac' were given to verify"), logical.ErrInvalidRequest

	case hmac != "":
		return b.pathHMACVerify(req, d, hmac)
	}

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

	var hf hash.Hash
	switch algorithm {
	case "sha2-224":
		hf = sha256.New224()
	case "sha2-256":
		hf = sha256.New()
	case "sha2-384":
		hf = sha512.New384()
	case "sha2-512":
		hf = sha512.New()
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}
	hf.Write(input)
	hashedInput := hf.Sum(nil)

	// Get the policy
	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
	}

	valid, err := p.VerifySignature(hashedInput, sig)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		case errutil.InternalError:
			return nil, err
		default:
			return nil, err
		}
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"valid": valid,
		},
	}
	return resp, nil
}

const pathSignHelpSyn = `Generate a signature for input data using the named key`

const pathSignHelpDesc = `
Generates a signature of the input data using the named key and the given hash algorithm.
`
const pathVerifyHelpSyn = `Verify a signature or HMAC for input data created using the named key`

const pathVerifyHelpDesc = `
Verifies a signature or HMAC of the input data using the named key and the given hash algorithm.
`
