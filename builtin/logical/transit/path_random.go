package transit

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/vault/sdk/helper/xor"
	"io"
	"k8s.io/utils/strings/slices"
	"strconv"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const maxBytes = 128 * 1024

func (b *backend) pathRandom() *framework.Path {
	return &framework.Path{
		Pattern: "random(/" + framework.GenericNameRegex("source") + ")?" + framework.OptionalParamRegex("urlbytes"),
		Fields: map[string]*framework.FieldSchema{
			"urlbytes": {
				Type:        framework.TypeString,
				Description: "The number of bytes to generate (POST URL parameter)",
			},

			"bytes": {
				Type:        framework.TypeInt,
				Default:     32,
				Description: "The number of bytes to generate (POST body parameter). Defaults to 32 (256 bits).",
			},

			"format": {
				Type:        framework.TypeString,
				Default:     "base64",
				Description: `Encoding format to use. Can be "hex" or "base64". Defaults to "base64".`,
			},

			"source": {
				Type:        framework.TypeString,
				Default:     "platform",
				Description: `Which system to source entropy from, ether "platform", "seal", or "all".`,
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

	//Parsing is convoluted here, but allows operators to ACL both bytes and entropy source
	maybeUrlBytes := d.Raw["urlbytes"]
	maybeSource := d.Raw["source"]
	source := "platform"
	var err error
	if maybeSource == "" {
		bytes = d.Get("bytes").(int)
	} else if maybeUrlBytes == "" && slices.Contains([]string{"", "platform", "seal", "all"}, maybeSource.(string)) {
		source = maybeSource.(string)
		bytes = d.Get("bytes").(int)
	} else if maybeUrlBytes == "" {
		bytes, err = strconv.Atoi(maybeSource.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error parsing url-set byte count: %s", err)), nil
		}
	} else {
		source = maybeSource.(string)
		bytes, err = strconv.Atoi(maybeUrlBytes.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error parsing url-set byte count: %s", err)), nil
		}
	}
	format := d.Get("format").(string)

	if bytes < 1 {
		return logical.ErrorResponse(`"bytes" cannot be less than 1`), nil
	}

	if bytes > maxBytes {
		return logical.ErrorResponse(`"bytes" should be less than %d`, maxBytes), nil
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse("unsupported encoding format %q; must be \"hex\" or \"base64\"", format), nil
	}

	var randBytes []byte
	var warning string
	switch source {
	case "", "platform":
		randBytes, err = uuid.GenerateRandomBytes(bytes)
		if err != nil {
			return nil, err
		}
	case "seal":
		if rand.Reader == b.GetRandomReader() {
			warning = "no seal/entropy augmentation available, using platform entropy source"
		}
		randBytes = make([]byte, bytes)
		_, err = io.ReadFull(b.GetRandomReader(), randBytes)
		if err != nil {
			return nil, err
		}
	case "all":
		sealBytes := make([]byte, bytes)
		_, err = io.ReadFull(b.GetRandomReader(), sealBytes)
		if err != nil {
			return nil, err
		}
		randBytes, err = uuid.GenerateRandomBytes(bytes)
		if err != nil {
			return nil, err
		}
		randBytes, err = xor.XORBytes(sealBytes, randBytes)
		if err != nil {
			return nil, err
		}
	default:
		return logical.ErrorResponse("unsupported entropy source %q; must be \"platform\" or \"seal\", or \"all\"", source), nil
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
	if warning != "" {
		resp.Warnings = []string{warning}
	}
	return resp, nil
}

const pathRandomHelpSyn = `Generate random bytes`

const pathRandomHelpDesc = `
This function can be used to generate high-entropy random bytes.
`
