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
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type PackagesClient struct {
	client *client.Client
}

type PackageDisk struct {
	Size       interface{}
	SizeInMiB  int64
	Remaining  bool
	OSDiskSize bool
}

func (d *PackageDisk) UnmarshalJSON(data []byte) error {
	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		log.Fatal(err)
		return err
	}
	if decoded["size"] != nil {
		if err := json.Unmarshal(decoded["size"], &d.Size); err != nil {
			log.Fatal(err)
			return err
		}
	}
	switch d.Size.(type) {
	case string:
		d.Remaining = true
		d.OSDiskSize = false
	case nil:
		d.Remaining = false
		d.OSDiskSize = true
	default:
		d.Remaining = false
		d.OSDiskSize = false
		d.SizeInMiB = int64(d.Size.(float64))
	}
	return nil
}

type Package struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Memory       int64         `json:"memory"`
	Disk         int64         `json:"disk"`
	Swap         int64         `json:"swap"`
	LWPs         int64         `json:"lwps"`
	VCPUs        int64         `json:"vcpus"`
	Version      string        `json:"version"`
	Group        string        `json:"group"`
	Description  string        `json:"description"`
	Default      bool          `json:"default"`
	Brand        string        `json:"brand"`
	FlexibleDisk bool          `json:"flexible_disk,omitempty"`
	Disks        []PackageDisk `json:"disks,omitempty"`
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
	Brand   string `json:"brand"`
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
	if input.Brand != "" {
		query.Set("brand", input.Brand)
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
