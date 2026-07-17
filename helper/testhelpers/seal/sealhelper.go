// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package sealhelper

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/stretchr/testify/require"
)

type TransitSealServer struct {
	*vault.TestCluster
}

func NewTransitSealServer(t testing.TB, idx int) *TransitSealServer {
	conf := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	}
	opts := &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: http.Handler,
		Logger:      corehelpers.NewTestLogger(t).Named("transit-seal" + strconv.Itoa(idx)),
	}
	teststorage.InmemBackendSetup(conf, opts)
	cluster := vault.NewTestCluster(t, conf, opts)

	if err := cluster.Cores[0].Client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	}); err != nil {
		t.Fatal(err)
	}

	return &TransitSealServer{cluster}
}

func (tss *TransitSealServer) MakeKey(t testing.TB, key string) {
	client := tss.Cores[0].Client
	if _, err := client.Logical().Write(path.Join("transit", "keys", key), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write(path.Join("transit", "keys", key, "config"), map[string]interface{}{
		"deletion_allowed": true,
	}); err != nil {
		t.Fatal(err)
	}
}

func (tss *TransitSealServer) MakeSeal(t testing.TB, key string) (vault.Seal, error) {
	client := tss.Cores[0].Client
	wrapperConfig := map[string]string{
		"address":     client.Address(),
		"token":       client.Token(),
		"mount_path":  "transit",
		"key_name":    key,
		"tls_ca_cert": tss.CACertPEMFile,
	}
	transitSealWrapper, _, err := configutil.GetTransitKMSFunc(&configutil.KMS{Config: wrapperConfig})
	if err != nil {
		t.Fatalf("error setting wrapper config: %v", err)
	}

	access, err := seal.NewAccessFromWrapper(tss.Logger, transitSealWrapper, vault.SealConfigTypeTransit.String())
	if err != nil {
		return nil, err
	}
	return vault.NewAutoSeal(access), nil
}

type TransitDockerSealServer struct {
	cluster *docker.DockerCluster
	t       *testing.T
}

func NewTransitDockerSealServer(t *testing.T) *TransitDockerSealServer {
	opts := docker.DefaultOptions(t)
	opts.NumCores = 1
	opts.ImageRepo, opts.ImageTag = "hashicorp/vault", "latest"
	opts.VaultNodeConfig.StorageOptions = map[string]string{
		"performance_multiplier": "1",
	}
	opts.DisableTLS = true // simplify, this way we don't have to deal with ca
	opts.ClusterName = strings.ReplaceAll(t.Name()+"-transit", "/", "-")
	return &TransitDockerSealServer{t: t, cluster: docker.NewTestDockerCluster(t, opts)}
}

func (tc *TransitDockerSealServer) APIClient() *api.Client {
	return tc.cluster.Nodes()[0].APIClient()
}

func (tc *TransitDockerSealServer) SealWithPriorityAndDisabled(name string, idx int, disabled bool, priority int) testcluster.VaultNodeSealConfig {
	seal := tc.Seal(name, idx)
	seal.Config["disabled"] = strconv.FormatBool(disabled)
	seal.Config["priority"] = strconv.Itoa(priority)
	return seal
}

// Seal creates a seal using the given mount name and an idx that identifies a key.
// The mount and key will be created.
func (tc *TransitDockerSealServer) Seal(name string, idx int) testcluster.VaultNodeSealConfig {
	client := tc.cluster.Nodes()[0].APIClient()
	if m, _ := client.Sys().GetMount(name); m == nil {
		require.NoError(tc.t, client.Sys().Mount(name, &api.MountInput{
			Type: "transit",
		}))
	}

	keyName := fmt.Sprintf("transit-seal-%d", idx+1)

	_, err := client.Logical().Write(path.Join(name, "keys", keyName), nil)
	require.NoError(tc.t, err)

	return testcluster.VaultNodeSealConfig{
		Type: "transit",
		Config: map[string]string{
			// For another docker container to talk to this cluster they
			// must use the real api address, not the remapped localhost
			// address test code uses.
			"address":    tc.cluster.Nodes()[0].(*docker.DockerClusterNode).RealAPIAddr,
			"token":      tc.cluster.GetRootToken(),
			"mount_path": name,
			"key_name":   keyName,
			"name":       strings.ReplaceAll(name, " ", "_") + "-" + keyName,
			"priority":   "1",
		},
	}
}
