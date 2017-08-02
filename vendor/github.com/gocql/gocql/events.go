package gocql

import (
	"net"
	"sync"
	"time"
)

type eventDebouncer struct {
	name   string
	timer  *time.Timer
	mu     sync.Mutex
	events []frame

	callback func([]frame)
	quit     chan struct{}
}

func newEventDebouncer(name string, eventHandler func([]frame)) *eventDebouncer {
	e := &eventDebouncer{
		name:     name,
		quit:     make(chan struct{}),
		timer:    time.NewTimer(eventDebounceTime),
		callback: eventHandler,
	}
	e.timer.Stop()
	go e.flusher()

	return e
}

func (e *eventDebouncer) stop() {
	e.quit <- struct{}{} // sync with flusher
	close(e.quit)
}

func (e *eventDebouncer) flusher() {
	for {
		select {
		case <-e.timer.C:
			e.mu.Lock()
			e.flush()
			e.mu.Unlock()
		case <-e.quit:
			return
		}
	}
}

const (
	eventBufferSize   = 1000
	eventDebounceTime = 1 * time.Second
)

// flush must be called with mu locked
func (e *eventDebouncer) flush() {
	if len(e.events) == 0 {
		return
	}

	// if the flush interval is faster than the callback then we will end up calling
	// the callback multiple times, probably a bad idea. In this case we could drop
	// frames?
	go e.callback(e.events)
	e.events = make([]frame, 0, eventBufferSize)
}

func (e *eventDebouncer) debounce(frame frame) {
	e.mu.Lock()
	e.timer.Reset(eventDebounceTime)

	// TODO: probably need a warning to track if this threshold is too low
	if len(e.events) < eventBufferSize {
		e.events = append(e.events, frame)
	} else {
		Logger.Printf("%s: buffer full, dropping event frame: %s", e.name, frame)
	}

	e.mu.Unlock()
}

func (s *Session) handleEvent(framer *framer) {
	// TODO(zariel): need to debounce events frames, and possible also events
	defer framerPool.Put(framer)

	frame, err := framer.parseFrame()
	if err != nil {
		// TODO: logger
		Logger.Printf("gocql: unable to parse event frame: %v\n", err)
		return
	}

	if gocqlDebug {
		Logger.Printf("gocql: handling frame: %v\n", frame)
	}

	// TODO: handle medatadata events
	switch f := frame.(type) {
	case *schemaChangeKeyspace, *schemaChangeFunction, *schemaChangeTable:
		s.schemaEvents.debounce(frame)
	case *topologyChangeEventFrame, *statusChangeEventFrame:
		s.nodeEvents.debounce(frame)
	default:
		Logger.Printf("gocql: invalid event frame (%T): %v\n", f, f)
	}
}

func (s *Session) handleSchemaEvent(frames []frame) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.schemaDescriber == nil {
		return
	}
	for _, frame := range frames {
		switch f := frame.(type) {
		case *schemaChangeKeyspace:
			s.schemaDescriber.clearSchema(f.keyspace)
		case *schemaChangeTable:
			s.schemaDescriber.clearSchema(f.keyspace)
		}
	}
}

func (s *Session) handleNodeEvent(frames []frame) {
	type nodeEvent struct {
		change string
		host   net.IP
		port   int
	}

	events := make(map[string]*nodeEvent)

	for _, frame := range frames {
		// TODO: can we be sure the order of events in the buffer is correct?
		switch f := frame.(type) {
		case *topologyChangeEventFrame:
			event, ok := events[f.host.String()]
			if !ok {
				event = &nodeEvent{change: f.change, host: f.host, port: f.port}
				events[f.host.String()] = event
			}
			event.change = f.change

		case *statusChangeEventFrame:
			event, ok := events[f.host.String()]
			if !ok {
				event = &nodeEvent{change: f.change, host: f.host, port: f.port}
				events[f.host.String()] = event
			}
			event.change = f.change
		}
	}

	for _, f := range events {
		if gocqlDebug {
			Logger.Printf("gocql: dispatching event: %+v\n", f)
		}

		switch f.change {
		case "NEW_NODE":
			s.handleNewNode(f.host, f.port, true)
		case "REMOVED_NODE":
			s.handleRemovedNode(f.host, f.port)
		case "MOVED_NODE":
		// java-driver handles this, not mentioned in the spec
		// TODO(zariel): refresh token map
		case "UP":
			s.handleNodeUp(f.host, f.port, true)
		case "DOWN":
			s.handleNodeDown(f.host, f.port)
		}
	}
}

func (s *Session) handleNewNode(ip net.IP, port int, waitForBinary bool) {
	// Get host info and apply any filters to the host
	hostInfo, err := s.hostSource.GetHostInfo(ip, port)
	if err != nil {
		Logger.Printf("gocql: events: unable to fetch host info for (%s:%d): %v\n", ip, port, err)
		return
	}

	// If hostInfo is nil, this host was filtered out by cfg.HostFilter
	if hostInfo == nil {
		return
	}

	if t := hostInfo.Version().nodeUpDelay(); t > 0 && waitForBinary {
		time.Sleep(t)
	}

	// should this handle token moving?
	if existing, ok := s.ring.addHostIfMissing(hostInfo); ok {
		existing.update(hostInfo)
		hostInfo = existing
	}

	s.pool.addHost(hostInfo)
	s.policy.AddHost(hostInfo)
	hostInfo.setState(NodeUp)
	if s.control != nil && !s.cfg.IgnorePeerAddr {
		s.hostSource.refreshRing()
	}
}

func (s *Session) handleRemovedNode(ip net.IP, port int) {
	// we remove all nodes but only add ones which pass the filter
	host := s.ring.getHost(ip)
	if host == nil {
		host = &HostInfo{connectAddress: ip, port: port}
	}

	if s.cfg.HostFilter != nil && !s.cfg.HostFilter.Accept(host) {
		return
	}

	host.setState(NodeDown)
	s.policy.RemoveHost(host)
	s.pool.removeHost(ip)
	s.ring.removeHost(ip)

	if !s.cfg.IgnorePeerAddr {
		s.hostSource.refreshRing()
	}
}

func (s *Session) handleNodeUp(ip net.IP, port int, waitForBinary bool) {
	if gocqlDebug {
		Logger.Printf("gocql: Session.handleNodeUp: %s:%d\n", ip.String(), port)
	}

	host := s.ring.getHost(ip)
	if host != nil {
		// If we receive a node up event and user has asked us to ignore the peer address use
		// the address provide by the event instead the address provide by peer the table.
		if s.cfg.IgnorePeerAddr && !host.ConnectAddress().Equal(ip) {
			host.SetConnectAddress(ip)
		}

		if s.cfg.HostFilter != nil && !s.cfg.HostFilter.Accept(host) {
			return
		}

		if t := host.Version().nodeUpDelay(); t > 0 && waitForBinary {
			time.Sleep(t)
		}

		s.pool.hostUp(host)
		s.policy.HostUp(host)
		host.setState(NodeUp)
		return
	}

	s.handleNewNode(ip, port, waitForBinary)
}

func (s *Session) handleNodeDown(ip net.IP, port int) {
	if gocqlDebug {
		Logger.Printf("gocql: Session.handleNodeDown: %s:%d\n", ip.String(), port)
	}

	host := s.ring.getHost(ip)
	if host == nil {
		host = &HostInfo{connectAddress: ip, port: port}
	}

	if s.cfg.HostFilter != nil && !s.cfg.HostFilter.Accept(host) {
		return
	}

	host.setState(NodeDown)
	s.policy.HostDown(host)
	s.pool.hostDown(ip)
}
