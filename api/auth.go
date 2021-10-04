package api

import "fmt"

// Auth is used to perform credential backend related operations.
type Auth struct {
	c *Client
}

type AuthMethod interface {
	Login(client *Client) (*Secret, error)
}

// Auth is used to return the client for credential-backend API calls.
func (c *Client) Auth() *Auth {
	return &Auth{c: c}
}

// Logs in to the given auth method and sets the client token.
func (a *Auth) Login(authMethod AuthMethod) (*Secret, error) {
	if authMethod == nil {
		return nil, fmt.Errorf("no auth method provided for login")
	}

	authSecret, err := authMethod.Login(a.c)
	if err != nil {
		return nil, fmt.Errorf("unable to log in to auth method: %w", err)
	}

	a.c.SetToken(authSecret.Auth.ClientToken)

	return authSecret, nil
}
