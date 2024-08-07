// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const sudoKey = "x-vault-sudo"

// Tests that the static list of sudo paths in the api package matches what's in the current OpenAPI spec.
func TestSudoPaths(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		EnableRaw:           true,
		EnableIntrospection: true,
		BuiltinRegistry:     builtinplugins.Registry,
		LogicalBackends: map[string]logical.Factory{
			"pki": pki.Factory,
		},
	}

	cluster := minimal.NewTestSoloCluster(t, coreConfig)
	client := cluster.Cores[0].Client

	// NOTE: At present there are no auth methods with sudo paths, except for the
	// automatically mounted token backend. If this changes this test should be
	// updated to iterate over the auth methods and test them as we are doing below.

	// Each secrets engine that contains sudo paths (other than automatically mounted ones) must be mounted here
	for _, logicalBackendName := range []string{"pki"} {
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
		pathTrimmed := strings.TrimRight(path, "/")
		if _, ok := sudoPathsInCode[pathTrimmed]; !ok {
			t.Fatalf(
				"A path in the OpenAPI spec is missing from the static list of "+
					"sudo paths in the api module (%s). Please reconcile the two "+
					"accordingly.", pathTrimmed)
		}
	}

	// check for extra paths
	for path := range sudoPathsInCode {
		if _, ok := sudoPathsFromSpec[path]; !ok {
			if _, ok := sudoPathsFromSpec[path+"/"]; !ok {
				t.Fatalf(
					"A path in the static list of sudo paths in the api module "+
						"is not marked as a sudo path in the OpenAPI spec (%s). Please reconcile the two "+
						"accordingly.", path)
			}
		}
	}
}

func getSudoPathsFromSpec(client *api.Client) (map[string]struct{}, error) {
	resp, err := client.Logical().ReadRaw("sys/internal/specs/openapi")
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
