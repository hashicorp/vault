// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"crypto/x509"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
)

const (
	minUnifiedTransferDelay = 30 * time.Minute
)

type unifiedTransferStatus struct {
	isRunning  atomic.Bool
	lastRun    time.Time
	forceRerun atomic.Bool
}

func (uts *unifiedTransferStatus) forceRun() {
	uts.forceRerun.Store(true)
}

func newUnifiedTransferStatus() *unifiedTransferStatus {
	return &unifiedTransferStatus{}
}

// runUnifiedTransfer meant to run as a background, this will process all and
// send all missing local revocation entries to the unified space if the feature
// is enabled.
func runUnifiedTransfer(sc *storageContext) {
	b := sc.Backend
	status := b.unifiedTransferStatus

	isPerfStandby := b.System().ReplicationState().HasState(consts.ReplicationDRSecondary | consts.ReplicationPerformanceStandby)

	if isPerfStandby || b.System().LocalMount() {
		// We only do this on active enterprise nodes, when we aren't a local mount
		return
	}

	config, err := b.crlBuilder.getConfigWithUpdate(sc)
	if err != nil {
		b.Logger().Error("failed to retrieve crl config from storage for unified transfer background process",
			"error", err)
		return
	}

	if !status.lastRun.IsZero() {
		// We have run before, we only run again if we have
		// been requested to forceRerun, and we haven't run since our
		// minimum delay
		if !(status.forceRerun.Load() && time.Since(status.lastRun) < minUnifiedTransferDelay) {
			return
		}
	}

	if !config.UnifiedCRL {
		// Feature is disabled, no need to run
		return
	}

	clusterId, err := b.System().ClusterID(sc.Context)
	if err != nil {
		b.Logger().Error("failed to fetch cluster id for unified transfer background process",
			"error", err)
		return
	}

	if !status.isRunning.CompareAndSwap(false, true) {
		b.Logger().Debug("an existing unified transfer process is already running")
		return
	}
	defer status.isRunning.Store(false)

	// Reset our flag before we begin, we do this before we start as
	// we can't guarantee that we can properly parse/fix the error from an
	// error that comes in from the revoke API after that. This will
	// force another run, which worst case, we will fix it on the next
	// periodic function call that passes our min delay.
	status.forceRerun.Store(false)

	err = doUnifiedTransferMissingLocalSerials(sc, clusterId)
	if err != nil {
		b.Logger().Error("an error occurred running unified transfer", "error", err.Error())
		status.forceRerun.Store(true)
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
			config, _ := sc.Backend.crlBuilder.getConfigWithUpdate(sc)
			if config != nil && !config.UnifiedCRL {
				return errors.New("unified crl has been disabled after we started, stopping")
			}
		}
		if _, ok := unifiedCertLookup[serialNum]; !ok {
			err := readRevocationEntryAndTransfer(sc, serialNum)
			if err != nil {
				errCount++
				sc.Backend.Logger().Debug("Failed transferring local revocation to unified space",
					"serial", serialNum, "error", err)
			}
		}
	}

	if errCount > 0 {
		sc.Backend.Logger().Warn(fmt.Sprintf("Failed transfering %d local serials to unified storage", errCount))
	}

	return nil
}

func readRevocationEntryAndTransfer(sc *storageContext, serial string) error {
	hyphenSerial := normalizeSerial(serial)
	revInfo, err := sc.fetchRevocationInfo(hyphenSerial)
	if err != nil {
		return fmt.Errorf("failed loading revocation entry for serial: %s: %w", serial, err)
	}
	if revInfo == nil {
		sc.Backend.Logger().Debug("no certificate revocation entry for serial", "serial", serial)
		return nil
	}
	cert, err := x509.ParseCertificate(revInfo.CertificateBytes)
	if err != nil {
		sc.Backend.Logger().Debug("failed parsing certificate stored in revocation entry for serial",
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

	entry := &unifiedRevocationEntry{
		SerialNumber:      hyphenSerial,
		CertExpiration:    cert.NotAfter,
		RevocationTimeUTC: revocationTime,
		CertificateIssuer: revInfo.CertificateIssuer,
	}

	return writeUnifiedRevocationEntry(sc, entry)
}
