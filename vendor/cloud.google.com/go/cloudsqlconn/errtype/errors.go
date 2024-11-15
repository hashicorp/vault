// Copyright 2021 Google LLC
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

// Package errtype provides a number of concrete types which are used by the
// cloudsqlconn package.
package errtype

import "fmt"

type genericError struct {
	Message  string
	ConnName string
}

func (e *genericError) Error() string {
	return fmt.Sprintf("%v (connection name = %q)", e.Message, e.ConnName)
}

// NewConfigError initializes a ConfigError.
func NewConfigError(msg, cn string) *ConfigError {
	return &ConfigError{
		genericError: &genericError{Message: "Config error: " + msg, ConnName: cn},
	}
}

// ConfigError represents an incorrect request by the user. Config errors
// usually indicate a semantic error (e.g., the instance connection name is
// malformed, the SQL instance does not support the requested IP type, etc.)
// ConfigError's should not be retried.
type ConfigError struct{ *genericError }

// NewRefreshError initializes a RefreshError.
func NewRefreshError(msg, cn string, err error) *RefreshError {
	return &RefreshError{
		genericError: &genericError{Message: msg, ConnName: cn},
		Err:          err,
	}
}

// RefreshError means that an error occurred during the background
// refresh operation. In general, this is an unexpected error caused by
// an interaction with the API itself (e.g., missing certificates,
// invalid certificate encoding, region mismatch with the requested
// instance connection name, etc.). RefreshError's usually can be retried.
type RefreshError struct {
	*genericError
	// Err is the underlying error and may be nil.
	Err error
}

func (e *RefreshError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("Refresh error: %v", e.genericError)
	}
	return fmt.Sprintf("Refresh error: %v: %v", e.genericError, e.Err)
}

func (e *RefreshError) Unwrap() error { return e.Err }

// NewDialError initializes a DialError.
func NewDialError(msg, cn string, err error) *DialError {
	return &DialError{
		genericError: &genericError{Message: msg, ConnName: cn},
		Err:          err,
	}
}

// DialError represents a problem that occurred when trying to dial a SQL
// instance (e.g., a failure to set the keep-alive property, a TLS handshake
// failure, a missing certificate, etc.). DialError's are often network-related
// and can be retried.
type DialError struct {
	*genericError
	// Err is the underlying error and may be nil.
	Err error
}

func (e *DialError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("Dial error: %v", e.genericError)
	}
	return fmt.Sprintf("Dial error: %v: %v", e.genericError, e.Err)
}

func (e *DialError) Unwrap() error { return e.Err }
