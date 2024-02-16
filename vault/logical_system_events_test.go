// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/plugins/event"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

var (
	noopLock     sync.Mutex
	noopEvents   []string
	unsubscribed atomic.Pointer[string]
)

type noopEventBackend struct{}

func (n *noopEventBackend) Subscribe(_ context.Context, _ *event.SubscribeRequest) error {
	return nil
}

func (n *noopEventBackend) Send(_ context.Context, request *event.SendRequest) error {
	noopLock.Lock()
	defer noopLock.Unlock()
	noopEvents = append(noopEvents, request.EventJSON)
	return nil
}

func (n *noopEventBackend) Unsubscribe(_ context.Context, req *event.UnsubscribeRequest) error {
	unsubscribed.Store(&req.SubscriptionID)
	return nil
}

func (n *noopEventBackend) PluginMetadata() *event.PluginMetadata {
	return nil
}

func (n *noopEventBackend) Close(_ context.Context) error {
	return nil
}

func noopFactory(_ context.Context) (event.SubscriptionPlugin, error) {
	return &noopEventBackend{}, nil
}

// TestSystemBackendEvents_BasicSubscription tests that a basic subscription and send events work.
func TestSystemBackendEvents_BasicSubscription(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	core.eventBackends["noop"] = noopFactory
	ctx := namespace.RootContext(nil)
	req := logical.TestRequest(t, logical.UpdateOperation, "events/subscriptions")
	req.ClientToken = root
	req.Data["config"] = map[string]interface{}{
		"noop": "nothing",
	}
	req.Data["event_type"] = "a*"
	req.Data["plugin"] = "noop"
	resp, err := core.systemBackend.HandleRequest(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Data["http_status_code"])

	// generate an event
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = core.events.SendEventInternal(nil, namespace.RootNamespace, nil, "abc", ev)
	if err != nil {
		t.Fatal(err)
	}
	assert.Eventually(t, func() bool {
		noopLock.Lock()
		defer noopLock.Unlock()
		for _, e := range noopEvents {
			if strings.Contains(e, ev.Id) {
				return true
			}
		}
		return false
	}, 5*time.Second, 10*time.Millisecond)
}

// TestSystemBackendEvents_Unsubscribe tests that unsubscribe terminates the stream of events.
func TestSystemBackendEvents_Unsubscribe(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	core.eventBackends["noop"] = noopFactory
	ctx := namespace.RootContext(nil)
	req := logical.TestRequest(t, logical.UpdateOperation, "events/subscriptions")
	req.ClientToken = root
	req.Data["config"] = map[string]interface{}{
		"noop": "nothing",
	}
	req.Data["event_type"] = "a*"
	req.Data["plugin"] = "noop"
	resp, err := core.systemBackend.HandleRequest(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Data["http_status_code"])
	raw := resp.Data["http_raw_body"].(string)
	data := map[string]interface{}{}
	err = json.Unmarshal([]byte(raw), &data)
	if err != nil {
		t.Fatal(err)
	}
	subID := data["data"].(map[string]interface{})["id"].(string)

	ureq := logical.TestRequest(t, logical.DeleteOperation, "events/subscriptions/noop/"+subID)
	ureq.ClientToken = root
	resp, err = core.systemBackend.HandleRequest(ctx, ureq)
	assert.Nil(t, err)
	fmt.Printf("resp = %+v\n", resp)
	assert.Equal(t, 204, resp.Data["http_status_code"])
	assert.Eventually(t, func() bool {
		x := unsubscribed.Load()
		if x == nil {
			return false
		}
		return *x == subID
	}, 5*time.Second, 10*time.Millisecond)

	// make sure that we don't get an event
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = core.events.SendEventInternal(nil, namespace.RootNamespace, nil, "abc", ev)
	if err != nil {
		t.Fatal(err)
	}
	assert.Never(t, func() bool {
		noopLock.Lock()
		defer noopLock.Unlock()
		for _, e := range noopEvents {
			if strings.Contains(e, ev.Id) {
				return true
			}
		}
		return false
	}, 1*time.Second, 10*time.Millisecond)
}
