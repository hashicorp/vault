// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
)

func GetRaft(t testing.TB, bootstrap bool, noStoreState bool) (*RaftBackend, string) {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)

	return getRaftWithDirAndLogOutput(t, bootstrap, noStoreState, raftDir, nil)
}

func GetRaftWithLogOutput(t testing.TB, bootstrap bool, noStoreState bool, logOutput io.Writer) (*RaftBackend, string) {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)

	return getRaftWithDirAndLogOutput(t, bootstrap, noStoreState, raftDir, logOutput)
}

func getRaftWithDirAndLogOutput(t testing.TB, bootstrap bool, noStoreState bool, raftDir string, logOutput io.Writer) (*RaftBackend, string) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("raft-%s", id),
		Level:  hclog.Trace,
		Output: logOutput,
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

	return backend, raftDir
}
