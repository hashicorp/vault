// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/metadata"
	slogctx "github.com/veqryn/slog-context"
)

// VersionCadence defines the type of version increment for n-minus calculations.
type VersionCadence string

const (
	// CadenceMinor increments by minor version (e.g., 1.20 -> 1.17 for n-minus 3)
	CadenceMinor VersionCadence = "minor"

	// CadenceMajor increments by major version (e.g., 5.0.0 -> 2.0.0 for n-minus 3)
	CadenceMajor VersionCadence = "major"
)

// ListVersionsReq is a request to list versions from the releases API.
type ListVersionsReq struct {
	UpperBound string         // The upper bound of our range
	LowerBound string         // The lower bound of our range
	NMinus     uint           // Another way of specifying our lower bound. Calculate our lower as upper-N
	Cadence    VersionCadence // The release candence to use when calculating the lower bound from n-minus
	// If we transitioned from a different cadence, this is the _last_ version in that cadence
	TransitionVersion string
	// If we've transitioned from one release candence to another, eg: minor releases to major, which version is the transition
	PriorCadence VersionCadence

	Skip         []string // Specifically exclude some versions from the list
	LicenseClass string   // Which license class to look for, either

	ProductName string // Which product to search for

	VersionLister
}

// ListVersionsRes is a list versions response.
type ListVersionsRes struct {
	Versions []string `json:"versions,omitempty"`
}

// NewListVersionsReq returns a new releases API version lister request.
func NewListVersionsReq() *ListVersionsReq {
	return &ListVersionsReq{}
}

func (req *ListVersionsReq) Validate(ctx context.Context) error {
	if req == nil {
		return errors.New("releases list versions req: unitialized")
	}

	if req.ProductName == "" {
		return errors.New("no product name provided")
	}

	// Allow callers to pass in "oss" or "ce" but always rewrite it to what the releases API expects.
	if slices.Contains(metadata.CeEditions, req.LicenseClass) {
		req.LicenseClass = string(LicenseClassCE)
	}

	// Allow callers to pass in any enterprise edition but always rewrite it to what the releases API
	// expects.
	if slices.Contains(metadata.EntEditions, req.LicenseClass) {
		req.LicenseClass = string(LicenseClassEnt)
	}

	if req.LicenseClass != "oss" && req.LicenseClass != "enterprise" {
		return fmt.Errorf("releases list versions req: validate: invalid license class: %s: must be 'ce' or 'enterprise'", req.LicenseClass)
	}

	if req.VersionLister == nil {
		return errors.New("releases list versions req: no version lister has been configured")
	}

	if req.LowerBound != "" && req.NMinus != 0 {
		return errors.New("releases list versions req: only one of a lower bound floor or nminus option can be configured")
	}

	// Validate cadence configuration
	if req.Cadence != "" && req.Cadence != CadenceMinor && req.Cadence != CadenceMajor {
		return fmt.Errorf("releases list versions req: invalid cadence: %s: must be 'minor' or 'major'", req.Cadence)
	}

	// If transition version is specified, transition cadence must also be specified
	if req.TransitionVersion != "" {
		if req.PriorCadence == "" {
			return errors.New("releases list versions req: transition cadence must be specified when transition version is set")
		}
		if req.PriorCadence != CadenceMinor && req.PriorCadence != CadenceMajor {
			return fmt.Errorf("releases list versions req: invalid transition cadence: %s: must be 'minor' or 'major'", req.PriorCadence)
		}
		if req.PriorCadence == req.Cadence {
			return errors.New("releases list versions req: transition cadence must be different from current cadence")
		}
		// Validate transition version is a valid semver
		_, err := semver.NewVersion(req.TransitionVersion)
		if err != nil {
			return fmt.Errorf("releases list versions req: invalid transition version: %w", err)
		}
	}

	// If transition cadence is specified, transition version must also be specified
	if req.PriorCadence != "" && req.TransitionVersion == "" {
		return errors.New("releases list versions req: transition version must be specified when transition cadence is set")
	}

	return nil
}

func (req *ListVersionsReq) VersionRange() (*semver.Version, *semver.Version, error) {
	ceil, err := semver.NewVersion(req.UpperBound)
	if err != nil {
		return nil, nil, fmt.Errorf("releases list versions req: invalid upper bound version: %w", err)
	}

	var floor *semver.Version

	if req.LowerBound != "" {
		floor, err = semver.NewVersion(req.LowerBound)
		if err != nil {
			return nil, nil, fmt.Errorf("releases list versions req: invalid lower bound version: %w", err)
		}
	} else if req.NMinus != 0 {
		// Calculate floor based on cadence
		if req.TransitionVersion != "" {
			return req.calculateFloorWithTransition(ceil)
		}
		return req.calculateFloor(ceil)
	} else {
		return nil, nil, fmt.Errorf("releases list versions req: no floor version or nminus has been specified")
	}

	return floor, ceil, nil
}

// calculateFloor calculates the floor version based on the current cadence
func (req *ListVersionsReq) calculateFloor(ceil *semver.Version) (*semver.Version, *semver.Version, error) {
	// Default to minor cadence for backward compatibility
	cadence := req.Cadence
	if cadence == "" {
		cadence = CadenceMinor
	}

	switch cadence {
	case CadenceMajor:
		return req.calculateMajorCadenceFloor(ceil)
	case CadenceMinor:
		return req.calculateMinorCadenceFloor(ceil)
	default:
		return nil, nil, fmt.Errorf("releases list versions req: unknown cadence: %s", cadence)
	}
}

// calculateMinorCadenceFloor calculates the floor version using minor version increments
func (req *ListVersionsReq) calculateMinorCadenceFloor(ceil *semver.Version) (*semver.Version, *semver.Version, error) {
	minor := ceil.Minor() - int64(req.NMinus)
	if minor < 0 {
		return nil, nil, fmt.Errorf("releases list versions req: impossible nminus: cannot subtract %d from minor version %d", req.NMinus, ceil.Minor())
	}

	floorStr := fmt.Sprintf("%d.%d.0", ceil.Major(), minor)
	floor, err := semver.NewVersion(floorStr)
	if err != nil {
		return nil, nil, fmt.Errorf("releases list versions req: invalid calculated floor: %s: %w", floorStr, err)
	}

	return floor, ceil, nil
}

// calculateMajorCadenceFloor calculates the floor version using major version increments
func (req *ListVersionsReq) calculateMajorCadenceFloor(ceil *semver.Version) (*semver.Version, *semver.Version, error) {
	major := ceil.Major() - int64(req.NMinus)
	if major < 0 {
		return nil, nil, fmt.Errorf("releases list versions req: impossible nminus: cannot subtract %d from major version %d", req.NMinus, ceil.Major())
	}

	floorStr := fmt.Sprintf("%d.0.0", major)
	floor, err := semver.NewVersion(floorStr)
	if err != nil {
		return nil, nil, fmt.Errorf("releases list versions req: invalid calculated floor: %s: %w", floorStr, err)
	}

	return floor, ceil, nil
}

// decrement takes in a version and a cadence and returns a new version that has been decremented
// based on the cadence.
func (req *ListVersionsReq) decrement(ver *semver.Version, cadence VersionCadence) (*semver.Version, error) {
	var vs string

	switch cadence {
	case CadenceMajor:
		vs = fmt.Sprintf("%d.0.0", ver.Major()-1)
	case CadenceMinor:
		vs = fmt.Sprintf("%d.%d.0", ver.Major(), ver.Minor()-1)
	default:
		return nil, fmt.Errorf("unsupported cadence: %s", cadence)
	}

	v, err := semver.NewVersion(vs)
	if err != nil {
		err = fmt.Errorf("decrementing version %s from candence %s to %s: %w", ver.String(), cadence, vs, err)
	}

	return v, err
}

// calculateFloorWithTransition calculates the floor version when transitioning between cadences
func (req *ListVersionsReq) calculateFloorWithTransition(ceil *semver.Version) (*semver.Version, *semver.Version, error) {
	transition, err := semver.NewVersion(req.TransitionVersion)
	if err != nil {
		return nil, nil, fmt.Errorf("releases list versions req: invalid transition version: %w", err)
	}

	// Default to minor cadence for backward compatibility
	if req.Cadence == "" {
		req.Cadence = CadenceMinor
	}

	floor, err := semver.NewVersion(ceil.String())
	if err != nil {
		return nil, nil, fmt.Errorf("releases list versions req: invalid ceiling version: %w", err)
	}

	// Interate backwards with the current cadence unless we're equal to our transition version
	for range req.NMinus {
		// If our floor is less than the transition or equal to it then we're using the prior candence
		if floor.LessThan(transition) || floor.Equal(transition) {
			floor, err = req.decrement(floor, req.PriorCadence)
			if err != nil {
				return nil, nil, fmt.Errorf("releases list versions req: calculating nminus floor with transition: %w", err)
			}
		} else {
			// Our floor is greater than the transition. Test if we decrement the current cadence. If we
			// Go below the transition then our current floor is the transition.
			floor, err = req.decrement(floor, req.Cadence)
			if err != nil {
				return nil, nil, fmt.Errorf("releases list versions req: calculating nminus floor with transition: %w", err)
			}

			if floor.LessThan(transition) {
				floor, err = semver.NewVersion(transition.String())
				if err != nil {
					return nil, nil, fmt.Errorf("releases list versions req: calculating nminus floor with transition: %w", err)
				}
			}
		}
	}

	return floor, ceil, nil
}

// Run the versions between request by determining our upper and lower version boundaries, using
// them to get a list of versions from the configured VersionLister, and then filtering out any
// skipped versions.
func (req *ListVersionsReq) Run(ctx context.Context) (*ListVersionsRes, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ctx = slogctx.Append(ctx,
		slog.String("upper-bound", req.UpperBound),
		slog.String("lower-bound", req.LowerBound),
		slog.Uint64("n-minus", uint64(req.NMinus)),
		slog.String("edition", string(req.LicenseClass)),
		slog.String("cadence", string(req.Cadence)),
		slog.String("transition-version", req.TransitionVersion),
		slog.String("prior-cadence", string(req.PriorCadence)),
		"skip", req.Skip,
	)
	slog.Default().DebugContext(ctx, "running releases list version request")

	err := req.Validate(ctx)
	if err != nil {
		return nil, err
	}

	slog.Default().DebugContext(ctx, "determining version request range")
	floor, ceil, err := req.VersionRange()
	if err != nil {
		return nil, err
	}

	versions, err := req.VersionLister.ListVersions(
		ctx, req.ProductName, LicenseClass(req.LicenseClass), ceil, floor,
	)
	if err != nil {
		return nil, err
	}

	res := &ListVersionsRes{Versions: []string{}}
	seen := map[string]struct{}{}

	for _, v := range versions {
		rv, err := semver.NewVersion(v)
		if err != nil {
			return nil, err
		}

		// The releases API will list all editions as seperate release versions, as it should. However,
		// we don't make that distinction here. For our purposes we neeed a singular list of all versions
		// on a license class basis. As such, we'll drop metadata and only focus on major, minor, patch,
		// and prerelease.
		nv, err := rv.SetMetadata("")
		if err != nil {
			return nil, fmt.Errorf("failed to unset metadata: %v", err)
		}

		// Since each enterprise version can be listed many times due to the metadata. If we've already
		// seen this version we can move on.
		if _, ok := seen[nv.String()]; ok {
			continue
		}
		seen[nv.String()] = struct{}{}

		if len(req.Skip) > 0 {
			if slices.Contains(req.Skip, nv.String()) || slices.Contains(req.Skip, rv.String()) {
				// We're skipping this version
				continue
			}
		}

		// Add it to the versions slice
		res.Versions = append(res.Versions, nv.String())
	}

	res.Versions, err = sortVersions(res.Versions)
	slog.Default().DebugContext(ctx, "found versions", "versions", res.Versions)

	return res, err
}
