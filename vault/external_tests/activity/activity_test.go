package activity

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func validateClientCounts(t *testing.T, resp *api.Secret, expectedEntities, expectedTokens int) {
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Data == nil {
		t.Fatal("no data")
	}

	expectedClients := expectedEntities + expectedTokens

	entityCountJSON, ok := resp.Data["distinct_entities"]
	if !ok {
		t.Fatalf("no entity count: %v", resp.Data)
	}
	entityCount, err := entityCountJSON.(json.Number).Int64()
	if err != nil {
		t.Fatal(err)
	}
	if entityCount != int64(expectedEntities) {
		t.Errorf("bad entity count. expected %v, got %v", expectedEntities, entityCount)
	}

	tokenCountJSON, ok := resp.Data["non_entity_tokens"]
	if !ok {
		t.Fatalf("no token count: %v", resp.Data)
	}
	tokenCount, err := tokenCountJSON.(json.Number).Int64()
	if err != nil {
		t.Fatal(err)
	}
	if tokenCount != int64(expectedTokens) {
		t.Errorf("bad token count. expected %v, got %v", expectedTokens, tokenCount)
	}

	clientCountJSON, ok := resp.Data["clients"]
	if !ok {
		t.Fatalf("no client count: %v", resp.Data)
	}
	clientCount, err := clientCountJSON.(json.Number).Int64()
	if err != nil {
		t.Fatal(err)
	}
	if clientCount != int64(expectedClients) {
		t.Errorf("bad client count. expected %v, got %v", expectedClients, clientCount)
	}
}

func TestActivityLog_MonthlyActivityApi(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	validateClientCounts(t, resp, 0, 0)

	// inject some data and query the API
	entities, tokens := core.InjectActivityLogDataThisMonth(t)
	expectedEntities := len(entities)
	var expectedTokens int
	for _, tokenCount := range tokens {
		expectedTokens += int(tokenCount)
	}

	resp, err = client.Logical().Read("sys/internal/counters/activity/monthly")
	if err != nil {
		t.Fatal(err)
	}
	validateClientCounts(t, resp, expectedEntities, expectedTokens)

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
