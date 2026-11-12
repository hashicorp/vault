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

// ErrBoardNotFound is returned when no board exists with the requested id.
var ErrBoardNotFound = errors.New("board not found")

// MainWorkspaceName is reported for boards in the Main workspace, which
// Monday.com returns as a null workspace.
const MainWorkspaceName = "Main workspace"

// board is the subset of a Monday.com board that commands read before acting
// on it, for example to confirm the target before a destructive change.
type board struct {
	ID            string
	Name          string
	State         string
	WorkspaceID   string // empty for the Main workspace
	WorkspaceName string
	ItemsCount    int
	UpdatedAt     string
}

// getBoard reads a single board by id in any state, so archived and deleted
// boards can still be inspected. Transient errors are retried the same way as
// the mutations.
func getBoard(ctx context.Context, client *graphql.Client, boardID string, retryDelay time.Duration) (*board, error) {
	var query struct {
		Boards []struct {
			ID         string `graphql:"id"`
			Name       string `graphql:"name"`
			State      string `graphql:"state"`
			ItemsCount int    `graphql:"items_count"`
			UpdatedAt  string `graphql:"updated_at"`
			Workspace  *struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"workspace"`
		} `graphql:"boards(ids: [$boardId], state: all)"`
	}

	vars := map[string]any{
		"boardId": graphql.ID(boardID),
	}

	err := retry.Do(
		func() error {
			if err := client.Query(ctx, &query, vars); err != nil {
				slog.Default().DebugContext(ctx, "boards query failed",
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
		return nil, fmt.Errorf("boards query: %w", err)
	}

	if len(query.Boards) == 0 {
		return nil, fmt.Errorf("%w: board %s", ErrBoardNotFound, boardID)
	}

	found := query.Boards[0]
	b := &board{
		ID:            found.ID,
		Name:          found.Name,
		State:         found.State,
		WorkspaceName: MainWorkspaceName,
		ItemsCount:    found.ItemsCount,
		UpdatedAt:     found.UpdatedAt,
	}
	if found.Workspace != nil {
		b.WorkspaceID = found.Workspace.ID
		b.WorkspaceName = found.Workspace.Name
	}

	return b, nil
}
