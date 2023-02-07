package http

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
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

var eventTypeRegex = regexp.MustCompile(`.*/events/subscribe/(.*)`)

// handleEventsSubscribeWebsocket runs forever, returning a websocket error code and reason
// only if the connection closes or there was an error.
func handleEventsSubscribeWebsocket(ctx context.Context, core *vault.Core, ns *namespace.Namespace, eventType logical.EventType, conn *websocket.Conn, json bool) (websocket.StatusCode, string, error) {
	events := core.Events()
	ch, err := events.Subscribe(ctx, ns, eventType)
	if err != nil {
		core.Logger().Info("Error subscribing", "error", err)
		return websocket.StatusUnsupportedData, "Error subscribing", nil
	}

	defer close(ch)

	for {
		select {
		case <-ctx.Done():
			core.Logger().Info("Websocket context is done, closing the connection")
			return websocket.StatusNormalClosure, "", nil
		case message := <-ch:
			core.Logger().Info("Got websocket message", "message", message)
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
		core.Logger().Info("Got request to", "url", r.URL, "version", r.Proto)
		matches := eventTypeRegex.FindStringSubmatch(r.URL.Path)
		if len(matches) < 2 {
			respondError(w, http.StatusBadRequest, fmt.Errorf("did not specify eventType to subscribe to"))
			return
		}
		eventType := logical.EventType(strings.Join(matches[1:], "/"))
		if eventType == "" {
			respondError(w, http.StatusBadRequest, fmt.Errorf("eventType cannot be blank"))
			return
		}

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

		ctx := r.Context()
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			if err != nil {
				core.Logger().Info("Could not find namespace", "error", err)
				respondError(w, http.StatusInternalServerError, fmt.Errorf("could not find namespace"))
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
		conn.CloseRead(ctx)
		// start the pinger
		go func() {
			for {
				time.Sleep(1 * time.Second)
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
			closeReason = closeReason[:123]
		}
		err = conn.Close(closeStatus, closeReason)
		if err != nil {
			core.Logger().Debug("Error closing websocket", "error", err)
		}
	})
}
