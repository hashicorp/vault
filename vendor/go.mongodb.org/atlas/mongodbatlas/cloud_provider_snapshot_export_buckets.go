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
	cloudProviderSnapshotExportBucketsPath = "api/atlas/v1.0/groups/%s/backup/exportBuckets"
)

// CloudProviderSnapshotExportBucketsService is an interface for interfacing with the Cloud Provider Snapshots Export Buckets
// endpoints of the MongoDB Atlas API.
type CloudProviderSnapshotExportBucketsService interface {
	List(context.Context, string, *ListOptions) (*CloudProviderSnapshotExportBuckets, *Response, error)
	Get(context.Context, string, string) (*CloudProviderSnapshotExportBucket, *Response, error)
	Create(context.Context, string, *CloudProviderSnapshotExportBucket) (*CloudProviderSnapshotExportBucket, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// CloudProviderSnapshotExportBucketsServiceOp handles communication with the CloudProviderSnapshotExportBucketsService related methods of the
// MongoDB Atlas API.
type CloudProviderSnapshotExportBucketsServiceOp service

// List Retrieve all the buckets for the specified project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/get-all-export-buckets/
func (c CloudProviderSnapshotExportBucketsServiceOp) List(ctx context.Context, projectID string, options *ListOptions) (*CloudProviderSnapshotExportBuckets, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	path := fmt.Sprintf(cloudProviderSnapshotExportBucketsPath, projectID)

	path, err := setListOptions(path, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotExportBuckets)
	resp, err := c.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get Allows you to retrieve one bucket specified by the bucket ID.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/get-one-export-bucket/
func (c CloudProviderSnapshotExportBucketsServiceOp) Get(ctx context.Context, projectID, bucketID string) (*CloudProviderSnapshotExportBucket, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if bucketID == "" {
		return nil, nil, NewArgError("bucketID", "must be set")
	}

	path := fmt.Sprintf("api/atlas/v1.0/groups/%s/backup/exportBuckets/%s", projectID, bucketID)

	req, err := c.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotExportBucket)
	resp, err := c.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create Allows you to grant Atlas access to the specified bucket for exporting backup snapshots.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/create-one-export-bucket/
func (c CloudProviderSnapshotExportBucketsServiceOp) Create(ctx context.Context, projectID string, bucket *CloudProviderSnapshotExportBucket) (*CloudProviderSnapshotExportBucket, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	path := fmt.Sprintf(cloudProviderSnapshotExportBucketsPath, projectID)

	req, err := c.Client.NewRequest(ctx, http.MethodPost, path, bucket)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshotExportBucket)
	resp, err := c.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete Allows you to remove one bucket specified by the bucket ID.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/export/delete-one-export-bucket/
func (c CloudProviderSnapshotExportBucketsServiceOp) Delete(ctx context.Context, projectID, bucketID string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}
	if bucketID == "" {
		return nil, NewArgError("bucketID", "must be set")
	}

	path := fmt.Sprintf("api/atlas/v1.0/groups/%s/backup/exportBuckets/%s", projectID, bucketID)

	req, err := c.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(ctx, req, nil)

	return resp, err
}

var _ CloudProviderSnapshotExportBucketsService = &CloudProviderSnapshotExportBucketsServiceOp{}

// CloudProviderSnapshotExportBuckets represents all cloud provider snapshot export buckets.
type CloudProviderSnapshotExportBuckets struct {
	Results    []*CloudProviderSnapshotExportBucket `json:"results,omitempty"`    // Includes one CloudProviderSnapshotExportBucket object for each item detailed in the results array section.
	Links      []*Link                              `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int                                  `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

// CloudProviderSnapshotExportBucket represents one cloud provider snapshot export buckets.
type CloudProviderSnapshotExportBucket struct {
	ID            string `json:"_id,omitempty"`           // Unique identifier of the S3 bucket.
	BucketName    string `json:"bucketName,omitempty"`    // Name of the bucket that the role ID is authorized to access.
	CloudProvider string `json:"cloudProvider,omitempty"` // Name of the provider of the cloud service where Atlas can access the S3 bucket. Atlas only supports AWS.
	IAMRoleID     string `json:"iamRoleId,omitempty"`     // Unique identifier of the role that Atlas can use to access the bucket. If necessary, use the UI or API to retrieve the role ID. You must also specify the bucketName.
}
