// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package testhelpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	raftlib "github.com/hashicorp/raft"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/xor"
	"github.com/hashicorp/vault/vault"
)

//go:generate enumer -type=GenerateRootKind -trimprefix=GenerateRoot
type GenerateRootKind int

const (
	GenerateRootRegular GenerateRootKind = iota
	GenerateRootDR
	GenerateRecovery
)

// GenerateRoot generates a root token on the target cluster.
func GenerateRoot(t testing.TB, cluster *vault.TestCluster, kind GenerateRootKind) string {
	t.Helper()
	token, err := GenerateRootWithError(t, cluster, kind)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GenerateRootWithError(t testing.TB, cluster *vault.TestCluster, kind GenerateRootKind) (string, error) {
	t.Helper()
	// If recovery keys supported, use those to perform root token generation instead
	var keys [][]byte
	if cluster.Cores[0].SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}
	client := cluster.Cores[0].Client
	oldNS := client.Namespace()
	defer client.SetNamespace(oldNS)
	client.ClearNamespace()

	var err error
	var status *api.GenerateRootStatusResponse
	switch kind {
	case GenerateRootRegular:
		status, err = client.Sys().GenerateRootInit("", "")
	case GenerateRootDR:
		status, err = client.Sys().GenerateDROperationTokenInit("", "")
	case GenerateRecovery:
		status, err = client.Sys().GenerateRecoveryOperationTokenInit("", "")
	}
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

		strKey := base64.StdEncoding.EncodeToString(key)
		switch kind {
		case GenerateRootRegular:
			status, err = client.Sys().GenerateRootUpdate(strKey, status.Nonce)
		case GenerateRootDR:
			status, err = client.Sys().GenerateDROperationTokenUpdate(strKey, status.Nonce)
		case GenerateRecovery:
			status, err = client.Sys().GenerateRecoveryOperationTokenUpdate(strKey, status.Nonce)
		}
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

func EnsureCoresSealed(t testing.TB, c *vault.TestCluster) {
	t.Helper()
	for _, core := range c.Cores {
		EnsureCoreSealed(t, core)
	}
}

func EnsureCoreSealed(t testing.TB, core *vault.TestClusterCore) {
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

func EnsureCoresUnsealed(t testing.TB, c *vault.TestCluster) {
	t.Helper()
	for i, core := range c.Cores {
		err := AttemptUnsealCore(c, core)
		if err != nil {
			t.Fatalf("failed to unseal core %d: %v", i, err)
		}
	}
}

func EnsureCoreUnsealed(t testing.TB, c *vault.TestCluster, core *vault.TestClusterCore) {
	t.Helper()
	err := AttemptUnsealCore(c, core)
	if err != nil {
		t.Fatalf("failed to unseal core: %v", err)
	}
}

func AttemptUnsealCores(c *vault.TestCluster) error {
	for i, core := range c.Cores {
		err := AttemptUnsealCore(c, core)
		if err != nil {
			return fmt.Errorf("failed to unseal core %d: %v", i, err)
		}
	}
	return nil
}

func AttemptUnsealCore(c *vault.TestCluster, core *vault.TestClusterCore) error {
	if !core.Sealed() {
		return nil
	}

	core.SealAccess().ClearCaches(context.Background())
	if err := core.UnsealWithStoredKeys(context.Background()); err != nil {
		return err
	}

	client := core.Client
	oldNS := client.Namespace()
	defer client.SetNamespace(oldNS)
	client.ClearNamespace()

	client.Sys().ResetUnsealProcess()
	for j := 0; j < len(c.BarrierKeys); j++ {
		statusResp, err := client.Sys().Unseal(base64.StdEncoding.EncodeToString(c.BarrierKeys[j]))
		if err != nil {
			// Sometimes when we get here it's already unsealed on its own
			// and then this fails for DR secondaries so check again
			if core.Sealed() {
				return err
			} else {
				return nil
			}
		}
		if statusResp == nil {
			return fmt.Errorf("nil status response during unseal")
		}
		if !statusResp.Sealed {
			break
		}
	}
	if core.Sealed() {
		return fmt.Errorf("core is still sealed")
	}
	return nil
}

func EnsureStableActiveNode(t testing.TB, cluster *vault.TestCluster) {
	t.Helper()
	deriveStableActiveCore(t, cluster)
}

func DeriveStableActiveCore(t testing.TB, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	return deriveStableActiveCore(t, cluster)
}

func deriveStableActiveCore(t testing.TB, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	activeCore := DeriveActiveCore(t, cluster)
	minDuration := time.NewTimer(3 * time.Second)

	for i := 0; i < 60; i++ {
		leaderResp, err := activeCore.Client.Sys().Leader()
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			minDuration.Reset(3 * time.Second)
		}
		time.Sleep(200 * time.Millisecond)
	}

	select {
	case <-minDuration.C:
	default:
		if stopped := minDuration.Stop(); stopped {
			t.Fatal("unstable active node")
		}
		// Drain the value
		<-minDuration.C
	}

	return activeCore
}

func DeriveActiveCore(t testing.TB, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	for i := 0; i < 60; i++ {
		for _, core := range cluster.Cores {
			oldNS := core.Client.Namespace()
			core.Client.ClearNamespace()
			leaderResp, err := core.Client.Sys().Leader()
			core.Client.SetNamespace(oldNS)
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

func DeriveStandbyCores(t testing.TB, cluster *vault.TestCluster) []*vault.TestClusterCore {
	t.Helper()
	cores := make([]*vault.TestClusterCore, 0, 2)
	for _, core := range cluster.Cores {
		oldNS := core.Client.Namespace()
		core.Client.ClearNamespace()
		leaderResp, err := core.Client.Sys().Leader()
		core.Client.SetNamespace(oldNS)
		if err != nil {
			t.Fatal(err)
		}
		if !leaderResp.IsSelf {
			cores = append(cores, core)
		}
	}

	return cores
}

func WaitForNCoresUnsealed(t testing.TB, cluster *vault.TestCluster, n int) {
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

	t.Fatalf("%d cores were not unsealed", n)
}

func SealCores(t testing.TB, cluster *vault.TestCluster) {
	t.Helper()
	for _, core := range cluster.Cores {
		if err := core.Shutdown(); err != nil {
			t.Fatal(err)
		}
		timeout := time.Now().Add(3 * time.Second)
		for {
			if time.Now().After(timeout) {
				t.Fatal("timeout waiting for core to seal")
			}
			if core.Sealed() {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func WaitForNCoresSealed(t testing.TB, cluster *vault.TestCluster, n int) {
	t.Helper()
	for i := 0; i < 60; i++ {
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

func WaitForActiveNode(t testing.TB, cluster *vault.TestCluster) *vault.TestClusterCore {
	t.Helper()
	for i := 0; i < 60; i++ {
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

func WaitForStandbyNode(t testing.TB, core *vault.TestClusterCore) {
	t.Helper()
	for i := 0; i < 30; i++ {
		if isLeader, _, clusterAddr, _ := core.Core.Leader(); isLeader != true && clusterAddr != "" {
			return
		}
		if core.Core.ActiveNodeReplicationState() == 0 {
			return
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not become standby")
}

func RekeyCluster(t testing.TB, cluster *vault.TestCluster, recovery bool) [][]byte {
	t.Helper()
	cluster.Logger.Info("rekeying cluster", "recovery", recovery)
	client := cluster.Cores[0].Client

	initFunc := client.Sys().RekeyInit
	if recovery {
		initFunc = client.Sys().RekeyRecoveryKeyInit
	}
	init, err := initFunc(&api.RekeyInitRequest{
		SecretShares:    5,
		SecretThreshold: 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	var statusResp *api.RekeyUpdateResponse
	keys := cluster.BarrierKeys
	if cluster.Cores[0].Core.SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	}

	updateFunc := client.Sys().RekeyUpdate
	if recovery {
		updateFunc = client.Sys().RekeyRecoveryKeyUpdate
	}
	for j := 0; j < len(keys); j++ {
		statusResp, err = updateFunc(base64.StdEncoding.EncodeToString(keys[j]), init.Nonce)
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
	cluster.Logger.Info("cluster rekeyed", "recovery", recovery)

	if cluster.Cores[0].Core.SealAccess().RecoveryKeySupported() && !recovery {
		return nil
	}
	if len(statusResp.KeysB64) != 5 {
		t.Fatal("wrong number of keys")
	}

	newKeys := make([][]byte, 5)
	for i, key := range statusResp.KeysB64 {
		newKeys[i], err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Fatal(err)
		}
	}
	return newKeys
}

// HardcodedServerAddressProvider is a ServerAddressProvider that uses
// a hardcoded map of raft node addresses.
//
// It is useful in cases where the raft configuration is known ahead of time,
// but some of the cores have not yet had startClusterListener() called (via
// either unsealing or raft joining), and thus do not yet have a ClusterAddr()
// assigned.
type HardcodedServerAddressProvider struct {
	Entries map[raftlib.ServerID]raftlib.ServerAddress
}

func (p *HardcodedServerAddressProvider) ServerAddr(id raftlib.ServerID) (raftlib.ServerAddress, error) {
	if addr, ok := p.Entries[id]; ok {
		return addr, nil
	}
	return "", errors.New("could not find cluster addr")
}

// NewHardcodedServerAddressProvider is a convenience function that makes a
// ServerAddressProvider from a given cluster address base port.
func NewHardcodedServerAddressProvider(numCores, baseClusterPort int) raftlib.ServerAddressProvider {
	entries := make(map[raftlib.ServerID]raftlib.ServerAddress)

	for i := 0; i < numCores; i++ {
		id := fmt.Sprintf("core-%d", i)
		addr := fmt.Sprintf("127.0.0.1:%d", baseClusterPort+i)
		entries[raftlib.ServerID(id)] = raftlib.ServerAddress(addr)
	}

	return &HardcodedServerAddressProvider{
		entries,
	}
}

// VerifyRaftConfiguration checks that we have a valid raft configuration, i.e.
// the correct number of servers, having the correct NodeIDs, and exactly one
// leader.
func VerifyRaftConfiguration(core *vault.TestClusterCore, numCores int) error {
	backend := core.UnderlyingRawStorage.(*raft.RaftBackend)
	ctx := namespace.RootContext(context.Background())
	config, err := backend.GetConfiguration(ctx)
	if err != nil {
		return err
	}

	servers := config.Servers
	if len(servers) != numCores {
		return fmt.Errorf("Found %d servers, not %d", len(servers), numCores)
	}

	leaders := 0
	for i, s := range servers {
		if s.NodeID != fmt.Sprintf("core-%d", i) {
			return fmt.Errorf("Found unexpected node ID %q", s.NodeID)
		}
		if s.Leader {
			leaders++
		}
	}

	if leaders != 1 {
		return fmt.Errorf("Found %d leaders", leaders)
	}

	return nil
}

func RaftAppliedIndex(core *vault.TestClusterCore) uint64 {
	return core.UnderlyingRawStorage.(*raft.RaftBackend).AppliedIndex()
}

func WaitForRaftApply(t testing.TB, core *vault.TestClusterCore, index uint64) {
	t.Helper()

	backend := core.UnderlyingRawStorage.(*raft.RaftBackend)
	for i := 0; i < 30; i++ {
		if backend.AppliedIndex() >= index {
			return
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not apply index")
}

// AwaitLeader waits for one of the cluster's nodes to become leader.
func AwaitLeader(t testing.TB, cluster *vault.TestCluster) (int, error) {
	timeout := time.Now().Add(60 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}

		for i, core := range cluster.Cores {
			if core.Core.Sealed() {
				continue
			}

			isLeader, _, _, _ := core.Leader()
			if isLeader {
				return i, nil
			}
		}

		time.Sleep(time.Second)
	}

	return 0, fmt.Errorf("timeout waiting leader")
}

func GenerateDebugLogs(t testing.TB, client *api.Client) chan struct{} {
	t.Helper()

	stopCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				err := client.Sys().Mount("foo", &api.MountInput{
					Type: "kv",
					Options: map[string]string{
						"version": "1",
					},
				})
				if err != nil {
					t.Fatal(err)
				}

				err = client.Sys().Unmount("foo")
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}()

	return stopCh
}

// VerifyRaftPeers verifies that the raft configuration contains a given set of peers.
// The `expected` contains a map of expected peers. Existing entries are deleted
// from the map by removing entries whose keys are in the raft configuration.
// Remaining entries result in an error return so that the caller can poll for
// an expected configuration.
func VerifyRaftPeers(t testing.TB, client *api.Client, expected map[string]bool) error {
	t.Helper()

	resp, err := client.Logical().Read("sys/storage/raft/configuration")
	if err != nil {
		t.Fatalf("error reading raft config: %v", err)
	}

	if resp == nil || resp.Data == nil {
		t.Fatal("missing response data")
	}

	config, ok := resp.Data["config"].(map[string]interface{})
	if !ok {
		t.Fatal("missing config in response data")
	}

	servers, ok := config["servers"].([]interface{})
	if !ok {
		t.Fatal("missing servers in response data config")
	}

	// Iterate through the servers and remove the node found in the response
	// from the expected collection
	for _, s := range servers {
		server := s.(map[string]interface{})
		delete(expected, server["node_id"].(string))
	}

	// If the collection is non-empty, it means that the peer was not found in
	// the response.
	if len(expected) != 0 {
		return fmt.Errorf("failed to read configuration successfully, expected peers not found in configuration list: %v", expected)
	}

	return nil
}

func TestMetricSinkProvider(gaugeInterval time.Duration) func(string) (*metricsutil.ClusterMetricSink, *metricsutil.MetricsHelper) {
	return func(clusterName string) (*metricsutil.ClusterMetricSink, *metricsutil.MetricsHelper) {
		inm := metrics.NewInmemSink(1000000*time.Hour, 2000000*time.Hour)
		clusterSink := metricsutil.NewClusterMetricSink(clusterName, inm)
		clusterSink.GaugeInterval = gaugeInterval
		return clusterSink, metricsutil.NewMetricsHelper(inm, false)
	}
}

func SysMetricsReq(client *api.Client, cluster *vault.TestCluster, unauth bool) (*SysMetricsJSON, error) {
	r := client.NewRequest("GET", "/v1/sys/metrics")
	if !unauth {
		r.Headers.Set("X-Vault-Token", cluster.RootToken)
	}
	var data SysMetricsJSON
	resp, err := client.RawRequestWithContext(context.Background(), r)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Response.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, errors.New("failed to unmarshal:" + err.Error())
	}
	return &data, nil
}

type SysMetricsJSON struct {
	Gauges   []gaugeJSON   `json:"Gauges"`
	Counters []counterJSON `json:"Counters"`

	// note: this is referred to as a "Summary" type in our telemetry docs, but
	// the field name in the JSON is "Samples"
	Summaries []summaryJSON `json:"Samples"`
}

type baseInfoJSON struct {
	Name   string                 `json:"Name"`
	Labels map[string]interface{} `json:"Labels"`
}

type gaugeJSON struct {
	baseInfoJSON
	Value int `json:"Value"`
}

type counterJSON struct {
	baseInfoJSON
	Count  int     `json:"Count"`
	Rate   float64 `json:"Rate"`
	Sum    int     `json:"Sum"`
	Min    int     `json:"Min"`
	Max    int     `json:"Max"`
	Mean   float64 `json:"Mean"`
	Stddev float64 `json:"Stddev"`
}

type summaryJSON struct {
	baseInfoJSON
	Count  int     `json:"Count"`
	Rate   float64 `json:"Rate"`
	Sum    float64 `json:"Sum"`
	Min    float64 `json:"Min"`
	Max    float64 `json:"Max"`
	Mean   float64 `json:"Mean"`
	Stddev float64 `json:"Stddev"`
}

// SetNonRootToken sets a token on :client: with a fairly generic policy.
// This is useful if a test needs to examine differing behavior based on if a
// root token is passed with the request.
func SetNonRootToken(client *api.Client) error {
	policy := `path "*" { capabilities = ["create", "update", "read"] }`
	if err := client.Sys().PutPolicy("policy", policy); err != nil {
		return fmt.Errorf("error putting policy: %v", err)
	}

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"policy"},
		TTL:      "30m",
	})
	if err != nil {
		return fmt.Errorf("error creating token secret: %v", err)
	}

	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
		return fmt.Errorf("missing token auth data")
	}

	client.SetToken(secret.Auth.ClientToken)
	return nil
}

// RetryUntilAtCadence runs f until it returns a nil result or the timeout is reached.
// If a nil result hasn't been obtained by timeout, calls t.Fatal.
func RetryUntilAtCadence(t testing.TB, timeout, sleepTime time.Duration, f func() error) {
	t.Helper()
	fail := func(err error) {
		t.Helper()
		t.Fatalf("did not complete before deadline, err: %v", err)
	}
	RetryUntilAtCadenceWithHandler(t, timeout, sleepTime, fail, f)
}

// RetryUntilAtCadenceWithHandler runs f until it returns a nil result or the timeout is reached.
// If a nil result hasn't been obtained by timeout, onFailure is called.
func RetryUntilAtCadenceWithHandler(t testing.TB, timeout, sleepTime time.Duration, onFailure func(error), f func() error) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var err error
	for time.Now().Before(deadline) {
		if err = f(); err == nil {
			return
		}
		time.Sleep(sleepTime)
	}
	onFailure(err)
}

// RetryUntil runs f with a 100ms pause between calls, until f returns a nil result
// or the timeout is reached.
// If a nil result hasn't been obtained by timeout, calls t.Fatal.
// NOTE: See RetryUntilAtCadence if you want to specify a different wait/sleep
// duration between calls.
func RetryUntil(t testing.TB, timeout time.Duration, f func() error) {
	t.Helper()
	RetryUntilAtCadence(t, timeout, 100*time.Millisecond, f)
}

// CreateEntityAndAlias clones an existing client and creates an entity/alias, uses userpass mount path
// It returns the cloned client, entityID, and aliasID.
func CreateEntityAndAlias(t testing.TB, client *api.Client, mountAccessor, entityName, aliasName string) (*api.Client, string, string) {
	return CreateEntityAndAliasWithinMount(t, client, mountAccessor, "userpass", entityName, aliasName)
}

// CreateEntityAndAliasWithinMount clones an existing client and creates an entity/alias, within the specified mountPath
// It returns the cloned client, entityID, and aliasID.
func CreateEntityAndAliasWithinMount(t testing.TB, client *api.Client, mountAccessor, mountPath, entityName, aliasName string) (*api.Client, string, string) {
	t.Helper()
	userClient, err := client.Clone()
	if err != nil {
		t.Fatalf("failed to clone the client:%v", err)
	}
	userClient.SetToken(client.Token())

	resp, err := client.Logical().WriteWithContext(context.Background(), "identity/entity", map[string]interface{}{
		"name": entityName,
	})
	if err != nil {
		t.Fatalf("failed to create an entity:%v", err)
	}
	entityID := resp.Data["id"].(string)

	aliasResp, err := client.Logical().WriteWithContext(context.Background(), "identity/entity-alias", map[string]interface{}{
		"name":           aliasName,
		"canonical_id":   entityID,
		"mount_accessor": mountAccessor,
	})
	if err != nil {
		t.Fatalf("failed to create an entity alias:%v", err)
	}
	aliasID := aliasResp.Data["id"].(string)
	if aliasID == "" {
		t.Fatal("Alias ID not present in response")
	}
	path := fmt.Sprintf("auth/%s/users/%s", mountPath, aliasName)
	_, err = client.Logical().WriteWithContext(context.Background(), path, map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("failed to configure userpass backend: %v", err)
	}

	return userClient, entityID, aliasID
}

// SetupTOTPMount enables the totp secrets engine by mounting it. This requires
// that the test cluster has a totp backend available.
func SetupTOTPMount(t testing.TB, client *api.Client) {
	t.Helper()
	// Mount the TOTP backend
	mountInfo := &api.MountInput{
		Type: "totp",
	}
	if err := client.Sys().Mount("totp", mountInfo); err != nil {
		t.Fatalf("failed to mount totp backend: %v", err)
	}
}

// SetupTOTPMethod configures the TOTP secrets engine with a provided config map.
func SetupTOTPMethod(t testing.TB, client *api.Client, config map[string]interface{}) string {
	t.Helper()

	resp1, err := client.Logical().Write("identity/mfa/method/totp", config)

	if err != nil || (resp1 == nil) {
		t.Fatalf("bad: resp: %#v\n err: %v", resp1, err)
	}

	methodID := resp1.Data["method_id"].(string)
	if methodID == "" {
		t.Fatalf("method ID is empty")
	}

	return methodID
}

// SetupMFALoginEnforcement configures a single enforcement method using the
// provided config map. "name" field is required in the config map.
func SetupMFALoginEnforcement(t testing.TB, client *api.Client, config map[string]interface{}) {
	t.Helper()
	enfName, ok := config["name"]
	if !ok {
		t.Fatalf("couldn't find name in login-enforcement config")
	}
	_, err := client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("identity/mfa/login-enforcement/%s", enfName), config)
	if err != nil {
		t.Fatalf("failed to configure MFAEnforcementConfig: %v", err)
	}
}

// SetupUserpassMountAccessor sets up userpass auth and returns its mount
// accessor. This requires that the test cluster has a "userpass" auth method
// available.
func SetupUserpassMountAccessor(t testing.TB, client *api.Client) string {
	t.Helper()
	// Enable Userpass authentication
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}

	auths, err := client.Sys().ListAuthWithContext(context.Background())
	if err != nil {
		t.Fatalf("failed to list auth methods: %v", err)
	}
	if auths == nil || auths["userpass/"] == nil {
		t.Fatalf("failed to get userpass mount accessor")
	}

	return auths["userpass/"].Accessor
}

// RegisterEntityInTOTPEngine registers an entity with a methodID and returns
// the generated name.
func RegisterEntityInTOTPEngine(t testing.TB, client *api.Client, entityID, methodID string) string {
	t.Helper()
	totpGenName := fmt.Sprintf("%s-%s", entityID, methodID)
	secret, err := client.Logical().WriteWithContext(context.Background(), "identity/mfa/method/totp/admin-generate", map[string]interface{}{
		"entity_id": entityID,
		"method_id": methodID,
	})
	if err != nil {
		t.Fatalf("failed to generate a TOTP secret on an entity: %v", err)
	}
	totpURL := secret.Data["url"].(string)
	if totpURL == "" {
		t.Fatalf("failed to get TOTP url in secret response: %+v", secret)
	}
	_, err = client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("totp/keys/%s", totpGenName), map[string]interface{}{
		"url": totpURL,
	})
	if err != nil {
		t.Fatalf("failed to register a TOTP URL: %v", err)
	}
	enfPath := fmt.Sprintf("identity/mfa/login-enforcement/%s", methodID[0:4])
	_, err = client.Logical().WriteWithContext(context.Background(), enfPath, map[string]interface{}{
		"name":                methodID[0:4],
		"identity_entity_ids": []string{entityID},
		"mfa_method_ids":      []string{methodID},
	})
	if err != nil {
		t.Fatalf("failed to create login enforcement")
	}

	return totpGenName
}

// GetTOTPCodeFromEngine requests a TOTP code from the specified enginePath.
func GetTOTPCodeFromEngine(t testing.TB, client *api.Client, enginePath string) string {
	t.Helper()
	totpPath := fmt.Sprintf("totp/code/%s", enginePath)
	secret, err := client.Logical().ReadWithContext(context.Background(), totpPath)
	if err != nil {
		t.Fatalf("failed to create totp passcode: %v", err)
	}
	if secret == nil || secret.Data == nil {
		t.Fatalf("bad secret returned from %s", totpPath)
	}
	return secret.Data["code"].(string)
}

// SetupLoginMFATOTP setups up a TOTP MFA using some basic configuration and
// returns all relevant information to the client.
func SetupLoginMFATOTP(t testing.TB, client *api.Client, methodName string, waitPeriod int) (*api.Client, string, string) {
	t.Helper()
	// Mount the totp secrets engine
	SetupTOTPMount(t, client)

	// Create a mount accessor to associate with an entity
	mountAccessor := SetupUserpassMountAccessor(t, client)

	// Create a test entity and alias
	entityClient, entityID, _ := CreateEntityAndAlias(t, client, mountAccessor, "entity1", "testuser1")

	// Configure a default TOTP method
	totpConfig := map[string]interface{}{
		"issuer":                  "yCorp",
		"period":                  waitPeriod,
		"algorithm":               "SHA256",
		"digits":                  6,
		"skew":                    1,
		"key_size":                20,
		"qr_size":                 200,
		"max_validation_attempts": 5,
		"method_name":             methodName,
	}
	methodID := SetupTOTPMethod(t, client, totpConfig)

	// Configure a default login enforcement
	enforcementConfig := map[string]interface{}{
		"auth_method_types": []string{"userpass"},
		"name":              methodID[0:4],
		"mfa_method_ids":    []string{methodID},
	}

	SetupMFALoginEnforcement(t, client, enforcementConfig)
	return entityClient, entityID, methodID
}

func SkipUnlessEnvVarsSet(t testing.TB, envVars []string) {
	t.Helper()

	for _, i := range envVars {
		if os.Getenv(i) == "" {
			t.Skipf("%s must be set for this test to run", strings.Join(envVars, " "))
		}
	}
}

// WaitForNodesExcludingSelectedStandbys is variation on WaitForActiveNodeAndStandbys.
// It waits for the active node before waiting for standby nodes, however
// it will not wait for cores with indexes that match those specified as arguments.
// Whilst you could specify index 0 which is likely to be the leader node, the function
// checks for the leader first regardless of the indexes to skip, so it would be redundant to do so.
// The intention/use case for this function is to allow a cluster to start and become active with one
// or more nodes not joined, so that we can test scenarios where a node joins later.
// e.g. 4 nodes in the cluster, only 3 nodes in cluster 'active', 1 node can be joined later in tests.
func WaitForNodesExcludingSelectedStandbys(t testing.TB, cluster *vault.TestCluster, indexesToSkip ...int) {
	WaitForActiveNode(t, cluster)

	contains := func(elems []int, e int) bool {
		for _, v := range elems {
			if v == e {
				return true
			}
		}

		return false
	}
	for i, core := range cluster.Cores {
		if contains(indexesToSkip, i) {
			continue
		}

		if standby, _ := core.Core.Standby(); standby {
			WaitForStandbyNode(t, core)
		}
	}
}

// IsLocalOrRegressionTests returns true when the tests are running locally (not in CI), or when
// the regression test env var (VAULT_REGRESSION_TESTS) is provided.
func IsLocalOrRegressionTests() bool {
	return os.Getenv("CI") == "" || os.Getenv("VAULT_REGRESSION_TESTS") == "true"
}

func RaftDataDir(t testing.TB, core *vault.TestClusterCore) string {
	t.Helper()
	r, ok := core.UnderlyingStorage.(*raft.RaftBackend)
	if !ok {
		r, ok = core.UnderlyingHAStorage.(*raft.RaftBackend)
		if !ok {
			t.Fatal("no raft backend")
		}
	}
	return r.DataDir(t)
}
