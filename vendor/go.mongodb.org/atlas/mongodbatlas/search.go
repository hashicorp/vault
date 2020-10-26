package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	searchBasePath = "groups/%s/clusters/%s/fts"
)

// SearchService provides access to the search related functions in the Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/atlas-search/
type SearchService interface {
	ListIndexes(context.Context, string, string, string, string, *ListOptions) ([]*SearchIndex, *Response, error)
	GetIndex(context.Context, string, string, string) (*SearchIndex, *Response, error)
	CreateIndex(context.Context, string, string, *SearchIndex) (*SearchIndex, *Response, error)
	UpdateIndex(context.Context, string, string, string, *SearchIndex) (*SearchIndex, *Response, error)
	DeleteIndex(context.Context, string, string, string) (*Response, error)
	ListAnalyzers(context.Context, string, string, *ListOptions) ([]*SearchAnalyzer, *Response, error)
}

// SearchServiceOp provides an implementation of the SearchService interface
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

// SearchIndex index definition.
type SearchIndex struct {
	Analyzer       string        `json:"analyzer,omitempty"`
	CollectionName string        `json:"collectionName"`
	Database       string        `json:"database"`
	IndexID        string        `json:"indexID,omitempty"`
	Mappings       *IndexMapping `json:"mappings,omitempty"`
	Name           string        `json:"name"`
	SearchAnalyzer string        `json:"searchAnalyzer,omitempty"`
}

// IndexMapping containing index specifications for the collection fields.
type IndexMapping struct {
	Dynamic bool                   `json:"dynamic"`
	Fields  *map[string]IndexField `json:"fields,omitempty"`
}

// IndexField field specifications.
type IndexField struct {
	Analyzer       string                 `json:"analyzer,omitempty"`
	Type           string                 `json:"type"`
	Tokenization   string                 `json:"tokenization,omitempty"` // edgeGram|nGram
	MinGrams       *int                   `json:"minGrams,omitempty"`
	MaxGrams       *int                   `json:"maxGrams,omitempty"`
	FoldDiacritics *bool                  `json:"foldDiacritics,omitempty"`
	Fields         *map[string]IndexField `json:"fields,omitempty"`
	SearchAnalyzer string                 `json:"searchAnalyzer,omitempty"`
	IndexOptions   string                 `json:"indexOptions,omitempty"` // docs|freqs|positions
	Store          *bool                  `json:"store,omitempty"`
	IgnoreAbove    *int                   `json:"ignoreAbove,omitempty"`
	Norms          string                 `json:"norms,omitempty"` // include|omit
	Dynamic        *bool                  `json:"dynamic,omitempty"`
	Representation string                 `json:"representation,omitempty"`
	IndexIntegers  *bool                  `json:"indexIntegers,omitempty"`
	IndexDoubles   *bool                  `json:"indexDoubles,omitempty"`
	IndexShapes    *bool                  `json:"indexShapes,omitempty"`
}

// SearchAnalyzer custom analyzer definition.
type SearchAnalyzer struct {
	BaseAnalyzer     string   `json:"baseAnalyzer"`
	MaxTokenLength   *float64 `json:"maxTokenLength,omitempty"`
	Name             string   `json:"name"`
	StemExclusionSet []string `json:"stemExclusionSet,omitempty"`
	Stopwords        []string `json:"stopwords,omitempty"`
}
