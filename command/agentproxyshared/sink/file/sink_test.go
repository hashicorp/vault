// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package file

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestSinkServer(t *testing.T) {
	log := logging.NewVaultLogger(hclog.Trace)

	fs1, path1 := testFileSink(t, log)
	defer os.RemoveAll(path1)
	fs2, path2 := testFileSink(t, log)
	defer os.RemoveAll(path2)

	ctx, cancelFunc := context.WithCancel(context.Background())

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: log.Named("sink.server"),
	})

	uuidStr, _ := uuid.GenerateUUID()
	in := make(chan string)
	sinks := []*sink.SinkConfig{fs1, fs2}
	errCh := make(chan error)
	tokenRenewalInProgress := &atomic.Bool{}
	tokenRenewalInProgress.Store(true)
	go func() {
		errCh <- ss.Run(ctx, in, sinks, tokenRenewalInProgress)
	}()

	// Seed a token
	in <- uuidStr

	// Tell it to shut down and give it time to do so
	timer := time.AfterFunc(3*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, path := range []string{path1, path2} {
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/token", path))
		if err != nil {
			t.Fatal(err)
		}

		if string(fileBytes) != uuidStr {
			t.Fatalf("expected %s, got %s", uuidStr, string(fileBytes))
		}
	}

	if tokenRenewalInProgress.Load() {
		t.Fatal("should have reset tokenRenewalInProgress to false")
	}
}

type badSink struct {
	tryCount uint32
	logger   hclog.Logger
}

func (b *badSink) WriteToken(token string) error {
	switch token {
	case "bad":
		atomic.AddUint32(&b.tryCount, 1)
		b.logger.Info("got bad")
		return errors.New("bad")
	case "good":
		atomic.StoreUint32(&b.tryCount, 0)
		b.logger.Info("got good")
		return nil
	default:
		return errors.New("unknown case")
	}
}

func TestSinkServerRetry(t *testing.T) {
	log := logging.NewVaultLogger(hclog.Trace)

	b1 := &badSink{logger: log.Named("b1")}
	b2 := &badSink{logger: log.Named("b2")}

	ctx, cancelFunc := context.WithCancel(context.Background())

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: log.Named("sink.server"),
	})

	in := make(chan string)
	sinks := []*sink.SinkConfig{{Sink: b1}, {Sink: b2}}
	errCh := make(chan error)
	tokenRenewalInProgress := &atomic.Bool{}
	tokenRenewalInProgress.Store(true)

	go func() {
		errCh <- ss.Run(ctx, in, sinks, tokenRenewalInProgress)
	}()

	// Seed a token
	in <- "bad"

	// During this time we should see it retry multiple times
	time.Sleep(10 * time.Second)
	if atomic.LoadUint32(&b1.tryCount) < 2 {
		t.Fatal("bad try count")
	}
	if atomic.LoadUint32(&b2.tryCount) < 2 {
		t.Fatal("bad try count")
	}

	if !tokenRenewalInProgress.Load() {
		t.Fatal("token renewal should still be in progress, sink server has not exited")
	}

	in <- "good"

	time.Sleep(2 * time.Second)
	if atomic.LoadUint32(&b1.tryCount) != 0 {
		t.Fatal("bad try count")
	}
	if atomic.LoadUint32(&b2.tryCount) != 0 {
		t.Fatal("bad try count")
	}

	// Tell it to shut down and give it time to do so
	cancelFunc()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	}

	if tokenRenewalInProgress.Load() {
		t.Fatal("should have reset tokenRenewalInProgress to false")
	}
}
