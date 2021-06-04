package diagnose

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/physical/raft"
)

// RaftFilePermsChecks inputs the raft file path and outputs:
// 1. A bool, which dictates whether Owner has rw permissions
// 2. A list of errors
func RaftFilePermsChecks(ctx context.Context, path string) {
	owner := "owner"
	group := "group"
	other := "other"
	var overPermissions []string

	// Note: Stat does not return information about the symlink itself, in the case where we are dealing with one.
	info, err := os.Stat(path)
	if err != nil {
		SpotError(ctx, "raft file permission checks", fmt.Errorf("error computing file permissions: %w", err))
	}

	mode := info.Mode()

	canRead := false
	canWrite := false
	for i := 1; i < 4; i++ {
		perm := string(mode.String()[i])
		fmt.Println("hi")
		fmt.Printf("%s", perm)
		if perm != "w" && perm != "r" && perm != "-" {
			overPermissions = append(overPermissions, owner+":"+perm)
		} else if perm == "w" {
			canWrite = true
		} else if perm == "r" {
			canRead = true
		}
	}

	for i := 4; i < 7; i++ {
		perm := string(mode.String()[i])
		fmt.Printf("%s", perm)
		if perm != "-" {
			overPermissions = append(overPermissions, group+":"+perm)
		}
	}

	for i := 7; i < 10; i++ {
		perm := string(mode.String()[i])
		fmt.Printf("%s", perm)
		if perm != "-" {
			overPermissions = append(overPermissions, other+":"+perm)
		}
	}

	okFlag := true
	if mode&os.ModeSymlink != 0 {
		okFlag = false
		SpotWarn(ctx, "raft file permission checks", "raft storage backend file is a symlink")
	}
	if len(overPermissions) > 0 {
		okFlag = false
		SpotWarn(ctx, "raft file permission checks", fmt.Sprintf("too many permissions -- %s", strings.Join(overPermissions, ",")))
	}
	if !canRead || !canWrite {
		okFlag = false
		SpotError(ctx, "raft file permission checks", fmt.Errorf("owner read: %v, owner write: %v", canRead, canWrite))
	}
	if okFlag {
		SpotOk(ctx, "raft file permission checks", "correct permissions to raft storage backend file path")
	}
}

// RaftStorageQuorum checks that there is an odd number of voters present
// It returns the status message for testing purposes
func RaftStorageQuorum(ctx context.Context, b *raft.RaftBackend) string {
	var conf *raft.RaftConfigurationResponse
	var err error
	conf, err = b.GetConfiguration(ctx)
	if err != nil {
		SpotError(ctx, "raft quorum", fmt.Errorf("error retrieving server configuration: %w", err))
		return "error retrieving server configuration"
	}
	voterCount := 0
	for _, s := range conf.Servers {
		if s.Voter {
			voterCount++
		}
	}
	if voterCount == 1 || voterCount == 3 || voterCount == 5 || voterCount == 7 {
		okMsg := fmt.Sprintf("voter quorum exists: %d voters", voterCount)
		SpotOk(ctx, "raft quorum", okMsg)
		return okMsg
	}
	warnMsg := fmt.Sprintf("even number of voters found: %d", voterCount)
	SpotWarn(ctx, "raft quorum", warnMsg)
	return warnMsg
}
