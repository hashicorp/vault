package monitor

import (
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/hashicorp/go-hclog"
)

// Monitor provides a mechanism to stream logs using go-hclog
// InterceptLogger and SinkAdapter. It allows streaming of logs
// at a different log level than what is set on the logger.
type Monitor interface {
	// Start returns a channel of log messages which are sent
	// ever time a log message occurs
	Start() <-chan []byte

	// Stop de-registers the sink from the InterceptLogger
	// and closes the log channels
	Stop()
}

// monitor implements the Monitor interface
type monitor struct {
	sink log.SinkAdapter

	// logger is the logger we will be monitoring
	logger log.InterceptLogger

	// logCh is a buffered chan where we send logs when streaming
	logCh chan []byte

	// doneCh coordinates the shutdown of logCh
	doneCh chan struct{}

	// droppedCount is the current count of messages
	// that were dropped from the logCh buffer.
	droppedCount uint64
	bufSize      int
	// dropCheckInterval is the amount of time we should
	// wait to check for dropped messages. Defaults
	// to 3 seconds
	dropCheckInterval time.Duration
}

// NewMonitor creates a new Monitor. Start must be called in order to actually start
// streaming logs. buf is the buffer size of the channel that sends log messages.
func NewMonitor(buf int, logger log.InterceptLogger, opts *log.LoggerOptions) Monitor {
	return newMonitor(buf, logger, opts)
}

func newMonitor(buf int, logger log.InterceptLogger, opts *log.LoggerOptions) *monitor {
	sw := &monitor{
		logger:            logger,
		logCh:             make(chan []byte, buf),
		doneCh:            make(chan struct{}),
		bufSize:           buf,
		dropCheckInterval: 3 * time.Second,
	}

	opts.Output = sw
	sink := log.NewSinkAdapter(opts)
	sw.sink = sink

	return sw
}

// Stop deregisters the sink and stops the monitoring process
func (d *monitor) Stop() {
	d.logger.DeregisterSink(d.sink)
	close(d.doneCh)
}

// Start registers a sink on the monitor's logger and starts sending
// received log messages over the returned channel.
func (d *monitor) Start() <-chan []byte {
	// register our sink with the logger
	d.logger.RegisterSink(d.sink)

	streamCh := make(chan []byte, d.bufSize)

	// Run a go routine that listens for streamed
	// log messages and sends them to streamCh.
	//
	// It also periodically checks for dropped
	// messages and makes room on the logCh to add
	// a dropped message count warning
	go func() {
		defer close(streamCh)

		ticker := time.NewTicker(d.dropCheckInterval)
		defer ticker.Stop()

		var logMessage []byte
		for {
			// Reset the byte slice on every loop iteration, which is what makes
			// the below for loop work.
			logMessage = nil

			select {
			case <-ticker.C:
				// Check if there have been any dropped messages.
				dc := atomic.LoadUint64(&d.droppedCount)

				if dc > 0 {
					logMessage = []byte(fmt.Sprintf("[WARN] Monitor dropped %d logs during monitor request\n", dc))
					atomic.SwapUint64(&d.droppedCount, 0)
				}
			case logMessage = <-d.logCh:
			case <-d.doneCh:
				return
			}

			if len(logMessage) > 0 {
				select {
				case <-d.doneCh:
					return
				case streamCh <- logMessage:
				}
			}
		}
	}()

	return streamCh
}

// Write attempts to send latest log to logCh
// it drops the log if channel is unavailable to receive
func (d *monitor) Write(p []byte) (n int, err error) {
	// ensure logCh is still open
	select {
	case <-d.doneCh:
		return
	default:
	}

	bytes := make([]byte, len(p))
	copy(bytes, p)

	select {
	case d.logCh <- bytes:
	default:
		atomic.AddUint64(&d.droppedCount, 1)
	}

	return len(p), nil
}
