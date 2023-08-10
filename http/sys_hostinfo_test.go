// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package http

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/vault/helper/hostutil"
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
	secret, err := cores[0].Client.Logical().Read("sys/host-info")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data == nil {
		t.Fatal("expected data in the response")
	}

	dataBytes, err := json.Marshal(secret.Data)
	if err != nil {
		t.Fatal(err)
	}

	var info hostutil.HostInfo
	if err := json.Unmarshal(dataBytes, &info); err != nil {
		t.Fatal(err)
	}

	if info.Timestamp.IsZero() {
		t.Fatal("expected non-zero Timestamp")
	}
	if info.CPU == nil {
		t.Fatal("expected non-nil CPU value")
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

	// Query against a standby, should error
	secret, err = cores[1].Client.Logical().Read("sys/host-info")
	if err == nil || secret != nil {
		t.Fatalf("expected error on standby node, HostInfo: %v", secret)
	}
}
