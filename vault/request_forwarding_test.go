// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

// Test_haIDFromContext verifies that the HA node ID gets correctly extracted
// from a gRPC context
func Test_haIDFromContext(t *testing.T) {
	testCases := []struct {
		name   string
		md     metadata.MD
		wantID string
		wantOk bool
	}{
		{
			name:   "no ID",
			md:     metadata.MD{},
			wantID: "",
			wantOk: false,
		},
		{
			name:   "with ID",
			md:     metadata.MD{haNodeIDKey: {"node_id"}},
			wantID: "node_id",
			wantOk: true,
		},
		{
			name:   "with empty string ID",
			md:     metadata.MD{haNodeIDKey: {""}},
			wantID: "",
			wantOk: true,
		},
		{
			name:   "with empty ID",
			md:     metadata.MD{haNodeIDKey: {}},
			wantID: "",
			wantOk: false,
		},

		{
			name:   "with multiple IDs",
			md:     metadata.MD{haNodeIDKey: {"1", "2"}},
			wantID: "1",
			wantOk: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), tc.md)
			id, ok := haIDFromContext(ctx)
			require.Equal(t, tc.wantID, id)
			require.Equal(t, tc.wantOk, ok)
		})
	}
}

type mockHARemovableNodeBackend struct {
	physical.RemovableNodeHABackend
	isRemoved func(context.Context, string) (bool, error)
}

func (m *mockHARemovableNodeBackend) IsNodeRemoved(ctx context.Context, nodeID string) (bool, error) {
	return m.isRemoved(ctx, nodeID)
}

func newMockHARemovableNodeBackend(isRemoved func(context.Context, string) (bool, error)) physical.RemovableNodeHABackend {
	return &mockHARemovableNodeBackend{isRemoved: isRemoved}
}

// Test_haMembershipServerCheck verifies that the correct error is returned
// when the context contains a removed node ID
func Test_haMembershipServerCheck(t *testing.T) {
	nodeIDCtx := metadata.NewIncomingContext(context.Background(), metadata.MD{haNodeIDKey: {"node_id"}})
	otherErr := errors.New("error checking")
	testCases := []struct {
		name      string
		nodeIDCtx context.Context
		haBackend physical.RemovableNodeHABackend
		wantError error
	}{
		{
			name:      "nil backend",
			haBackend: nil,
			nodeIDCtx: nodeIDCtx,
		}, {
			name: "no node ID context",
			haBackend: newMockHARemovableNodeBackend(func(ctx context.Context, s string) (bool, error) {
				return false, nil
			}),
			nodeIDCtx: context.Background(),
		}, {
			name: "node removed",
			haBackend: newMockHARemovableNodeBackend(func(ctx context.Context, s string) (bool, error) {
				return true, nil
			}),
			nodeIDCtx: nodeIDCtx,
			wantError: StatusNotHAMember,
		}, {
			name: "node removed err",
			haBackend: newMockHARemovableNodeBackend(func(ctx context.Context, s string) (bool, error) {
				return false, otherErr
			}),
			nodeIDCtx: nodeIDCtx,
			wantError: otherErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Core{
				logger: hclog.NewNullLogger(),
			}
			err := haMembershipServerCheck(tc.nodeIDCtx, c, tc.haBackend)
			if tc.wantError != nil {
				require.EqualError(t, err, tc.wantError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
