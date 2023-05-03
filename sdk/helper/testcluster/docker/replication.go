// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
)

type ReplicationDockerOptions struct {
	NumCores    int
	ClusterName string
}

func NewReplicationSetDocker(t *testing.T, opt *ReplicationDockerOptions) (*testcluster.ReplicationSet, error) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}

	r := &testcluster.ReplicationSet{
		Clusters: map[string]testcluster.VaultCluster{},
		Logger:   logging.NewVaultLogger(hclog.Trace).Named(t.Name()),
	}

	if opt == nil {
		opt = &ReplicationDockerOptions{}
	}

	var nc int
	if opt.NumCores > 0 {
		nc = opt.NumCores
	}

	clusterName := t.Name()
	if opt.ClusterName != "" {
		clusterName = opt.ClusterName
	}
	// clusterName is used for container name as well.
	// A container name should not exceed 64 chars.
	// There are additional chars that are added to the name as well
	// like "-A-core0". So, setting a max limit for a cluster name.
	if len(clusterName) > MaxClusterNameLength {
		return nil, fmt.Errorf("cluster name length exceeded the maximum allowed length of %v", MaxClusterNameLength)
	}

	r.Builder = func(ctx context.Context, name string, baseLogger hclog.Logger) (testcluster.VaultCluster, error) {
		cluster := NewTestDockerCluster(t, &DockerClusterOptions{
			ImageRepo:   "hashicorp/vault",
			ImageTag:    "latest",
			VaultBinary: os.Getenv("VAULT_BINARY"),
			ClusterOptions: testcluster.ClusterOptions{
				ClusterName: strings.ReplaceAll(clusterName+"-"+name, "/", "-"),
				Logger:      baseLogger.Named(name),
				VaultNodeConfig: &testcluster.VaultNodeConfig{
					LogLevel: "TRACE",
					// If you want the test to run faster locally, you could
					// uncomment this performance_multiplier change.
					//StorageOptions: map[string]string{
					//	"performance_multiplier": "1",
					//},
				},
				NumCores: nc,
			},
			CA: r.CA,
		})
		return cluster, nil
	}

	a, err := r.Builder(context.TODO(), "A", r.Logger)
	if err != nil {
		return nil, err
	}
	r.Clusters["A"] = a
	r.CA = a.(*DockerCluster).CA

	return r, err
}
