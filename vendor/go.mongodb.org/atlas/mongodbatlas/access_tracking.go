// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const accessTrackingBasePath = "api/atlas/v1.0/groups/%s/dbAccessHistory"

// AccessTrackingService is an interface for interfacing with the Access Tracking endpoints of the MongoDB Atlas API.
//
// See more: https://docs.atlas.mongodb.com/reference/api/access-tracking/
type AccessTrackingService interface {
	ListByCluster(context.Context, string, string, *AccessLogOptions) (*AccessLogSettings, *Response, error)
	ListByHostname(context.Context, string, string, *AccessLogOptions) (*AccessLogSettings, *Response, error)
}

// AccessTrackingServiceOp handles communication with the AccessTrackingService related methods of the
// MongoDB Atlas API.
type AccessTrackingServiceOp service

var _ AccessTrackingService = &AccessTrackingServiceOp{}

// AccessLogOptions represents the query options of AccessTrackingService.List.
type AccessLogOptions struct {
	Start      string `url:"start,omitempty"`      // Start is the timestamp in the number of milliseconds that have elapsed since the UNIX epoch for the first entry that Atlas returns from the database access logs.
	End        string `url:"end,omitempty"`        // End is the timestamp in the number of milliseconds that have elapsed since the UNIX epoch for the last entry that Atlas returns from the database access logs.
	NLogs      int    `url:"nLogs,omitempty"`      // NLogs is the maximum number of log entries to return. Atlas accepts values between 0 and 20000, inclusive.
	IPAddress  string `url:"ipAddress,omitempty"`  // IPAddress is the single IP address that attempted to authenticate with the database. Atlas filters the returned logs to include documents with only this IP address.
	AuthResult *bool  `url:"authResult,omitempty"` // AuthResult indicates whether to return either successful or failed authentication attempts. When set to true, Atlas filters the log to return only successful authentication attempts. When set to false, Atlas filters the log to return only failed authentication attempts.
}

// AccessLogs represents authentication attempts made against the cluster.
type AccessLogs struct {
	GroupID       string `json:"groupId,omitempty"`       // GroupID is the unique identifier for the project.
	Hostname      string `json:"hostname,omitempty"`      // Hostname is the hostname of the target node that received the authentication attempt.
	ClusterName   string `json:"clusterName,omitempty"`   // ClusterName is the name associated with the cluster.
	IPAddress     string `json:"ipAddress,omitempty"`     // IPAddress is the IP address that the authentication attempt originated from.
	AuthResult    *bool  `json:"authResult,omitempty"`    // AuthResult is the result of the authentication attempt. Returns true if the authentication request was successful. Returns false if the authentication request resulted in failure.
	LogLine       string `json:"logLine,omitempty"`       // LogLine is the text of the server log concerning the authentication attempt.
	Timestamp     string `json:"timestamp,omitempty"`     // Timestamp is the UTC timestamp of the authentication attempt.
	Username      string `json:"username,omitempty"`      // Username is the username that attempted to authenticate.
	FailureReason string `json:"failureReason,omitempty"` // FailureReason is the reason that the request failed to authenticate. Returns null if the authentication request was successful.
	AuthSource    string `json:"authSource,omitempty"`    // AuthSource is the database that the request attempted to authenticate against. Returns admin if the authentication source for the user is SCRAM-SHA. Returns $external if the authentication source for the user is LDAP.
}

// AccessLogSettings represents database access history settings.
type AccessLogSettings struct {
	AccessLogs []*AccessLogs `json:"accessLogs,omitempty"` // AccessLogs contains the authentication attempts made against the cluster.
}

// ListByCluster retrieves the access logs of a cluster by hostname.
//
// See more: https://docs.atlas.mongodb.com/reference/api/access-tracking-get-database-history-hostname/
func (s *AccessTrackingServiceOp) ListByCluster(ctx context.Context, groupID, clusterName string, opts *AccessLogOptions) (*AccessLogSettings, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(accessTrackingBasePath, groupID)
	path := fmt.Sprintf("%s/clusters/%s", basePath, clusterName)
	path, err := setListOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root *AccessLogSettings
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// ListByHostname retrieves the access logs of a cluster by hostname.
//
// See more: https://docs.atlas.mongodb.com/reference/api/access-tracking-get-database-history-hostname/
func (s *AccessTrackingServiceOp) ListByHostname(ctx context.Context, groupID, hostname string, opts *AccessLogOptions) (*AccessLogSettings, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if hostname == "" {
		return nil, nil, NewArgError("hostname", "must be set")
	}

	basePath := fmt.Sprintf(accessTrackingBasePath, groupID)
	path := fmt.Sprintf("%s/processes/%s", basePath, hostname)
	path, err := setListOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root *AccessLogSettings
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
