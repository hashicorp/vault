// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/jedib0t/go-pretty/v6/table"
)

// GetItemReq is a request to read an item or subitem from a Monday.com board.
// It is read-only. A subitem lives on its own subitem board, so pass that
// board's id as BoardID when reading a subitem.
type GetItemReq struct {
	// BoardID is the board the item must be on.
	BoardID string
	// ItemID selects the item to read.
	ItemID string

	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

// ColumnValue is one column value of an item. Text is the display text. Value
// is the raw JSON string Monday.com stores for the column, and is empty for a
// column with no value.
type ColumnValue struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Text  string `json:"text"`
	Value string `json:"value"`
}

// GetItemRes describes the item that was read.
type GetItemRes struct {
	ItemID       string        `json:"item_id"`
	ItemName     string        `json:"item_name"`
	BoardID      string        `json:"board_id"`
	GroupID      string        `json:"group_id"`
	GroupTitle   string        `json:"group_title"`
	State        string        `json:"state"`
	ParentItemID string        `json:"parent_item_id,omitempty"`
	ColumnValues []ColumnValue `json:"column_values"`
	SubitemCount int           `json:"subitem_count"`
	SubitemIDs   []string      `json:"subitem_ids"`
	UpdatedAt    string        `json:"updated_at"`
}

// Run reads the item and checks that it is on BoardID. items(ids:) is not
// scoped to a board, so a wrong board and item pair is rejected with
// ErrBoardMismatch rather than silently accepted. A missing item returns
// ErrItemNotFound.
//
// Retryable errors are retried for up to 5 attempts.
func (r *GetItemReq) Run(ctx context.Context, client *graphql.Client) (*GetItemRes, error) {
	slog.Default().DebugContext(ctx, "Getting Monday board item",
		slog.String("board_id", r.BoardID),
		slog.String("item_id", r.ItemID),
	)

	it, err := getItem(ctx, client, r.ItemID, r.retryDelay)
	if err != nil {
		return nil, err
	}

	if it.BoardID != r.BoardID {
		return nil, fmt.Errorf("%w: item %s is on board %s, not board %s",
			ErrBoardMismatch, r.ItemID, it.BoardID, r.BoardID)
	}

	res := &GetItemRes{
		ItemID:       it.ID,
		ItemName:     it.Name,
		BoardID:      it.BoardID,
		GroupID:      it.GroupID,
		GroupTitle:   it.GroupTitle,
		State:        it.State,
		ParentItemID: it.ParentItemID,
		ColumnValues: make([]ColumnValue, 0, len(it.ColumnValues)),
		SubitemCount: it.SubitemCount,
		SubitemIDs:   it.SubitemIDs,
		UpdatedAt:    it.UpdatedAt,
	}
	if res.SubitemIDs == nil {
		res.SubitemIDs = []string{} // marshal as [] rather than null
	}
	for _, cv := range it.ColumnValues {
		res.ColumnValues = append(res.ColumnValues, ColumnValue(cv))
	}

	return res, nil
}

// ToJSON marshals the response to JSON. Column values include both text and
// the raw value.
func (r *GetItemRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling get item response to JSON: %w", err)
	}
	return b, nil
}

// ToTable marshals the response to a text table of field and value rows. Each
// column value is a row that shows the column's text.
func (r *GetItemRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"field", "value"})
	t.AppendRow(table.Row{"item_id", r.ItemID})
	t.AppendRow(table.Row{"item_name", r.ItemName})
	t.AppendRow(table.Row{"board_id", r.BoardID})
	t.AppendRow(table.Row{"group_id", r.GroupID})
	t.AppendRow(table.Row{"group_title", r.GroupTitle})
	t.AppendRow(table.Row{"state", r.State})
	if r.ParentItemID != "" {
		t.AppendRow(table.Row{"parent_item_id", r.ParentItemID})
	}
	t.AppendRow(table.Row{"subitem_count", strconv.Itoa(r.SubitemCount)})
	if len(r.SubitemIDs) > 0 {
		t.AppendRow(table.Row{"subitem_ids", strings.Join(r.SubitemIDs, ",")})
	}
	t.AppendRow(table.Row{"updated_at", r.UpdatedAt})
	for _, cv := range r.ColumnValues {
		t.AppendRow(table.Row{fmt.Sprintf("column %s (%s, %s)", cv.ID, cv.Title, cv.Type), cv.Text})
	}
	return t
}
