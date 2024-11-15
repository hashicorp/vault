// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instance

import (
	"context"
	"fmt"
	"regexp"

	"cloud.google.com/go/cloudsqlconn/errtype"
)

var (
	// Instance connection name is the format <PROJECT>:<REGION>:<INSTANCE>
	// Additionally, we have to support legacy "domain-scoped" projects
	// (e.g. "google.com:PROJECT")
	connNameRegex = regexp.MustCompile("([^:]+(:[^:]+)?):([^:]+):([^:]+)")
)

// ConnName represents the "instance connection name", in the format
// "project:region:name".
type ConnName struct {
	project    string
	region     string
	name       string
	domainName string
}

func (c *ConnName) String() string {
	if c.domainName != "" {
		return fmt.Sprintf("%s -> %s:%s:%s", c.domainName, c.project, c.region, c.name)
	}
	return fmt.Sprintf("%s:%s:%s", c.project, c.region, c.name)
}

// Project returns the project within which the Cloud SQL instance runs.
func (c *ConnName) Project() string {
	return c.project
}

// Region returns the region where the Cloud SQL instance runs.
func (c *ConnName) Region() string {
	return c.region
}

// Name returns the Cloud SQL instance name
func (c *ConnName) Name() string {
	return c.name
}

// DomainName returns the domain name for this instance
func (c *ConnName) DomainName() string {
	return c.domainName
}

// HasDomainName returns the Cloud SQL domain name
func (c *ConnName) HasDomainName() bool {
	return c.domainName != ""
}

// ParseConnName initializes a new ConnName struct.
func ParseConnName(cn string) (ConnName, error) {
	return ParseConnNameWithDomainName(cn, "")
}

// ParseConnNameWithDomainName initializes a new ConnName struct,
// also setting the domain name.
func ParseConnNameWithDomainName(cn string, dn string) (ConnName, error) {
	b := []byte(cn)
	m := connNameRegex.FindSubmatch(b)
	if m == nil {
		err := errtype.NewConfigError(
			"invalid instance connection name, expected PROJECT:REGION:INSTANCE",
			cn,
		)
		return ConnName{}, err
	}

	c := ConnName{
		project:    string(m[1]),
		region:     string(m[3]),
		name:       string(m[4]),
		domainName: dn,
	}
	return c, nil
}

// ConnectionNameResolver resolves the connection name string into a valid
// instance name. This allows an application to replace the default
// resolver with a custom implementation.
type ConnectionNameResolver interface {
	// Resolve accepts a name, and returns a ConnName with the instance
	// connection string for the name. If the name cannot be resolved, returns
	// an error.
	Resolve(ctx context.Context, name string) (ConnName, error)
}
