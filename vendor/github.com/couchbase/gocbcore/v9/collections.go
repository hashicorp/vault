package gocbcore

import (
	"encoding/json"
	"strconv"
)

const (
	unknownCid = uint32(0xFFFFFFFF)
	pendingCid = uint32(0xFFFFFFFE)
)

// ManifestCollection is the representation of a collection within a manifest.
type ManifestCollection struct {
	UID  uint32
	Name string
}

// UnmarshalJSON is a custom implementation of json unmarshaling.
func (item *ManifestCollection) UnmarshalJSON(data []byte) error {
	decData := struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(data, &decData); err != nil {
		return err
	}

	decUID, err := strconv.ParseUint(decData.UID, 16, 32)
	if err != nil {
		return err
	}

	item.UID = uint32(decUID)
	item.Name = decData.Name
	return nil
}

// ManifestScope is the representation of a scope within a manifest.
type ManifestScope struct {
	UID         uint32
	Name        string
	Collections []ManifestCollection
}

// UnmarshalJSON is a custom implementation of json unmarshaling.
func (item *ManifestScope) UnmarshalJSON(data []byte) error {
	decData := struct {
		UID         string               `json:"uid"`
		Name        string               `json:"name"`
		Collections []ManifestCollection `json:"collections"`
	}{}
	if err := json.Unmarshal(data, &decData); err != nil {
		return err
	}

	decUID, err := strconv.ParseUint(decData.UID, 16, 32)
	if err != nil {
		return err
	}

	item.UID = uint32(decUID)
	item.Name = decData.Name
	item.Collections = decData.Collections
	return nil
}

// Manifest is the representation of a collections manifest.
type Manifest struct {
	UID    uint64
	Scopes []ManifestScope
}

// UnmarshalJSON is a custom implementation of json unmarshaling.
func (item *Manifest) UnmarshalJSON(data []byte) error {
	decData := struct {
		UID    string          `json:"uid"`
		Scopes []ManifestScope `json:"scopes"`
	}{}
	if err := json.Unmarshal(data, &decData); err != nil {
		return err
	}

	decUID, err := strconv.ParseUint(decData.UID, 16, 64)
	if err != nil {
		return err
	}

	item.UID = decUID
	item.Scopes = decData.Scopes
	return nil
}

// GetCollectionManifestOptions are the options available to the GetCollectionManifest command.
type GetCollectionManifestOptions struct {
	// Volatile: Tracer API is subject to change.
	TraceContext  RequestSpanContext
	RetryStrategy RetryStrategy
}

// GetCollectionIDOptions are the options available to the GetCollectionID command.
type GetCollectionIDOptions struct {
	RetryStrategy RetryStrategy
	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// GetCollectionIDResult encapsulates the result of a GetCollectionID operation.
type GetCollectionIDResult struct {
	ManifestID   uint64
	CollectionID uint32
}

// GetCollectionManifestResult encapsulates the result of a GetCollectionManifest operation.
type GetCollectionManifestResult struct {
	Manifest []byte
}
