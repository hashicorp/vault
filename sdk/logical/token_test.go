// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestJSONSerialization(t *testing.T) {
	tt := TokenTypeDefaultBatch
	s, err := json.Marshal(tt)
	if err != nil {
		t.Fatal(err)
	}

	var utt TokenType
	err = json.Unmarshal(s, &utt)
	if err != nil {
		t.Fatal(err)
	}

	if tt != utt {
		t.Fatalf("expected %v, got %v", tt, utt)
	}

	utt = TokenTypeDefault
	err = json.Unmarshal([]byte(`"default-batch"`), &utt)
	if err != nil {
		t.Fatal(err)
	}
	if tt != utt {
		t.Fatalf("expected %v, got %v", tt, utt)
	}

	// Test on an empty value, which should unmarshal into TokenTypeDefault
	tt = TokenTypeDefault
	err = json.Unmarshal([]byte(`""`), &utt)
	if err != nil {
		t.Fatal(err)
	}
	if tt != utt {
		t.Fatalf("expected %v, got %v", tt, utt)
	}
}

// TestCreateClientID verifies that CreateClientID uses the entity ID for a token
// entry if one exists, and creates an appropriate client ID otherwise.
func TestCreateClientID(t *testing.T) {
	entry := TokenEntry{NamespaceID: "namespaceFoo", Policies: []string{"bar", "baz", "foo", "banana"}}
	id, isTWE := entry.CreateClientID()
	if !isTWE {
		t.Fatalf("TWE token should return true value in isTWE bool")
	}
	expectedIDPlaintext := "banana" + string(SortedPoliciesTWEDelimiter) + "bar" +
		string(SortedPoliciesTWEDelimiter) + "baz" +
		string(SortedPoliciesTWEDelimiter) + "foo" + string(ClientIDTWEDelimiter) + "namespaceFoo"

	hashed := sha256.Sum256([]byte(expectedIDPlaintext))
	expectedID := base64.StdEncoding.EncodeToString(hashed[:])
	if expectedID != id {
		t.Fatalf("wrong ID: expected %s, found %s", expectedID, id)
	}
	// Test with entityID
	entry = TokenEntry{EntityID: "entityFoo", NamespaceID: "namespaceFoo", Policies: []string{"bar", "baz", "foo", "banana"}}
	id, isTWE = entry.CreateClientID()
	if isTWE {
		t.Fatalf("token with entity should return false value in isTWE bool")
	}
	if id != "entityFoo" {
		t.Fatalf("client ID should be entity ID")
	}

	// Test without namespace
	entry = TokenEntry{Policies: []string{"bar", "baz", "foo", "banana"}}
	id, isTWE = entry.CreateClientID()
	if !isTWE {
		t.Fatalf("TWE token should return true value in isTWE bool")
	}
	expectedIDPlaintext = "banana" + string(SortedPoliciesTWEDelimiter) + "bar" +
		string(SortedPoliciesTWEDelimiter) + "baz" +
		string(SortedPoliciesTWEDelimiter) + "foo" + string(ClientIDTWEDelimiter)

	hashed = sha256.Sum256([]byte(expectedIDPlaintext))
	expectedID = base64.StdEncoding.EncodeToString(hashed[:])
	if expectedID != id {
		t.Fatalf("wrong ID: expected %s, found %s", expectedID, id)
	}

	// Test without policies
	entry = TokenEntry{NamespaceID: "namespaceFoo"}
	id, isTWE = entry.CreateClientID()
	if !isTWE {
		t.Fatalf("TWE token should return true value in isTWE bool")
	}
	expectedIDPlaintext = "namespaceFoo"

	hashed = sha256.Sum256([]byte(expectedIDPlaintext))
	expectedID = base64.StdEncoding.EncodeToString(hashed[:])
	if expectedID != id {
		t.Fatalf("wrong ID: expected %s, found %s", expectedID, id)
	}
}
