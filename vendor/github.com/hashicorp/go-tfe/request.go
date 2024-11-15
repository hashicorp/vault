// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
)

// ClientRequest encapsulates a request sent by the Client
type ClientRequest struct {
	retryableRequest *retryablehttp.Request
	http             *retryablehttp.Client
	limiter          *rate.Limiter

	// Header are the headers that will be sent in this request
	Header http.Header
}

func (r ClientRequest) Do(ctx context.Context, model interface{}) error {
	// Wait will block until the limiter can obtain a new token
	// or returns an error if the given context is canceled.
	if r.limiter != nil {
		if err := r.limiter.Wait(ctx); err != nil {
			return err
		}
	}

	// If the caller provided a response header hook then we'll call it
	// once we have a response.
	respHeaderHook := contextResponseHeaderHook(ctx)

	// Add the context to the request.
	reqWithCxt := r.retryableRequest.WithContext(ctx)

	// Execute the request and check the response.
	resp, err := r.http.Do(reqWithCxt)
	if resp != nil {
		// We call the callback whenever there's any sort of response,
		// even if it's returned in conjunction with an error.
		respHeaderHook(resp.StatusCode, resp.Header)
	}
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	// Basic response checking.
	if err := checkResponseCode(resp); err != nil {
		return err
	}

	// Return here if decoding the response isn't needed.
	if model == nil {
		return nil
	}

	// If v implements io.Writer, write the raw response body.
	if w, ok := model.(io.Writer); ok {
		_, err := io.Copy(w, resp.Body)
		return err
	}

	return unmarshalResponse(resp.Body, model)
}

// DoJSON is similar to Do except that it should be used when a plain JSON response is expected
// as opposed to json-api.
func (r *ClientRequest) DoJSON(ctx context.Context, model any) error {
	// Wait will block until the limiter can obtain a new token
	// or returns an error if the given context is canceled.
	if r.limiter != nil {
		if err := r.limiter.Wait(ctx); err != nil {
			return err
		}
	}

	// Add the context to the request.
	contextReq := r.retryableRequest.WithContext(ctx)

	// If the caller provided a response header hook then we'll call it
	// once we have a response.
	respHeaderHook := contextResponseHeaderHook(ctx)

	// Execute the request and check the response.
	resp, err := r.http.Do(contextReq)
	if resp != nil {
		// We call the callback whenever there's any sort of response,
		// even if it's returned in conjunction with an error.
		respHeaderHook(resp.StatusCode, resp.Header)
	}
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("error HTTP response: %d", resp.StatusCode)
	} else if resp.StatusCode == 304 {
		// Got a "Not Modified" response, but we can't return a model because there is no response body.
		// This is necessary to support the IPRanges endpoint, which has the peculiar behavior
		// of not returning content but allowing a 304 response by optionally sending an
		// If-Modified-Since header.
		return nil
	}

	// Return here if decoding the response isn't needed.
	if model == nil {
		return nil
	}

	// If v implements io.Writer, write the raw response body.
	if w, ok := model.(io.Writer); ok {
		_, err := io.Copy(w, resp.Body)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(model)
}
