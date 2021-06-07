package diagnose

import (
	"context"
	"strings"
	"testing"
)

func TestRaftStorageQuorum(t *testing.T) {
	m := mockStorageBackend{}
	m.raftServerQuorumType = 0
	twoVoterCluster := RaftStorageQuorum(context.Background(), m)

	if !strings.Contains(twoVoterCluster, "error") {
		t.Fatalf("two voter cluster yielded wrong error: %+s", twoVoterCluster)
	}

	m.raftServerQuorumType = 1
	threeVoterCluster := RaftStorageQuorum(context.Background(), m)
	if !strings.Contains(threeVoterCluster, "voter quorum exists") {
		t.Fatalf("three voter cluster yielded incorrect error: %s", threeVoterCluster)
	}

	m.raftServerQuorumType = 2
	threeNodeTwoVoterCluster := RaftStorageQuorum(context.Background(), m)
	if !strings.Contains(threeNodeTwoVoterCluster, "error") {
		t.Fatalf("two voter cluster yielded wrong error: %+s", threeNodeTwoVoterCluster)
	}

	m.raftServerQuorumType = 3
	errClusterInfo := RaftStorageQuorum(context.Background(), m)
	if !strings.Contains(errClusterInfo, "error") {
		t.Fatalf("two voter cluster yielded wrong error: %+s", errClusterInfo)
	}
}
