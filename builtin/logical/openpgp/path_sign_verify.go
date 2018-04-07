package openpgp

import (
	"bytes"
	"context"
	"crypto"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"strings"
)

func pathSign(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("urlalgorithm"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The key to use",
			},
			"input": {
				Type:        framework.TypeString,
				Description: "The base64-encoded input data",
			},
			"urlalgorithm": {
				Type:        framework.TypeString,
				Description: "Hash algorithm to use (POST URL parameter)",
			},
			"algorithm": {
				Type:    framework.TypeString,
				Default: "sha2-256",
				Description: `Hash algorithm to use (POST body parameter). Valid values are:

* sha2-224
* sha2-256
* sha2-384
* sha2-512

Defaults to "sha2-256".`,
			},
			"format": {
				Type:        framework.TypeString,
				Default:     "base64",
				Description: `Encoding format to use. Can be "base64" or "ascii-armor". Defaults to "base64".`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSignWrite,
		},
		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}
}

func pathVerify(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "verify/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The key to use",
			},
			"input": {
				Type:        framework.TypeString,
				Description: "The base64-encoded input data to verify",
			},
			"signature": {
				Type:        framework.TypeString,
				Description: "The signature",
			},
			"format": {
				Type:        framework.TypeString,
				Default:     "base64",
				Description: `Encoding format the signature use. Can be "base64" or "ascii-armor". Defaults to "base64".`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathVerifyWrite,
		},
		HelpSynopsis:    pathVerifyHelpSyn,
		HelpDescription: pathVerifyHelpDesc,
	}
}

func (b *backend) pathSignWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	inputB64 := data.Get("input").(string)
	input, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode input as base64: %s", err)), logical.ErrInvalidRequest
	}

	config := packet.Config{}

	algorithm := data.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = data.Get("algorithm").(string)
	}
	switch algorithm {
	case "sha2-224":
		config.DefaultHash = crypto.SHA224
	case "sha2-256":
		config.DefaultHash = crypto.SHA256
	case "sha2-384":
		config.DefaultHash = crypto.SHA384
	case "sha2-512":
		config.DefaultHash = crypto.SHA512
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}

	format := data.Get("format").(string)
	switch format {
	case "base64":
	case "ascii-armor":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"base64\" or \"ascii-armor\"", format)), nil
	}

	entry, err := b.key(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("key not found"), logical.ErrInvalidRequest
	}
	entity, err := b.entity(entry)
	if err != nil {
		return nil, err
	}

	message := bytes.NewReader(input)
	var signature bytes.Buffer
	switch format {
	case "ascii-armor":
		err = openpgp.ArmoredDetachSign(&signature, entity, message, &config)
		if err != nil {
			return nil, err
		}
	case "base64":
		encoder := base64.NewEncoder(base64.StdEncoding, &signature)
		err = openpgp.DetachSign(encoder, entity, message, &config)
		if err != nil {
			return nil, err
		}
		err = encoder.Close()
		if err != nil {
			return nil, err
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": signature.String(),
		},
	}, nil
}

func (b *backend) pathVerifyWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	inputB64 := data.Get("input").(string)
	input, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode input as base64: %s", err)), logical.ErrInvalidRequest
	}

	format := data.Get("format").(string)
	switch format {
	case "base64":
	case "ascii-armor":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"base64\" or \"ascii-armor\"", format)), nil
	}

	keyEntry, err := b.key(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if keyEntry == nil {
		return logical.ErrorResponse("key not found"), logical.ErrInvalidRequest
	}

	r := bytes.NewReader(keyEntry.SerializedKey)
	keyring, err := openpgp.ReadKeyRing(r)
	if err != nil {
		return nil, err
	}

	signature := strings.NewReader(data.Get("signature").(string))
	message := bytes.NewReader(input)
	switch format {
	case "base64":
		decoder := base64.NewDecoder(base64.StdEncoding, signature)
		_, err = openpgp.CheckDetachedSignature(keyring, message, decoder)
	case "ascii-armor":
		_, err = openpgp.CheckArmoredDetachedSignature(keyring, message, signature)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"valid": err == nil,
		},
	}

	return resp, nil
}

const pathSignHelpSyn = "Generate a signature for input data using the named PGP key"
const pathSignHelpDesc = "Generates a signature of the input data using the named PGP key."
const pathVerifyHelpSyn = "Verify a signature for input data created using the named PGP key"
const pathVerifyHelpDesc = "Verifies a signature of the input data using the named PGP key."
