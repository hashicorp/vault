// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const (
	defaultNumberOfRetries = 5
	defaultDelayTimeBase   = time.Second
	defaultMaxBackoffDelay = 5 * time.Minute
)

type retryOptions struct {
	maxRetries int64 // Optional, defaults to 5
	// maxBackoffDelay  sets a capping value for the delay between calls, to avoid it growing infinitely
	maxBackoffDelay time.Duration // Optional, defaults to 5 min
	// maxToLastCall sets a capping value for all the retry process, in case there is a deadline to make the call.
	maxToLastCall time.Duration // Optional, defaults to 0, meaning no time cap
	// fixedDelay is used in case an uniform distribution of the calls is preferred.
	fixedDelay time.Duration // Optional, defaults to 0, meaning Delay is exponential, starting at 1sec
	// delayBase is used to calculate the starting value at which the delay starts to grow,
	// When left empty, a value of 1 sec will be used as base and then the delays will
	// grow exponentially with every attempt: starting at 1s, then 2s, 4s, 8s...
	delayBase time.Duration // Optional, defaults to 1sec

	// maxValidAttempt is used to ensure that a big attempts number or a big delayBase number will not cause
	// a negative delay by overflowing the delay increase. Every attempt after the
	// maxValid will use the maxBackoffDelay if configured, or the defaultMaxBackoffDelay if not.
	maxValidAttempt int64
}

func (c *Client) retryPut(ctx context.Context, endpoint string, in, out any, q *WriteOptions) (*WriteMeta, error) {
	var err error
	var wm *WriteMeta

	attemptDelay := 100 * time.Second // Avoid a tick before starting
	startTime := time.Now()

	t := time.NewTimer(attemptDelay)
	defer t.Stop()

	for attempt := int64(0); attempt < c.config.retryOptions.maxRetries+1; attempt++ {
		attemptDelay = c.calculateDelay(attempt)

		t.Reset(attemptDelay)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:

		}

		wm, err = c.put(endpoint, in, out, q)

		// Maximum retry period is up, don't retry
		if c.config.retryOptions.maxToLastCall != 0 && time.Since(startTime) > c.config.retryOptions.maxToLastCall {
			break
		}

		// The put function only returns WriteMetadata if the call was successful
		// don't retry
		if wm != nil {
			break
		}

		// If WriteMetadata is nil, we need to process the error to decide if a retry is
		// necessary or not
		var callErr UnexpectedResponseError
		ok := errors.As(err, &callErr)

		// If is not UnexpectedResponseError, it is an error while performing the call
		// don't retry
		if !ok {
			break
		}

		// Only 500+ or 429 status calls may be retried, otherwise
		// don't retry
		if !isCallRetriable(callErr.StatusCode()) {
			break
		}
	}

	return wm, err
}

// According to the HTTP protocol, it only makes sense to retry calls
// when the error is caused by a temporary situation, like a server being down
// (500s+) or the call being rate limited (429), this function checks if the
// statusCode is between the errors worth retrying.
func isCallRetriable(statusCode int) bool {
	return statusCode > http.StatusInternalServerError &&
		statusCode < http.StatusNetworkAuthenticationRequired ||
		statusCode == http.StatusTooManyRequests
}

func (c *Client) calculateDelay(attempt int64) time.Duration {
	if c.config.retryOptions.fixedDelay != 0 {
		return c.config.retryOptions.fixedDelay
	}

	if attempt == 0 {
		return 0
	}

	if attempt > c.config.retryOptions.maxValidAttempt {
		return c.config.retryOptions.maxBackoffDelay
	}

	newDelay := c.config.retryOptions.delayBase << (attempt - 1)
	if c.config.retryOptions.maxBackoffDelay != defaultMaxBackoffDelay &&
		newDelay > c.config.retryOptions.maxBackoffDelay {
		return c.config.retryOptions.maxBackoffDelay
	}

	return newDelay
}
