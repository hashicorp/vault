// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package seal

import (
	"bytes"
	"testing"
)

func TestEnvelope(t *testing.T) {
	input := []byte("test")
	env, err := NewEnvelope().Encrypt(input, nil)
	if err != nil {
		t.Fatal(err)
	}

	output, err := NewEnvelope().Decrypt(env, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(input, output) {
		t.Fatalf("expected the same text: expected %s, got %s", string(input), string(output))
	}
}
