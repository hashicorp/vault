// Copyright (c) HashiCorp, Inc.
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

// ListVersionsReq is a request to list versions from the releases API.
type ListVersionsReq struct {
	UpperBound   string
	LowerBound   string
	NMinus       uint
	Skip         []string
	LicenseClass string
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
		// This is quite naive. We only consider minor versions and pay no attention to whether or not
		// going backwards should bump us back to a different major/minor version. We also do not
		// consider preleases here at all so RC's will still go back two minor versions. In the event
		// that we bump major versions we'll have to revisit this.
		floor, err = semver.NewVersion(req.UpperBound)
		if err != nil {
			return nil, nil, fmt.Errorf("releases list versions req: invalid upper bound version: %w", err)
		}

		minor := floor.Minor() - int64(req.NMinus)
		if minor < 0 {
			return nil, nil, fmt.Errorf("releases list versions req: impossible nminus version, cannot subtract %d from %d", req.NMinus, floor.Minor())
		}

		// Create a new version with the new minor. Always set the patch to zero to allow for the full
		// range.
		nv := fmt.Sprintf("%d.%d.0", floor.Major(), minor)
		floor, err = semver.NewVersion(nv)
		if err != nil {
			return nil, nil, fmt.Errorf("releases list versions req: invalid nminus version: %s", nv)
		}
	} else {
		return nil, nil, fmt.Errorf("releases list versions req: no floor version or nminus has been specified")
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
		ctx, "vault", LicenseClass(req.LicenseClass), ceil, floor,
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
