// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cacheboltdb"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/command/agentproxyshared/cache/keymanager"
	"github.com/hashicorp/vault/helper/useragent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func testNewLeaseCache(t *testing.T, responses []*SendResponse) *LeaseCache {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	lc, err := NewLeaseCache(&LeaseCacheConfig{
		Client:              client,
		BaseContext:         context.Background(),
		Proxier:             NewMockProxier(responses),
		Logger:              logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
		CacheStaticSecrets:  true,
		CacheDynamicSecrets: true,
		UserAgentToUse:      "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	return lc
}

func testNewLeaseCacheWithDelay(t *testing.T, cacheable bool, delay int) *LeaseCache {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	lc, err := NewLeaseCache(&LeaseCacheConfig{
		Client:              client,
		BaseContext:         context.Background(),
		Proxier:             &mockDelayProxier{cacheable, delay},
		Logger:              logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
		CacheStaticSecrets:  true,
		CacheDynamicSecrets: true,
		UserAgentToUse:      "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	return lc
}

func testNewLeaseCacheWithPersistence(t *testing.T, responses []*SendResponse, storage *cacheboltdb.BoltStorage) *LeaseCache {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	require.NoError(t, err)

	lc, err := NewLeaseCache(&LeaseCacheConfig{
		Client:              client,
		BaseContext:         context.Background(),
		Proxier:             NewMockProxier(responses),
		Logger:              logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
		Storage:             storage,
		CacheStaticSecrets:  true,
		CacheDynamicSecrets: true,
		UserAgentToUse:      "test",
	})
	require.NoError(t, err)

	return lc
}

func TestCache_ComputeIndexID(t *testing.T) {
	tests := []struct {
		name    string
		req     *SendRequest
		want    string
		wantErr bool
	}{
		{
			"basic",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "test",
					},
				},
			},
			"7b5db388f211fd9edca8c6c254831fb01ad4e6fe624dbb62711f256b5e803717",
			false,
		},
		{
			"ignore consistency headers",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "test",
					},
					Header: http.Header{
						vaulthttp.VaultIndexHeaderName:        []string{"foo"},
						vaulthttp.VaultInconsistentHeaderName: []string{"foo"},
						vaulthttp.VaultForwardHeaderName:      []string{"foo"},
					},
				},
			},
			"7b5db388f211fd9edca8c6c254831fb01ad4e6fe624dbb62711f256b5e803717",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeIndexID(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("actual_error: %v, expected_error: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, string(tt.want)) {
				t.Errorf("bad: index id; actual: %q, expected: %q", got, string(tt.want))
			}
		})
	}
}

// TestCache_ComputeStaticSecretIndexID ensures that
// computeStaticSecretCacheIndex works correctly. If this test breaks, then our
// hashing algorithm has changed, and we risk breaking backwards compatibility.
func TestCache_ComputeStaticSecretIndexID(t *testing.T) {
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/foo/bar",
			},
		},
	}

	index := computeStaticSecretCacheIndex(req)
	// We expect this to be "", as it doesn't start with /v1
	expectedIndex := ""
	require.Equal(t, expectedIndex, index)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/foo/bar",
			},
		},
	}

	expectedIndex = "b117a962f19f17fa372c8681cadcd6fd370d28ee6e0a7012196b780bef601b53"
	index2 := computeStaticSecretCacheIndex(req)
	require.Equal(t, expectedIndex, index2)
}

// Test_GetStaticSecretPathFromRequestNoNamespaces tests that getStaticSecretPathFromRequest
// behaves as expected when no namespaces are involved.
func Test_GetStaticSecretPathFromRequestNoNamespaces(t *testing.T) {
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/foo/bar",
			},
		},
	}

	path := getStaticSecretPathFromRequest(req)
	require.Equal(t, "foo/bar", path)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				// Paths like this are not static secrets, so we should return ""
				Path: "foo/bar",
			},
		},
	}

	path = getStaticSecretPathFromRequest(req)
	require.Equal(t, "", path)
}

// Test_GetStaticSecretPathFromRequestNamespaces tests that getStaticSecretPathFromRequest
// behaves as expected when namespaces are involved.
func Test_GetStaticSecretPathFromRequestNamespaces(t *testing.T) {
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/foo/bar",
			},
			Header: map[string][]string{api.NamespaceHeaderName: {"ns1"}},
		},
	}

	path := getStaticSecretPathFromRequest(req)
	require.Equal(t, "ns1/foo/bar", path)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/ns1/foo/bar",
			},
		},
	}

	path = getStaticSecretPathFromRequest(req)
	require.Equal(t, "ns1/foo/bar", path)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				// Paths like this are not static secrets, so we should return ""
				Path: "ns1/foo/bar",
			},
		},
	}

	path = getStaticSecretPathFromRequest(req)
	require.Equal(t, "", path)
}

// TestCache_CanonicalizeStaticSecretPath ensures that
// canonicalizeStaticSecretPath works as expected with all kinds of inputs.
func TestCache_CanonicalizeStaticSecretPath(t *testing.T) {
	expected := "foo/bar"
	actual := canonicalizeStaticSecretPath("/v1/foo/bar", "")
	require.Equal(t, expected, actual)

	actual = canonicalizeStaticSecretPath("foo/bar", "")
	require.Equal(t, expected, actual)
	actual = canonicalizeStaticSecretPath("/foo/bar", "")
	require.Equal(t, expected, actual)

	expected = "ns1/foo/bar"
	actual = canonicalizeStaticSecretPath("/v1/ns1/foo/bar", "")
	require.Equal(t, expected, actual)

	actual = canonicalizeStaticSecretPath("ns1/foo/bar", "")
	require.Equal(t, expected, actual)
	actual = canonicalizeStaticSecretPath("/ns1/foo/bar", "")
	require.Equal(t, expected, actual)

	expected = "ns1/foo/bar"
	actual = canonicalizeStaticSecretPath("/v1/foo/bar", "ns1")
	require.Equal(t, expected, actual)

	actual = canonicalizeStaticSecretPath("/foo/bar", "ns1")
	require.Equal(t, expected, actual)
	actual = canonicalizeStaticSecretPath("foo/bar", "ns1")
	require.Equal(t, expected, actual)

	expected = "ns1/foo/bar"
	actual = canonicalizeStaticSecretPath("/v1/foo/bar", "ns1/")
	require.Equal(t, expected, actual)

	actual = canonicalizeStaticSecretPath("/foo/bar", "ns1/")
	require.Equal(t, expected, actual)
	actual = canonicalizeStaticSecretPath("foo/bar", "ns1/")
	require.Equal(t, expected, actual)

	expected = "ns1/foo/bar"
	actual = canonicalizeStaticSecretPath("/v1/foo/bar", "/ns1/")
	require.Equal(t, expected, actual)

	actual = canonicalizeStaticSecretPath("/foo/bar", "/ns1/")
	require.Equal(t, expected, actual)
	actual = canonicalizeStaticSecretPath("foo/bar", "/ns1/")
	require.Equal(t, expected, actual)
}

// TestCache_ComputeStaticSecretIndexIDNamespaces ensures that
// computeStaticSecretCacheIndex correctly identifies that a request
// with a namespace header and a request specifying the namespace in the path
// are equivalent.
func TestCache_ComputeStaticSecretIndexIDNamespaces(t *testing.T) {
	req := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "foo/bar",
			},
			Header: map[string][]string{api.NamespaceHeaderName: {"ns1"}},
		},
	}

	index := computeStaticSecretCacheIndex(req)
	// Paths like this are not static secrets, so we should expect ""
	require.Equal(t, "", index)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "ns1/foo/bar",
			},
		},
	}

	// Paths like this are not static secrets, so we should expect ""
	index2 := computeStaticSecretCacheIndex(req)
	require.Equal(t, "", index2)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/ns1/foo/bar",
			},
		},
	}

	expectedIndex := "a4605679d269aa1bebac7079a471a33403413f388f63bf0da3c771b225857932"
	// We expect that computeStaticSecretCacheIndex will compute the same index
	index3 := computeStaticSecretCacheIndex(req)
	require.Equal(t, expectedIndex, index3)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/foo/bar",
			},
			Header: map[string][]string{api.NamespaceHeaderName: {"ns1"}},
		},
	}

	index4 := computeStaticSecretCacheIndex(req)
	require.Equal(t, expectedIndex, index4)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/foo/bar",
			},
			Header: map[string][]string{api.NamespaceHeaderName: {"ns1/"}},
		},
	}

	// Paths like this are not static secrets, so we should expect ""
	index5 := computeStaticSecretCacheIndex(req)
	require.Equal(t, "", index5)

	req = &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/foo/bar",
			},
			Header: map[string][]string{api.NamespaceHeaderName: {"ns1/"}},
		},
	}

	index6 := computeStaticSecretCacheIndex(req)
	require.Equal(t, expectedIndex, index6)
}

func TestLeaseCache_EmptyToken(t *testing.T) {
	responses := []*SendResponse{
		newTestSendResponse(http.StatusCreated, `{"value": "invalid", "auth": {"client_token": "testtoken"}}`),
	}
	lc := testNewLeaseCache(t, responses)

	// Even if the send request doesn't have a token on it, a successful
	// cacheable response should result in the index properly getting populated
	// with a token and memdb shouldn't complain while inserting the index.
	urlPath := "http://example.com/v1/sample/api"
	sendReq := &SendRequest{
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
	}
	resp, err := lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("expected a non empty response")
	}
}

func TestLeaseCache_SendCacheable(t *testing.T) {
	// Emulate 2 responses from the api proxy. One returns a new token and the
	// other returns a lease.
	responses := []*SendResponse{
		newTestSendResponse(http.StatusCreated, `{"auth": {"client_token": "testtoken", "renewable": true}}`),
		newTestSendResponse(http.StatusOK, `{"lease_id": "foo", "renewable": true, "data": {"value": "foo"}}`),
	}

	lc := testNewLeaseCache(t, responses)
	// Register a token so that the token and lease requests are cached
	require.NoError(t, lc.RegisterAutoAuthToken("autoauthtoken"))

	// Make a request. A response with a new token is returned to the lease
	// cache and that will be cached.
	urlPath := "http://example.com/v1/sample/api"
	sendReq := &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
	}
	resp, err := lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, responses[0].Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Send the same request again to get the cached response
	sendReq = &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, responses[0].Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Check TokenParent
	cachedItem, err := lc.db.Get(cachememdb.IndexNameToken, "testtoken")
	if err != nil {
		t.Fatal(err)
	}
	if cachedItem == nil {
		t.Fatalf("expected token entry from cache")
	}
	if cachedItem.TokenParent != "autoauthtoken" {
		t.Fatalf("unexpected value for tokenparent: %s", cachedItem.TokenParent)
	}

	// Modify the request a little bit to ensure the second response is
	// returned to the lease cache.
	sendReq = &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input_changed"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, responses[1].Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Make the same request again and ensure that the same response is returned
	// again.
	sendReq = &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input_changed"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, responses[1].Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}
}

// TestLeaseCache_StoreCacheableStaticSecret tests that cacheStaticSecret works
// as expected, creating the two expected cache entries, and also ensures
// that we can evict the cache entry with the cache clear API afterwards.
func TestLeaseCache_StoreCacheableStaticSecret(t *testing.T) {
	request := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/secrets/foo/bar",
			},
		},
		Token: "token",
	}
	response := newTestSendResponse(http.StatusCreated, `{"data": {"foo": "bar"}, "mount_type": "kvv2"}`)
	responses := []*SendResponse{
		response,
	}
	index := &cachememdb.Index{
		Type:        cacheboltdb.StaticSecretType,
		RequestPath: request.Request.URL.Path,
		Namespace:   "root/",
		Token:       "token",
		ID:          computeStaticSecretCacheIndex(request),
	}

	lc := testNewLeaseCache(t, responses)

	// We expect two entries to be stored by this:
	// 1. The actual static secret
	// 2. The capabilities index
	err := lc.cacheStaticSecret(context.Background(), request, response, index)
	if err != nil {
		return
	}

	indexFromDB, err := lc.db.Get(cachememdb.IndexNameID, index.ID)
	if err != nil {
		return
	}

	require.NotNil(t, indexFromDB)
	require.Equal(t, "token", indexFromDB.Token)
	require.Equal(t, map[string]struct{}{"token": {}}, indexFromDB.Tokens)
	require.Equal(t, cacheboltdb.StaticSecretType, indexFromDB.Type)
	require.Equal(t, request.Request.URL.Path, indexFromDB.RequestPath)
	require.Equal(t, "root/", indexFromDB.Namespace)

	capabilitiesIndexFromDB, err := lc.db.GetCapabilitiesIndex(cachememdb.IndexNameID, hex.EncodeToString(cryptoutil.Blake2b256Hash(index.Token)))
	if err != nil {
		return
	}

	require.NotNil(t, capabilitiesIndexFromDB)
	require.Equal(t, "token", capabilitiesIndexFromDB.Token)
	require.Equal(t, map[string]struct{}{"secrets/foo/bar": {}}, capabilitiesIndexFromDB.ReadablePaths)

	err = lc.handleCacheClear(context.Background(), &cacheClearInput{
		Type:        "request_path",
		RequestPath: request.Request.URL.Path,
	})
	require.NoError(t, err)

	expectedClearedIndex, err := lc.db.Get(cachememdb.IndexNameID, index.ID)
	require.Equal(t, cachememdb.ErrCacheItemNotFound, err)
	require.Nil(t, expectedClearedIndex)
}

// TestLeaseCache_StaticSecret_CacheClear_All tests that static secrets are
// stored correctly, as well as removed from the cache by a cache clear with
// "all" specified as the type.
func TestLeaseCache_StaticSecret_CacheClear_All(t *testing.T) {
	request := &SendRequest{
		Request: &http.Request{
			URL: &url.URL{
				Path: "/v1/secrets/foo/bar",
			},
		},
		Token: "token",
	}
	response := newTestSendResponse(http.StatusCreated, `{"data": {"foo": "bar"}, "mount_type": "kvv2"}`)
	responses := []*SendResponse{
		response,
	}
	index := &cachememdb.Index{
		Type:        cacheboltdb.StaticSecretType,
		RequestPath: request.Request.URL.Path,
		Namespace:   "root/",
		Token:       "token",
		ID:          computeStaticSecretCacheIndex(request),
	}

	lc := testNewLeaseCache(t, responses)

	// We expect two entries to be stored by this:
	// 1. The actual static secret
	// 2. The capabilities index
	err := lc.cacheStaticSecret(context.Background(), request, response, index)
	if err != nil {
		return
	}

	indexFromDB, err := lc.db.Get(cachememdb.IndexNameID, index.ID)
	if err != nil {
		return
	}

	require.NotNil(t, indexFromDB)
	require.Equal(t, "token", indexFromDB.Token)
	require.Equal(t, map[string]struct{}{"token": {}}, indexFromDB.Tokens)
	require.Equal(t, cacheboltdb.StaticSecretType, indexFromDB.Type)
	require.Equal(t, request.Request.URL.Path, indexFromDB.RequestPath)
	require.Equal(t, "root/", indexFromDB.Namespace)

	capabilitiesIndexFromDB, err := lc.db.GetCapabilitiesIndex(cachememdb.IndexNameID, hex.EncodeToString(cryptoutil.Blake2b256Hash(index.Token)))
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, capabilitiesIndexFromDB)
	require.Equal(t, "token", capabilitiesIndexFromDB.Token)
	require.Equal(t, map[string]struct{}{"secrets/foo/bar": {}}, capabilitiesIndexFromDB.ReadablePaths)

	err = lc.handleCacheClear(context.Background(), &cacheClearInput{
		Type: "all",
	})
	require.NoError(t, err)

	expectedClearedIndex, err := lc.db.Get(cachememdb.IndexNameID, index.ID)
	require.Equal(t, cachememdb.ErrCacheItemNotFound, err)
	require.Nil(t, expectedClearedIndex)

	expectedClearedCapabilitiesIndex, err := lc.db.GetCapabilitiesIndex(cachememdb.IndexNameID, capabilitiesIndexFromDB.ID)
	require.Equal(t, cachememdb.ErrCacheItemNotFound, err)
	require.Nil(t, expectedClearedCapabilitiesIndex)
}

// TestLeaseCache_SendCacheableStaticSecret tests that the cache has no issue returning
// static secret style responses. It's similar to TestLeaseCache_SendCacheable in that it
// only tests the surface level of the functionality, but there are other tests that
// test the rest.
func TestLeaseCache_SendCacheableStaticSecret(t *testing.T) {
	response := newTestSendResponse(http.StatusCreated, `{"data": {"foo": "bar"}, "mount_type": "kvv2"}`)
	responses := []*SendResponse{
		response,
		response,
		response,
		response,
	}

	lc := testNewLeaseCache(t, responses)

	// Register a token
	require.NoError(t, lc.RegisterAutoAuthToken("autoauthtoken"))

	// Make a request. A response with a new token is returned to the lease
	// cache and that will be cached.
	urlPath := "http://example.com/v1/sample/api"
	sendReq := &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
	}
	resp, err := lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, response.Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Send the same request again to get the cached response
	sendReq = &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, responses[0].Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Modify the request a little to ensure the second response is
	// returned to the lease cache.
	sendReq = &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input_changed"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, response.Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Make the same request again and ensure that the same response is returned
	// again.
	sendReq = &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input_changed"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response.StatusCode, response.Response.StatusCode); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}
}

func TestLeaseCache_SendNonCacheable(t *testing.T) {
	responses := []*SendResponse{
		newTestSendResponse(http.StatusOK, `{"value": "output"}`),
		newTestSendResponse(http.StatusNotFound, `{"value": "invalid"}`),
		newTestSendResponse(http.StatusOK, `<html>Hello</html>`),
		newTestSendResponse(http.StatusTemporaryRedirect, ""),
	}

	lc := testNewLeaseCache(t, responses)

	// Send a request through the lease cache which is not cacheable (there is
	// no lease information or auth information in the response)
	sendReq := &SendRequest{
		Request: httptest.NewRequest("GET", "http://example.com", strings.NewReader(`{"value": "input"}`)),
	}
	resp, err := lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response, responses[0].Response); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Since the response is non-cacheable, the second response will be
	// returned.
	sendReq = &SendRequest{
		Token:   "foo",
		Request: httptest.NewRequest("GET", "http://example.com", strings.NewReader(`{"value": "input"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response, responses[1].Response); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Since the response is non-cacheable, the third response will be
	// returned.
	sendReq = &SendRequest{
		Token:   "foo",
		Request: httptest.NewRequest("GET", "http://example.com", nil),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response, responses[2].Response); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	// Since the response is non-cacheable, the fourth response will be
	// returned.
	sendReq = &SendRequest{
		Token:   "foo",
		Request: httptest.NewRequest("GET", "http://example.com", nil),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response, responses[3].Response); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}
}

func TestLeaseCache_SendNonCacheableNonTokenLease(t *testing.T) {
	// Create the cache
	responses := []*SendResponse{
		newTestSendResponse(http.StatusOK, `{"value": "output", "lease_id": "foo"}`),
		newTestSendResponse(http.StatusCreated, `{"value": "invalid", "auth": {"client_token": "testtoken"}}`),
	}
	lc := testNewLeaseCache(t, responses)

	// Send a request through lease cache which returns a response containing
	// lease_id. Response will not be cached because it doesn't belong to a
	// token that is managed by the lease cache.
	urlPath := "http://example.com/v1/sample/api"
	sendReq := &SendRequest{
		Token:   "foo",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
	}
	resp, err := lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response, responses[0].Response); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	_, err = lc.db.Get(cachememdb.IndexNameRequestPath, "root/", urlPath)
	if err != cachememdb.ErrCacheItemNotFound {
		t.Fatal("expected entry to be nil, got", err)
	}

	// Verify that the response is not cached by sending the same request and
	// by expecting a different response.
	sendReq = &SendRequest{
		Token:   "foo",
		Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
	}
	resp, err = lc.Send(context.Background(), sendReq)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(resp.Response, responses[1].Response); diff != nil {
		t.Fatalf("expected getting proxied response: got %v", diff)
	}

	_, err = lc.db.Get(cachememdb.IndexNameRequestPath, "root/", urlPath)
	if err != cachememdb.ErrCacheItemNotFound {
		t.Fatal("expected entry to be nil, got", err)
	}
}

func TestLeaseCache_HandleCacheClear(t *testing.T) {
	lc := testNewLeaseCache(t, nil)

	handler := lc.HandleCacheClear(context.Background())
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Test missing body, should return 400
	resp, err := http.Post(ts.URL, "application/json", nil)
	if err != nil {
		t.Fatal()
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status code mismatch: expected = %v, got = %v", http.StatusBadRequest, resp.StatusCode)
	}

	testCases := []struct {
		name               string
		reqType            string
		reqValue           string
		expectedStatusCode int
	}{
		{
			"invalid_type",
			"foo",
			"",
			http.StatusBadRequest,
		},
		{
			"invalid_value",
			"",
			"bar",
			http.StatusBadRequest,
		},
		{
			"all",
			"all",
			"",
			http.StatusOK,
		},
		{
			"by_request_path",
			"request_path",
			"foo",
			http.StatusOK,
		},
		{
			"by_token",
			"token",
			"foo",
			http.StatusOK,
		},
		{
			"by_lease",
			"lease",
			"foo",
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := fmt.Sprintf("{\"type\": \"%s\", \"value\": \"%s\"}", tc.reqType, tc.reqValue)
			resp, err := http.Post(ts.URL, "application/json", strings.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			if tc.expectedStatusCode != resp.StatusCode {
				t.Fatalf("status code mismatch: expected = %v, got = %v", tc.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestCache_DeriveNamespaceAndRevocationPath(t *testing.T) {
	tests := []struct {
		name             string
		req              *SendRequest
		wantNamespace    string
		wantRelativePath string
	}{
		{
			"non_revocation_full_path",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/ns1/sys/mounts",
					},
				},
			},
			"root/",
			"/v1/ns1/sys/mounts",
		},
		{
			"non_revocation_relative_path",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/sys/mounts",
					},
					Header: http.Header{
						consts.NamespaceHeaderName: []string{"ns1/"},
					},
				},
			},
			"ns1/",
			"/v1/sys/mounts",
		},
		{
			"non_revocation_relative_path",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/ns2/sys/mounts",
					},
					Header: http.Header{
						consts.NamespaceHeaderName: []string{"ns1/"},
					},
				},
			},
			"ns1/",
			"/v1/ns2/sys/mounts",
		},
		{
			"revocation_full_path",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/ns1/sys/leases/revoke",
					},
				},
			},
			"ns1/",
			"/v1/sys/leases/revoke",
		},
		{
			"revocation_relative_path",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/sys/leases/revoke",
					},
					Header: http.Header{
						consts.NamespaceHeaderName: []string{"ns1/"},
					},
				},
			},
			"ns1/",
			"/v1/sys/leases/revoke",
		},
		{
			"revocation_relative_partial_ns",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/ns2/sys/leases/revoke",
					},
					Header: http.Header{
						consts.NamespaceHeaderName: []string{"ns1/"},
					},
				},
			},
			"ns1/ns2/",
			"/v1/sys/leases/revoke",
		},
		{
			"revocation_prefix_full_path",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/ns1/sys/leases/revoke-prefix/foo",
					},
				},
			},
			"ns1/",
			"/v1/sys/leases/revoke-prefix/foo",
		},
		{
			"revocation_prefix_relative_path",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/sys/leases/revoke-prefix/foo",
					},
					Header: http.Header{
						consts.NamespaceHeaderName: []string{"ns1/"},
					},
				},
			},
			"ns1/",
			"/v1/sys/leases/revoke-prefix/foo",
		},
		{
			"revocation_prefix_partial_ns",
			&SendRequest{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/v1/ns2/sys/leases/revoke-prefix/foo",
					},
					Header: http.Header{
						consts.NamespaceHeaderName: []string{"ns1/"},
					},
				},
			},
			"ns1/ns2/",
			"/v1/sys/leases/revoke-prefix/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNamespace, gotRelativePath := deriveNamespaceAndRevocationPath(tt.req)
			if gotNamespace != tt.wantNamespace {
				t.Errorf("deriveNamespaceAndRevocationPath() gotNamespace = %v, want %v", gotNamespace, tt.wantNamespace)
			}
			if gotRelativePath != tt.wantRelativePath {
				t.Errorf("deriveNamespaceAndRevocationPath() gotRelativePath = %v, want %v", gotRelativePath, tt.wantRelativePath)
			}
		})
	}
}

func TestLeaseCache_Concurrent_NonCacheable(t *testing.T) {
	lc := testNewLeaseCacheWithDelay(t, false, 50)

	// We are going to send 100 requests, each taking 50ms to process. If these
	// requests are processed serially, it will take ~5seconds to finish. we
	// use a ContextWithTimeout to tell us if this is the case by giving ample
	// time for it process them concurrently but time out if they get processed
	// serially.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	wgDoneCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
		var wg sync.WaitGroup
		// 100 concurrent requests
		for i := 0; i < 100; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				// Send a request through the lease cache which is not cacheable (there is
				// no lease information or auth information in the response)
				sendReq := &SendRequest{
					Request: httptest.NewRequest("GET", "http://example.com", nil),
				}

				_, err := lc.Send(ctx, sendReq)
				if err != nil {
					errCh <- err
				}
			}()
		}

		wg.Wait()
		close(wgDoneCh)
	}()

	select {
	case <-ctx.Done():
		t.Fatalf("request timed out: %s", ctx.Err())
	case <-wgDoneCh:
	case err := <-errCh:
		t.Fatal(err)
	}
}

func TestLeaseCache_Concurrent_Cacheable(t *testing.T) {
	lc := testNewLeaseCacheWithDelay(t, true, 50)

	if err := lc.RegisterAutoAuthToken("autoauthtoken"); err != nil {
		t.Fatal(err)
	}

	// We are going to send 100 requests, each taking 50ms to process. If these
	// requests are processed serially, it will take ~5seconds to finish, so we
	// use a ContextWithTimeout to tell us if this is the case.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var cacheCount atomic.Uint32
	wgDoneCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
		var wg sync.WaitGroup
		// Start 100 concurrent requests
		for i := 0; i < 100; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				sendReq := &SendRequest{
					Token:   "autoauthtoken",
					Request: httptest.NewRequest("GET", "http://example.com/v1/sample/api", nil),
				}

				resp, err := lc.Send(ctx, sendReq)
				if err != nil {
					errCh <- err
				}

				if resp.CacheMeta != nil && resp.CacheMeta.Hit {
					cacheCount.Inc()
				}
			}()
		}

		wg.Wait()
		close(wgDoneCh)
	}()

	select {
	case <-ctx.Done():
		t.Fatalf("request timed out: %s", ctx.Err())
	case <-wgDoneCh:
	case err := <-errCh:
		t.Fatal(err)
	}

	// Ensure that all but one request got proxied. The other 99 should be
	// returned from the cache.
	if cacheCount.Load() != 99 {
		t.Fatalf("Should have returned a cached response 99 times, got %d", cacheCount.Load())
	}
}

func setupBoltStorage(t *testing.T) (tempCacheDir string, boltStorage *cacheboltdb.BoltStorage) {
	t.Helper()

	km, err := keymanager.NewPassthroughKeyManager(context.Background(), nil)
	require.NoError(t, err)

	tempCacheDir, err = ioutil.TempDir("", "agent-cache-test")
	require.NoError(t, err)
	boltStorage, err = cacheboltdb.NewBoltStorage(&cacheboltdb.BoltStorageConfig{
		Path:    tempCacheDir,
		Logger:  hclog.Default(),
		Wrapper: km.Wrapper(),
	})
	require.NoError(t, err)
	require.NotNil(t, boltStorage)
	// The calling function should `defer boltStorage.Close()` and `defer os.RemoveAll(tempCacheDir)`
	return tempCacheDir, boltStorage
}

func compareBeforeAndAfter(t *testing.T, before, after *LeaseCache, beforeLen, afterLen int) {
	beforeDB, err := before.db.GetByPrefix(cachememdb.IndexNameID)
	require.NoError(t, err)
	assert.Len(t, beforeDB, beforeLen)
	afterDB, err := after.db.GetByPrefix(cachememdb.IndexNameID)
	require.NoError(t, err)
	assert.Len(t, afterDB, afterLen)
	for _, cachedItem := range beforeDB {
		if strings.Contains(cachedItem.RequestPath, "expect-missing") {
			continue
		}
		restoredItem, err := after.db.Get(cachememdb.IndexNameID, cachedItem.ID)
		require.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, cachedItem.ID, restoredItem.ID)
		assert.Equal(t, cachedItem.Lease, restoredItem.Lease)
		assert.Equal(t, cachedItem.LeaseToken, restoredItem.LeaseToken)
		assert.Equal(t, cachedItem.Namespace, restoredItem.Namespace)
		assert.EqualValues(t, cachedItem.RequestHeader, restoredItem.RequestHeader)
		assert.Equal(t, cachedItem.RequestMethod, restoredItem.RequestMethod)
		assert.Equal(t, cachedItem.RequestPath, restoredItem.RequestPath)
		assert.Equal(t, cachedItem.RequestToken, restoredItem.RequestToken)
		assert.Equal(t, cachedItem.Response, restoredItem.Response)
		assert.Equal(t, cachedItem.Token, restoredItem.Token)
		assert.Equal(t, cachedItem.TokenAccessor, restoredItem.TokenAccessor)
		assert.Equal(t, cachedItem.TokenParent, restoredItem.TokenParent)

		// check what we can in the renewal context
		assert.NotEmpty(t, restoredItem.RenewCtxInfo.CancelFunc)
		assert.NotZero(t, restoredItem.RenewCtxInfo.DoneCh)
		require.NotEmpty(t, restoredItem.RenewCtxInfo.Ctx)
		assert.Equal(t,
			cachedItem.RenewCtxInfo.Ctx.Value(contextIndexID),
			restoredItem.RenewCtxInfo.Ctx.Value(contextIndexID),
		)
	}
}

func TestLeaseCache_PersistAndRestore(t *testing.T) {
	// Emulate responses from the api proxy. The first two use the auto-auth
	// token, and the others use another token.
	// The test re-sends each request to ensure that the response is cached
	// so the number of responses and cacheTests specified should always be equal.
	responses := []*SendResponse{
		newTestSendResponse(200, `{"auth": {"client_token": "testtoken", "renewable": true, "lease_duration": 600}}`),
		newTestSendResponse(201, `{"lease_id": "foo", "renewable": true, "data": {"value": "foo"}, "lease_duration": 600}`),
		// The auth token will get manually deleted from the bolt DB storage, causing both of the following two responses
		// to be missing from the cache after a restore, because the lease is a child of the auth token.
		newTestSendResponse(202, `{"auth": {"client_token": "testtoken2", "renewable": true, "orphan": true, "lease_duration": 600}}`),
		newTestSendResponse(203, `{"lease_id": "secret2-lease", "renewable": true, "data": {"number": "two"}, "lease_duration": 600}`),
		// 204 No content gets special handling - avoid.
		newTestSendResponse(250, `{"auth": {"client_token": "testtoken3", "renewable": true, "orphan": true, "lease_duration": 600}}`),
		newTestSendResponse(251, `{"lease_id": "secret3-lease", "renewable": true, "data": {"number": "three"}, "lease_duration": 600}`),
		newTestSendResponse(http.StatusCreated, `{"data": {"foo": "bar"}, "mount_type": "kvv2"}`),
	}

	tempDir, boltStorage := setupBoltStorage(t)
	defer os.RemoveAll(tempDir)
	defer boltStorage.Close()
	lc := testNewLeaseCacheWithPersistence(t, responses, boltStorage)

	// Register an auto-auth token so that the token and lease requests are cached
	err := lc.RegisterAutoAuthToken("autoauthtoken")
	require.NoError(t, err)

	cacheTests := []struct {
		token                     string
		method                    string
		urlPath                   string
		body                      string
		deleteFromPersistentStore bool // If true, will be deleted from bolt DB to induce an error on restore
		expectMissingAfterRestore bool // If true, the response is not expected to be present in the restored cache
	}{
		{
			// Make a request. A response with a new token is returned to the
			// lease cache and that will be cached.
			token:   "autoauthtoken",
			method:  "GET",
			urlPath: "http://example.com/v1/sample/api",
			body:    `{"value": "input"}`,
		},
		{
			// Modify the request a little bit to ensure the second response is
			// returned to the lease cache.
			token:   "autoauthtoken",
			method:  "GET",
			urlPath: "http://example.com/v1/sample/api",
			body:    `{"value": "input_changed"}`,
		},
		{
			// Simulate an approle login to get another token
			method:                    "PUT",
			urlPath:                   "http://example.com/v1/auth/approle-expect-missing/login",
			body:                      `{"role_id": "my role", "secret_id": "my secret"}`,
			deleteFromPersistentStore: true,
			expectMissingAfterRestore: true,
		},
		{
			// Test caching with the token acquired from the approle login
			token:   "testtoken2",
			method:  "GET",
			urlPath: "http://example.com/v1/sample-expect-missing/api",
			body:    `{"second": "input"}`,
			// This will be missing from the restored cache because its parent token was deleted
			expectMissingAfterRestore: true,
		},
		{
			// Simulate another approle login to get another token
			method:  "PUT",
			urlPath: "http://example.com/v1/auth/approle/login",
			body:    `{"role_id": "my role", "secret_id": "my secret"}`,
		},
		{
			// Test caching with the token acquired from the latest approle login
			token:   "testtoken3",
			method:  "GET",
			urlPath: "http://example.com/v1/sample3/api",
			body:    `{"third": "input"}`,
		},
	}

	var deleteIDs []string
	for i, ct := range cacheTests {
		// Send once to cache
		req := httptest.NewRequest(ct.method, ct.urlPath, strings.NewReader(ct.body))
		req.Header.Set("User-Agent", useragent.AgentProxyString())

		sendReq := &SendRequest{
			Token:   ct.token,
			Request: req,
		}
		if ct.deleteFromPersistentStore {
			deleteID, err := computeIndexID(sendReq)
			require.NoError(t, err)
			deleteIDs = append(deleteIDs, deleteID)
			// Now reset the body after calculating the index
			req = httptest.NewRequest(ct.method, ct.urlPath, strings.NewReader(ct.body))
			req.Header.Set("User-Agent", useragent.AgentProxyString())
			sendReq.Request = req
		}
		resp, err := lc.Send(context.Background(), sendReq)
		require.NoError(t, err)
		assert.Equal(t, responses[i].Response.StatusCode, resp.Response.StatusCode, "expected proxied response")
		assert.Nil(t, resp.CacheMeta)

		// Send again to test cache. If this isn't cached, the response returned
		// will be the next in the list and the status code will not match.
		req = httptest.NewRequest(ct.method, ct.urlPath, strings.NewReader(ct.body))
		req.Header.Set("User-Agent", useragent.AgentProxyString())
		sendCacheReq := &SendRequest{
			Token:   ct.token,
			Request: req,
		}
		respCached, err := lc.Send(context.Background(), sendCacheReq)
		require.NoError(t, err, "failed to send request %+v", ct)
		assert.Equal(t, responses[i].Response.StatusCode, respCached.Response.StatusCode, "expected proxied response")
		require.NotNil(t, respCached.CacheMeta)
		assert.True(t, respCached.CacheMeta.Hit)
	}

	require.NotEmpty(t, deleteIDs)
	for _, deleteID := range deleteIDs {
		err = boltStorage.Delete(deleteID, cacheboltdb.LeaseType)
		require.NoError(t, err)
	}

	// Now we know the cache is working, so try restoring from the persisted
	// cache's storage. Responses 3 and 4 have been cleared from the cache, so
	// re-send those.
	restoredCache := testNewLeaseCache(t, responses[2:4])

	err = restoredCache.Restore(context.Background(), boltStorage)
	errors, ok := err.(*multierror.Error)
	require.True(t, ok)
	assert.Len(t, errors.Errors, 1)
	assert.Contains(t, errors.Error(), "could not find parent Token testtoken2")

	// Now compare the cache contents before and after
	compareBeforeAndAfter(t, lc, restoredCache, 7, 5)

	// And finally send the cache requests once to make sure they're all being
	// served from the restoredCache unless they were intended to be missing after restore.
	for i, ct := range cacheTests {
		req := httptest.NewRequest(ct.method, ct.urlPath, strings.NewReader(ct.body))
		req.Header.Set("User-Agent", useragent.AgentProxyString())
		sendCacheReq := &SendRequest{
			Token:   ct.token,
			Request: req,
		}
		respCached, err := restoredCache.Send(context.Background(), sendCacheReq)
		require.NoError(t, err, "failed to send request %+v", ct)
		assert.Equal(t, responses[i].Response.StatusCode, respCached.Response.StatusCode, "expected proxied response")
		if ct.expectMissingAfterRestore {
			require.Nil(t, respCached.CacheMeta)
		} else {
			require.NotNil(t, respCached.CacheMeta)
			assert.True(t, respCached.CacheMeta.Hit)
		}
	}
}

func TestLeaseCache_PersistAndRestore_WithManyDependencies(t *testing.T) {
	tempDir, boltStorage := setupBoltStorage(t)
	defer os.RemoveAll(tempDir)
	defer boltStorage.Close()

	var requests []*SendRequest
	var responses []*SendResponse
	var orderedRequestPaths []string

	// helper func to generate new auth leases with a child secret lease attached
	authAndSecretLease := func(id int, parentToken, newToken string) {
		t.Helper()
		path := fmt.Sprintf("/v1/auth/approle-%d/login", id)
		orderedRequestPaths = append(orderedRequestPaths, path)
		requests = append(requests, &SendRequest{
			Token:   parentToken,
			Request: httptest.NewRequest("PUT", "http://example.com"+path, strings.NewReader("")),
		})
		responses = append(responses, newTestSendResponse(200, fmt.Sprintf(`{"auth": {"client_token": "%s", "renewable": true, "lease_duration": 600}}`, newToken)))

		// Fetch a leased secret using the new token
		path = fmt.Sprintf("/v1/kv/%d", id)
		orderedRequestPaths = append(orderedRequestPaths, path)
		requests = append(requests, &SendRequest{
			Token:   newToken,
			Request: httptest.NewRequest("GET", "http://example.com"+path, strings.NewReader("")),
		})
		responses = append(responses, newTestSendResponse(200, fmt.Sprintf(`{"lease_id": "secret-%d-lease", "renewable": true, "data": {"number": %d}, "lease_duration": 600}`, id, id)))
	}

	// Pathological case: a long chain of child tokens
	authAndSecretLease(0, "autoauthtoken", "many-ancestors-token;0")
	for i := 1; i <= 50; i++ {
		// Create a new generation of child token
		authAndSecretLease(i, fmt.Sprintf("many-ancestors-token;%d", i-1), fmt.Sprintf("many-ancestors-token;%d", i))
	}

	// Lots of sibling tokens with auto auth token as their parent
	for i := 51; i <= 100; i++ {
		authAndSecretLease(i, "autoauthtoken", fmt.Sprintf("many-siblings-token;%d", i))
	}

	// Also create some extra siblings for an auth token further down the chain
	for i := 101; i <= 110; i++ {
		authAndSecretLease(i, "many-ancestors-token;25", fmt.Sprintf("many-siblings-for-ancestor-token;%d", i))
	}

	lc := testNewLeaseCacheWithPersistence(t, responses, boltStorage)

	// Register an auto-auth token so that the token and lease requests are cached
	err := lc.RegisterAutoAuthToken("autoauthtoken")
	require.NoError(t, err)

	for _, req := range requests {
		// Send once to cache
		resp, err := lc.Send(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.Response.StatusCode, "expected success")
		assert.Nil(t, resp.CacheMeta)
	}

	// Ensure leases are retrieved in the correct order
	var processed int

	leases, err := boltStorage.GetByType(context.Background(), cacheboltdb.LeaseType)
	require.NoError(t, err)
	for _, lease := range leases {
		index, err := cachememdb.Deserialize(lease)
		require.NoError(t, err)
		require.Equal(t, orderedRequestPaths[processed], index.RequestPath)
		processed++
	}

	assert.Equal(t, len(orderedRequestPaths), processed)

	restoredCache := testNewLeaseCache(t, nil)
	err = restoredCache.Restore(context.Background(), boltStorage)
	require.NoError(t, err)

	// Now compare the cache contents before and after
	compareBeforeAndAfter(t, lc, restoredCache, 223, 223)
}

func TestEvictPersistent(t *testing.T) {
	ctx := context.Background()

	responses := []*SendResponse{
		newTestSendResponse(201, `{"lease_id": "foo", "renewable": true, "data": {"value": "foo"}}`),
	}

	tempDir, boltStorage := setupBoltStorage(t)
	defer os.RemoveAll(tempDir)
	defer boltStorage.Close()
	lc := testNewLeaseCacheWithPersistence(t, responses, boltStorage)

	require.NoError(t, lc.RegisterAutoAuthToken("autoauthtoken"))

	// populate cache by sending request through
	sendReq := &SendRequest{
		Token:   "autoauthtoken",
		Request: httptest.NewRequest("GET", "http://example.com/v1/sample/api", strings.NewReader(`{"value": "some_input"}`)),
	}
	resp, err := lc.Send(context.Background(), sendReq)
	require.NoError(t, err)
	assert.Equal(t, resp.Response.StatusCode, 201, "expected proxied response")
	assert.Nil(t, resp.CacheMeta)

	// Check bolt for the cached lease
	secrets, err := lc.ps.GetByType(ctx, cacheboltdb.LeaseType)
	require.NoError(t, err)
	assert.Len(t, secrets, 1)

	// Call clear for the request path
	err = lc.handleCacheClear(context.Background(), &cacheClearInput{
		Type:        "request_path",
		RequestPath: "/v1/sample/api",
	})
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	// Check that cached item is gone
	secrets, err = lc.ps.GetByType(ctx, cacheboltdb.LeaseType)
	require.NoError(t, err)
	assert.Len(t, secrets, 0)
}

func TestRegisterAutoAuth_sameToken(t *testing.T) {
	// If the auto-auth token already exists in the cache, it should not be
	// stored again in a new index.
	lc := testNewLeaseCache(t, nil)
	err := lc.RegisterAutoAuthToken("autoauthtoken")
	assert.NoError(t, err)

	oldTokenIndex, err := lc.db.Get(cachememdb.IndexNameToken, "autoauthtoken")
	assert.NoError(t, err)
	oldTokenID := oldTokenIndex.ID

	// register the same token again
	err = lc.RegisterAutoAuthToken("autoauthtoken")
	assert.NoError(t, err)

	// check that there's only one index for autoauthtoken
	entries, err := lc.db.GetByPrefix(cachememdb.IndexNameToken, "autoauthtoken")
	assert.NoError(t, err)
	assert.Len(t, entries, 1)

	newTokenIndex, err := lc.db.Get(cachememdb.IndexNameToken, "autoauthtoken")
	assert.NoError(t, err)

	// compare the ID's since those are randomly generated when an index for a
	// token is added to the cache, so if a new token was added, the id's will
	// not match.
	assert.Equal(t, oldTokenID, newTokenIndex.ID)
}

func Test_hasExpired(t *testing.T) {
	responses := []*SendResponse{
		newTestSendResponse(200, `{"auth": {"client_token": "testtoken", "renewable": true, "lease_duration": 60}}`),
		newTestSendResponse(201, `{"lease_id": "foo", "renewable": true, "data": {"value": "foo"}, "lease_duration": 60}`),
	}
	lc := testNewLeaseCache(t, responses)
	require.NoError(t, lc.RegisterAutoAuthToken("autoauthtoken"))

	cacheTests := []struct {
		token          string
		urlPath        string
		leaseType      string
		wantStatusCode int
	}{
		{
			// auth lease
			token:          "autoauthtoken",
			urlPath:        "/v1/sample/auth",
			leaseType:      cacheboltdb.LeaseType,
			wantStatusCode: responses[0].Response.StatusCode,
		},
		{
			// secret lease
			token:          "autoauthtoken",
			urlPath:        "/v1/sample/secret",
			leaseType:      cacheboltdb.LeaseType,
			wantStatusCode: responses[1].Response.StatusCode,
		},
	}

	for _, ct := range cacheTests {
		// Send once to cache
		urlPath := "http://example.com" + ct.urlPath
		sendReq := &SendRequest{
			Token:   ct.token,
			Request: httptest.NewRequest("GET", urlPath, strings.NewReader(`{"value": "input"}`)),
		}
		resp, err := lc.Send(context.Background(), sendReq)
		require.NoError(t, err)
		assert.Equal(t, resp.Response.StatusCode, ct.wantStatusCode, "expected proxied response")
		assert.Nil(t, resp.CacheMeta)

		// get the Index out of the mem cache
		index, err := lc.db.Get(cachememdb.IndexNameRequestPath, "root/", ct.urlPath)
		require.NoError(t, err)
		assert.Equal(t, ct.leaseType, index.Type)

		// The lease duration is 60 seconds, so time.Now() should be within that
		notExpired, err := lc.hasExpired(time.Now().UTC(), index)
		require.NoError(t, err)
		assert.False(t, notExpired)

		// In 90 seconds the index should be "expired"
		futureTime := time.Now().UTC().Add(time.Second * 90)
		expired, err := lc.hasExpired(futureTime, index)
		require.NoError(t, err)
		assert.True(t, expired)
	}
}

func TestLeaseCache_hasExpired_wrong_type(t *testing.T) {
	index := &cachememdb.Index{
		Type: cacheboltdb.TokenType,
		Response: []byte(`HTTP/0.0 200 OK
Content-Type: application/json
Date: Tue, 02 Mar 2021 17:54:16 GMT

{}`),
	}

	lc := testNewLeaseCache(t, nil)
	expired, err := lc.hasExpired(time.Now().UTC(), index)
	assert.False(t, expired)
	assert.EqualError(t, err, `secret without lease encountered in expiration check`)
}

func TestLeaseCacheRestore_expired(t *testing.T) {
	// Emulate 2 responses from the api proxy, both expired
	responses := []*SendResponse{
		newTestSendResponse(200, `{"auth": {"client_token": "testtoken", "renewable": true, "lease_duration": -600}}`),
		newTestSendResponse(201, `{"lease_id": "foo", "renewable": true, "data": {"value": "foo"}, "lease_duration": -600}`),
	}

	tempDir, boltStorage := setupBoltStorage(t)
	defer os.RemoveAll(tempDir)
	defer boltStorage.Close()
	lc := testNewLeaseCacheWithPersistence(t, responses, boltStorage)

	// Register an auto-auth token so that the token and lease requests are cached in mem
	require.NoError(t, lc.RegisterAutoAuthToken("autoauthtoken"))

	cacheTests := []struct {
		token          string
		method         string
		urlPath        string
		body           string
		wantStatusCode int
	}{
		{
			// Make a request. A response with a new token is returned to the
			// lease cache and that will be cached.
			token:          "autoauthtoken",
			method:         "GET",
			urlPath:        "http://example.com/v1/sample/api",
			body:           `{"value": "input"}`,
			wantStatusCode: responses[0].Response.StatusCode,
		},
		{
			// Modify the request a little bit to ensure the second response is
			// returned to the lease cache.
			token:          "autoauthtoken",
			method:         "GET",
			urlPath:        "http://example.com/v1/sample/api",
			body:           `{"value": "input_changed"}`,
			wantStatusCode: responses[1].Response.StatusCode,
		},
	}

	for _, ct := range cacheTests {
		// Send once to cache
		sendReq := &SendRequest{
			Token:   ct.token,
			Request: httptest.NewRequest(ct.method, ct.urlPath, strings.NewReader(ct.body)),
		}
		resp, err := lc.Send(context.Background(), sendReq)
		require.NoError(t, err)
		assert.Equal(t, resp.Response.StatusCode, ct.wantStatusCode, "expected proxied response")
		assert.Nil(t, resp.CacheMeta)
	}

	// Restore from the persisted cache's storage
	restoredCache := testNewLeaseCache(t, nil)

	err := restoredCache.Restore(context.Background(), boltStorage)
	assert.NoError(t, err)

	// The original mem cache should between one-to-three items.
	// This will usually be three, but could be less if any renewals
	// happens before this check, which will evict the expired cache entries.
	// e.g. you add a time.Sleep before this, it will be 1. We check
	// between the range to reduce flakiness.
	beforeDB, err := lc.db.GetByPrefix(cachememdb.IndexNameID)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(beforeDB), 3)
	assert.LessOrEqual(t, 1, len(beforeDB))

	// There should only be one item in the restored cache: the autoauth token
	afterDB, err := restoredCache.db.GetByPrefix(cachememdb.IndexNameID)
	require.NoError(t, err)
	assert.Len(t, afterDB, 1)

	// Just verify that the one item in the restored mem cache matches one in the original mem cache, and that it's the auto-auth token
	beforeItem, err := lc.db.Get(cachememdb.IndexNameID, afterDB[0].ID)
	require.NoError(t, err)
	assert.NotNil(t, beforeItem)

	assert.Equal(t, "autoauthtoken", afterDB[0].Token)
	assert.Equal(t, cacheboltdb.TokenType, afterDB[0].Type)
}
