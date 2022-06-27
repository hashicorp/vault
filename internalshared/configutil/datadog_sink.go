package configutil

import (
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/datadog"
	"github.com/mitchellh/cli"
)

type dogStatsdSink interface {
	SetTags(tags []string)
	EnableHostNamePropagation()
	SetGauge(key []string, val float32)
	IncrCounter(key []string, val float32)
	EmitKey(key []string, val float32)
	AddSample(key []string, val float32)
	SetGaugeWithLabels(key []string, val float32, labels []metrics.Label)
	IncrCounterWithLabels(key []string, val float32, labels []metrics.Label)
	AddSampleWithLabels(key []string, val float32, labels []metrics.Label)
	Shutdown()
}

type DatadogSink struct {
	tags                 []string
	addr                 string
	hostName             string
	propagateHostname    bool
	sink                 dogStatsdSink
	logger               cli.Ui
	attemptedToConnectAt *time.Time
	lock                 sync.RWMutex
	creator              func(addr string, hostName string) (dogStatsdSink, error)
}

func NewDatadogSink(addr string, hostName string, logger cli.Ui) *DatadogSink {
	return &DatadogSink{
		addr:              addr,
		hostName:          hostName,
		propagateHostname: false,
		logger:            logger,
		creator: func(addr string, hostName string) (dogStatsdSink, error) {
			return datadog.NewDogStatsdSink(addr, hostName)
		},
	}
}

func (s *DatadogSink) SetTags(tags []string) {
	sink := s.getSink()

	if sink == nil {
		s.lock.Lock()
		defer s.lock.Unlock()
		sink = s.sink
		if sink == nil { // store the tags in object if we can't connect to the server
			s.tags = tags
			return
		}
	}

	sink.SetTags(tags)
}

func (s *DatadogSink) EnableHostNamePropagation() {
	sink := s.getSink()

	if sink == nil {
		s.lock.Lock()
		defer s.lock.Unlock()
		sink = s.sink
		if sink == nil { // store the flag in object if we can't connect to the server
			s.propagateHostname = true
			return
		}
	}

	sink.EnableHostNamePropagation()
}

// Implementation of methods in the MetricSink interface

func (s *DatadogSink) SetGauge(key []string, val float32) {
	s.SetGaugeWithLabels(key, val, nil)
}

func (s *DatadogSink) IncrCounter(key []string, val float32) {
	s.IncrCounterWithLabels(key, val, nil)
}

func (s *DatadogSink) EmitKey(key []string, val float32) {
	sink := s.getSink()
	if sink == nil {
		return
	}
	sink.EmitKey(key, val)
}

func (s *DatadogSink) AddSample(key []string, val float32) {
	s.AddSampleWithLabels(key, val, nil)
}

func (s *DatadogSink) initSinkIfNil() bool {
	if s.sink != nil {
		return false
	}

	now := time.Now()
	if s.attemptedToConnectAt != nil && now.Sub(*s.attemptedToConnectAt).Minutes() < 5 {
		return false
	}

	s.attemptedToConnectAt = &now

	sink, err := s.creator(s.addr, s.hostName)
	if err != nil {
		s.logger.Warn("failed to connect to datadog: " + err.Error())
		return false
	}

	s.sink = sink
	return true
}

func (s *DatadogSink) getSink() dogStatsdSink {
	s.lock.RLock()
	sink := s.sink
	s.lock.RUnlock()

	if sink != nil {
		return sink
	}

	s.lock.Lock()

	if !s.initSinkIfNil() {
		s.lock.Unlock()
		return s.sink
	}

	tags := s.tags
	s.tags = nil

	s.lock.Unlock()

	s.logger.Info("connected to datadog")

	s.sink.SetTags(tags)

	if s.propagateHostname {
		s.sink.EnableHostNamePropagation()
	}

	return s.sink
}

func (s *DatadogSink) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
	sink := s.getSink()
	if sink == nil {
		return
	}
	sink.SetGaugeWithLabels(key, val, labels)
}

func (s *DatadogSink) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
	sink := s.getSink()
	if sink == nil {
		return
	}
	sink.IncrCounterWithLabels(key, val, labels)
}

func (s *DatadogSink) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
	sink := s.getSink()
	if sink == nil {
		return
	}
	sink.IncrCounterWithLabels(key, val, labels)
}

func (s *DatadogSink) Shutdown() {
	sink := s.getSink()
	if sink == nil {
		return
	}
	sink.Shutdown()
}
