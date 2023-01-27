package sealmigration

import (
	"encoding/hex"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"

	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
)

type testFunc func(t *testing.T, logger hclog.Logger, storage teststorage.ReusableStorage, basePort int)

func testVariousBackends(t *testing.T, tf testFunc, basePort int, includeRaft bool) {
	logger := logging.NewVaultLogger(hclog.Trace).Named(t.Name())

	t.Run("inmem", func(t *testing.T) {
		t.Parallel()

		logger := logger.Named("inmem")
		storage, cleanup := teststorage.MakeReusableStorage(
			t, logger, teststorage.MakeInmemBackend(t, logger))
		defer cleanup()
		tf(t, logger, storage, basePort+100)
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

// TestSealMigration_ShamirToTransit_Post14 tests shamir-to-transit seal
// migration, using the post-1.4 method of bring individual nodes in the cluster
// to do the migration.
func TestSealMigration_ShamirToTransit_Post14(t *testing.T) {
	t.Parallel()
	testVariousBackends(t, ParamTestSealMigrationShamirToTransit_Post14, BasePort_ShamirToTransit_Post14, true)
}

// TestSealMigration_TransitToShamir_Post14 tests transit-to-shamir seal
// migration, using the post-1.4 method of bring individual nodes in the
// cluster to do the migration.
func TestSealMigration_TransitToShamir_Post14(t *testing.T) {
	t.Parallel()
	testVariousBackends(t, ParamTestSealMigrationTransitToShamir_Post14, BasePort_TransitToShamir_Post14, true)
}

// TestSealMigration_TransitToTransit tests transit-to-shamir seal
// migration, using the post-1.4 method of bring individual nodes in the
// cluster to do the migration.
func TestSealMigration_TransitToTransit(t *testing.T) {
	testVariousBackends(t, ParamTestSealMigration_TransitToTransit, BasePort_TransitToTransit, true)
}

func TestSysRekey_Update_ChangingRecoveryMode(t *testing.T) {
	cases := []struct {
		name             string
		recoveryModeInit bool
		recoveryModeAck  bool
		failureExpected  bool
	}{
		{"rekey-successful-unseal-recovery", true, true, false},
		{"rekey-fail-not-acked", true, false, true},
		{"rekey-fail-ack-unexpected", false, true, true},
		{"rekey-unseal-recovery-off", false, false, false},
	}

	// Create the transit server.
	tss := sealhelper.NewTransitSealServer(t, 0)
	defer func() {
		if tss != nil {
			tss.Cleanup()
		}
	}()
	sealKeyName := "transit-seal-key-1"
	tss.MakeKey(t, sealKeyName)
	logger := logging.NewVaultLogger(hclog.Trace)
	storage, cleanup := teststorage.MakeReusableStorage(
		t, logger, teststorage.MakeInmemBackend(t, logger))
	defer cleanup()
	// Initialize the backend with transit.
	cluster, _ := InitializeTransit(t, logger, storage, BasePort_TransitOnly, tss, sealKeyName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client
	keys := make([]string, len(cluster.RecoveryKeys))
	for i, k := range cluster.RecoveryKeys {
		keys[i] = hex.EncodeToString(k)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			initConfig := &api.RekeyInitRequest{
				SecretShares:         5,
				SecretThreshold:      3,
				EnableUnsealRecovery: c.recoveryModeInit,
			}
			resp, err := client.Sys().RekeyRecoveryKeyInit(initConfig)
			assert.NoError(t, err)

			for _, key := range keys[:3] {
				resp, err := client.Sys().RekeyRecoveryKeyUpdate(key, resp.Nonce, c.recoveryModeAck)

				if c.failureExpected {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					if resp.Complete {
						assert.True(t, resp.UnsealRecoveryEnabled == c.recoveryModeInit)
						keys = resp.Keys
					}
				}

			}
		})
	}
}
