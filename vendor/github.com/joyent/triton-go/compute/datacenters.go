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
	"path"
	"sort"

	"github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/errors"
	pkgerrors "github.com/pkg/errors"
)

type DataCentersClient struct {
	client *client.Client
}

type DataCenter struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ListDataCentersInput struct{}

func (c *DataCentersClient) List(ctx context.Context, _ *ListDataCentersInput) ([]*DataCenter, error) {
	fullPath := path.Join("/", c.client.AccountName, "datacenters")

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to list data centers")
	}

	var intermediate map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&intermediate); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode list data centers response")
	}

	keys := make([]string, len(intermediate))
	i := 0
	for k := range intermediate {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	result := make([]*DataCenter, len(intermediate))
	i = 0
	for _, key := range keys {
		result[i] = &DataCenter{
			Name: key,
			URL:  intermediate[key],
		}
		i++
	}

	return result, nil
}

type GetDataCenterInput struct {
	Name string
}

func (c *DataCentersClient) Get(ctx context.Context, input *GetDataCenterInput) (*DataCenter, error) {
	dcs, err := c.List(ctx, &ListDataCentersInput{})
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to get data center")
	}

	for _, dc := range dcs {
		if dc.Name == input.Name {
			return &DataCenter{
				Name: input.Name,
				URL:  dc.URL,
			}, nil
		}
	}

	return nil, &errors.APIError{
		StatusCode: http.StatusNotFound,
		Code:       "ResourceNotFound",
		Message:    fmt.Sprintf("data center %q not found", input.Name),
	}
}
