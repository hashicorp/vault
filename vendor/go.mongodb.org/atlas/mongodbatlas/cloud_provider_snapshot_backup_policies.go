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
	cloudProviderSnapshotBackupPolicesBasePath = "api/atlas/v1.0/groups/%s/clusters/%s/backup/schedule"
)

// CloudProviderSnapshotBackupPoliciesService is an interface for interfacing with the Cloud Provider Snapshots Backup Policy
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-schedule/
type CloudProviderSnapshotBackupPoliciesService interface {
	Get(context.Context, string, string) (*CloudProviderSnapshotBackupPolicy, *Response, error)
	Update(context.Context, string, string, *CloudProviderSnapshotBackupPolicy) (*CloudProviderSnapshotBackupPolicy, *Response, error)
	Delete(context.Context, string, string) (*CloudProviderSnapshotBackupPolicy, *Response, error)
}

// CloudProviderSnapshotBackupPoliciesServiceOp handles communication with the CloudProviderSnapshotBackupPoliciesService related methods of the
// MongoDB Atlas API.
type CloudProviderSnapshotBackupPoliciesServiceOp service

var _ CloudProviderSnapshotBackupPoliciesService = &CloudProviderSnapshotBackupPoliciesServiceOp{}

// CloudProviderSnapshotBackupPolicy represents a cloud provider snapshot schedule.
type CloudProviderSnapshotBackupPolicy struct {
	ClusterID                         string               `json:"clusterId,omitempty"`                         //	Unique identifier of the Atlas cluster.
	ClusterName                       string               `json:"clusterName,omitempty"`                       //	Name of the Atlas cluster.
	ReferenceHourOfDay                *int64               `json:"referenceHourOfDay,omitempty"`                // UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
	ReferenceMinuteOfHour             *int64               `json:"referenceMinuteOfHour,omitempty"`             // UTC Minutes after referenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive. Number of days back in time you can restore to with point-in-time accuracy.
	RestoreWindowDays                 *int64               `json:"restoreWindowDays,omitempty"`                 // Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
	UpdateSnapshots                   *bool                `json:"updateSnapshots,omitempty"`                   // Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.
	NextSnapshot                      string               `json:"nextSnapshot,omitempty"`                      // UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
	Policies                          []Policy             `json:"policies,omitempty"`                          // A list of policy definitions for the cluster.
	AutoExportEnabled                 *bool                `json:"autoExportEnabled,omitempty"`                 // Specify true to enable automatic export of cloud backup snapshots to the AWS bucket. You must also define the export policy using export. Specify false to disable automatic export.
	Export                            *Export              `json:"export,omitempty"`                            // Export struct that represents a policy for automatically exporting cloud backup snapshots to AWS bucket.
	UseOrgAndGroupNamesInExportPrefix *bool                `json:"useOrgAndGroupNamesInExportPrefix,omitempty"` // Specifies whether to use organization and project names instead of organization and project UUIDs in the path to the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots
	Links                             []*Link              `json:"links,omitempty"`                             // One or more links to sub-resources and/or related resources.
	CopySettings                      []CopySetting        `json:"copySettings"`                                // List that contains a document for each copy setting item in the desired backup policy.
	DeleteCopiedBackups               []DeleteCopiedBackup `json:"deleteCopiedBackups,omitempty"`               // List that contains a document for each deleted copy setting whose backup copies you want to delete.
}

// Policy represents for the snapshot and an array of backup policy items.
type Policy struct {
	ID          string       `json:"id,omitempty"`          // Unique identifier of the backup policy.
	PolicyItems []PolicyItem `json:"policyItems,omitempty"` // A list of specifications for a policy.
}

// PolicyItem represents a specifications for a backup policy.
type PolicyItem struct {
	ID                string `json:"id,omitempty"`                // Unique identifier of the backup policy item.
	FrequencyInterval int    `json:"frequencyInterval,omitempty"` // Desired frequency of the new backup policy item specified by frequencyType.
	FrequencyType     string `json:"frequencyType,omitempty"`     // Frequency associated with the backup policy item. One of the following values: hourly, daily, weekly or monthly.
	RetentionUnit     string `json:"retentionUnit,omitempty"`     // Metric of duration of the backup policy item: days, weeks, or months.
	RetentionValue    int    `json:"retentionValue,omitempty"`    // Duration for which the backup is kept. Associated with retentionUnit.
}

// Export represents a policy for automatically exporting cloud backup snapshots to AWS bucket.
type Export struct {
	ExportBucketID string `json:"exportBucketId,omitempty"` // Unique identifier of the AWS bucket to export the cloud backup snapshot to.
	FrequencyType  string `json:"frequencyType,omitempty"`  // Frequency associated with the export policy.
}

// CopySetting is autogenerated from the json schema.
type CopySetting struct {
	CloudProvider     *string  `json:"cloudProvider,omitempty"`     // Identifies the cloud provider that stores the snapshot copy.
	RegionName        *string  `json:"regionName,omitempty"`        // Target region to copy snapshots belonging to replicationSpecId to.
	ReplicationSpecID *string  `json:"replicationSpecId,omitempty"` // Unique identifier that identifies the replication object for a zone in a cluster.
	ShouldCopyOplogs  *bool    `json:"shouldCopyOplogs,omitempty"`  // Flag that indicates whether to copy the oplogs to the target region.
	Frequencies       []string `json:"frequencies,omitempty"`       // List that describes which types of snapshots to copy.
}

// DeleteCopiedBackup is autogenerated from the json schema.
type DeleteCopiedBackup struct {
	CloudProvider     *string `json:"cloudProvider,omitempty"`     // Identifies the cloud provider that stores the snapshot copy.
	RegionName        *string `json:"regionName,omitempty"`        // Target region to copy snapshots belonging to replicationSpecId to.
	ReplicationSpecID *string `json:"replicationSpecId,omitempty"` // Unique identifier that identifies the replication object for a zone in a cluster.
}

// Get gets the current snapshot schedule and retention settings for the cluster with {CLUSTER-NAME}.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-schedule-get-all/
func (s *CloudProviderSnapshotBackupPoliciesServiceOp) Get(ctx context.Context, groupID, clusterName string) (*CloudProviderSnapshotBackupPolicy, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(cloudProviderSnapshotBackupPolicesBasePath, groupID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotBackupPolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, err
}

// Update updates the snapshot schedule or retention settings for the cluster with {CLUSTER-NAME}.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-schedule-modify-one/
func (s *CloudProviderSnapshotBackupPoliciesServiceOp) Update(ctx context.Context, groupID, clusterName string, createRequest *CloudProviderSnapshotBackupPolicy) (*CloudProviderSnapshotBackupPolicy, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(cloudProviderSnapshotBackupPolicesBasePath, groupID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotBackupPolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, err
}

// Delete deletes all cloud backup schedules.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/delete-all-schedules/
func (s *CloudProviderSnapshotBackupPoliciesServiceOp) Delete(ctx context.Context, groupID, clusterName string) (*CloudProviderSnapshotBackupPolicy, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(cloudProviderSnapshotBackupPolicesBasePath, groupID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotBackupPolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
