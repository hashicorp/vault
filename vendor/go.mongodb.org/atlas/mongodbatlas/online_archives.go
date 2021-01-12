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
	List(context.Context, string, string) ([]*OnlineArchive, *Response, error)
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
func (s *OnlineArchiveServiceOp) List(ctx context.Context, projectID, clusterName string) ([]*OnlineArchive, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	path := fmt.Sprintf(onlineArchiveBasePath, projectID, clusterName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []*OnlineArchive
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
	ExpireAfterDays float64 `json:"expireAfterDays"`
}

// PartitionFields fields to use to partition data
type PartitionFields struct {
	FieldName string   `json:"fieldName,omitempty"`
	FieldType string   `json:"fieldType,omitempty"`
	Order     *float64 `json:"order,omitempty"`
}
