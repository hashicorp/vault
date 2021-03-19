package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ UserTokens = (*userTokens)(nil)

// UserTokens describes all the user token related methods that the
// Terraform Cloud/Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/enterprise/api/user-tokens.html
type UserTokens interface {
	// List all the tokens of the given user ID.
	List(ctx context.Context, userID string) (*UserTokenList, error)

	// Generate a new user token
	Generate(ctx context.Context, userID string, options UserTokenGenerateOptions) (*UserToken, error)

	// Read a user token by its ID.
	Read(ctx context.Context, tokenID string) (*UserToken, error)

	// Delete a user token by its ID.
	Delete(ctx context.Context, tokenID string) error
}

// userTokens implements UserTokens.
type userTokens struct {
	client *Client
}

// UserTokenList is a list of tokens for the given user ID.
type UserTokenList struct {
	*Pagination
	Items []*UserToken
}

// UserToken represents a Terraform Enterprise user token.
type UserToken struct {
	ID          string    `jsonapi:"primary,authentication-tokens"`
	CreatedAt   time.Time `jsonapi:"attr,created-at,iso8601"`
	Description string    `jsonapi:"attr,description"`
	LastUsedAt  time.Time `jsonapi:"attr,last-used-at,iso8601"`
	Token       string    `jsonapi:"attr,token"`
}

// UserTokenGenerateOptions the options for creating a user token.
type UserTokenGenerateOptions struct {
	// Description of the token
	Description string `jsonapi:"attr,description,omitempty"`
}

// Generate a new user token
func (s *userTokens) Generate(ctx context.Context, userID string, options UserTokenGenerateOptions) (*UserToken, error) {
	if !validStringID(&userID) {
		return nil, errors.New("invalid value for user ID")
	}

	u := fmt.Sprintf("users/%s/authentication-tokens", url.QueryEscape(userID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	ut := &UserToken{}
	err = s.client.do(ctx, req, ut)
	if err != nil {
		return nil, err
	}

	return ut, err
}

// List shows existing user tokens
func (s *userTokens) List(ctx context.Context, userID string) (*UserTokenList, error) {
	if !validStringID(&userID) {
		return nil, errors.New("invalid value for user ID")
	}

	u := fmt.Sprintf("users/%s/authentication-tokens", url.QueryEscape(userID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tl := &UserTokenList{}
	err = s.client.do(ctx, req, tl)
	if err != nil {
		return nil, err
	}

	return tl, err
}

// Read a user token by its ID.
func (s *userTokens) Read(ctx context.Context, tokenID string) (*UserToken, error) {
	if !validStringID(&tokenID) {
		return nil, errors.New("invalid value for token ID")
	}

	u := fmt.Sprintf("authentication-tokens/%s", url.QueryEscape(tokenID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tt := &UserToken{}
	err = s.client.do(ctx, req, tt)
	if err != nil {
		return nil, err
	}

	return tt, err
}

// Delete a user token by its ID.
func (s *userTokens) Delete(ctx context.Context, tokenID string) error {
	if !validStringID(&tokenID) {
		return errors.New("invalid value for token ID")
	}

	u := fmt.Sprintf("authentication-tokens/%s", url.QueryEscape(tokenID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
