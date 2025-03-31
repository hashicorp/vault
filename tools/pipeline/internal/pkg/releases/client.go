// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/releases-api/pkg/api"
	"github.com/hashicorp/releases-api/pkg/client"
	"github.com/hashicorp/releases-api/pkg/models"
	slogctx "github.com/veqryn/slog-context"
)

// Client is an api.releases.hashicorp.com API client.
type Client struct {
	Addr string
}

// VersionLister lists versions from the releases API.
type VersionLister interface {
	ListVersions(ctx context.Context, product string, edition LicenseClass, ceil, floor *semver.Version) ([]string, error)
}

var _ VersionLister = (*Client)(nil)

type LicenseClass string

// These map to the licenses classes defined by the releases API.
const (
	LicenseClassNone LicenseClass = ""
	LicenseClassCE   LicenseClass = "oss"
	LicenseClassEnt  LicenseClass = "enterprise"
	LicenseClassHCP  LicenseClass = "hcp"
)

// NewClient returns a new releases API client.
func NewClient() *Client {
	return &Client{
		Addr: "api.releases.hashicorp.com",
	}
}

// GetRelease takes a context, product, edition, and version and returns the relase information.
func (c *Client) GetRelease(ctx context.Context, product string, edition LicenseClass, version string) (*models.Release, error) {
	ctx = slogctx.Append(ctx,
		slog.String("product", product),
		slog.String("edition", string(edition)),
		slog.String("version", version),
	)
	slog.Default().DebugContext(ctx, "getting release")

	rc := client.New(c.Addr, "", api.V1MimeType)
	lc := string(edition)
	release, _, err := rc.GetRelease(&client.GetReleaseParams{
		Version:      version,
		Product:      product,
		LicenseClass: &lc,
	})

	slog.Default().DebugContext(ctx, "got release", "release", release, "error", err)

	return release, err
}

// ListVersions takes a context, product, edition, a ceiling version and floor version and returns
// all unique versions between the ceiling and floor range.
func (c *Client) ListVersions(ctx context.Context, product string, edition LicenseClass, ceil, floor *semver.Version) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ctx = slogctx.Append(ctx,
		slog.String("product", product),
		slog.String("edition", string(edition)),
		slog.String("ceil", ceil.String()),
		slog.String("floor", floor.String()),
		slog.String("addr", c.Addr),
	)
	slog.Default().DebugContext(ctx, "listing release versions")

	// The releases API lists releases by their upload order starting by most recent, not sequentially
	// by version. It also uses upload timestamps as a method of pagination. Since version upload
	// lineages are all mixed up we'll have to start by listing the latest versions and work back
	// until we've reached an upload timestamp at our floor version. To do so we'll get the created
	// time of our floor and subtract a day which ought to give us a reasonable time floor.
	// NOTE: This requires our floor version to actually exist.
	floorR, err := c.GetRelease(ctx, product, edition, floor.String())
	if err != nil {
		return nil, fmt.Errorf("invalid floor version: %s, %w", floor.String(), err)
	}
	timeFloor := floorR.TimestampCreated.Add(time.Duration(-24) * time.Hour)

	// We'll do the same for our ceiling but be less restrictive as to allow an infinite ceiling.
	// If the ceiling versions exists we'll use it as our initial time ceiling, otherwise we'll
	// set it to the current time.
	var after time.Time
	ceilR, err := c.GetRelease(ctx, product, edition, ceil.String())
	if err == nil {
		after = ceilR.TimestampCreated
	} else {
		after = time.Now()
	}

	versions := []string{}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		slog.Default().DebugContext(ctx, "listing releases", "after", after, "limit", 20)
		rc := client.New(c.Addr, "", api.V1MimeType)
		releases, _, err := rc.ListReleases(&client.ListReleasesParams{
			Product:        product,
			AfterTimestamp: after,
			Limit:          20, // 20 is the max limit
			LicenseClass:   string(edition),
		})
		if err != nil {
			return nil, err
		}

		// We've returned all that we can return
		if len(releases) < 1 {
			break
		}

		// Select any release verions that are within our desired range
		releaseVersions, err := selectReleaseVersions(releases, ceil, floor)
		if err != nil {
			return nil, err
		}
		versions = append(versions, releaseVersions...)

		// Reset our pagination and do another round
		after = releases[len(releases)-1].TimestampCreated

		// Short circuit if we've hit our time floor.
		if timeFloor.After(after) {
			break
		}
	}

	versions, err = sortVersions(versions)
	slog.Default().DebugContext(ctx, "found release versions", "versions", versions, "error", err)

	return versions, err
}

func sortVersions(in []string) ([]string, error) {
	if len(in) < 2 {
		return in, nil
	}

	c := semver.Collection{}
	out := []string{}
	for _, ver := range in {
		v, err := semver.NewVersion(ver)
		if err != nil {
			return nil, err
		}
		c = append(c, v)
	}

	sort.Sort(c)
	for _, v := range c {
		out = append(out, v.String())
	}

	return out, nil
}

func selectReleaseVersions(releases []*models.Release, ceil *semver.Version, floor *semver.Version) ([]string, error) {
	versions := []string{}
	if len(releases) < 1 {
		return versions, nil
	}

	for _, release := range releases {
		rv, err := semver.NewVersion(release.Version)
		if err != nil {
			return nil, err
		}

		if rv.GreaterThan(ceil) {
			// We've found releases that are too new
			continue
		}

		if rv.LessThan(floor) {
			// We've found a release after our floor so we can skip it.
			continue
		}

		// We're in-between
		versions = append(versions, rv.String())
	}

	return versions, nil
}
