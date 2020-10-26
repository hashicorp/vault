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
	"net/http"
	"net/url"
	"path"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type VolumesClient struct {
	client *client.Client
}

type Volume struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Owner          string   `json:"owner_uuid"`
	Type           string   `json:"type"`
	FileSystemPath string   `json:"filesystem_path"`
	Size           int64    `json:"size"`
	State          string   `json:"state"`
	Networks       []string `json:"networks"`
	Refs           []string `json:"refs"`
}

type ListVolumesInput struct {
	Name  string
	Size  string
	State string
	Type  string
}

func (c *VolumesClient) List(ctx context.Context, input *ListVolumesInput) ([]*Volume, error) {
	fullPath := path.Join("/", c.client.AccountName, "volumes")

	query := &url.Values{}
	if input.Name != "" {
		query.Set("name", input.Name)
	}
	if input.Size != "" {
		query.Set("size", input.Size)
	}
	if input.State != "" {
		query.Set("state", input.State)
	}
	if input.Type != "" {
		query.Set("type", input.Type)
	}

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
		Query:  query,
	}
	resp, err := c.client.ExecuteRequest(ctx, reqInputs)
	if resp != nil {
		defer resp.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list volumes")
	}

	var result []*Volume
	decoder := json.NewDecoder(resp)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode list volumes response")
	}

	return result, nil
}

type CreateVolumeInput struct {
	Name     string
	Size     int64
	Networks []string
	Type     string
}

func (input *CreateVolumeInput) toAPI() map[string]interface{} {
	result := make(map[string]interface{}, 0)

	if input.Name != "" {
		result["name"] = input.Name
	}

	if input.Size != 0 {
		result["size"] = input.Size
	}

	if input.Type != "" {
		result["type"] = input.Type
	}

	if len(input.Networks) > 0 {
		result["networks"] = input.Networks
	}

	return result
}

func (c *VolumesClient) Create(ctx context.Context, input *CreateVolumeInput) (*Volume, error) {
	fullPath := path.Join("/", c.client.AccountName, "volumes")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input.toAPI(),
	}
	resp, err := c.client.ExecuteRequest(ctx, reqInputs)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create volume")
	}
	if resp != nil {
		defer resp.Close()
	}

	var result *Volume
	decoder := json.NewDecoder(resp)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode create volume response")
	}

	return result, nil
}

type DeleteVolumeInput struct {
	ID string
}

func (c *VolumesClient) Delete(ctx context.Context, input *DeleteVolumeInput) error {
	fullPath := path.Join("/", c.client.AccountName, "volumes", input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	resp, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return errors.Wrap(err, "unable to delete volume")
	}
	if resp == nil {
		return errors.Wrap(err, "unable to delete volume")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return nil
	}

	return nil
}

type GetVolumeInput struct {
	ID string
}

func (c *VolumesClient) Get(ctx context.Context, input *GetVolumeInput) (*Volume, error) {
	fullPath := path.Join("/", c.client.AccountName, "volumes", input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	resp, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get volume")
	}
	if resp == nil {
		return nil, errors.Wrap(err, "unable to get volume")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	var result *Volume
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode get volume volume")
	}

	return result, nil
}

type UpdateVolumeInput struct {
	ID   string `json:"-"`
	Name string `json:"name"`
}

func (c *VolumesClient) Update(ctx context.Context, input *UpdateVolumeInput) error {
	fullPath := path.Join("/", c.client.AccountName, "volumes", input.ID)

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input,
	}
	resp, err := c.client.ExecuteRequest(ctx, reqInputs)
	if err != nil {
		return errors.Wrap(err, "unable to update volume")
	}
	if resp != nil {
		defer resp.Close()
	}

	return nil
}
