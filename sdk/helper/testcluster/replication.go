// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testcluster

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/mapstructure"
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
		return fmt.Errorf("error enabling perf primary: %w", err)
	}

	err = WaitForPerfReplicationState(ctx, pri, consts.ReplicationPerformancePrimary)
	if err != nil {
		return fmt.Errorf("error waiting for perf primary to have the correct state: %w", err)
	}
	return WaitForActiveNodeAndPerfStandbys(ctx, pri)
}

func WaitForPerfReplicationState(ctx context.Context, cluster VaultCluster, state consts.ReplicationState) error {
	client := cluster.Nodes()[0].APIClient()
	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		health, err = client.Sys().HealthWithContext(ctx)
		if err == nil && health.ReplicationPerformanceMode == state.GetPerformanceString() {
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
		"token":   perfToken,
		"ca_file": pri.GetCACertPEMFile(),
	}
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
	return WaitForMatchingMerkleRootsClients(ctx, endpoint, pri.Nodes()[0].APIClient(), sec.Nodes()[0].APIClient())
}

func WaitForMatchingMerkleRootsClients(ctx context.Context, endpoint string, pri, sec *api.Client) error {
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

	var priRoot, secRoot string
	var err error
	genRet := func() error {
		return fmt.Errorf("unequal merkle roots, pri=%s sec=%s, err=%w", priRoot, secRoot, err)
	}
	for ctx.Err() == nil {
		secRoot, err = getRoot("secondary", sec)
		if err != nil {
			return genRet()
		}
		priRoot, err = getRoot("primary", pri)
		if err != nil {
			return genRet()
		}

		if reflect.DeepEqual(priRoot, secRoot) {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("roots did not become equal")
}

func WaitForPerformanceWAL(ctx context.Context, pri, sec VaultCluster) error {
	endpoint := "sys/replication/performance/"
	if err := WaitForMatchingMerkleRoots(ctx, endpoint, pri, sec); err != nil {
		return nil
	}
	getWAL := func(mode, walKey string, cli *api.Client) (int64, error) {
		status, err := cli.Logical().Read(endpoint + "status")
		if err != nil {
			return 0, err
		}
		if status == nil || status.Data == nil || status.Data["mode"] == nil {
			return 0, fmt.Errorf("got nil secret or data")
		}
		if status.Data["mode"].(string) != mode {
			return 0, fmt.Errorf("expected mode=%s, got %s", mode, status.Data["mode"].(string))
		}
		return status.Data[walKey].(json.Number).Int64()
	}

	secClient := sec.Nodes()[0].APIClient()
	priClient := pri.Nodes()[0].APIClient()
	for ctx.Err() == nil {
		secLastRemoteWAL, err := getWAL("secondary", "last_remote_wal", secClient)
		if err != nil {
			return err
		}
		priLastPerfWAL, err := getWAL("primary", "last_performance_wal", priClient)
		if err != nil {
			return err
		}

		if secLastRemoteWAL >= priLastPerfWAL {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("performance WALs on the secondary did not catch up with the primary, context err: %w", ctx.Err())
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
	priActiveIdx, err := WaitForActiveNode(ctx, pri)
	if err != nil {
		return err
	}
	secActiveIdx, err := WaitForActiveNode(ctx, sec)
	if err != nil {
		return err
	}

	priClient, secClient := pri.Nodes()[priActiveIdx].APIClient(), sec.Nodes()[secActiveIdx].APIClient()
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
		return fmt.Errorf("unable to write KV on primary, path=%s", path)
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
	if err == nil {
		err = ctx.Err()
	}
	return fmt.Errorf("unable to read replicated KV on secondary, path=%s, err=%v", path, err)
}

func SetupTwoClusterPerfReplication(ctx context.Context, pri, sec VaultCluster) error {
	if err := EnablePerfPrimary(ctx, pri); err != nil {
		return fmt.Errorf("failed to enable perf primary: %w", err)
	}
	perfToken, err := GetPerformanceToken(pri, sec.ClusterID(), "")
	if err != nil {
		return fmt.Errorf("failed to get performance token from perf primary: %w", err)
	}

	_, err = EnablePerformanceSecondary(ctx, perfToken, pri, sec, false, false)
	if err != nil {
		return fmt.Errorf("failed to enable perf secondary: %w", err)
	}
	return nil
}

// PassiveWaitForActiveNodeAndPerfStandbys should be used instead of
// WaitForActiveNodeAndPerfStandbys when you don't want to do any writes
// as a side-effect. This returns perfStandby nodes in the cluster and
// an error.
func PassiveWaitForActiveNodeAndPerfStandbys(ctx context.Context, pri VaultCluster) (VaultClusterNode, []VaultClusterNode, error) {
	leaderNode, standbys, err := GetActiveAndStandbys(ctx, pri)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive standby nodes, %w", err)
	}

	for i, node := range standbys {
		client := node.APIClient()
		// Make sure we get perf standby nodes
		if err = EnsureCoreIsPerfStandby(ctx, client); err != nil {
			return nil, nil, fmt.Errorf("standby node %d is not a perfStandby, %w", i, err)
		}
	}

	return leaderNode, standbys, nil
}

func GetActiveAndStandbys(ctx context.Context, cluster VaultCluster) (VaultClusterNode, []VaultClusterNode, error) {
	var leaderIndex int
	var err error
	if leaderIndex, err = WaitForActiveNode(ctx, cluster); err != nil {
		return nil, nil, err
	}

	var leaderNode VaultClusterNode
	var nodes []VaultClusterNode
	for i, node := range cluster.Nodes() {
		if i == leaderIndex {
			leaderNode = node
			continue
		}
		nodes = append(nodes, node)
	}

	return leaderNode, nodes, nil
}

func EnsureCoreIsPerfStandby(ctx context.Context, client *api.Client) error {
	var err error
	var health *api.HealthResponse
	for ctx.Err() == nil {
		health, err = client.Sys().HealthWithContext(ctx)
		if err == nil && health.PerformanceStandby {
			return nil
		}
		time.Sleep(time.Millisecond * 500)
	}
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func WaitForDRReplicationState(ctx context.Context, cluster VaultCluster, state consts.ReplicationState) error {
	client := cluster.Nodes()[0].APIClient()
	var health *api.HealthResponse
	var err error
	for ctx.Err() == nil {
		health, err = client.Sys().HealthWithContext(ctx)
		if err == nil && health.ReplicationDRMode == state.GetDRString() {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func EnableDrPrimary(ctx context.Context, pri VaultCluster) error {
	client := pri.Nodes()[0].APIClient()
	_, err := client.Logical().Write("sys/replication/dr/primary/enable", nil)
	if err != nil {
		return err
	}

	err = WaitForDRReplicationState(ctx, pri, consts.ReplicationDRPrimary)
	if err != nil {
		return err
	}
	return WaitForActiveNodeAndPerfStandbys(ctx, pri)
}

func GenerateDRActivationToken(pri VaultCluster, id, secondaryPublicKey string) (string, error) {
	client := pri.Nodes()[0].APIClient()
	req := map[string]interface{}{
		"id": id,
	}
	if secondaryPublicKey != "" {
		req["secondary_public_key"] = secondaryPublicKey
	}
	secret, err := client.Logical().Write("sys/replication/dr/primary/secondary-token", req)
	if err != nil {
		return "", err
	}

	if secondaryPublicKey != "" {
		return secret.Data["token"].(string), nil
	}
	return secret.WrapInfo.Token, nil
}

func WaitForDRSecondary(ctx context.Context, pri, sec VaultCluster, skipPoisonPill bool) error {
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
					return err
				}
			}
		} else {
			// We want to make sure we unseal all the nodes so we first need to wait
			// until two of the nodes seal due to the poison pill being written
			if err := WaitForNCoresSealed(ctx, sec, len(sec.Nodes())-1); err != nil {
				return err
			}
		}
	}
	if _, err := WaitForActiveNode(ctx, sec); err != nil {
		return err
	}

	// unseal nodes
	for i := range sec.Nodes() {
		if err := UnsealNode(ctx, sec, i); err != nil {
			// Sometimes when we get here it's already unsealed on its own
			// and then this fails for DR secondaries so check again
			// The error is "path disabled in replication DR secondary mode".
			if healthErr := NodeHealthy(ctx, sec, i); healthErr != nil {
				// return the original error
				return err
			}
		}
	}

	sec.SetRootToken(pri.GetRootToken())

	if _, err := WaitForActiveNode(ctx, sec); err != nil {
		return err
	}

	return nil
}

func EnableDRSecondaryNoWait(ctx context.Context, sec VaultCluster, drToken string) error {
	postData := map[string]interface{}{
		"token":   drToken,
		"ca_file": sec.GetCACertPEMFile(),
	}

	_, err := sec.Nodes()[0].APIClient().Logical().Write("sys/replication/dr/secondary/enable", postData)
	if err != nil {
		return err
	}

	return WaitForDRReplicationState(ctx, sec, consts.ReplicationDRSecondary)
}

func WaitForReplicationStatus(ctx context.Context, client *api.Client, dr bool, accept func(map[string]interface{}) error) error {
	url := "sys/replication/performance/status"
	if dr {
		url = "sys/replication/dr/status"
	}

	var err error
	var secret *api.Secret
	for ctx.Err() == nil {
		secret, err = client.Logical().Read(url)
		if err == nil && secret != nil && secret.Data != nil {
			if err = accept(secret.Data); err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err == nil {
		err = ctx.Err()
	}

	return fmt.Errorf("unable to get acceptable replication status: error=%v secret=%#v", err, secret)
}

func WaitForDRReplicationWorking(ctx context.Context, pri, sec VaultCluster) error {
	priClient := pri.Nodes()[0].APIClient()
	secClient := sec.Nodes()[0].APIClient()

	// Make sure we've entered stream-wals mode
	err := WaitForReplicationStatus(ctx, secClient, true, func(secret map[string]interface{}) error {
		state := secret["state"]
		if state == string("stream-wals") {
			return nil
		}
		return fmt.Errorf("expected stream-wals replication state, got %v", state)
	})
	if err != nil {
		return err
	}

	// Now write some data and make sure that we see last_remote_wal nonzero, i.e.
	// at least one WAL has been streamed.
	secret, err := priClient.Auth().Token().Create(&api.TokenCreateRequest{})
	if err != nil {
		return err
	}

	// Revoke the token since some tests won't be happy to see it.
	err = priClient.Auth().Token().RevokeTree(secret.Auth.ClientToken)
	if err != nil {
		return err
	}

	err = WaitForReplicationStatus(ctx, secClient, true, func(secret map[string]interface{}) error {
		state := secret["state"]
		if state != string("stream-wals") {
			return fmt.Errorf("expected stream-wals replication state, got %v", state)
		}

		if secret["last_remote_wal"] != nil {
			lastRemoteWal, _ := secret["last_remote_wal"].(json.Number).Int64()
			if lastRemoteWal <= 0 {
				return fmt.Errorf("expected last_remote_wal to be greater than zero")
			}
			return nil
		}

		return fmt.Errorf("replication seems to be still catching up, maybe need to wait more")
	})
	if err != nil {
		return err
	}
	return nil
}

func EnableDrSecondary(ctx context.Context, pri, sec VaultCluster, drToken string) error {
	err := EnableDRSecondaryNoWait(ctx, sec, drToken)
	if err != nil {
		return err
	}

	if err = WaitForMatchingMerkleRoots(ctx, "sys/replication/dr/", pri, sec); err != nil {
		return err
	}

	err = WaitForDRSecondary(ctx, pri, sec, false)
	if err != nil {
		return err
	}

	if err = WaitForDRReplicationWorking(ctx, pri, sec); err != nil {
		return err
	}
	return nil
}

func SetupTwoClusterDRReplication(ctx context.Context, pri, sec VaultCluster) error {
	if err := EnableDrPrimary(ctx, pri); err != nil {
		return err
	}

	drToken, err := GenerateDRActivationToken(pri, sec.ClusterID(), "")
	if err != nil {
		return err
	}
	err = EnableDrSecondary(ctx, pri, sec, drToken)
	if err != nil {
		return err
	}
	return nil
}

func DemoteDRPrimary(client *api.Client) error {
	_, err := client.Logical().Write("sys/replication/dr/primary/demote", map[string]interface{}{})
	return err
}

func createBatchToken(client *api.Client, path string) (string, error) {
	// TODO: should these be more random in case more than one batch token needs to be created?
	suffix := strings.Replace(path, "/", "", -1)
	policyName := "path-batch-policy-" + suffix
	roleName := "path-batch-role-" + suffix

	rules := fmt.Sprintf(`path "%s" { capabilities = [ "read", "update" ] }`, path)

	// create policy
	_, err := client.Logical().Write("sys/policy/"+policyName, map[string]interface{}{
		"policy": rules,
	})
	if err != nil {
		return "", err
	}

	// create a role
	_, err = client.Logical().Write("auth/token/roles/"+roleName, map[string]interface{}{
		"allowed_policies": policyName,
		"orphan":           true,
		"renewable":        false,
		"token_type":       "batch",
	})
	if err != nil {
		return "", err
	}

	// create batch token
	secret, err := client.Logical().Write("auth/token/create/"+roleName, nil)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}

// PromoteDRSecondaryWithBatchToken creates a batch token for DR promotion
// before promotion, it demotes the primary cluster. The primary cluster needs
// to be functional for the generation of the batch token
func PromoteDRSecondaryWithBatchToken(ctx context.Context, pri, sec VaultCluster) error {
	client := pri.Nodes()[0].APIClient()
	drToken, err := createBatchToken(client, "sys/replication/dr/secondary/promote")
	if err != nil {
		return err
	}

	err = DemoteDRPrimary(client)
	if err != nil {
		return err
	}

	return promoteDRSecondaryInternal(ctx, sec, drToken)
}

// PromoteDRSecondary generates a DR operation token on the secondary using
// unseal/recovery keys. Therefore, the primary cluster could potentially
// be out of service.
func PromoteDRSecondary(ctx context.Context, sec VaultCluster) error {
	// generate DR operation token to do update primary on vC to point to
	// the new perfSec primary vD
	drToken, err := GenerateRoot(sec, GenerateRootDR)
	if err != nil {
		return err
	}
	return promoteDRSecondaryInternal(ctx, sec, drToken)
}

func promoteDRSecondaryInternal(ctx context.Context, sec VaultCluster, drToken string) error {
	secClient := sec.Nodes()[0].APIClient()

	// Allow retries of 503s, e.g.: replication is still catching up,
	// try again later or provide the "force" argument
	oldMaxRetries := secClient.MaxRetries()
	secClient.SetMaxRetries(10)
	defer secClient.SetMaxRetries(oldMaxRetries)
	resp, err := secClient.Logical().Write("sys/replication/dr/secondary/promote", map[string]interface{}{
		"dr_operation_token": drToken,
	})
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil status response during DR promotion")
	}

	if _, err := WaitForActiveNode(ctx, sec); err != nil {
		return err
	}

	return WaitForDRReplicationState(ctx, sec, consts.ReplicationDRPrimary)
}

func checkClusterAddr(ctx context.Context, pri, sec VaultCluster) error {
	priClient := pri.Nodes()[0].APIClient()
	priLeader, err := priClient.Sys().LeaderWithContext(ctx)
	if err != nil {
		return err
	}
	secClient := sec.Nodes()[0].APIClient()
	endpoint := "sys/replication/dr/"
	status, err := secClient.Logical().Read(endpoint + "status")
	if err != nil {
		return err
	}
	if status == nil || status.Data == nil {
		return fmt.Errorf("got nil secret or data")
	}

	var priAddrs []string
	err = mapstructure.Decode(status.Data["known_primary_cluster_addrs"], &priAddrs)
	if err != nil {
		return err
	}
	if !strutil.StrListContains(priAddrs, priLeader.LeaderClusterAddress) {
		return fmt.Errorf("failed to fine the expected primary cluster address %v in known_primary_cluster_addrs", priLeader.LeaderClusterAddress)
	}

	return nil
}

func UpdatePrimary(ctx context.Context, pri, sec VaultCluster) error {
	// generate DR operation token to do update primary on vC to point to
	// the new perfSec primary vD
	rootToken, err := GenerateRoot(sec, GenerateRootDR)
	if err != nil {
		return err
	}

	// secondary activation token
	drToken, err := GenerateDRActivationToken(pri, sec.ClusterID(), "")
	if err != nil {
		return err
	}

	// update-primary on vC (new perfSec Dr secondary) to point to
	// the new perfSec Dr primary
	secClient := sec.Nodes()[0].APIClient()
	resp, err := secClient.Logical().Write("sys/replication/dr/secondary/update-primary", map[string]interface{}{
		"dr_operation_token": rootToken,
		"token":              drToken,
		"ca_file":            sec.GetCACertPEMFile(),
	})
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil status response during update primary")
	}

	if _, err = WaitForActiveNode(ctx, sec); err != nil {
		return err
	}

	if err = WaitForDRReplicationState(ctx, sec, consts.ReplicationDRSecondary); err != nil {
		return err
	}

	if err = checkClusterAddr(ctx, pri, sec); err != nil {
		return err
	}

	return nil
}

func SetupFourClusterReplication(ctx context.Context, pri, sec, pridr, secdr VaultCluster) error {
	err := SetupTwoClusterPerfReplication(ctx, pri, sec)
	if err != nil {
		return err
	}
	err = SetupTwoClusterDRReplication(ctx, pri, pridr)
	if err != nil {
		return err
	}
	err = SetupTwoClusterDRReplication(ctx, sec, secdr)
	if err != nil {
		return err
	}
	return nil
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := SetupTwoClusterPerfReplication(ctx, r.Clusters["A"], r.Clusters["C"])
	if err != nil {
		return err
	}

	return nil
}

func (r *ReplicationSet) StandardDRReplication(ctx context.Context) error {
	for _, name := range []string{"A", "B"} {
		if _, ok := r.Clusters[name]; !ok {
			cluster, err := r.Builder(ctx, name, r.Logger)
			if err != nil {
				return err
			}
			r.Clusters[name] = cluster
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := SetupTwoClusterDRReplication(ctx, r.Clusters["A"], r.Clusters["B"])
	if err != nil {
		return err
	}

	return nil
}

func (r *ReplicationSet) GetFourReplicationCluster(ctx context.Context) error {
	for _, name := range []string{"A", "B", "C", "D"} {
		if _, ok := r.Clusters[name]; !ok {
			cluster, err := r.Builder(ctx, name, r.Logger)
			if err != nil {
				return err
			}
			r.Clusters[name] = cluster
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := SetupFourClusterReplication(ctx, r.Clusters["A"], r.Clusters["C"], r.Clusters["B"], r.Clusters["D"])
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

func WaitForPerfReplicationConnectionStatus(ctx context.Context, client *api.Client) error {
	type Primary struct {
		APIAddress       string `mapstructure:"api_address"`
		ConnectionStatus string `mapstructure:"connection_status"`
		ClusterAddress   string `mapstructure:"cluster_address"`
		LastHeartbeat    string `mapstructure:"last_heartbeat"`
	}
	type Status struct {
		Primaries []Primary `mapstructure:"primaries"`
	}
	return WaitForPerfReplicationStatus(ctx, client, func(m map[string]interface{}) error {
		var status Status
		err := mapstructure.Decode(m, &status)
		if err != nil {
			return err
		}
		if len(status.Primaries) == 0 {
			return fmt.Errorf("primaries is zero")
		}
		for _, v := range status.Primaries {
			if v.ConnectionStatus == "connected" {
				return nil
			}
		}
		return fmt.Errorf("no primaries connected")
	})
}

func WaitForPerfReplicationStatus(ctx context.Context, client *api.Client, accept func(map[string]interface{}) error) error {
	var err error
	var secret *api.Secret
	for ctx.Err() == nil {
		secret, err = client.Logical().Read("sys/replication/performance/status")
		if err == nil && secret != nil && secret.Data != nil {
			if err = accept(secret.Data); err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("unable to get acceptable replication status within allotted time: error=%v secret=%#v", err, secret)
}
