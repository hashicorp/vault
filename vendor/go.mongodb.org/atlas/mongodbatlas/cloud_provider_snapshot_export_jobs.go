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

const (
	cloudProviderSnapshotExportJobsPath = "api/atlas/v1.0/groups/%s/clusters/%s/backup/exports"
)

// CloudProviderSnapshotExportJobsService is an interface for interfacing with the Cloud Provider Snapshots Export Jobs
// of the MongoDB Atlas API.
type CloudProviderSnapshotExportJobsService interface {
	List(context.Context, string, string, *ListOptions) (*CloudProviderSnapshotExportJobs, *Response, error)
	Get(context.Context, string, string, string) (*CloudProviderSnapshotExportJob, *Response, error)
	Create(context.Context, string, string, *CloudProviderSnapshotExportJob) (*CloudProviderSnapshotExportJob, *Response, error)
}

// CloudProviderSnapshotExportJobsServiceOp handles communication with the CloudProviderSnapshotExportJobsService related methods of the
// MongoDB Atlas API.
type CloudProviderSnapshotExportJobsServiceOp service

// List Retrieve all the export jobs for the specified project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/get-all-export-jobs/
func (c CloudProviderSnapshotExportJobsServiceOp) List(ctx context.Context, projectID, clusterName string, options *ListOptions) (*CloudProviderSnapshotExportJobs, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(cloudProviderSnapshotExportJobsPath, projectID, clusterName)

	path, err := setListOptions(path, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotExportJobs)
	resp, err := c.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get Allows you to retrieve one export job specified by the export job ID.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/get-one-export-job/
func (c CloudProviderSnapshotExportJobsServiceOp) Get(ctx context.Context, projectID, clusterName, exportID string) (*CloudProviderSnapshotExportJob, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if exportID == "" {
		return nil, nil, NewArgError("exportID", "must be set")
	}

	path := fmt.Sprintf("api/atlas/v1.0/groups/%s/clusters/%s/backup/exports/%s", projectID, clusterName, exportID)

	req, err := c.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotExportJob)
	resp, err := c.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create Allows you to grant Atlas access to the specified export job for exporting backup snapshots.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/create-one-export-job/
func (c CloudProviderSnapshotExportJobsServiceOp) Create(ctx context.Context, projectID, clusterName string, bucket *CloudProviderSnapshotExportJob) (*CloudProviderSnapshotExportJob, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(cloudProviderSnapshotExportJobsPath, projectID, clusterName)

	req, err := c.Client.NewRequest(ctx, http.MethodPost, path, bucket)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotExportJob)
	resp, err := c.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

var _ CloudProviderSnapshotExportJobsService = &CloudProviderSnapshotExportJobsServiceOp{}

// CloudProviderSnapshotExportJobs represents all cloud provider snapshot export jobs.
type CloudProviderSnapshotExportJobs struct {
	Results    []*CloudProviderSnapshotExportJob `json:"results,omitempty"`    // Includes one CloudProviderSnapshotExportJob object for each item detailed in the results array section.
	Links      []*Link                           `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int                               `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

type CloudProviderSnapshotExportJobComponent struct {
	ExportID       string `json:"exportId,omitempty"`       // Returned for sharded clusters only. Unique identifier of the export job for the replica set.
	ReplicaSetName string `json:"replicaSetName,omitempty"` // Returned for sharded clusters only. Name of the replica set.
}

type CloudProviderSnapshotExportJobCustomData struct {
	Key   string `json:"key,omitempty"`   // Custom data specified as key in the key and value pair.
	Value string `json:"value,omitempty"` // Value for the key specified using CloudProviderSnapshotExportJobCustomData.key.
}

type CloudProviderSnapshotExportJobStatus struct {
	ExportedCollections int `json:"exportedCollections,omitempty"` // Returned for replica set only. Number of collections that have been exported.
	TotalCollections    int `json:"totalCollections,omitempty"`    // Returned for replica set only. Total number of collections to export.
}

// CloudProviderSnapshotExportJob represents one cloud provider snapshot export jobs.
type CloudProviderSnapshotExportJob struct {
	ID             string                                      `json:"id,omitempty"`             // Unique identifier of the export job.
	Components     []*CloudProviderSnapshotExportJobComponent  `json:"components,omitempty"`     // Returned for sharded clusters only. Export job details for each replica set in the sharded cluster.
	CreatedAt      string                                      `json:"createdAt,omitempty"`      // Timestamp in ISO 8601 date and time format in UTC when the export job was created.
	CustomData     []*CloudProviderSnapshotExportJobCustomData `json:"customData,omitempty"`     // Custom data for the metadata file named .complete that Atlas uploads to the bucket when the export job finishes.
	ErrMsg         string                                      `json:"errMsg,omitempty"`         // Error message, only if the export job failed.
	ExportBucketID string                                      `json:"exportBucketId,omitempty"` // Unique identifier of the bucket.
	ExportStatus   *CloudProviderSnapshotExportJobStatus       `json:"exportStatus,omitempty"`   // Returned for replica set only. Status of the export job.
	FinishedAt     string                                      `json:"finishedAt,omitempty"`     // Timestamp in ISO 8601 date and time format in UTC when the export job completes.
	Prefix         string                                      `json:"prefix,omitempty"`         // Full path on the cloud provider bucket to the folder where the snapshot is exported. The path is in the following format: /exported_snapshots/{ORG-NAME}/{PROJECT-NAME}/{CLUSTER-NAME}/{SNAPSHOT-INITIATION-DATE}/{TIMESTAMP}
	SnapshotID     string                                      `json:"snapshotId,omitempty"`     // Unique identifier of the snapshot.
	State          string                                      `json:"state,omitempty"`          // Status of the export job.
}
