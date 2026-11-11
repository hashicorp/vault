// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoginResponseHasNoData ensures a successful Okta login does not carry
// an empty, non-nil Data map: it is serialized as "data": {}, which the UI's
// generated client prefers over the "auth" block (json.data || json.auth),
// discarding the issued token and logging the user out immediately (#32193).
func TestLoginResponseHasNoData(t *testing.T) {
	resp := newOktaLoginResponse()
	require.Nil(t, resp.Data)

	raw, err := json.Marshal(resp)
	require.NoError(t, err)
	require.NotContains(t, string(raw), `"data":{}`)
}
