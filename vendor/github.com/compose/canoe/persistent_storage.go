package canoe

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/snap"
	"github.com/coreos/etcd/wal"
	"github.com/coreos/etcd/wal/walpb"
)

type walMetadata struct {
	NodeID    uint64 `json:"node_id"`
	ClusterID uint64 `json:"cluster_id"`
}

func (rn *Node) initPersistentStorage() error {
	if err := rn.initSnap(); err != nil {
		return errors.Wrap(err, "Error initializing snapshot")
	}

	raftSnap, err := rn.ss.Load()
	if err != nil {
		if err != snap.ErrNoSnapshot && err != snap.ErrEmptySnapshot {
			return errors.Wrap(err, "Error loading latest snapshot")
		}
	}

	var walSnap walpb.Snapshot

	if raftSnap != nil {
		walSnap.Index, walSnap.Term = raftSnap.Metadata.Index, raftSnap.Metadata.Term
	}

	if err := rn.initWAL(walSnap); err != nil {
		return errors.Wrap(err, "Error initializing WAL")
	}

	return nil
}

// Correct order of ops
// 1: Restore Metadata from WAL
// 2: Apply any persisted snapshot to FSM
// 3: Apply any Snapshot to raft storage
// 4: Apply any hardstate to raft storage
// 5: Apply and WAL Entries to raft storage
func (rn *Node) restoreRaft() error {
	raftSnap, err := rn.ss.Load()
	if err != nil {
		if err != snap.ErrNoSnapshot && err != snap.ErrEmptySnapshot {
			return errors.Wrap(err, "Error loading latest snapshot")
		}
	}

	var walSnap walpb.Snapshot

	if raftSnap != nil {
		walSnap.Index, walSnap.Term = raftSnap.Metadata.Index, raftSnap.Metadata.Term
	} else {
		raftSnap = &raftpb.Snapshot{}
	}

	wMetadata, hState, ents, err := rn.wal.ReadAll()
	if err != nil {
		return errors.Wrap(err, "Error reading WAL")
	}

	// NOTE: Step 1
	if err := rn.restoreMetadata(wMetadata); err != nil {
		return errors.Wrap(err, "Error restoring from WAL metadata")
	}

	// We can do this now that we restored the metadata
	if err := rn.attachTransport(); err != nil {
		return errors.Wrap(err, "Error attaching raft Transport layer")
	}

	if err := rn.transport.Start(); err != nil {
		return errors.Wrap(err, "Error starting raft transport layer")
	}

	// NOTE: Step 2
	if err := rn.restoreFSMFromSnapshot(*raftSnap); err != nil {
		return errors.Wrap(err, "Error restoring FSM from snapshot")
	}

	// NOTE: Step 3, 4, 5
	if err := rn.restoreMemoryStorage(*raftSnap, hState, ents); err != nil {
		return errors.Wrap(err, "Error restoring raft memory storage")
	}

	return nil
}

func (rn *Node) initSnap() error {
	if rn.snapDir() == "" {
		return nil
	}

	if err := os.MkdirAll(rn.snapDir(), 0750); err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "Error trying to create directory for snapshots")
	}

	rn.ss = snap.New(rn.snapDir())

	return nil
}

func (rn *Node) persistSnapshot(raftSnap raftpb.Snapshot) error {

	if rn.ss != nil {
		if err := rn.ss.SaveSnap(raftSnap); err != nil {
			return errors.Wrap(err, "Error saving snapshot to persistent storage")
		}
	}

	if rn.wal != nil {
		var walSnap walpb.Snapshot
		walSnap.Index, walSnap.Term = raftSnap.Metadata.Index, raftSnap.Metadata.Term

		if err := rn.wal.SaveSnapshot(walSnap); err != nil {
			return errors.Wrap(err, "Error updating WAL with snapshot")
		}
		if err := rn.wal.ReleaseLockTo(raftSnap.Metadata.Index); err != nil {
			return errors.Wrap(err, "Error releasing WAL locks")
		}
	}
	return nil
}

func (rn *Node) initWAL(walSnap walpb.Snapshot) error {
	if rn.walDir() == "" {
		return nil
	}

	if !wal.Exist(rn.walDir()) {

		if err := os.MkdirAll(rn.walDir(), 0750); err != nil && !os.IsExist(err) {
			return errors.Wrap(err, "Error creating directory for raft WAL")
		}

		metaStruct := &walMetadata{
			NodeID:    rn.id,
			ClusterID: rn.cid,
		}

		metaData, err := json.Marshal(metaStruct)
		if err != nil {
			return errors.Wrap(err, "Error marshaling WAL metadata")
		}

		w, err := wal.Create(rn.walDir(), metaData)
		if err != nil {
			return errors.Wrap(err, "Error creating new WAL")
		}
		rn.wal = w
	} else {
		// This assumes we WILL be reading this once elsewhere
		w, err := wal.Open(rn.walDir(), walSnap)
		if err != nil {
			return errors.Wrap(err, "Error opening existing WAL")
		}
		rn.wal = w
	}

	return nil
}

func (rn *Node) restoreMetadata(wMetadata []byte) error {
	var metaData walMetadata
	if err := json.Unmarshal(wMetadata, &metaData); err != nil {
		return errors.Wrap(err, "Error unmarshaling WAL metadata")
	}

	rn.id, rn.cid = metaData.NodeID, metaData.ClusterID
	rn.raftConfig.ID = metaData.NodeID
	return nil
}

func (rn *Node) restoreMemoryStorage(raftSnap raftpb.Snapshot, hState raftpb.HardState, ents []raftpb.Entry) error {
	if !raft.IsEmptySnap(raftSnap) {
		if err := rn.raftStorage.ApplySnapshot(raftSnap); err != nil {
			return errors.Wrap(err, "Error applying snapshot to raft memory storage")
		}
	}

	if rn.wal != nil {
		if err := rn.raftStorage.SetHardState(hState); err != nil {
			return errors.Wrap(err, "Error setting memory hardstate")
		}

		if err := rn.raftStorage.Append(ents); err != nil {
			return errors.Wrap(err, "Error appending entries to memory storage")
		}
	}

	return nil
}

func (rn *Node) deletePersistentData() error {
	if rn.snapDir() != "" {
		if err := os.RemoveAll(rn.snapDir()); err != nil {
			return errors.Wrap(err, "Error deleting snapshot directory")
		}
	}
	if rn.walDir() != "" {
		//TODO: Should be delete walDir or snapDir()?
		if err := os.RemoveAll(rn.walDir()); err != nil {
			return errors.Wrap(err, "Error deleting WAL directory")
		}
	}
	return nil
}

func (rn *Node) walDir() string {
	if rn.dataDir == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", rn.dataDir, walDirExtension)
}

func (rn *Node) snapDir() string {
	if rn.dataDir == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", rn.dataDir, snapDirExtension)
}
