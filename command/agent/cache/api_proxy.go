package cache

import (
	"bytes"
	"context"
	"io/ioutil"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
)

// APIProxy is an implementation of the proxier interface that is used to
// forward the request to Vault and get the response.
type APIProxy struct {
	logger hclog.Logger
}

type APIProxyConfig struct {
	Logger hclog.Logger
}

func NewAPIProxy(config *APIProxyConfig) Proxier {
	return &APIProxy{
		logger: config.Logger,
	}
}

func (ap *APIProxy) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	client.SetToken(req.Token)
	client.SetHeaders(req.Request.Header)

	fwReq := client.NewRequest(req.Request.Method, req.Request.URL.Path)
	fwReq.BodyBytes = req.RequestBody

	// Make the request to Vault and get the response
	ap.logger.Info("forwarding request", "path", req.Request.URL.Path, "method", req.Request.Method)
	resp, err := client.RawRequestWithContext(ctx, fwReq)
	if err != nil {
		return nil, err
	}

	// Parse and reset response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ap.logger.Error("failed to read request body", "error", err)
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

	return &SendResponse{
		Response:     resp,
		ResponseBody: respBody,
	}, nil
}
