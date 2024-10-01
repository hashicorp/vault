// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
)

const (
	unifiedRevocationReadPathPrefix = revocation.UnifiedRevocationReadPathPrefix
)

func getUnifiedRevocationBySerial(sc *storageContext, serial string) (*revocation.UnifiedRevocationEntry, error) {
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
			var revEntry revocation.UnifiedRevocationEntry
			if err := entryRaw.DecodeJSON(&revEntry); err != nil {
				return nil, fmt.Errorf("failed json decoding of unified entry at path %s: %w", serialPath, err)
			}
			revEntry.SerialNumber = serial
			return &revEntry, nil
		}
	}

	return nil, nil
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
