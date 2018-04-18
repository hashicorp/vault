package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical/inmem"
)

func TestConfig_Enabled(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
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
	logger := logging.NewVaultLogger(log.Trace)
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
	if len(headers) != len(config.defaultHeaders) {
		t.Fatalf("expected %d headers, got %d", len(config.defaultHeaders), len(headers))
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

func TestConfig_DefaultHeaders(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
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
	if len(headers) != len(config.defaultHeaders) {
		t.Fatalf("expected %d headers, got %d", len(config.defaultHeaders), len(headers))
	}

	headers, err = config.Headers(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defaultCSP := config.defaultHeaders.Get("Content-security-Policy")
	head := headers.Get("Content-Security-Policy")
	if head != defaultCSP {
		t.Fatalf("header does not match: expected %s, got %s", defaultCSP, head)
	}

	err = config.SetHeader(context.Background(), "Content-security-Policy", "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	headers, err = config.Headers(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	head = headers.Get("Content-Security-Policy")
	if head != "test" {
		t.Fatalf("header does not match: expected %s, got %s", "test", head)
	}

	err = config.DeleteHeader(context.Background(), "Content-Security-Policy")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	headers, err = config.Headers(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	head = headers.Get("Content-Security-Policy")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if head != defaultCSP {
		t.Fatalf("header does not match: expected %s, got %s", defaultCSP, head)
	}
}
