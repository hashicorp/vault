// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build testonly

package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/hashicorp/vault/vault/activity/generation"
	"github.com/stretchr/testify/require"
)

// TestSystemBackend_handleActivityWriteData calls the activity log write endpoint and confirms that the inputs are
// correctly validated
func TestSystemBackend_handleActivityWriteData(t *testing.T) {
	testCases := []struct {
		name      string
		operation logical.Operation
		input     map[string]interface{}
		wantError error
	}{
		{
			name:      "read fails",
			operation: logical.ReadOperation,
			wantError: logical.ErrUnsupportedOperation,
		},
		{
			name:      "empty write fails",
			operation: logical.CreateOperation,
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "wrong key fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"other": "data"},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "incorrectly formatted data fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": "data"},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "incorrect json data fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"other":"json"}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "empty write value fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"write":[],"data":[]}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "empty data value fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"write":["WRITE_PRECOMPUTED_QUERIES"],"data":[]}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "correctly formatted data succeeds",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"write":["WRITE_PRECOMPUTED_QUERIES"],"data":[{"current_month":true,"all":{"clients":[{"count":5}]}}]}`},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := testSystemBackend(t)
			req := logical.TestRequest(t, tc.operation, "internal/counters/activity/write")
			req.Data = tc.input
			resp, err := b.HandleRequest(namespace.RootContext(nil), req)
			if tc.wantError != nil {
				require.Equal(t, tc.wantError, err, resp.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test_singleMonthActivityClients_addNewClients verifies that new clients are created correctly, adhering to the
// requested parameters. The clients should have the correct namespace and mount, replaced with the default if the input
// values are empty. The clients should have a generated ID if one is not supplied. The new client should be added to
// the month's `clients` slice and `allClients` map the correct number of times
func Test_singleMonthActivityClients_addNewClients(t *testing.T) {
	tests := []struct {
		name          string
		namespace     string
		mount         string
		clients       *generation.Client
		wantNamespace string
		wantMount     string
		wantID        string
	}{
		{
			name:          "default mount and namespace are used",
			namespace:     "default_ns",
			mount:         "default_mount",
			wantNamespace: "default_ns",
			wantMount:     "default_mount",
			clients:       &generation.Client{},
		},
		{
			name:          "record namespace is used, default mount is used",
			namespace:     "default_ns",
			mount:         "default_mount",
			wantNamespace: "ns",
			wantMount:     "default_mount",
			clients: &generation.Client{
				Namespace: "ns",
				Mount:     "mount",
			},
		},
		{
			name: "predefined ID is used",
			clients: &generation.Client{
				Id: "client_id",
			},
			wantID: "client_id",
		},
		{
			name: "non zero times seen",
			clients: &generation.Client{
				TimesSeen: 5,
			},
		},
		{
			name: "non zero count",
			clients: &generation.Client{
				Count: 5,
			},
		},
		{
			name: "non zero times seen and count",
			clients: &generation.Client{
				Count:     5,
				TimesSeen: 3,
			},
		},
		{
			name: "non entity client",
			clients: &generation.Client{
				NonEntity: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &singleMonthActivityClients{
				allClients: make(map[string]*activity.EntityRecord),
			}
			err := m.addNewClients(tt.clients, tt.namespace, tt.mount)
			require.NoError(t, err)
			numNew := tt.clients.Count
			if numNew == 0 {
				numNew = 1
			}
			numSeen := tt.clients.TimesSeen
			if numSeen == 0 {
				numSeen = 1
			}
			require.Len(t, m.allClients, int(numNew))
			require.Len(t, m.clients, int(numNew*numSeen))
			for _, c := range m.clients {
				rec := m.allClients[c]
				require.NotNil(t, rec)
				require.Equal(t, tt.wantNamespace, rec.NamespaceID)
				require.Equal(t, tt.wantMount, rec.MountAccessor)
				require.Equal(t, tt.clients.NonEntity, rec.NonEntity)
				if tt.wantID != "" {
					require.Equal(t, tt.wantID, rec.ClientID)
				} else {
					require.NotEqual(t, "", rec.ClientID)
				}
			}
		})
	}
}

// Test_multipleMonthsActivityClients_processMonth verifies that a month of data is added correctly. The test checks
// that default values are handled correctly for mounts and namespaces.
func Test_multipleMonthsActivityClients_processMonth(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	tests := []struct {
		name      string
		clients   *generation.Data
		wantError bool
		numMonths int
	}{
		{
			name: "specified namespace and mount exist",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: namespace.RootNamespaceID,
					Mount:     "identity/",
				}}}},
			},
			numMonths: 1,
		},
		{
			name: "specified namespace exists, mount empty",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: namespace.RootNamespaceID,
				}}}},
			},
			numMonths: 1,
		},
		{
			name: "empty namespace and mount",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{}}}},
			},
			numMonths: 1,
		},
		{
			name: "namespace doesn't exist",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: "abcd",
				}}}},
			},
			wantError: true,
			numMonths: 1,
		},
		{
			name: "namespace exists, mount doesn't exist",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: namespace.RootNamespaceID,
					Mount:     "mount",
				}}}},
			},
			wantError: true,
			numMonths: 1,
		},
		{
			name: "older month",
			clients: &generation.Data{
				Month:   &generation.Data_MonthsAgo{MonthsAgo: 4},
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{}}}},
			},
			numMonths: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMultipleMonthsActivityClients(tt.numMonths)
			err := m.processMonth(context.Background(), core, tt.clients)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, m.allClients, len(tt.clients.GetAll().Clients))
				require.Len(t, m.months[tt.clients.GetMonthsAgo()].clients, len(tt.clients.GetAll().Clients))
				for _, c := range m.allClients {
					require.NotEmpty(t, c.NamespaceID)
					require.NotEmpty(t, c.MountAccessor)
				}
			}
		})
	}
}
