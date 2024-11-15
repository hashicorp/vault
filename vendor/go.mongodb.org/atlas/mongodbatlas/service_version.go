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
	"net/http"
)

const versionPath = "api/private/unauth/version"

// ServiceVersionService is an interface for the version private endpoint of the MongoDB Atlas API.
//
// We currently make no promise to support or document this service or endpoint
// beyond what can be seen here.
type ServiceVersionService interface {
	Get(context.Context) (*ServiceVersion, *Response, error)
}

type ServiceVersionServiceOp struct {
	Client PlainRequestDoer
}

// Get gets the version information and parses it.
func (s *ServiceVersionServiceOp) Get(ctx context.Context) (*ServiceVersion, *Response, error) {
	req, err := s.Client.NewPlainRequest(ctx, http.MethodGet, versionPath)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return nil, nil, err
	}

	version := resp.ServiceVersion()

	return version, resp, nil
}
