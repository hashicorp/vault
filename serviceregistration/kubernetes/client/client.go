package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	// Retry configuration
	retryWaitMin = 500 * time.Millisecond
	retryWaitMax = 30 * time.Second
	retryMax     = 10
)

var (
	ErrNamespaceUnset = errors.New(`"namespace" is unset`)
	ErrPodNameUnset   = errors.New(`"podName" is unset`)
	ErrNotInCluster   = errors.New("unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined")
)

// New instantiates a Client. The stopCh is used for exiting retry loops
// when closed.
func New(logger hclog.Logger, stopCh <-chan struct{}) (*Client, error) {
	config, err := inClusterConfig()
	if err != nil {
		return nil, err
	}
	return &Client{
		logger: logger,
		config: config,
		stopCh: stopCh,
	}, nil
}

// Client is a minimal Kubernetes client. We rolled our own because the existing
// Kubernetes client-go library available externally has a high number of dependencies
// and we thought it wasn't worth it for only two API calls. If at some point they break
// the client into smaller modules, or if we add quite a few methods to this client, it may
// be worthwhile to revisit that decision.
type Client struct {
	logger hclog.Logger
	config *Config
	stopCh <-chan struct{}
}

// GetPod gets a pod from the Kubernetes API.
func (c *Client) GetPod(namespace, podName string) (*Pod, error) {
	endpoint := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName)
	method := http.MethodGet

	// Validate that we received required parameters.
	if namespace == "" {
		return nil, ErrNamespaceUnset
	}
	if podName == "" {
		return nil, ErrPodNameUnset
	}

	req, err := http.NewRequest(method, c.config.Host+endpoint, nil)
	if err != nil {
		return nil, err
	}
	pod := &Pod{}
	if err := c.do(req, pod); err != nil {
		return nil, err
	}
	return pod, nil
}

// PatchPod updates the pod's tags to the given ones.
// It does so non-destructively, or in other words, without tearing down
// the pod.
func (c *Client) PatchPod(namespace, podName string, patches ...*Patch) error {
	endpoint := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName)
	method := http.MethodPatch

	// Validate that we received required parameters.
	if namespace == "" {
		return ErrNamespaceUnset
	}
	if podName == "" {
		return ErrPodNameUnset
	}
	if len(patches) == 0 {
		// No work to perform.
		return nil
	}

	var jsonPatches []map[string]interface{}
	for _, patch := range patches {
		if patch.Operation == Unset {
			return errors.New("patch operation must be set")
		}
		jsonPatches = append(jsonPatches, map[string]interface{}{
			"op":    patch.Operation,
			"path":  patch.Path,
			"value": patch.Value,
		})
	}
	body, err := json.Marshal(jsonPatches)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.config.Host+endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json-patch+json")
	return c.do(req, nil)
}

// do executes the given request, retrying if necessary.
func (c *Client) do(req *http.Request, ptrToReturnObj interface{}) error {
	// Finish setting up a valid request.
	retryableReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return err
	}
	retryableReq.WithContext(newContext(c.stopCh))
	retryableReq.Header.Set("Authorization", "Bearer "+c.config.BearerToken)
	retryableReq.Header.Set("Accept", "application/json")

	client := &retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultClient(),
		RetryWaitMin: retryWaitMin,
		RetryWaitMax: retryWaitMax,
		RetryMax:     retryMax,
		CheckRetry:   c.getCheckRetry(req),
		Backoff:      retryablehttp.DefaultBackoff,
	}
	client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: c.config.CACertPool,
		},
	}

	// Execute and retry the request. This client comes with exponential backoff and
	// jitter already rolled in.
	resp, err := client.Do(retryableReq)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			if c.logger.IsWarn() {
				// Failing to close response bodies can present as a memory leak so it's
				// important to surface it.
				c.logger.Warn(fmt.Sprintf("unable to close response body: %s", err))
			}
		}
	}()

	// If we're not supposed to read out the body, we have nothing further
	// to do here.
	if ptrToReturnObj == nil {
		return nil
	}

	// Attempt to read out the body into the given return object.
	return json.NewDecoder(resp.Body).Decode(ptrToReturnObj)
}

func (c *Client) getCheckRetry(req *http.Request) retryablehttp.CheckRetry {
	return func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		switch resp.StatusCode {
		case 200, 201, 202, 204:
			// Success.
			return false, nil
		case 401, 403:
			// Perhaps the token from our bearer token file has been refreshed.
			config, err := inClusterConfig()
			if err != nil {
				return false, err
			}
			if config.BearerToken == c.config.BearerToken {
				// It's the same token.
				return false, fmt.Errorf("bad status code: %s", sanitizedDebuggingInfo(req, resp.StatusCode))
			}
			c.config = config
			// Continue to try again, but return the error too in case the caller would rather read it out.
			return true, fmt.Errorf("bad status code: %s", sanitizedDebuggingInfo(req, resp.StatusCode))
		case 404:
			return false, &ErrNotFound{debuggingInfo: sanitizedDebuggingInfo(req, resp.StatusCode)}
		case 500, 502, 503, 504:
			// Could be transient.
			return true, fmt.Errorf("unexpected status code: %s", sanitizedDebuggingInfo(req, resp.StatusCode))
		}
		// Unexpected.
		return false, fmt.Errorf("unexpected status code: %s", sanitizedDebuggingInfo(req, resp.StatusCode))
	}
}

type Pod struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Name string `json:"name,omitempty"`

	// This map will be nil if no "labels" key was provided.
	// It will be populated but have a length of zero if the
	// key was provided, but no values.
	Labels map[string]string `json:"labels,omitempty"`
}

type PatchOperation string

const (
	Unset   PatchOperation = "unset"
	Add                    = "add"
	Replace                = "replace"
)

type Patch struct {
	Operation PatchOperation
	Path      string
	Value     interface{}
}

type ErrNotFound struct {
	debuggingInfo string
}

func (e *ErrNotFound) Error() string {
	return e.debuggingInfo
}

// sanitizedDebuggingInfo provides a returnable string that can be used for debugging. This is intentionally somewhat vague
// because we don't want to leak secrets that may be in a request or response body.
func sanitizedDebuggingInfo(req *http.Request, respStatus int) string {
	return fmt.Sprintf("req method: %s, req url: %s, resp statuscode: %d", req.Method, req.URL, respStatus)
}

// stopChContext allows us to cause requests to stop retrying if we receive
// a stop from the caller.
func newContext(stopCh <-chan struct{}) context.Context {
	ctx := &stopChContext{
		stopCh: stopCh,
	}
	go ctx.monitorStopCh()
	return ctx
}

type stopChContext struct {
	stopCh  <-chan struct{}
	errLock sync.RWMutex
	err     error
}

func (c *stopChContext) Deadline() (deadline time.Time, ok bool) {
	// Return false to indicate that no deadline is set.
	return time.Time{}, false
}

func (c *stopChContext) Done() <-chan struct{} {
	return c.stopCh
}

func (c *stopChContext) Err() error {
	c.errLock.RLock()
	defer c.errLock.RUnlock()
	return c.err
}

func (c *stopChContext) Value(key interface{}) interface{} {
	return nil
}

func (c *stopChContext) monitorStopCh() {
	<-c.stopCh
	c.errLock.Lock()
	defer c.errLock.Unlock()
	c.err = errors.New("stop channel has been closed")
}
