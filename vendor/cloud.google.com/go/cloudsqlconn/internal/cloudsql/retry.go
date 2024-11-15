// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudsql

import (
	"context"
	"math"
	"math/rand"
	"time"

	"google.golang.org/api/googleapi"
)

// exponentialBackoff calculates a duration based on the attempt i.
//
// The formula is:
//
//	base * multi^(attempt + 1 + random)
//
// With base = 200ms and multi = 1.618, and random = [0.0, 1.0),
// the backoff values would fall between the following low and high ends:
//
// Attempt  Low (ms)  High (ms)
//
//	0         324	     524
//	1         524	     847
//	2         847	    1371
//	3        1371	    2218
//	4        2218	    3588
//
// The theoretical worst case scenario would have a client wait 8.5s in total
// for an API request to complete (with the first four attempts failing, and
// the fifth succeeding).
//
// This backoff strategy matches the behavior of the Cloud SQL Proxy v1.
func exponentialBackoff(attempt int) time.Duration {
	const (
		base  = float64(200 * time.Millisecond)
		multi = 1.618
	)
	exp := float64(attempt+1) + rand.Float64()
	return time.Duration(base * math.Pow(multi, exp))
}

// retry50x will retry any 50x HTTP response up to maxRetries times. The
// backoffFunc determines the duration to wait between attempts.
func retry50x[T any](
	ctx context.Context,
	f func(context.Context) (*T, error),
	waitDuration func(int) time.Duration,
) (*T, error) {
	const maxRetries = 5
	var (
		resp *T
		err  error
	)
	for i := 0; i < maxRetries; i++ {
		resp, err = f(ctx)
		// If err is nil, break and return the response.
		if err == nil {
			break
		}

		gErr, ok := err.(*googleapi.Error)
		// If err is not a googleapi.Error, don't retry.
		if !ok {
			return nil, err
		}
		// If the error code is not a 50x error, don't retry.
		if gErr.Code < 500 {
			return nil, err
		}

		if wErr := wait(ctx, waitDuration(i)); wErr != nil {
			err = wErr
			break
		}

	}
	return resp, err
}

// wait will block until the provided duration passes or the context is
// canceled, whatever happens first.
func wait(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
