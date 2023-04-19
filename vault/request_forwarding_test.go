package vault

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_RequestForwarding_ForwardingHeaders exercises the functionality which adds metadata to an HTTP
// request received by a standby node, this would happen before we forward it to the primary node.
func Test_RequestForwarding_ForwardingHeaders(t *testing.T) {
	cluster := NewTestCluster(t, nil, nil)
	cluster.Start()
	defer cluster.Cleanup()

	leader := cluster.Cores[0]
	standby := cluster.Cores[1]
	standbyURL, err := url.Parse(standby.redirectAddr)
	require.NoError(t, err)
	require.NotNil(t, standbyURL)

	isLeader, leaderAddr, _, err := leader.Leader()
	require.NoError(t, err)
	require.True(t, isLeader)
	leaderURL, err := url.Parse(leaderAddr)
	require.NoError(t, err)
	require.NotNil(t, leaderURL)

	req := &http.Request{}
	err = standby.addForwardingHeaders(req)
	require.NoError(t, err)

	from := req.Header.Get("X-Vault-Forwarded-From")
	to := req.Header.Get("X-Vault-Forwarded-To")

	require.Equal(t, standbyURL.Host, from)
	require.Equal(t, leaderURL.Host, to)
}
