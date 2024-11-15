// Copyright 2014 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"

	"cloud.google.com/go/internal"
	"cloud.google.com/go/internal/version"
	sinternal "cloud.google.com/go/storage/internal"
	"github.com/google/uuid"
	gax "github.com/googleapis/gax-go/v2"
	"github.com/googleapis/gax-go/v2/callctx"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var defaultRetry *retryConfig = &retryConfig{}
var xGoogDefaultHeader = fmt.Sprintf("gl-go/%s gccl/%s", version.Go(), sinternal.Version)

const (
	xGoogHeaderKey       = "x-goog-api-client"
	idempotencyHeaderKey = "x-goog-gcs-idempotency-token"
)

// run determines whether a retry is necessary based on the config and
// idempotency information. It then calls the function with or without retries
// as appropriate, using the configured settings.
func run(ctx context.Context, call func(ctx context.Context) error, retry *retryConfig, isIdempotent bool) error {
	attempts := 1
	invocationID := uuid.New().String()

	if retry == nil {
		retry = defaultRetry
	}
	if (retry.policy == RetryIdempotent && !isIdempotent) || retry.policy == RetryNever {
		ctxWithHeaders := setInvocationHeaders(ctx, invocationID, attempts)
		return call(ctxWithHeaders)
	}
	bo := gax.Backoff{}
	if retry.backoff != nil {
		bo.Multiplier = retry.backoff.Multiplier
		bo.Initial = retry.backoff.Initial
		bo.Max = retry.backoff.Max
	}
	var errorFunc func(err error) bool = ShouldRetry
	if retry.shouldRetry != nil {
		errorFunc = retry.shouldRetry
	}

	return internal.Retry(ctx, bo, func() (stop bool, err error) {
		ctxWithHeaders := setInvocationHeaders(ctx, invocationID, attempts)
		err = call(ctxWithHeaders)
		if err != nil && retry.maxAttempts != nil && attempts >= *retry.maxAttempts {
			return true, fmt.Errorf("storage: retry failed after %v attempts; last error: %w", *retry.maxAttempts, err)
		}
		attempts++
		retryable := errorFunc(err)
		// Explicitly check context cancellation so that we can distinguish between a
		// DEADLINE_EXCEEDED error from the server and a user-set context deadline.
		// Unfortunately gRPC will codes.DeadlineExceeded (which may be retryable if it's
		// sent by the server) in both cases.
		if ctxErr := ctx.Err(); errors.Is(ctxErr, context.Canceled) || errors.Is(ctxErr, context.DeadlineExceeded) {
			retryable = false
		}
		return !retryable, err
	})
}

// Sets invocation ID headers on the context which will be propagated as
// headers in the call to the service (for both gRPC and HTTP).
func setInvocationHeaders(ctx context.Context, invocationID string, attempts int) context.Context {
	invocationHeader := fmt.Sprintf("gccl-invocation-id/%v gccl-attempt-count/%v", invocationID, attempts)
	xGoogHeader := strings.Join([]string{invocationHeader, xGoogDefaultHeader}, " ")

	ctx = callctx.SetHeaders(ctx, xGoogHeaderKey, xGoogHeader)
	ctx = callctx.SetHeaders(ctx, idempotencyHeaderKey, invocationID)
	return ctx
}

// ShouldRetry returns true if an error is retryable, based on best practice
// guidance from GCS. See
// https://cloud.google.com/storage/docs/retry-strategy#go for more information
// on what errors are considered retryable.
//
// If you would like to customize retryable errors, use the WithErrorFunc to
// supply a RetryOption to your library calls. For example, to retry additional
// errors, you can write a custom func that wraps ShouldRetry and also specifies
// additional errors that should return true.
func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	if errors.Is(err, net.ErrClosed) {
		return true
	}

	switch e := err.(type) {
	case *googleapi.Error:
		// Retry on 408, 429, and 5xx, according to
		// https://cloud.google.com/storage/docs/exponential-backoff.
		return e.Code == 408 || e.Code == 429 || (e.Code >= 500 && e.Code < 600)
	case *net.OpError, *url.Error:
		// Retry socket-level errors ECONNREFUSED and ECONNRESET (from syscall).
		// Unfortunately the error type is unexported, so we resort to string
		// matching.
		retriable := []string{"connection refused", "connection reset", "broken pipe"}
		for _, s := range retriable {
			if strings.Contains(e.Error(), s) {
				return true
			}
		}
	case *net.DNSError:
		if e.IsTemporary {
			return true
		}
	case interface{ Temporary() bool }:
		if e.Temporary() {
			return true
		}
	}
	// UNAVAILABLE, RESOURCE_EXHAUSTED, INTERNAL, and DEADLINE_EXCEEDED codes are all retryable for gRPC.
	if st, ok := status.FromError(err); ok {
		if code := st.Code(); code == codes.Unavailable || code == codes.ResourceExhausted || code == codes.Internal || code == codes.DeadlineExceeded {
			return true
		}
	}
	// Unwrap is only supported in go1.13.x+
	if e, ok := err.(interface{ Unwrap() error }); ok {
		return ShouldRetry(e.Unwrap())
	}
	return false
}
