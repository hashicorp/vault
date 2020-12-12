package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SnapshotDefinition represents a  snapshot definition
type SnapshotDefinition struct {
	VolumeIDentifier string `json:"volume_id"`
	Name             string `json:"name,omitempty"`
	Organization     string `json:"organization"`
}

// Snapshot represents a  Snapshot
type Snapshot struct {
	// Identifier is a unique identifier for the snapshot
	Identifier string `json:"id,omitempty"`

	// Name is a user-defined name for the snapshot
	Name string `json:"name,omitempty"`

	// CreationDate is the creation date of the snapshot
	CreationDate string `json:"creation_date,omitempty"`

	// ModificationDate is the date of the last modification of the snapshot
	ModificationDate string `json:"modification_date,omitempty"`

	// Size is the allocated size of the volume
	Size uint64 `json:"size,omitempty"`

	// Organization is the owner of the snapshot
	Organization string `json:"organization"`

	// State is the current state of the snapshot
	State string `json:"state"`

	// VolumeType is the kind of volume behind the snapshot
	VolumeType string `json:"volume_type"`

	// BaseVolume is the volume from which the snapshot inherits
	BaseVolume Volume `json:"base_volume,omitempty"`
}

// oneSnapshot represents the response of a GET /snapshots/UUID API call
type oneSnapshot struct {
	Snapshot Snapshot `json:"snapshot,omitempty"`
}

// Snapshots represents a group of  snapshots
type Snapshots struct {
	// Snapshots holds  snapshots of the response
	Snapshots []Snapshot `json:"snapshots,omitempty"`
}

// CreateSnapshot creates a new snapshot
func (s *API) CreateSnapshot(volumeID string, name string) (*Snapshot, error) {
	definition := SnapshotDefinition{
		VolumeIDentifier: volumeID,
		Name:             name,
		Organization:     s.Organization,
	}
	resp, err := s.PostResponse(s.computeAPI, "snapshots", definition)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return nil, err
	}
	var snapshot oneSnapshot

	if err = json.Unmarshal(body, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot.Snapshot, nil
}

// DeleteSnapshot deletes a snapshot
func (s *API) DeleteSnapshot(snapshotID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("snapshots/%s", snapshotID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}

// GetSnapshots gets the list of snapshots from the API
func (s *API) GetSnapshots() ([]Snapshot, error) {
	query := url.Values{}

	resp, err := s.GetResponsePaginate(s.computeAPI, "snapshots", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var snapshots Snapshots

	if err = json.Unmarshal(body, &snapshots); err != nil {
		return nil, err
	}
	return snapshots.Snapshots, nil
}

// GetSnapshot gets a snapshot from the API
func (s *API) GetSnapshot(snapshotID string) (*Snapshot, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "snapshots/"+snapshotID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var oneSnapshot oneSnapshot

	if err = json.Unmarshal(body, &oneSnapshot); err != nil {
		return nil, err
	}
	// FIXME region, arch, owner, title
	return &oneSnapshot.Snapshot, nil
}
