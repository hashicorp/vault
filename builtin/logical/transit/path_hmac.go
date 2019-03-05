package transit

import (
	"context"
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

// BatchRequestHMACItem represents a request item for batch processing.
// A map type allows us to distinguish between empty and missing values.
type batchRequestHMACItem map[string]string

// BatchResponseItem represents a response item for batch processing
type batchResponseHMACItem struct {
	// HMAC for the input present in the corresponding batch request item
	HMAC string `json:"hmac,omitempty" mapstructure:"hmac"`

	// Valid indicates whether signature matches the signature derived from the input string
	Valid bool `json:"valid,omitempty" mapstructure:"valid"`

	// Error, if set represents a failure encountered while encrypting a
	// corresponding batch request item
	Error string `json:"error,omitempty" mapstructure:"error"`

	// The return paths in some cases are (nil, err) and others
	// (logical.ErrorResponse(..),nil), and others (logical.ErrorResponse(..),err).
	// For batch processing to successfully mimic previous handling for simple 'input',
	// both output values are needed - though 'err' should never be serialized.
	err error
}

func (b *backend) pathHMAC() *framework.Path {
	return &framework.Path{
		Pattern: "hmac/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("urlalgorithm"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The key to use for the HMAC function",
			},

			"input": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The base64-encoded input data",
			},

			"algorithm": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "sha2-256",
				Description: `Algorithm to use (POST body parameter). Valid values are:

* sha2-224
* sha2-256
* sha2-384
* sha2-512

Defaults to "sha2-256".`,
			},

			"urlalgorithm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Algorithm to use (POST URL parameter)`,
			},

			"key_version": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `The version of the key to use for generating the HMAC.
Must be 0 (for latest) or a value greater than or equal
to the min_encryption_version configured on the key.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathHMACWrite,
		},

		HelpSynopsis:    pathHMACHelpSyn,
		HelpDescription: pathHMACHelpDesc,
	}
}

func (b *backend) pathHMACWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	ver := d.Get("key_version").(int)

	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}

	// Get the policy
	p, _, err := b.lm.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
	}

	switch {
	case ver == 0:
		// Allowed, will use latest; set explicitly here to ensure the string
		// is generated properly
		ver = p.LatestVersion
	case ver == p.LatestVersion:
		// Allowed
	case p.MinEncryptionVersion > 0 && ver < p.MinEncryptionVersion:
		p.Unlock()
		return logical.ErrorResponse("cannot generate HMAC: version is too old (disallowed by policy)"), logical.ErrInvalidRequest
	}

	key, err := p.HMACKey(ver)
	if err != nil {
		p.Unlock()
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if key == nil {
		p.Unlock()
		return nil, fmt.Errorf("HMAC key value could not be computed")
	}

	hashAlgorithm, ok := keysutil.HashTypeMap[algorithm]
	if !ok {
		p.Unlock()
		return logical.ErrorResponse("unsupported algorithm %q", hashAlgorithm), nil
	}

	hashAlg := keysutil.HashFuncMap[hashAlgorithm]

	batchInputRaw := d.Raw["batch_input"]
	var batchInputItems []batchRequestHMACItem
	if batchInputRaw != nil {
		err = mapstructure.Decode(batchInputRaw, &batchInputItems)
		if err != nil {
			p.Unlock()
			return nil, errwrap.Wrapf("failed to parse batch input: {{err}}", err)
		}

		if len(batchInputItems) == 0 {
			p.Unlock()
			return logical.ErrorResponse("missing batch input to process"), logical.ErrInvalidRequest
		}
	} else {
		valueRaw, ok := d.GetOk("input")
		if !ok {
			p.Unlock()
			return logical.ErrorResponse("missing input for HMAC"), logical.ErrInvalidRequest
		}

		batchInputItems = make([]batchRequestHMACItem, 1)
		batchInputItems[0] = batchRequestHMACItem{
			"input": valueRaw.(string),
		}
	}

	response := make([]batchResponseHMACItem, len(batchInputItems))

	for i, item := range batchInputItems {
		rawInput, ok := item["input"]
		if !ok {
			response[i].Error = "missing input for HMAC"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		input, err := base64.StdEncoding.DecodeString(rawInput)
		if err != nil {
			response[i].Error = fmt.Sprintf("unable to decode input as base64: %s", err)
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		var hf = hmac.New(hashAlg, key)
		hf.Write(input)
		retBytes := hf.Sum(nil)

		retStr := base64.StdEncoding.EncodeToString(retBytes)
		retStr = fmt.Sprintf("vault:v%s:%s", strconv.Itoa(ver), retStr)
		response[i].HMAC = retStr
	}

	p.Unlock()

	// Generate the response
	resp := &logical.Response{}
	if batchInputRaw != nil {
		resp.Data = map[string]interface{}{
			"batch_results": response,
		}
	} else {
		if response[0].Error != "" || response[0].err != nil {
			if response[0].Error != "" {
				return logical.ErrorResponse(response[0].Error), response[0].err
			} else {
				return nil, response[0].err
			}
		}
		resp.Data = map[string]interface{}{
			"hmac": response[0].HMAC,
		}
	}

	return resp, nil
}

func (b *backend) pathHMACVerify(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	name := d.Get("name").(string)
	algorithm := d.Get("urlalgorithm").(string)
	if algorithm == "" {
		algorithm = d.Get("algorithm").(string)
	}

	// Get the policy
	p, _, err := b.lm.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("encryption key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
	}

	hashAlgorithm, ok := keysutil.HashTypeMap[algorithm]
	if !ok {
		p.Unlock()
		return logical.ErrorResponse("unsupported algorithm %q", hashAlgorithm), nil
	}

	hashAlg := keysutil.HashFuncMap[hashAlgorithm]

	batchInputRaw := d.Raw["batch_input"]
	var batchInputItems []batchRequestHMACItem
	if batchInputRaw != nil {
		err := mapstructure.Decode(batchInputRaw, &batchInputItems)
		if err != nil {
			p.Unlock()
			return nil, errwrap.Wrapf("failed to parse batch input: {{err}}", err)
		}

		if len(batchInputItems) == 0 {
			p.Unlock()
			return logical.ErrorResponse("missing batch input to process"), logical.ErrInvalidRequest
		}
	} else {
		// use empty string if input is missing - not an error
		inputB64 := d.Get("input").(string)
		hmac := d.Get("hmac").(string)

		batchInputItems = make([]batchRequestHMACItem, 1)
		batchInputItems[0] = batchRequestHMACItem{
			"input": inputB64,
			"hmac":  hmac,
		}
	}

	response := make([]batchResponseHMACItem, len(batchInputItems))

	for i, item := range batchInputItems {
		rawInput, ok := item["input"]
		if !ok {
			response[i].Error = "missing input"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		input, err := base64.StdEncoding.DecodeString(rawInput)
		if err != nil {
			response[i].Error = fmt.Sprintf("unable to decode input as base64: %s", err)
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		verificationHMAC, ok := item["hmac"]
		if !ok {
			response[i].Error = "missing hmac"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		// Verify the prefix
		if !strings.HasPrefix(verificationHMAC, "vault:v") {
			response[i].Error = "invalid HMAC to verify: no prefix"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		splitVerificationHMAC := strings.SplitN(strings.TrimPrefix(verificationHMAC, "vault:v"), ":", 2)
		if len(splitVerificationHMAC) != 2 {
			response[i].Error = "invalid HMAC: wrong number of fields"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		ver, err := strconv.Atoi(splitVerificationHMAC[0])
		if err != nil {
			response[i].Error = "invalid HMAC: version number could not be decoded"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		verBytes, err := base64.StdEncoding.DecodeString(splitVerificationHMAC[1])
		if err != nil {
			response[i].Error = fmt.Sprintf("unable to decode verification HMAC as base64: %s", err)
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		if ver > p.LatestVersion {
			response[i].Error = "invalid HMAC: version is too new"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		if p.MinDecryptionVersion > 0 && ver < p.MinDecryptionVersion {
			response[i].Error = "cannot verify HMAC: version is too old (disallowed by policy)"
			response[i].err = logical.ErrInvalidRequest
			continue
		}

		key, err := p.HMACKey(ver)
		if err != nil {
			response[i].Error = err.Error()
			response[i].err = logical.ErrInvalidRequest
			continue
		}
		if key == nil {
			response[i].Error = ""
			response[i].err = fmt.Errorf("HMAC key value could not be computed")
			continue
		}

		var hf = hmac.New(hashAlg, key)
		hf.Write(input)
		retBytes := hf.Sum(nil)
		response[i].Valid = hmac.Equal(retBytes, verBytes)
	}

	p.Unlock()

	// Generate the response
	resp := &logical.Response{}
	if batchInputRaw != nil {
		resp.Data = map[string]interface{}{
			"batch_results": response,
		}
	} else {
		if response[0].Error != "" || response[0].err != nil {
			if response[0].Error != "" {
				return logical.ErrorResponse(response[0].Error), response[0].err
			} else {
				return nil, response[0].err
			}
		}
		resp.Data = map[string]interface{}{
			"valid": response[0].Valid,
		}
	}

	return resp, nil
}

const pathHMACHelpSyn = `Generate an HMAC for input data using the named key`

const pathHMACHelpDesc = `
Generates an HMAC sum of the given algorithm and key against the given input data.
`
