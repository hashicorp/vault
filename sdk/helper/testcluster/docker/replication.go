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

func DefaultOptions(t *testing.T) *DockerClusterOptions {
	return &DockerClusterOptions{
		ImageRepo:   "hashicorp/vault",
		ImageTag:    "latest",
		VaultBinary: os.Getenv("VAULT_BINARY"),
		ClusterOptions: testcluster.ClusterOptions{
			NumCores:    3,
			ClusterName: strings.ReplaceAll(t.Name(), "/", "-"),
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
			},
		},
	}
}

func NewReplicationSetDocker(t *testing.T, opts *DockerClusterOptions) (*testcluster.ReplicationSet, error) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}

	r := &testcluster.ReplicationSet{
		Clusters: map[string]testcluster.VaultCluster{},
		Logger:   logging.NewVaultLogger(hclog.Trace).Named(t.Name()),
	}

	// clusterName is used for container name as well.
	// A container name should not exceed 64 chars.
	// There are additional chars that are added to the name as well
	// like "-A-core0". So, setting a max limit for a cluster name.
	if len(opts.ClusterName) > MaxClusterNameLength {
		return nil, fmt.Errorf("cluster name length exceeded the maximum allowed length of %v", MaxClusterNameLength)
	}

	r.Builder = func(ctx context.Context, name string, baseLogger hclog.Logger) (testcluster.VaultCluster, error) {
		myOpts := *opts
		myOpts.Logger = baseLogger.Named(name)
		if myOpts.ClusterName == "" {
			myOpts.ClusterName = strings.ReplaceAll(t.Name(), "/", "-")
		}
		myOpts.ClusterName += "-" + strings.ReplaceAll(name, "/", "-")
		myOpts.CA = r.CA
		return NewTestDockerCluster(t, &myOpts), nil
	}

	a, err := r.Builder(context.TODO(), "A", r.Logger)
	if err != nil {
		return nil, err
	}
	r.Clusters["A"] = a
	r.CA = a.(*DockerCluster).CA

	return r, err
}
