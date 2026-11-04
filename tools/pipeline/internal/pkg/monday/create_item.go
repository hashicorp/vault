// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/avast/retry-go/v4"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/jedib0t/go-pretty/v6/table"
)

// CreateItemReq is a request to create a new top-level item on a Monday.com board.
type CreateItemReq struct {
	// BoardID is the Monday.com ID of the board the item is created on.
	BoardID string
	// GroupID is the group on the board that receives the item. Optional: when
	// empty, Monday.com places the item in the board's top group.
	GroupID string
	// ItemName is the name of the new item.
	ItemName string

	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

// CreateItemRes is the result of a CreateItemReq.
type CreateItemRes struct {
	// ItemID is the Monday.com ID of the newly created item.
	ItemID string `json:"item_id"`
	// ItemName is the name of the newly created item.
	ItemName string `json:"item_name"`
}

// Run executes the create_item mutation against the Monday.com API.
// Transient errors (rate limits, complexity budget) are retried with a fixed
// delay for up to 5 attempts. 5xx gateway errors are not retried because
// Monday.com may have already committed the item before the gateway failed —
// retrying would create a duplicate. See isCreateRetryableError.
func (r *CreateItemReq) Run(ctx context.Context, client *graphql.Client) (*CreateItemRes, error) {
	slog.Default().DebugContext(ctx, "Creating Monday board item",
		slog.String("board_id", r.BoardID),
		slog.String("group_id", r.GroupID),
		slog.String("item_name", r.ItemName),
	)

	var mutation struct {
		CreateItem struct {
			ID   graphql.ID `graphql:"id"`
			Name string     `graphql:"name"`
		} `graphql:"create_item(board_id: $boardId, group_id: $groupId, item_name: $itemName)"`
	}

	// An empty GroupID is sent as null instead of "", so Monday.com leaves
	// group_id unset and uses the board's top group.
	var groupID *graphql.String
	if r.GroupID != "" {
		g := graphql.String(r.GroupID)
		groupID = &g
	}

	vars := map[string]any{
		"boardId":  graphql.ID(r.BoardID),
		"groupId":  groupID,
		"itemName": graphql.String(r.ItemName),
	}

	var res *CreateItemRes
	err := retry.Do(
		func() error {
			if err := client.Mutate(ctx, &mutation, vars); err != nil {
				slog.Default().DebugContext(ctx, "create_item mutation failed",
					slog.String("error", err.Error()),
				)
				return err
			}
			res = &CreateItemRes{
				ItemID:   fmt.Sprintf("%v", mutation.CreateItem.ID),
				ItemName: mutation.CreateItem.Name,
			}
			return nil
		},
		retry.Attempts(5),
		retry.Delay(retryDelayOrDefault(r.retryDelay)),
		retry.DelayType(retry.FixedDelay),
		retry.Context(ctx),
		retry.RetryIf(isCreateRetryableError),
		retry.WrapContextErrorWithLastError(true),
	)
	if err != nil {
		return nil, fmt.Errorf("create_item mutation: %w", err)
	}

	slog.Default().DebugContext(ctx, "Monday board item created",
		slog.String("item_id", res.ItemID),
		slog.String("item_name", res.ItemName),
	)

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *CreateItemRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling create item response to JSON: %w", err)
	}
	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CreateItemRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"item_id", "item_name"})
	t.AppendRow(table.Row{r.ItemID, r.ItemName})
	return t
}
