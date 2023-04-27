package testcluster

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

func GetPerformanceToken(pri VaultCluster, id, secondaryPublicKey string) (string, error) {
	client := pri.Nodes()[0].APIClient()
	req := map[string]interface{}{
		"id": id,
	}
	if secondaryPublicKey != "" {
		req["secondary_public_key"] = secondaryPublicKey
	}
	secret, err := client.Logical().Write("sys/replication/performance/primary/secondary-token", req)
	if err != nil {
		return "", err
	}

	if secondaryPublicKey != "" {
		return secret.Data["token"].(string), nil
	}
	return secret.WrapInfo.Token, nil
}

func EnablePerfPrimary(ctx context.Context, pri VaultCluster) error {
	client := pri.Nodes()[0].APIClient()
	_, err := client.Logical().WriteWithContext(ctx, "sys/replication/performance/primary/enable", nil)
	if err != nil {
		return err
	}

	err = WaitForPerfReplicationState(ctx, pri, consts.ReplicationPerformancePrimary)
	if err != nil {
		return err
	}
	return WaitForActiveNodeAndPerfStandbys(ctx, pri)
}

func WaitForPerfReplicationState(ctx context.Context, cluster VaultCluster, state consts.ReplicationState) error {
	client := cluster.Nodes()[0].APIClient()
	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		health, err = client.Sys().HealthWithContext(ctx)
		if health.ReplicationPerformanceMode == state.GetPerformanceString() {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func EnablePerformanceSecondaryNoWait(ctx context.Context, perfToken string, pri, sec VaultCluster, updatePrimary bool) error {
	postData := map[string]interface{}{
		"token": perfToken,
		//"ca_file": pri.GetCACertPEMFile(),
		"ca_file": "/vault/config/ca.pem",
	}
	//if pri.ClientAuthRequired {
	//	p := pri.Cores[0]
	//	postData["client_cert_pem"] = string(p.ServerCertPEM)
	//	postData["client_key_pem"] = string(p.ServerKeyPEM)
	//}
	path := "sys/replication/performance/secondary/enable"
	if updatePrimary {
		path = "sys/replication/performance/secondary/update-primary"
	}
	err := WaitForActiveNodeAndPerfStandbys(ctx, sec)
	if err != nil {
		return err
	}
	_, err = sec.Nodes()[0].APIClient().Logical().Write(path, postData)
	if err != nil {
		return err
	}

	return WaitForPerfReplicationState(ctx, sec, consts.ReplicationPerformanceSecondary)
}

func EnablePerformanceSecondary(ctx context.Context, perfToken string, pri, sec VaultCluster, updatePrimary, skipPoisonPill bool) (string, error) {
	if err := EnablePerformanceSecondaryNoWait(ctx, perfToken, pri, sec, updatePrimary); err != nil {
		return "", err
	}
	if err := WaitForMatchingMerkleRoots(ctx, "sys/replication/performance/", pri, sec); err != nil {
		return "", err
	}
	root, err := WaitForPerformanceSecondary(ctx, pri, sec, skipPoisonPill)
	if err != nil {
		return "", err
	}
	if err := WaitForPerfReplicationWorking(ctx, pri, sec); err != nil {
		return "", err
	}
	return root, nil
}

func WaitForMatchingMerkleRoots(ctx context.Context, endpoint string, pri, sec VaultCluster) error {
	getRoot := func(mode string, cli *api.Client) (string, error) {
		status, err := cli.Logical().Read(endpoint + "status")
		if err != nil {
			return "", err
		}
		if status == nil || status.Data == nil || status.Data["mode"] == nil {
			return "", fmt.Errorf("got nil secret or data")
		}
		if status.Data["mode"].(string) != mode {
			return "", fmt.Errorf("expected mode=%s, got %s", mode, status.Data["mode"].(string))
		}
		return status.Data["merkle_root"].(string), nil
	}

	secClient := sec.Nodes()[0].APIClient()
	priClient := pri.Nodes()[0].APIClient()
	for i := 0; i < 30; i++ {
		secRoot, err := getRoot("secondary", secClient)
		if err != nil {
			return err
		}
		priRoot, err := getRoot("primary", priClient)
		if err != nil {
			return err
		}

		if reflect.DeepEqual(priRoot, secRoot) {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("roots did not become equal")
}

func WaitForPerformanceSecondary(ctx context.Context, pri, sec VaultCluster, skipPoisonPill bool) (string, error) {
	if len(pri.GetRecoveryKeys()) > 0 {
		sec.SetBarrierKeys(pri.GetRecoveryKeys())
		sec.SetRecoveryKeys(pri.GetRecoveryKeys())
	} else {
		sec.SetBarrierKeys(pri.GetBarrierKeys())
		sec.SetRecoveryKeys(pri.GetBarrierKeys())
	}

	if len(sec.Nodes()) > 1 {
		if skipPoisonPill {
			// As part of prepareSecondary on the active node the keyring is
			// deleted from storage.  Its absence can cause standbys to seal
			// themselves. But it's not reliable, so we'll seal them
			// ourselves to force the issue.
			for i := range sec.Nodes()[1:] {
				if err := SealNode(ctx, sec, i+1); err != nil {
					return "", err
				}
			}
		} else {
			// We want to make sure we unseal all the nodes so we first need to wait
			// until two of the nodes seal due to the poison pill being written
			if err := WaitForNCoresSealed(ctx, sec, len(sec.Nodes())-1); err != nil {
				return "", err
			}
		}
	}
	if _, err := WaitForActiveNode(ctx, sec); err != nil {
		return "", err
	}
	if err := UnsealAllNodes(ctx, sec); err != nil {
		return "", err
	}

	perfSecondaryRootToken, err := GenerateRoot(sec, GenerateRootRegular)
	if err != nil {
		return "", err
	}
	sec.SetRootToken(perfSecondaryRootToken)
	if err := WaitForActiveNodeAndPerfStandbys(ctx, sec); err != nil {
		return "", err
	}

	return perfSecondaryRootToken, nil
}

func WaitForPerfReplicationWorking(ctx context.Context, pri, sec VaultCluster) error {
	priClient, secClient := pri.Nodes()[0].APIClient(), sec.Nodes()[0].APIClient()
	mountPoint, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	err = priClient.Sys().Mount(mountPoint, &api.MountInput{
		Type:  "kv",
		Local: false,
	})
	if err != nil {
		return fmt.Errorf("unable to mount KV engine on primary")
	}

	path := mountPoint + "/foo"
	_, err = priClient.Logical().Write(path, map[string]interface{}{
		"bar": 1,
	})
	if err != nil {
		return fmt.Errorf("unable to write KV on primary", "path", path)
	}

	for ctx.Err() == nil {
		var secret *api.Secret
		secret, err = secClient.Logical().Read(path)
		if err == nil && secret != nil {
			err = priClient.Sys().Unmount(mountPoint)
			if err != nil {
				return fmt.Errorf("unable to unmount KV engine on primary")
			}
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("unable to read replicated KV on secondary", "path", path, "err", err)
}

func SetupTwoClusterPerfReplication(ctx context.Context, pri, sec VaultCluster) error {
	if err := EnablePerfPrimary(ctx, pri); err != nil {
		return err
	}
	perfToken, err := GetPerformanceToken(pri, sec.ClusterID(), "")
	if err != nil {
		return err
	}

	_, err = EnablePerformanceSecondary(ctx, perfToken, pri, sec, false, false)
	return err
}

type ReplicationSet struct {
	// By convention, we recommend the following naming scheme for
	// clusters in this map:
	// A: perf primary
	// B: primary's DR
	// C: first perf secondary of A
	// D: C's DR
	// E: second perf secondary of A
	// F: E's DR
	// ... etc.
	//
	// We use generic names rather than role-specific names because
	// that's less confusing when promotions take place that result in role
	// changes. In other words, if D gets promoted to replace C as a perf
	// secondary, and C gets demoted and updated to become D's DR secondary,
	// they should maintain their initial names of D and C throughout.
	Clusters map[string]VaultCluster
	Builder  ClusterBuilder
	Logger   hclog.Logger
	CA       *CA
}

type ClusterBuilder func(ctx context.Context, name string, logger hclog.Logger) (VaultCluster, error)

func NewReplicationSet(b ClusterBuilder) (*ReplicationSet, error) {
	return &ReplicationSet{
		Clusters: map[string]VaultCluster{},
		Builder:  b,
		Logger:   hclog.NewNullLogger(),
	}, nil
}

func (r *ReplicationSet) StandardPerfReplication(ctx context.Context) error {
	for _, name := range []string{"A", "C"} {
		if _, ok := r.Clusters[name]; !ok {
			cluster, err := r.Builder(ctx, name, r.Logger)
			if err != nil {
				return err
			}
			r.Clusters[name] = cluster
		}
	}

	err := SetupTwoClusterPerfReplication(ctx, r.Clusters["A"], r.Clusters["C"])
	if err != nil {
		return err
	}

	return nil
}

func (r *ReplicationSet) Cleanup() {
	for _, cluster := range r.Clusters {
		cluster.Cleanup()
	}
}
