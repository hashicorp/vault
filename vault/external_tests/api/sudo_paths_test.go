// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const sudoKey = "x-vault-sudo"

// Tests that the static list of sudo paths in the api package matches what's in the current OpenAPI spec.
func TestSudoPaths(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		DisableMlock:        true,
		DisableCache:        true,
		EnableRaw:           true,
		EnableIntrospection: true,
		Logger:              log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"database":       database.Factory,
			"generic-leased": vault.LeasedPassthroughBackendFactory,
			"pki":            pki.Factory,
			"transit":        transit.Factory,
		},
		BuiltinRegistry: builtinplugins.Registry,
	}
	client, _, closer := testVaultServerCoreConfig(t, coreConfig)
	defer closer()

	for credBackendName := range coreConfig.CredentialBackends {
		err := client.Sys().EnableAuthWithOptions(credBackendName, &api.EnableAuthOptions{
			Type: credBackendName,
		})
		if err != nil {
			t.Fatalf("error enabling auth backend for test: %v", err)
		}
	}

	for logicalBackendName := range coreConfig.LogicalBackends {
		err := client.Sys().Mount(logicalBackendName, &api.MountInput{
			Type: logicalBackendName,
		})
		if err != nil {
			t.Fatalf("error enabling logical backend for test: %v", err)
		}
	}

	sudoPathsFromSpec, err := getSudoPathsFromSpec(client)
	if err != nil {
		t.Fatalf("error getting list of paths that require sudo from OpenAPI endpoint: %v", err)
	}

	sudoPathsInCode := api.SudoPaths()

	// check for missing paths
	for path := range sudoPathsFromSpec {
		if _, ok := sudoPathsInCode[path]; !ok {
			t.Fatalf(
				"A path in the OpenAPI spec is missing from the static list of "+
					"sudo paths in the api module (%s). Please reconcile the two "+
					"accordingly.", path)
		}
	}
}

func getSudoPathsFromSpec(client *api.Client) (map[string]struct{}, error) {
	r := client.NewRequest("GET", "/v1/sys/internal/specs/openapi")
	resp, err := client.RawRequest(r)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve sudo endpoints: %v", err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	oasInfo := make(map[string]interface{})
	if err := jsonutil.DecodeJSONFromReader(resp.Body, &oasInfo); err != nil {
		return nil, fmt.Errorf("unable to decode JSON from OpenAPI response: %v", err)
	}

	paths, ok := oasInfo["paths"]
	if !ok {
		return nil, fmt.Errorf("OpenAPI response did not include paths")
	}

	pathsMap, ok := paths.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("OpenAPI response did not return valid paths")
	}

	sudoPaths := make(map[string]struct{})
	for pathName, pathInfo := range pathsMap {
		pathInfoMap, ok := pathInfo.(map[string]interface{})
		if !ok {
			continue
		}

		if sudo, ok := pathInfoMap[sudoKey]; ok {
			if sudo == true {
				sudoPaths[pathName] = struct{}{}
			}
		}
	}

	return sudoPaths, nil
}
