// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/jedib0t/go-pretty/v6/table"
)

// BoardAction is what ArchiveBoardReq.Run did to the board.
type BoardAction string

const (
	// BoardActionNone means the board was only read (a dry run).
	BoardActionNone BoardAction = "none"
	// BoardActionArchived means the board was archived. Archiving can be
	// undone in the Monday.com UI.
	BoardActionArchived BoardAction = "archived"
	// BoardActionDeleted means the board was deleted.
	BoardActionDeleted BoardAction = "deleted"
)

// ArchiveBoardReq is a request to archive or delete a Monday.com board, for
// example a throwaway board made by create-board or duplicate-board.
type ArchiveBoardReq struct {
	// BoardID selects the board. Names are not unique across workspaces, so
	// the name is never used to select.
	BoardID string
	// HardDelete deletes the board instead of archiving it.
	HardDelete bool
	// DryRun reads and returns the board without changing it.
	DryRun bool

	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

// ArchiveBoardRes describes the targeted board and what was done to it.
type ArchiveBoardRes struct {
	BoardID       string      `json:"board_id"`
	BoardName     string      `json:"board_name"`
	State         string      `json:"state"`
	WorkspaceID   string      `json:"workspace_id,omitempty"`
	WorkspaceName string      `json:"workspace_name"`
	ItemsCount    int         `json:"items_count"`
	UpdatedAt     string      `json:"updated_at"`
	Action        BoardAction `json:"action"`
}

// Run reads the board and then archives it, or deletes it when HardDelete is
// set.
//
// Transient errors are retried for up to 5 attempts.
func (r *ArchiveBoardReq) Run(ctx context.Context, client *graphql.Client) (*ArchiveBoardRes, error) {
	slog.Default().DebugContext(ctx, "Archiving Monday board",
		slog.String("board_id", r.BoardID),
		slog.Bool("hard_delete", r.HardDelete),
		slog.Bool("dry_run", r.DryRun),
	)

	b, err := getBoard(ctx, client, r.BoardID, r.retryDelay)
	if err != nil {
		return nil, err
	}

	res := &ArchiveBoardRes{
		BoardID:       b.ID,
		BoardName:     b.Name,
		State:         b.State,
		WorkspaceID:   b.WorkspaceID,
		WorkspaceName: b.WorkspaceName,
		ItemsCount:    b.ItemsCount,
		UpdatedAt:     b.UpdatedAt,
		Action:        BoardActionNone,
	}

	if r.DryRun {
		return res, nil
	}

	action, err := r.mutate(ctx, client)
	if err != nil {
		return nil, err
	}
	res.Action = action

	slog.Default().DebugContext(ctx, "Monday board changed",
		slog.String("board_id", res.BoardID),
		slog.String("board_name", res.BoardName),
		slog.String("action", string(res.Action)),
	)

	return res, nil
}

// mutate sends archive_board, or delete_board when HardDelete is set, and
// reports which action was taken.
func (r *ArchiveBoardReq) mutate(ctx context.Context, client *graphql.Client) (BoardAction, error) {
	var archiveMutation struct {
		ArchiveBoard struct {
			ID string `graphql:"id"`
		} `graphql:"archive_board(board_id: $boardId)"`
	}
	var deleteMutation struct {
		DeleteBoard struct {
			ID string `graphql:"id"`
		} `graphql:"delete_board(board_id: $boardId)"`
	}

	name, action, mutation := "archive_board", BoardActionArchived, any(&archiveMutation)
	if r.HardDelete {
		name, action, mutation = "delete_board", BoardActionDeleted, any(&deleteMutation)
	}

	vars := map[string]any{
		"boardId": graphql.ID(r.BoardID),
	}

	err := retry.Do(
		func() error {
			if err := client.Mutate(ctx, mutation, vars); err != nil {
				slog.Default().DebugContext(ctx, name+" mutation failed",
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
		return "", fmt.Errorf("%s mutation: %w", name, err)
	}

	return action, nil
}

// ToJSON marshals the response to JSON.
func (r *ArchiveBoardRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling archive board response to JSON: %w", err)
	}
	return b, nil
}

// ToTable marshals the response to a text table.
func (r *ArchiveBoardRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"board_id", "board_name", "state", "workspace_id", "workspace_name", "items_count", "updated_at", "action"})
	t.AppendRow(table.Row{r.BoardID, r.BoardName, r.State, r.WorkspaceID, r.WorkspaceName, strconv.Itoa(r.ItemsCount), r.UpdatedAt, r.Action})
	return t
}
