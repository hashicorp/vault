package physical

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
)

const (
	keyShares    = 5
	keyThreshold = 3
)

func TestReusableStorage(t *testing.T) {

	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())

	//t.Run("inmem", func(t *testing.T) {
	//	t.Parallel()

	//	logger := logger.Named("inmem")
	//	storage, cleanup := teststorage.MakeReusableStorage(
	//		t, logger, teststorage.MakeInmemBackend(t, logger))
	//	defer cleanup()
	//	testReusableStorage(t, logger, storage)
	//})

	t.Run("raft", func(t *testing.T) {
		t.Parallel()

		logger := logger.Named("raft")
		storage, cleanup := teststorage.MakeReusableRaftStorage(t, logger)
		defer cleanup()
		testReusableStorage(t, logger, storage)
	})
}

func testReusableStorage(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage) {
	rootToken, keys := initializeStorage(t, logger, storage)
	fmt.Printf("rootToken %v, keys %v\n", rootToken, keys)
}

// initializeStorage initializes a brand new backend.
func initializeStorage(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage) (string, [][]byte) {

	var conf = vault.CoreConfig{
		Logger: logger.Named("initializeStorage"),
	}
	var opts = vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		//SkipInit:    true,
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		storage.Cleanup(t, cluster)
		cluster.Cleanup()
	}()

	//// Initialize leader
	//leader := cluster.Cores[0]
	//client := leader.Client
	//resp, err := client.Sys().Init(&api.InitRequest{
	//	SecretShares:    keyShares,
	//	SecretThreshold: keyThreshold,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//client.SetToken(resp.RootToken)
	//cluster.BarrierKeys = decodeKeys(t, resp.KeysB64)
	//// Unseal leader
	//cluster.UnsealCore(t, leader)

	rootToken := cluster.RootToken
	keys := cluster.BarrierKeys

	// Join raft cluster
	testhelpers.RaftClusterJoinNodes(t, cluster)
	time.Sleep(15 * time.Second)

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	testhelpers.WaitForNCoresUnsealed(t, cluster, vault.DefaultNumCores)

	return rootToken, keys
}

func decodeKeys(t *testing.T, keysB64 []string) [][]byte {
	keys := make([][]byte, len(keysB64))
	for i, k := range keysB64 {
		b, err := base64.RawStdEncoding.DecodeString(k)
		if err != nil {
			t.Fatal(err)
		}
		keys[i] = b
	}
	return keys
}

//////////////////////////////////////////////////////////////////////////////

//import (
//	"encoding/base64"
//	"testing"
//	"time"
//
//	"github.com/go-test/deep"
//
//	hclog "github.com/hashicorp/go-hclog"
//	"github.com/hashicorp/vault/api"
//	"github.com/hashicorp/vault/helper/testhelpers"
//	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
//	"github.com/hashicorp/vault/http"
//	"github.com/hashicorp/vault/sdk/helper/logging"
// 	"github.com/hashicorp/vault/vault"
// )
//
//const (
//	keyShares    = 5
//	keyThreshold = 3
//)
//
//func TestReusableStorage(t *testing.T) {
//
//	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())
//
//	t.Run("inmem", func(t *testing.T) {
//		t.Parallel()
//
//		logger := logger.Named("inmem")
//		storage, cleanup := teststorage.MakeReusableStorage(
//			t, logger, teststorage.MakeInmemBackend(t, logger))
//		defer cleanup()
//		testReusableStorage(t, logger, storage)
//	})
//
//	//t.Run("raft", func(t *testing.T) {
//	//	t.Parallel()
//
//	//	logger := logger.Named("raft")
//	//	storage, cleanup := teststorage.MakeReusableRaftStorage(t, logger)
//	//	defer cleanup()
//	//	testReusableStorage(t, logger, storage)
//	//})
//}
//
//func testReusableStorage(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage) {
//	//initializeStorage(t, logger, storage)
//	rootToken, keys := initializeStorage(t, logger, storage)
//	reuseStorage(t, logger, storage, rootToken, keys)
//}
//
//// initializeStorage initializes a brand new backend.
//func initializeStorage(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage) (string, [][]byte) {
//
//	var conf = vault.CoreConfig{
//		Logger: logger.Named("initializeStorage"),
//	}
//	var opts = vault.TestClusterOptions{
//		// TODO don't forget to handle BaseListenAddress correctly with
//		// parallelized tests.
//		BaseListenAddress: "127.0.0.1:50000",
//		HandlerFunc:       http.Handler,
//		SkipInit:          true,
//	}
//	storage.Setup(&conf, &opts)
//	cluster := vault.NewTestCluster(t, &conf, &opts)
//	cluster.Start()
//	defer func() {
//		storage.Cleanup(t, cluster)
//		cluster.Cleanup()
//	}()
//
//	leader := cluster.Cores[0]
//	client := leader.Client
//
//	// Initialize leader
//	resp, err := client.Sys().Init(&api.InitRequest{
//		SecretShares:    keyShares,
//		SecretThreshold: keyThreshold,
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	client.SetToken(resp.RootToken)
//
//	// Unseal
//	cluster.BarrierKeys = decodeKeys(t, resp.KeysB64)
//	if storage.IsRaft {
//
//		// Unseal leader
//		cluster.UnsealCore(t, leader)
//		time.Sleep(10 * time.Second)
//		//testhelpers.WaitForCoreUnsealed(t, leader)
//		//testhelpers.WaitForActiveNode(t, cluster)
//
//		// Join the followers to the raft cluster
//		for i := 1; i < vault.DefaultNumCores; i++ {
//			follower := cluster.Cores[i]
//			teststorage.JoinRaftFollower(t, cluster, leader, follower)
//
//			cluster.UnsealCore(t, follower)
//			//testhelpers.WaitForActiveNode(t, follower)
//			//testhelpers.WaitForCoreUnsealed(t, follower)
//		}
//		time.Sleep(10 * time.Second)
//	} else {
//		cluster.UnsealCores(t)
//	}
//	testhelpers.WaitForNCoresUnsealed(t, cluster, vault.DefaultNumCores)
//	//testhelpers.WaitForActiveNode(t, cluster)
//
//	// Mount kv
//	err = client.Sys().Mount("secret", &api.MountInput{
//		Type:    "kv",
//		Options: map[string]string{"version": "2"},
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Write a secret that we will read back out later.
//	_, err = client.Logical().Write(
//		"secret/foo",
//		map[string]interface{}{"zork": "quux"})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	cluster.EnsureCoresSealed(t)
//
//	return client.Token(), cluster.BarrierKeys
//}
//
//// reuseStorage re-uses a pre-populated backend.
//func reuseStorage(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, rootToken string, keys [][]byte) {
//
//	var conf = vault.CoreConfig{
//		Logger: logger.Named("reuseStorage"),
//	}
//	var opts = vault.TestClusterOptions{
//		BaseListenAddress: "127.0.0.1:50000",
//		HandlerFunc:       http.Handler,
//		SkipInit:          true,
//	}
//	storage.Setup(&conf, &opts)
//	cluster := vault.NewTestCluster(t, &conf, &opts)
//	cluster.Start()
//	defer func() {
//		storage.Cleanup(t, cluster)
//		cluster.Cleanup()
//	}()
//
//	leader := cluster.Cores[0]
//	client := leader.Client
//	client.SetToken(rootToken)
//
//	// Unseal
//	cluster.BarrierKeys = keys
//	cluster.UnsealCores(t)
//	testhelpers.WaitForNCoresUnsealed(t, cluster, vault.DefaultNumCores)
//
//	// Read the secret
//	secret, err := client.Logical().Read("secret/foo")
//	if err != nil {
//		t.Fatal(err)
//	}
//	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
//		t.Fatal(diff)
//	}
//
//	// Seal the cluster
//	cluster.EnsureCoresSealed(t)
//}
//
//func decodeKeys(t *testing.T, keysB64 []string) [][]byte {
//	keys := make([][]byte, len(keysB64))
//	for i, k := range keysB64 {
//		b, err := base64.RawStdEncoding.DecodeString(k)
//		if err != nil {
//			t.Fatal(err)
//		}
//		keys[i] = b
//	}
//	return keys
//}
