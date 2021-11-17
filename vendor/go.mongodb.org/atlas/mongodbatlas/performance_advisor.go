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
	performanceAdvisorPath                     = "api/atlas/v1.0/groups/%s/processes/%s/performanceAdvisor"
	performanceAdvisorNamespacesPath           = performanceAdvisorPath + "/namespaces"
	performanceAdvisorSlowQueryLogsPath        = performanceAdvisorPath + "/slowQueryLogs"
	performanceAdvisorSuggestedIndexesLogsPath = performanceAdvisorPath + "/suggestedIndexes"
	performanceAdvisorManagedSlowMs            = "api/atlas/v1.0/groups/%s/managedSlowMs"
)

// PerformanceAdvisorService is an interface of the Performance Advisor
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/performance-advisor/
type PerformanceAdvisorService interface {
	GetNamespaces(context.Context, string, string, *NamespaceOptions) (*Namespaces, *Response, error)
	GetSlowQueries(context.Context, string, string, *SlowQueryOptions) (*SlowQueries, *Response, error)
	GetSuggestedIndexes(context.Context, string, string, *SuggestedIndexOptions) (*SuggestedIndexes, *Response, error)
	EnableManagedSlowOperationThreshold(context.Context, string) (*Response, error)
	DisableManagedSlowOperationThreshold(context.Context, string) (*Response, error)
}

// PerformanceAdvisorServiceOp handles communication with the Performance Advisor related methods of the MongoDB Atlas API.
type PerformanceAdvisorServiceOp service

var _ PerformanceAdvisorService = &PerformanceAdvisorServiceOp{}

// Namespace represents a Namespace.
type Namespace struct {
	Namespace string `json:"namespace,omitempty"` // A namespace on the specified host.
	Type      string `json:"type,omitempty"`      // The type of namespace.
}

// Namespaces represents a list of Namespace.
type Namespaces struct {
	Namespaces []*Namespace `json:"namespaces,omitempty"`
}

// SlowQuery represents a slow query.
type SlowQuery struct {
	Namespace string `json:"namespace,omitempty"` // The namespace in which the slow query ran.
	Line      string `json:"line,omitempty"`      // The raw log line pertaining to the slow query.
}

// SlowQueries represents a list of SlowQuery.
type SlowQueries struct {
	SlowQuery []*SlowQuery `json:"slowQueries,omitempty"` // A list of documents with information about slow queries as detected by the Performance Advisor.
}

// Shape represents a document with information about the query shapes that are served by the suggested indexes.
type Shape struct {
	AvgMs             float64      `json:"avgMs,omitempty"`             // Average duration in milliseconds for the queries examined that match this shape.
	Count             int64        `json:"count,omitempty"`             // Number of queries examined that match this shape.
	ID                string       `json:"id,omitempty"`                // Unique id for this shape. Exists only for the duration of the API request.
	InefficiencyScore int64        `json:"inefficiencyScore,omitempty"` //  Average number of documents read for every document returned by the query.
	Namespace         string       `json:"namespace,omitempty"`         // The namespace in which the slow query ran.
	Operations        []*Operation `json:"operations,omitempty"`        // It represents documents with specific information and log lines for individual queries.
}

// Operation represents a document with specific information and log lines for individual queries.
type Operation struct {
	Raw        string                   `json:"raw,omitempty"`        // Raw log line produced by the query.
	Stats      Stats                    `json:"stats,omitempty"`      // Query statistics.
	Predicates []map[string]interface{} `json:"predicates,omitempty"` // Documents containing the search criteria used by the query.
}

// Stats represents query statistics.
type Stats struct {
	MS        float64 `json:"ms,omitempty"`        // Duration in milliseconds of the query.
	NReturned int64   `json:"nReturned,omitempty"` // Number of results returned by the query.
	NScanned  int64   `json:"nScanned,omitempty"`  // Number of documents read by the query.
	TS        int64   `json:"ts,omitempty"`        // Query timestamp, in seconds since epoch.
}

// SuggestedIndex represents a suggested index.
type SuggestedIndex struct {
	ID        string           `json:"id,omitempty"`        // Unique id for this suggested index.
	Impact    []string         `json:"impact,omitempty"`    // List of unique identifiers which correspond the query shapes in this response which pertain to this suggested index.
	Namespace string           `json:"namespace,omitempty"` // 	Namespace of the suggested index.
	Weight    float64          `json:"weight,omitempty"`    // Estimated percentage performance improvement that the suggested index would provide.
	Index     []map[string]int `json:"index,omitempty"`     // Array of documents that specifies a key in the index and its sort order, ascending or descending.
}

// SuggestedIndexes represents an array of suggested indexes.
type SuggestedIndexes struct {
	SuggestedIndexes []*SuggestedIndex `json:"suggestedIndexes,omitempty"` // Documents with information about the indexes suggested by the Performance Advisor.
	Shapes           []*Shape          `json:"shapes,omitempty"`           // Documents with information about the query shapes that are served by the suggested indexes.

}

// SlowQueryOptions contains the request query parameters for the API request.
type SlowQueryOptions struct {
	Namespaces string `url:"namespaces,omitempty"` // Namespaces from which to retrieve slow query logs. A namespace consists of the database and collection resource separated by a ., such as <database>.<collection>.
	NLogs      int64  `url:"nLogs,omitempty"`      // Maximum number of log lines to return. Defaults to 20000.
	NamespaceOptions
}

// NamespaceOptions contains the request query parameters for the API request.
type NamespaceOptions struct {
	Since    int64 `url:"since,omitempty"`    // Point in time, specified as milliseconds since the Unix Epoch, from which you want to receive results.
	Duration int64 `url:"duration,omitempty"` // 	Length of time from the since parameter, in milliseconds, for which you want to receive results.
}

// SuggestedIndexOptions contains the request query parameters for the API request.
type SuggestedIndexOptions struct {
	Namespaces string `url:"namespaces,omitempty"` // Namespaces from which to retrieve slow query logs. A namespace consists of the database and collection resource separated by a ., such as <database>.<collection>.
	NIndexes   int64  `url:"nIndexes,omitempty"`   // Maximum number of indexes to suggest. Defaults to unlimited.
	NExamples  int64  `url:"NExamples,omitempty"`  // Maximum number of examples queries to provide that will be improved by a suggested index. Defaults to 5.
	NamespaceOptions
}

// GetNamespaces retrieves the namespaces for collections experiencing slow queries for a specified host.
//
// See more: https://docs.atlas.mongodb.com/reference/api/pa-namespaces-get-all/
func (s *PerformanceAdvisorServiceOp) GetNamespaces(ctx context.Context, groupID, processName string, opts *NamespaceOptions) (*Namespaces, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if processName == "" {
		return nil, nil, NewArgError("processName", "must be set")
	}

	path := fmt.Sprintf(performanceAdvisorNamespacesPath, groupID, processName)
	path, err := setListOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Namespaces)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetSlowQueries gets log lines for slow queries as determined by the Performance Advisor.
//
// See more: https://docs.atlas.mongodb.com/reference/api/pa-get-slow-query-logs/
func (s *PerformanceAdvisorServiceOp) GetSlowQueries(ctx context.Context, groupID, processName string, opts *SlowQueryOptions) (*SlowQueries, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if processName == "" {
		return nil, nil, NewArgError("processName", "must be set")
	}

	path := fmt.Sprintf(performanceAdvisorSlowQueryLogsPath, groupID, processName)
	path, err := setListOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SlowQueries)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetSuggestedIndexes gets suggested indexes as determined by the Performance Advisor.
//
// See more: https://docs.atlas.mongodb.com/reference/api/pa-suggested-indexes-get-all/
func (s *PerformanceAdvisorServiceOp) GetSuggestedIndexes(ctx context.Context, groupID, processName string, opts *SuggestedIndexOptions) (*SuggestedIndexes, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if processName == "" {
		return nil, nil, NewArgError("processName", "must be set")
	}

	path := fmt.Sprintf(performanceAdvisorSuggestedIndexesLogsPath, groupID, processName)
	path, err := setListOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SuggestedIndexes)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// EnableManagedSlowOperationThreshold enables the Atlas managed slow operation threshold for your project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/pa-managed-slow-ms-enable/
func (s *PerformanceAdvisorServiceOp) EnableManagedSlowOperationThreshold(ctx context.Context, groupID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(performanceAdvisorManagedSlowMs, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "enable")

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// DisableManagedSlowOperationThreshold disables the Atlas managed slow operation threshold for your project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/pa-managed-slow-ms-disable/
func (s *PerformanceAdvisorServiceOp) DisableManagedSlowOperationThreshold(ctx context.Context, groupID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(performanceAdvisorManagedSlowMs, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "disable")

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
