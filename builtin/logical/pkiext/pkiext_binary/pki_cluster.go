// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pkiext_binary

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki/dnstest"
	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
)

type VaultPkiCluster struct {
	cluster *docker.DockerCluster
	Dns     *dnstest.TestServer
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
			NumCores: 3,
		},
	}

	cluster := docker.NewTestDockerCluster(t, opts)

	return &VaultPkiCluster{cluster: cluster}
}

func NewVaultPkiClusterWithDNS(t *testing.T) *VaultPkiCluster {
	cluster := NewVaultPkiCluster(t)
	dns := dnstest.SetupResolverOnNetwork(t, "dadgarcorp.com", cluster.GetContainerNetworkName())
	cluster.Dns = dns
	return cluster
}

func (vpc *VaultPkiCluster) Cleanup() {
	vpc.cluster.Cleanup()
	if vpc.Dns != nil {
		vpc.Dns.Cleanup()
	}
}

func (vpc *VaultPkiCluster) GetActiveClusterNode() *docker.DockerClusterNode {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	node, err := testcluster.WaitForActiveNode(ctx, vpc.cluster)
	if err != nil {
		panic(fmt.Sprintf("no cluster node became active in timeout window: %v", err))
	}

	return vpc.cluster.ClusterNodes[node]
}

func (vpc *VaultPkiCluster) GetNonActiveNodes() []*docker.DockerClusterNode {
	nodes := []*docker.DockerClusterNode{}
	for _, node := range vpc.cluster.ClusterNodes {
		leader, err := node.APIClient().Sys().Leader()
		if err != nil {
			continue
		}

		if !leader.IsSelf {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (vpc *VaultPkiCluster) GetActiveContainerHostPort() string {
	return vpc.GetActiveClusterNode().HostPort
}

func (vpc *VaultPkiCluster) GetContainerNetworkName() string {
	return vpc.cluster.ClusterNodes[0].ContainerNetworkName
}

func (vpc *VaultPkiCluster) GetActiveContainerIP() string {
	return vpc.GetActiveClusterNode().ContainerIPAddress
}

func (vpc *VaultPkiCluster) GetActiveContainerID() string {
	return vpc.GetActiveClusterNode().Container.ID
}

func (vpc *VaultPkiCluster) GetActiveNode() *api.Client {
	return vpc.GetActiveClusterNode().APIClient()
}

// GetListenerCACertPEM returns the Vault cluster's PEM-encoded CA certificate.
func (vpc *VaultPkiCluster) GetListenerCACertPEM() []byte {
	return vpc.cluster.CACertPEM
}

func (vpc *VaultPkiCluster) AddHostname(hostname, ip string) error {
	if vpc.Dns != nil {
		vpc.Dns.AddRecord(hostname, "A", ip)
		vpc.Dns.PushConfig()
		return nil
	} else {
		return vpc.AddNameToHostFiles(hostname, ip)
	}
}

func (vpc *VaultPkiCluster) AddNameToHostFiles(hostname, ip string) error {
	updateHostsCmd := []string{
		"sh", "-c",
		"echo '" + ip + " " + hostname + "' >> /etc/hosts",
	}
	for _, node := range vpc.cluster.ClusterNodes {
		containerID := node.Container.ID
		_, _, retcode, err := dockhelper.RunCmdWithOutput(vpc.cluster.DockerAPI, context.Background(), containerID, updateHostsCmd)
		if err != nil {
			return fmt.Errorf("failed updating container %s host file: %w", containerID, err)
		}

		if retcode != 0 {
			return fmt.Errorf("expected zero retcode from updating vault host file in container %s got: %d", containerID, retcode)
		}
	}

	return nil
}

func (vpc *VaultPkiCluster) AddDNSRecord(hostname, recordType, ip string) error {
	if vpc.Dns == nil {
		return fmt.Errorf("no DNS server was provisioned on this cluster group; unable to provision custom records")
	}

	vpc.Dns.AddRecord(hostname, recordType, ip)
	vpc.Dns.PushConfig()
	return nil
}

func (vpc *VaultPkiCluster) RemoveDNSRecord(domain string, record string, value string) error {
	if vpc.Dns == nil {
		return fmt.Errorf("no DNS server was provisioned on this cluster group; unable to remove specific record")
	}

	vpc.Dns.RemoveRecord(domain, record, value)
	return nil
}

func (vpc *VaultPkiCluster) RemoveDNSRecordsOfTypeForDomain(domain string, record string) error {
	if vpc.Dns == nil {
		return fmt.Errorf("no DNS server was provisioned on this cluster group; unable to remove all records of type")
	}

	vpc.Dns.RemoveRecordsOfTypeForDomain(domain, record)
	return nil
}

func (vpc *VaultPkiCluster) RemoveDNSRecordsForDomain(domain string) error {
	if vpc.Dns == nil {
		return fmt.Errorf("no DNS server was provisioned on this cluster group; unable to remove records for domain")
	}

	vpc.Dns.RemoveRecordsForDomain(domain)
	return nil
}

func (vpc *VaultPkiCluster) RemoveAllDNSRecords() error {
	if vpc.Dns == nil {
		return fmt.Errorf("no DNS server was provisioned on this cluster group; unable to remove all records")
	}

	vpc.Dns.RemoveAllRecords()
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

	cfg := map[string]interface{}{
		"eab_policy": "not-required",
	}
	if vpc.Dns != nil {
		cfg["dns_resolver"] = vpc.Dns.GetRemoteAddr()
	}

	err = pki.UpdateAcmeConfig(true, cfg)
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
