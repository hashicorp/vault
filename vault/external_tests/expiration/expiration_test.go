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

	countPerMountRaw, ok := resp.Data["counts"]
	if !ok {
		t.Fatalf("expected 'counts' response, got %#v", resp.Data)
	}
	countPerMount := countPerMountRaw.(map[string]interface{})
	if len(countPerMount) != 0 {
		t.Errorf("expected no mounts with counts, got %#v", countPerMount)
	}

	expectedNumLeases := 50
	expectedCountPerMount := core.InjectIrrevocableLeases(t, namespace.RootContext(nil), expectedNumLeases)

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

	countPerMountRaw, ok = resp.Data["counts"]
	if !ok {
		t.Fatalf("expected 'counts' response, got %#v", resp.Data)
	}

	countPerMount = countPerMountRaw.(map[string]interface{})
	for mount, expectedCount := range expectedCountPerMount {
		gotCountRaw, ok := countPerMount[mount]
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

func TestExpiration_irrevocableLeaseListAPI(t *testing.T) {
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
	// TODO update RFC to say that it's a get operation
	resp, err := client.Logical().ReadWithData("sys/leases", params)
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

	leasesPerMountRaw, ok := resp.Data["leases"]
	if !ok {
		t.Fatalf("expected 'leases' response, got %#v", resp.Data)
	}
	leasesPerMount := leasesPerMountRaw.(map[string]interface{})
	if len(leasesPerMount) != 0 {
		t.Errorf("expected no mounts with leases, got %#v", leasesPerMount)
	}

	// test with a low enough number to not give an error without force flag
	expectedNumLeases := 50
	expectedCountsPerMount := core.InjectIrrevocableLeases(t, namespace.RootContext(nil), expectedNumLeases)

	resp, err = client.Logical().ReadWithData("sys/leases", params)
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

	leasesPerMountRaw, ok = resp.Data["leases"]
	if !ok {
		t.Fatalf("expected 'leases' response, got %#v", resp.Data)
	}

	leasesPerMount = leasesPerMountRaw.(map[string]interface{})
	for mount, expectedCount := range expectedCountsPerMount {
		leaseCount := len(leasesPerMount[mount].([]interface{}))
		if leaseCount != expectedCount {
			t.Errorf("bad count for mount %q, expected %d, got %d", mount, expectedCount, leaseCount)
		}
	}
}

func TestExpiration_irrevocableLeaseListAPI_force(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	core := cluster.Cores[0].Core

	// test with a low enough number to not give an error without force flag
	expectedNumLeases := vault.MaxIrrevocableLeasesToReturn + 50
	expectedCountsPerMount := core.InjectIrrevocableLeases(t, namespace.RootContext(nil), expectedNumLeases)

	params := make(map[string][]string)
	params["type"] = []string{"irrevocable"}

	resp, err := client.Logical().ReadWithData("sys/leases", params)
	if err == nil {
		t.Fatalf("expected error without force flag")
	}
	if resp != nil {
		t.Errorf("expected nil response, got: %#v", resp)
	}

	// now try it with the force flag - we expect no errors and many results
	params["force"] = []string{"true"}
	resp, err = client.Logical().ReadWithData("sys/leases", params)
	if err != nil {
		t.Fatalf("unexpected error when using force flag: %v", err)
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
	if totalLeaseCount != int64(expectedNumLeases) {
		t.Errorf("expected %d leases, got %d", expectedNumLeases, totalLeaseCount)
	}

	leasesPerMountRaw, ok := resp.Data["leases"]
	if !ok {
		t.Fatalf("expected 'leases' response, got %#v", resp.Data)
	}

	leasesPerMount := leasesPerMountRaw.(map[string]interface{})
	for mount, expectedCount := range expectedCountsPerMount {
		leaseCount := len(leasesPerMount[mount].([]interface{}))
		if leaseCount != expectedCount {
			t.Errorf("bad count for mount %q, expected %d, got %d", mount, expectedCount, leaseCount)
		}
	}
}
