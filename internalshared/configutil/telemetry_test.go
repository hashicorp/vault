// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"testing"

	armonmetrics "github.com/armon/go-metrics"
	hclog "github.com/hashicorp/go-hclog"
	hcmetrics "github.com/hashicorp/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePrefixFilters(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputFilters            []string
		expectedErrStr          string
		expectedAllowedPrefixes []string
		expectedBlockedPrefixes []string
	}{
		{
			[]string{""},
			"Cannot have empty filter rule in prefix_filter",
			[]string(nil),
			[]string(nil),
		},
		{
			[]string{"vault.abc"},
			"Filter rule must begin with either '+' or '-': \"vault.abc\"",
			[]string(nil),
			[]string(nil),
		},
		{
			[]string{"+vault.abc", "-vault.bcd"},
			"",
			[]string{"vault.abc"},
			[]string{"vault.bcd"},
		},
	}
	t.Run("validate metric filter configs", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {

			allowedPrefixes, blockedPrefixes, err := parsePrefixFilter(tc.inputFilters)

			if err != nil {
				assert.EqualError(t, err, tc.expectedErrStr)
			} else {
				assert.Equal(t, "", tc.expectedErrStr)
				assert.Equal(t, tc.expectedAllowedPrefixes, allowedPrefixes)

				assert.Equal(t, tc.expectedBlockedPrefixes, blockedPrefixes)
			}
		}
	})
}

// TestNormalizeTelemetryAddresses ensures that any telemetry configuration that
// can be a URL, IP Address, or host:port address is conformant with RFC-5942 §4
// See: https://rfc-editor.org/rfc/rfc5952.html
func TestNormalizeTelemetryAddresses(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		given    *Telemetry
		expected *Telemetry
	}{
		"ipv6-conformance": {
			given: &Telemetry{
				// RFC-5952 4.1 leading zeroes
				CirconusAPIURL: "https://[2001:0db8::0001]:443",
				// RFC-5952 4.2.3 longest run of 0 bits shortened
				CirconusCheckSubmissionURL: "https://[2001:0:0:1:0:0:0:1]:443",
				// RFC-5952 4.2.3 equal runs of 0 bits shortened
				DogStatsDAddr: "https://[2001:db8:0:0:1:0:0:1]:443",
				// 	RFC-5952 4.3 downcase hex letters
				StatsdAddr:   "https://[2001:DB8:AC3:FE4::1]:443",
				StatsiteAddr: "https://[2001:DB8:AC3:FE4::1]:443",
			},
			expected: &Telemetry{
				CirconusAPIURL:             "https://[2001:db8::1]:443",
				CirconusCheckSubmissionURL: "https://[2001:0:0:1::1]:443",
				DogStatsDAddr:              "https://[2001:db8::1:0:0:1]:443",
				StatsdAddr:                 "https://[2001:db8:ac3:fe4::1]:443",
				StatsiteAddr:               "https://[2001:db8:ac3:fe4::1]:443",
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			normalizeTelemetryAddresses(tc.given)
			require.EqualValues(t, tc.expected, tc.given)
		})
	}
}

// TestConvertLabels proves that convertLabels produces a hcmetrics.Label slice
// with identical Name/Value fields to the input armonmetrics.Label slice, and
// that an empty input yields an empty (not nil) output.
func TestConvertLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in       []armonmetrics.Label
		expected []hcmetrics.Label
	}{
		"nil input": {
			in:       nil,
			expected: []hcmetrics.Label{},
		},
		"single label": {
			in:       []armonmetrics.Label{{Name: "region", Value: "us-east-1"}},
			expected: []hcmetrics.Label{{Name: "region", Value: "us-east-1"}},
		},
		"multiple labels": {
			in: []armonmetrics.Label{
				{Name: "region", Value: "us-east-1"},
				{Name: "cluster", Value: "primary"},
			},
			expected: []hcmetrics.Label{
				{Name: "region", Value: "us-east-1"},
				{Name: "cluster", Value: "primary"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, convertLabels(tc.in))
		})
	}
}

// recordingSink is a hcmetrics.MetricSink that records every call made to it so
// tests can assert that hcSinkAdapter forwards each method correctly.
type recordingSink struct {
	gauges          [][]string
	gaugesWithLabel [][]string
	gaugeLabels     [][]hcmetrics.Label
	keys            [][]string
	counters        [][]string
	countersLabeled [][]string
	counterLabels   [][]hcmetrics.Label
	samples         [][]string
	samplesLabeled  [][]string
	sampleLabels    [][]hcmetrics.Label
}

func (r *recordingSink) SetGauge(key []string, _ float32) { r.gauges = append(r.gauges, key) }
func (r *recordingSink) SetGaugeWithLabels(key []string, _ float32, labels []hcmetrics.Label) {
	r.gaugesWithLabel = append(r.gaugesWithLabel, key)
	r.gaugeLabels = append(r.gaugeLabels, labels)
}
func (r *recordingSink) EmitKey(key []string, _ float32)     { r.keys = append(r.keys, key) }
func (r *recordingSink) IncrCounter(key []string, _ float32) { r.counters = append(r.counters, key) }
func (r *recordingSink) IncrCounterWithLabels(key []string, _ float32, labels []hcmetrics.Label) {
	r.countersLabeled = append(r.countersLabeled, key)
	r.counterLabels = append(r.counterLabels, labels)
}
func (r *recordingSink) AddSample(key []string, _ float32) { r.samples = append(r.samples, key) }
func (r *recordingSink) AddSampleWithLabels(key []string, _ float32, labels []hcmetrics.Label) {
	r.samplesLabeled = append(r.samplesLabeled, key)
	r.sampleLabels = append(r.sampleLabels, labels)
}

// TestHcSinkAdapter proves that every method on hcSinkAdapter delegates to the
// underlying hcmetrics.MetricSink and that label conversion is applied on the
// WithLabels variants.
func TestHcSinkAdapter(t *testing.T) {
	t.Parallel()

	key := []string{"vault", "core", "requests"}
	armonLabels := []armonmetrics.Label{
		{Name: "namespace", Value: "root"},
		{Name: "mount", Value: "secret"},
	}
	wantHCLabels := []hcmetrics.Label{
		{Name: "namespace", Value: "root"},
		{Name: "mount", Value: "secret"},
	}

	tests := []struct {
		name       string
		call       func(*hcSinkAdapter)
		gotKeys    func(*recordingSink) [][]string
		gotLabels  func(*recordingSink) [][]hcmetrics.Label
		wantLabels [][]hcmetrics.Label
	}{
		{
			name:      "SetGauge",
			call:      func(a *hcSinkAdapter) { a.SetGauge(key, 1.0) },
			gotKeys:   func(r *recordingSink) [][]string { return r.gauges },
			gotLabels: nil,
		},
		{
			name:       "SetGaugeWithLabels",
			call:       func(a *hcSinkAdapter) { a.SetGaugeWithLabels(key, 1.0, armonLabels) },
			gotKeys:    func(r *recordingSink) [][]string { return r.gaugesWithLabel },
			gotLabels:  func(r *recordingSink) [][]hcmetrics.Label { return r.gaugeLabels },
			wantLabels: [][]hcmetrics.Label{wantHCLabels},
		},
		{
			name:      "EmitKey",
			call:      func(a *hcSinkAdapter) { a.EmitKey(key, 1.0) },
			gotKeys:   func(r *recordingSink) [][]string { return r.keys },
			gotLabels: nil,
		},
		{
			name:      "IncrCounter",
			call:      func(a *hcSinkAdapter) { a.IncrCounter(key, 1.0) },
			gotKeys:   func(r *recordingSink) [][]string { return r.counters },
			gotLabels: nil,
		},
		{
			name:       "IncrCounterWithLabels",
			call:       func(a *hcSinkAdapter) { a.IncrCounterWithLabels(key, 1.0, armonLabels) },
			gotKeys:    func(r *recordingSink) [][]string { return r.countersLabeled },
			gotLabels:  func(r *recordingSink) [][]hcmetrics.Label { return r.counterLabels },
			wantLabels: [][]hcmetrics.Label{wantHCLabels},
		},
		{
			name:      "AddSample",
			call:      func(a *hcSinkAdapter) { a.AddSample(key, 1.0) },
			gotKeys:   func(r *recordingSink) [][]string { return r.samples },
			gotLabels: nil,
		},
		{
			name:       "AddSampleWithLabels",
			call:       func(a *hcSinkAdapter) { a.AddSampleWithLabels(key, 1.0, armonLabels) },
			gotKeys:    func(r *recordingSink) [][]string { return r.samplesLabeled },
			gotLabels:  func(r *recordingSink) [][]hcmetrics.Label { return r.sampleLabels },
			wantLabels: [][]hcmetrics.Label{wantHCLabels},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := &recordingSink{}
			tc.call(&hcSinkAdapter{sink: rec})
			require.Equal(t, [][]string{key}, tc.gotKeys(rec))
			if tc.gotLabels != nil {
				assert.Equal(t, tc.wantLabels, tc.gotLabels(rec))
			}
		})
	}
}

// TestSetupTelemetry_LoggerWired proves that SetupTelemetry accepts a Logger on
// the opts struct and does not error when statsd/statsite addresses are omitted.
// It also verifies that providing a non-nil logger with no sink addresses still
// initialises the in-memory sink successfully.
func TestSetupTelemetry_LoggerWired(t *testing.T) {
	tests := map[string]struct {
		logger hclog.Logger
	}{
		"nil logger": {
			logger: nil,
		},
		"named logger": {
			logger: hclog.NewNullLogger().Named("telemetry"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inmem, clusterSink, _, err := SetupTelemetry(&SetupTelemetryOpts{
				Config:      &Telemetry{},
				Ui:          noopUI{},
				ServiceName: "vault",
				DisplayName: "Vault",
				Logger:      tc.logger,
			})
			require.NoError(t, err)
			assert.NotNil(t, inmem)
			assert.NotNil(t, clusterSink)
		})
	}
}

// noopUI satisfies the cli.Ui interface with no-op implementations so
// SetupTelemetry can be called in tests without a real terminal.
type noopUI struct{}

func (noopUI) Ask(_ string) (string, error)       { return "", nil }
func (noopUI) AskSecret(_ string) (string, error) { return "", nil }
func (noopUI) Output(_ string)                    {}
func (noopUI) Info(_ string)                      {}
func (noopUI) Error(_ string)                     {}
func (noopUI) Warn(_ string)                      {}
