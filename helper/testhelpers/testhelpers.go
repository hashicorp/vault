package testhelpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
	"github.com/hashicorp/vault/vault"
	testing "github.com/mitchellh/go-testing-interface"
)

type ReplicatedTestClusters struct {
	PerfPrimaryCluster     *vault.TestCluster
	PerfSecondaryCluster   *vault.TestCluster
	PerfPrimaryDRCluster   *vault.TestCluster
	PerfSecondaryDRCluster *vault.TestCluster
}

func (r *ReplicatedTestClusters) Cleanup() {
	r.PerfPrimaryCluster.Cleanup()
	r.PerfSecondaryCluster.Cleanup()
	if r.PerfPrimaryDRCluster != nil {
		r.PerfPrimaryDRCluster.Cleanup()
	}
	if r.PerfSecondaryDRCluster != nil {
		r.PerfSecondaryDRCluster.Cleanup()
	}
}

// Generates a root token on the target cluster.
func GenerateRoot(t testing.T, cluster *vault.TestCluster, drToken bool) string {
	token, err := GenerateRootWithError(t, cluster, drToken)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GenerateRootWithError(t testing.T, cluster *vault.TestCluster, drToken bool) (string, error) {
	// If recovery keys supported, use those to perform root token generation instead
	var keys [][]byte
	if cluster.Cores[0].SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}

	client := cluster.Cores[0].Client
	f := client.Sys().GenerateRootInit
	if drToken {
		f = client.Sys().GenerateDROperationTokenInit
	}
	status, err := f("", "")
	if err != nil {
		return "", err
	}

	if status.Required > len(keys) {
		return "", fmt.Errorf("need more keys than have, need %d have %d", status.Required, len(keys))
	}

	otp := status.OTP

	for i, key := range keys {
		if i >= status.Required {
			break
		}
		f := client.Sys().GenerateRootUpdate
		if drToken {
			f = client.Sys().GenerateDROperationTokenUpdate
		}
		status, err = f(base64.StdEncoding.EncodeToString(key), status.Nonce)
		if err != nil {
			return "", err
		}
	}
	if !status.Complete {
		return "", errors.New("generate root operation did not end successfully")
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(status.EncodedToken)
	if err != nil {
		return "", err
	}
	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

// RandomWithPrefix is used to generate a unique name with a prefix, for
// randomizing names in acceptance tests
func RandomWithPrefix(name string) string {
	return fmt.Sprintf("%s-%d", name, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func EnsureCoresSealed(t testing.T, c *vault.TestCluster) {
	t.Helper()
	for _, core := range c.Cores {
		EnsureCoreSealed(t, core)
	}
}

func EnsureCoreSealed(t testing.T, core *vault.TestClusterCore) error {
	core.Seal(t)
	timeout := time.Now().Add(60 * time.Second)
	for {
		if time.Now().After(timeout) {
			return fmt.Errorf("timeout waiting for core to seal")
		}
		if core.Core.Sealed() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	return nil
}

func EnsureCoresUnsealed(t testing.T, c *vault.TestCluster) {
	t.Helper()
	for _, core := range c.Cores {
		EnsureCoreUnsealed(t, c, core)
	}
}
func EnsureCoreUnsealed(t testing.T, c *vault.TestCluster, core *vault.TestClusterCore) {
	if !core.Sealed() {
		return
	}

	client := core.Client
	client.Sys().ResetUnsealProcess()
	for j := 0; j < len(c.BarrierKeys); j++ {
		statusResp, err := client.Sys().Unseal(base64.StdEncoding.EncodeToString(c.BarrierKeys[j]))
		if err != nil {
			// Sometimes when we get here it's already unsealed on its own
			// and then this fails for DR secondaries so check again
			if core.Sealed() {
				t.Fatal(err)
			}
			break
		}
		if statusResp == nil {
			t.Fatal("nil status response during unseal")
		}
		if !statusResp.Sealed {
			break
		}
	}
	if core.Sealed() {
		t.Fatal("core is still sealed")
	}
}

func EnsureCoreIsPerfStandby(t testing.T, core *vault.TestClusterCore) {
	t.Helper()
	start := time.Now()
	for {
		health, err := core.Client.Sys().Health()
		if err != nil {
			t.Fatal(err)
		}
		if health.PerformanceStandby {
			break
		}
		time.Sleep(time.Millisecond * 500)
		if time.Now().After(start.Add(time.Second * 60)) {
			t.Fatal("did not become a perf standby")
		}
	}
}

func WaitForReplicationState(t testing.T, c *vault.Core, state consts.ReplicationState) {
	timeout := time.Now().Add(10 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatalf("timeout waiting for core to have state %d", uint32(state))
		}
		state := c.ReplicationState()
		if state.HasState(state) {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func ConfClusterAndCore(t testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions) (*vault.TestCluster, *vault.TestClusterCore) {
	if conf.Physical != nil || conf.HAPhysical != nil {
		t.Fatalf("conf.Physical and conf.HAPhysical cannot be specified")
	}
	if opts.Logger == nil {
		t.Fatalf("opts.Logger must be specified")
	}

	inm, err := inmem.NewTransactionalInmem(nil, opts.Logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, opts.Logger)
	if err != nil {
		t.Fatal(err)
	}

	coreConfig := *conf
	coreConfig.Physical = inm
	coreConfig.HAPhysical = inmha.(physical.HABackend)

	cluster := vault.NewTestCluster(t, &coreConfig, opts)
	cluster.Start()

	cores := cluster.Cores
	core := cores[0]

	vault.TestWaitActive(t, core.Core)

	return cluster, core
}

func GetClusterAndCore(t testing.T, logger log.Logger, handlerFunc func(*vault.HandlerProperties) http.Handler) (*vault.TestCluster, *vault.TestClusterCore) {
	return ConfClusterAndCore(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		Logger:      logger,
		HandlerFunc: handlerFunc,
	})
}

func GetPerfReplicatedClusters(t testing.T, handlerFunc func(*vault.HandlerProperties) http.Handler) *ReplicatedTestClusters {
	ret := &ReplicatedTestClusters{}

	logger := log.New(&log.LoggerOptions{
		Mutex: &sync.Mutex{},
		Level: log.Trace,
	})
	// Set this lower so that state populates quickly to standby nodes
	vault.HeartbeatInterval = 2 * time.Second

	ret.PerfPrimaryCluster, _ = GetClusterAndCore(t, logger.Named("perf-pri"), handlerFunc)

	ret.PerfSecondaryCluster, _ = GetClusterAndCore(t, logger.Named("perf-sec"), handlerFunc)

	SetupTwoClusterPerfReplication(t, ret.PerfPrimaryCluster, ret.PerfSecondaryCluster)

	// Wait until poison pills have been read
	time.Sleep(45 * time.Second)
	EnsureCoresUnsealed(t, ret.PerfPrimaryCluster)
	EnsureCoresUnsealed(t, ret.PerfSecondaryCluster)

	return ret
}

func GetFourReplicatedClusters(t testing.T, handlerFunc func(*vault.HandlerProperties) http.Handler) *ReplicatedTestClusters {
	ret := &ReplicatedTestClusters{}

	logger := log.New(&log.LoggerOptions{
		Mutex: &sync.Mutex{},
		Level: log.Trace,
	})
	// Set this lower so that state populates quickly to standby nodes
	vault.HeartbeatInterval = 2 * time.Second

	ret.PerfPrimaryCluster, _ = GetClusterAndCore(t, logger.Named("perf-pri"), handlerFunc)

	ret.PerfSecondaryCluster, _ = GetClusterAndCore(t, logger.Named("perf-sec"), handlerFunc)

	ret.PerfPrimaryDRCluster, _ = GetClusterAndCore(t, logger.Named("perf-pri-dr"), handlerFunc)

	ret.PerfSecondaryDRCluster, _ = GetClusterAndCore(t, logger.Named("perf-sec-dr"), handlerFunc)

	SetupFourClusterReplication(t, ret.PerfPrimaryCluster, ret.PerfSecondaryCluster, ret.PerfPrimaryDRCluster, ret.PerfSecondaryDRCluster)

	// Wait until poison pills have been read
	time.Sleep(45 * time.Second)
	EnsureCoresUnsealed(t, ret.PerfPrimaryCluster)
	EnsureCoresUnsealed(t, ret.PerfSecondaryCluster)
	EnsureCoresUnsealed(t, ret.PerfPrimaryDRCluster)
	EnsureCoresUnsealed(t, ret.PerfSecondaryDRCluster)

	return ret
}

func SetupTwoClusterPerfReplication(t testing.T, perfPrimary, perfSecondary *vault.TestCluster) {
	// Enable performance primary
	_, err := perfPrimary.Cores[0].Client.Logical().Write("sys/replication/performance/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, perfPrimary.Cores[0].Core, consts.ReplicationPerformancePrimary)

	// get performance token
	secret, err := perfPrimary.Cores[0].Client.Logical().Write("sys/replication/performance/primary/secondary-token", map[string]interface{}{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	token := secret.WrapInfo.Token

	// enable performace secondary
	secret, err = perfSecondary.Cores[0].Client.Logical().Write("sys/replication/performance/secondary/enable", map[string]interface{}{
		"token":   token,
		"ca_file": perfPrimary.CACertPEMFile,
	})
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, perfSecondary.Cores[0].Core, consts.ReplicationPerformanceSecondary)
	time.Sleep(time.Second * 3)
	perfSecondary.BarrierKeys = perfPrimary.BarrierKeys

	EnsureCoresUnsealed(t, perfSecondary)
	rootToken := GenerateRoot(t, perfSecondary, false)
	perfSecondary.Cores[0].Client.SetToken(rootToken)
	for _, core := range perfSecondary.Cores {
		core.Client.SetToken(rootToken)
	}
}

func SetupFourClusterReplication(t testing.T, perfPrimary, perfSecondary, perfDRSecondary, perfSecondaryDRSecondary *vault.TestCluster) {
	var perfToken string
	var drToken string

	// Setup perf-primary
	{
		// Enable performance primary
		_, err := perfPrimary.Cores[0].Client.Logical().Write("sys/replication/primary/enable", nil)
		if err != nil {
			t.Fatal(err)
		}

		WaitForReplicationState(t, perfPrimary.Cores[0].Core, consts.ReplicationPerformancePrimary)
		// get performance token
		secret, err := perfPrimary.Cores[0].Client.Logical().Write("sys/replication/primary/secondary-token", map[string]interface{}{
			"id": "perf-secondary",
		})
		if err != nil {
			t.Fatal(err)
		}

		perfToken = secret.WrapInfo.Token

		// Enable dr primary
		_, err = perfPrimary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/enable", nil)
		if err != nil {
			t.Fatal(err)
		}

		WaitForReplicationState(t, perfPrimary.Cores[0].Core, consts.ReplicationDRPrimary)

		// get dr token
		secret, err = perfPrimary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/secondary-token", map[string]interface{}{
			"id": "primary-dr-secondary",
		})
		if err != nil {
			t.Fatal(err)
		}
		drToken = secret.WrapInfo.Token
		if drToken == "" {
			t.Fatal("empty token retrieved")
		}
	}

	WaitForActiveNode(t, perfPrimary)

	// Setup perf-secondary
	var perfSecondaryRootToken string
	var perfSecondaryDRToken string
	{
		// enable performace secondary
		_, err := perfSecondary.Cores[0].Client.Logical().Write("sys/replication/secondary/enable", map[string]interface{}{
			"token":   perfToken,
			"ca_file": perfPrimary.CACertPEMFile,
		})
		if err != nil {
			t.Fatal(err)
		}

		WaitForReplicationState(t, perfSecondary.Cores[0].Core, consts.ReplicationPerformanceSecondary)
		perfSecondary.BarrierKeys = perfPrimary.BarrierKeys

		// We want to make sure we unseal all the nodes so we first need to wait
		// until two of the nodes seal due to the poison pill being written
		WaitForNCoresSealed(t, perfSecondary, 2)
		EnsureCoresUnsealed(t, perfSecondary)
		perfSecondaryRootToken = GenerateRoot(t, perfSecondary, false)
		perfSecondary.Cores[0].Client.SetToken(perfSecondaryRootToken)

		// Enable dr primary on perf secondary
		_, err = perfSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/enable", nil)
		if err != nil {
			t.Fatal(err)
		}

		WaitForReplicationState(t, perfSecondary.Cores[0].Core, consts.ReplicationDRPrimary)

		// get dr token from perf secondary
		secret, err := perfSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/secondary-token", map[string]interface{}{
			"id": "secondary-dr-secondary",
		})
		if err != nil {
			t.Fatal(err)
		}

		perfSecondaryDRToken = secret.WrapInfo.Token
		if perfSecondaryDRToken == "" {
			t.Fatal("empty token retrieved")
		}
	}

	WaitForActiveNode(t, perfSecondary)
	// Setup pref-primary's dr secondary using "drToken"
	{
		// enable dr secondary
		_, err := perfDRSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/secondary/enable", map[string]interface{}{
			"token":   drToken,
			"ca_file": perfPrimary.CACertPEMFile,
		})
		if err != nil {
			t.Fatal(err)
		}

		WaitForReplicationState(t, perfDRSecondary.Cores[0].Core, consts.ReplicationDRSecondary)
		perfDRSecondary.BarrierKeys = perfPrimary.BarrierKeys

		// We want to make sure we unseal all the nodes so we first need to wait
		// until two of the nodes seal due to the poison pill being written
		WaitForNCoresSealed(t, perfDRSecondary, 2)
		EnsureCoresUnsealed(t, perfDRSecondary)

		perfDRSecondary.Cores[0].Client.SetToken(perfPrimary.Cores[0].Client.Token())
	}

	WaitForActiveNode(t, perfDRSecondary)
	time.Sleep(1 * time.Second)
	// Setup pref-primary's dr secondary using "perfSecondaryDRToken"
	{
		// enable dr secondary
		_, err := perfSecondaryDRSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/secondary/enable", map[string]interface{}{
			"token":   perfSecondaryDRToken,
			"ca_file": perfSecondary.CACertPEMFile,
		})
		if err != nil {
			t.Fatal(err)
		}

		WaitForReplicationState(t, perfSecondaryDRSecondary.Cores[0].Core, consts.ReplicationDRSecondary)
		perfSecondaryDRSecondary.BarrierKeys = perfPrimary.BarrierKeys

		// We want to make sure we unseal all the nodes so we first need to wait
		// until two of the nodes seal due to the poison pill being written
		WaitForNCoresSealed(t, perfSecondaryDRSecondary, 2)
		EnsureCoresUnsealed(t, perfSecondaryDRSecondary)

		perfSecondaryDRSecondary.Cores[0].Client.SetToken(perfSecondaryRootToken)
	}

	WaitForActiveNode(t, perfSecondaryDRSecondary)
}

func DeriveActiveCore(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	for i := 0; i < 10; i++ {
		for _, core := range cluster.Cores {
			leaderResp, err := core.Client.Sys().Leader()
			if err != nil {
				t.Fatal(err)
			}
			if leaderResp.IsSelf {
				return core
			}
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatal("could not derive the active core")
	return nil
}

func WaitForNCoresSealed(t testing.T, cluster *vault.TestCluster, n int) {
	for i := 0; i < 10; i++ {
		sealed := 0
		for _, core := range cluster.Cores {
			if core.Core.Sealed() {
				sealed++
			}
		}

		if sealed >= n {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("%d cores were not sealed", n)
}

func WaitForActiveNode(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	for i := 0; i < 10; i++ {
		for _, core := range cluster.Cores {
			if standby, _ := core.Core.Standby(); !standby {
				return core
			}
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not become active")
	return nil
}
