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

// TestDuplicateBoardReq_Run verifies DuplicateBoardReq.Run against a minimal
// httptest server standing in for Monday.com. It asserts the exact mutation
// and variables sent on the wire (including the DuplicateBoardType enum type
// and the default duplicate type), that the nested board in the response is
// mapped into DuplicateBoardRes, and the retry behavior: 429 and 5xx responses
// and complexity-budget errors are retried and then succeed, retryable errors
// give up after 5 attempts, non-retryable HTTP and GraphQL errors (HTTP 200
// with an errors array) are sent exactly once, and cancelling the context
// stops further retries.
func TestDuplicateBoardReq_Run(t *testing.T) {
	t.Parallel()

	const (
		// The mutation that Run() should be sending
		duplicateBoardMutation = `mutation ($boardId:ID!$boardName:String!$dupType:DuplicateBoardType!){duplicate_board(board_id: $boardId, duplicate_type: $dupType, board_name: $boardName){board{id,name}}}`

		// Response bodies captured from the Monday.com API on a
		// throwaway board, except complexityResponse, which is a placeholder.
		// IDs and request IDs were replaced with fake values.
		successResponse      = `{"data":{"duplicate_board":{"board":{"id":"4444444444","name":"vault-51558 copy"}}},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		unauthorizedResponse = `{"errors":["Not Authenticated"]}`
		boardNotFoundResp    = `{"errors":[{"message":"Board not found","locations":[{"line":1,"column":57}],"path":["duplicate_board"],"extensions":{"code":"ResourceNotFoundException","status_code":404,"error_data":{"resource_type":"board","board_id":999999999999,"error_reason":"store.monday.automation.error.missing_board"}}}],"data":{"duplicate_board":null},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		complexityResponse   = `{"errors":[{"message":"Complexity budget exhausted, retry after 10 seconds"}]}`

		duplicateBoardErrorPrefix = "duplicate_board mutation"
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
	run := func(t *testing.T, req *DuplicateBoardReq, replies []serverReply, cancelOnRequest bool) (*DuplicateBoardRes, []wireRequest, error) {
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

	dupReq := func(dupType DuplicateBoardType) *DuplicateBoardReq {
		return &DuplicateBoardReq{
			TemplateID:    "1111111111",
			BoardName:     "vault-51558 copy",
			DuplicateType: dupType,
			retryDelay:    time.Millisecond,
		}
	}
	dupVars := func(dupType string) map[string]any {
		return map[string]any{"boardId": "1111111111", "dupType": dupType, "boardName": "vault-51558 copy"}
	}
	structureVars := dupVars("duplicate_board_with_structure")
	dupRes := &DuplicateBoardRes{BoardID: "4444444444", BoardName: "vault-51558 copy"}

	for name, tc := range map[string]struct {
		req             *DuplicateBoardReq
		replies         []serverReply // the last reply repeats
		cancelOnRequest bool
		wantCalls       int
		wantVars        map[string]any
		wantRes         *DuplicateBoardRes
		wantErr         string
		wantCanceled    bool
	}{
		"type: default is structure only": {
			req: dupReq(""), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: structureVars, wantRes: dupRes,
		},
		"type: structure": {
			req: dupReq(DuplicateBoardWithStructure), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: structureVars, wantRes: dupRes,
		},
		"type: pulses": {
			req: dupReq(DuplicateBoardWithPulses), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: dupVars("duplicate_board_with_pulses"), wantRes: dupRes,
		},
		"type: pulses and updates": {
			req: dupReq(DuplicateBoardWithPulsesAndUpdates), replies: []serverReply{ok}, wantCalls: 1,
			wantVars: dupVars("duplicate_board_with_pulses_and_updates"), wantRes: dupRes,
		},
		"retry: HTTP 429 then success": {
			req: dupReq(""), replies: []serverReply{httpErr(429), ok}, wantCalls: 2,
			wantVars: structureVars, wantRes: dupRes,
		},
		"retry: HTTP 500 then success": {
			req: dupReq(""), replies: []serverReply{httpErr(500), ok}, wantCalls: 2,
			wantVars: structureVars, wantRes: dupRes,
		},
		"retry: HTTP 503 then success": {
			req: dupReq(""), replies: []serverReply{httpErr(503), ok}, wantCalls: 2,
			wantVars: structureVars, wantRes: dupRes,
		},
		"retry: complexity budget error then success": {
			req: dupReq(""), replies: []serverReply{complexity, ok}, wantCalls: 2,
			wantVars: structureVars, wantRes: dupRes,
		},
		"retry: two failures then success": {
			req: dupReq(""), replies: []serverReply{httpErr(502), httpErr(504), ok}, wantCalls: 3,
			wantVars: structureVars, wantRes: dupRes,
		},
		"retry: retryable HTTP error exhausts 5 attempts": {
			req: dupReq(""), replies: []serverReply{httpErr(503)}, wantCalls: 5,
			wantVars: structureVars, wantErr: duplicateBoardErrorPrefix,
		},
		"retry: retryable GraphQL error exhausts 5 attempts": {
			req: dupReq(""), replies: []serverReply{complexity}, wantCalls: 5,
			wantVars: structureVars, wantErr: "omplexity",
		},
		"no retry: HTTP 400": {
			req: dupReq(""), replies: []serverReply{httpErr(400)}, wantCalls: 1,
			wantVars: structureVars, wantErr: duplicateBoardErrorPrefix,
		},
		"no retry: HTTP 401": {
			req: dupReq(""), replies: []serverReply{{status: 401, body: unauthorizedResponse}}, wantCalls: 1,
			wantVars: structureVars, wantErr: duplicateBoardErrorPrefix,
		},
		"no retry: HTTP 403": {
			req: dupReq(""), replies: []serverReply{httpErr(403)}, wantCalls: 1,
			wantVars: structureVars, wantErr: duplicateBoardErrorPrefix,
		},
		"no retry: GraphQL board not found in 200 response": {
			req: dupReq(""), replies: []serverReply{{status: 200, body: boardNotFoundResp}}, wantCalls: 1,
			wantVars: structureVars, wantErr: "Board not found",
		},
		"cancel: cancelled context stops retrying": {
			// The server cancels the context as soon as it receives the
			// request and replies with a retryable 503. Run must return the
			// cancellation without a second attempt and without waiting out
			// the hour-long retry delay (otherwise the test would hang).
			req:             &DuplicateBoardReq{TemplateID: "1111111111", BoardName: "vault-51558 copy", retryDelay: time.Hour},
			replies:         []serverReply{httpErr(503)},
			cancelOnRequest: true,
			wantCalls:       1,
			wantVars:        structureVars,
			wantErr:         duplicateBoardErrorPrefix,
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
				require.Equal(t, duplicateBoardMutation, r.Query)
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
