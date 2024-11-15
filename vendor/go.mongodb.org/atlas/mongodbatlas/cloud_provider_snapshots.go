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
	cloudProviderSnapshotsBasePath = "api/atlas/v1.0/groups"
)

// CloudProviderSnapshotsService is an interface for interfacing with the Cloud Provider Snapshots
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot/
type CloudProviderSnapshotsService interface {
	GetAllCloudProviderSnapshots(context.Context, *SnapshotReqPathParameters, *ListOptions) (*CloudProviderSnapshots, *Response, error)
	GetOneCloudProviderSnapshot(context.Context, *SnapshotReqPathParameters) (*CloudProviderSnapshot, *Response, error)
	Create(context.Context, *SnapshotReqPathParameters, *CloudProviderSnapshot) (*CloudProviderSnapshot, *Response, error)
	Delete(context.Context, *SnapshotReqPathParameters) (*Response, error)
	GetOneServerlessSnapshot(context.Context, *SnapshotReqPathParameters) (*CloudProviderSnapshot, *Response, error)
	GetAllServerlessSnapshots(context.Context, *SnapshotReqPathParameters, *ListOptions) (*CloudProviderSnapshots, *Response, error)
}

// CloudProviderSnapshotsServiceOp handles communication with the CloudProviderSnapshotsService related methods of the
// MongoDB Atlas API.
type CloudProviderSnapshotsServiceOp service

var _ CloudProviderSnapshotsService = &CloudProviderSnapshotsServiceOp{}

// CloudProviderSnapshot represents a cloud provider snapshot.
type CloudProviderSnapshot struct {
	ID                     string    `json:"id,omitempty"`                     // Unique identifier of the snapshot.
	RetentionInDays        int       `json:"retentionInDays,omitempty"`        // The number of days that Atlas should retain the on-demand snapshot. Must be at least 1 .
	CreatedAt              string    `json:"createdAt,omitempty"`              // UTC ISO 8601 formatted point in time when Atlas took the snapshot.
	ExpiresAt              string    `json:"expiresAt,omitempty"`              // UTC ISO 8601 formatted point in time when Atlas will delete the snapshot.
	Description            string    `json:"description,omitempty"`            // Description of the on-demand snapshot.
	Links                  []*Link   `json:"links,omitempty"`                  // One or more links to sub-resources and/or related resources.
	MasterKeyUUID          string    `json:"masterKeyUUID,omitempty"`          // Unique ID of the AWS KMS Customer Master Key used to encrypt the snapshot. Only visible for clusters using Encryption at Rest via Customer KMS.
	MongodVersion          string    `json:"mongodVersion,omitempty"`          // Version of the MongoDB server.
	SnapshotType           string    `json:"snapshotType,omitempty"`           // Specified the type of snapshot. Valid values are onDemand and scheduled.
	Status                 string    `json:"status,omitempty"`                 // Current status of the snapshot. One of the following values: queued, inProgress, completed, failed
	StorageSizeBytes       int       `json:"storageSizeBytes,omitempty"`       // Specifies the size of the snapshot in bytes.
	Type                   string    `json:"type,omitempty"`                   // Specifies the type of cluster: replicaSet or shardedCluster.
	CloudProvider          string    `json:"cloudProvider,omitempty"`          // Cloud provider that stores this snapshot. Atlas returns this parameter when "type": "replicaSet".
	Members                []*Member `json:"members,omitempty"`                // List of ReplicaSetName and the cloud provider where the snapshots are stored. Atlas returns this parameter when "type": "shardedCluster".
	ReplicaSetName         string    `json:"replicaSetName,omitempty"`         // Label given to the replica set from which Atlas took this snapshot. Atlas returns this parameter when "type": "replicaSet".
	SnapshotsIds           []string  `json:"snapshotIds,omitempty"`            // Unique identifiers of the snapshots created for the shards and config server for a sharded cluster.
	ServerlessInstanceName string    `json:"serverlessInstanceName,omitempty"` // Human-readable label given to the serverless instance from which Atlas took this snapshot.
}

// CloudProviderSnapshots represents all cloud provider snapshots.
type CloudProviderSnapshots struct {
	Results    []*CloudProviderSnapshot `json:"results,omitempty"`    // Includes one CloudProviderSnapshot object for each item detailed in the results array section.
	Links      []*Link                  `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int                      `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

// SnapshotReqPathParameters represents all the pissible parameters to make the request.
type SnapshotReqPathParameters struct {
	GroupID      string `json:"groupId,omitempty"`        // The unique identifier of the project for the Atlas cluster.
	SnapshotID   string `json:"snapshotId,omitempty"`     // The unique identifier of the snapshot you want to retrieve.
	ClusterName  string `json:"clusterName,omitempty"`    // The name of the Atlas cluster that contains the snapshots you want to retrieve.
	InstanceName string `json:"instanceName,omitempty"`   // Human-readable label that identifies your serverless instance.
	JobID        string `json:"jobId,omitempty"`          // The unique identifier of the restore job to retrieve.
	RestoreJobID string `json:"restore_job_id,omitempty"` // The unique identifier of the restore job id to retrieve.
}

// Member represents all member of cloud provider snapshot.
type Member struct {
	ID             string `json:"id,omitempty"`             // Cloud provider that stores this snapshot.
	CloudProvider  string `json:"cloudProvider,omitempty"`  // Unique identifier for the sharded cluster snapshot.
	ReplicaSetName string `json:"replicaSetName,omitempty"` // Label given to a shard or config server from which Atlas took this snapshot.
}

// GetAllCloudProviderSnapshots gets all cloud provider snapshots for the specified cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-get-all/
func (s *CloudProviderSnapshotsServiceOp) GetAllCloudProviderSnapshots(ctx context.Context, requestParameters *SnapshotReqPathParameters, listOptions *ListOptions) (*CloudProviderSnapshots, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/snapshots", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.ClusterName)

	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshots)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// GetOneCloudProviderSnapshot gets the snapshot associated to {SNAPSHOT-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-get-one/
func (s *CloudProviderSnapshotsServiceOp) GetOneCloudProviderSnapshot(ctx context.Context, requestParameters *SnapshotReqPathParameters) (*CloudProviderSnapshot, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if requestParameters.SnapshotID == "" {
		return nil, nil, NewArgError("snapshotId", "must be set")
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/snapshots/%s", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.ClusterName, requestParameters.SnapshotID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshot)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create takes one on-demand snapshot. Atlas takes on-demand snapshots immediately, unlike scheduled snapshots which occur at regular intervals.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-take-one-ondemand/
func (s *CloudProviderSnapshotsServiceOp) Create(ctx context.Context, requestParameters *SnapshotReqPathParameters, createRequest *CloudProviderSnapshot) (*CloudProviderSnapshot, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/snapshots", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.ClusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshot)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes the snapshot associated to {SNAPSHOT-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-delete-one/
func (s *CloudProviderSnapshotsServiceOp) Delete(ctx context.Context, requestParameters *SnapshotReqPathParameters) (*Response, error) {
	if requestParameters.GroupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.ClusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}
	if requestParameters.SnapshotID == "" {
		return nil, NewArgError("snapshotId", "must be set")
	}

	path := fmt.Sprintf("%s/%s/clusters/%s/backup/snapshots/%s", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.ClusterName, requestParameters.SnapshotID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// GetAllServerlessSnapshots Return all snapshots of one serverless instance from the specified project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cloud-backup/backup/return-all-snapshots-of-one-serverless-instance/
func (s *CloudProviderSnapshotsServiceOp) GetAllServerlessSnapshots(ctx context.Context, requestParameters *SnapshotReqPathParameters, listOptions *ListOptions) (*CloudProviderSnapshots, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.InstanceName == "" {
		return nil, nil, NewArgError("InstanceName", "must be set")
	}

	path := fmt.Sprintf("%s/%s/serverless/%s/backup/snapshots", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.InstanceName)

	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshots)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// GetOneServerlessSnapshot Return one snapshot of one serverless instance from the specified project.
//
// See more: // GetOneCloudProviderSnapshot https://docs.atlas.mongodb.com/reference/api/cloud-backup/backup/return-one-snapshot-of-one-serverless-instance/
func (s *CloudProviderSnapshotsServiceOp) GetOneServerlessSnapshot(ctx context.Context, requestParameters *SnapshotReqPathParameters) (*CloudProviderSnapshot, *Response, error) {
	if requestParameters.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if requestParameters.InstanceName == "" {
		return nil, nil, NewArgError("InstanceName", "must be set")
	}
	if requestParameters.SnapshotID == "" {
		return nil, nil, NewArgError("snapshotId", "must be set")
	}

	path := fmt.Sprintf("%s/%s/serverless/%s/backup/snapshots/%s", cloudProviderSnapshotsBasePath, requestParameters.GroupID, requestParameters.InstanceName, requestParameters.SnapshotID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviderSnapshot)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
