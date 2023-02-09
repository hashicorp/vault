package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"google.golang.org/protobuf/encoding/protojson"
	"nhooyr.io/websocket"
)

// handleEventsSubscribeWebsocket runs forever, returning a websocket error code and reason
// only if the connection closes or there was an error.
func handleEventsSubscribeWebsocket(ctx context.Context, core *vault.Core, ns *namespace.Namespace, eventType logical.EventType, conn *websocket.Conn, json bool) (websocket.StatusCode, string, error) {
	events := core.Events()
	ch, cancel, err := events.Subscribe(ctx, ns, eventType)
	if err != nil {
		core.Logger().Info("Error subscribing", "error", err)
		return websocket.StatusUnsupportedData, "Error subscribing", nil
	}
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			core.Logger().Info("Websocket context is done, closing the connection")
			return websocket.StatusNormalClosure, "", nil
		case message := <-ch:
			core.Logger().Debug("Sending message to websocket", "message", message)
			var messageBytes []byte
			if json {
				messageBytes, err = protojson.Marshal(message)
			} else {
				messageBytes, err = proto.Marshal(message)
			}
			if err != nil {
				core.Logger().Warn("Could not serialize websocket event", "error", err)
				return 0, "", err
			}
			messageString := string(messageBytes) + "\n"
			err = conn.Write(ctx, websocket.MessageText, []byte(messageString))
			if err != nil {
				return 0, "", err
			}
		}
	}
}

func handleEventsSubscribe(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		core.Logger().Debug("Got request to", "url", r.URL, "version", r.Proto)

		ctx := r.Context()
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			core.Logger().Info("Could not find namespace", "error", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("could not find namespace"))
			return
		}

		prefix := "/v1/sys/events/subscribe/"
		if ns.ID != "root" {
			prefix = fmt.Sprintf("/v1/%s/sys/events/subscribe/", ns.Path)
		}
		eventTypeStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
		if eventTypeStr == "" {
			respondError(w, http.StatusBadRequest, fmt.Errorf("did not specify eventType to subscribe to"))
			return
		}
		eventType := logical.EventType(eventTypeStr)

		json := false
		jsonRaw := r.URL.Query().Get("json")
		if jsonRaw != "" {
			var err error
			json, err = strconv.ParseBool(jsonRaw)
			if err != nil {
				respondError(w, http.StatusBadRequest, fmt.Errorf("invalid parameter for JSON: %v", jsonRaw))
				return
			}
		}

		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			core.Logger().Info("Could not accept as websocket", "error", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("could not accept as websocket"))
			return
		}

		// we don't expect any incoming messages
		ctx = conn.CloseRead(ctx)
		// start the pinger
		go func() {
			for {
				time.Sleep(30 * time.Second) // not too aggressive, but keep the HTTP connection alive
				err := conn.Ping(ctx)
				if err != nil {
					return
				}
			}
		}()

		closeStatus, closeReason, err := handleEventsSubscribeWebsocket(ctx, core, ns, eventType, conn, json)
		if err != nil {
			closeStatus = websocket.CloseStatus(err)
			if closeStatus == -1 {
				closeStatus = websocket.StatusInternalError
			}
			closeReason = fmt.Sprintf("Internal error: %v", err)
			core.Logger().Debug("Error from websocket handler", "error", err)
		}
		// Close() will panic if the reason is greater than this length
		if len(closeReason) > 123 {
			core.Logger().Debug("Truncated close reason", "closeReason", closeReason)
			closeReason = closeReason[:123]
		}
		err = conn.Close(closeStatus, closeReason)
		if err != nil {
			core.Logger().Debug("Error closing websocket", "error", err)
		}
	})
}
