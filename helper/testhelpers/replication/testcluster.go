// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package replication

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// SetCorePerf returns a ReplicationSet using NewTestCluster,
// i.e. core-based rather than subprocess- or docker-based clusters.
// The set will contain two clusters A and C connected using perf replication.
func SetCorePerf(t *testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *testcluster.ReplicationSet {
	r := NewReplicationSetCore(t, conf, opts, teststorage.InmemBackendSetup)
	t.Cleanup(r.Cleanup)

	// By default NewTestCluster will mount a kv under secret/.  This isn't
	// done by docker-based clusters, so remove this to make us more like that.
	require.Nil(t, r.Clusters["A"].Nodes()[0].APIClient().Sys().Unmount("secret"))

	err := r.StandardPerfReplication(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func NewReplicationSetCore(t *testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions, setup teststorage.ClusterSetupMutator) *testcluster.ReplicationSet {
	r := &testcluster.ReplicationSet{
		Clusters: map[string]testcluster.VaultCluster{},
		Logger:   logging.NewVaultLogger(hclog.Trace).Named(t.Name()),
	}

	r.Builder = func(ctx context.Context, name string, baseLogger hclog.Logger) (testcluster.VaultCluster, error) {
		conf, opts := teststorage.ClusterSetup(conf, opts, setup)
		opts.Logger = baseLogger.Named(name)
		return vault.NewTestCluster(t, conf, opts), nil
	}

	a, err := r.Builder(context.TODO(), "A", r.Logger)
	if err != nil {
		t.Fatal(err)
	}
	r.Clusters["A"] = a

	return r
}
