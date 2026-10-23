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

// TestUpdateItemReq_Run verifies UpdateItemReq.Run against a minimal httptest
// server standing in for Monday.com. It asserts the exact mutation and
// variables sent on the wire, that the response is mapped into UpdateItemRes
// (with the requested value echoed as NewValue), and the retry behavior: 429
// and 5xx responses and complexity-budget errors are retried and then succeed,
// retryable errors give up after 5 attempts, non-retryable HTTP and GraphQL
// errors (HTTP 200 with an errors array, including a nonexistent status label
// reported as ColumnValueException/missingLabel) are sent exactly once, and
// cancelling the context stops further retries.
func TestUpdateItemReq_Run(t *testing.T) {
	t.Parallel()

	const (
		// The mutation that Run() should be sending
		updateItemMutation = `mutation ($boardId:ID!$columnId:String!$itemId:ID!$value:String!){change_simple_column_value(board_id: $boardId, item_id: $itemId, column_id: $columnId, value: $value){id,name}}`

		// Response bodies captured from the Monday.com API on a
		// throwaway board, except complexityResponse, which is a placeholder.
		// IDs and request IDs were replaced with fake values.
		successResponse      = `{"data":{"change_simple_column_value":{"id":"2222222222","name":"Item 1"}},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		unauthorizedResponse = `{"errors":["Not Authenticated"]}`
		missingLabelResp     = `{"errors":[{"message":"This status label doesn't exist, possible statuses are: {0: Working on it, 1: Done, 2: Stuck}","locations":[{"line":1,"column":54}],"path":["change_simple_column_value"],"extensions":{"code":"ColumnValueException","status_code":200,"error_data":{"column_type":"color","column_id":"status","column_value":"{\"label\" => \"NoSuchLabel\"}","column_validation_error_code":"missingLabel","column_name":"Status"}}}],"data":{"change_simple_column_value":null},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		itemNotFoundResp     = `{"errors":[{"message":"Item not found","locations":[{"line":1,"column":54}],"path":["change_simple_column_value"],"extensions":{"code":"InvalidItemIdException","status_code":200,"error_data":{"item_id":999999999999}}}],"data":{"change_simple_column_value":null},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		complexityResponse   = `{"errors":[{"message":"Complexity budget exhausted, retry after 10 seconds"}]}`

		updateItemErrorPrefix = "change_simple_column_value mutation"
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
	run := func(t *testing.T, req *UpdateItemReq, replies []serverReply, cancelOnRequest bool) (*UpdateItemRes, []wireRequest, error) {
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

	updateReq := &UpdateItemReq{
		BoardID:    "1111111111",
		ItemID:     "2222222222",
		ColumnID:   "status",
		Value:      "Done",
		retryDelay: time.Millisecond,
	}
	updateVars := map[string]any{
		"boardId":  "1111111111",
		"itemId":   "2222222222",
		"columnId": "status",
		"value":    "Done",
	}
	updateRes := &UpdateItemRes{ItemID: "2222222222", ItemName: "Item 1", NewValue: "Done"}

	for name, tc := range map[string]struct {
		req             *UpdateItemReq
		replies         []serverReply // the last reply repeats
		cancelOnRequest bool
		wantCalls       int
		wantRes         *UpdateItemRes
		wantErr         string
		wantCanceled    bool
	}{
		"success: response mapped into result": {
			req: updateReq, replies: []serverReply{ok}, wantCalls: 1, wantRes: updateRes,
		},
		"retry: HTTP 429 then success": {
			req: updateReq, replies: []serverReply{httpErr(429), ok}, wantCalls: 2, wantRes: updateRes,
		},
		"retry: HTTP 500 then success": {
			req: updateReq, replies: []serverReply{httpErr(500), ok}, wantCalls: 2, wantRes: updateRes,
		},
		"retry: HTTP 503 then success": {
			req: updateReq, replies: []serverReply{httpErr(503), ok}, wantCalls: 2, wantRes: updateRes,
		},
		"retry: complexity budget error then success": {
			req: updateReq, replies: []serverReply{complexity, ok}, wantCalls: 2, wantRes: updateRes,
		},
		"retry: two failures then success": {
			req: updateReq, replies: []serverReply{httpErr(502), httpErr(504), ok}, wantCalls: 3, wantRes: updateRes,
		},
		"retry: retryable HTTP error exhausts 5 attempts": {
			req: updateReq, replies: []serverReply{httpErr(503)}, wantCalls: 5, wantErr: updateItemErrorPrefix,
		},
		"retry: retryable GraphQL error exhausts 5 attempts": {
			req: updateReq, replies: []serverReply{complexity}, wantCalls: 5, wantErr: "omplexity",
		},
		"no retry: HTTP 400": {
			req: updateReq, replies: []serverReply{httpErr(400)}, wantCalls: 1, wantErr: updateItemErrorPrefix,
		},
		"no retry: HTTP 401": {
			req: updateReq, replies: []serverReply{{status: 401, body: unauthorizedResponse}}, wantCalls: 1,
			wantErr: updateItemErrorPrefix,
		},
		"no retry: HTTP 403": {
			req: updateReq, replies: []serverReply{httpErr(403)}, wantCalls: 1, wantErr: updateItemErrorPrefix,
		},
		"no retry: GraphQL missing status label in 200 response": {
			req: updateReq, replies: []serverReply{{status: 200, body: missingLabelResp}}, wantCalls: 1,
			wantErr: "This status label doesn't exist",
		},
		"no retry: GraphQL item not found in 200 response": {
			req: updateReq, replies: []serverReply{{status: 200, body: itemNotFoundResp}}, wantCalls: 1,
			wantErr: "Item not found",
		},
		"cancel: cancelled context stops retrying": {
			// The server cancels the context as soon as it receives the
			// request and replies with a retryable 503. Run must return the
			// cancellation without a second attempt and without waiting out
			// the hour-long retry delay (otherwise the test would hang).
			req: &UpdateItemReq{
				BoardID:    "1111111111",
				ItemID:     "2222222222",
				ColumnID:   "status",
				Value:      "Done",
				retryDelay: time.Hour,
			},
			replies:         []serverReply{httpErr(503)},
			cancelOnRequest: true,
			wantCalls:       1,
			wantErr:         updateItemErrorPrefix,
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
				require.Equal(t, updateItemMutation, r.Query)
				require.Equal(t, updateVars, r.Variables)
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
