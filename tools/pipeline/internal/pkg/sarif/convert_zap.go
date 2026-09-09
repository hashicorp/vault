// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	gosarif "github.com/owenrumney/go-sarif/v3/pkg/report/v22/sarif"
)

// ConvertZapReq is the request to convert SARIF to ZAP format
type ConvertZapReq struct {
	SarifPath  string // Path to SARIF file
	OutputPath string // Output file path (empty for stdout)
}

// ConvertZapRes is the response for ZAP format conversion
type ConvertZapRes struct {
	OutputPath string     // Output file path
	Report     *ZapReport // Converted report
}

// ZapReport represents the top-level ZAP format output
type ZapReport struct {
	ProgramName string     `json:"@programName"`
	Version     string     `json:"@version"`
	Generated   string     `json:"@generated"`
	Sites       []*ZapSite `json:"site"`
}

// ZapSite represents a site in the ZAP format
type ZapSite struct {
	Name   string      `json:"@name"`
	Host   string      `json:"@host"`
	Port   string      `json:"@port"`
	SSL    string      `json:"@ssl"`
	Alerts []*ZapAlert `json:"alerts"`
}

// ZapAlert represents an alert/finding in the ZAP format
type ZapAlert struct {
	PluginID   string         `json:"pluginid"`
	AlertRef   string         `json:"alertRef"`
	Alert      string         `json:"alert"`
	Name       string         `json:"name"`
	RiskCode   string         `json:"riskcode"`
	Confidence string         `json:"confidence"`
	RiskDesc   string         `json:"riskdesc"`
	Desc       string         `json:"desc"`
	Instances  []*ZapInstance `json:"instances"`
	Count      string         `json:"count"`
	Solution   string         `json:"solution"`
	OtherInfo  string         `json:"otherinfo"`
	Reference  string         `json:"reference"`
	CWEID      string         `json:"cweid"`
	WASCID     string         `json:"wascid"`
	SourceID   string         `json:"sourceid"`
}

// ZapInstance represents an instance of an alert
type ZapInstance struct {
	URI       string `json:"uri"`
	Method    string `json:"method"`
	Param     string `json:"param"`
	Attack    string `json:"attack"`
	Evidence  string `json:"evidence"`
	OtherInfo string `json:"otherinfo"`
}

// siteInfo holds parsed site information
type siteInfo struct {
	name string
	host string
	port string
	ssl  string
}

// riskInfo holds risk code and label
type riskInfo struct {
	code  string
	label string
}

// zapSiteBuilder is a helper for building ZAP sites
type zapSiteBuilder struct {
	info   siteInfo
	alerts map[string]*zapAlertBuilder
}

// zapAlertBuilder is a helper for building ZAP alerts
type zapAlertBuilder struct {
	PluginID   string
	AlertRef   string
	Alert      string
	Name       string
	RiskCode   string
	Confidence string
	RiskDesc   string
	Desc       string
	Solution   string
	OtherInfo  string
	Reference  string
	CWEID      string
	WASCID     string
	SourceID   string
	Instances  []*ZapInstance
}

// Run executes the ZAP format conversion
func (r *ConvertZapReq) Run(ctx context.Context) (*ConvertZapRes, error) {
	// Parse SARIF file using go-sarif library
	report, err := gosarif.Open(r.SarifPath)
	if err != nil {
		return nil, fmt.Errorf("opening SARIF file: %w", err)
	}

	return &ConvertZapRes{
		OutputPath: r.OutputPath,
		Report:     sarifToZap(report),
	}, nil
}

// ToJSON marshals the response to JSON
func (r *ConvertZapRes) ToJSON() ([]byte, error) {
	b, err := json.MarshalIndent(r.Report, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling ZAP report to JSON: %w", err)
	}
	// Add trailing newline to match our golden examples from the node.js
	// example.
	b = append(b, '\n')
	return b, nil
}

// ToTable marshals the response to a text table
func (r *ConvertZapRes) ToTable(err error) (table.Writer, error) {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	defer t.SuppressEmptyColumns()
	defer t.SuppressTrailingSpaces()

	if err != nil {
		t.AppendHeader(table.Row{"error"})
		t.AppendRow(table.Row{err.Error()})
		return t, err
	}

	if r == nil || r.Report == nil || len(r.Report.Sites) == 0 {
		t.AppendHeader(table.Row{"message"})
		t.AppendRow(table.Row{"No sites found in SARIF report."})
		return t, nil
	}

	// Summary table
	t.AppendHeader(table.Row{"SITE", "HOST", "PORT", "SSL", "ALERTS", "INSTANCES"})
	for _, site := range r.Report.Sites {
		totalInstances := 0
		for _, alert := range site.Alerts {
			totalInstances += len(alert.Instances)
		}
		t.AppendRow(table.Row{
			site.Name,
			site.Host,
			site.Port,
			site.SSL,
			len(site.Alerts),
			totalInstances,
		})
	}

	return t, nil
}

// sarifToZap converts SARIF report to ZAP format
func sarifToZap(report *gosarif.Report) *ZapReport {
	if report == nil {
		return nil
	}

	// Group results by site, preserving insertion order
	siteMap := make(map[string]*zapSiteBuilder)
	siteOrder := []string{}

	for _, run := range report.Runs {
		if run.Results == nil {
			continue
		}

		for _, result := range run.Results {
			// Get target URI from webRequest or physicalLocation
			target := sarifResultURI(result)
			if target == "" {
				target = "http://unknown/"
			}

			// Parse site information
			site := parseSite(target)
			siteKey := fmt.Sprintf("%s|%s|%s|%s", site.name, site.host, site.port, site.ssl)

			// Get or create site builder
			if _, exists := siteMap[siteKey]; !exists {
				siteMap[siteKey] = &zapSiteBuilder{
					info:   site,
					alerts: make(map[string]*zapAlertBuilder),
				}
				siteOrder = append(siteOrder, siteKey)
			}
			siteBuilder := siteMap[siteKey]

			// Get alert information
			pluginID := "0"
			if result.RuleID != nil {
				pluginID = *result.RuleID
			}

			alertText := fmt.Sprintf("Rule %s", pluginID)
			if result.Message.Text != nil {
				alertText = *result.Message.Text
			}

			// Get risk information
			level := result.Level
			if level == "" {
				level = "note"
			}
			risk := getRiskInfo(level)

			alertKey := fmt.Sprintf("%s::%s", pluginID, alertText)

			// Get or create alert builder
			if _, exists := siteBuilder.alerts[alertKey]; !exists {
				siteBuilder.alerts[alertKey] = &zapAlertBuilder{
					PluginID:   pluginID,
					AlertRef:   pluginID,
					Alert:      alertText,
					Name:       alertText,
					RiskCode:   risk.code,
					Confidence: "2",
					RiskDesc:   fmt.Sprintf("%s (Medium)", risk.label),
					Desc:       fmt.Sprintf("<p>%s</p>", alertText),
					Solution:   "",
					OtherInfo:  "",
					Reference:  "",
					CWEID:      "-1",
					WASCID:     "-1",
					SourceID:   "0",
					Instances:  []*ZapInstance{},
				}
			}
			alertBuilder := siteBuilder.alerts[alertKey]

			// Process locations
			locations := result.Locations
			if len(locations) == 0 {
				locations = []*gosarif.Location{nil}
			}

			for _, loc := range locations {
				instance := sarifResultLocationToZapInstance(result, loc, target)
				alertBuilder.Instances = append(alertBuilder.Instances, &instance)
			}
		}
	}

	// Build final report
	zapReport := &ZapReport{
		ProgramName: "ZAP",
		Version:     "2.15.0",
		Generated:   time.Now().UTC().Format(time.RFC1123),
		Sites:       []*ZapSite{},
	}

	// Iterate sites in order of first appearance
	for _, siteKey := range siteOrder {
		siteBuilder := siteMap[siteKey]
		site := &ZapSite{
			Name:   siteBuilder.info.name,
			Host:   siteBuilder.info.host,
			Port:   siteBuilder.info.port,
			SSL:    siteBuilder.info.ssl,
			Alerts: []*ZapAlert{},
		}

		// Collect and sort alerts by pluginID to ensure consistent ordering
		var alertKeys []string
		for key := range siteBuilder.alerts {
			alertKeys = append(alertKeys, key)
		}
		sort.Strings(alertKeys)

		for _, key := range alertKeys {
			alertBuilder := siteBuilder.alerts[key]
			alert := &ZapAlert{
				PluginID:   alertBuilder.PluginID,
				AlertRef:   alertBuilder.AlertRef,
				Alert:      alertBuilder.Alert,
				Name:       alertBuilder.Name,
				RiskCode:   alertBuilder.RiskCode,
				Confidence: alertBuilder.Confidence,
				RiskDesc:   alertBuilder.RiskDesc,
				Desc:       alertBuilder.Desc,
				Instances:  alertBuilder.Instances,
				Count:      fmt.Sprintf("%d", len(alertBuilder.Instances)),
				Solution:   alertBuilder.Solution,
				OtherInfo:  alertBuilder.OtherInfo,
				Reference:  alertBuilder.Reference,
				CWEID:      alertBuilder.CWEID,
				WASCID:     alertBuilder.WASCID,
				SourceID:   alertBuilder.SourceID,
			}
			site.Alerts = append(site.Alerts, alert)
		}

		zapReport.Sites = append(zapReport.Sites, site)
	}

	return zapReport
}

// sarifResultURI extracts the target URI from a result
func sarifResultURI(result *gosarif.Result) string {
	// Try webRequest.target first (SARIF 2.1.0 §3.14.21)
	if result.WebRequest != nil && result.WebRequest.Target != nil {
		return *result.WebRequest.Target
	}

	// Fallback to first location's URI
	if len(result.Locations) > 0 {
		loc := result.Locations[0]
		if loc.PhysicalLocation != nil &&
			loc.PhysicalLocation.ArtifactLocation != nil &&
			loc.PhysicalLocation.ArtifactLocation.URI != nil {
			return *loc.PhysicalLocation.ArtifactLocation.URI
		}
	}

	return ""
}

// sarifResultLocationToZapInstance creates a ZapInstance from a result and location
func sarifResultLocationToZapInstance(result *gosarif.Result, loc *gosarif.Location, defaultURI string) ZapInstance {
	instance := ZapInstance{
		URI:       defaultURI,
		Method:    "",
		Param:     "",
		Attack:    "",
		Evidence:  "",
		OtherInfo: "",
	}

	// Get URI from location or use default
	if loc != nil {
		if loc.PhysicalLocation != nil &&
			loc.PhysicalLocation.ArtifactLocation != nil &&
			loc.PhysicalLocation.ArtifactLocation.URI != nil {
			instance.URI = *loc.PhysicalLocation.ArtifactLocation.URI
		}

		// Get attack from properties
		if loc.Properties != nil && loc.Properties.Properties != nil {
			if attackVal, ok := loc.Properties.Properties["attack"]; ok {
				if attackStr, ok := attackVal.(string); ok {
					instance.Attack = attackStr
				}
			}
		}

		// Get evidence from snippet
		if loc.PhysicalLocation != nil &&
			loc.PhysicalLocation.Region != nil &&
			loc.PhysicalLocation.Region.Snippet != nil &&
			loc.PhysicalLocation.Region.Snippet.Text != nil {
			instance.Evidence = *loc.PhysicalLocation.Region.Snippet.Text
		}
	}

	// Get method from webRequest
	if result != nil && result.WebRequest != nil && result.WebRequest.Method != nil {
		instance.Method = *result.WebRequest.Method
	}

	return instance
}

// parseSite parses a URI string and extracts site information
func parseSite(uriString string) siteInfo {
	u, err := url.Parse(uriString)
	if err != nil {
		return siteInfo{
			name: "unknown",
			host: "unknown",
			port: "80",
			ssl:  "false",
		}
	}

	host := u.Hostname()
	if host == "" {
		host = "unknown"
	}

	port := u.Port()
	ssl := "false"

	if u.Scheme == "https" {
		ssl = "true"
		if port == "" {
			port = "443"
		}
	} else {
		if port == "" {
			port = "80"
		}
	}

	name := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if name == "://" {
		name = "unknown"
	}

	return siteInfo{
		name: name,
		host: host,
		port: port,
		ssl:  ssl,
	}
}

// getRiskInfo maps SARIF level to Concert risk code and label
func getRiskInfo(level string) riskInfo {
	levelLower := strings.ToLower(level)
	switch levelLower {
	case "note":
		return riskInfo{code: "0", label: "Informational"}
	case "warning":
		return riskInfo{code: "1", label: "Low"}
	case "error":
		return riskInfo{code: "2", label: "Medium"}
	case "critical":
		return riskInfo{code: "3", label: "High"}
	default:
		return riskInfo{code: "0", label: "Informational"}
	}
}
