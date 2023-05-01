// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pkiext_binary

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"

	"github.com/hashicorp/vault/sdk/helper/testcluster"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

type VaultPkiCluster struct {
	cluster *docker.DockerCluster
}

func NewVaultPkiCluster(t *testing.T) *VaultPkiCluster {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test when $VAULT_BINARY present")
	}

	opts := &docker.DockerClusterOptions{
		ImageRepo: "docker.mirror.hashicorp.services/hashicorp/vault",
		// We're replacing the binary anyway, so we're not too particular about
		// the docker image version tag.
		ImageTag:    "latest",
		VaultBinary: binary,
		ClusterOptions: testcluster.ClusterOptions{
			VaultNodeConfig: &testcluster.VaultNodeConfig{
				LogLevel: "TRACE",
			},
			NumCores: 1,
		},
	}

	cluster := docker.NewTestDockerCluster(t, opts)

	return &VaultPkiCluster{cluster: cluster}
}

func (vpc *VaultPkiCluster) Cleanup() {
	vpc.cluster.Cleanup()
}

func (vpc *VaultPkiCluster) GetActiveContainerHostPort() string {
	return vpc.cluster.ClusterNodes[0].HostPort
}

func (vpc *VaultPkiCluster) GetContainerNetworkName() string {
	return vpc.cluster.ClusterNodes[0].ContainerNetworkName
}

func (vpc *VaultPkiCluster) GetActiveContainerIP() string {
	return vpc.cluster.ClusterNodes[0].ContainerIPAddress
}

func (vpc *VaultPkiCluster) GetActiveContainerID() string {
	return vpc.cluster.ClusterNodes[0].Container.ID
}

func (vpc *VaultPkiCluster) GetActiveNode() *api.Client {
	return vpc.cluster.Nodes()[0].APIClient()
}

func (vpc *VaultPkiCluster) AddNameToHostsFile(ip, hostname string, logConsumer func(string), logStdout, logStderr io.Writer) error {
	updateHostsCmd := []string{
		"sh", "-c",
		"echo '" + ip + " " + hostname + "' >> /etc/hosts",
	}
	for _, node := range vpc.cluster.ClusterNodes {
		containerID := node.Container.ID
		runner, err := dockhelper.NewServiceRunner(dockhelper.RunOptions{
			ImageRepo:     node.ImageRepo,
			ImageTag:      node.ImageTag,
			ContainerName: containerID,
			NetworkName:   node.ContainerNetworkName,
			LogConsumer:   logConsumer,
			LogStdout:     logStdout,
			LogStderr:     logStderr,
		})
		if err != nil {
			return err
		}

		_, _, retcode, err := runner.RunCmdWithOutput(context.Background(), containerID, updateHostsCmd)
		if err != nil {
			return fmt.Errorf("failed updating container %s host file: %w", containerID, err)
		}

		if retcode != 0 {
			return fmt.Errorf("expected zero retcode from updating vault host file in container %s got: %d", containerID, retcode)
		}
	}

	return nil
}

func (vpc *VaultPkiCluster) CreateMount(name string) (*VaultPkiMount, error) {
	err := vpc.GetActiveNode().Sys().Mount(name, &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
			AllowedResponseHeaders: []string{
				"Last-Modified", "Replay-Nonce",
				"Link", "Location",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &VaultPkiMount{
		vpc,
		name,
	}, nil
}

func (vpc *VaultPkiCluster) CreateAcmeMount(mountName string) (*VaultPkiMount, error) {
	pki, err := vpc.CreateMount(mountName)
	if err != nil {
		return nil, fmt.Errorf("failed creating mount %s: %w", mountName, err)
	}

	err = pki.UpdateClusterConfig(nil)
	if err != nil {
		return nil, fmt.Errorf("failed updating cluster config: %w", err)
	}

	err = pki.UpdateAcmeConfig(true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed updating acme config: %w", err)
	}

	// Setup root+intermediate CA hierarchy within this mount.
	resp, err := pki.GenerateRootInternal(map[string]interface{}{
		"common_name":  "Root X1",
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     "ec",
		"key_bits":     256,
		"use_pss":      false,
		"issuer_name":  "root",
	})
	if err != nil {
		return nil, fmt.Errorf("failed generating root internal: %w", err)
	}
	if resp == nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("failed generating root internal: nil or empty response but no error")
	}

	resp, err = pki.GenerateIntermediateInternal(map[string]interface{}{
		"common_name":  "Intermediate I1",
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     "ec",
		"key_bits":     256,
		"use_pss":      false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed generating int csr: %w", err)
	}
	if resp == nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("failed generating int csr: nil or empty response but no error")
	}

	resp, err = pki.SignIntermediary("default", resp.Data["csr"], map[string]interface{}{
		"common_name":  "Intermediate I1",
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     "ec",
		"csr":          resp.Data["csr"],
	})
	if err != nil {
		return nil, fmt.Errorf("failed signing int csr: %w", err)
	}
	if resp == nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("failed signing int csr: nil or empty response but no error")
	}
	intCert := resp.Data["certificate"].(string)

	resp, err = pki.ImportBundle(intCert, nil)
	if err != nil {
		return nil, fmt.Errorf("failed importing signed cert: %w", err)
	}
	if resp == nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("failed importing signed cert: nil or empty response but no error")
	}

	err = pki.UpdateDefaultIssuer(resp.Data["imported_issuers"].([]interface{})[0].(string), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to set intermediate as default: %w", err)
	}

	err = pki.UpdateIssuer("default", map[string]interface{}{
		"leaf_not_after_behavior": "truncate",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update intermediate ttl behavior: %w", err)
	}

	err = pki.UpdateIssuer("root", map[string]interface{}{
		"leaf_not_after_behavior": "truncate",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update root ttl behavior: %w", err)
	}

	return pki, nil
}
