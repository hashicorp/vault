/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"errors"
	"fmt"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/grpc/codes"
	proto3 "google.golang.org/protobuf/types/known/structpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// A Statement is a SQL query with named parameters.
//
// A parameter placeholder consists of '@' followed by the parameter name.
// The parameter name is an identifier which must conform to the naming
// requirements in https://cloud.google.com/spanner/docs/lexical#identifiers.
// Parameters may appear anywhere that a literal value is expected. The same
// parameter name may be used more than once.  It is an error to execute a
// statement with unbound parameters. On the other hand, it is allowable to
// bind parameter names that are not used.
//
// See the documentation of the Row type for how Go types are mapped to Cloud
// Spanner types.
type Statement struct {
	SQL    string
	Params map[string]interface{}
}

// NewStatement returns a Statement with the given SQL and an empty Params map.
func NewStatement(sql string) Statement {
	return Statement{SQL: sql, Params: map[string]interface{}{}}
}

// convertParams converts a statement's parameters into proto Param and
// ParamTypes.
func (s *Statement) convertParams() (*structpb.Struct, map[string]*sppb.Type, error) {
	params := &proto3.Struct{
		Fields: map[string]*proto3.Value{},
	}
	paramTypes := map[string]*sppb.Type{}
	for k, v := range s.Params {
		val, t, err := encodeValue(v)
		if err != nil {
			return nil, nil, errBindParam(k, v, err)
		}
		params.Fields[k] = val
		if t != nil {
			paramTypes[k] = t
		}
	}

	return params, paramTypes, nil
}

// errBindParam returns error for not being able to bind parameter to query
// request.
func errBindParam(k string, v interface{}, err error) error {
	if err == nil {
		return nil
	}
	var se *Error
	if !errors.As(err, &se) {
		return spannerErrorf(codes.InvalidArgument, "failed to bind query parameter(name: %q, value: %v), error = <%v>", k, v, err)
	}
	se.decorate(fmt.Sprintf("failed to bind query parameter(name: %q, value: %v)", k, v))
	return se
}
