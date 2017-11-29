package transit

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

// BatchRequestItem represents a request item for batch processing
type BatchRequestItem struct {
	// Context for key derivation. This is required for derived keys.
	Context string `json:"context" structs:"context" mapstructure:"context"`

	// DecodedContext is the base64 decoded version of Context
	DecodedContext []byte

	// Plaintext for encryption
	Plaintext string `json:"plaintext" structs:"plaintext" mapstructure:"plaintext"`

	// Ciphertext for decryption
	Ciphertext string `json:"ciphertext" structs:"ciphertext" mapstructure:"ciphertext"`

	// Nonce to be used when v1 convergent encryption is used
	Nonce string `json:"nonce" structs:"nonce" mapstructure:"nonce"`

	// The key version to be used for encryption
	KeyVersion int `json:"key_version" structs:"key_version" mapstructure:"key_version"`

	// DecodedNonce is the base64 decoded version of Nonce
	DecodedNonce []byte
}

// BatchResponseItem represents a response item for batch processing
type BatchResponseItem struct {
	// Ciphertext for the plaintext present in the corresponding batch
	// request item
	Ciphertext string `json:"ciphertext,omitempty" structs:"ciphertext" mapstructure:"ciphertext"`

	// Plaintext for the ciphertext present in the corresponsding batch
	// request item
	Plaintext string `json:"plaintext,omitempty" structs:"plaintext" mapstructure:"plaintext"`

	// Error, if set represents a failure encountered while encrypting a
	// corresponding batch request item
	Error string `json:"error,omitempty" structs:"error" mapstructure:"error"`
}

func (b *backend) pathEncrypt() *framework.Path {
	return &framework.Path{
		Pattern: "encrypt/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},

			"plaintext": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Base64 encoded plaintext value to be encrypted",
			},

			"context": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Base64 encoded context for key derivation. Required if key derivation is enabled",
			},

			"nonce": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
Base64 encoded nonce value. Must be provided if convergent encryption is
enabled for this key and the key was generated with Vault 0.6.1. Not required
for keys created in 0.6.2+. The value must be exactly 96 bits (12 bytes) long
and the user must ensure that for any given context (and thus, any given
encryption key) this nonce value is **never reused**.
`,
			},

			"type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "aes256-gcm96",
				Description: `
This parameter is required when encryption key is expected to be created.
When performing an upsert operation, the type of key to create. Currently,
"aes256-gcm96" (symmetric) is the only type supported. Defaults to
"aes256-gcm96".`,
			},

			"convergent_encryption": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
This parameter will only be used when a key is expected to be created.  Whether
to support convergent encryption. This is only supported when using a key with
key derivation enabled and will require all requests to carry both a context
and 96-bit (12-byte) nonce. The given nonce will be used in place of a randomly
generated nonce. As a result, when the same context and nonce are supplied, the
same ciphertext is generated. It is *very important* when using this mode that
you ensure that all nonces are unique for a given context.  Failing to do so
will severely impact the ciphertext's security.`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `The version of the key to use for encryption.
Must be 0 (for latest) or a value greater than or equal
to the min_encryption_version configured on the key.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathEncryptWrite,
			logical.UpdateOperation: b.pathEncryptWrite,
		},

		ExistenceCheck: b.pathEncryptExistenceCheck,

		HelpSynopsis:    pathEncryptHelpSyn,
		HelpDescription: pathEncryptHelpDesc,
	}
}

func (b *backend) pathEncryptExistenceCheck(
	req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)
	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return false, err
	}
	return p != nil, nil
}

func (b *backend) pathEncryptWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	var err error

	batchInputRaw := d.Raw["batch_input"]
	var batchInputItems []BatchRequestItem
	if batchInputRaw != nil {
		err = mapstructure.Decode(batchInputRaw, &batchInputItems)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch input: %v", err)
		}

		if len(batchInputItems) == 0 {
			return logical.ErrorResponse("missing batch input to process"), logical.ErrInvalidRequest
		}
	} else {
		valueRaw, ok := d.GetOk("plaintext")
		if !ok {
			return logical.ErrorResponse("missing plaintext to encrypt"), logical.ErrInvalidRequest
		}

		batchInputItems = make([]BatchRequestItem, 1)
		batchInputItems[0] = BatchRequestItem{
			Plaintext:  valueRaw.(string),
			Context:    d.Get("context").(string),
			Nonce:      d.Get("nonce").(string),
			KeyVersion: d.Get("key_version").(int),
		}
	}

	batchResponseItems := make([]BatchResponseItem, len(batchInputItems))
	contextSet := len(batchInputItems[0].Context) != 0

	// Before processing the batch request items, get the policy. If the
	// policy is supposed to be upserted, then determine if 'derived' is to
	// be set or not, based on the presence of 'context' field in all the
	// input items.
	for i, item := range batchInputItems {
		if (len(item.Context) == 0 && contextSet) || (len(item.Context) != 0 && !contextSet) {
			return logical.ErrorResponse("context should be set either in all the request blocks or in none"), logical.ErrInvalidRequest
		}

		_, err := base64.StdEncoding.DecodeString(item.Plaintext)
		if err != nil {
			batchResponseItems[i].Error = err.Error()
			continue
		}

		// Decode the context
		if len(item.Context) != 0 {
			batchInputItems[i].DecodedContext, err = base64.StdEncoding.DecodeString(item.Context)
			if err != nil {
				batchResponseItems[i].Error = err.Error()
				continue
			}
		}

		// Decode the nonce
		if len(item.Nonce) != 0 {
			batchInputItems[i].DecodedNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
			if err != nil {
				batchResponseItems[i].Error = err.Error()
				continue
			}
		}
	}

	// Get the policy
	var p *keysutil.Policy
	var lock *sync.RWMutex
	var upserted bool
	if req.Operation == logical.CreateOperation {
		convergent := d.Get("convergent_encryption").(bool)
		if convergent && !contextSet {
			return logical.ErrorResponse("convergent encryption requires derivation to be enabled, so context is required"), nil
		}

		polReq := keysutil.PolicyRequest{
			Storage:    req.Storage,
			Name:       name,
			Derived:    contextSet,
			Convergent: convergent,
		}

		keyType := d.Get("type").(string)
		switch keyType {
		case "aes256-gcm96":
			polReq.KeyType = keysutil.KeyType_AES256_GCM96
		case "ecdsa-p256":
			return logical.ErrorResponse(fmt.Sprintf("key type %v not supported for this operation", keyType)), logical.ErrInvalidRequest
		default:
			return logical.ErrorResponse(fmt.Sprintf("unknown key type %v", keyType)), logical.ErrInvalidRequest
		}

		p, lock, upserted, err = b.lm.GetPolicyUpsert(polReq)

	} else {
		p, lock, err = b.lm.GetPolicyShared(req.Storage, name)
	}
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}

	// Process batch request items. If encryption of any request
	// item fails, respectively mark the error in the response
	// collection and continue to process other items.
	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			continue
		}

		ciphertext, err := p.Encrypt(item.KeyVersion, item.DecodedContext, item.DecodedNonce, item.Plaintext)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				batchResponseItems[i].Error = err.Error()
				continue
			default:
				return nil, err
			}
		}

		if ciphertext == "" {
			return nil, fmt.Errorf("empty ciphertext returned for input item %d", i)
		}

		batchResponseItems[i].Ciphertext = ciphertext
	}

	resp := &logical.Response{}
	if batchInputRaw != nil {
		resp.Data = map[string]interface{}{
			"batch_results": batchResponseItems,
		}
	} else {
		if batchResponseItems[0].Error != "" {
			return logical.ErrorResponse(batchResponseItems[0].Error), logical.ErrInvalidRequest
		}
		resp.Data = map[string]interface{}{
			"ciphertext": batchResponseItems[0].Ciphertext,
		}
	}

	if req.Operation == logical.CreateOperation && !upserted {
		resp.AddWarning("Attempted creation of the key during the encrypt operation, but it was created beforehand")
	}
	return resp, nil
}

const pathEncryptHelpSyn = `Encrypt a plaintext value or a batch of plaintext
blocks using a named key`

const pathEncryptHelpDesc = `
This path uses the named key from the request path to encrypt a user provided
plaintext or a batch of plaintext blocks. The plaintext must be base64 encoded.
`
