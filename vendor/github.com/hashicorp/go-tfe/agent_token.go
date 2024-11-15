// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ AgentTokens = (*agentTokens)(nil)

// AgentTokens describes all the agent token related methods that the
// HCP Terraform API supports.
//
// TFE API docs:
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/agent-tokens
type AgentTokens interface {
	// List all the agent tokens of the given agent pool.
	List(ctx context.Context, agentPoolID string) (*AgentTokenList, error)

	// Create a new agent token with the given options.
	Create(ctx context.Context, agentPoolID string, options AgentTokenCreateOptions) (*AgentToken, error)

	// Read an agent token by its ID.
	Read(ctx context.Context, agentTokenID string) (*AgentToken, error)

	// Delete an agent token by its ID.
	Delete(ctx context.Context, agentTokenID string) error
}

// agentTokens implements AgentTokens.
type agentTokens struct {
	client *Client
}

// AgentToken represents a HCP Terraform agent token.
type AgentToken struct {
	ID          string    `jsonapi:"primary,authentication-tokens"`
	CreatedAt   time.Time `jsonapi:"attr,created-at,iso8601"`
	Description string    `jsonapi:"attr,description"`
	LastUsedAt  time.Time `jsonapi:"attr,last-used-at,iso8601"`
	Token       string    `jsonapi:"attr,token"`
}

// AgentTokenList represents a list of agent tokens.
type AgentTokenList struct {
	*Pagination
	Items []*AgentToken
}

// AgentTokenCreateOptions represents the options for creating an agent token.
type AgentTokenCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,agent-tokens"`

	// Description of the token
	Description *string `jsonapi:"attr,description"`
}

// List all the agent tokens of the given agent pool.
func (s *agentTokens) List(ctx context.Context, agentPoolID string) (*AgentTokenList, error) {
	if !validStringID(&agentPoolID) {
		return nil, ErrInvalidAgentPoolID
	}

	u := fmt.Sprintf("agent-pools/%s/authentication-tokens", url.PathEscape(agentPoolID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	tokenList := &AgentTokenList{}
	err = req.Do(ctx, tokenList)
	if err != nil {
		return nil, err
	}

	return tokenList, nil
}

// Create a new agent token with the given options.
func (s *agentTokens) Create(ctx context.Context, agentPoolID string, options AgentTokenCreateOptions) (*AgentToken, error) {
	if !validStringID(&agentPoolID) {
		return nil, ErrInvalidAgentPoolID
	}

	if !validString(options.Description) {
		return nil, ErrAgentTokenDescription
	}

	u := fmt.Sprintf("agent-pools/%s/authentication-tokens", url.PathEscape(agentPoolID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	at := &AgentToken{}
	err = req.Do(ctx, at)
	if err != nil {
		return nil, err
	}

	return at, err
}

// Read an agent token by its ID.
func (s *agentTokens) Read(ctx context.Context, agentTokenID string) (*AgentToken, error) {
	if !validStringID(&agentTokenID) {
		return nil, ErrInvalidAgentTokenID
	}

	u := fmt.Sprintf("authentication-tokens/%s", url.PathEscape(agentTokenID))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	at := &AgentToken{}
	err = req.Do(ctx, at)
	if err != nil {
		return nil, err
	}

	return at, err
}

// Delete an agent token by its ID.
func (s *agentTokens) Delete(ctx context.Context, agentTokenID string) error {
	if !validStringID(&agentTokenID) {
		return ErrInvalidAgentTokenID
	}

	u := fmt.Sprintf("authentication-tokens/%s", url.PathEscape(agentTokenID))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}
