package testhelpers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
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

// StrPtr returns a pointer to a string value
func StrPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a boolean value
func BoolPtr(b bool) *bool {
	return &b
}

// TimeDurationPtr returns a pointer to a time duration value
func TimeDurationPtr(duration string) *time.Duration {
	d, _ := time.ParseDuration(duration)

	return &d
}

// FileMode returns a pointer to the given os.FileMode
func FileModePtr(o os.FileMode) *os.FileMode {
	return &o
}
