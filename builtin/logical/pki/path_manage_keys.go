package pki

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerateKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/generate/(internal|exported)",

		Fields: map[string]*framework.FieldSchema{
			keyNameParam: {
				Type:        framework.TypeString,
				Description: "Optional name to be used for this key",
			},
			keyTypeParam: {
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
	pathGenerateKeyHelpSyn  = `Generate a new private key used for signing.`
	pathGenerateKeyHelpDesc = `This endpoint will generate a new key pair of the specified type (internal, exported, or kms).`
)

func (b *backend) pathGenerateKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	keyName, err := getKeyName(ctx, req.Storage, data)
	if err != nil { // Fail Immediately if Key Name is in Use, etc...
		return nil, err
	}
	keyType := data.Get(keyTypeParam).(string)
	keyBits := data.Get("key_bits").(int)

	switch {
	case strings.HasSuffix(req.Path, "/internal"):
		// Internal key generation, stored in storage
		keyBundle, err := certutil.GetKeyBundleFromKeyGenerator(keyType, keyBits, nil)
		if err != nil {
			return nil, err
		}
		privateKeyPemString, err := keyBundle.ToPrivateKeyPemString()
		if err != nil {
			return nil, err
		}
		key, _, err := importKey(ctx, req.Storage, privateKeyPemString, keyName)
		if err != nil {
			return nil, err
		}
		resp := logical.Response{
			Data: map[string]interface{}{
				keyIdParam:   key.ID,
				keyNameParam: key.Name,
				keyTypeParam: key.PrivateKeyType,
			},
		}
		return &resp, nil
	case strings.HasSuffix(req.Path, "/exported"):
		// Same as internal key generation but we return the generated key
		keyBundle, err := certutil.GetKeyBundleFromKeyGenerator(keyType, keyBits, nil)
		if err != nil {
			return nil, err
		}
		privateKeyPemString, err := keyBundle.ToPrivateKeyPemString()
		if err != nil {
			return nil, err
		}
		key, _, err := importKey(ctx, req.Storage, privateKeyPemString, keyName)
		if err != nil {
			return nil, err
		}
		resp := logical.Response{
			Data: map[string]interface{}{
				keyIdParam:    key.ID,
				keyNameParam:  key.Name,
				keyTypeParam:  key.PrivateKeyType,
				"private_key": privateKeyPemString,
			},
		}
		return &resp, nil
	case strings.HasSuffix(req.Path, "/kms"):
		return nil, errEntOnly
	default:
		return logical.ErrorResponse("Unknown type of key to generate"), nil
	}
}

func pathImportKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/import",

		Fields: map[string]*framework.FieldSchema{
			keyNameParam: {
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
	pathImportKeyHelpSyn  = `Import the specified key.`
	pathImportKeyHelpDesc = `This endpoint allows importing a specified issuer key from a pem bundle.
If name is set, that will be set on the key.`
)

func (b *backend) pathImportKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	keyValueInterface, isOk := data.GetOk("pem_bundle")
	if !isOk {
		return logical.ErrorResponse("keyValue must be set"), nil
	}
	keyValue := keyValueInterface.(string)
	keyName := data.Get(keyNameParam).(string)

	key, existed, err := importKey(ctx, req.Storage, keyValue, keyName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	resp := logical.Response{
		Data: map[string]interface{}{
			keyIdParam:   key.ID,
			keyNameParam: key.Name,
			keyTypeParam: key.PrivateKeyType,
			"backing":    "", // This would show up as "Managed" in "type"
		},
	}

	if existed {
		resp.AddWarning("Key already imported, use key/ endpoint to update name.")
	}

	return &resp, nil
}
