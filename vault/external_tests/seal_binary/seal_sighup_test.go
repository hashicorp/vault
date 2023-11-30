// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package seal_binary

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"
)

func TestSealReloadSIGHUP(t *testing.T) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test with $VAULT_BINARY present")
	}

	cleanup, firstTransitKeyConfig, secondTransitKeyConfig, err := createTransitTestContainer("hasicorp/vault", "latest")
	if err != nil {
		t.Fatalf("error creating vault container: %s", err)
	}
	defer cleanup()

	testCases := map[string]struct {
		sealStanzas       []string
		expectedSealTypes []string
	}{
		"migrate transit to transit": {
			sealStanzas: []string{
				fmt.Sprintf(transitSealStanza, firstTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitSealStanza, firstTransitKeyConfig, 2, "true") + "," +
					fmt.Sprintf(transitSealStanza, secondTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitSealStanza, secondTransitKeyConfig, 1, "false"),
			},
			expectedSealTypes: []string{
				"transit",
				"transit",
				"transit",
			},
		},
		"migrate shamir to transit fails": {
			sealStanzas: []string{
				"",
				fmt.Sprintf(transitSealStanza, firstTransitKeyConfig, 1, "false"),
			},
			expectedSealTypes: []string{
				"shamir",
				"shamir",
			},
		},
		"migrate transit to shamir fails": {
			sealStanzas: []string{
				fmt.Sprintf(transitSealStanza, firstTransitKeyConfig, 1, "false"),
				"",
			},
			expectedSealTypes: []string{
				"transit",
				"transit",
			},
		},
		"replacing seal fails": {
			sealStanzas: []string{
				fmt.Sprintf(transitSealStanza, firstTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitSealStanza, secondTransitKeyConfig, 1, "false"),
			},
			expectedSealTypes: []string{
				"transit",
				"transit",
			},
		},
		"more than one seal fails": {
			sealStanzas: []string{
				fmt.Sprintf(transitSealStanza, firstTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitSealStanza, firstTransitKeyConfig, 1, "false") + "," +
					fmt.Sprintf(transitSealStanza, secondTransitKeyConfig, 2, "false"),
			},
			expectedSealTypes: []string{
				"transit",
				"transit",
			},
		},
	}

	err = createDockerImage("hashicorp/vault", "test-image", os.Getenv("VAULT_BINARY"))
	if err != nil {
		t.Fatalf("error creating docker image: %s", err)
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			var sealConfig string
			if test.sealStanzas[0] != "" {
				sealConfig = fmt.Sprintf(sealList, test.sealStanzas[0])
			}

			vaultConfig := fmt.Sprintf(testContainerConfig, sealConfig)

			svc, runner, err := createContainerWithConfig(vaultConfig, "hashicorp/vault", "test-image", func(s string) { t.Log(s) })
			defer svc.Cleanup()
			if err != nil {
				t.Fatalf("error creating container: %s", err)
			}

			time.Sleep(5 * time.Second)

			clientConfig := api.DefaultConfig()
			clientConfig.Address = svc.Config.URL().String()
			testClient, err := api.NewClient(clientConfig)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if test.expectedSealTypes[0] == "shamir" {
				initResp, err := testClient.Sys().Init(&api.InitRequest{
					SecretThreshold: 1,
					SecretShares:    1,
				})
				if err != nil {
					t.Fatalf("error initializing vault: %s", err)
				}

				_, err = testClient.Sys().Unseal(initResp.Keys[0])
				if err != nil {
					t.Fatalf("error unsealing vault: %s", err)
				}
			} else {
				_, err = testClient.Sys().Init(&api.InitRequest{
					RecoveryShares:    1,
					RecoveryThreshold: 1,
				})
				if err != nil {
					t.Fatalf("error initializing vault: %s", err)
				}
			}

			for i := range test.sealStanzas {
				if test.sealStanzas[i] != "" {
					sealConfig = fmt.Sprintf(sealList, test.sealStanzas[i])
				}

				vaultConfig = fmt.Sprintf(testContainerConfig, sealConfig)

				bCtx := dockhelper.NewBuildContext()
				bCtx["local.json"] = &dockhelper.FileContents{
					Data: []byte(vaultConfig),
					Mode: 0o644,
				}

				tar, err := bCtx.ToTarball()
				if err != nil {
					t.Fatalf("error creating config tarball: %s", err)
				}

				err = runner.DockerAPI.CopyToContainer(context.Background(), svc.Container.ID, "/vault/config", tar, types.CopyToContainerOptions{})
				if err != nil {
					t.Fatalf("error copying config to container: %s", err)
				}

				err = runner.DockerAPI.ContainerKill(context.Background(), svc.Container.ID, "SIGHUP")
				if err != nil {
					t.Fatalf("error sending SIGHUP: %s", err)
				}

				err = checkVaultSealType(testClient, test.expectedSealTypes[i])
				if err != nil {
					t.Fatalf("seal type check failed: %s", err)
				}
			}
		})
	}
}
