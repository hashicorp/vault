// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
)

func GetRaft(t testing.TB, bootstrap bool, noStoreState bool) (*RaftBackend, string) {
	raftDir := t.TempDir()
	t.Logf("raft dir: %s", raftDir)

	conf := map[string]string{
		"path":          raftDir,
		"trailing_logs": "100",
	}

	if noStoreState {
		conf["doNotStoreLatestState"] = ""
	}

	return getRaftWithDirAndConfig(t, bootstrap, conf)
}

func GetRaftWithConfig(t testing.TB, bootstrap bool, noStoreState bool, conf map[string]string) (*RaftBackend, string) {
	raftDir := t.TempDir()
	t.Logf("raft dir: %s", raftDir)

	conf["path"] = raftDir
	if noStoreState {
		conf["doNotStoreLatestState"] = ""
	}

	return getRaftWithDirAndConfig(t, bootstrap, conf)
}

func getRaftWithDirAndConfig(t testing.TB, bootstrap bool, conf map[string]string) (*RaftBackend, string) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  fmt.Sprintf("raft-%s", id),
		Level: hclog.Trace,
	})

	conf["node_id"] = id

	backendRaw, err := NewRaftBackend(conf, logger)
	if err != nil {
		t.Fatal(err)
	}
	backend := backendRaw.(*RaftBackend)

	if bootstrap {
		err = backend.Bootstrap([]Peer{
			{
				ID:      backend.NodeID(),
				Address: backend.NodeID(),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		err = backend.SetupCluster(context.Background(), SetupOpts{})
		if err != nil {
			t.Fatal(err)
		}

		for {
			if backend.raft.AppliedIndex() >= 2 {
				break
			}
		}

	}

	backend.DisableAutopilot()

	return backend, conf["path"]
}
