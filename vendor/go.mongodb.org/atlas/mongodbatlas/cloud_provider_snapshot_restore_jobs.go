package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudProviderSnapshotRestoreJobBasePath = "groups"
)

// CloudProviderSnapshotRestoreJobsService is an interface for interfacing with the CloudProviderSnapshotRestoreJobs
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/cloudProviderSnapshotRestoreJobs/
type CloudProviderSnapshotRestoreJobsService interface {
	List(context.Context, *SnapshotReqPathParameters, *ListOptions) (*CloudProviderSnapshotRestoreJobs, *Response, error)
	Get(context.Context, *SnapshotReqPathParameters) (*CloudProviderSnapshotRestoreJob, *Response, error)
	Create(context.Context, *SnapshotReqPathParameters, *CloudProviderSnapshotRestoreJob) (*CloudProviderSnapshotRestoreJob, *Response, error)
	Delete(context.Context, *SnapshotReqPathParameters) (*Response, error)
}

// CloudProviderSnapshotRestoreJobsServiceOp handles communication with the CloudProviderSnapshotRestoreJobs related methods of the
// MongoDB Atlas API
type CloudProviderSnapshotRestoreJobsServiceOp service

var _ CloudProviderSnapshotRestoreJobsService = &CloudProviderSnapshotRestoreJobsServiceOp{}

// CloudProviderSnapshotRestoreJob represents the structure of a cloudProviderSnapshotRestoreJob.
type CloudProviderSnapshotRestoreJob struct {
	ID                    string   `json:"id,omitempty"`                    // The unique identifier of the restore job.
	SnapshotID            string   `json:"snapshotId,omitempty"`            // Unique identifier of the snapshot to restore.
	DeliveryType          string   `json:"deliveryType,omitempty"`          // Type of restore job to create. Possible values are: automated or download or pointInTime
	DeliveryURL           []string `json:"deliveryUrl,omitempty"`           // One or more URLs for the compressed snapshot files for manual download. Only visible if deliveryType is download.
	TargetClusterName     string   `json:"targetClusterName,omitempty"`     // Name of the target Atlas cluster to which the restore job restores the snapshot. Only required if deliveryType is automated.
	TargetGroupID         string   `json:"targetGroupId,omitempty"`         // Unique ID of the target Atlas project for the specified targetClusterName. Only required if deliveryType is automated.
	Cancelled             bool     `json:"cancelled,omitempty"`             // Indicates whether the restore job was canceled.
	CreatedAt             string   `json:"createdAt,omitempty"`             // UTC ISO 8601 formatted point in time when Atlas created the restore job.
	Expired               bool     `json:"expired,omitempty"`               // Indicates whether the restore job expired.
	ExpiresAt             string   `json:"expiresAt,omitempty"`             // UTC ISO 8601 formatted point in time when the restore job expires.
	FinishedAt            string   `json:"finishedAt,omitempty"`            // UTC ISO 8601 formatted point in time when the restore job completed.
	Links                 []*Link  `json:"links,omitempty"`                 // One or more links to sub-resources and/or related resources. The relations between URLs are explained in the Web Linking Specification.
	Timestamp             string   `json:"timestamp,omitempty"`             // Timestamp in ISO 8601 date and time format in UTC when the snapshot associated to snapshotId was taken.
	OplogTs               int64    `json:"oplogTs,omitempty"`               //nolint:stylecheck // not changing this // Timestamp in the number of seconds that have elapsed since the UNIX epoch from which to you want to restore this snapshot. This is the first part of an Oplog timestamp.
	OplogInc              int64    `json:"oplogInc,omitempty"`              // Oplog operation number from which to you want to restore this snapshot. This is the second part of an Oplog timestamp.
	PointInTimeUTCSeconds int64    `json:"pointInTimeUTCSeconds,omitempty"` // Timestamp in the number of seconds that have elapsed since the UNIX epoch from which you want to restore this snapshot.
}

// CloudProviderSnapshotRestoreJobs represents an array of cloudProviderSnapshotRestoreJob
type CloudProviderSnapshotRestoreJobs struct {
	Links      []*Link                            `json:"links"`
	Results    []*CloudProviderSnapshotRestoreJob `json:"results"`
	TotalCount int                                `json:"totalCount"`
}

// List gets all cloud provider snapshot restore jobs for the specified cluster.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs-get-all/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) List(ctx context.Context, requestParameters *SnapshotReqPathParameters, listOptions *ListOptions) (*CloudProviderSnapshotRestoreJobs, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/restoreJobs", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.ClusterName)
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotRestoreJobs)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get gets one cloud provider snapshot restore jobs for the specified cluster.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs-get-one/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) Get(ctx context.Context, requestParameters *SnapshotReqPathParameters) (*CloudProviderSnapshotRestoreJob, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if requestParameters.JobID == "" {
		return nil, nil, NewArgError("jobId", "must be set")
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/restoreJobs/%s", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.ClusterName, requestParameters.JobID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotRestoreJob)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates a new restore job from a cloud provider snapshot associated to the specified cluster.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs-create-one/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) Create(ctx context.Context, requestParameters *SnapshotReqPathParameters, createRequest *CloudProviderSnapshotRestoreJob) (*CloudProviderSnapshotRestoreJob, *Response, error) {
	// Verify if is download or automated
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if createRequest.DeliveryType == "download" {
		createRequest.TargetClusterName = ""
		createRequest.TargetGroupID = ""
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/restoreJobs", cloudProviderSnapshotRestoreJobBasePath, requestParameters.GroupID, requestParameters.ClusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotRestoreJob)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, err
}

// Delete cancels the cloud provider snapshot manual download restore job associated to {JOB-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs-delete-one/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) Delete(ctx context.Context, requestParameters *SnapshotReqPathParameters) (*Response, error) {
	if requestParameters.GroupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}
	if requestParameters.JobID == "" {
		return nil, NewArgError("jobId", "must be set")
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/restoreJobs/%s", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.ClusterName, requestParameters.JobID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
