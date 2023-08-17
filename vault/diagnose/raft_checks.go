// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package diagnose

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/hashicorp/vault/physical/raft"
)

const (
	DatabaseFilename = "vault.db"
	owner            = "owner"
	group            = "group"
	other            = "other"
)

const (
	ownershipTestName   = "Check Raft Folder Ownership"
	permissionsTestName = "Check Raft Folder Permissions"
	raftQuorumTestName  = "Check For Raft Quorum"
)

func RaftFileChecks(ctx context.Context, path string) {
	// Note: Stat does not return information about the symlink itself, in the case where we are dealing with one.
	info, err := os.Stat(path)
	if err != nil {
		SpotError(ctx, permissionsTestName, fmt.Errorf("Error computing file permissions: %w.", err))
	}

	if !IsDir(info) {
		SpotError(ctx, ownershipTestName, fmt.Errorf("Error: Raft storage path variable does not point to a folder."))
	}

	if !HasDB(path) {
		SpotWarn(ctx, ownershipTestName, "Raft boltDB file has not been created.")
	}

	hasOnlyOwnerRW, errs := CheckFilePerms(info)
	for _, err := range errs {
		switch {
		case strings.Contains(err, FileIsSymlinkWarning) || strings.Contains(err, FileTooPermissiveWarning):
			SpotWarn(ctx, permissionsTestName, err)
		case strings.Contains(err, FilePermissionsMissingWarning):
			SpotError(ctx, permissionsTestName, errors.New(err))
		}
	}
	ownedByRoot := IsOwnedByRoot(info)
	requiresRoot := ownedByRoot && hasOnlyOwnerRW
	if requiresRoot {
		SpotWarn(ctx, ownershipTestName, "Raft backend files are owned by root and are only accessible as root or with overpermissive file permissions. This prevents Vault from running as a non-privileged user.")
		Advise(ctx, "Please change raft path permissions to allow for non-root access.")
	}

	if runtime.GOOS == "windows" {
		SpotWarn(ctx, permissionsTestName, "Diagnose cannot determine if Vault needs to run as root to open boltDB file. Please check these permissions manually.")
	} else if errs == nil && !requiresRoot {
		SpotOk(ctx, permissionsTestName, "Raft BoltDB file has correct set of permissions.")
	}
}

// RaftStorageQuorum checks that there is an odd number of voters present
// It returns the status message for testing purposes
func RaftStorageQuorum(ctx context.Context, b RaftConfigurableStorageBackend) string {
	var conf *raft.RaftConfigurationResponse
	var err error
	conf, err = b.GetConfigurationOffline()
	if err != nil {
		SpotError(ctx, raftQuorumTestName, fmt.Errorf("Error retrieving server configuration: %w.", err))
		return fmt.Sprintf("Error retrieving server configuration: %s.", err.Error())
	}
	voterCount := 0
	for _, s := range conf.Servers {
		if s.Voter {
			voterCount++
		}
	}
	if voterCount == 1 {
		nonHAWarning := "Only one server node found. Vault is not running in high availability mode."
		SpotWarn(ctx, raftQuorumTestName, nonHAWarning)
		return nonHAWarning
	}
	var warnMsg string
	if voterCount%2 == 0 {
		warnMsg = fmt.Sprintf("%d voters found. Please ensure that Vault has access to an odd number of voter nodes.", voterCount)
		SpotWarn(ctx, raftQuorumTestName, warnMsg)
		return warnMsg
	}
	if voterCount > 7 {
		warnMsg = fmt.Sprintf("Very large cluster detected: %d voters found. ", voterCount)
		SpotWarn(ctx, raftQuorumTestName, warnMsg)
		return warnMsg
	}

	okMsg := fmt.Sprintf("Voter quorum exists: %d voters.", voterCount)
	SpotOk(ctx, raftQuorumTestName, okMsg)
	return okMsg
}
