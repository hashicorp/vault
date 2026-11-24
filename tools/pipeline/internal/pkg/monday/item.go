// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/avast/retry-go/v4"
	graphql "github.com/hasura/go-graphql-client"
)

// ErrItemNotFound is returned when no item exists with the requested id.
var ErrItemNotFound = errors.New("item not found")

// itemColumnValue is one column value of an item. Text is the display text
// and Value is the raw JSON string Monday.com stores for the column.
type itemColumnValue struct {
	ID    string
	Title string
	Type  string
	Text  string
	Value string
}

// item is the subset of a Monday.com item that commands read before acting
// on it, for example to confirm the target before a destructive change.
type item struct {
	ID           string
	Name         string
	State        string
	UpdatedAt    string
	BoardID      string
	GroupID      string
	GroupTitle   string
	ParentItemID string // empty unless the item is a subitem
	SubitemIDs   []string
	SubitemCount int
	ColumnValues []itemColumnValue
}

// getItem reads a single item by id. items(ids:) is not scoped to a board, so
// callers must compare item.BoardID with the board they expect.
// Transient errors are retried the same way as the mutations.
func getItem(ctx context.Context, client *graphql.Client, itemID string, retryDelay time.Duration) (*item, error) {
	var query struct {
		Items []struct {
			ID        string `graphql:"id"`
			Name      string `graphql:"name"`
			State     string `graphql:"state"`
			UpdatedAt string `graphql:"updated_at"`
			Board     struct {
				ID string `graphql:"id"`
			} `graphql:"board"`
			Group struct {
				ID    string `graphql:"id"`
				Title string `graphql:"title"`
			} `graphql:"group"`
			ParentItem *struct {
				ID string `graphql:"id"`
			} `graphql:"parent_item"`
			Subitems []struct {
				ID string `graphql:"id"`
			} `graphql:"subitems"`
			ColumnValues []struct {
				ID     string `graphql:"id"`
				Type   string `graphql:"type"`
				Text   string `graphql:"text"`
				Value  string `graphql:"value"`
				Column struct {
					Title string `graphql:"title"`
				} `graphql:"column"`
			} `graphql:"column_values"`
		} `graphql:"items(ids: [$itemId])"`
	}

	vars := map[string]any{
		"itemId": graphql.ID(itemID),
	}

	err := retry.Do(
		func() error {
			if err := client.Query(ctx, &query, vars); err != nil {
				slog.Default().DebugContext(ctx, "items query failed",
					slog.String("error", err.Error()),
				)
				return err
			}
			return nil
		},
		retry.Attempts(5),
		retry.Delay(retryDelayOrDefault(retryDelay)),
		retry.DelayType(retry.FixedDelay),
		retry.Context(ctx),
		retry.RetryIf(isRetryableError),
		retry.WrapContextErrorWithLastError(true),
	)
	if err != nil {
		return nil, fmt.Errorf("items query: %w", err)
	}

	if len(query.Items) == 0 {
		return nil, fmt.Errorf("%w: item %s", ErrItemNotFound, itemID)
	}

	found := query.Items[0]
	it := &item{
		ID:           found.ID,
		Name:         found.Name,
		State:        found.State,
		UpdatedAt:    found.UpdatedAt,
		BoardID:      found.Board.ID,
		GroupID:      found.Group.ID,
		GroupTitle:   found.Group.Title,
		SubitemCount: len(found.Subitems),
	}
	if found.ParentItem != nil {
		it.ParentItemID = found.ParentItem.ID
	}
	for _, s := range found.Subitems {
		it.SubitemIDs = append(it.SubitemIDs, s.ID)
	}
	for _, cv := range found.ColumnValues {
		it.ColumnValues = append(it.ColumnValues, itemColumnValue{
			ID: cv.ID, Title: cv.Column.Title, Type: cv.Type, Text: cv.Text, Value: cv.Value,
		})
	}

	return it, nil
}
