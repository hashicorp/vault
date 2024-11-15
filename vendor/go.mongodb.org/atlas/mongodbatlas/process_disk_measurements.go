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

const processDiskMeasurementsPath = processesDisksPath + "/%s/measurements"

// ProcessDiskMeasurementsService is an interface for interfacing with the Process Disk Measurements
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/process-disks-measurements/#get-measurements-of-a-disk-for-a-mongodb-process
type ProcessDiskMeasurementsService interface {
	List(context.Context, string, string, int, string, *ProcessMeasurementListOptions) (*ProcessDiskMeasurements, *Response, error)
}

// ProcessDiskMeasurementsServiceOp handles communication with the Process Disk Measurements related methods of the
// MongoDB Atlas API.
type ProcessDiskMeasurementsServiceOp service

// ProcessDiskMeasurements represents a MongoDB Process Disk Measurements.
type ProcessDiskMeasurements struct {
	*ProcessMeasurements
	PartitionName string `json:"partitionName"`
}

var _ ProcessDiskMeasurementsService = &ProcessDiskMeasurementsServiceOp{}

// List lists measurements for a specific Atlas MongoDB disk.
//
// See more: https://docs.atlas.mongodb.com/reference/api/process-disks-measurements/#get-measurements-of-a-disk-for-a-mongodb-process
func (s *ProcessDiskMeasurementsServiceOp) List(ctx context.Context, groupID, hostName string, port int, diskName string, opts *ProcessMeasurementListOptions) (*ProcessDiskMeasurements, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if hostName == "" {
		return nil, nil, NewArgError("hostName", "must be set")
	}

	if diskName == "" {
		return nil, nil, NewArgError("diskName", "must be set")
	}

	basePath := fmt.Sprintf(processDiskMeasurementsPath, groupID, hostName, port, diskName)

	// Add query params from listOptions
	path, err := setListOptions(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDiskMeasurements)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
