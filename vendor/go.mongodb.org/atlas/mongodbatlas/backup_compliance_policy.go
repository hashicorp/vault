// Copyright 2023 MongoDB Inc
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
	BackupCompliancePolicyBasePath = "api/atlas/v1.0/groups/%s/backupCompliancePolicy"
)

// BackupCompliancePolicyService is an interface for interfacing with the Backup Compliance Policy
// endpoints of the MongoDB Atlas API.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cloud-Backups/operation/updateDataProtectionSettings
type BackupCompliancePolicyService interface {
	Get(context.Context, string) (*BackupCompliancePolicy, *Response, error)
	Update(context.Context, string, *BackupCompliancePolicy) (*BackupCompliancePolicy, *Response, error)
}

// CloudProviderSnapshotBackupPolicyServiceOp handles communication with the BackupCompliancePolicyService related methods of the
// MongoDB Atlas API.
type BackupCompliancePolicyServiceOp service

var _ BackupCompliancePolicyService = &BackupCompliancePolicyServiceOp{}

// BackupCompliancePolicy represents a backup compiance policy.
type BackupCompliancePolicy struct {
	AuthorizedEmail         string                `json:"authorizedEmail,omitempty"`
	AuthorizedUserFirstName string                `json:"authorizedUserFirstName,omitempty"`
	AuthorizedUserLastName  string                `json:"authorizedUserLastName,omitempty"`
	CopyProtectionEnabled   *bool                 `json:"copyProtectionEnabled,omitempty"`
	EncryptionAtRestEnabled *bool                 `json:"encryptionAtRestEnabled,omitempty"`
	OnDemandPolicyItem      PolicyItem            `json:"onDemandPolicyItem,omitempty"`
	PitEnabled              *bool                 `json:"pitEnabled,omitempty"`
	ProjectID               string                `json:"projectId,omitempty"`
	RestoreWindowDays       *int64                `json:"restoreWindowDays,omitempty"`
	ScheduledPolicyItems    []ScheduledPolicyItem `json:"scheduledPolicyItems,omitempty"`
	State                   string                `json:"state,omitempty"`
	UpdatedDate             string                `json:"updatedDate,omitempty"`
	UpdatedUser             string                `json:"updatedUser,omitempty"`
}

// PolicyItem represents a specifications for a scheduled backup policy and on demand policy.
type ScheduledPolicyItem struct {
	ID                string `json:"id,omitempty"`                // Unique identifier of the backup policy item.
	FrequencyInterval int    `json:"frequencyInterval,omitempty"` // Desired frequency of the new backup policy item specified by frequencyType.
	FrequencyType     string `json:"frequencyType,omitempty"`     // Frequency associated with the backup policy item. One of the following values: hourly, daily, weekly or monthly.
	RetentionUnit     string `json:"retentionUnit,omitempty"`     // Metric of duration of the backup policy item: days, weeks, or months.
	RetentionValue    int    `json:"retentionValue,omitempty"`    // Duration for which the backup is kept. Associated with retentionUnit.
}

// Get gets the current snapshot schedule and retention settings for the cluster with {CLUSTER-NAME}.
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cloud-Backups/operation/getDataProtectionSettings
func (s *BackupCompliancePolicyServiceOp) Get(ctx context.Context, groupID string) (*BackupCompliancePolicy, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}

	path := fmt.Sprintf(BackupCompliancePolicyBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupCompliancePolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates the snapshot schedule or retention settings for the cluster with {CLUSTER-NAME}.
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cloud-Backups/operation/updateDataProtectionSettings
func (s *BackupCompliancePolicyServiceOp) Update(ctx context.Context, groupID string, createRequest *BackupCompliancePolicy) (*BackupCompliancePolicy, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}

	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(BackupCompliancePolicyBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupCompliancePolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
