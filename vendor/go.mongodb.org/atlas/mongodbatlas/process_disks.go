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
	"fmt"
	"net/http"
)

const processesDisksPath = "api/atlas/v1.0/groups/%s/processes/%s:%d/disks"

// ProcessDisksService is an interface for interfacing with the Process Measurements
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-disks/
type ProcessDisksService interface {
	List(context.Context, string, string, int, *ListOptions) (*ProcessDisksResponse, *Response, error)
}

// ProcessDisksServiceOp handles communication with the process disks related methods of the
// MongoDB Atlas API.
type ProcessDisksServiceOp service

var _ ProcessDisksService = &ProcessDisksServiceOp{}

// ProcessDisksResponse is the response from the ProcessDisksService.List.
type ProcessDisksResponse struct {
	Links      []*Link        `json:"links"`
	Results    []*ProcessDisk `json:"results"`
	TotalCount int            `json:"totalCount"`
}

// ProcessDisk is the partition information of a process.
type ProcessDisk struct {
	Links         []*Link `json:"links"`
	PartitionName string  `json:"partitionName"`
}

// List gets partitions for a specific Atlas MongoDB process.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-disks/
func (s *ProcessDisksServiceOp) List(ctx context.Context, groupID, host string, port int, opts *ListOptions) (*ProcessDisksResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if host == "" {
		return nil, nil, NewArgError("host", "must be set")
	}

	if port <= 0 {
		return nil, nil, NewArgError("port", "must be valid")
	}

	basePath := fmt.Sprintf(processesDisksPath, groupID, host, port)

	// Add query params from listOptions
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDisksResponse)
	resp, err := s.Client.Do(ctx, req, root)
	return root, resp, err
}
