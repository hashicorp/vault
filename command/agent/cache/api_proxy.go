package cache

import (
	"context"
	"fmt"
	"sync"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/vault"
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
	client                 *api.Client
	logger                 hclog.Logger
	enforceConsistency     EnforceConsistency
	whenInconsistentAction WhenInconsistentAction
	l                      sync.RWMutex
	lastIndexStates        []string
}

var _ Proxier = &APIProxy{}

type APIProxyConfig struct {
	Client                 *api.Client
	Logger                 hclog.Logger
	EnforceConsistency     EnforceConsistency
	WhenInconsistentAction WhenInconsistentAction
}

func NewAPIProxy(config *APIProxyConfig) (Proxier, error) {
	if config.Client == nil {
		return nil, fmt.Errorf("nil API client")
	}
	return &APIProxy{
		client:                 config.Client,
		logger:                 config.Logger,
		enforceConsistency:     config.EnforceConsistency,
		whenInconsistentAction: config.WhenInconsistentAction,
	}, nil
}

// compareStates returns 1 if s1 is newer or identical, -1 if s1 is older, and 0
// if neither s1 or s2 is strictly greater.  An error is returned if s1 or s2
// are invalid or from different clusters.
func compareStates(s1, s2 string) (int, error) {
	w1, err := vault.ParseRequiredState(s1, nil)
	if err != nil {
		return 0, err
	}
	w2, err := vault.ParseRequiredState(s2, nil)
	if err != nil {
		return 0, err
	}

	if w1.ClusterID != w2.ClusterID {
		return 0, fmt.Errorf("don't know how to compare states with different ClusterIDs")
	}

	switch {
	case w1.LocalIndex >= w2.LocalIndex && w1.ReplicatedIndex >= w2.ReplicatedIndex:
		return 1, nil
	// We've already handled the case where both are equal above, so really we're
	// asking here if one or both are lesser.
	case w1.LocalIndex <= w2.LocalIndex && w1.ReplicatedIndex <= w2.ReplicatedIndex:
		return -1, nil
	}

	return 0, nil
}

func mergeStates(old []string, new string) []string {
	if len(old) == 0 || len(old) > 2 {
		return []string{new}
	}

	var ret []string
	for _, o := range old {
		c, err := compareStates(o, new)
		if err != nil {
			return []string{new}
		}
		switch c {
		case 1:
			ret = append(ret, o)
		case -1:
			ret = append(ret, new)
		case 0:
			ret = append(ret, o, new)
		}
	}
	return strutil.RemoveDuplicates(ret, false)
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
	client.SetHeaders(req.Request.Header)

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
	ap.logger.Info("forwarding request", "method", req.Request.Method, "path", req.Request.URL.Path)

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
		ap.lastIndexStates = mergeStates(ap.lastIndexStates, newState)
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
