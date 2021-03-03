package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AgentPools = (*agentPools)(nil)

// AgentPools describes all the agent pool related methods that the Terraform
// Cloud API supports. Note that agents are not available in Terraform Enterprise.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/agents.html
type AgentPools interface {
	// List all the agent pools of the given organization.
	List(ctx context.Context, organization string, options AgentPoolListOptions) (*AgentPoolList, error)

	// Create a new agent pool with the given options.
	Create(ctx context.Context, organization string, options AgentPoolCreateOptions) (*AgentPool, error)

	// Read a agent pool by its ID.
	Read(ctx context.Context, agentPoolID string) (*AgentPool, error)

	// Update an agent pool by its ID.
	Update(ctx context.Context, agentPool string, options AgentPoolUpdateOptions) (*AgentPool, error)

	// Delete an agent pool by its ID.
	Delete(ctx context.Context, agentPoolID string) error
}

// agentPools implements AgentPools.
type agentPools struct {
	client *Client
}

// AgentPoolList represents a list of agent pools.
type AgentPoolList struct {
	*Pagination
	Items []*AgentPool
}

// AgentPool represents a Terraform Cloud agent pool.
type AgentPool struct {
	ID   string `jsonapi:"primary,agent-pools"`
	Name string `jsonapi:"attr,name"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
}

// AgentPoolListOptions represents the options for listing agent pools.
type AgentPoolListOptions struct {
	ListOptions
}

// List all the agent pools of the given organization.
func (s *agentPools) List(ctx context.Context, organization string, options AgentPoolListOptions) (*AgentPoolList, error) {
	if !validStringID(&organization) {
		return nil, errors.New("invalid value for organization")
	}

	u := fmt.Sprintf("organizations/%s/agent-pools", url.QueryEscape(organization))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	poolList := &AgentPoolList{}
	err = s.client.do(ctx, req, poolList)
	if err != nil {
		return nil, err
	}

	return poolList, nil
}

// AgentPoolCreateOptions represents the options for creating an agent pool.
type AgentPoolCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,agent-pools"`

	// A name to identify the agent pool.
	Name *string `jsonapi:"attr,name"`
}

func (o AgentPoolCreateOptions) valid() error {
	if !validString(o.Name) {
		return errors.New("name is required")
	}
	if !validStringID(o.Name) {
		return errors.New("invalid value for name")
	}
	return nil
}

// Create a new agent pool with the given options.
func (s *agentPools) Create(ctx context.Context, organization string, options AgentPoolCreateOptions) (*AgentPool, error) {
	if !validStringID(&organization) {
		return nil, errors.New("invalid value for organization")
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("organizations/%s/agent-pools", url.QueryEscape(organization))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	pool := &AgentPool{}
	err = s.client.do(ctx, req, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// Read a single agent pool by its ID.
func (s *agentPools) Read(ctx context.Context, agentpoolID string) (*AgentPool, error) {
	if !validStringID(&agentpoolID) {
		return nil, errors.New("invalid value for agent pool ID")
	}

	u := fmt.Sprintf("agent-pools/%s", url.QueryEscape(agentpoolID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	pool := &AgentPool{}
	err = s.client.do(ctx, req, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// AgentPoolUpdateOptions represents the options for updating an agent pool.
type AgentPoolUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,agent-pools"`

	// A new name to identify the agent pool.
	Name *string `jsonapi:"attr,name"`
}

func (o AgentPoolUpdateOptions) valid() error {
	if o.Name != nil && !validStringID(o.Name) {
		return errors.New("invalid value for name")
	}
	return nil
}

// Update an agent pool by its ID.
func (s *agentPools) Update(ctx context.Context, agentPoolID string, options AgentPoolUpdateOptions) (*AgentPool, error) {
	if !validStringID(&agentPoolID) {
		return nil, errors.New("invalid value for agent pool ID")
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("agent-pools/%s", url.QueryEscape(agentPoolID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	k := &AgentPool{}
	err = s.client.do(ctx, req, k)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Delete an agent pool by its ID.
func (s *agentPools) Delete(ctx context.Context, agentPoolID string) error {
	if !validStringID(&agentPoolID) {
		return errors.New("invalid value for agent pool ID")
	}

	u := fmt.Sprintf("agent-pools/%s", url.QueryEscape(agentPoolID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
