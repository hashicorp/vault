// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package seal_binary

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"
)

func TestSealReloadSIGHUP(t *testing.T) {
	binary := os.Getenv("VAULT_BINARY")
	if binary == "" {
		t.Skip("only running docker test with $VAULT_BINARY present")
	}

	transitContainer, transitConfig, err := createTransitTestContainer("hashicorp/vault", "latest", 2)
	if err != nil {
		t.Fatalf("error creating vault container: %s", err)
	}
	defer transitContainer.Cleanup()

	firstTransitKeyConfig := fmt.Sprintf(transitParameters,
		transitConfig.Address,
		transitConfig.Token,
		transitConfig.MountPaths[0],
		transitConfig.KeyNames[0],
		"transit-seal-1",
	)

	secondTransitKeyConfig := fmt.Sprintf(transitParameters,
		transitConfig.Address,
		transitConfig.Token,
		transitConfig.MountPaths[1],
		transitConfig.KeyNames[1],
		"transit-seal-2",
	)

	testCases := map[string]struct {
		sealStanzas       []string
		expectedSealTypes []string
	}{
		"migrate transit to transit": {
			sealStanzas: []string{
				fmt.Sprintf(transitStanza, firstTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitStanza, firstTransitKeyConfig, 2, "true") + "," +
					fmt.Sprintf(transitStanza, secondTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitStanza, secondTransitKeyConfig, 1, "false"),
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
				fmt.Sprintf(transitStanza, firstTransitKeyConfig, 1, "false"),
			},
			expectedSealTypes: []string{
				"shamir",
				"shamir",
			},
		},
		"migrate transit to shamir fails": {
			sealStanzas: []string{
				fmt.Sprintf(transitStanza, firstTransitKeyConfig, 1, "false"),
				"",
			},
			expectedSealTypes: []string{
				"transit",
				"transit",
			},
		},
		"replacing seal fails": {
			sealStanzas: []string{
				fmt.Sprintf(transitStanza, firstTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitStanza, secondTransitKeyConfig, 1, "false"),
			},
			expectedSealTypes: []string{
				"transit",
				"transit",
			},
		},
		"more than one seal fails": {
			sealStanzas: []string{
				fmt.Sprintf(transitStanza, firstTransitKeyConfig, 1, "false"),
				fmt.Sprintf(transitStanza, firstTransitKeyConfig, 1, "false") + "," +
					fmt.Sprintf(transitStanza, secondTransitKeyConfig, 2, "false"),
			},
			expectedSealTypes: []string{
				"transit",
				"transit",
			},
		},
	}

	containerFile := `
FROM hashicorp/vault:latest
COPY vault /bin/vault
`
	bCtx, err := createBuildContextWithBinary(os.Getenv("VAULT_BINARY"))
	if err != nil {
		t.Fatalf("error creating build context: %s", err)
	}
	err = createDockerImage("hashicorp/vault", "test-image", containerFile, bCtx)
	if err != nil {
		t.Fatalf("error creating docker image: %s", err)
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			var sealList string
			if test.sealStanzas[0] != "" {
				sealList = fmt.Sprintf(sealConfig, test.sealStanzas[0])
			}

			vaultConfig := fmt.Sprintf(containerConfig, sealList)

			svc, runner, err := createContainerWithConfig(vaultConfig, "hashicorp/vault", "test-image", func(s string) { t.Log(s) })
			if err != nil {
				t.Fatalf("error creating container: %s", err)
			}
			defer svc.Cleanup()

			time.Sleep(5 * time.Second)

			client, err := testClient(svc.Config.URL().String())
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			_, token, err := initializeVault(client, test.expectedSealTypes[0])
			if err != nil {
				t.Fatalf("error initializing vault: %s", err)
			}
			client.SetToken(token)

			for i := range test.sealStanzas {
				if test.sealStanzas[i] != "" {
					sealList = fmt.Sprintf(sealList, test.sealStanzas[i])
				}

				vaultConfig = fmt.Sprintf(containerConfig, sealList)
				configCtx := dockhelper.NewBuildContext()
				configCtx["local.json"] = &dockhelper.FileContents{
					Data: []byte(vaultConfig),
					Mode: 0o644,
				}

				err = copyConfigToContainer(svc.Container.ID, bCtx, runner)
				if err != nil {
					t.Fatalf("error copying over config file: %s", err)
				}

				err = runner.DockerAPI.ContainerKill(context.Background(), svc.Container.ID, "SIGHUP")
				if err != nil {
					t.Fatalf("error sending SIGHUP: %s", err)
				}

				err = validateVaultStatusAndSealType(client, test.expectedSealTypes[i])
				if err != nil {
					t.Fatalf("seal type check failed: %s", err)
				}
			}
		})
	}
}
