// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/require"
)

//	TestListVersionsReq_calculateMinorCadenceFloor tests calculating the floor
//
// when given a version and nminus with a minor cadence.
func TestListVersionsReq_calculateMinorCadenceFloor(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		upperBound string
		nminus     uint
		floor      string
		shouldFail bool
	}{
		"basic minor cadence": {
			upperBound: "1.20.5",
			nminus:     3,
			floor:      "1.17.0",
		},
		"minor cadence with patch version": {
			upperBound: "1.18.3",
			nminus:     2,
			floor:      "1.16.0",
		},
		"minor cadence nminus 1": {
			upperBound: "1.15.0",
			nminus:     1,
			floor:      "1.14.0",
		},
		"minor cadence cannot negatively traverse past zero": {
			upperBound: "1.2.0",
			nminus:     5,
			shouldFail: true,
		},
		"nminus equals minor version": {
			upperBound: "1.3.0",
			nminus:     3,
			floor:      "1.0.0",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &ListVersionsReq{
				UpperBound: test.upperBound,
				NMinus:     test.nminus,
				Cadence:    CadenceMinor,
			}

			ceil, err := semver.NewVersion(test.upperBound)
			require.NoError(t, err)

			floor, gotCeil, err := req.calculateMinorCadenceFloor(ceil)

			if test.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.floor, floor.String())
			require.Equal(t, ceil, gotCeil)
		})
	}
}

// TestListVersionsReq_calculateMajorCadenceFloor tests calculating the floor
// when given a version and nminus with a major cadence.
func TestListVersionsReq_calculateMajorCadenceFloor(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		upperBound string
		nminus     uint
		floor      string
		shouldFail bool
	}{
		"basic major cadence": {
			upperBound: "5.1.3",
			nminus:     3,
			floor:      "2.0.0",
		},
		"major cadence with minor and patch": {
			upperBound: "10.5.2",
			nminus:     2,
			floor:      "8.0.0",
		},
		"major cadence nminus 1": {
			upperBound: "7.0.0",
			nminus:     1,
			floor:      "6.0.0",
		},
		"major cadence cannot negatively traverse past zero": {
			upperBound: "2.0.0",
			nminus:     5,
			shouldFail: true,
		},
		"nminus equals major version": {
			upperBound: "3.0.0",
			nminus:     3,
			floor:      "0.0.0",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &ListVersionsReq{
				UpperBound: test.upperBound,
				NMinus:     test.nminus,
				Cadence:    CadenceMajor,
			}

			ceil, err := semver.NewVersion(test.upperBound)
			require.NoError(t, err)

			floor, gotCeil, err := req.calculateMajorCadenceFloor(ceil)

			if test.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.floor, floor.String())
			require.Equal(t, ceil, gotCeil)
		})
	}
}

// TestListVersionsReq_calculateMajorCadenceFloor tests calculating the floor
// when given a version and nminus with a cadence transition.
func TestListVersionsReq_calculateFloorWithTransition(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		upperBound        string
		transitionVersion string
		nminus            uint
		cadence           VersionCadence
		priorCadence      VersionCadence
		floor             string
		shouldFail        bool
	}{
		"minor to major transition - basic": {
			upperBound:        "3.0.0",
			transitionVersion: "1.22.0",
			nminus:            3,
			cadence:           CadenceMajor,
			priorCadence:      CadenceMinor,
			floor:             "1.21.0",
		},
		"minor to major transition - all in major": {
			upperBound:        "10.0.0",
			transitionVersion: "1.22.0",
			nminus:            3,
			cadence:           CadenceMajor,
			priorCadence:      CadenceMinor,
			floor:             "7.0.0",
		},
		"minor to major transition - split evenly": {
			upperBound:        "3.0.0",
			transitionVersion: "1.21.0",
			nminus:            4,
			cadence:           CadenceMajor,
			priorCadence:      CadenceMinor,
			floor:             "1.19.0",
		},
		"impossible major to minor transition": {
			upperBound:        "1.20.0",
			transitionVersion: "5.0.0",
			nminus:            3,
			cadence:           CadenceMinor,
			priorCadence:      CadenceMajor,
			shouldFail:        true,
		},
		"major to minor transition - all in major": {
			upperBound:        "5.2.0",
			transitionVersion: "5.0.0",
			nminus:            3,
			cadence:           CadenceMinor,
			priorCadence:      CadenceMajor,
			floor:             "4.0.0",
		},
		"minor to major transition - upper bound in prior cadence": {
			upperBound:        "1.20.0",
			transitionVersion: "3.0.0",
			nminus:            3,
			cadence:           CadenceMajor,
			priorCadence:      CadenceMinor,
			floor:             "1.17.0",
		},
		"invalid transition - cannot transition to the same cadence": {
			upperBound:        "3.0.0",
			transitionVersion: "1.22.0",
			nminus:            3,
			cadence:           CadenceMinor,
			priorCadence:      CadenceMinor,
			shouldFail:        true,
		},
		"impossible nminus in transition cadence": {
			upperBound:        "3.0.0",
			transitionVersion: "1.1.0",
			nminus:            5,
			cadence:           CadenceMajor,
			priorCadence:      CadenceMinor,
			shouldFail:        true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &ListVersionsReq{
				UpperBound:        test.upperBound,
				TransitionVersion: test.transitionVersion,
				NMinus:            test.nminus,
				Cadence:           test.cadence,
				PriorCadence:      test.priorCadence,
			}

			ceil, err := semver.NewVersion(test.upperBound)
			require.NoError(t, err)

			floor, gotCeil, err := req.calculateFloorWithTransition(ceil)

			if test.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.floor, floor.String())
			require.Equal(t, ceil, gotCeil)
		})
	}
}

// TestListVersionsReq_calculateFloor tests the outer floor calculation func.
func TestListVersionsReq_calculateFloor(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		upperBound string
		nminus     uint
		cadence    VersionCadence
		floor      string
		shouldFail bool
	}{
		"default to minor cadence when not specified": {
			upperBound: "1.20.0",
			nminus:     3,
			cadence:    "",
			floor:      "1.17.0",
		},
		"explicit minor cadence": {
			upperBound: "1.20.0",
			nminus:     3,
			cadence:    CadenceMinor,
			floor:      "1.17.0",
		},
		"explicit major cadence": {
			upperBound: "5.0.0",
			nminus:     3,
			cadence:    CadenceMajor,
			floor:      "2.0.0",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &ListVersionsReq{
				UpperBound: test.upperBound,
				NMinus:     test.nminus,
				Cadence:    test.cadence,
			}

			ceil, err := semver.NewVersion(test.upperBound)
			require.NoError(t, err)

			floor, gotCeil, err := req.calculateFloor(ceil)

			if test.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.floor, floor.String())
			require.Equal(t, ceil, gotCeil)
		})
	}
}

// TestListVersionsReq_Validate_Cadence tests the request validation.
func TestListVersionsReq_Validate_Cadence(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		req        *ListVersionsReq
		shouldFail bool
	}{
		"valid minor cadence": {
			req: &ListVersionsReq{
				UpperBound:    "1.20.0",
				NMinus:        3,
				LicenseClass:  "oss",
				Cadence:       CadenceMinor,
				VersionLister: NewMockClient([]string{}),
			},
		},
		"valid major cadence": {
			req: &ListVersionsReq{
				UpperBound:    "5.0.0",
				NMinus:        3,
				LicenseClass:  "oss",
				Cadence:       CadenceMajor,
				VersionLister: NewMockClient([]string{}),
			},
		},
		"invalid cadence value": {
			req: &ListVersionsReq{
				UpperBound:    "1.20.0",
				NMinus:        3,
				LicenseClass:  "oss",
				Cadence:       "invalid",
				VersionLister: NewMockClient([]string{}),
			},
			shouldFail: true,
		},
		"valid transition with both cadences": {
			req: &ListVersionsReq{
				UpperBound:        "3.0.0",
				NMinus:            3,
				LicenseClass:      "oss",
				Cadence:           CadenceMajor,
				TransitionVersion: "1.22.0",
				PriorCadence:      CadenceMinor,
				VersionLister:     NewMockClient([]string{}),
			},
		},
		"transition version without transition cadence": {
			req: &ListVersionsReq{
				UpperBound:        "3.0.0",
				NMinus:            3,
				LicenseClass:      "oss",
				Cadence:           CadenceMajor,
				TransitionVersion: "1.22.0",
				VersionLister:     NewMockClient([]string{}),
			},
			shouldFail: true,
		},
		"transition cadence without transition version": {
			req: &ListVersionsReq{
				UpperBound:    "3.0.0",
				NMinus:        3,
				LicenseClass:  "oss",
				Cadence:       CadenceMajor,
				PriorCadence:  CadenceMinor,
				VersionLister: NewMockClient([]string{}),
			},
			shouldFail: true,
		},
		"same cadence for current and transition": {
			req: &ListVersionsReq{
				UpperBound:        "3.0.0",
				NMinus:            3,
				LicenseClass:      "oss",
				Cadence:           CadenceMinor,
				TransitionVersion: "1.22.0",
				PriorCadence:      CadenceMinor,
				VersionLister:     NewMockClient([]string{}),
			},
			shouldFail: true,
		},
		"invalid transition version format": {
			req: &ListVersionsReq{
				UpperBound:        "3.0.0",
				NMinus:            3,
				LicenseClass:      "oss",
				Cadence:           CadenceMajor,
				TransitionVersion: "invalid",
				PriorCadence:      CadenceMinor,
				VersionLister:     NewMockClient([]string{}),
			},
			shouldFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			test.req.ProductName = "vault"
			err := test.req.Validate(t.Context())

			if test.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
