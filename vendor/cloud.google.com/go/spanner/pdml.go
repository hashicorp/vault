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
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/googleapis/gax-go/v2"
	"go.opencensus.io/tag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	return c.partitionedUpdate(ctx, statement, c.qo)
}

// PartitionedUpdateWithOptions executes a DML statement in parallel across the database,
// using separate, internal transactions that commit independently. The sql
// query execution will be optimized based on the given query options.
func (c *Client) PartitionedUpdateWithOptions(ctx context.Context, statement Statement, opts QueryOptions) (count int64, err error) {
	return c.partitionedUpdate(ctx, statement, c.qo.merge(opts))
}

func (c *Client) partitionedUpdate(ctx context.Context, statement Statement, options QueryOptions) (count int64, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.PartitionedUpdate")
	defer func() { trace.EndSpan(ctx, err) }()
	if err := checkNestedTxn(ctx); err != nil {
		return 0, err
	}

	sh, err := c.idleSessions.take(ctx)
	if err != nil {
		return 0, ToSpannerError(err)
	}
	if sh != nil {
		defer sh.recycle()
	}
	// Mark isLongRunningTransaction to true, as the session in case of partitioned dml can be long-running
	sh.mu.Lock()
	sh.eligibleForLongRunning = true
	sh.mu.Unlock()

	// Create the parameters and the SQL request, but without a transaction.
	// The transaction reference will be added by the executePdml method.
	params, paramTypes, err := statement.convertParams()
	if err != nil {
		return 0, ToSpannerError(err)
	}
	req := &sppb.ExecuteSqlRequest{
		Session:        sh.getID(),
		Sql:            statement.SQL,
		Params:         params,
		ParamTypes:     paramTypes,
		QueryOptions:   options.Options,
		RequestOptions: createRequestOptions(options.Priority, options.RequestTag, ""),
	}

	// Make a retryer for Aborted and certain Internal errors.
	retryer := onCodes(DefaultRetryBackoff, codes.Aborted, codes.Internal)
	// Execute the PDML and retry if the transaction is aborted.
	executePdmlWithRetry := func(ctx context.Context) (int64, error) {
		for {
			count, err := executePdml(contextWithOutgoingMetadata(ctx, sh.getMetadata(), c.disableRouteToLeader), sh, req, options)
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
func executePdml(ctx context.Context, sh *sessionHandle, req *sppb.ExecuteSqlRequest, options QueryOptions) (count int64, err error) {
	var md metadata.MD
	sh.updateLastUseTime()
	// Begin transaction.
	res, err := sh.getClient().BeginTransaction(ctx, &sppb.BeginTransactionRequest{
		Session: sh.getID(),
		Options: &sppb.TransactionOptions{
			Mode:                        &sppb.TransactionOptions_PartitionedDml_{PartitionedDml: &sppb.TransactionOptions_PartitionedDml{}},
			ExcludeTxnFromChangeStreams: options.ExcludeTxnFromChangeStreams,
		},
	})
	if err != nil {
		return 0, ToSpannerError(err)
	}
	// Add a reference to the PDML transaction on the ExecuteSql request.
	req.Transaction = &sppb.TransactionSelector{
		Selector: &sppb.TransactionSelector_Id{Id: res.Id},
	}

	sh.updateLastUseTime()
	resultSet, err := sh.getClient().ExecuteSql(ctx, req, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil && sh.session.pool != nil {
		err := captureGFELatencyStats(tag.NewContext(ctx, sh.session.pool.tagMap), md, "executePdml_ExecuteSql")
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	if metricErr := recordGFELatencyMetricsOT(ctx, md, "executePdml_ExecuteSql", sh.session.pool.otConfig); metricErr != nil {
		trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
	}
	if err != nil {
		return 0, err
	}

	if resultSet.Stats == nil {
		return 0, spannerErrorf(codes.InvalidArgument, "query passed to Update: %q", req.Sql)
	}
	return extractRowCount(resultSet.Stats)
}
