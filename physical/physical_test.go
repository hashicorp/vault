package physical

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
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

	initializeStorage(t, logger, storage)

	//rootToken, keys := initializeStorage(t, logger, storage)
	//fmt.Printf("=======================================================================================\n")
	//fmt.Printf("=======================================================================================\n")
	//fmt.Printf("=======================================================================================\n")
	//reuseStorage(t, logger, storage, rootToken, keys)
}

// initializeStorage initializes a brand new backend storage.
func initializeStorage(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage) (string, [][]byte) {

	var conf = vault.CoreConfig{
		Logger: logger.Named("initializeStorage"),
	}
	var opts = vault.TestClusterOptions{
		HandlerFunc:           vaulthttp.Handler,
		BaseListenAddress:     "127.0.0.1:50000",
		BaseClusterListenPort: 50100,
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

	// Join raft cluster
	testhelpers.RaftClusterJoinNodes(t, cluster)
	time.Sleep(15 * time.Second)
	verifyRaftConfiguration(t, leader)

	// Wait until unsealed
	vault.TestWaitActive(t, leader.Core)
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

// reuseStorage uses a pre-populated backend storage.
func reuseStorage(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, rootToken string, keys [][]byte) {

	var conf = vault.CoreConfig{
		Logger: logger.Named("reuseStorage"),
	}
	var opts = vault.TestClusterOptions{
		HandlerFunc:           vaulthttp.Handler,
		BaseListenAddress:     "127.0.0.1:50000",
		BaseClusterListenPort: 50100,
		SkipInit:              true,
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		storage.Cleanup(t, cluster)
		cluster.Cleanup()
	}()

	for i, c := range cluster.Cores {
		if !c.Core.Sealed() {
			t.Fatalf("core is not sealed %d", i)
		}
	}

	//leader := cluster.Cores[0]
	//client := leader.Client
	//client.SetToken(rootToken)

	// Set Raft address providers
	testhelpers.RaftClusterSetAddressProviders(t, cluster)

	// Unseal cores
	cluster.BarrierKeys = keys
	for _, core := range cluster.Cores {
		cluster.UnsealCore(t, core)
		verifyRaftConfiguration(t, core)
		vault.TestWaitActive(t, core.Core)
	}

	// Wait until unsealed
	testhelpers.WaitForNCoresUnsealed(t, cluster, vault.DefaultNumCores)
}

//func getRaftConfiguration(t *testing.T, client *api.Client) []*raft.RaftServer {
//
//	resp, err := client.Logical().Read("sys/storage/raft/configuration")
//	if err != nil {
//		t.Fatal(err)
//	}
//	raw := resp.Data["config"].(map[string]interface{})["servers"].([]interface{})
//
//	servers := []*raft.RaftServer{}
//	for _, r := range raw {
//		rs := r.(map[string]interface{})
//		servers = append(servers, &raft.RaftServer{
//			NodeID:          rs["node_id"].(string),
//			Address:         rs["address"].(string),
//			Leader:          rs["leader"].(bool),
//			ProtocolVersion: rs["protocol_version"].(string),
//			Voter:           rs["voter"].(bool),
//		})
//	}
//	return servers
//}

func printRaftConfiguration(servers []*raft.RaftServer) string {
	var b strings.Builder
	for _, server := range servers {
		fmt.Fprintf(&b, "{\n")
		fmt.Fprintf(&b, "	NodeID:          %q,\n", server.NodeID)
		fmt.Fprintf(&b, "	Address:         %q,\n", server.Address)
		fmt.Fprintf(&b, "	Leader:          %t,\n", server.Leader)
		fmt.Fprintf(&b, "	ProtocolVersion: %q,\n", server.ProtocolVersion)
		fmt.Fprintf(&b, "	Voter:           %t,\n", server.Voter)
		fmt.Fprintf(&b, "},\n")
	}
	return b.String()
}

func verifyRaftConfiguration(t *testing.T, core *vault.TestClusterCore) {

	//servers := getRaftConfiguration(t, client)

	backend := core.UnderlyingRawStorage.(*raft.RaftBackend)
	ctx := namespace.RootContext(context.Background())
	config, err := backend.GetConfiguration(ctx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("-----------------------------------------------------------------\n")
	fmt.Printf("%s\n", printRaftConfiguration(config.Servers))

	//resp, err := client.Logical().Read("sys/storage/raft/configuration")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//servers := resp.Data["config"].(map[string]interface{})["servers"].([]interface{})

	//actual := []config{}
	//for _, s := range servers {
	//	server := s.(map[string]interface{})
	//	actual = append(actual, config{
	//		nodeID:   server[NodeID].(string),
	//		isLeader: server[Leader].(bool),
	//	})
	//}

	//var expected = []raft.RaftServer{
	//	{
	//		NodeID:          "node1",
	//		Address:         "127.0.0.1:8201",
	//		Leader:          false,
	//		ProtocolVersion: "3",
	//		Voter:           true,
	//	},
	//	{
	//		NodeID:          "node2",
	//		Address:         "127.0.0.2:8201",
	//		Leader:          true,
	//		ProtocolVersion: "3",
	//		Voter:           true,
	//	},
	//	{
	//		NodeID:          "node3",
	//		Address:         "127.0.0.3:8201",
	//		Leader:          false,
	//		ProtocolVersion: "3",
	//		Voter:           true,
	//	},
	//}

	//if diff := deep.Equal(actual, expected); len(diff) > 0 {
	//	t.Fatal(diff)
	//}
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
