// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hcp_link

import (
	"os"
	"testing"
	"time"

	sdkResource "github.com/hashicorp/hcp-sdk-go/resource"
	"github.com/hashicorp/vault/api"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/hcp_link"
)

type VaultHCPLinkInstances struct {
	instances []*hcp_link.HCPLinkVault
}

func NewVaultHCPLinkInstances() *VaultHCPLinkInstances {
	i := &VaultHCPLinkInstances{
		instances: make([]*hcp_link.HCPLinkVault, 0),
	}

	return i
}

func (v *VaultHCPLinkInstances) Cleanup() {
	for _, inst := range v.instances {
		inst.Shutdown()
	}
}

func getHCPConfig(t *testing.T, clientID, clientSecret string) *configutil.HCPLinkConfig {
	resourceIDRaw, ok := os.LookupEnv("HCP_RESOURCE_ID")
	if !ok {
		t.Skip("failed to find the HCP resource ID")
	}
	res, err := sdkResource.FromString(resourceIDRaw)
	if err != nil {
		t.Fatalf("failed to parse the resource ID, %v", err.Error())
	}
	return &configutil.HCPLinkConfig{
		ResourceIDRaw: resourceIDRaw,
		Resource:      &res,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
	}
}

func getTestCluster(t *testing.T, numCores int) *vault.TestCluster {
	t.Helper()
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}

	if numCores <= 0 {
		numCores = 1
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    numCores,
	})

	return cluster
}

func TestClusterWithHCPLinkEnabled(t *testing.T, cluster *vault.TestCluster, enableAPICap, enablePassthroughCap bool) (*VaultHCPLinkInstances, *configutil.HCPLinkConfig) {
	t.Helper()
	clientID, ok := os.LookupEnv("HCP_CLIENT_ID")
	if !ok {
		t.Skip("HCP client ID not found in env")
	}
	clientSecret, ok := os.LookupEnv("HCP_CLIENT_SECRET")
	if !ok {
		t.Skip("HCP client secret not found in env")
	}

	if _, ok := os.LookupEnv("HCP_API_ADDRESS"); !ok {
		t.Skip("failed to find HCP_API_ADDRESS in the environment")
	}
	if _, ok := os.LookupEnv("HCP_SCADA_ADDRESS"); !ok {
		t.Skip("failed to find HCP_SCADA_ADDRESS in the environment")
	}
	if _, ok := os.LookupEnv("HCP_AUTH_URL"); !ok {
		t.Skip("failed to find HCP_AUTH_URL in the environment")
	}

	hcpConfig := getHCPConfig(t, clientID, clientSecret)
	if enableAPICap {
		hcpConfig.EnableAPICapability = true
	}
	if enablePassthroughCap {
		hcpConfig.EnablePassThroughCapability = true
	}
	hcpLinkIns := NewVaultHCPLinkInstances()

	cluster.Start()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	for _, c := range cluster.Cores {
		logger := c.Logger().Named("hcpLink")
		vaultHCPLink, err := hcp_link.NewHCPLink(hcpConfig, c.Core, logger)
		if err != nil {
			t.Fatalf("failed to start HCP link, %v", err)
		}
		hcpLinkIns.instances = append(hcpLinkIns.instances, vaultHCPLink)
	}

	return hcpLinkIns, hcpConfig
}

func checkLinkStatus(client *api.Client, expectedStatus string, t *testing.T) {
	deadline := time.Now().Add(10 * time.Second)
	var status *api.SealStatusResponse
	var err error
	for time.Now().Before(deadline) {
		status, err = client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if status.HCPLinkStatus == expectedStatus {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if status.HCPLinkStatus != expectedStatus {
		t.Fatalf("HCP link did not behave as expected. expected status %v, actual status %v", expectedStatus, status.HCPLinkStatus)
	}
}
