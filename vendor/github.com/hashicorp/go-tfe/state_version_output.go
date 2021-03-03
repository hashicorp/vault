package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ StateVersionOutputs = (*stateVersionOutputs)(nil)

//State version outputs are the output values from a Terraform state file.
//They include the name and value of the output, as well as a sensitive boolean
//if the value should be hidden by default in UIs.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/state-version-outputs.html
type StateVersionOutputs interface {
	Read(ctx context.Context, outputID string) (*StateVersionOutput, error)
}

type stateVersionOutputs struct {
	client *Client
}

type StateVersionOutput struct {
	ID        string `jsonapi:"primary,state-version-outputs"`
	Name      string `jsonapi:"attr,name"`
	Sensitive bool   `jsonapi:"attr,sensitive"`
	Type      string `jsonapi:"attr,type"`
	Value     string `jsonapi:"attr,value"`
}

func (s *stateVersionOutputs) Read(ctx context.Context, outputID string) (*StateVersionOutput, error) {
	if !validStringID(&outputID) {
		return nil, errors.New("invalid value for run ID")
	}

	u := fmt.Sprintf("state-version-outputs/%s", url.QueryEscape(outputID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	so := &StateVersionOutput{}
	err = s.client.do(ctx, req, so)
	if err != nil {
		return nil, err
	}

	return so, nil
}
