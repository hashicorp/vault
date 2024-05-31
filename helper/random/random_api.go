// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package random

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/xor"
	"github.com/hashicorp/vault/sdk/logical"
)

const APIMaxBytes = 128 * 1024

func HandleRandomAPI(d *framework.FieldData, additionalSource io.Reader) (*logical.Response, error) {
	bytes := 0
	// Parsing is convoluted here, but allows operators to ACL both source and byte count
	maybeUrlBytes := d.Raw["urlbytes"]
	maybeSource := d.Raw["source"]
	source := "platform"
	var err error
	if maybeSource == "" {
		bytes = d.Get("bytes").(int)
	} else if maybeUrlBytes == "" && isValidSource(maybeSource.(string)) {
		source = maybeSource.(string)
		bytes = d.Get("bytes").(int)
	} else if maybeUrlBytes == "" {
		bytes, err = strconv.Atoi(maybeSource.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error parsing url-set byte count: %s", err)), nil
		}
	} else {
		source = maybeSource.(string)
		bytes, err = strconv.Atoi(maybeUrlBytes.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("error parsing url-set byte count: %s", err)), nil
		}
	}
	format := d.Get("format").(string)

	if bytes < 1 {
		return logical.ErrorResponse(`"bytes" cannot be less than 1`), nil
	}

	if bytes > APIMaxBytes {
		return logical.ErrorResponse(`"bytes" should be less than %d`, APIMaxBytes), nil
	}

	switch format {
	case "hex":
	case "base64":
	default:
		return logical.ErrorResponse("unsupported encoding format %q; must be \"hex\" or \"base64\"", format), nil
	}

	var randBytes []byte
	var warning string
	switch source {
	case "", "platform":
		randBytes, err = uuid.GenerateRandomBytes(bytes)
		if err != nil {
			return nil, err
		}
	case "seal":
		if rand.Reader == additionalSource {
			warning = "no seal/entropy augmentation available, using platform entropy source"
		}
		randBytes, err = uuid.GenerateRandomBytesWithReader(bytes, additionalSource)
	case "all":
		randBytes, err = uuid.GenerateRandomBytes(bytes)
		if err == nil && rand.Reader != additionalSource {
			var sealBytes []byte
			sealBytes, err = uuid.GenerateRandomBytesWithReader(bytes, additionalSource)
			if err == nil {
				randBytes, err = xor.XORBytes(sealBytes, randBytes)
			}
		}
	default:
		return logical.ErrorResponse("unsupported entropy source %q; must be \"platform\" or \"seal\", or \"all\"", source), nil
	}
	if err != nil {
		return nil, err
	}

	var retStr string
	switch format {
	case "hex":
		retStr = hex.EncodeToString(randBytes)
	case "base64":
		retStr = base64.StdEncoding.EncodeToString(randBytes)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"random_bytes": retStr,
		},
	}
	if warning != "" {
		resp.Warnings = []string{warning}
	}
	return resp, nil
}

func isValidSource(s string) bool {
	switch s {
	case "", "platform", "seal", "all":
		return true
	}
	return false
}
