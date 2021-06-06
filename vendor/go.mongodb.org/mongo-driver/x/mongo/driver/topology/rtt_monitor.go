// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

const (
	rttAlphaValue = 0.2
)

type rttConfig struct {
	interval           time.Duration
	createConnectionFn func() (*connection, error)
	createOperationFn  func(driver.Connection) *operation.IsMaster
}

type rttMonitor struct {
	sync.Mutex
	conn          *connection
	averageRTT    time.Duration
	averageRTTSet bool
	closeWg       sync.WaitGroup
	cfg           *rttConfig
	ctx           context.Context
	cancelFn      context.CancelFunc
}

func newRttMonitor(cfg *rttConfig) *rttMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &rttMonitor{
		cfg:      cfg,
		ctx:      ctx,
		cancelFn: cancel,
	}
}

func (r *rttMonitor) connect() {
	r.closeWg.Add(1)
	go r.start()
}

func (r *rttMonitor) disconnect() {
	// Signal for the routine to stop.
	r.cancelFn()

	var conn *connection
	r.Lock()
	conn = r.conn
	r.Unlock()

	if conn != nil {
		// If the connection exists, we need to wait for it to be connected. We can ignore the error from conn.wait().
		// If the connection wasn't successfully opened, its state was set back to disconnected, so calling conn.close()
		// will be a noop.
		conn.closeConnectContext()
		_ = conn.wait()
		_ = conn.close()
	}

	r.closeWg.Wait()
}

func (r *rttMonitor) start() {
	defer r.closeWg.Done()
	ticker := time.NewTicker(r.cfg.interval)
	defer ticker.Stop()

	for {
		// The context is only cancelled in disconnect() so if there's an error on it, the monitor is shutting down.
		if r.ctx.Err() != nil {
			return
		}

		r.pingServer()

		select {
		case <-ticker.C:
		case <-r.ctx.Done():
			// Shutting down
			return
		}
	}
}

// reset sets the average RTT to 0. This should only be called from the server monitor when an error occurs during a
// server check. Errors in the RTT monitor should not reset the average RTT.
func (r *rttMonitor) reset() {
	r.Lock()
	defer r.Unlock()

	r.averageRTT = 0
	r.averageRTTSet = false
}

func (r *rttMonitor) setupRttConnection() error {
	conn, err := r.cfg.createConnectionFn()
	if err != nil {
		return err
	}

	r.Lock()
	r.conn = conn
	r.Unlock()

	r.conn.connect(r.ctx)
	return r.conn.wait()
}

func (r *rttMonitor) pingServer() {
	if r.conn == nil || r.conn.closed() {
		if err := r.setupRttConnection(); err != nil {
			return
		}

		// Add the initial connection handshake time as an RTT sample.
		r.addSample(r.conn.isMasterRTT)
		return
	}

	// We're using an already established connection. Issue an isMaster command to get a new RTT sample.
	rttConn := initConnection{r.conn}
	start := time.Now()
	err := r.cfg.createOperationFn(rttConn).Execute(r.ctx)
	if err != nil {
		// Errors from the RTT monitor do not reset the average RTT or update the topology, so we close the existing
		// connection and recreate it on the next check.
		_ = r.conn.close()
		return
	}

	r.addSample(time.Since(start))
}

func (r *rttMonitor) addSample(rtt time.Duration) {
	// Lock for the duration of this method. We're doing compuationally inexpensive work very infrequently, so lock
	// contention isn't expected.
	r.Lock()
	defer r.Unlock()

	if !r.averageRTTSet {
		r.averageRTT = rtt
		r.averageRTTSet = true
		return
	}

	r.averageRTT = time.Duration(rttAlphaValue*float64(rtt) + (1-rttAlphaValue)*float64(r.averageRTT))
}

func (r *rttMonitor) getRTT() time.Duration {
	r.Lock()
	defer r.Unlock()

	return r.averageRTT
}
