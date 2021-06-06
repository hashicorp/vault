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
	onlineArchiveBasePath = "groups/%s/clusters/%s/onlineArchives"
)

// OnlineArchiveService provides access to the online archive related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/online-archive/
type OnlineArchiveService interface {
	List(context.Context, string, string, *ListOptions) (*OnlineArchives, *Response, error)
	Get(context.Context, string, string, string) (*OnlineArchive, *Response, error)
	Create(context.Context, string, string, *OnlineArchive) (*OnlineArchive, *Response, error)
	Update(context.Context, string, string, string, *OnlineArchive) (*OnlineArchive, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
}

// OnlineArchiveServiceOp provides an implementation of the OnlineArchiveService interface
type OnlineArchiveServiceOp service

var _ OnlineArchiveService = &OnlineArchiveServiceOp{}

// List gets all online archives.
//
// See more: https://docs.atlas.mongodb.com/reference/api/online-archive-get-all-for-cluster/#api-online-archive-get-all-for-clstr
func (s *OnlineArchiveServiceOp) List(ctx context.Context, projectID, clusterName string, listOptions *ListOptions) (*OnlineArchives, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(onlineArchiveBasePath, projectID, clusterName)
	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root *OnlineArchives
	resp, err := s.Client.Do(ctx, req, &root)
	return root, resp, err
}

// Get gets a single online archive.
//
// See more: https://docs.atlas.mongodb.com/reference/api/online-archive-get-one/
func (s *OnlineArchiveServiceOp) Get(ctx context.Context, projectID, clusterName, archiveID string) (*OnlineArchive, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if archiveID == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(onlineArchiveBasePath, projectID, clusterName)
	path = fmt.Sprintf("%s/%s", path, archiveID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(OnlineArchive)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create creates a new online archive.
//
// See more: https://docs.atlas.mongodb.com/reference/api/online-archive-create-one/
func (s *OnlineArchiveServiceOp) Create(ctx context.Context, projectID, clusterName string, r *OnlineArchive) (*OnlineArchive, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(onlineArchiveBasePath, projectID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, r)
	if err != nil {
		return nil, nil, err
	}

	root := new(OnlineArchive)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update let's you pause or resume archiving for an online archive or modify the archiving criteria.
//
// See more: https://docs.atlas.mongodb.com/reference/api/online-archive-update-one/
func (s *OnlineArchiveServiceOp) Update(ctx context.Context, projectID, clusterName, archiveID string, r *OnlineArchive) (*OnlineArchive, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if archiveID == "" {
		return nil, nil, NewArgError("archiveID", "must be set")
	}

	path := fmt.Sprintf(onlineArchiveBasePath, projectID, clusterName)
	path = fmt.Sprintf("%s/%s", path, archiveID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, r)
	if err != nil {
		return nil, nil, err
	}

	root := new(OnlineArchive)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete deletes an online archive.
//
// See more: https://docs.atlas.mongodb.com/reference/api/online-archive-delete-one/
func (s *OnlineArchiveServiceOp) Delete(ctx context.Context, projectID, clusterName, archiveID string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}
	if archiveID == "" {
		return nil, NewArgError("archiveID", "must be set")
	}

	path := fmt.Sprintf(onlineArchiveBasePath, projectID, clusterName)
	path = fmt.Sprintf("%s/%s", path, archiveID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// OnlineArchives is a collection of OnlineArchive
type OnlineArchives struct {
	Links      []*Link          `json:"links,omitempty"`
	Results    []*OnlineArchive `json:"results,omitempty"`
	TotalCount int              `json:"totalCount,omitempty"`
}

// OnlineArchive represents the structure of an online archive.
type OnlineArchive struct {
	ID              string                 `json:"_id,omitempty"`
	ClusterName     string                 `json:"clusterName,omitempty"`
	CollName        string                 `json:"collName,omitempty"`
	Criteria        *OnlineArchiveCriteria `json:"criteria,omitempty"`
	DBName          string                 `json:"dbName,omitempty"`
	GroupID         string                 `json:"groupId,omitempty"`
	PartitionFields []*PartitionFields     `json:"partitionFields,omitempty"`
	Paused          *bool                  `json:"paused,omitempty"`
	State           string                 `json:"state,omitempty"`
}

// OnlineArchiveCriteria criteria to use for archiving data.
type OnlineArchiveCriteria struct {
	DateField       string  `json:"dateField,omitempty"`
	DateFormat      string  `json:"dateFormat,omitempty"`
	ExpireAfterDays float64 `json:"expireAfterDays"`
	Type            string  `json:"type,omitempty"`
}

// PartitionFields fields to use to partition data
type PartitionFields struct {
	FieldName string   `json:"fieldName,omitempty"`
	FieldType string   `json:"fieldType,omitempty"`
	Order     *float64 `json:"order,omitempty"`
}
