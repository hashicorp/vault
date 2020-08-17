package gocb

import (
	"encoding/json"
	"time"
)

// EndpointPingReport represents a single entry in a ping report.
type EndpointPingReport struct {
	ID        string
	Local     string
	Remote    string
	State     PingState
	Error     string
	Namespace string
	Latency   time.Duration
}

// PingResult encapsulates the details from a executed ping operation.
type PingResult struct {
	ID       string
	Services map[ServiceType][]EndpointPingReport

	sdk string
}

type jsonEndpointPingReport struct {
	ID        string `json:"id,omitempty"`
	Local     string `json:"local,omitempty"`
	Remote    string `json:"remote,omitempty"`
	State     string `json:"state,omitempty"`
	Error     string `json:"error,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	LatencyUs uint64 `json:"latency_us"`
}

type jsonPingReport struct {
	Version  uint16                              `json:"version"`
	SDK      string                              `json:"sdk,omitempty"`
	ID       string                              `json:"id,omitempty"`
	Services map[string][]jsonEndpointPingReport `json:"services,omitempty"`
}

// MarshalJSON generates a JSON representation of this ping report.
func (report *PingResult) MarshalJSON() ([]byte, error) {
	jsonReport := jsonPingReport{
		Version:  2,
		SDK:      report.sdk,
		ID:       report.ID,
		Services: make(map[string][]jsonEndpointPingReport),
	}

	for serviceType, serviceInfo := range report.Services {
		serviceStr := serviceTypeToString(serviceType)
		if _, ok := jsonReport.Services[serviceStr]; !ok {
			jsonReport.Services[serviceStr] = make([]jsonEndpointPingReport, 0)
		}

		for _, service := range serviceInfo {
			jsonReport.Services[serviceStr] = append(jsonReport.Services[serviceStr], jsonEndpointPingReport{
				ID:        service.ID,
				Local:     service.Local,
				Remote:    service.Remote,
				State:     pingStateToString(service.State),
				Error:     service.Error,
				Namespace: service.Namespace,
				LatencyUs: uint64(service.Latency / time.Nanosecond),
			})
		}
	}

	return json.Marshal(&jsonReport)
}

// PingOptions are the options available to the Ping operation.
type PingOptions struct {
	ServiceTypes []ServiceType
	ReportID     string
	Timeout      time.Duration
}

// Ping will ping a list of services and verify they are active and
// responding in an acceptable period of time.
func (b *Bucket) Ping(opts *PingOptions) (*PingResult, error) {
	if opts == nil {
		opts = &PingOptions{}
	}

	provider, err := b.connectionManager.getDiagnosticsProvider(b.bucketName)
	if err != nil {
		return nil, err
	}

	return ping(provider, opts, b.timeoutsConfig)
}
