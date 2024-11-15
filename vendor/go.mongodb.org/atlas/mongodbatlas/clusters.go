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

type ChangeStatus string

const (
	// ChangeStatusApplied signals when changes to the deployments have completed.
	ChangeStatusApplied ChangeStatus = "APPLIED"
	// ChangeStatusPending signals when changes to the deployments are still pending.
	ChangeStatusPending          ChangeStatus = "PENDING"
	clustersPath                              = "api/atlas/v1.0/groups/%s/clusters"
	sampleDatasetLoadPath                     = "api/atlas/v1.0/groups/%s/sampleDatasetLoad"
	cloudProviderRegionsBasePath              = clustersPath + "/provider/regions"
)

// ClustersService is an interface for interfacing with the Clusters
// endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/clusters/
type ClustersService interface {
	List(ctx context.Context, groupID string, options *ListOptions) ([]Cluster, *Response, error)
	Get(ctx context.Context, groupID, clusterName string) (*Cluster, *Response, error)
	Create(ctx context.Context, groupID string, cluster *Cluster) (*Cluster, *Response, error)
	Update(ctx context.Context, groupID, clusterName string, cluster *Cluster) (*Cluster, *Response, error)
	Delete(ctx context.Context, groupID, clusterName string, options *DeleteAdvanceClusterOptions) (*Response, error)
	UpdateProcessArgs(ctx context.Context, groupID, clusterName string, args *ProcessArgs) (*ProcessArgs, *Response, error)
	GetProcessArgs(ctx context.Context, groupID, clusterName string) (*ProcessArgs, *Response, error)
	Status(ctx context.Context, groupID, clusterName string) (ClusterStatus, *Response, error)
	LoadSampleDataset(ctx context.Context, groupID, clusterName string) (*SampleDatasetJob, *Response, error)
	GetSampleDatasetStatus(ctx context.Context, groupID, id string) (*SampleDatasetJob, *Response, error)
	ListCloudProviderRegions(context.Context, string, *CloudProviderRegionsOptions) (*CloudProviders, *Response, error)
	Upgrade(ctx context.Context, groupID string, cluster *Cluster) (*Cluster, *Response, error)
}

// ClustersServiceOp handles communication with the Cluster related methods
// of the MongoDB Atlas API.
type ClustersServiceOp service

var _ ClustersService = &ClustersServiceOp{}

// AutoScaling configures your cluster to automatically scale its storage.
type AutoScaling struct {
	AutoIndexingEnabled *bool    `json:"autoIndexingEnabled,omitempty"` // Autopilot mode is only available if you are enrolled in the Auto Pilot Early Access program.
	Compute             *Compute `json:"compute,omitempty"`
	DiskGBEnabled       *bool    `json:"diskGBEnabled,omitempty"`
}

// Compute Specifies whether the cluster automatically scales its cluster tier and whether the cluster can scale down.
type Compute struct {
	Enabled          *bool  `json:"enabled,omitempty"`
	ScaleDownEnabled *bool  `json:"scaleDownEnabled,omitempty"`
	MinInstanceSize  string `json:"minInstanceSize,omitempty"`
	MaxInstanceSize  string `json:"maxInstanceSize,omitempty"`
}

// BiConnector specifies BI Connector for Atlas configuration on this cluster.
type BiConnector struct {
	Enabled        *bool  `json:"enabled,omitempty"`
	ReadPreference string `json:"readPreference,omitempty"`
}

// ProviderSettings configuration for the provisioned servers on which MongoDB runs. The available options are specific to the cloud service provider.
type ProviderSettings struct {
	BackingProviderName string       `json:"backingProviderName,omitempty"`
	DiskIOPS            *int64       `json:"diskIOPS,omitempty"`
	DiskTypeName        string       `json:"diskTypeName,omitempty"`
	EncryptEBSVolume    *bool        `json:"encryptEBSVolume,omitempty"`
	InstanceSizeName    string       `json:"instanceSizeName,omitempty"`
	ProviderName        string       `json:"providerName,omitempty"`
	RegionName          string       `json:"regionName,omitempty"`
	VolumeType          string       `json:"volumeType,omitempty"`
	AutoScaling         *AutoScaling `json:"autoScaling,omitempty"`
}

// RegionsConfig describes the region’s priority in elections and the number and type of MongoDB nodes Atlas deploys to the region.
type RegionsConfig struct {
	AnalyticsNodes *int64 `json:"analyticsNodes,omitempty"`
	ElectableNodes *int64 `json:"electableNodes,omitempty"`
	Priority       *int64 `json:"priority,omitempty"`
	ReadOnlyNodes  *int64 `json:"readOnlyNodes,omitempty"`
}

// ReplicationSpec represents a configuration for cluster regions.
type ReplicationSpec struct {
	ID            string                   `json:"id,omitempty"`
	NumShards     *int64                   `json:"numShards,omitempty"`
	ZoneName      string                   `json:"zoneName,omitempty"`
	RegionsConfig map[string]RegionsConfig `json:"regionsConfig,omitempty"`
}

// PrivateEndpoint connection strings. Each object describes the connection strings
// you can use to connect to this cluster through a private endpoint.
// Atlas returns this parameter only if you deployed a private endpoint to all regions
// to which you deployed this cluster's nodes.
type PrivateEndpoint struct {
	ConnectionString                  string     `json:"connectionString,omitempty"`
	Endpoints                         []Endpoint `json:"endpoints,omitempty"`
	SRVConnectionString               string     `json:"srvConnectionString,omitempty"`
	SRVShardOptimizedConnectionString string     `json:"srvShardOptimizedConnectionString,omitempty"`
	Type                              string     `json:"type,omitempty"`
}

// Endpoint through which you connect to Atlas.
type Endpoint struct {
	EndpointID   string `json:"endpointId,omitempty"`
	ProviderName string `json:"providerName,omitempty"`
	Region       string `json:"region,omitempty"`
}

// ConnectionStrings configuration for applications use to connect to this cluster.
type ConnectionStrings struct {
	Standard          string            `json:"standard,omitempty"`
	StandardSrv       string            `json:"standardSrv,omitempty"`
	PrivateEndpoint   []PrivateEndpoint `json:"privateEndpoint,omitempty"`
	AwsPrivateLink    map[string]string `json:"awsPrivateLink,omitempty"`    // Deprecated: Use connectionStrings.PrivateEndpoint[n].ConnectionString
	AwsPrivateLinkSrv map[string]string `json:"awsPrivateLinkSrv,omitempty"` // Deprecated: Use ConnectionStrings.privateEndpoint[n].SRVConnectionString
	Private           string            `json:"private,omitempty"`
	PrivateSrv        string            `json:"privateSrv,omitempty"`
}

// Cluster represents MongoDB cluster.
type Cluster struct {
	AcceptDataRisksAndForceReplicaSetReconfig string                   `json:"acceptDataRisksAndForceReplicaSetReconfig,omitempty"`
	AutoScaling                               *AutoScaling             `json:"autoScaling,omitempty"`
	BackupEnabled                             *bool                    `json:"backupEnabled,omitempty"` // Deprecated: Use ProviderBackupEnabled instead
	BiConnector                               *BiConnector             `json:"biConnector,omitempty"`
	ClusterType                               string                   `json:"clusterType,omitempty"`
	DiskSizeGB                                *float64                 `json:"diskSizeGB,omitempty"`
	EncryptionAtRestProvider                  string                   `json:"encryptionAtRestProvider,omitempty"`
	Labels                                    []Label                  `json:"labels,omitempty"`
	ID                                        string                   `json:"id,omitempty"`
	GroupID                                   string                   `json:"groupId,omitempty"`
	MongoDBVersion                            string                   `json:"mongoDBVersion,omitempty"`
	MongoDBMajorVersion                       string                   `json:"mongoDBMajorVersion,omitempty"`
	MongoURI                                  string                   `json:"mongoURI,omitempty"`
	MongoURIUpdated                           string                   `json:"mongoURIUpdated,omitempty"`
	MongoURIWithOptions                       string                   `json:"mongoURIWithOptions,omitempty"`
	Name                                      string                   `json:"name,omitempty"`
	CreateDate                                string                   `json:"createDate,omitempty"`
	NumShards                                 *int64                   `json:"numShards,omitempty"`
	Paused                                    *bool                    `json:"paused,omitempty"`
	PitEnabled                                *bool                    `json:"pitEnabled,omitempty"`
	ProviderBackupEnabled                     *bool                    `json:"providerBackupEnabled,omitempty"`
	ProviderSettings                          *ProviderSettings        `json:"providerSettings,omitempty"`
	ReplicationFactor                         *int64                   `json:"replicationFactor,omitempty"`
	ReplicationSpec                           map[string]RegionsConfig `json:"replicationSpec,omitempty"`
	ReplicationSpecs                          []ReplicationSpec        `json:"replicationSpecs,omitempty"`
	SrvAddress                                string                   `json:"srvAddress,omitempty"`
	StateName                                 string                   `json:"stateName,omitempty"`
	ServerlessBackupOptions                   *ServerlessBackupOptions `json:"serverlessBackupOptions,omitempty"`
	ConnectionStrings                         *ConnectionStrings       `json:"connectionStrings,omitempty"`
	Links                                     []*Link                  `json:"links,omitempty"`
	VersionReleaseSystem                      string                   `json:"versionReleaseSystem,omitempty"`
	RootCertType                              string                   `json:"rootCertType,omitempty"`
	TerminationProtectionEnabled              *bool                    `json:"terminationProtectionEnabled,omitempty"`
	Tags                                      *[]*Tag                  `json:"tags,omitempty"`
}

// ProcessArgs represents the advanced configuration options for the cluster.
type ProcessArgs struct {
	DefaultReadConcern                                    string   `json:"defaultReadConcern,omitempty"`
	DefaultWriteConcern                                   string   `json:"defaultWriteConcern,omitempty"`
	MinimumEnabledTLSProtocol                             string   `json:"minimumEnabledTlsProtocol,omitempty"`
	FailIndexKeyTooLong                                   *bool    `json:"failIndexKeyTooLong,omitempty"`
	JavascriptEnabled                                     *bool    `json:"javascriptEnabled,omitempty"`
	NoTableScan                                           *bool    `json:"noTableScan,omitempty"`
	OplogSizeMB                                           *int64   `json:"oplogSizeMB,omitempty"`
	SampleSizeBIConnector                                 *int64   `json:"sampleSizeBIConnector,omitempty"`
	SampleRefreshIntervalBIConnector                      *int64   `json:"sampleRefreshIntervalBIConnector,omitempty"`
	TransactionLifetimeLimitSeconds                       *int64   `json:"transactionLifetimeLimitSeconds,omitempty"`
	OplogMinRetentionHours                                *float64 `json:"oplogMinRetentionHours,omitempty"`
	ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds *int64   `json:"changeStreamOptionsPreAndPostImagesExpireAfterSeconds,omitempty"`
}

type Tag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// ClusterStatus is the status of the operations on the cluster.
type ClusterStatus struct {
	ChangeStatus ChangeStatus `json:"changeStatus"`
}

// clustersResponse is the response from the ClustersService.List.
type clustersResponse struct {
	Links      []*Link   `json:"links,omitempty"`
	Results    []Cluster `json:"results,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

// SampleDatasetJob represents a sample dataset job.
type SampleDatasetJob struct {
	ClusterName  string `json:"clusterName"`
	CompleteDate string `json:"completeDate,omitempty"`
	CreateDate   string `json:"createDate,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ID           string `json:"_id"`
	State        string `json:"state"`
}

// CloudProvider represents a cloud provider of the MongoDB Atlas API.
type CloudProvider struct {
	Provider      string          `json:"provider,omitempty"`
	InstanceSizes []*InstanceSize `json:"instanceSizes,omitempty"`
}

// InstanceSize represents an instance size of the MongoDB Atlas API.
type InstanceSize struct {
	Name             string             `json:"name,omitempty"`
	AvailableRegions []*AvailableRegion `json:"availableRegions,omitempty"`
}

// AvailableRegion represents an available region of the MongoDB Atlas API.
type AvailableRegion struct {
	Name    string `json:"name,omitempty"`
	Default bool   `json:"default,omitempty"`
}

// CloudProviders represents the response from CloudProviderRegionsService.Get.
type CloudProviders struct {
	Links      []*Link          `json:"links,omitempty"`
	Results    []*CloudProvider `json:"results,omitempty"`
	TotalCount int              `json:"totalCount,omitempty"`
}

// CloudProviderRegionsOptions specifies the optional parameters to the CloudProviderRegions Get method.
type CloudProviderRegionsOptions struct {
	Providers []*string `url:"providers,omitempty"`
	Tier      string    `url:"tier,omitempty"`
}

// DefaultDiskSizeGB represents the Tier and the default disk size for each one
// it can be use like: DefaultDiskSizeGB["AWS"]["M10"].
var DefaultDiskSizeGB = map[string]map[string]float64{
	"TENANT": {
		"M2": 2,
		"M5": 5,
	},
	"AWS": {
		"M10":       10,
		"M20":       20,
		"M30":       40,
		"M40":       80,
		"R40":       80,
		"M40_NVME":  380,
		"M50":       160,
		"R50":       160,
		"M50_NVME":  760,
		"M60":       320,
		"R60":       320,
		"M60_NVME":  1600,
		"M80":       750,
		"R80":       750,
		"M80_NVME":  1600,
		"M140":      1000,
		"M200":      1500,
		"R200":      1500,
		"M200_NVME": 3100,
		"M300":      2000,
		"R300":      2000,
		"R400":      3000,
		"M400_NVME": 4000,
	},
	"GCP": {
		"M10":  10,
		"M20":  20,
		"M30":  40,
		"M40":  80,
		"M50":  160,
		"M60":  320,
		"M80":  750,
		"M200": 1500,
		"M300": 2200,
	},
	"AZURE": {
		"M10":  32,
		"M20":  32,
		"M30":  32,
		"M40":  128,
		"M50":  128,
		"M60":  128,
		"M80":  256,
		"M200": 256,
	},
}

// List all clusters in the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/clusters-get-all/
func (s *ClustersServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]Cluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	path := fmt.Sprintf(clustersPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(clustersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets the cluster specified to {ClUSTER-NAME} from the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/clusters-get-one/
func (s *ClustersServiceOp) Get(ctx context.Context, groupID, clusterName string) (*Cluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if err := checkClusterNameParam(clusterName); err != nil {
		return nil, nil, err
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create adds a cluster to the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/clusters-create-one/
func (s *ClustersServiceOp) Create(ctx context.Context, groupID string, createRequest *Cluster) (*Cluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(clustersPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update a cluster in the project associated to {GROUP-ID}
//
// See more: https://docs.atlas.mongodb.com/reference/api/clusters-modify-one/
func (s *ClustersServiceOp) Update(ctx context.Context, groupID, clusterName string, updateRequest *Cluster) (*Cluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the cluster specified to {CLUSTER-NAME} from the project associated to {GROUP-ID}.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Clusters/operation/deleteLegacyCluster
func (s *ClustersServiceOp) Delete(ctx context.Context, groupID, clusterName string, options *DeleteAdvanceClusterOptions) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
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

// UpdateProcessArgs Modifies Advanced Configuration Options for One Cluster
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/updateAdvancedConfigurationOptionsForOneCluster
func (s *ClustersServiceOp) UpdateProcessArgs(ctx context.Context, groupID, clusterName string, updateRequest *ProcessArgs) (*ProcessArgs, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	path := fmt.Sprintf("%s/%s/processArgs", basePath, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessArgs)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetProcessArgs gets the Advanced Configuration Options for One Cluster
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/returnOneAdvancedConfigurationOptionsForOneCluster
func (s *ClustersServiceOp) GetProcessArgs(ctx context.Context, groupID, clusterName string) (*ProcessArgs, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if err := checkClusterNameParam(clusterName); err != nil {
		return nil, nil, err
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s/processArgs", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessArgs)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// LoadSampleDataset loads the sample dataset into your cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/cluster/load-dataset/
func (s *ClustersServiceOp) LoadSampleDataset(ctx context.Context, groupID, clusterName string) (*SampleDatasetJob, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if err := checkClusterNameParam(clusterName); err != nil {
		return nil, nil, err
	}

	basePath := fmt.Sprintf(sampleDatasetLoadPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SampleDatasetJob)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetSampleDatasetStatus gets the Sample Dataset job
//
// See more: https://docs.atlas.mongodb.com/reference/api/cluster/check-dataset-status/
func (s *ClustersServiceOp) GetSampleDatasetStatus(ctx context.Context, groupID, id string) (*SampleDatasetJob, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	basePath := fmt.Sprintf(sampleDatasetLoadPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, id)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SampleDatasetJob)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Status gets the status of the operation on the Cluster.
//
// See more: https://docs.atlas.mongodb.com/reference/api/clusters-check-operation-status/
func (s *ClustersServiceOp) Status(ctx context.Context, groupID, clusterName string) (ClusterStatus, *Response, error) {
	var root ClusterStatus
	if groupID == "" {
		return root, nil, NewArgError("groupId", "must be set")
	}
	if err := checkClusterNameParam(clusterName); err != nil {
		return root, nil, err
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s/status", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return root, nil, err
	}

	resp, err := s.Client.Do(ctx, req, &root)
	return root, resp, err
}

// ListCloudProviderRegions gets the available regions for each cloud provider
//
// See more: https://docs.atlas.mongodb.com/reference/api/cluster-get-regions/
func (s *ClustersServiceOp) ListCloudProviderRegions(ctx context.Context, groupID string, options *CloudProviderRegionsOptions) (*CloudProviders, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	path := fmt.Sprintf(cloudProviderRegionsBasePath, groupID)

	path, err := setListOptions(path, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CloudProviders)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func checkClusterNameParam(clusterName string) error {
	if clusterName == "" {
		return NewArgError("name", "must be set")
	}
	return nil
}

// Upgrade a cluster in the project associated to {GROUP-ID}
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#operation/upgradeOneTenantCluster
func (s *ClustersServiceOp) Upgrade(ctx context.Context, groupID string, upgradeRequest *Cluster) (*Cluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if upgradeRequest == nil {
		return nil, nil, NewArgError("upgradeRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	path := fmt.Sprintf("%s/tenantUpgrade", basePath)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, upgradeRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
