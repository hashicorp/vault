package expiration

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func TestExpiration_irrevocableLeaseCountsAPI(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	core := cluster.Cores[0].Core

	params := make(map[string][]string)
	params["type"] = []string{"irrevocable"}
	resp, err := client.Logical().ReadWithData("sys/leases/count", params)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	totalLeaseCountRaw, ok := resp.Data["lease_count"]
	if !ok {
		t.Fatalf("expected 'lease_count' response, got: %#v", resp.Data)
	}

	totalLeaseCount, err := totalLeaseCountRaw.(json.Number).Int64()
	if err != nil {
		t.Fatalf("error extracting lease count: %v", err)
	}
	if totalLeaseCount != 0 {
		t.Errorf("expected no leases, got %d", totalLeaseCount)
	}

	expectedNumLeases := 50
	expectedCountsPerMount := core.InjectIrrevocableLeases(t, namespace.RootContext(nil), expectedNumLeases)

	resp, err = client.Logical().ReadWithData("sys/leases/count", params)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	totalLeaseCountRaw, ok = resp.Data["lease_count"]
	if !ok {
		t.Fatalf("expected 'lease_count' response, got: %#v", resp.Data)
	}

	totalLeaseCount, err = totalLeaseCountRaw.(json.Number).Int64()
	if err != nil {
		t.Fatalf("error extracting lease count: %v", err)
	}
	if totalLeaseCount != int64(expectedNumLeases) {
		t.Errorf("expected %d leases, got %d", expectedNumLeases, totalLeaseCount)
	}

	countsPerMountRaw, ok := resp.Data["counts"]
	if !ok {
		t.Fatalf("expected 'counts' response, got %#v", resp.Data)
	}

	countsPerMount := countsPerMountRaw.(map[string]interface{})
	for mount, expectedCount := range expectedCountsPerMount {
		gotCountRaw, ok := countsPerMount[mount]
		if !ok {
			t.Errorf("missing mount %q", mount)
			continue
		}

		gotCount, err := gotCountRaw.(json.Number).Int64()
		if err != nil {
			t.Errorf("error extracting lease count for mount %q: %v", mount, err)
			continue
		}
		if gotCount != int64(expectedCount) {
			t.Errorf("bad count for mount %q: expected: %d, got: %d", mount, expectedCount, gotCount)
		}
	}
}
