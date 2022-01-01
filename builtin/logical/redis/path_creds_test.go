package redis

import (
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_NoCreds(t *testing.T) {
	b, ctx, s, _, stop := getBackendAndSetConfig(t)
	defer stop()

	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/missing",
		Storage:   s,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("Expected error, got: %#v", resp)
	}
	if resp.Error().Error() != `no role named "missing" found` {
		t.Fatalf("Wrong error: %s", resp.Error())
	}
}

func TestBackend_Creds(t *testing.T) {
	b, ctx, s, addr, stop := getBackendAndSetConfig(t)
	defer stop()

	resp, err := b.HandleRequest(ctx, &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/test",
		Storage:   s,
		Data: map[string]interface{}{
			"rules": []string{"on", "allkeys", "+set"},
		},
	})
	if err != nil || resp.IsError() {
		t.Fatal("failed to create role")
	}

	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "creds/test",
		Storage:     s,
		DisplayName: "token",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}

	// Try to use the credentials
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: resp.Data["username"].(string),
		Password: resp.Data["password"].(string),
	})

	if _, err := client.Set(ctx, "hello", "world", 0).Result(); err != nil {
		t.Fatal(err)
	}

	if err := client.Get(ctx, "hello").Err(); err == nil {
		t.Fatal("Getting a value should have raised an error")
	} else if !strings.Contains(err.Error(), "NOPERM") {
		t.Fatalf("Wrong error: %s", err)
	}

	// Revoke the credentials
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.RevokeOperation,
		Storage:   s,
		Secret:    resp.Secret,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}

	if _, err := client.Set(ctx, "hello", "world", 0).Result(); err == nil {
		t.Fatal("Setting a key should have raised an error")
	}
}
