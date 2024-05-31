// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/sha3"
)

func (b *backend) pathHash() *framework.Path {
	return &framework.Path{
		Pattern: "hash" + framework.OptionalParamRegex("urlalgorithm"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "hash",
			OperationSuffix: "|with-algorithm",
		},

		Fields: map[string]*framework.FieldSchema{
			"input": {
				Type:        framework.TypeString,
				Description: "The base64-encoded input data",
			},

			"algorithm": {
				Type:    framework.TypeString,
				Default: "sha2-256",
				Description: `Algorithm to use (POST body parameter). Valid values are:

* sha2-224
* sha2-256
* sha2-384
* sha2-512
* sha3-224
* sha3-256
* sha3-384
* sha3-512

Defaults to "sha2-256".`,
			},

			"urlalgorithm": {
				Type:        framework.TypeString,
				Description: `Algorithm to use (POST URL parameter)`,
			},

			"format": {
				Type:        framework.TypeString,
				Default:     "hex",
				Description: `Encoding format to use. Can be "hex" or "base64". Defaults to "hex".`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathHashWrite,
		},

		HelpSynopsis:    pathHashHelpSyn,
		HelpDescription: pathHashHelpDesc,
	}
}

func (b *backend) pathHashWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rawInput, ok, err := d.GetOkErr("input")
	if err != nil {
		return nil, err
	}
	if !ok {
		return logical.ErrorResponse("input missing"), logical.ErrInvalidRequest
	}

	inputB64 := rawInput.(string)
	format := d.Get("format").(string)
	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}

	input, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode input as base64: %s", err)), logical.ErrInvalidRequest
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"hex\" or \"base64\"", format)), nil
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
	case "sha3-224":
		hf = sha3.New224()
	case "sha3-256":
		hf = sha3.New256()
	case "sha3-384":
		hf = sha3.New384()
	case "sha3-512":
		hf = sha3.New512()
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported algorithm %s", algorithm)), nil
	}
	hf.Write(input)
	retBytes := hf.Sum(nil)

	var retStr string
	switch format {
	case "hex":
		retStr = hex.EncodeToString(retBytes)
	case "base64":
		retStr = base64.StdEncoding.EncodeToString(retBytes)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"sum": retStr,
		},
	}
	return resp, nil
}

const pathHashHelpSyn = `Generate a hash sum for input data`

const pathHashHelpDesc = `
Generates a hash sum of the given algorithm against the given input data.
`
