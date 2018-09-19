package testhelpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/xor"
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

func EnsureCoresUnsealed(t testing.T, c *vault.TestCluster) {
	t.Helper()
	for _, core := range c.Cores {
		if !core.Sealed() {
			continue
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

func SetupFourClusterReplication(t testing.T, perfPrimary, perfSecondary, perfDRSecondary, perfSecondaryDRSecondary *vault.TestCluster) {
	// Enable dr primary
	_, err := perfPrimary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, perfPrimary.Cores[0].Core, consts.ReplicationDRPrimary)

	// Enable performance primary
	_, err = perfPrimary.Cores[0].Client.Logical().Write("sys/replication/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, perfPrimary.Cores[0].Core, consts.ReplicationPerformancePrimary)

	// get dr token
	secret, err := perfPrimary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/secondary-token", map[string]interface{}{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	token := secret.WrapInfo.Token

	// enable dr secondary
	secret, err = perfDRSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/secondary/enable", map[string]interface{}{
		"token":   token,
		"ca_file": perfPrimary.CACertPEMFile,
	})
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, perfDRSecondary.Cores[0].Core, consts.ReplicationDRSecondary)
	perfDRSecondary.BarrierKeys = perfPrimary.BarrierKeys
	EnsureCoresUnsealed(t, perfDRSecondary)

	// get performance token
	secret, err = perfPrimary.Cores[0].Client.Logical().Write("sys/replication/primary/secondary-token", map[string]interface{}{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	token = secret.WrapInfo.Token

	// enable performace secondary
	secret, err = perfSecondary.Cores[0].Client.Logical().Write("sys/replication/secondary/enable", map[string]interface{}{
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

	// Enable dr primary on perf secondary
	_, err = perfSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, perfSecondary.Cores[0].Core, consts.ReplicationDRPrimary)

	// get dr token from perf secondary
	secret, err = perfSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/primary/secondary-token", map[string]interface{}{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	token = secret.WrapInfo.Token

	// enable dr secondary
	secret, err = perfSecondaryDRSecondary.Cores[0].Client.Logical().Write("sys/replication/dr/secondary/enable", map[string]interface{}{
		"token":   token,
		"ca_file": perfSecondary.CACertPEMFile,
	})
	if err != nil {
		t.Fatal(err)
	}

	WaitForReplicationState(t, perfSecondaryDRSecondary.Cores[0].Core, consts.ReplicationDRSecondary)
	perfSecondaryDRSecondary.BarrierKeys = perfPrimary.BarrierKeys
	EnsureCoresUnsealed(t, perfSecondaryDRSecondary)

	perfDRSecondary.Cores[0].Client.SetToken(perfPrimary.Cores[0].Client.Token())
	perfSecondaryDRSecondary.Cores[0].Client.SetToken(rootToken)
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
