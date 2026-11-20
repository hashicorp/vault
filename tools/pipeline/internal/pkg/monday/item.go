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

// item is the subset of a Monday.com item that commands read before acting
// on it, for example to confirm the target before a destructive change.
type item struct {
	ID           string
	Name         string
	BoardID      string
	GroupID      string
	GroupTitle   string
	ParentItemID string // empty unless the item is a subitem
	SubitemCount int
}

// getItem reads a single item by id. items(ids:) is not scoped to a board, so
// callers must compare item.BoardID with the board they expect.
// Transient errors are retried the same way as the mutations.
func getItem(ctx context.Context, client *graphql.Client, itemID string, retryDelay time.Duration) (*item, error) {
	var query struct {
		Items []struct {
			ID    string `graphql:"id"`
			Name  string `graphql:"name"`
			Board struct {
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
		BoardID:      found.Board.ID,
		GroupID:      found.Group.ID,
		GroupTitle:   found.Group.Title,
		SubitemCount: len(found.Subitems),
	}
	if found.ParentItem != nil {
		it.ParentItemID = found.ParentItem.ID
	}

	return it, nil
}
