package gocb

import (
	"encoding/json"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
	"github.com/google/uuid"
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

	cli := b.getCachedClient()
	provider, err := cli.getDiagnosticsProvider()
	if err != nil {
		return nil, err
	}

	services := opts.ServiceTypes
	if services == nil {
		services = []ServiceType{
			ServiceTypeKeyValue,
			ServiceTypeViews,
			ServiceTypeQuery,
			ServiceTypeSearch,
			ServiceTypeAnalytics,
		}
	}

	gocbcoreServices := make([]gocbcore.ServiceType, len(services))
	for i, svc := range services {
		gocbcoreServices[i] = gocbcore.ServiceType(svc)
	}

	coreopts := gocbcore.PingOptions{
		ServiceTypes: gocbcoreServices,
	}
	now := time.Now()
	timeout := opts.Timeout
	if timeout == 0 {
		coreopts.KVDeadline = now.Add(b.timeoutsConfig.KVTimeout)
		coreopts.CapiDeadline = now.Add(b.timeoutsConfig.ViewTimeout)
		coreopts.N1QLDeadline = now.Add(b.timeoutsConfig.QueryTimeout)
		coreopts.CbasDeadline = now.Add(b.timeoutsConfig.AnalyticsTimeout)
		coreopts.FtsDeadline = now.Add(b.timeoutsConfig.SearchTimeout)
	} else {
		coreopts.KVDeadline = now.Add(timeout)
		coreopts.CapiDeadline = now.Add(timeout)
		coreopts.N1QLDeadline = now.Add(timeout)
		coreopts.CbasDeadline = now.Add(timeout)
		coreopts.FtsDeadline = now.Add(timeout)
	}

	id := opts.ReportID
	if id == "" {
		id = uuid.New().String()
	}

	result, err := provider.Ping(coreopts)
	if err != nil {
		return nil, err
	}

	reportSvcs := make(map[ServiceType][]EndpointPingReport)
	for svcType, svc := range result.Services {
		st := ServiceType(svcType)

		svcs := make([]EndpointPingReport, len(svc))
		for i, rep := range svc {
			var errStr string
			if rep.Error != nil {
				errStr = rep.Error.Error()
			}
			svcs[i] = EndpointPingReport{
				ID:        rep.ID,
				Remote:    rep.Endpoint,
				State:     PingState(rep.State),
				Error:     errStr,
				Namespace: rep.Scope,
				Latency:   rep.Latency,
			}
		}

		reportSvcs[st] = svcs
	}

	return &PingResult{
		ID:       id,
		sdk:      Identifier() + " " + "gocbcore/" + gocbcore.Version(),
		Services: reportSvcs,
	}, nil
}
