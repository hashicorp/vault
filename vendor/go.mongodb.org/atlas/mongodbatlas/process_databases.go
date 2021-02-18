package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const processesDatabasesPath = "groups/%s/processes/%s:%d/databases"

// ProcessDatabasesService is an interface for interfacing with the Process Measurements
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-databases/
type ProcessDatabasesService interface {
	List(context.Context, string, string, int, *ListOptions) (*ProcessDatabasesResponse, *Response, error)
}

// ProcessDatabasesServiceOp handles communication with the process disks related methods of the
// MongoDB Atlas API
type ProcessDatabasesServiceOp service

var _ ProcessDatabasesService = &ProcessDatabasesServiceOp{}

// ProcessDatabasesResponse is the response from the ProcessDatabasesService.List.
type ProcessDatabasesResponse struct {
	Links      []*Link            `json:"links"`
	Results    []*ProcessDatabase `json:"results"`
	TotalCount int                `json:"totalCount"`
}

// ProcessDatabase is the database information of a process
type ProcessDatabase struct {
	Links        []*Link `json:"links"`
	DatabaseName string  `json:"databaseName"`
}

// List gets databases for a specific Atlas MongoDB process.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-databases/
func (s *ProcessDatabasesServiceOp) List(ctx context.Context, groupID, host string, port int, opts *ListOptions) (*ProcessDatabasesResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if host == "" {
		return nil, nil, NewArgError("host", "must be set")
	}

	if port <= 0 {
		return nil, nil, NewArgError("port", "must be valid")
	}

	basePath := fmt.Sprintf(processesDatabasesPath, groupID, host, port)

	// Add query params from listOptions
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDatabasesResponse)
	resp, err := s.Client.Do(ctx, req, root)
	return root, resp, err
}
