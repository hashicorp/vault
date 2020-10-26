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
	"path"
	"sort"

	"github.com/joyent/triton-go/client"
	"github.com/pkg/errors"
)

type ServicesClient struct {
	client *client.Client
}

type Service struct {
	Name     string
	Endpoint string
}

type ListServicesInput struct{}

func (c *ServicesClient) List(ctx context.Context, _ *ListServicesInput) ([]*Service, error) {
	fullPath := path.Join("/", c.client.AccountName, "services")
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list services")
	}

	var intermediate map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&intermediate); err != nil {
		return nil, errors.Wrap(err, "unable to decode list services response")
	}

	keys := make([]string, len(intermediate))
	i := 0
	for k := range intermediate {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	result := make([]*Service, len(intermediate))
	i = 0
	for _, key := range keys {
		result[i] = &Service{
			Name:     key,
			Endpoint: intermediate[key],
		}
		i++
	}

	return result, nil
}
