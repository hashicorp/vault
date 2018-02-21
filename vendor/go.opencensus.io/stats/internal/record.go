package internal

import (
	"time"

	"go.opencensus.io/tag"
)

type Recorder func(*tag.Map, time.Time, interface{})

// DefaultRecorder will be called for each Record call.
var DefaultRecorder Recorder = nil
