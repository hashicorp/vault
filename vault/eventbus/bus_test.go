package eventbus

import (
	"context"
	"testing"
	"time"
)

func TestBusBasics(t *testing.T) {
	bus := NewEventBus()
	ctx := context.Background()

	eventType := "someType"

	err := bus.Send(ctx, eventType, "message")
	if err != ErrNotStarted {
		t.Errorf("Expected not started error but got: %v", err)
	}

	bus.Start()

	err = bus.Send(ctx, eventType, "sent but never received")
	if err != nil {
		t.Errorf("Expected no error sending: %v", err)
	}

	ch, err := bus.Subscribe(ctx, eventType)
	if err != nil {
		t.Fatal(err)
	}

	err = bus.Send(ctx, eventType, "message2")
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch:
		if message != "message2" {
			t.Errorf("Got unexpected message: %v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for message")
	}
}
