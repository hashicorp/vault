// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package seal_binary

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/hashicorp/vault/api"
	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"
)

const (
	testContainerConfig = `
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

	sealList = `
"seal": [
	%s
]
`

	transitSealParameters = `
  "address": "%s",
  "token": "%s",
  "mount_path": "%s",
  "key_name": "%s",
  "name": "%s"
`

	transitSealStanza = `
{
	"transit": {
	  %s,
	  "priority": %d,
	  "disabled": %s
	}
}
`
)

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
		c := &DockerVaultConfig{
			ServiceURL: *dockhelper.NewServiceURL(url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, port)}),
			tlsConfig: &api.TLSConfig{
				Insecure: true,
			},
		}
		return c, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start docker vault: %s", err)
	}

	return svc, runner, nil
}

func createTransitTestContainer(imageRepo, imageTag string) (func(), string, string, error) {
	rootToken, err := uuid.GenerateUUID()
	if err != nil {
		return nil, "", "", fmt.Errorf("err: %s", err)
	}
	testMountPath, err := uuid.GenerateUUID()
	if err != nil {
		return nil, "", "", fmt.Errorf("err: %s", err)
	}
	firstTestKeyName, err := uuid.GenerateUUID()
	if err != nil {
		return nil, "", "", fmt.Errorf("err: %s", err)
	}
	secondTransitKeyName, err := uuid.GenerateUUID()
	if err != nil {
		return nil, "", "", fmt.Errorf("err: %s", err)
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
		return nil, "", "", fmt.Errorf("could not create runner: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (dockhelper.ServiceConfig, error) {
		c := &DockerVaultConfig{
			ServiceURL: *dockhelper.NewServiceURL(url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, port)}),
			tlsConfig: &api.TLSConfig{
				Insecure: true,
			},
		}
		clientConfig := api.DefaultConfig()
		clientConfig.Address = c.ServiceURL.URL().String()
		vault, err := api.NewClient(clientConfig)
		if err != nil {
			return nil, err
		}
		vault.SetToken(rootToken)

		// Set up transit
		if err := vault.Sys().Mount(testMountPath, &api.MountInput{
			Type: "transit",
		}); err != nil {
			return nil, err
		}

		// Create two transit keys
		if _, err := vault.Logical().Write(path.Join(testMountPath, "keys", firstTestKeyName), map[string]interface{}{}); err != nil {
			return nil, err
		}

		if _, err := vault.Logical().Write(path.Join(testMountPath, "keys", secondTransitKeyName), map[string]interface{}{}); err != nil {
			return nil, fmt.Errorf("error creating transit key: %s", err)
		}

		return c, nil
	})
	if err != nil {
		return nil, "", "", fmt.Errorf("could not start docker vault: %s", err)
	}

	mapping, err := runner.GetNetworkAndAddresses(svc.Container.Name)
	if err != nil {
		svc.Cleanup()
		return nil, "", "", fmt.Errorf("failed to get container network information: %s", err)
	}

	if len(mapping) != 1 {
		svc.Cleanup()
		return nil, "", "", fmt.Errorf("expected 1 network mapping, got %d", len(mapping))
	}

	var ip string
	for _, ip = range mapping {
		// capture the container IP address from the map
	}

	return svc.Cleanup,
		fmt.Sprintf(transitSealParameters,
			fmt.Sprintf("http://%s:8200", ip),
			rootToken,
			testMountPath,
			firstTestKeyName,
			"transit-seal-1",
		),
		fmt.Sprintf(transitSealParameters,
			fmt.Sprintf("http://%s:8200", ip),
			rootToken,
			testMountPath,
			firstTestKeyName,
			"transit-seal-2",
		), nil
}

func checkVaultSealType(client *api.Client, expectedSealType string) error {
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

type DockerVaultConfig struct {
	dockhelper.ServiceURL
	token     string
	mountPath string
	keyName   string
	tlsConfig *api.TLSConfig
}

func (c *DockerVaultConfig) apiConfig() *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = c.URL().String()
	if err := vaultConfig.ConfigureTLS(c.tlsConfig); err != nil {
		panic("unable to configure TLS")
	}

	return vaultConfig
}
