package eventbus

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestBusBasics tests that basic event sending and subscribing function.
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

	ch, cancel, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

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

// TestNamespaceFiltering verifies that events for other namespaces are filtered out by the bus.
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

	ch, cancel, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

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

// TestBus2Subscriptions verifies that events of different types are successfully routed to the correct subscribers.
func TestBus2Subscriptions(t *testing.T) {
	bus, err := NewEventBus(nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	eventType1 := logical.EventType("someType1")
	eventType2 := logical.EventType("someType2")
	bus.Start()

	ch1, cancel1, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType1)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel1()

	ch2, cancel2, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType2)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel2()

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

// TestBusSubscriptionsCancel verifies that canceled subscriptions are cleaned up.
func TestBusSubscriptionsCancel(t *testing.T) {
	subscriptions.Store(0)
	bus, err := NewEventBus(nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	bus.Start()

	// create and cancel a bunch of subscriptions
	const create = 100
	const cancel = 50

	eventType := logical.EventType("someType")

	var channels []<-chan *logical.EventReceived
	var cancels []context.CancelFunc
	stopped := atomic.Int32{}
	canceled := atomic.Int32{}
	stopped.Store(-1)

	received := atomic.Int32{}

	for i := 0; i < create; i++ {
		ch, cancelFunc, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(cancelFunc)
		channels = append(channels, ch)
		cancels = append(cancels, cancelFunc)

		go func(i int32) {
			<-ch
			received.Add(1)
			if i < int32(cancel) {
				<-ch
				received.Add(1)
			}
			cancelFunc() // cancel explicitly to unsubscribe
			canceled.Add(1)
		}(int32(i))
	}

	// check that all channels receive a message
	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}
	waitFor(t, 1*time.Second, func() bool { return received.Load() == int32(create) })

	if canceled.Load() != int32(cancel) {
		t.Errorf("Expected %d canceled goroutines, but got %d", cancel, canceled.Load())
	}

	// send another message, but half should stop receiving
	event, err = logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}
	waitFor(t, 1*time.Second, func() bool { return received.Load() == int32(create*2-cancel) })
	// the sends should time out and the subscriptions should drop when the cancelFunc was called
	waitFor(t, 1*time.Second, func() bool { return subscriptions.Load() == int64(cancel) })
}

// TestBusSubscriptionsBreak verifies that timed out subscriptions are cleaned up.
func TestBusSubscriptionsBreak(t *testing.T) {
	subscriptions.Store(0)
	bus, err := NewEventBus(nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	// set the timeout very short to make the test faster
	bus.SetSendTimeout(100 * time.Millisecond)
	bus.Start()

	// create and cancel a bunch of subscriptions
	create := 100
	cancel := 50

	eventType := logical.EventType("someType")

	var channels []<-chan *logical.EventReceived
	var cancels []context.CancelFunc
	stopped := atomic.Int32{}
	canceled := atomic.Int32{}
	stopped.Store(-1)

	received := atomic.Int32{}

	for i := 0; i < create; i++ {
		ch, cancelFunc, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(cancelFunc)
		channels = append(channels, ch)
		cancels = append(cancels, cancelFunc)

		go func(i int32) {
			<-ch
			received.Add(1)
			for i < int32(cancel) {
				<-ch
				received.Add(1)
			}
			canceled.Add(1)
		}(int32(i))
	}

	// check that all channels receive a message
	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}
	waitFor(t, 1*time.Second, func() bool { return received.Load() == int32(create) })

	if canceled.Load() != int32(cancel) {
		t.Errorf("Expected %d canceled goroutines, but got %d", cancel, canceled.Load())
	}

	// send another message, but half should stop receiving
	event, err = logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}
	waitFor(t, 1*time.Second, func() bool { return received.Load() == int32(create*2-cancel) })
	// the sends should time out and the subscriptions should drop when the send context cancels
	waitFor(t, 1*time.Second, func() bool { return subscriptions.Load() == int64(cancel) })
}

// waitFor waits for a condition to be true, up to the maximum timeout.
// It waits with a capped exponential backoff starting at 1ms.
// It is guaranteed to try f() at least once.
func waitFor(t *testing.T, maxWait time.Duration, f func() bool) {
	t.Helper()
	start := time.Now()

	if f() {
		return
	}
	sleepAmount := 1 * time.Millisecond
	for time.Now().Sub(start) <= maxWait {
		left := time.Now().Sub(start)
		sleepAmount = sleepAmount * 2
		if sleepAmount > left {
			sleepAmount = left
		}
		time.Sleep(sleepAmount)
		if f() {
			return
		}
	}
	t.Error("Timeout waiting for condition")
}
