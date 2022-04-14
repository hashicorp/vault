package pki

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerateKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/generate/",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback:                    b.pathGenerateKeyHandler,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathGenerateKeyHelpSyn,
		HelpDescription: pathGenerateKeyHelpDesc,
	}
}

const pathGenerateKeyHelpSyn = ``
const pathGenerateKeyHelpDesc = ``

func (b *backend) pathGenerateKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}

func pathImportKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/import",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback:                    b.pathImportKeyHandler,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathImportKeyHelpSyn,
		HelpDescription: pathImportKeyHelpDesc,
	}
}

const pathImportKeyHelpSyn = ``
const pathImportKeyHelpDesc = ``

func (b *backend) pathImportKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
