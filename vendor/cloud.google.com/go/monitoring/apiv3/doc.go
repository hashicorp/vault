// Copyright 2020 Google LLC
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

// Package monitoring is an auto-generated package for the
// Cloud Monitoring API.
//
// Manages your Cloud Monitoring data and configurations. Most projects must
// be associated with a Workspace, with a few exceptions as noted on the
// individual method pages. The table entries below are presented in
// alphabetical order, not in order of common use. For explanations of the
// concepts found in the table entries, read the [Cloud Monitoring
// documentation](https://cloud.google.com/monitoring/docs).
//
// # Use of Context
//
// The ctx passed to NewClient is used for authentication requests and
// for creating the underlying connection, but is not used for subsequent calls.
// Individual methods on the client use the ctx given to them.
//
// To close the open connection, use the Close() method.
//
// For information about setting deadlines, reusing contexts, and more
// please visit godoc.org/cloud.google.com/go.
//
// Deprecated: Please use cloud.google.com/go/monitoring/apiv3/v2.
package monitoring // import "cloud.google.com/go/monitoring/apiv3"

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/grpc/metadata"
)

// For more information on implementing a client constructor hook, see
// https://github.com/googleapis/google-cloud-go/wiki/Customizing-constructors.
type clientHookParams struct{}
type clientHook func(context.Context, clientHookParams) ([]option.ClientOption, error)

var versionClient = "20220222"

func insertMetadata(ctx context.Context, mds ...metadata.MD) context.Context {
	out, _ := metadata.FromOutgoingContext(ctx)
	out = out.Copy()
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return metadata.NewOutgoingContext(ctx, out)
}

// DefaultAuthScopes reports the default set of authentication scopes to use with this package.
func DefaultAuthScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/monitoring",
		"https://www.googleapis.com/auth/monitoring.read",
		"https://www.googleapis.com/auth/monitoring.write",
	}
}
