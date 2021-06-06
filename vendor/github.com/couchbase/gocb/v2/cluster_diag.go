package gocb

import (
	"encoding/json"
	"time"

	"github.com/couchbase/gocbcore/v9"

	"github.com/google/uuid"
)

// EndPointDiagnostics represents a single entry in a diagnostics report.
type EndPointDiagnostics struct {
	Type         ServiceType
	ID           string
	Local        string
	Remote       string
	LastActivity time.Time
	State        EndpointState
	Namespace    string
}

// DiagnosticsResult encapsulates the results of a Diagnostics operation.
type DiagnosticsResult struct {
	ID       string
	Services map[string][]EndPointDiagnostics
	sdk      string
	State    ClusterState
}

type jsonDiagnosticEntry struct {
	ID             string `json:"id,omitempty"`
	LastActivityUs uint64 `json:"last_activity_us,omitempty"`
	Remote         string `json:"remote,omitempty"`
	Local          string `json:"local,omitempty"`
	State          string `json:"state,omitempty"`
	Details        string `json:"details,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
}

type jsonDiagnosticReport struct {
	Version  int16                            `json:"version"`
	SDK      string                           `json:"sdk,omitempty"`
	ID       string                           `json:"id,omitempty"`
	Services map[string][]jsonDiagnosticEntry `json:"services"`
	State    string                           `json:"state"`
}

// MarshalJSON generates a JSON representation of this diagnostics report.
func (report *DiagnosticsResult) MarshalJSON() ([]byte, error) {
	jsonReport := jsonDiagnosticReport{
		Version:  2,
		SDK:      report.sdk,
		ID:       report.ID,
		Services: make(map[string][]jsonDiagnosticEntry),
		State:    clusterStateToString(report.State),
	}

	for _, serviceType := range report.Services {
		for _, service := range serviceType {
			serviceStr := serviceTypeToString(service.Type)
			stateStr := endpointStateToString(service.State)

			jsonReport.Services[serviceStr] = append(jsonReport.Services[serviceStr], jsonDiagnosticEntry{
				ID:             service.ID,
				LastActivityUs: uint64(time.Since(service.LastActivity).Nanoseconds()),
				Remote:         service.Remote,
				Local:          service.Local,
				State:          stateStr,
				Details:        "",
				Namespace:      service.Namespace,
			})
		}
	}

	return json.Marshal(&jsonReport)
}

// DiagnosticsOptions are the options that are available for use with the Diagnostics operation.
type DiagnosticsOptions struct {
	ReportID string
}

// Diagnostics returns information about the internal state of the SDK.
func (c *Cluster) Diagnostics(opts *DiagnosticsOptions) (*DiagnosticsResult, error) {
	if opts == nil {
		opts = &DiagnosticsOptions{}
	}

	if opts.ReportID == "" {
		opts.ReportID = uuid.New().String()
	}

	provider, err := c.getDiagnosticsProvider()
	if err != nil {
		return nil, err
	}

	agentReport, err := provider.Diagnostics(gocbcore.DiagnosticsOptions{})
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
