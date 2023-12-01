// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	UICustomMessageKey string = "/custom-messages"

	MaximumCustomMessageCount int = 100
)

// customMessageBarrierView determines the appropriate logical.Storage to return
// depending on whether the logical.Storage is being used to access entries for
// the root namespace or any other namespace based on the provided
// context.Context.
func (c *UIConfig) customMessageBarrierView(ctx context.Context) logical.Storage {
	// If nsBarrierView is nil, which occurs in the non-enterprise edition, then
	// simply use the barrierStorage.
	if c.nsBarrierView == nil {
		return c.barrierStorage
	}

	// Retrieve the namespace from the context.Context
	// namespace.FromContext returns an error when:
	// 1. ctx is nil
	// 2. there's no Namespace value in ctx
	// 3. the Namespace value in ctx is nil
	// In each of those cases, returning the barrierStorage is an appropriate
	// course of action.
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return c.barrierStorage
	}

	return NewBarrierView(c.nsBarrierView, ns.ID)
}

// ListUICustomMessagesFilters is a struct that captures the different filtering
// criteria that can be provided to the (*UIConfig).ListCustomMessages method.
type ListUICustomMessagesFilters struct {
	authenticated *bool
	active        *bool
	messageType   *string
}

// Authenticated adds the authenticated filter criterion to the receiver
// ListUICustomMessagesFilters.
func (f *ListUICustomMessagesFilters) Authenticated(value bool) {
	f.authenticated = &value
}

// Active adds the active filter criterion to the receiver
// ListUICustomMessagesFilters.
func (f *ListUICustomMessagesFilters) Active(value bool) {
	f.active = &value
}

// MessageType adds the messageType filter criterion to the receiver
// ListUICustomMessagesFilters.
func (f *ListUICustomMessagesFilters) MessageType(value string) {
	f.messageType = &value
}

// UICustomMessageEntry is a struct that contains all of the details of a
// custom message. This type is used to encode the information into a
// logical.StorageEntry as well as to transmit custom messages between
// the logical layer and the request handling layer.
type UICustomMessageEntry struct {
	Id            string         `json:"id"`
	Title         string         `json:"title"`
	Message       string         `json:"message"`
	StartTime     time.Time      `json:"start_time"`
	EndTime       time.Time      `json:"end_time"`
	Options       map[string]any `json:"options"`
	Link          map[string]any `json:"link"`
	Authenticated bool           `json:"authenticated"`
	MessageType   string         `json:"type"`
	active        bool
}

// isTimeNowBetween is a function that determines if the current time, returned
// by time.Now() is after the provided startTime and before the provided
// endTime.
func isTimeNowBetween(startTime, endTime time.Time) bool {
	now := time.Now()

	return !(startTime.After(now) || endTime.Before(now))
}

// ListCustomMessages retrieves all of the custom messages for the appropriate
// namespace.Namespace based on the provided context.Context and the receiver's
// configuration. The provided ListUICustomMessagesFilters is then used to
// determine which custom messages satisfy the filter criteria.
func (c *UIConfig) ListCustomMessages(ctx context.Context, filters ListUICustomMessagesFilters) ([]UICustomMessageEntry, error) {
	entries, err := c.retrieveCustomMessagesInternal(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]UICustomMessageEntry, 0)

	// Calculate Active property and apply filters
	for _, entry := range entries {
		entry.active = isTimeNowBetween(entry.StartTime, entry.EndTime)

		if filters.authenticated != nil && *filters.authenticated != entry.Authenticated {
			continue
		}

		if filters.messageType != nil && *filters.messageType != entry.MessageType {
			continue
		}

		if filters.active != nil && *filters.active != entry.active {
			continue
		}

		results = append(results, entry)
	}

	return results, nil
}

// retrieveCustomMessagesInternal handles the internal logic of retrieving all
// of the custom messages stored in the current namespace. If there are no
// custom messages, an empty slice of UICustomMessageEntry is returned.
func (c *UIConfig) retrieveCustomMessagesInternal(ctx context.Context) ([]UICustomMessageEntry, error) {
	c.customMessageLock.RLock()
	defer c.customMessageLock.RUnlock()

	keys, err := c.customMessageBarrierView(ctx).List(ctx, fmt.Sprintf("%s/", UICustomMessageKey))
	if err != nil {
		return nil, err
	}

	results := make([]UICustomMessageEntry, len(keys))

	for idx, key := range keys {
		storageEntry, err := c.customMessageBarrierView(ctx).Get(ctx, fmt.Sprintf("%s/%s", UICustomMessageKey, key))
		if err != nil {
			return nil, err
		}

		customMessageEntry := UICustomMessageEntry{}
		if err = storageEntry.DecodeJSON(&customMessageEntry); err != nil {
			return nil, err
		}

		results[idx] = customMessageEntry
	}

	return results, nil
}

// ReadCustomMessage reads a specific custom message from the underlying storage
// based on the provided messageId value.
func (c *UIConfig) ReadCustomMessage(ctx context.Context, messageId string) (*UICustomMessageEntry, error) {
	customMessageEntry, err := c.retrieveCustomMessageInternal(ctx, messageId)
	if err != nil {
		return nil, err
	}

	if customMessageEntry != nil {
		customMessageEntry.active = isTimeNowBetween(customMessageEntry.StartTime, customMessageEntry.EndTime)
	}

	return customMessageEntry, nil
}

// retrieveCustomMessageInternal handles the internal logic to retrieve a specific
// custom message. If no custom message exists with the provided messageId,
// nil, nil is returned
func (c *UIConfig) retrieveCustomMessageInternal(ctx context.Context, messageId string) (*UICustomMessageEntry, error) {
	c.customMessageLock.RLock()
	defer c.customMessageLock.RUnlock()

	storageEntry, err := c.customMessageBarrierView(ctx).Get(ctx, fmt.Sprintf("%s/%s", UICustomMessageKey, messageId))
	if err != nil {
		return nil, err
	}

	var customMessageEntry *UICustomMessageEntry

	if storageEntry != nil {
		customMessageEntry = &UICustomMessageEntry{}
		if err = storageEntry.DecodeJSON(customMessageEntry); err != nil {
			return nil, err
		}
	}

	return customMessageEntry, nil
}

// DeleteCustomMessage removes a specific custom message from the underlying
// storage. The custom message is specified by the messageId argument. If no
// custom message exists with the provided messageId, no error is returned.
func (c *UIConfig) DeleteCustomMessage(ctx context.Context, messageId string) error {
	c.customMessageLock.Lock()
	defer c.customMessageLock.Unlock()

	return c.customMessageBarrierView(ctx).Delete(ctx, fmt.Sprintf("%s/%s", UICustomMessageKey, messageId))
}

// CreateCustomMessage stores the provided UICustomMessageEntry into the
// underlying storage.
func (c *UIConfig) CreateCustomMessage(ctx context.Context, entry UICustomMessageEntry) (*UICustomMessageEntry, error) {
	if err := validateStartAndEndTimes(entry.StartTime, entry.EndTime); err != nil {
		return nil, err
	}

	count, err := c.countCustomMessagesInternal(ctx)
	if err != nil {
		return nil, err
	}

	if count >= MaximumCustomMessageCount {
		return nil, errors.New("maximum number of Custom Messages already exists")
	}

	messageId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	entry.Id = messageId

	err = c.saveCustomMessageInternal(ctx, entry)
	if err != nil {
		return nil, err
	}

	entry.active = isTimeNowBetween(entry.StartTime, entry.EndTime)

	return &entry, nil
}

// countCustomMessagesInternal returns a count of existing custom messages. It's used to
// detect if the maximum number of custom messages has been met.
func (c *UIConfig) countCustomMessagesInternal(ctx context.Context) (int, error) {
	c.customMessageLock.RLock()
	defer c.customMessageLock.RUnlock()

	keys, err := c.customMessageBarrierView(ctx).List(ctx, fmt.Sprintf("%s/", UICustomMessageKey))
	if err != nil {
		return 0, err
	}

	return len(keys), nil
}

// validateStartAndEndTimes tests the provided startTime and endTime to ensure
// that startTime occurs BEFORE endTime, otherwise an error is returned.
func validateStartAndEndTimes(startTime, endTime time.Time) error {
	if !startTime.Before(endTime) {
		return errors.New("start_time must occur before end_time")
	}

	return nil
}

// UpdateCustomMessage modifies the properties of an existing custom message.
func (c *UIConfig) UpdateCustomMessage(ctx context.Context, entry UICustomMessageEntry) (*UICustomMessageEntry, error) {
	if err := validateStartAndEndTimes(entry.StartTime, entry.EndTime); err != nil {
		return nil, err
	}

	err := c.saveCustomMessageInternal(ctx, entry)
	if err != nil {
		return nil, err
	}

	entry.active = isTimeNowBetween(entry.StartTime, entry.EndTime)

	return &entry, nil
}

// saveCustomMessageInternal handles the internal logic of storing a new or
// updated custom message in the underlying storage.
func (c *UIConfig) saveCustomMessageInternal(ctx context.Context, customMessage UICustomMessageEntry) error {
	updatedValue, err := json.Marshal(&customMessage)
	if err != nil {
		return err
	}

	storageEntry := &logical.StorageEntry{
		Key:   fmt.Sprintf("%s/%s", UICustomMessageKey, customMessage.Id),
		Value: updatedValue,
	}

	c.customMessageLock.Lock()
	defer c.customMessageLock.Unlock()

	return c.customMessageBarrierView(ctx).Put(ctx, storageEntry)
}
