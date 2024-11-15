package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Volume represents a  Volume
type Volume struct {
	// Identifier is a unique identifier for the volume
	Identifier string `json:"id,omitempty"`

	// Size is the allocated size of the volume
	Size uint64 `json:"size,omitempty"`

	// CreationDate is the creation date of the volume
	CreationDate string `json:"creation_date,omitempty"`

	// ModificationDate is the date of the last modification of the volume
	ModificationDate string `json:"modification_date,omitempty"`

	// Organization is the organization owning the volume
	Organization string `json:"organization,omitempty"`

	// Name is the name of the volume
	Name string `json:"name,omitempty"`

	// Server is the server using this image
	Server *struct {
		Identifier string `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"server,omitempty"`

	// VolumeType is a  identifier for the kind of volume (default: l_ssd)
	VolumeType string `json:"volume_type,omitempty"`

	// ExportURI represents the url used by initrd/scripts to attach the volume
	ExportURI string `json:"export_uri,omitempty"`
}

type volumeResponse struct {
	Volume Volume `json:"volume,omitempty"`
}

// VolumeDefinition represents a  volume definition
type VolumeDefinition struct {
	// Name is the user-defined name of the volume
	Name string `json:"name"`

	// Image is the image used by the volume
	Size uint64 `json:"size"`

	// Bootscript is the bootscript used by the volume
	Type string `json:"volume_type"`

	// Organization is the owner of the volume
	Organization string `json:"organization"`
}

// VolumePutDefinition represents a  volume with nullable fields (for PUT)
type VolumePutDefinition struct {
	Identifier       *string `json:"id,omitempty"`
	Size             *uint64 `json:"size,omitempty"`
	CreationDate     *string `json:"creation_date,omitempty"`
	ModificationDate *string `json:"modification_date,omitempty"`
	Organization     *string `json:"organization,omitempty"`
	Name             *string `json:"name,omitempty"`
	Server           struct {
		Identifier *string `json:"id,omitempty"`
		Name       *string `json:"name,omitempty"`
	} `json:"server,omitempty"`
	VolumeType *string `json:"volume_type,omitempty"`
	ExportURI  *string `json:"export_uri,omitempty"`
}

// CreateVolume creates a new volume
func (s *API) CreateVolume(definition VolumeDefinition) (*Volume, error) {
	definition.Organization = s.Organization
	if definition.Type == "" {
		definition.Type = "l_ssd"
	}

	resp, err := s.PostResponse(s.computeAPI, "volumes", definition)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return nil, err
	}
	var volume volumeResponse

	if err = json.Unmarshal(body, &volume); err != nil {
		return nil, err
	}
	return &volume.Volume, nil
}

// UpdateVolume updates a volume
func (s *API) UpdateVolume(volumeID string, definition VolumePutDefinition) (*Volume, error) {
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("volumes/%s", volumeID), definition)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var volume volumeResponse

	if err = json.Unmarshal(body, &volume); err != nil {
		return nil, err
	}
	return &volume.Volume, nil
}

// DeleteVolume deletes a volume
func (s *API) DeleteVolume(volumeID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("volumes/%s", volumeID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}

type volumesResponse struct {
	Volumes []Volume `json:"volumes,omitempty"`
}

// GetVolumes gets the list of volumes from the API
func (s *API) GetVolumes() (*[]Volume, error) {
	query := url.Values{}

	resp, err := s.GetResponsePaginate(s.computeAPI, "volumes", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}

	var volumes volumesResponse

	if err = json.Unmarshal(body, &volumes); err != nil {
		return nil, err
	}
	return &volumes.Volumes, nil
}

// GetVolume gets a volume from the API
func (s *API) GetVolume(volumeID string) (*Volume, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "volumes/"+volumeID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var volume volumeResponse

	if err = json.Unmarshal(body, &volume); err != nil {
		return nil, err
	}
	// FIXME region, arch, owner, title
	return &volume.Volume, nil
}
