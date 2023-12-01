// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal_binary

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"
)

const (
	containerConfig = `
{
	"storage": {
		"file": {
			"path": "/tmp",
		}
	},

	"disable_mlock": true,

	"listener": [{
		"tcp": {
			"address": "0.0.0.0:8200",
			"tls_disable": "true"
		}
	}],

	"api_addr": "http://0.0.0.0:8200",
	"cluster_addr": "http://0.0.0.0:8201",
	%s
}`

	sealConfig = `
"seal": [
	%s
]
`

	transitParameters = `
  "address": "%s",
  "token": "%s",
  "mount_path": "%s",
  "key_name": "%s",
  "name": "%s"
`

	transitStanza = `
{
	"transit": {
	  %s,
	  "priority": %d,
	  "disabled": %s
	}
}
`
)

type transitContainerConfig struct {
	Address    string
	Token      string
	MountPaths []string
	KeyNames   []string
}

func createDockerImage(imageRepo, imageTag, vaultBinary string) error {
	runner, err := dockhelper.NewServiceRunner(dockhelper.RunOptions{
		ContainerName: "vault",
		ImageRepo:     imageRepo,
		ImageTag:      "latest",
	})
	if err != nil {
		return fmt.Errorf("error creating runner: %s", err)
	}

	f, err := os.Open(vaultBinary)
	if err != nil {
		return fmt.Errorf("error opening vault binary file: %s", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading vault binary file: %s", err)
	}
	bCtx := dockhelper.NewBuildContext()
	bCtx["vault"] = &dockhelper.FileContents{
		Data: data,
		Mode: 0o755,
	}

	containerFile := fmt.Sprintf(`
FROM %s:latest
COPY vault /bin/vault
`, imageRepo)

	_, err = runner.BuildImage(context.Background(), containerFile, bCtx,
		dockhelper.BuildRemove(true), dockhelper.BuildForceRemove(true),
		dockhelper.BuildPullParent(true),
		dockhelper.BuildTags([]string{fmt.Sprintf("hashicorp/vault:%s", imageTag)}))
	if err != nil {
		return fmt.Errorf("error building docker image: %s", err)
	}

	return nil
}

func createContainerWithConfig(config string, imageRepo, imageTag string, logConsumer func(s string)) (*dockhelper.Service, *dockhelper.Runner, error) {
	runner, err := dockhelper.NewServiceRunner(dockhelper.RunOptions{
		ContainerName: "vault",
		ImageRepo:     imageRepo,
		ImageTag:      imageTag,
		Cmd: []string{
			"server", "-log-level=trace",
		},
		Ports:       []string{"8200/tcp"},
		Env:         []string{fmt.Sprintf("VAULT_LICENSE=%s", os.Getenv("VAULT_LICENSE")), fmt.Sprintf("VAULT_LOCAL_CONFIG=%s", config)},
		LogConsumer: logConsumer,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating runner: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (dockhelper.ServiceConfig, error) {
		return *dockhelper.NewServiceURL(url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, port)}), nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start docker vault: %s", err)
	}

	return svc, runner, nil
}

func createTransitTestContainer(imageRepo, imageTag string, numKeys int) (*dockhelper.Service, *transitContainerConfig, error) {
	rootToken, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, fmt.Errorf("err: %s", err)
	}

	mountPaths := make([]string, numKeys)
	keyNames := make([]string, numKeys)

	for i := range mountPaths {
		mountPaths[i], err = uuid.GenerateUUID()
		if err != nil {
			return nil, nil, fmt.Errorf("error generating UUID: %s", err)
		}

		keyNames[i], err = uuid.GenerateUUID()
		if err != nil {
			return nil, nil, fmt.Errorf("error generating UUID: %s", err)
		}
	}

	runner, err := dockhelper.NewServiceRunner(dockhelper.RunOptions{
		ContainerName: "vault",
		ImageRepo:     imageRepo,
		ImageTag:      imageTag,
		Cmd: []string{
			"server", "-log-level=trace", "-dev", fmt.Sprintf("-dev-root-token-id=%s", rootToken),
			"-dev-listen-address=0.0.0.0:8200",
		},
		Env:   []string{fmt.Sprintf("VAULT_LICENSE=%s", os.Getenv("VAULT_LICENSE"))},
		Ports: []string{"8200/tcp"},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not create runner: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (dockhelper.ServiceConfig, error) {
		c := *dockhelper.NewServiceURL(url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, port)})

		clientConfig := api.DefaultConfig()
		clientConfig.Address = c.URL().String()
		vault, err := api.NewClient(clientConfig)
		if err != nil {
			return nil, err
		}
		vault.SetToken(rootToken)

		// Set up transit mounts and keys
		for i := range mountPaths {
			if err := vault.Sys().Mount(mountPaths[i], &api.MountInput{
				Type: "transit",
			}); err != nil {
				return nil, err
			}

			if _, err := vault.Logical().Write(path.Join(mountPaths[i], "keys", keyNames[i]), map[string]interface{}{}); err != nil {
				return nil, err
			}
		}

		return c, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start docker vault: %s", err)
	}

	mapping, err := runner.GetNetworkAndAddresses(svc.Container.Name)
	if err != nil {
		svc.Cleanup()
		return nil, nil, fmt.Errorf("failed to get container network information: %s", err)
	}

	if len(mapping) != 1 {
		svc.Cleanup()
		return nil, nil, fmt.Errorf("expected 1 network mapping, got %d", len(mapping))
	}

	var ip string
	for _, ip = range mapping {
		// capture the container IP address from the map
	}

	return svc,
		&transitContainerConfig{
			Address:    fmt.Sprintf("http://%s:8200", ip),
			Token:      rootToken,
			MountPaths: mountPaths,
			KeyNames:   keyNames,
		}, nil
}

func validateVaultStatusAndSealType(client *api.Client, expectedSealType string) error {
	statusResp, err := client.Sys().SealStatus()
	if err != nil {
		return fmt.Errorf("error getting vault status: %s", err)
	}

	if statusResp.Sealed {
		return fmt.Errorf("expected vault to be unsealed, but it is sealed")
	}

	if statusResp.Type != expectedSealType {
		return fmt.Errorf("unexpected seal type: expected transit, got %s", statusResp.Type)
	}

	return nil
}
