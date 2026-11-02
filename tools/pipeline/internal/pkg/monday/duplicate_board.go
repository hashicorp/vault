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

// DuplicateBoardReq is a request to create a new Monday.com board by
// duplicating an existing template board.
type DuplicateBoardReq struct {
	// TemplateID is the board ID of the template to duplicate.
	TemplateID string
	// BoardName is the name to assign to the new board.
	BoardName string
	// DuplicateType controls what is copied from the template.
	// Defaults to DuplicateBoardWithStructure when zero.
	DuplicateType DuplicateBoardType

	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

// DuplicateBoardRes is the result of a DuplicateBoardReq.
type DuplicateBoardRes struct {
	// BoardID is the Monday.com ID of the newly created board.
	BoardID string `json:"board_id"`
	// BoardName is the name of the newly created board.
	BoardName string `json:"board_name"`
}

// Run executes the duplicate_board mutation against the Monday.com API.
// Transient errors (rate limits, server errors, complexity budget) are
// retried with a fixed 5-second delay for up to 5 attempts.
func (r *DuplicateBoardReq) Run(ctx context.Context, client *graphql.Client) (*DuplicateBoardRes, error) {
	dupType := r.DuplicateType
	if dupType == "" {
		dupType = DuplicateBoardWithStructure
	}

	slog.Default().DebugContext(ctx, "duplicating Monday.com board",
		slog.String("template_id", r.TemplateID),
		slog.String("board_name", r.BoardName),
		slog.String("duplicate_type", string(dupType)),
	)

	var mutation struct {
		DuplicateBoard struct {
			Board struct {
				ID   graphql.ID `graphql:"id"`
				Name string     `graphql:"name"`
			} `graphql:"board"`
		} `graphql:"duplicate_board(board_id: $boardId, duplicate_type: $dupType, board_name: $boardName)"`
	}

	vars := map[string]any{
		"boardId":   graphql.ID(r.TemplateID),
		"dupType":   dupType,
		"boardName": graphql.String(r.BoardName),
	}

	var res *DuplicateBoardRes
	err := retry.Do(
		func() error {
			if err := client.Mutate(ctx, &mutation, vars); err != nil {
				slog.Default().DebugContext(ctx, "duplicate_board mutation failed",
					slog.String("error", err.Error()),
				)
				return err
			}
			res = &DuplicateBoardRes{
				BoardID:   fmt.Sprintf("%v", mutation.DuplicateBoard.Board.ID),
				BoardName: mutation.DuplicateBoard.Board.Name,
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
		return nil, fmt.Errorf("duplicate_board mutation: %w", err)
	}

	slog.Default().DebugContext(ctx, "Monday.com board created",
		slog.String("board_id", res.BoardID),
		slog.String("board_name", res.BoardName),
	)

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *DuplicateBoardRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling duplicate board response to JSON: %w", err)
	}
	return b, nil
}

// ToTable marshals the response to a text table.
func (r *DuplicateBoardRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"board_id", "board_name"})
	t.AppendRow(table.Row{r.BoardID, r.BoardName})
	return t
}
