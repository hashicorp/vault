package diagnose

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/physical"
)

func TestStorageTimeout(t *testing.T) {

	testCases := []struct {
		errSubString string
		mb           physical.Backend
	}{
		{
			errSubString: timeOutErr + "operation: Put",
			mb:           mockStorageBackend{callType: timeoutCallWrite},
		},
		{
			errSubString: timeOutErr + "operation: Get",
			mb:           mockStorageBackend{callType: timeoutCallRead},
		},
		{
			errSubString: timeOutErr + "operation: Delete",
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
		{
			errSubString: "",
			mb:           mockStorageBackend{callType: ""},
		},
	}

	for _, tc := range testCases {
		outErr := StorageEndToEndLatencyCheck(context.Background(), tc.mb)
		if tc.errSubString == "" && outErr == nil {
			// this is the success case where the Storage Latency check passes
			continue
		}
		if !strings.Contains(outErr.Error(), tc.errSubString) {
			t.Errorf("wrong error: expected %s to be contained in %s", tc.errSubString, outErr)
		}
	}
}

// mb := mockStorageBackend{callType: timeoutCallRead}
// err := StorageEndToEndLatencyCheck(context.Background(), mb)
