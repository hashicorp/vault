// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package logger

import (
	"encoding/json"
	"io"
	"math"
	"sync"
	"time"
)

// IOSink writes a JSON-encoded message to the io.Writer.
type IOSink struct {
	enc *json.Encoder

	// encMu protects the encoder from concurrent writes. While the logger
	// itself does not concurrently write to the sink, the sink may be used
	// concurrently within the driver.
	encMu sync.Mutex
}

// Compile-time check to ensure IOSink implements the LogSink interface.
var _ LogSink = &IOSink{}

// NewIOSink will create an IOSink object that writes JSON messages to the
// provided io.Writer.
func NewIOSink(out io.Writer) *IOSink {
	return &IOSink{
		enc: json.NewEncoder(out),
	}
}

// Info will write a JSON-encoded message to the io.Writer.
func (sink *IOSink) Info(_ int, msg string, keysAndValues ...interface{}) {
	mapSize := len(keysAndValues) / 2
	if math.MaxInt-mapSize >= 2 {
		mapSize += 2
	}
	kvMap := make(map[string]interface{}, mapSize)

	kvMap[KeyTimestamp] = time.Now().UnixNano()
	kvMap[KeyMessage] = msg

	for i := 0; i < len(keysAndValues); i += 2 {
		kvMap[keysAndValues[i].(string)] = keysAndValues[i+1]
	}

	sink.encMu.Lock()
	defer sink.encMu.Unlock()

	_ = sink.enc.Encode(kvMap)
}

// Error will write a JSON-encoded error message to the io.Writer.
func (sink *IOSink) Error(err error, msg string, kv ...interface{}) {
	kv = append(kv, KeyError, err.Error())
	sink.Info(0, msg, kv...)
}
