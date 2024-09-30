// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"slices"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/require"
)

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

func Test_VersionBetweenReq_Run(t *testing.T) {
	t.Parallel()

	for desc, test := range map[string]struct {
		req  *ListVersionsReq
		res  func() *ListVersionsRes
		fail bool
	}{
		"no lister": {
			req: &ListVersionsReq{
				UpperBound: "1.18.0",
				LowerBound: "1.15.0",
			},
			fail: true,
		},
		"all": {
			req: &ListVersionsReq{
				UpperBound:    "1.18.0",
				LowerBound:    "1.15.0",
				VersionLister: NewMockClient(testAPIVersions),
			},
			res: func() *ListVersionsRes { return &ListVersionsRes{Versions: testAllVersions} },
		},
		"skips": {
			req: &ListVersionsReq{
				UpperBound:    "1.18.0",
				LowerBound:    "1.15.0",
				Skip:          []string{"1.16.0", "1.17.1"},
				VersionLister: NewMockClient(testAPIVersions),
			},
			res: func() *ListVersionsRes {
				versions := slices.Clone(testAllVersions)
				return &ListVersionsRes{Versions: slices.DeleteFunc(versions, func(v string) bool {
					if v == "1.17.1" || v == "1.16.0" {
						return true
					}
					return false
				})}
			},
		},
		"correct range upper and lower bound": {
			req: &ListVersionsReq{
				UpperBound:    "1.17.5",
				LowerBound:    "1.16.4",
				VersionLister: NewMockClient(testAPIVersions),
			},
			res: func() *ListVersionsRes {
				versions := slices.Clone(testAllVersions)
				return &ListVersionsRes{Versions: slices.DeleteFunc(versions, func(ver string) bool {
					u, err := semver.NewVersion("1.17.5")
					require.NoError(t, err)
					l, err := semver.NewVersion("1.16.4")
					require.NoError(t, err)
					v, err := semver.NewVersion(ver)
					require.NoError(t, err)
					if v.GreaterThan(u) || v.LessThan(l) {
						return true
					}
					return false
				})}
			},
		},
		"correct range nminus 1": {
			req: &ListVersionsReq{
				UpperBound:    "1.18.0",
				NMinus:        1,
				VersionLister: NewMockClient(testAPIVersions),
			},
			res: func() *ListVersionsRes {
				versions := slices.Clone(testAllVersions)
				return &ListVersionsRes{Versions: slices.DeleteFunc(versions, func(ver string) bool {
					u, err := semver.NewVersion("1.18.0")
					require.NoError(t, err)
					l, err := semver.NewVersion("1.17.0")
					require.NoError(t, err)
					v, err := semver.NewVersion(ver)
					require.NoError(t, err)
					if v.GreaterThan(u) || v.LessThan(l) {
						return true
					}
					return false
				})}
			},
		},
		"correct range nminus 2": {
			req: &ListVersionsReq{
				UpperBound:    "1.18.0",
				NMinus:        2,
				VersionLister: NewMockClient(testAPIVersions),
			},
			res: func() *ListVersionsRes { return &ListVersionsRes{Versions: testAllVersions} },
		},
		"lower and nminus": {
			req: &ListVersionsReq{
				UpperBound:    "1.18.0",
				LowerBound:    "1.15.0",
				NMinus:        2,
				VersionLister: NewMockClient(testAPIVersions),
			},
			res:  func() *ListVersionsRes { return &ListVersionsRes{Versions: testAllVersions} },
			fail: true,
		},
		"invalid upper": {
			req: &ListVersionsReq{
				UpperBound:    "1_18_0",
				LowerBound:    "1.15.0",
				VersionLister: NewMockClient(testAPIVersions),
			},
			res:  func() *ListVersionsRes { return &ListVersionsRes{Versions: testAllVersions} },
			fail: true,
		},
		"invalid lower": {
			req: &ListVersionsReq{
				UpperBound:    "1.18.0",
				LowerBound:    "1_15_0",
				VersionLister: NewMockClient(testAPIVersions),
			},
			res:  func() *ListVersionsRes { return &ListVersionsRes{Versions: testAllVersions} },
			fail: true,
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
			require.EqualValues(t, test.res().Versions, res.Versions)
		})
	}
}
