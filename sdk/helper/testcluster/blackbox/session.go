// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"slices"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// Session holds the test context and Vault client
type Session struct {
	t         *testing.T
	NoCleanup bool
	Client    *api.Client
	Namespace string
}

func (s *Session) T() *testing.T {
	return s.t
}

type SessionOpts func(s *Session)

func WithNoCleanup() SessionOpts {
	return func(s *Session) {
		s.NoCleanup = true
	}
}

func New(t *testing.T, opts ...SessionOpts) *Session {
	t.Helper()

	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

	// detect the parent namespace, e.g. "admin" in HVD
	parentNS := os.Getenv("VAULT_NAMESPACE")

	if addr == "" || token == "" {
		t.Fatal("VAULT_ADDR and VAULT_TOKEN are required")
	}

	config := api.DefaultConfig()
	config.Address = addr
	config.Timeout = 120 * time.Second // Increase timeout for LDAP operations that verify service accounts

	privClient, err := api.NewClient(config)
	require.NoError(t, err)
	privClient.SetToken(token)

	nsName := fmt.Sprintf("bbsdk-%s", randomString(8))
	nsURLPath := fmt.Sprintf("sys/namespaces/%s", nsName)

	_, err = privClient.Logical().Write(nsURLPath, nil)
	require.NoError(t, err)

	// session client should get the full namespace of parent + test
	fullNSPath := nsName
	if parentNS != "" {
		fullNSPath = path.Join(parentNS, nsName)
	}

	sessionConfig := privClient.CloneConfig()
	sessionClient, err := api.NewClient(sessionConfig)
	require.NoError(t, err)
	sessionClient.SetToken(token)
	sessionClient.SetNamespace(fullNSPath)

	session := &Session{
		t:         t,
		Client:    sessionClient,
		Namespace: nsName,
	}

	for opt := range slices.Values(opts) {
		opt(session)
	}

	t.Cleanup(func() {
		if session.NoCleanup {
			t.Logf("WARN: NoDebug has been set, not cleaning up namespace")
			return
		}
		_, err = privClient.Logical().Delete(nsURLPath)
		require.NoError(t, err)
		t.Logf("Cleaned up namespace %s", nsName)
	})

	// make sure the namespace has been created
	session.Eventually(func() error {
		// this runs inside the new namespace, so if it succeeds, we're good
		_, err := sessionClient.Auth().Token().LookupSelf()
		return err
	})

	return session
}

func randomString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
