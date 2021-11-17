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
	backupCheckpoints = "api/atlas/v1.0/groups/%s/clusters/%s/backupCheckpoints"
)

// CheckpointsService is an interface for interfacing with the Checkpoint
// endpoints of the MongoDB Atlas API.
type CheckpointsService interface {
	List(context.Context, string, string, *ListOptions) (*Checkpoints, *Response, error)
	Get(context.Context, string, string, string) (*Checkpoint, *Response, error)
}

// CheckpointsServiceOp handles communication with the checkpoint related methods of the
// MongoDB Atlas API.
type CheckpointsServiceOp service

var _ CheckpointsService = &CheckpointsServiceOp{}

// Checkpoint represents MongoDB Checkpoint.
type Checkpoint struct {
	ClusterID  string  `json:"clusterId"`
	Completed  string  `json:"completed,omitempty"`
	GroupID    string  `json:"groupId"`
	ID         string  `json:"id,omitempty"`    // Unique identifier of the checkpoint.
	Links      []*Link `json:"links,omitempty"` // One or more links to sub-resources and/or related resources.
	Parts      []*Part `json:"parts,omitempty"`
	Restorable bool    `json:"restorable"`
	Started    string  `json:"started"`
	Timestamp  string  `json:"timestamp"`
}

// CheckpointPart represents the individual parts that comprise the complete checkpoint.
type CheckpointPart struct {
	ShardName       string            `json:"shardName"`
	TokenDiscovered bool              `json:"tokenDiscovered"`
	TokenTimestamp  SnapshotTimestamp `json:"tokenTimestamp"`
}

// Checkpoints represents all the backup checkpoints related to a cluster.
type Checkpoints struct {
	Results    []*Checkpoint `json:"results,omitempty"`    // Includes one Checkpoint object for each item detailed in the results array section.
	Links      []*Link       `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int           `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

// List all checkpoints for the specified sharded cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/checkpoints-get-all/
func (s CheckpointsServiceOp) List(ctx context.Context, groupID, clusterName string, listOptions *ListOptions) (*Checkpoints, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(backupCheckpoints, groupID, clusterName)
	path, err := setListOptions(basePath, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Checkpoints)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get one checkpoint for the specified sharded cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/checkpoints-get-one/
func (s CheckpointsServiceOp) Get(ctx context.Context, groupID, clusterName, checkpointID string) (*Checkpoint, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if checkpointID == "" {
		return nil, nil, NewArgError("checkpointID", "must be set")
	}

	basePath := fmt.Sprintf(backupCheckpoints, groupID, clusterName)
	path := fmt.Sprintf("%s/%s", basePath, checkpointID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(Checkpoint)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
