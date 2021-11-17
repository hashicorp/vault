// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlas

import (
	"context"
	"net/http"
)

const ipInfoPath = "api/private/ipinfo"

// IPInfoService is used to determine the public ip address of the client
//
// We currently make no promise to support or document this service or endpoint
// beyond what can be seen here.
type IPInfoService interface {
	Get(context.Context) (*IPInfo, *Response, error)
}

// IPInfoServiceOp is an implementation of IPInfoService.
type IPInfoServiceOp service

var _ IPInfoService = &IPInfoServiceOp{}

type IPInfo struct {
	CurrentIPv4Address string `json:"currentIpv4Address"`
}

// Get gets the public ip address of the client.
func (s *IPInfoServiceOp) Get(ctx context.Context) (*IPInfo, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, ipInfoPath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(IPInfo)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
