// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	// ErrExpectNameRequired is returned when a real delete is requested
	// without the name the caller expects the item to have.
	ErrExpectNameRequired = errors.New("expect name is required unless dry run is set")
	// ErrBoardMismatch is returned when the item is not on the requested board.
	ErrBoardMismatch = errors.New("item is not on the requested board")
	// ErrNameMismatch is returned when the item's name differs from the expected name.
	ErrNameMismatch = errors.New("item name does not match the expected name")
)

// DeleteItemReq is a request to delete an item or subitem from a Monday.com
// board. A subitem lives on its own subitem board, so pass that board's id as
// BoardID when deleting a subitem.
//
// The request is a compare-and-act check: Run reads the item first and
// refuses to delete it unless it is on BoardID and its name equals
// ExpectName. Monday.com has no conditional delete, so a small race window
// remains between the read and the delete.
type DeleteItemReq struct {
	// BoardID is the board the item must be on.
	BoardID string
	// ItemID selects the item to delete. Names are not unique on a board, so
	// the name is never used to select.
	ItemID string
	// ExpectName is the name the caller expects the item to have. Required
	// unless DryRun is set. When set on a dry run, it is still checked.
	ExpectName string
	// DryRun reads and returns the item without deleting it.
	DryRun bool

	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

// DeleteItemRes describes the targeted item and whether it was deleted.
type DeleteItemRes struct {
	ItemID       string `json:"item_id"`
	ItemName     string `json:"item_name"`
	BoardID      string `json:"board_id"`
	GroupID      string `json:"group_id"`
	GroupTitle   string `json:"group_title"`
	ParentItemID string `json:"parent_item_id,omitempty"`
	// SubitemCount is the number of subitems that are deleted along with the item.
	SubitemCount int  `json:"subitem_count"`
	Deleted      bool `json:"deleted"`
}

// Run checks that the item matches the request and then deletes it. Deleting
// an item also deletes its subitems.
//
// Transient errors are retried for up to 5 attempts. If a delete succeeds on
// the server but the response is lost, the retry fails with an item-not-found
// error even though the item is gone.
func (r *DeleteItemReq) Run(ctx context.Context, client *graphql.Client) (*DeleteItemRes, error) {
	slog.Default().DebugContext(ctx, "Deleting Monday board item",
		slog.String("board_id", r.BoardID),
		slog.String("item_id", r.ItemID),
		slog.String("expect_name", r.ExpectName),
		slog.Bool("dry_run", r.DryRun),
	)

	if !r.DryRun && r.ExpectName == "" {
		return nil, ErrExpectNameRequired
	}

	it, err := getItem(ctx, client, r.ItemID, r.retryDelay)
	if err != nil {
		return nil, err
	}

	if it.BoardID != r.BoardID {
		return nil, fmt.Errorf("%w: item %s is on board %s, not board %s",
			ErrBoardMismatch, r.ItemID, it.BoardID, r.BoardID)
	}

	if r.ExpectName != "" && it.Name != r.ExpectName {
		return nil, fmt.Errorf("%w: item %s is named %q, expected %q",
			ErrNameMismatch, r.ItemID, it.Name, r.ExpectName)
	}

	res := &DeleteItemRes{
		ItemID:       it.ID,
		ItemName:     it.Name,
		BoardID:      it.BoardID,
		GroupID:      it.GroupID,
		GroupTitle:   it.GroupTitle,
		ParentItemID: it.ParentItemID,
		SubitemCount: it.SubitemCount,
	}

	if r.DryRun {
		return res, nil
	}

	if err := r.deleteItem(ctx, client); err != nil {
		return nil, err
	}
	res.Deleted = true

	slog.Default().DebugContext(ctx, "Monday board item deleted",
		slog.String("item_id", res.ItemID),
		slog.String("item_name", res.ItemName),
	)

	return res, nil
}

// deleteItem sends the delete_item mutation. It takes no board id.
func (r *DeleteItemReq) deleteItem(ctx context.Context, client *graphql.Client) error {
	var mutation struct {
		DeleteItem struct {
			ID string `graphql:"id"`
		} `graphql:"delete_item(item_id: $itemId)"`
	}

	vars := map[string]any{
		"itemId": graphql.ID(r.ItemID),
	}

	err := retry.Do(
		func() error {
			if err := client.Mutate(ctx, &mutation, vars); err != nil {
				slog.Default().DebugContext(ctx, "delete_item mutation failed",
					slog.String("error", err.Error()),
				)
				return err
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
		return fmt.Errorf("delete_item mutation: %w", err)
	}

	return nil
}

// ToJSON marshals the response to JSON.
func (r *DeleteItemRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling delete item response to JSON: %w", err)
	}
	return b, nil
}

// ToTable marshals the response to a text table.
func (r *DeleteItemRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"item_id", "item_name", "board_id", "group_id", "group_title", "parent_item_id", "subitem_count", "deleted"})
	t.AppendRow(table.Row{r.ItemID, r.ItemName, r.BoardID, r.GroupID, r.GroupTitle, r.ParentItemID, strconv.Itoa(r.SubitemCount), r.Deleted})
	return t
}
