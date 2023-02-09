package http

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
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

	stop := atomic.Bool{}

	eventType := "abc"

	// send some events
	go func() {
		for !stop.Load() {
			id, err := uuid.GenerateUUID()
			if err != nil {
				core.Logger().Info("Error generating UUID, exiting sender", "error", err)
			}
			err = core.Events().SendInternal(namespace.RootContext(context.Background()), namespace.RootNamespace, nil, logical.EventType(eventType), &logical.EventData{
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

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancelFunc)

	wsAddr := strings.Replace(addr, "http", "ws", 1)
	conn, _, err := websocket.Dial(ctx, wsAddr+"/v1/sys/events/subscribe/"+eventType+"?json=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, msg, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	msgJson := strings.TrimSpace(string(msg))
	if !strings.HasPrefix(msgJson, "{") || !strings.HasSuffix(msgJson, "}") {
		t.Errorf("Expected to get JSON event but got: %v", msgJson)
	}
}
