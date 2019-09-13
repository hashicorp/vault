package cache

import (
	"context"
	"fmt"
	"net/http"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

// APIProxy is an implementation of the proxier interface that is used to
// forward the request to Vault and get the response.
type APIProxy struct {
	client *api.Client
	logger hclog.Logger
}

type APIProxyConfig struct {
	Client *api.Client
	Logger hclog.Logger
}

func NewAPIProxy(config *APIProxyConfig) (Proxier, error) {
	if config.Client == nil {
		return nil, fmt.Errorf("nil API client")
	}
	return &APIProxy{
		client: config.Client,
		logger: config.Logger,
	}, nil
}

func (ap *APIProxy) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	client, err := ap.client.Clone()
	if err != nil {
		return nil, err
	}
	client.SetToken(req.Token)

	// http.Transport will transparently request gzip and decompress the response, but only if
	// the client doesn't manually set the header (in which case the client has to handle
	// the decompression. If gzip has already been set, remove it to avoid triggering the manual
	// handling requirement.
	h := clone(req.Request.Header)
	if v, ok := h["Accept-Encoding"]; ok {
		h["Accept-Encoding"] = strutil.StrListDelete(v, "gzip")
	}
	client.SetHeaders(h)

	fwReq := client.NewRequest(req.Request.Method, req.Request.URL.Path)
	fwReq.BodyBytes = req.RequestBody

	query := req.Request.URL.Query()
	if len(query) != 0 {
		fwReq.Params = query
	}

	// Make the request to Vault and get the response
	ap.logger.Info("forwarding request", "method", req.Request.Method, "path", req.Request.URL.Path)

	resp, err := client.RawRequestWithContext(ctx, fwReq)
	if resp == nil && err != nil {
		// We don't want to cache nil responses, so we simply return the error
		return nil, err
	}

	// Before error checking from the request call, we'd want to initialize a SendResponse to
	// potentially return
	sendResponse, newErr := NewSendResponse(resp, nil)
	if newErr != nil {
		return nil, newErr
	}

	// Bubble back the api.Response as well for error checking/handling at the handler layer.
	return sendResponse, err
}

// clone returns a copy of h or nil if h is nil.
//
// TODO: This is a copy of the new (Header) Clone method in Go 1.13. We should remove this once we update.
func clone(h http.Header) http.Header {
	if h == nil {
		return nil
	}

	// Find total number of values.
	nv := 0
	for _, vv := range h {
		nv += len(vv)
	}
	sv := make([]string, nv) // shared backing array for headers' values
	h2 := make(http.Header, len(h))

	for k, vv := range h {
		n := copy(sv, vv)
		h2[k] = sv[:n:n]
		sv = sv[n:]

	}

	return h2
}
