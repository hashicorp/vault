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

const processesPath = "api/atlas/v1.0/groups/%s/processes"

// ProcessesService provides access to the alert processes related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/monitoring-and-logs/
type ProcessesService interface {
	Get(context.Context, string, string, int) (*Process, *Response, error)
	List(context.Context, string, *ProcessesListOptions) ([]*Process, *Response, error)
}

// ProcessesServiceOp handles communication with the Process related methods
// of the MongoDB Atlas API.
type ProcessesServiceOp service

var _ ProcessesService = &ProcessesServiceOp{}

// Process represents a MongoDB process.
type Process struct {
	Created        string  `json:"created"`
	GroupID        string  `json:"groupId"`
	Hostname       string  `json:"hostname"`
	ID             string  `json:"id"`
	LastPing       string  `json:"lastPing"`
	Links          []*Link `json:"links"`
	Port           int     `json:"port"`
	ShardName      string  `json:"shardName"`
	ReplicaSetName string  `json:"replicaSetName"`
	TypeName       string  `json:"typeName"`
	Version        string  `json:"version"`
	UserAlias      string  `json:"userAlias"`
}

// processesResponse is the response from Processes.List.
type processesResponse struct {
	Links      []*Link    `json:"links,omitempty"`
	Results    []*Process `json:"results,omitempty"`
	TotalCount int        `json:"totalCount,omitempty"`
}

// ProcessesListOptions filter options for the processes API.
type ProcessesListOptions struct {
	ListOptions
	ClusterID string `url:"clusterId,omitempty"` // ClusterID is only available for Ops Manager and CLoud Manager.
}

// Get information for the specified Atlas MongoDB process in the specified project.
// An Atlas MongoDB process can be either a mongod or a mongos.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/processes-get-one/
func (s *ProcessesServiceOp) Get(ctx context.Context, groupID, hostname string, port int) (*Process, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if hostname == "" {
		return nil, nil, NewArgError("hostname", "must be set")
	}
	if port == 0 {
		return nil, nil, NewArgError("port", "must be set")
	}
	path := fmt.Sprintf(processesPath, groupID)
	path = fmt.Sprintf("%s/%s:%d", path, hostname, port)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Process)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// List lists all processes in the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/processes-get-all/
func (s *ProcessesServiceOp) List(ctx context.Context, groupID string, listOptions *ProcessesListOptions) ([]*Process, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(processesPath, groupID)

	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(processesResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}
