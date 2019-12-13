package transit_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault/seal/transit"
)

func TestTransitSeal_Lifecycle(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Skip()
	}
	vault := docker.PrepareTestVaultContainer(t)
	defer vault.Cleanup()
	vault.MountTransit(t)

	sealConfig := map[string]string{
		"address":    vault.RetAddress,
		"token":      vault.Token,
		"mount_path": vault.MountPath,
		"key_name":   vault.KeyName,
	}
	s := transit.NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(sealConfig)
	if err != nil {
		t.Fatalf("error setting seal config: %v", err)
	}

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := s.Decrypt(context.Background(), swi)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}

func TestTransitSeal_TokenRenewal(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Skip()
	}
	vault := docker.PrepareTestVaultContainer(t)
	defer vault.Cleanup()
	vault.MountTransit(t)

	clientConfig := &api.Config{
		Address: vault.RetAddress,
	}
	if err := clientConfig.ConfigureTLS(vault.TLSConfig); err != nil {
		t.Fatalf("err: %s", err)
	}

	remoteClient, err := api.NewClient(clientConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	remoteClient.SetToken(vault.Token)

	req := &api.TokenCreateRequest{
		Period: "5s",
	}
	rsp, err := remoteClient.Auth().Token().Create(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	sealConfig := map[string]string{
		"address":    vault.RetAddress,
		"token":      rsp.Auth.ClientToken,
		"mount_path": vault.MountPath,
		"key_name":   vault.KeyName,
	}
	s := transit.NewSeal(logging.NewVaultLogger(log.Trace))
	_, err = s.SetConfig(sealConfig)
	if err != nil {
		t.Fatalf("error setting seal config: %v", err)
	}

	time.Sleep(7 * time.Second)

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := s.Decrypt(context.Background(), swi)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}
