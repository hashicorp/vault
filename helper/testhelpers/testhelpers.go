package testhelpers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sync/atomic"
	"time"

	raftlib "github.com/hashicorp/raft"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/go-testing-interface"
)

type GenerateRootKind int

const (
	GenerateRootRegular GenerateRootKind = iota
	GenerateRootDR
	GenerateRecovery
)

// Generates a root token on the target cluster.
func GenerateRoot(t testing.T, cluster *vault.TestCluster, kind GenerateRootKind) string {
	t.Helper()
	token, err := GenerateRootWithError(t, cluster, kind)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func GenerateRootWithError(t testing.T, cluster *vault.TestCluster, kind GenerateRootKind) (string, error) {
	t.Helper()
	// If recovery keys supported, use those to perform root token generation instead
	var keys [][]byte
	if cluster.Cores[0].SealAccess().RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}
	client := cluster.Cores[0].Client

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
	for i, core := range c.Cores {
		err := AttemptUnsealCore(c, core)
		if err != nil {
			t.Fatalf("failed to unseal core %d: %v", i, err)
		}
	}
}

func EnsureCoreUnsealed(t testing.T, c *vault.TestCluster, core *vault.TestClusterCore) {
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

func EnsureStableActiveNode(t testing.T, cluster *vault.TestCluster) {
	deriveStableActiveCore(t, cluster)
}

func DeriveStableActiveCore(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	return deriveStableActiveCore(t, cluster)
}

func deriveStableActiveCore(t testing.T, cluster *vault.TestCluster) *vault.TestClusterCore {
	activeCore := DeriveActiveCore(t, cluster)
	minDuration := time.NewTimer(3 * time.Second)

	for i := 0; i < 30; i++ {
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

func WaitForStandbyNode(t testing.T, core *vault.TestClusterCore) {
	t.Helper()
	for i := 0; i < 30; i++ {
		if isLeader, _, clusterAddr, _ := core.Core.Leader(); isLeader != true && clusterAddr != "" {
			return
		}

		time.Sleep(time.Second)
	}

	t.Fatalf("node did not become standby")
}

func RekeyCluster(t testing.T, cluster *vault.TestCluster, recovery bool) [][]byte {
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
	var keys = cluster.BarrierKeys
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

	leaderInfo := &raft.LeaderJoinInfo{
		LeaderAPIAddr: leaderAPI,
		TLSConfig:     leaderCore.TLSConfig,
	}

	// Join core1
	{
		core := cluster.Cores[1]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		leaderInfos := []*raft.LeaderJoinInfo{
			leaderInfo,
		}
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
		if err != nil {
			t.Fatal(err)
		}

		cluster.UnsealCore(t, core)
	}

	// Join core2
	{
		core := cluster.Cores[2]
		core.UnderlyingRawStorage.(*raft.RaftBackend).SetServerAddressProvider(addressProvider)
		leaderInfos := []*raft.LeaderJoinInfo{
			leaderInfo,
		}
		_, err := core.JoinRaftCluster(namespace.RootContext(context.Background()), leaderInfos, false)
		if err != nil {
			t.Fatal(err)
		}

		cluster.UnsealCore(t, core)
	}

	WaitForNCoresUnsealed(t, cluster, 3)
}

func GenerateDebugLogs(t testing.T, stopCh chan struct{}, client *api.Client) {
	t.Helper()

	// This ticker value was chosen after some trial and error trying to figure
	// out why the ticker wasn't stopping when stopCh was closed. It turns out
	// there's some kind of delay between closing a channel in a different
	// goroutine and detecting that it was closed in this one. This ticker interval
	// turned out to be the right value in combination with a 4 second delay
	// in my other goroutine to let this helper produce log output. I'm not sure
	// if I'm doing something wrong here, or if that's just the way this goes.
	ticker := time.NewTicker(1500 * time.Millisecond)
	var err error

	for {
		select {
		case <-stopCh:
			ticker.Stop()
			return
		case <-ticker.C:
			err = client.Sys().Mount("foo", &api.MountInput{
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
}