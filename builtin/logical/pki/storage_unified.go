// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	unifiedRevocationReadPathPrefix  = "unified-revocation/"
	unifiedRevocationWritePathPrefix = unifiedRevocationReadPathPrefix + "{{clusterId}}/"
)

type unifiedRevocationEntry struct {
	SerialNumber      string    `json:"-"`
	CertExpiration    time.Time `json:"certificate_expiration_utc"`
	RevocationTimeUTC time.Time `json:"revocation_time_utc"`
	CertificateIssuer issuerID  `json:"issuer_id"`
}

func getUnifiedRevocationBySerial(sc *storageContext, serial string) (*unifiedRevocationEntry, error) {
	clusterPaths, err := lookupUnifiedClusterPaths(sc)
	if err != nil {
		return nil, err
	}

	for _, path := range clusterPaths {
		serialPath := path + serial
		entryRaw, err := sc.Storage.Get(sc.Context, serialPath)
		if err != nil {
			return nil, err
		}

		if entryRaw != nil {
			var revEntry unifiedRevocationEntry
			if err := entryRaw.DecodeJSON(&revEntry); err != nil {
				return nil, fmt.Errorf("failed json decoding of unified entry at path %s: %w", serialPath, err)
			}
			revEntry.SerialNumber = serial
			return &revEntry, nil
		}
	}

	return nil, nil
}

func writeUnifiedRevocationEntry(sc *storageContext, ure *unifiedRevocationEntry) error {
	json, err := logical.StorageEntryJSON(unifiedRevocationWritePathPrefix+normalizeSerial(ure.SerialNumber), ure)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, json)
}

// listClusterSpecificUnifiedRevokedCerts returns a list of revoked certificates from a given cluster
func listClusterSpecificUnifiedRevokedCerts(sc *storageContext, clusterId string) ([]string, error) {
	path := unifiedRevocationReadPathPrefix + clusterId + "/"
	serials, err := sc.Storage.List(sc.Context, path)
	if err != nil {
		return nil, err
	}

	return serials, nil
}

// lookupUnifiedClusterPaths returns a map of cluster id to the prefix storage path for that given cluster's
// unified revoked certificates
func lookupUnifiedClusterPaths(sc *storageContext) (map[string]string, error) {
	fullPaths := map[string]string{}

	clusterPaths, err := sc.Storage.List(sc.Context, unifiedRevocationReadPathPrefix)
	if err != nil {
		return nil, err
	}

	for _, clusterIdWithSlash := range clusterPaths {
		// Only include folder listings, if a file were to be stored under this path ignore it.
		if strings.HasSuffix(clusterIdWithSlash, "/") {
			clusterId := clusterIdWithSlash[:len(clusterIdWithSlash)-1] // remove trailing /
			fullPaths[clusterId] = unifiedRevocationReadPathPrefix + clusterIdWithSlash
		}
	}

	return fullPaths, nil
}
