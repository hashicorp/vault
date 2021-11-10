package transit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
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

// EncryptBatchResponseItem represents a response item for batch processing
type EncryptBatchResponseItem struct {
	// Ciphertext for the plaintext present in the corresponding batch
	// request item
	Ciphertext string `json:"ciphertext,omitempty" structs:"ciphertext" mapstructure:"ciphertext"`

	// KeyVersion defines the key version used to encrypt plaintext.
	KeyVersion int `json:"key_version,omitempty" structs:"key_version" mapstructure:"key_version"`

	// Error, if set represents a failure encountered while encrypting a
	// corresponding batch request item
	Error string `json:"error,omitempty" structs:"error" mapstructure:"error"`
}

func (b *backend) pathEncrypt() *framework.Path {
	return &framework.Path{
		Pattern: "encrypt/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},

			"plaintext": {
				Type:        framework.TypeString,
				Description: "Base64 encoded plaintext value to be encrypted",
			},

			"context": {
				Type:        framework.TypeString,
				Description: "Base64 encoded context for key derivation. Required if key derivation is enabled",
			},

			"nonce": {
				Type: framework.TypeString,
				Description: `
Base64 encoded nonce value. Must be provided if convergent encryption is
enabled for this key and the key was generated with Vault 0.6.1. Not required
for keys created in 0.6.2+. The value must be exactly 96 bits (12 bytes) long
and the user must ensure that for any given context (and thus, any given
encryption key) this nonce value is **never reused**.
`,
			},

			"type": {
				Type:    framework.TypeString,
				Default: "aes256-gcm96",
				Description: `
This parameter is required when encryption key is expected to be created.
When performing an upsert operation, the type of key to create. Currently,
"aes128-gcm96" (symmetric) and "aes256-gcm96" (symmetric) are the only types supported. Defaults to "aes256-gcm96".`,
			},

			"convergent_encryption": {
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

			"key_version": {
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

// decodeBatchRequestItems is a fast path alternative to mapstructure.Decode to decode []BatchRequestItem.
// It aims to behave as closely possible to the original mapstructure.Decode and will return the same errors.
// https://github.com/hashicorp/vault/pull/8775/files#r437709722
func decodeBatchRequestItems(src interface{}, dst *[]BatchRequestItem) error {
	if src == nil || dst == nil {
		return nil
	}

	items, ok := src.([]interface{})
	if !ok {
		return fmt.Errorf("source data must be an array or slice, got %T", src)
	}

	// Early return should happen before allocating the array if the batch is empty.
	// However to comply with mapstructure output it's needed to allocate an empty array.
	sitems := len(items)
	*dst = make([]BatchRequestItem, sitems)
	if sitems == 0 {
		return nil
	}

	// To comply with mapstructure output the same error type is needed.
	var errs mapstructure.Error

	for i, iitem := range items {
		item, ok := iitem.(map[string]interface{})
		if !ok {
			return fmt.Errorf("[%d] expected a map, got '%T'", i, iitem)
		}

		if v, has := item["context"]; has {
			if !reflect.ValueOf(v).IsValid() {
			} else if casted, ok := v.(string); ok {
				(*dst)[i].Context = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].context' expected type 'string', got unconvertible type '%T'", i, item["context"]))
			}
		}

		if v, has := item["ciphertext"]; has {
			if !reflect.ValueOf(v).IsValid() {
			} else if casted, ok := v.(string); ok {
				(*dst)[i].Ciphertext = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].ciphertext' expected type 'string', got unconvertible type '%T'", i, item["ciphertext"]))
			}
		}

		// don't allow "null" to be passed in for the plaintext value
		if v, has := item["plaintext"]; has {
			if casted, ok := v.(string); ok {
				(*dst)[i].Plaintext = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].plaintext' expected type 'string', got unconvertible type '%T'", i, item["plaintext"]))
			}
		}

		if v, has := item["nonce"]; has {
			if !reflect.ValueOf(v).IsValid() {
			} else if casted, ok := v.(string); ok {
				(*dst)[i].Nonce = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].nonce' expected type 'string', got unconvertible type '%T'", i, item["nonce"]))
			}
		}

		if v, has := item["key_version"]; has {
			if !reflect.ValueOf(v).IsValid() {
			} else if casted, ok := v.(int); ok {
				(*dst)[i].KeyVersion = casted
			} else if js, ok := v.(json.Number); ok {
				// https://github.com/hashicorp/vault/issues/10232
				// Because API server parses json request with UseNumber=true, logical.Request.Data can include json.Number for a number field.
				if casted, err := js.Int64(); err == nil {
					(*dst)[i].KeyVersion = int(casted)
				} else {
					errs.Errors = append(errs.Errors, fmt.Sprintf(`error decoding %T into [%d].key_version: strconv.ParseInt: parsing "%s": invalid syntax`, v, i, v))
				}
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].key_version' expected type 'int', got unconvertible type '%T'", i, item["key_version"]))
			}
		}
	}

	if len(errs.Errors) > 0 {
		return &errs
	}

	return nil
}

func (b *backend) pathEncryptExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	}, b.GetRandomReader())
	if err != nil {
		return false, err
	}
	if p != nil && b.System().CachingDisabled() {
		p.Unlock()
	}

	return p != nil, nil
}

func (b *backend) pathEncryptWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	var err error
	batchInputRaw := d.Raw["batch_input"]
	var batchInputItems []BatchRequestItem
	if batchInputRaw != nil {
		err = decodeBatchRequestItems(batchInputRaw, &batchInputItems)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch input: %w", err)
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

	batchResponseItems := make([]EncryptBatchResponseItem, len(batchInputItems))
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
	var upserted bool
	var polReq keysutil.PolicyRequest

	if req.Operation == logical.CreateOperation {
		convergent := d.Get("convergent_encryption").(bool)
		if convergent && !contextSet {
			return logical.ErrorResponse("convergent encryption requires derivation to be enabled, so context is required"), nil
		}

		polReq = keysutil.PolicyRequest{
			Upsert:     true,
			Storage:    req.Storage,
			Name:       name,
			Derived:    contextSet,
			Convergent: convergent,
		}

		keyType := d.Get("type").(string)
		switch keyType {
		case "aes128-gcm96":
			polReq.KeyType = keysutil.KeyType_AES128_GCM96
		case "aes256-gcm96":
			polReq.KeyType = keysutil.KeyType_AES256_GCM96
		case "chacha20-poly1305":
			polReq.KeyType = keysutil.KeyType_ChaCha20_Poly1305
		case "ecdsa-p256", "ecdsa-p384", "ecdsa-p521":
			return logical.ErrorResponse(fmt.Sprintf("key type %v not supported for this operation", keyType)), logical.ErrInvalidRequest
		default:
			return logical.ErrorResponse(fmt.Sprintf("unknown key type %v", keyType)), logical.ErrInvalidRequest
		}
	} else {
		polReq = keysutil.PolicyRequest{
			Storage: req.Storage,
			Name:    name,
		}
	}

	p, upserted, err = b.GetPolicy(ctx, polReq, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
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
			batchResponseItems[i].Error = err.Error()
			continue
		}

		if ciphertext == "" {
			batchResponseItems[i].Error = fmt.Errorf("empty ciphertext returned for input item %d", i).Error()
			continue
		}

		keyVersion := item.KeyVersion
		if keyVersion == 0 {
			keyVersion = p.LatestVersion
		}

		batchResponseItems[i].Ciphertext = ciphertext
		batchResponseItems[i].KeyVersion = keyVersion
	}

	resp := &logical.Response{}
	if batchInputRaw != nil {
		resp.Data = map[string]interface{}{
			"batch_results": batchResponseItems,
		}
	} else {
		if batchResponseItems[0].Error != "" {
			p.Unlock()
			return logical.ErrorResponse(batchResponseItems[0].Error), logical.ErrInvalidRequest
		}

		resp.Data = map[string]interface{}{
			"ciphertext":  batchResponseItems[0].Ciphertext,
			"key_version": batchResponseItems[0].KeyVersion,
		}
	}

	if req.Operation == logical.CreateOperation && !upserted {
		resp.AddWarning("Attempted creation of the key during the encrypt operation, but it was created beforehand")
	}

	p.Unlock()
	return resp, nil
}

const pathEncryptHelpSyn = `Encrypt a plaintext value or a batch of plaintext
blocks using a named key`

const pathEncryptHelpDesc = `
This path uses the named key from the request path to encrypt a user provided
plaintext or a batch of plaintext blocks. The plaintext must be base64 encoded.
`
