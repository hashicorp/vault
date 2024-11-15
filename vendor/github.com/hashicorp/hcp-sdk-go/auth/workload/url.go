// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// URLCredentialSource sources credentials by making an HTTP request to the
// given URL.
type URLCredentialSource struct {
	// URL reads the credentials by invoking the given URL with the headers.
	URL string `json:"url,omitempty"`

	// Headers are included when invoking the given URL.
	Headers map[string]string `json:"headers,omitempty"`

	// CredentialFormat configures how the credentials are extracted from the HTTP
	// response body.
	CredentialFormat
}

// Validate validates the config.
func (uc *URLCredentialSource) Validate() error {
	if uc.URL == "" {
		return fmt.Errorf("non-empty URL is required")
	}

	_, err := url.Parse(uc.URL)
	if err != nil {
		return fmt.Errorf("failed to parse url %q: %v", uc.URL, err)
	}

	return uc.CredentialFormat.Validate()
}

// token retrieves the token by making a HTTP request to the configured URL.
func (uc *URLCredentialSource) token() (string, error) {
	// Build the request
	req, err := http.NewRequest("GET", uc.URL, nil)
	if err != nil {
		return "", fmt.Errorf("failed creating an HTTP request for workload access_token: %v", err)
	}

	// Make the request with a timeout
	reqCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)

	// Add the headers to the request
	for key, val := range uc.Headers {
		req.Header.Add(key, val)
	}

	// Make the request
	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("invalid response retrieving token: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("failed reading body in subject token response: %v", err)
	}
	if c := resp.StatusCode; c < 200 || c > 299 {
		return "", fmt.Errorf("subject token response failed with status code %d: %s", c, respBody)
	}

	if len(respBody) == 0 {
		return "", fmt.Errorf("response body is empty")
	}

	// Extract the value
	return uc.CredentialFormat.get(respBody)
}
