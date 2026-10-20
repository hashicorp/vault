// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/internalshared/agentproxyshared/sink"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

// Avoiding a circular dependency in the test.
type mockSink struct {
	token *atomic.String
}

func (m *mockSink) Token() string {
	return m.token.Load()
}

func (m *mockSink) WriteToken(token string) error {
	m.token.Store(token)
	return nil
}

func newMockSink(t *testing.T) sink.Sink {
	t.Helper()

	return &mockSink{
		token: atomic.NewString(""),
	}
}

// testNewStaticSecretCacheUpdater returns a new StaticSecretCacheUpdater
// for use in tests.
func testNewStaticSecretCacheUpdater(t *testing.T, client *api.Client) *StaticSecretCacheUpdater {
	t.Helper()

	lc := testNewLeaseCache(t, []*SendResponse{})
	tokenSink := newMockSink(t)
	tokenSink.WriteToken(client.Token())

	updater, err := NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.updater"),
		TokenSink:  tokenSink,
	})
	require.NoError(t, err)
	return updater
}

// TestNewStaticSecretCacheUpdater tests the NewStaticSecretCacheUpdater method,
// to ensure it errors out when appropriate.
func TestNewStaticSecretCacheUpdater(t *testing.T) {
	t.Parallel()

	lc := testNewLeaseCache(t, []*SendResponse{})
	config := api.DefaultConfig()
	logger := logging.NewVaultLogger(hclog.Trace).Named("cache.updater")
	client, err := api.NewClient(config)
	require.NoError(t, err)
	tokenSink := newMockSink(t)

	// Expect an error if any of the arguments are nil:
	updater, err := NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     nil,
		LeaseCache: lc,
		Logger:     logger,
		TokenSink:  tokenSink,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: nil,
		Logger:     logger,
		TokenSink:  tokenSink,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     nil,
		TokenSink:  tokenSink,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.updater"),
		TokenSink:  nil,
	})
	require.Error(t, err)
	require.Nil(t, updater)

	// Don't expect an error if the arguments are as expected
	updater, err = NewStaticSecretCacheUpdater(&StaticSecretCacheUpdaterConfig{
		Client:     client,
		LeaseCache: lc,
		Logger:     logging.NewVaultLogger(hclog.Trace).Named("cache.updater"),
		TokenSink:  tokenSink,
	})
	require.NoError(t, err)
	require.NotNil(t, updater)
}

// TestCheckForDeleteOrDestroyEvent tests the behaviour of checkForDeleteOrDestroyEvent
// and assures it gives the right responses for different events.
func TestCheckForDeleteOrDestroyEvent(t *testing.T) {
	t.Parallel()

	expectedVersions := []int{1, 3, 5}
	jsonFormatExpectedVersions := "[1,3,5]"
	expectedPath := "secret-v2/data/my-secret"
	deletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":  "2",
					"deleted_versions": jsonFormatExpectedVersions,
					"modified":         true,
					"operation":        "delete",
					"path":             "secret-v2/delete/my-secret",
				},
			},
			"event_type": "kv-v2/delete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	undeletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"undeleted_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "undelete",
					"path":               "secret-v2/undelete/my-secret",
				},
			},
			"event_type": "kv-v2/undelete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	destroyedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"destroyed_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "destroy",
					"path":               "secret-v2/destroy/my-secret",
				},
			},
			"event_type": "kv-v2/destroy",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	actualVersions, actualPath := checkForDeleteOrDestroyEvent(deletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(undeletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(destroyedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)
}

// TestCheckForDeleteOrDestroyNamespacedEvent tests the behaviour of checkForDeleteOrDestroyEvent
// with namespaces in paths.
func TestCheckForDeleteOrDestroyNamespacedEvent(t *testing.T) {
	t.Parallel()

	expectedVersions := []int{1, 3, 5}
	jsonFormatExpectedVersions := "[1,3,5]"
	expectedPath := "ns/secret-v2/data/my-secret"
	deletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":  "2",
					"deleted_versions": jsonFormatExpectedVersions,
					"modified":         true,
					"operation":        "delete",
					"data_path":        "secret-v2/data/my-secret",
					"path":             "secret-v2/delete/my-secret",
				},
			},
			"namespace":  "ns/",
			"event_type": "kv-v2/delete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	undeletedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"undeleted_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "undelete",
					"data_path":          "secret-v2/data/my-secret",
					"path":               "secret-v2/undelete/my-secret",
				},
			},
			"namespace":  "ns/",
			"event_type": "kv-v2/undelete",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	destroyedVersionEventMap := map[string]interface{}{
		"id":     "abc",
		"source": "abc",
		"data": map[string]interface{}{
			"event": map[string]interface{}{
				"id": "bar",
				"metadata": map[string]interface{}{
					"current_version":    "2",
					"destroyed_versions": jsonFormatExpectedVersions,
					"modified":           true,
					"operation":          "destroy",
					"data_path":          "secret-v2/data/my-secret",
					"path":               "secret-v2/destroy/my-secret",
				},
			},
			"namespace":  "ns/",
			"event_type": "kv-v2/destroy",
			"plugin_info": map[string]interface{}{
				"mount_path": "secret-v2/",
				"plugin":     "kv",
				"version":    2,
			},
		},
	}

	actualVersions, actualPath := checkForDeleteOrDestroyEvent(deletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(undeletedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)

	actualVersions, actualPath = checkForDeleteOrDestroyEvent(destroyedVersionEventMap)
	require.Equal(t, expectedVersions, actualVersions)
	require.Equal(t, expectedPath, actualPath)
}
