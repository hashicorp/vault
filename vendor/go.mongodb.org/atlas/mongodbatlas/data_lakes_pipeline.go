// Copyright 2023 MongoDB Inc
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
	dataLakesPipelineBasePath = "api/atlas/v1.0/groups/%s/pipelines"
	dataLakesPipelineRunPath  = dataLakesPipelineBasePath + "/%s/runs"
)

// DataLakePipelineService is an interface for interfacing with the Data Lake Pipeline endpoints of the MongoDB Atlas API.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines
type DataLakePipelineService interface {
	List(context.Context, string) ([]*DataLakePipeline, *Response, error)
	ListSnapshots(context.Context, string, string, *ListDataLakePipelineSnapshotOptions) (*DataLakePipelineSnapshotsResponse, *Response, error)
	ListIngestionSchedules(context.Context, string, string) ([]*DataLakePipelineIngestionSchedule, *Response, error)
	ListRuns(context.Context, string, string) (*DataLakePipelineRunsResponse, *Response, error)
	Get(context.Context, string, string) (*DataLakePipeline, *Response, error)
	GetRun(context.Context, string, string, string) (*DataLakePipelineRun, *Response, error)
	Create(context.Context, string, *DataLakePipeline) (*DataLakePipeline, *Response, error)
	Update(context.Context, string, string, *DataLakePipeline) (*DataLakePipeline, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// DataLakePipelineServiceOp handles communication with the DataLakePipelineService related methods of the
// MongoDB Atlas API.
type DataLakePipelineServiceOp service

var _ DataLakePipelineService = &DataLakePipelineServiceOp{}

// DataLakePipeline represents a store of data lake data. Docs: https://docs.mongodb.com/datalake/reference/format/data-lake-configuration/#stores
type DataLakePipeline struct {
	ID              string                            `json:"_id,omitempty"`             // Unique 24-hexadecimal digit string that identifies the Data Lake Pipeline.
	GroupID         string                            `json:"groupId,omitempty"`         // Unique identifier for the project.
	Name            string                            `json:"name,omitempty"`            // Name of this Data Lake Pipeline.
	CreatedDate     string                            `json:"createdDate,omitempty"`     // Timestamp that indicates when the Data Lake Pipeline was created.
	LastUpdatedDate string                            `json:"lastUpdatedDate,omitempty"` // Timestamp that indicates the last time that the Data Lake Pipeline was updated.
	State           string                            `json:"state,omitempty"`           // State of this Data Lake Pipeline.
	Sink            *DataLakePipelineSink             `json:"sink,omitempty"`            // Ingestion destination of a Data Lake Pipeline.
	Source          *DataLakePipelineSource           `json:"source,omitempty"`          // Ingestion Source of a Data Lake Pipeline.
	Transformations []*DataLakePipelineTransformation `json:"transformations,omitempty"` // Fields to be excluded for this Data Lake Pipeline.
}

// DataLakePipelineTransformation represents fields to be excluded for this Data Lake Pipeline.
type DataLakePipelineTransformation struct {
	Field string `json:"field,omitempty"` // Key in the document.
	Type  string `json:"type,omitempty"`  // Type of transformation applied during the export of the namespace in a Data Lake Pipeline.
}

// DataLakePipelineSink represents ingestion destination of a Data Lake Pipeline.
type DataLakePipelineSink struct {
	Type             string                            `json:"type,omitempty"`             // Type of ingestion destination of this Data Lake Pipeline.
	MetadataProvider string                            `json:"metadataProvider,omitempty"` // Target cloud provider for this Data Lake Pipeline.
	MetadataRegion   string                            `json:"metadataRegion,omitempty"`   // Target cloud provider region for this Data Lake Pipeline.
	PartitionFields  []*DataLakePipelinePartitionField `json:"partitionFields,omitempty"`  // Ordered fields used to physically organize data in the destination.
}

// DataLakePipelinePartitionField represents ordered fields used to physically organize data in the destination.
type DataLakePipelinePartitionField struct {
	FieldName string `json:"fieldName,omitempty"`
	Order     int32  `json:"order,omitempty"`
}

// DataLakePipelineSource represents the storage configuration for a data lake.
type DataLakePipelineSource struct {
	Type           string `json:"type,omitempty"`           // Type of ingestion source of this Data Lake Pipeline.
	ClusterName    string `json:"clusterName,omitempty"`    // Human-readable name that identifies the cluster.
	CollectionName string `json:"collectionName,omitempty"` // Human-readable name that identifies the collection.
	DatabaseName   string `json:"databaseName,omitempty"`   // Human-readable name that identifies the database.
	PolicyItemID   string `json:"policyItemId,omitempty"`   // Unique 24-hexadecimal character string that identifies a policy item.
	GroupID        string `json:"groupId,omitempty"`        // Unique 24-hexadecimal character string that identifies the project.
}

// ListDataLakePipelineSnapshotOptions specifies the optional parameters to ListSnapshots method.
type ListDataLakePipelineSnapshotOptions struct {
	*ListOptions
	CompletedAfter string `url:"completedAfter,omitempty"` // Date and time after which MongoDB Cloud created the snapshot.
}

// DataLakePipelineSnapshotsResponse represents the response of DataLakePipelineService.ListSnapshots.
type DataLakePipelineSnapshotsResponse struct {
	Links      []*Link                     `json:"links,omitempty"`      // List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both.
	Results    []*DataLakePipelineSnapshot `json:"results,omitempty"`    // List of returned documents that MongoDB Cloud providers when completing this request.
	TotalCount int                         `json:"totalCount,omitempty"` // Number of documents returned in this response.
}

// DataLakePipelineSnapshot represents a snapshot that you can use to trigger an on demand pipeline run.
type DataLakePipelineSnapshot struct {
	ID               string   `json:"id,omitempty"`               // Unique 24-hexadecimal digit string that identifies the snapshot.
	CloudProvider    string   `json:"cloudProvider,omitempty"`    // Human-readable label that identifies the cloud provider that stores this snapshot.
	CreatedAt        string   `json:"createdAt,omitempty"`        // Date and time when MongoDB Cloud took the snapshot.
	Description      string   `json:"description,omitempty"`      // Human-readable phrase or sentence that explains the purpose of the snapshot.
	ExpiresAt        string   `json:"expiresAt,omitempty"`        // Date and time when MongoDB Cloud deletes the snapshot.
	FrequencyType    string   `json:"frequencyType,omitempty"`    // Human-readable label that identifies how often this snapshot triggers.
	MasterKeyUUID    string   `json:"masterKeyUUID,omitempty"`    // Unique string that identifies the Amazon Web Services (AWS) Key Management Service (KMS) Customer Master Key (CMK) used to encrypt the snapshot.
	MongodVersion    string   `json:"mongodVersion,omitempty"`    // Version of the MongoDB host that this snapshot backs up.
	ReplicaSetName   string   `json:"replicaSetName,omitempty"`   // Human-readable label that identifies the replica set from which MongoDB Cloud took this snapshot.
	SnapshotType     string   `json:"snapshotType,omitempty"`     // Human-readable label that identifies when this snapshot triggers.
	Status           string   `json:"status,omitempty"`           // Human-readable label that indicates the stage of the backup process for this snapshot.
	Type             string   `json:"type,omitempty"`             // Human-readable label that categorizes the cluster as a replica set or sharded cluster.
	StorageSizeBytes int64    `json:"storageSizeBytes,omitempty"` // Number of bytes taken to store the backup snapshot.
	CopyRegions      []string `json:"copyRegions,omitempty"`      // List that identifies the regions to which MongoDB Cloud copies the snapshot.
	PolicyItems      []string `json:"policyItems,omitempty"`      // List that contains unique identifiers for the policy items.
	Links            []*Link  `json:"links,omitempty"`            // List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both.
}

// DataLakePipelineIngestionSchedule represents a backup schedule policy item that you can use as a Data Lake Pipeline source.
type DataLakePipelineIngestionSchedule struct {
	ID                string `json:"id,omitempty"`                // Unique 24-hexadecimal digit string that identifies this backup policy item.
	FrequencyType     string `json:"frequencyType,omitempty"`     // Human-readable label that identifies the frequency type associated with the backup policy.
	RetentionUnit     string `json:"retentionUnit,omitempty"`     // Unit of time in which MongoDB Cloud measures snapshot retention.
	FrequencyInterval int32  `json:"frequencyInterval,omitempty"` // Number that indicates the frequency interval for a set of snapshots.
	RetentionValue    int32  `json:"retentionValue,omitempty"`    // Duration in days, weeks, or months that MongoDB Cloud retains the snapshot.
}

// DataLakePipelineRunsResponse represents the response of DataLakePipelineService.ListRuns.
type DataLakePipelineRunsResponse struct {
	Links      []*Link                `json:"links,omitempty"`      // List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both.
	Results    []*DataLakePipelineRun `json:"results,omitempty"`    // List of returned documents that MongoDB Cloud providers when completing this request.
	TotalCount int                    `json:"totalCount,omitempty"` // Number of documents returned in this response.
}

// DataLakePipelineRun represents a DataLake Pipeline Run.
type DataLakePipelineRun struct {
	ID                  string                    `json:"_id,omitempty"`                 // Unique 24-hexadecimal character string that identifies a Data Lake Pipeline run.
	BackupFrequencyType string                    `json:"backupFrequencyType,omitempty"` // Backup schedule interval of the Data Lake Pipeline.
	CreatedDate         string                    `json:"createdDate,omitempty"`         // Timestamp that indicates when the pipeline run was created.
	DatasetName         string                    `json:"datasetName,omitempty"`         // Human-readable label that identifies the dataset that Atlas generates during this pipeline run.
	GroupID             string                    `json:"groupId,omitempty"`             // Unique 24-hexadecimal character string that identifies the project.
	LastUpdatedDate     string                    `json:"lastUpdatedDate,omitempty"`     // Timestamp that indicates the last time that the pipeline run was updated.
	Phase               string                    `json:"phase,omitempty"`               // Processing phase of the Data Lake Pipeline.
	PipelineID          string                    `json:"pipelineId,omitempty"`          // Unique 24-hexadecimal character string that identifies a Data Lake Pipeline.
	SnapshotID          string                    `json:"snapshotId,omitempty"`          // Unique 24-hexadecimal character string that identifies the snapshot of a cluster.
	State               string                    `json:"state,omitempty"`               // State of the pipeline run.
	Stats               *DataLakePipelineRunStats `json:"stats,omitempty"`               // Runtime statistics for this Data Lake Pipeline run.
}

// DataLakePipelineRunStats represents runtime statistics for this Data Lake Pipeline run.
type DataLakePipelineRunStats struct {
	BytesExported int64 `json:"bytesExported,omitempty"` // Total data size in bytes exported for this pipeline run.
	NumDocs       int64 `json:"numDocs,omitempty"`       // Number of docs ingested for a this pipeline run.
}

// List gets a list of Data Lake Pipelines.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/listPipelines
func (s *DataLakePipelineServiceOp) List(ctx context.Context, groupID string) ([]*DataLakePipeline, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(dataLakesPipelineBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*DataLakePipeline
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// ListSnapshots gets a list of backup snapshots that you can use to trigger an on demand pipeline run.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/listPipelineSchedules
func (s *DataLakePipelineServiceOp) ListSnapshots(ctx context.Context, groupID, name string, options *ListDataLakePipelineSnapshotOptions) (*DataLakePipelineSnapshotsResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(dataLakesPipelineBasePath, groupID)
	path := fmt.Sprintf("%s/%s/availableSnapshots", basePath, name)

	// Add query params from ListDataLakePipelineSnapshotOptions
	pathWithOptions, err := setListOptions(path, options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, pathWithOptions, nil)
	if err != nil {
		return nil, nil, err
	}

	var root *DataLakePipelineSnapshotsResponse
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// ListIngestionSchedules gets a list of backup schedule policy items that you can use as a Data Lake Pipeline source.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/listPipelineSchedules
func (s *DataLakePipelineServiceOp) ListIngestionSchedules(ctx context.Context, groupID, name string) ([]*DataLakePipelineIngestionSchedule, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(dataLakesPipelineBasePath, groupID)
	path := fmt.Sprintf("%s/%s/availableSchedules", basePath, name)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*DataLakePipelineIngestionSchedule
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// ListRuns gets a list of past Data Lake Pipeline runs.
//
// https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/listPipelineRuns
func (s *DataLakePipelineServiceOp) ListRuns(ctx context.Context, groupID, name string) (*DataLakePipelineRunsResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	path := fmt.Sprintf(dataLakesPipelineRunPath, groupID, name)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLakePipelineRunsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Get gets the details of one Data Lake Pipeline within the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/getPipeline
func (s *DataLakePipelineServiceOp) Get(ctx context.Context, groupID, name string) (*DataLakePipeline, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(dataLakesPipelineBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, name)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLakePipeline)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetRun gets the details of one Data Lake Pipeline run within the specified project.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/getPipelineRun
func (s *DataLakePipelineServiceOp) GetRun(ctx context.Context, groupID, name, id string) (*DataLakePipelineRun, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	if id == "" {
		return nil, nil, NewArgError("id", "must be set")
	}

	basePath := fmt.Sprintf(dataLakesPipelineRunPath, groupID, name)
	path := fmt.Sprintf("%s/%s", basePath, id)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLakePipelineRun)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates one Data Lake Pipeline.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/createPipeline
func (s *DataLakePipelineServiceOp) Create(ctx context.Context, groupID string, createRequest *DataLakePipeline) (*DataLakePipeline, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "must be set")
	}

	path := fmt.Sprintf(dataLakesPipelineBasePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLakePipeline)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates one Data Lake Pipeline.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/updatePipeline
func (s *DataLakePipelineServiceOp) Update(ctx context.Context, groupID, name string, updateRequest *DataLakePipeline) (*DataLakePipeline, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(dataLakesPipelineBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, name)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(DataLakePipeline)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes one Data Lake Pipeline.
//
// See more: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/deletePipeline
func (s *DataLakePipelineServiceOp) Delete(ctx context.Context, groupID, name string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if name == "" {
		return nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(dataLakesPipelineBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, name)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
