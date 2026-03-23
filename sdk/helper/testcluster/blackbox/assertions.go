// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

type SecretAssertion struct {
	t      *testing.T
	secret *api.Secret
}

type MapAssertion struct {
	t    *testing.T
	data map[string]any
	path string
}

type SliceAssertion struct {
	t    *testing.T
	data []any
	path string
}

func (s *Session) AssertSecret(secret *api.Secret) *SecretAssertion {
	s.t.Helper()
	require.NotNil(s.t, secret)
	return &SecretAssertion{t: s.t, secret: secret}
}

func (sa *SecretAssertion) Data() *MapAssertion {
	sa.t.Helper()
	require.NotNil(sa.t, sa.secret.Data)
	return &MapAssertion{t: sa.t, data: sa.secret.Data, path: "Data"}
}

func (sa *SecretAssertion) KV2() *MapAssertion {
	sa.t.Helper()

	require.NotNil(sa.t, sa.secret.Data)
	inner, ok := sa.secret.Data["data"]
	if !ok {
		sa.t.Fatal("data not found in secret")
	}

	innerMap, ok := inner.(map[string]any)
	if !ok {
		sa.t.Fatalf("expected 'data' to be a map, got %T", inner)
	}

	return &MapAssertion{t: sa.t, data: innerMap, path: "Data.data"}
}

func (ma *MapAssertion) HasKey(key string, expected any) *MapAssertion {
	ma.t.Helper()

	val, ok := ma.data[key]
	if !ok {
		ma.t.Fatalf("[%s] missing expected key: %q", ma.path, key)
	}

	if !smartCompare(val, expected) {
		ma.t.Fatalf("[%s] key %q:\n\texpected: %v\n\tgot: %v", ma.path, key, expected, val)
	}

	return ma
}

func (ma *MapAssertion) HasKeyCustom(key string, f func(val any) bool) *MapAssertion {
	ma.t.Helper()

	val, ok := ma.data[key]
	if !ok {
		ma.t.Fatalf("[%s] missing expected key: %q", ma.path, key)
	}

	okAgain := f(val)
	if !okAgain {
		ma.t.Fatalf("[%s] key %q failed custom check", ma.path, key)
	}

	return ma
}

func (ma *MapAssertion) HasKeyExists(key string) *MapAssertion {
	ma.t.Helper()

	if _, ok := ma.data[key]; !ok {
		ma.t.Fatalf("[%s] missing expected key: %q", ma.path, key)
	}

	return ma
}

func (ma *MapAssertion) GetMap(key string) *MapAssertion {
	ma.t.Helper()

	val, ok := ma.data[key]
	if !ok {
		ma.t.Fatalf("[%s] missing expected key: %q", ma.path, key)
	}

	nestedMap, ok := val.(map[string]any)
	if !ok {
		ma.t.Fatalf("[%s] key %q is not a map, it is %T", ma.path, key, val)
	}

	return &MapAssertion{
		t:    ma.t,
		data: nestedMap,
		path: ma.path + "." + key,
	}
}

func (ma *MapAssertion) GetSlice(key string) *SliceAssertion {
	ma.t.Helper()

	val, ok := ma.data[key]
	if !ok {
		ma.t.Fatalf("[%s] missing expected key: %q", ma.path, key)
	}

	slice, ok := val.([]any)
	if !ok {
		ma.t.Fatalf("[%s] key %q is not a slice, it is %T", ma.path, key, val)
	}

	return &SliceAssertion{
		t:    ma.t,
		data: slice,
		path: ma.path,
	}
}

func (sa *SliceAssertion) Length(expected int) *SliceAssertion {
	sa.t.Helper()

	if len(sa.data) != expected {
		sa.t.Fatalf("[%s] expected slice length %d, got %d", sa.path, expected, len(sa.data))
	}

	return sa
}

func (sa *SliceAssertion) FindMap(key string, expectedValue any) *MapAssertion {
	sa.t.Helper()

	for i, item := range sa.data {
		// we expect the slice to contain maps
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		// check if this map has the key/value we're looking for
		if val, exists := m[key]; exists {
			if smartCompare(val, expectedValue) {
				return &MapAssertion{
					t:    sa.t,
					data: m,
					path: fmt.Sprintf("%s[%d]", sa.path, i),
				}
			}
		}
	}

	sa.t.Fatalf("[%s] could not find element with %q == %v", sa.path, key, expectedValue)
	return nil
}

func (sa *SliceAssertion) AllHaveKey(key string, expectedValue any) *SliceAssertion {
	sa.t.Helper()

	for i, item := range sa.data {
		m, ok := item.(map[string]any)
		if !ok {
			sa.t.Fatalf("[%s[%d]] expected element to be a map, got %T", sa.path, i, item)
		}

		val, exists := m[key]
		if !exists {
			sa.t.Fatalf("[%s[%d]] missing expected key: %q", sa.path, i, key)
		}

		if !smartCompare(val, expectedValue) {
			sa.t.Fatalf("[%s[%d]] key %q mismatch:\n\texpected: %v\n\tgot: %v", sa.path, i, key, expectedValue, val)
		}
	}

	return sa
}

// AllHaveKeyCustom asserts that every element in the slice is a map
// containing the key, and that the provided function returns true for the value.
func (sa *SliceAssertion) AllHaveKeyCustom(key string, check func(val any) bool) *SliceAssertion {
	sa.t.Helper()

	for i, item := range sa.data {
		m, ok := item.(map[string]any)
		if !ok {
			sa.t.Fatalf("[%s[%d]] expected element to be a map, got %T", sa.path, i, item)
		}

		val, exists := m[key]
		if !exists {
			sa.t.Fatalf("[%s[%d]] missing expected key: %q", sa.path, i, key)
		}

		if !check(val) {
			sa.t.Fatalf("[%s[%d]] key %q failed custom check. Value was: %v", sa.path, i, key, val)
		}
	}

	return sa
}

// NoneHaveKeyVal asserts that NO element in the slice contains the specific key/value pair.
// It succeeds if the key is missing, or if the key is present but has a different value.
func (sa *SliceAssertion) NoneHaveKeyVal(key string, restrictedValue any) *SliceAssertion {
	sa.t.Helper()

	for i, item := range sa.data {
		m, ok := item.(map[string]any)
		if !ok {
			sa.t.Fatalf("[%s[%d]] expected element to be a map, got %T", sa.path, i, item)
		}

		if val, exists := m[key]; exists {
			if smartCompare(val, restrictedValue) {
				sa.t.Fatalf("[%s[%d]] found restricted key/value pair: %q: %v", sa.path, i, key, val)
			}
		}
	}

	return sa
}

// smartCompare is designed to get around the weird stuff that happens when Vault's API sometimes
// returns json.Number, sometimes strings with numbers, sometimes actual numbers. It's a mess.
func smartCompare(actual, expected any) bool {
	// if they match exactly (type and value), we are done.
	if reflect.DeepEqual(actual, expected) {
		return true
	}

	// if actual is NOT a json.Number, and step 1 failed, they aren't equal.
	jNum, isJSON := actual.(json.Number)
	if !isJSON {
		return false
	}

	switch v := expected.(type) {
	case int:
		// user expects an int (e.g., HasKey("count", 5))
		// json.Number stores as string. we convert to int64, then cast to int.
		i64, err := jNum.Int64()
		if err != nil {
			return false // not a valid integer
		}
		return int(i64) == v

	case int64:
		i64, err := jNum.Int64()
		if err != nil {
			return false
		}
		return i64 == v

	case float64:
		// user expects float (e.g., HasKey("ttl", 1.5))
		f64, err := jNum.Float64()
		if err != nil {
			return false
		}
		return f64 == v

	case string:
		// user expects string (e.g. huge ID), just compare string to string
		return jNum.String() == v
	}

	return false
}
