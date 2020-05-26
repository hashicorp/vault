// +build !enterprise

package seal_migration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/helper/testhelpers"
	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

// TestSealMigration_TransitToShamir_Pre14 tests transit-to-shamir seal
// migration, using the pre-1.4 method of bring down the whole cluster to do
// the migration.
func TestSealMigration_TransitToShamir_Pre14(t *testing.T) {
	// Note that we do not test integrated raft storage since this is
	// a pre-1.4 test.
	testVariousBackends(t, testSealMigrationTransitToShamir_Pre14, false)
}

func testSealMigrationTransitToShamir_Pre14(
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

	// Initialize the backend with transit.
	rootToken, recoveryKeys, transitSeal := initializeTransit(t, logger, storage, basePort, tss)

	// Migrate the backend from transit to shamir
	migrateFromTransitToShamir_Pre14(t, logger, storage, basePort, tss, transitSeal, rootToken, recoveryKeys)

	// Now that migration is done, we can nuke the transit server, since we
	// can unseal without it.
	tss.Cleanup()
	tss = nil

	// Run the backend with shamir.  Note that the recovery keys are now the
	// barrier keys.
	runShamir(t, logger, storage, basePort, rootToken, recoveryKeys)
}

func migrateFromTransitToShamir_Pre14(
	t *testing.T, logger hclog.Logger,
	storage teststorage.ReusableStorage, basePort int,
	tss *sealhelper.TransitSealServer, transitSeal vault.Seal,
	rootToken string, recoveryKeys [][]byte) {

	var baseClusterPort = basePort + 10

	var conf = vault.CoreConfig{
		Logger: logger.Named("migrateFromTransitToShamir"),
		// N.B. Providing an UnwrapSeal puts us in migration mode. This is the
		// equivalent of doing the following in HCL:
		//     seal "transit" {
		//       // ...
		//       disabled = "true"
		//     }
		UnwrapSeal: transitSeal,
	}
	var opts = vault.TestClusterOptions{
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
		storage.Cleanup(t, cluster)
		cluster.Cleanup()
	}()

	leader := cluster.Cores[0]
	client := leader.Client
	client.SetToken(rootToken)

	// Attempt to unseal while the transit server is unreachable.  Although
	// we're unsealing using the recovery keys, this is still an
	// autounseal, so it should fail.
	tss.EnsureCoresSealed(t)
	unsealMigrate(t, client, recoveryKeys, false)
	tss.UnsealCores(t)
	testhelpers.WaitForActiveNode(t, tss.TestCluster)

	// Unseal and migrate to Shamir. Although we're unsealing using the
	// recovery keys, this is still an autounseal.
	unsealMigrate(t, client, recoveryKeys, true)
	testhelpers.WaitForActiveNode(t, cluster)

	// Wait for migration to finish.  Sadly there is no callback, and the
	// test will fail later on if we don't do this.
	time.Sleep(10 * time.Second)

	// Read the secret
	secret, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(secret.Data, map[string]interface{}{"zork": "quux"}); len(diff) > 0 {
		t.Fatal(diff)
	}

	// Make sure the seal configs were updated correctly.
	b, r, err := cluster.Cores[0].Core.PhysicalSealConfigs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	verifyBarrierConfig(t, b, wrapping.Shamir, keyShares, keyThreshold, 1)
	if r != nil {
		t.Fatalf("expected nil recovery config, got: %#v", r)
	}

	cluster.EnsureCoresSealed(t)
}
