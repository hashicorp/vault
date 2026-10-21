// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/require"
)

// TestGetItemIDReq_Run verifies GetItemIDReq.Run against a minimal httptest
// server standing in for Monday.com, for both the --board-id and --parent-id
// lookups. It asserts the exact query and variables sent on the wire, that the
// response is mapped into ItemID and BoardID (the board path echoes the
// requested board, the parent path returns the subitem board), that duplicate
// names return the first match, that names are matched exactly (including
// whitespace), and that a missing board, parent or item, a GraphQL error in a
// 200 response, and an HTTP error each return an error after a single request.
func TestGetItemIDReq_Run(t *testing.T) {
	t.Parallel()

	const (
		// The queries that Run() should be sending for the different cases
		boardQuery  = `query ($boardId:ID!){boards(ids:[$boardId]){items_page(limit:100){items{id,name}}}}`
		parentQuery = `query ($parentId:ID!){items(ids:[$parentId]){subitems{id,name,board{id}}}}`

		// The responses Run() should receive from Monday for the different cases
		boardResponse = `{"data":{"boards":[{"items_page":{"items":[
			{"id":"111","name":"Other item"},
			{"id":"222","name":"Release Checklist"},
			{"id":"333","name":"Release Checklist"}
		]}}]}}`

		parentResponse = `{"data":{"items":[{"subitems":[
			{"id":"11858417779","name":"Another subitem","board":{"id":"18410441433"}},
			{"id":"11858417780","name":"If updating,  update the Go Version for CE","board":{"id":"18410441433"}},
			{"id":"11858417781","name":"If updating,  update the Go Version for CE","board":{"id":"18410441433"}}
		]}]}}`
	)

	// wireRequest is the GraphQL request body that Run() sends to the server.
	type wireRequest struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}

	// lookup runs req against an httptest server that replies with status and
	// body, and returns the result along with every request the server saw.
	// get-item-id does not retry, so callers expect exactly one request.
	lookup := func(t *testing.T, req *GetItemIDReq, status int, body string) (*GetItemIDRes, []wireRequest, error) {
		t.Helper()
		var reqs []wireRequest
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			reqs = append(reqs, got)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		}))
		t.Cleanup(srv.Close)
		c := graphql.NewClient(srv.URL, srv.Client())
		res, err := req.Run(context.Background(), c)
		return res, reqs, err
	}

	boardReq := func(itemName string) *GetItemIDReq {
		return &GetItemIDReq{BoardID: "1010222", ItemName: itemName}
	}
	parentReq := func(itemName string) *GetItemIDReq {
		return &GetItemIDReq{ParentID: "5050606", ItemName: itemName}
	}
	boardVars := map[string]any{"boardId": "1010222"}
	parentVars := map[string]any{"parentId": "5050606"}

	for name, tc := range map[string]struct {
		req       *GetItemIDReq
		status    int
		body      string
		wantQuery string
		wantVars  map[string]any
		wantItem  string
		wantBoard string
		wantErr   string
	}{
		"board: matching item": {
			req: boardReq("Other item"), status: 200, body: boardResponse,
			wantQuery: boardQuery, wantVars: boardVars, wantItem: "111", wantBoard: "1010222",
		},
		"board: duplicate names return first match": {
			req: boardReq("Release Checklist"), status: 200, body: boardResponse,
			wantQuery: boardQuery, wantVars: boardVars, wantItem: "222", wantBoard: "1010222",
		},
		"board: item not found": {
			req: boardReq("Missing"), status: 200, body: boardResponse,
			wantQuery: boardQuery, wantVars: boardVars, wantErr: "couldn't find id for item Missing",
		},
		"board: different whitespace does not match": {
			req: boardReq("Release  Checklist"), status: 200, body: boardResponse,
			wantQuery: boardQuery, wantVars: boardVars, wantErr: "couldn't find id for item Release  Checklist",
		},
		"board: board not found": {
			req: boardReq("Release Checklist"), status: 200, body: `{"data":{"boards":[]}}`,
			wantQuery: boardQuery, wantVars: boardVars, wantErr: "board 1010222 not found",
		},
		"board: graphql error in 200 response": {
			req: boardReq("Release Checklist"), status: 200, body: `{"errors":[{"message":"Not authenticated"}]}`,
			wantQuery: boardQuery, wantVars: boardVars, wantErr: "Not authenticated",
		},
		"board: HTTP 401": {
			req: boardReq("Release Checklist"), status: 401, body: `{}`,
			wantQuery: boardQuery, wantVars: boardVars, wantErr: "querying item by item-name",
		},
		"board: HTTP 500": {
			req: boardReq("Release Checklist"), status: 500, body: `{}`,
			wantQuery: boardQuery, wantVars: boardVars, wantErr: "querying item by item-name",
		},
		"parent: matching subitem": {
			req: parentReq("Another subitem"), status: 200, body: parentResponse,
			wantQuery: parentQuery, wantVars: parentVars, wantItem: "11858417779", wantBoard: "18410441433",
		},
		"parent: duplicate names return first match": {
			req: parentReq("If updating,  update the Go Version for CE"), status: 200, body: parentResponse,
			wantQuery: parentQuery, wantVars: parentVars, wantItem: "11858417780", wantBoard: "18410441433",
		},
		"parent: different whitespace does not match": {
			req: parentReq("If updating, update the Go Version for CE"), status: 200, body: parentResponse,
			wantQuery: parentQuery, wantVars: parentVars, wantErr: "couldn't find subitem",
		},
		"parent: subitem not found": {
			req: parentReq("Missing"), status: 200, body: parentResponse,
			wantQuery: parentQuery, wantVars: parentVars, wantErr: "couldn't find subitem Missing under parent item 5050606",
		},
		"parent: no subitems": {
			req: parentReq("If updating,  update the Go Version for CE"), status: 200, body: `{"data":{"items":[{"subitems":[]}]}}`,
			wantQuery: parentQuery, wantVars: parentVars, wantErr: "couldn't find subitem",
		},
		"parent: parent not found": {
			req: parentReq("If updating,  update the Go Version for CE"), status: 200, body: `{"data":{"items":[]}}`,
			wantQuery: parentQuery, wantVars: parentVars, wantErr: "parent item 5050606 not found",
		},
		"parent: graphql error in 200 response": {
			req: parentReq("If updating,  update the Go Version for CE"), status: 200, body: `{"errors":[{"message":"Item not found"}]}`,
			wantQuery: parentQuery, wantVars: parentVars, wantErr: "Item not found",
		},
		"parent: HTTP 401": {
			req: parentReq("If updating,  update the Go Version for CE"), status: 401, body: `{}`,
			wantQuery: parentQuery, wantVars: parentVars, wantErr: "querying subitem by parent-id",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, reqs, err := lookup(t, tc.req, tc.status, tc.body)

			// Make sure no retries on query
			require.Len(t, reqs, 1, "expected exactly one request")
			// Make sure the correct query and variables are sent given either the board-id or parent-id request
			require.Equal(t, tc.wantQuery, reqs[0].Query)
			require.Equal(t, tc.wantVars, reqs[0].Variables)

			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			// Make sure server sends back correct ItemID
			require.Equal(t, tc.wantItem, res.ItemID)
			// Make sure server sends back correct BoardID
			require.Equal(t, tc.wantBoard, res.BoardID)
		})
	}
}
