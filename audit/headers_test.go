// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockStorage is a struct that is used to mock barrier storage.
type mockStorage struct {
	mock.Mock
	v map[string][]byte
}

// List implements List from BarrierStorage interface.
// ignore-nil-nil-function-check.
func (m *mockStorage) List(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

// Get implements Get from BarrierStorage interface.
// ignore-nil-nil-function-check.
func (m *mockStorage) Get(_ context.Context, key string) (*logical.StorageEntry, error) {
	b, ok := m.v[key]
	if !ok {
		return nil, nil
	}

	var entry *logical.StorageEntry
	err := json.Unmarshal(b, &entry)

	return entry, err
}

// Put implements Put from BarrierStorage interface.
func (m *mockStorage) Put(_ context.Context, entry *logical.StorageEntry) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	m.v[entry.Key] = b

	return nil
}

// Delete implements Delete from BarrierStorage interface.
func (m *mockStorage) Delete(_ context.Context, _ string) error {
	return nil
}

func newMockStorage(t *testing.T) *mockStorage {
	t.Helper()

	return &mockStorage{
		Mock: mock.Mock{},
		v:    make(map[string][]byte),
	}
}

func mockAuditedHeadersConfig(t *testing.T) *HeadersConfig {
	return &HeadersConfig{
		headerSettings: make(map[string]*HeaderSettings),
		view:           newMockStorage(t),
	}
}

func TestAuditedHeadersConfig_CRUD(t *testing.T) {
	t.Parallel()

	conf := mockAuditedHeadersConfig(t)

	testAddHeaders(t, conf)
	testRemoveHeaders(t, conf)
}

func testAddHeaders(t *testing.T, conf *HeadersConfig) {
	t.Helper()

	err := conf.Add(context.Background(), "X-Test-Header", false)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok := conf.headerSettings["x-test-header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if settings.HMAC {
		t.Fatal("Expected HMAC to be set to false, got true")
	}

	out, err := conf.view.Get(context.Background(), auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}
	if out == nil {
		t.Fatal("nil value")
	}

	headers := make(map[string]*HeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*HeaderSettings{
		"x-test-header": {
			HMAC: false,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.Add(context.Background(), "X-Vault-Header", true)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok = conf.headerSettings["x-vault-header"]
	if !ok {
		t.Fatal("Expected header to be found in config")
	}

	if !settings.HMAC {
		t.Fatal("Expected HMAC to be set to true, got false")
	}

	out, err = conf.view.Get(context.Background(), auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}
	if out == nil {
		t.Fatal("nil value")
	}

	headers = make(map[string]*HeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected["x-vault-header"] = &HeaderSettings{
		HMAC: true,
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}
}

func testRemoveHeaders(t *testing.T, conf *HeadersConfig) {
	t.Helper()

	err := conf.Remove(context.Background(), "X-Test-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok := conf.headerSettings["x-Test-HeAder"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err := conf.view.Get(context.Background(), auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}
	if out == nil {
		t.Fatal("nil value")
	}

	headers := make(map[string]*HeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*HeaderSettings{
		"x-vault-header": {
			HMAC: true,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.Remove(context.Background(), "x-VaulT-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok = conf.headerSettings["x-vault-header"]
	if ok {
		t.Fatal("Expected header to not be found in config")
	}

	out, err = conf.view.Get(context.Background(), auditedHeadersEntry)
	if err != nil {
		t.Fatalf("Could not retrieve headers entry from config: %s", err)
	}
	if out == nil {
		t.Fatal("nil value")
	}

	headers = make(map[string]*HeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected = make(map[string]*HeaderSettings)

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}
}

func TestAuditedHeadersConfig_ApplyConfig(t *testing.T) {
	t.Parallel()

	conf := mockAuditedHeadersConfig(t)

	err := conf.Add(context.Background(), "X-TesT-Header", false)
	require.NoError(t, err)
	err = conf.Add(context.Background(), "X-Vault-HeAdEr", true)
	require.NoError(t, err)

	reqHeaders := map[string][]string{
		"X-Test-Header":  {"foo"},
		"X-Vault-Header": {"bar", "bar"},
		"Content-Type":   {"json"},
	}

	salter := &TestSalter{}

	result, err := conf.ApplyConfig(context.Background(), reqHeaders, salter)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string][]string{
		"x-test-header":  {"foo"},
		"x-vault-header": {"hmac-sha256:", "hmac-sha256:"},
	}

	if len(expected) != len(result) {
		t.Fatalf("Expected headers count did not match actual count: Expected count %d\n Got %d\n", len(expected), len(result))
	}

	for resultKey, resultValues := range result {
		expectedValues := expected[resultKey]

		if len(expectedValues) != len(resultValues) {
			t.Fatalf("Expected header values count did not match actual values count: Expected count: %d\n Got %d\n", len(expectedValues), len(resultValues))
		}

		for i, e := range expectedValues {
			if e == "hmac-sha256:" {
				if !strings.HasPrefix(resultValues[i], e) {
					t.Fatalf("Expected headers did not match actual: Expected %#v...\n Got %#v\n", e, resultValues[i])
				}
			} else {
				if e != resultValues[i] {
					t.Fatalf("Expected headers did not match actual: Expected %#v\n Got %#v\n", e, resultValues[i])
				}
			}
		}
	}

	// Make sure we didn't edit the reqHeaders map
	reqHeadersCopy := map[string][]string{
		"X-Test-Header":  {"foo"},
		"X-Vault-Header": {"bar", "bar"},
		"Content-Type":   {"json"},
	}

	if !reflect.DeepEqual(reqHeaders, reqHeadersCopy) {
		t.Fatalf("Req headers were changed, expected %#v\n got %#v", reqHeadersCopy, reqHeaders)
	}
}

// TestAuditedHeadersConfig_ApplyConfig_NoHeaders tests the case where there are
// no headers in the request.
func TestAuditedHeadersConfig_ApplyConfig_NoRequestHeaders(t *testing.T) {
	t.Parallel()

	conf := mockAuditedHeadersConfig(t)

	err := conf.Add(context.Background(), "X-TesT-Header", false)
	require.NoError(t, err)
	err = conf.Add(context.Background(), "X-Vault-HeAdEr", true)
	require.NoError(t, err)

	salter := &TestSalter{}

	// Test sending in nil headers first.
	result, err := conf.ApplyConfig(context.Background(), nil, salter)
	require.NoError(t, err)
	require.NotNil(t, result)

	result, err = conf.ApplyConfig(context.Background(), map[string][]string{}, salter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 0)
}

func TestAuditedHeadersConfig_ApplyConfig_NoConfiguredHeaders(t *testing.T) {
	t.Parallel()

	conf := mockAuditedHeadersConfig(t)

	reqHeaders := map[string][]string{
		"X-Test-Header":  {"foo"},
		"X-Vault-Header": {"bar", "bar"},
		"Content-Type":   {"json"},
	}

	salter := &TestSalter{}

	result, err := conf.ApplyConfig(context.Background(), reqHeaders, salter)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("Expected no headers but actually got: %d\n", len(result))
	}

	// Make sure we didn't edit the reqHeaders map
	reqHeadersCopy := map[string][]string{
		"X-Test-Header":  {"foo"},
		"X-Vault-Header": {"bar", "bar"},
		"Content-Type":   {"json"},
	}

	if !reflect.DeepEqual(reqHeaders, reqHeadersCopy) {
		t.Fatalf("Req headers were changed, expected %#v\n got %#v", reqHeadersCopy, reqHeaders)
	}
}

// FailingSalter is an implementation of the Salter interface where the Salt
// method always returns an error.
type FailingSalter struct{}

// Salt always returns an error.
func (s *FailingSalter) Salt(context.Context) (*salt.Salt, error) {
	return nil, errors.New("testing error")
}

// TestAuditedHeadersConfig_ApplyConfig_HashStringError tests the case where
// an error is returned from HashString instead of a map of headers.
func TestAuditedHeadersConfig_ApplyConfig_HashStringError(t *testing.T) {
	t.Parallel()

	conf := mockAuditedHeadersConfig(t)

	err := conf.Add(context.Background(), "X-TesT-Header", false)
	require.NoError(t, err)
	err = conf.Add(context.Background(), "X-Vault-HeAdEr", true)
	require.NoError(t, err)

	reqHeaders := map[string][]string{
		"X-Test-Header":  {"foo"},
		"X-Vault-Header": {"bar", "bar"},
		"Content-Type":   {"json"},
	}

	salter := &FailingSalter{}

	_, err = conf.ApplyConfig(context.Background(), reqHeaders, salter)
	if err == nil {
		t.Fatal("expected error from ApplyConfig")
	}
}

func BenchmarkAuditedHeaderConfig_ApplyConfig(b *testing.B) {
	conf := &HeadersConfig{
		headerSettings: make(map[string]*HeaderSettings),
		view:           nil,
	}

	conf.headerSettings = map[string]*HeaderSettings{
		"X-Test-Header":  {false},
		"X-Vault-Header": {true},
	}

	reqHeaders := map[string][]string{
		"X-Test-Header":  {"foo"},
		"X-Vault-Header": {"bar", "bar"},
		"Content-Type":   {"json"},
	}

	salter := &TestSalter{}

	// Reset the timer since we did a lot above
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := conf.ApplyConfig(context.Background(), reqHeaders, salter)
		require.NoError(b, err)
	}
}

// TestAuditedHeaders_auditedHeadersKey is used to check the key we use to handle
// invalidation doesn't change when we weren't expecting it to.
func TestAuditedHeaders_auditedHeadersKey(t *testing.T) {
	t.Parallel()

	require.Equal(t, "audited-headers-config/audited-headers", AuditedHeadersKey())
}

// TestAuditedHeaders_NewAuditedHeadersConfig checks supplying incorrect params to
// the constructor for HeadersConfig returns an error.
func TestAuditedHeaders_NewAuditedHeadersConfig(t *testing.T) {
	t.Parallel()

	ac, err := NewHeadersConfig(nil)
	require.Error(t, err)
	require.Nil(t, ac)

	ac, err = NewHeadersConfig(newMockStorage(t))
	require.NoError(t, err)
	require.NotNil(t, ac)
}

// TestAuditedHeaders_invalidate ensures that we can update the headers on HeadersConfig
// when we invalidate, and load the updated headers from the view/storage.
func TestAuditedHeaders_invalidate(t *testing.T) {
	t.Parallel()

	view := newMockStorage(t)
	ahc, err := NewHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 0)

	// Store some data using the view.
	fakeHeaders1 := map[string]*HeaderSettings{"x-magic-header": {}}
	fakeBytes1, err := json.Marshal(fakeHeaders1)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes1})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.Invalidate(context.Background())
	require.NoError(t, err)
	require.Equal(t, len(ahc.DefaultHeaders())+1, len(ahc.headerSettings)) // (defaults + 1).
	_, ok := ahc.headerSettings["x-magic-header"]
	require.True(t, ok)

	// Do it again with more headers and random casing.
	fakeHeaders2 := map[string]*HeaderSettings{
		"x-magic-header":           {},
		"x-even-MORE-magic-header": {},
	}
	fakeBytes2, err := json.Marshal(fakeHeaders2)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes2})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.Invalidate(context.Background())
	require.NoError(t, err)
	require.Equal(t, len(ahc.DefaultHeaders())+2, len(ahc.headerSettings)) // (defaults + 2 new headers)
	_, ok = ahc.headerSettings["x-magic-header"]
	require.True(t, ok)
	_, ok = ahc.headerSettings["x-even-more-magic-header"]
	require.True(t, ok)
}

// TestAuditedHeaders_invalidate_nil_view ensures that we invalidate the headers
// correctly (clear them) when we get nil for the storage entry from the view.
func TestAuditedHeaders_invalidate_nil_view(t *testing.T) {
	t.Parallel()

	view := newMockStorage(t)
	ahc, err := NewHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 0)

	// Store some data using the view.
	fakeHeaders1 := map[string]*HeaderSettings{"x-magic-header": {}}
	fakeBytes1, err := json.Marshal(fakeHeaders1)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes1})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.Invalidate(context.Background())
	require.NoError(t, err)
	require.Equal(t, len(ahc.DefaultHeaders())+1, len(ahc.headerSettings)) // defaults + 1
	_, ok := ahc.headerSettings["x-magic-header"]
	require.True(t, ok)

	// Swap out the view with a mock that returns nil when we try to invalidate.
	// This should mean we end up just clearing the headers (no errors).
	mockStorageBarrier := newMockStorage(t)
	mockStorageBarrier.On("Get", mock.Anything, mock.Anything).Return(nil, nil)
	ahc.view = mockStorageBarrier
	// ahc.view = NewBarrierView(mockStorageBarrier, AuditedHeadersSubPath)

	// Invalidate should clear out the existing headers without error
	err = ahc.Invalidate(context.Background())
	require.NoError(t, err)
	require.Equal(t, len(ahc.DefaultHeaders()), len(ahc.headerSettings)) // defaults
}

// TestAuditedHeaders_invalidate_bad_data ensures that we correctly error if the
// underlying data cannot be parsed as expected.
func TestAuditedHeaders_invalidate_bad_data(t *testing.T) {
	t.Parallel()

	view := newMockStorage(t)
	ahc, err := NewHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 0)

	// Store some bad data using the view.
	badBytes, err := json.Marshal("i am bad")
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: badBytes})
	require.NoError(t, err)

	// Invalidate should
	err = ahc.Invalidate(context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to parse config")
}

// TestAuditedHeaders_header checks we can return a copy of settings associated with
// an existing header, and we also know when a header wasn't found.
func TestAuditedHeaders_header(t *testing.T) {
	t.Parallel()

	view := newMockStorage(t)
	ahc, err := NewHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 0)

	err = ahc.Add(context.Background(), "juan", true)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 1)

	s, ok := ahc.Header("juan")
	require.True(t, ok)
	require.Equal(t, true, s.HMAC)

	s, ok = ahc.Header("x-magic-token")
	require.False(t, ok)
}

// TestAuditedHeaders_headers checks we are able to return a copy of the existing
// configured headers.
func TestAuditedHeaders_headers(t *testing.T) {
	t.Parallel()

	view := newMockStorage(t)
	ahc, err := NewHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 0)

	err = ahc.Add(context.Background(), "juan", true)
	require.NoError(t, err)
	err = ahc.Add(context.Background(), "john", false)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 2)

	s := ahc.Headers()
	require.Len(t, s, 2)
	require.Equal(t, true, s["juan"].HMAC)
	require.Equal(t, false, s["john"].HMAC)
}

// TestAuditedHeaders_invalidate_defaults checks that we ensure any 'default' headers
// are present after invalidation, and if they were loaded from storage then they
// do not get overwritten with our defaults.
func TestAuditedHeaders_invalidate_defaults(t *testing.T) {
	t.Parallel()

	view := newMockStorage(t)
	ahc, err := NewHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.headerSettings, 0)

	// Store some data using the view.
	fakeHeaders1 := map[string]*HeaderSettings{"x-magic-header": {}}
	fakeBytes1, err := json.Marshal(fakeHeaders1)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes1})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.Invalidate(context.Background())
	require.NoError(t, err)
	require.Equal(t, len(ahc.DefaultHeaders())+1, len(ahc.headerSettings)) // (defaults + 1 new header)
	_, ok := ahc.headerSettings["x-magic-header"]
	require.True(t, ok)
	s, ok := ahc.headerSettings["x-correlation-id"]
	require.True(t, ok)
	require.False(t, s.HMAC)

	// Add correlation ID specifically with HMAC and make sure it doesn't get blasted away.
	fakeHeaders1 = map[string]*HeaderSettings{"x-magic-header": {}, "X-Correlation-ID": {HMAC: true}}
	fakeBytes1, err = json.Marshal(fakeHeaders1)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes1})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.Invalidate(context.Background())
	require.NoError(t, err)
	require.Equal(t, len(ahc.DefaultHeaders())+1, len(ahc.headerSettings)) // (defaults + 1 new header, 1 is also a default)
	_, ok = ahc.headerSettings["x-magic-header"]
	require.True(t, ok)
	s, ok = ahc.headerSettings["x-correlation-id"]
	require.True(t, ok)
	require.True(t, s.HMAC)
}
