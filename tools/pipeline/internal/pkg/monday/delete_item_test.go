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

// TestDeleteItemReq_Run verifies DeleteItemReq.Run against a minimal httptest
// server standing in for Monday.com. Each case asserts the exact sequence of
// requests the server received, which proves when the delete_item mutation is
// sent and, just as importantly, when it is not.
func TestDeleteItemReq_Run(t *testing.T) {
	t.Parallel()

	const (
		// The read query and mutation that Run() should be sending.
		readItemQuery      = `query ($itemId:ID!){items(ids: [$itemId]){id,name,state,updated_at,board{id},group{id,title},parent_item{id},subitems{id},column_values{id,type,text,value,column{title}}}}`
		deleteItemMutation = `mutation ($itemId:ID!){delete_item(item_id: $itemId){id}}`

		boardID        = "1111111111"
		itemID         = "2222222222"
		sameNameItemID = "3333333333"
		subitemBoardID = "4444444444"
		subitemID      = "6666666666"

		// Response bodies follow the Monday.com API shape. IDs are fake.
		itemResponse         = `{"data":{"items":[{"id":"2222222222","name":"Test item","board":{"id":"1111111111"},"group":{"id":"topics","title":"Group Title"},"parent_item":null,"subitems":[{"id":"5555555555"}]}]}}`
		sameNameItemResponse = `{"data":{"items":[{"id":"3333333333","name":"Test item","board":{"id":"1111111111"},"group":{"id":"topics","title":"Group Title"},"parent_item":null,"subitems":[]}]}}`
		subitemResponse      = `{"data":{"items":[{"id":"6666666666","name":"Test subitem","board":{"id":"4444444444"},"group":{"id":"topics","title":"Subitems"},"parent_item":{"id":"2222222222"},"subitems":[]}]}}`
		notFoundResponse     = `{"data":{"items":[]}}`
		deleteResponse       = `{"data":{"delete_item":{"id":"2222222222"}}}`
		deleteNotFoundResp   = `{"errors":[{"message":"Item not found","locations":[{"line":1,"column":21}],"path":["delete_item"],"extensions":{"code":"InvalidItemIdException","status_code":200,"error_data":{"item_id":2222222222}}}],"data":{"delete_item":null}}`
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
		return wireRequest{Query: readItemQuery, Variables: map[string]any{"itemId": id}}
	}
	del := func(id string) wireRequest {
		return wireRequest{Query: deleteItemMutation, Variables: map[string]any{"itemId": id}}
	}

	// run executes req against an httptest server that answers each request
	// with the next reply (the last reply repeats), and returns the result
	// along with every request the server saw.
	run := func(t *testing.T, req *DeleteItemReq, replies []serverReply) (*DeleteItemRes, []wireRequest, error) {
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

	itemRes := &DeleteItemRes{
		ItemID:       itemID,
		ItemName:     "Test item",
		BoardID:      boardID,
		GroupID:      "topics",
		GroupTitle:   "Group Title",
		SubitemCount: 1,
	}
	deletedItemRes := *itemRes
	deletedItemRes.Deleted = true

	for name, tc := range map[string]struct {
		req          *DeleteItemReq
		replies      []serverReply // the last reply repeats
		wantRequests []wireRequest
		wantRes      *DeleteItemRes
		wantErr      error  // checked with errors.Is
		wantErrMsg   string // checked with ErrorContains, for errors with no sentinel
	}{
		"refused: real run without expect name sends nothing": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID},
			replies:      []serverReply{ok(itemResponse)},
			wantRequests: nil,
			wantErr:      ErrExpectNameRequired,
		},
		"dry run: sends only the read query": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, DryRun: true},
			replies:      []serverReply{ok(itemResponse)},
			wantRequests: []wireRequest{read(itemID)},
			wantRes:      itemRes,
		},
		"dry run: a wrong expect name is still refused": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Other item", DryRun: true},
			replies:      []serverReply{ok(itemResponse)},
			wantRequests: []wireRequest{read(itemID)},
			wantErr:      ErrNameMismatch,
		},
		"refused: name mismatch sends no mutation": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Other item"},
			replies:      []serverReply{ok(itemResponse)},
			wantRequests: []wireRequest{read(itemID)},
			wantErr:      ErrNameMismatch,
		},
		"refused: item on another board sends no mutation": {
			req:          &DeleteItemReq{BoardID: "9999999999", ItemID: itemID, ExpectName: "Test item"},
			replies:      []serverReply{ok(itemResponse)},
			wantRequests: []wireRequest{read(itemID)},
			wantErr:      ErrBoardMismatch,
		},
		"refused: item not found sends no mutation": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Test item"},
			replies:      []serverReply{ok(notFoundResponse)},
			wantRequests: []wireRequest{read(itemID)},
			wantErr:      ErrItemNotFound,
		},
		"same name: deletes only the given id": {
			// Items 2222222222 and 3333333333 are both named "Test item".
			// The id selects the target; the name only confirms it.
			req:          &DeleteItemReq{BoardID: boardID, ItemID: sameNameItemID, ExpectName: "Test item"},
			replies:      []serverReply{ok(sameNameItemResponse), ok(deleteResponse)},
			wantRequests: []wireRequest{read(sameNameItemID), del(sameNameItemID)},
			wantRes: &DeleteItemRes{
				ItemID:     sameNameItemID,
				ItemName:   "Test item",
				BoardID:    boardID,
				GroupID:    "topics",
				GroupTitle: "Group Title",
				Deleted:    true,
			},
		},
		"delete: matching name sends delete_item": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Test item"},
			replies:      []serverReply{ok(itemResponse), ok(deleteResponse)},
			wantRequests: []wireRequest{read(itemID), del(itemID)},
			wantRes:      &deletedItemRes,
		},
		"subitem: board id is the subitem board": {
			req:          &DeleteItemReq{BoardID: subitemBoardID, ItemID: subitemID, ExpectName: "Test subitem"},
			replies:      []serverReply{ok(subitemResponse), ok(deleteResponse)},
			wantRequests: []wireRequest{read(subitemID), del(subitemID)},
			wantRes: &DeleteItemRes{
				ItemID:       subitemID,
				ItemName:     "Test subitem",
				BoardID:      subitemBoardID,
				GroupID:      "topics",
				GroupTitle:   "Subitems",
				ParentItemID: itemID,
				Deleted:      true,
			},
		},
		"retry: HTTP 503 on the read then success": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Test item", retryDelay: time.Millisecond},
			replies:      []serverReply{httpErr(503), ok(itemResponse), ok(deleteResponse)},
			wantRequests: []wireRequest{read(itemID), read(itemID), del(itemID)},
			wantRes:      &deletedItemRes,
		},
		"retry: HTTP 503 on the delete then success": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Test item", retryDelay: time.Millisecond},
			replies:      []serverReply{ok(itemResponse), httpErr(503), ok(deleteResponse)},
			wantRequests: []wireRequest{read(itemID), del(itemID), del(itemID)},
			wantRes:      &deletedItemRes,
		},
		"no retry: HTTP 401 on the read sends no mutation": {
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Test item", retryDelay: time.Millisecond},
			replies:      []serverReply{httpErr(401)},
			wantRequests: []wireRequest{read(itemID)},
			wantErrMsg:   "items query",
		},
		"delete fails: response lost, then the retry finds no item": {
			// The first delete succeeds on Monday.com but the client sees a
			// 503, so it retries. The item is already gone, so the retry gets
			// "Item not found", which is not retryable. Run reports the
			// failure even though the item was deleted.
			req:          &DeleteItemReq{BoardID: boardID, ItemID: itemID, ExpectName: "Test item", retryDelay: time.Millisecond},
			replies:      []serverReply{ok(itemResponse), httpErr(503), ok(deleteNotFoundResp)},
			wantRequests: []wireRequest{read(itemID), del(itemID), del(itemID)},
			wantErrMsg:   "delete_item mutation",
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
