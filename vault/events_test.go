// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/experiments"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestCanSendEventsFromBuiltinPlugin(t *testing.T) {
	c, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{
		Experiments: []string{experiments.VaultExperimentEventsAlpha1},
	})

	ctx := namespace.RootContext(nil)

	// subscribe to an event type
	eventType, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	ch, cancel, err := c.events.Subscribe(ctx, namespace.RootNamespace, eventType, "")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	// generate the event in a plugin
	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = c.cubbyholeBackend.SendEvent(ctx, logical.EventType(eventType), event)
	if err != nil {
		t.Fatal(err)
	}

	// check that the event is routed to the subscription
	select {
	case receivedEvent := <-ch:
		received := receivedEvent.Payload.(*logical.EventReceived)
		if event.Id != received.Event.Id {
			t.Errorf("Got wrong event: %+v, expected %+v", received, event)
		}
		if received.PluginInfo == nil {
			t.Error("Expected plugin info but got nil")
		} else {
			if received.PluginInfo.Plugin != "cubbyhole" {
				t.Errorf("Expected cubbyhole plugin but got %s", received.PluginInfo.Plugin)
			}
		}

	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for event")
	}
}
