// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/stretchr/testify/require"
)

func TestCheckSealConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name:   "no-seals",
			config: Config{SharedConfig: &configutil.SharedConfig{Seals: []*configutil.KMS{}}},
		},
		{
			name: "one-seal",
			config: Config{SharedConfig: &configutil.SharedConfig{Seals: []*configutil.KMS{
				{
					Disabled: false,
				},
			}}},
		},
		{
			name: "one-disabled-seal",
			config: Config{SharedConfig: &configutil.SharedConfig{Seals: []*configutil.KMS{
				{
					Disabled: true,
				},
			}}},
		},
		{
			name: "two-seals-one-disabled",
			config: Config{SharedConfig: &configutil.SharedConfig{Seals: []*configutil.KMS{
				{
					Disabled: false,
				},
				{
					Disabled: true,
				},
			}}},
		},
		{
			name: "two-seals-enabled",
			config: Config{SharedConfig: &configutil.SharedConfig{Seals: []*configutil.KMS{
				{
					Disabled: false,
				},
				{
					Disabled: false,
				},
			}}},
			expectError: true,
		},
		{
			name: "two-disabled-seals",
			config: Config{SharedConfig: &configutil.SharedConfig{Seals: []*configutil.KMS{
				{
					Disabled: true,
				},
				{
					Disabled: true,
				},
			}}},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.checkSealConfig()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestRequestLimiterConfig verifies that the census config is correctly instantiated from HCL
func TestRequestLimiterConfig(t *testing.T) {
	testCases := []struct {
		name              string
		inConfig          string
		outErr            bool
		outRequestLimiter *configutil.RequestLimiter
	}{
		{
			name:              "empty",
			outRequestLimiter: nil,
		},
		{
			name: "disabled",
			inConfig: `
request_limiter {
	disable = true
}`,
			outRequestLimiter: &configutil.RequestLimiter{Disable: true},
		},
		{
			name: "invalid disable",
			inConfig: `
request_limiter {
	disable = "whywouldyoudothis"
}`,
			outErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := fmt.Sprintf(`
ui = false
storage "file" {
	path = "/tmp/test"
}

listener "tcp" {
	address = "0.0.0.0:8200"
}
%s`, tc.inConfig)
			gotConfig, err := ParseConfig(config, "")
			if tc.outErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.outRequestLimiter, gotConfig.RequestLimiter)
			}
		})
	}
}
