package cache

import (
	"context"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/cache/cacheboltdb"
	"github.com/hashicorp/vault/command/agent/cache/cachememdb"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func testNewLeaseCache(t *testing.T, responses []*SendResponse) *LeaseCache {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	lc, err := NewLeaseCache(&LeaseCacheConfig{
		Client:      client,
		BaseContext: context.Background(),
		Proxier:     newMockProxier(responses),
		Logger:      logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
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
		Client:      client,
		BaseContext: context.Background(),
		Proxier:     &mockDelayProxier{cacheable, delay},
		Logger:      logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
	})
	if err != nil {
		t.Fatal(err)
	}

	return lc
}

func testNewLeaseCacheWithPersistence(t *testing.T, responses []*SendResponse, storage cacheboltdb.Storage) *LeaseCache {
	t.Helper()

	client, err := api.NewClient(api.DefaultConfig())
	require.NoError(t, err)

	lc, err := NewLeaseCache(&LeaseCacheConfig{
		Client:      client,
		BaseContext: context.Background(),
		Proxier:     newMockProxier(responses),
		Logger:      logging.NewVaultLogger(hclog.Trace).Named("cache.leasecache"),
		Storage:     storage,
	})
	require.NoError(t, err)

	return lc
}

func TestCache_ComputeIndexID(t *testing.T) {
	type args struct {
		req *http.Request
	}
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
	// Register an token so that the token and lease requests are cached
	lc.RegisterAutoAuthToken("autoauthtoken")

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

	idx, err := lc.db.Get(cachememdb.IndexNameRequestPath, "root/", urlPath)
	if err != nil {
		t.Fatal(err)
	}
	if idx != nil {
		t.Fatalf("expected nil entry, got: %#v", idx)
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

	idx, err = lc.db.Get(cachememdb.IndexNameRequestPath, "root/", urlPath)
	if err != nil {
		t.Fatal(err)
	}
	if idx != nil {
		t.Fatalf("expected nil entry, got: %#v", idx)
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
					t.Fatal(err)
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
					t.Fatal(err)
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
	}

	// Ensure that all but one request got proxied. The other 99 should be
	// returned from the cache.
	if cacheCount.Load() != 99 {
		t.Fatalf("Should have returned a cached response 99 times, got %d", cacheCount.Load())
	}
}

func setupBoltStorage(t *testing.T) (tempCacheDir string, boltStorage *cacheboltdb.BoltStorage) {
	t.Helper()

	e, err := cacheboltdb.NewAES(&cacheboltdb.AESConfig{
		Key:    []byte("thisisafakekey!!thisisafakekey!!"),
		AAD:    []byte("extra-data"),
		Logger: hclog.NewNullLogger(),
	})
	require.NoError(t, err)

	tempCacheDir, err = ioutil.TempDir("", "agent-cache-test")
	require.NoError(t, err)
	boltStorage, err = cacheboltdb.NewBoltStorage(&cacheboltdb.BoltStorageConfig{
		Path:       tempCacheDir,
		RootBucket: "topbucketname",
		Logger:     hclog.Default(),
		Encrypter:  e,
	})
	require.NoError(t, err)
	require.NotNil(t, boltStorage)
	// The calling function should `defer boltStorage.Close()` and `defer os.RemoveAll(tempCacheDir)`
	return tempCacheDir, boltStorage
}

func TestLeaseCache_PersistAndRestore(t *testing.T) {
	// Emulate 4 responses from the api proxy. The first two use the auto-auth
	// token, and the last two use another token.
	responses := []*SendResponse{
		newTestSendResponse(200, `{"auth": {"client_token": "testtoken", "renewable": true}}`),
		newTestSendResponse(201, `{"lease_id": "foo", "renewable": true, "data": {"value": "foo"}}`),
		newTestSendResponse(202, `{"auth": {"client_token": "testtoken2", "renewable": true, "orphan": true}}`),
		newTestSendResponse(203, `{"lease_id": "secret2-lease", "renewable": true, "data": {"number": "two"}}`),
	}

	tempDir, boltStorage := setupBoltStorage(t)
	defer os.RemoveAll(tempDir)
	defer boltStorage.Close()
	lc := testNewLeaseCacheWithPersistence(t, responses, boltStorage)

	// Register an auto-auth token so that the token and lease requests are cached
	lc.RegisterAutoAuthToken("autoauthtoken")

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
		{
			// Simulate an approle login to get another token
			method:         "PUT",
			urlPath:        "http://example.com/v1/auth/approle/login",
			body:           `{"role_id": "my role", "secret_id": "my secret"}`,
			wantStatusCode: responses[2].Response.StatusCode,
		},
		{
			// Test caching with the token acquired from the approle login
			token:          "testtoken2",
			method:         "GET",
			urlPath:        "http://example.com/v1/sample2/api",
			body:           `{"second": "input"}`,
			wantStatusCode: responses[3].Response.StatusCode,
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

		// Send again to test cache. If this isn't cached, the response returned
		// will be the next in the list and the status code will not match.
		sendCacheReq := &SendRequest{
			Token:   ct.token,
			Request: httptest.NewRequest(ct.method, ct.urlPath, strings.NewReader(ct.body)),
		}
		respCached, err := lc.Send(context.Background(), sendCacheReq)
		require.NoError(t, err, "failed to send request %+v", ct)
		assert.Equal(t, respCached.Response.StatusCode, ct.wantStatusCode, "expected proxied response")
		require.NotNil(t, respCached.CacheMeta)
		assert.True(t, respCached.CacheMeta.Hit)
	}

	// Now we know the cache is working, so try restoring from the persisted
	// cache's storage
	restoredCache := testNewLeaseCache(t, nil)

	err := restoredCache.Restore(boltStorage)
	assert.NoError(t, err)

	// Now compare before and after
	beforeDB, err := lc.db.GetByPrefix(cachememdb.IndexNameID)
	require.NoError(t, err)
	assert.Len(t, beforeDB, 5)

	for _, cachedItem := range beforeDB {
		restoredItem, err := restoredCache.db.Get(cachememdb.IndexNameID, cachedItem.ID)
		require.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, cachedItem.ID, restoredItem.ID)
		assert.Equal(t, cachedItem.Lease, restoredItem.Lease)
		assert.Equal(t, cachedItem.LeaseToken, restoredItem.LeaseToken)
		assert.Equal(t, cachedItem.Namespace, restoredItem.Namespace)
		assert.Equal(t, cachedItem.RequestHeader, restoredItem.RequestHeader)
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
	afterDB, err := restoredCache.db.GetByPrefix(cachememdb.IndexNameID)
	require.NoError(t, err)
	assert.Len(t, afterDB, 5)

	// And finally send the cache requests once to make sure they're all being
	// served from the restoredCache
	for _, ct := range cacheTests {
		sendCacheReq := &SendRequest{
			Token:   ct.token,
			Request: httptest.NewRequest(ct.method, ct.urlPath, strings.NewReader(ct.body)),
		}
		respCached, err := restoredCache.Send(context.Background(), sendCacheReq)
		require.NoError(t, err, "failed to send request %+v", ct)
		assert.Equal(t, respCached.Response.StatusCode, ct.wantStatusCode, "expected proxied response")
		require.NotNil(t, respCached.CacheMeta)
		assert.True(t, respCached.CacheMeta.Hit)
	}
}

func TestEvictPersistent(t *testing.T) {
	responses := []*SendResponse{
		newTestSendResponse(201, `{"lease_id": "foo", "renewable": true, "data": {"value": "foo"}}`),
	}

	tempDir, boltStorage := setupBoltStorage(t)
	defer os.RemoveAll(tempDir)
	defer boltStorage.Close()
	lc := testNewLeaseCacheWithPersistence(t, responses, boltStorage)

	lc.RegisterAutoAuthToken("autoauthtoken")

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
	secrets, err := lc.ps.GetByType(cacheboltdb.SecretLeaseType)
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
	secrets, err = lc.ps.GetByType(cacheboltdb.SecretLeaseType)
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
