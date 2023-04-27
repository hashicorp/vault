package vault

import (
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// Test_AuditBroker_ForwardingHeaders tests that we are able to extract the relevant headers
// from a logical.LogInput Request's headers, and move the values of the headers to explicit
// fields within the logical.LogInput struct.
// The headers originate from the HTTP request headers that the node received, and should have
// been augmented with metadata about forwarding from/to hosts if a request is forwarded from
// a standby to a primary node.
func Test_AuditBroker_ForwardingHeaders(t *testing.T) {
	tests := map[string]struct {
		headers  map[string][]string
		wantNil  bool
		wantFrom string
		wantTo   string
	}{
		"none": {
			headers: map[string][]string{},
			wantNil: true,
		},
		"from": {
			headers: map[string][]string{
				HTTPHeaderVaultForwardFrom: {"juan:8080"},
			},
			wantNil:  false,
			wantFrom: "juan:8080",
		},
		"to": {
			headers: map[string][]string{
				HTTPHeaderVaultForwardTo: {"john:8000"},
			},
			wantNil: false,
			wantTo:  "john:8000",
		},
		"from-and-to": {
			headers: map[string][]string{
				HTTPHeaderVaultForwardFrom: {"juan:8080"},
				HTTPHeaderVaultForwardTo:   {"john:8000"},
			},
			wantNil:  false,
			wantFrom: "juan:8080",
			wantTo:   "john:8000",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logInput := &logical.LogInput{Request: &logical.Request{Headers: tc.headers}}

			extractForwardingHeaders(logInput)

			if tc.wantNil {
				require.Nil(t, logInput.Forwarding)
				return
			}

			require.NotNil(t, logInput.Forwarding)
			require.Equal(t, tc.wantFrom, logInput.Forwarding.From)
			if tc.wantFrom != "" {
				require.NotContains(t, logInput.Request.Headers, HTTPHeaderVaultForwardFrom)
			}

			require.Equal(t, tc.wantTo, logInput.Forwarding.To)
			if tc.wantTo != "" {
				require.NotContains(t, logInput.Request.Headers, HTTPHeaderVaultForwardTo)
			}
		})
	}
}
