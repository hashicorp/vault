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

// TestCreateBoardReq_Run verifies CreateBoardReq.Run against a minimal httptest
// server standing in for Monday.com. It asserts the exact mutation and
// variables sent on the wire (including the BoardKind enum type and the
// default kind), that the response is mapped into CreateBoardRes, and the
// retry behavior: 429 and 5xx responses and complexity-budget errors are
// retried and then succeed, retryable errors give up after 5 attempts,
// non-retryable HTTP and GraphQL errors (HTTP 200 with an errors array) are
// sent exactly once, and cancelling the context stops further retries.
func TestCreateBoardReq_Run(t *testing.T) {
	t.Parallel()

	const (
		// The mutation that Run() should be sending
		createBoardMutation = `mutation ($boardKind:BoardKind!$boardName:String!){create_board(board_name: $boardName, board_kind: $boardKind){id,name}}`

		// Response bodies captured from the Monday.com API, except
		// complexityResponse, which is a placeholder. IDs and request IDs were
		// replaced with fake values. successResponse came from a create_board
		// call that also passed a workspace_id, which this mutation does not
		// send; the account got userUnauthorizedResp without one.
		successResponse        = `{"data":{"create_board":{"id":"5555555555","name":"vault-51558 new board"}},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		unauthorizedResponse   = `{"errors":["Not Authenticated"]}`
		userUnauthorizedResp   = `{"errors":[{"message":"User unauthorized to perform action","locations":[{"line":1,"column":39}],"path":["create_board"],"extensions":{"code":"UserUnauthorizedException","status_code":403,"error_data":{}}}],"data":{"create_board":null},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		complexityResponse     = `{"errors":[{"message":"Complexity budget exhausted, retry after 10 seconds"}]}`
		createBoardErrorPrefix = "create_board mutation"
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

	ok := serverReply{status: 200, body: successResponse}
	complexity := serverReply{status: 200, body: complexityResponse}
	httpErr := func(status int) serverReply { return serverReply{status: status, body: `{}`} }

	// run executes req against an httptest server that answers each request
	// with the next reply (the last reply repeats), and returns the result
	// along with every request the server saw. When cancelOnRequest is set the
	// server cancels the request's context as soon as it receives a request.
	run := func(t *testing.T, req *CreateBoardReq, replies []serverReply, cancelOnRequest bool) (*CreateBoardRes, []wireRequest, error) {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		var mu sync.Mutex
		var reqs []wireRequest
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			mu.Lock()
			reqs = append(reqs, got)
			n := len(reqs)
			mu.Unlock()

			if cancelOnRequest {
				cancel()
			}
			rep := replies[min(n-1, len(replies)-1)]
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(rep.status)
			_, _ = w.Write([]byte(rep.body))
		}))
		t.Cleanup(srv.Close)
		c := graphql.NewClient(srv.URL, srv.Client())

		res, err := req.Run(ctx, c)
		mu.Lock()
		defer mu.Unlock()
		return res, append([]wireRequest(nil), reqs...), err
	}

	boardReq := func(kind BoardKind) *CreateBoardReq {
		return &CreateBoardReq{BoardName: "vault-51558 new board", BoardKind: kind, retryDelay: time.Millisecond}
	}
	boardVars := func(kind string) map[string]any {
		return map[string]any{"boardName": "vault-51558 new board", "boardKind": kind}
	}
	publicVars := boardVars("public")
	boardRes := &CreateBoardRes{BoardID: "5555555555", BoardName: "vault-51558 new board"}

	for name, tc := range map[string]struct {
		req             *CreateBoardReq
		replies         []serverReply // the last reply repeats
		cancelOnRequest bool
		wantCalls       int
		wantVars        map[string]any
		wantRes         *CreateBoardRes
		wantErr         string
		wantCanceled    bool
	}{
		"kind: default is public": {
			req: boardReq(""), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: publicVars, wantRes: boardRes,
		},
		"kind: public": {
			req: boardReq(BoardKindPublic), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: publicVars, wantRes: boardRes,
		},
		"kind: private": {
			req: boardReq(BoardKindPrivate), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: boardVars("private"), wantRes: boardRes,
		},
		"kind: share": {
			req: boardReq(BoardKindShare), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: boardVars("share"), wantRes: boardRes,
		},
		"retry: HTTP 429 then success": {
			req: boardReq(""), replies: []serverReply{httpErr(429), ok}, wantCalls: 2,
			wantVars: publicVars, wantRes: boardRes,
		},
		"retry: HTTP 500 then success": {
			req: boardReq(""), replies: []serverReply{httpErr(500), ok}, wantCalls: 2,
			wantVars: publicVars, wantRes: boardRes,
		},
		"retry: HTTP 503 then success": {
			req: boardReq(""), replies: []serverReply{httpErr(503), ok}, wantCalls: 2,
			wantVars: publicVars, wantRes: boardRes,
		},
		"retry: complexity budget error then success": {
			req: boardReq(""), replies: []serverReply{complexity, ok}, wantCalls: 2,
			wantVars: publicVars, wantRes: boardRes,
		},
		"retry: two failures then success": {
			req: boardReq(""), replies: []serverReply{httpErr(502), httpErr(504), ok}, wantCalls: 3,
			wantVars: publicVars, wantRes: boardRes,
		},
		"retry: retryable HTTP error exhausts 5 attempts": {
			req: boardReq(""), replies: []serverReply{httpErr(503)}, wantCalls: 5,
			wantVars: publicVars, wantErr: createBoardErrorPrefix,
		},
		"retry: retryable GraphQL error exhausts 5 attempts": {
			req: boardReq(""), replies: []serverReply{complexity}, wantCalls: 5,
			wantVars: publicVars, wantErr: "omplexity",
		},
		"no retry: HTTP 400": {
			req: boardReq(""), replies: []serverReply{httpErr(400)}, wantCalls: 1,
			wantVars: publicVars, wantErr: createBoardErrorPrefix,
		},
		"no retry: HTTP 401": {
			req: boardReq(""), replies: []serverReply{{status: 401, body: unauthorizedResponse}}, wantCalls: 1,
			wantVars: publicVars, wantErr: createBoardErrorPrefix,
		},
		"no retry: HTTP 403": {
			req: boardReq(""), replies: []serverReply{httpErr(403)}, wantCalls: 1,
			wantVars: publicVars, wantErr: createBoardErrorPrefix,
		},
		"no retry: GraphQL user unauthorized in 200 response": {
			req: boardReq(""), replies: []serverReply{{status: 200, body: userUnauthorizedResp}}, wantCalls: 1,
			wantVars: publicVars, wantErr: "User unauthorized to perform action",
		},
		"cancel: cancelled context stops retrying": {
			// The server cancels the context as soon as it receives the
			// request and replies with a retryable 503. Run must return the
			// cancellation without a second attempt and without waiting out
			// the hour-long retry delay (otherwise the test would hang).
			req:             &CreateBoardReq{BoardName: "vault-51558 new board", retryDelay: time.Hour},
			replies:         []serverReply{httpErr(503)},
			cancelOnRequest: true,
			wantCalls:       1,
			wantVars:        publicVars,
			wantErr:         createBoardErrorPrefix,
			wantCanceled:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, reqs, err := run(t, tc.req, tc.replies, tc.cancelOnRequest)

			// Make sure retryable errors are retried and everything else is sent once
			require.Len(t, reqs, tc.wantCalls, "unexpected number of requests")
			// Make sure every attempt sends the correct mutation and variables
			for _, r := range reqs {
				require.Equal(t, createBoardMutation, r.Query)
				require.Equal(t, tc.wantVars, r.Variables)
			}

			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				require.Nil(t, res)
				if tc.wantCanceled {
					require.ErrorIs(t, err, context.Canceled)
				}
				return
			}
			require.NoError(t, err)
			// Make sure the response is mapped into the result
			require.Equal(t, tc.wantRes, res)
		})
	}
}
