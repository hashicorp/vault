package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudProviderSnapshotBackupPolicesBasePath = "groups/%s/clusters/%s/backup/schedule"
)

// CloudProviderSnapshotBackupPoliciesService is an interface for interfacing with the Cloud Provider Snapshots Backup Policy
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-schedule/
type CloudProviderSnapshotBackupPoliciesService interface {
	Get(context.Context, string, string) (*CloudProviderSnapshotBackupPolicy, *Response, error)
	Update(context.Context, string, string, *CloudProviderSnapshotBackupPolicy) (*CloudProviderSnapshotBackupPolicy, *Response, error)
}

// CloudProviderSnapshotBackupPoliciesServiceOp handles communication with the CloudProviderSnapshotBackupPoliciesService related methods of the
// MongoDB Atlas API
type CloudProviderSnapshotBackupPoliciesServiceOp service

var _ CloudProviderSnapshotBackupPoliciesService = &CloudProviderSnapshotBackupPoliciesServiceOp{}

// CloudProviderSnapshotBackupPolicy represents a cloud provider snapshot schedule.
type CloudProviderSnapshotBackupPolicy struct {
	ClusterID             string   `json:"clusterId,omitempty"`             //	Unique identifier of the Atlas cluster.
	ClusterName           string   `json:"clusterName,omitempty"`           //	Name of the Atlas cluster.
	ReferenceHourOfDay    *int64   `json:"referenceHourOfDay,omitempty"`    // UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
	ReferenceMinuteOfHour *int64   `json:"referenceMinuteOfHour,omitempty"` // UTC Minutes after referenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive. Number of days back in time you can restore to with point-in-time accuracy.
	RestoreWindowDays     *int64   `json:"restoreWindowDays,omitempty"`     // Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
	UpdateSnapshots       *bool    `json:"updateSnapshots,omitempty"`       // Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.
	NextSnapshot          string   `json:"nextSnapshot,omitempty"`          // UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
	Policies              []Policy `json:"policies,omitempty"`              // A list of policy definitions for the cluster.
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

	return root, resp, err
}
