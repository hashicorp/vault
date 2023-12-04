// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
)

type EventSubscriptionPlugin interface {
	Initialize(ctx context.Context) error
	Subscribe(ctx context.Context, request *SubscribeRequest) error
	SendSubscriptionEvent(subscriptionID string, eventJson string) error
	Unsubscribe(ctx context.Context, subscriptionID string) error
	// Type returns the name and version for the particular event subscription plugin.
	// This type name is usually set as a constant the backend, e.g., "sqs" for the
	// AWS SQS backend.
	Type() (string, string)
	Close(ctx context.Context) error
}

type SubscribeRequest struct {
	SubscriptionID   string
	Config           map[string]interface{}
	VerifyConnection bool
}
