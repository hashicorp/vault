// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cachememdb

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerializeDeserialize(t *testing.T) {
	testIndex := &Index{
		ID:            "testid",
		Token:         "testtoken",
		TokenParent:   "parent token",
		TokenAccessor: "test accessor",
		Namespace:     "test namespace",
		RequestPath:   "/test/path",
		Lease:         "lease id",
		LeaseToken:    "lease token id",
		Response:      []byte(`{"something": "here"}`),
		RenewCtxInfo:  NewContextInfo(context.Background()),
		RequestMethod: "GET",
		RequestToken:  "request token",
		RequestHeader: http.Header{
			"X-Test": []string{"vault", "agent"},
		},
		LastRenewed: time.Now().UTC(),
	}
	indexBytes, err := testIndex.Serialize()
	require.NoError(t, err)
	assert.True(t, len(indexBytes) > 0)
	assert.NotNil(t, testIndex.RenewCtxInfo, "Serialize should not modify original Index object")

	restoredIndex, err := Deserialize(indexBytes)
	require.NoError(t, err)

	testIndex.RenewCtxInfo = nil
	assert.Equal(t, testIndex, restoredIndex, "They should be equal without RenewCtxInfo set on the original")
}
