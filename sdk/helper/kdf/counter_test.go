package kdf

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"testing"
)

func TestSubsetNISTCAVPCounterKBKDF(t *testing.T) {
	for index, test := range counterModeTestCases {
		key, err := hex.DecodeString(test.KI)
		if err != nil {
			t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.KI)
		}

		prf := hmac.New(test.Hash, key)
		var params []KBKDFParameter

		if test.DataBefore != "" {
			data, err := hex.DecodeString(test.DataBefore)
			if err != nil {
				t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.DataBefore)
			}
			params = append(params, ByteArray(data))
		}

		// NIST CAVP tests are always big-endian.
		params = append(params, CounterVariable{false, test.RLen})

		if test.DataAfter != "" {
			data, err := hex.DecodeString(test.DataAfter)
			if err != nil {
				t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.DataAfter)
			}
			params = append(params, ByteArray(data))
		}

		kbkdf, err := NewCounter(prf, params, []int{test.L})
		if err != nil {
			t.Fatalf("unexpected error creating Counter KDF: %v", err)
		}

		expected, err := hex.DecodeString(test.KO)
		if err != nil {
			t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.KO)
		}

		var output []byte = make([]byte, test.L/8)
		length, err := kbkdf.Read(output)
		if err != nil {
			t.Fatalf("unexpected error reading from Counter KDF: %v", err)
		}

		if length != test.L/8 {
			t.Fatalf("unexpected return size from read: got %d, expected %d on test %d", length, test.L/8, index)
		}

		if !bytes.Equal(output, expected) {
			t.Fatalf("Test failed: %d/%d, %s != %s", index, test.Count, hex.EncodeToString(output), hex.EncodeToString(expected))
		}
	}
}
