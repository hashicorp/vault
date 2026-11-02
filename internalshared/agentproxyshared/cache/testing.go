// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/internalshared/agentproxyshared/cache/cachememdb"
	"github.com/hashicorp/vault/internalshared/agentproxyshared/sink"
	"github.com/hashicorp/vault/internalshared/useragent"
)

// mockProxier is a mock implementation of the Proxier interface, used for testing purposes.
// The mock will return the provided responses every time it reaches its Send method, up to
// the last provided response. This lets tests control what the next/underlying Proxier layer
// might expect to return.
type mockProxier struct {
	proxiedResponses []*SendResponse
	responseIndex    int
}

func NewMockProxier(responses []*SendResponse) *mockProxier {
	return &mockProxier{
		proxiedResponses: responses,
	}
}

func (p *mockProxier) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	if p.responseIndex >= len(p.proxiedResponses) {
		return nil, fmt.Errorf("index out of bounds: responseIndex = %d, responses = %d", p.responseIndex, len(p.proxiedResponses))
	}
	resp := p.proxiedResponses[p.responseIndex]

	p.responseIndex++

	return resp, nil
}

func (p *mockProxier) ResponseIndex() int {
	return p.responseIndex
}

// NewTestSendResponse returns a canned [SendResponse] for tests.
func NewTestSendResponse(status int, body string) *SendResponse {
	headers := make(http.Header)
	headers.Add("User-Agent", useragent.AgentProxyString())
	resp := &SendResponse{
		Response: &api.Response{
			Response: &http.Response{
				StatusCode: status,
				Header:     headers,
			},
		},
	}
	resp.Response.Header.Set("Date", time.Now().Format(http.TimeFormat))

	if body != "" {
		resp.Response.Body = io.NopCloser(strings.NewReader(body))
		resp.ResponseBody = []byte(body)
	}

	if json.Valid([]byte(body)) {
		resp.Response.Header.Set("content-type", "application/json")
	}

	return resp
}

// MockTokenVerifierProxier is a [Proxier] that records the token of the last
// request and echoes it back as the response data id.
type MockTokenVerifierProxier struct {
	currentToken string
}

func (p *MockTokenVerifierProxier) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	p.currentToken = req.Token
	resp := NewTestSendResponse(http.StatusOK,
		`{"data": {"id": "`+p.currentToken+`"}}`)

	return resp, nil
}

func (p *MockTokenVerifierProxier) GetCurrentRequestToken() string {
	return p.currentToken
}

type mockDelayProxier struct {
	cacheableResp bool
	delay         int
}

func (p *mockDelayProxier) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	if p.delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(p.delay) * time.Millisecond):
		}
	}

	// If this is a cacheable response, we return a unique response every time
	if p.cacheableResp {
		rand.Seed(time.Now().Unix())
		s := fmt.Sprintf(`{"lease_id": "%d", "renewable": true, "data": {"foo": "bar"}}`, rand.Int())
		return NewTestSendResponse(http.StatusOK, s), nil
	}

	return NewTestSendResponse(http.StatusOK, `{"value": "output"}`), nil
}

// The ForTest accessors below expose internals to the cluster-backed cache
// tests in command/server, which can't live in this module because they
// depend on the Vault server.

// DBForTest returns the cache's in-memory index database.
func (c *LeaseCache) DBForTest() *cachememdb.CacheMemDB {
	return c.db
}

// BaseContextInfoForTest returns the cache's base context info.
func (c *LeaseCache) BaseContextInfoForTest() *cachememdb.ContextInfo {
	return c.baseCtxInfo
}

// HashStaticSecretIndexForTest exposes hashStaticSecretIndex.
func HashStaticSecretIndexForTest(unhashedIndex string) string {
	return hashStaticSecretIndex(unhashedIndex)
}

// ComputeStaticSecretCacheIndexForTest exposes computeStaticSecretCacheIndex.
func ComputeStaticSecretCacheIndexForTest(req *SendRequest) string {
	return computeStaticSecretCacheIndex(req)
}

// LeaseCacheForTest returns the updater's lease cache.
func (updater *StaticSecretCacheUpdater) LeaseCacheForTest() *LeaseCache {
	return updater.leaseCache
}

// ClientForTest returns the updater's Vault client.
func (updater *StaticSecretCacheUpdater) ClientForTest() *api.Client {
	return updater.client
}

// TokenSinkForTest returns the updater's token sink.
func (updater *StaticSecretCacheUpdater) TokenSinkForTest() sink.Sink {
	return updater.tokenSink
}

// UpdateStaticSecretForTest exposes updateStaticSecret.
func (updater *StaticSecretCacheUpdater) UpdateStaticSecretForTest(ctx context.Context, path string, requiredVaultIndex string) error {
	return updater.updateStaticSecret(ctx, path, requiredVaultIndex)
}

// OpenWebSocketConnectionForTest exposes openWebSocketConnection.
func (updater *StaticSecretCacheUpdater) OpenWebSocketConnectionForTest(ctx context.Context) (*websocket.Conn, error) {
	return updater.openWebSocketConnection(ctx)
}

// StreamStaticSecretEventsForTest exposes streamStaticSecretEvents.
func (updater *StaticSecretCacheUpdater) StreamStaticSecretEventsForTest(ctx context.Context) error {
	return updater.streamStaticSecretEvents(ctx)
}

// PreEventStreamUpdateForTest exposes preEventStreamUpdate.
func (updater *StaticSecretCacheUpdater) PreEventStreamUpdateForTest(ctx context.Context) error {
	return updater.preEventStreamUpdate(ctx)
}

// LeaseCacheForTest returns the capability manager's lease cache.
func (sscm *StaticSecretCapabilityManager) LeaseCacheForTest() *LeaseCache {
	return sscm.leaseCache
}

// GetCapabilitiesForTest exposes getCapabilities.
func GetCapabilitiesForTest(paths []string, client *api.Client) (map[string][]string, error) {
	return getCapabilities(paths, client)
}
