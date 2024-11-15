/*
Copyright 2018 Google LLC

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
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"cloud.google.com/go/internal/trace"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// BatchReadOnlyTransaction is a ReadOnlyTransaction that allows for exporting
// arbitrarily large amounts of data from Cloud Spanner databases.
// BatchReadOnlyTransaction partitions a read/query request. Read/query request
// can then be executed independently over each partition while observing the
// same snapshot of the database. BatchReadOnlyTransaction can also be shared
// across multiple clients by passing around the BatchReadOnlyTransactionID and
// then recreating the transaction using Client.BatchReadOnlyTransactionFromID.
//
// Note: if a client is used only to run partitions, you can
// create it using a ClientConfig with both MinOpened and MaxIdle set to
// zero to avoid creating unnecessary sessions. You can also avoid excess
// gRPC channels by setting ClientConfig.NumChannels to the number of
// concurrently active BatchReadOnlyTransactions you expect to have.
type BatchReadOnlyTransaction struct {
	ReadOnlyTransaction
	ID BatchReadOnlyTransactionID
}

// BatchReadOnlyTransactionID is a unique identifier for a
// BatchReadOnlyTransaction. It can be used to re-create a
// BatchReadOnlyTransaction on a different machine or process by calling
// Client.BatchReadOnlyTransactionFromID.
type BatchReadOnlyTransactionID struct {
	// unique ID for the transaction.
	tid transactionID
	// sid is the id of the Cloud Spanner session used for this transaction.
	sid string
	// rts is the read timestamp of this transaction.
	rts time.Time
}

// Partition defines a segment of data to be read in a batch read or query. A
// partition can be serialized and processed across several different machines
// or processes.
type Partition struct {
	pt   []byte
	qreq *sppb.ExecuteSqlRequest
	rreq *sppb.ReadRequest
}

// PartitionOptions specifies options for a PartitionQueryRequest and
// PartitionReadRequest. See
// https://godoc.org/google.golang.org/genproto/googleapis/spanner/v1#PartitionOptions
// for more details.
type PartitionOptions struct {
	// The desired data size for each partition generated.
	PartitionBytes int64
	// The desired maximum number of partitions to return.
	MaxPartitions int64
}

// toProto converts a spanner.PartitionOptions into a sppb.PartitionOptions
func (opt PartitionOptions) toProto() *sppb.PartitionOptions {
	return &sppb.PartitionOptions{
		PartitionSizeBytes: opt.PartitionBytes,
		MaxPartitions:      opt.MaxPartitions,
	}
}

// PartitionRead returns a list of Partitions that can be used to read rows from
// the database. These partitions can be executed across multiple processes,
// even across different machines. The partition size and count hints can be
// configured using PartitionOptions.
func (t *BatchReadOnlyTransaction) PartitionRead(ctx context.Context, table string, keys KeySet, columns []string, opt PartitionOptions) ([]*Partition, error) {
	return t.PartitionReadUsingIndex(ctx, table, "", keys, columns, opt)
}

// PartitionReadWithOptions returns a list of Partitions that can be used to
// read rows from the database. These partitions can be executed across multiple
// processes, even across different machines. The partition size and count hints
// can be configured using PartitionOptions. Pass a ReadOptions to modify the
// read operation.
func (t *BatchReadOnlyTransaction) PartitionReadWithOptions(ctx context.Context, table string, keys KeySet, columns []string, opt PartitionOptions, readOptions ReadOptions) ([]*Partition, error) {
	return t.PartitionReadUsingIndexWithOptions(ctx, table, "", keys, columns, opt, t.ReadOnlyTransaction.txReadOnly.ro.merge(readOptions))
}

// PartitionReadUsingIndex returns a list of Partitions that can be used to read
// rows from the database using an index.
func (t *BatchReadOnlyTransaction) PartitionReadUsingIndex(ctx context.Context, table, index string, keys KeySet, columns []string, opt PartitionOptions) ([]*Partition, error) {
	return t.PartitionReadUsingIndexWithOptions(ctx, table, index, keys, columns, opt, t.ReadOnlyTransaction.txReadOnly.ro)
}

// PartitionReadUsingIndexWithOptions returns a list of Partitions that can be
// used to read rows from the database using an index. Pass a ReadOptions to
// modify the read operation.
func (t *BatchReadOnlyTransaction) PartitionReadUsingIndexWithOptions(ctx context.Context, table, index string, keys KeySet, columns []string, opt PartitionOptions, readOptions ReadOptions) ([]*Partition, error) {
	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, err
	}
	sid, client := sh.getID(), sh.getClient()
	var (
		kset       *sppb.KeySet
		resp       *sppb.PartitionResponse
		partitions []*Partition
	)
	kset, err = keys.keySetProto()
	// Request partitions.
	if err != nil {
		return nil, err
	}
	var md metadata.MD
	sh.updateLastUseTime()
	resp, err = client.PartitionRead(contextWithOutgoingMetadata(ctx, sh.getMetadata(), t.disableRouteToLeader), &sppb.PartitionReadRequest{
		Session:          sid,
		Transaction:      ts,
		Table:            table,
		Index:            index,
		Columns:          columns,
		KeySet:           kset,
		PartitionOptions: opt.toProto(),
	}, gax.WithGRPCOptions(grpc.Header(&md)))

	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "PartitionReadUsingIndexWithOptions"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	if metricErr := recordGFELatencyMetricsOT(ctx, md, "PartitionReadUsingIndexWithOptions", t.otConfig); metricErr != nil {
		trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
	}
	// Prepare ReadRequest.
	req := &sppb.ReadRequest{
		Session:             sid,
		Transaction:         ts,
		Table:               table,
		Index:               index,
		Columns:             columns,
		KeySet:              kset,
		RequestOptions:      createRequestOptions(readOptions.Priority, readOptions.RequestTag, ""),
		DataBoostEnabled:    readOptions.DataBoostEnabled,
		DirectedReadOptions: readOptions.DirectedReadOptions,
	}
	// Generate partitions.
	for _, p := range resp.GetPartitions() {
		partitions = append(partitions, &Partition{
			pt:   p.PartitionToken,
			rreq: req,
		})
	}
	return partitions, err
}

// PartitionQuery returns a list of Partitions that can be used to execute a
// query against the database.
func (t *BatchReadOnlyTransaction) PartitionQuery(ctx context.Context, statement Statement, opt PartitionOptions) ([]*Partition, error) {
	return t.partitionQuery(ctx, statement, opt, t.ReadOnlyTransaction.txReadOnly.qo)
}

// PartitionQueryWithOptions returns a list of Partitions that can be used to
// execute a query against the database. The sql query execution will be
// optimized based on the given query options.
func (t *BatchReadOnlyTransaction) PartitionQueryWithOptions(ctx context.Context, statement Statement, opt PartitionOptions, qOpts QueryOptions) ([]*Partition, error) {
	return t.partitionQuery(ctx, statement, opt, t.ReadOnlyTransaction.txReadOnly.qo.merge(qOpts))
}

func (t *BatchReadOnlyTransaction) partitionQuery(ctx context.Context, statement Statement, opt PartitionOptions, qOpts QueryOptions) ([]*Partition, error) {
	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, err
	}
	sid, client := sh.getID(), sh.getClient()
	params, paramTypes, err := statement.convertParams()
	if err != nil {
		return nil, err
	}
	var md metadata.MD

	// request Partitions
	req := &sppb.PartitionQueryRequest{
		Session:          sid,
		Transaction:      ts,
		Sql:              statement.SQL,
		PartitionOptions: opt.toProto(),
		Params:           params,
		ParamTypes:       paramTypes,
	}
	sh.updateLastUseTime()
	resp, err := client.PartitionQuery(contextWithOutgoingMetadata(ctx, sh.getMetadata(), t.disableRouteToLeader), req, gax.WithGRPCOptions(grpc.Header(&md)))

	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "partitionQuery"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	if metricErr := recordGFELatencyMetricsOT(ctx, md, "partitionQuery", t.otConfig); metricErr != nil {
		trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
	}

	// prepare ExecuteSqlRequest
	r := &sppb.ExecuteSqlRequest{
		Session:             sid,
		Transaction:         ts,
		Sql:                 statement.SQL,
		Params:              params,
		ParamTypes:          paramTypes,
		QueryOptions:        qOpts.Options,
		RequestOptions:      createRequestOptions(qOpts.Priority, qOpts.RequestTag, ""),
		DataBoostEnabled:    qOpts.DataBoostEnabled,
		DirectedReadOptions: qOpts.DirectedReadOptions,
	}

	// generate Partitions
	var partitions []*Partition
	for _, p := range resp.GetPartitions() {
		partitions = append(partitions, &Partition{
			pt:   p.PartitionToken,
			qreq: r,
		})
	}
	return partitions, err
}

// release implements txReadEnv.release, noop.
func (t *BatchReadOnlyTransaction) release(err error) {
}

// setTimestamp implements txReadEnv.setTimestamp, noop.
//
// read timestamp is ready on txn initialization, avoid contending writing to it
// with future partitions.
func (t *BatchReadOnlyTransaction) setTimestamp(ts time.Time) {
}

// Close marks the txn as closed.
func (t *BatchReadOnlyTransaction) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.state = txClosed
}

// Cleanup cleans up all the resources used by this transaction and makes
// it unusable. Once this method is invoked, the transaction is no longer
// usable anywhere, including other clients/processes with which this
// transaction was shared.
//
// Calling Cleanup is optional, but recommended. If Cleanup is not called, the
// transaction's resources will be freed when the session expires on the backend
// and is deleted. For more information about recycled sessions, see
// https://cloud.google.com/spanner/docs/sessions.
func (t *BatchReadOnlyTransaction) Cleanup(ctx context.Context) {
	t.Close()
	t.mu.Lock()
	defer t.mu.Unlock()
	sh := t.sh
	if sh == nil {
		return
	}
	t.sh = nil
	sid, client := sh.getID(), sh.getClient()

	var md metadata.MD
	err := client.DeleteSession(contextWithOutgoingMetadata(ctx, sh.getMetadata(), true), &sppb.DeleteSessionRequest{Name: sid}, gax.WithGRPCOptions(grpc.Header(&md)))

	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "Cleanup"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	if metricErr := recordGFELatencyMetricsOT(ctx, md, "Cleanup", t.otConfig); metricErr != nil {
		trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
	}

	if err != nil {
		var logger *log.Logger
		if sh.session != nil {
			logger = sh.session.logger
		}
		logf(logger, "Failed to delete session %v. Error: %v", sid, err)
	}
}

// Execute runs a single Partition obtained from PartitionRead or
// PartitionQuery.
func (t *BatchReadOnlyTransaction) Execute(ctx context.Context, p *Partition) *RowIterator {
	var (
		sh  *sessionHandle
		err error
		rpc func(ct context.Context, resumeToken []byte) (streamingReceiver, error)
	)
	if sh, _, err = t.acquire(ctx); err != nil {
		return &RowIterator{err: err}
	}
	client := sh.getClient()
	if client == nil {
		// Might happen if transaction is closed in the middle of a API call.
		return &RowIterator{err: errSessionClosed(sh)}
	}
	sh.updateLastUseTime()
	// Read or query partition.
	if p.rreq != nil {
		rpc = func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
			client, err := client.StreamingRead(ctx, &sppb.ReadRequest{
				Session:             p.rreq.Session,
				Transaction:         p.rreq.Transaction,
				Table:               p.rreq.Table,
				Index:               p.rreq.Index,
				Columns:             p.rreq.Columns,
				KeySet:              p.rreq.KeySet,
				PartitionToken:      p.pt,
				RequestOptions:      p.rreq.RequestOptions,
				ResumeToken:         resumeToken,
				DataBoostEnabled:    p.rreq.DataBoostEnabled,
				DirectedReadOptions: p.rreq.DirectedReadOptions,
			})
			if err != nil {
				return client, err
			}
			md, err := client.Header()
			if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
				if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "Execute"); err != nil {
					trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
				}
			}
			if metricErr := recordGFELatencyMetricsOT(ctx, md, "Execute", t.otConfig); metricErr != nil {
				trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
			}
			return client, err
		}
	} else {
		rpc = func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
			client, err := client.ExecuteStreamingSql(ctx, &sppb.ExecuteSqlRequest{
				Session:             p.qreq.Session,
				Transaction:         p.qreq.Transaction,
				Sql:                 p.qreq.Sql,
				Params:              p.qreq.Params,
				ParamTypes:          p.qreq.ParamTypes,
				QueryOptions:        p.qreq.QueryOptions,
				PartitionToken:      p.pt,
				RequestOptions:      p.qreq.RequestOptions,
				ResumeToken:         resumeToken,
				DataBoostEnabled:    p.qreq.DataBoostEnabled,
				DirectedReadOptions: p.qreq.DirectedReadOptions,
			})
			if err != nil {
				return client, err
			}
			md, err := client.Header()

			if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
				if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "Execute"); err != nil {
					trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
				}
			}
			if metricErr := recordGFELatencyMetricsOT(ctx, md, "Execute", t.otConfig); metricErr != nil {
				trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
			}
			return client, err
		}
	}
	return stream(
		contextWithOutgoingMetadata(ctx, sh.getMetadata(), t.disableRouteToLeader),
		sh.session.logger,
		t.sp.sc.metricsTracerFactory,
		rpc,
		t.setTimestamp,
		t.release)
}

// MarshalBinary implements BinaryMarshaler.
func (tid BatchReadOnlyTransactionID) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(tid.tid); err != nil {
		return nil, err
	}
	if err := enc.Encode(tid.sid); err != nil {
		return nil, err
	}
	if err := enc.Encode(tid.rts); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements BinaryUnmarshaler.
func (tid *BatchReadOnlyTransactionID) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&tid.tid); err != nil {
		return err
	}
	if err := dec.Decode(&tid.sid); err != nil {
		return err
	}
	return dec.Decode(&tid.rts)
}

// MarshalBinary implements BinaryMarshaler.
func (p Partition) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p.pt); err != nil {
		return nil, err
	}
	var isReadPartition bool
	var req proto.Message
	if p.rreq != nil {
		isReadPartition = true
		req = p.rreq
	} else {
		isReadPartition = false
		req = p.qreq
	}
	if err := enc.Encode(isReadPartition); err != nil {
		return nil, err
	}
	if data, err = proto.Marshal(req); err != nil {
		return nil, err
	}
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements BinaryUnmarshaler.
func (p *Partition) UnmarshalBinary(data []byte) error {
	var (
		isReadPartition bool
		d               []byte
		err             error
	)
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&p.pt); err != nil {
		return err
	}
	if err := dec.Decode(&isReadPartition); err != nil {
		return err
	}
	if err := dec.Decode(&d); err != nil {
		return err
	}
	if isReadPartition {
		p.rreq = &sppb.ReadRequest{}
		err = proto.Unmarshal(d, p.rreq)
	} else {
		p.qreq = &sppb.ExecuteSqlRequest{}
		err = proto.Unmarshal(d, p.qreq)
	}
	return err
}

// GetPartitionToken returns partition token
func (p *Partition) GetPartitionToken() []byte {
	if p != nil {
		return p.pt
	}
	return nil
}
