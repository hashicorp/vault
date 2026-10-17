// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/jedib0t/go-pretty/v6/table"
)

type GetItemIDReq struct {
	BoardID  string
	ItemName string
}

type GetItemIDRes struct {
	ItemID string `json:"item_id"`
}

func (r *GetItemIDReq) Run(ctx context.Context, client *graphql.Client) (*GetItemIDRes, error) {
	slog.Default().DebugContext(ctx, "Getting Monday board item_id",
		slog.String("board_id", r.BoardID),
		slog.String("item_name", r.ItemName),
	)

	var query struct {
		Boards []struct {
			ItemsPage struct {
				Items []struct {
					ID   graphql.ID `graphql:"id"`
					Name string     `graphql:"name"`
				}
			} `graphql:"items_page(limit:100)"`
		} `graphql:"boards(ids:[$boardId])"`
	}

	vars := map[string]any{
		"boardId": graphql.ID(r.BoardID),
	}

	err := client.Query(ctx, &query, vars)
	if err != nil {
		return nil, fmt.Errorf("querying item by item-name: %w", err)
	}

	// range over the board we identified in our query and look for
	// the item has a matching name to our request item-name.
	if len(query.Boards) < 1 {
		return nil, fmt.Errorf("board %s not found", r.BoardID)
	}

	items := query.Boards[0].ItemsPage.Items

	for _, item := range items {
		if item.Name == r.ItemName {
			slog.Default().DebugContext(ctx, "item ID located", slog.String("item_id", string(item.ID)))
			return &GetItemIDRes{ItemID: string(item.ID)}, nil
		}
	}

	return nil, fmt.Errorf("couldn't find id for item %s", r.ItemName)
}

// marshal the response data to JSON
func (r *GetItemIDRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshalling response to JSON: %w", err)
	}

	return b, nil
}

func (r *GetItemIDRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"item_id"})
	t.AppendRow(table.Row{r.ItemID})
	return t
}
