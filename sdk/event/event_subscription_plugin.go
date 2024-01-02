// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/helper/backoff"
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

// SubscribeConfigDefaults defines configuration map keys for common default options.
// Embed this in your own config struct to pick up these default options.
type SubscribeConfigDefaults struct {
	Retries         *int           `mapstructure:"retries"`
	RetryMinBackoff *time.Duration `mapstructure:"retry_min_backoff"`
	RetryMaxBackoff *time.Duration `mapstructure:"retry_max_backoff"`
}

// default values for common configuration keys
const (
	DefaultRetries         = 3
	DefaultRetryMinBackoff = 100 * time.Millisecond
	DefaultRetryMaxBackoff = 5 * time.Second
)

func (c *SubscribeConfigDefaults) GetRetries() int {
	if c.Retries == nil {
		return DefaultRetries
	}
	return *c.Retries
}

func (c *SubscribeConfigDefaults) GetRetryMinBackoff() time.Duration {
	if c.RetryMinBackoff == nil {
		return DefaultRetryMinBackoff
	}
	return *c.RetryMinBackoff
}

func (c *SubscribeConfigDefaults) GetRetryMaxBackoff() time.Duration {
	if c.RetryMaxBackoff == nil {
		return DefaultRetryMaxBackoff
	}
	return *c.RetryMaxBackoff
}

func (c *SubscribeConfigDefaults) NewRetryBackoff() *backoff.Backoff {
	return backoff.NewBackoff(c.GetRetries(), c.GetRetryMinBackoff(), c.GetRetryMaxBackoff())
}
