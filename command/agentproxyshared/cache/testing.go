// Copyright (c) HashiCorp, Inc.
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

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/useragent"
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

func newTestSendResponse(status int, body string) *SendResponse {
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

type mockTokenVerifierProxier struct {
	currentToken string
}

func (p *mockTokenVerifierProxier) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	p.currentToken = req.Token
	resp := newTestSendResponse(http.StatusOK,
		`{"data": {"id": "`+p.currentToken+`"}}`)

	return resp, nil
}

func (p *mockTokenVerifierProxier) GetCurrentRequestToken() string {
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
		return newTestSendResponse(http.StatusOK, s), nil
	}

	return newTestSendResponse(http.StatusOK, `{"value": "output"}`), nil
}
