// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testVersionConfig = `
schema = 1

active_versions {
	version "1.19.x" {
		ce_active = true
		lts       = true
	}
	version "1.18.x" {
		ce_active = true
	}
	version "1.17.x" {
		ce_active = false
	}
	version "1.16.x" {
		ce_active = false
		lts       = true
	}
}
`

func TestListActiveVersionReq_unmarshalConfig(t *testing.T) {
	t.Parallel()

	versionsConfig, err := DecodeBytes([]byte(testVersionConfig))
	require.NoError(t, err)
	require.EqualValues(t, &VersionsConfig{
		Schema: 1,
		ActiveVersion: &ActiveVersion{
			Versions: map[string]*Version{
				"1.19.x": {CEActive: true, LTS: true},
				"1.18.x": {CEActive: true, LTS: false},
				"1.17.x": {CEActive: false, LTS: false},
				"1.16.x": {CEActive: false, LTS: true},
			},
		},
	}, versionsConfig)
}
