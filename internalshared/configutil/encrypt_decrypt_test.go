// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configutil

import (
	"bytes"
	"context"
	"encoding/base64"
	"testing"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"google.golang.org/protobuf/proto"
)

func getAEADTestKMS(t *testing.T) {
}

func TestEncryptParams(t *testing.T) {
	rawStr := `
storage "consul" {
	api_key = "{{encrypt(foobar)}}"
}

telemetry {
	some_param = "something"
	circonus_api_key = "{{encrypt(barfoo)}}"
}
`

	finalStr := `
storage "consul" {
	api_key = "foobar"
}

telemetry {
	some_param = "something"
	circonus_api_key = "barfoo"
}
`

	reverser := new(reversingWrapper)
	out, err := EncryptDecrypt(rawStr, false, false, reverser)
	if err != nil {
		t.Fatal(err)
	}

	first := true
	locs := decryptRegex.FindAllIndex([]byte(out), -1)
	for _, match := range locs {
		matchBytes := []byte(out)[match[0]:match[1]]
		matchBytes = bytes.TrimSuffix(bytes.TrimPrefix(matchBytes, []byte("{{decrypt(")), []byte(")}}"))
		inMsg, err := base64.RawURLEncoding.DecodeString(string(matchBytes))
		if err != nil {
			t.Fatal(err)
		}
		inBlob := new(wrapping.BlobInfo)
		if err := proto.Unmarshal(inMsg, inBlob); err != nil {
			t.Fatal(err)
		}
		ct := string(inBlob.Ciphertext)
		if first {
			if ct != "raboof" {
				t.Fatal(ct)
			}
			first = false
		} else {
			if ct != "oofrab" {
				t.Fatal(ct)
			}
		}
	}

	decOut, err := EncryptDecrypt(out, true, false, reverser)
	if err != nil {
		t.Fatal(err)
	}

	if decOut != rawStr {
		t.Fatal(decOut)
	}

	decOut, err = EncryptDecrypt(out, true, true, reverser)
	if err != nil {
		t.Fatal(err)
	}

	if decOut != finalStr {
		t.Fatal(decOut)
	}
}

type reversingWrapper struct{}

func (r *reversingWrapper) Type(_ context.Context) (wrapping.WrapperType, error) {
	return "reverser", nil
}
func (r *reversingWrapper) KeyId(_ context.Context) (string, error) { return "reverser", nil }
func (r *reversingWrapper) HMACKeyID() string                       { return "" }
func (r *reversingWrapper) Init(_ context.Context) error            { return nil }
func (r *reversingWrapper) Finalize(_ context.Context) error        { return nil }
func (r *reversingWrapper) SetConfig(_ context.Context, opts ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	return &wrapping.WrapperConfig{}, nil
}

func (r *reversingWrapper) Encrypt(_ context.Context, input []byte, _ ...wrapping.Option) (*wrapping.BlobInfo, error) {
	return &wrapping.BlobInfo{
		Ciphertext: r.reverse(input),
	}, nil
}

func (r *reversingWrapper) Decrypt(_ context.Context, input *wrapping.BlobInfo, _ ...wrapping.Option) ([]byte, error) {
	return r.reverse(input.Ciphertext), nil
}

func (r *reversingWrapper) reverse(input []byte) []byte {
	output := make([]byte, len(input))
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = input[j], input[i]
	}
	return output
}
