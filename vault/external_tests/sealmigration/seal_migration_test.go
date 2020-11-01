package sealmigration

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers"
	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage/consul"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
)

const (
	numTestCores = 3
	keyShares    = 3
	keyThreshold = 3

	basePort_ShamirToTransit_Pre14  = 20000
	basePort_TransitToShamir_Pre14  = 21000
	basePort_ShamirToTransit_Post14 = 22000
	basePort_TransitToShamir_Post14 = 23000
	basePort_TransitToTransit       = 24000
)

type testFunc func(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int)

func testVariousBackends(t *testing.T, tf testFunc, basePort int, includeRaft bool) {

	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())

	t.Run("inmem", func(t *testing.T) {
		t.Parallel()

		logger := logger.Named("inmem")
		storage, cleanup := teststorage.MakeReusableStorage(
			t, logger, teststorage.MakeInmemBackend(t, logger))
		defer cleanup()
		tf(t, logger, storage, basePort+100)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		logger := logger.Named("file")
		storage, cleanup := teststorage.MakeReusableStorage(
			t, logger, teststorage.MakeFileBackend(t, logger))
		defer cleanup()
		tf(t, logger, storage, basePort+200)
	})

	t.Run("consul", func(t *testing.T) {
		t.Parallel()

		logger := logger.Named("consul")
		storage, cleanup := teststorage.MakeReusableStorage(
			t, logger, consul.MakeConsulBackend(t, logger))
		defer cleanup()
		tf(t, logger, storage, basePort+300)
	})

	if includeRaft {
		t.Run("raft", func(t *testing.T) {
			t.Parallel()

			logger := logger.Named("raft")
			raftBasePort := basePort + 400

			atomic.StoreUint32(&vault.TestingUpdateClusterAddr, 1)
			addressProvider := testhelpers.NewHardcodedServerAddressProvider(numTestCores, raftBasePort+10)

			storage, cleanup := teststorage.MakeReusableRaftStorage(t, logger, numTestCores, addressProvider)
			defer cleanup()
			tf(t, logger, storage, raftBasePort)
		})
	}
}

// TestSealMigration_ShamirToTransit_Pre14 tests shamir-to-transit seal
// migration, using the pre-1.4 method of bring down the whole cluster to do
// the migration.
func TestSealMigration_ShamirToTransit_Pre14(t *testing.T) {
	t.Parallel()
	// Note that we do not test integrated raft storage since this is
	// a pre-1.4 test.
	testVariousBackends(t, testSealMigrationShamirToTransit_Pre14, basePort_ShamirToTransit_Pre14, false)
}

func testSealMigrationShamirToTransit_Pre14(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int) {

	// Initialize the backend using shamir
	cluster, _ := initializeShamir(t, logger, storage, basePort)
	rootToken, barrierKeys := cluster.RootToken, cluster.BarrierKeys
	cluster.Cleanup()
	storage.Cleanup(t, cluster)

	// Create the transit server.
	tss := sealhelper.NewTransitSealServer(t, 0)
	defer func() {
		tss.EnsureCoresSealed(t)
		tss.Cleanup()
	}()
	tss.MakeKey(t, "transit-seal-key-1")

	// Migrate the backend from shamir to transit.  Note that the barrier keys
	// are now the recovery keys.
	sealFunc := migrateFromShamirToTransit_Pre14(t, logger, storage, basePort, tss, rootToken, barrierKeys)

	// Run the backend with transit.
	runAutoseal(t, logger, storage, basePort, rootToken, sealFunc)
}

func migrateFromShamirToTransit_Pre14(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int, tss *sealhelper.TransitSealServer, rootToken string, recoveryKeys [][]byte) func() vault.Seal {
	var baseClusterPort = basePort + 10

	var conf = vault.CoreConfig{}
	var opts = vault.TestClusterOptions{
		Logger:                logger.Named("migrateFromShamirToTransit"),
		HandlerFunc:           vaulthttp.Handler,
		NumCores:              numTestCores,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
		SkipInit:              true,
		// N.B. Providing a transit seal puts us in migration mode.
		SealFunc: func() vault.Seal {
			return tss.MakeSeal(t, "transit-seal-key")
		},
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		cluster.Cleanup()
		storage.Cleanup(t, cluster)
	}()

	leader := cluster.Cores[0]
	leader.Client.SetToken(rootToken)

	// Unseal and migrate to Transit.
	unsealMigrate(t, leader.Client, recoveryKeys, true)

	// Wait for migration to finish.
	awaitMigration(t, leader.Client)

	verifySealConfigTransit(t, leader)

	// Read the secrets
	secret, err := leader.Client.Logical().Read("kv-wrapped/foo")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}

	// Write a new secret
	_, err = leader.Client.Logical().Write("kv-wrapped/test", map[string]interface{}{
		"zork": "quux",
	})
	if err != nil {
		t.Fatal(err)
	}

	return opts.SealFunc
}

// TestSealMigration_ShamirToTransit_Post14 tests shamir-to-transit seal
// migration, using the post-1.4 method of bring individual nodes in the cluster
// to do the migration.
func TestSealMigration_ShamirToTransit_Post14(t *testing.T) {
	t.Parallel()
	testVariousBackends(t, testSealMigrationShamirToTransit_Post14, basePort_ShamirToTransit_Post14, true)
}

func testSealMigrationShamirToTransit_Post14(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int) {
	// Initialize the backend using shamir
	cluster, opts := initializeShamir(t, logger, storage, basePort)

	// Create the transit server.
	tss := sealhelper.NewTransitSealServer(t, 0)
	defer tss.Cleanup()
	sealKeyName := "transit-seal-key-1"
	tss.MakeKey(t, sealKeyName)

	// Migrate the backend from shamir to transit.
	opts.SealFunc = func() vault.Seal {
		return tss.MakeSeal(t, sealKeyName)
	}

	// Restart each follower with the new config, and migrate to Transit.
	// Note that the barrier keys are being used as recovery keys.
	leaderIdx := migratePost14(t, storage, cluster, opts, cluster.BarrierKeys)
	validateMigration(t, storage, cluster, leaderIdx, verifySealConfigTransit)
	cluster.Cleanup()
	storage.Cleanup(t, cluster)

	// Run the backend with transit.
	runAutoseal(t, logger, storage, basePort, cluster.RootToken, opts.SealFunc)
}

// TestSealMigration_TransitToShamir_Post14 tests transit-to-shamir seal
// migration, using the post-1.4 method of bring individual nodes in the
// cluster to do the migration.
func TestSealMigration_TransitToShamir_Post14(t *testing.T) {
	t.Parallel()
	testVariousBackends(t, testSealMigrationTransitToShamir_Post14, basePort_TransitToShamir_Post14, true)
}

func testSealMigrationTransitToShamir_Post14(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int) {
	// Create the transit server.
	tss := sealhelper.NewTransitSealServer(t, 0)
	defer func() {
		if tss != nil {
			tss.Cleanup()
		}
	}()
	sealKeyName := "transit-seal-key-1"
	tss.MakeKey(t, sealKeyName)

	// Initialize the backend with transit.
	cluster, opts := initializeTransit(t, logger, storage, basePort, tss, sealKeyName)
	rootToken, recoveryKeys := cluster.RootToken, cluster.RecoveryKeys

	// Migrate the backend from transit to shamir
	opts.UnwrapSealFunc = opts.SealFunc
	opts.SealFunc = func() vault.Seal { return nil }
	leaderIdx := migratePost14(t, storage, cluster, opts, cluster.RecoveryKeys)
	validateMigration(t, storage, cluster, leaderIdx, verifySealConfigShamir)

	cluster.Cleanup()
	storage.Cleanup(t, cluster)

	// Now that migration is done, we can nuke the transit server, since we
	// can unseal without it.
	tss.Cleanup()
	tss = nil

	// Run the backend with shamir.  Note that the recovery keys are now the
	// barrier keys.
	runShamir(t, logger, storage, basePort, rootToken, recoveryKeys)
}

func validateMigration(t *testing.T, storage teststorage.ReusableStorage,
	cluster *vault.TestCluster, leaderIdx int, f func(t *testing.T, core *vault.TestClusterCore)) {
	t.Helper()

	leader := cluster.Cores[leaderIdx]

	secret, err := leader.Client.Logical().Read("kv-wrapped/foo")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}

	var appliedIndex uint64
	if storage.IsRaft {
		appliedIndex = testhelpers.RaftAppliedIndex(leader)
	}

	for _, core := range cluster.Cores {
		if storage.IsRaft {
			testhelpers.WaitForRaftApply(t, core, appliedIndex)
		}

		f(t, core)
	}
}

// TestSealMigration_TransitToTransit tests transit-to-shamir seal
// migration, using the post-1.4 method of bring individual nodes in the
// cluster to do the migration.
func TestSealMigration_TransitToTransit(t *testing.T) {
	testVariousBackends(t, testSealMigration_TransitToTransit, basePort_TransitToTransit, true)
}

func testSealMigration_TransitToTransit(t *testing.T, logger hclog.Logger,
	storage teststorage.ReusableStorage, basePort int) {

	// Create the transit server.
	tss1 := sealhelper.NewTransitSealServer(t, 0)
	defer func() {
		if tss1 != nil {
			tss1.Cleanup()
		}
	}()
	sealKeyName := "transit-seal-key-1"
	tss1.MakeKey(t, sealKeyName)

	// Initialize the backend with transit.
	cluster, opts := initializeTransit(t, logger, storage, basePort, tss1, sealKeyName)
	rootToken := cluster.RootToken

	// Create the transit server.
	tss2 := sealhelper.NewTransitSealServer(t, 1)
	defer func() {
		tss2.Cleanup()
	}()
	tss2.MakeKey(t, "transit-seal-key-2")

	// Migrate the backend from transit to transit.
	opts.UnwrapSealFunc = opts.SealFunc
	opts.SealFunc = func() vault.Seal {
		return tss2.MakeSeal(t, "transit-seal-key-2")
	}
	leaderIdx := migratePost14(t, storage, cluster, opts, cluster.RecoveryKeys)
	validateMigration(t, storage, cluster, leaderIdx, verifySealConfigTransit)
	cluster.Cleanup()
	storage.Cleanup(t, cluster)

	// Now that migration is done, we can nuke the transit server, since we
	// can unseal without it.
	tss1.Cleanup()
	tss1 = nil

	// Run the backend with transit.
	runAutoseal(t, logger, storage, basePort, rootToken, opts.SealFunc)
}

func migratePost14(t *testing.T, storage teststorage.ReusableStorage, cluster *vault.TestCluster,
	opts *vault.TestClusterOptions, unsealKeys [][]byte) int {

	cluster.Logger = cluster.Logger.Named("migration")
	// Restart each follower with the new config, and migrate.
	for i := 1; i < len(cluster.Cores); i++ {
		cluster.StopCore(t, i)
		if storage.IsRaft {
			teststorage.CloseRaftStorage(t, cluster, i)
		}
		cluster.StartCore(t, i, opts)

		unsealMigrate(t, cluster.Cores[i].Client, unsealKeys, true)
	}
	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

	// Step down the active node which will kick off the migration on one of the
	// other nodes.
	err := cluster.Cores[0].Client.Sys().StepDown()
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the followers to establish a new leader
	var leaderIdx int
	for i := 0; i < 30; i++ {
		leaderIdx, err = testhelpers.AwaitLeader(t, cluster)
		if err != nil {
			t.Fatal(err)
		}
		if leaderIdx != 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if leaderIdx == 0 {
		t.Fatalf("Core 0 cannot be the leader right now")
	}
	leader := cluster.Cores[leaderIdx]

	// Wait for migration to occur on the leader
	awaitMigration(t, leader.Client)

	var appliedIndex uint64
	if storage.IsRaft {
		appliedIndex = testhelpers.RaftAppliedIndex(leader)
		testhelpers.WaitForRaftApply(t, cluster.Cores[0], appliedIndex)
	}

	// Bring down the leader
	cluster.StopCore(t, 0)
	if storage.IsRaft {
		teststorage.CloseRaftStorage(t, cluster, 0)
	}

	// Bring core 0 back up; we still have the seal migration config in place,
	// but now that migration has been performed we should be able to unseal
	// with the new seal and without using the `migrate` unseal option.
	cluster.StartCore(t, 0, opts)
	unseal(t, cluster.Cores[0].Client, unsealKeys)

	// Write a new secret
	_, err = leader.Client.Logical().Write("kv-wrapped/test", map[string]interface{}{
		"zork": "quux",
	})
	if err != nil {
		t.Fatal(err)
	}

	return leaderIdx
}

func unsealMigrate(t *testing.T, client *api.Client, keys [][]byte, transitServerAvailable bool) {
	t.Helper()
	if err := attemptUnseal(client, keys); err == nil {
		t.Fatal("expected error due to lack of migrate parameter")
	}
	if err := attemptUnsealMigrate(client, keys, transitServerAvailable); err != nil {
		t.Fatal(err)
	}
}

func attemptUnsealMigrate(client *api.Client, keys [][]byte, transitServerAvailable bool) error {
	for i, key := range keys {
		resp, err := client.Sys().UnsealWithOptions(&api.UnsealOpts{
			Key:     base64.StdEncoding.EncodeToString(key),
			Migrate: true,
		})

		if i < keyThreshold-1 {
			// Not enough keys have been provided yet.
			if err != nil {
				return err
			}
		} else {
			if transitServerAvailable {
				// The transit server is running.
				if err != nil {
					return err
				}
				if resp == nil || resp.Sealed {
					return fmt.Errorf("expected unsealed state; got %#v", resp)
				}
			} else {
				// The transit server is stopped.
				if err == nil {
					return fmt.Errorf("expected error due to transit server being stopped.")
				}
			}
			break
		}
	}
	return nil
}

// awaitMigration waits for migration to finish.
func awaitMigration(t *testing.T, client *api.Client) {
	timeout := time.Now().Add(60 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}

		resp, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if !resp.Migration {
			return
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("migration did not complete.")
}

func unseal(t *testing.T, client *api.Client, keys [][]byte) {
	t.Helper()
	if err := attemptUnseal(client, keys); err != nil {
		t.Fatal(err)
	}
}

func attemptUnseal(client *api.Client, keys [][]byte) error {
	for i, key := range keys {

		resp, err := client.Sys().UnsealWithOptions(&api.UnsealOpts{
			Key: base64.StdEncoding.EncodeToString(key),
		})
		if i < keyThreshold-1 {
			// Not enough keys have been provided yet.
			if err != nil {
				return err
			}
		} else {
			if err != nil {
				return err
			}
			if resp == nil || resp.Sealed {
				return fmt.Errorf("expected unsealed state; got %#v", resp)
			}
			break
		}
	}
	return nil
}

func verifySealConfigShamir(t *testing.T, core *vault.TestClusterCore) {
	t.Helper()
	b, r, err := core.PhysicalSealConfigs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	verifyBarrierConfig(t, b, wrapping.Shamir, keyShares, keyThreshold, 1)
	if r != nil {
		t.Fatal("should not have recovery config for shamir")
	}
}

func verifySealConfigTransit(t *testing.T, core *vault.TestClusterCore) {
	t.Helper()
	b, r, err := core.PhysicalSealConfigs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	verifyBarrierConfig(t, b, wrapping.Transit, 1, 1, 1)
	verifyBarrierConfig(t, r, wrapping.Shamir, keyShares, keyThreshold, 0)
}

// verifyBarrierConfig verifies that a barrier configuration is correct.
func verifyBarrierConfig(t *testing.T, cfg *vault.SealConfig, sealType string, shares, threshold, stored int) {
	t.Helper()
	if cfg.Type != sealType {
		t.Fatalf("bad seal config: %#v, expected type=%q", cfg, sealType)
	}
	if cfg.SecretShares != shares {
		t.Fatalf("bad seal config: %#v, expected SecretShares=%d", cfg, shares)
	}
	if cfg.SecretThreshold != threshold {
		t.Fatalf("bad seal config: %#v, expected SecretThreshold=%d", cfg, threshold)
	}
	if cfg.StoredShares != stored {
		t.Fatalf("bad seal config: %#v, expected StoredShares=%d", cfg, stored)
	}
}

// initializeShamir initializes a brand new backend storage with Shamir.
func initializeShamir(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int) (*vault.TestCluster, *vault.TestClusterOptions) {
	t.Helper()

	var baseClusterPort = basePort + 10

	// Start the cluster
	var conf = vault.CoreConfig{}
	var opts = vault.TestClusterOptions{
		Logger:                logger.Named("initializeShamir"),
		HandlerFunc:           vaulthttp.Handler,
		NumCores:              numTestCores,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()

	leader := cluster.Cores[0]
	client := leader.Client

	// Unseal
	if storage.IsRaft {
		joinRaftFollowers(t, cluster, false)
		if err := testhelpers.VerifyRaftConfiguration(leader, len(cluster.Cores)); err != nil {
			t.Fatal(err)
		}
	} else {
		cluster.UnsealCores(t)
	}
	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

	err := client.Sys().Mount("kv-wrapped", &api.MountInput{
		SealWrap: true,
		Type:     "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write a secret that we will read back out later.
	_, err = client.Logical().Write("kv-wrapped/foo", map[string]interface{}{
		"zork": "quux",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cluster, &opts
}

// runShamir uses a pre-populated backend storage with Shamir.
func runShamir(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int, rootToken string, barrierKeys [][]byte) {
	t.Helper()
	var baseClusterPort = basePort + 10

	// Start the cluster
	var conf = vault.CoreConfig{}
	var opts = vault.TestClusterOptions{
		Logger:                logger.Named("runShamir"),
		HandlerFunc:           vaulthttp.Handler,
		NumCores:              numTestCores,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
		SkipInit:              true,
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		cluster.Cleanup()
		storage.Cleanup(t, cluster)
	}()

	leader := cluster.Cores[0]

	// Unseal
	cluster.BarrierKeys = barrierKeys
	if storage.IsRaft {
		for _, core := range cluster.Cores {
			cluster.UnsealCore(t, core)
		}
		// This is apparently necessary for the raft cluster to get itself
		// situated.
		time.Sleep(15 * time.Second)
		if err := testhelpers.VerifyRaftConfiguration(leader, len(cluster.Cores)); err != nil {
			t.Fatal(err)
		}
	} else {
		cluster.UnsealCores(t)
	}
	testhelpers.WaitForNCoresUnsealed(t, cluster, len(cluster.Cores))

	// Ensure that we always use the leader's client for this read check
	leaderIdx, err := testhelpers.AwaitLeader(t, cluster)
	if err != nil {
		t.Fatal(err)
	}
	client := cluster.Cores[leaderIdx].Client
	client.SetToken(rootToken)

	// Read the secrets
	secret, err := client.Logical().Read("kv-wrapped/foo")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}
	secret, err = client.Logical().Read("kv-wrapped/test")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}
}

// initializeTransit initializes a brand new backend storage with Transit.
func initializeTransit(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int,
	tss *sealhelper.TransitSealServer, sealKeyName string) (*vault.TestCluster, *vault.TestClusterOptions) {
	t.Helper()

	var baseClusterPort = basePort + 10

	// Start the cluster
	var conf = vault.CoreConfig{}
	var opts = vault.TestClusterOptions{
		Logger:                logger.Named("initializeTransit"),
		HandlerFunc:           vaulthttp.Handler,
		NumCores:              numTestCores,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
		SealFunc: func() vault.Seal {
			return tss.MakeSeal(t, sealKeyName)
		},
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()

	leader := cluster.Cores[0]
	client := leader.Client

	// Join raft
	if storage.IsRaft {
		joinRaftFollowers(t, cluster, true)

		if err := testhelpers.VerifyRaftConfiguration(leader, len(cluster.Cores)); err != nil {
			t.Fatal(err)
		}
	}
	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)

	err := client.Sys().Mount("kv-wrapped", &api.MountInput{
		SealWrap: true,
		Type:     "kv",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write a secret that we will read back out later.
	_, err = client.Logical().Write("kv-wrapped/foo", map[string]interface{}{
		"zork": "quux",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cluster, &opts
}

func runAutoseal(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int, rootToken string, sealFunc func() vault.Seal) {

	var baseClusterPort = basePort + 10

	// Start the cluster
	var conf = vault.CoreConfig{}
	var opts = vault.TestClusterOptions{
		Logger:                logger.Named("runTransit"),
		HandlerFunc:           vaulthttp.Handler,
		NumCores:              numTestCores,
		BaseListenAddress:     fmt.Sprintf("127.0.0.1:%d", basePort),
		BaseClusterListenPort: baseClusterPort,
		SkipInit:              true,
		SealFunc:              sealFunc,
	}
	storage.Setup(&conf, &opts)
	cluster := vault.NewTestCluster(t, &conf, &opts)
	cluster.Start()
	defer func() {
		cluster.Cleanup()
		storage.Cleanup(t, cluster)
	}()

	for _, c := range cluster.Cores {
		c.Client.SetToken(rootToken)
	}

	// Unseal.  Even though we are using autounseal, we have to unseal
	// explicitly because we are using SkipInit.
	if storage.IsRaft {
		for _, core := range cluster.Cores {
			cluster.UnsealCoreWithStoredKeys(t, core)
		}
		// This is apparently necessary for the raft cluster to get itself
		// situated.
		time.Sleep(15 * time.Second)
		// We're taking the first core, but we're not assuming it's the leader here.
		if err := testhelpers.VerifyRaftConfiguration(cluster.Cores[0], len(cluster.Cores)); err != nil {
			t.Fatal(err)
		}
	} else {
		if err := cluster.UnsealCoresWithError(true); err != nil {
			t.Fatal(err)
		}
	}
	testhelpers.WaitForNCoresUnsealed(t, cluster, len(cluster.Cores))

	// Preceding code may have stepped down the leader, so we're not sure who it is
	// at this point.
	leaderCore := testhelpers.DeriveActiveCore(t, cluster)
	client := leaderCore.Client

	// Read the secrets
	secret, err := client.Logical().Read("kv-wrapped/foo")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}
	secret, err = client.Logical().Read("kv-wrapped/test")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil {
		t.Fatal("secret is nil")
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}
}

// joinRaftFollowers unseals the leader, and then joins-and-unseals the
// followers one at a time.  We assume that the ServerAddressProvider has
// already been installed on all the nodes.
func joinRaftFollowers(t *testing.T, cluster *vault.TestCluster, useStoredKeys bool) {
	leader := cluster.Cores[0]

	cluster.UnsealCore(t, leader)
	vault.TestWaitActive(t, leader.Core)

	leaderInfos := []*raft.LeaderJoinInfo{
		&raft.LeaderJoinInfo{
			LeaderAPIAddr: leader.Client.Address(),
			TLSConfig:     leader.TLSConfig,
		},
	}

	// Join followers
	for i := 1; i < len(cluster.Cores); i++ {
		core := cluster.Cores[i]
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
		if err != nil {
			t.Fatal(err)
		}

		if useStoredKeys {
			// For autounseal, the raft backend is not initialized right away
			// after the join.  We need to wait briefly before we can unseal.
			awaitUnsealWithStoredKeys(t, core)
		} else {
			cluster.UnsealCore(t, core)
		}
	}

	testhelpers.WaitForNCoresUnsealed(t, cluster, len(cluster.Cores))
}

func awaitUnsealWithStoredKeys(t *testing.T, core *vault.TestClusterCore) {
	timeout := time.Now().Add(30 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatal("raft join: timeout waiting for core to unseal")
		}
		// Its actually ok for an error to happen here the first couple of
		// times -- it means the raft join hasn't gotten around to initializing
		// the backend yet.
		err := core.UnsealWithStoredKeys(context.Background())
		if err == nil {
			return
		}
		core.Logger().Warn("raft join: failed to unseal core", "error", err)
		time.Sleep(time.Second)
	}
}
