// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"

	"github.com/hashicorp/vault/vault/seal"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/stretchr/testify/require"
)

func init() {
	if signed := os.Getenv("VAULT_LICENSE_CI"); signed != "" {
		os.Setenv(EnvVaultLicense, signed)
	}
}

func TestMultiSealCases(t *testing.T) {
	cases := []struct {
		name                     string
		existingSealGenInfo      *seal.SealGenerationInfo
		allSealKmsConfigs        []*configutil.KMS
		expectedSealGenInfo      *seal.SealGenerationInfo
		isRewrapped              bool
		hasPartiallyWrappedPaths bool
		sealHaBetaEnabled        bool
		isErrorExpected          bool
		expectedErrorMsg         string
	}{
		// none_to_shamir
		{
			name:                "none_to_shamir",
			existingSealGenInfo: nil,
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "shamir",
					Name:     "shamirSeal1",
					Priority: 1,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "shamir",
						Name:     "shamirSeal1",
						Priority: 1,
					},
				},
			},
			sealHaBetaEnabled: true,
		},
		// none_to_auto
		{
			name:                "none_to_auto",
			existingSealGenInfo: nil,
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			sealHaBetaEnabled: true,
		},
		// none_to_multi
		{
			name:                "none_to_multi",
			existingSealGenInfo: nil,
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 2,
				},
			},
			isErrorExpected:   true,
			expectedErrorMsg:  "Initializing a cluster or enabling multi-seal on an existing cluster must occur with a single configured seal",
			sealHaBetaEnabled: true,
		},
		// none_to_multi_with_disabled_seals_with_beta
		{
			name:                "none_to_multi_with_disabled_seals_with_beta",
			existingSealGenInfo: nil,
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 2,
					Disabled: true,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
						Disabled: true,
					},
				},
			},
			isErrorExpected:   false,
			sealHaBetaEnabled: true,
		},
		// none_to_multi_with_disabled_seals_no_beta
		{
			name:                "none_to_multi_with_disabled_seals_no_beta",
			existingSealGenInfo: nil,
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 2,
					Disabled: true,
				},
			},
			isErrorExpected:   false,
			sealHaBetaEnabled: false,
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
						Disabled: true,
					},
				},
			},
		},
		// shamir_to_auto
		{
			name: "shamir_to_auto",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "shamir",
						Name:     "shamirSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 3,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			isRewrapped:       false,
			sealHaBetaEnabled: true,
		},
		// shamir_to_multi
		{
			name: "shamir_to_multi",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "shamir",
						Name:     "shamirSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 2,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 3,
				},
			},
			isRewrapped:       false,
			sealHaBetaEnabled: true,
			isErrorExpected:   true,
			expectedErrorMsg:  "cannot add more than one seal",
		},
		// auto_to_shamir_no_common_seal
		{
			name: "auto_to_shamir_no_common_seal",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "shamir",
					Name:     "shamirSeal1",
					Priority: 1,
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
			isErrorExpected:          true,
			expectedErrorMsg:         "must have at least one seal in common with the old generation",
		},
		// auto_to_shamir_with_common_seal
		{
			name: "auto_to_shamir_with_common_seal",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "shamir",
					Name:     "shamirSeal1",
					Priority: 1,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
					Disabled: true,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "shamir",
						Name:     "shamirSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
						Disabled: true,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
		},
		// auto_to_auto_no_common_seal
		{
			name: "auto_to_auto_no_common_seal",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 1,
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
			isErrorExpected:          true,
			expectedErrorMsg:         "must have at least one seal in common with the old generation",
		},
		// auto_to_auto_with_common_seal
		{
			name: "auto_to_auto_with_common_seal",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
					Disabled: true,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 2,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
						Disabled: true,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
		},
		// auto_to_multi_add_one
		{
			name: "auto_to_multi_add_one",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 2,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
		},
		// auto_to_multi_add_two
		{
			name: "auto_to_multi_add_two",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 2,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal3",
					Priority: 3,
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
			isErrorExpected:          true,
			expectedErrorMsg:         "cannot add more than one seal",
		},
		// multi_to_auto_delete_one
		{
			name: "multi_to_auto_delete_one",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
		},
		// multi_to_auto_delete_two
		{
			name: "multi_to_auto_delete_two",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal3",
						Priority: 3,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
			isErrorExpected:          true,
			expectedErrorMsg:         "cannot delete more than one seal",
		},
		// disable_two_auto
		{
			name: "disable_two_auto",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal3",
						Priority: 3,
					},
				},
			},
			allSealKmsConfigs: []*configutil.KMS{
				{
					Type:     "pkcs11",
					Name:     "autoSeal1",
					Priority: 1,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal2",
					Priority: 2,
					Disabled: true,
				},
				{
					Type:     "pkcs11",
					Name:     "autoSeal3",
					Priority: 3,
					Disabled: true,
				},
			},
			expectedSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
						Disabled: true,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal3",
						Priority: 3,
						Disabled: true,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			sealHaBetaEnabled:        true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &ServerCommand{}
			cmd.logger = corehelpers.NewTestLogger(t)
			if tc.existingSealGenInfo != nil {
				tc.existingSealGenInfo.SetRewrapped(tc.isRewrapped)
			}
			sealGenInfo, err := cmd.computeSealGenerationInfo(tc.existingSealGenInfo, tc.allSealKmsConfigs, tc.hasPartiallyWrappedPaths, tc.sealHaBetaEnabled)
			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErrorMsg)
				require.Nil(t, sealGenInfo)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.expectedSealGenInfo, sealGenInfo)
			}
		})
	}

	cases2 := []struct {
		name                     string
		existingSealGenInfo      *seal.SealGenerationInfo
		newSealGenInfo           *seal.SealGenerationInfo
		isRewrapped              bool
		hasPartiallyWrappedPaths bool
		isErrorExpected          bool
		expectedErrorMsg         string
	}{
		// same_generation_different_seals
		{
			name: "same_generation_different_seals",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal3",
						Priority: 2,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          true,
			expectedErrorMsg:         "existing seal generation is the same, but the configured seals are different",
		},

		// same_generation_same_seals
		{
			name: "same_generation_same_seals",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          false,
		},
		// existing seal gen info rewrapped is set to false
		{
			name: "existing_sgi_rewrapped_false",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			isRewrapped:              false,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          true,
			expectedErrorMsg:         "cannot make seal config changes while seal re-wrap is in progress, please revert any seal configuration changes",
		},
		// single seal migration use-case
		{
			name: "single_seal_migration",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "transit",
						Name:     "transit",
						Priority: 1,
					},
				},
			},
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "transit",
						Name:     "transit-disabled",
						Priority: 1,
						Disabled: true,
					},
					{
						Type:     "shamir",
						Name:     "shamir",
						Priority: 1,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          false,
		},
		// migrate from non-beta single seal to single seal
		{
			name:                "none_to_single_seal",
			existingSealGenInfo: nil,
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "shamir",
						Name:     "shamir",
						Priority: 1,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          false,
		},
		// migrate from non-beta single seal to multi seal, with one disabled, so perform an old style migration
		{
			name:                "none_to_multiple_seals_one_disabled",
			existingSealGenInfo: nil,
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type: "pkcs11",
						Name: "autoSeal",
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal",
						Disabled: true,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          false,
		},
		// migrate from non-beta single seal to multi seal
		{
			name:                "none_to_multiple_seals",
			existingSealGenInfo: nil,
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          true,
			expectedErrorMsg:         "Initializing a cluster or enabling multi-seal on an existing cluster must occur with a single configured seal",
		},
		// have partially wrapped paths
		{
			name: "have_partially_wrapped_paths",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: true,
			isErrorExpected:          true,
			expectedErrorMsg:         "cannot make seal config changes while seal re-wrap is in progress, please revert any seal configuration changes",
		},
		// no partially wrapped paths
		{
			name: "no_partially_wrapped_paths",
			existingSealGenInfo: &seal.SealGenerationInfo{
				Generation: 2,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
					{
						Type:     "pkcs11",
						Name:     "autoSeal2",
						Priority: 2,
					},
				},
			},
			newSealGenInfo: &seal.SealGenerationInfo{
				Generation: 1,
				Seals: []*configutil.KMS{
					{
						Type:     "pkcs11",
						Name:     "autoSeal1",
						Priority: 1,
					},
				},
			},
			isRewrapped:              true,
			hasPartiallyWrappedPaths: false,
			isErrorExpected:          false,
		},
	}
	for _, tc := range cases2 {
		t.Run(tc.name, func(t *testing.T) {
			if tc.existingSealGenInfo != nil {
				tc.existingSealGenInfo.SetRewrapped(tc.isRewrapped)
			}
			err := tc.newSealGenInfo.Validate(tc.existingSealGenInfo, tc.hasPartiallyWrappedPaths)
			switch {
			case tc.isErrorExpected:
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErrorMsg)
			default:
				require.NoError(t, err)
			}
		})
	}
}
