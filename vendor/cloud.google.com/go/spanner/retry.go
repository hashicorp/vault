/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"context"
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/internal/trace"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	retryInfoKey = "google.rpc.retryinfo-bin"
)

// DefaultRetryBackoff is used for retryers as a fallback value when the server
// did not return any retry information.
var DefaultRetryBackoff = gax.Backoff{
	Initial:    20 * time.Millisecond,
	Max:        32 * time.Second,
	Multiplier: 1.3,
}

// spannerRetryer extends the generic gax Retryer, but also checks for any
// retry info returned by Cloud Spanner and uses that if present.
type spannerRetryer struct {
	gax.Retryer
}

// onCodes returns a spannerRetryer that will retry on the specified error
// codes. For Internal errors, only errors that have one of a list of known
// descriptions should be retried.
func onCodes(bo gax.Backoff, cc ...codes.Code) gax.Retryer {
	return &spannerRetryer{
		Retryer: gax.OnCodes(cc, bo),
	}
}

// Retry returns the retry delay returned by Cloud Spanner if that is present.
// Otherwise it returns the retry delay calculated by the generic gax Retryer.
func (r *spannerRetryer) Retry(err error) (time.Duration, bool) {
	if status.Code(err) == codes.Internal &&
		!strings.Contains(err.Error(), "stream terminated by RST_STREAM") &&
		// See b/25451313.
		!strings.Contains(err.Error(), "HTTP/2 error code: INTERNAL_ERROR") &&
		// See b/27794742.
		!strings.Contains(err.Error(), "Connection closed with unknown cause") &&
		!strings.Contains(err.Error(), "Received unexpected EOS on DATA frame from server") {
		return 0, false
	}

	delay, shouldRetry := r.Retryer.Retry(err)
	if !shouldRetry {
		return 0, false
	}
	if serverDelay, hasServerDelay := ExtractRetryDelay(err); hasServerDelay {
		delay = serverDelay
	}
	return delay, true
}

// runWithRetryOnAbortedOrFailedInlineBeginOrSessionNotFound executes the given function and
// retries it if it returns an Aborted, Session not found error or certain Internal errors. The retry
// is delayed if the error was Aborted or Internal error. The delay between retries is the delay
// returned by Cloud Spanner, or if none is returned, the calculated delay with
// a minimum of 10ms and maximum of 32s. There is no delay before the retry if
// the error was Session not found or failed inline begin transaction.
func runWithRetryOnAbortedOrFailedInlineBeginOrSessionNotFound(ctx context.Context, f func(context.Context) error) error {
	retryer := onCodes(DefaultRetryBackoff, codes.Aborted, codes.ResourceExhausted, codes.Internal)
	funcWithRetry := func(ctx context.Context) error {
		for {
			err := f(ctx)
			if err == nil {
				return nil
			}
			// Get Spanner or GRPC status error.
			// TODO(loite): Refactor to unwrap Status error instead of Spanner
			// error when statusError implements the (errors|xerrors).Wrapper
			// interface.
			var retryErr error
			var se *Error
			if errors.As(err, &se) {
				// It is a (wrapped) Spanner error. Use that to check whether
				// we should retry.
				retryErr = se
			} else {
				// It's not a Spanner error, check if it is a status error.
				_, ok := status.FromError(err)
				if !ok {
					return err
				}
				retryErr = err
			}
			if isSessionNotFoundError(retryErr) {
				trace.TracePrintf(ctx, nil, "Retrying after Session not found")
				continue
			}
			if isFailedInlineBeginTransaction(retryErr) {
				trace.TracePrintf(ctx, nil, "Retrying after failed inline begin transaction")
				continue
			}
			delay, shouldRetry := retryer.Retry(retryErr)
			if !shouldRetry {
				return err
			}
			trace.TracePrintf(ctx, nil, "Backing off after ABORTED for %s, then retrying", delay)
			if err := gax.Sleep(ctx, delay); err != nil {
				return err
			}
		}
	}
	return funcWithRetry(ctx)
}

// ExtractRetryDelay extracts retry backoff from a *spanner.Error if present.
func ExtractRetryDelay(err error) (time.Duration, bool) {
	var se *Error
	var s *status.Status
	// Unwrap status error.
	if errors.As(err, &se) {
		s = status.Convert(se.Unwrap())
	} else {
		s = status.Convert(err)
	}
	if s == nil {
		return 0, false
	}
	for _, detail := range s.Details() {
		if retryInfo, ok := detail.(*errdetails.RetryInfo); ok {
			if !retryInfo.GetRetryDelay().IsValid() {
				return 0, false
			}
			return retryInfo.GetRetryDelay().AsDuration(), true
		}
	}
	return 0, false
}
