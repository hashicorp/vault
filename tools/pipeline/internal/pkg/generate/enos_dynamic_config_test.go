// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package generate

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
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

var testAllVersions = []string{
	"1.16.6",
	"1.16.7",
	"1.16.8",
	"1.16.9",
	"1.16.10",
	"1.17.2",
	"1.17.3",
	"1.17.4",
	"1.17.5",
	"1.17.6",
	"1.18.0-rc1",
}

func Test_EnosDynamicConfigReq_Validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		in   *EnosDynamicConfigReq
		fail bool
	}{
		"ce edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition:  "ce",
				VaultVersion:  "1.18.0",
				EnosDir:       t.TempDir(),
				FileName:      "test.hcl",
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
		},
		"oss edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition:  "oss",
				VaultVersion:  "1.18.0",
				EnosDir:       t.TempDir(),
				FileName:      "test.hcl",
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
		},
		"ent edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition:  "ent",
				VaultVersion:  "1.18.0",
				EnosDir:       t.TempDir(),
				FileName:      "test.hcl",
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
		},
		"enterprise edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition:  "enterprise",
				VaultVersion:  "1.18.0",
				EnosDir:       t.TempDir(),
				FileName:      "test.hcl",
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
		},
		"ent.hsm edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition:  "ent.hsm",
				VaultVersion:  "1.18.0",
				EnosDir:       t.TempDir(),
				FileName:      "test.hcl",
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
		},
		"ent.fips1402 edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition:  "ent.fips1402",
				VaultVersion:  "1.18.0",
				EnosDir:       t.TempDir(),
				FileName:      "test.hcl",
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
		},
		"ent.hsm.fips1402 edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition:  "ent.hsm.fips1402",
				VaultVersion:  "1.18.0",
				EnosDir:       t.TempDir(),
				FileName:      "test.hcl",
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
		},
		"unknown edition": {
			in: &EnosDynamicConfigReq{
				VaultEdition: "ent.nope",
				VaultVersion: "1.18.0",
				EnosDir:      t.TempDir(),
				FileName:     "test.hcl",
			},
			fail: true,
		},
		"invalid version": {
			in: &EnosDynamicConfigReq{
				VaultEdition: "ent.hsm.fips1402",
				VaultVersion: "vault-1.18.0",
				EnosDir:      t.TempDir(),
				FileName:     "test.hcl",
			},
			fail: true,
		},
		"target dir doesn't exist": {
			in: &EnosDynamicConfigReq{
				VaultEdition: "ent.hsm.fips1402",
				VaultVersion: "1.18.0",
			},
			fail: true,
		},
		"no file name": {
			in: &EnosDynamicConfigReq{
				VaultEdition: "ent.hsm.fips1402",
				VaultVersion: "1.18.0",
				EnosDir:      t.TempDir(),
			},
			fail: true,
		},
		"no version lister": {
			in: &EnosDynamicConfigReq{
				VaultEdition: "ent.hsm.fips1402",
				VaultVersion: "1.18.0",
				EnosDir:      t.TempDir(),
				FileName:     "test.hcl",
			},
			fail: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.in.Validate(context.Background())
			if test.fail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_EnosDynamicConfigReq_Run(t *testing.T) {
	t.Parallel()

	for desc, test := range map[string]struct {
		req  *EnosDynamicConfigReq
		res  func() *EnosDynamicConfigRes
		hcl  []byte
		fail bool
	}{
		"default config": {
			req: &EnosDynamicConfigReq{
				FileName:      "test.hcl",
				VaultEdition:  "ent.hsm.fips1402",
				VaultVersion:  "1.18.0",
				Skip:          []string{"1.17.2", "1.17.5"},
				NMinus:        2,
				EnosDir:       t.TempDir(),
				VersionLister: releases.NewMockClient(testAPIVersions),
			},
			res: func() *EnosDynamicConfigRes {
				versions := testAllVersions
				versions = slices.DeleteFunc(versions, func(v string) bool {
					return v == "1.17.2" || v == "1.17.5"
				})
				return &EnosDynamicConfigRes{
					Globals: &Globals{
						SampleAttributes: &SampleAttrs{
							AWSRegion:             []string{"us-east-1", "us-west-2"},
							DistroVersionAmzn:     []string{"2023"},
							DistroVersionLeap:     []string{"15.6"},
							DistroVersionRhel:     []string{"8.10, 9.4"},
							DistroVersionSles:     []string{"15.6"},
							DistroVersionUbuntu:   []string{"20.04", "24.04"},
							UpgradeInitialVersion: versions,
						},
					},
				}
			},
			hcl: []byte(`
globals {
  sample_attributes = {
    aws_region              = ["us-east-1", "us-west-2"]
    distro_version_amzn     = ["2023"]
    distro_version_leap     = ["15.6"]
    distro_version_rhel     = ["8.10, 9.4"]
    distro_version_sles     = ["15.6"]
    distro_version_ubuntu   = ["20.04", "24.04"]
    upgrade_initial_version = ["1.16.6", "1.16.7", "1.16.8", "1.16.9", "1.16.10", "1.17.3", "1.17.4", "1.17.6", "1.18.0-rc1"]
  }
}
`),
		},
	} {
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			res, err := test.req.Run(context.Background())
			if test.fail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, test.res(), res)
			b, err := os.ReadFile(filepath.Join(test.req.EnosDir, test.req.FileName))
			require.NoError(t, err)
			require.EqualValuesf(t, test.hcl, b, string(b))
		})
	}
}
