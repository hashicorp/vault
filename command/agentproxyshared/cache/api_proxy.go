// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"fmt"
	gohttp "net/http"
	"sync"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/http"
)

type EnforceConsistency int

const (
	EnforceConsistencyNever EnforceConsistency = iota
	EnforceConsistencyAlways
)

type WhenInconsistentAction int

const (
	WhenInconsistentFail WhenInconsistentAction = iota
	WhenInconsistentRetry
	WhenInconsistentForward
)

// APIProxy is an implementation of the proxier interface that is used to
// forward the request to Vault and get the response.
type APIProxy struct {
	client                  *api.Client
	logger                  hclog.Logger
	enforceConsistency      EnforceConsistency
	whenInconsistentAction  WhenInconsistentAction
	l                       sync.RWMutex
	lastIndexStates         []string
	userAgentString         string
	userAgentStringFunction func(string) string
	// clientNamespace is a one-time set representation of the namespace of the client
	// (i.e. client.Namespace()) to avoid repeated calls and lock usage.
	clientNamespace            string
	prependConfiguredNamespace bool
}

var _ Proxier = &APIProxy{}

type APIProxyConfig struct {
	Client                 *api.Client
	Logger                 hclog.Logger
	EnforceConsistency     EnforceConsistency
	WhenInconsistentAction WhenInconsistentAction
	// UserAgentString is used as the User Agent when the proxied client
	// does not have a user agent of its own.
	UserAgentString string
	// UserAgentStringFunction is the function to transform the proxied client's
	// user agent into one that includes Vault-specific information.
	UserAgentStringFunction func(string) string
	// PrependConfiguredNamespace configures whether the client's namespace
	// should be prepended to proxied requests
	PrependConfiguredNamespace bool
}

func NewAPIProxy(config *APIProxyConfig) (Proxier, error) {
	if config.Client == nil {
		return nil, fmt.Errorf("nil API client")
	}
	return &APIProxy{
		client:                     config.Client,
		logger:                     config.Logger,
		enforceConsistency:         config.EnforceConsistency,
		whenInconsistentAction:     config.WhenInconsistentAction,
		userAgentString:            config.UserAgentString,
		userAgentStringFunction:    config.UserAgentStringFunction,
		prependConfiguredNamespace: config.PrependConfiguredNamespace,
		clientNamespace:            namespace.Canonicalize(config.Client.Namespace()),
	}, nil
}

func (ap *APIProxy) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	client, err := ap.client.Clone()
	if err != nil {
		return nil, err
	}
	client.SetToken(req.Token)

	// Derive and set a logger for the client
	clientLogger := ap.logger.Named("client")
	client.SetLogger(clientLogger)

	// http.Transport will transparently request gzip and decompress the response, but only if
	// the client doesn't manually set the header. Removing any Accept-Encoding header allows the
	// transparent compression to occur.
	req.Request.Header.Del("Accept-Encoding")

	if req.Request.Header == nil {
		req.Request.Header = make(gohttp.Header)
	}

	// Set our User-Agent to be one indicating we are Vault Agent's API proxy.
	// If the sending client had one, preserve it.
	if req.Request.Header.Get("User-Agent") != "" {
		initialUserAgent := req.Request.Header.Get("User-Agent")
		req.Request.Header.Set("User-Agent", ap.userAgentStringFunction(initialUserAgent))
	} else {
		req.Request.Header.Set("User-Agent", ap.userAgentString)
	}

	client.SetHeaders(req.Request.Header)
	if ap.prependConfiguredNamespace && ap.clientNamespace != "" {
		currentNamespace := namespace.Canonicalize(client.Namespace())
		newNamespace := namespace.Canonicalize(ap.clientNamespace + currentNamespace)
		client.SetNamespace(newNamespace)
	}

	fwReq := client.NewRequest(req.Request.Method, req.Request.URL.Path)
	fwReq.BodyBytes = req.RequestBody

	query := req.Request.URL.Query()
	if len(query) != 0 {
		fwReq.Params = query
	}

	var newState string
	manageState := ap.enforceConsistency == EnforceConsistencyAlways &&
		req.Request.Header.Get(http.VaultIndexHeaderName) == "" &&
		req.Request.Header.Get(http.VaultForwardHeaderName) == "" &&
		req.Request.Header.Get(http.VaultInconsistentHeaderName) == ""

	if manageState {
		client = client.WithResponseCallbacks(api.RecordState(&newState))
		ap.l.RLock()
		lastStates := ap.lastIndexStates
		ap.l.RUnlock()
		if len(lastStates) != 0 {
			client = client.WithRequestCallbacks(api.RequireState(lastStates...))
			switch ap.whenInconsistentAction {
			case WhenInconsistentFail:
				// In this mode we want to delegate handling of inconsistency
				// failures to the external client talking to Agent.
				client.SetCheckRetry(retryablehttp.DefaultRetryPolicy)
			case WhenInconsistentRetry:
				// In this mode we want to handle retries due to inconsistency
				// internally.  This is the default api.Client behaviour so
				// we needn't do anything.
			case WhenInconsistentForward:
				fwReq.Headers.Set(http.VaultInconsistentHeaderName, http.VaultInconsistentForward)
			}
		}
	}

	// Make the request to Vault and get the response
	ap.logger.Info("forwarding request to Vault", "method", req.Request.Method, "path", req.Request.URL.Path)

	resp, err := client.RawRequestWithContext(ctx, fwReq)
	if resp == nil && err != nil {
		// We don't want to cache nil responses, so we simply return the error
		return nil, err
	}

	if newState != "" {
		ap.l.Lock()
		// We want to be using the "newest" states seen, but newer isn't well
		// defined here.  There can be two states S1 and S2 which aren't strictly ordered:
		// S1 could have a newer localindex and S2 could have a newer replicatedindex.  So
		// we need to merge them.  But we can't merge them because we wouldn't be able to
		// "sign" the resulting header because we don't have access to the HMAC key that
		// Vault uses to do so.  So instead we compare any of the 0-2 saved states
		// we have to the new header, keeping the newest 1-2 of these, and sending
		// them to Vault to evaluate.
		ap.lastIndexStates = api.MergeReplicationStates(ap.lastIndexStates, newState)
		ap.l.Unlock()
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
