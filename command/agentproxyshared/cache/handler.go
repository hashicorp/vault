// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

func ProxyHandler(ctx context.Context, logger hclog.Logger, proxier Proxier, inmemSink sink.Sink, forceAutoAuthToken bool, useAutoAuthToken bool, authInProgress *atomic.Bool, invalidTokenErrCh chan error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received request", "method", r.Method, "path", r.URL.Path)

		if forceAutoAuthToken {
			r.Header.Del(consts.AuthHeaderName)
		}

		token := r.Header.Get(consts.AuthHeaderName)

		var autoAuthToken string
		if inmemSink != nil {
			autoAuthToken = inmemSink.(sink.SinkReader).Token()
			if token == "" && useAutoAuthToken {
				logger.Debug("using auto auth token", "method", r.Method, "path", r.URL.Path)
				token = autoAuthToken
			}
		}

		// Parse and reset body.
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body")
			logical.RespondError(w, http.StatusInternalServerError, errors.New("failed to read request body"))
			return
		}
		if r.Body != nil {
			r.Body.Close()
		}
		r.Body = io.NopCloser(bytes.NewReader(reqBody))
		req := &SendRequest{
			Token:       token,
			Request:     r,
			RequestBody: reqBody,
		}

		resp, err := proxier.Send(ctx, req)
		if err != nil {
			// If this is an api.Response error, don't wrap the response.
			if resp != nil && resp.Response.Error() != nil {
				responseErrMessage := resp.Response.Error()
				copyHeader(w.Header(), resp.Response.Header)
				w.WriteHeader(resp.Response.StatusCode)
				io.Copy(w, resp.Response.Body)
				metrics.IncrCounter([]string{"agent", "proxy", "client_error"}, 1)
				// Re-trigger auto auth if the token is the same as the auto auth token
				if resp.Response.StatusCode == 403 && strings.Contains(responseErrMessage.Error(), logical.ErrInvalidToken.Error()) &&
					autoAuthToken == token && !authInProgress.Load() {
					// Drain the error channel first
					logger.Info("proxy received an invalid token error")
					select {
					case <-invalidTokenErrCh:
					default:
					}
					invalidTokenErrCh <- resp.Response.Error()
				}
			} else {
				metrics.IncrCounter([]string{"agent", "proxy", "error"}, 1)
				logical.RespondError(w, http.StatusInternalServerError, fmt.Errorf("failed to get the response: %w", err))
			}
			return
		}

		err = sanitizeAutoAuthTokenResponse(ctx, logger, inmemSink, req, resp)
		if err != nil {
			logical.RespondError(w, http.StatusInternalServerError, fmt.Errorf("failed to process token lookup response: %w", err))
			return
		}

		defer resp.Response.Body.Close()

		metrics.IncrCounter([]string{"agent", "proxy", "success"}, 1)
		if resp.CacheMeta != nil {
			if resp.CacheMeta.Hit {
				metrics.IncrCounter([]string{"agent", "cache", "hit"}, 1)
			} else {
				metrics.IncrCounter([]string{"agent", "cache", "miss"}, 1)
			}
		}

		// Set headers
		setHeaders(w, resp)

		// Set response body
		io.Copy(w, resp.Response.Body)
		return
	})
}

// setHeaders is a helper that sets the header values based on SendResponse. It
// copies over the headers from the original response and also includes any
// cache-related headers.
func setHeaders(w http.ResponseWriter, resp *SendResponse) {
	// Set header values
	copyHeader(w.Header(), resp.Response.Header)
	if resp.CacheMeta != nil {
		xCacheVal := "MISS"

		if resp.CacheMeta.Hit {
			xCacheVal = "HIT"

			// If this is a cache hit, we also set the Age header
			age := fmt.Sprintf("%.0f", resp.CacheMeta.Age.Seconds())
			w.Header().Set("Age", age)

			// Update the date value
			w.Header().Set("Date", time.Now().Format(http.TimeFormat))
		}

		w.Header().Set("X-Cache", xCacheVal)
	}

	// Set status code
	w.WriteHeader(resp.Response.StatusCode)
}

// sanitizeAutoAuthTokenResponse checks if the request was a lookup or renew
// and if the auto-auth token was used to perform lookup-self, the identifier
// of the token and its accessor same will be stripped off of the response.
func sanitizeAutoAuthTokenResponse(ctx context.Context, logger hclog.Logger, inmemSink sink.Sink, req *SendRequest, resp *SendResponse) error {
	// If auto-auth token is not being used, there is nothing to do.
	if inmemSink == nil {
		return nil
	}
	autoAuthToken := inmemSink.(sink.SinkReader).Token()

	// If lookup responded with non 200 status, there is nothing to do.
	if resp.Response.StatusCode != http.StatusOK {
		return nil
	}

	_, path := deriveNamespaceAndRevocationPath(req)
	switch path {
	case vaultPathTokenLookupSelf, vaultPathTokenRenewSelf:
		if req.Token != autoAuthToken {
			return nil
		}
	case vaultPathTokenLookup, vaultPathTokenRenew:
		jsonBody := map[string]interface{}{}
		if err := json.Unmarshal(req.RequestBody, &jsonBody); err != nil {
			return err
		}
		tokenRaw, ok := jsonBody["token"]
		if !ok {
			// Input error will be caught by the API
			return nil
		}
		token, ok := tokenRaw.(string)
		if !ok {
			// Input error will be caught by the API
			return nil
		}
		if token != "" && token != autoAuthToken {
			// Lookup is performed on the non-auto-auth token
			return nil
		}
	default:
		return nil
	}

	logger.Info("stripping auto-auth token from the response", "method", req.Request.Method, "path", req.Request.URL.Path)
	secret, err := api.ParseSecret(bytes.NewReader(resp.ResponseBody))
	if err != nil {
		return fmt.Errorf("failed to parse token lookup response: %v", err)
	}
	if secret == nil {
		return nil
	} else if secret.Data != nil {
		// lookup endpoints
		if secret.Data["id"] == nil && secret.Data["accessor"] == nil {
			return nil
		}
		delete(secret.Data, "id")
		delete(secret.Data, "accessor")
	} else if secret.Auth != nil {
		// renew endpoints
		if secret.Auth.Accessor == "" && secret.Auth.ClientToken == "" {
			return nil
		}
		secret.Auth.Accessor = ""
		secret.Auth.ClientToken = ""
	} else {
		// nothing to redact
		return nil
	}

	bodyBytes, err := json.Marshal(secret)
	if err != nil {
		return err
	}
	if resp.Response.Body != nil {
		resp.Response.Body.Close()
	}
	resp.Response.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	resp.Response.ContentLength = int64(len(bodyBytes))

	// Serialize and re-read the response
	var respBytes bytes.Buffer
	err = resp.Response.Write(&respBytes)
	if err != nil {
		return fmt.Errorf("failed to serialize the updated response: %v", err)
	}

	updatedResponse, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(respBytes.Bytes())), nil)
	if err != nil {
		return fmt.Errorf("failed to deserialize the updated response: %v", err)
	}

	resp.Response = &api.Response{
		Response: updatedResponse,
	}
	resp.ResponseBody = bodyBytes

	return nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
