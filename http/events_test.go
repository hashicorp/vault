package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"nhooyr.io/websocket"
)

// TestEventsSubscribe tests the websocket endpoint for subscribing to events
// by generating some events.
func TestEventsSubscribe(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	// unseal the core
	keys, token := vault.TestCoreInit(t, core)
	for _, key := range keys {
		_, err := core.Unseal(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	stop := atomic.Bool{}

	const eventType = "abc"

	// send some events
	go func() {
		for !stop.Load() {
			id, err := uuid.GenerateUUID()
			if err != nil {
				core.Logger().Info("Error generating UUID, exiting sender", "error", err)
			}
			pluginInfo := &logical.EventPluginInfo{
				MountPath: "secret",
			}
			err = core.Events().SendInternal(namespace.RootContext(context.Background()), namespace.RootNamespace, pluginInfo, logical.EventType(eventType), &logical.EventData{
				Id:        id,
				Metadata:  nil,
				EntityIds: nil,
				Note:      "testing",
			})
			if err != nil {
				core.Logger().Info("Error sending event, exiting sender", "error", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	t.Cleanup(func() {
		stop.Store(true)
	})

	ctx := context.Background()
	wsAddr := strings.Replace(addr, "http", "ws", 1)

	conn, _, err := websocket.Dial(ctx, wsAddr+"/v1/sys/events/subscribe/"+eventType+"?json=true", &websocket.DialOptions{
		HTTPHeader: http.Header{"x-vault-token": []string{token}},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, msg, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(msg, &data)
	if err != nil {
		t.Fatal(err)
	}
	if actualType := data["event_type"].(string); actualType != eventType {
		t.Fatalf("Expeced event type %s, got %s", eventType, actualType)
	}
	pluginInfo, ok := data["plugin_info"].(map[string]interface{})
	if !ok || pluginInfo == nil {
		t.Fatalf("No plugin_info object: %v", data)
	}
	mountPath, ok := pluginInfo["mount_path"].(string)
	if !ok || mountPath != "secret" {
		t.Fatalf("Wrong mount_path: %v", data)
	}
}

func TestEventsSubscribeAuth(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	// unseal the core
	keys, root := vault.TestCoreInit(t, core)
	for _, key := range keys {
		_, err := core.Unseal(key)
		if err != nil {
			t.Fatal(err)
		}
	}

	var nonPrivilegedToken string
	// Fetch a valid non privileged token.
	{
		config := api.DefaultConfig()
		config.Address = addr

		client, err := api.NewClient(config)
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(root)

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{Policies: []string{"default"}})
		if err != nil {
			t.Fatal(err)
		}
		if secret.Auth.ClientToken == "" {
			t.Fatal("Failed to fetch a non privileged token")
		}
		nonPrivilegedToken = secret.Auth.ClientToken
	}

	ctx := context.Background()
	wsAddr := strings.Replace(addr, "http", "ws", 1)

	// Get a 403 with no token.
	_, resp, err := websocket.Dial(ctx, wsAddr+"/v1/sys/events/subscribe/abc", nil)
	if err == nil {
		t.Error("Expected websocket error but got none")
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %+v", resp)
	}

	// Get a 403 with a non privileged token.
	_, resp, err = websocket.Dial(ctx, wsAddr+"/v1/sys/events/subscribe/abc", &websocket.DialOptions{
		HTTPHeader: http.Header{"x-vault-token": []string{nonPrivilegedToken}},
	})
	if err == nil {
		t.Error("Expected websocket error but got none")
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %+v", resp)
	}
}
