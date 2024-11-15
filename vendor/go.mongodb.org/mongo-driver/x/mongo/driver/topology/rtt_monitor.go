// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

const (
	rttAlphaValue = 0.2
	minSamples    = 10
	maxSamples    = 500
)

type rttConfig struct {
	// The minimum interval between RTT measurements. The actual interval may be greater if running
	// the operation takes longer than the interval.
	interval time.Duration

	// The timeout applied to running the "hello" operation. If the timeout is reached while running
	// the operation, the RTT sample is discarded. The default is 1 minute.
	timeout time.Duration

	minRTTWindow       time.Duration
	createConnectionFn func() *connection
	createOperationFn  func(driver.Connection) *operation.Hello
}

type rttMonitor struct {
	mu sync.RWMutex // mu guards samples, offset, minRTT, averageRTT, and averageRTTSet

	// connMu guards connecting and disconnecting. This is necessary since
	// disconnecting will await the cancellation of a started connection. The
	// use case for rttMonitor.connect needs to be goroutine safe.
	connMu        sync.Mutex
	samples       []time.Duration
	offset        int
	minRTT        time.Duration
	rtt90         time.Duration
	averageRTT    time.Duration
	averageRTTSet bool

	closeWg  sync.WaitGroup
	cfg      *rttConfig
	ctx      context.Context
	cancelFn context.CancelFunc
}

var _ driver.RTTMonitor = &rttMonitor{}

func newRTTMonitor(cfg *rttConfig) *rttMonitor {
	if cfg.interval <= 0 {
		panic("RTT monitor interval must be greater than 0")
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Determine the number of samples we need to keep to store the minWindow of RTT durations. The
	// number of samples must be between [10, 500].
	numSamples := int(math.Max(minSamples, math.Min(maxSamples, float64((cfg.minRTTWindow)/cfg.interval))))

	return &rttMonitor{
		samples:  make([]time.Duration, numSamples),
		cfg:      cfg,
		ctx:      ctx,
		cancelFn: cancel,
	}
}

func (r *rttMonitor) connect() {
	r.connMu.Lock()
	defer r.connMu.Unlock()

	r.closeWg.Add(1)

	go func() {
		defer r.closeWg.Done()

		r.start()
	}()
}

func (r *rttMonitor) disconnect() {
	r.connMu.Lock()
	defer r.connMu.Unlock()

	r.cancelFn()

	// Wait for the existing connection to complete.
	r.closeWg.Wait()
}

func (r *rttMonitor) start() {
	var conn *connection
	defer func() {
		if conn != nil {
			// If the connection exists, we need to wait for it to be connected because
			// conn.connect() and conn.close() cannot be called concurrently. If the connection
			// wasn't successfully opened, its state was set back to disconnected, so calling
			// conn.close() will be a no-op.
			conn.closeConnectContext()
			conn.wait()
			_ = conn.close()
		}
	}()

	ticker := time.NewTicker(r.cfg.interval)
	defer ticker.Stop()

	for {
		conn := r.cfg.createConnectionFn()
		err := conn.connect(r.ctx)

		// Add an RTT sample from the new connection handshake and start a runHellos() loop if we
		// successfully established the new connection. Otherwise, close the connection and try to
		// create another new connection.
		if err == nil {
			r.addSample(conn.helloRTT)
			r.runHellos(conn)
		}

		// Close any connection here because we're either about to try to create another new
		// connection or we're about to exit the loop.
		_ = conn.close()

		// If a connection error happens quickly, always wait for the monitoring interval to try
		// to create a new connection to prevent creating connections too quickly.
		select {
		case <-ticker.C:
		case <-r.ctx.Done():
			return
		}
	}
}

// runHellos runs "hello" operations in a loop using the provided connection, measuring and
// recording the operation durations as RTT samples. If it encounters any errors, it returns.
func (r *rttMonitor) runHellos(conn *connection) {
	ticker := time.NewTicker(r.cfg.interval)
	defer ticker.Stop()

	for {
		// Assume that the connection establishment recorded the first RTT sample, so wait for the
		// first tick before trying to record another RTT sample.
		select {
		case <-ticker.C:
		case <-r.ctx.Done():
			return
		}

		// Create a Context with the operation timeout specified in the RTT monitor config. If a
		// timeout is not set in the RTT monitor config, default to the connection's
		// "connectTimeoutMS". The purpose of the timeout is to allow the RTT monitor to continue
		// monitoring server RTTs after an operation gets stuck. An operation can get stuck if the
		// server or a proxy stops responding to requests on the RTT connection but does not close
		// the TCP socket, effectively creating an operation that will never complete. We expect
		// that "connectTimeoutMS" provides at least enough time for a single round trip.
		timeout := r.cfg.timeout
		if timeout <= 0 {
			timeout = conn.config.connectTimeout
		}
		ctx, cancel := context.WithTimeout(r.ctx, timeout)

		start := time.Now()
		err := r.cfg.createOperationFn(initConnection{conn}).Execute(ctx)
		cancel()
		if err != nil {
			return
		}
		// Only record a sample if the "hello" operation was successful. If it was not successful,
		// the operation may not have actually performed a complete round trip, so the duration may
		// be artificially short.
		r.addSample(time.Since(start))
	}
}

// reset sets the average and min RTT to 0. This should only be called from the server monitor when an error
// occurs during a server check. Errors in the RTT monitor should not reset the RTTs.
func (r *rttMonitor) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.samples {
		r.samples[i] = 0
	}
	r.offset = 0
	r.minRTT = 0
	r.rtt90 = 0
	r.averageRTT = 0
	r.averageRTTSet = false
}

func (r *rttMonitor) addSample(rtt time.Duration) {
	// Lock for the duration of this method. We're doing compuationally inexpensive work very infrequently, so lock
	// contention isn't expected.
	r.mu.Lock()
	defer r.mu.Unlock()

	r.samples[r.offset] = rtt
	r.offset = (r.offset + 1) % len(r.samples)
	// Set the minRTT and 90th percentile RTT of all collected samples. Require at least 10 samples before
	// setting these to prevent noisy samples on startup from artificially increasing RTT and to allow the
	// calculation of a 90th percentile.
	r.minRTT = min(r.samples, minSamples)
	r.rtt90 = percentile(90.0, r.samples, minSamples)

	if !r.averageRTTSet {
		r.averageRTT = rtt
		r.averageRTTSet = true
		return
	}

	r.averageRTT = time.Duration(rttAlphaValue*float64(rtt) + (1-rttAlphaValue)*float64(r.averageRTT))
}

// min returns the minimum value of the slice of duration samples. Zero values are not considered
// samples and are ignored. If no samples or fewer than minSamples are found in the slice, min
// returns 0.
func min(samples []time.Duration, minSamples int) time.Duration {
	count := 0
	min := time.Duration(math.MaxInt64)
	for _, d := range samples {
		if d > 0 {
			count++
		}
		if d > 0 && d < min {
			min = d
		}
	}
	if count == 0 || count < minSamples {
		return 0
	}
	return min
}

// percentile returns the specified percentile value of the slice of duration samples. Zero values
// are not considered samples and are ignored. If no samples or fewer than minSamples are found
// in the slice, percentile returns 0.
func percentile(perc float64, samples []time.Duration, minSamples int) time.Duration {
	// Convert Durations to float64s.
	floatSamples := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if sample > 0 {
			floatSamples = append(floatSamples, float64(sample))
		}
	}
	if len(floatSamples) == 0 || len(floatSamples) < minSamples {
		return 0
	}

	p, err := stats.Percentile(floatSamples, perc)
	if err != nil {
		panic(fmt.Errorf("x/mongo/driver/topology: error calculating %f percentile RTT: %w for samples:\n%v", perc, err, floatSamples))
	}
	return time.Duration(p)
}

// EWMA returns the exponentially weighted moving average observed round-trip time.
func (r *rttMonitor) EWMA() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.averageRTT
}

// Min returns the minimum observed round-trip time over the window period.
func (r *rttMonitor) Min() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.minRTT
}

// P90 returns the 90th percentile observed round-trip time over the window period.
func (r *rttMonitor) P90() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.rtt90
}

// Stats returns stringified stats of the current state of the monitor.
func (r *rttMonitor) Stats() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Calculate standard deviation and average (non-EWMA) of samples.
	var sum float64
	floatSamples := make([]float64, 0, len(r.samples))
	for _, sample := range r.samples {
		if sample > 0 {
			floatSamples = append(floatSamples, float64(sample))
			sum += float64(sample)
		}
	}

	var avg, stdDev float64
	if len(floatSamples) > 0 {
		avg = sum / float64(len(floatSamples))

		var err error
		stdDev, err = stats.StandardDeviation(floatSamples)
		if err != nil {
			panic(fmt.Errorf("x/mongo/driver/topology: error calculating standard deviation RTT: %w for samples:\n%v", err, floatSamples))
		}
	}

	return fmt.Sprintf(
		"network round-trip time stats: avg: %v, min: %v, 90th pct: %v, stddev: %v",
		time.Duration(avg),
		r.minRTT,
		r.rtt90,
		time.Duration(stdDev))
}
