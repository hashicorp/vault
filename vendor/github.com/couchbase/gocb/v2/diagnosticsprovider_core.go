package gocb

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type diagnosticsProviderCoreProvider interface {
	Diagnostics(opts gocbcore.DiagnosticsOptions) (*gocbcore.DiagnosticInfo, error)
	Ping(ctx context.Context, opts gocbcore.PingOptions) (*gocbcore.PingResult, error)
}

type diagnosticsProviderCore struct {
	provider diagnosticsProviderCoreProvider

	tracer   RequestTracer
	meter    *meterWrapper
	timeouts TimeoutsConfig
}

func (d *diagnosticsProviderCore) Diagnostics(opts *DiagnosticsOptions) (*DiagnosticsResult, error) {
	if opts == nil {
		opts = &DiagnosticsOptions{}
	}

	if opts.ReportID == "" {
		opts.ReportID = uuid.New().String()
	}

	agentReport, err := d.provider.Diagnostics(gocbcore.DiagnosticsOptions{})
	if err != nil {
		return nil, err
	}

	report := &DiagnosticsResult{
		ID:       opts.ReportID,
		Services: make(map[string][]EndPointDiagnostics),
		sdk:      Identifier(),
		State:    ClusterState(agentReport.State),
	}

	report.Services["kv"] = make([]EndPointDiagnostics, 0)

	for _, conn := range agentReport.MemdConns {
		state := EndpointState(conn.State)

		report.Services["kv"] = append(report.Services["kv"], EndPointDiagnostics{
			Type:         ServiceTypeKeyValue,
			State:        state,
			Local:        conn.LocalAddr,
			Remote:       conn.RemoteAddr,
			LastActivity: conn.LastActivity,
			Namespace:    conn.Scope,
			ID:           conn.ID,
		})
	}

	return report, nil
}

func (d *diagnosticsProviderCore) Ping(opts *PingOptions) (*PingResult, error) {
	startTime := time.Now()
	defer d.meter.ValueRecord(meterValueServiceKV, "ping", startTime)

	span := createSpan(d.tracer, opts.ParentSpan, "ping", "kv")
	defer span.End()

	services := opts.ServiceTypes

	gocbcoreServices := make([]gocbcore.ServiceType, len(services))
	for i, svc := range services {
		gocbcoreServices[i] = gocbcore.ServiceType(svc)
	}

	coreopts := gocbcore.PingOptions{
		ServiceTypes: gocbcoreServices,
		TraceContext: span.Context(),
	}

	now := time.Now()
	timeout := opts.Timeout
	if timeout == 0 {
		coreopts.KVDeadline = now.Add(d.timeouts.KVTimeout)
		coreopts.CapiDeadline = now.Add(d.timeouts.ViewTimeout)
		coreopts.N1QLDeadline = now.Add(d.timeouts.QueryTimeout)
		coreopts.CbasDeadline = now.Add(d.timeouts.AnalyticsTimeout)
		coreopts.FtsDeadline = now.Add(d.timeouts.SearchTimeout)
		coreopts.MgmtDeadline = now.Add(d.timeouts.ManagementTimeout)
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

	result, err := d.provider.Ping(opts.Context, coreopts)
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
