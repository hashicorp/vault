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
	"time"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type ImagesClient struct {
	client *client.Client
}

type ImageFile struct {
	Compression string `json:"compression"`
	SHA1        string `json:"sha1"`
	Size        int64  `json:"size"`
}

type Image struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	OS           string                 `json:"os"`
	Description  string                 `json:"description"`
	Version      string                 `json:"version"`
	Type         string                 `json:"type"`
	Requirements map[string]interface{} `json:"requirements"`
	Homepage     string                 `json:"homepage"`
	Files        []*ImageFile           `json:"files"`
	PublishedAt  time.Time              `json:"published_at"`
	Owner        string                 `json:"owner"`
	Public       bool                   `json:"public"`
	State        string                 `json:"state"`
	Tags         map[string]string      `json:"tags"`
	EULA         string                 `json:"eula"`
	ACL          []string               `json:"acl"`
}

type ListImagesInput struct {
	Name    string
	OS      string
	Version string
	Public  bool
	State   string
	Owner   string
	Type    string
}

func (c *ImagesClient) List(ctx context.Context, input *ListImagesInput) ([]*Image, error) {
	fullPath := path.Join("/", c.client.AccountName, "images")

	query := &url.Values{}
	if input.Name != "" {
		query.Set("name", input.Name)
	}
	if input.OS != "" {
		query.Set("os", input.OS)
	}
	if input.Version != "" {
		query.Set("version", input.Version)
	}
	if input.Public {
		query.Set("public", "true")
	}
	if input.State != "" {
		query.Set("state", input.State)
	}
	if input.Owner != "" {
		query.Set("owner", input.Owner)
	}
	if input.Type != "" {
		query.Set("type", input.Type)
	}

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
		Query:  query,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list images")
	}

	var result []*Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode list images response")
	}

	return result, nil
}

type GetImageInput struct {
	ImageID string
}

func (c *ImagesClient) Get(ctx context.Context, input *GetImageInput) (*Image, error) {
	fullPath := path.Join("/", c.client.AccountName, "images", input.ImageID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get image")
	}

	var result *Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode get image response")
	}

	return result, nil
}

type DeleteImageInput struct {
	ImageID string
}

func (c *ImagesClient) Delete(ctx context.Context, input *DeleteImageInput) error {
	fullPath := path.Join("/", c.client.AccountName, "images", input.ImageID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to delete image")
	}

	return nil
}

type ExportImageInput struct {
	ImageID   string
	MantaPath string
}

type MantaLocation struct {
	MantaURL     string `json:"manta_url"`
	ImagePath    string `json:"image_path"`
	ManifestPath string `json:"manifest_path"`
}

func (c *ImagesClient) Export(ctx context.Context, input *ExportImageInput) (*MantaLocation, error) {
	fullPath := path.Join("/", c.client.AccountName, "images", input.ImageID)
	query := &url.Values{}
	query.Set("action", "export")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  query,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to export image")
	}

	var result *MantaLocation
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode export image response")
	}

	return result, nil
}

type CreateImageFromMachineInput struct {
	MachineID   string            `json:"machine"`
	Name        string            `json:"name"`
	Version     string            `json:"version,omitempty"`
	Description string            `json:"description,omitempty"`
	HomePage    string            `json:"homepage,omitempty"`
	EULA        string            `json:"eula,omitempty"`
	ACL         []string          `json:"acl,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

func (c *ImagesClient) CreateFromMachine(ctx context.Context, input *CreateImageFromMachineInput) (*Image, error) {
	fullPath := path.Join("/", c.client.AccountName, "images")
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to create machine from image")
	}

	var result *Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode create machine from image response")
	}

	return result, nil
}

type UpdateImageInput struct {
	ImageID     string            `json:"-"`
	Name        string            `json:"name,omitempty"`
	Version     string            `json:"version,omitempty"`
	Description string            `json:"description,omitempty"`
	HomePage    string            `json:"homepage,omitempty"`
	EULA        string            `json:"eula,omitempty"`
	ACL         []string          `json:"acl,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

func (c *ImagesClient) Update(ctx context.Context, input *UpdateImageInput) (*Image, error) {
	fullPath := path.Join("/", c.client.AccountName, "images", input.ImageID)
	query := &url.Values{}
	query.Set("action", "update")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  query,
		Body:   input,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to update image")
	}

	var result *Image
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode update image response")
	}

	return result, nil
}
