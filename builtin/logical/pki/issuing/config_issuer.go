// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const StorageIssuerConfig = "config/issuers"

type IssuerConfigEntry struct {
	// This new fetchedDefault field allows us to detect if the default
	// issuer was modified, in turn dispatching the timestamp updater
	// if necessary.
	fetchedDefault             IssuerID `json:"-"`
	DefaultIssuerId            IssuerID `json:"default"`
	DefaultFollowsLatestIssuer bool     `json:"default_follows_latest_issuer"`
}

func GetIssuersConfig(ctx context.Context, s logical.Storage) (*IssuerConfigEntry, error) {
	entry, err := s.Get(ctx, StorageIssuerConfig)
	if err != nil {
		return nil, err
	}

	issuerConfig := &IssuerConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(issuerConfig); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode issuer configuration: %v", err)}
		}
	}
	issuerConfig.fetchedDefault = issuerConfig.DefaultIssuerId

	return issuerConfig, nil
}

func SetIssuersConfig(ctx context.Context, s logical.Storage, config *IssuerConfigEntry) error {
	json, err := logical.StorageEntryJSON(StorageIssuerConfig, config)
	if err != nil {
		return err
	}

	if err := s.Put(ctx, json); err != nil {
		return err
	}

	if err := changeDefaultIssuerTimestamps(ctx, s, config.fetchedDefault, config.DefaultIssuerId); err != nil {
		return err
	}

	return nil
}

func changeDefaultIssuerTimestamps(ctx context.Context, s logical.Storage, oldDefault IssuerID, newDefault IssuerID) error {
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
	for _, thisId := range []IssuerID{oldDefault, newDefault} {
		if len(thisId) == 0 {
			continue
		}

		// 1 & 2 above.
		issuer, err := FetchIssuerById(ctx, s, thisId)
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
				hclog.L().Warn(err.Error())
				continue
			}

			return err
		}

		issuer.LastModified = now
		err = WriteIssuer(ctx, s, issuer)
		if err != nil {
			return fmt.Errorf("unable to update issuer (%v)'s modification time: error persisting issuer: %w", thisId, err)
		}
	}

	// Fetch and update the internalCRLConfigEntry (3&4).
	cfg, err := GetLocalCRLConfig(ctx, s)
	if err != nil {
		return fmt.Errorf("unable to update local CRL config's modification time: error fetching local CRL config: %w", err)
	}

	cfg.LastModified = now
	cfg.DeltaLastModified = now
	err = SetLocalCRLConfig(ctx, s, cfg)
	if err != nil {
		return fmt.Errorf("unable to update local CRL config's modification time: error persisting local CRL config: %w", err)
	}

	return nil
}
