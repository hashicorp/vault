package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

// mockProxier is a mock implementation of the Proxier interface, used for testing purposes.
// The mock will return the provided responses every time it reaches its Send method, up to
// the last provided response. This lets tests control what the next/underlying Proxier layer
// might expect to return.
type mockProxier struct {
	proxiedResponses []*SendResponse
	responseIndex    int
}

func newMockProxier(responses []*SendResponse) *mockProxier {
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
	resp := &SendResponse{
		Response: &api.Response{
			Response: &http.Response{
				StatusCode: status,
				Header:     http.Header{},
			},
		},
	}
	resp.Response.Header.Set("Date", time.Now().Format(http.TimeFormat))

	if body != "" {
		resp.Response.Body = ioutil.NopCloser(strings.NewReader(body))
		resp.ResponseBody = []byte(body)
	}

	if json.Valid([]byte(body)) {
		resp.Response.Header.Set("content-type", "application/json")
	}

	return resp
}
