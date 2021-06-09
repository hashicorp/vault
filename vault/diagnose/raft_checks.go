package diagnose

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/physical/raft"
)

const (
	FileIsSymlinkWarning          = "raft storage backend file is a symlink"
	FileTooPermissiveWarning      = "too many permissions"
	FilePermissionsMissingWarning = "owner needs read and write permissions"
)

const DatabaseFilename = "vault.db"
const owner = "owner"
const group = "group"
const other = "other"

func RaftFileChecks(ctx context.Context, path string) {

	// Note: Stat does not return information about the symlink itself, in the case where we are dealing with one.
	info, err := os.Stat(path)
	if err != nil {
		SpotError(ctx, "raft folder permission checks", fmt.Errorf("error computing file permissions: %w", err))
	}

	if !IsDir(info) {
		SpotError(ctx, "raft folder ownership checks", fmt.Errorf("error: path does not point to folder"))
	}

	if !HasDB(path) {
		SpotWarn(ctx, "raft folder ownership checks", "boltDB file has not been created")
	}

	correctPerms, errs := HasCorrectFilePerms(info)
	if errs != nil {
		for _, err := range errs {
			switch {
			case strings.Contains(err, FileIsSymlinkWarning) || strings.Contains(err, FileTooPermissiveWarning):
				SpotWarn(ctx, "raft folder permission checks", err)
			case strings.Contains(err, FilePermissionsMissingWarning):
				SpotError(ctx, "raft folder permission checks", errors.New(err))
			}
		}
	} else if correctPerms {
		SpotOk(ctx, "raft folder permission checks", "boltDB file has correct set of permissions")
	}

	ownedByRoot, err := IsOwnedByRoot(info)
	if err != nil {
		SpotError(ctx, "raft folder ownership checks", fmt.Errorf("vault could not determine file owner for boltDB storage file"))
	}
	if ownedByRoot {
		SpotWarn(ctx, "raft folder ownership checks", "raft backend files owned by root and only accessible as root or with overpermissive file perms")
		Advise(ctx, "this prevents Vault from running as a non-privileged user")
	}
}

// RaftStorageQuorum checks that there is an odd number of voters present
// It returns the status message for testing purposes
func RaftStorageQuorum(ctx context.Context, b RaftConfigurableStorageBackend) string {
	var conf *raft.RaftConfigurationResponse
	var err error
	conf, err = b.GetConfiguration(ctx)
	if err != nil {
		SpotError(ctx, "raft quorum", fmt.Errorf("error retrieving server configuration: %w", err))
		return fmt.Sprintf("error retrieving server configuration: %s", err.Error())
	}
	voterCount := 0
	for _, s := range conf.Servers {
		if s.Voter {
			voterCount++
		}
	}
	if voterCount == 1 {
		nonHAWarning := "warning: only one server node found. Vault is not running in high availability mode"
		SpotWarn(ctx, "raft quorum", nonHAWarning)
		return nonHAWarning
	}
	if voterCount == 3 || voterCount == 5 || voterCount == 7 {
		okMsg := fmt.Sprintf("voter quorum exists: %d voters", voterCount)
		SpotOk(ctx, "raft quorum", okMsg)
		return okMsg
	}
	warnMsg := fmt.Sprintf("error: even number of voters found: %d", voterCount)
	SpotWarn(ctx, "raft quorum", warnMsg)
	return warnMsg
}
