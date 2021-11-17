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

const processMeasurementsPath = "api/atlas/v1.0/groups/%s/processes/%s:%d/measurements"

// ProcessMeasurementsService is an interface for interfacing with the Process Measurements
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-measurements/
type ProcessMeasurementsService interface {
	List(context.Context, string, string, int, *ProcessMeasurementListOptions) (*ProcessMeasurements, *Response, error)
}

// ProcessMeasurementsServiceOp handles communication with the Process Measurements related methods of the
// MongoDB Atlas API.
type ProcessMeasurementsServiceOp service

var _ ProcessMeasurementsService = &ProcessMeasurementsServiceOp{}

// ProcessMeasurements represents a MongoDB Process Measurements.
type ProcessMeasurements struct {
	End          string          `json:"end"`
	Granularity  string          `json:"granularity"`
	GroupID      string          `json:"groupId"`
	HostID       string          `json:"hostId"`
	Links        []*Link         `json:"links,omitempty"`
	Measurements []*Measurements `json:"measurements"`
	ProcessID    string          `json:"processId"`
	Start        string          `json:"start"`
}

// Measurements represents a MongoDB Measurement.
type Measurements struct {
	DataPoints []*DataPoints `json:"dataPoints,omitempty"`
	Name       string        `json:"name"`
	Units      string        `json:"units"`
}

// DataPoints represents a MongoDB DataPoints.
type DataPoints struct {
	Timestamp string   `json:"timestamp"`
	Value     *float32 `json:"value"`
}

// ProcessMeasurementListOptions contains the list of options for Process Measurements.
type ProcessMeasurementListOptions struct {
	*ListOptions
	Granularity string   `url:"granularity"`
	Period      string   `url:"period,omitempty"`
	Start       string   `url:"start,omitempty"`
	End         string   `url:"end,omitempty"`
	M           []string `url:"m,omitempty"`
}

// List lists measurements for a specific Atlas MongoDB process.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-measurements/
func (s *ProcessMeasurementsServiceOp) List(ctx context.Context, groupID, host string, port int, opts *ProcessMeasurementListOptions) (*ProcessMeasurements, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if host == "" {
		return nil, nil, NewArgError("host", "must be set")
	}

	basePath := fmt.Sprintf(processMeasurementsPath, groupID, host, port)

	// Add query params from listOptions
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessMeasurements)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
