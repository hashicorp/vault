// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import (
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestFileBackend_newFileBackend_fallback ensures that we get the correct errors
// in CE when we try to enable a FileBackend with enterprise options like fallback
// and filter.
func TestFileBackend_newFileBackend_fallback(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		backendConfig        *BackendConfig
		isErrorExpected      bool
		expectedErrorMessage string
	}{
		"non-fallback-device-with-filter": {
			backendConfig: &BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback":  "false",
					"file_path": discard,
					"filter":    "mount_type == kv",
				},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "enterprise-only options supplied: invalid configuration",
		},
		"fallback-device-with-filter": {
			backendConfig: &BackendConfig{
				MountPath:  "discard",
				SaltConfig: &salt.Config{},
				SaltView:   &logical.InmemStorage{},
				Logger:     hclog.NewNullLogger(),
				Config: map[string]string{
					"fallback":  "true",
					"file_path": discard,
					"filter":    "mount_type == kv",
				},
			},
			isErrorExpected:      true,
			expectedErrorMessage: "enterprise-only options supplied: invalid configuration",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			be, err := newFileBackend(tc.backendConfig, &NoopHeaderFormatter{})

			if tc.isErrorExpected {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrorMessage)
			} else {
				require.NoError(t, err)
				require.NotNil(t, be)
			}
		})
	}
}

// TestFileBackend_newFileBackend_FilterFormatterSink ensures that when configuring
// a backend in community edition we cannot configure a filter node.
// We can verify that we have formatter and sink nodes added to the backend.
// The order of calls influences the slice of IDs on the Backend.
func TestFileBackend_newFileBackend_FilterFormatterSink(t *testing.T) {
	t.Parallel()

	cfg := map[string]string{
		"file_path": "/tmp/foo",
		"mode":      "0777",
		"format":    "json",
		"filter":    "mount_type == \"kv\"",
	}

	backendConfig := &BackendConfig{
		SaltView:   &logical.InmemStorage{},
		SaltConfig: &salt.Config{},
		Config:     cfg,
		MountPath:  "bar",
		Logger:     hclog.NewNullLogger(),
	}

	b, err := newFileBackend(backendConfig, &NoopHeaderFormatter{})
	require.Error(t, err)
	require.EqualError(t, err, "enterprise-only options supplied: invalid configuration")

	// Try without filter option
	delete(cfg, "filter")
	b, err = newFileBackend(backendConfig, &NoopHeaderFormatter{})
	require.NoError(t, err)

	require.Len(t, b.nodeIDList, 2)
	require.Len(t, b.nodeMap, 2)

	id := b.nodeIDList[0]
	node := b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeFormatter, node.Type())

	id = b.nodeIDList[1]
	node = b.nodeMap[id]
	require.Equal(t, eventlogger.NodeTypeSink, node.Type())
}

// TestBackend_IsFallback ensures that no CE audit device can be a fallback.
func TestBackend_IsFallback(t *testing.T) {
	t.Parallel()

	cfg := &BackendConfig{
		MountPath:  "discard",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
		Config: map[string]string{
			"fallback":  "true",
			"file_path": discard,
		},
	}

	be, err := newFileBackend(cfg, &NoopHeaderFormatter{})
	require.Error(t, err)
	require.EqualError(t, err, "enterprise-only options supplied: invalid configuration")

	// Remove the option and try again
	delete(cfg.Config, "fallback")

	be, err = newFileBackend(cfg, &NoopHeaderFormatter{})
	require.NoError(t, err)
	require.NotNil(t, be)
	require.Equal(t, false, be.IsFallback())
}
