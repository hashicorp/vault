// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/version"
)

// TestVersionStore_StoreMultipleVaultVersions writes multiple versions of 1.9.0 and verifies that only
// the original timestamp is stored.
func TestVersionStore_StoreMultipleVaultVersions(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	upgradeTimePlusEpsilon := time.Now().UTC()
	vaultVersion := &VaultVersion{
		Version:            version.Version,
		TimestampInstalled: upgradeTimePlusEpsilon.Add(30 * time.Hour),
	}
	wasStored, err := c.storeVersionEntry(context.Background(), vaultVersion, false)
	if err != nil || wasStored {
		t.Fatalf("vault version was re-stored: %v, err is: %s", wasStored, err.Error())
	}
	versionEntry, ok := c.versionHistory[version.Version]
	if !ok {
		t.Fatalf("no %s version timestamp found", version.Version)
	}
	if versionEntry.TimestampInstalled.After(upgradeTimePlusEpsilon) {
		t.Fatalf("upgrade time for %s is incorrect: got %+v, expected less than %+v", version.Version, versionEntry.TimestampInstalled, upgradeTimePlusEpsilon)
	}
}

// TestVersionStore_GetOldestVersion verifies that FindOldestVersionTimestamp finds the oldest
// (in time) vault version stored.
func TestVersionStore_GetOldestVersion(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	upgradeTimePlusEpsilon := time.Now().UTC()

	// 1.6.2 is stored before 1.6.1, so even though it is a higher number, it should be returned.
	versionEntries := []VaultVersion{
		{Version: "1.6.2", TimestampInstalled: upgradeTimePlusEpsilon.Add(-4 * time.Hour)},
		{Version: "1.6.1", TimestampInstalled: upgradeTimePlusEpsilon.Add(2 * time.Hour)},
	}

	for _, entry := range versionEntries {
		_, err := c.storeVersionEntry(context.Background(), &entry, false)
		if err != nil {
			t.Fatalf("failed to write version entry %#v, err: %s", entry, err.Error())
		}
	}

	err := c.loadVersionHistory(c.activeContext)
	if err != nil {
		t.Fatalf("failed to populate version history cache, err: %s", err.Error())
	}

	if len(c.versionHistory) != 3 {
		t.Fatalf("expected 3 entries in timestamps map after refresh, found: %d", len(c.versionHistory))
	}
	v, tm, err := c.FindOldestVersionTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if v != "1.6.2" {
		t.Fatalf("expected 1.6.2, found: %s", v)
	}
	if tm.Before(upgradeTimePlusEpsilon.Add(-6*time.Hour)) || tm.After(upgradeTimePlusEpsilon.Add(-2*time.Hour)) {
		t.Fatalf("incorrect upgrade time logged: %v", tm)
	}
}

// TestVersionStore_GetNewestVersion verifies that FindNewestVersionTimestamp finds the newest
// (in time) vault version stored.
func TestVersionStore_GetNewestVersion(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	upgradeTimePlusEpsilon := time.Now().UTC()

	// 1.6.1 is stored after 1.6.2, so even though it is a lower number, it should be returned.
	versionEntries := []VaultVersion{
		{Version: "1.6.2", TimestampInstalled: upgradeTimePlusEpsilon.Add(-4 * time.Hour)},
		{Version: "1.6.1", TimestampInstalled: upgradeTimePlusEpsilon.Add(2 * time.Hour)},
	}

	for _, entry := range versionEntries {
		_, err := c.storeVersionEntry(context.Background(), &entry, false)
		if err != nil {
			t.Fatalf("failed to write version entry %#v, err: %s", entry, err.Error())
		}
	}

	err := c.loadVersionHistory(c.activeContext)
	if err != nil {
		t.Fatalf("failed to populate version history cache, err: %s", err.Error())
	}

	if len(c.versionHistory) != 3 {
		t.Fatalf("expected 3 entries in timestamps map after refresh, found: %d", len(c.versionHistory))
	}
	v, tm, err := c.FindNewestVersionTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if v != "1.6.1" {
		t.Fatalf("expected 1.6.1, found: %s", v)
	}
	if tm.Before(upgradeTimePlusEpsilon.Add(1*time.Hour)) || tm.After(upgradeTimePlusEpsilon.Add(3*time.Hour)) {
		t.Fatalf("incorrect upgrade time logged: %v", tm)
	}
}

func TestVersionStore_SelfHealUTC(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	estLoc, err := time.LoadLocation("EST")
	if err != nil {
		t.Fatalf("failed to load location, err: %s", err.Error())
	}

	nowEST := time.Now().In(estLoc)

	versionEntries := []VaultVersion{
		{Version: "1.9.0", TimestampInstalled: nowEST.Add(24 * time.Hour)},
		{Version: "1.9.1", TimestampInstalled: nowEST.Add(48 * time.Hour)},
	}

	for _, entry := range versionEntries {
		_, err := c.storeVersionEntry(context.Background(), &entry, false)
		if err != nil {
			t.Fatalf("failed to write version entry %#v, err: %s", entry, err.Error())
		}
	}

	err = c.loadVersionHistory(c.activeContext)
	if err != nil {
		t.Fatalf("failed to load version timestamps, err: %s", err.Error())
	}

	for _, entry := range c.versionHistory {
		if entry.TimestampInstalled.Location() != time.UTC {
			t.Fatalf("failed to convert %s timestamp %s to UTC", entry.Version, entry.TimestampInstalled)
		}
	}
}
