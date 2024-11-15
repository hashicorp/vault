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
	"net/url"
)

const advancedClustersPath = "api/atlas/v1.5/groups/%s/clusters"

// AdvancedClustersService is an interface for interfacing with the Clusters (Advanced)
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/clusters-advanced/
type AdvancedClustersService interface {
	List(ctx context.Context, groupID string, options *ListOptions) (*AdvancedClustersResponse, *Response, error)
	Get(ctx context.Context, groupID, clusterName string) (*AdvancedCluster, *Response, error)
	Create(ctx context.Context, groupID string, cluster *AdvancedCluster) (*AdvancedCluster, *Response, error)
	Update(ctx context.Context, groupID, clusterName string, cluster *AdvancedCluster) (*AdvancedCluster, *Response, error)
	Delete(ctx context.Context, groupID, clusterName string, options *DeleteAdvanceClusterOptions) (*Response, error)
	TestFailover(ctx context.Context, groupID, clusterName string) (*Response, error)
}

// AdvancedClustersServiceOp handles communication with the Cluster (Advanced) related methods
// of the MongoDB Atlas API.
type AdvancedClustersServiceOp service

var _ AdvancedClustersService = &AdvancedClustersServiceOp{}

// AdvancedCluster represents MongoDB cluster.
type AdvancedCluster struct {
	AcceptDataRisksAndForceReplicaSetReconfig string                     `json:"acceptDataRisksAndForceReplicaSetReconfig,omitempty"`
	BackupEnabled                             *bool                      `json:"backupEnabled,omitempty"`
	BiConnector                               *BiConnector               `json:"biConnector,omitempty"`
	ClusterType                               string                     `json:"clusterType,omitempty"`
	ConnectionStrings                         *ConnectionStrings         `json:"connectionStrings,omitempty"`
	DiskSizeGB                                *float64                   `json:"diskSizeGB,omitempty"`
	EncryptionAtRestProvider                  string                     `json:"encryptionAtRestProvider,omitempty"`
	GroupID                                   string                     `json:"groupId,omitempty"`
	ID                                        string                     `json:"id,omitempty"`
	Labels                                    []Label                    `json:"labels,omitempty"`
	MongoDBMajorVersion                       string                     `json:"mongoDBMajorVersion,omitempty"`
	MongoDBVersion                            string                     `json:"mongoDBVersion,omitempty"`
	Name                                      string                     `json:"name,omitempty"`
	Paused                                    *bool                      `json:"paused,omitempty"`
	PitEnabled                                *bool                      `json:"pitEnabled,omitempty"`
	StateName                                 string                     `json:"stateName,omitempty"`
	ReplicationSpecs                          []*AdvancedReplicationSpec `json:"replicationSpecs,omitempty"`
	CreateDate                                string                     `json:"createDate,omitempty"`
	RootCertType                              string                     `json:"rootCertType,omitempty"`
	VersionReleaseSystem                      string                     `json:"versionReleaseSystem,omitempty"`
	TerminationProtectionEnabled              *bool                      `json:"terminationProtectionEnabled,omitempty"`
	Tags                                      []*Tag                     `json:"tags,omitempty"`
}

type AdvancedReplicationSpec struct {
	NumShards     int                     `json:"numShards,omitempty"`
	ID            string                  `json:"id,omitempty"`
	ZoneName      string                  `json:"zoneName,omitempty"`
	RegionConfigs []*AdvancedRegionConfig `json:"regionConfigs,omitempty"`
}

type AdvancedRegionConfig struct {
	AnalyticsAutoScaling *AdvancedAutoScaling `json:"analyticsAutoScaling,omitempty"`
	AnalyticsSpecs       *Specs               `json:"analyticsSpecs,omitempty"`
	ElectableSpecs       *Specs               `json:"electableSpecs,omitempty"`
	ReadOnlySpecs        *Specs               `json:"readOnlySpecs,omitempty"`
	AutoScaling          *AdvancedAutoScaling `json:"autoScaling,omitempty"`
	BackingProviderName  string               `json:"backingProviderName,omitempty"`
	Priority             *int                 `json:"priority,omitempty"`
	ProviderName         string               `json:"providerName,omitempty"`
	RegionName           string               `json:"regionName,omitempty"`
}

type AdvancedAutoScaling struct {
	DiskGB  *DiskGB  `json:"diskGB,omitempty"`
	Compute *Compute `json:"compute,omitempty"`
}

type DiskGB struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type Specs struct {
	DiskIOPS      *int64 `json:"diskIOPS,omitempty"`
	EbsVolumeType string `json:"ebsVolumeType,omitempty"`
	InstanceSize  string `json:"instanceSize,omitempty"`
	NodeCount     *int   `json:"nodeCount,omitempty"`
}

// AdvancedClustersResponse is the response from the AdvancedClustersService.List.
type AdvancedClustersResponse struct {
	Links      []*Link            `json:"links,omitempty"`
	Results    []*AdvancedCluster `json:"results,omitempty"`
	TotalCount int                `json:"totalCount,omitempty"`
}

type DeleteAdvanceClusterOptions struct {
	// Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster.
	RetainBackups *bool `url:"retainBackups,omitempty"`
}

// List all clusters in the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cluster-advanced/get-all-cluster-advanced/
func (s *AdvancedClustersServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) (*AdvancedClustersResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	path := fmt.Sprintf(advancedClustersPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AdvancedClustersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get gets the cluster specified to {ClUSTER-NAME} from the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cluster-advanced/get-one-cluster-advanced/
func (s *AdvancedClustersServiceOp) Get(ctx context.Context, groupID, clusterName string) (*AdvancedCluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if err := checkClusterNameParam(clusterName); err != nil {
		return nil, nil, err
	}

	basePath := fmt.Sprintf(advancedClustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AdvancedCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create adds a cluster to the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cluster-advanced/create-one-cluster-advanced/
func (s *AdvancedClustersServiceOp) Create(ctx context.Context, groupID string, createRequest *AdvancedCluster) (*AdvancedCluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(advancedClustersPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(AdvancedCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update a cluster in the project associated to {GROUP-ID}
//
// See more: https://docs.atlas.mongodb.com/reference/api/cluster-advanced/modify-one-cluster-advanced/
func (s *AdvancedClustersServiceOp) Update(ctx context.Context, groupID, clusterName string, updateRequest *AdvancedCluster) (*AdvancedCluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(advancedClustersPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(AdvancedCluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the cluster specified to {CLUSTER-NAME} from the project associated to {GROUP-ID}.
func (s *AdvancedClustersServiceOp) Delete(ctx context.Context, groupID, clusterName string, options *DeleteAdvanceClusterOptions) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(advancedClustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	// Add query params from options
	path, err := setListOptions(path, options)
	if err != nil {
		return nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// TestFailover starts a failover test for the specified cluster in the specified project
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Multi-Cloud-Clusters/operation/testFailover
func (s *AdvancedClustersServiceOp) TestFailover(ctx context.Context, groupID, clusterName string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(advancedClustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s/restartPrimaries", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
