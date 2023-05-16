// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
)

func NewReplicationSetDocker(t *testing.T) (*testcluster.ReplicationSet, error) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}

	r := &testcluster.ReplicationSet{
		Clusters: map[string]testcluster.VaultCluster{},
		Logger:   logging.NewVaultLogger(hclog.Trace).Named(t.Name()),
	}

	r.Builder = func(ctx context.Context, name string, baseLogger hclog.Logger) (testcluster.VaultCluster, error) {
		cluster := NewTestDockerCluster(t, &DockerClusterOptions{
			ImageRepo:   "hashicorp/vault",
			ImageTag:    "latest",
			VaultBinary: os.Getenv("VAULT_BINARY"),
			ClusterOptions: testcluster.ClusterOptions{
				NumCores:    5,
				ClusterName: strings.ReplaceAll(t.Name()+"-"+name, "/", "-"),
				Logger:      baseLogger.Named(name),
				VaultNodeConfig: &testcluster.VaultNodeConfig{
					LogLevel: "TRACE",
					// If you want the test to run faster locally, you could
					// uncomment this performance_multiplier change.
					//StorageOptions: map[string]string{
					//	"performance_multiplier": "1",
					//},
				},
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
