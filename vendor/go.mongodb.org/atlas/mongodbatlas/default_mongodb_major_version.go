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
	"bytes"
	"context"
	"net/http"
)

const defaultMongoDBMajorVersionPath = "api/private/unauth/nds/defaultMongoDBMajorVersion"

// DefaultMongoDBMajorVersionService this service is to be used by other MongoDB tools
// to determine the current default major version of MongoDB Server in Atlas.
//
// We currently make no promise to support or document this service or endpoint
// beyond what can be seen here.
type DefaultMongoDBMajorVersionService interface {
	Get(context.Context) (string, *Response, error)
}

// DefaultMongoDBMajorVersionServiceOp is an implementation of DefaultMongoDBMajorVersionService.
type DefaultMongoDBMajorVersionServiceOp struct {
	Client PlainRequestDoer
}

var _ DefaultMongoDBMajorVersionService = &DefaultMongoDBMajorVersionServiceOp{}

// Get gets the current major MongoDB version in Atlas.
func (s *DefaultMongoDBMajorVersionServiceOp) Get(ctx context.Context) (string, *Response, error) {
	req, err := s.Client.NewPlainRequest(ctx, http.MethodGet, defaultMongoDBMajorVersionPath)
	if err != nil {
		return "", nil, err
	}
	root := new(bytes.Buffer)

	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}

	return root.String(), resp, err
}
