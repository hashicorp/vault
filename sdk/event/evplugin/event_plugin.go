package evplugin

import "context"

// these are just stubs for now

// EventPlugin allows plugins to enable new ways to subscribe to events.
// Plugins should implement the server and Vault should implement the client.
type EventPlugin interface {
	// Initialize the plugin.
	Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error)
	// Subscribe creates a new subscription. Only a single subscription will be created per plugin.
	Subscribe(ctx context.Context, req SubscribeRequest) (SubscribeResponse, error)
	// ReceiveEvents is called with a stream of events for a single subscription.
	ReceiveEvents(ctx context.Context, subscriptionID uint32, events <-chan SubscriptionEvent) (ReceiveEventsResponse, error)
	// Type is used to return the type of the plugin, e.g., "vault-plugin-event-sqs".
	Type(ctx context.Context) (string, error)
	// Close is used to close the subscription.
	Close() error
}

type InitializeRequest struct{}

type InitializeResponse struct{}

type SubscribeRequest struct{}

type SubscribeResponse struct{}

type SubscriptionEvent struct{}

type ReceiveEventsResponse struct{}
