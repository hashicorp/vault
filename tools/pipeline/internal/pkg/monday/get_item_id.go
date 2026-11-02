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
	ParentID string
}

// GetItemIDRes holds the located item ID and the ID of the board it lives on.
// For a regular item lookup, BoardID echoes the request's BoardID,
// but for a subitem lookup BoardID echoes the subitem board ID.
type GetItemIDRes struct {
	ItemID  string `json:"item_id"`
	BoardID string `json:"board_id"`
}

func (r *GetItemIDReq) Run(ctx context.Context, client *graphql.Client) (*GetItemIDRes, error) {
	slog.Default().DebugContext(ctx, "Getting Monday board item_id",
		slog.String("board_id", r.BoardID),
		slog.String("item_name", r.ItemName),
		slog.String("parent_id", r.ParentID),
	)

	if r.ParentID != "" {
		return r.getByParent(ctx, client)
	}

	return r.getByBoard(ctx, client)
}

// getByBoard returns the first item on the board whose name matches ItemName.
// The result's BoardID is the requested BoardID.
func (r *GetItemIDReq) getByBoard(ctx context.Context, client *graphql.Client) (*GetItemIDRes, error) {
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
			return &GetItemIDRes{ItemID: string(item.ID), BoardID: r.BoardID}, nil
		}
	}

	return nil, fmt.Errorf("couldn't find id for item %s", r.ItemName)
}

// getByParent returns the first subitem under the parent item whose name
// matches ItemName exactly (copy names from the template to make sure whitespace matches).
// The result's BoardID is the ID of the subitem board.
func (r *GetItemIDReq) getByParent(ctx context.Context, client *graphql.Client) (*GetItemIDRes, error) {
	var query struct {
		Items []struct {
			Subitems []struct {
				ID    graphql.ID `graphql:"id"`
				Name  string     `graphql:"name"`
				Board struct {
					ID graphql.ID `graphql:"id"`
				} `graphql:"board"`
			} `graphql:"subitems"`
		} `graphql:"items(ids:[$parentId])"`
	}

	vars := map[string]any{
		"parentId": graphql.ID(r.ParentID),
	}

	err := client.Query(ctx, &query, vars)
	if err != nil {
		return nil, fmt.Errorf("querying subitem by parent-id: %w", err)
	}

	if len(query.Items) < 1 {
		return nil, fmt.Errorf("parent item %s not found", r.ParentID)
	}

	for _, sub := range query.Items[0].Subitems {
		if sub.Name == r.ItemName {
			slog.Default().DebugContext(ctx, "subitem ID located", slog.String("item_id", string(sub.ID)),
				slog.String("board_id", string(sub.Board.ID)))
			return &GetItemIDRes{ItemID: string(sub.ID), BoardID: string(sub.Board.ID)}, nil
		}
	}

	return nil, fmt.Errorf("couldn't find subitem %s under parent item %s", r.ItemName, r.ParentID)
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
	t.AppendHeader(table.Row{"item_id", "board_id"})
	t.AppendRow(table.Row{r.ItemID, r.BoardID})
	return t
}
