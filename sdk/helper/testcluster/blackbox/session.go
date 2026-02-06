// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// Session holds the test context and Vault client
type Session struct {
	t         *testing.T
	Client    *api.Client
	Namespace string
}

func New(t *testing.T) *Session {
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

	privClient, err := api.NewClient(config)
	require.NoError(t, err)
	privClient.SetToken(token)

	nsName := fmt.Sprintf("bbsdk-%s", randomString(8))
	nsURLPath := fmt.Sprintf("sys/namespaces/%s", nsName)

	_, err = privClient.Logical().Write(nsURLPath, nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = privClient.Logical().Delete(nsURLPath)
		require.NoError(t, err)
		t.Logf("Cleaned up namespace %s", nsName)
	})

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
