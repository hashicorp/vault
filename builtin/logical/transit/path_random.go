// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"

	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathRandom() *framework.Path {
	return &framework.Path{
		Pattern: "random(/" + framework.GenericNameRegex("source") + ")?" + framework.OptionalParamRegex("urlbytes"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "generate",
			OperationSuffix: "random|random-with-source|random-with-bytes|random-with-source-and-bytes",
		},

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
				Description: `Which system to source random data from, ether "platform", "seal", or "all".`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRandomWrite,
		},

		HelpSynopsis:    pathRandomHelpSyn,
		HelpDescription: pathRandomHelpDesc,
	}
}

func (b *backend) pathRandomWrite(_ context.Context, _ *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return random.HandleRandomAPI(d, b.GetRandomReader())
}

const pathRandomHelpSyn = `Generate random bytes`

const pathRandomHelpDesc = `
This function can be used to generate high-entropy random bytes.
`
