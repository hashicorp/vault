package activity

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/timeutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func TestActivityLog_MonthlyActivityApi(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
		ActivityLogConfig: vault.ActivityLogCoreConfig{
			ForceEnable: true,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	core := cluster.Cores[0].Core

	resp, err := client.Logical().Read("sys/internal/counters/activity/monthly")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	c, err := resp.Data["clients"].(json.Number).Int64()
	require.NoError(t, err)
	require.EqualValues(t, 0, c)

	core.InjectActivityLogDataThisMonth(t)

	resp, err = client.Logical().Read("sys/internal/counters/activity/monthly")
	require.NoError(t, err)

	var byNS []*vault.ResponseNamespace
	err = mapstructure.WeakDecode(resp.Data["by_namespace"], &byNS)
	require.NoError(t, err)

	var byMonth []*vault.ResponseMonth
	err = mapstructure.WeakDecode(resp.Data["months"], &byMonth)
	require.NoError(t, err)

	clients := 0
	err = mapstructure.WeakDecode(resp.Data["clients"], &clients)
	require.NoError(t, err)

	// Overall clients
	require.Equal(t, 7, clients)

	require.NotNil(t, byNS[0])
	require.NotNil(t, byNS[0].Counts)

	// Root namespace should have the count of 3
	require.Equal(t, 3, byNS[0].Counts.Clients)

	// Other namespaces should have 2 each
	require.NotNil(t, byNS[1])
	require.NotNil(t, byNS[1].Counts)
	require.Equal(t, 2, byNS[1].Counts.Clients)

	require.NotNil(t, byNS[2])
	require.NotNil(t, byNS[2].Counts)
	require.Equal(t, 2, byNS[2].Counts.Clients)

	// Root ns should have 3 mount entries, other namespaces should have 2 each
	require.Len(t, byNS[0].Mounts, 3)
	require.Len(t, byNS[1].Mounts, 2)
	require.Len(t, byNS[2].Mounts, 2)

	// Months section should also report 7 clients
	require.NotNil(t, byMonth[0])
	require.NotNil(t, byMonth[0].Counts)
	require.Equal(t, 7, byMonth[0].Counts.Clients)

	// New clients within month should also report 7 clients
	require.NotNil(t, byMonth[0].NewClients)
	require.NotNil(t, byMonth[0].NewClients.Counts)
	require.Equal(t, 7, byMonth[0].NewClients.Counts.Clients)

	// we expect a 204 if activity log is disabled
	core.GetActivityLog().SetEnable(false)
	req := client.NewRequest("GET", "/v1/sys/internal/counters/activity/monthly")
	rawResp, err := client.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if rawResp == nil {
		t.Error("nil response")
	}
	if rawResp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code %v, got %v", http.StatusNoContent, rawResp.StatusCode)
	}
}
