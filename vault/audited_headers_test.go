// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

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
}

// List implements List from BarrierStorage interface.
// ignore-nil-nil-function-check.
func (m *mockStorage) List(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

// Get implements Get from BarrierStorage interface.
// ignore-nil-nil-function-check.
func (m *mockStorage) Get(_ context.Context, _ string) (*logical.StorageEntry, error) {
	return nil, nil
}

// Put implements Put from BarrierStorage interface.
func (m *mockStorage) Put(_ context.Context, _ *logical.StorageEntry) error {
	return nil
}

// Delete implements Delete from BarrierStorage interface.
func (m *mockStorage) Delete(_ context.Context, _ string) error {
	return nil
}

func mockAuditedHeadersConfig(t *testing.T) *AuditedHeadersConfig {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	return &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
		view:    view,
	}
}

func TestAuditedHeadersConfig_CRUD(t *testing.T) {
	conf := mockAuditedHeadersConfig(t)

	testAuditedHeadersConfig_Add(t, conf)
	testAuditedHeadersConfig_Remove(t, conf)
}

func testAuditedHeadersConfig_Add(t *testing.T, conf *AuditedHeadersConfig) {
	err := conf.add(context.Background(), "X-Test-Header", false)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok := conf.Headers["x-test-header"]
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

	headers := make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*auditedHeaderSettings{
		"x-test-header": {
			HMAC: false,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.add(context.Background(), "X-Vault-Header", true)
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	settings, ok = conf.Headers["x-vault-header"]
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

	headers = make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected["x-vault-header"] = &auditedHeaderSettings{
		HMAC: true,
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}
}

func testAuditedHeadersConfig_Remove(t *testing.T, conf *AuditedHeadersConfig) {
	err := conf.remove(context.Background(), "X-Test-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok := conf.Headers["x-Test-HeAder"]
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

	headers := make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected := map[string]*auditedHeaderSettings{
		"x-vault-header": {
			HMAC: true,
		},
	}

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}

	err = conf.remove(context.Background(), "x-VaulT-Header")
	if err != nil {
		t.Fatalf("Error when adding header to config: %s", err)
	}

	_, ok = conf.Headers["x-vault-header"]
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

	headers = make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		t.Fatalf("Error decoding header view: %s", err)
	}

	expected = make(map[string]*auditedHeaderSettings)

	if !reflect.DeepEqual(headers, expected) {
		t.Fatalf("Expected config didn't match actual. Expected: %#v, Got: %#v", expected, headers)
	}
}

type TestSalter struct{}

func (*TestSalter) Salt(ctx context.Context) (*salt.Salt, error) {
	return salt.NewSalt(ctx, nil, nil)
}

func TestAuditedHeadersConfig_ApplyConfig(t *testing.T) {
	conf := mockAuditedHeadersConfig(t)

	conf.add(context.Background(), "X-TesT-Header", false)
	conf.add(context.Background(), "X-Vault-HeAdEr", true)

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
	conf := mockAuditedHeadersConfig(t)

	conf.add(context.Background(), "X-TesT-Header", false)
	conf.add(context.Background(), "X-Vault-HeAdEr", true)

	reqHeaders := map[string][]string{}

	salter := &TestSalter{}

	result, err := conf.ApplyConfig(context.Background(), reqHeaders, salter)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("Expected no headers but actually got: %d\n", len(result))
	}
}

func TestAuditedHeadersConfig_ApplyConfig_NoConfiguredHeaders(t *testing.T) {
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
	conf := mockAuditedHeadersConfig(t)

	conf.add(context.Background(), "X-TesT-Header", false)
	conf.add(context.Background(), "X-Vault-HeAdEr", true)

	reqHeaders := map[string][]string{
		"X-Test-Header":  {"foo"},
		"X-Vault-Header": {"bar", "bar"},
		"Content-Type":   {"json"},
	}

	salter := &FailingSalter{}

	_, err := conf.ApplyConfig(context.Background(), reqHeaders, salter)
	if err == nil {
		t.Fatal("expected error from ApplyConfig")
	}
}

func BenchmarkAuditedHeaderConfig_ApplyConfig(b *testing.B) {
	conf := &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
		view:    nil,
	}

	conf.Headers = map[string]*auditedHeaderSettings{
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
		conf.ApplyConfig(context.Background(), reqHeaders, salter)
	}
}

// TestAuditedHeaders_auditedHeadersKey is used to check the key we use to handle
// invalidation doesn't change when we weren't expecting it to.
func TestAuditedHeaders_auditedHeadersKey(t *testing.T) {
	require.Equal(t, "audited-headers-config/audited-headers", auditedHeadersKey())
}

// TestAuditedHeaders_NewAuditedHeadersConfig checks supplying incorrect params to
// the constructor for AuditedHeadersConfig returns an error.
func TestAuditedHeaders_NewAuditedHeadersConfig(t *testing.T) {
	ac, err := NewAuditedHeadersConfig(nil)
	require.Error(t, err)
	require.Nil(t, ac)

	ac, err = NewAuditedHeadersConfig(&BarrierView{})
	require.NoError(t, err)
	require.NotNil(t, ac)
}

// TestAuditedHeaders_invalidate ensures that we can update the headers on AuditedHeadersConfig
// when we invalidate, and load the updated headers from the view/storage.
func TestAuditedHeaders_invalidate(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, auditedHeadersSubPath)
	ahc, err := NewAuditedHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.Headers, 0)

	// Store some data using the view.
	fakeHeaders1 := map[string]*auditedHeaderSettings{"x-magic-header": {}}
	fakeBytes1, err := json.Marshal(fakeHeaders1)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes1})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.invalidate(context.Background())
	require.NoError(t, err)
	require.Len(t, ahc.Headers, 1)
	_, ok := ahc.Headers["x-magic-header"]
	require.True(t, ok)

	// Do it again with more headers and random casing.
	fakeHeaders2 := map[string]*auditedHeaderSettings{
		"x-magic-header":           {},
		"x-even-MORE-magic-header": {},
	}
	fakeBytes2, err := json.Marshal(fakeHeaders2)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes2})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.invalidate(context.Background())
	require.NoError(t, err)
	require.Len(t, ahc.Headers, 2)
	_, ok = ahc.Headers["x-magic-header"]
	require.True(t, ok)
	_, ok = ahc.Headers["x-even-more-magic-header"]
	require.True(t, ok)
}

// TestAuditedHeaders_invalidate_nil_view ensures that we invalidate the headers
// correctly (clear them) when we get nil for the storage entry from the view.
func TestAuditedHeaders_invalidate_nil_view(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, auditedHeadersSubPath)
	ahc, err := NewAuditedHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.Headers, 0)

	// Store some data using the view.
	fakeHeaders1 := map[string]*auditedHeaderSettings{"x-magic-header": {}}
	fakeBytes1, err := json.Marshal(fakeHeaders1)
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: fakeBytes1})
	require.NoError(t, err)

	// Invalidate and check we now see the header we stored
	err = ahc.invalidate(context.Background())
	require.NoError(t, err)
	require.Len(t, ahc.Headers, 1)
	_, ok := ahc.Headers["x-magic-header"]
	require.True(t, ok)

	// Swap out the view with a mock that returns nil when we try to invalidate.
	// This should mean we end up just clearing the headers (no errors).
	mockStorageBarrier := new(mockStorage)
	mockStorageBarrier.On("Get", mock.Anything, mock.Anything).Return(nil, nil)
	ahc.view = NewBarrierView(mockStorageBarrier, auditedHeadersSubPath)

	// Invalidate should clear out the existing headers without error
	err = ahc.invalidate(context.Background())
	require.NoError(t, err)
	require.Len(t, ahc.Headers, 0)
}

// TestAuditedHeaders_invalidate_bad_data ensures that we correctly error if the
// underlying data cannot be parsed as expected.
func TestAuditedHeaders_invalidate_bad_data(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, auditedHeadersSubPath)
	ahc, err := NewAuditedHeadersConfig(view)
	require.NoError(t, err)
	require.Len(t, ahc.Headers, 0)

	// Store some bad data using the view.
	badBytes, err := json.Marshal("i am bad")
	require.NoError(t, err)
	err = view.Put(context.Background(), &logical.StorageEntry{Key: auditedHeadersEntry, Value: badBytes})
	require.NoError(t, err)

	// Invalidate should
	err = ahc.invalidate(context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to parse config")
}
