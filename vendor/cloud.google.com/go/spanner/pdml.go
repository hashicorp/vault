// Copyright 2018 Google LLC
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

package spanner

import (
	"context"

	"cloud.google.com/go/internal/trace"
	"github.com/googleapis/gax-go/v2"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
)

// PartitionedUpdate executes a DML statement in parallel across the database,
// using separate, internal transactions that commit independently. The DML
// statement must be fully partitionable: it must be expressible as the union
// of many statements each of which accesses only a single row of the table. The
// statement should also be idempotent, because it may be applied more than once.
//
// PartitionedUpdate returns an estimated count of the number of rows affected.
// The actual number of affected rows may be greater than the estimate.
func (c *Client) PartitionedUpdate(ctx context.Context, statement Statement) (count int64, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.PartitionedUpdate")
	defer func() { trace.EndSpan(ctx, err) }()
	if err := checkNestedTxn(ctx); err != nil {
		return 0, err
	}
	var (
		s  *session
		sh *sessionHandle
	)
	// Create session.
	s, err = c.sc.createSession(ctx)
	if err != nil {
		return 0, toSpannerError(err)
	}
	// Delete the session at the end of the request. If the PDML statement
	// timed out or was cancelled, the DeleteSession request might not succeed,
	// but the session will eventually be garbage collected by the server.
	defer s.delete(ctx)
	sh = &sessionHandle{session: s}
	// Create the parameters and the SQL request, but without a transaction.
	// The transaction reference will be added by the executePdml method.
	params, paramTypes, err := statement.convertParams()
	if err != nil {
		return 0, toSpannerError(err)
	}
	req := &sppb.ExecuteSqlRequest{
		Session:    sh.getID(),
		Sql:        statement.SQL,
		Params:     params,
		ParamTypes: paramTypes,
	}

	// Make a retryer for Aborted errors.
	// TODO: use generic Aborted retryer when merged with master
	retryer := gax.OnCodes([]codes.Code{codes.Aborted}, DefaultRetryBackoff)
	// Execute the PDML and retry if the transaction is aborted.
	executePdmlWithRetry := func(ctx context.Context) (int64, error) {
		for {
			count, err := executePdml(ctx, sh, req)
			if err == nil {
				return count, nil
			}
			delay, shouldRetry := retryer.Retry(err)
			if !shouldRetry {
				return 0, err
			}
			if err := gax.Sleep(ctx, delay); err != nil {
				return 0, err
			}
		}
	}
	return executePdmlWithRetry(ctx)
}

// executePdml executes the following steps:
// 1. Begin a PDML transaction
// 2. Add the ID of the PDML transaction to the SQL request.
// 3. Execute the update statement on the PDML transaction
//
// Note that PDML transactions cannot be committed or rolled back.
func executePdml(ctx context.Context, sh *sessionHandle, req *sppb.ExecuteSqlRequest) (count int64, err error) {
	// Begin transaction.
	res, err := sh.getClient().BeginTransaction(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.BeginTransactionRequest{
		Session: sh.getID(),
		Options: &sppb.TransactionOptions{
			Mode: &sppb.TransactionOptions_PartitionedDml_{PartitionedDml: &sppb.TransactionOptions_PartitionedDml{}},
		},
	})
	if err != nil {
		return 0, toSpannerError(err)
	}
	// Add a reference to the PDML transaction on the ExecuteSql request.
	req.Transaction = &sppb.TransactionSelector{
		Selector: &sppb.TransactionSelector_Id{Id: res.Id},
	}
	resultSet, err := sh.getClient().ExecuteSql(ctx, req)
	if err != nil {
		return 0, err
	}
	if resultSet.Stats == nil {
		return 0, spannerErrorf(codes.InvalidArgument, "query passed to Update: %q", req.Sql)
	}
	return extractRowCount(resultSet.Stats)
}
