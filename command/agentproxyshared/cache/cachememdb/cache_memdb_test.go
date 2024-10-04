// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cachememdb

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

func testContextInfo() *ContextInfo {
	ctx, cancelFunc := context.WithCancel(context.Background())

	return &ContextInfo{
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	}
}

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCacheMemDB_Get(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Test invalid index name
	_, err = cache.Get("foo", "bar")
	if err == nil {
		t.Fatal("expected error")
	}

	// Test on empty cache
	index, err := cache.Get(IndexNameID, "foo")
	if err != ErrCacheItemNotFound {
		t.Fatal("expected cache item to be not found", err)
	}
	if index != nil {
		t.Fatalf("expected nil index, got: %v", index)
	}

	// Populate cache
	in := &Index{
		ID:            "test_id",
		Namespace:     "test_ns/",
		RequestPath:   "/v1/request/path",
		Token:         "test_token",
		TokenAccessor: "test_accessor",
		Lease:         "test_lease",
		Response:      []byte("hello world"),
		Tokens:        map[string]struct{}{},
	}

	if err := cache.Set(in); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		indexName   string
		indexValues []interface{}
	}{
		{
			"by_index_id",
			"id",
			[]interface{}{in.ID},
		},
		{
			"by_request_path",
			"request_path",
			[]interface{}{in.Namespace, in.RequestPath},
		},
		{
			"by_lease",
			"lease",
			[]interface{}{in.Lease},
		},
		{
			"by_token",
			"token",
			[]interface{}{in.Token},
		},
		{
			"by_token_accessor",
			"token_accessor",
			[]interface{}{in.TokenAccessor},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := cache.Get(tc.indexName, tc.indexValues...)
			if err != nil && err != ErrCacheItemNotFound {
				t.Fatal(err)
			}
			if diff := deep.Equal(in, out); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}

func TestCacheMemDB_GetByPrefix(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Test invalid index name
	_, err = cache.GetByPrefix("foo", "bar", "baz")
	if err == nil {
		t.Fatal("expected error")
	}

	// Test on empty cache
	index, err := cache.GetByPrefix(IndexNameRequestPath, "foo", "bar")
	if err != nil {
		t.Fatal(err)
	}
	if index != nil {
		t.Fatalf("expected nil index, got: %v", index)
	}

	// Populate cache
	in := &Index{
		ID:            "test_id",
		Namespace:     "test_ns/",
		RequestPath:   "/v1/request/path/1",
		Token:         "test_token",
		TokenParent:   "test_token_parent",
		TokenAccessor: "test_accessor",
		Lease:         "path/to/test_lease/1",
		LeaseToken:    "test_lease_token",
		Response:      []byte("hello world"),
	}

	if err := cache.Set(in); err != nil {
		t.Fatal(err)
	}

	// Populate cache
	in2 := &Index{
		ID:            "test_id_2",
		Namespace:     "test_ns/",
		RequestPath:   "/v1/request/path/2",
		Token:         "test_token2",
		TokenParent:   "test_token_parent",
		TokenAccessor: "test_accessor2",
		Lease:         "path/to/test_lease/2",
		LeaseToken:    "test_lease_token",
		Response:      []byte("hello world"),
	}

	if err := cache.Set(in2); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		indexName   string
		indexValues []interface{}
	}{
		{
			"by_request_path",
			IndexNameRequestPath,
			[]interface{}{"test_ns/", "/v1/request/path"},
		},
		{
			"by_lease",
			IndexNameLease,
			[]interface{}{"path/to/test_lease"},
		},
		{
			"by_token_parent",
			IndexNameTokenParent,
			[]interface{}{"test_token_parent"},
		},
		{
			"by_lease_token",
			IndexNameLeaseToken,
			[]interface{}{"test_lease_token"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := cache.GetByPrefix(tc.indexName, tc.indexValues...)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal([]*Index{in, in2}, out); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}

func TestCacheMemDB_Set(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		index   *Index
		wantErr bool
	}{
		{
			"nil",
			nil,
			true,
		},
		{
			"empty_fields",
			&Index{},
			true,
		},
		{
			"missing_required_fields",
			&Index{
				Lease: "foo",
			},
			true,
		},
		{
			"all_fields",
			&Index{
				ID:            "test_id",
				Namespace:     "test_ns/",
				RequestPath:   "/v1/request/path",
				Token:         "test_token",
				TokenAccessor: "test_accessor",
				Lease:         "test_lease",
				RenewCtxInfo:  testContextInfo(),
			},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := cache.Set(tc.index); (err != nil) != tc.wantErr {
				t.Fatalf("CacheMemDB.Set() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestCacheMemDB_Evict(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Test on empty cache
	if err := cache.Evict(IndexNameID, "foo"); err != nil {
		t.Fatal(err)
	}

	testIndex := &Index{
		ID:            "test_id",
		Namespace:     "test_ns/",
		RequestPath:   "/v1/request/path",
		Token:         "test_token",
		TokenAccessor: "test_token_accessor",
		Lease:         "test_lease",
		RenewCtxInfo:  testContextInfo(),
	}

	testCases := []struct {
		name        string
		indexName   string
		indexValues []interface{}
		insertIndex *Index
		wantErr     bool
	}{
		{
			"empty_params",
			"",
			[]interface{}{""},
			nil,
			true,
		},
		{
			"invalid_params",
			"foo",
			[]interface{}{"bar"},
			nil,
			true,
		},
		{
			"by_id",
			"id",
			[]interface{}{"test_id"},
			testIndex,
			false,
		},
		{
			"by_request_path",
			"request_path",
			[]interface{}{"test_ns/", "/v1/request/path"},
			testIndex,
			false,
		},
		{
			"by_token",
			"token",
			[]interface{}{"test_token"},
			testIndex,
			false,
		},
		{
			"by_token_accessor",
			"token_accessor",
			[]interface{}{"test_accessor"},
			testIndex,
			false,
		},
		{
			"by_lease",
			"lease",
			[]interface{}{"test_lease"},
			testIndex,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.insertIndex != nil {
				if err := cache.Set(tc.insertIndex); err != nil {
					t.Fatal(err)
				}
			}

			if err := cache.Evict(tc.indexName, tc.indexValues...); (err != nil) != tc.wantErr {
				t.Fatal(err)
			}

			// Verify that the cache doesn't contain the entry any more
			index, err := cache.Get(tc.indexName, tc.indexValues...)
			if err != ErrCacheItemNotFound && !tc.wantErr {
				t.Fatal("expected cache item to be not found", err)
			}
			if index != nil {
				t.Fatalf("expected nil entry, got = %#v", index)
			}
		})
	}
}

func TestCacheMemDB_Flush(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatal(err)
	}

	// Populate cache
	in := &Index{
		ID:          "test_id",
		Token:       "test_token",
		Lease:       "test_lease",
		Namespace:   "test_ns/",
		RequestPath: "/v1/request/path",
		Response:    []byte("hello world"),
	}

	if err := cache.Set(in); err != nil {
		t.Fatal(err)
	}

	// Reset the cache
	if err := cache.Flush(); err != nil {
		t.Fatal(err)
	}

	// Check the cache doesn't contain inserted index
	out, err := cache.Get(IndexNameID, "test_id")
	if err != ErrCacheItemNotFound {
		t.Fatal("expected cache item to be not found", err)
	}
	if out != nil {
		t.Fatalf("expected cache to be empty, got = %v", out)
	}
}

// TestCacheMemDB_EvictCapabilitiesIndex tests EvictCapabilitiesIndex works as expected.
func TestCacheMemDB_EvictCapabilitiesIndex(t *testing.T) {
	cache, err := New()
	require.Nil(t, err)

	// Test on empty cache
	err = cache.EvictCapabilitiesIndex(IndexNameID, "foo")
	require.Nil(t, err)

	capabilitiesIndex := &CapabilitiesIndex{
		ID:    "id",
		Token: "token",
	}

	err = cache.SetCapabilitiesIndex(capabilitiesIndex)
	require.Nil(t, err)

	err = cache.EvictCapabilitiesIndex(IndexNameID, capabilitiesIndex.ID)
	require.Nil(t, err)

	// Verify that the cache doesn't contain the entry anymore
	index, err := cache.GetCapabilitiesIndex(IndexNameID, capabilitiesIndex.ID)
	require.Equal(t, ErrCacheItemNotFound, err)
	require.Nil(t, index)
}

// TestCacheMemDB_GetCapabilitiesIndex tests GetCapabilitiesIndex works as expected.
func TestCacheMemDB_GetCapabilitiesIndex(t *testing.T) {
	cache, err := New()
	require.Nil(t, err)

	capabilitiesIndex := &CapabilitiesIndex{
		ID:    "id",
		Token: "token",
	}

	err = cache.SetCapabilitiesIndex(capabilitiesIndex)
	require.Nil(t, err)

	// Verify that we can retrieve the index
	index, err := cache.GetCapabilitiesIndex(IndexNameID, capabilitiesIndex.ID)
	require.Nil(t, err)
	require.Equal(t, capabilitiesIndex, index)

	// Verify behaviour on a non-existing ID
	index, err = cache.GetCapabilitiesIndex(IndexNameID, "not a real id")
	require.Equal(t, ErrCacheItemNotFound, err)
	require.Nil(t, index)

	// Verify behaviour with a non-existing index name
	index, err = cache.GetCapabilitiesIndex("not a real name", capabilitiesIndex.ID)
	require.NotNil(t, err)
}

// TestCacheMemDB_SetCapabilitiesIndex tests SetCapabilitiesIndex works as expected.
func TestCacheMemDB_SetCapabilitiesIndex(t *testing.T) {
	cache, err := New()
	require.Nil(t, err)

	capabilitiesIndex := &CapabilitiesIndex{
		ID:    "id",
		Token: "token",
	}

	err = cache.SetCapabilitiesIndex(capabilitiesIndex)
	require.Nil(t, err)

	// Verify we can retrieve the index
	index, err := cache.GetCapabilitiesIndex(IndexNameID, capabilitiesIndex.ID)
	require.Nil(t, err)
	require.Equal(t, capabilitiesIndex, index)

	// Verify behaviour on a nil index
	err = cache.SetCapabilitiesIndex(nil)
	require.NotNil(t, err)

	// Verify behaviour on an index without id
	err = cache.SetCapabilitiesIndex(&CapabilitiesIndex{
		Token: "token",
	})
	require.NotNil(t, err)

	// Verify behaviour on an index with only ID
	err = cache.SetCapabilitiesIndex(&CapabilitiesIndex{
		ID: "id",
	})
	require.Nil(t, err)
}
