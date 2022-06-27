package configutil

import (
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/datadog"
	"github.com/mitchellh/cli"
)

type DatadogSink struct {
	tags                 []string
	addr                 string
	hostName             string
	propagateHostname    bool
	sink                 *datadog.DogStatsdSink
	logger               cli.Ui
	attemptedToConnectAt *time.Time
	lock                 sync.Mutex
}

func NewDatadogSink(addr string, hostName string, logger cli.Ui) *DatadogSink {
	return &DatadogSink{
		addr:              addr,
		hostName:          hostName,
		propagateHostname: false,
		logger:            logger,
	}
}

func (s *DatadogSink) SetTags(tags []string) {
	s.lock.Lock()
	s.tags = tags
	s.lock.Unlock()

	sink := s.getSink()
	if sink == nil {
		return
	}
	sink.SetTags(tags)
}

func (s *DatadogSink) EnableHostNamePropagation() {
	s.lock.Lock()
	s.propagateHostname = true
	s.lock.Unlock()

	sink := s.getSink()
	if sink == nil {
		return
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

func (s *DatadogSink) getSink() *datadog.DogStatsdSink {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.sink != nil {
		return s.sink
	}

	now := time.Now()
	if s.attemptedToConnectAt != nil && now.Sub(*s.attemptedToConnectAt).Minutes() < 5 {
		return nil
	}

	s.attemptedToConnectAt = &now

	sink, err := datadog.NewDogStatsdSink(s.addr, s.hostName)
	if err != nil {
		s.logger.Warn("failed to connect to datadog: " + err.Error())
		return nil
	}

	s.logger.Info("connected to datadog")

	sink.SetTags(s.tags)
	if s.propagateHostname {
		sink.EnableHostNamePropagation()
	}

	s.sink = sink
	return sink
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
