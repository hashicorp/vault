// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/eventbus"
	"nhooyr.io/websocket"
)

type eventSubscribeArgs struct {
	ctx               context.Context
	logger            hclog.Logger
	events            *eventbus.EventBus
	namespacePatterns []string
	pattern           string
	conn              *websocket.Conn
	json              bool
}

// handleEventsSubscribeWebsocket runs forever, returning a websocket error code and reason
// only if the connection closes or there was an error.
func handleEventsSubscribeWebsocket(args eventSubscribeArgs) (websocket.StatusCode, string, error) {
	ctx := args.ctx
	logger := args.logger
	ch, cancel, err := args.events.SubscribeMultipleNamespaces(ctx, args.namespacePatterns, args.pattern)
	if err != nil {
		logger.Info("Error subscribing", "error", err)
		return websocket.StatusUnsupportedData, "Error subscribing", nil
	}
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Websocket context is done, closing the connection")
			return websocket.StatusNormalClosure, "", nil
		case message := <-ch:
			logger.Debug("Sending message to websocket", "message", message.Payload)
			var messageBytes []byte
			var messageType websocket.MessageType
			if args.json {
				var ok bool
				messageBytes, ok = message.Format("cloudevents-json")
				if !ok {
					logger.Warn("Could not get cloudevents JSON format")
					return 0, "", errors.New("could not get cloudevents JSON format")
				}
				messageType = websocket.MessageText
			} else {
				messageBytes, err = proto.Marshal(message.Payload.(*logical.EventReceived))
				messageType = websocket.MessageBinary
			}
			if err != nil {
				logger.Warn("Could not serialize websocket event", "error", err)
				return 0, "", err
			}
			err = args.conn.Write(ctx, messageType, messageBytes)
			if err != nil {
				return 0, "", err
			}
		}
	}
}

func handleEventsSubscribe(core *vault.Core, req *logical.Request) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := core.Logger().Named("events-subscribe")
		logger.Debug("Got request to", "url", r.URL, "version", r.Proto)

		ctx := r.Context()

		// ACL check
		_, _, err := core.CheckToken(ctx, req, false)
		if err != nil {
			if errors.Is(err, logical.ErrPermissionDenied) {
				respondError(w, http.StatusForbidden, logical.ErrPermissionDenied)
				return
			}
			logger.Debug("Error validating token", "error", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("error validating token"))
			return
		}

		ns, err := namespace.FromContext(ctx)
		if err != nil {
			logger.Info("Could not find namespace", "error", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("could not find namespace"))
			return
		}

		prefix := "/v1/sys/events/subscribe/"
		if ns.ID != namespace.RootNamespaceID {
			prefix = fmt.Sprintf("/v1/%ssys/events/subscribe/", ns.Path)
		}
		pattern := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
		if pattern == "" {
			respondError(w, http.StatusBadRequest, fmt.Errorf("did not specify eventType to subscribe to"))
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

		namespacePatterns := r.URL.Query()["namespaces"]
		namespacePatterns = prependNamespacePatterns(namespacePatterns, ns)
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			logger.Info("Could not accept as websocket", "error", err)
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

		closeStatus, closeReason, err := handleEventsSubscribeWebsocket(eventSubscribeArgs{ctx, logger, core.Events(), namespacePatterns, pattern, conn, json})
		if err != nil {
			closeStatus = websocket.CloseStatus(err)
			if closeStatus == -1 {
				closeStatus = websocket.StatusInternalError
			}
			closeReason = fmt.Sprintf("Internal error: %v", err)
			logger.Debug("Error from websocket handler", "error", err)
		}
		// Close() will panic if the reason is greater than this length
		if len(closeReason) > 123 {
			logger.Debug("Truncated close reason", "closeReason", closeReason)
			closeReason = closeReason[:123]
		}
		err = conn.Close(closeStatus, closeReason)
		if err != nil {
			logger.Debug("Error closing websocket", "error", err)
		}
	})
}

// prependNamespacePatterns prepends the request namespace to the namespace patterns,
// and also adds the request namespace to the list.
func prependNamespacePatterns(patterns []string, requestNamespace *namespace.Namespace) []string {
	prepend := strings.Trim(requestNamespace.Path, "/")
	newPatterns := make([]string, 0, len(patterns)+1)
	newPatterns = append(newPatterns, prepend)
	for _, pattern := range patterns {
		if strings.Trim(strings.TrimSpace(pattern), "/") == "" {
			continue
		}
		newPatterns = append(newPatterns, path.Join(prepend, pattern, "/"))
	}
	return newPatterns
}
