package seal_migration

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
)

func TestShamir(t *testing.T) {
	testVariousBackends(t, testShamir)
}

func TestTransit(t *testing.T) {
	testVariousBackends(t, testTransit)
}

type testFunc func(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int)

func testVariousBackends(t *testing.T, tf testFunc) {

	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())

	t.Run("inmem", func(t *testing.T) {
		t.Parallel()

		logger := logger.Named("inmem")
		storage, cleanup := teststorage.MakeReusableStorage(
			t, logger, teststorage.MakeInmemBackend(t, logger))
		defer cleanup()
		tf(t, logger, storage, 51000)
	})

	//t.Run("file", func(t *testing.T) {
	//	t.Parallel()

	//	logger := logger.Named("file")
	//	storage, cleanup := teststorage.MakeReusableStorage(
	//		t, logger, teststorage.MakeFileBackend(t, logger))
	//	defer cleanup()
	//	tf(t, logger, storage, 52000)
	//})

	//t.Run("consul", func(t *testing.T) {
	//	t.Parallel()

	//	logger := logger.Named("consul")
	//	storage, cleanup := teststorage.MakeReusableStorage(
	//		t, logger, teststorage.MakeConsulBackend(t, logger))
	//	defer cleanup()
	//	tf(t, logger, storage, 53000)
	//})

	//t.Run("raft", func(t *testing.T) {
	//	t.Parallel()

	//	logger := logger.Named("raft")
	//	storage, cleanup := teststorage.MakeReusableRaftStorage(t, logger)
	//	defer cleanup()
	//	tf(t, logger, storage, 54000)
	//})
}

func testShamir(
	t *testing.T, logger hclog.Logger,
	storage teststorage.ReusableStorage, basePort int) {

	rootToken, barrierKeys := initializeShamir(t, logger, storage, basePort)
	reuseShamir(t, logger, storage, basePort, rootToken, barrierKeys)
}

func testTransit(
	t *testing.T, logger hclog.Logger,
	storage teststorage.ReusableStorage, basePort int) {

	// Create the transit server.
	tss := sealhelper.NewTransitSealServer(t)
	defer func() {
		if tss != nil {
			tss.Cleanup()
		}
	}()
	tss.MakeKey(t, "transit-seal-key")

	rootToken, recoveryKeys, transitSeal := initializeTransit(t, logger, storage, basePort, tss)
	println("rootToken, recoveryKeys, transitSeal", rootToken, recoveryKeys, transitSeal)
	//reuseShamir(t, logger, storage, basePort, rootToken, barrierKeys)
}

// initializeShamir initializes a brand new backend storage with Shamir.
func initializeShamir(
	t *testing.T, logger hclog.Logger,
	storage teststorage.ReusableStorage, basePort int) (string, [][]byte) {

	var baseClusterPort = basePort + 10

	// Start the cluster
	var conf = vault.CoreConfig{
		Logger: logger.Named("initializeShamir"),
	}
	var opts = vault.TestClusterOptions{
		HandlerFunc:           vaulthttp.Handler,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		storage.Cleanup(t, cluster)
		cluster.Cleanup()
	}()

	leader := cluster.Cores[0]
	client := leader.Client

	// Unseal
	if storage.IsRaft {
		testhelpers.RaftClusterJoinNodes(t, cluster)
		if err := testhelpers.VerifyRaftConfiguration(t, leader); err != nil {
			t.Fatal(err)
		}
	} else {
		cluster.UnsealCores(t)
	}
	testhelpers.WaitForNCoresUnsealed(t, cluster, vault.DefaultNumCores)

	// Write a secret that we will read back out later.
	_, err := client.Logical().Write(
		"secret/foo",
		map[string]interface{}{"zork": "quux"})
	if err != nil {
		t.Fatal(err)
	}

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	return cluster.RootToken, cluster.BarrierKeys
}

// reuseShamir uses a pre-populated backend storage with Shamir.
func reuseShamir(
	t *testing.T, logger hclog.Logger,
	storage teststorage.ReusableStorage, basePort int,
	rootToken string, barrierKeys [][]byte) {

	var baseClusterPort = basePort + 10

	// Start the cluster
	var conf = vault.CoreConfig{
		Logger: logger.Named("reuseShamir"),
	}
	var opts = vault.TestClusterOptions{
		HandlerFunc:           vaulthttp.Handler,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
		SkipInit:              true,
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		storage.Cleanup(t, cluster)
		cluster.Cleanup()
	}()

	leader := cluster.Cores[0]
	client := leader.Client
	client.SetToken(rootToken)

	// Unseal
	cluster.BarrierKeys = barrierKeys
	if storage.IsRaft {
		provider := testhelpers.NewHardcodedServerAddressProvider(baseClusterPort)
		testhelpers.SetRaftAddressProviders(t, cluster, provider)

		for _, core := range cluster.Cores {
			cluster.UnsealCore(t, core)
		}
		time.Sleep(15 * time.Second)

		if err := testhelpers.VerifyRaftConfiguration(t, leader); err != nil {
			t.Fatal(err)
		}
	} else {
		cluster.UnsealCores(t)
	}
	testhelpers.WaitForNCoresUnsealed(t, cluster, vault.DefaultNumCores)

	// Read the secret
	secret, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}

	// Seal the cluster
	cluster.EnsureCoresSealed(t)
}

// initializeTransit initializes a brand new backend storage with Transit.
func initializeTransit(
	t *testing.T, logger hclog.Logger,
	storage teststorage.ReusableStorage, basePort int,
	tss *sealhelper.TransitSealServer) (string, [][]byte, vault.Seal) {

	var transitSeal vault.Seal

	var baseClusterPort = basePort + 10

	// Start the cluster
	var conf = vault.CoreConfig{
		Logger: logger.Named("initializeTransit"),
	}
	var opts = vault.TestClusterOptions{
		HandlerFunc:           vaulthttp.Handler,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
		SealFunc: func() vault.Seal {
			// Each core will create its own transit seal here.  Later
			// on it won't matter which one of these we end up using, since
			// they were all created from the same transit key.
			transitSeal = tss.MakeSeal(t, "transit-seal-key")
			return transitSeal
		},
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		storage.Cleanup(t, cluster)
		cluster.Cleanup()
	}()

	leader := cluster.Cores[0]
	client := leader.Client

	// Unseal
	if storage.IsRaft {
		testhelpers.RaftClusterJoinNodes(t, cluster)
		if err := testhelpers.VerifyRaftConfiguration(t, leader); err != nil {
			t.Fatal(err)
		}
	} else {
		cluster.UnsealCores(t)
	}
	testhelpers.WaitForNCoresUnsealed(t, cluster, vault.DefaultNumCores)

	// Write a secret that we will read back out later.
	_, err := client.Logical().Write(
		"secret/foo",
		map[string]interface{}{"zork": "quux"})
	if err != nil {
		t.Fatal(err)
	}

	// Seal the cluster
	cluster.EnsureCoresSealed(t)

	return cluster.RootToken, cluster.RecoveryKeys, transitSeal
}
