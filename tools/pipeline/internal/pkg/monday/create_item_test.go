// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/require"
)

// TestCreateItemReq_Run verifies CreateItemReq.Run against a minimal httptest
// server standing in for Monday.com. It validates that the expected GraphQL
// mutation and variables are sent (including handling of optional group ID),
// that successful responses map into CreateItemRes, that non-retryable errors
// fail immediately after a single attempt, and that transient retryable errors
// (rate limits and complexity limits only) retry up to the maximum attempts.
//
// This test suite also directly validates the non-idempotency safety invariant:
// 5xx gateway codes (500, 502, 503, 504) must NOT be retried for create_item
// because Monday.com may have already committed the item before the gateway
// returned the error — retrying would create a duplicate.
func TestCreateItemReq_Run(t *testing.T) {
	t.Parallel()

	const wantMutation = `mutation ($boardId:ID!$groupId:String$itemName:String!){create_item(board_id: $boardId, group_id: $groupId, item_name: $itemName){id,name}}`

	type wireRequest struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}

	t.Run("success with group ID", func(t *testing.T) {
		t.Parallel()

		var reqs []wireRequest
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			reqs = append(reqs, got)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"create_item":{"id":"987654","name":"New Item"}}}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:  "12345",
			GroupID:  "topics",
			ItemName: "New Item",
		}

		res, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.Len(t, reqs, 1, "expected exactly one request")
		require.Equal(t, wantMutation, reqs[0].Query)
		require.Equal(t, map[string]any{
			"boardId":  "12345",
			"groupId":  "topics",
			"itemName": "New Item",
		}, reqs[0].Variables)

		require.Equal(t, "987654", res.ItemID)
		require.Equal(t, "New Item", res.ItemName)
	})

	t.Run("success without group ID sends null groupId", func(t *testing.T) {
		t.Parallel()

		var reqs []wireRequest
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			reqs = append(reqs, got)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"create_item":{"id":"987654","name":"New Item"}}}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:  "12345",
			ItemName: "New Item",
		}

		res, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.Len(t, reqs, 1, "expected exactly one request")
		require.Equal(t, wantMutation, reqs[0].Query)
		require.Equal(t, map[string]any{
			"boardId":  "12345",
			"groupId":  nil,
			"itemName": "New Item",
		}, reqs[0].Variables)

		require.Equal(t, "987654", res.ItemID)
		require.Equal(t, "New Item", res.ItemName)
	})

	// Without a status the request must be byte-for-byte what it was before
	// the Status field existed: no column_values in the mutation or variables.
	// The other success tests assert the same through wantMutation, so this
	// one pins the legacy shape explicitly.
	t.Run("empty status sends the legacy mutation without column_values", func(t *testing.T) {
		t.Parallel()

		const legacyMutation = `mutation ($boardId:ID!$groupId:String$itemName:String!){create_item(board_id: $boardId, group_id: $groupId, item_name: $itemName){id,name}}`

		var reqs []wireRequest
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			reqs = append(reqs, got)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"create_item":{"id":"987654","name":"New Item"}}}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{BoardID: "12345", ItemName: "New Item"}

		_, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.Len(t, reqs, 1)
		require.Equal(t, legacyMutation, reqs[0].Query)
		require.NotContains(t, reqs[0].Variables, "columnValues")
	})

	// A status is sent in a single call as column_values, a JSON string that
	// sets the label on the "status" column.
	t.Run("status is sent as column_values in a single call", func(t *testing.T) {
		t.Parallel()

		const statusMutation = `mutation ($boardId:ID!$columnValues:JSON!$groupId:String$itemName:String!){create_item(board_id: $boardId, group_id: $groupId, item_name: $itemName, column_values: $columnValues){id,name}}`

		for name, tc := range map[string]struct {
			groupID    string
			status     string
			wantGroup  any
			wantValues string
		}{
			"label with group":              {groupID: "topics", status: "Not started", wantGroup: "topics", wantValues: `{"status":{"label":"Not started"}}`},
			"label without group":           {status: "Not started", wantGroup: nil, wantValues: `{"status":{"label":"Not started"}}`},
			"label case is preserved":       {status: "Not Started", wantGroup: nil, wantValues: `{"status":{"label":"Not Started"}}`},
			"quotes in label are escaped":   {status: `Say "hi"`, wantGroup: nil, wantValues: `{"status":{"label":"Say \"hi\""}}`},
			"backslash in label is escaped": {status: `a\b`, wantGroup: nil, wantValues: `{"status":{"label":"a\\b"}}`},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				var reqs []wireRequest
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var got wireRequest
					_ = json.NewDecoder(r.Body).Decode(&got)
					reqs = append(reqs, got)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"data":{"create_item":{"id":"987654","name":"New Item"}}}`))
				}))
				t.Cleanup(srv.Close)

				client := graphql.NewClient(srv.URL, srv.Client())
				req := &CreateItemReq{BoardID: "12345", GroupID: tc.groupID, ItemName: "New Item", Status: tc.status}

				res, err := req.Run(context.Background(), client)
				require.NoError(t, err)
				require.Len(t, reqs, 1, "expected a single call, not create then update")
				require.Equal(t, statusMutation, reqs[0].Query)
				require.Equal(t, map[string]any{
					"boardId":      "12345",
					"columnValues": tc.wantValues,
					"groupId":      tc.wantGroup,
					"itemName":     "New Item",
				}, reqs[0].Variables)
				require.Equal(t, "987654", res.ItemID)
				require.Equal(t, "New Item", res.ItemName)
			})
		}
	})

	// Setting a status must not widen the retry policy: a gateway error can
	// follow a committed create, and retrying would duplicate the item.
	t.Run("HTTP 503 with a status is still not retried", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{BoardID: "12345", ItemName: "New Item", Status: "Not started", retryDelay: time.Millisecond}

		res, err := req.Run(context.Background(), client)
		require.Nil(t, res)
		require.ErrorContains(t, err, "create_item mutation")
		require.Equal(t, int32(1), attempts.Load(), "expected HTTP 503 to fail after 1 attempt")
	})

	// A rate limit is rejected before execution, so it is safe to retry, and
	// every attempt must carry the same status.
	t.Run("HTTP 429 with a status is retried and keeps the status", func(t *testing.T) {
		t.Parallel()

		var (
			attempts atomic.Int32
			mu       sync.Mutex
			reqs     []wireRequest
		)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var got wireRequest
			_ = json.NewDecoder(r.Body).Decode(&got)
			mu.Lock()
			reqs = append(reqs, got)
			mu.Unlock()
			if attempts.Add(1) < 3 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"create_item":{"id":"987654","name":"New Item"}}}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{BoardID: "12345", ItemName: "New Item", Status: "Not started", retryDelay: time.Millisecond}

		res, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.Equal(t, "987654", res.ItemID)

		mu.Lock()
		defer mu.Unlock()
		require.Len(t, reqs, 3)
		for _, r := range reqs {
			require.Equal(t, `{"status":{"label":"Not started"}}`, r.Variables["columnValues"])
		}
	})

	// Monday.com rejects an unknown or wrong-case label with a GraphQL error
	// (for example missingLabel). Retrying cannot fix it, so it must surface
	// after a single attempt.
	t.Run("rejected status label fails immediately without retrying", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"errors":[{"message":"missingLabel","extensions":{"code":"ColumnValueException"}}],"data":{"create_item":null}}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{BoardID: "12345", ItemName: "New Item", Status: "not started", retryDelay: time.Millisecond}

		res, err := req.Run(context.Background(), client)
		require.Nil(t, res)
		require.ErrorContains(t, err, "create_item mutation")
		require.ErrorContains(t, err, "missingLabel")
		require.Equal(t, int32(1), attempts.Load())
	})

	t.Run("non-retryable graphql error fails immediately without retrying", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"errors":[{"message":"Not authenticated"}]}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:    "12345",
			ItemName:   "New Item",
			retryDelay: time.Millisecond,
		}

		res, err := req.Run(context.Background(), client)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorContains(t, err, "create_item mutation")
		require.ErrorContains(t, err, "Not authenticated")
		require.Equal(t, int32(1), attempts.Load(), "expected non-retryable error to fail after 1 attempt")
	})

	t.Run("non-retryable HTTP 400 fails immediately without retrying", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusBadRequest)
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:    "12345",
			ItemName:   "New Item",
			retryDelay: time.Millisecond,
		}

		res, err := req.Run(context.Background(), client)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorContains(t, err, "create_item mutation")
		require.Equal(t, int32(1), attempts.Load(), "expected HTTP 400 to fail after 1 attempt")
	})

	t.Run("retryable HTTP 429 recovers on subsequent attempt", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := attempts.Add(1)
			if count < 3 {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"create_item":{"id":"987654","name":"New Item"}}}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:    "12345",
			ItemName:   "New Item",
			retryDelay: time.Millisecond,
		}

		res, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "987654", res.ItemID)
		require.Equal(t, "New Item", res.ItemName)
		require.Equal(t, int32(3), attempts.Load(), "expected 3 attempts before succeeding")
	})

	// HTTP 503 is included here to confirm that create_item uses the narrow
	// isCreateRetryableError policy, not the broader isRetryableError.  A
	// gateway error is not safe to retry for a non-idempotent mutation because
	// Monday.com may have already committed the item before the gateway
	// returned the error.  The request must fail immediately on the first
	// attempt — it must not retry and potentially create a duplicate item.
	t.Run("HTTP 503 is not retried for create_item and fails immediately", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:    "12345",
			ItemName:   "New Item",
			retryDelay: time.Millisecond,
		}

		res, err := req.Run(context.Background(), client)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorContains(t, err, "create_item mutation")
		// Exactly one attempt: 503 must not be retried for a create mutation.
		require.Equal(t, int32(1), attempts.Load(), "expected HTTP 503 to fail after 1 attempt, not be retried")
	})

	// Regression test for the committed-then-gateway-error scenario described
	// in the code review: Monday.com commits the item and then a 504 is
	// returned by the gateway.  Run must surface the 504 as an error and must
	// not retry, because the item has already been created.  Retrying would
	// result in a second, duplicate item being created.
	t.Run("HTTP 504 after committed item is not retried and returns error", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			// Simulate Monday.com committing the item but the gateway timing
			// out before the success response reaches the client.
			w.WriteHeader(http.StatusGatewayTimeout)
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:    "12345",
			ItemName:   "New Item",
			retryDelay: time.Millisecond,
		}

		res, err := req.Run(context.Background(), client)
		require.Nil(t, res, "no result should be returned when the gateway reports a timeout")
		require.Error(t, err)
		require.ErrorContains(t, err, "create_item mutation")
		// Exactly one attempt: 504 must not be retried for a create mutation.
		require.Equal(t, int32(1), attempts.Load(), "expected HTTP 504 to fail after 1 attempt, not be retried")
	})

	// HTTP 429 is the one 4xx code that is safe to retry for create_item:
	// Monday.com rejects the request before executing it, so no item is
	// created.  Confirm that 429 still exhausts all retry attempts.
	t.Run("HTTP 429 exhausts all 5 attempts and returns error", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:    "12345",
			ItemName:   "New Item",
			retryDelay: time.Millisecond,
		}

		res, err := req.Run(context.Background(), client)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorContains(t, err, "create_item mutation")
		// 429 is safe to retry for create mutations: Monday.com rejects before executing.
		require.Equal(t, int32(5), attempts.Load(), "expected HTTP 429 to exhaust all 5 attempts")
	})

	t.Run("retryable GraphQL complexity error recovers on subsequent attempt", func(t *testing.T) {
		t.Parallel()

		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := attempts.Add(1)
			w.Header().Set("Content-Type", "application/json")
			if count == 1 {
				_, _ = w.Write([]byte(`{"errors":[{"message":"Complexity budget exhausted, retry after 10 seconds"}]}`))
				return
			}
			_, _ = w.Write([]byte(`{"data":{"create_item":{"id":"987654","name":"New Item"}}}`))
		}))
		t.Cleanup(srv.Close)

		client := graphql.NewClient(srv.URL, srv.Client())
		req := &CreateItemReq{
			BoardID:    "12345",
			ItemName:   "New Item",
			retryDelay: time.Millisecond,
		}

		res, err := req.Run(context.Background(), client)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "987654", res.ItemID)
		require.Equal(t, "New Item", res.ItemName)
		require.Equal(t, int32(2), attempts.Load(), "expected 2 attempts before succeeding")
	})
}

// TestCreateItemRes_ToJSON verifies that ToJSON produces valid JSON
// containing the expected item_id and item_name fields.
func TestCreateItemRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &CreateItemRes{
		ItemID:   "987654",
		ItemName: "Release Checklist",
	}

	b, err := res.ToJSON()
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&got))
	require.Equal(t, "987654", got["item_id"])
	require.Equal(t, "Release Checklist", got["item_name"])
}

// TestCreateItemRes_ToTable verifies that ToTable renders the item ID and
// item name in the table output.
func TestCreateItemRes_ToTable(t *testing.T) {
	t.Parallel()

	res := &CreateItemRes{
		ItemID:   "987654",
		ItemName: "Release Checklist",
	}

	output := res.ToTable().Render()
	require.Contains(t, output, "ITEM_ID")
	require.Contains(t, output, "ITEM_NAME")
	require.Contains(t, output, "987654")
	require.Contains(t, output, "Release Checklist")
}
