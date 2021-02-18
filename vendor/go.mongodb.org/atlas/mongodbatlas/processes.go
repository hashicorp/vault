package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const processesPath = "groups/%s/processes"

// ProcessesService provides access to the alert processes related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/monitoring-and-logs/
type ProcessesService interface {
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
}

// processesResponse is the response from Processes.List.
type processesResponse struct {
	Links      []*Link    `json:"links,omitempty"`
	Results    []*Process `json:"results,omitempty"`
	TotalCount int        `json:"totalCount,omitempty"`
}

// ProcessesListOptions filter options for the processes API
type ProcessesListOptions struct {
	ListOptions
	ClusterID string `url:"clusterId,omitempty"`
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
