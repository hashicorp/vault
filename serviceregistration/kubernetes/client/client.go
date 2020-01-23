package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
)

// maxRetries is the maximum number of times the client
// should retry.
const maxRetries = 10

var (
	ErrNamespaceUnset = errors.New(`"namespace" is unset`)
	ErrPodNameUnset   = errors.New(`"podName" is unset`)
	ErrNotFound       = errors.New("not found")
	ErrNotInCluster   = errors.New("unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined")
)

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

type PatchOperation int

const (
	// When adding support for more PatchOperations in the future,
	// DO NOT alphebetize them because it will change the underlying
	// int representing a user's intent. If that's stored anywhere,
	// it will cause storage reads to map to the incorrect operation.
	Unset PatchOperation = iota
	Add
	Replace
)

func Parse(s string) PatchOperation {
	switch s {
	case "add":
		return Add
	case "replace":
		return Replace
	default:
		return Unset
	}
}

func (p PatchOperation) String() string {
	switch p {
	case Unset:
		// This is an invalid choice, and will be shown on a patch
		// where the PatchOperation is unset. That's because ints
		// default to 0, and Unset corresponds to 0.
		return "unset"
	case Add:
		return "add"
	case Replace:
		return "replace"
	default:
		// Should never arrive here.
		return ""
	}
}

type Patch struct {
	Operation PatchOperation
	Path      string
	Value     interface{}
}

func New(logger hclog.Logger) (*Client, error) {
	config, err := inClusterConfig()
	if err != nil {
		return nil, err
	}
	return &Client{
		logger: logger,
		config: config,
	}, nil
}

type Client struct {
	logger hclog.Logger
	config *Config
}

// GetPod merely verifies a pod's existence, returning an
// error if the pod doesn't exist.
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

// PatchPod updates the pod's tags to the given ones,
// overwriting previous values for a given tag key. It does so
// non-destructively, or in other words, without tearing down
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

	var jsonPatches []interface{}
	for _, patch := range patches {
		if patch.Operation == Unset {
			return errors.New("patch operation must be set")
		}
		jsonPatches = append(jsonPatches, map[string]interface{}{
			"op":    patch.Operation.String(),
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

func (c *Client) do(req *http.Request, ptrToReturnObj interface{}) error {
	// Finish setting up a valid request.
	req.Header.Set("Authorization", "Bearer "+c.config.BearerToken)
	req.Header.Set("Accept", "application/json")
	client := cleanhttp.DefaultClient()
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: c.config.CACertPool,
		},
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		shouldRetry, err := c.attemptRequest(client, req, ptrToReturnObj)
		if !shouldRetry {
			// The error may be nil or populated depending on whether the
			// request was successful.
			return err
		}
		lastErr = err
	}
	return lastErr
}

func (c *Client) attemptRequest(client *http.Client, req *http.Request, ptrToReturnObj interface{}) (shouldRetry bool, err error) {
	// Preserve the original request body so it can be viewed for debugging if needed.
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = ioutil.ReadAll(req.Body)
		reqBodyReader := bytes.NewReader(reqBody)
		req.Body = ioutil.NopCloser(reqBodyReader)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
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

	// Check for success.
	switch resp.StatusCode {
	case 200, 201, 202:
		// Pass.
	case 401, 403:
		// Perhaps the token from our bearer token file has been refreshed.
		config, err := inClusterConfig()
		if err != nil {
			return false, err
		}
		if config.BearerToken == c.config.BearerToken {
			// It's the same token.
			return false, fmt.Errorf("bad status code: %s", sanitizedDebuggingInfo(req, reqBody, resp))
		}
		c.config = config
		// Continue to try again, but return the error too in case the caller would rather read it out.
		return true, fmt.Errorf("bad status code: %s", sanitizedDebuggingInfo(req, reqBody, resp))
	case 404:
		return false, ErrNotFound
	default:
		return false, fmt.Errorf("unexpected status code: %s", sanitizedDebuggingInfo(req, reqBody, resp))
	}

	// If we're not supposed to read out the body, we have nothing further
	// to do here.
	if ptrToReturnObj == nil {
		return false, nil
	}

	// Attempt to read out the body into the given return object.
	return false, json.NewDecoder(resp.Body).Decode(ptrToReturnObj)
}

// sanitizedDebuggingInfo converts an http response to a string without
// including its headers to avoid leaking authorization
// headers.
func sanitizedDebuggingInfo(req *http.Request, reqBody []byte, resp *http.Response) string {
	respBody, _ := ioutil.ReadAll(resp.Body)
	return fmt.Sprintf("req method: %s, req url: %s, req body: %s, resp statuscode: %d, resp respBody: %s", req.Method, req.URL, reqBody, resp.StatusCode, respBody)
}
