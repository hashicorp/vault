// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package server

import (
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
