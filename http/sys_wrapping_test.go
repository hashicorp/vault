package http

import (
	"fmt"
	"net/http"
	"testing"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/vault"
)

// Test wrapping functionality
func TestHTTP_Wrapping(t *testing.T) {
	handler1 := http.NewServeMux()
	handler2 := http.NewServeMux()
	handler3 := http.NewServeMux()

	coreConfig := &vault.CoreConfig{}

	// Chicken-and-egg: Handler needs a core. So we create handlers first, then
	// add routes chained to a Handler-created handler.
	cores := vault.TestCluster(t, []http.Handler{handler1, handler2, handler3}, coreConfig, true)
	for _, core := range cores {
		defer core.CloseListeners()
	}
	handler1.Handle("/", Handler(cores[0].Core))
	handler2.Handle("/", Handler(cores[1].Core))
	handler3.Handle("/", Handler(cores[2].Core))

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	root := cores[0].Root

	transport := cleanhttp.DefaultTransport()
	transport.TLSClientConfig = cores[0].TLSConfig
	httpClient := &http.Client{
		Transport: transport,
	}
	addr := fmt.Sprintf("https://127.0.0.1:%d", cores[0].Listeners[0].Address.Port)
	config := api.DefaultConfig()
	config.Address = addr
	config.HttpClient = httpClient
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(root)

	// Write a value that we will use with wrapping for lookup
	_, err = client.Logical().Write("secret/foo", map[string]interface{}{
		"zip": "zap",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Set a wrapping lookup function for reads on that path
	client.SetWrappingLookupFunc(func(operation, path string) string {
		if operation == "GET" && path == "secret/foo" {
			return "5m"
		}

		return ""
	})

	// First test: basic things that should fail
	// Root token isn't a wrapping token
	_, err = client.Logical().Unwrap(root)
	if err == nil {
		t.Fatal("expected error")
	}
	// Root token isn't a wrapping token
	_, err = client.Logical().Write("sys/wrapping/lookup", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	// Not supplied
	_, err = client.Logical().Write("sys/wrapping/lookup", map[string]interface{}{
		"foo": "bar",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	// Nonexistent token isn't a wrapping token
	_, err = client.Logical().Write("sys/wrapping/lookup", map[string]interface{}{
		"foo": "bar",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Test lookup
	// Create a wrapping token
	secret, err := client.Logical.Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.WrapInfo == nil {
		t.Fatal("secret or wrap info is nil")
	}

}
