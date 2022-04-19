package pki

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/certutil"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerateKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/generate/(internal|exported)",

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Optional name to be used for this key",
			},
			"key_type": {
				Type:        framework.TypeString,
				Default:     "rsa",
				Description: `Type of the secret key to generate`,
			},
			"key_bits": {
				Type:        framework.TypeInt,
				Default:     2048,
				Description: `Type of the secret key to generate`,
			},
		},

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

const (
	pathGenerateKeyHelpSyn  = ``
	pathGenerateKeyHelpDesc = ``
)

func (b *backend) pathGenerateKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keyName := data.Get("name").(string)
	keyType := data.Get("key_type").(string)
	keyBits := data.Get("key_bits").(int)
	switch {
	case strings.HasSuffix(req.Path, "/internal"):
		// Internal key generation, stored in storage
		keyBundle, err := certutil.GetKeyBundleFromKeyGenerator(keyType, keyBits, nil)
		if err != nil {
			return nil, err
		}
		importKey(ctx, req.Storage, string(keyBundle.PrivateKeyBytes), keyName)
	case strings.HasSuffix(req.Path, "/exported"):
		keyRef := data.Get("key_ref").(string)
		keyBundle, err := certutil.GetKeyBundleFromKeyGenerator(keyType, keyBits, existingGeneratePrivateKey(ctx, req.Storage, keyRef))
		if err != nil {
			return nil, err
		}
		importKey(ctx, req.Storage, string(keyBundle.PrivateKeyBytes), keyName)
	default:
		return logical.ErrorResponse("Unknown type of key to generate"), nil
	}

	return nil, nil
}

func pathImportKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/import",

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Optional name to be used for this key",
			},
			"pem_bundle": {
				Type:        framework.TypeString,
				Description: `PEM-format, unencrypted secret key`,
			},
		},

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

const (
	pathImportKeyHelpSyn  = ``
	pathImportKeyHelpDesc = ``
)

func (b *backend) pathImportKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keyValueInterface, isOk := data.GetOk("pem_bundle")
	if !isOk {
		return logical.ErrorResponse("keyValue must be set"), nil
	}
	keyValue := keyValueInterface.(string)
	keyName := data.Get("name").(string)

	key, existed, err := importKey(ctx, req.Storage, keyValue, keyName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	resp := logical.Response{
		Data: map[string]interface{}{
			"id":      key.ID,
			"name":    key.Name,
			"type":    key.PrivateKeyType,
			"backing": "", // This would show up as "Managed" in "type"
		},
	}

	if existed {
		resp.AddWarning("Key already imported, use key/ endpoint to update name.")
	}

	return &resp, nil
}
