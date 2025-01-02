// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
)

func (b *RaftBackend) DataDir(t testing.TB) string {
	t.Helper()
	return b.dataDir
}

func GetRaft(t testing.TB, bootstrap bool, noStoreState bool) (*RaftBackend, string) {
	return getRaftInternal(t, bootstrap, defaultRaftConfig(t, bootstrap, noStoreState), nil, nil, nil)
}

func GetRaftWithConfig(t testing.TB, bootstrap bool, noStoreState bool, conf map[string]string) (*RaftBackend, string) {
	defaultConf := defaultRaftConfig(t, bootstrap, noStoreState)
	conf["path"] = defaultConf["path"]
	conf["doNotStoreLatestState"] = defaultConf["doNotStoreLatestState"]
	return getRaftInternal(t, bootstrap, conf, nil, nil, nil)
}

func GetRaftWithConfigAndSetupOpts(t testing.TB, bootstrap bool, noStoreState bool, conf map[string]string, setupOpts *SetupOpts) (*RaftBackend, string) {
	defaultConf := defaultRaftConfig(t, bootstrap, noStoreState)
	conf["path"] = defaultConf["path"]
	conf["doNotStoreLatestState"] = defaultConf["doNotStoreLatestState"]
	return getRaftInternal(t, bootstrap, conf, setupOpts, nil, nil)
}

func GetRaftWithConfigAndInitFn(t testing.TB, bootstrap bool, noStoreState bool, conf map[string]string, initFn func(b *RaftBackend)) (*RaftBackend, string) {
	defaultConf := defaultRaftConfig(t, bootstrap, noStoreState)
	conf["path"] = defaultConf["path"]
	conf["doNotStoreLatestState"] = defaultConf["doNotStoreLatestState"]
	return getRaftInternal(t, bootstrap, conf, nil, nil, initFn)
}

func GetRaftWithLogOutput(t testing.TB, bootstrap bool, noStoreState bool, logOutput io.Writer) (*RaftBackend, string) {
	return getRaftInternal(t, bootstrap, defaultRaftConfig(t, bootstrap, noStoreState), nil, logOutput, nil)
}

func defaultRaftConfig(t testing.TB, bootstrap bool, noStoreState bool) map[string]string {
	raftDir := t.TempDir()
	t.Logf("raft dir: %s", raftDir)

	conf := map[string]string{
		"path":          raftDir,
		"trailing_logs": "100",
	}

	if noStoreState {
		conf["doNotStoreLatestState"] = ""
	}

	return conf
}

func getRaftInternal(t testing.TB, bootstrap bool, conf map[string]string, setupOpts *SetupOpts, logOutput io.Writer, initFn func(b *RaftBackend)) (*RaftBackend, string) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("raft-%s", id),
		Level:  hclog.Trace,
		Output: logOutput,
	})

	conf["node_id"] = id

	backendRaw, err := NewRaftBackend(conf, logger)
	if err != nil {
		t.Fatal(err)
	}
	backend := backendRaw.(*RaftBackend)
	if initFn != nil {
		initFn(backend)
	}

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

		so := SetupOpts{}
		if setupOpts != nil {
			so = *setupOpts
		}

		err = backend.SetupCluster(context.Background(), so)
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
