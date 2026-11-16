// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/require"
)

// TestArchiveBoardReq_Run verifies ArchiveBoardReq.Run against a minimal
// httptest server standing in for Monday.com. Each case asserts the exact
// sequence of requests the server received, which proves when a mutation is
// sent and, just as importantly, when it is not.
func TestArchiveBoardReq_Run(t *testing.T) {
	t.Parallel()

	const (
		// The read query and mutations that Run() should be sending.
		readBoardQuery       = `query ($boardId:ID!){boards(ids: [$boardId], state: all){id,name,state,items_count,updated_at,workspace{id,name}}}`
		archiveBoardMutation = `mutation ($boardId:ID!){archive_board(board_id: $boardId){id}}`
		deleteBoardMutation  = `mutation ($boardId:ID!){delete_board(board_id: $boardId){id}}`

		boardID         = "1111111111"
		sameNameBoardID = "2222222222"

		// Response bodies follow the Monday.com API shape. IDs are fake.
		boardResponse         = `{"data":{"boards":[{"id":"1111111111","name":"vault-52079 copy","state":"active","items_count":12,"updated_at":"2026-11-10T15:04:05Z","workspace":{"id":"3333333333","name":"Vault Releases"}}]}}`
		sameNameBoardResponse = `{"data":{"boards":[{"id":"2222222222","name":"vault-52079 copy","state":"active","items_count":0,"updated_at":"2026-11-09T08:00:00Z","workspace":null}]}}`
		notFoundResponse      = `{"data":{"boards":[]}}`
		archiveResponse       = `{"data":{"archive_board":{"id":"1111111111"}}}`
		deleteResponse        = `{"data":{"delete_board":{"id":"1111111111"}}}`
		archiveNotAllowedResp = `{"errors":[{"message":"User unauthorized to perform action","path":["archive_board"],"extensions":{"code":"UserUnauthorizedException","status_code":403}}],"data":{"archive_board":null}}`
	)

	// wireRequest is the GraphQL request body that Run() sends to the server.
	type wireRequest struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}

	// serverReply is one canned response from the fake Monday.com server.
	type serverReply struct {
		status int
		body   string
	}

	ok := func(body string) serverReply { return serverReply{status: 200, body: body} }
	httpErr := func(status int) serverReply { return serverReply{status: status, body: `{}`} }

	read := func(id string) wireRequest {
		return wireRequest{Query: readBoardQuery, Variables: map[string]any{"boardId": id}}
	}
	archive := func(id string) wireRequest {
		return wireRequest{Query: archiveBoardMutation, Variables: map[string]any{"boardId": id}}
	}
	hardDelete := func(id string) wireRequest {
		return wireRequest{Query: deleteBoardMutation, Variables: map[string]any{"boardId": id}}
	}

	// run executes req against an httptest server that answers each request
	// with the next reply (the last reply repeats), and returns the result
	// along with every request the server saw.
	run := func(t *testing.T, req *ArchiveBoardReq, replies []serverReply) (*ArchiveBoardRes, []wireRequest, error) {
		t.Helper()

		var mu sync.Mutex
		var reqs []wireRequest
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			mu.Lock()
			reqs = append(reqs, got)
			n := len(reqs)
			mu.Unlock()

			rep := replies[min(n-1, len(replies)-1)]
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(rep.status)
			_, _ = w.Write([]byte(rep.body))
		}))
		t.Cleanup(srv.Close)
		c := graphql.NewClient(srv.URL, srv.Client())

		res, err := req.Run(context.Background(), c)
		mu.Lock()
		defer mu.Unlock()
		return res, append([]wireRequest(nil), reqs...), err
	}

	// boardRes is the result for boardResponse with the given action.
	boardRes := func(action BoardAction) *ArchiveBoardRes {
		return &ArchiveBoardRes{
			BoardID:       boardID,
			BoardName:     "vault-52079 copy",
			State:         "active",
			WorkspaceID:   "3333333333",
			WorkspaceName: "Vault Releases",
			ItemsCount:    12,
			UpdatedAt:     "2026-11-10T15:04:05Z",
			Action:        action,
		}
	}

	for name, tc := range map[string]struct {
		req          *ArchiveBoardReq
		replies      []serverReply // the last reply repeats
		wantRequests []wireRequest
		wantRes      *ArchiveBoardRes
		wantErr      error  // checked with errors.Is
		wantErrMsg   string // checked with ErrorContains, for errors with no sentinel
	}{
		"dry run: sends only the read query": {
			req:          &ArchiveBoardReq{BoardID: boardID, DryRun: true},
			replies:      []serverReply{ok(boardResponse)},
			wantRequests: []wireRequest{read(boardID)},
			wantRes:      boardRes(BoardActionNone),
		},
		"dry run: with hard delete still sends only the read query": {
			req:          &ArchiveBoardReq{BoardID: boardID, HardDelete: true, DryRun: true},
			replies:      []serverReply{ok(boardResponse)},
			wantRequests: []wireRequest{read(boardID)},
			wantRes:      boardRes(BoardActionNone),
		},
		"refused: board not found sends no mutation": {
			req:          &ArchiveBoardReq{BoardID: boardID},
			replies:      []serverReply{ok(notFoundResponse)},
			wantRequests: []wireRequest{read(boardID)},
			wantErr:      ErrBoardNotFound,
		},
		"same name: archives only the given id": {
			// Boards 1111111111 and 2222222222 are both named
			// "vault-52079 copy" in different workspaces. The id selects the
			// target. 2222222222 is in the Main workspace, which Monday.com
			// reports as a null workspace.
			req:          &ArchiveBoardReq{BoardID: sameNameBoardID},
			replies:      []serverReply{ok(sameNameBoardResponse), ok(archiveResponse)},
			wantRequests: []wireRequest{read(sameNameBoardID), archive(sameNameBoardID)},
			wantRes: &ArchiveBoardRes{
				BoardID:       sameNameBoardID,
				BoardName:     "vault-52079 copy",
				State:         "active",
				WorkspaceName: MainWorkspaceName,
				UpdatedAt:     "2026-11-09T08:00:00Z",
				Action:        BoardActionArchived,
			},
		},
		"archive: sends archive_board": {
			req:          &ArchiveBoardReq{BoardID: boardID},
			replies:      []serverReply{ok(boardResponse), ok(archiveResponse)},
			wantRequests: []wireRequest{read(boardID), archive(boardID)},
			wantRes:      boardRes(BoardActionArchived),
		},
		"hard delete: sends delete_board": {
			req:          &ArchiveBoardReq{BoardID: boardID, HardDelete: true},
			replies:      []serverReply{ok(boardResponse), ok(deleteResponse)},
			wantRequests: []wireRequest{read(boardID), hardDelete(boardID)},
			wantRes:      boardRes(BoardActionDeleted),
		},
		"retry: HTTP 503 on the read then success": {
			req:          &ArchiveBoardReq{BoardID: boardID, retryDelay: time.Millisecond},
			replies:      []serverReply{httpErr(503), ok(boardResponse), ok(archiveResponse)},
			wantRequests: []wireRequest{read(boardID), read(boardID), archive(boardID)},
			wantRes:      boardRes(BoardActionArchived),
		},
		"retry: HTTP 503 on archive_board then success": {
			req:          &ArchiveBoardReq{BoardID: boardID, retryDelay: time.Millisecond},
			replies:      []serverReply{ok(boardResponse), httpErr(503), ok(archiveResponse)},
			wantRequests: []wireRequest{read(boardID), archive(boardID), archive(boardID)},
			wantRes:      boardRes(BoardActionArchived),
		},
		"retry: HTTP 503 on delete_board then success": {
			req:          &ArchiveBoardReq{BoardID: boardID, HardDelete: true, retryDelay: time.Millisecond},
			replies:      []serverReply{ok(boardResponse), httpErr(503), ok(deleteResponse)},
			wantRequests: []wireRequest{read(boardID), hardDelete(boardID), hardDelete(boardID)},
			wantRes:      boardRes(BoardActionDeleted),
		},
		"no retry: HTTP 401 on the read sends no mutation": {
			req:          &ArchiveBoardReq{BoardID: boardID, retryDelay: time.Millisecond},
			replies:      []serverReply{httpErr(401)},
			wantRequests: []wireRequest{read(boardID)},
			wantErrMsg:   "boards query",
		},
		"no retry: GraphQL error on archive_board is returned": {
			req:          &ArchiveBoardReq{BoardID: boardID, retryDelay: time.Millisecond},
			replies:      []serverReply{ok(boardResponse), ok(archiveNotAllowedResp)},
			wantRequests: []wireRequest{read(boardID), archive(boardID)},
			wantErrMsg:   "archive_board mutation",
		},
		"no retry: HTTP 403 on delete_board is returned": {
			req:          &ArchiveBoardReq{BoardID: boardID, HardDelete: true, retryDelay: time.Millisecond},
			replies:      []serverReply{ok(boardResponse), httpErr(403)},
			wantRequests: []wireRequest{read(boardID), hardDelete(boardID)},
			wantErrMsg:   "delete_board mutation",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, reqs, err := run(t, tc.req, tc.replies)

			// The exact sequence of requests proves whether a mutation was sent.
			require.Equal(t, tc.wantRequests, reqs)

			if tc.wantErr != nil || tc.wantErrMsg != "" {
				require.Nil(t, res)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
				}
				require.ErrorContains(t, err, tc.wantErrMsg)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantRes, res)
		})
	}
}
