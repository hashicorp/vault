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

const continuousRestoreJobsPath = "api/atlas/v1.0/groups/%s/clusters/%s/restoreJobs"

// ContinuousRestoreJobsService provides access to the restore jobs related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/legacy-backup/restore/restores/
type ContinuousRestoreJobsService interface {
	List(context.Context, string, string, *ListOptions) (*ContinuousJobs, *Response, error)
	Get(context.Context, string, string, string) (*ContinuousJob, *Response, error)
	Create(context.Context, string, string, *ContinuousJobRequest) (*ContinuousJobs, *Response, error)
}

// ContinuousRestoreJobsServiceOp handles communication with the Continuous Backup Restore Jobs related methods
// of the MongoDB Atlas API.
type ContinuousRestoreJobsServiceOp service

var _ ContinuousRestoreJobsService = &ContinuousRestoreJobsServiceOp{}

type ContinuousJob struct {
	BatchID           string            `json:"batchId,omitempty"`
	ClusterID         string            `json:"clusterId,omitempty"`
	Created           string            `json:"created"`
	ClusterName       string            `json:"clusterName,omitempty"`
	Delivery          *Delivery         `json:"delivery,omitempty"`
	EncryptionEnabled bool              `json:"encryptionEnabled"`
	GroupID           string            `json:"groupId"`
	Hashes            []*Hash           `json:"hashes,omitempty"`
	ID                string            `json:"id"`
	Links             []*Link           `json:"links,omitempty"`
	MasterKeyUUID     string            `json:"masterKeyUUID,omitempty"`
	SnapshotID        string            `json:"snapshotId"`
	StatusName        string            `json:"statusName"`
	PointInTime       *bool             `json:"pointInTime,omitempty"`
	Timestamp         SnapshotTimestamp `json:"timestamp"`
}

type ContinuousJobs struct {
	Results    []*ContinuousJob `json:"results,omitempty"`
	Links      []*Link          `json:"links,omitempty"`
	TotalCount int64            `json:"totalCount,omitempty"`
}

type SnapshotTimestamp struct {
	Date      string `json:"date"`
	Increment int64  `json:"increment"`
}

type Delivery struct {
	Expires           string `json:"expires,omitempty"`
	ExpirationHours   int64  `json:"expirationHours,omitempty"`
	MaxDownloads      int64  `json:"maxDownloads,omitempty"`
	MethodName        string `json:"methodName"`
	StatusName        string `json:"statusName,omitempty"`
	URL               string `json:"url,omitempty"`
	TargetClusterID   string `json:"targetClusterId,omitempty"`
	TargetClusterName string `json:"targetClusterName,omitempty"`
	TargetGroupID     string `json:"targetGroupId,omitempty"`
}

type Hash struct {
	TypeName string `json:"typeName"`
	FileName string `json:"fileName"`
	Hash     string `json:"hash"`
}

type ContinuousJobRequest struct {
	CheckPointID         string   `json:"checkPointId,omitempty"`
	Delivery             Delivery `json:"delivery"`
	OplogTS              string   `json:"oplogTs,omitempty"`
	OplogInc             int64    `json:"oplogInc,omitempty"`
	PointInTimeUTCMillis float64  `json:"pointInTimeUTCMillis,omitempty"`
	SnapshotID           string   `json:"snapshotId,omitempty"`
}

// List lists all continuous backup jobs in Atlas
//
// See more: https://docs.atlas.mongodb.com/reference/api/restore-jobs-get-all/
func (s *ContinuousRestoreJobsServiceOp) List(ctx context.Context, groupID, clusterID string, opts *ListOptions) (*ContinuousJobs, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(continuousRestoreJobsPath, groupID, clusterID)

	path, err := setListOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousJobs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets a continuous backup job in Atlas
//
// See more: https://docs.atlas.mongodb.com/reference/api/restore-jobs-get-one/
func (s *ContinuousRestoreJobsServiceOp) Get(ctx context.Context, groupID, clusterID, jobID string) (*ContinuousJob, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if jobID == "" {
		return nil, nil, NewArgError("jobID", "must be set")
	}
	defaultPath := fmt.Sprintf(continuousRestoreJobsPath, groupID, clusterID)

	path := fmt.Sprintf("%s/%s", defaultPath, jobID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousJob)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create creates a continuous backup job in Atlas
//
// See more: https://docs.atlas.mongodb.com/reference/api/restore-jobs-create-one/
func (s *ContinuousRestoreJobsServiceOp) Create(ctx context.Context, groupID, clusterID string, request *ContinuousJobRequest) (*ContinuousJobs, *Response, error) {
	if request == nil {
		return nil, nil, NewArgError("request", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(continuousRestoreJobsPath, groupID, clusterID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, request)

	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousJobs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
