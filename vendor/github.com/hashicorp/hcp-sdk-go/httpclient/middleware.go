// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"fmt"
	"net/http"
	"strings"
)

// MiddlewareOption is a function that modifies an HTTP request.
type MiddlewareOption = func(req *http.Request) error

// roundTripperWithMiddleware takes a plain Roundtripper and an array of MiddlewareOptions to apply to the Roundtripper's request.
type roundTripperWithMiddleware struct {
	OriginalRoundTripper http.RoundTripper
	MiddlewareOptions    []MiddlewareOption
}

// withSourceChannel updates the request header to include the HCP Go SDK source channel stamp.
func withSourceChannel(sourceChannel string) MiddlewareOption {
	return func(req *http.Request) error {
		req.Header.Set("X-HCP-Source-Channel", sourceChannel)
		return nil
	}
}

// withProfile takes the user profile's org ID and project ID and sets them in the request path if needed.
func withOrgAndProjectIDs(orgID, projID string) MiddlewareOption {
	return func(req *http.Request) error {
		path := req.URL.Path
		path = strings.Replace(path, "organizations//", fmt.Sprintf("organizations/%s/", orgID), 1)
		path = strings.Replace(path, "projects//", fmt.Sprintf("projects/%s/", projID), 1)
		req.URL.Path = path
		return nil
	}
}

// RoundTrip attaches MiddlewareOption modifications to the request before sending along.
func (rt *roundTripperWithMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {

	for _, mw := range rt.MiddlewareOptions {
		if err := mw(req); err != nil {
			// Failure to apply middleware should not fail the request
			fmt.Printf("failed to apply middleware: %#v", mw(req))
			continue
		}
	}

	return rt.OriginalRoundTripper.RoundTrip(req)
}
