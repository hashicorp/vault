package http

import (
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysHostInfo(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	// Query against the active node, should get host information back
	info, err := cores[0].Client.Sys().HostInfo()
	if err != nil {
		t.Fatal(err)
	}
	if info == nil {
		t.Fatal("expected non-nil HostInfo")
	}

	if info.Timestamp == "" {
		t.Fatalf("expected a valid timestamp")
	}
	if info.CPU == nil {
		t.Fatal("expected CPU info")
	}
	if info.Disk == nil {
		t.Fatal("expected disk info")
	}
	if info.Host == nil {
		t.Fatal("expected host info")
	}
	if info.Memory == nil {
		t.Fatal("expected memory info")
	}

	// Query against the standby, should error
	info, err = cores[1].Client.Sys().HostInfo()
	if err == nil || info != nil {
		t.Fatalf("expected error on standby node, HostInfo: %v", info)
	}
}
