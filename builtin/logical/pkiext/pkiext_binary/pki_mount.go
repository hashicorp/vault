// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pkiext_binary

import (
	"context"
	"encoding/base64"
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
)

type VaultPkiMount struct {
	*VaultPkiCluster
	mount string
}

func (vpm *VaultPkiMount) UpdateClusterConfig(config map[string]interface{}) error {
	defaultPath := "https://" + vpm.cluster.ClusterNodes[0].ContainerIPAddress + ":8200/v1/" + vpm.mount
	defaults := map[string]interface{}{
		"path":     defaultPath,
		"aia_path": defaultPath,
	}

	_, err := vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/config/cluster", mergeWithDefaults(config, defaults))
	return err
}

func (vpm *VaultPkiMount) UpdateClusterConfigLocalAddr() (string, error) {
	basePath := fmt.Sprintf("https://%s/v1/%s", vpm.GetActiveContainerHostPort(), vpm.mount)
	return basePath, vpm.UpdateClusterConfig(map[string]interface{}{
		"path": basePath,
	})
}

func (vpm *VaultPkiMount) UpdateAcmeConfig(enable bool, config map[string]interface{}) error {
	defaults := map[string]interface{}{
		"enabled": enable,
	}

	_, err := vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/config/acme", mergeWithDefaults(config, defaults))
	return err
}

func (vpm *VaultPkiMount) GenerateRootInternal(props map[string]interface{}) (*api.Secret, error) {
	defaults := map[string]interface{}{
		"common_name": "root-test.com",
		"key_type":    "ec",
		"issuer_name": "root",
	}

	return vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/root/generate/internal", mergeWithDefaults(props, defaults))
}

func (vpm *VaultPkiMount) GenerateIntermediateInternal(props map[string]interface{}) (*api.Secret, error) {
	defaults := map[string]interface{}{
		"common_name": "intermediary-test.com",
		"key_type":    "ec",
		"issuer_name": "intermediary",
	}

	return vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/intermediate/generate/internal", mergeWithDefaults(props, defaults))
}

func (vpm *VaultPkiMount) SignIntermediary(signingIssuer string, csr interface{}, props map[string]interface{}) (*api.Secret, error) {
	defaults := map[string]interface{}{
		"csr": csr,
	}

	return vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/issuer/"+signingIssuer+"/sign-intermediate",
		mergeWithDefaults(props, defaults))
}

func (vpm *VaultPkiMount) ImportBundle(pemBundle interface{}, props map[string]interface{}) (*api.Secret, error) {
	defaults := map[string]interface{}{
		"pem_bundle": pemBundle,
	}

	return vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/issuers/import/bundle", mergeWithDefaults(props, defaults))
}

func (vpm *VaultPkiMount) UpdateDefaultIssuer(issuerId string, props map[string]interface{}) error {
	defaults := map[string]interface{}{
		"default": issuerId,
	}

	_, err := vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/config/issuers", mergeWithDefaults(props, defaults))

	return err
}

func (vpm *VaultPkiMount) UpdateIssuer(issuerRef string, props map[string]interface{}) error {
	defaults := map[string]interface{}{}

	_, err := vpm.GetActiveNode().Logical().JSONMergePatch(context.Background(),
		vpm.mount+"/issuer/"+issuerRef, mergeWithDefaults(props, defaults))

	return err
}

func (vpm *VaultPkiMount) UpdateRole(roleName string, config map[string]interface{}) error {
	defaults := map[string]interface{}{}

	_, err := vpm.GetActiveNode().Logical().WriteWithContext(context.Background(),
		vpm.mount+"/roles/"+roleName, mergeWithDefaults(config, defaults))

	return err
}

func (vpm *VaultPkiMount) GetEabKey(acmeDirectory string) (string, string, error) {
	eabPath := path.Join(vpm.mount, acmeDirectory, "/new-eab")
	resp, err := vpm.GetActiveNode().Logical().WriteWithContext(context.Background(), eabPath, map[string]interface{}{})
	if err != nil {
		return "", "", fmt.Errorf("failed fetching eab from %s: %w", eabPath, err)
	}
	eabId := resp.Data["id"].(string)
	base64EabKey := resp.Data["key"].(string)
	// just make sure we get something valid back from the server, we still want to pass back the base64 version
	// to the caller...
	_, err = base64.RawURLEncoding.DecodeString(base64EabKey)
	if err != nil {
		return "", "", fmt.Errorf("failed decoding key response field: %s: %w", base64EabKey, err)
	}
	return eabId, base64EabKey, nil
}

// GetCACertPEM retrieves the PKI mount's PEM-encoded CA certificate.
func (vpm *VaultPkiMount) GetCACertPEM() (string, error) {
	caCertPath := path.Join(vpm.mount, "/cert/ca")
	resp, err := vpm.GetActiveNode().Logical().ReadWithContext(context.Background(), caCertPath)
	if err != nil {
		return "", err
	}
	return resp.Data["certificate"].(string), nil
}

func mergeWithDefaults(config map[string]interface{}, defaults map[string]interface{}) map[string]interface{} {
	myConfig := config
	if myConfig == nil {
		myConfig = map[string]interface{}{}
	}
	for key, value := range defaults {
		if origVal, exists := config[key]; !exists {
			myConfig[key] = value
		} else {
			myConfig[key] = origVal
		}
	}

	return myConfig
}
