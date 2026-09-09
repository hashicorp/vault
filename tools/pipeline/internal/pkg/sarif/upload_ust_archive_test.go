// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"
)

// TestDoUpload verifies the retry behaviour of doUpload across combinations of
// initial-push success/failure, pull success/failure, and exhausted retries.
// The primary regression case is a pull failure mid-retry: we used to return
// immediately on a pull error instead of continuing to the next attempt.
func TestDoUpload(t *testing.T) {
	t.Parallel()

	var (
		errPush = errors.New("push rejected")
		errPull = errors.New("pull rebase conflict")
	)

	for name, test := range map[string]struct {
		maxRetries     int
		pushFails      int
		pullFails      int
		shouldFail     bool
		expectedErrors []error // all errors that should be reachable via errors.Is
		expectedPushes int
		expectedPulls  int
	}{
		// initial push succeeds — no retry path entered at all
		"first push succeeds": {
			maxRetries:     1,
			expectedPushes: 1,
		},
		// pull fails on attempt 1, succeeds on attempt 2, push then succeeds.
		// This is the primary regression case: previously a pull failure
		// returned immediately instead of continuing to the next attempt.
		// Call sequence: push(fail) → pull(fail,skip push) → pull(ok) → push(ok)
		"pull fail then success": {
			maxRetries:     2,
			pushFails:      1,
			pullFails:      1,
			expectedPushes: 2, // initial + push on attempt 2 (attempt 1 skips push after pull failure)
			expectedPulls:  2,
		},
		// pull always succeeds but push always fails — exhausts all retries
		"pull ok push always fails": {
			maxRetries:     3,
			pushFails:      4, // initial + 3 retries all fail
			shouldFail:     true,
			expectedErrors: []error{errPush},
			expectedPushes: 4,
			expectedPulls:  3,
		},
		// both pull and push always fail — exhausts all retries, joined errors
		// contain both pull and push errors
		"pull and push always fail": {
			maxRetries:     2,
			pushFails:      1, // only the initial push; pull failures skip the retry push
			pullFails:      2,
			shouldFail:     true,
			expectedErrors: []error{errPull},
			expectedPushes: 1,
			expectedPulls:  2,
		},
		// UploadRetries=0 is treated as 1 so callers that omit the field still
		// get one pull+push retry after the initial push fails.
		"zero retries defaults to one": {
			maxRetries:     0,
			pushFails:      1,
			expectedPushes: 2,
			expectedPulls:  1,
		},
		// pull fails on the final (only) attempt — the accumulated error
		// contains the pull error
		"pull fail on last attempt": {
			maxRetries:     1,
			pushFails:      1,
			pullFails:      1,
			shouldFail:     true,
			expectedErrors: []error{errPull},
			expectedPushes: 1,
			expectedPulls:  1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				pushCalls, pullCalls := 0, 0
				pushFn := func(_ context.Context) error {
					pushCalls++
					if pushCalls <= test.pushFails {
						return errPush
					}
					return nil
				}
				pullFn := func(_ context.Context) error {
					pullCalls++
					if pullCalls <= test.pullFails {
						return errPull
					}
					return nil
				}

				req := &UploadUSTArchiveReq{UploadRetries: test.maxRetries}
				err := req.doUpload(t.Context(), pushFn, pullFn)

				if test.shouldFail {
					require.Error(t, err)
					require.ErrorContains(t, err, "failed to upload results")
					for _, target := range test.expectedErrors {
						require.ErrorIs(t, err, target)
					}
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, test.expectedPushes, pushCalls, "push call count")
				require.Equal(t, test.expectedPulls, pullCalls, "pull call count")
			})
		})
	}
}
