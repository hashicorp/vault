package testhelpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	realtesting "testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	raftlib "github.com/hashicorp/raft"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/builtin/credential/ldap"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	"github.com/hashicorp/vault/helper/xor"
	physConsul "github.com/hashicorp/vault/physical/consul"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	physFile "github.com/hashicorp/vault/sdk/physical/file"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/mitchellh/go-testing-interface"
)

type ReplicatedTestClusters struct {
	PerfPrimaryCluster     *vault.TestCluster
	PerfSecondaryCluster   *vault.TestCluster
	PerfPrimaryDRCluster   *vault.TestCluster
	PerfSecondaryDRCluster *vault.TestCluster
}

func (r *ReplicatedTestClusters) nonNilClusters() []*vault.TestCluster {
	all := []*vault.TestCluster{r.PerfPrimaryCluster, r.PerfSecondaryCluster,
		r.PerfPrimaryDRCluster, r.PerfSecondaryDRCluster}

	var ret []*vault.TestCluster
	for _, cluster := range all {
		if cluster != nil {
			ret = append(ret, cluster)
		}
	}
	return ret
}

func (r *ReplicatedTestClusters) Cleanup() {
	for _, cluster := range r.nonNilClusters() {
		cluster.Cleanup()
	}
}

func (r *ReplicatedTestClusters) Primary() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfPrimaryCluster, r.PerfPrimaryCluster.Cores[0], r.PerfPrimaryCluster.Cores[0].Client
}

func (r *ReplicatedTestClusters) Secondary() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfSecondaryCluster, r.PerfSecondaryCluster.Cores[0], r.PerfSecondaryCluster.Cores[0].Client
}

func (r *ReplicatedTestClusters) PrimaryDR() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfPrimaryDRCluster, r.PerfPrimaryDRCluster.Cores[0], r.PerfPrimaryDRCluster.Cores[0].Client
}

func (r *ReplicatedTestClusters) SecondaryDR() (*vault.TestCluster, *vault.TestClusterCore, *api.Client) {
	return r.PerfSecondaryDRCluster, r.PerfSecondaryDRCluster.Cores[0], r.PerfSecondaryDRCluster.Cores[0].Client
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

func EnsureCoreSealed(t testing.T, core *vault.TestClusterCore) {
	t.Helper()
	core.Seal(t)
	timeout := time.Now().Add(60 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatal("timeout waiting for core to seal")
		}
		if core.Core.Sealed() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
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

	core.SealAccess().ClearCaches(context.Background())
	if err := core.UnsealWithStoredKeys(context.Background()); err != nil {
		t.Fatal(err)
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

func EnsureCoreIsPerfStandby(t testing.T, client *api.Client) {
	t.Helper()
	logger := logging.NewVaultLogger(hclog.Info).Named(t.Name())
	start := time.Now()
	for {
		health, err := client.Sys().Health()
		if err != nil {
			t.Fatal(err)
		}
		if health.PerformanceStandby {
			break
		}

		logger.Info("waiting for performance standby", "health", health)
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

type PassthroughWithLocalPaths struct {
	logical.Backend
}

func (p *PassthroughWithLocalPaths) SpecialPaths() *logical.Paths {
	return &logical.Paths{
		LocalStorage: []string{
			"*",
		},
	}
}

func PassthroughWithLocalPathsFactory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b, err := vault.PassthroughBackendFactory(ctx, c)
	if err != nil {
		return nil, err
	}

	return &PassthroughWithLocalPaths{b}, nil
}

func ConfClusterAndCore(t testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions) (*vault.TestCluster, *vault.TestClusterCore) {
	cluster := MakeClusterAndStorage(t, InmemBackendSetup, conf, opts)
	return cluster, cluster.Cores[0]
}

func MakeClusterAndStorage(t testing.T, setup *BackendSetup, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *vault.TestCluster {
	{
		var coreConfig vault.CoreConfig
		if conf != nil {
			if conf.Physical != nil || conf.HAPhysical != nil {
				t.Fatalf("conf.Physical and conf.HAPhysical cannot be specified")
			}
			coreConfig = *conf
		}
		conf = &coreConfig
	}
	if setup == nil {
		setup = InmemBackendSetup
	}

	if setup.DisablePerformanceStandby {
		conf.DisablePerformanceStandby = true
	}
	conf.CredentialBackends = map[string]logical.Factory{
		"approle":  credAppRole.Factory,
		"userpass": credUserpass.Factory,
		"ldap":     ldap.Factory,
	}

	opts = getClusterDefaultsOpts(t, opts, "")
	if setup.KeepStandbysSealed {
		opts.KeepStandbysSealed = true
	}
	opts.SetupFunc = setup.Setup
	opts.CleanupFunc = setup.Cleanup

	switch {
	case setup.PhysicalFactory != nil && setup.CorePhysicalFactory == nil:
		conf.Physical, conf.HAPhysical, opts.CleanupFunc = setup.PhysicalFactory(t.(*realtesting.T), opts.Logger)
	case setup.PhysicalFactory == nil && setup.CorePhysicalFactory != nil:
		opts.PhysicalFactory = setup.CorePhysicalFactory
	default:
		t.Fatal("setup must provide exactly one of PhysicalFactory or CorePhysicalFactory")
	}

	vault.AddNoopAudit(conf)

	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	vault.TestWaitActive(t, cluster.Cores[0].Core)
	return cluster
}

// GetPerfReplicatedClusters returns a ReplicatedTestClusters containing both
// a perf primary and a pref secondary cluster, with replication enabled.
func GetPerfReplicatedClusters(t testing.T, setup SetupGetter, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	rc := PrepPerfReplicatedClusters(t, setup, conf, opts)
	rc.SetupTwoClusterPerfReplication(t, false)
	return rc
}

// getClusterDefaultsOpts returns a non-nil TestClusterOptions, based on opts
// if it is non-nil.  The Logger option will be populated.  If name is given,
// the logger will be created using the Named logger method, such that the string
// will appear as part of every log entry.
func getClusterDefaultsOpts(t testing.T, opts *vault.TestClusterOptions, name string) *vault.TestClusterOptions {
	if opts == nil {
		opts = &vault.TestClusterOptions{}
	}

	localOpts := *opts
	opts = &localOpts

	if opts.Logger == nil {
		opts.Logger = logging.NewVaultLogger(hclog.Trace).Named(t.Name())
	}
	if name != "" {
		opts.Logger = opts.Logger.Named(name)
	}

	return opts
}

// GetPerfPrimaryCluster returns a ReplicatedTestClusters containing only a
// single cluster.  Normally you would use NewTestCluster directly, but this
// helper may make sense if you want to test cluster replication but first do
// something with a standalone cluster.
func GetPerfPrimaryCluster(t testing.T, setup SetupGetter, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	if setup == nil {
		setup = InmemBackendSetup
	}
	opts = getClusterDefaultsOpts(t, opts, "")
	ret := &ReplicatedTestClusters{}

	// Set this lower so that state populates quickly to standby nodes
	cluster.HeartbeatInterval = 2 * time.Second

	ret.PerfPrimaryCluster = MakeClusterAndStorage(t, setup.Get(), conf, getClusterDefaultsOpts(t, opts, "perf-pri"))
	return ret
}

// AddPerfSecondaryCluster spins up a Perf Secondary cluster and adds it to
// the receiver.  Replication is not enabled.
func (r *ReplicatedTestClusters) AddPerfSecondaryCluster(t testing.T, setup SetupGetter, conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	if r.PerfSecondaryCluster != nil {
		t.Fatal("adding a perf secondary cluster when one is already present")
	}
	if setup == nil {
		setup = InmemBackendSetup
	}
	opts = getClusterDefaultsOpts(t, opts, "perf-sec")
	opts.FirstCoreNumber += len(r.PerfPrimaryCluster.Cores)
	r.PerfSecondaryCluster = MakeClusterAndStorage(t, setup.Get(), conf, opts)
}

// PrepPerfReplicatedClusters returns a ReplicatedTestClusters containing both
// a perf primary and a pref secondary cluster.  Replication is not enabled.
func PrepPerfReplicatedClusters(t testing.T, setup SetupGetter, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	ret := GetPerfPrimaryCluster(t, setup, conf, opts)
	ret.AddPerfSecondaryCluster(t, setup, conf, opts)
	return ret
}

// GetFourReplicatedClusters returns an inmem ReplicatedTestClusters with all
// clusters populated and replication enabled.
func GetFourReplicatedClusters(t testing.T, handlerFunc func(*vault.HandlerProperties) http.Handler) *ReplicatedTestClusters {
	return GetFourReplicatedClustersWithConf(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: handlerFunc,
	})
}

// GetFourReplicatedClustersWithConf returns a ReplicatedTestClusters with all
// clusters populated and replication enabled.
func GetFourReplicatedClustersWithConf(t testing.T, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	ret := &ReplicatedTestClusters{}

	opts = getClusterDefaultsOpts(t, opts, "")
	// Set this lower so that state populates quickly to standby nodes
	cluster.HeartbeatInterval = 2 * time.Second

	localopts := *opts
	localopts.Logger = opts.Logger.Named("perf-pri")
	ret.PerfPrimaryCluster = MakeClusterAndStorage(t, InmemBackendSetup, conf, &localopts)

	localopts.Logger = opts.Logger.Named("perf-sec")
	localopts.FirstCoreNumber += len(ret.PerfPrimaryCluster.Cores)
	ret.PerfSecondaryCluster = MakeClusterAndStorage(t, InmemBackendSetup, conf, &localopts)

	localopts.Logger = opts.Logger.Named("perf-pri-dr")
	localopts.FirstCoreNumber += len(ret.PerfSecondaryCluster.Cores)
	ret.PerfPrimaryDRCluster = MakeClusterAndStorage(t, InmemBackendSetup, conf, &localopts)

	localopts.Logger = opts.Logger.Named("perf-sec-dr")
	localopts.FirstCoreNumber += len(ret.PerfPrimaryDRCluster.Cores)
	ret.PerfSecondaryDRCluster = MakeClusterAndStorage(t, InmemBackendSetup, conf, &localopts)

	SetupFourClusterReplication(t, ret.PerfPrimaryCluster, ret.PerfSecondaryCluster, ret.PerfPrimaryDRCluster, ret.PerfSecondaryDRCluster)

	return ret
}

func (r *ReplicatedTestClusters) SetupTwoClusterPerfReplication(t testing.T, maskSecondaryToken bool) {
	SetupTwoClusterPerfReplication(t, r.PerfPrimaryCluster, r.PerfSecondaryCluster, maskSecondaryToken)
}

func SetupTwoClusterPerfReplication(t testing.T, pri, sec *vault.TestCluster, maskSecondaryToken bool) {
	EnablePerfPrimary(t, pri)

	var publicKey string
	if maskSecondaryToken {
		publicKey = generatePublicKey(t, sec)
	}
	perfToken := GetPerformanceToken(t, pri, sec.ID, publicKey)

	EnablePerformanceSecondary(t, perfToken, pri, sec, false, false)
}

func GetDRReplicatedClusters(t testing.T, setup SetupGetter, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	clusters := PrepDRReplicatedClusters(t, setup, conf, opts)
	SetupTwoClusterDRReplication(t, clusters.PerfPrimaryCluster, clusters.PerfPrimaryDRCluster, false)
	return clusters
}

func (r *ReplicatedTestClusters) AddDRSecondaryCluster(t testing.T, setup SetupGetter, conf *vault.CoreConfig, opts *vault.TestClusterOptions) {
	if setup == nil {
		setup = InmemBackendSetup
	}
	opts = getClusterDefaultsOpts(t, opts, "perf-dr-pri")
	opts.FirstCoreNumber += len(r.PerfPrimaryCluster.Cores)
	r.PerfPrimaryDRCluster = MakeClusterAndStorage(t, setup.Get(), conf, opts)
}

func PrepDRReplicatedClusters(t testing.T, setup SetupGetter, conf *vault.CoreConfig, opts *vault.TestClusterOptions) *ReplicatedTestClusters {
	ret := GetPerfPrimaryCluster(t, setup, conf, opts)
	ret.AddDRSecondaryCluster(t, setup, conf, opts)
	return ret
}

func SetupTwoClusterDRReplication(t testing.T, pri, sec *vault.TestCluster, maskSecondaryToken bool) {
	EnableDrPrimary(t, pri)
	setupDRReplication(t, pri, sec, maskSecondaryToken)
}

func setupDRReplication(t testing.T, pri, sec *vault.TestCluster, maskSecondaryToken bool) {
	var publicKey string
	if maskSecondaryToken {
		publicKey = generatePublicKey(t, sec)
	}
	drToken := getDrToken(t, pri, sec.ID, publicKey)

	EnableDrSecondary(t, pri, sec, drToken)
	for _, core := range sec.Cores {
		core.Client.SetToken(pri.Cores[0].Client.Token())
	}
	WaitForActiveNode(t, sec)
	WaitForMatchingMerkleRoots(t, "sys/replication/dr/", pri.Cores[0].Client, sec.Cores[0].Client)
	WaitForDRReplicationWorking(t, pri, sec)
}

func SetupFourClusterReplication(t testing.T, pri, sec, pridr, secdr *vault.TestCluster) {
	SetupTwoClusterPerfReplication(t, pri, sec, false)
	SetupTwoClusterDRReplication(t, pri, pridr, false)
	SetupTwoClusterDRReplication(t, sec, secdr, false)
}

func EnablePerfPrimary(t testing.T, cluster *vault.TestCluster) {
	cluster.Logger.Info("enabling performance primary")
	c := cluster.Cores[0]
	_, err := c.Client.Logical().Write("sys/replication/performance/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, c.Core, consts.ReplicationPerformancePrimary)
	WaitForActiveNodeAndPerfStandbys(t, cluster)
	cluster.Logger.Info("enabled performance primary")
}

func generatePublicKey(t testing.T, cluster *vault.TestCluster) string {
	generateKeyPath := "sys/replication/performance/secondary/generate-public-key"
	secret, err := cluster.Cores[0].Client.Logical().Write(generateKeyPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data == nil {
		t.Fatal("secret or secret data is nil")
	}

	return secret.Data["secondary_public_key"].(string)
}

func GetPerformanceToken(t testing.T, pri *vault.TestCluster, id, secondaryPublicKey string) string {
	client := pri.Cores[0].Client
	req := map[string]interface{}{
		"id": id,
	}
	if secondaryPublicKey != "" {
		req["secondary_public_key"] = secondaryPublicKey
	}
	secret, err := client.Logical().Write("sys/replication/performance/primary/secondary-token", req)
	if err != nil {
		t.Fatal(err)
	}
	return secret.WrapInfo.Token
}

func EnableDrPrimary(t testing.T, tc *vault.TestCluster) {
	tc.Logger.Info("enabling dr primary")
	c := tc.Cores[0]
	_, err := c.Client.Logical().Write("sys/replication/dr/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationStatus(t, c.Client, true, func(secret map[string]interface{}) bool {
		return secret["mode"] != nil && secret["mode"] == "primary"
	})
	tc.Logger.Info("enabled dr primary")
}

func getDrToken(t testing.T, tc *vault.TestCluster, id, secondaryPublicKey string) string {
	req := map[string]interface{}{
		"id": id,
	}
	if secondaryPublicKey != "" {
		req["secondary_public_key"] = secondaryPublicKey
	}
	secret, err := tc.Cores[0].Client.Logical().Write("sys/replication/dr/primary/secondary-token", req)
	if err != nil {
		t.Fatal(err)
	}
	return secret.WrapInfo.Token
}

func EnablePerformanceSecondary(t testing.T, perfToken string, pri, sec *vault.TestCluster, updatePrimary, skipPoisonPill bool) string {
	postData := map[string]interface{}{
		"token":   perfToken,
		"ca_file": pri.CACertPEMFile,
	}
	if pri.ClientAuthRequired {
		p := pri.Cores[0]
		postData["client_cert_pem"] = string(p.ServerCertPEM)
		postData["client_key_pem"] = string(p.ServerKeyPEM)
	}
	path := "sys/replication/performance/secondary/enable"
	if updatePrimary {
		path = "sys/replication/performance/secondary/update-primary"
	}
	_, err := sec.Cores[0].Client.Logical().Write(path, postData)
	if err != nil {
		t.Fatal(err)
	}

	sec.Logger.Info("enabled perf secondary, waiting for its replication state")
	WaitForReplicationState(t, sec.Cores[0].Core, consts.ReplicationPerformanceSecondary)

	var perfSecondaryRootToken string
	if !updatePrimary {
		sec.BarrierKeys = pri.BarrierKeys
		if !pri.Cores[0].SealAccess().RecoveryKeySupported() {
			sec.RecoveryKeys = pri.BarrierKeys
		} else {
			sec.RecoveryKeys = pri.RecoveryKeys
		}

		if len(sec.Cores) > 1 {
			if skipPoisonPill {
				// As part of prepareSecondary on the active node the keyring is
				// deleted from storage.  Its absence can cause standbys to seal
				// themselves. But it's not reliable, so we'll seal them
				// ourselves to force the issue.
				for _, core := range sec.Cores[1:] {
					EnsureCoreSealed(t, core)
				}
			} else {
				sec.Logger.Info("waiting for perf secondary standbys to seal")
				// We want to make sure we unseal all the nodes so we first need to wait
				// until two of the nodes seal due to the poison pill being written
				WaitForNCoresSealed(t, sec, len(sec.Cores)-1)
			}
		}
		sec.Logger.Info("waiting for perf secondary standbys to be unsealed")
		EnsureCoresUnsealed(t, sec)
		sec.Logger.Info("waiting for perf secondary active node")
		WaitForActiveNode(t, sec)
		sec.Logger.Info("generating new perf secondary root")

		perfSecondaryRootToken = GenerateRoot(t, sec, false)
		for _, core := range sec.Cores {
			core.Client.SetToken(perfSecondaryRootToken)
		}
		WaitForActiveNodeAndPerfStandbys(t, sec)
	}

	WaitForPerfReplicationWorking(t, pri, sec)
	return perfSecondaryRootToken
}

func EnableDrSecondary(t testing.T, pri, sec *vault.TestCluster, token string) {
	sec.Logger.Info("enabling dr secondary")
	_, err := sec.Cores[0].Client.Logical().Write("sys/replication/dr/secondary/enable", map[string]interface{}{
		"token":   token,
		"ca_file": pri.CACertPEMFile,
	})
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, sec.Cores[0].Core, consts.ReplicationDRSecondary)
	sec.BarrierKeys = pri.BarrierKeys

	// We want to make sure we unseal all the nodes so we first need to wait
	// until two of the nodes seal due to the poison pill being written
	WaitForNCoresSealed(t, sec, len(sec.Cores)-1)
	EnsureCoresUnsealed(t, sec)
	WaitForReplicationStatus(t, sec.Cores[0].Client, true, func(secret map[string]interface{}) bool {
		return secret["mode"] != nil && secret["mode"] == "secondary"
	})
}

func EnsureStableActiveNode(t testing.T, cluster *vault.TestCluster) {
	activeCore := DeriveActiveCore(t, cluster)

	for i := 0; i < 30; i++ {
		leaderResp, err := activeCore.Client.Sys().Leader()
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			t.Fatal("unstable active node")
		}
		time.Sleep(200 * time.Millisecond)
	}
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

func DeriveStandbyCores(t testing.T, cluster *vault.TestCluster) []*vault.TestClusterCore {
	cores := make([]*vault.TestClusterCore, 0, 2)
	for _, core := range cluster.Cores {
		leaderResp, err := core.Client.Sys().Leader()
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			cores = append(cores, core)
		}
	}

	return cores
}

func WaitForNCoresUnsealed(t testing.T, cluster *vault.TestCluster, n int) {
	t.Helper()
	for i := 0; i < 30; i++ {
		unsealed := 0
		for _, core := range cluster.Cores {
			if !core.Core.Sealed() {
				unsealed++
			}
		}

		if unsealed >= n {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("%d cores were not sealed", n)
}

func WaitForNCoresSealed(t testing.T, cluster *vault.TestCluster, n int) {
	t.Helper()
	for i := 0; i < 30; i++ {
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

func WaitForActiveNodeAndPerfStandbys(t testing.T, cluster *vault.TestCluster) {
	t.Helper()
	mountPoint, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = cluster.Cores[0].Client.Sys().Mount(mountPoint, &api.MountInput{
		Type:  "kv",
		Local: true,
	})
	if err != nil {
		t.Fatal("unable to mount KV engine")
	}
	path := mountPoint + "/foo"
	var standbys, actives int64
	var wg sync.WaitGroup
	deadline := time.Now().Add(30 * time.Second)
	for _, c := range cluster.Cores {
		wg.Add(1)
		go func(client *api.Client) {
			defer wg.Done()
			val := 1
			for time.Now().Before(deadline) {
				_, err = cluster.Cores[0].Client.Logical().Write(path, map[string]interface{}{
					"bar": val,
				})
				if err != nil {
					t.Fatal("unable to write KV", "path", path)
				}
				val++
				time.Sleep(250 * time.Millisecond)
				leader, err := client.Sys().Leader()
				if err != nil {
					if strings.Contains(err.Error(), "Vault is sealed") {
						continue
					}
					t.Fatal(err)
				}
				if leader.IsSelf {
					atomic.AddInt64(&actives, 1)
					return
				}
				if leader.PerfStandby && leader.PerfStandbyLastRemoteWAL > 0 {
					atomic.AddInt64(&standbys, 1)
					return
				}
			}
		}(c.Client)
	}
	wg.Wait()
	if actives != 1 || int(standbys) != len(cluster.Cores)-1 {
		t.Fatalf("expected 1 active core and %d standbys, got %d active and %d standbys",
			len(cluster.Cores)-1, actives, standbys)
	}
	err = cluster.Cores[0].Client.Sys().Unmount(mountPoint)
	if err != nil {
		t.Fatal("unable to unmount KV engine on primary")
	}
}

func WaitForActiveNode(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	for i := 0; i < 30; i++ {
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

func WaitForMatchingMerkleRoots(t testing.T, endpoint string, primary, secondary *api.Client) {
	getRoot := func(mode string, cli *api.Client) string {
		status, err := cli.Logical().Read(endpoint + "status")
		if err != nil {
			t.Fatal(err)
		}
		if status == nil || status.Data == nil || status.Data["mode"] == nil {
			t.Fatal("got nil secret or data")
		}
		if status.Data["mode"].(string) != mode {
			t.Fatalf("expected mode=%s, got %s", mode, status.Data["mode"].(string))
		}
		return status.Data["merkle_root"].(string)
	}

	t.Helper()
	for i := 0; i < 30; i++ {
		secRoot := getRoot("secondary", secondary)
		priRoot := getRoot("primary", primary)

		if reflect.DeepEqual(priRoot, secRoot) {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("roots did not become equal")
}

func WaitForMatchingMerkleRootsCore(t testing.T, pri, sec *vault.TestClusterCore, dr bool) {
	rootFunc := vault.PerformanceMerkleRoot
	if dr {
		rootFunc = vault.DRMerkleRoot
	}

	t.Helper()
	for i := 0; i < 30; i++ {
		secRoot := rootFunc(pri.Core)
		priRoot := rootFunc(pri.Core)

		if reflect.DeepEqual(priRoot, secRoot) {
			return
		}
		time.Sleep(time.Second)
	}

	t.Fatalf("roots did not become equal")
}

func WaitForReplicationStatus(t testing.T, client *api.Client, dr bool, accept func(map[string]interface{}) bool) {
	t.Helper()
	url := "sys/replication/performance/status"
	if dr {
		url = "sys/replication/dr/status"
	}

	var err error
	var secret *api.Secret
	for i := 0; i < 30; i++ {
		secret, err = client.Logical().Read(url)
		if err == nil && secret != nil && secret.Data != nil {
			if accept(secret.Data) {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("unable to get acceptable replication status: error=%v secret=%#v", err, secret)
}

func WaitForWAL(t testing.T, c *vault.TestClusterCore, wal uint64) {
	t.Helper()
	timeout := time.Now().Add(3 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatal("timeout waiting for WAL", "segment", wal, "lastrmtwal", vault.LastRemoteWAL(c.Core))
		}
		if vault.LastRemoteWAL(c.Core) >= wal {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func RekeyCluster(t testing.T, cluster *vault.TestCluster) {
	client := cluster.Cores[0].Client

	init, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	var statusResp *api.RekeyUpdateResponse
	for j := 0; j < len(cluster.BarrierKeys); j++ {
		statusResp, err = client.Sys().RekeyUpdate(base64.StdEncoding.EncodeToString(cluster.BarrierKeys[j]), init.Nonce)
		if err != nil {
			t.Fatal(err)
		}
		if statusResp == nil {
			t.Fatal("nil status response during unseal")
		}
		if statusResp.Complete {
			break
		}
	}

	if len(statusResp.KeysB64) != 5 {
		t.Fatal("wrong number of keys")
	}

	newBarrierKeys := make([][]byte, 5)
	for i, key := range statusResp.KeysB64 {
		newBarrierKeys[i], err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	cluster.BarrierKeys = newBarrierKeys
}

func CreateRaftBackend(t testing.T, logger hclog.Logger, nodeID string) (physical.Backend, func(), error) {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)
	cleanupFunc := func() {
		os.RemoveAll(raftDir)
	}

	logger.Info("raft dir", "dir", raftDir)

	conf := map[string]string{
		"path":                   raftDir,
		"node_id":                nodeID,
		"performance_multiplier": "8",
	}

	backend, err := raft.NewRaftBackend(conf, logger)
	if err != nil {
		cleanupFunc()
		t.Fatal(err)
	}

	return backend, cleanupFunc, nil
}

type TestRaftServerAddressProvider struct {
	Cluster *vault.TestCluster
}

func (p *TestRaftServerAddressProvider) ServerAddr(id raftlib.ServerID) (raftlib.ServerAddress, error) {
	for _, core := range p.Cluster.Cores {
		if core.NodeID == string(id) {
			parsed, err := url.Parse(core.ClusterAddr())
			if err != nil {
				return "", err
			}

			return raftlib.ServerAddress(parsed.Host), nil
		}
	}

	return "", errors.New("could not find cluster addr")
}

func RaftClusterJoinNodes(t testing.T, cluster *vault.TestCluster) {
	addressProvider := &TestRaftServerAddressProvider{Cluster: cluster}

	leaderCore := cluster.Cores[0]
	leaderAPI := leaderCore.Client.Address()
	atomic.StoreUint32(&vault.UpdateClusterAddrForTests, 1)

	// Seal the leader so we can install an address provider
	{
		EnsureCoreSealed(t, leaderCore)
		leaderCore.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		cluster.UnsealCore(t, leaderCore)
		vault.TestWaitActive(t, leaderCore.Core)
	}

	// Join core1
	{
		core := cluster.Cores[1]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderAPI, leaderCore.TLSConfig, false)
		if err != nil {
			t.Fatal(err)
		}

		cluster.UnsealCore(t, core)
	}

	// Join core2
	{
		core := cluster.Cores[2]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderAPI, leaderCore.TLSConfig, false)
		if err != nil {
			t.Fatal(err)
		}

		cluster.UnsealCore(t, core)
	}

	WaitForNCoresUnsealed(t, cluster, 3)
}

// WaitForPerfReplicationWorking mounts a KV non-locally, writes to it on pri, and waits for the value to be readable on sec.
func WaitForPerfReplicationWorking(t testing.T, pri, sec *vault.TestCluster) {
	t.Helper()
	priClient, secClient := pri.Cores[0].Client, sec.Cores[0].Client
	mountPoint, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = priClient.Sys().Mount(mountPoint, &api.MountInput{
		Type:  "kv",
		Local: false,
	})
	if err != nil {
		t.Fatal("unable to mount KV engine on primary")
	}

	path := mountPoint + "/foo"
	_, err = priClient.Logical().Write(path, map[string]interface{}{
		"bar": 1,
	})
	if err != nil {
		t.Fatal("unable to write KV on primary", "path", path)
	}

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		var secret *api.Secret
		secret, err = secClient.Logical().Read(path)
		if err == nil && secret != nil {
			err = priClient.Sys().Unmount(mountPoint)
			if err != nil {
				t.Fatal("unable to unmount KV engine on primary")
			}
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("unable to read replicated KV on secondary", "path", path, "err", err)
}

func WaitForDRReplicationWorking(t testing.T, pri, sec *vault.TestCluster) {
	priClient, secClient := pri.Cores[0].Client, sec.Cores[0].Client
	mountPoint, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = priClient.Sys().Mount(mountPoint, &api.MountInput{
		Type:  "kv",
		Local: false,
	})
	if err != nil {
		t.Fatal("unable to mount KV engine on primary")
	}

	path := mountPoint + "/foo"
	_, err = priClient.Logical().Write(path, map[string]interface{}{
		"bar": 1,
	})
	if err != nil {
		t.Fatal("unable to write KV on primary", "path", path)
	}

	WaitForReplicationStatus(t, secClient, true, func(secret map[string]interface{}) bool {
		if secret["last_remote_wal"] != nil {
			lastRemoteWal, _ := secret["last_remote_wal"].(json.Number).Int64()
			return lastRemoteWal > 0
		}

		return false
	})

	err = priClient.Sys().Unmount(mountPoint)
	if err != nil {
		t.Fatal("unable to unmount KV engine on primary")
	}
}

func MakeInmemBackend(t *realtesting.T, logger hclog.Logger) (physical.Backend, physical.HABackend, CleanupFunc) {
	inm, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	return inm, inmha.(physical.HABackend), func() {}
}

func MakeInmemNonTransactionalBackend(t *realtesting.T, logger hclog.Logger) (physical.Backend, physical.HABackend, CleanupFunc) {
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	return inm, inmha.(physical.HABackend), func() {}
}

func MakeFileBackend(t *realtesting.T, logger hclog.Logger) (physical.Backend, physical.HABackend, CleanupFunc) {
	path, err := ioutil.TempDir("", "vault-integ-file-")
	if err != nil {
		t.Fatal(err)
	}
	fileConf := map[string]string{
		"path": path,
	}
	fileBackend, err := physFile.NewTransactionalFileBackend(fileConf, logger)
	if err != nil {
		t.Fatal(err)
	}

	return fileBackend, nil, func() {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func MakeConsulBackend(t *realtesting.T, logger hclog.Logger) (physical.Backend, physical.HABackend, CleanupFunc) {
	cleanup, consulAddress, consulToken := consul.PrepareTestContainer(t, "1.4.0-rc1")
	consulConf := map[string]string{
		"address":      consulAddress,
		"token":        consulToken,
		"max_parallel": "32",
	}
	consulBackend, err := physConsul.NewConsulBackend(consulConf, logger)
	if err != nil {
		t.Fatal(err)
	}
	return consulBackend, nil, cleanup
}

type CleanupFunc func()

type SetupFunc func(t testing.T, c *vault.TestCluster)

type SetupGetter interface {
	Get() *BackendSetup
}

func (b *BackendSetup) Get() *BackendSetup {
	return b
}

type BackendSetup struct {
	KeepStandbysSealed        bool
	DisablePerformanceStandby bool
	Cleanup                   CleanupFunc
	Setup                     SetupFunc
	PhysicalFactory           func(*realtesting.T, hclog.Logger) (physical.Backend, physical.HABackend, CleanupFunc)
	CorePhysicalFactory       func(hclog.Logger) (physical.Backend, error)
}

var InmemBackendSetup = &BackendSetup{PhysicalFactory: MakeInmemBackend}
var InmemNonTransactionalBackendSetup = &BackendSetup{PhysicalFactory: MakeInmemNonTransactionalBackend}
var FileBackendSetup = &BackendSetup{PhysicalFactory: MakeFileBackend}
var ConsulBackendSetup = &BackendSetup{PhysicalFactory: MakeConsulBackend}

type RaftSetupGetter struct {
	T testing.T
}

func (r *RaftSetupGetter) Get() *BackendSetup {
	var cleanupFuncs []func()
	var i int
	return &BackendSetup{
		KeepStandbysSealed:        true,
		DisablePerformanceStandby: true, // TODO: remove this later
		CorePhysicalFactory: func(logger hclog.Logger) (physical.Backend, error) {
			backend, cleanupFunc, err := CreateRaftBackend(r.T, logger.Named("raft"), fmt.Sprintf("core-%d", i))
			i++
			cleanupFuncs = append(cleanupFuncs, cleanupFunc)
			return backend, err
		},
		Cleanup: func() {
			for _, c := range cleanupFuncs {
				c()
			}
		},
		Setup: func(t testing.T, c *vault.TestCluster) {
			RaftClusterJoinNodes(t, c)
			time.Sleep(15 * time.Second)
		},
	}
}
