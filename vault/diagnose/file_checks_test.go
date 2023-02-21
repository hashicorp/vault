package diagnose

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestRaftFolderPerms(t *testing.T) {
	// Make sure overpermissive permissions are caught
	err := os.Mkdir("diagnose", 0o777)
	if err != nil {
		t.Fatal(err)
	}

	info, _ := os.Stat("diagnose")

	if !IsDir(info) {
		t.Fatal("directory was reported to not be a directory")
	}

	// Create a boltDB formatted file and make sure isDB returns true
	fullDBPath := "diagnose/" + DatabaseFilename
	_, err = os.Create(fullDBPath)
	if err != nil {
		t.Fatal(err)
	}
	if !HasDB(fullDBPath) {
		t.Fatal("well-formatted database path is not accepted by DB check function")
	}

	hasOnlyOwnerRW, errs := CheckFilePerms(info)
	if hasOnlyOwnerRW {
		t.Fatal("folder has more than owner rw")
	}
	if len(errs) != 1 && !strings.Contains(errs[0], FileTooPermissiveWarning) {
		t.Fatalf("wrong error or number of errors or wrong error returned: %v", errs)
	}

	// Make sure underpermissiveness is caught
	err = os.Chmod("diagnose", 0o100)
	if err != nil {
		t.Fatal(err)
	}
	info, _ = os.Stat("diagnose")
	hasOnlyOwnerRW, errs = CheckFilePerms(info)
	if hasOnlyOwnerRW {
		t.Fatal("folder should not have owner write")
	}
	if len(errs) != 1 || !strings.Contains(errs[0], FilePermissionsMissingWarning) {
		t.Fatalf("wrong error or number of errors returned: %v", errs)
	}

	// Make sure actually setting owner rw returns properly
	err = os.Chmod("diagnose", 0o600)
	if err != nil {
		t.Fatal(err)
	}
	info, _ = os.Stat("diagnose")
	hasOnlyOwnerRW, errs = CheckFilePerms(info)
	if errs != nil || !hasOnlyOwnerRW {
		t.Fatal("folder with correct perms returns error")
	}

	// Make sure we can clean up the diagnose folder
	os.Chmod("diagnose", 0o777)

	// Clean up test diagnose folder
	err = os.RemoveAll("diagnose")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRaftStorageQuorum(t *testing.T) {
	m := mockStorageBackend{}
	m.raftServerQuorumType = 0
	twoVoterCluster := RaftStorageQuorum(context.Background(), m)

	if !strings.Contains(twoVoterCluster, "Please ensure that Vault has access to an odd number of voter nodes.") {
		t.Fatalf("two voter cluster yielded wrong error: %+s", twoVoterCluster)
	}

	m.raftServerQuorumType = 1
	threeVoterCluster := RaftStorageQuorum(context.Background(), m)
	if !strings.Contains(threeVoterCluster, "Voter quorum exists") {
		t.Fatalf("three voter cluster yielded incorrect error: %s", threeVoterCluster)
	}

	m.raftServerQuorumType = 2
	threeNodeTwoVoterCluster := RaftStorageQuorum(context.Background(), m)
	if !strings.Contains(threeNodeTwoVoterCluster, "Please ensure that Vault has access to an odd number of voter nodes.") {
		t.Fatalf("two voter cluster yielded wrong error: %+s", threeNodeTwoVoterCluster)
	}

	m.raftServerQuorumType = 3
	errClusterInfo := RaftStorageQuorum(context.Background(), m)
	if !strings.Contains(errClusterInfo, "error") {
		t.Fatalf("two voter cluster yielded wrong error: %+s", errClusterInfo)
	}
}
