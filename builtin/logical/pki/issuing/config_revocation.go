// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	StorageLocalCRLConfig   = "crls/config"
	StorageUnifiedCRLConfig = "unified-crls/config"
)

type InternalCRLConfigEntry struct {
	IssuerIDCRLMap        map[IssuerID]CrlID  `json:"issuer_id_crl_map"`
	CRLNumberMap          map[CrlID]int64     `json:"crl_number_map"`
	LastCompleteNumberMap map[CrlID]int64     `json:"last_complete_number_map"`
	CRLExpirationMap      map[CrlID]time.Time `json:"crl_expiration_map"`
	LastModified          time.Time           `json:"last_modified"`
	DeltaLastModified     time.Time           `json:"delta_last_modified"`
	UseGlobalQueue        bool                `json:"cross_cluster_revocation"`
}

type CrlID string

func (p CrlID) String() string {
	return string(p)
}

func GetLocalCRLConfig(ctx context.Context, s logical.Storage) (*InternalCRLConfigEntry, error) {
	return _getInternalCRLConfig(ctx, s, StorageLocalCRLConfig)
}

func GetUnifiedCRLConfig(ctx context.Context, s logical.Storage) (*InternalCRLConfigEntry, error) {
	return _getInternalCRLConfig(ctx, s, StorageUnifiedCRLConfig)
}

func _getInternalCRLConfig(ctx context.Context, s logical.Storage, path string) (*InternalCRLConfigEntry, error) {
	entry, err := s.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	mapping := &InternalCRLConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(mapping); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode cluster-local CRL configuration: %v", err)}
		}
	}

	if len(mapping.IssuerIDCRLMap) == 0 {
		mapping.IssuerIDCRLMap = make(map[IssuerID]CrlID)
	}

	if len(mapping.CRLNumberMap) == 0 {
		mapping.CRLNumberMap = make(map[CrlID]int64)
	}

	if len(mapping.LastCompleteNumberMap) == 0 {
		mapping.LastCompleteNumberMap = make(map[CrlID]int64)

		// Since this might not exist on migration, we want to guess as
		// to the last full CRL number was. This was likely the last
		// Value from CRLNumberMap if it existed, since we're just adding
		// the mapping here in this block.
		//
		// After the next full CRL build, we will have set this Value
		// correctly, so it doesn't really matter in the long term if
		// we're off here.
		for id, number := range mapping.CRLNumberMap {
			// Decrement by one, since CRLNumberMap is the future number,
			// not the last built number.
			mapping.LastCompleteNumberMap[id] = number - 1
		}
	}

	if len(mapping.CRLExpirationMap) == 0 {
		mapping.CRLExpirationMap = make(map[CrlID]time.Time)
	}

	return mapping, nil
}

func SetLocalCRLConfig(ctx context.Context, s logical.Storage, mapping *InternalCRLConfigEntry) error {
	return _setInternalCRLConfig(ctx, s, mapping, StorageLocalCRLConfig)
}

func SetUnifiedCRLConfig(ctx context.Context, s logical.Storage, mapping *InternalCRLConfigEntry) error {
	return _setInternalCRLConfig(ctx, s, mapping, StorageUnifiedCRLConfig)
}

func _setInternalCRLConfig(ctx context.Context, s logical.Storage, mapping *InternalCRLConfigEntry, path string) error {
	if err := _cleanupInternalCRLMapping(ctx, s, mapping, path); err != nil {
		return fmt.Errorf("failed to clean up internal CRL mapping: %w", err)
	}

	json, err := logical.StorageEntryJSON(path, mapping)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func _cleanupInternalCRLMapping(ctx context.Context, s logical.Storage, mapping *InternalCRLConfigEntry, path string) error {
	// Track which CRL IDs are presently referred to by issuers; any other CRL
	// IDs are subject to cleanup.
	//
	// Unused IDs both need to be removed from this map (cleaning up the size
	// of this storage entry) but also the full CRLs removed from disk.
	presentMap := make(map[CrlID]bool)
	for _, id := range mapping.IssuerIDCRLMap {
		presentMap[id] = true
	}

	// Identify which CRL IDs exist and are candidates for removal;
	// theoretically these three maps should be in sync, but were added
	// at different times.
	toRemove := make(map[CrlID]bool)
	for id := range mapping.CRLNumberMap {
		if !presentMap[id] {
			toRemove[id] = true
		}
	}
	for id := range mapping.LastCompleteNumberMap {
		if !presentMap[id] {
			toRemove[id] = true
		}
	}
	for id := range mapping.CRLExpirationMap {
		if !presentMap[id] {
			toRemove[id] = true
		}
	}

	// Depending on which path we're writing this config to, we need to
	// remove CRLs from the relevant folder too.
	isLocal := path == StorageLocalCRLConfig
	baseCRLPath := PathCrls
	if !isLocal {
		baseCRLPath = "unified-crls/"
	}

	for id := range toRemove {
		// Clean up space in this mapping...
		delete(mapping.CRLNumberMap, id)
		delete(mapping.LastCompleteNumberMap, id)
		delete(mapping.CRLExpirationMap, id)

		// And clean up space on disk from the fat CRL mapping.
		crlPath := baseCRLPath + string(id)
		deltaCRLPath := crlPath + "-delta"
		if err := s.Delete(ctx, crlPath); err != nil {
			return fmt.Errorf("failed to delete unreferenced CRL %v: %w", id, err)
		}
		if err := s.Delete(ctx, deltaCRLPath); err != nil {
			return fmt.Errorf("failed to delete unreferenced delta CRL %v: %w", id, err)
		}
	}

	// Lastly, some CRLs could've been partially removed from the map but
	// not from disk. Check to see if we have any dangling CRLs and remove
	// them too.
	list, err := s.List(ctx, baseCRLPath)
	if err != nil {
		return fmt.Errorf("failed listing all CRLs: %w", err)
	}
	for _, crl := range list {
		if crl == "config" || strings.HasSuffix(crl, "/") {
			continue
		}

		if presentMap[CrlID(crl)] {
			continue
		}

		if err := s.Delete(ctx, baseCRLPath+"/"+crl); err != nil {
			return fmt.Errorf("failed cleaning up orphaned CRL %v: %w", crl, err)
		}
	}

	return nil
}
