// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfc

import (
	"errors"

	"github.com/hashicorp/go-tfe"
)

type client struct {
	*tfe.Client
}

type terraformToken struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Token       string `json:"token"`
}

func newClient(config *tfConfig) (*client, error) {
	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	cfg := &tfe.Config{
		Address:  config.Address,
		BasePath: config.BasePath,
		Token:    config.Token,
	}

	tfc, err := tfe.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &client{
		tfc,
	}, nil
}
