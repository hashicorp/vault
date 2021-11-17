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

const (
	maintenanceWindowsPath = "api/atlas/v1.0/groups/%s/maintenanceWindow"
	// Sunday day of the week when you would like the maintenance window to start.
	Sunday = 1
	// Monday day of the week when you would like the maintenance window to start.
	Monday = 2
	// Tuesday day of the week when you would like the maintenance window to start.
	Tuesday = 3
	// Wednesday day of the week when you would like the maintenance window to start.
	Wednesday = 4
	// Thursday day of the week when you would like the maintenance window to start.
	Thursday = 5
	// Friday day of the week when you would like the maintenance window to start.
	Friday = 6
	// Saturday day of the week when you would like the maintenance window to start.
	Saturday = 7
)

// MaintenanceWindowsService is an interface for interfacing with the Maintenance Windows
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/maintenance-windows/
type MaintenanceWindowsService interface {
	Get(context.Context, string) (*MaintenanceWindow, *Response, error)
	Update(context.Context, string, *MaintenanceWindow) (*Response, error)
	Defer(context.Context, string) (*Response, error)
	AutoDefer(context.Context, string) (*Response, error)
	Reset(context.Context, string) (*Response, error)
}

// MaintenanceWindowsServiceOp handles communication with the MaintenanceWindows related methods
// of the MongoDB Atlas API.
type MaintenanceWindowsServiceOp service

var _ MaintenanceWindowsService = &MaintenanceWindowsServiceOp{}

// MaintenanceWindow represents MongoDB Maintenance Windows.
type MaintenanceWindow struct {
	DayOfWeek            int   `json:"dayOfWeek,omitempty"`            // Day of the week when you would like the maintenance window to start as a 1-based integer.Sunday 	1, Monday 	2, Tuesday 	3, Wednesday 	4, Thursday 5, Friday 6, Saturday 7
	HourOfDay            *int  `json:"hourOfDay,omitempty"`            // Hour of the day when you would like the maintenance window to start. This parameter uses the 24-hour clock, where midnight is 0, noon is 12.
	StartASAP            *bool `json:"startASAP,omitempty"`            // Flag indicating whether project maintenance has been directed to start immediately.
	NumberOfDeferrals    int   `json:"numberOfDeferrals,omitempty"`    // Number of times the current maintenance event for this project has been deferred.
	AutoDeferOnceEnabled *bool `json:"autoDeferOnceEnabled,omitempty"` // Flag that indicates whether you want to defer all maintenance windows one week they would be triggered.
}

// Get gets the current user-defined maintenance window for the given project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/maintenance-windows-view-in-one-project/
func (s *MaintenanceWindowsServiceOp) Get(ctx context.Context, groupID string) (*MaintenanceWindow, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(maintenanceWindowsPath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(MaintenanceWindow)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update the current maintenance window for the given project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/maintenance-window-update/
func (s *MaintenanceWindowsServiceOp) Update(ctx context.Context, groupID string, updateRequest *MaintenanceWindow) (*Response, error) {
	if updateRequest == nil {
		return nil, NewArgError("updateRequest", "cannot be nil")
	}
	if groupID == "" {
		return nil, NewArgError("groupID", "cannot be nil")
	}

	path := fmt.Sprintf(maintenanceWindowsPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Defer maintenance for the given project for one week.
//
// See more: https://docs.atlas.mongodb.com/reference/api/maintenance-window-defer/
func (s *MaintenanceWindowsServiceOp) Defer(ctx context.Context, groupID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "cannot be nil")
	}

	path := fmt.Sprintf(maintenanceWindowsPath+"/defer", groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// AutoDefer any scheduled maintenance for the given project for one week.
//
// See more: https://docs.atlas.mongodb.com/reference/api/maintenance-window-auto-defer/
func (s *MaintenanceWindowsServiceOp) AutoDefer(ctx context.Context, groupID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "cannot be nil")
	}

	path := fmt.Sprintf(maintenanceWindowsPath+"/autoDefer", groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// Reset clears the current maintenance window for the given project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/maintenance-window-clear/
func (s *MaintenanceWindowsServiceOp) Reset(ctx context.Context, groupID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(maintenanceWindowsPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
