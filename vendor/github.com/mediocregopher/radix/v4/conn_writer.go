package radix

import (
	"context"
	"fmt"
	"time"

	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

type connWriter struct {
	wCh    <-chan connMarshalerUnmarshaler
	rCh    chan<- connMarshalerUnmarshaler
	doneCh <-chan struct{}

	bw   resp.BufferedWriter
	opts *resp.Opts

	flushInterval time.Duration
	flushTicker   interface {
		Stop()
		Reset(time.Duration)
	}
	flushTickerCh <-chan time.Time

	// populated during run()
	flushBuf          []connMarshalerUnmarshaler
	flushTickerPaused bool

	// only used for tests, will be written to in the event loop
	eventLoopCh chan struct{}
}

func (c *conn) writer(ctx context.Context, flushInterval time.Duration) {
	cw := &connWriter{
		wCh:           c.wCh,
		rCh:           c.rCh,
		doneCh:        ctx.Done(),
		bw:            c.bw,
		opts:          c.wOpts,
		flushInterval: flushInterval,
	}

	if cw.flushInterval > 0 {
		ticker := time.NewTicker(cw.flushInterval)
		cw.flushTicker = ticker
		cw.flushTickerCh = ticker.C
	}

	cw.run()
}

func (cw *connWriter) pauseTicker() {
	if cw.flushTicker == nil || cw.flushTickerPaused {
		return
	}
	cw.flushTicker.Stop()
	cw.flushTickerPaused = true
}

func (cw *connWriter) resumeTicker() {
	if cw.flushTicker == nil || !cw.flushTickerPaused {
		return
	}
	cw.flushTicker.Reset(cw.flushInterval)
	cw.flushTickerPaused = false
}

// forwardToReader returns false if doneCh is closed.
func (cw *connWriter) forwardToReader(mu connMarshalerUnmarshaler) bool {
	select {
	case <-cw.doneCh:
		return false
	case cw.rCh <- mu:
		return true
	}
}

// write returns true if the write was successful.
func (cw *connWriter) write(mu connMarshalerUnmarshaler) bool {
	if err := mu.ctx.Err(); err != nil {
		mu.errCh <- resp.ErrConnUsable{Err: fmt.Errorf("checking context before write: %w", err)}
		return false

	} else if err := resp3.Marshal(cw.bw, mu.marshal, cw.opts); err != nil {
		mu.errCh <- fmt.Errorf("marshaling message to Conn: %w", err)
		return false
	}

	return true
}

// flush returns false if doneCh is closed.
func (cw *connWriter) flush() bool {
	if len(cw.flushBuf) == 0 {
		cw.pauseTicker()
		return true
	}

	flushBuf := cw.flushBuf[:0]
	for _, mu := range cw.flushBuf {
		if cw.write(mu) {
			flushBuf = append(flushBuf, mu)
		}

	}
	cw.flushBuf = cw.flushBuf[:0]

	if err := cw.bw.Flush(); err != nil {
		err = fmt.Errorf("flushing write buffer: %w", err)
		for _, mu := range flushBuf {
			mu.errCh <- err
		}

	} else {
		for _, mu := range flushBuf {
			// if there's no unmarshaler then don't forward to the reader
			if mu.unmarshalInto == nil {
				mu.errCh <- nil
			} else if !cw.forwardToReader(mu) {
				return false
			}
		}
	}

	cw.resumeTicker()
	return true
}

func (cw *connWriter) run() {
	cw.pauseTicker()
	for {
		select {
		case <-cw.doneCh:
			return
		case <-cw.flushTickerCh:
			if !cw.flush() {
				return
			}
		case <-cw.eventLoopCh:
			// do nothing, only used for tests
		case mu := <-cw.wCh:
			if mu.marshal != nil {
				cw.flushBuf = append(cw.flushBuf, mu)
				if (cw.flushInterval == 0 || cw.flushTickerPaused) && !cw.flush() {
					return
				}

				// if there's no marshaler then flush what's there so far before
				// forwarding, so that order can be preserved
			} else if !cw.flush() || !cw.forwardToReader(mu) {
				return
			}
		}
	}
}
