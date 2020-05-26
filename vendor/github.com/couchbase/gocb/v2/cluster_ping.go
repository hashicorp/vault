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

	services := opts.ServiceTypes
	if services == nil {
		services = []ServiceType{
			ServiceTypeQuery,
			ServiceTypeSearch,
			ServiceTypeAnalytics,
		}
	}

	gocbcoreServices := make([]gocbcore.ServiceType, len(services))
	for i, svc := range services {
		if svc == ServiceTypeKeyValue {
			return nil, invalidArgumentsError{
				message: "keyvalue service is not a valid service type for cluster level ping",
			}
		}
		if svc == ServiceTypeViews {
			return nil, invalidArgumentsError{
				message: "view service is not a valid service type for cluster level ping",
			}
		}
		gocbcoreServices[i] = gocbcore.ServiceType(svc)
	}

	coreopts := gocbcore.PingOptions{
		ServiceTypes: gocbcoreServices,
	}
	now := time.Now()
	timeout := opts.Timeout
	if timeout == 0 {
		coreopts.N1QLDeadline = now.Add(c.timeoutsConfig.QueryTimeout)
		coreopts.CbasDeadline = now.Add(c.timeoutsConfig.AnalyticsTimeout)
		coreopts.FtsDeadline = now.Add(c.timeoutsConfig.SearchTimeout)
	} else {
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
