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

// TestCreateSubitemReq_Run verifies CreateSubitemReq.Run against a minimal
// httptest server standing in for Monday.com. It asserts the exact mutation
// and variables sent on the wire (ID!, String! and JSON! variable types, with
// column_values always an empty JSON object), that the response is mapped into
// CreateSubitemRes, and the retry behavior: 429 and 5xx responses and
// complexity-budget errors are retried and then succeed, retryable errors give
// up after 5 attempts, non-retryable HTTP and GraphQL errors (HTTP 200 with an
// errors array, such as a missing parent item) are sent exactly once, and
// cancelling the context stops further retries.
func TestCreateSubitemReq_Run(t *testing.T) {
	t.Parallel()

	const (
		// The mutation that Run() should be sending
		createSubitemMutation = `mutation ($columnValues:JSON!$itemName:String!$parentItemId:ID!){create_subitem(parent_item_id: $parentItemId, item_name: $itemName, column_values: $columnValues){id}}`

		// Response bodies captured from the Monday.com API on a
		// throwaway board, except complexityResponse, which is a placeholder.
		// IDs and request IDs were replaced with fake values.
		successResponse      = `{"data":{"create_subitem":{"id":"3333333333"}},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		unauthorizedResponse = `{"errors":["Not Authenticated"]}`
		itemNotFoundResp     = `{"errors":[{"message":"Item not found","locations":[{"line":1,"column":43}],"path":["create_subitem"],"extensions":{"code":"InvalidItemIdException","status_code":200,"error_data":{"item_id":null}}}],"data":{"create_subitem":null},"extensions":{"request_id":"00000000-0000-0000-0000-000000000000"}}`
		complexityResponse   = `{"errors":[{"message":"Complexity budget exhausted, retry after 10 seconds"}]}`

		createSubitemErrorPrefix = "create_subitem mutation"
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
	run := func(t *testing.T, req *CreateSubitemReq, replies []serverReply, cancelOnRequest bool) (*CreateSubitemRes, []wireRequest, error) {
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

	subitemReq := &CreateSubitemReq{
		ParentItemID: "1210003",
		ItemName:     "Release Checklist",
		retryDelay:   time.Millisecond,
	}
	subitemVars := map[string]any{
		"parentItemId": "1210003",
		"itemName":     "Release Checklist",
		"columnValues": "{}",
	}
	subitemRes := &CreateSubitemRes{SubitemID: "3333333333"}

	for name, tc := range map[string]struct {
		req             *CreateSubitemReq
		replies         []serverReply // the last reply repeats
		cancelOnRequest bool
		wantCalls       int
		wantRes         *CreateSubitemRes
		wantErr         string
		wantCanceled    bool
	}{
		"success: response mapped into result": {
			req: subitemReq, replies: []serverReply{ok}, wantCalls: 1, wantRes: subitemRes,
		},
		"retry: HTTP 429 then success": {
			req: subitemReq, replies: []serverReply{httpErr(429), ok}, wantCalls: 2, wantRes: subitemRes,
		},
		"retry: HTTP 500 then success": {
			req: subitemReq, replies: []serverReply{httpErr(500), ok}, wantCalls: 2, wantRes: subitemRes,
		},
		"retry: HTTP 503 then success": {
			req: subitemReq, replies: []serverReply{httpErr(503), ok}, wantCalls: 2, wantRes: subitemRes,
		},
		"retry: complexity budget error then success": {
			req: subitemReq, replies: []serverReply{complexity, ok}, wantCalls: 2, wantRes: subitemRes,
		},
		"retry: two failures then success": {
			req: subitemReq, replies: []serverReply{httpErr(502), httpErr(504), ok}, wantCalls: 3, wantRes: subitemRes,
		},
		"retry: retryable HTTP error exhausts 5 attempts": {
			req: subitemReq, replies: []serverReply{httpErr(503)}, wantCalls: 5, wantErr: createSubitemErrorPrefix,
		},
		"retry: retryable GraphQL error exhausts 5 attempts": {
			req: subitemReq, replies: []serverReply{complexity}, wantCalls: 5, wantErr: "omplexity",
		},
		"no retry: HTTP 400": {
			req: subitemReq, replies: []serverReply{httpErr(400)}, wantCalls: 1, wantErr: createSubitemErrorPrefix,
		},
		"no retry: HTTP 401": {
			req: subitemReq, replies: []serverReply{{status: 401, body: unauthorizedResponse}}, wantCalls: 1,
			wantErr: createSubitemErrorPrefix,
		},
		"no retry: HTTP 403": {
			req: subitemReq, replies: []serverReply{httpErr(403)}, wantCalls: 1, wantErr: createSubitemErrorPrefix,
		},
		"no retry: GraphQL item not found in 200 response": {
			req: subitemReq, replies: []serverReply{{status: 200, body: itemNotFoundResp}}, wantCalls: 1,
			wantErr: "Item not found",
		},
		"cancel: cancelled context stops retrying": {
			// The server cancels the context as soon as it receives the
			// request and replies with a retryable 503. Run must return the
			// cancellation without a second attempt and without waiting out
			// the hour-long retry delay (otherwise the test would hang).
			req:             &CreateSubitemReq{ParentItemID: "1210003", ItemName: "Release Checklist", retryDelay: time.Hour},
			replies:         []serverReply{httpErr(503)},
			cancelOnRequest: true,
			wantCalls:       1,
			wantErr:         createSubitemErrorPrefix,
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
				require.Equal(t, createSubitemMutation, r.Query)
				require.Equal(t, subitemVars, r.Variables)
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
