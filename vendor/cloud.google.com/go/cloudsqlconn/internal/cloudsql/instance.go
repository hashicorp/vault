// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudsql

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/cloudsqlconn/debug"
	"cloud.google.com/go/cloudsqlconn/errtype"
	"cloud.google.com/go/cloudsqlconn/instance"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

const (
	// the refresh buffer is the amount of time before a refresh operation's
	// certificate expires that a new refresh operation begins.
	refreshBuffer = 4 * time.Minute

	// refreshInterval is the amount of time between refresh attempts as
	// enforced by the rate limiter.
	refreshInterval = 30 * time.Second

	// RefreshTimeout is the maximum amount of time to wait for a refresh
	// cycle to complete. This value should be greater than the
	// refreshInterval.
	RefreshTimeout = 60 * time.Second

	// FailoverPeriod is the frequency with which the dialer will check
	// if the DNS record has changed for connections configured using
	// a DNS name.
	FailoverPeriod = 30 * time.Second

	// refreshBurst is the initial burst allowed by the rate limiter.
	refreshBurst = 2
)

// refreshOperation is a pending result of a refresh operation of data used to
// connect securely. It should only be initialized by the Instance struct as
// part of a refresh cycle.
type refreshOperation struct {
	// indicates the struct is ready to read from
	ready chan struct{}
	// timer that triggers refresh, can be used to cancel.
	timer  *time.Timer
	result ConnectionInfo
	err    error
}

// cancel prevents the the refresh operation from starting, if it hasn't
// already started. Returns true if timer was stopped successfully, or false if
// it has already started.
func (r *refreshOperation) cancel() bool {
	return r.timer.Stop()
}

// isValid returns true if this result is complete, successful, and is still
// valid.
func (r *refreshOperation) isValid() bool {
	// verify the refreshOperation has finished running
	select {
	default:
		return false
	case <-r.ready:
		if r.err != nil || time.Now().After(r.result.Expiration.Round(0)) {
			return false
		}
		return true
	}
}

// RefreshAheadCache manages the information used to connect to the Cloud SQL
// instance by periodically calling the Cloud SQL Admin API. It automatically
// refreshes the required information approximately 4 minutes before the
// previous certificate expires (every ~56 minutes).
type RefreshAheadCache struct {
	// openConns is the number of open connections to the instance.
	openConns uint64

	connName instance.ConnName
	logger   debug.ContextLogger

	// refreshTimeout sets the maximum duration a refresh cycle can run
	// for.
	refreshTimeout time.Duration
	// l controls the rate at which refresh cycles are run.
	l *rate.Limiter
	r adminAPIClient

	mu              sync.RWMutex
	useIAMAuthNDial bool
	// cur represents the current refreshOperation that will be used to
	// create connections. If a valid complete refreshOperation isn't
	// available it's possible for cur to be equal to next.
	cur *refreshOperation
	// next represents a future or ongoing refreshOperation. Once complete,
	// it will replace cur and schedule a replacement to occur.
	next *refreshOperation

	// ctx is the default ctx for refresh operations. Canceling it prevents
	// new refresh operations from being triggered.
	ctx    context.Context
	cancel context.CancelFunc
}

// NewRefreshAheadCache initializes a new Instance given an instance connection name
func NewRefreshAheadCache(
	cn instance.ConnName,
	l debug.ContextLogger,
	client *sqladmin.Service,
	key *rsa.PrivateKey,
	refreshTimeout time.Duration,
	ts oauth2.TokenSource,
	dialerID string,
	useIAMAuthNDial bool,
) *RefreshAheadCache {
	ctx, cancel := context.WithCancel(context.Background())
	i := &RefreshAheadCache{
		connName: cn,
		logger:   l,
		l:        rate.NewLimiter(rate.Every(refreshInterval), refreshBurst),
		r: newAdminAPIClient(
			l,
			client,
			key,
			ts,
			dialerID,
		),
		refreshTimeout:  refreshTimeout,
		useIAMAuthNDial: useIAMAuthNDial,
		ctx:             ctx,
		cancel:          cancel,
	}
	// For the initial refresh operation, set cur = next so that connection
	// requests block until the first refresh is complete.
	i.mu.Lock()
	i.cur = i.scheduleRefresh(0)
	i.next = i.cur
	i.mu.Unlock()
	return i
}

// Close closes the instance; it stops the refresh cycle and prevents it from
// making additional calls to the Cloud SQL Admin API.
func (i *RefreshAheadCache) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.cancel()
	i.cur.cancel()
	i.next.cancel()
	return nil
}

// ConnectionInfo contains all necessary information to connect securely to the
// server-side Proxy running on a Cloud SQL instance.
type ConnectionInfo struct {
	ConnectionName    instance.ConnName
	ClientCertificate tls.Certificate
	ServerCACert      []*x509.Certificate
	ServerCAMode      string
	DBVersion         string
	// The DNSName is from the ConnectSettings API.
	// It is used to validate the server identity of the CAS instances.
	DNSName    string
	Expiration time.Time

	addrs map[string]string
}

// NewConnectionInfo initializes a ConnectionInfo struct.
func NewConnectionInfo(
	cn instance.ConnName,
	dnsName string,
	serverCAMode string,
	version string,
	ipAddrs map[string]string,
	serverCACert []*x509.Certificate,
	clientCert tls.Certificate,
) ConnectionInfo {
	return ConnectionInfo{
		addrs:             ipAddrs,
		ClientCertificate: clientCert,
		ServerCACert:      serverCACert,
		ServerCAMode:      serverCAMode,
		Expiration:        clientCert.Leaf.NotAfter,
		DBVersion:         version,
		ConnectionName:    cn,
		DNSName:           dnsName,
	}
}

// Addr returns the IP address or DNS name for the given IP type.
func (c ConnectionInfo) Addr(ipType string) (string, error) {
	var (
		addr string
		ok   bool
	)
	switch ipType {
	case AutoIP:
		// Try Public first
		addr, ok = c.addrs[PublicIP]
		if !ok {
			// Try Private second
			addr, ok = c.addrs[PrivateIP]
		}
	default:
		addr, ok = c.addrs[ipType]
	}
	if !ok {
		err := errtype.NewConfigError(
			fmt.Sprintf("instance does not have IP of type %q", ipType),
			c.ConnectionName.String(),
		)
		return "", err
	}
	return addr, nil
}

// TLSConfig constructs a TLS configuration for the given connection info.
func (c ConnectionInfo) TLSConfig() *tls.Config {
	pool := x509.NewCertPool()
	for _, caCert := range c.ServerCACert {
		pool.AddCert(caCert)
	}
	if c.ServerCAMode == "GOOGLE_MANAGED_CAS_CA" {
		// For CAS instances, we can rely on the DNS name to verify the server identity.
		return &tls.Config{
			ServerName:   c.DNSName,
			Certificates: []tls.Certificate{c.ClientCertificate},
			RootCAs:      pool,
			MinVersion:   tls.VersionTLS13,
		}
	}
	return &tls.Config{
		ServerName:   c.ConnectionName.String(),
		Certificates: []tls.Certificate{c.ClientCertificate},
		RootCAs:      pool,
		// We need to set InsecureSkipVerify to true due to
		// https://github.com/GoogleCloudPlatform/cloudsql-proxy/issues/194
		// https://tip.golang.org/doc/go1.11#crypto/x509
		//
		// Since we have a secure channel to the Cloud SQL API which we use to
		// retrieve the certificates, we instead need to implement our own
		// VerifyPeerCertificate function that will verify that the certificate
		// is OK.
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyPeerCertificateFunc(c.ConnectionName, pool),
		MinVersion:            tls.VersionTLS13,
	}
}

// verifyPeerCertificateFunc creates a VerifyPeerCertificate func that
// verifies that the peer certificate is in the cert pool. We need to define
// our own because CloudSQL instances use the instance name (e.g.,
// my-project:my-instance) instead of a valid domain name for the certificate's
// Common Name.
func verifyPeerCertificateFunc(
	cn instance.ConnName, pool *x509.CertPool,
) func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
		if len(rawCerts) == 0 {
			return errtype.NewDialError(
				"no certificate to verify", cn.String(), nil,
			)
		}

		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return errtype.NewDialError(
				"failed to parse X.509 certificate", cn.String(), err,
			)
		}

		opts := x509.VerifyOptions{Roots: pool}
		if _, err = cert.Verify(opts); err != nil {
			return errtype.NewDialError(
				"failed to verify certificate", cn.String(), err,
			)
		}

		certInstanceName := fmt.Sprintf("%s:%s", cn.Project(), cn.Name())
		if cert.Subject.CommonName != certInstanceName {
			return errtype.NewDialError(
				fmt.Sprintf(
					"certificate had CN %q, expected %q",
					cert.Subject.CommonName, certInstanceName,
				),
				cn.String(),
				nil,
			)
		}
		return nil
	}
}

// ConnectionInfo returns an IP address specified by ipType (i.e., public or
// private) and a TLS config that can be used to connect to a Cloud SQL
// instance.
func (i *RefreshAheadCache) ConnectionInfo(ctx context.Context) (ConnectionInfo, error) {
	op, err := i.refreshOperation(ctx)
	if err != nil {
		return ConnectionInfo{}, err
	}
	return op.result, nil
}

// UpdateRefresh cancels all existing refresh attempts and schedules new
// attempts with the provided config only if it differs from the current
// configuration.
func (i *RefreshAheadCache) UpdateRefresh(useIAMAuthNDial *bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if useIAMAuthNDial != nil && *useIAMAuthNDial != i.useIAMAuthNDial {
		// Cancel any pending refreshes
		i.cur.cancel()
		i.next.cancel()

		i.useIAMAuthNDial = *useIAMAuthNDial
		// reschedule a new refresh immediately
		i.cur = i.scheduleRefresh(0)
		i.next = i.cur
	}
}

// ForceRefresh triggers an immediate refresh operation to be scheduled and
// used for future connection attempts. Until the refresh completes, the
// existing connection info will be available for use if valid.
func (i *RefreshAheadCache) ForceRefresh() {
	i.mu.Lock()
	defer i.mu.Unlock()
	// If the next refresh hasn't started yet, we can cancel it and start an
	// immediate one
	if i.next.cancel() {
		i.next = i.scheduleRefresh(0)
	}
	// block all sequential connection attempts on the next refresh operation
	// if current is invalid
	if !i.cur.isValid() {
		i.cur = i.next
	}
}

// refreshOperation returns the most recent refresh operation
// waiting for it to complete if necessary
func (i *RefreshAheadCache) refreshOperation(ctx context.Context) (*refreshOperation, error) {
	i.mu.RLock()
	cur := i.cur
	i.mu.RUnlock()
	var err error
	select {
	case <-cur.ready:
		err = cur.err
	case <-ctx.Done():
		err = ctx.Err()
	case <-i.ctx.Done():
		err = i.ctx.Err()
	}
	if err != nil {
		return nil, err
	}
	return cur, nil
}

// refreshDuration returns the duration to wait before starting the next
// refresh. Usually that duration will be half of the time until certificate
// expiration.
func refreshDuration(now, certExpiry time.Time) time.Duration {
	d := certExpiry.Sub(now.Round(0))
	if d < time.Hour {
		// Something is wrong with the certificate, refresh now.
		if d < refreshBuffer {
			return 0
		}
		// Otherwise wait until 4 minutes before expiration for next
		// refresh cycle.
		return d - refreshBuffer
	}
	return d / 2
}

// scheduleRefresh schedules a refresh operation to be triggered after a given
// duration. The returned refreshOperation can be used to either Cancel or Wait
// for the operation's completion.
func (i *RefreshAheadCache) scheduleRefresh(d time.Duration) *refreshOperation {
	r := &refreshOperation{}
	r.ready = make(chan struct{})
	r.timer = time.AfterFunc(d, func() {
		// instance has been closed, don't schedule anything
		if err := i.ctx.Err(); err != nil {
			i.logger.Debugf(
				context.Background(),
				"[%v] Instance is closed, stopping refresh operations",
				i.connName.String(),
			)
			r.err = err
			close(r.ready)
			return
		}
		i.logger.Debugf(
			context.Background(),
			"[%v] Connection info refresh operation started",
			i.connName.String(),
		)

		ctx, cancel := context.WithTimeout(i.ctx, i.refreshTimeout)
		defer cancel()

		// avoid refreshing too often to try not to tax the SQL Admin
		// API quotas
		err := i.l.Wait(ctx)
		if err != nil {
			r.err = errtype.NewDialError(
				"context was canceled or expired before refresh completed",
				i.connName.String(),
				nil,
			)
		} else {
			var useIAMAuthN bool
			i.mu.Lock()
			useIAMAuthN = i.useIAMAuthNDial
			i.mu.Unlock()
			r.result, r.err = i.r.ConnectionInfo(
				ctx, i.connName, useIAMAuthN,
			)
		}
		switch r.err {
		case nil:
			i.logger.Debugf(
				ctx,
				"[%v] Connection info refresh operation complete",
				i.connName.String(),
			)
			i.logger.Debugf(
				ctx,
				"[%v] Current certificate expiration = %v",
				i.connName.String(),
				r.result.Expiration.UTC().Format(time.RFC3339),
			)
		default:
			i.logger.Debugf(
				ctx,
				"[%v] Connection info refresh operation failed, err = %v",
				i.connName.String(),
				r.err,
			)
		}

		close(r.ready)

		// Once the refresh is complete, update "current" with working
		// refreshOperation and schedule a new refresh
		i.mu.Lock()
		defer i.mu.Unlock()

		// if failed, scheduled the next refresh immediately
		if r.err != nil {
			i.logger.Debugf(
				ctx,
				"[%v] Connection info refresh operation scheduled immediately",
				i.connName.String(),
			)
			i.next = i.scheduleRefresh(0)
			// If the latest refreshOperation is bad, avoid replacing the
			// used refreshOperation while it's still valid and potentially
			// able to provide successful connections. TODO: This
			// means that errors while the current refreshOperation is still
			// valid are suppressed. We should try to surface
			// errors in a more meaningful way.
			if !i.cur.isValid() {
				i.cur = r
			}
			return
		}

		// Update the current results, and schedule the next refresh in
		// the future
		i.cur = r
		t := refreshDuration(time.Now(), i.cur.result.Expiration)
		i.logger.Debugf(
			ctx,
			"[%v] Connection info refresh operation scheduled at %v (now + %v)",
			i.connName.String(),
			time.Now().Add(t).UTC().Format(time.RFC3339),
			t.Round(time.Minute),
		)
		i.next = i.scheduleRefresh(t)
	})
	return r
}
