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

const processDatabaseMeasurementsPath = processesDatabasesPath + "/%s/measurements"

// ProcessDatabaseMeasurementsService is an interface for interfacing with the process database measurements
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-databases-measurements/
type ProcessDatabaseMeasurementsService interface {
	List(context.Context, string, string, int, string, *ProcessMeasurementListOptions) (*ProcessDatabaseMeasurements, *Response, error)
}

// ProcessDatabaseMeasurementsServiceOp handles communication with the process database measurements related methods of the
// MongoDB Atlas API.
type ProcessDatabaseMeasurementsServiceOp service

// ProcessDatabaseMeasurements represents a MongoDB process database measurements.
type ProcessDatabaseMeasurements struct {
	*ProcessMeasurements
	DatabaseName string `json:"databaseName"`
}

var _ ProcessDatabaseMeasurementsService = &ProcessDatabaseMeasurementsServiceOp{}

// List list measurements for a specific Atlas MongoDB database.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-databases-measurements/
func (s *ProcessDatabaseMeasurementsServiceOp) List(ctx context.Context, groupID, hostName string, port int, databaseName string, opts *ProcessMeasurementListOptions) (*ProcessDatabaseMeasurements, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if hostName == "" {
		return nil, nil, NewArgError("hostName", "must be set")
	}

	if databaseName == "" {
		return nil, nil, NewArgError("databaseName", "must be set")
	}

	basePath := fmt.Sprintf(processDatabaseMeasurementsPath, groupID, hostName, port, databaseName)

	// Add query params from listOptions
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDatabaseMeasurements)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
