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

			"context": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Base64 encoded context for key derivation. Required if key
derivation is enabled; currently only available with ed25519 keys.`,
			},

			"algorithm": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "sha2-256",
				Description: `Hash algorithm to use (POST body parameter). Valid values are:

* sha2-224
* sha2-256
* sha2-384
* sha2-512

Defaults to "sha2-256". Not valid for all key types,
including ed25519.`,
			},

			"urlalgorithm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Hash algorithm to use (POST URL parameter)`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `The version of the key to use for signing.
Must be 0 (for latest) or a value greater than or equal
to the min_encryption_version configured on the key.`,
			},

			"prehashed": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: `Set to 'true' when the input is already hashed. If the key type is 'rsa-2048' or 'rsa-4096', then the algorithm used to hash the input should be indicated by the 'algorithm' parameter.`,
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

			"context": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Base64 encoded context for key derivation. Required if key
derivation is enabled; currently only available with ed25519 keys.`,
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

Defaults to "sha2-256". Not valid for all key types.`,
			},

			"prehashed": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: `Set to 'true' when the input is already hashed. If the key type is 'rsa-2048' or 'rsa-4096', then the algorithm used to hash the input should be indicated by the 'algorithm' parameter.`,
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
	ver := d.Get("key_version").(int)
	inputB64 := d.Get("input").(string)
	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}
	prehashed := d.Get("prehashed").(bool)

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

	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type %v does not support signing", p.Type)), logical.ErrInvalidRequest
	}

	contextRaw := d.Get("context").(string)
	var context []byte
	if len(contextRaw) != 0 {
		context, err = base64.StdEncoding.DecodeString(contextRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode context"), logical.ErrInvalidRequest
		}
	}

	if p.Type.HashSignatureInput() && !prehashed {
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
		input = hf.Sum(nil)
	}

	sig, err := p.Sign(ver, context, input, algorithm)
	if err != nil {
		return nil, err
	}
	if sig == nil {
		return nil, fmt.Errorf("signature could not be computed")
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"signature": sig.Signature,
		},
	}

	if len(sig.PublicKey) > 0 {
		resp.Data["public_key"] = sig.PublicKey
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
	prehashed := d.Get("prehashed").(bool)

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

	if !p.Type.SigningSupported() {
		return logical.ErrorResponse(fmt.Sprintf("key type %v does not support verification", p.Type)), logical.ErrInvalidRequest
	}

	contextRaw := d.Get("context").(string)
	var context []byte
	if len(contextRaw) != 0 {
		context, err = base64.StdEncoding.DecodeString(contextRaw)
		if err != nil {
			return logical.ErrorResponse("failed to base64-decode context"), logical.ErrInvalidRequest
		}
	}

	if p.Type.HashSignatureInput() && !prehashed {
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
		input = hf.Sum(nil)
	}

	valid, err := p.VerifySignature(context, input, sig, algorithm)
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
