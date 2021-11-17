package api

import (
	"context"
	"fmt"
)

// Auth is used to perform credential backend related operations.
type Auth struct {
	c *Client
}

type AuthMethod interface {
	Login(ctx context.Context, client *Client) (*Secret, error)
}

// Auth is used to return the client for credential-backend API calls.
func (c *Client) Auth() *Auth {
	return &Auth{c: c}
}

// Login sets up the required request body for login requests to the given auth
// method's /login API endpoint, and then performs a write to it. After a
// successful login, this method will automatically set the client's token to
// the login response's ClientToken as well.
//
// The Secret returned is the authentication secret, which if desired can be
// passed as input to the NewLifetimeWatcher method in order to start
// automatically renewing the token.
func (a *Auth) Login(ctx context.Context, authMethod AuthMethod) (*Secret, error) {
	if authMethod == nil {
		return nil, fmt.Errorf("no auth method provided for login")
	}

	authSecret, err := authMethod.Login(ctx, a.c)
	if err != nil {
		return nil, fmt.Errorf("unable to log in to auth method: %w", err)
	}
	if authSecret == nil || authSecret.Auth == nil || authSecret.Auth.ClientToken == "" {
		return nil, fmt.Errorf("login response from auth method did not return client token")
	}

	a.c.SetToken(authSecret.Auth.ClientToken)

	return authSecret, nil
}
