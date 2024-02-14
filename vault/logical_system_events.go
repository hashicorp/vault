// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/syncmap"
	"github.com/hashicorp/vault/plugins/event"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type eventSubscriptions struct {
	subscriptions *syncmap.SyncMap[string, *eventSubscription]
}

type eventSubscription struct {
	ctx     context.Context
	cancel  context.CancelFunc
	id      string
	plugin  string
	config  map[string]interface{}
	backend event.SubscriptionPlugin
}

func (sub *eventSubscription) ID() string {
	return sub.id
}

func newEventSubscriptions() *eventSubscriptions {
	return &eventSubscriptions{
		subscriptions: syncmap.NewSyncMap[string, *eventSubscription](),
	}
}

// handleEventsSubscribe
func (b *SystemBackend) handleEventsSubscribe(requestCtx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// TODO: check policy
	eventTypePattern := data.Get("event_type").(string)
	if eventTypePattern == "" {
		return logical.ErrorResponse("event_type is required"), logical.ErrInvalidRequest
	}
	bexprFilter := data.Get("filter").(string)
	plugin := data.Get("plugin").(string)
	if plugin == "" {
		return logical.ErrorResponse("plugin is required"), logical.ErrInvalidRequest
	}
	factory, ok := b.Core.eventBackends[plugin]
	if !ok {
		return logical.ErrorResponse("unsupported plugin type"), logical.ErrInvalidRequest
	}
	config := data.Get("config").(map[string]interface{})
	if config == nil || len(config) == 0 {
		return logical.ErrorResponse("config is required"), logical.ErrInvalidRequest
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		b.logger.Error("Error generating UUID", "error", err)
		return logical.ErrorResponse("error generating UUID"), logical.ErrUnrecoverable
	}
	backend, err := factory(requestCtx)
	if err != nil {
		b.logger.Error("Error initializing plugin", "error", err)
		return logical.ErrorResponse("error initializing plugin"), logical.ErrUnrecoverable
	}
	if err != nil {
		b.logger.Warn("Error subscribing to events", "error", err)
		return logical.ErrorResponse("error subscribing to events"), logical.ErrUnrecoverable
	}
	ns, err := namespace.FromContext(requestCtx)
	if err != nil {
		b.logger.Error("No namespace found", "error", err)
		return logical.ErrorResponse("no namespace found"), logical.ErrUnrecoverable
	}
	ctx, ctxCancel := context.WithCancel(context.Background())
	eventCh, eventCancel, err := b.Core.events.Subscribe(ctx, ns, eventTypePattern, bexprFilter)
	sub := &eventSubscription{
		ctx:     ctx,
		id:      id,
		plugin:  plugin,
		config:  config,
		backend: backend,
	}
	sub.cancel = func() {
		sub.cancel = func() {} // for safety
		ctxCancel()
		eventCancel()
	}
	err = sub.backend.Subscribe(requestCtx, &event.SubscribeRequest{
		SubscriptionID:   id,
		Config:           config,
		VerifyConnection: false, // TODO: read this from parameters
	})
	if err != nil {
		b.logger.Warn("Error starting subscription", "error", err)
		return logical.ErrorResponse(fmt.Sprintf("error starting subscription: %v", err)), logical.ErrUnrecoverable
	}
	b.eventSubscriptions.subscriptions.Put(id, sub)
	go b.subscriptionLoop(sub, eventCh)
	resp := &logical.Response{
		Data: map[string]interface{}{
			"id":     id,
			"plugin": plugin,
		},
	}
	return logical.RespondWithStatusCode(resp, req, http.StatusOK)
}

func (b *SystemBackend) unsubscribe(sub *eventSubscription) {
	sub.cancel()
	go func() {
		// subscription context might have been canceled already
		err := sub.backend.Unsubscribe(context.Background(), &event.UnsubscribeRequest{
			SubscriptionID: sub.id,
		})
		if err != nil {
			b.logger.Warn("Error unsubscribing", "subscription_id", sub.id, "error", err)
		}
	}()
}

func (b *SystemBackend) subscriptionLoop(sub *eventSubscription, eventCh <-chan *eventlogger.Event) {
	doneCh := sub.ctx.Done()
	defer b.unsubscribe(sub)
	for {
		select {
		case <-doneCh:
			return
		case payload := <-eventCh:
			jsonBytes, ok := payload.Format("cloudevents-json")
			if !ok {
				b.logger.Error("Event system is not supplying 'cloudevents-json' format")
				return
			}
			err := sub.backend.Send(sub.ctx, &event.SendRequest{
				SubscriptionID: sub.id,
				EventJSON:      string(jsonBytes),
			})
			if err != nil {
				b.logger.Warn("Error sending event", "subscription_id", sub.id, "error", err)
				return
			}
		}
	}
}

// handleEventsUnsubscribe
func (b *SystemBackend) handleEventsUnsubscribe(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	plugin := data.Get("plugin").(string)
	id := data.Get("id").(string)
	if plugin == "" || id == "" {
		return logical.ErrorResponse("no subscription specified"), logical.ErrNotFound
	}
	sub := b.eventSubscriptions.subscriptions.Pop(id) // Pop avoids race condition
	if sub == nil {
		return logical.ErrorResponse("no subscription found"), logical.ErrNotFound
	}
	if sub.plugin != plugin {
		b.eventSubscriptions.subscriptions.Put(id, sub)
		return logical.ErrorResponse("wrong plugin type specified"), logical.ErrNotFound
	}
	b.unsubscribe(sub)
	return logical.RespondWithStatusCode(nil, req, http.StatusNoContent)
}

// handleEventsListSubscriptions
func (b *SystemBackend) handleEventsListSubscriptions(_ context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	allSubs := b.eventSubscriptions.subscriptions.Values()
	var redactedSubs []map[string]any
	for _, sub := range allSubs {
		newSub := map[string]any{
			"id":     sub.id,
			"plugin": sub.plugin,
		}
		redactedSubs = append(redactedSubs, newSub)
	}
	resp := &logical.Response{
		Data: map[string]interface{}{
			"subscriptions": redactedSubs,
		},
	}
	return logical.RespondWithStatusCode(resp, req, http.StatusOK)
}

// handleReadSubscription
func (b *SystemBackend) handleEventsReadSubscription(_ context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	plugin := data.Get("plugin").(string)
	id := data.Get("id").(string)
	if plugin == "" || id == "" {
		return logical.ErrorResponse("plugin and id are required"), logical.ErrInvalidRequest
	}
	sub := b.eventSubscriptions.subscriptions.Get(id)
	if sub == nil {
		return logical.ErrorResponse("subscription not found"), logical.ErrNotFound
	}
	if sub.plugin != plugin {
		return logical.ErrorResponse(fmt.Sprintf("subscription found, but with plugin %s instead of %s", sub.plugin, plugin)), logical.ErrNotFound
	}
	resp := &logical.Response{
		Data: map[string]interface{}{
			"id":     sub.id,
			"plugin": sub.plugin,
		},
	}
	return logical.RespondWithStatusCode(resp, req, http.StatusOK)
}
