package transit

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathRandom() *framework.Path {
	return &framework.Path{
		Pattern: "random" + framework.OptionalParamRegex("urlbytes"),
		Fields: map[string]*framework.FieldSchema{
			"urlbytes": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The number of bytes to generate (POST URL parameter)",
			},

			"bytes": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Default:     32,
				Description: "The number of bytes to generate (POST body parameter). Defaults to 32 (256 bits).",
			},

			"format": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "base64",
				Description: `Encoding format to use. Can be "hex" or "base64". Defaults to "base64".`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRandomWrite,
		},

		HelpSynopsis:    pathRandomHelpSyn,
		HelpDescription: pathRandomHelpDesc,
	}
}

func (b *backend) pathRandomWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	bytes := 0
	var err error
	strBytes := d.Get("urlbytes").(string)
	if strBytes != "" {
		bytes, err = strconv.Atoi(strBytes)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error parsing url-set byte count: %s", err)), nil
		}
	} else {
		bytes = d.Get("bytes").(int)
	}
	format := d.Get("format").(string)

	if bytes < 1 {
		return logical.ErrorResponse(`"bytes" cannot be less than 1`), nil
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"hex\" or \"base64\"", format)), nil
	}

	randBytes, err := uuid.GenerateRandomBytes(bytes)
	if err != nil {
		return nil, err
	}

	var retStr string
	switch format {
	case "hex":
		retStr = hex.EncodeToString(randBytes)
	case "base64":
		retStr = base64.StdEncoding.EncodeToString(randBytes)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"random_bytes": retStr,
		},
	}
	return resp, nil
}

const pathRandomHelpSyn = `Generate random bytes`

const pathRandomHelpDesc = `
This function can be used to generate high-entropy random bytes.
`
