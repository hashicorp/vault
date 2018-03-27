package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical/inmem"
	log "github.com/mgutz/logxi/v1"
)

func TestConfig_Enabled(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	phys, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	logl := &logical.InmemStorage{}

	config := NewUIConfig(true, phys, logl)
	if !config.Enabled() {
		t.Fatal("ui should be enabled")
	}

	config = NewUIConfig(false, phys, logl)
	if config.Enabled() {
		t.Fatal("ui should not be enabled")
	}
}

func TestConfig_Headers(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	phys, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	logl := &logical.InmemStorage{}

	config := NewUIConfig(true, phys, logl)
	headers, err := config.Headers(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(headers) != len(staticHeaders) {
		t.Fatalf("expected %d headers, got %d", len(staticHeaders), len(headers))
	}

	head, err := config.GetHeader(context.Background(), "Test-Header")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if head != "" {
		t.Fatal("header returned found, should not be found")
	}
	err = config.SetHeader(context.Background(), "Test-Header", "123")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	head, err = config.GetHeader(context.Background(), "Test-Header")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if head == "" {
		t.Fatal("header not found when it should be")
	}
	if head != "123" {
		t.Fatalf("expected: %s, got: %s", "123", head)
	}

	head, err = config.GetHeader(context.Background(), "tEST-hEADER")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if head == "" {
		t.Fatal("header not found when it should be")
	}
	if head != "123" {
		t.Fatalf("expected: %s, got: %s", "123", head)
	}

	keys, err := config.HeaderKeys(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}

	err = config.SetHeader(context.Background(), "Test-Header-2", "321")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, err = config.HeaderKeys(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	err = config.DeleteHeader(context.Background(), "Test-Header-2")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = config.DeleteHeader(context.Background(), "Test-Header")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	head, err = config.GetHeader(context.Background(), "Test-Header")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if head != "" {
		t.Fatal("header returned found, should not be found")
	}
	keys, err = config.HeaderKeys(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("expected 0 key, got %d", len(keys))
	}
}
