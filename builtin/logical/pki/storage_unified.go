package pki

import (
	"fmt"
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
	clusterPaths, err := lookupClusterPaths(sc)
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

func listUnifiedRevokedCerts(sc *storageContext) ([]string, error) {
	allSerials := []string{}

	clusterPaths, err := lookupClusterPaths(sc)
	if err != nil {
		return nil, err
	}

	for _, path := range clusterPaths {
		clusterSerials, err := sc.Storage.List(sc.Context, path)
		if err != nil {
			return nil, fmt.Errorf("failed listing revoked certs for path %s: %w", path, err)
		}

		allSerials = append(allSerials, clusterSerials...)
	}
	return allSerials, nil
}

func lookupClusterPaths(sc *storageContext) ([]string, error) {
	fullPaths := []string{}

	clusterPaths, err := sc.Storage.List(sc.Context, unifiedRevocationReadPathPrefix)
	if err != nil {
		return fullPaths, err
	}

	for _, clusterId := range clusterPaths {
		fullPaths = append(fullPaths, unifiedRevocationReadPathPrefix+clusterId)
	}

	return fullPaths, nil
}
