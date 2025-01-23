package identity

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/activationflags"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type logFn struct {
	fn func(msg string, args []interface{})
}

// Accept implements hclog.SinkAdapter.
func (f *logFn) Accept(name string, level hclog.Level, msg string, args ...interface{}) {
	f.fn(msg, args)
}

// TestIdentityDeduplication_Invalidate tests invalidation of the
// `force_identity_deduplication` activation flag by checking for a log line and
// ensuring the conflictResolver is modified on all nodes of a performance
// replication topology.
func TestIdentityDeduplication_Invalidate(t *testing.T) {
	testCluster := func(t *testing.T) (*corehelpers.TestLogger, *api.Client, *api.Client, func()) {
		var activeState, perfStandbyState string
		l := corehelpers.NewTestLogger(t)
		inm, err := inmem.NewTransactionalInmem(nil, l)
		require.NoError(t, err)
		inmha, err := inmem.NewInmemHA(nil, l)
		require.NoError(t, err)
		coreConfig := &vault.CoreConfig{
			Logger:                    l,
			Physical:                  inm,
			HAPhysical:                inmha.(physical.HABackend),
			DisablePerformanceStandby: false,
		}
		cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
			NumCores:    2,
			HandlerFunc: vaulthttp.Handler,
		})
		cluster.Start()
		cores := cluster.Cores

		// Track the active client's writes to keep the perf standby in sync
		active := cores[0].Client.WithResponseCallbacks(api.RecordState(&activeState))

		// Compose a new perfstandby client which is consistent with the active
		perfStandby := cores[1].Client.WithResponseCallbacks(api.RecordState(&perfStandbyState))
		perfStandby = perfStandby.WithRequestCallbacks(api.RequireState(activeState))

		vault.TestWaitActive(t, cores[0].Core)
		vault.TestWaitPerfStandby(t, cores[1].Core)
		return l, active, perfStandby, cluster.Cleanup
	}

	featurePath := "sys/activation-flags/" + activationflags.IdentityDeduplication
	enablePath := fmt.Sprintf("%s/%s", featurePath, "activate")
	disablePath := fmt.Sprintf("%s/%s", featurePath, "deactivate")
	readPath := "sys/activation-flags"

	// Always invalid
	moveToInactiveState := func(t *testing.T, active *api.Client, timeout time.Duration) {
		t.Helper()
		resp, err := active.Logical().Write(disablePath, nil)
		require.Error(t, err)
		require.Nil(t, resp)
	}

	// Always valid
	moveToActiveState := func(t *testing.T, active, perfStandby *api.Client, timeout time.Duration) {
		t.Helper()
		resp, err := active.Logical().Write(enablePath, nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			resp, err := perfStandby.Logical().Read(readPath)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Contains(t, resp.Data["activated"], activationflags.IdentityDeduplication)
		}, timeout, 100*time.Millisecond, "flag should be enabled")
	}

	tsts := map[string]struct {
		toState         string
		validTransition bool
	}{
		"Invalid transition to inactive": {
			"inactive",
			false,
		},
		"Valid transition to active": {
			"active",
			true,
		},
		"Invalid transition from active to inactive": {
			"inactive",
			false,
		},
	}

	for tn, tst := range tsts {
		t.Run(tn, func(t *testing.T) {
			l, active, perfStandby, clusterCleanup := testCluster(t)
			defer clusterCleanup()

			var loggerState *logFn
			stateChan := make(chan struct{}, 1)
			// Setup a logger we can use to capture sys backend logs
			loggerState = &logFn{
				fn: func(msg string, _ []interface{}) {
					if strings.Contains(msg, "groups restored") {
						select {
						case stateChan <- struct{}{}:
						default:
						}
					}
				},
			}
			l.RegisterSink(loggerState)

			timeout := 1 * time.Second
			if tst.validTransition {
				timeout = 20 * time.Second
			}

			switch tst.toState {
			case "inactive":
				moveToInactiveState(t, active, timeout)
			case "active":
				moveToActiveState(t, active, perfStandby, timeout)
			default:
				// noop
			}

			l.DeregisterSink(loggerState)

			waitForInvalidate := func(t *testing.T, timeout time.Duration) {
				t.Helper()
				select {
				case <-stateChan:
					require.True(t, tst.validTransition)
				case <-time.After(timeout):
					require.False(t, tst.validTransition)
				}
			}

			waitForInvalidate(t, timeout)
		})
	}
}

type DupeType int

const (
	UnDefined DupeType = iota
	Entities
	EntityAliases
	LocalAliases
	Groups
)

func (d DupeType) String() string {
	switch d {
	case Entities:
		return "entities"
	case EntityAliases:
		return "entity-aliases"
	case LocalAliases:
		return "local-aliases"
	case Groups:
		return "groups"
	default:
		return "invalid"
	}
}

// TestIdentityDeduplication_BucketInvalidation tests that as a result of
// deduplication renaming logic being introduced, there was a possibility
// that as part of normal identity engine operations we could lose the renames
// as they were only done in memory and not persisted to storage. We started
// persisting to storage so this test validates that after renames and
// bucket invalidation, the renames remain and are consistent.
func TestIdentityDeduplication_BucketInvalidation(t *testing.T) {
	t.Parallel()

	// getActive returns a core and client for the current active node of the TestCluster
	getActive := func(cluster *vault.TestCluster) (*vault.TestClusterCore, *api.Client) {
		activeIdx, err := testcluster.WaitForActiveNode(context.Background(), cluster)
		if err != nil {
			t.Fatal(err)
		}

		activeCore := cluster.Cores[activeIdx]
		activeClient := activeCore.Client
		return activeCore, activeClient
	}

	testCases := []struct {
		name                string
		dupeType            DupeType
		storagePackerPrefix string
		testData            map[string]interface{}
	}{
		{
			name:                "renamed entities",
			dupeType:            Entities,
			storagePackerPrefix: storagepacker.StoragePackerBucketsPrefix,
			testData: map[string]interface{}{
				"name":      "test-dupe",
				"count":     2,
				"random_id": true,
			},
		},
		{
			name:                "renamed groups",
			dupeType:            Groups,
			storagePackerPrefix: "packer/group/buckets/",
			testData: map[string]interface{}{
				"name":  "test-dupe",
				"count": 2,
			},
		},
	}

	for _, tc := range testCases {
		dt := tc.dupeType.String()
		if dt == "invalid" {
			t.Fatal("invalid duplicate type")
		}

		var listEndpoint string

		switch tc.dupeType {
		case Entities:
			listEndpoint = "identity/entity/name"
		case Groups:
			listEndpoint = "identity/group/name"
		default:
			t.Fatal("invalid duplicate type")
		}

		coreConfig := &vault.CoreConfig{}
		opts := &vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
		}
		cluster := vault.NewTestCluster(t, coreConfig, opts)
		cluster.Start()
		activeCore, activeClient := getActive(cluster)

		resp, err := activeClient.Logical().Write("identity/duplicate/"+dt, tc.testData)
		if err != nil {
			cluster.Cleanup()
			t.Fatal("error creating duplicate "+dt, err)
		}

		// Enabling the activation flag will trigger a reload of the identity engine
		resp, err = activeClient.Logical().Write("sys/activation-flags/"+activationflags.IdentityDeduplication+"/activate", map[string]interface{}{})
		if err != nil {
			cluster.Cleanup()
			t.Fatal("error enabling activation flag", err)
		}
		require.Contains(t, resp.Data["activated"], activationflags.IdentityDeduplication)

		// Get list of dupe type
		resp, err = activeClient.Logical().List(listEndpoint)
		if err != nil {
			cluster.Cleanup()
			t.Fatal("error listing "+dt+" names", err)
		}

		preInvalidationKeys := resp.Data["keys"]

		// Invalidate storage buckets
		for i := 0; i < 256; i++ {
			bucketNum := strconv.Itoa(i)
			activeCore.IdentityStore().Invalidate(context.Background(), tc.storagePackerPrefix+bucketNum)
		}

		// Get list of duplicate type
		resp, err = activeClient.Logical().List(listEndpoint)
		if err != nil {
			cluster.Cleanup()
			t.Fatal("error listing "+dt+" names", err)
		}

		postInvalidationKeys := resp.Data["keys"]

		require.Equal(t, preInvalidationKeys, postInvalidationKeys)
		cluster.Cleanup()
	}
}

func TestIdentityDeduplication_PR(t *testing.T) {
	// Create clusters
	conf, opts := teststorage.ClusterSetup(nil, nil, nil)
	clusters := testhelpers.GetFourReplicatedClustersWithConf(t, conf, opts)
	defer clusters.Cleanup()

	// Set up clients for PR Primary and PR Secondary
	_, prPrimaryCore, prPrimaryClient := clusters.Primary()
	_, prSecondaryCore, prSecondaryClient := clusters.Secondary()

	// Check prPrimary replication status
	status, err := prPrimaryClient.Logical().Read("sys/replication/status")
	require.NoError(t, err)
	prStatus := status.Data["performance"].(map[string]interface{})
	require.Equal(t, "primary", prStatus["mode"])

	// Check prSecondary replication status
	status, err = prSecondaryClient.Logical().Read("sys/replication/status")
	require.NoError(t, err)
	prStatus = status.Data["performance"].(map[string]interface{})
	require.Equal(t, "secondary", prStatus["mode"])

	// Create Duplicate Entities
	resp, err := prPrimaryClient.Logical().Write("identity/duplicate/entities", map[string]interface{}{
		"name":      "test-dupe",
		"random_id": true,
		"count":     2,
	})
	require.NoError(t, err)
	t.Log("IDs written to Primary: ", resp.Data["ids"])

	// Wait for PR Primary and Secondary to match
	testhelpers.WaitForMatchingMerkleRoots(t, "sys/replication/performance/", prPrimaryClient, prSecondaryClient)

	resp, err = prSecondaryClient.Logical().List("identity/entity/name")
	require.NoError(t, err)
	secondaryList := resp.Data["keys"]
	t.Log("Secondary List Pre-Activation:", secondaryList)

	// Enabling the activation flag will trigger a reload of the identity engine
	resp, err = prPrimaryClient.Logical().Write("sys/activation-flags/"+activationflags.IdentityDeduplication+"/activate", map[string]interface{}{})
	require.NoError(t, err)
	require.Contains(t, resp.Data["activated"], activationflags.IdentityDeduplication)

	// Get list of entity names from PR Primary and Secondary and ensure they're equal
	resp, err = prPrimaryClient.Logical().List("identity/entity/name")
	require.NoError(t, err)
	primaryList := resp.Data["keys"]

	resp, err = prSecondaryClient.Logical().List("identity/entity/name")
	require.NoError(t, err)
	secondaryList = resp.Data["keys"]

	t.Log("Primary List:", primaryList, "Secondary List:", secondaryList)
	require.Equal(t, primaryList, secondaryList)

	// Invalidate Identity Storage Buckets on PR Primary
	for i := 0; i < 256; i++ {
		bucketNum := strconv.Itoa(i)
		prPrimaryCore.Core.IdentityStore().Invalidate(context.Background(), storagepacker.StoragePackerBucketsPrefix+bucketNum)
	}

	// Get list of entity names from PR Primary and Secondary and ensure they're equal
	resp, err = prPrimaryClient.Logical().List("identity/entity/name")
	require.NoError(t, err)
	primaryList = resp.Data["keys"]

	resp, err = prSecondaryClient.Logical().List("identity/entity/name")
	require.NoError(t, err)
	secondaryList = resp.Data["keys"]

	t.Log("Primary List:", primaryList, "Secondary List:", secondaryList)
	require.Equal(t, primaryList, secondaryList)

	// Invalidate Identity Storage Buckets on PR Secondary
	for i := 0; i < 256; i++ {
		bucketNum := strconv.Itoa(i)
		prSecondaryCore.Core.IdentityStore().Invalidate(context.Background(), storagepacker.StoragePackerBucketsPrefix+bucketNum)
	}

	// Get list of entity names from PR Primary and Secondary and ensure they're equal
	resp, err = prPrimaryClient.Logical().List("identity/entity/name")
	require.NoError(t, err)
	primaryList = resp.Data["keys"]

	resp, err = prSecondaryClient.Logical().List("identity/entity/name")
	require.NoError(t, err)
	secondaryList = resp.Data["keys"]

	t.Log("Primary List:", primaryList, "Secondary List:", secondaryList)
	require.Equal(t, primaryList, secondaryList)
}
