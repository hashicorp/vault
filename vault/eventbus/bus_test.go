// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

// TestBusBasics tests that basic event sending and subscribing function.
func TestBusBasics(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	eventType := logical.EventType("someType")

	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if !errors.Is(err, ErrNotStarted) {
		t.Errorf("Expected not started error but got: %v", err)
	}

	bus.Start()

	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Errorf("Expected no error sending: %v", err)
	}

	ch, cancel, err := bus.Subscribe(ctx, namespace.RootNamespace, string(eventType), "")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	event, err = logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch:
		if message.Payload.(*logical.EventReceived).Event.Id != event.Id {
			t.Errorf("Got unexpected message: %+v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for message")
	}
}

// TestBusIgnoresSendContext tests that the context is ignored when sending to an event,
// so that we do not give up too quickly.
func TestBusIgnoresSendContext(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}
	eventType := logical.EventType("someType")

	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	ch, subCancel, err := bus.Subscribe(context.Background(), namespace.RootNamespace, string(eventType), "")
	if err != nil {
		t.Fatal(err)
	}
	defer subCancel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Errorf("Expected no error sending: %v", err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch:
		if message.Payload.(*logical.EventReceived).Event.Id != event.Id {
			t.Errorf("Got unexpected message: %+v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for message")
	}
}

// TestSubscribeNonRootNamespace verifies that events for non-root namespaces
// aren't filtered out by the bus.
func TestSubscribeNonRootNamespace(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}
	bus.Start()
	ctx := context.Background()

	eventType := logical.EventType("someType")

	ns := &namespace.Namespace{
		ID:   "abc",
		Path: "abc/",
	}

	ch, cancel, err := bus.Subscribe(ctx, ns, string(eventType), "")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendEventInternal(ctx, ns, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch:
		if message.Payload.(*logical.EventReceived).Event.Id != event.Id {
			t.Errorf("Got unexpected message: %+v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for message")
	}
}

// TestNamespaceFiltering verifies that events for other namespaces are filtered out by the bus.
func TestNamespaceFiltering(t *testing.T) {
	bus, err := NewEventBus("", nil)
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

	ch, cancel, err := bus.Subscribe(ctx, namespace.RootNamespace, string(eventType), "")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	event, err = logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}

	err = bus.SendEventInternal(ctx, &namespace.Namespace{
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

	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout = time.After(1 * time.Second)
	select {
	case message := <-ch:
		if message.Payload.(*logical.EventReceived).Event.Id != event.Id {
			t.Errorf("Got unexpected message %+v but was waiting for %+v", message, event)
		}

	case <-timeout:
		t.Error("Timed out waiting for message")
	}
}

// TestBus2Subscriptions verifies that events of different types are successfully routed to the correct subscribers.
func TestBus2Subscriptions(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	eventType1 := logical.EventType("someType1")
	eventType2 := logical.EventType("someType2")
	bus.Start()

	ch1, cancel1, err := bus.Subscribe(ctx, namespace.RootNamespace, string(eventType1), "")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel1()

	ch2, cancel2, err := bus.Subscribe(ctx, namespace.RootNamespace, string(eventType2), "")
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

	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType2, event2)
	if err != nil {
		t.Error(err)
	}
	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType1, event1)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch1:
		if message.Payload.(*logical.EventReceived).Event.Id != event1.Id {
			t.Errorf("Got unexpected message: %v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for event1")
	}
	select {
	case message := <-ch2:
		if message.Payload.(*logical.EventReceived).Event.Id != event2.Id {
			t.Errorf("Got unexpected message: %v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for event2")
	}
}

// TestBusSubscriptionsCancel verifies that canceled subscriptions are cleaned up.
func TestBusSubscriptionsCancel(t *testing.T) {
	testCases := []struct {
		cancel bool
	}{
		{cancel: true},
		{cancel: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("cancel=%v", tc.cancel), func(t *testing.T) {
			subscriptions.Store(0)
			bus, err := NewEventBus("", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			if !tc.cancel {
				// set the timeout very short to make the test faster if we aren't canceling explicitly
				bus.SetSendTimeout(100 * time.Millisecond)
			}
			bus.Start()

			// create and stop a bunch of subscriptions
			const create = 100
			const stop = 50

			eventType := logical.EventType("someType")

			var channels []<-chan *eventlogger.Event
			var cancels []context.CancelFunc
			stopped := atomic.Int32{}

			received := atomic.Int32{}

			for i := 0; i < create; i++ {
				ch, cancelFunc, err := bus.Subscribe(ctx, namespace.RootNamespace, string(eventType), "")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(cancelFunc)
				channels = append(channels, ch)
				cancels = append(cancels, cancelFunc)

				go func(i int32) {
					<-ch // always receive one message
					received.Add(1)
					// continue receiving messages as long as are not stopped
					for i < int32(stop) {
						<-ch
						received.Add(1)
					}
					if tc.cancel {
						cancelFunc() // stop explicitly to unsubscribe
					}
					stopped.Add(1)
				}(int32(i))
			}

			// check that all channels receive a message
			event, err := logical.NewEvent()
			if err != nil {
				t.Fatal(err)
			}
			err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType, event)
			if err != nil {
				t.Error(err)
			}
			waitFor(t, 1*time.Second, func() bool { return received.Load() == int32(create) })
			waitFor(t, 1*time.Second, func() bool { return stopped.Load() == int32(stop) })

			// send another message, but half should stop receiving
			event, err = logical.NewEvent()
			if err != nil {
				t.Fatal(err)
			}
			err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, eventType, event)
			if err != nil {
				t.Error(err)
			}
			waitFor(t, 1*time.Second, func() bool { return received.Load() == int32(create*2-stop) })
			// the sends should time out and the subscriptions should drop when cancelFunc is called or the context cancels
			waitFor(t, 1*time.Second, func() bool { return subscriptions.Load() == int64(create-stop) })
		})
	}
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

// TestBusWildcardSubscriptions tests that a single subscription can receive
// multiple event types using * for glob patterns.
func TestBusWildcardSubscriptions(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	fooEventType := logical.EventType("kv/foo")
	barEventType := logical.EventType("kv/bar")
	bus.Start()

	ch1, cancel1, err := bus.Subscribe(ctx, namespace.RootNamespace, "kv/*", "")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel1()

	ch2, cancel2, err := bus.Subscribe(ctx, namespace.RootNamespace, "*/bar", "")
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

	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, barEventType, event2)
	if err != nil {
		t.Error(err)
	}
	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, fooEventType, event1)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	// Expect to receive both events on ch1, which subscribed to kv/*
	var ch1Seen []string
	for i := 0; i < 2; i++ {
		select {
		case message := <-ch1:
			ch1Seen = append(ch1Seen, message.Payload.(*logical.EventReceived).Event.Id)
		case <-timeout:
			t.Error("Timeout waiting for event1")
		}
	}
	if len(ch1Seen) != 2 {
		t.Errorf("Expected 2 events but got: %v", ch1Seen)
	} else {
		if !strutil.StrListContains(ch1Seen, event1.Id) {
			t.Errorf("Did not find %s event1 ID in ch1seen", event1.Id)
		}
		if !strutil.StrListContains(ch1Seen, event2.Id) {
			t.Errorf("Did not find %s event2 ID in ch1seen", event2.Id)
		}
	}
	// Expect to receive just kv/bar on ch2, which subscribed to */bar
	select {
	case message := <-ch2:
		if message.Payload.(*logical.EventReceived).Event.Id != event2.Id {
			t.Errorf("Got unexpected message: %v", message)
		}
	case <-timeout:
		t.Error("Timeout waiting for event2")
	}
}

// TestDataPathIsPrependedWithMount tests that "data_path", if present in the
// metadata, is prepended with the plugin's mount.
func TestDataPathIsPrependedWithMount(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	fooEventType := logical.EventType("kv/foo")
	bus.Start()

	ch, cancel, err := bus.Subscribe(ctx, namespace.RootNamespace, "kv/*", "")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	event, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	metadata := map[string]string{
		logical.EventMetadataDataPath: "my/secret/path",
		"not_touched":                 "xyz",
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		t.Fatal(err)
	}
	event.Metadata = &structpb.Struct{}
	if err := event.Metadata.UnmarshalJSON(metadataBytes); err != nil {
		t.Fatal(err)
	}

	// no plugin info means nothing should change
	err = bus.SendEventInternal(ctx, namespace.RootNamespace, nil, fooEventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case message := <-ch:
		metadata := message.Payload.(*logical.EventReceived).Event.Metadata.AsMap()
		assert.Contains(t, metadata, "not_touched")
		assert.Equal(t, "xyz", metadata["not_touched"])
		assert.Contains(t, metadata, "data_path")
		assert.Equal(t, "my/secret/path", metadata["data_path"])
	case <-timeout:
		t.Error("Timeout waiting for event")
	}

	// send with a secrets plugin mounted
	pluginInfo := logical.EventPluginInfo{
		MountClass:    "secrets",
		MountAccessor: "kv_abc",
		MountPath:     "secret/",
		Plugin:        "kv",
		PluginVersion: "v1.13.1+builtin",
		Version:       "2",
	}
	err = bus.SendEventInternal(ctx, namespace.RootNamespace, &pluginInfo, fooEventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout = time.After(1 * time.Second)
	select {
	case message := <-ch:
		metadata := message.Payload.(*logical.EventReceived).Event.Metadata.AsMap()
		assert.Contains(t, metadata, "not_touched")
		assert.Equal(t, "xyz", metadata["not_touched"])
		assert.Contains(t, metadata, "data_path")
		assert.Equal(t, "secret/my/secret/path", metadata["data_path"])
	case <-timeout:
		t.Error("Timeout waiting for event")
	}

	// send with an auth plugin mounted
	pluginInfo = logical.EventPluginInfo{
		MountClass:    "auth",
		MountAccessor: "kubernetes_abc",
		MountPath:     "kubernetes/",
		Plugin:        "vault-plugin-auth-kubernetes",
		PluginVersion: "v1.13.1+builtin",
	}
	event, err = logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	metadata = map[string]string{
		logical.EventMetadataDataPath: "my/secret/path",
		"not_touched":                 "xyz",
	}
	metadataBytes, err = json.Marshal(metadata)
	if err != nil {
		t.Fatal(err)
	}
	event.Metadata = &structpb.Struct{}
	if err := event.Metadata.UnmarshalJSON(metadataBytes); err != nil {
		t.Fatal(err)
	}
	err = bus.SendEventInternal(ctx, namespace.RootNamespace, &pluginInfo, fooEventType, event)
	if err != nil {
		t.Error(err)
	}

	timeout = time.After(1 * time.Second)
	select {
	case message := <-ch:
		metadata := message.Payload.(*logical.EventReceived).Event.Metadata.AsMap()
		assert.Contains(t, metadata, "not_touched")
		assert.Equal(t, "xyz", metadata["not_touched"])
		assert.Contains(t, metadata, "data_path")
		assert.Equal(t, "auth/kubernetes/my/secret/path", metadata["data_path"])
	case <-timeout:
		t.Error("Timeout waiting for event")
	}
}

// TestBexpr tests go-bexpr filters are evaluated on an event.
func TestBexpr(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	bus.Start()

	sendEvent := func(eventType string) error {
		event, err := logical.NewEvent()
		if err != nil {
			return err
		}
		metadata := map[string]string{
			logical.EventMetadataDataPath:  "my/secret/path",
			logical.EventMetadataOperation: "write",
		}
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		event.Metadata = &structpb.Struct{}
		if err := event.Metadata.UnmarshalJSON(metadataBytes); err != nil {
			return err
		}
		// send with a secrets plugin mounted
		pluginInfo := logical.EventPluginInfo{
			MountClass:    "secrets",
			MountAccessor: "kv_abc",
			MountPath:     "secret/",
			Plugin:        "kv",
			PluginVersion: "v1.13.1+builtin",
			Version:       "2",
		}
		return bus.SendEventInternal(ctx, namespace.RootNamespace, &pluginInfo, logical.EventType(eventType), event)
	}

	testCases := []struct {
		name             string
		filter           string
		shouldPassFilter bool
	}{
		{"empty expression", "", true},
		{"non-matching expression", "data_path == nothing", false},
		{"matching expression", "data_path == secret/my/secret/path", true},
		{"full matching expression", "data_path == secret/my/secret/path and operation != read and source_plugin_mount == secret/ and source_plugin_mount != somethingelse", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			eventType, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatal(err)
			}
			ch, cancel, err := bus.Subscribe(ctx, namespace.RootNamespace, eventType, testCase.filter)
			if err != nil {
				t.Fatal(err)
			}
			defer cancel()
			err = sendEvent(eventType)
			if err != nil {
				t.Fatal(err)
			}

			timer := time.NewTimer(5 * time.Second)
			defer timer.Stop()
			got := false
			select {
			case <-ch:
				got = true
			case <-timer.C:
			}
			assert.Equal(t, testCase.shouldPassFilter, got)
		})
	}
}

// TestPipelineCleanedUp ensures pipelines are properly cleaned up after
// subscriptions are closed.
func TestPipelineCleanedUp(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	eventType := logical.EventType("someType")
	bus.Start()

	_, cancel, err := bus.Subscribe(context.Background(), namespace.RootNamespace, string(eventType), "")
	if err != nil {
		t.Fatal(err)
	}
	// check that the filters are set
	if !bus.filters.anyMatch(namespace.RootNamespace, eventType) {
		t.Fatal()
	}
	if !bus.broker.IsAnyPipelineRegistered(eventTypeAll) {
		cancel()
		t.Fatal()
	}

	cancel()

	if bus.broker.IsAnyPipelineRegistered(eventTypeAll) {
		t.Fatal()
	}

	// and that the filters are cleaned up
	if bus.filters.anyMatch(namespace.RootNamespace, eventType) {
		t.Fatal()
	}
}

// TestSubscribeGlobal tests that the global filter subscription mechanism works.
func TestSubscribeGlobal(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	bus.filters.addGlobalPattern([]string{""}, "abc*")
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	ch, cancel2, err := bus.NewGlobalSubscription(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendEventInternal(nil, namespace.RootNamespace, nil, "abcd", ev)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case recv := <-ch:
		// ok
		event := recv.Payload.(*logical.EventReceived)
		assert.Equal(t, "abcd", event.EventType)
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for event")
	}
}

// TestSubscribeGlobal_WithApply tests that the global filter subscription mechanism works when using ApplyGlobalFilterChanges.
func TestSubscribeGlobal_WithApply(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()
	assert.False(t, bus.GlobalMatch(namespace.RootNamespace, "abcd"))
	bus.ApplyGlobalFilterChanges([]FilterChange{
		{
			Operation:         FilterChangeAdd,
			NamespacePatterns: []string{""},
			EventTypePattern:  "abc*",
		},
	})
	assert.True(t, bus.GlobalMatch(namespace.RootNamespace, "abcd"))

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	ch, cancel2, err := bus.NewGlobalSubscription(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendEventInternal(nil, namespace.RootNamespace, nil, "abcd", ev)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case recv := <-ch:
		// ok
		event := recv.Payload.(*logical.EventReceived)
		assert.Equal(t, "abcd", event.EventType)
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for event")
	}
}

// TestSubscribeCluster tests that the cluster filter subscription mechanism works.
func TestSubscribeCluster(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	bus.filters.addPattern("somecluster", []string{""}, "abc*")
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	ch, cancel2, err := bus.NewClusterSubscription(ctx, "somecluster")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendEventInternal(nil, namespace.RootNamespace, nil, "abcd", ev)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case recv := <-ch:
		// ok
		event := recv.Payload.(*logical.EventReceived)
		assert.Equal(t, "abcd", event.EventType)
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for event")
	}
}

// TestSubscribeCluster_WithApply tests that the cluster filter subscription mechanism works when using ApplyClusterFilterChanges.
func TestSubscribeCluster_WithApply(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	bus.ApplyClusterFilterChanges("somecluster", []FilterChange{
		{
			Operation:         FilterChangeAdd,
			NamespacePatterns: []string{""},
			EventTypePattern:  "abc*",
		},
	})
	ch, cancel2, err := bus.NewClusterSubscription(ctx, "somecluster")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendEventInternal(nil, namespace.RootNamespace, nil, "abcd", ev)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case recv := <-ch:
		// ok
		event := recv.Payload.(*logical.EventReceived)
		assert.Equal(t, "abcd", event.EventType)
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for event")
	}
}

// TestClearGlobalFilter tests that clearing the global filter means no messages get through.
func TestClearGlobalFilter(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	bus.ApplyGlobalFilterChanges([]FilterChange{
		{
			Operation:         FilterChangeAdd,
			NamespacePatterns: []string{""},
			EventTypePattern:  "abc*",
		},
	})
	bus.ClearGlobalFilter()
	ch, cancel2, err := bus.NewGlobalSubscription(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendEventInternal(nil, namespace.RootNamespace, nil, "abcd", ev)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-ch:
		t.Fatal("We should not have gotten an event")
	case <-time.After(100 * time.Millisecond):
		// ok
	}
}

// TestClearClusterFilter tests that clearing a cluster filter means no messages get through.
func TestClearClusterFilter(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	bus.ApplyClusterFilterChanges("somecluster", []FilterChange{
		{
			Operation:         FilterChangeAdd,
			NamespacePatterns: []string{""},
			EventTypePattern:  "abc*",
		},
	})
	bus.ClearClusterFilter("somecluster")
	ch, cancel2, err := bus.NewClusterSubscription(ctx, "somecluster")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	ev, err := logical.NewEvent()
	if err != nil {
		t.Fatal(err)
	}
	err = bus.SendEventInternal(nil, namespace.RootNamespace, nil, "abcd", ev)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-ch:
		t.Fatal("We should not have gotten an event")
	case <-time.After(100 * time.Millisecond):
		// ok
	}
}

// TestNotifyOnGlobalFilterChanges tests that notifications on global filter changes are sent.
func TestNotifyOnGlobalFilterChanges(t *testing.T) {
	bus, err := NewEventBus("", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	ch, cancel2, err := bus.NotifyOnGlobalFilterChanges(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	bus.ApplyGlobalFilterChanges([]FilterChange{
		{
			Operation:         FilterChangeAdd,
			NamespacePatterns: []string{""},
			EventTypePattern:  "abc*",
		},
	})

	select {
	case changes := <-ch:
		if len(changes) == 2 {
			assert.Equal(t, []FilterChange{{Operation: FilterChangeClear}, {Operation: FilterChangeAdd, NamespacePatterns: []string{""}, EventTypePattern: "abc*"}}, changes)
		} else {
			// could be split into two updates
			assert.Len(t, changes, 1)
			assert.Equal(t, []FilterChange{{Operation: FilterChangeClear}}, changes)
			changes := <-ch
			assert.Len(t, changes, 1)
			assert.Equal(t, []FilterChange{{Operation: FilterChangeAdd, NamespacePatterns: []string{""}, EventTypePattern: "abc*"}}, changes)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("We expected to get a global filter notification")
	}
}

// TestNotifyOnLocalFilterChanges tests that notifications on local cluster filter changes are sent.
func TestNotifyOnLocalFilterChanges(t *testing.T) {
	bus, err := NewEventBus("somecluster", nil)
	if err != nil {
		t.Fatal(err)
	}

	bus.Start()

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	ch, cancel2, err := bus.NotifyOnLocalFilterChanges(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cancel2)
	bus.ApplyClusterFilterChanges("somecluster", []FilterChange{
		{
			Operation:         FilterChangeAdd,
			NamespacePatterns: []string{""},
			EventTypePattern:  "abc*",
		},
	})

	select {
	case changes := <-ch:
		if len(changes) == 2 {
			assert.Equal(t, []FilterChange{{Operation: FilterChangeClear}, {Operation: FilterChangeAdd, NamespacePatterns: []string{""}, EventTypePattern: "abc*"}}, changes)
		} else {
			// could be split into two updates
			assert.Len(t, changes, 1)
			assert.Equal(t, []FilterChange{{Operation: FilterChangeClear}}, changes)
			changes := <-ch
			assert.Len(t, changes, 1)
			assert.Equal(t, []FilterChange{{Operation: FilterChangeAdd, NamespacePatterns: []string{""}, EventTypePattern: "abc*"}}, changes)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("We expected to get a global filter notification")
	}
}
