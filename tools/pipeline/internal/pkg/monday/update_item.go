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

type UpdateItemReq struct {
	BoardID  string
	ItemID   string
	ColumnID string
	Value    string

	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

type UpdateItemRes struct {
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name"`
	NewValue string `json:"new_value"`
}

func (r *UpdateItemReq) Run(ctx context.Context, client *graphql.Client) (*UpdateItemRes, error) {
	slog.Default().DebugContext(ctx, "Updating Monday board column",
		slog.String("board_id", r.BoardID),
		slog.String("item_id", r.ItemID),
		slog.String("column_id", r.ColumnID),
		slog.String("value", r.Value),
	)

	var mutation struct {
		ChangeSimpleColumnValue struct {
			ID   graphql.ID `graphql:"id"`
			Name string     `graphql:"name"`
		} `graphql:"change_simple_column_value(board_id: $boardId, item_id: $itemId, column_id: $columnId, value: $value)"`
	}

	vars := map[string]any{
		"boardId":  graphql.ID(r.BoardID),
		"itemId":   graphql.ID(r.ItemID),
		"columnId": graphql.String(r.ColumnID),
		"value":    graphql.String(r.Value),
	}

	var res *UpdateItemRes
	err := retry.Do(
		func() error {
			if err := client.Mutate(ctx, &mutation, vars); err != nil {
				slog.Default().DebugContext(ctx, "change_simple_column_value mutation failed",
					slog.String("error", err.Error()),
				)
				return err
			}
			res = &UpdateItemRes{
				ItemID:   fmt.Sprintf("%v", mutation.ChangeSimpleColumnValue.ID),
				ItemName: mutation.ChangeSimpleColumnValue.Name,
				NewValue: r.Value,
			}
			return nil
		},
		retry.Attempts(5),
		retry.Delay(retryDelayOrDefault(r.retryDelay)),
		retry.DelayType(retry.FixedDelay),
		retry.Context(ctx),
		retry.RetryIf(isRetryableError),
		retry.WrapContextErrorWithLastError(true),
	)
	if err != nil {
		return nil, fmt.Errorf("change_simple_column_value mutation: %w", err)
	}

	slog.Default().DebugContext(ctx, "Simple column value updated",
		slog.String("item_id", res.ItemID),
		slog.String("item_name", res.ItemName),
		slog.String("value", res.NewValue),
	)

	return res, nil
}

// marshal the response data to JSON
func (r *UpdateItemRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshalling response to JSON: %w", err)
	}

	return b, nil
}

// output response to table
func (r *UpdateItemRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"item_id", "item_name", "new_value"})
	t.AppendRow(table.Row{r.ItemID, r.ItemName, r.NewValue})
	return t
}
