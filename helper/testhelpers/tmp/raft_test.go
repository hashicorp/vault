package tmp

import (
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/vault"
	"testing"
)

func TestRaftSnapshot(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(nil, nil, teststorage.RaftBackendSetup)
	opts.SetupFunc = nil
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
}
