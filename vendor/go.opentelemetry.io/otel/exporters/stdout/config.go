// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stdout

import (
	"io"
	"os"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

var (
	defaultWriter              = os.Stdout
	defaultPrettyPrint         = false
	defaultTimestamps          = true
	defaultQuantiles           = []float64{0.5, 0.9, 0.99}
	defaultLabelEncoder        = label.DefaultEncoder()
	defaultDisableTraceExport  = false
	defaultDisableMetricExport = false
)

// Config contains options for the STDOUT exporter.
type Config struct {
	// Writer is the destination.  If not set, os.Stdout is used.
	Writer io.Writer

	// PrettyPrint will encode the output into readable JSON. Default is
	// false.
	PrettyPrint bool

	// Timestamps specifies if timestamps should be pritted. Default is
	// true.
	Timestamps bool

	// Quantiles are the desired aggregation quantiles for distribution
	// summaries, used when the configured aggregator supports
	// quantiles.
	//
	// Note: this exporter is meant as a demonstration; a real
	// exporter may wish to configure quantiles on a per-metric
	// basis.
	Quantiles []float64

	// LabelEncoder encodes the labels.
	LabelEncoder label.Encoder

	// DisableTraceExport prevents any export of trace telemetry.
	DisableTraceExport bool

	// DisableMetricExport prevents any export of metric telemetry.
	DisableMetricExport bool
}

// NewConfig creates a validated Config configured with options.
func NewConfig(options ...Option) (Config, error) {
	config := Config{
		Writer:              defaultWriter,
		PrettyPrint:         defaultPrettyPrint,
		Timestamps:          defaultTimestamps,
		Quantiles:           defaultQuantiles,
		LabelEncoder:        defaultLabelEncoder,
		DisableTraceExport:  defaultDisableTraceExport,
		DisableMetricExport: defaultDisableMetricExport,
	}
	for _, opt := range options {
		opt.Apply(&config)

	}
	for _, q := range config.Quantiles {
		if q < 0 || q > 1 {
			return config, aggregation.ErrInvalidQuantile
		}
	}
	return config, nil
}

// Option sets the value of an option for a Config.
type Option interface {
	// Apply option value to Config.
	Apply(*Config)
}

// WithWriter sets the export stream destination.
func WithWriter(w io.Writer) Option {
	return writerOption{w}
}

type writerOption struct {
	W io.Writer
}

func (o writerOption) Apply(config *Config) {
	config.Writer = o.W
}

// WithPrettyPrint sets the export stream format to use JSON.
func WithPrettyPrint() Option {
	return prettyPrintOption(true)
}

type prettyPrintOption bool

func (o prettyPrintOption) Apply(config *Config) {
	config.PrettyPrint = bool(o)
}

// WithoutTimestamps sets the export stream to not include timestamps.
func WithoutTimestamps() Option {
	return timestampsOption(false)
}

type timestampsOption bool

func (o timestampsOption) Apply(config *Config) {
	config.Timestamps = bool(o)
}

// WithQuantiles sets the quantile values to export.
func WithQuantiles(quantiles []float64) Option {
	return quantilesOption(quantiles)
}

type quantilesOption []float64

func (o quantilesOption) Apply(config *Config) {
	config.Quantiles = []float64(o)
}

// WithLabelEncoder sets the label encoder used in export.
func WithLabelEncoder(enc label.Encoder) Option {
	return labelEncoderOption{enc}
}

type labelEncoderOption struct {
	LabelEncoder label.Encoder
}

func (o labelEncoderOption) Apply(config *Config) {
	config.LabelEncoder = o.LabelEncoder
}

// WithoutTraceExport disables all trace exporting.
func WithoutTraceExport() Option {
	return disableTraceExportOption(true)
}

type disableTraceExportOption bool

func (o disableTraceExportOption) Apply(config *Config) {
	config.DisableTraceExport = bool(o)
}

// WithoutMetricExport disables all metric exporting.
func WithoutMetricExport() Option {
	return disableMetricExportOption(true)
}

type disableMetricExportOption bool

func (o disableMetricExportOption) Apply(config *Config) {
	config.DisableMetricExport = bool(o)
}
