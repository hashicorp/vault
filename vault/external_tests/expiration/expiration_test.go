// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package expiration

import (
	"encoding/json"
	"reflect"
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

	if len(resp.Warnings) > 0 {
		t.Errorf("expected no warnings, got: %v", resp.Warnings)
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
	expectedCountPerMount, err := core.InjectIrrevocableLeases(namespace.RootContext(nil), expectedNumLeases)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().ReadWithData("sys/leases/count", params)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}

	if len(resp.Warnings) > 0 {
		t.Errorf("expected no warnings, got: %v", resp.Warnings)
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
	if len(countPerMount) != len(expectedCountPerMount) {
		t.Fatalf("expected %d mounts, got %d: %#v", len(expectedCountPerMount), len(countPerMount), countPerMount)
	}

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
	resp, err := client.Logical().ReadWithData("sys/leases", params)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}

	if len(resp.Warnings) > 0 {
		t.Errorf("expected no warnings, got: %v", resp.Warnings)
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

	leasesRaw, ok := resp.Data["leases"]
	if !ok {
		t.Fatalf("expected 'leases' response, got %#v", resp.Data)
	}
	leases := leasesRaw.([]interface{})
	if len(leases) != 0 {
		t.Errorf("expected no mounts with leases, got %#v", leases)
	}

	// test with a low enough number to not give an error without limit set to none
	expectedNumLeases := 50
	expectedCountPerMount, err := core.InjectIrrevocableLeases(namespace.RootContext(nil), expectedNumLeases)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().ReadWithData("sys/leases", params)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}

	if len(resp.Warnings) > 0 {
		t.Errorf("expected no warnings, got: %v", resp.Warnings)
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

	leasesRaw, ok = resp.Data["leases"]
	if !ok {
		t.Fatalf("expected 'leases' response, got %#v", resp.Data)
	}

	leases = leasesRaw.([]interface{})
	countPerMount := make(map[string]int)
	for _, leaseRaw := range leases {
		lease := leaseRaw.(map[string]interface{})
		mount := lease["mount_id"].(string)

		if _, ok := countPerMount[mount]; !ok {
			countPerMount[mount] = 0
		}

		countPerMount[mount]++
	}

	if !reflect.DeepEqual(countPerMount, expectedCountPerMount) {
		t.Errorf("bad mount count. expected %v, got %v", expectedCountPerMount, countPerMount)
	}
}

func TestExpiration_irrevocableLeaseListAPI_includeAll(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	core := cluster.Cores[0].Core

	// test with a low enough number to not give an error with the default limit
	expectedNumLeases := vault.MaxIrrevocableLeasesToReturn + 50
	expectedCountPerMount, err := core.InjectIrrevocableLeases(namespace.RootContext(nil), expectedNumLeases)
	if err != nil {
		t.Fatal(err)
	}

	params := make(map[string][]string)
	params["type"] = []string{"irrevocable"}

	resp, err := client.Logical().ReadWithData("sys/leases", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("unexpected nil response")
	}

	if len(resp.Warnings) != 1 {
		t.Errorf("expected one warning (%q), got: %v", vault.MaxIrrevocableLeasesWarning, resp.Warnings)
	}

	// now try it with the no limit on return size - we expect no errors and many results
	params["limit"] = []string{"none"}
	resp, err = client.Logical().ReadWithData("sys/leases", params)
	if err != nil {
		t.Fatalf("unexpected error when using limit=none: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}

	if len(resp.Warnings) > 0 {
		t.Errorf("expected no warnings, got: %v", resp.Warnings)
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

	leasesRaw, ok := resp.Data["leases"]
	if !ok {
		t.Fatalf("expected 'leases' response, got %#v", resp.Data)
	}

	leases := leasesRaw.([]interface{})
	countPerMount := make(map[string]int)
	for _, leaseRaw := range leases {
		lease := leaseRaw.(map[string]interface{})
		mount := lease["mount_id"].(string)

		if _, ok := countPerMount[mount]; !ok {
			countPerMount[mount] = 0
		}

		countPerMount[mount]++
	}

	if !reflect.DeepEqual(countPerMount, expectedCountPerMount) {
		t.Errorf("bad mount count. expected %v, got %v", expectedCountPerMount, countPerMount)
	}
}
