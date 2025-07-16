// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/require"
)

var testAPIVersions = []string{
	"1.16.10+ent.hsm.fips1402",
	"1.16.10+ent.fips1402",
	"1.16.10+ent.hsm",
	"1.16.10+ent",
	"1.17.6+ent.hsm.fips1402",
	"1.17.6+ent.fips1402",
	"1.17.6+ent.hsm",
	"1.17.6+ent",
	"1.18.0-rc1+ent.hsm.fips1402",
	"1.18.0-rc1+ent.fips1402",
	"1.18.0-rc1+ent.hsm",
	"1.18.0-rc1+ent",
	"1.17.5+ent.hsm.fips1402",
	"1.17.5+ent.fips1402",
	"1.17.5+ent.hsm",
	"1.17.5+ent",
	"1.16.9+ent.hsm.fips1402",
	"1.16.9+ent.fips1402",
	"1.16.9+ent.hsm",
	"1.16.9+ent",
	"1.17.4+ent.hsm.fips1402",
	"1.17.4+ent.fips1402",
	"1.17.4+ent.hsm",
	"1.17.4+ent",
	"1.16.8+ent.hsm.fips1402",
	"1.16.8+ent.fips1402",
	"1.16.8+ent.hsm",
	"1.16.8+ent",
	"1.17.3+ent.hsm.fips1402",
	"1.17.3+ent.fips1402",
	"1.17.3+ent.hsm",
	"1.16.7+ent.hsm.fips1402",
	"1.17.3+ent",
	"1.16.7+ent.fips1402",
	"1.16.7+ent.hsm",
	"1.16.7+ent",
	"1.17.2+ent.hsm.fips1402",
	"1.17.2+ent.fips1402",
	"1.17.2+ent.hsm",
	"1.17.2+ent",
	"1.16.6+ent.hsm.fips1402",
	"1.16.6+ent.fips1402",
	"1.16.6+ent.hsm",
	"1.16.6+ent",
}

func Test_Client_ListVersions(t *testing.T) {
	t.Parallel()

	for desc, test := range map[string]struct {
		rf func() (*semver.Version, *semver.Version)
		e  []string
	}{
		"all": {
			rf: func() (*semver.Version, *semver.Version) {
				ceil, err := semver.NewVersion("1.18.0")
				require.NoError(t, err)
				floor, err := semver.NewVersion("1.15.0")
				require.NoError(t, err)
				return ceil, floor
			},
			e: testAPIVersions,
		},
		"high": {
			rf: func() (*semver.Version, *semver.Version) {
				ceil, err := semver.NewVersion("1.18.0")
				require.NoError(t, err)
				floor, err := semver.NewVersion("1.17.0")
				require.NoError(t, err)
				return ceil, floor
			},
			e: []string{
				"1.17.6+ent.hsm.fips1402",
				"1.17.6+ent.fips1402",
				"1.17.6+ent.hsm",
				"1.17.6+ent",
				"1.18.0-rc1+ent.hsm.fips1402",
				"1.18.0-rc1+ent.fips1402",
				"1.18.0-rc1+ent.hsm",
				"1.18.0-rc1+ent",
				"1.17.5+ent.hsm.fips1402",
				"1.17.5+ent.fips1402",
				"1.17.5+ent.hsm",
				"1.17.5+ent",
				"1.17.4+ent.hsm.fips1402",
				"1.17.4+ent.fips1402",
				"1.17.4+ent.hsm",
				"1.17.4+ent",
				"1.17.3+ent.hsm.fips1402",
				"1.17.3+ent.fips1402",
				"1.17.3+ent.hsm",
				"1.17.3+ent",
				"1.17.2+ent.hsm.fips1402",
				"1.17.2+ent.fips1402",
				"1.17.2+ent.hsm",
				"1.17.2+ent",
			},
		},
		"middle": {
			rf: func() (*semver.Version, *semver.Version) {
				ceil, err := semver.NewVersion("1.17.4")
				require.NoError(t, err)
				floor, err := semver.NewVersion("1.16.7")
				require.NoError(t, err)
				return ceil, floor
			},
			e: []string{
				"1.16.10+ent.hsm.fips1402",
				"1.16.10+ent.fips1402",
				"1.16.10+ent.hsm",
				"1.16.10+ent",
				"1.16.9+ent.hsm.fips1402",
				"1.16.9+ent.fips1402",
				"1.16.9+ent.hsm",
				"1.16.9+ent",
				"1.17.4+ent.hsm.fips1402",
				"1.17.4+ent.fips1402",
				"1.17.4+ent.hsm",
				"1.17.4+ent",
				"1.16.8+ent.hsm.fips1402",
				"1.16.8+ent.fips1402",
				"1.16.8+ent.hsm",
				"1.16.8+ent",
				"1.17.3+ent.hsm.fips1402",
				"1.17.3+ent.fips1402",
				"1.17.3+ent.hsm",
				"1.16.7+ent.hsm.fips1402",
				"1.17.3+ent",
				"1.16.7+ent.fips1402",
				"1.16.7+ent.hsm",
				"1.16.7+ent",
				"1.17.2+ent.hsm.fips1402",
				"1.17.2+ent.fips1402",
				"1.17.2+ent.hsm",
				"1.17.2+ent",
			},
		},
		"low": {
			rf: func() (*semver.Version, *semver.Version) {
				ceil, err := semver.NewVersion("1.16.9")
				require.NoError(t, err)
				floor, err := semver.NewVersion("1.15.3")
				require.NoError(t, err)
				return ceil, floor
			},
			e: []string{
				"1.16.9+ent.hsm.fips1402",
				"1.16.9+ent.fips1402",
				"1.16.9+ent.hsm",
				"1.16.9+ent",
				"1.16.8+ent.hsm.fips1402",
				"1.16.8+ent.fips1402",
				"1.16.8+ent.hsm",
				"1.16.8+ent",
				"1.16.7+ent.hsm.fips1402",
				"1.16.7+ent.fips1402",
				"1.16.7+ent.hsm",
				"1.16.7+ent",
				"1.16.6+ent.hsm.fips1402",
				"1.16.6+ent.fips1402",
				"1.16.6+ent.hsm",
				"1.16.6+ent",
			},
		},
	} {
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			client := NewMockClient(testAPIVersions)
			ceil, floor := test.rf()
			res, err := client.ListVersions(context.Background(), "vault", "enterprise", ceil, floor)
			require.NoError(t, err)
			require.EqualValues(t, test.e, res)
		})
	}
}
