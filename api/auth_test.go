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

func TestAuth_MFALoginSinglePhase(t *testing.T) {
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
}

func TestAuth_MFALoginTwoPhase(t *testing.T) {
	tests := []struct {
		name    string
		a       *Auth
		m       *mockAuthMethod
		creds   *string
		wantErr bool
	}{
		{
			name: "return MFARequirements",
			a: &Auth{
				c: &Client{},
			},
			m: &mockAuthMethod{
				mockedSecret: &Secret{
					Auth: &SecretAuth{
						MFARequirement: &logical.MFARequirement{
							MFARequestID:   "a-req-id",
							MFAConstraints: nil,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error if no MFARequirements",
			a: &Auth{
				c: &Client{},
			},
			m: &mockAuthMethod{
				mockedSecret: &Secret{
					Auth: &SecretAuth{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := tt.a.MFALogin(context.Background(), tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("MFALogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if secret.Auth.MFARequirement != tt.m.mockedSecret.Auth.MFARequirement {
				t.Errorf("MFALogin() returned %v, expected %v", secret.Auth.MFARequirement, tt.m.mockedSecret.Auth.MFARequirement)
				return
			}
		})
	}
}
