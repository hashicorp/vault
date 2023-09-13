// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/eventbus"
	"github.com/patrickmn/go-cache"
	"github.com/ryanuber/go-glob"
	"nhooyr.io/websocket"
)

// webSocketRevalidationTime is how often we re-check access to the
// events that the websocket requested access to.
var webSocketRevalidationTime = 5 * time.Minute

type eventSubscriber struct {
	ctx               context.Context
	cancelCtx         context.CancelFunc
	clientToken       string
	logger            hclog.Logger
	events            *eventbus.EventBus
	namespacePatterns []string
	pattern           string
	bexprFilter       string
	json              bool
	checkCache        *cache.Cache
	isRootToken       bool
	core              *vault.Core
	w                 http.ResponseWriter
	r                 *http.Request
	req               *logical.Request
}

// handleEventsSubscribeWebsocket subscribes to the events, accepts the websocket connection, and then runs forever,
// serving events to the websocket connection.
func (sub *eventSubscriber) handleEventsSubscribeWebsocket() {
	ctx := sub.ctx
	logger := sub.logger
	// subscribe before accept to avoid race conditions
	ch, cancel, err := sub.events.SubscribeMultipleNamespaces(ctx, sub.namespacePatterns, sub.pattern, sub.bexprFilter)
	if err != nil {
		logger.Info("Error subscribing", "error", err)
		sub.w.WriteHeader(400)
		sub.w.Write([]byte("Error subscribing"))
		return
	}
	defer cancel()
	logger.Debug("WebSocket is subscribed to messages", "namespaces", sub.namespacePatterns, "event_types", sub.pattern, "bexpr_filter", sub.bexprFilter)

	conn, err := websocket.Accept(sub.w, sub.r, nil)
	if err != nil {
		logger.Info("Could not accept as websocket", "error", err)
		respondError(sub.w, http.StatusInternalServerError, fmt.Errorf("could not accept as websocket"))
		return
	}

	// continually validate subscribe access while the websocket is running
	// this has to be done after accepting the websocket to avoid a race condition
	go sub.validateSubscribeAccessLoop()

	// make sure to close the websocket
	closeStatus := websocket.StatusNormalClosure
	closeReason := ""
	var closeErr error = nil

	defer func() {
		if closeErr != nil {
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
	}()

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

	for {
		select {
		case <-ctx.Done():
			logger.Info("Websocket context is done, closing the connection")
			return
		case message := <-ch:
			// Perform one last check that the message is allowed to be received.
			// For example, if a new namespace was created that matches the namespace patterns,
			// but the token doesn't have access to it, we don't want to accidentally send it to
			// the websocket.
			if !sub.allowMessageCached(message.Payload.(*logical.EventReceived)) {
				continue
			}

			logger.Debug("Sending message to websocket", "message", message.Payload)
			var messageBytes []byte
			var messageType websocket.MessageType
			if sub.json {
				var ok bool
				messageBytes, ok = message.Format("cloudevents-json")
				if !ok {
					logger.Warn("Could not get cloudevents JSON format")
					closeErr = errors.New("could not get cloudevents JSON format")
					return
				}
				messageType = websocket.MessageText
			} else {
				messageBytes, err = proto.Marshal(message.Payload.(*logical.EventReceived))
				messageType = websocket.MessageBinary
			}
			if err != nil {
				logger.Warn("Could not serialize websocket event", "error", err)
				closeErr = err
				return
			}
			err = conn.Write(ctx, messageType, messageBytes)
			if err != nil {
				closeErr = err
				return
			}
		}
	}
}

// allowMessageCached checks that the message is allowed to received by the websocket.
// It caches results for specific namespaces, data paths, and event types.
func (sub *eventSubscriber) allowMessageCached(message *logical.EventReceived) bool {
	if sub.isRootToken {
		// fast-path root tokens
		return true
	}

	messageNs := strings.Trim(message.Namespace, "/")
	dataPath := ""
	if message.Event.Metadata != nil {
		dataPathField := message.Event.Metadata.GetFields()[logical.EventMetadataDataPath]
		if dataPathField != nil {
			dataPath = dataPathField.GetStringValue()
		}
	}
	if dataPath == "" {
		// Only allow root tokens to subscribe to events with no data path, for now.
		return false
	}
	cacheKey := fmt.Sprintf("%v!%v!%v", messageNs, dataPath, message.EventType)
	_, ok := sub.checkCache.Get(cacheKey)
	if ok {
		return true
	}

	// perform the actual check and cache it if true
	ok = sub.allowMessage(messageNs, dataPath, message.EventType)
	if ok {
		err := sub.checkCache.Add(cacheKey, ok, webSocketRevalidationTime)
		if err != nil {
			sub.logger.Debug("Error adding to policy check cache for websocket", "error", err)
			// still return the right value, but we can't guarantee it was cached
		}
	}
	return ok
}

// allowMessage checks that the message is allowed to received by the websocket
func (sub *eventSubscriber) allowMessage(eventNs, dataPath, eventType string) bool {
	// does this even match the requested namespaces
	matchedNs := false
	for _, nsPattern := range sub.namespacePatterns {
		if glob.Glob(nsPattern, eventNs) {
			matchedNs = true
			break
		}
	}
	if !matchedNs {
		return false
	}

	// next check for specific access to the namespace and event types
	nsDataPath := dataPath
	if eventNs != "" {
		nsDataPath = path.Join(eventNs, dataPath)
	}
	capabilities, allowedEventTypes, err := sub.core.CapabilitiesAndSubscribeEventTypes(sub.ctx, sub.clientToken, nsDataPath)
	if err != nil {
		sub.logger.Debug("Error checking capabilities and event types for token", "error", err, "namespace", eventNs)
		return false
	}
	if !(slices.Contains(capabilities, vault.RootCapability) || slices.Contains(capabilities, vault.SubscribeCapability)) {
		return false
	}
	for _, pattern := range allowedEventTypes {
		if glob.Glob(pattern, eventType) {
			return true
		}
	}
	// no event types matched, so return false
	return false
}

func handleEventsSubscribe(core *vault.Core, req *logical.Request) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := core.Logger().Named("events-subscribe")
		logger.Debug("Got request to", "url", r.URL, "version", r.Proto)

		ctx := r.Context()

		// ACL check
		auth, entry, err := core.CheckToken(ctx, req, false)
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

		bexprFilter := strings.TrimSpace(r.URL.Query().Get("filter"))
		namespacePatterns := r.URL.Query()["namespaces"]
		namespacePatterns = prependNamespacePatterns(namespacePatterns, ns)
		isRoot := entry.IsRoot()
		ctx, cancelCtx := context.WithCancel(ctx)
		defer cancelCtx()

		sub := &eventSubscriber{
			ctx:               ctx,
			cancelCtx:         cancelCtx,
			logger:            logger,
			events:            core.Events(),
			namespacePatterns: namespacePatterns,
			pattern:           pattern,
			bexprFilter:       bexprFilter,
			json:              json,
			checkCache:        cache.New(webSocketRevalidationTime, webSocketRevalidationTime),
			clientToken:       auth.ClientToken,
			isRootToken:       isRoot,
			core:              core,
			w:                 w,
			r:                 r,
			req:               req,
		}
		sub.handleEventsSubscribeWebsocket()
	})
}

// prependNamespacePatterns prepends the request namespace to the namespace patterns,
// and also adds the request namespace to the list.
func prependNamespacePatterns(patterns []string, requestNamespace *namespace.Namespace) []string {
	prepend := strings.Trim(requestNamespace.Path, "/")
	newPatterns := make([]string, 0, len(patterns)+1)
	newPatterns = append(newPatterns, prepend)
	for _, pattern := range patterns {
		if strings.Trim(pattern, "/") != "" {
			newPatterns = append(newPatterns, path.Join(prepend, pattern))
		}
	}
	return newPatterns
}

// validateSubscribeAccessLoop continually checks if the request has access to the subscribe endpoint in
// its namespace. If the access check ever fails, then the cancel function is called and  the function returns.
func (sub *eventSubscriber) validateSubscribeAccessLoop() {
	// if something breaks, default to canceling the websocket
	defer sub.cancelCtx()
	for {
		_, _, err := sub.core.CheckTokenWithLock(sub.ctx, sub.req, false)
		if err != nil {
			sub.core.Logger().Debug("Token does not have access to subscription path in its own namespace, terminating WebSocket subscription", "path", sub.req.Path, "error", err)
			return
		}
		// wait a while and try again, but quit the loop if the context finishes early
		finished := func() bool {
			ticker := time.NewTicker(webSocketRevalidationTime)
			defer ticker.Stop()
			select {
			case <-sub.ctx.Done():
				return true
			case <-ticker.C:
				return false
			}
		}()
		if finished {
			return
		}
	}
}
