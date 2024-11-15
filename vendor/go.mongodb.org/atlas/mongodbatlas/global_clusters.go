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

const globalClustersBasePath = "api/atlas/v1.5/groups/%s/clusters/%s/globalWrites/%s"

// GlobalClustersService is an interface for interfacing with the Global Clusters
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/global-clusters/
type GlobalClustersService interface {
	Get(context.Context, string, string) (*GlobalCluster, *Response, error)
	AddManagedNamespace(context.Context, string, string, *ManagedNamespace) (*GlobalCluster, *Response, error)
	DeleteManagedNamespace(context.Context, string, string, *ManagedNamespace) (*GlobalCluster, *Response, error)
	AddCustomZoneMappings(context.Context, string, string, *CustomZoneMappingsRequest) (*GlobalCluster, *Response, error)
	DeleteCustomZoneMappings(context.Context, string, string) (*GlobalCluster, *Response, error)
}

// GlobalClustersServiceOp handles communication with the GlobalClusters related methods of the
// MongoDB Atlas API.
type GlobalClustersServiceOp service

var _ GlobalClustersService = &GlobalClustersServiceOp{}

// GlobalCluster represents MongoDB Global Cluster Configuration in your Global Cluster.
type GlobalCluster struct {
	CustomZoneMapping map[string]string  `json:"customZoneMapping"`
	ManagedNamespaces []ManagedNamespace `json:"managedNamespaces"`
}

// ManagedNamespace represents the information about managed namespace configuration.
type ManagedNamespace struct {
	Db                     string `json:"db"` //nolint:stylecheck // not changing this as is a breaking change
	Collection             string `json:"collection"`
	CustomShardKey         string `json:"customShardKey,omitempty"`
	IsCustomShardKeyHashed *bool  `json:"isCustomShardKeyHashed,omitempty"` // Flag that specifies whether the custom shard key for the collection is hashed.
	IsShardKeyUnique       *bool  `json:"isShardKeyUnique,omitempty"`       // Flag that specifies whether the underlying index enforces a unique constraint.
	NumInitialChunks       int    `json:"numInitialChunks,omitempty"`
	PresplitHashedZones    *bool  `json:"presplitHashedZones,omitempty"`
}

// CustomZoneMappingsRequest represents the request related to add custom zone mappings to a global cluster.
type CustomZoneMappingsRequest struct {
	CustomZoneMappings []CustomZoneMapping `json:"customZoneMappings"`
}

// CustomZoneMapping represents the custom zone mapping.
type CustomZoneMapping struct {
	Location string `json:"location"`
	Zone     string `json:"zone"`
}

// Get retrieves all managed namespaces and custom zone mappings associated with the specified Global Cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/global-clusters-retrieve-namespaces/
func (s *GlobalClustersServiceOp) Get(ctx context.Context, groupID, clusterName string) (*GlobalCluster, *Response, error) {
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(globalClustersBasePath, groupID, clusterName, "")

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// AddManagedNamespace adds a managed namespace to the specified Global Cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/
func (s *GlobalClustersServiceOp) AddManagedNamespace(ctx context.Context, groupID, clusterName string, createRequest *ManagedNamespace) (*GlobalCluster, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(globalClustersBasePath, groupID, clusterName, "managedNamespaces")

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeleteManagedNamespace deletes the managed namespace configuration of the global cluster given.
//
// See more: https://docs.atlas.mongodb.com/reference/api/global-clusters-delete-namespace/
func (s *GlobalClustersServiceOp) DeleteManagedNamespace(ctx context.Context, groupID, clusterName string, deleteRequest *ManagedNamespace) (*GlobalCluster, *Response, error) {
	if deleteRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(globalClustersBasePath, groupID, clusterName, "managedNamespaces")

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	q.Add("collection", deleteRequest.Collection)
	q.Add("db", deleteRequest.Db)
	req.URL.RawQuery = q.Encode()

	root := new(GlobalCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// AddCustomZoneMappings adds an entry to the list of custom zone mappings for the specified Global Cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/global-clusters-add-customzonemapping/
func (s *GlobalClustersServiceOp) AddCustomZoneMappings(ctx context.Context, groupID, clusterName string, createRequest *CustomZoneMappingsRequest) (*GlobalCluster, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(globalClustersBasePath, groupID, clusterName, "customZoneMapping")

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeleteCustomZoneMappings removes all custom zone mappings from the specified Global Cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/global-clusters-delete-namespace/
func (s *GlobalClustersServiceOp) DeleteCustomZoneMappings(ctx context.Context, groupID, clusterName string) (*GlobalCluster, *Response, error) {
	path := fmt.Sprintf(globalClustersBasePath, groupID, clusterName, "customZoneMapping")

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root, resp, err
}
