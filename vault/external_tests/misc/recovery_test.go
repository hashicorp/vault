package misc

import (
	"github.com/hashicorp/vault/helper/testhelpers"
	"go.uber.org/atomic"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
)

func TestRecovery(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())
	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	conf := vault.CoreConfig{
		Physical: inm,
		Logger:   logger,
	}
	opts := vault.TestClusterOptions{
		HandlerFunc: http.Handler,
		NumCores:    1,
	}

	cluster := vault.NewTestCluster(t, &conf, &opts)
	defer cluster.Cleanup()
	cluster.EnsureCoresSealed(t)

	// Now bring it up in recovery mode.
	var tokenRef atomic.String
	opts.DefaultHandlerProperties.RecoveryMode = true
	opts.DefaultHandlerProperties.RecoveryToken = &tokenRef
	opts.SkipInit = true
	conf.RecoveryMode = true
	keys := cluster.BarrierKeys
	cluster = vault.NewTestCluster(t, &conf, &opts)
	cluster.BarrierKeys = keys
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	client.SetToken(testhelpers.GenerateRoot(t, cluster, testhelpers.GenerateRecovery))
	client.Logical().List("sys/raw")

}
