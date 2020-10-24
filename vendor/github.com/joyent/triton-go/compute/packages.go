//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type PackagesClient struct {
	client *client.Client
}

type Package struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Memory      int64  `json:"memory"`
	Disk        int64  `json:"disk"`
	Swap        int64  `json:"swap"`
	LWPs        int64  `json:"lwps"`
	VCPUs       int64  `json:"vcpus"`
	Version     string `json:"version"`
	Group       string `json:"group"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

type ListPackagesInput struct {
	Name    string `json:"name"`
	Memory  int64  `json:"memory"`
	Disk    int64  `json:"disk"`
	Swap    int64  `json:"swap"`
	LWPs    int64  `json:"lwps"`
	VCPUs   int64  `json:"vcpus"`
	Version string `json:"version"`
	Group   string `json:"group"`
}

func (c *PackagesClient) List(ctx context.Context, input *ListPackagesInput) ([]*Package, error) {
	fullPath := path.Join("/", c.client.AccountName, "packages")

	query := &url.Values{}
	if input.Name != "" {
		query.Set("name", input.Name)
	}
	if input.Memory != 0 {
		query.Set("memory", fmt.Sprintf("%d", input.Memory))
	}
	if input.Disk != 0 {
		query.Set("disk", fmt.Sprintf("%d", input.Disk))
	}
	if input.Swap != 0 {
		query.Set("swap", fmt.Sprintf("%d", input.Swap))
	}
	if input.LWPs != 0 {
		query.Set("lwps", fmt.Sprintf("%d", input.LWPs))
	}
	if input.VCPUs != 0 {
		query.Set("vcpus", fmt.Sprintf("%d", input.VCPUs))
	}
	if input.Version != "" {
		query.Set("version", input.Version)
	}
	if input.Group != "" {
		query.Set("group", input.Group)
	}

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
		Query:  query,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list packages")
	}

	var result []*Package
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode list packages response")
	}

	return result, nil
}

type GetPackageInput struct {
	ID string
}

func (c *PackagesClient) Get(ctx context.Context, input *GetPackageInput) (*Package, error) {
	fullPath := path.Join("/", c.client.AccountName, "packages", input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get package")
	}

	var result *Package
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode get package response")
	}

	return result, nil
}
