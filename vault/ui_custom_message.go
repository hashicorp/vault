package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	UICustomMessageKey string = "/ui-custom-messages"

	MaximumCustomMessageCount int = 100
)

type ListUICustomMessagesFilters struct {
	authenticated *bool
	active        *bool
	messageType   *string
}

func (f *ListUICustomMessagesFilters) Authenticated(value bool) {
	f.authenticated = &value
}

func (f *ListUICustomMessagesFilters) Active(value bool) {
	f.active = &value
}

func (f *ListUICustomMessagesFilters) MessageType(value string) {
	f.messageType = &value
}

type UICustomMessagesEntry struct {
	Id            string         `json:"id"`
	Title         string         `json:"title"`
	Message       string         `json:"message"`
	StartTime     time.Time      `json:"start-time"`
	EndTime       time.Time      `json:"end-time"`
	Options       map[string]any `json:"options"`
	Link          map[string]any `json:"link"`
	Authenticated bool           `json:"authenticated"`
	MessageType   string         `json:"type"`
	active        bool
}

func isTimeNowBetween(startTime, endTime time.Time) bool {
	now := time.Now()

	return !(startTime.After(now) || endTime.Before(now))
}

func (c *UIConfig) ListCustomMessages(ctx context.Context, filters ListUICustomMessagesFilters) ([]*UICustomMessagesEntry, error) {
	entries, err := c.retrieveCustomMessages(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*UICustomMessagesEntry, 0)

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

func (c *UIConfig) retrieveCustomMessages(ctx context.Context) ([]*UICustomMessagesEntry, error) {
	c.customMessageLock.RLock()
	defer c.customMessageLock.RUnlock()

	keys, err := c.barrierStorage.List(ctx, fmt.Sprintf("%s/", UICustomMessageKey))
	if err != nil {
		return nil, err
	}

	results := make([]*UICustomMessagesEntry, len(keys))

	for idx, key := range keys {
		storageEntry, err := c.barrierStorage.Get(ctx, fmt.Sprintf("%s/%s", UICustomMessageKey, key))
		if err != nil {
			return nil, err
		}

		customMessageEntry := &UICustomMessagesEntry{}
		if err = storageEntry.DecodeJSON(customMessageEntry); err != nil {
			return nil, err
		}

		results[idx] = customMessageEntry
	}

	return results, nil
}

func (c *UIConfig) ReadCustomMessage(ctx context.Context, messageId string) (*UICustomMessagesEntry, error) {
	customMessageEntry, err := c.retrieveCustomMessage(ctx, messageId)
	if err != nil {
		return nil, err
	}

	customMessageEntry.active = isTimeNowBetween(customMessageEntry.StartTime, customMessageEntry.EndTime)

	return customMessageEntry, nil
}

func (c *UIConfig) retrieveCustomMessage(ctx context.Context, messageId string) (*UICustomMessagesEntry, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	c.customMessageLock.RLock()
	defer c.customMessageLock.RUnlock()

	storageEntry, err := c.barrierStorage.Get(ctx, fmt.Sprintf("%s%s%s", ns.Path, UICustomMessageKey, messageId))
	if err != nil {
		return nil, err
	}

	var customMessageEntry *UICustomMessagesEntry
	if err = storageEntry.DecodeJSON(customMessageEntry); err != nil {
		return nil, err
	}

	return customMessageEntry, nil
}

func (c *UIConfig) DeleteCustomMessage(ctx context.Context, messageId string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	c.customMessageLock.Lock()
	defer c.customMessageLock.Unlock()

	return c.barrierStorage.Delete(ctx, fmt.Sprintf("%s%s/%s", ns.Path, UICustomMessageKey, messageId))
}

func (c *UIConfig) CreateCustomMessage(ctx context.Context, entry UICustomMessagesEntry) (*UICustomMessagesEntry, error) {
	count, err := c.countCustomMessages(ctx)
	if err != nil {
		return nil, err
	}

	if count >= MaximumCustomMessageCount {
		return nil, logical.ErrUnrecoverable
	}

	messageId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	entry.Id = messageId

	err = c.saveCustomMessage(ctx, entry)
	if err != nil {
		return nil, err
	}

	entry.active = isTimeNowBetween(entry.StartTime, entry.EndTime)

	return &entry, nil
}

func (c *UIConfig) countCustomMessages(ctx context.Context) (int, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	c.customMessageLock.RLock()
	defer c.customMessageLock.RUnlock()

	keys, err := c.barrierStorage.List(ctx, fmt.Sprintf("%s%s/", ns.Path, UICustomMessageKey))
	if err != nil {
		return 0, err
	}

	return len(keys), nil
}

func (c *UIConfig) UpdateCustomMessage(ctx context.Context, entry UICustomMessagesEntry) (*UICustomMessagesEntry, error) {
	err := c.saveCustomMessage(ctx, entry)
	if err != nil {
		return nil, err
	}

	entry.active = isTimeNowBetween(entry.StartTime, entry.EndTime)

	return &entry, nil
}

func (c *UIConfig) saveCustomMessage(ctx context.Context, entry UICustomMessagesEntry) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	c.customMessageLock.Lock()
	defer c.customMessageLock.Unlock()

	customMessageRaw, err := json.Marshal(&entry)
	if err != nil {
		return err
	}

	storageEntry := &logical.StorageEntry{
		Key:   fmt.Sprintf("%s%s/%s", ns.Path, UICustomMessageKey, entry.Id),
		Value: customMessageRaw,
	}

	return c.barrierStorage.Put(ctx, storageEntry)
}
