// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package replication_binary

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

func TestReplication_FailoverPrimaryActive(t *testing.T) {
	//if !constants.IsEnterprise {
	//	// Disable on OSS since this needs an ent binary (or docker image) to work
	//	t.Skip()
	//}

	opts := docker.DefaultOptions(t)
	opts.NumCores = 5
	r, err := docker.NewReplicationSetDocker(t, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cleanup()

	err = r.StandardPerfReplication(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	a, c := r.Clusters["A"], r.Clusters["C"]
	c0 := c.Nodes()[0]
	err = testcluster.WaitForPerfReplicationStatus(ctx, c0.APIClient(), func(data map[string]interface{}) error {
		found := data["known_primary_cluster_addrs"]
		if len(found.([]interface{})) == len(a.Nodes()) {
			return nil
		}
		return fmt.Errorf("expected %d known_primary_cluster_addrs, got: %#v", len(a.Nodes()), found)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = testcluster.WaitForPerfReplicationConnectionStatus(ctx, c0.APIClient())
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		idx, err := testcluster.LeaderNode(ctx, a)
		if err != nil {
			t.Fatal(err)
		}
		leader := a.Nodes()[idx]

		priAPIAddrRaw := leader.(*docker.DockerClusterNode).RealAPIAddr
		priAPIAddr, err := url.Parse(priAPIAddrRaw)
		if err != nil {
			t.Fatalf("bad api addr %q: %v", priAPIAddrRaw, err)
		}
		err = c0.(*docker.DockerClusterNode).AddNetworkDelay(ctx, 10*time.Second, strings.Split(priAPIAddr.Host, ":")[0])
		if err != nil {
			t.Fatal(fmt.Sprintf("delaying sec node 0 traffic to pri node 0: %s", err))
		}
		err = leader.(*docker.DockerClusterNode).Pause(ctx)
		if err != nil {
			t.Fatal(fmt.Sprintf("pausing node 0: %s", err))
		}
		time.Sleep(5 * time.Second)

		err = testcluster.WaitForPerfReplicationWorking(ctx, a, c)
		if err != nil {
			t.Fatal(err)
		}

		err = testcluster.WaitForPerfReplicationConnectionStatus(ctx, c0.APIClient())
		if err != nil {
			t.Fatal(err)
		}

	}
}
