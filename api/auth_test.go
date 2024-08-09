// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"testing"
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

	t.Run("Login with nil AuthMethod should return error", func(t *testing.T) {
		_, err := a.Login(context.Background(), nil)
		if err == nil {
			t.Errorf("expected error when AuthMethod is nil, got none")
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

	t.Run("MFALogin() with empty credentials should assume two-phase", func(t *testing.T) {
		a := &Auth{
			c: &Client{},
		}

		m := mockAuthMethod{
			mockedSecret: &Secret{
				Auth: &SecretAuth{
					MFARequirement: &MFARequirement{
						MFARequestID: "a-req-id",
					},
				},
			},
			mockedError: nil,
		}

		secret, err := a.MFALogin(context.Background(), &m)
		if err != nil {
			t.Errorf("MFALogin() error %v", err)
			return
		}

		if secret.Auth.MFARequirement.MFARequestID != "a-req-id" {
			t.Errorf("MFARequirement ID was %v expected %v", secret.Auth.MFARequirement.MFARequestID, "a-req-id")
		}
	})
}

func TestAuth_MFALoginTwoPhase(t *testing.T) {
	tests := []struct {
		name    string
		a       *Auth
		m       *mockAuthMethod
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
						MFARequirement: &MFARequirement{
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

			if tt.wantErr {
				if secret != nil {
					t.Errorf("MFALogin() returned non-nil secret when error was expected")
				}
				return
			}

			if secret.Auth.MFARequirement != tt.m.mockedSecret.Auth.MFARequirement {
				t.Errorf("MFALogin() returned %v, expected %v", secret.Auth.MFARequirement, tt.m.mockedSecret.Auth.MFARequirement)
				return
			}
		})
	}
}

func TestAuth_MFAValidate(t *testing.T) {
	a := &Auth{
		c: &Client{},
	}

	t.Run("MFAValidate should return error if secret is nil", func(t *testing.T) {
		_, err := a.MFAValidate(context.Background(), nil, nil)
		if err == nil {
			t.Errorf("expected error when secret is nil, got none")
		}
	})

	t.Run("MFAValidate should return error if MFARequirement is nil", func(t *testing.T) {
		secret := &Secret{
			Auth: &SecretAuth{},
		}
		_, err := a.MFAValidate(context.Background(), secret, nil)
		if err == nil {
			t.Errorf("expected error when MFARequirement is nil, got none")
		}
	})
}

func TestCheckAndSetToken(t *testing.T) {
	a := &Auth{
		c: &Client{},
	}

	t.Run("checkAndSetToken should return error if secret is nil", func(t *testing.T) {
		_, err := a.checkAndSetToken(nil)
		if err == nil {
			t.Errorf("expected error when secret is nil, got none")
		}
	})

	t.Run("checkAndSetToken should return error if ClientToken is missing", func(t *testing.T) {
		secret := &Secret{
			Auth: &SecretAuth{},
		}
		_, err := a.checkAndSetToken(secret)
		if err == nil {
			t.Errorf("expected error when ClientToken is missing, got none")
		}
	})
}
