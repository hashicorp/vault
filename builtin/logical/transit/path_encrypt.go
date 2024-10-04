// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
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

	// PaddingScheme for encryption/decryption
	PaddingScheme string `json:"padding_scheme" structs:"padding_scheme" mapstructure:"padding_scheme"`

	// Nonce to be used when v1 convergent encryption is used
	Nonce string `json:"nonce" structs:"nonce" mapstructure:"nonce"`

	// The key version to be used for encryption
	KeyVersion int `json:"key_version" structs:"key_version" mapstructure:"key_version"`

	// DecodedNonce is the base64 decoded version of Nonce
	DecodedNonce []byte

	// Associated Data for AEAD ciphers
	AssociatedData string `json:"associated_data" struct:"associated_data" mapstructure:"associated_data"`

	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference" structs:"reference" mapstructure:"reference"`
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

	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference"`
}

type AssocDataFactory struct {
	Encoded string
}

func (a AssocDataFactory) GetAssociatedData() ([]byte, error) {
	return base64.StdEncoding.DecodeString(a.Encoded)
}

type ManagedKeyFactory struct {
	managedKeyParams keysutil.ManagedKeyParameters
}

func (m ManagedKeyFactory) GetManagedKeyParameters() keysutil.ManagedKeyParameters {
	return m.managedKeyParams
}

func (b *backend) pathEncrypt() *framework.Path {
	return &framework.Path{
		Pattern: "encrypt/" + framework.GenericNameRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "encrypt",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"plaintext": {
				Type:        framework.TypeString,
				Description: "Base64 encoded plaintext value to be encrypted",
			},

			"padding_scheme": {
				Type: framework.TypeString,
				Description: `The padding scheme to use for decrypt. Currently only applies to RSA key types.
Options are 'oaep' or 'pkcs1v15'. Defaults to 'oaep'`,
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

			"partial_failure_response_code": {
				Type: framework.TypeInt,
				Description: `
Ordinarily, if a batch item fails to encrypt due to a bad input, but other batch items succeed, 
the HTTP response code is 400 (Bad Request).  Some applications may want to treat partial failures differently.
Providing the parameter returns the given response code integer instead of a 400 in this case. If all values fail
HTTP 400 is still returned.`,
			},

			"associated_data": {
				Type: framework.TypeString,
				Description: `
When using an AEAD cipher mode, such as AES-GCM, this parameter allows
passing associated data (AD/AAD) into the encryption function; this data
must be passed on subsequent decryption requests but can be transited in
plaintext. On successful decryption, both the ciphertext and the associated
data are attested not to have been tampered with.
				`,
			},

			"batch_input": {
				Type: framework.TypeSlice,
				Description: `
Specifies a list of items to be encrypted in a single batch. When this parameter
is set, if the parameters 'plaintext', 'context' and 'nonce' are also set, they
will be ignored. Any batch output will preserve the order of the batch input.`,
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

func decodeEncryptBatchRequestItems(src interface{}, dst *[]BatchRequestItem) error {
	return decodeBatchRequestItems(src, true, false, dst)
}

func decodeDecryptBatchRequestItems(src interface{}, dst *[]BatchRequestItem) error {
	return decodeBatchRequestItems(src, false, true, dst)
}

// decodeBatchRequestItems is a fast path alternative to mapstructure.Decode to decode []BatchRequestItem.
// It aims to behave as closely possible to the original mapstructure.Decode and will return the same errors.
// Note, however, that an error will also be returned if one of the required fields is missing.
// https://github.com/hashicorp/vault/pull/8775/files#r437709722
func decodeBatchRequestItems(src interface{}, requirePlaintext bool, requireCiphertext bool, dst *[]BatchRequestItem) error {
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
		} else if requireCiphertext {
			errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].ciphertext' missing ciphertext to decrypt", i))
		}

		if v, has := item["plaintext"]; has {
			if casted, ok := v.(string); ok {
				(*dst)[i].Plaintext = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].plaintext' expected type 'string', got unconvertible type '%T'", i, item["plaintext"]))
			}
		} else if requirePlaintext {
			errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].plaintext' missing plaintext to encrypt", i))
		}
		if v, has := item["padding_scheme"]; has {
			if casted, ok := v.(string); ok {
				(*dst)[i].PaddingScheme = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].padding_scheme' expected type 'string', got unconvertible type '%T'", i, item["padding_scheme"]))
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

		if v, has := item["associated_data"]; has {
			if !reflect.ValueOf(v).IsValid() {
			} else if casted, ok := v.(string); ok {
				(*dst)[i].AssociatedData = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].associated_data' expected type 'string', got unconvertible type '%T'", i, item["associated_data"]))
			}
		}
		if v, has := item["reference"]; has {
			if !reflect.ValueOf(v).IsValid() {
			} else if casted, ok := v.(string); ok {
				(*dst)[i].Reference = casted
			} else {
				errs.Errors = append(errs.Errors, fmt.Sprintf("'[%d].reference' expected type 'string', got unconvertible type '%T'", i, item["reference"]))
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
		err = decodeEncryptBatchRequestItems(batchInputRaw, &batchInputItems)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch input: %w", err)
		}

		if len(batchInputItems) == 0 {
			return logical.ErrorResponse("missing batch input to process"), logical.ErrInvalidRequest
		}
	} else {
		valueRaw, ok, err := d.GetOkErr("plaintext")
		if err != nil {
			return nil, err
		}
		if !ok {
			return logical.ErrorResponse("missing plaintext to encrypt"), logical.ErrInvalidRequest
		}

		batchInputItems = make([]BatchRequestItem, 1)
		batchInputItems[0] = BatchRequestItem{
			Plaintext:      valueRaw.(string),
			Context:        d.Get("context").(string),
			Nonce:          d.Get("nonce").(string),
			KeyVersion:     d.Get("key_version").(int),
			AssociatedData: d.Get("associated_data").(string),
		}
		if psRaw, ok := d.GetOk("padding_scheme"); ok {
			if ps, ok := psRaw.(string); ok {
				batchInputItems[0].PaddingScheme = ps
			} else {
				return logical.ErrorResponse("padding_scheme was not a string"), logical.ErrInvalidRequest
			}
		}
	}

	batchResponseItems := make([]EncryptBatchResponseItem, len(batchInputItems))
	contextSet := len(batchInputItems[0].Context) != 0

	userErrorInBatch := false
	internalErrorInBatch := false

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
			userErrorInBatch = true
			batchResponseItems[i].Error = err.Error()
			continue
		}

		// Decode the context
		if len(item.Context) != 0 {
			batchInputItems[i].DecodedContext, err = base64.StdEncoding.DecodeString(item.Context)
			if err != nil {
				userErrorInBatch = true
				batchResponseItems[i].Error = err.Error()
				continue
			}
		}

		// Decode the nonce
		if len(item.Nonce) != 0 {
			batchInputItems[i].DecodedNonce, err = base64.StdEncoding.DecodeString(item.Nonce)
			if err != nil {
				userErrorInBatch = true
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

		cfg, err := b.readConfigKeys(ctx, req)
		if err != nil {
			return nil, err
		}

		polReq = keysutil.PolicyRequest{
			Upsert:     !cfg.DisableUpsert,
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
		case "rsa-2048":
			polReq.KeyType = keysutil.KeyType_RSA2048
		case "rsa-3072":
			polReq.KeyType = keysutil.KeyType_RSA3072
		case "rsa-4096":
			polReq.KeyType = keysutil.KeyType_RSA4096
		case "ecdsa-p256", "ecdsa-p384", "ecdsa-p521":
			return logical.ErrorResponse(fmt.Sprintf("key type %v not supported for this operation", keyType)), logical.ErrInvalidRequest
		case "managed_key":
			polReq.KeyType = keysutil.KeyType_MANAGED_KEY
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
	defer p.Unlock()

	// Process batch request items. If encryption of any request
	// item fails, respectively mark the error in the response
	// collection and continue to process other items.
	warnAboutNonceUsage := false
	successesInBatch := false
	for i, item := range batchInputItems {
		if batchResponseItems[i].Error != "" {
			userErrorInBatch = true
			continue
		}

		if item.Nonce != "" && !nonceAllowed(p) {
			userErrorInBatch = true
			batchResponseItems[i].Error = ErrNonceNotAllowed.Error()
			continue
		}

		if !warnAboutNonceUsage && shouldWarnAboutNonceUsage(p, item.DecodedNonce) {
			warnAboutNonceUsage = true
		}

		var factories []any
		if item.PaddingScheme != "" {
			paddingScheme, err := parsePaddingSchemeArg(p.Type, item.PaddingScheme)
			if err != nil {
				batchResponseItems[i].Error = fmt.Sprintf("'[%d].padding_scheme' invalid: %s", i, err.Error())
				continue
			}
			factories = append(factories, paddingScheme)
		}
		if item.AssociatedData != "" {
			if !p.Type.AssociatedDataSupported() {
				batchResponseItems[i].Error = fmt.Sprintf("'[%d].associated_data' provided for non-AEAD cipher suite %v", i, p.Type.String())
				continue
			}

			factories = append(factories, AssocDataFactory{item.AssociatedData})
		}

		if p.Type == keysutil.KeyType_MANAGED_KEY {
			managedKeySystemView, ok := b.System().(logical.ManagedKeySystemView)
			if !ok {
				batchResponseItems[i].Error = errors.New("unsupported system view").Error()
			}

			factories = append(factories, ManagedKeyFactory{
				managedKeyParams: keysutil.ManagedKeyParameters{
					ManagedKeySystemView: managedKeySystemView,
					BackendUUID:          b.backendUUID,
					Context:              ctx,
				},
			})
		}

		ciphertext, err := p.EncryptWithFactory(item.KeyVersion, item.DecodedContext, item.DecodedNonce, item.Plaintext, factories...)
		if err != nil {
			switch err.(type) {
			case errutil.InternalError:
				internalErrorInBatch = true
			default:
				userErrorInBatch = true
			}
			batchResponseItems[i].Error = err.Error()
			continue
		}

		if ciphertext == "" {
			userErrorInBatch = true
			batchResponseItems[i].Error = fmt.Sprintf("empty ciphertext returned for input item %d", i)
			continue
		}

		successesInBatch = true
		keyVersion := item.KeyVersion
		if keyVersion == 0 {
			keyVersion = p.LatestVersion
		}

		batchResponseItems[i].Ciphertext = ciphertext
		batchResponseItems[i].KeyVersion = keyVersion
	}

	resp := &logical.Response{}
	if batchInputRaw != nil {
		// Copy the references
		for i := range batchInputItems {
			batchResponseItems[i].Reference = batchInputItems[i].Reference
		}
		resp.Data = map[string]interface{}{
			"batch_results": batchResponseItems,
		}
	} else {
		if batchResponseItems[0].Error != "" {
			if internalErrorInBatch {
				return nil, errutil.InternalError{Err: batchResponseItems[0].Error}
			}

			return logical.ErrorResponse(batchResponseItems[0].Error), logical.ErrInvalidRequest
		}

		resp.Data = map[string]interface{}{
			"ciphertext":  batchResponseItems[0].Ciphertext,
			"key_version": batchResponseItems[0].KeyVersion,
		}
	}

	if constants.IsFIPS() && warnAboutNonceUsage {
		resp.AddWarning("A provided nonce value was used within FIPS mode, this violates FIPS 140 compliance.")
	}

	if req.Operation == logical.CreateOperation && !upserted {
		resp.AddWarning("Attempted creation of the key during the encrypt operation, but it was created beforehand")
	}

	return batchRequestResponse(d, resp, req, successesInBatch, userErrorInBatch, internalErrorInBatch)
}

func nonceAllowed(p *keysutil.Policy) bool {
	var supportedKeyType bool
	switch p.Type {
	case keysutil.KeyType_MANAGED_KEY:
		return true
	case keysutil.KeyType_AES128_GCM96, keysutil.KeyType_AES256_GCM96, keysutil.KeyType_ChaCha20_Poly1305:
		supportedKeyType = true
	default:
		supportedKeyType = false
	}

	if supportedKeyType && p.ConvergentEncryption && p.ConvergentVersion == 1 {
		// We only use the user supplied nonce for v1 convergent encryption keys
		return true
	}

	return false
}

// Depending on the errors in the batch, different status codes should be returned. User errors
// will return a 400 and precede internal errors which return a 500. The reasoning behind this is
// that user errors are non-retryable without making changes to the request, and should be surfaced
// to the user first.
func batchRequestResponse(d *framework.FieldData, resp *logical.Response, req *logical.Request, successesInBatch, userErrorInBatch, internalErrorInBatch bool) (*logical.Response, error) {
	if userErrorInBatch || internalErrorInBatch {
		var code int
		switch {
		case userErrorInBatch:
			code = http.StatusBadRequest
		case internalErrorInBatch:
			code = http.StatusInternalServerError
		}
		if codeRaw, ok := d.GetOk("partial_failure_response_code"); ok && successesInBatch {
			newCode := codeRaw.(int)
			if newCode < 1 || newCode > 599 {
				resp.AddWarning(fmt.Sprintf("invalid HTTP response code override from partial_failure_response_code, reverting to %d", code))
			} else {
				code = newCode
			}
		}
		return logical.RespondWithStatusCode(resp, req, code)
	}

	return resp, nil
}

// shouldWarnAboutNonceUsage attempts to determine if we will use a provided nonce or not. Ideally this
// would be information returned through p.Encrypt but that would require an SDK api change and this is
// transit specific
func shouldWarnAboutNonceUsage(p *keysutil.Policy, userSuppliedNonce []byte) bool {
	if len(userSuppliedNonce) == 0 {
		return false
	}

	var supportedKeyType bool
	switch p.Type {
	case keysutil.KeyType_AES128_GCM96, keysutil.KeyType_AES256_GCM96, keysutil.KeyType_ChaCha20_Poly1305:
		supportedKeyType = true
	default:
		supportedKeyType = false
	}

	if supportedKeyType && p.ConvergentEncryption && p.ConvergentVersion == 1 {
		// We only use the user supplied nonce for v1 convergent encryption keys
		return true
	}

	if supportedKeyType && !p.ConvergentEncryption {
		return true
	}

	return false
}

const pathEncryptHelpSyn = `Encrypt a plaintext value or a batch of plaintext
blocks using a named key`

const pathEncryptHelpDesc = `
This path uses the named key from the request path to encrypt a user provided
plaintext or a batch of plaintext blocks. The plaintext must be base64 encoded.
`
