// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// Login authenticates against the current namespace and returns a new session object
// acting as that user.
func (s *Session) Login(path string, data map[string]any) *Session {
	s.t.Helper()

	newClient := s.newClient()
	secret, err := newClient.Logical().Write(path, data)
	require.NoError(s.t, err)

	if secret == nil || secret.Auth == nil {
		s.t.Fatal("failed to login")
	}

	newClient.SetToken(secret.Auth.ClientToken)

	return &Session{
		t:         s.t,
		Client:    newClient,
		Namespace: s.Namespace,
	}
}

func (s *Session) NewClientFromToken(token string) *Session {
	s.t.Helper()

	newClient := s.newClient()
	newClient.SetToken(token)

	return &Session{
		t:         s.t,
		Client:    newClient,
		Namespace: s.Namespace,
	}
}

func (s *Session) LoginUserpass(username, password string) *Session {
	s.t.Helper()

	path := fmt.Sprintf("auth/userpass/login/%s", username)
	payload := map[string]any{
		"password": password,
	}

	return s.Login(path, payload)
}

// TryLoginUserpass attempts to login with userpass but returns error instead of failing test
// This is useful for testing in environments where auth may not be available (e.g., managed HCP)
func (s *Session) TryLoginUserpass(username, password string) (*Session, error) {
	s.t.Helper()

	path := fmt.Sprintf("auth/userpass/login/%s", username)
	payload := map[string]any{
		"password": password,
	}

	secret, err := s.Client.Logical().Write(path, payload)
	if err != nil {
		return nil, err
	}

	clientToken, ok := secret.Auth.ClientToken, secret.Auth != nil
	if !ok {
		return nil, fmt.Errorf("login response missing client token")
	}

	newClient, err := s.Client.Clone()
	if err != nil {
		return nil, err
	}

	newClient.SetToken(clientToken)
	return &Session{
		t:         s.t,
		Client:    newClient,
		Namespace: s.Namespace,
	}, nil
}

func (s *Session) AssertWriteFails(path string, data map[string]any) {
	s.t.Helper()

	_, err := s.Client.Logical().Write(path, data)
	require.NotNil(s.t, err)
}

func (s *Session) AssertReadFails(path string) {
	s.t.Helper()

	_, err := s.Client.Logical().Read(path)
	require.NotNil(s.t, err)
}

func (s *Session) newClient() *api.Client {
	s.t.Helper()

	parentConfig := s.Client.CloneConfig()
	newClient, err := api.NewClient(parentConfig)
	require.NoError(s.t, err)
	newClient.SetNamespace(s.Namespace)

	return newClient
}
