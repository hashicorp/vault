// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"
	"strings"
	"time"
)

func (sc *storageContext) isDefaultKeySet() (bool, error) {
	config, err := sc.getKeysConfig()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultKeyId.String()) != "", nil
}

func (sc *storageContext) isDefaultIssuerSet() (bool, error) {
	config, err := sc.getIssuersConfig()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultIssuerId.String()) != "", nil
}

func (sc *storageContext) updateDefaultKeyId(id keyID) error {
	config, err := sc.getKeysConfig()
	if err != nil {
		return err
	}

	if config.DefaultKeyId != id {
		return sc.setKeysConfig(&keyConfigEntry{
			DefaultKeyId: id,
		})
	}

	return nil
}

func (sc *storageContext) updateDefaultIssuerId(id issuerID) error {
	config, err := sc.getIssuersConfig()
	if err != nil {
		return err
	}

	if config.DefaultIssuerId != id {
		config.DefaultIssuerId = id
		return sc.setIssuersConfig(config)
	}

	return nil
}

func (sc *storageContext) changeDefaultIssuerTimestamps(oldDefault issuerID, newDefault issuerID) error {
	if newDefault == oldDefault {
		return nil
	}

	now := time.Now().UTC()

	// When the default issuer changes, we need to modify four
	// pieces of information:
	//
	// 1. The old default issuer's modification time, as it no
	//    longer works for the /cert/ca path.
	// 2. The new default issuer's modification time, as it now
	//    works for the /cert/ca path.
	// 3. & 4. Both issuer's CRLs, as they behave the same, under
	//    the /cert/crl path!
	for _, thisId := range []issuerID{oldDefault, newDefault} {
		if len(thisId) == 0 {
			continue
		}

		// 1 & 2 above.
		issuer, err := sc.fetchIssuerById(thisId)
		if err != nil {
			// Due to the lack of transactions, if we deleted the default
			// issuer (successfully), but the subsequent issuer config write
			// (to clear the default issuer's old id) failed, we might have
			// an inconsistent config. If we later hit this loop (and flush
			// these timestamps again -- perhaps because the operator
			// selected a new default), we'd have erred out here, because
			// the since-deleted default issuer doesn't exist. In this case,
			// skip the issuer instead of bailing.
			err := fmt.Errorf("unable to update issuer (%v)'s modification time: error fetching issuer: %w", thisId, err)
			if strings.Contains(err.Error(), "does not exist") {
				sc.Backend.Logger().Warn(err.Error())
				continue
			}

			return err
		}

		issuer.LastModified = now
		err = sc.writeIssuer(issuer)
		if err != nil {
			return fmt.Errorf("unable to update issuer (%v)'s modification time: error persisting issuer: %w", thisId, err)
		}
	}

	// Fetch and update the internalCRLConfigEntry (3&4).
	cfg, err := sc.getLocalCRLConfig()
	if err != nil {
		return fmt.Errorf("unable to update local CRL config's modification time: error fetching local CRL config: %w", err)
	}

	cfg.LastModified = now
	cfg.DeltaLastModified = now
	err = sc.setLocalCRLConfig(cfg)
	if err != nil {
		return fmt.Errorf("unable to update local CRL config's modification time: error persisting local CRL config: %w", err)
	}

	return nil
}
