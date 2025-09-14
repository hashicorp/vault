// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
)

// GetLatestProductVersionReq is an HCP image service request to get the latest
// product version. It can also be used to get information for other images
// when configured with different constraints.
type GetLatestProductVersionReq struct {
	ProductName                  string
	ProductVersionConstraint     string
	HostManagerVersionConstraint string
	CloudProvider                string
	CloudRegion                  string
	Availability                 GetLatestProductVersionAvailability
	ExcludeReleaseCandidates     bool
}

// GetLatestProductVersionAvailability describes the availability state of an
// image.
type GetLatestProductVersionAvailability string

// GetLatestProductVersionRes is a response from a request to get the latest
// image from the HCP image service.
type GetLatestProductVersionRes struct {
	Response *http.Response `json:"-"`
	Image    *HCPImage      `json:"image,omitempty"`
}

// HCPRegion is a cloud region for the image.
type HCPRegion struct {
	Provider string `json:"provider,omitempty"`
	Region   string `json:"region,omitempty"`
}

// HCPImageReference is the image reference information.
type HCPImageReference struct {
	ImageID string     `json:"image_id,omitempty"`
	Region  *HCPRegion `json:"region,omitempty"`
}

// HCPImage is an image in the HCP image service.
type HCPImage struct {
	ID                 string             `json:"id,omitempty"`
	ProductName        string             `json:"product_name,omitempty"`
	ProductVersion     string             `json:"product_version,omitempty"`
	HostManagerVersion string             `json:"host_manager_version,omitempty"`
	Region             *HCPRegion         `json:"region,omitempty"`
	Availability       string             `json:"availability,omitempty"`
	AWS                *HCPImageReference `json:"aws,omitempty"`
	Azure              *HCPImageReference `json:"azure,omitempty"`
	OSVersion          string             `json:"os_version,omitempty"`
	CreatedAt          time.Time          `json:"created_at,omitempty"`
	UpdatedAt          time.Time          `json:"updated_at,omitempty"`
}

const (
	GetLatestProductVersionAvailabilityUnknown  GetLatestProductVersionAvailability = ""
	GetLatestProductVersionAvailabilityDisabled GetLatestProductVersionAvailability = "disabled"
	GetLatestProductVersionAvailabilityInternal GetLatestProductVersionAvailability = "internal"
	GetLatestProductVersionAvailabilityPublic   GetLatestProductVersionAvailability = "public"
	GetLatestProductVersionAvailabilityBeta     GetLatestProductVersionAvailability = "beta"
)

// ID returns the availability into the corresponding integer enum.
func (g GetLatestProductVersionAvailability) ID() string {
	switch g {
	case GetLatestProductVersionAvailabilityDisabled:
		return "1"
	case GetLatestProductVersionAvailabilityInternal:
		return "2"
	case GetLatestProductVersionAvailabilityPublic:
		return "3"
	case GetLatestProductVersionAvailabilityBeta:
		return "4"
	default:
		return "0"
	}
}

const imageServicePath = "image/2009-12-19/.internal/latestproductversion"

// Request takes an environment and produces an HTTP request that the client
// can execute.
func (r *GetLatestProductVersionReq) Request(env Environment) (*http.Request, error) {
	reqURL, err := url.JoinPath(env.Addr(), imageServicePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if r.ProductName != "" {
		query.Add("product_name", r.ProductName)
	}
	if r.ProductVersionConstraint != "" {
		query.Add("product_version_constraint", r.ProductVersionConstraint)
	}
	if r.HostManagerVersionConstraint != "" {
		query.Add("host_manager_version_constraint", r.HostManagerVersionConstraint)
	}
	if r.CloudProvider != "" {
		query.Add("region.provider", r.CloudProvider)
	}
	if r.CloudRegion != "" {
		query.Add("region.region", r.CloudRegion)
	}
	if r.Availability != "" {
		query.Add("availability", r.Availability.ID())
	}
	if r.ExcludeReleaseCandidates {
		query.Add("exclude_release_candidates", fmt.Sprintf("%t", r.ExcludeReleaseCandidates))
	}

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// Run runs the request to find an HCP image that matches the request criteria.
func (r *GetLatestProductVersionReq) Run(ctx context.Context, client *Client) (*GetLatestProductVersionRes, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	res := &GetLatestProductVersionRes{}

	ctx = slogctx.Append(ctx,
		slog.String("availability", string(r.Availability)),
		slog.String("availability-id", r.Availability.ID()),
		slog.String("cloud", r.CloudProvider),
		slog.Bool("exclude-release-candidates", r.ExcludeReleaseCandidates),
		slog.String("host-manager-version-constraint", r.HostManagerVersionConstraint),
		slog.String("product", r.ProductName),
		slog.String("product-version-constraint", r.ProductVersionConstraint),
		slog.String("region", r.CloudRegion),
	)
	slog.Default().DebugContext(ctx, "getting latest HCP product version")

	var err error
	res.Response, err = client.Do(ctx, r)
	if err != nil {
		return res, err
	}

	defer res.Response.Body.Close()
	bytes, err := io.ReadAll(res.Response.Body)
	if err != nil {
		return res, err
	}

	if err = json.Unmarshal(bytes, res); err != nil {
		return res, err
	}

	if res.Response.StatusCode > 299 {
		return res, fmt.Errorf("received unexpected http response code: %d", res.Response.StatusCode)
	}

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *GetLatestProductVersionRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling latest HCP image response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *GetLatestProductVersionRes) ToTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"name", "id", "cloud", "region", "version", "created_at"})

	var imageID string
	if aws := r.Image.AWS; aws != nil {
		imageID = aws.ImageID
	}
	if azure := r.Image.Azure; azure != nil {
		imageID = azure.ImageID
	}

	t.AppendRow(table.Row{
		r.Image.ProductName,
		imageID,
		r.Image.Region.Provider,
		r.Image.Region.Region,
		r.Image.ProductVersion,
		r.Image.CreatedAt.String(),
	})

	return t
}
