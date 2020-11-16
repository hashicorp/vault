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

package jaeger

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

// Environment variable names
const (
	// The service name.
	envServiceName = "JAEGER_SERVICE_NAME"
	// Whether the exporter is disabled or not. (default false).
	envDisabled = "JAEGER_DISABLED"
	// A comma separated list of name=value tracer-level tags, which get added to all reported spans.
	// The value can also refer to an environment variable using the format ${envVarName:defaultValue}.
	envTags = "JAEGER_TAGS"
	// The HTTP endpoint for sending spans directly to a collector,
	// i.e. http://jaeger-collector:14268/api/traces.
	envEndpoint = "JAEGER_ENDPOINT"
	// Username to send as part of "Basic" authentication to the collector endpoint.
	envUser = "JAEGER_USER"
	// Password to send as part of "Basic" authentication to the collector endpoint.
	envPassword = "JAEGER_PASSWORD"
)

// CollectorEndpointFromEnv return environment variable value of JAEGER_ENDPOINT
func CollectorEndpointFromEnv() string {
	return os.Getenv(envEndpoint)
}

// WithCollectorEndpointOptionFromEnv uses environment variables to set the username and password
// if basic auth is required.
func WithCollectorEndpointOptionFromEnv() CollectorEndpointOption {
	return func(o *CollectorEndpointOptions) {
		if e := os.Getenv(envUser); e != "" {
			o.username = e
		}
		if e := os.Getenv(envPassword); e != "" {
			o.password = os.Getenv(envPassword)
		}
	}
}

// WithDisabledFromEnv uses environment variables and overrides disabled field.
func WithDisabledFromEnv() Option {
	return func(o *options) {
		if e := os.Getenv(envDisabled); e != "" {
			if v, err := strconv.ParseBool(e); err == nil {
				o.Disabled = v
			}
		}
	}
}

// ProcessFromEnv parse environment variables into jaeger exporter's Process.
// It will return a nil tag slice if the environment variable JAEGER_TAGS is malformed.
func ProcessFromEnv() Process {
	var p Process
	if e := os.Getenv(envServiceName); e != "" {
		p.ServiceName = e
	}
	if e := os.Getenv(envTags); e != "" {
		tags, err := parseTags(e)
		if err != nil {
			global.Handle(err)
		} else {
			p.Tags = tags
		}
	}

	return p
}

// WithProcessFromEnv uses environment variables and overrides jaeger exporter's Process.
func WithProcessFromEnv() Option {
	return func(o *options) {
		p := ProcessFromEnv()
		if p.ServiceName != "" {
			o.Process.ServiceName = p.ServiceName
		}
		if len(p.Tags) != 0 {
			o.Process.Tags = p.Tags
		}
	}
}

var errTagValueNotFound = errors.New("missing tag value")
var errTagEnvironmentDefaultValueNotFound = errors.New("missing default value for tag environment value")

// parseTags parses the given string into a collection of Tags.
// Spec for this value:
// - comma separated list of key=value
// - value can be specified using the notation ${envVar:defaultValue}, where `envVar`
// is an environment variable and `defaultValue` is the value to use in case the env var is not set
func parseTags(sTags string) ([]label.KeyValue, error) {
	pairs := strings.Split(sTags, ",")
	tags := make([]label.KeyValue, len(pairs))
	for i, p := range pairs {
		field := strings.SplitN(p, "=", 2)
		if len(field) != 2 {
			return nil, errTagValueNotFound
		}
		k, v := strings.TrimSpace(field[0]), strings.TrimSpace(field[1])

		if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
			ed := strings.SplitN(v[2:len(v)-1], ":", 2)
			if len(ed) != 2 {
				return nil, errTagEnvironmentDefaultValueNotFound
			}
			e, d := ed[0], ed[1]
			v = os.Getenv(e)
			if v == "" && d != "" {
				v = d
			}
		}

		tags[i] = parseKeyValue(k, v)
	}

	return tags, nil
}

func parseKeyValue(k, v string) label.KeyValue {
	return label.KeyValue{
		Key:   label.Key(k),
		Value: parseValue(v),
	}
}

func parseValue(str string) label.Value {
	if v, err := strconv.ParseInt(str, 10, 64); err == nil {
		return label.Int64Value(v)
	}
	if v, err := strconv.ParseFloat(str, 64); err == nil {
		return label.Float64Value(v)
	}
	if v, err := strconv.ParseBool(str); err == nil {
		return label.BoolValue(v)
	}

	// Fallback
	return label.StringValue(str)
}
