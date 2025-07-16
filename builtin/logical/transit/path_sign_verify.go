// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

// policySignArgs are the arguments required to pass through to the SDK policy.SignWithOptions
type policySignArgs struct {
	keyVersion int
	keyContext []byte
	input      []byte
	options    keysutil.SigningOptions
}

type policyVerifyArgs struct {
	keyContext []byte
	input      []byte
	options    keysutil.SigningOptions
	sig        string
}

type commonSignVerifyApiArgs struct {
	keyName       string
	hashAlgorithm keysutil.HashType
	marshaling    keysutil.MarshalingType
	prehashed     bool
	sigAlgorithm  string
	saltLength    int
}

// signApiArgs represents the input arguments that apply to all members of the batch
type signApiArgs struct {
	commonSignVerifyApiArgs
	keyVersion int
}

// verifyApiArgs represents the input arguments that apply to all members of the batch
type verifyApiArgs struct {
	commonSignVerifyApiArgs
}

// BatchRequestSignItem represents a request item for batch processing.
// A map type allows us to distinguish between empty and missing values.
type batchRequestSignItem map[string]string

// BatchResponseSignItem represents a response item for batch processing
type batchResponseSignItem struct {
	// signature for the input present in the corresponding batch
	// request item
	Signature string `json:"signature,omitempty" mapstructure:"signature"`

	// The key version to be used for signing
	KeyVersion int `json:"key_version" mapstructure:"key_version"`

	PublicKey []byte `json:"publickey,omitempty" mapstructure:"publickey"`

	// Error, if set represents a failure encountered while signing a
	// corresponding batch request item
	Error string `json:"error,omitempty" mapstructure:"error"`

	// The return paths through WriteSign in some cases are (nil, err) and others
	// (logical.ErrorResponse(..),nil), and others (logical.ErrorResponse(..),err).
	// For batch processing to successfully mimic previous handling for simple 'input',
	// both output values are needed - though 'err' should never be serialized.
	err error

	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference" mapstructure:"reference"`
}

// BatchRequestVerifyItem represents a request item for batch processing.
// A map type allows us to distinguish between empty and missing values.
type batchRequestVerifyItem map[string]interface{}

// BatchResponseVerifyItem represents a response item for batch processing
type batchResponseVerifyItem struct {
	// Valid indicates whether signature matches the signature derived from the input string
	Valid bool `json:"valid" mapstructure:"valid"`

	// Error, if set represents a failure encountered while verifying a
	// corresponding batch request item
	Error string `json:"error,omitempty" mapstructure:"error"`

	// The return paths through WriteSign in some cases are (nil, err) and others
	// (logical.ErrorResponse(..),nil), and others (logical.ErrorResponse(..),err).
	// For batch processing to successfully mimic previous handling for simple 'input',
	// both output values are needed - though 'err' should never be serialized.
	err error

	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference" mapstructure:"reference"`
}

const defaultHashAlgorithm = "sha2-256"

func (b *backend) pathSign() *framework.Path {
	return &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("urlalgorithm"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "sign",
			OperationSuffix: "|with-algorithm",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The key to use",
			},

			"input": {
				Type:        framework.TypeString,
				Description: "The base64-encoded input data",
			},

			"context": {
				Type: framework.TypeString,
				Description: `Base64 encoded context for key derivation. Required if key
derivation is enabled; currently only available with ed25519 keys.`,
			},

			"signature_context": {
				Type: framework.TypeString,
				Description: `Base64 encoded context for Ed25519ph and Ed25519ctx signatures.
Currently only available with Ed25519 keys. (Enterprise Only)`,
			},

			"hash_algorithm": {
				Type:    framework.TypeString,
				Default: defaultHashAlgorithm,
				Description: `Hash algorithm to use (POST body parameter). Valid values are:

* sha1
* sha2-224
* sha2-256
* sha2-384
* sha2-512
* sha3-224
* sha3-256
* sha3-384
* sha3-512
* none

Defaults to "sha2-256". Not valid for all key types,
including ed25519. Using none requires setting prehashed=true and
signature_algorithm=pkcs1v15, yielding a PKCSv1_5_NoOID instead of
the usual PKCSv1_5_DERnull signature.`,
			},

			"algorithm": {
				Type:        framework.TypeString,
				Default:     defaultHashAlgorithm,
				Description: `Deprecated: use "hash_algorithm" instead.`,
			},

			"urlalgorithm": {
				Type:        framework.TypeString,
				Description: `Hash algorithm to use (POST URL parameter)`,
			},

			"key_version": {
				Type: framework.TypeInt,
				Description: `The version of the key to use for signing.
Must be 0 (for latest) or a value greater than or equal
to the min_encryption_version configured on the key.`,
			},

			"prehashed": {
				Type:        framework.TypeBool,
				Description: `Set to 'true' when the input is already hashed. If the key type is 'rsa-2048', 'rsa-3072' or 'rsa-4096', then the algorithm used to hash the input should be indicated by the 'algorithm' parameter.`,
			},

			"signature_algorithm": {
				Type: framework.TypeString,
				Description: `The signature algorithm to use for signing. Currently only applies to RSA key types.
Options are 'pss' or 'pkcs1v15'. Defaults to 'pss'`,
			},

			"marshaling_algorithm": {
				Type:        framework.TypeString,
				Default:     "asn1",
				Description: `The method by which to marshal the signature. The default is 'asn1' which is used by openssl and X.509. It can also be set to 'jws' which is used for JWT signatures; setting it to this will also cause the encoding of the signature to be url-safe base64 instead of using standard base64 encoding. Currently only valid for ECDSA P-256 key types".`,
			},

			"salt_length": {
				Type:    framework.TypeString,
				Default: "auto",
				Description: `The salt length used to sign. Currently only applies to the RSA PSS signature scheme.
Options are 'auto' (the default used by Golang, causing the salt to be as large as possible when signing), 'hash' (causes the salt length to equal the length of the hash used in the signature), or an integer between the minimum and the maximum permissible salt lengths for the given RSA key size. Defaults to 'auto'.`,
			},

			"batch_input": {
				Type: framework.TypeSlice,
				Description: `Specifies a list of items for processing. When this parameter is set,
any supplied 'input' or 'context' parameters will be ignored. Responses are returned in the
'batch_results' array component of the 'data' element of the response. Any batch output will
preserve the order of the batch input`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSignWrite,
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}
}

func (b *backend) pathVerify() *framework.Path {
	return &framework.Path{
		Pattern: "verify/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("urlalgorithm"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "verify",
			OperationSuffix: "|with-algorithm",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The key to use",
			},

			"context": {
				Type: framework.TypeString,
				Description: `Base64 encoded context for key derivation. Required if key
derivation is enabled; currently only available with ed25519 keys.`,
			},

			"signature_context": {
				Type: framework.TypeString,
				Description: `Base64 encoded context for Ed25519ph and Ed25519ctx signatures. 
Currently only available with Ed25519 keys. (Enterprise Only)`,
			},

			"signature": {
				Type:        framework.TypeString,
				Description: "The signature, including vault header/key version",
			},

			"hmac": {
				Type:        framework.TypeString,
				Description: "The HMAC, including vault header/key version",
			},

			"cmac": {
				Type:        framework.TypeString,
				Description: "The CMAC, including vault header/key version (Enterprise only)",
			},

			"input": {
				Type:        framework.TypeString,
				Description: "The base64-encoded input data to verify",
			},

			"urlalgorithm": {
				Type:        framework.TypeString,
				Description: `Hash algorithm to use (POST URL parameter)`,
			},

			"mac_length": {
				Type:        framework.TypeInt,
				Description: `MAC length to use (POST body parameter). Valid values are:`,
			},

			"hash_algorithm": {
				Type:    framework.TypeString,
				Default: defaultHashAlgorithm,
				Description: `Hash algorithm to use (POST body parameter). Valid values are:

* sha1
* sha2-224
* sha2-256
* sha2-384
* sha2-512
* sha3-224
* sha3-256
* sha3-384
* sha3-512
* none

Defaults to "sha2-256". Not valid for all key types. See note about
none on signing path.`,
			},

			"algorithm": {
				Type:        framework.TypeString,
				Default:     defaultHashAlgorithm,
				Description: `Deprecated: use "hash_algorithm" instead.`,
			},

			"prehashed": {
				Type:        framework.TypeBool,
				Description: `Set to 'true' when the input is already hashed. If the key type is 'rsa-2048', 'rsa-3072' or 'rsa-4096', then the algorithm used to hash the input should be indicated by the 'algorithm' parameter.`,
			},

			"signature_algorithm": {
				Type: framework.TypeString,
				Description: `The signature algorithm to use for signature verification. Currently only applies to RSA key types.
Options are 'pss' or 'pkcs1v15'. Defaults to 'pss'`,
			},

			"marshaling_algorithm": {
				Type:        framework.TypeString,
				Default:     "asn1",
				Description: `The method by which to unmarshal the signature when verifying. The default is 'asn1' which is used by openssl and X.509; can also be set to 'jws' which is used for JWT signatures in which case the signature is also expected to be url-safe base64 encoding instead of standard base64 encoding. Currently only valid for ECDSA P-256 key types".`,
			},

			"salt_length": {
				Type:    framework.TypeString,
				Default: "auto",
				Description: `The salt length used to sign. Currently only applies to the RSA PSS signature scheme.
Options are 'auto' (the default used by Golang, causing the salt to be as large as possible when signing), 'hash' (causes the salt length to equal the length of the hash used in the signature), or an integer between the minimum and the maximum permissible salt lengths for the given RSA key size. Defaults to 'auto'.`,
			},

			"batch_input": {
				Type: framework.TypeSlice,
				Description: `Specifies a list of items for processing. When this parameter is set,
any supplied  'input', 'hmac', 'cmac' or 'signature' parameters will be ignored. Responses are returned in the
'batch_results' array component of the 'data' element of the response. Any batch output will
preserve the order of the batch input`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathVerifyWrite,
		},

		HelpSynopsis:    pathVerifyHelpSyn,
		HelpDescription: pathVerifyHelpDesc,
	}
}

func getSaltLength(d *framework.FieldData) (int, error) {
	rawSaltLength, ok := d.GetOk("salt_length")
	// This should only happen when something is wrong with the schema,
	// so this is a reasonable default.
	if !ok {
		return rsa.PSSSaltLengthAuto, nil
	}

	rawSaltLengthStr := rawSaltLength.(string)
	lowerSaltLengthStr := strings.ToLower(rawSaltLengthStr)
	switch lowerSaltLengthStr {
	case "auto":
		return rsa.PSSSaltLengthAuto, nil
	case "hash":
		return rsa.PSSSaltLengthEqualsHash, nil
	default:
		saltLengthInt, err := strconv.Atoi(lowerSaltLengthStr)
		if err != nil {
			return rsa.PSSSaltLengthEqualsHash - 1, fmt.Errorf("salt length neither 'auto', 'hash', nor an int: %s", rawSaltLength)
		}
		if saltLengthInt < rsa.PSSSaltLengthEqualsHash {
			return rsa.PSSSaltLengthEqualsHash - 1, fmt.Errorf("salt length is invalid: %d", saltLengthInt)
		}
		return saltLengthInt, nil
	}
}

func (b *backend) pathSignWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Fetch the top-level arguments for the sign api. These do not include
	// the values that are within the batch input parameters, only those that
	// apply globally.
	apiArgs, err := getSignApiArgs(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Get the policy
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    apiArgs.keyName,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("signing key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
	}
	defer p.Unlock()

	if err := validateCommonSignVerifyApiArgs(p, apiArgs.commonSignVerifyApiArgs); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	batchInputRaw := d.Raw["batch_input"]
	var batchInputItems []batchRequestSignItem
	if batchInputRaw != nil {
		err = mapstructure.Decode(batchInputRaw, &batchInputItems)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch input: %w", err)
		}

		if len(batchInputItems) == 0 {
			return logical.ErrorResponse("missing batch input to process"), logical.ErrInvalidRequest
		}
	} else {
		// use empty string if input is missing - not an error
		batchInputItems = make([]batchRequestSignItem, 1)
		batchInputItems[0] = batchRequestSignItem{
			"input":   d.Get("input").(string),
			"context": d.Get("context").(string),
		}

		// Only defined within the ENT schema, so this is useless on CE.
		if sigContext, ok := d.GetOk("signature_context"); ok {
			batchInputItems[0]["signature_context"] = sigContext.(string)
		}
	}

	response := make([]batchResponseSignItem, len(batchInputItems))
	for i, item := range batchInputItems {
		psa, err := b.getPolicySignArgs(ctx, p, apiArgs, item)
		if err != nil {
			response[i].Error = err.Error()
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		sig, err := p.SignWithOptions(psa.keyVersion, psa.keyContext, psa.input, &psa.options)
		if err != nil {
			if batchInputRaw != nil {
				response[i].Error = err.Error()
			}
			response[i].err = err
		} else if sig == nil {
			response[i].err = fmt.Errorf("signature could not be computed")
		} else {
			keyVersion := apiArgs.keyVersion
			if keyVersion == 0 {
				keyVersion = p.LatestVersion
			}

			response[i].Signature = sig.Signature
			response[i].PublicKey = sig.PublicKey
			response[i].KeyVersion = keyVersion
		}
	}

	// Generate the response
	resp := &logical.Response{}
	if batchInputRaw != nil {
		// Copy the references
		for i := range batchInputItems {
			response[i].Reference = batchInputItems[i]["reference"]
		}
		resp.Data = map[string]interface{}{
			"batch_results": response,
		}
	} else {
		if response[0].Error != "" || response[0].err != nil {
			if response[0].Error != "" {
				return logical.ErrorResponse(response[0].Error), response[0].err
			}

			return nil, response[0].err
		}

		resp.Data = map[string]interface{}{
			"signature":   response[0].Signature,
			"key_version": response[0].KeyVersion,
		}

		if len(response[0].PublicKey) > 0 {
			resp.Data["public_key"] = response[0].PublicKey
		}
	}

	return resp, nil
}

func (b *backend) getPolicySignArgs(ctx context.Context, p *keysutil.Policy, args signApiArgs, item batchRequestSignItem) (policySignArgs, error) {
	rawInput, ok := item["input"]
	if !ok {
		return policySignArgs{}, fmt.Errorf("missing input")
	}

	input, err := base64.StdEncoding.DecodeString(rawInput)
	if err != nil {
		return policySignArgs{}, fmt.Errorf("unable to decode input: %s", err)
	}

	if p.Type.HashSignatureInput() && !args.prehashed {
		hf := keysutil.HashFuncMap[args.hashAlgorithm]()
		if hf != nil {
			hf.Write(input)
			input = hf.Sum(nil)
		}
	}

	keyContext, err := decodeBase64Arg(item["context"])
	if err != nil {
		return policySignArgs{}, fmt.Errorf("failed to base64-decode context")
	}

	psa := policySignArgs{
		keyVersion: args.keyVersion,
		keyContext: keyContext,
		input:      input,
		options: keysutil.SigningOptions{
			HashAlgorithm: args.hashAlgorithm,
			Marshaling:    args.marshaling,
			SaltLength:    args.saltLength,
			SigAlgorithm:  args.sigAlgorithm,
		},
	}

	if err := b.populateEntPolicySigningOptions(ctx, p, args, item, &psa); err != nil {
		return policySignArgs{}, fmt.Errorf("failed to parse batch item: %s", err)
	}
	return psa, nil
}

func validateCommonSignVerifyApiArgs(p *keysutil.Policy, apiArgs commonSignVerifyApiArgs) error {
	if !p.Type.SigningSupported() {
		return fmt.Errorf("key type %v does not support signing", p.Type)
	}

	// Perform Vault version specific checks (CE vs ENT)
	return validateSignApiArgsVersionSpecific(p, apiArgs)
}

func getSignApiArgs(d *framework.FieldData) (signApiArgs, error) {
	keyVersion := d.Get("key_version").(int)
	commonArgs, err := getCommonSignVerifyApiArgs(d)
	if err != nil {
		return signApiArgs{}, err
	}

	return signApiArgs{
		commonSignVerifyApiArgs: commonArgs,
		keyVersion:              keyVersion,
	}, nil
}

func getCommonSignVerifyApiArgs(d *framework.FieldData) (commonSignVerifyApiArgs, error) {
	keyName := d.Get("name").(string)
	hashAlgorithmStr, hashAlgorithm, ok := getHashAlgorithmFromArgs(d)
	if !ok {
		return commonSignVerifyApiArgs{}, fmt.Errorf("invalid hash algorithm %q", hashAlgorithmStr)
	}

	marshalingStr := d.Get("marshaling_algorithm").(string)
	marshaling, ok := keysutil.MarshalingTypeMap[marshalingStr]
	if !ok {
		return commonSignVerifyApiArgs{}, fmt.Errorf("invalid marshaling type %q", marshalingStr)
	}

	prehashed := d.Get("prehashed").(bool)
	sigAlgorithm := d.Get("signature_algorithm").(string)
	saltLength, err := getSaltLength(d)
	if err != nil {
		return commonSignVerifyApiArgs{}, err
	}

	return commonSignVerifyApiArgs{
		keyName:       keyName,
		hashAlgorithm: hashAlgorithm,
		marshaling:    marshaling,
		prehashed:     prehashed,
		sigAlgorithm:  sigAlgorithm,
		saltLength:    saltLength,
	}, nil
}

func getHashAlgorithmFromArgs(d *framework.FieldData) (string, keysutil.HashType, bool) {
	hashAlgorithmStr := d.Get("urlalgorithm").(string)
	if hashAlgorithmStr == "" {
		hashAlgorithmStr = d.Get("hash_algorithm").(string)
		if hashAlgorithmStr == "" {
			hashAlgorithmStr = d.Get("algorithm").(string)
			if hashAlgorithmStr == "" {
				hashAlgorithmStr = defaultHashAlgorithm
			}
		}
	}

	hashAlgorithm, ok := keysutil.HashTypeMap[hashAlgorithmStr]
	return hashAlgorithmStr, hashAlgorithm, ok
}

func decodeBase64Arg(fieldVal interface{}) ([]byte, error) {
	parsedStr, err := parseutil.ParseString(fieldVal)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(parsedStr)
}

func (b *backend) pathVerifyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	batchInputItems, isBatchInput, err := getVerifyBatchItems(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// For simplicity, 'signature' and 'hmac' cannot be mixed across batch_input elements.
	// If one batch_input item is 'signature', they all must be 'signature'.
	// If one batch_input item is 'hmac', they all must be 'hmac'.
	sigFound := false
	hmacFound := false
	cmacFound := false
	missing := false
	for _, v := range batchInputItems {
		if _, ok := v["signature"]; ok {
			sigFound = true
		} else if _, ok := v["hmac"]; ok {
			hmacFound = true
		} else if _, ok := v["cmac"]; ok {
			cmacFound = true
		} else {
			missing = true
		}
	}
	optionsSet := numBooleansTrue(sigFound, hmacFound, cmacFound)

	switch {
	case !isBatchInput && optionsSet > 1:
		return logical.ErrorResponse("provide one of 'signature', 'hmac' or 'cmac'"), logical.ErrInvalidRequest

	case !isBatchInput && optionsSet == 0:
		return logical.ErrorResponse("missing one of 'signature', 'hmac' or 'cmac' arguments to verify"), logical.ErrInvalidRequest

	case optionsSet > 1:
		return logical.ErrorResponse("elements of batch_input must all provide either 'signature', 'hmac' or 'cmac'"), logical.ErrInvalidRequest

	case missing && sigFound:
		return logical.ErrorResponse("some elements of batch_input are missing 'signature'"), logical.ErrInvalidRequest

	case missing && hmacFound:
		return logical.ErrorResponse("some elements of batch_input are missing 'hmac'"), logical.ErrInvalidRequest

	case missing && cmacFound:
		return logical.ErrorResponse("some elements of batch_input are missing 'cmac'"), logical.ErrInvalidRequest

	case optionsSet == 0:
		return logical.ErrorResponse("no batch_input elements have 'signature', 'hmac' or 'cmac'"), logical.ErrInvalidRequest

	case hmacFound:
		return b.pathHMACVerify(ctx, req, d)

	case cmacFound:
		return b.pathCMACVerify(ctx, req, d)
	}

	apiArgs, err := getVerifyApiArgs(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Get the policy
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    apiArgs.keyName,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("signature verification key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
	}
	defer p.Unlock()

	if err := validateCommonSignVerifyApiArgs(p, apiArgs.commonSignVerifyApiArgs); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	response := make([]batchResponseVerifyItem, len(batchInputItems))

	for i, item := range batchInputItems {
		pva, err := b.getPolicyVerifyArgs(ctx, p, apiArgs, item)
		if err != nil {
			response[i].Error = fmt.Sprintf("failed to parse item: %s", err)
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		valid, err := p.VerifySignatureWithOptions(pva.keyContext, pva.input, pva.sig, &pva.options)
		if err != nil {
			switch err.(type) {
			case errutil.UserError:
				response[i].Error = err.Error()
				response[i].err = logical.ErrInvalidRequest
			default:
				if isBatchInput {
					response[i].Error = err.Error()
				}
				response[i].err = err
			}
		} else {
			response[i].Valid = valid
		}
	}

	// Generate the response
	resp := &logical.Response{}
	if isBatchInput {
		// Copy the references
		for i := range batchInputItems {
			if ref, err := parseutil.ParseString(batchInputItems[i]["reference"]); err == nil {
				response[i].Reference = ref
			}
		}
		resp.Data = map[string]interface{}{
			"batch_results": response,
		}
	} else {
		if response[0].Error != "" || response[0].err != nil {
			if response[0].Error != "" {
				return logical.ErrorResponse(response[0].Error), response[0].err
			}
			return nil, response[0].err
		}
		resp.Data = map[string]interface{}{
			"valid": response[0].Valid,
		}
	}

	return resp, nil
}

func getVerifyApiArgs(d *framework.FieldData) (verifyApiArgs, error) {
	commonArgs, err := getCommonSignVerifyApiArgs(d)
	if err != nil {
		return verifyApiArgs{}, err
	}

	return verifyApiArgs{
		commonSignVerifyApiArgs: commonArgs,
	}, nil
}

func getVerifyBatchItems(d *framework.FieldData) ([]batchRequestVerifyItem, bool, error) {
	if batchInputRaw, ok := d.Raw["batch_input"]; ok && batchInputRaw != nil {
		var batchInputItems []batchRequestVerifyItem
		err := mapstructure.Decode(batchInputRaw, &batchInputItems)
		if err != nil {
			return nil, false, fmt.Errorf("failed to parse batch input: %w", err)
		}

		if len(batchInputItems) == 0 {
			return nil, false, fmt.Errorf("missing batch input to process")
		}

		return batchInputItems, true, nil
	}

	// use empty string if input is missing - not an error
	item := batchRequestVerifyItem{
		"input": d.Get("input").(string),
	}
	if sig, ok := d.GetOk("signature"); ok {
		item["signature"] = sig.(string)
	}
	if hmac, ok := d.GetOk("hmac"); ok {
		item["hmac"] = hmac.(string)
	}
	if cmac, ok := d.GetOk("cmac"); ok {
		item["cmac"] = cmac.(string)
	}

	item["context"] = d.Get("context").(string)

	if sigContext, ok := d.GetOk("signature_context"); ok {
		item["signature_context"] = sigContext.(string)
	}

	return []batchRequestVerifyItem{item}, false, nil
}

func (b *backend) getPolicyVerifyArgs(ctx context.Context, p *keysutil.Policy, apiArgs verifyApiArgs, item batchRequestVerifyItem) (policyVerifyArgs, error) {
	input, err := decodeBase64Arg(item["input"])
	if err != nil {
		return policyVerifyArgs{}, fmt.Errorf("failed to parse input: %s", err)
	}

	sigRaw, ok := item["signature"]
	if !ok {
		return policyVerifyArgs{}, fmt.Errorf("missing signature")
	}
	sig, err := parseutil.ParseString(sigRaw)
	if err != nil {
		return policyVerifyArgs{}, fmt.Errorf("failed to parse signature as a string: %s", err)
	}

	if p.Type.HashSignatureInput() && !apiArgs.prehashed {
		hf := keysutil.HashFuncMap[apiArgs.hashAlgorithm]()
		if hf != nil {
			hf.Write(input)
			input = hf.Sum(nil)
		}
	}

	keyContext, err := decodeBase64Arg(item["context"])
	if err != nil {
		return policyVerifyArgs{}, fmt.Errorf("failed to parse context: %s", err)
	}

	vsa := policyVerifyArgs{
		keyContext: keyContext,
		input:      input,
		sig:        sig,
		options: keysutil.SigningOptions{
			HashAlgorithm: apiArgs.hashAlgorithm,
			Marshaling:    apiArgs.marshaling,
			SaltLength:    apiArgs.saltLength,
			SigAlgorithm:  apiArgs.sigAlgorithm,
		},
	}

	if err := b.populateEntPolicyVerifyOptions(ctx, p, apiArgs, item, &vsa); err != nil {
		return policyVerifyArgs{}, fmt.Errorf("failed to parse batch item: %s", err)
	}

	return vsa, nil
}

func numBooleansTrue(bools ...bool) int {
	numSet := 0
	for _, value := range bools {
		if value {
			numSet++
		}
	}
	return numSet
}

const pathSignHelpSyn = `Generate a signature for input data using the named key`

const pathSignHelpDesc = `
Generates a signature of the input data using the named key and the given hash algorithm.
`
const pathVerifyHelpSyn = `Verify a signature or HMAC for input data created using the named key`

const pathVerifyHelpDesc = `
Verifies a signature or HMAC of the input data using the named key and the given hash algorithm.
`
