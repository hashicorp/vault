// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/require"
)

// TestGetItemReq_Run verifies GetItemReq.Run against an httptest server
// standing in for Monday.com. A parent item returns its column values and
// subitem ids, and a subitem returns its parent id on the subitem board. Because
// items(ids:) is not scoped to a board, an item that is on a different board is
// rejected with ErrBoardMismatch instead of being returned. A missing item
// returns ErrItemNotFound. A 503 is retried using the retryDelay override, and
// a 401 is not retried.
func TestGetItemReq_Run(t *testing.T) {
	t.Parallel()

	const (
		readItemQuery = `query ($itemId:ID!){items(ids: [$itemId]){id,name,state,updated_at,board{id},group{id,title},parent_item{id},subitems{id},column_values{id,type,text,value,column{title}}}}`

		boardID        = "1111111111"
		itemID         = "2222222222"
		subitemBoardID = "4444444444"
		subitemID      = "6666666666"

		parentResponse   = `{"data":{"items":[{"id":"2222222222","name":"Test item","state":"active","updated_at":"2026-05-01T10:00:00Z","board":{"id":"1111111111"},"group":{"id":"topics","title":"Group Title"},"parent_item":null,"subitems":[{"id":"5555555555"},{"id":"5555555556"}],"column_values":[{"id":"status","type":"status","text":"Done","value":"{\"index\":1}","column":{"title":"Status"}},{"id":"owner","type":"people","text":"","value":null,"column":{"title":"Owner"}}]}]}}`
		subitemResponse  = `{"data":{"items":[{"id":"6666666666","name":"Test subitem","state":"active","updated_at":"2026-05-02T10:00:00Z","board":{"id":"4444444444"},"group":{"id":"topics","title":"Subitems"},"parent_item":{"id":"2222222222"},"subitems":[],"column_values":[]}]}}`
		notFoundResponse = `{"data":{"items":[]}}`
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

	// run executes req against an httptest server that answers each request
	// with the next reply (the last reply repeats), and returns the result
	// along with every request the server saw.
	run := func(t *testing.T, req *GetItemReq, replies []serverReply) (*GetItemRes, []wireRequest, error) {
		t.Helper()

		var reqs []wireRequest
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			reqs = append(reqs, got)
			rep := replies[min(len(reqs)-1, len(replies)-1)]
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(rep.status)
			_, _ = w.Write([]byte(rep.body))
		}))
		t.Cleanup(srv.Close)
		c := graphql.NewClient(srv.URL, srv.Client())
		res, err := req.Run(context.Background(), c)
		return res, reqs, err
	}

	parentRes := &GetItemRes{
		ItemID:       itemID,
		ItemName:     "Test item",
		BoardID:      boardID,
		GroupID:      "topics",
		GroupTitle:   "Group Title",
		State:        "active",
		SubitemCount: 2,
		SubitemIDs:   []string{"5555555555", "5555555556"},
		UpdatedAt:    "2026-05-01T10:00:00Z",
		ColumnValues: []ColumnValue{
			{ID: "status", Title: "Status", Type: "status", Text: "Done", Value: `{"index":1}`},
			{ID: "owner", Title: "Owner", Type: "people"},
		},
	}

	for name, tc := range map[string]struct {
		req          *GetItemReq
		replies      []serverReply // the last reply repeats
		wantRequests int
		wantRes      *GetItemRes
		wantErr      error  // checked with errors.Is
		wantErrMsg   string // checked with ErrorContains, for errors with no sentinel
	}{
		"parent item: returns column values and subitem ids": {
			req:          &GetItemReq{BoardID: boardID, ItemID: itemID},
			replies:      []serverReply{ok(parentResponse)},
			wantRequests: 1,
			wantRes:      parentRes,
		},
		"subitem: returns the parent id on the subitem board": {
			req:          &GetItemReq{BoardID: subitemBoardID, ItemID: subitemID},
			replies:      []serverReply{ok(subitemResponse)},
			wantRequests: 1,
			wantRes: &GetItemRes{
				ItemID:       subitemID,
				ItemName:     "Test subitem",
				BoardID:      subitemBoardID,
				GroupID:      "topics",
				GroupTitle:   "Subitems",
				State:        "active",
				ParentItemID: itemID,
				SubitemIDs:   []string{},
				UpdatedAt:    "2026-05-02T10:00:00Z",
				ColumnValues: []ColumnValue{},
			},
		},
		"board mismatch: item is on another board": {
			req:          &GetItemReq{BoardID: "9999999999", ItemID: itemID},
			replies:      []serverReply{ok(parentResponse)},
			wantRequests: 1,
			wantErr:      ErrBoardMismatch,
		},
		"board mismatch: subitem id with the parent board id": {
			req:          &GetItemReq{BoardID: boardID, ItemID: subitemID},
			replies:      []serverReply{ok(subitemResponse)},
			wantRequests: 1,
			wantErr:      ErrBoardMismatch,
		},
		"not found: no items returned": {
			req:          &GetItemReq{BoardID: boardID, ItemID: itemID},
			replies:      []serverReply{ok(notFoundResponse)},
			wantRequests: 1,
			wantErr:      ErrItemNotFound,
		},
		"retry: HTTP 503 then success": {
			req:          &GetItemReq{BoardID: boardID, ItemID: itemID, retryDelay: time.Millisecond},
			replies:      []serverReply{httpErr(503), ok(parentResponse)},
			wantRequests: 2,
			wantRes:      parentRes,
		},
		"no retry: HTTP 401": {
			req:          &GetItemReq{BoardID: boardID, ItemID: itemID, retryDelay: time.Millisecond},
			replies:      []serverReply{httpErr(401)},
			wantRequests: 1,
			wantErrMsg:   "items query",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, reqs, err := run(t, tc.req, tc.replies)

			// Every request must be the same read query for the requested item.
			require.Len(t, reqs, tc.wantRequests)
			for _, got := range reqs {
				require.Equal(t, readItemQuery, got.Query)
				require.Equal(t, map[string]any{"itemId": tc.req.ItemID}, got.Variables)
			}

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
