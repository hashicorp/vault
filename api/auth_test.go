package api

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

type mockAuthMethod struct {
	mockedSecret *Secret
	mockedError  error
}

func (m *mockAuthMethod) Login(_ context.Context, _ *Client) (*Secret, error) {
	return m.mockedSecret, m.mockedError
}

func TestAuth_Login(t *testing.T) {
	a := &Auth{
		c: &Client{},
	}

	m := mockAuthMethod{
		mockedSecret: &Secret{
			Auth: &SecretAuth{
				ClientToken: "a-client-token",
			},
		},
		mockedError: nil,
	}

	t.Run("Login should set token on success", func(t *testing.T) {
		if a.c.Token() != "" {
			t.Errorf("client token was %v expected to be unset", a.c.Token())
		}

		_, err := a.Login(context.Background(), &m)
		if err != nil {
			t.Errorf("Login() error = %v", err)
			return
		}

		if a.c.Token() != m.mockedSecret.Auth.ClientToken {
			t.Errorf("client token was %v expected %v", a.c.Token(), m.mockedSecret.Auth.ClientToken)
			return
		}
	})
}

func TestAuth_MFALogin(t *testing.T) {
	t.Parallel()

	t.Run("MFALogin() should succeed if credentials are passed in", func(t *testing.T) {
		a := &Auth{
			c: &Client{},
		}
		m := mockAuthMethod{
			mockedSecret: &Secret{
				Auth: &SecretAuth{
					ClientToken: "a-client-token",
				},
			},
			mockedError: nil,
		}

		_, err := a.MFALogin(context.Background(), &m, "testMethod:testPasscode")
		if err != nil {
			t.Errorf("MFALogin() error %v", err)
			return
		}
		if a.c.Token() != m.mockedSecret.Auth.ClientToken {
			t.Errorf("client token was %v expected %v", a.c.Token(), m.mockedSecret.Auth.ClientToken)
			return
		}
	})

	t.Run("MFALogin() should return requirements if no creds are provided", func(t *testing.T) {
		a := &Auth{
			c: &Client{},
		}
		m := mockAuthMethod{
			mockedSecret: &Secret{
				Auth: &SecretAuth{
					MFARequirement: &logical.MFARequirement{
						MFARequestID:   "a-req-id",
						MFAConstraints: nil,
					},
				},
			},
			mockedError: nil,
		}

		secret, err := a.MFALogin(context.Background(), &m)
		if err != nil {
			t.Errorf("MFALogin() returned an error: %v", err)
			return
		}
		if secret.Auth.MFARequirement != m.mockedSecret.Auth.MFARequirement {
			t.Errorf("MFALogin() returned %v, expected %v", secret.Auth.MFARequirement, m.mockedSecret.Auth.MFARequirement)
			return
		}
	})

	t.Run("MFALogin() should error if no creds provided and no requirements returned", func(t *testing.T) {
		a := &Auth{
			c: &Client{},
		}
		m := mockAuthMethod{
			mockedSecret: &Secret{
				Auth: &SecretAuth{},
			},
			mockedError: nil,
		}
		if _, err := a.MFALogin(context.Background(), &m); err == nil {
			t.Errorf("MFALogin() should error if no credentials are set and no MFARequirements are returned")
			return
		}
	})
}
