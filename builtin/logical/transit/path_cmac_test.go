package transit

import (
	"context"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCmacValidOptions validates various successful CMAC generation and verification options
// work together using the normal non-batch input parameters
func TestCmacValidOptions(t *testing.T) {
	b, s, keyNames := createBackendWithCmacKeys(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		keyVersion    int
		macLength     int
		urlMacLength  int
		expectWarning bool
	}{
		{"defaults", -1, -1, -1, false},
		{"force-key-version-0", 0, -1, -1, false},
		{"force-key-version-1", 1, -1, -1, false},
		{"force-key-version-2", 2, -1, -1, false},
		{"with-mac-length", -1, 2, -1, false},
		{"with-url-mac-length", -1, -1, 8, false},
		{"with-override-url-mac-length", -1, 3, aes.BlockSize, true},
	}
	for _, keyName := range keyNames {
		for _, tc := range tests {
			testName := fmt.Sprintf("%s-%s", keyName, tc.name)
			t.Run(testName, func(t *testing.T) {
				input := "dGhlIHF1aWNrIGJyb3duIGZveA==" // "the quick brown fox"
				cmacPath := fmt.Sprintf("cmac/%s", keyName)
				if tc.urlMacLength != -1 {
					cmacPath = fmt.Sprintf("cmac/%s/%d", keyName, tc.urlMacLength)
				}
				cmacReq := &logical.Request{
					Path:      cmacPath,
					Operation: logical.UpdateOperation,
					Storage:   s,
					Data: map[string]interface{}{
						"input": input,
					},
				}

				if tc.keyVersion != -1 {
					cmacReq.Data["key_version"] = tc.keyVersion
				}

				if tc.macLength != -1 {
					cmacReq.Data["mac_length"] = tc.macLength
				}

				resp, err := b.HandleRequest(ctx, cmacReq)
				if err != nil || (resp != nil && resp.IsError()) {
					t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
				}

				cmac := resp.Data["cmac"].(string)
				assert.Equal(t, tc.expectWarning, len(resp.Warnings) > 0, "expect warnings was %v got %v", tc.expectWarning, resp.Warnings)

				_, myKeyVersion, err := decodeTransitSignature(cmac)
				require.NoError(t, err, "failed decoding cmac signature")

				if tc.keyVersion == -1 || tc.keyVersion == 0 {
					// If we didn't specify the key version or specified v0, we should have the
					// latest key version within our signature output prefix
					require.Equal(t, 2, myKeyVersion)
				} else {
					require.Equal(t, tc.keyVersion, myKeyVersion)
				}

				verifyPath := fmt.Sprintf("verify/%s", keyName)
				if tc.urlMacLength != -1 {
					verifyPath = fmt.Sprintf("verify/%s/%d", keyName, tc.urlMacLength)
				}
				cmacVerifyReq := &logical.Request{
					Path:      verifyPath,
					Operation: logical.UpdateOperation,
					Storage:   s,
					Data: map[string]interface{}{
						"input": input,
						"cmac":  cmac,
					},
				}
				if tc.macLength != -1 {
					cmacVerifyReq.Data["mac_length"] = tc.macLength
				}

				resp, err = b.HandleRequest(ctx, cmacVerifyReq)
				if err != nil || (resp != nil && resp.IsError()) {
					t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
				}

				require.True(t, resp.Data["valid"].(bool), "verification of cmac failed")
				assert.Equal(t, tc.expectWarning, len(resp.Warnings) > 0, "expect warnings was %v got %v", tc.expectWarning, resp.Warnings)
			})
		}
	}
}

// TestCmacValidBatchInput validates various successful CMAC generation and verification options
// work together using the batch input parameter
func TestCmacValidBatchInput(t *testing.T) {
	b, s, keyNames := createBackendWithCmacKeys(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		keyVersion    int
		urlMacLength  int
		expectWarning bool
	}{
		{"defaults", -1, -1, false},
		{"force-key-version-0", 0, -1, false},
		{"force-key-version-1", 1, -1, false},
		{"force-key-version-2", 2, -1, false},
		{"with-override-url-mac-length", -1, aes.BlockSize, true},
	}
	for _, keyName := range keyNames {
		for _, tc := range tests {
			testName := fmt.Sprintf("%s-%s", keyName, tc.name)
			t.Run(testName, func(t *testing.T) {
				batchInput := make([]map[string]interface{}, 0, 10)
				for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
					index := strconv.Itoa(i)
					item := map[string]interface{}{
						"input":      base64.StdEncoding.EncodeToString([]byte("the quick brown fox " + index)),
						"reference":  index,
						"mac_length": i,
					}

					batchInput = append(batchInput, item)
				}
				cmacPath := fmt.Sprintf("cmac/%s", keyName)
				if tc.urlMacLength != -1 {
					cmacPath = fmt.Sprintf("cmac/%s/%d", keyName, tc.urlMacLength)
				}
				cmacReq := &logical.Request{
					Path:      cmacPath,
					Operation: logical.UpdateOperation,
					Storage:   s,
					Data: map[string]interface{}{
						"batch_input": batchInput,
					},
				}

				if tc.keyVersion != -1 {
					cmacReq.Data["key_version"] = tc.keyVersion
				}

				resp, err := b.HandleRequest(ctx, cmacReq)
				if err != nil || (resp != nil && resp.IsError()) {
					t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
				}
				assert.Equal(t, tc.expectWarning, len(resp.Warnings) > 0, "expect warnings was %v got %v", tc.expectWarning, resp.Warnings)

				verifyPath := fmt.Sprintf("verify/%s", keyName)
				if tc.urlMacLength != -1 {
					verifyPath = fmt.Sprintf("verify/%s/%d", keyName, tc.urlMacLength)
				}

				batchResp := resp.Data["batch_results"].([]cmacWriteResponseItem)
				cmacByRef := make(map[string]string)
				for _, result := range batchResp {
					ref := result.Reference
					cmac := result.CMAC

					require.NotContains(t, cmacByRef, ref, "duplicated reference value %v in batch: %v", ref, batchResp)
					cmacByRef[ref] = cmac
				}

				batchVerifyInput := make([]map[string]interface{}, 0, 10)
				for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
					index := strconv.Itoa(i)
					item := map[string]interface{}{
						"input":      base64.StdEncoding.EncodeToString([]byte("the quick brown fox " + index)),
						"cmac":       cmacByRef[index],
						"reference":  index,
						"mac_length": i,
					}

					batchVerifyInput = append(batchVerifyInput, item)
				}

				cmacVerifyReq := &logical.Request{
					Path:      verifyPath,
					Operation: logical.UpdateOperation,
					Storage:   s,
					Data: map[string]interface{}{
						"batch_input": batchVerifyInput,
					},
				}

				resp, err = b.HandleRequest(ctx, cmacVerifyReq)
				if err != nil || (resp != nil && resp.IsError()) {
					t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
				}

				require.Contains(t, resp.Data, "batch_results", "verify response is missing batch_results: %v", resp.Data)
				verifyRes := resp.Data["batch_results"].([]cmacVerifyResponseItem)
				for i, res := range verifyRes {
					require.True(t, res.Valid, "value was not considered valid: %v: %v", res, resp.Data)
					require.Equal(t, strconv.Itoa(i), res.Reference, "reference values did not match")
				}
				assert.Equal(t, tc.expectWarning, len(resp.Warnings) > 0, "expect warnings was %v got %v", tc.expectWarning, resp.Warnings)
			})
		}
	}
}

// TestCmacVerifyInvalidCmacEntries verifies that we return an output of not valid for different
// scenarios that CMAC validation should fail, but not due to invalid input parameters.
func TestCmacVerifyInvalidCmacEntries(t *testing.T) {
	b, s, keyNames := createBackendWithCmacKeys(t)
	ctx := context.Background()

	// Get a valid CMAC first
	keyName := keyNames[0]
	input := "dGhlIHF1aWNrIGJyb3duIGZveA==" // "the quick brown fox"
	cmacReq := &logical.Request{
		Path:      fmt.Sprintf("cmac/%s", keyName),
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"input": input,
		},
	}
	resp, err := b.HandleRequest(ctx, cmacReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	cmac := resp.Data["cmac"].(string)

	garbageCmac := base64.StdEncoding.EncodeToString([]byte("garbage-cmac"))

	tests := []struct {
		name      string
		keyName   string
		macLength int
		cmac      string
	}{
		{"different-key-version", keyName, -1, strings.ReplaceAll(cmac, ":v2:", ":v1:")},
		{"different-key", keyNames[1], -1, cmac},
		{"bad-base64-cmac", keyName, -1, "vault:v2:" + garbageCmac},
		{"bad-mac-length", keyName, 8, cmac},
	}
	for _, tc := range tests {
		testName := fmt.Sprintf("%s-%s", tc.keyName, tc.name)
		t.Run(testName, func(t *testing.T) {
			cmacVerifyReq := &logical.Request{
				Path:      fmt.Sprintf("verify/%s", tc.keyName),
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data: map[string]interface{}{
					"input": input,
					"cmac":  tc.cmac,
				},
			}
			if tc.macLength != -1 {
				cmacVerifyReq.Data["mac_length"] = tc.macLength
			}

			resp, err = b.HandleRequest(ctx, cmacVerifyReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
			}

			require.False(t, resp.Data["valid"].(bool), "verification of cmac passed, should have failed")
		})
	}
}

// TestCmacWriteInvalidOptions verifies we properly error out when invalid inputs are provided to
// the CMAC Write API call
func TestCmacWriteInvalidOptions(t *testing.T) {
	b, s, keyNames := createBackendWithCmacKeys(t)
	ctx := context.Background()

	// Create a non-cmac key
	keyReq := &logical.Request{
		Path:      "keys/encrypt-key",
		Operation: logical.UpdateOperation,
		Data:      map[string]interface{}{},
		Storage:   s,
	}

	resp, err := b.HandleRequest(ctx, keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to create encrypt-key: err: %v\nresp: %#v", err, resp)
	}

	keyName := keyNames[0]
	input := "dGhlIHF1aWNrIGJyb3duIGZveA==" // "the quick brown fox"

	tests := []struct {
		name         string
		keyName      string
		input        string
		keyVersion   int
		macLength    int
		urlMacLength int
	}{
		{"bad-keyname", "a-bad-key", input, -1, -1, -1},
		{"non-cmac-key", "encrypt-key", input, -1, -1, -1},
		{"missing-input", keyName, "", -1, -1, -1},
		{"non-base64-input", keyName, "garbage-input", -1, -1, -1},
		{"bad-neg-keyversion", keyName, input, -10, -1, -1},
		{"bad-keyversion", keyName, input, 100, -1, -1},
		{"bad-mac-length", keyName, input, -1, 100, -1},
		{"bad-url-mac-length", keyName, input, -1, -1, 100},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmacPath := fmt.Sprintf("cmac/%s", tc.keyName)
			if tc.urlMacLength != -1 {
				cmacPath = fmt.Sprintf("cmac/%s/%d", tc.keyName, tc.urlMacLength)
			}
			data := map[string]interface{}{}
			if tc.input != "" {
				data["input"] = tc.input
			}
			if tc.keyVersion != -1 {
				data["key_version"] = tc.keyVersion
			}
			if tc.macLength != -1 {
				data["mac_length"] = tc.macLength
			}

			cmacReq := &logical.Request{
				Path:      cmacPath,
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data:      data,
			}

			resp, err := b.HandleRequest(ctx, cmacReq)
			if err == nil && (resp != nil && !resp.IsError()) {
				t.Fatalf("expected an error got none: resp: %v", resp)
			}
		})
	}

	batchTests := []struct {
		name string
		data interface{}
	}{
		{name: "empty-input-batch", data: []map[string]interface{}{}},
		{name: "nil-input-batch", data: nil},
	}
	for _, tc := range batchTests {
		t.Run(tc.name, func(t *testing.T) {
			// Test empty batch_input
			cmacReq := &logical.Request{
				Path:      fmt.Sprintf("cmac/%s", keyName),
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data: map[string]interface{}{
					"batch_input": tc.data,
				},
			}

			resp, err := b.HandleRequest(ctx, cmacReq)
			if err == nil && (resp != nil && !resp.IsError()) {
				t.Fatalf("expected an error got none: resp: %v", resp)
			}
		})
	}
}

// TestCmacVerifyInvalidOptions verifies we properly error out when invalid inputs are provided to
// the verify API in a cmac mode.
func TestCmacVerifyInvalidOptions(t *testing.T) {
	b, s, keyNames := createBackendWithCmacKeys(t)
	ctx := context.Background()

	// Create a non-cmac key
	keyReq := &logical.Request{
		Path:      "keys/encrypt-key",
		Operation: logical.UpdateOperation,
		Data:      map[string]interface{}{},
		Storage:   s,
	}

	resp, err := b.HandleRequest(ctx, keyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to create encrypt-key: err: %v\nresp: %#v", err, resp)
	}

	// Get a valid CMAC first
	keyName := keyNames[0]
	input := "dGhlIHF1aWNrIGJyb3duIGZveA==" // "the quick brown fox"
	cmacReq := &logical.Request{
		Path:      fmt.Sprintf("cmac/%s", keyName),
		Operation: logical.UpdateOperation,
		Storage:   s,
		Data: map[string]interface{}{
			"input": input,
		},
	}
	resp, err = b.HandleRequest(ctx, cmacReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	cmac := resp.Data["cmac"].(string)

	tests := []struct {
		name         string
		keyName      string
		input        string
		cmac         string
		macLength    int
		urlMacLength int
	}{
		{"bad-keyname", "a-bad-key", input, cmac, -1, -1},
		{"non-cmac-key", "encrypt-key", input, cmac, -1, -1},
		{"missing-input", keyName, "", cmac, -1, -1},
		{"non-base64-input", keyName, "garbage-input", cmac, -1, -1},
		{"non-base64-cmac", keyName, input, "vault:v2:garbage-input", -1, -1},
		{"no-base64-cmac", keyName, input, "vault:v2:", -1, -1},
		{"non-digit-keyversion", keyName, input, strings.ReplaceAll(cmac, ":v2:", ":vblah:"), -1, -1},
		{"bad-neg-keyversion", keyName, input, strings.ReplaceAll(cmac, ":v2:", ":v-10:"), -1, -1},
		{"bad-keyversion", keyName, input, strings.ReplaceAll(cmac, ":v2:", ":v10:"), -1, -1},
		{"bad-mac-length", keyName, input, cmac, 100, -1},
		{"bad-url-mac-length", keyName, input, cmac, -1, 100},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmacPath := fmt.Sprintf("verify/%s", tc.keyName)
			if tc.urlMacLength != -1 {
				cmacPath = fmt.Sprintf("verify/%s/%d", tc.keyName, tc.urlMacLength)
			}
			data := map[string]interface{}{}
			if tc.input != "" {
				data["input"] = tc.input
			}
			if tc.macLength != -1 {
				data["mac_length"] = tc.macLength
			}
			if tc.cmac != "" {
				data["cmac"] = tc.cmac
			}

			cmacReq := &logical.Request{
				Path:      cmacPath,
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data:      data,
			}

			resp, err := b.HandleRequest(ctx, cmacReq)
			if err == nil && (resp != nil && !resp.IsError()) {
				t.Fatalf("expected an error got none: resp: %v", resp)
			}
			t.Logf("ERR: %v", err)
			if resp != nil {
				t.Logf("Response Error: %v", resp.Error())
			}
		})
	}

	batchTests := []struct {
		name string
		data interface{}
	}{
		{name: "empty-input-batch", data: []map[string]interface{}{}},
		{name: "nil-input-batch", data: nil},
	}
	for _, tc := range batchTests {
		t.Run(tc.name, func(t *testing.T) {
			// Test empty batch_input
			cmacReq := &logical.Request{
				Path:      fmt.Sprintf("verify/%s", keyName),
				Operation: logical.UpdateOperation,
				Storage:   s,
				Data: map[string]interface{}{
					"batch_input": tc.data,
				},
			}

			resp, err := b.HandleRequest(ctx, cmacReq)
			if err == nil && (resp != nil && !resp.IsError()) {
				t.Fatalf("expected an error got none: resp: %v", resp)
			}
			t.Logf("ERR: %v", err)
			if resp != nil {
				t.Logf("Response Error: %v", resp.Error())
			}
		})
	}
}

func createBackendWithCmacKeys(t *testing.T) (*backend, logical.Storage, []string) {
	b, s := createBackendWithStorage(t)
	ctx := context.Background()

	keyNames := []string{"aes128-cmac", "aes256-cmac"}
	for _, name := range keyNames {
		keyReq := &logical.Request{
			Path:      "keys/" + name,
			Operation: logical.UpdateOperation,
			Data: map[string]interface{}{
				"type": name,
			},
			Storage: s,
		}

		resp, err := b.HandleRequest(ctx, keyReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}

		// Add another key version
		rotateReq := &logical.Request{
			Path:      fmt.Sprintf("keys/%s/rotate", name),
			Operation: logical.UpdateOperation,
			Data:      map[string]interface{}{},
			Storage:   s,
		}

		resp, err = b.HandleRequest(ctx, rotateReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
		}
	}
	return b, s, keyNames
}
