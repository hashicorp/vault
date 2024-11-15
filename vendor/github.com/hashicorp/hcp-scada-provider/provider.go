// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	msgpackrpc "github.com/hashicorp/net-rpc-msgpackrpc/v2"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/hashicorp/hcp-scada-provider/internal/client"
	"github.com/hashicorp/hcp-scada-provider/internal/client/dialer/tcp"
	"github.com/hashicorp/hcp-scada-provider/internal/listener"
	"github.com/hashicorp/hcp-scada-provider/types"
)

const (
	// defaultBackoff is the amount of time we back off if we encounter an
	// error, and no specific backoff is available.
	defaultBackoff = 10 * time.Second

	// disconnectDelay is the amount of time to wait between the moment
	// the disconnect RPC call is received and actually disconnecting the provider.
	disconnectDelay = time.Second

	// expiryDefault sets up a default time for the session expiry ticker
	// in the run() loop.
	expiryDefault = 60 * time.Minute
	// expiryFactor is the value to multiply the
	// the Expiry duration with and reduce it's value to
	// rehanshake within a good time margin, before the broker
	// closes the session.
	expiryFactor = 0.9
)

var (
	errNoRetry    = errors.New("provider is configured to not retry a connection")
	errNotRunning = errors.New("provider is not running")
)

type handler struct {
	provider listener.Provider
	listener net.Listener
}

// Provider is a high-level interface to SCADA by which instances declare
// themselves as a Service providing capabilities. Provider manages the
// client/server interactions required, making it simpler to integrate.
type Provider struct {
	config     *Config
	configLock sync.RWMutex

	logger hclog.Logger

	handlers     map[string]handler
	handlersLock sync.RWMutex

	noRetry     bool          // set when the server instructs us to not retry
	backoff     time.Duration // set when the server provides a longer backoff
	backoffLock sync.Mutex

	meta     map[string]string
	metaLock sync.RWMutex

	running     bool
	runningLock sync.Mutex

	sessionStatuses chan SessionStatus

	actions chan action

	cancel context.CancelFunc

	lastErrors chan timeError
}

// New creates a new SCADA provider instance using the configuration in config.
func New(config *Config) (SCADAProvider, error) {
	return newProvider(config)
}

func newProvider(config *Config) (*Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	p := &Provider{
		config:          config,
		logger:          config.Logger.Named("scada-provider"),
		meta:            make(map[string]string),
		handlers:        make(map[string]handler),
		sessionStatuses: make(chan SessionStatus),
		actions:         make(chan action),
		lastErrors:      make(chan timeError),
	}

	return p, nil
}

// UpdateMeta overwrites the internal map of meta-data values
// and performs a re-handshake to update the remote broker.
// If the provider isn't running, updated meta will be applied
// after the provoder starts.
func (p *Provider) UpdateMeta(m map[string]string) {
	p.metaLock.Lock()
	defer p.metaLock.Unlock()

	// reset the current metadata
	p.meta = make(map[string]string, len(m))

	// Update the map
	for k, v := range m {
		p.meta[k] = v
	}

	// tell the run loop to re-handshake and update the broker
	p.action(actionRehandshake)
}

// AddMeta upserts keys and values in the internal map of meta-data
// and performs a re-handshake to update the remote broker.
func (p *Provider) AddMeta(metas ...Meta) {
	p.metaLock.Lock()
	defer p.metaLock.Unlock()

	// Update the map
	for _, v := range metas {
		p.meta[v.Key] = v.Value
	}

	// tell the run loop to re-handshake and update the broker
	p.action(actionRehandshake)
}

// DeleteMeta delete keys from the meta-date values and then perform a
// re-handshake to update the remote broker.
func (p *Provider) DeleteMeta(keys ...string) {
	p.metaLock.Lock()
	defer p.metaLock.Unlock()

	// Update the map
	for _, v := range keys {
		delete(p.meta, v)
	}

	// tell the run loop to re-handshake and update the broker
	p.action(actionRehandshake)
}

// GetMeta returns the provider's current meta-data.
// The returned map is a copy and can be updated or modified.
func (p *Provider) GetMeta() map[string]string {
	p.metaLock.RLock()
	defer p.metaLock.RUnlock()

	// copy the map
	var meta = make(map[string]string, len(p.meta))
	for k, v := range p.meta {
		meta[k] = v
	}

	return meta
}

// Listen will expose the provided capability and make new connections
// available through the returned listener. Closing the listener will stop
// exposing the provided capability.
//
// The method will return an existing listener if the capability already existed.
// Listeners can be retrieved even when the provider is stopped (e.g. before it is
// started). New capabilities and new meta data can be added at any time.
//
// The listener will only be closed, if it is closed explicitly by calling Close().
// The listener will not be closed due to errors or when the provider is stopped.
// The listener can hence be used after a restart of the provider.
func (p *Provider) Listen(capability string) (net.Listener, error) {
	// Check if the capability already exists
	p.handlersLock.RLock()
	capHandler, ok := p.handlers[capability]
	p.handlersLock.RUnlock()

	if ok {
		return capHandler.listener, nil
	}

	// Get write lock
	p.handlersLock.Lock()
	defer p.handlersLock.Unlock()

	// Ensure that no concurrent call has set the listener in the meantime
	if capHandler, ok = p.handlers[capability]; ok {
		return capHandler.listener, nil
	}

	// Generate a provider and listener for the new capability
	capProvider, capListener, err := listener.New(capability)
	if err != nil {
		return nil, err
	}

	// Assign an OnClose callback on a listener, to make sure the handler is removed for the capacity.
	capListenerProxy := listener.WithCloseCallback(capListener, func() {
		p.handlersLock.Lock()
		defer p.handlersLock.Unlock()

		delete(p.handlers, capability)
	})

	p.handlers[capability] = handler{
		provider: capProvider,
		listener: capListenerProxy,
	}

	// re-handshake to update the broker
	p.action(actionRehandshake)

	return capListenerProxy, nil
}

// Start will register the provider on the SCADA broker and expose the
// registered capabilities.
func (p *Provider) Start() error {
	p.runningLock.Lock()
	defer p.runningLock.Unlock()

	// Check if the provider is already running
	if p.running {
		return nil
	}

	p.logger.Info("starting")

	// Set the provider to its running state
	p.running = true
	// Run the provider
	p.cancel = p.run()

	return nil
}

// Stop will gracefully close the currently active SCADA session. This will
// not close the capability listeners.
func (p *Provider) Stop() error {
	p.runningLock.Lock()
	defer p.runningLock.Unlock()

	// Check if the provider is already stopped
	if !p.running {
		return nil
	}

	p.logger.Info("stopping")

	// Stop the provider
	p.cancel()
	// Set the provider to its non-running state
	p.running = false

	return nil
}

// SessionStatus returns the status of the SCADA connection.
//
// The possibles statuses are:
//   - SessionStatusDisconnected: the provider is stopped
//   - SessionStatusConnecting:   in the connect/handshake cycle
//   - SessionStatusConnected:    connected and serving scada consumers
//   - SessionStatusWaiting:      disconnected and waiting to retry a connection to the broker
//
// The full lifecycle is: connecting -> connected -> waiting -> connecting -> ... -> disconnected.
func (p *Provider) SessionStatus() SessionStatus {
	p.runningLock.Lock()
	defer p.runningLock.Unlock()

	// Check if the provider is running
	if !p.running {
		return SessionStatusDisconnected
	}

	// get the status from the run() loop
	return <-p.sessionStatuses
}

// LastError returns the last error recorded in the provider
// connection state engine as well as the time at which the error occured.
// That record is erased at each occasion when the provider achieves a new connection.
//
// A few common internal error will return a known type:
//   - ErrProviderNotStarted: the provider is not started
//   - ErrInvalidCredentials: could not obtain a token with the supplied credentials
//   - ErrPermissionDenied:   principal does not have the permision to register as a provider
//
// Any other internal error will be returned directly and unchanged.
func (p *Provider) LastError() (time.Time, error) {
	p.runningLock.Lock()
	defer p.runningLock.Unlock()

	// Check if the provider is running
	if !p.running {
		return time.Now(), ErrProviderNotStarted
	}

	lastError := <-p.lastErrors
	return lastError.Time, lastError.error
}

/////

// run is a long running routine to manage the provider.
func (p *Provider) run() context.CancelFunc {
	// setup a statuses and errors channel to communicate with ourselves
	var statuses = make(chan SessionStatus)
	var errors = make(chan error)

	// setup a ticker for session's expiry
	var ticker = time.NewTicker(expiryDefault)

	// setup a context that will
	// cancel on stop
	ctx, cancel := context.WithCancel(context.Background())

	// setup done and ret to sync ctx.Done() with SessionStatusDisconnected
	var done, ret = make(chan bool), make(chan bool)

	go func() {
		defer cancel()
		var cl *client.Client
		// locally hold the current session status so we can communicate it
		var sessionStatus SessionStatus = SessionStatusDisconnected
		// locally hold the last error so we can communicate it
		var lastError timeError = NewTimeError(nil)

		// engage in running the provider
		for {
			select {
			case p.sessionStatuses <- sessionStatus:
				// p.SessionStatus() was waiting to read
			case p.lastErrors <- lastError:
				// p.LastError() was waiting to read
			case err := <-errors:
				// async receive an error from one of the handlers
				lastError = NewTimeError(err)

			case status := <-statuses:
				switch status {
				case SessionStatusWaiting:
					sessionStatus = SessionStatusWaiting
					// backoff
					go func() {
						if err := p.wait(ctx); err != nil {
							// wait returns an error if we shouldn't retry
							// or if ctx is canceled()
							statuses <- SessionStatusDisconnected
						} else {
							statuses <- SessionStatusConnecting
						}
					}()

				case SessionStatusConnecting:
					sessionStatus = SessionStatusConnecting
					// Try to connect a session
					go func() {
						// if we get canceled() during this,
						// connect will error out and we go to SessionStatusWaiting
						if client, err := p.connect(ctx); err != nil {
							// make a note of the error
							errors <- err
							// not connected
							statuses <- SessionStatusWaiting
						} else if response, err := p.handshake(ctx, client); err != nil {
							// make a note of the error
							errors <- err
							// connect closes client if any error
							// occured at handshake() except for resp.Authenticated == false
							statuses <- SessionStatusWaiting
						} else {
							// reset the ticker
							tickerReset(time.Now(), response.Expiry, ticker)
							// assigned the newly created client to this routine's cl
							cl = client
							statuses <- SessionStatusConnected
						}
					}()

				case SessionStatusConnected:
					sessionStatus = SessionStatusConnected
					// reset the error
					lastError = NewTimeError(nil)
					// reset any longer backoff period set by the Disconnect RPC call
					p.backoffReset()
					go func(client *client.Client) {
						// Handle the session
						if err := p.handleSession(ctx, client); err != nil {
							// make a note of the error
							errors <- err
							// handleSession will always close client
							// on errors or if the ctx is canceled().
							// go to the waiting state
							statuses <- SessionStatusWaiting
						}
					}(cl)

				case SessionStatusDisconnected:
					sessionStatus = SessionStatusDisconnected
					// after officially disconnecting, reset the backoff period for this provider
					p.backoffReset()
					close(done)
				}

			case <-ticker.C:
				// it's time to refresh the session with the broker
				// by issuing a re-handshake
				go func() {
					p.actions <- actionRehandshake
				}()

			case action := <-p.actions:
				// if sessionStatus is not SessionStatusConnected,
				// none of these actions can proceed
				if sessionStatus != SessionStatusConnected {
					continue
				}

				// these actions always close `cl` directly, or when they error out.
				// this affects the state engine in the following ways:
				// * connect, handshake will return with an error and continue to the next state
				// * handleSession will return with an error and continue to the next state
				switch action {
				case actionDisconnect:
					cl.Close()

				case actionRehandshake:
					if response, err := p.handshake(ctx, cl); err == nil {
						// reset the ticker
						tickerReset(time.Now(), response.Expiry, ticker)
					} else {
						// make a note of the error
						lastError = NewTimeError(err)
					}
				}

			case <-done:
				// exit the run() loop only when done is closed and ctx is canceled.
				// we don't want to stop processing events here even if the
				// session is SessionStatusDisconnected, until we are told to Stop().
				// * cancel will eventually close the done channel
				// * the Disconnect RPC call with NoRetry = true will eventually close the done channel
				//   but it will fire that action long after (disconnectDelay) we received the RPC call.
				//   The run() loop must still be running when it does, unless we are explicitely Stop()
				//   in which case we are protected by the running mutex.
				ticker.Stop()
				done = nil
				go func() {
					<-ctx.Done()
					close(ret)
				}()

			case <-ret:
				return
				// ¯\_(ツ)_/¯
			}
		}
	}()

	// initialize the for loop
	statuses <- SessionStatusConnecting
	return cancel
}

// connect sets up a new connection to a broker.
func (p *Provider) connect(ctx context.Context) (*client.Client, error) {
	// Dial a new connection
	p.configLock.RLock()
	tlsConfig := p.config.HCPConfig.SCADATLSConfig()
	scadaAddress := p.config.HCPConfig.SCADAAddress()
	p.configLock.RUnlock()

	opts := client.Opts{
		Dialer: &tcp.Dialer{
			TLSConfig: tlsConfig,
		},
		LogOutput: p.logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true}),
	}
	client, err := client.DialOptsContext(ctx, scadaAddress, &opts)
	if err != nil {
		p.logger.Error("failed to dial SCADA endpoint", "error", err)
		return nil, err
	}

	return client, nil
}

// handshake does the initial handshake. Handshake will return prefixed errors in certain scenarios:
//   - if HCPConfig.Token() returns *oauth2.RetrieveError, it will prefix the error with ErrInvalidCredentials
//   - client.RPC("Session.Handshake") might prefix the error with ErrPermissionDenied
//
// The prefixes are processed in NewTimeError called from the run() loop.
func (p *Provider) handshake(ctx context.Context, client *client.Client) (resp *types.HandshakeResponse, err error) {
	defer func() {
		if err != nil {
			p.logger.Error("handshake failed", "error", err)
		}
	}()

	// Build the set of capabilities based on the registered handlers.
	p.handlersLock.RLock()
	capabilities := make(map[string]int, len(p.handlers))
	for h := range p.handlers {
		capabilities[h] = 1
	}
	p.handlersLock.RUnlock()

	// determine configuration values
	p.configLock.RLock()
	service := p.config.Service
	resource := &p.config.Resource
	var oauthToken *oauth2.Token
	oauthToken, err = p.config.HCPConfig.Token()
	p.configLock.RUnlock()
	if err != nil {
		client.Close()
		err = PrefixError("failed to get access token", err)
		return nil, err
	}

	// make sure nobody is writing to the
	// meta map while client.RPC is reading from it
	p.metaLock.RLock()
	defer p.metaLock.RUnlock()

	req := types.HandshakeRequest{
		Service:  service,
		Resource: resource,

		AccessToken: oauthToken.AccessToken,

		// TODO: remove once it is not required anymore.
		ServiceVersion: "0.0.1",

		Capabilities: capabilities,
		Meta:         p.meta,
	}
	resp = new(types.HandshakeResponse)
	if err := client.RPC("Session.Handshake", &req, resp); err != nil {
		client.Close()
		return nil, err
	}

	if resp != nil && resp.SessionID != "" {
		p.logger.Debug("assigned session ID", "id", resp.SessionID)
	}
	if resp != nil && !resp.Authenticated {
		p.logger.Warn("authentication failed", "reason", resp.Reason)
	}

	return resp, nil
}

// handleSession is used to handle an established session.
func (p *Provider) handleSession(ctx context.Context, yamux net.Listener) error {
	var done = make(chan bool)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		// make the other go routine return
		// if yamux.Accept() errors out
		defer close(done)
		defer yamux.Close()
		for {
			if conn, err := yamux.Accept(); err != nil {
				select {
				case <-ctx.Done():
					// Do not log an error if we are shutting down
				default:
					p.logger.Error("failed to accept connection", "error", err)
				}
				return err
			} else {
				p.logger.Debug("accepted connection")
				go p.handleConnection(ctx, conn)
			}
		}
	})

	g.Go(func() error {
		// return nil here so that g.Wait()
		// always picks the error the Accept() routine
		// returned.
		for {
			select {
			case <-done:
				// the other go routine returned with an error
				// and closed the yamux client
				return nil

			case <-ctx.Done():
				// make the other go routine return
				// if ctx is canceled()
				yamux.Close()
				return nil
			}
		}
	})

	return g.Wait()
}

// handleConnection handles an incoming connection.
func (p *Provider) handleConnection(ctx context.Context, conn net.Conn) {
	// Create an RPC server to handle inbound
	pe := &providerEndpoint{p: p}
	rpcServer := rpc.NewServer()
	_ = rpcServer.RegisterName("Provider", pe)
	rpcCodec := msgpackrpc.NewCodec(false, false, conn)

	defer func() {
		if !pe.hijacked() {
			conn.Close()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := rpcServer.ServeRequest(rpcCodec); err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "closed") {
				p.logger.Error("RPC error", "error", err)
			}
			return
		}

		// Handle potential hijack in Provider.Connect
		if pe.hijacked() {
			cb := pe.getHijack()
			cb(conn)
			return
		}
	}
}

// wait is used to delay dialing on an error.
// it will return an error if the connection should not be
// retried.
func (p *Provider) wait(ctx context.Context) error {
	// Compute the backoff time
	backoff, noRetry := p.backoffDuration()
	// is this a no retry situation?
	if noRetry {
		return errNoRetry
	}

	// Setup a wait timer
	var wait <-chan time.Time
	if backoff > 0 {
		backoff = backoff + time.Duration(rand.Uint32())%backoff
		p.logger.Debug("backing off", "seconds", backoff.Seconds())
		wait = time.After(backoff)
	}

	// Wait until timer or shutdown
	select {
	case <-wait:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// backoffDuration is used to compute the next backoff duration.
// it returns the backoff time to wait for and a bool that will be
// set to true if no retries should be attempted.
func (p *Provider) backoffDuration() (time.Duration, bool) {
	// Use the default backoff
	backoff := defaultBackoff

	// Check for a server specified backoff
	p.backoffLock.Lock()
	providerBackoff := p.backoff
	noRetry := p.noRetry
	p.backoffLock.Unlock()
	if providerBackoff != 0 {
		backoff = providerBackoff
	}
	if noRetry {
		backoff = 0
	}

	// Use the test backoff
	p.configLock.RLock()
	testBackoff := p.config.TestBackoff
	p.configLock.RUnlock()
	if testBackoff != 0 {
		backoff = testBackoff
	}

	return backoff, noRetry
}

func (p *Provider) backoffReset() {
	// Reset the previous backoff
	p.backoffLock.Lock()
	p.noRetry = false
	p.backoff = 0
	p.backoffLock.Unlock()
}

// tickerReset resets ticker's period's to expiry-time.Now(). If the value of expiry is zero, it
// will return expiryDefault. If the value of expiry is before now, it will return expiryDefault.
// It applies expiryFactor to calculated duration before returning.
// for example, duration = 60s will return 54s with an expiryFactor of 0.90.
// note that this function will return incorrect results for expiry times smaller than 2 seconds.
func tickerReset(now, expiry time.Time, ticker *time.Ticker) time.Duration {
	// reject expiry time zero
	if expiry.IsZero() {
		return calculateExpiryFactor(expiryDefault)
	}
	// reject expiry time in the past
	if expiry.Before(now) {
		return calculateExpiryFactor(expiryDefault)
	}
	// calculate expiry-time.Now()
	d := expiry.Sub(now)
	// calculate d after expiryFactor
	d = calculateExpiryFactor(d)
	// reset the ticker
	ticker.Reset(d)

	return d
}

// calculateExpiryFactor multiplies d by expiryFactor and
// returns the multiplied time.Duration.
func calculateExpiryFactor(d time.Duration) time.Duration {
	var seconds = d.Seconds()
	var factored = seconds * expiryFactor
	d = time.Duration(factored) * time.Second
	return d
}

// UpdateConfig overwrites the provider's configuration
// with the given configuration.
func (p *Provider) UpdateConfig(config *Config) error {
	p.configLock.Lock()
	defer p.configLock.Unlock()

	if err := config.Validate(); err != nil {
		return err
	}
	p.config = config
	return nil
}
