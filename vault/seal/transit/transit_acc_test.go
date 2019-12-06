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
	cleanup, retAddress, token, mountPath, keyName, _ := docker.PrepareTestContainer(t)
	defer cleanup()

	sealConfig := map[string]string{
		"address":    retAddress,
		"token":      token,
		"mount_path": mountPath,
		"key_name":   keyName,
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
	cleanup, retAddress, token, mountPath, keyName, tlsConfig := docker.PrepareTestContainer(t)
	defer cleanup()

	clientConfig := &api.Config{
		Address: retAddress,
	}
	if err := clientConfig.ConfigureTLS(tlsConfig); err != nil {
		t.Fatalf("err: %s", err)
	}

	remoteClient, err := api.NewClient(clientConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	remoteClient.SetToken(token)

	req := &api.TokenCreateRequest{
		Period: "5s",
	}
	rsp, err := remoteClient.Auth().Token().Create(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	sealConfig := map[string]string{
		"address":    retAddress,
		"token":      rsp.Auth.ClientToken,
		"mount_path": mountPath,
		"key_name":   keyName,
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
