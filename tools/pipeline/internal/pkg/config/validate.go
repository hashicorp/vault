// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	slogctx "github.com/veqryn/slog-context"
)

// ValidateReq holds the state and configuration for validating a pipeline
// config
type ValidateReq struct {
	// DecodeRes is the result of decoding the pipeline configuration.
	DecodeRes *DecodeRes
}

// ValidateRes represents the response from validating a pipeline config
type ValidateRes struct {
	Valid  bool    `json:"valid,omitempty"`
	Config *Config `json:"config,omitempty"`
}

// Run executes the config validation operation
func (v *ValidateReq) Run(ctx context.Context) (*ValidateRes, error) {
	slog.Default().DebugContext(ctx, "validating pipeline configuration")

	if err := v.validate(); err != nil {
		return nil, err
	}

	res := &ValidateRes{}

	// Use the decoded config
	if v.DecodeRes.Config != nil {
		slog.Default().DebugContext(
			slogctx.Append(ctx,
				slog.String("path", v.DecodeRes.Path),
			),
			"validating configuration",
		)
		res.Config = v.DecodeRes.Config
		res.Valid = true

		return res, nil
	}

	return nil, errors.New("no configuration available to validate")
}

// validate checks that the request is properly configured
func (v *ValidateReq) validate() error {
	if v == nil {
		return errors.New("uninitialized")
	}

	if v.DecodeRes == nil {
		return errors.New("no decode result provided")
	}

	if err := v.DecodeRes.Err; err != nil {
		return fmt.Errorf("config decode error: %w", err)
	}

	return nil
}

// ToJSON marshals the response to JSON
func (r *ValidateRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling validate response to JSON: %w", err)
	}

	return b, nil
}

// String returns a string representation of the response
func (r *ValidateRes) String() string {
	if r == nil || r.Config == nil {
		return "No configuration has been provided"
	}

	if r.Valid {
		return "Configuration is valid"
	}

	return "Configuration is invalid"
}
