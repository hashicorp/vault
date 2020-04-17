package seal_migration

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	vaultseal "github.com/hashicorp/vault/vault/seal"
)

func TestSealMigration_TransitToShamir(t *testing.T) {
	t.Parallel()

	t.Run("inmem", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())
		inm, err := inmem.NewTransactionalInmemHA(nil, logger)
		if err != nil {
			t.Fatal(err)
		}
		testSealMigrationTransitToShamir(t, logger, inm)
	})

	////t.Run("file", func(t *testing.T) {
	////	t.Parallel()
	////	testSealMigrationTransitToShamir(t, teststorage.FileBackendSetup)
	////})

	////t.Run("consul", func(t *testing.T) {
	////	t.Parallel()
	////	testSealMigrationTransitToShamir(t, teststorage.ConsulBackendSetup)
	////})

	////t.Run("raft", func(t *testing.T) {
	////	t.Parallel()
	////	testSealMigrationTransitToShamir(t, teststorage.RaftBackendSetup)
	////})
}

func testSealMigrationTransitToShamir(t *testing.T, logger log.Logger, backend physical.Backend) {

	var testEntry = map[string]interface{}{"bar": "quux"}

	// Create the transit server.
	tss := sealhelper.NewTransitSealServer(t)
	defer func() {
		if tss != nil {
			tss.Cleanup()
		}
	}()

	// Create a transit seal
	tss.MakeKey(t, "key1")
	transitSeal := tss.MakeSeal(t, "key1")

	// Create a cluster that uses transit
	var rootToken string
	var recoveryKeys [][]byte
	{
		cluster := vault.NewTestCluster(t,
			&vault.CoreConfig{
				Physical: backend,
				Logger:   logger.Named("transit_cluster"),
				Seal:     transitSeal,
			},
			&vault.TestClusterOptions{
				HandlerFunc: http.Handler,
				NumCores:    1,
			})
		cluster.Start()
		defer cluster.Cleanup()

		// save the root token and recovery keys
		client := cluster.Cores[0].Client
		rootToken = client.Token()
		recoveryKeys = cluster.RecoveryKeys

		// Write a secret that we will read back out later.
		_, err := client.Logical().Write("secret/foo", testEntry)
		if err != nil {
			t.Fatal(err)
		}

		// Seal the cluster
		cluster.EnsureCoresSealed(t)
	}

	// Create a shamir seal
	shamirSeal := vault.NewDefaultSeal(&vaultseal.Access{
		Wrapper: aeadwrapper.NewWrapper(&wrapping.WrapperOptions{}),
	})

	// Create a cluster that migrates from transit to shamir
	{
		cluster := vault.NewTestCluster(t,
			&vault.CoreConfig{
				Physical: backend,
				Logger:   logger.Named("transit_to_shamir_cluster"),
				Seal:     shamirSeal,
				// Setting an UnwrapSeal puts us in migration mode. This is the
				// equivalent of doing the following in HCL:
				//
				//     seal "transit" {
				//       // ...
				//       disabled = "true"
				//     }
				//
				UnwrapSeal: transitSeal,
			},
			&vault.TestClusterOptions{
				HandlerFunc: http.Handler,
				NumCores:    1,
				SkipInit:    true,
			})
		cluster.Start()
		defer cluster.Cleanup()

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		// Unseal and migrate to Shamir. Although we're unsealing using the
		// recovery keys, this is still an autounseal; if we stopped the
		// transit server this would fail.
		var resp *api.SealStatusResponse
		var err error
		for _, key := range recoveryKeys {
			strKey := base64.RawStdEncoding.EncodeToString(key)

			resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: strKey})
			if err == nil {
				t.Fatal("expected error due to lack of migrate parameter")
			}

			resp, err = client.Sys().UnsealWithOptions(&api.UnsealOpts{Key: strKey, Migrate: true})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil || !resp.Sealed {
				break
			}
		}
		if resp == nil || resp.Sealed {
			t.Fatalf("expected unsealed state; got %#v", resp)
		}
		testhelpers.WaitForActiveNode(t, cluster)

		// Await migration to finish.  Sadly there is no callback, and the test
		// will fail later on if we don't do this.
		// TODO maybe try to read?
		time.Sleep(10 * time.Second)

		// Read our secret
		secret, err := client.Logical().Read("secret/foo")
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(secret.Data, testEntry); len(diff) > 0 {
			t.Fatal(diff)
		}

		// Seal the cluster
		cluster.EnsureCoresSealed(t)
	}

	// Nuke the transit server
	tss.EnsureCoresSealed(t)
	tss.Cleanup()
	tss = nil

	// Create a cluster that uses shamir
	{
		cluster := vault.NewTestCluster(t,
			&vault.CoreConfig{
				Physical: backend,
				Logger:   logger.Named("shamir_cluster"),
				Seal:     shamirSeal,
			},
			&vault.TestClusterOptions{
				HandlerFunc: http.Handler,
				NumCores:    1,
				SkipInit:    true,
			})
		cluster.Start()
		defer cluster.Cleanup()

		// Note that the recovery keys are now the barrier keys.
		cluster.BarrierKeys = recoveryKeys
		testhelpers.EnsureCoresUnsealed(t, cluster)

		client := cluster.Cores[0].Client
		client.SetToken(rootToken)

		// Read our secret
		secret, err := client.Logical().Read("secret/foo")
		if err != nil {
			t.Fatal(err)
		}
		if diff := deep.Equal(secret.Data, testEntry); len(diff) > 0 {
			t.Fatal(diff)
		}

		// Seal the cluster
		cluster.EnsureCoresSealed(t)
	}
}
