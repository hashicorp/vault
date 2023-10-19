// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
)

func GetRaft(t testing.TB, bootstrap bool, noStoreState bool) (*RaftBackend, string) {
	return GetRaftWithOpts(t, bootstrap, noStoreState, SetupOpts{})
}

func GetRaftWithOpts(t testing.TB, bootstrap bool, noStoreState bool, setupOpts SetupOpts) (*RaftBackend, string) {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)

	return getRaftWithDir(t, bootstrap, noStoreState, raftDir, setupOpts)
}

func getRaftWithDir(t testing.TB, bootstrap bool, noStoreState bool, raftDir string, setupOpts SetupOpts) (*RaftBackend, string) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  fmt.Sprintf("raft-%s", id),
		Level: hclog.Trace,
	})
	logger.Info("raft dir", "dir", raftDir)

	conf := map[string]string{
		"path":          raftDir,
		"trailing_logs": "100",
		"node_id":       id,
	}

	if noStoreState {
		conf["doNotStoreLatestState"] = ""
	}

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

		err = backend.SetupCluster(context.Background(), setupOpts)
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

	return backend, raftDir
}
