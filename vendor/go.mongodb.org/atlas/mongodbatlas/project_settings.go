// Copyright 2022 MongoDB Inc
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

const projectSettingsBasePath = projectBasePath + "/%s/settings"

type ProjectSettings struct {
	IsCollectDatabaseSpecificsStatisticsEnabled *bool `json:"isCollectDatabaseSpecificsStatisticsEnabled,omitempty"`
	IsDataExplorerEnabled                       *bool `json:"isDataExplorerEnabled,omitempty"`
	IsExtendedStorageSizesEnabled               *bool `json:"isExtendedStorageSizesEnabled,omitempty"`
	IsPerformanceAdvisorEnabled                 *bool `json:"isPerformanceAdvisorEnabled,omitempty"`
	IsRealtimePerformancePanelEnabled           *bool `json:"isRealtimePerformancePanelEnabled,omitempty"`
	IsSchemaAdvisorEnabled                      *bool `json:"isSchemaAdvisorEnabled,omitempty"`
}

// GetProjectSettings gets details about the settings for specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/project-settings-get-one/
func (s *ProjectsServiceOp) GetProjectSettings(ctx context.Context, groupID string) (*ProjectSettings, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(projectSettingsBasePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root *ProjectSettings
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// UpdateProjectSettings updates the settings for the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api/project-settings-update-one/
func (s *ProjectsServiceOp) UpdateProjectSettings(ctx context.Context, groupID string, projectSettings *ProjectSettings) (*ProjectSettings, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(projectSettingsBasePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, projectSettings)
	if err != nil {
		return nil, nil, err
	}

	var root *ProjectSettings
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
