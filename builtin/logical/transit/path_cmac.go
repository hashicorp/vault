// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/aes"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"hash"
	"strings"

	aesCmac "github.com/aead/cmac/aes"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

// cmacWriteInput captures all the arguments for a cmac write api call
type cmacWriteInput struct {
	IsBatch    bool
	KeyName    string
	KeyVersion int
	Items      []cmacWriteItem
}

// cmacWriteItem represents a single request item to CMAC
type cmacWriteItem struct {
	Input     string `json:"input,omitempty" mapstructure:"input"`
	MacLength int    `json:"mac_length" mapstructure:"mac_length"`
	Reference string `json:"reference,omitempty" mapstructure:"reference,omitempty"`
}

// cmacWriteResponseItem represents an associated response to a request item
type cmacWriteResponseItem struct {
	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference" mapstructure:"reference"`

	// CMAC for the input present in the corresponding batch request item
	CMAC string `json:"cmac,omitempty" mapstructure:"cmac"`

	// Error, if set represents a failure encountered while encrypting a
	// corresponding batch request item
	Error string `json:"error,omitempty" mapstructure:"error"`
}

// cmacVerifyInput captures all arguments related to a CMAC verify API call
type cmacVerifyInput struct {
	IsBatch bool
	KeyName string
	Items   []cmacVerifyItem
}

// cmacVerifyItem captures inputs related to a single work item
type cmacVerifyItem struct {
	Input     string `json:"input,omitempty" mapstructure:"input"`
	Cmac      string `json:"cmac,omitempty" mapstructure:"cmac,omitempty"`
	MacLength int    `json:"mac_length" mapstructure:"mac_length"`
	Reference string `json:"reference,omitempty" mapstructure:"reference,omitempty"`
}

// cmacVerifyResponseItem represents an associated response to a request item
type cmacVerifyResponseItem struct {
	// Valid indicates whether signature matches the signature derived from the input string
	Valid bool `json:"valid,omitempty" mapstructure:"valid"`

	// Error, if set represents a failure encountered while encrypting a
	// corresponding batch request item
	Error string `json:"error,omitempty" mapstructure:"error"`

	// Reference is an arbitrary caller supplied string value that will be placed on the
	// batch response to ease correlation between inputs and outputs
	Reference string `json:"reference" mapstructure:"reference"`
}

func (b *backend) pathCMAC() *framework.Path {
	return &framework.Path{
		Pattern: "cmac/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("url_mac_length"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "generate",
			OperationSuffix: "cmac|cmac-with-mac-length",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The key to use for the CMAC function",
			},

			"input": {
				Type:        framework.TypeString,
				Description: "The base64-encoded input data",
			},

			"mac_length": {
				Type:        framework.TypeInt,
				Description: `MAC length to use (POST body parameter). Valid values are:`,
			},

			"url_mac_length": {
				Type:        framework.TypeString,
				Description: `MAC length to use (POST URL parameter), overrides mac_length`,
			},

			"key_version": {
				Type: framework.TypeInt,
				Description: `The version of the key to use for generating the CMAC.
Must be 0 (for latest) or a value greater than or equal
to the min_encryption_version configured on the key.`,
			},

			"batch_input": {
				Type: framework.TypeSlice,
				Description: `
Specifies a list of items to be processed in a single batch. When this parameter
is set, if the parameter 'input' is also set, it will be ignored.
Any batch output will preserve the order of the batch input.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCMACWrite,
			},
		},

		HelpSynopsis:    pathCMACHelpSyn,
		HelpDescription: pathCMACHelpDesc,
	}
}

func (b *backend) pathCMACWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	input, warnings, err := parseCmacWriteInput(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return b.runWithReadLockedPolicy(ctx, req.Storage, input.KeyName, func(p *keysutil.Policy) (*logical.Response, error) {
		if p.Type == keysutil.KeyType_MANAGED_KEY {
			return logical.ErrorResponse("CMAC creation is not supported with managed keys"), logical.ErrInvalidRequest
		}

		if !p.Type.CMACSupported() {
			return logical.ErrorResponse("key %s is not a supported CMAC key type: %s", p.Name, p.Type), logical.ErrInvalidRequest
		}

		input.KeyVersion, err = validateKeyVersion(p, input.KeyVersion)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		cmacKey, err := p.CMACKey(input.KeyVersion)
		if err != nil {
			return logical.ErrorResponse("failed fetching CMAC key for %s with version %d: %s", input.KeyName, input.KeyVersion, err.Error()), logical.ErrInvalidRequest
		}

		cmacResponseItems := make([]cmacWriteResponseItem, len(input.Items))
		for i, item := range input.Items {
			itemRes := cmacWriteResponseItem{
				Reference: item.Reference,
			}

			cmacResp, err := performCmac(cmacKey, item.MacLength, item.Input)
			if err != nil {
				errMsg := fmt.Sprintf("failed processing item %d: %s", i, err.Error())
				itemRes.Error = errMsg
			} else {
				itemRes.CMAC = encodeTransitSignature(cmacResp, input.KeyVersion)
			}

			cmacResponseItems[i] = itemRes
		}

		dataResponse, err := buildCmacWriteResponse(input.IsBatch, cmacResponseItems)
		if err != nil {
			return nil, err
		}

		resp := &logical.Response{
			Data:     dataResponse,
			Warnings: warnings,
		}

		return resp, nil
	})
}

func parseCmacWriteInput(d *framework.FieldData) (cmacWriteInput, []string, error) {
	keyName := d.Get("name").(string)
	keyVersion := d.Get("key_version").(int)

	if strings.TrimSpace(keyName) == "" {
		return cmacWriteInput{}, nil, fmt.Errorf("name parameter must be provided")
	}

	urlMacLengthRaw := d.Get("url_mac_length").(string)
	overrideMacLength := false
	urlMacLength := 0
	if strings.TrimSpace(urlMacLengthRaw) != "" {
		var err error
		overrideMacLength = true
		urlMacLength, err = parseutil.SafeParseInt(urlMacLengthRaw)
		if err != nil {
			return cmacWriteInput{}, nil, fmt.Errorf("the url mac-length parameter failed parsing: %w", err)
		}
		err = validateMacLength(urlMacLength)
		if err != nil {
			return cmacWriteInput{}, nil, fmt.Errorf("the url mac-length parameter is invalid: %w", err)
		}
	}

	batchInputRaw, isBatchInput := d.GetOk("batch_input")
	var batchInputItems []cmacWriteItem
	if isBatchInput {
		err := mapstructure.Decode(batchInputRaw, &batchInputItems)
		if err != nil {
			return cmacWriteInput{}, nil, fmt.Errorf("failed to parse batch input: %w", err)
		}
	} else {
		inputB64 := d.Get("input").(string)
		macLength := d.Get("mac_length").(int)

		batchInputItems = make([]cmacWriteItem, 1)
		batchInputItems[0] = cmacWriteItem{
			Input:     inputB64,
			MacLength: macLength,
		}
	}

	if len(batchInputItems) == 0 {
		return cmacWriteInput{}, nil, fmt.Errorf("no inputs to process")
	}

	var warnings []string
	for i := range batchInputItems {
		if strings.TrimSpace(batchInputItems[i].Input) == "" {
			return cmacWriteInput{}, nil, fmt.Errorf("input item %d was blank", i)
		}

		if overrideMacLength {
			if batchInputItems[i].MacLength != 0 {
				msg := fmt.Sprintf("input item %d mac_length of %d overridden by url_mac_length %d", i, batchInputItems[i].MacLength, urlMacLength)
				warnings = append(warnings, msg)
			}
			batchInputItems[i].MacLength = urlMacLength
		}
	}

	input := cmacWriteInput{
		IsBatch:    isBatchInput,
		KeyName:    keyName,
		KeyVersion: keyVersion,
		Items:      batchInputItems,
	}

	return input, warnings, nil
}

func performCmac(cmacKey []byte, tagSize int, base64Input string) ([]byte, error) {
	rawInput, err := base64.StdEncoding.DecodeString(base64Input)
	if err != nil {
		return nil, fmt.Errorf("failed base 64 decoding: %w", err)
	}

	var hf hash.Hash
	switch tagSize {
	case 0:
		hf, err = aesCmac.New(cmacKey)
	default:
		hf, err = aesCmac.NewWithTagSize(cmacKey, tagSize)
	}

	if err != nil {
		return nil, fmt.Errorf("failed generating cmac hash function: %w", err)
	}

	return hf.Sum(rawInput), nil
}

func buildCmacWriteResponse(wasBatch bool, responses []cmacWriteResponseItem) (map[string]interface{}, error) {
	numResponses := len(responses)
	if numResponses < 1 {
		return nil, fmt.Errorf("no CMAC responses were generated")
	}

	if wasBatch {
		return map[string]interface{}{
			"batch_results": responses,
		}, nil
	}

	response := responses[0]
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}

	return map[string]interface{}{
		"cmac": response.CMAC,
	}, nil
}

func (b *backend) pathCMACVerify(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	input, warnings, err := parseCmacVerifyInput(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return b.runWithReadLockedPolicy(ctx, req.Storage, input.KeyName, func(p *keysutil.Policy) (*logical.Response, error) {
		if p.Type == keysutil.KeyType_MANAGED_KEY {
			return logical.ErrorResponse("CMAC verification is not supported with managed keys"), logical.ErrInvalidRequest
		}

		if !p.Type.CMACSupported() {
			return logical.ErrorResponse("key %s is not a supported CMAC key type: %s", p.Name, p.Type), logical.ErrInvalidRequest
		}

		verifyCmacResponses := make([]cmacVerifyResponseItem, len(input.Items))
		for i, item := range input.Items {
			itemRes := cmacVerifyResponseItem{
				Reference: item.Reference,
			}

			isValid, err := verifyCmacInput(p, item)
			if err != nil {
				itemRes.Error = err.Error()
			} else {
				itemRes.Valid = isValid
			}

			verifyCmacResponses[i] = itemRes
		}

		dataResponse, err := buildCmacVerifyResponse(input.IsBatch, verifyCmacResponses)
		if err != nil {
			return nil, err
		}

		resp := &logical.Response{
			Data:     dataResponse,
			Warnings: warnings,
		}

		return resp, nil
	})
}

func parseCmacVerifyInput(d *framework.FieldData) (cmacVerifyInput, []string, error) {
	keyName := d.Get("name").(string)

	if strings.TrimSpace(keyName) == "" {
		return cmacVerifyInput{}, nil, fmt.Errorf("name parameter must be provided")
	}

	urlAlgo := d.Get("urlalgorithm").(string) // Reuse the existing value for hmac in the url
	urlMacLength := 0
	overrideMacLength := false
	if strings.TrimSpace(urlAlgo) != "" {
		var err error
		urlMacLength, err = parseutil.SafeParseInt(urlAlgo)
		if err != nil {
			return cmacVerifyInput{}, nil, fmt.Errorf("the url mac-length parameter failed parsing: %w", err)
		}
		overrideMacLength = true
		err = validateMacLength(urlMacLength)
		if err != nil {
			return cmacVerifyInput{}, nil, fmt.Errorf("the url mac-length parameter is invalid: %w", err)
		}
	}

	batchInputRaw, isBatchInput := d.GetOk("batch_input")
	var batchInputItems []cmacVerifyItem
	if isBatchInput {
		err := mapstructure.Decode(batchInputRaw, &batchInputItems)
		if err != nil {
			return cmacVerifyInput{}, nil, fmt.Errorf("failed to parse batch input: %w", err)
		}
	} else {
		inputB64 := d.Get("input").(string)
		macLength := d.Get("mac_length").(int)
		cmac := d.Get("cmac").(string)

		batchInputItems = make([]cmacVerifyItem, 1)
		batchInputItems[0] = cmacVerifyItem{
			Input:     inputB64,
			MacLength: macLength,
			Cmac:      cmac,
		}
	}

	if len(batchInputItems) == 0 {
		return cmacVerifyInput{}, nil, fmt.Errorf("no inputs to process")
	}

	var warnings []string
	for i := range batchInputItems {
		if strings.TrimSpace(batchInputItems[i].Input) == "" {
			return cmacVerifyInput{}, nil, fmt.Errorf("input field on item %d was blank", i)
		}
		if strings.TrimSpace(batchInputItems[i].Cmac) == "" {
			return cmacVerifyInput{}, nil, fmt.Errorf("cmac field on item %d was blank", i)
		}

		if overrideMacLength {
			if batchInputItems[i].MacLength != 0 {
				msg := fmt.Sprintf("input item %d mac_length of %d overridden by url_mac_length %d", i, batchInputItems[i].MacLength, urlMacLength)
				warnings = append(warnings, msg)
			}
			batchInputItems[i].MacLength = urlMacLength
		} else {
			err := validateMacLength(batchInputItems[i].MacLength)
			if err != nil {
				return cmacVerifyInput{}, nil, fmt.Errorf("mac_length field on item %d was not valid: %w", i, err)
			}
		}
	}

	input := cmacVerifyInput{
		IsBatch: isBatchInput,
		KeyName: keyName,
		Items:   batchInputItems,
	}

	return input, warnings, nil
}

func validateMacLength(length int) error {
	if length < 0 {
		return fmt.Errorf("must not be less than 0")
	}

	if length > aes.BlockSize {
		return fmt.Errorf("must not be greater than %d", aes.BlockSize)
	}

	return nil
}

func verifyCmacInput(p *keysutil.Policy, item cmacVerifyItem) (bool, error) {
	origCmac, keyVersion, err := decodeTransitSignature(item.Cmac)
	if err != nil {
		return false, err
	}

	newKeyVersion, err := validateKeyVersion(p, keyVersion)
	if err != nil {
		return false, fmt.Errorf("invalid key version for key %s: %w", p.Name, err)
	}

	if newKeyVersion != keyVersion {
		return false, fmt.Errorf("%s: %w", p.Name, err)
	}

	cmacKey, err := p.CMACKey(keyVersion)
	if err != nil {
		return false, fmt.Errorf("failed fetching CMAC key for %s with version %d: %w", p.Name, keyVersion, err)
	}

	cmac, err := performCmac(cmacKey, item.MacLength, item.Input)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(origCmac, cmac) == 1, nil
}

func buildCmacVerifyResponse(isBatch bool, responses []cmacVerifyResponseItem) (map[string]interface{}, error) {
	numResponses := len(responses)
	if numResponses < 1 {
		return nil, fmt.Errorf("no CMAC verification responses were generated")
	}

	if isBatch {
		return map[string]interface{}{
			"batch_results": responses,
		}, nil
	}

	response := responses[0]
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}

	return map[string]interface{}{
		"valid": response.Valid,
	}, nil
}

const pathCMACHelpSyn = `Generate a CMAC for input data using the named key`

const pathCMACHelpDesc = `
Generates a CMAC against the given input data and the named key.
`
