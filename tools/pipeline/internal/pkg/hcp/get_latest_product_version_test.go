// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hcp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_GetLatestProductVersionReq_Request(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req         *GetLatestProductVersionReq
		expectedURL string
	}{
		"no query args": {
			&GetLatestProductVersionReq{},
			"https://api.hcp.dev/image/2009-12-19/.internal/latestproductversion",
		},
		"all query args": {
			&GetLatestProductVersionReq{
				Availability:             GetLatestProductVersionAvailabilityPublic,
				CloudRegion:              "us-east-1",
				CloudProvider:            "aws",
				ExcludeReleaseCandidates: true,
				ProductName:              "vault",
				ProductVersionConstraint: "1.21.0-beta1+ent-2cf0b2f",
			},
			"https://api.hcp.dev/image/2009-12-19/.internal/latestproductversion?availability=3&exclude_release_candidates=true&product_name=vault&product_version_constraint=1.21.0-beta1%2Bent-2cf0b2f&region.provider=aws&region.region=us-east-1",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := test.req.Request(EnvironmentDev)
			require.NoError(t, err)
			require.Equal(t, test.expectedURL, req.URL.String())
		})
	}
}

const hcpImageJSON = `{"image":{"id":"c613a76a-7c7f-484d-b21f-f99603dacee7","product_name":"vault","product_version":"v1.21.0-beta1+ent-bf26069","host_manager_version":"0.2.1806001022+1f8c65e9","region":{"provider":"aws","region":"us-west-2"},"availability":"PUBLIC","aws":{"image_id":"ami-037fc9428b5fc8d6a","region":{"provider":"aws","region":"us-west-2"}},"os_version":"","created_at":"2025-07-23T06:45:04.029Z","updated_at":"2025-07-23T06:45:04.008Z"}}`

func Test_GetLatestProductVersionRes_Unmarshal(t *testing.T) {
	t.Parallel()

	createdAt, err := time.Parse(time.RFC3339, "2025-07-23T06:45:04.029Z")
	require.NoError(t, err)
	updatedAt, err := time.Parse(time.RFC3339, "2025-07-23T06:45:04.008Z")
	require.NoError(t, err)

	res := &GetLatestProductVersionRes{}
	require.NoError(t, json.Unmarshal([]byte(hcpImageJSON), res))
	expect := &GetLatestProductVersionRes{
		Image: &HCPImage{
			ID:                 "c613a76a-7c7f-484d-b21f-f99603dacee7",
			ProductName:        "vault",
			ProductVersion:     "v1.21.0-beta1+ent-bf26069",
			HostManagerVersion: "0.2.1806001022+1f8c65e9",
			Region: &HCPRegion{
				Provider: "aws",
				Region:   "us-west-2",
			},
			Availability: "PUBLIC",
			AWS: &HCPImageReference{
				ImageID: "ami-037fc9428b5fc8d6a",
				Region: &HCPRegion{
					Provider: "aws",
					Region:   "us-west-2",
				},
			},
			OSVersion: "",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}
	require.EqualValues(t, expect, res)
}
