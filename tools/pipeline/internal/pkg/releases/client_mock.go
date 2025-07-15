// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/releases-api/pkg/models"
)

var _ VersionLister = (*MockClient)(nil)

// MockClient is an in-memory mock of a releases API client for use in testing.
type MockClient struct {
	Versions []string
}

// NewMockClient takes a list of versions and returns a new mock releases API client.
func NewMockClient(versions []string) *MockClient {
	return &MockClient{Versions: versions}
}

// ListVersions takes a context, product, edition, a ceiling version and floor version and returns
// all unique versions between the ceiling and floor range.
func (m *MockClient) ListVersions(ctx context.Context, product string, edition LicenseClass, ceil, floor *semver.Version) ([]string, error) {
	releaseVersions := []*models.Release{}
	for _, v := range m.Versions {
		releaseVersions = append(releaseVersions, &models.Release{Version: v})
	}

	return selectReleaseVersions(releaseVersions, ceil, floor)
}
