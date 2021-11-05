package kdf

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"testing"
)

func TestSubsetNISTCAVPFeedbackKBKDF(t *testing.T) {
	for index, test := range feedbackModeTestCases {
		key, err := hex.DecodeString(test.KI)
		if err != nil {
			t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.KI)
		}

		prf := hmac.New(test.Hash, key)
		var params []KBKDFParameter

		if test.Counter == ctrBeforeIter {
			params = append(params, CounterVariable{false, test.RLen})
			params = append(params, ChainingVariable{})
			data, err := hex.DecodeString(test.FixedData)
			if err != nil {
				t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.FixedData)
			}
			params = append(params, ByteArray{data})
		} else if test.Counter == ctrAfterIter {
			params = append(params, ChainingVariable{})
			params = append(params, CounterVariable{false, test.RLen})
			data, err := hex.DecodeString(test.FixedData)
			if err != nil {
				t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.FixedData)
			}
			params = append(params, ByteArray{data})
		} else {
			params = append(params, ChainingVariable{})
			data, err := hex.DecodeString(test.FixedData)
			if err != nil {
				t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.FixedData)
			}
			params = append(params, ByteArray{data})
			params = append(params, CounterVariable{false, test.RLen})
		}

		iv, err := hex.DecodeString(test.IV)
		if err != nil {
			t.Fatalf("unexpected error decoding hex string: %v / %s", err, test.IV)
		}

		kbkdf, err := NewFeedback(prf, params, iv, []int{test.L})
		if err != nil {
			t.Fatalf("unexpected error creating Feedback KDF: %v", err)
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
