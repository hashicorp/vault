// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/x509"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	minUnifiedTransferDelay = 30 * time.Minute
)

type UnifiedTransferStatus struct {
	isRunning  atomic.Bool
	lastRun    time.Time
	forceRerun atomic.Bool
}

func (uts *UnifiedTransferStatus) forceRun() {
	uts.forceRerun.Store(true)
}

func newUnifiedTransferStatus() *UnifiedTransferStatus {
	return &UnifiedTransferStatus{}
}

// runUnifiedTransfer meant to run as a background, this will process all and
// send all missing local revocation entries to the unified space if the feature
// is enabled.
func runUnifiedTransfer(sc *storageContext) {
	status := sc.GetUnifiedTransferStatus()

	isPerfStandby := sc.System().ReplicationState().HasState(consts.ReplicationDRSecondary | consts.ReplicationPerformanceStandby)

	if isPerfStandby || sc.System().LocalMount() {
		// We only do this on active enterprise nodes, when we aren't a local mount
		return
	}

	config, err := sc.CrlBuilder().GetConfigWithUpdate(sc)
	if err != nil {
		sc.Logger().Error("failed to retrieve crl config from storage for unified transfer background process",
			"error", err)
		return
	}

	if !config.UnifiedCRL {
		// Feature is disabled, no need to run
		return
	}

	clusterId, err := sc.System().ClusterID(sc.Context)
	if err != nil {
		sc.Logger().Error("failed to fetch cluster id for unified transfer background process",
			"error", err)
		return
	}

	if !status.isRunning.CompareAndSwap(false, true) {
		sc.Logger().Debug("an existing unified transfer process is already running")
		return
	}
	defer status.isRunning.Store(false)

	// Because access to lastRun is not locked, we need to delay this check
	// until after we grab the isRunning CAS lock.
	if !status.lastRun.IsZero() {
		// We have run before, we only run again if we have
		// been requested to forceRerun, and we haven't run since our
		// minimum delay.
		if !(status.forceRerun.Load() && time.Since(status.lastRun) < minUnifiedTransferDelay) {
			return
		}
	}

	// Reset our flag before we begin, we do this before we start as
	// we can't guarantee that we can properly parse/fix the error from an
	// error that comes in from the revoke API after that. This will
	// force another run, which worst case, we will fix it on the next
	// periodic function call that passes our min delay.
	status.forceRerun.Store(false)

	err = doUnifiedTransferMissingLocalSerials(sc, clusterId)
	if err != nil {
		sc.Logger().Error("an error occurred running unified transfer", "error", err.Error())
		status.forceRerun.Store(true)
	} else {
		if config.EnableDelta {
			err = doUnifiedTransferMissingDeltaWALSerials(sc, clusterId)
			if err != nil {
				sc.Logger().Error("an error occurred running unified transfer", "error", err.Error())
				status.forceRerun.Store(true)
			}
		}
	}

	status.lastRun = time.Now()
}

func doUnifiedTransferMissingLocalSerials(sc *storageContext, clusterId string) error {
	localRevokedSerialNums, err := sc.listRevokedCerts()
	if err != nil {
		return err
	}
	if len(localRevokedSerialNums) == 0 {
		// No local certs to transfer, no further work to do.
		return nil
	}

	unifiedSerials, err := listClusterSpecificUnifiedRevokedCerts(sc, clusterId)
	if err != nil {
		return err
	}
	unifiedCertLookup := sliceToMapKey(unifiedSerials)

	errCount := 0
	for i, serialNum := range localRevokedSerialNums {
		if i%25 == 0 {
			config, _ := sc.CrlBuilder().GetConfigWithUpdate(sc)
			if config != nil && !config.UnifiedCRL {
				return errors.New("unified crl has been disabled after we started, stopping")
			}
		}
		if _, ok := unifiedCertLookup[serialNum]; !ok {
			err := readRevocationEntryAndTransfer(sc, serialNum)
			if err != nil {
				errCount++
				sc.Logger().Error("Failed transferring local revocation to unified space",
					"serial", serialNum, "error", err)
			}
		}
	}

	if errCount > 0 {
		sc.Logger().Warn(fmt.Sprintf("Failed transfering %d local serials to unified storage", errCount))
	}

	return nil
}

func doUnifiedTransferMissingDeltaWALSerials(sc *storageContext, clusterId string) error {
	// We need to do a similar thing for Delta WAL entry certificates.
	// When the delta WAL failed to write for one or more entries,
	// we'll need to replicate these up to the primary cluster. When it
	// has performed a new delta WAL build, it will empty storage and
	// update to a last written WAL entry that exceeds what we've seen
	// locally.
	thisUnifiedWALEntryPath := unifiedDeltaWALPath + deltaWALLastRevokedSerialName
	lastUnifiedWALEntry, err := getLastWALSerial(sc, thisUnifiedWALEntryPath)
	if err != nil {
		return fmt.Errorf("failed to fetch last cross-cluster unified revoked delta WAL serial number: %w", err)
	}

	lastLocalWALEntry, err := getLastWALSerial(sc, localDeltaWALLastRevokedSerial)
	if err != nil {
		return fmt.Errorf("failed to fetch last locally revoked delta WAL serial number: %w", err)
	}

	// We now need to transfer all the entries and then write the last WAL
	// entry at the end. Start by listing all certificates; any missing
	// certificates will be copied over and then the WAL entry will be
	// updated once.
	//
	// We do not delete entries either locally or remotely, as either
	// cluster could've rebuilt delta CRLs with out-of-sync information,
	// removing some entries (and, we cannot differentiate between these
	// two cases). On next full CRL rebuild (on either cluster), the state
	// should get synchronized, and future delta CRLs after this function
	// returns without issue will see the remaining entries.
	//
	// Lastly, we need to ensure we don't accidentally write any unified
	// delta WAL entries that aren't present in the main cross-cluster
	// revoked storage location. This would mean the above function failed
	// to copy them for some reason, despite them presumably appearing
	// locally.
	_unifiedWALEntries, err := sc.Storage.List(sc.Context, unifiedDeltaWALPath)
	if err != nil {
		return fmt.Errorf("failed to list cross-cluster unified delta WAL storage: %w", err)
	}
	unifiedWALEntries := sliceToMapKey(_unifiedWALEntries)

	_unifiedRevokedSerials, err := listClusterSpecificUnifiedRevokedCerts(sc, clusterId)
	if err != nil {
		return fmt.Errorf("failed to list cross-cluster revoked certificates: %w", err)
	}
	unifiedRevokedSerials := sliceToMapKey(_unifiedRevokedSerials)

	localWALEntries, err := sc.Storage.List(sc.Context, localDeltaWALPath)
	if err != nil {
		return fmt.Errorf("failed to list local delta WAL storage: %w", err)
	}

	if lastUnifiedWALEntry == lastLocalWALEntry && len(_unifiedWALEntries) == len(localWALEntries) {
		// Writing the last revoked WAL entry is the last thing that we do.
		// Because these entries match (across clusters) and we have the same
		// number of entries, assume we don't have anything to sync and exit
		// early.
		//
		// We need both checks as, in the event of PBPWF failing and then
		// returning while more revocations are happening, we could have
		// been schedule to run, but then skip running (if only the first
		// condition was checked) because a later revocation succeeded
		// in writing a unified WAL entry, before we started replicating
		// the rest back up.
		//
		// The downside of this approach is that, if the main cluster
		// does a full rebuild in the mean time, we could re-sync more
		// entries back up to the primary cluster that are already
		// included in the complete CRL. Users can manually rebuild the
		// full CRL (clearing these duplicate delta CRL entries) if this
		// affects them.
		return nil
	}

	errCount := 0
	for index, serial := range localWALEntries {
		if index%25 == 0 {
			config, _ := sc.CrlBuilder().GetConfigWithUpdate(sc)
			if config != nil && (!config.UnifiedCRL || !config.EnableDelta) {
				return errors.New("unified or delta CRLs have been disabled after we started, stopping")
			}
		}

		if serial == deltaWALLastBuildSerialName || serial == deltaWALLastRevokedSerialName {
			// Skip our special serial numbers.
			continue
		}

		_, isAlreadyPresent := unifiedWALEntries[serial]
		if isAlreadyPresent {
			// Serial exists on both local and unified cluster. We're
			// presuming we don't need to read and re-write these entries
			// and that only missing entries need to be updated.
			continue
		}

		_, isRevokedCopied := unifiedRevokedSerials[serial]
		if !isRevokedCopied {
			// We need to wait here to copy over.
			errCount += 1
			sc.Logger().Debug("Delta WAL exists locally, but corresponding cross-cluster full revocation entry is missing; skipping", "serial", serial)
			continue
		}

		// All good: read the local entry and write to the remote variant.
		localPath := localDeltaWALPath + serial
		unifiedPath := unifiedDeltaWALPath + serial

		entry, err := sc.Storage.Get(sc.Context, localPath)
		if err != nil || entry == nil {
			errCount += 1
			sc.Logger().Error("Failed reading local delta WAL entry to copy to cross-cluster", "serial", serial, "err", err)
			continue
		}

		entry.Key = unifiedPath
		err = sc.Storage.Put(sc.Context, entry)
		if err != nil {
			errCount += 1
			sc.Logger().Error("Failed sync local delta WAL entry to cross-cluster unified delta WAL location", "serial", serial, "err", err)
			continue
		}
	}

	if errCount > 0 {
		// See note above about why we don't fail here.
		sc.Logger().Warn(fmt.Sprintf("Failed transfering %d local delta WAL serials to unified storage", errCount))
		return nil
	}

	// Everything worked. Here, we can write over the delta WAL last revoked
	// value. By using the earlier value, even if new revocations have
	// occurred, we ensure any further missing entries can be handled in the
	// next round.
	lastRevSerial := lastWALInfo{Serial: lastLocalWALEntry}
	lastWALEntry, err := logical.StorageEntryJSON(thisUnifiedWALEntryPath, lastRevSerial)
	if err != nil {
		return fmt.Errorf("unable to create cross-cluster unified last delta CRL WAL entry: %w", err)
	}
	if err = sc.Storage.Put(sc.Context, lastWALEntry); err != nil {
		return fmt.Errorf("error saving cross-cluster unified last delta CRL WAL entry: %w", err)
	}

	return nil
}

func readRevocationEntryAndTransfer(sc *storageContext, serial string) error {
	hyphenSerial := normalizeSerial(serial)
	revInfo, err := fetchRevocationInfo(sc, hyphenSerial)
	if err != nil {
		return fmt.Errorf("failed loading revocation entry for serial: %s: %w", serial, err)
	}
	if revInfo == nil {
		sc.Logger().Debug("no certificate revocation entry for serial", "serial", serial)
		return nil
	}
	cert, err := x509.ParseCertificate(revInfo.CertificateBytes)
	if err != nil {
		sc.Logger().Debug("failed parsing certificate stored in revocation entry for serial",
			"serial", serial, "error", err)
		return nil
	}
	if revInfo.CertificateIssuer == "" {
		// No certificate issuer assigned to this serial yet, just drop it for now,
		// as a crl rebuild/tidy needs to happen
		return nil
	}

	revocationTime := revInfo.RevocationTimeUTC
	if revInfo.RevocationTimeUTC.IsZero() {
		// Legacy revocation entries only had this field and not revocationTimeUTC set...
		revocationTime = time.Unix(revInfo.RevocationTime, 0)
	}

	if time.Now().After(cert.NotAfter) {
		// ignore transferring this entry as it has already expired.
		return nil
	}

	entry := &revocation.UnifiedRevocationEntry{
		SerialNumber:      hyphenSerial,
		CertExpiration:    cert.NotAfter,
		RevocationTimeUTC: revocationTime,
		CertificateIssuer: revInfo.CertificateIssuer,
	}

	return revocation.WriteUnifiedRevocationEntry(sc.GetContext(), sc.GetStorage(), entry)
}
