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

// columnValuesJSON is the GraphQL JSON scalar used for the create_subitem
// mutation's column_values argument. Implements graphql.GraphQLType so that
// hasura emits "JSON!" in the operation variable declaration rather than
// inferring "String!" from the underlying Go string value.
type columnValuesJSON string

// GetGraphQLType implements graphql.GraphQLType.
func (columnValuesJSON) GetGraphQLType() string { return "JSON" }

// CreateSubitemReq is a request to create a sub-item nested under an existing
// parent item on a Monday.com board.
type CreateSubitemReq struct {
	// ParentItemID is the Monday.com ID of the item that will contain the sub-item.
	ParentItemID string
	// ItemName is the name of the new sub-item.
	ItemName string
	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

// CreateSubitemRes is the result of a CreateSubitemReq.
type CreateSubitemRes struct {
	// SubitemID is the Monday.com ID of the newly created sub-item.
	SubitemID string `json:"subitem_id"`
}

// Run executes the create_subitem mutation against the Monday.com API.
// Transient errors (rate limits, server errors, complexity budget) are
// retried with a fixed delay for up to 5 attempts.
//
// create_subitem is not idempotent: if a server error is returned after
// Monday.com has already created the sub-item, a retry can create a duplicate.
func (r *CreateSubitemReq) Run(ctx context.Context, client *graphql.Client) (*CreateSubitemRes, error) {
	slog.Default().DebugContext(ctx, "Creating Monday board subitem",
		slog.String("parent_item_id", r.ParentItemID),
		slog.String("item_name", r.ItemName),
	)

	var mutation struct {
		CreateSubitem struct {
			ID string `graphql:"id"`
		} `graphql:"create_subitem(parent_item_id: $parentItemId, item_name: $itemName, column_values: $columnValues)"`
	}

	// column_values is part of the mutation but is not exposed on the request,
	// so an empty JSON object is sent and no columns are set.
	vars := map[string]any{
		"parentItemId": graphql.ID(r.ParentItemID),
		"itemName":     graphql.String(r.ItemName),
		"columnValues": columnValuesJSON("{}"),
	}

	var res *CreateSubitemRes
	err := retry.Do(
		func() error {
			if err := client.Mutate(ctx, &mutation, vars); err != nil {
				slog.Default().DebugContext(ctx, "create_subitem mutation failed",
					slog.String("error", err.Error()),
				)
				return err
			}
			res = &CreateSubitemRes{
				SubitemID: mutation.CreateSubitem.ID,
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
		return nil, fmt.Errorf("create_subitem mutation: %w", err)
	}

	slog.Default().DebugContext(ctx, "Monday board subitem created",
		slog.String("subitem_id", res.SubitemID),
	)

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *CreateSubitemRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling create subitem response to JSON: %w", err)
	}
	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CreateSubitemRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"subitem_id"})
	t.AppendRow(table.Row{r.SubitemID})
	return t
}
