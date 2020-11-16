// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/api/global"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

const (
	DefaultMaxQueueSize       = 2048
	DefaultBatchTimeout       = 5000 * time.Millisecond
	DefaultMaxExportBatchSize = 512
)

type BatchSpanProcessorOption func(o *BatchSpanProcessorOptions)

type BatchSpanProcessorOptions struct {
	// MaxQueueSize is the maximum queue size to buffer spans for delayed processing. If the
	// queue gets full it drops the spans. Use BlockOnQueueFull to change this behavior.
	// The default value of MaxQueueSize is 2048.
	MaxQueueSize int

	// BatchTimeout is the maximum duration for constructing a batch. Processor
	// forcefully sends available spans when timeout is reached.
	// The default value of BatchTimeout is 5000 msec.
	BatchTimeout time.Duration

	// MaxExportBatchSize is the maximum number of spans to process in a single batch.
	// If there are more than one batch worth of spans then it processes multiple batches
	// of spans one batch after the other without any delay.
	// The default value of MaxExportBatchSize is 512.
	MaxExportBatchSize int

	// BlockOnQueueFull blocks onEnd() and onStart() method if the queue is full
	// AND if BlockOnQueueFull is set to true.
	// Blocking option should be used carefully as it can severely affect the performance of an
	// application.
	BlockOnQueueFull bool
}

// BatchSpanProcessor is a SpanProcessor that batches asynchronously received
// SpanData and sends it to a trace.Exporter when complete.
type BatchSpanProcessor struct {
	e export.SpanExporter
	o BatchSpanProcessorOptions

	queue   chan *export.SpanData
	dropped uint32

	batch      []*export.SpanData
	batchMutex sync.Mutex
	timer      *time.Timer
	stopWait   sync.WaitGroup
	stopOnce   sync.Once
	stopCh     chan struct{}
}

var _ SpanProcessor = (*BatchSpanProcessor)(nil)

// NewBatchSpanProcessor creates a new BatchSpanProcessor that will send
// SpanData batches to the exporters with the supplied options.
//
// The returned BatchSpanProcessor needs to be registered with the SDK using
// the RegisterSpanProcessor method for it to process spans.
//
// If the exporter is nil, the span processor will preform no action.
func NewBatchSpanProcessor(exporter export.SpanExporter, options ...BatchSpanProcessorOption) *BatchSpanProcessor {
	o := BatchSpanProcessorOptions{
		BatchTimeout:       DefaultBatchTimeout,
		MaxQueueSize:       DefaultMaxQueueSize,
		MaxExportBatchSize: DefaultMaxExportBatchSize,
	}
	for _, opt := range options {
		opt(&o)
	}
	bsp := &BatchSpanProcessor{
		e:      exporter,
		o:      o,
		batch:  make([]*export.SpanData, 0, o.MaxExportBatchSize),
		timer:  time.NewTimer(o.BatchTimeout),
		queue:  make(chan *export.SpanData, o.MaxQueueSize),
		stopCh: make(chan struct{}),
	}

	bsp.stopWait.Add(1)
	go func() {
		defer bsp.stopWait.Done()
		bsp.processQueue()
		bsp.drainQueue()
	}()

	return bsp
}

// OnStart method does nothing.
func (bsp *BatchSpanProcessor) OnStart(sd *export.SpanData) {}

// OnEnd method enqueues export.SpanData for later processing.
func (bsp *BatchSpanProcessor) OnEnd(sd *export.SpanData) {
	// Do not enqueue spans if we are just going to drop them.
	if bsp.e == nil {
		return
	}
	bsp.enqueue(sd)
}

// Shutdown flushes the queue and waits until all spans are processed.
// It only executes once. Subsequent call does nothing.
func (bsp *BatchSpanProcessor) Shutdown() {
	bsp.stopOnce.Do(func() {
		close(bsp.stopCh)
		bsp.stopWait.Wait()
	})
}

// ForceFlush exports all ended spans that have not yet been exported.
func (bsp *BatchSpanProcessor) ForceFlush() {
	bsp.exportSpans()
}

func WithMaxQueueSize(size int) BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.MaxQueueSize = size
	}
}

func WithMaxExportBatchSize(size int) BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.MaxExportBatchSize = size
	}
}

func WithBatchTimeout(delay time.Duration) BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.BatchTimeout = delay
	}
}

func WithBlocking() BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.BlockOnQueueFull = true
	}
}

// exportSpans is a subroutine of processing and draining the queue.
func (bsp *BatchSpanProcessor) exportSpans() {
	bsp.timer.Reset(bsp.o.BatchTimeout)

	bsp.batchMutex.Lock()
	defer bsp.batchMutex.Unlock()

	if len(bsp.batch) > 0 {
		if err := bsp.e.ExportSpans(context.Background(), bsp.batch); err != nil {
			global.Handle(err)
		}
		bsp.batch = bsp.batch[:0]
	}
}

// processQueue removes spans from the `queue` channel until processor
// is shut down. It calls the exporter in batches of up to MaxExportBatchSize
// waiting up to BatchTimeout to form a batch.
func (bsp *BatchSpanProcessor) processQueue() {
	defer bsp.timer.Stop()

	for {
		select {
		case <-bsp.stopCh:
			return
		case <-bsp.timer.C:
			bsp.exportSpans()
		case sd := <-bsp.queue:
			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) == bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()
			if shouldExport {
				if !bsp.timer.Stop() {
					<-bsp.timer.C
				}
				bsp.exportSpans()
			}
		}
	}
}

// drainQueue awaits the any caller that had added to bsp.stopWait
// to finish the enqueue, then exports the final batch.
func (bsp *BatchSpanProcessor) drainQueue() {
	for {
		select {
		case sd := <-bsp.queue:
			if sd == nil {
				bsp.exportSpans()
				return
			}

			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) == bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()

			if shouldExport {
				bsp.exportSpans()
			}
		default:
			close(bsp.queue)
		}
	}
}

func (bsp *BatchSpanProcessor) enqueue(sd *export.SpanData) {
	if !sd.SpanContext.IsSampled() {
		return
	}

	// This ensures the bsp.queue<- below does not panic as the
	// processor shuts down.
	defer func() {
		x := recover()
		switch err := x.(type) {
		case nil:
			return
		case runtime.Error:
			if err.Error() == "send on closed channel" {
				return
			}
		}
		panic(x)
	}()

	select {
	case <-bsp.stopCh:
		return
	default:
	}

	if bsp.o.BlockOnQueueFull {
		bsp.queue <- sd
		return
	}

	select {
	case bsp.queue <- sd:
	default:
		atomic.AddUint32(&bsp.dropped, 1)
	}
}
