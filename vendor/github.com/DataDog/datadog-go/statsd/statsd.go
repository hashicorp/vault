// Copyright 2013 Ooyala, Inc.

/*
Package statsd provides a Go dogstatsd client. Dogstatsd extends the popular statsd,
adding tags and histograms and pushing upstream to Datadog.

Refer to http://docs.datadoghq.com/guides/dogstatsd/ for information about DogStatsD.

statsd is based on go-statsd-client.
*/
package statsd

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
OptimalUDPPayloadSize defines the optimal payload size for a UDP datagram, 1432 bytes
is optimal for regular networks with an MTU of 1500 so datagrams don't get
fragmented. It's generally recommended not to fragment UDP datagrams as losing
a single fragment will cause the entire datagram to be lost.
*/
const OptimalUDPPayloadSize = 1432

/*
MaxUDPPayloadSize defines the maximum payload size for a UDP datagram.
Its value comes from the calculation: 65535 bytes Max UDP datagram size -
8byte UDP header - 60byte max IP headers
any number greater than that will see frames being cut out.
*/
const MaxUDPPayloadSize = 65467

// DefaultUDPBufferPoolSize is the default size of the buffer pool for UDP clients.
const DefaultUDPBufferPoolSize = 2048

// DefaultUDSBufferPoolSize is the default size of the buffer pool for UDS clients.
const DefaultUDSBufferPoolSize = 512

/*
DefaultMaxAgentPayloadSize is the default maximum payload size the agent
can receive. This can be adjusted by changing dogstatsd_buffer_size in the
agent configuration file datadog.yaml. This is also used as the optimal payload size
for UDS datagrams.
*/
const DefaultMaxAgentPayloadSize = 8192

/*
UnixAddressPrefix holds the prefix to use to enable Unix Domain Socket
traffic instead of UDP.
*/
const UnixAddressPrefix = "unix://"

/*
ddEnvTagsMapping is a mapping of each "DD_" prefixed environment variable
to a specific tag name.
*/
var ddEnvTagsMapping = map[string]string{
	// Client-side entity ID injection for container tagging.
	"DD_ENTITY_ID": "dd.internal.entity_id",
	// The name of the env in which the service runs.
	"DD_ENV": "env",
	// The name of the running service.
	"DD_SERVICE": "service",
	// The current version of the running service.
	"DD_VERSION": "version",
}

type metricType int

const (
	gauge metricType = iota
	count
	histogram
	histogramAggregated
	distribution
	distributionAggregated
	set
	timing
	timingAggregated
	event
	serviceCheck
)

type ReceivingMode int

const (
	MutexMode ReceivingMode = iota
	ChannelMode
)

const (
	WriterNameUDP string = "udp"
	WriterNameUDS string = "uds"
)

type metric struct {
	metricType metricType
	namespace  string
	globalTags []string
	name       string
	fvalue     float64
	fvalues    []float64
	ivalue     int64
	svalue     string
	evalue     *Event
	scvalue    *ServiceCheck
	tags       []string
	stags      string
	rate       float64
}

type noClientErr string

// ErrNoClient is returned if statsd reporting methods are invoked on
// a nil client.
const ErrNoClient = noClientErr("statsd client is nil")

func (e noClientErr) Error() string {
	return string(e)
}

// ClientInterface is an interface that exposes the common client functions for the
// purpose of being able to provide a no-op client or even mocking. This can aid
// downstream users' with their testing.
type ClientInterface interface {
	// Gauge measures the value of a metric at a particular time.
	Gauge(name string, value float64, tags []string, rate float64) error

	// Count tracks how many times something happened per second.
	Count(name string, value int64, tags []string, rate float64) error

	// Histogram tracks the statistical distribution of a set of values on each host.
	Histogram(name string, value float64, tags []string, rate float64) error

	// Distribution tracks the statistical distribution of a set of values across your infrastructure.
	Distribution(name string, value float64, tags []string, rate float64) error

	// Decr is just Count of -1
	Decr(name string, tags []string, rate float64) error

	// Incr is just Count of 1
	Incr(name string, tags []string, rate float64) error

	// Set counts the number of unique elements in a group.
	Set(name string, value string, tags []string, rate float64) error

	// Timing sends timing information, it is an alias for TimeInMilliseconds
	Timing(name string, value time.Duration, tags []string, rate float64) error

	// TimeInMilliseconds sends timing information in milliseconds.
	// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
	TimeInMilliseconds(name string, value float64, tags []string, rate float64) error

	// Event sends the provided Event.
	Event(e *Event) error

	// SimpleEvent sends an event with the provided title and text.
	SimpleEvent(title, text string) error

	// ServiceCheck sends the provided ServiceCheck.
	ServiceCheck(sc *ServiceCheck) error

	// SimpleServiceCheck sends an serviceCheck with the provided name and status.
	SimpleServiceCheck(name string, status ServiceCheckStatus) error

	// Close the client connection.
	Close() error

	// Flush forces a flush of all the queued dogstatsd payloads.
	Flush() error

	// SetWriteTimeout allows the user to set a custom write timeout.
	SetWriteTimeout(d time.Duration) error
}

// A Client is a handle for sending messages to dogstatsd.  It is safe to
// use one Client from multiple goroutines simultaneously.
type Client struct {
	// Sender handles the underlying networking protocol
	sender *sender
	// Namespace to prepend to all statsd calls
	Namespace string
	// Tags are global tags to be added to every statsd call
	Tags []string
	// skipErrors turns off error passing and allows UDS to emulate UDP behaviour
	SkipErrors  bool
	flushTime   time.Duration
	metrics     *ClientMetrics
	telemetry   *telemetryClient
	stop        chan struct{}
	wg          sync.WaitGroup
	workers     []*worker
	closerLock  sync.Mutex
	receiveMode ReceivingMode
	agg         *aggregator
	aggHistDist *aggregator
	options     []Option
	addrOption  string
}

// ClientMetrics contains metrics about the client
type ClientMetrics struct {
	TotalMetrics             uint64
	TotalMetricsGauge        uint64
	TotalMetricsCount        uint64
	TotalMetricsHistogram    uint64
	TotalMetricsDistribution uint64
	TotalMetricsSet          uint64
	TotalMetricsTiming       uint64
	TotalEvents              uint64
	TotalServiceChecks       uint64
	TotalDroppedOnReceive    uint64
}

// Verify that Client implements the ClientInterface.
// https://golang.org/doc/faq#guarantee_satisfies_interface
var _ ClientInterface = &Client{}

func resolveAddr(addr string) (statsdWriter, string, error) {
	if !strings.HasPrefix(addr, UnixAddressPrefix) {
		w, err := newUDPWriter(addr)
		return w, WriterNameUDP, err
	}

	w, err := newUDSWriter(addr[len(UnixAddressPrefix):])
	return w, WriterNameUDS, err
}

// New returns a pointer to a new Client given an addr in the format "hostname:port" or
// "unix:///path/to/socket".
func New(addr string, options ...Option) (*Client, error) {
	o, err := resolveOptions(options)
	if err != nil {
		return nil, err
	}

	w, writerType, err := resolveAddr(addr)
	if err != nil {
		return nil, err
	}

	client, err := newWithWriter(w, o, writerType)
	if err == nil {
		client.options = append(client.options, options...)
		client.addrOption = addr
	}
	return client, err
}

// NewWithWriter creates a new Client with given writer. Writer is a
// io.WriteCloser + SetWriteTimeout(time.Duration) error
func NewWithWriter(w statsdWriter, options ...Option) (*Client, error) {
	o, err := resolveOptions(options)
	if err != nil {
		return nil, err
	}
	return newWithWriter(w, o, "custom")
}

// CloneWithExtraOptions create a new Client with extra options
func CloneWithExtraOptions(c *Client, options ...Option) (*Client, error) {
	if c == nil {
		return nil, ErrNoClient
	}

	if c.addrOption == "" {
		return nil, fmt.Errorf("can't clone client with no addrOption")
	}
	opt := append(c.options, options...)
	return New(c.addrOption, opt...)
}

func newWithWriter(w statsdWriter, o *Options, writerName string) (*Client, error) {

	w.SetWriteTimeout(o.WriteTimeoutUDS)

	c := Client{
		Namespace: o.Namespace,
		Tags:      o.Tags,
		metrics:   &ClientMetrics{},
	}
	if o.Aggregation || o.ExtendedAggregation {
		c.agg = newAggregator(&c)
		c.agg.start(o.AggregationFlushInterval)

		if o.ExtendedAggregation {
			c.aggHistDist = c.agg
		}
	}

	// Inject values of DD_* environment variables as global tags.
	for envName, tagName := range ddEnvTagsMapping {
		if value := os.Getenv(envName); value != "" {
			c.Tags = append(c.Tags, fmt.Sprintf("%s:%s", tagName, value))
		}
	}

	if o.MaxBytesPerPayload == 0 {
		if writerName == WriterNameUDS {
			o.MaxBytesPerPayload = DefaultMaxAgentPayloadSize
		} else {
			o.MaxBytesPerPayload = OptimalUDPPayloadSize
		}
	}
	if o.BufferPoolSize == 0 {
		if writerName == WriterNameUDS {
			o.BufferPoolSize = DefaultUDSBufferPoolSize
		} else {
			o.BufferPoolSize = DefaultUDPBufferPoolSize
		}
	}
	if o.SenderQueueSize == 0 {
		if writerName == WriterNameUDS {
			o.SenderQueueSize = DefaultUDSBufferPoolSize
		} else {
			o.SenderQueueSize = DefaultUDPBufferPoolSize
		}
	}

	bufferPool := newBufferPool(o.BufferPoolSize, o.MaxBytesPerPayload, o.MaxMessagesPerPayload)
	c.sender = newSender(w, o.SenderQueueSize, bufferPool)
	c.receiveMode = o.ReceiveMode
	for i := 0; i < o.BufferShardCount; i++ {
		w := newWorker(bufferPool, c.sender)
		c.workers = append(c.workers, w)
		if c.receiveMode == ChannelMode {
			w.startReceivingMetric(o.ChannelModeBufferSize)
		}
	}

	c.flushTime = o.BufferFlushInterval
	c.stop = make(chan struct{}, 1)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.watch()
	}()

	if o.Telemetry {
		if o.TelemetryAddr == "" {
			c.telemetry = newTelemetryClient(&c, writerName, o.DevMode)
		} else {
			var err error
			c.telemetry, err = newTelemetryClientWithCustomAddr(&c, writerName, o.DevMode, o.TelemetryAddr, bufferPool)
			if err != nil {
				return nil, err
			}
		}
		c.telemetry.run(&c.wg, c.stop)
	}

	return &c, nil
}

// NewBuffered returns a Client that buffers its output and sends it in chunks.
// Buflen is the length of the buffer in number of commands.
//
// When addr is empty, the client will default to a UDP client and use the DD_AGENT_HOST
// and (optionally) the DD_DOGSTATSD_PORT environment variables to build the target address.
func NewBuffered(addr string, buflen int) (*Client, error) {
	return New(addr, WithMaxMessagesPerPayload(buflen))
}

// SetWriteTimeout allows the user to set a custom UDS write timeout. Not supported for UDP.
func (c *Client) SetWriteTimeout(d time.Duration) error {
	if c == nil {
		return ErrNoClient
	}
	return c.sender.transport.SetWriteTimeout(d)
}

func (c *Client) watch() {
	ticker := time.NewTicker(c.flushTime)

	for {
		select {
		case <-ticker.C:
			for _, w := range c.workers {
				w.flush()
			}
		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}

// Flush forces a flush of all the queued dogstatsd payloads This method is
// blocking and will not return until everything is sent through the network.
// In MutexMode, this will also block sampling new data to the client while the
// workers and sender are flushed.
func (c *Client) Flush() error {
	if c == nil {
		return ErrNoClient
	}
	if c.agg != nil {
		c.agg.sendMetrics()
	}
	for _, w := range c.workers {
		w.pause()
		defer w.unpause()
		w.flushUnsafe()
	}
	// Now that the worker are pause the sender can flush the queue between
	// worker and senders
	c.sender.flush()
	return nil
}

func (c *Client) FlushTelemetryMetrics() ClientMetrics {
	cm := ClientMetrics{
		TotalMetricsGauge:        atomic.SwapUint64(&c.metrics.TotalMetricsGauge, 0),
		TotalMetricsCount:        atomic.SwapUint64(&c.metrics.TotalMetricsCount, 0),
		TotalMetricsSet:          atomic.SwapUint64(&c.metrics.TotalMetricsSet, 0),
		TotalMetricsHistogram:    atomic.SwapUint64(&c.metrics.TotalMetricsHistogram, 0),
		TotalMetricsDistribution: atomic.SwapUint64(&c.metrics.TotalMetricsDistribution, 0),
		TotalMetricsTiming:       atomic.SwapUint64(&c.metrics.TotalMetricsTiming, 0),
		TotalEvents:              atomic.SwapUint64(&c.metrics.TotalEvents, 0),
		TotalServiceChecks:       atomic.SwapUint64(&c.metrics.TotalServiceChecks, 0),
		TotalDroppedOnReceive:    atomic.SwapUint64(&c.metrics.TotalDroppedOnReceive, 0),
	}

	cm.TotalMetrics = cm.TotalMetricsGauge + cm.TotalMetricsCount +
		cm.TotalMetricsSet + cm.TotalMetricsHistogram +
		cm.TotalMetricsDistribution + cm.TotalMetricsTiming

	return cm
}

func (c *Client) send(m metric) error {
	if c == nil {
		return ErrNoClient
	}

	m.globalTags = c.Tags
	m.namespace = c.Namespace

	h := hashString32(m.name)
	worker := c.workers[h%uint32(len(c.workers))]

	if c.receiveMode == ChannelMode {
		select {
		case worker.inputMetrics <- m:
		default:
			atomic.AddUint64(&c.metrics.TotalDroppedOnReceive, 1)
		}
		return nil
	}
	return worker.processMetric(m)
}

// Gauge measures the value of a metric at a particular time.
func (c *Client) Gauge(name string, value float64, tags []string, rate float64) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalMetricsGauge, 1)
	if c.agg != nil {
		return c.agg.gauge(name, value, tags)
	}
	return c.send(metric{metricType: gauge, name: name, fvalue: value, tags: tags, rate: rate})
}

// Count tracks how many times something happened per second.
func (c *Client) Count(name string, value int64, tags []string, rate float64) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalMetricsCount, 1)
	if c.agg != nil {
		return c.agg.count(name, value, tags)
	}
	return c.send(metric{metricType: count, name: name, ivalue: value, tags: tags, rate: rate})
}

// Histogram tracks the statistical distribution of a set of values on each host.
func (c *Client) Histogram(name string, value float64, tags []string, rate float64) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalMetricsHistogram, 1)
	if c.aggHistDist != nil {
		return c.agg.histogram(name, value, tags)
	}
	return c.send(metric{metricType: histogram, name: name, fvalue: value, tags: tags, rate: rate})
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func (c *Client) Distribution(name string, value float64, tags []string, rate float64) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalMetricsDistribution, 1)
	if c.aggHistDist != nil {
		return c.agg.distribution(name, value, tags)
	}
	return c.send(metric{metricType: distribution, name: name, fvalue: value, tags: tags, rate: rate})
}

// Decr is just Count of -1
func (c *Client) Decr(name string, tags []string, rate float64) error {
	return c.Count(name, -1, tags, rate)
}

// Incr is just Count of 1
func (c *Client) Incr(name string, tags []string, rate float64) error {
	return c.Count(name, 1, tags, rate)
}

// Set counts the number of unique elements in a group.
func (c *Client) Set(name string, value string, tags []string, rate float64) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalMetricsSet, 1)
	if c.agg != nil {
		return c.agg.set(name, value, tags)
	}
	return c.send(metric{metricType: set, name: name, svalue: value, tags: tags, rate: rate})
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func (c *Client) Timing(name string, value time.Duration, tags []string, rate float64) error {
	return c.TimeInMilliseconds(name, value.Seconds()*1000, tags, rate)
}

// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
func (c *Client) TimeInMilliseconds(name string, value float64, tags []string, rate float64) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalMetricsTiming, 1)
	if c.aggHistDist != nil {
		return c.agg.timing(name, value, tags)
	}
	return c.send(metric{metricType: timing, name: name, fvalue: value, tags: tags, rate: rate})
}

// Event sends the provided Event.
func (c *Client) Event(e *Event) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalEvents, 1)
	return c.send(metric{metricType: event, evalue: e, rate: 1})
}

// SimpleEvent sends an event with the provided title and text.
func (c *Client) SimpleEvent(title, text string) error {
	e := NewEvent(title, text)
	return c.Event(e)
}

// ServiceCheck sends the provided ServiceCheck.
func (c *Client) ServiceCheck(sc *ServiceCheck) error {
	if c == nil {
		return ErrNoClient
	}
	atomic.AddUint64(&c.metrics.TotalServiceChecks, 1)
	return c.send(metric{metricType: serviceCheck, scvalue: sc, rate: 1})
}

// SimpleServiceCheck sends an serviceCheck with the provided name and status.
func (c *Client) SimpleServiceCheck(name string, status ServiceCheckStatus) error {
	sc := NewServiceCheck(name, status)
	return c.ServiceCheck(sc)
}

// Close the client connection.
func (c *Client) Close() error {
	if c == nil {
		return ErrNoClient
	}

	// Acquire closer lock to ensure only one thread can close the stop channel
	c.closerLock.Lock()
	defer c.closerLock.Unlock()

	// Notify all other threads that they should stop
	select {
	case <-c.stop:
		return nil
	default:
	}
	close(c.stop)

	if c.receiveMode == ChannelMode {
		for _, w := range c.workers {
			w.stopReceivingMetric()
		}
	}

	// Wait for the threads to stop
	c.wg.Wait()

	// Finally flush any remaining metrics that may have come in at the last moment
	if c.agg != nil {
		c.agg.stop()
	}
	c.Flush()

	return c.sender.close()
}
