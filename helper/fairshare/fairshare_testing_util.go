// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package fairshare

import (
	"fmt"
	"testing"

	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
)

type testJob struct {
	id     string
	ex     func(id string) error
	onFail func(error)
}

func (t *testJob) Execute() error {
	return t.ex(t.id)
}

func (t *testJob) OnFailure(err error) {
	t.onFail(err)
}

func newTestJob(t *testing.T, id string, ex func(string) error, onFail func(error)) testJob {
	t.Helper()
	if ex == nil {
		t.Errorf("ex cannot be nil")
	}
	if onFail == nil {
		t.Errorf("onFail cannot be nil")
	}

	return testJob{
		id:     id,
		ex:     ex,
		onFail: onFail,
	}
}

func newDefaultTestJob(t *testing.T, id string) testJob {
	ex := func(_ string) error { return nil }
	onFail := func(_ error) {}
	return newTestJob(t, id, ex, onFail)
}

func newTestLogger(name string) log.Logger {
	guid, err := uuid.GenerateUUID()
	if err != nil {
		guid = "no-guid"
	}
	return log.New(&log.LoggerOptions{
		Name:  fmt.Sprintf("%s-%s", name, guid),
		Level: log.LevelFromString("TRACE"),
	})
}

func GetNumWorkers(j *JobManager) int {
	return j.workerPool.numWorkers
}
