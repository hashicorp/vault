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

// BoardKind is the GraphQL enum for the create_board mutation's board_kind
// argument. Implements graphql.GraphQLType so that hasura emits "BoardKind!"
// in the operation variable declaration rather than inferring "String!".
type BoardKind string

// GetGraphQLType implements graphql.GraphQLType.
func (BoardKind) GetGraphQLType() string { return "BoardKind" }

const (
	// BoardKindPublic creates a publicly visible board.
	BoardKindPublic BoardKind = "public"
	// BoardKindPrivate creates a private board visible only to members.
	BoardKindPrivate BoardKind = "private"
	// BoardKindShare creates a shareable board accessible via link.
	BoardKindShare BoardKind = "share"
)

// CreateBoardReq is a request to create a new Monday.com board from scratch.
type CreateBoardReq struct {
	// BoardName is the name to assign to the new board.
	BoardName string
	// BoardKind controls the visibility of the new board.
	// Defaults to BoardKindPublic when zero.
	BoardKind BoardKind

	// retryDelay overrides the delay between retries. Zero uses defaultRetryDelay.
	retryDelay time.Duration
}

// CreateBoardRes is the result of a CreateBoardReq.
type CreateBoardRes struct {
	// BoardID is the Monday.com ID of the newly created board.
	BoardID string `json:"board_id"`
	// BoardName is the name of the newly created board.
	BoardName string `json:"board_name"`
}

// Run executes the create_board mutation against the Monday.com API.
// Transient errors (rate limits, server errors, complexity budget) are
// retried with a fixed 5-second delay for up to 5 attempts.
func (r *CreateBoardReq) Run(ctx context.Context, client *graphql.Client) (*CreateBoardRes, error) {
	kind := r.BoardKind
	if kind == "" {
		kind = BoardKindPublic
	}

	slog.Default().DebugContext(ctx, "creating Monday.com board",
		slog.String("board_name", r.BoardName),
		slog.String("board_kind", string(kind)),
	)

	var mutation struct {
		CreateBoard struct {
			ID   graphql.ID `graphql:"id"`
			Name string     `graphql:"name"`
		} `graphql:"create_board(board_name: $boardName, board_kind: $boardKind)"`
	}

	vars := map[string]any{
		"boardName": graphql.String(r.BoardName),
		"boardKind": kind,
	}

	var res *CreateBoardRes
	err := retry.Do(
		func() error {
			if err := client.Mutate(ctx, &mutation, vars); err != nil {
				slog.Default().DebugContext(ctx, "create_board mutation failed",
					slog.String("error", err.Error()),
				)
				return err
			}
			res = &CreateBoardRes{
				BoardID:   fmt.Sprintf("%v", mutation.CreateBoard.ID),
				BoardName: mutation.CreateBoard.Name,
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
		return nil, fmt.Errorf("create_board mutation: %w", err)
	}

	slog.Default().DebugContext(ctx, "Monday.com board created",
		slog.String("board_id", res.BoardID),
		slog.String("board_name", res.BoardName),
	)

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *CreateBoardRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling create board response to JSON: %w", err)
	}
	return b, nil
}

// ToTable marshals the response to a text table.
func (r *CreateBoardRes) ToTable() table.Writer {
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
