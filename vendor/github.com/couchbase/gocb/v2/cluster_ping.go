package gocb

import (
	"time"

	"github.com/couchbase/gocbcore/v9"
	"github.com/google/uuid"
)

// Ping will ping a list of services and verify they are active and
// responding in an acceptable period of time.
func (c *Cluster) Ping(opts *PingOptions) (*PingResult, error) {
	if opts == nil {
		opts = &PingOptions{}
	}

	provider, err := c.getDiagnosticsProvider()
	if err != nil {
		return nil, err
	}

	return ping(provider, opts, c.timeoutsConfig)
}

func ping(provider diagnosticsProvider, opts *PingOptions, timeouts TimeoutsConfig) (*PingResult, error) {
	services := opts.ServiceTypes

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
		coreopts.KVDeadline = now.Add(timeouts.KVTimeout)
		coreopts.CapiDeadline = now.Add(timeouts.ViewTimeout)
		coreopts.N1QLDeadline = now.Add(timeouts.QueryTimeout)
		coreopts.CbasDeadline = now.Add(timeouts.AnalyticsTimeout)
		coreopts.FtsDeadline = now.Add(timeouts.SearchTimeout)
		coreopts.MgmtDeadline = now.Add(timeouts.ManagementTimeout)
	} else {
		coreopts.KVDeadline = now.Add(timeout)
		coreopts.CapiDeadline = now.Add(timeout)
		coreopts.N1QLDeadline = now.Add(timeout)
		coreopts.CbasDeadline = now.Add(timeout)
		coreopts.FtsDeadline = now.Add(timeout)
		coreopts.MgmtDeadline = now.Add(timeout)
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
