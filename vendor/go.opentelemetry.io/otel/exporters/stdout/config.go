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

package stdout // import "go.opentelemetry.io/otel/exporters/stdout"

import (
	"io"
	"os"

	"go.opentelemetry.io/otel/attribute"
)

var (
	defaultWriter              = os.Stdout
	defaultPrettyPrint         = false
	defaultTimestamps          = true
	defaultLabelEncoder        = attribute.DefaultEncoder()
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

	// LabelEncoder encodes the labels.
	LabelEncoder attribute.Encoder

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
		LabelEncoder:        defaultLabelEncoder,
		DisableTraceExport:  defaultDisableTraceExport,
		DisableMetricExport: defaultDisableMetricExport,
	}
	for _, opt := range options {
		opt.Apply(&config)

	}
	return config, nil
}

// Option sets the value of an option for a Config.
type Option interface {
	// Apply option value to Config.
	Apply(*Config)

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
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

func (writerOption) private() {}

// WithPrettyPrint sets the export stream format to use JSON.
func WithPrettyPrint() Option {
	return prettyPrintOption(true)
}

type prettyPrintOption bool

func (o prettyPrintOption) Apply(config *Config) {
	config.PrettyPrint = bool(o)
}

func (prettyPrintOption) private() {}

// WithoutTimestamps sets the export stream to not include timestamps.
func WithoutTimestamps() Option {
	return timestampsOption(false)
}

type timestampsOption bool

func (o timestampsOption) Apply(config *Config) {
	config.Timestamps = bool(o)
}

func (timestampsOption) private() {}

// WithLabelEncoder sets the label encoder used in export.
func WithLabelEncoder(enc attribute.Encoder) Option {
	return labelEncoderOption{enc}
}

type labelEncoderOption struct {
	LabelEncoder attribute.Encoder
}

func (o labelEncoderOption) Apply(config *Config) {
	config.LabelEncoder = o.LabelEncoder
}

func (labelEncoderOption) private() {}

// WithoutTraceExport disables all trace exporting.
func WithoutTraceExport() Option {
	return disableTraceExportOption(true)
}

type disableTraceExportOption bool

func (o disableTraceExportOption) Apply(config *Config) {
	config.DisableTraceExport = bool(o)
}

func (disableTraceExportOption) private() {}

// WithoutMetricExport disables all metric exporting.
func WithoutMetricExport() Option {
	return disableMetricExportOption(true)
}

type disableMetricExportOption bool

func (o disableMetricExportOption) Apply(config *Config) {
	config.DisableMetricExport = bool(o)
}

func (disableMetricExportOption) private() {}
