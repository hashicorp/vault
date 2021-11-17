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
	continuousSnapshotsBasePath = "api/atlas/v1.0/groups/%s/clusters/%s/snapshots"
)

// ContinuousSnapshotsService is an interface for interfacing with the Continuous Snapshots
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/snapshots/
type ContinuousSnapshotsService interface {
	List(context.Context, string, string, *ListOptions) (*ContinuousSnapshots, *Response, error)
	Get(context.Context, string, string, string) (*ContinuousSnapshot, *Response, error)
	ChangeExpiry(context.Context, string, string, string, *ContinuousSnapshot) (*ContinuousSnapshot, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
}

// ContinuousSnapshotsServiceOp handles communication with the Continuous Snapshots related methods of the
// MongoDB Atlas API.
type ContinuousSnapshotsServiceOp service

var _ ContinuousSnapshotsService = &ContinuousSnapshotsServiceOp{}

// ContinuousSnapshot represents a cloud provider snapshot.
type ContinuousSnapshot struct {
	ClusterID                 string             `json:"clusterId,omitempty"`
	Complete                  bool               `json:"complete,omitempty"`
	Created                   *SnapshotTimestamp `json:"created,omitempty"`
	DoNotDelete               *bool              `json:"doNotDelete,omitempty"`
	Expires                   string             `json:"expires,omitempty"`
	GroupID                   string             `json:"groupId,omitempty"`
	ID                        string             `json:"id,omitempty"` // Unique identifier of the snapshot.
	IsPossiblyInconsistent    *bool              `json:"isPossiblyInconsistent,omitempty"`
	LastOplogAppliedTimestamp *SnapshotTimestamp `json:"lastOplogAppliedTimestamp,omitempty"`
	Links                     []*Link            `json:"links,omitempty"` // One or more links to sub-resources and/or related resources.
	NamespaceFilterList       *NamespaceFilter   `json:"namespaceFilterList,omitempty"`
	MissingShards             []*MissingShard    `json:"missingShards,omitempty"`
	Parts                     []*Part            `json:"parts,omitempty"`
}

type Part struct {
	ReplicaSetName string `json:"replicaSetName"`
	TypeName       string `json:"typeName"`
	SnapshotPart
	CheckpointPart
}

type SnapshotPart struct {
	ClusterID          string `json:"clusterId"`
	CompressionSetting string `json:"compressionSetting"`
	DataSizeBytes      int64  `json:"dataSizeBytes"`
	EncryptionEnabled  bool   `json:"encryptionEnabled"`
	FileSizeBytes      int64  `json:"fileSizeBytes"`
	MasterKeyUUID      string `json:"masterKeyUUID,omitempty"`
	MongodVersion      string `json:"mongodVersion"`
	StorageSizeBytes   int64  `json:"storageSizeBytes"`
}

type NamespaceFilter struct {
	FilterList []string `json:"filterList"`
	FilterType string   `json:"filterType"`
}

type MissingShard struct {
	ID             string `json:"id"`
	GroupID        string `json:"groupId"`
	TypeName       string `json:"typeName"`
	ClusterName    string `json:"clusterName,omitempty"`
	ShardName      string `json:"shardName,omitempty"`
	ReplicaSetName string `json:"replicaSetName"`
	LastHeartbeat  string `json:"lastHeartbeat"`
}

// ContinuousSnapshots represents all cloud provider snapshots.
type ContinuousSnapshots struct {
	Results    []*ContinuousSnapshot `json:"results,omitempty"`    // Includes one ContinuousSnapshots object for each item detailed in the results array section.
	Links      []*Link               `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int                   `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

// List lists continuous snapshots for the given cluster
//
// See more: https://docs.atlas.mongodb.com/reference/api/snapshots-get-all/
func (s *ContinuousSnapshotsServiceOp) List(ctx context.Context, groupID, clusterID string, listOptions *ListOptions) (*ContinuousSnapshots, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}

	path := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)

	// Add query params
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousSnapshots)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets the continuous snapshot for the given cluster and snapshot ID
//
// See more: https://docs.atlas.mongodb.com/reference/api/snapshots-get-one/
func (s *ContinuousSnapshotsServiceOp) Get(ctx context.Context, groupID, clusterID, snapshotID string) (*ContinuousSnapshot, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if snapshotID == "" {
		return nil, nil, NewArgError("snapshotID", "must be set")
	}

	basePath := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, snapshotID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousSnapshot)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ChangeExpiry changes the expiry date for the given cluster and snapshot ID
//
// See more: https://docs.atlas.mongodb.com/reference/api/snapshots-change-expiration/
func (s *ContinuousSnapshotsServiceOp) ChangeExpiry(ctx context.Context, groupID, clusterID, snapshotID string, updateRequest *ContinuousSnapshot) (*ContinuousSnapshot, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if snapshotID == "" {
		return nil, nil, NewArgError("snapshotID", "must be set")
	}

	basePath := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, snapshotID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousSnapshot)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete deletes the given continuous snapshot
//
// See more: https://docs.atlas.mongodb.com/reference/api/snapshots-delete-one/
func (s *ContinuousSnapshotsServiceOp) Delete(ctx context.Context, groupID, clusterID, snapshotID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, NewArgError("clusterID", "must be set")
	}
	if snapshotID == "" {
		return nil, NewArgError("snapshotID", "must be set")
	}

	basePath := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, snapshotID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}
