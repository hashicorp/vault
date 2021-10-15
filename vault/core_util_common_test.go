package vault

import (
	"context"
	"testing"
	"time"
)

// TestStoreMultipleVaultVersions writes multiple versions of 1.9.0 and verifies that only
// the original timestamp is stored.
func TestStoreMultipleVaultVersions(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	upgradeTimePlusEpsilon := time.Now()
	wasStored, err := c.StoreVersionTimestamp(context.Background(), "1.9.0", upgradeTimePlusEpsilon.Add(30*time.Hour))
	if err != nil || wasStored {
		t.Fatalf("vault version was re-stored: %v, err is: %s", wasStored, err.Error())
	}
	upgradeTime, ok := c.VersionTimestamps["1.9.0"]
	if !ok {
		t.Fatalf("no 1.9.0 version timestamp found")
	}
	if upgradeTime.After(upgradeTimePlusEpsilon) {
		t.Fatalf("upgrade time for 1.9.0 is incorrect: got %+v, expected less than %+v", upgradeTime, upgradeTimePlusEpsilon)
	}
}

// TestGetOldestVersion verifies that FindOldestVersionTimestamp finds the oldest
// vault version stored.
func TestGetOldestVersion(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	upgradeTimePlusEpsilon := time.Now()
	c.StoreVersionTimestamp(context.Background(), "1.9.1", upgradeTimePlusEpsilon.Add(-4*time.Hour))
	c.StoreVersionTimestamp(context.Background(), "1.9.2", upgradeTimePlusEpsilon.Add(2*time.Hour))
	c.HandleLoadVersionTimestamps(c.activeContext)
	if len(c.VersionTimestamps) != 3 {
		t.Fatalf("expected 3 entries in timestamps map after refresh, found: %d", len(c.VersionTimestamps))
	}
	v, tm, err := c.FindOldestVersionTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if v != "1.9.1" {
		t.Fatalf("expected 1.9.1, found: %s", v)
	}
	if tm.Before(upgradeTimePlusEpsilon.Add(-6*time.Hour)) || tm.After(upgradeTimePlusEpsilon.Add(-2*time.Hour)) {
		t.Fatalf("incorrect upgrade time logged: %v", tm)
	}
}
