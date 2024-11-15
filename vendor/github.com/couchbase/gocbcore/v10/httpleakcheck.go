package gocbcore

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var leakTrackingEnabled uint32 = 0
var trackedRespsLock sync.Mutex
var trackedResps []*leakTrackingReadCloser

// EnableHttpResponseTracking enables tracking response bodies to ensure that they are
// eventually closed.
//
// # VOLATILE
//
// This API is VOLATILE is subject to change at any time.
func EnableHttpResponseTracking() {
	atomic.StoreUint32(&leakTrackingEnabled, 1)
}

func wrapHttpResponse(resp *http.Response) *http.Response {
	if atomic.LoadUint32(&leakTrackingEnabled) == 0 {
		return resp
	}

	trackingBody := &leakTrackingReadCloser{
		parent:     resp.Body,
		stackTrace: debug.Stack(),
	}

	trackedRespsLock.Lock()
	trackedResps = append(trackedResps, trackingBody)
	trackedRespsLock.Unlock()

	resp.Body = trackingBody
	return resp
}

func removeTrackedHttpBodyRecord(l *leakTrackingReadCloser) {
	trackedRespsLock.Lock()
	recordIdx := -1
	for i, tracked := range trackedResps {
		if tracked == l {
			recordIdx = i
		}
	}

	if recordIdx >= 0 {
		trackedResps = append(trackedResps[:recordIdx], trackedResps[recordIdx+1:]...)
	}
	trackedRespsLock.Unlock()
}

// ReportLeakedHttpResponses prints the stack traces of any response bodies that have not
// been closed. Returns true if all bodies have been closed, false otherwise.
//
// # VOLATILE
//
// This API is VOLATILE is subject to change at any time.
func ReportLeakedHttpResponses() bool {
	if len(trackedResps) == 0 {
		log.Printf("No leaked http requests")
		return true
	}

	log.Printf("Found %d leaked http requests", len(trackedResps))
	for _, leakRecord := range trackedResps {
		log.Printf("Leaked http request stack: %s", leakRecord.stackTrace)
	}

	return false
}

type leakTrackingReadCloser struct {
	parent     io.ReadCloser
	stackTrace []byte
}

func (l *leakTrackingReadCloser) Read(p []byte) (int, error) {
	n, err := l.parent.Read(p)
	if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
		removeTrackedHttpBodyRecord(l)
	}
	return n, err
}

func (l *leakTrackingReadCloser) Close() error {
	removeTrackedHttpBodyRecord(l)
	return l.parent.Close()
}

var _ io.ReadCloser = (*leakTrackingReadCloser)(nil)
