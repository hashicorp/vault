// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/docker"
)

func TestTransitWrapper_Lifecycle(t *testing.T) {
	cleanup, config := prepareTestContainer(t)
	defer cleanup()

	wrapperConfig := map[string]string{
		"address":    config.URL().String(),
		"token":      config.token,
		"mount_path": config.mountPath,
		"key_name":   config.keyName,
	}

	kms, _, err := configutil.GetTransitKMSFunc(&configutil.KMS{Config: wrapperConfig})
	if err != nil {
		t.Fatalf("error setting wrapper config: %v", err)
	}

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := kms.Encrypt(context.Background(), input, nil)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := kms.Decrypt(context.Background(), swi, nil)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}

func TestTransitSeal_TokenRenewal(t *testing.T) {
	cleanup, config := prepareTestContainer(t)
	defer cleanup()

	remoteClient, err := api.NewClient(config.apiConfig())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	remoteClient.SetToken(config.token)

	req := &api.TokenCreateRequest{
		Period: "5s",
	}
	rsp, err := remoteClient.Auth().Token().Create(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	wrapperConfig := map[string]string{
		"address":    config.URL().String(),
		"token":      rsp.Auth.ClientToken,
		"mount_path": config.mountPath,
		"key_name":   config.keyName,
	}
	kms, _, err := configutil.GetTransitKMSFunc(&configutil.KMS{Config: wrapperConfig})
	if err != nil {
		t.Fatalf("error setting wrapper config: %v", err)
	}

	time.Sleep(7 * time.Second)

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := kms.Encrypt(context.Background(), input, nil)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := kms.Decrypt(context.Background(), swi, nil)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}

type DockerVaultConfig struct {
	docker.ServiceURL
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

var _ docker.ServiceConfig = &DockerVaultConfig{}

func prepareTestContainer(t *testing.T) (func(), *DockerVaultConfig) {
	rootToken, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testMountPath, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testKeyName, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName: "vault",
		ImageRepo:     "docker.mirror.hashicorp.services/hashicorp/vault",
		ImageTag:      "latest",
		Cmd: []string{
			"server", "-log-level=trace", "-dev", fmt.Sprintf("-dev-root-token-id=%s", rootToken),
			"-dev-listen-address=0.0.0.0:8200",
		},
		Ports: []string{"8200/tcp"},
	})
	if err != nil {
		t.Fatalf("could not start docker vault: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		c := &DockerVaultConfig{
			ServiceURL: *docker.NewServiceURL(url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, port)}),
			tlsConfig: &api.TLSConfig{
				Insecure: true,
			},
			token:     rootToken,
			mountPath: testMountPath,
			keyName:   testKeyName,
		}
		vault, err := api.NewClient(c.apiConfig())
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

		// Create default aesgcm key
		if _, err := vault.Logical().Write(path.Join(testMountPath, "keys", testKeyName), map[string]interface{}{}); err != nil {
			return nil, err
		}

		return c, nil
	})
	if err != nil {
		t.Fatalf("could not start docker vault: %s", err)
	}
	return svc.Cleanup, svc.Config.(*DockerVaultConfig)
}
