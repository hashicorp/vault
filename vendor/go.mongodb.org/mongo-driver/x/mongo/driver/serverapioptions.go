// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

// TestServerAPIVersion is the most recent, stable variant of options.ServerAPIVersion.
// Only to be used in testing.
const TestServerAPIVersion = "1"

// ServerAPIOptions represents options used to configure the API version sent to the server
// when running commands.
type ServerAPIOptions struct {
	ServerAPIVersion  string
	Strict            *bool
	DeprecationErrors *bool
}

// NewServerAPIOptions creates a new ServerAPIOptions configured with the provided serverAPIVersion.
func NewServerAPIOptions(serverAPIVersion string) *ServerAPIOptions {
	return &ServerAPIOptions{ServerAPIVersion: serverAPIVersion}
}

// SetStrict specifies whether the server should return errors for features that are not part of the API version.
func (s *ServerAPIOptions) SetStrict(strict bool) *ServerAPIOptions {
	s.Strict = &strict
	return s
}

// SetDeprecationErrors specifies whether the server should return errors for deprecated features.
func (s *ServerAPIOptions) SetDeprecationErrors(deprecationErrors bool) *ServerAPIOptions {
	s.DeprecationErrors = &deprecationErrors
	return s
}
