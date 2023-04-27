// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
)

func TestStorageTimeout(t *testing.T) {
	testCases := []struct {
		errSubString string
		mb           physical.Backend
	}{
		{
			errSubString: LatencyWarning,
			mb:           mockStorageBackend{callType: timeoutCallWrite},
		},
		{
			errSubString: LatencyWarning,
			mb:           mockStorageBackend{callType: timeoutCallRead},
		},
		{
			errSubString: LatencyWarning,
			mb:           mockStorageBackend{callType: timeoutCallDelete},
		},
		{
			errSubString: storageErrStringWrite,
			mb:           mockStorageBackend{callType: errCallWrite},
		},
		{
			errSubString: storageErrStringDelete,
			mb:           mockStorageBackend{callType: errCallDelete},
		},
		{
			errSubString: storageErrStringRead,
			mb:           mockStorageBackend{callType: errCallRead},
		},
		{
			errSubString: wrongRWValsPrefix,
			mb:           mockStorageBackend{callType: badReadCall},
		},
	}

	for _, tc := range testCases {
		var outErr error
		var dur time.Duration
		uuid := "foo"
		backendCallType := tc.mb.(mockStorageBackend).callType
		if callTypeToOp(backendCallType) == readOp {
			dur, outErr = EndToEndLatencyCheckRead(context.Background(), uuid, tc.mb)
		}
		if callTypeToOp(backendCallType) == writeOp {
			dur, outErr = EndToEndLatencyCheckWrite(context.Background(), uuid, tc.mb)
		}
		if callTypeToOp(backendCallType) == deleteOp {
			dur, outErr = EndToEndLatencyCheckDelete(context.Background(), uuid, tc.mb)
		}

		if tc.errSubString == "" && outErr == nil {
			// this is the success case where the Storage Latency check passes
			continue
		}
		if tc.errSubString == LatencyWarning && dur > time.Duration(0) {
			// this is the success case where the Storage Latency check successfully returns nonzero duration
			continue
		}
		if !strings.Contains(outErr.Error(), tc.errSubString) {
			t.Errorf("wrong error: expected %s to be contained in %s", tc.errSubString, outErr)
		}
	}
}
