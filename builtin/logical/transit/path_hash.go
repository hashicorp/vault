package transit

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathHash() *framework.Path {
	return &framework.Path{
		Pattern: "hash" + framework.OptionalParamRegex("urlalgorithm"),
		Fields: map[string]*framework.FieldSchema{
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

			"format": &framework.FieldSchema{
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

func (b *backend) pathHashWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	inputB64 := d.Get("input").(string)
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
