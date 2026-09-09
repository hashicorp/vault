// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package seal_binary

import (
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/constants"
	sealhelper "github.com/hashicorp/vault/helper/testhelpers/seal"
	"github.com/hashicorp/vault/helper/testhelpers/testimages"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSealReloadSIGHUP(t *testing.T) {
	transit := sealhelper.NewTransitDockerSealServer(t)

	repo, tag := testimages.GetImageRepoAndTag(t, constants.IsEnterprise)

	type testCase struct {
		name             string
		steps            []step
		disableMultiseal bool
	}

	testCases := []testCase{
		{
			name: "transit to transit",
			steps: []step{
				{
					"transit", []seal{
						{base: transit.Seal, index: 0, priority: 1},
					},
				}, {
					"multiseal", []seal{
						{base: transit.Seal, index: 0, priority: 2, disabled: true},
						{base: transit.Seal, index: 1, priority: 1},
					},
				}, {
					"transit", []seal{
						{base: transit.Seal, index: 1, priority: 1},
					},
				},
			},
		}, {
			name: "transit to transit no multiseal",
			steps: []step{
				{
					"transit", []seal{
						{base: transit.Seal, index: 0, priority: 1},
					},
				}, {
					"transit", []seal{
						{base: transit.Seal, index: 0, priority: 2, disabled: true},
						{base: transit.Seal, index: 1, priority: 1},
					},
				}, {
					"transit", []seal{
						{base: transit.Seal, index: 1, priority: 1},
					},
				},
			},
			disableMultiseal: true,
		}, {
			name: "transit to pkcs11",
			steps: []step{
				{
					"transit", []seal{
						{base: transit.Seal, priority: 1},
					},
				}, {
					"multiseal", []seal{
						{base: transit.Seal, priority: 1, disabled: true},
						{base: pkcsWrapper, priority: 2},
					},
				}, {
					"pkcs11", []seal{
						{base: pkcsWrapper, priority: 1},
					},
				},
			},
		}, {
			name: "pkcs11 to transit",
			steps: []step{
				{
					"pkcs11", []seal{
						{base: pkcsWrapper, priority: 1},
					},
				}, {
					"multiseal", []seal{
						{base: pkcsWrapper, priority: 2, disabled: true},
						{base: transit.Seal, priority: 1},
					},
				}, {
					"transit", []seal{
						{base: transit.Seal, priority: 1},
					},
				},
			},
		}, {
			name: "two transit seals",
			steps: []step{
				{
					"transit", []seal{
						{base: transit.Seal, priority: 1},
					},
				}, {
					"multiseal", []seal{
						{base: transit.Seal, index: 0, priority: 1},
						{base: transit.Seal, index: 1, priority: 2},
					},
				},
			},
		}, {
			name: "pkcs11 seal and transit seal",
			steps: []step{
				{
					"transit", []seal{
						{base: transit.Seal, priority: 1},
					},
				}, {
					"multiseal", []seal{
						{base: transit.Seal, priority: 1},
						{base: pkcsWrapper, priority: 2},
					},
				},
			},
		}, {
			name: "three seals",
			steps: []step{
				{
					"transit", []seal{
						{base: transit.Seal, index: 0, priority: 1},
					},
				}, {
					"multiseal", []seal{
						{base: transit.Seal, index: 0, priority: 1},
						{base: pkcsWrapper, priority: 2},
					},
				}, {
					"multiseal", []seal{
						{base: transit.Seal, index: 0, priority: 1},
						{base: pkcsWrapper, priority: 2},
						{base: transit.Seal, index: 1, priority: 3},
					},
				},
			},
		}, {
			name: "remove enabled seal",
			steps: []step{
				{
					"transit", []seal{
						{base: transit.Seal, priority: 1},
					},
				}, {
					"multiseal", []seal{
						{base: transit.Seal, priority: 1},
						{base: pkcsWrapper, priority: 2},
					},
				}, {
					"pkcs11", []seal{
						{base: pkcsWrapper, priority: 2},
					},
				},
			},
		}, {
			name: "shamir to transit fails",
			steps: []step{
				{
					"shamir", nil,
				}, {
					"shamir", []seal{
						{base: transit.Seal, priority: 1},
					},
				},
			},
		}, {
			name: "transit to shamir fails",
			steps: []step{
				{
					"transit", []seal{{base: transit.Seal, priority: 1}},
				}, {
					"transit", nil,
				},
			},
		}, {
			name: "replacing seal fails",
			steps: []step{
				{
					"transit", []seal{{base: transit.Seal, index: 0, priority: 1}},
				}, {
					"transit", []seal{{base: transit.Seal, index: 1, priority: 1}},
				},
			},
		},
	}

	isEnterpriseCase := func(tc testCase) bool {
		for _, step := range tc.steps {
			if step.expectedSealType == "multiseal" || step.expectedSealType == "pkcs11" {
				return true
			}
		}
		return false
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if isEnterpriseCase(tc) && !constants.IsEnterprise {
				t.Skip("Skipping enterprise tests")
			}

			opts := dockerOptions(t, repo, tag)
			for _, seal := range tc.steps[0].seals {
				vncseal := withPriorityAndDisabled(seal.priority, seal.disabled, seal.base(tc.name, seal.index))
				opts.VaultNodeConfig.Seal = append(opts.VaultNodeConfig.Seal, vncseal)
			}
			if tc.steps[0].expectedSealType != "shamir" && !tc.disableMultiseal {
				opts.VaultNodeConfig.EnableMultiSeal = true
			}
			cluster := docker.NewTestDockerCluster(t, opts)
			node := cluster.Nodes()[0].(*docker.DockerClusterNode)
			client := node.APIClient()
			lastRewrappedEntryCount, err := getRewrappedEntryCount(client)
			require.NoError(t, err)

			// kv mounts are sealwrapped.  In order to make sure that we don't get fooled
			// by the rewrap status endpoint saying "not in progress" prior to a rewrap
			// being started, we're going to arrange for there to be an extra key to wrap
			// each iteration, by creating a new kv entry each iteration.
			require.NoError(t, client.Sys().Mount("kv", &api.MountInput{
				Type: "kv",
			}))
			client.Logical().Write("kv/0", map[string]any{"1": 1})

			expectFailure := len(tc.steps) < 3
			for i := 1; i < len(tc.steps); i++ {
				if tc.steps[i].expectedSealType != "shamir" && !tc.disableMultiseal {
					opts.VaultNodeConfig.EnableMultiSeal = true
				}
				opts.VaultNodeConfig.Seal = nil
				for _, seal := range tc.steps[i].seals {
					opts.VaultNodeConfig.Seal = append(opts.VaultNodeConfig.Seal,
						withPriorityAndDisabled(seal.priority, seal.disabled, seal.base(tc.name, seal.index)))
				}
				require.NoError(t, node.UpdateConfig(t.Context(), opts))
				require.NoError(t, node.Signal(t.Context(), "SIGHUP"))

				require.EventuallyWithT(t, func(ct *assert.CollectT) {
					if !tc.disableMultiseal && !expectFailure && tc.steps[i].expectedSealType != "shamir" {
						lastRewrappedEntryCount = verifyRewrappedEntryCount(ct, client, lastRewrappedEntryCount+1)
					}

					resp, err := client.Sys().SealStatusWithContext(t.Context())
					require.NoError(t, err)
					assert.Equal(ct, resp.Type, tc.steps[i].expectedSealType)
					assert.False(ct, resp.Sealed)
				}, 20*time.Second, time.Second/2)

				client.Logical().Write("kv/"+strconv.Itoa(i), map[string]any{"1": 1})
			}
		})
	}
}
