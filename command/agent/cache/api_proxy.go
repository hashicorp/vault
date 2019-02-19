package cache

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
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
	client.SetHeaders(req.Request.Header)

	fwReq := client.NewRequest(req.Request.Method, req.Request.URL.Path)
	fwReq.BodyBytes = req.RequestBody

	// Make the request to Vault and get the response
	ap.logger.Info("forwarding request", "path", req.Request.URL.Path, "method", req.Request.Method)
	resp, err := client.RawRequestWithContext(ctx, fwReq)
	if err != nil {
		return nil, err
	}

	sendResponse := &SendResponse{
		Response: resp,
	}

	// Set SendResponse.ResponseBody if the response body is non-nil
	if resp.Body != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ap.logger.Error("failed to read request body", "error", err)
			return nil, err
		}
		// Close the old body
		resp.Body.Close()

		// Re-set the response body for potential consumption on the way back up the
		// Proxier middleware chain.
		resp.Body = ioutil.NopCloser(bytes.NewReader(respBody))

		sendResponse.ResponseBody = respBody
	}

	return sendResponse, nil
}
