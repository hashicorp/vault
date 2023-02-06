package eventbus

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestBusBasics(t *testing.T) {
	bus, err := NewEventBus(nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	eventType := logical.EventType("someType")

	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != ErrNotStarted {
		t.Errorf("Expected not started error but got: %v", err)
	}

	bus.Start()

	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Errorf("Expected no error sending: %v", err)
	}

	ch, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType)
	if err != nil {
		t.Fatal(err)
	}

	event, err = logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch:
		if message.Event.ID() != event.ID() {
			t.Errorf("Got unexpected message: %+v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for message")
	}
}

func TestNamespaceFiltering(t *testing.T) {
	bus, err := NewEventBus(nil)
	if err != nil {
		t.Fatal(err)
	}
	bus.Start()
	ctx := context.Background()

	eventType := logical.EventType("someType")

	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	ch, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType)
	if err != nil {
		t.Fatal(err)
	}

	event, err = logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendInternal(ctx, &namespace.Namespace{
		ID:   "abc",
		Path: "/abc",
	}, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(100 * time.Millisecond)
	select {
	case <-ch:
		t.Errorf("Got abc namespace message when root namespace was specified")
	case <-timeout:
		// okay
	}

	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout = time.After(1 * time.Second)
	select {
	case message := <-ch:
		if message.Event.ID() != event.ID() {
			t.Errorf("Got unexpected message %+v but was waiting for %+v", message, event)
		}

	case <-timeout:
		t.Error("Timed out waiting for message")
	}
}

func TestBus2Subscriptions(t *testing.T) {
	bus, err := NewEventBus(nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	eventType1 := logical.EventType("someType1")
	eventType2 := logical.EventType("someType2")
	bus.Start()

	ch1, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType1)
	if err != nil {
		t.Fatal(err)
	}

	ch2, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType2)
	if err != nil {
		t.Fatal(err)
	}

	event1, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	event2, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType2, event2)
	if err != nil {
		t.Error(err)
	}
	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType1, event1)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch1:
		if message.Event.ID() != event1.ID() {
			t.Errorf("Got unexpected message: %v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for event1")
	}
	select {
	case message := <-ch2:
		if message.Event.ID() != event2.ID() {
			t.Errorf("Got unexpected message: %v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for event2")
	}
}
