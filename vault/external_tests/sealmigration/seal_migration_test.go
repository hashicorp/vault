package sealmigration

import (
	"sync/atomic"
	"testing"

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

func TestSealMigration_TransitToShamir_Recovery(t *testing.T) {
	t.Parallel()
	testVariousBackends(t, ParamTestSealMigrationTransitToShamir_Recovery, BasePort_TransitToShamir_Recovery, true)
}

// TestSealMigration_TransitToTransit tests transit-to-shamir seal
// migration, using the post-1.4 method of bring individual nodes in the
// cluster to do the migration.
func TestSealMigration_TransitToTransit(t *testing.T) {
	testVariousBackends(t, ParamTestSealMigration_TransitToTransit, BasePort_TransitToTransit, true)
}
