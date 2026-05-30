// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func (s *Session) MustSanitizedConfig() *api.Secret {
	s.t.Helper()

	var config *api.Secret
	err := s.Req(
		func(c *api.Client) error {
			var err error
			config, err = c.Logical().Read("sys/config/state/sanitized")
			return err
		},
		WithClientRootNamespace(),
		WithClientTimeout(2*time.Second),
	)

	require.NoError(s.t, err)
	require.NotNil(s.t, config)

	return config
}

func (s *Session) MustGetConfigStorageType() string {
	s.t.Helper()

	sanitizedConfig := s.MustSanitizedConfig()
	// Verify we have at least one server configured
	storageType := s.AssertSecret(sanitizedConfig).
		Data().
		GetMap("storage").
		GetKey("type")

	storage, ok := storageType.(string)
	if !ok {
		s.t.Fatalf("cluster storage is unknown: %v", storageType)
	}

	if storage == "" {
		s.t.Fatal("cluster storage is empty")
	}

	s.t.Logf("cluster is using storage type: %s", storage)

	return storage
}

func (s *Session) getConfigStorageType() (string, error) {
	s.t.Helper()

	var storageTypeStr string
	getStorageType := func(c *api.Client) error {
		secret, err := c.Logical().Read("sys/seal-status")
		if err != nil {
			return err
		}

		if secret == nil || len(secret.Data) < 1 {
			return fmt.Errorf("seal-status is empty")
		}

		storage, ok := secret.Data["storage_type"]
		if !ok {
			return fmt.Errorf("seal-status does not include storage_type")
		}

		storageTypeStr, ok = storage.(string)
		if !ok {
			return fmt.Errorf("malformed storage type")
		}

		return nil
	}

	return storageTypeStr, s.Req(
		getStorageType,
		WithClientTimeout(2*time.Second),
	)
}
