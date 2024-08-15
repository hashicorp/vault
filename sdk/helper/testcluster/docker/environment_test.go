// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docker

import (
	"testing"
)

func TestSettingEnvsToContainer(t *testing.T) {
	expectedEnv := "TEST_ENV=value1"
	expectedEnv2 := "TEST_ENV2=value2"
	opts := &DockerClusterOptions{
		ImageRepo: "hashicorp/vault",
		ImageTag:  "latest",
		Envs:      []string{expectedEnv, expectedEnv2},
	}
	cluster := NewTestDockerCluster(t, opts)
	defer cluster.Cleanup()

	envs := cluster.GetActiveClusterNode().Container.Config.Env

	if !findEnv(envs, expectedEnv) {
		t.Errorf("Missing ENV variable: %s", expectedEnv)
	}
	if !findEnv(envs, expectedEnv2) {
		t.Errorf("Missing ENV variable: %s", expectedEnv2)
	}
}

func findEnv(envs []string, env string) bool {
	for _, e := range envs {
		if e == env {
			return true
		}
	}
	return false
}
