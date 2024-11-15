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

const cloudProviderSnapshotRestoreJobBasePath = "api/atlas/v1.0/groups"

// CloudProviderSnapshotRestoreJobsService is an interface for interfacing with the CloudProviderSnapshotRestoreJobs
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs/
type CloudProviderSnapshotRestoreJobsService interface {
	List(context.Context, *SnapshotReqPathParameters, *ListOptions) (*CloudProviderSnapshotRestoreJobs, *Response, error)
	Get(context.Context, *SnapshotReqPathParameters) (*CloudProviderSnapshotRestoreJob, *Response, error)
	Create(context.Context, *SnapshotReqPathParameters, *CloudProviderSnapshotRestoreJob) (*CloudProviderSnapshotRestoreJob, *Response, error)
	Delete(context.Context, *SnapshotReqPathParameters) (*Response, error)
	ListForServerlessBackupRestore(context.Context, string, string, *ListOptions) (*CloudProviderSnapshotRestoreJobs, *Response, error)
	GetForServerlessBackupRestore(context.Context, string, string, string) (*CloudProviderSnapshotRestoreJob, *Response, error)
	CreateForServerlessBackupRestore(context.Context, string, string, *CloudProviderSnapshotRestoreJob) (*CloudProviderSnapshotRestoreJob, *Response, error)
}

// CloudProviderSnapshotRestoreJobsServiceOp handles communication with the CloudProviderSnapshotRestoreJobs related methods of the
// MongoDB Atlas API.
type CloudProviderSnapshotRestoreJobsServiceOp service

var _ CloudProviderSnapshotRestoreJobsService = &CloudProviderSnapshotRestoreJobsServiceOp{}

// CloudProviderSnapshotRestoreJob represents the structure of a cloudProviderSnapshotRestoreJob.
type CloudProviderSnapshotRestoreJob struct {
	ID                    string       `json:"id,omitempty"`                    // The unique identifier of the restore job.
	SnapshotID            string       `json:"snapshotId,omitempty"`            // Unique identifier of the snapshot to restore.
	Components            []*Component `json:"components,omitempty"`            // Collection of clusters to be downloaded. Atlas returns this parameter when restoring a sharded cluster and "deliveryType" : "download".
	DeliveryType          string       `json:"deliveryType,omitempty"`          // Type of restore job to create. Possible values are: automated or download or pointInTime
	DeliveryURL           []string     `json:"deliveryUrl,omitempty"`           // One or more URLs for the compressed snapshot files for manual download. Only visible if deliveryType is download.
	TargetClusterName     string       `json:"targetClusterName,omitempty"`     // Name of the target Atlas cluster to which the restore job restores the snapshot. Only required if deliveryType is automated.
	TargetGroupID         string       `json:"targetGroupId,omitempty"`         // Unique ID of the target Atlas project for the specified targetClusterName. Only required if deliveryType is automated.
	Cancelled             bool         `json:"cancelled,omitempty"`             // Indicates whether the restore job was canceled.
	CreatedAt             string       `json:"createdAt,omitempty"`             // UTC ISO 8601 formatted point in time when Atlas created the restore job.
	Expired               bool         `json:"expired,omitempty"`               // Indicates whether the restore job expired.
	ExpiresAt             string       `json:"expiresAt,omitempty"`             // UTC ISO 8601 formatted point in time when the restore job expires.
	FinishedAt            string       `json:"finishedAt,omitempty"`            // UTC ISO 8601 formatted point in time when the restore job completed.
	Links                 []*Link      `json:"links,omitempty"`                 // One or more links to sub-resources and/or related resources. The relations between URLs are explained in the Web Linking Specification.
	Timestamp             string       `json:"timestamp,omitempty"`             // Timestamp in ISO 8601 date and time format in UTC when the snapshot associated to snapshotId was taken.
	OplogTs               int64        `json:"oplogTs,omitempty"`               //nolint:stylecheck // not changing this // Timestamp in the number of seconds that have elapsed since the UNIX epoch from which to you want to restore this snapshot. This is the first part of an Oplog timestamp.
	OplogInc              int64        `json:"oplogInc,omitempty"`              // Oplog operation number from which to you want to restore this snapshot. This is the second part of an Oplog timestamp.
	PointInTimeUTCSeconds int64        `json:"pointInTimeUTCSeconds,omitempty"` // Timestamp in the number of seconds that have elapsed since the UNIX epoch from which you want to restore this snapshot.
	SourceClusterName     string       `json:"sourceClusterName,omitempty"`
	Failed                *bool        `json:"failed,omitempty"`
}

// CloudProviderSnapshotRestoreJobs represents an array of cloudProviderSnapshotRestoreJob.
type CloudProviderSnapshotRestoreJobs struct {
	Links      []*Link                            `json:"links"`
	Results    []*CloudProviderSnapshotRestoreJob `json:"results"`
	TotalCount int                                `json:"totalCount"`
}

type Component struct {
	DownloadURL    string `json:"downloadUrl"`    // DownloadURL from which the snapshot of the components.replicaSetName should be downloaded. Atlas returns null for this parameter if the download URL has expired, has been used, or hasn't been created.
	ReplicaSetName string `json:"replicaSetName"` // ReplicaSetName of the shard or config server included in the snapshot.
}

// List gets all cloud provider snapshot restore jobs for the specified cluster.
//
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
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-restore-jobs-get-one/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) Get(ctx context.Context, requestParameters *SnapshotReqPathParameters) (*CloudProviderSnapshotRestoreJob, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if requestParameters.JobID == "" {
		return nil, nil, NewArgError("JobID", "must be set")
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
//
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
//
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

// ListForServerlessBackupRestore gets all cloud provider snapshot serverless restore jobs for the specified cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/restore/return-all-restore-jobs-for-one-serverless-instance/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) ListForServerlessBackupRestore(ctx context.Context, projectID, instanceName string, listOptions *ListOptions) (*CloudProviderSnapshotRestoreJobs, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}

	path := fmt.Sprintf("%s/%s/serverless/%s/backup/restoreJobs", cloudProviderSnapshotsBasePath, projectID, instanceName)
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

// GetForServerlessBackupRestore gets one cloud provider serverless snapshot restore jobs for the specified cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/restore/return-one-restore-job-for-one-serverless-instance/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) GetForServerlessBackupRestore(ctx context.Context, projectID, instanceName, jobID string) (*CloudProviderSnapshotRestoreJob, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}
	if jobID == "" {
		return nil, nil, NewArgError("jobID", "must be set")
	}

	path := fmt.Sprintf("%s/%s/serverless/%s/backup/restoreJobs/%s", cloudProviderSnapshotsBasePath, projectID, instanceName, jobID)

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

// CreateForServerlessBackupRestore creates a new restore job from a serverless cloud provider snapshot associated to the specified cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/restore/restore-one-snapshot-of-one-serverless-instance/
func (s *CloudProviderSnapshotRestoreJobsServiceOp) CreateForServerlessBackupRestore(ctx context.Context, projectID, instanceName string, createRequest *CloudProviderSnapshotRestoreJob) (*CloudProviderSnapshotRestoreJob, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if instanceName == "" {
		return nil, nil, NewArgError("instanceName", "must be set")
	}

	if createRequest.DeliveryType == "download" {
		createRequest.TargetClusterName = ""
		createRequest.TargetGroupID = ""
	}

	path := fmt.Sprintf("%s/%s/serverless/%s/backup/restoreJobs", cloudProviderSnapshotRestoreJobBasePath, projectID, instanceName)

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
