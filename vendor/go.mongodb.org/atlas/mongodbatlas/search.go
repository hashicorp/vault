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
	searchBasePath = "api/atlas/v1.0/groups/%s/clusters/%s/fts"
)

// SearchService provides access to the search related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/atlas-search/
type SearchService interface {
	ListIndexes(ctx context.Context, groupID string, clusterName string, databaseName string, collectionName string, opts *ListOptions) ([]*SearchIndex, *Response, error)
	GetIndex(ctx context.Context, groupID, clusterName, indexID string) (*SearchIndex, *Response, error)
	CreateIndex(ctx context.Context, projectID, clusterName string, r *SearchIndex) (*SearchIndex, *Response, error)
	UpdateIndex(ctx context.Context, projectID, clusterName, indexID string, r *SearchIndex) (*SearchIndex, *Response, error)
	DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) (*Response, error)
	ListAnalyzers(ctx context.Context, groupID, clusterName string, listOptions *ListOptions) ([]*SearchAnalyzer, *Response, error)
	UpdateAllAnalyzers(ctx context.Context, groupID, clusterName string, analyzers []*SearchAnalyzer) ([]*SearchAnalyzer, *Response, error)
}

// SearchServiceOp provides an implementation of the SearchService interface.
type SearchServiceOp service

var _ SearchService = &SearchServiceOp{}

// ListIndexes Get all Atlas Search indexes for a specified collection.
//
// See more: https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/
func (s *SearchServiceOp) ListIndexes(ctx context.Context, groupID, clusterName, databaseName, collectionName string, opts *ListOptions) ([]*SearchIndex, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("GroupID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("ClusterName", "must be set")
	}
	if databaseName == "" {
		return nil, nil, NewArgError("databaseName", "must be set")
	}
	if collectionName == "" {
		return nil, nil, NewArgError("collectionName", "must be set")
	}

	path := fmt.Sprintf(searchBasePath, groupID, clusterName)
	path = fmt.Sprintf("%s/indexes/%s/%s", path, databaseName, collectionName)

	path, err := setListOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*SearchIndex
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

// GetIndex gets one Atlas Search index by its indexId.
//
// See more: https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-one/
func (s *SearchServiceOp) GetIndex(ctx context.Context, groupID, clusterName, indexID string) (*SearchIndex, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if indexID == "" {
		return nil, nil, NewArgError("indexID", "must be set")
	}

	path := fmt.Sprintf(searchBasePath, groupID, clusterName)
	path = fmt.Sprintf("%s/indexes/%s", path, indexID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root *SearchIndex
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

// CreateIndex creates an Atlas Search index.
//
// See more: https://docs.atlas.mongodb.com/reference/api/fts-indexes-create-one/
func (s *SearchServiceOp) CreateIndex(ctx context.Context, projectID, clusterName string, r *SearchIndex) (*SearchIndex, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(searchBasePath, projectID, clusterName)
	path = fmt.Sprintf("%s/indexes", path)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, r)
	if err != nil {
		return nil, nil, err
	}

	var root *SearchIndex
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

// UpdateIndex updates an Atlas Search index by its indexId.
//
// See more: https://docs.atlas.mongodb.com/reference/api/fts-indexes-update-one/
func (s *SearchServiceOp) UpdateIndex(ctx context.Context, projectID, clusterName, indexID string, r *SearchIndex) (*SearchIndex, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}
	if indexID == "" {
		return nil, nil, NewArgError("indexID", "must be set")
	}

	path := fmt.Sprintf(searchBasePath, projectID, clusterName)
	path = fmt.Sprintf("%s/indexes/%s", path, indexID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, r)
	if err != nil {
		return nil, nil, err
	}

	var root *SearchIndex
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

// DeleteIndex deletes one Atlas Search index by its indexId.
//
// See more: https://docs.atlas.mongodb.com/reference/api/fts-indexes-delete-one/
func (s *SearchServiceOp) DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}
	if indexID == "" {
		return nil, NewArgError("indexID", "must be set")
	}

	path := fmt.Sprintf(searchBasePath, projectID, clusterName)
	path = fmt.Sprintf("%s/indexes/%s", path, indexID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// ListAnalyzers gets all Atlas Search user-defined analyzers for a specified cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/fts-analyzers-get-all/
func (s *SearchServiceOp) ListAnalyzers(ctx context.Context, groupID, clusterName string, listOptions *ListOptions) ([]*SearchAnalyzer, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(searchBasePath, groupID, clusterName)
	path = fmt.Sprintf("%s/analyzers", path)

	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*SearchAnalyzer
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

// UpdateAllAnalyzers Update All User-Defined Analyzers for a specific Cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/fts-analyzers-update-all//
func (s *SearchServiceOp) UpdateAllAnalyzers(ctx context.Context, groupID, clusterName string, analyzers []*SearchAnalyzer) ([]*SearchAnalyzer, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(searchBasePath, groupID, clusterName)
	path = fmt.Sprintf("%s/analyzers", path)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, analyzers)
	if err != nil {
		return nil, nil, err
	}

	var root []*SearchAnalyzer
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

// SearchIndex index definition.
type SearchIndex struct {
	Analyzer       string                   `json:"analyzer,omitempty"`
	Analyzers      []map[string]interface{} `json:"analyzers,omitempty"` // Custom analyzers
	CollectionName string                   `json:"collectionName"`
	Database       string                   `json:"database"`
	IndexID        string                   `json:"indexID,omitempty"`
	Mappings       *IndexMapping            `json:"mappings,omitempty"`
	Name           string                   `json:"name"`
	SearchAnalyzer string                   `json:"searchAnalyzer,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Synonyms       []map[string]interface{} `json:"synonyms,omitempty"`
}

// IndexMapping containing index specifications for the collection fields.
type IndexMapping struct {
	Dynamic bool                    `json:"dynamic"`
	Fields  *map[string]interface{} `json:"fields,omitempty"`
}

// SearchAnalyzer search analyzer definition.
type SearchAnalyzer struct {
	BaseAnalyzer     string   `json:"baseAnalyzer"`
	MaxTokenLength   *int     `json:"maxTokenLength,omitempty"`
	IgnoreCase       *bool    `json:"ignoreCase,omitempty"`
	Name             string   `json:"name"`
	StemExclusionSet []string `json:"stemExclusionSet,omitempty"`
	Stopwords        []string `json:"stopwords,omitempty"`
}
