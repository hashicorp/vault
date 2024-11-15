/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2016, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

// Package gocql implements a fast and robust Cassandra driver for the
// Go programming language.
//
// # Connecting to the cluster
//
// Pass a list of initial node IP addresses to NewCluster to create a new cluster configuration:
//
//	cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
//
// Port can be specified as part of the address, the above is equivalent to:
//
//	cluster := gocql.NewCluster("192.168.1.1:9042", "192.168.1.2:9042", "192.168.1.3:9042")
//
// It is recommended to use the value set in the Cassandra config for broadcast_address or listen_address,
// an IP address not a domain name. This is because events from Cassandra will use the configured IP
// address, which is used to index connected hosts. If the domain name specified resolves to more than 1 IP address
// then the driver may connect multiple times to the same host, and will not mark the node being down or up from events.
//
// Then you can customize more options (see ClusterConfig):
//
//	cluster.Keyspace = "example"
//	cluster.Consistency = gocql.Quorum
//	cluster.ProtoVersion = 4
//
// The driver tries to automatically detect the protocol version to use if not set, but you might want to set the
// protocol version explicitly, as it's not defined which version will be used in certain situations (for example
// during upgrade of the cluster when some of the nodes support different set of protocol versions than other nodes).
//
// The driver advertises the module name and version in the STARTUP message, so servers are able to detect the version.
// If you use replace directive in go.mod, the driver will send information about the replacement module instead.
//
// When ready, create a session from the configuration. Don't forget to Close the session once you are done with it:
//
//	session, err := cluster.CreateSession()
//	if err != nil {
//		return err
//	}
//	defer session.Close()
//
// # Authentication
//
// CQL protocol uses a SASL-based authentication mechanism and so consists of an exchange of server challenges and
// client response pairs. The details of the exchanged messages depend on the authenticator used.
//
// To use authentication, set ClusterConfig.Authenticator or ClusterConfig.AuthProvider.
//
// PasswordAuthenticator is provided to use for username/password authentication:
//
//	 cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
//	 cluster.Authenticator = gocql.PasswordAuthenticator{
//			Username: "user",
//			Password: "password"
//	 }
//	 session, err := cluster.CreateSession()
//	 if err != nil {
//	 	return err
//	 }
//	 defer session.Close()
//
// # Transport layer security
//
// It is possible to secure traffic between the client and server with TLS.
//
// To use TLS, set the ClusterConfig.SslOpts field. SslOptions embeds *tls.Config so you can set that directly.
// There are also helpers to load keys/certificates from files.
//
// Warning: Due to historical reasons, the SslOptions is insecure by default, so you need to set EnableHostVerification
// to true if no Config is set. Most users should set SslOptions.Config to a *tls.Config.
// SslOptions and Config.InsecureSkipVerify interact as follows:
//
//	Config.InsecureSkipVerify | EnableHostVerification | Result
//	Config is nil             | false                  | do not verify host
//	Config is nil             | true                   | verify host
//	false                     | false                  | verify host
//	true                      | false                  | do not verify host
//	false                     | true                   | verify host
//	true                      | true                   | verify host
//
// For example:
//
//	cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
//	cluster.SslOpts = &gocql.SslOptions{
//		EnableHostVerification: true,
//	}
//	session, err := cluster.CreateSession()
//	if err != nil {
//		return err
//	}
//	defer session.Close()
//
// # Data-center awareness and query routing
//
// To route queries to local DC first, use DCAwareRoundRobinPolicy. For example, if the datacenter you
// want to primarily connect is called dc1 (as configured in the database):
//
//	cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
//	cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("dc1")
//
// The driver can route queries to nodes that hold data replicas based on partition key (preferring local DC).
//
//	cluster := gocql.NewCluster("192.168.1.1", "192.168.1.2", "192.168.1.3")
//	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("dc1"))
//
// Note that TokenAwareHostPolicy can take options such as gocql.ShuffleReplicas and gocql.NonLocalReplicasFallback.
//
// We recommend running with a token aware host policy in production for maximum performance.
//
// The driver can only use token-aware routing for queries where all partition key columns are query parameters.
// For example, instead of
//
//	session.Query("select value from mytable where pk1 = 'abc' AND pk2 = ?", "def")
//
// use
//
//	session.Query("select value from mytable where pk1 = ? AND pk2 = ?", "abc", "def")
//
// # Rack-level awareness
//
// The DCAwareRoundRobinPolicy can be replaced with RackAwareRoundRobinPolicy, which takes two parameters, datacenter and rack.
//
// Instead of dividing hosts with two tiers (local datacenter and remote datacenters) it divides hosts into three
// (the local rack, the rest of the local datacenter, and everything else).
//
// RackAwareRoundRobinPolicy can be combined with TokenAwareHostPolicy in the same way as DCAwareRoundRobinPolicy.
//
// # Executing queries
//
// Create queries with Session.Query. Query values must not be reused between different executions and must not be
// modified after starting execution of the query.
//
// To execute a query without reading results, use Query.Exec:
//
//	 err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
//			"me", gocql.TimeUUID(), "hello world").WithContext(ctx).Exec()
//
// Single row can be read by calling Query.Scan:
//
//	 err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
//			"me").WithContext(ctx).Consistency(gocql.One).Scan(&id, &text)
//
// Multiple rows can be read using Iter.Scanner:
//
//	 scanner := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`,
//	 	"me").WithContext(ctx).Iter().Scanner()
//	 for scanner.Next() {
//	 	var (
//	 		id gocql.UUID
//			text string
//	 	)
//	 	err = scanner.Scan(&id, &text)
//	 	if err != nil {
//	 		log.Fatal(err)
//	 	}
//	 	fmt.Println("Tweet:", id, text)
//	 }
//	 // scanner.Err() closes the iterator, so scanner nor iter should be used afterwards.
//	 if err := scanner.Err(); err != nil {
//	 	log.Fatal(err)
//	 }
//
// See Example for complete example.
//
// # Prepared statements
//
// The driver automatically prepares DML queries (SELECT/INSERT/UPDATE/DELETE/BATCH statements) and maintains a cache
// of prepared statements.
// CQL protocol does not support preparing other query types.
//
// When using CQL protocol >= 4, it is possible to use gocql.UnsetValue as the bound value of a column.
// This will cause the database to ignore writing the column.
// The main advantage is the ability to keep the same prepared statement even when you don't
// want to update some fields, where before you needed to make another prepared statement.
//
// # Executing multiple queries concurrently
//
// Session is safe to use from multiple goroutines, so to execute multiple concurrent queries, just execute them
// from several worker goroutines. Gocql provides synchronously-looking API (as recommended for Go APIs) and the queries
// are executed asynchronously at the protocol level.
//
//	results := make(chan error, 2)
//	go func() {
//		results <- session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
//			"me", gocql.TimeUUID(), "hello world 1").Exec()
//	}()
//	go func() {
//		results <- session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
//			"me", gocql.TimeUUID(), "hello world 2").Exec()
//	}()
//
// # Nulls
//
// Null values are are unmarshalled as zero value of the type. If you need to distinguish for example between text
// column being null and empty string, you can unmarshal into *string variable instead of string.
//
//	var text *string
//	err := scanner.Scan(&text)
//	if err != nil {
//		// handle error
//	}
//	if text != nil {
//		// not null
//	}
//	else {
//		// null
//	}
//
// See Example_nulls for full example.
//
// # Reusing slices
//
// The driver reuses backing memory of slices when unmarshalling. This is an optimization so that a buffer does not
// need to be allocated for every processed row. However, you need to be careful when storing the slices to other
// memory structures.
//
//	scanner := session.Query(`SELECT myints FROM table WHERE pk = ?`, "key").WithContext(ctx).Iter().Scanner()
//	var myInts []int
//	for scanner.Next() {
//		// This scan reuses backing store of myInts for each row.
//		err = scanner.Scan(&myInts)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
// When you want to save the data for later use, pass a new slice every time. A common pattern is to declare the
// slice variable within the scanner loop:
//
//	scanner := session.Query(`SELECT myints FROM table WHERE pk = ?`, "key").WithContext(ctx).Iter().Scanner()
//	for scanner.Next() {
//		var myInts []int
//		// This scan always gets pointer to fresh myInts slice, so does not reuse memory.
//		err = scanner.Scan(&myInts)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
// # Paging
//
// The driver supports paging of results with automatic prefetch, see ClusterConfig.PageSize, Session.SetPrefetch,
// Query.PageSize, and Query.Prefetch.
//
// It is also possible to control the paging manually with Query.PageState (this disables automatic prefetch).
// Manual paging is useful if you want to store the page state externally, for example in a URL to allow users
// browse pages in a result. You might want to sign/encrypt the paging state when exposing it externally since
// it contains data from primary keys.
//
// Paging state is specific to the CQL protocol version and the exact query used. It is meant as opaque state that
// should not be modified. If you send paging state from different query or protocol version, then the behaviour
// is not defined (you might get unexpected results or an error from the server). For example, do not send paging state
// returned by node using protocol version 3 to a node using protocol version 4. Also, when using protocol version 4,
// paging state between Cassandra 2.2 and 3.0 is incompatible (https://issues.apache.org/jira/browse/CASSANDRA-10880).
//
// The driver does not check whether the paging state is from the same protocol version/statement.
// You might want to validate yourself as this could be a problem if you store paging state externally.
// For example, if you store paging state in a URL, the URLs might become broken when you upgrade your cluster.
//
// Call Query.PageState(nil) to fetch just the first page of the query results. Pass the page state returned by
// Iter.PageState to Query.PageState of a subsequent query to get the next page. If the length of slice returned
// by Iter.PageState is zero, there are no more pages available (or an error occurred).
//
// Using too low values of PageSize will negatively affect performance, a value below 100 is probably too low.
// While Cassandra returns exactly PageSize items (except for last page) in a page currently, the protocol authors
// explicitly reserved the right to return smaller or larger amount of items in a page for performance reasons, so don't
// rely on the page having the exact count of items.
//
// See Example_paging for an example of manual paging.
//
// # Dynamic list of columns
//
// There are certain situations when you don't know the list of columns in advance, mainly when the query is supplied
// by the user. Iter.Columns, Iter.RowData, Iter.MapScan and Iter.SliceMap can be used to handle this case.
//
// See Example_dynamicColumns.
//
// # Batches
//
// The CQL protocol supports sending batches of DML statements (INSERT/UPDATE/DELETE) and so does gocql.
// Use Session.NewBatch to create a new batch and then fill-in details of individual queries.
// Then execute the batch with Session.ExecuteBatch.
//
// Logged batches ensure atomicity, either all or none of the operations in the batch will succeed, but they have
// overhead to ensure this property.
// Unlogged batches don't have the overhead of logged batches, but don't guarantee atomicity.
// Updates of counters are handled specially by Cassandra so batches of counter updates have to use CounterBatch type.
// A counter batch can only contain statements to update counters.
//
// For unlogged batches it is recommended to send only single-partition batches (i.e. all statements in the batch should
// involve only a single partition).
// Multi-partition batch needs to be split by the coordinator node and re-sent to
// correct nodes.
// With single-partition batches you can send the batch directly to the node for the partition without incurring the
// additional network hop.
//
// It is also possible to pass entire BEGIN BATCH .. APPLY BATCH statement to Query.Exec.
// There are differences how those are executed.
// BEGIN BATCH statement passed to Query.Exec is prepared as a whole in a single statement.
// Session.ExecuteBatch prepares individual statements in the batch.
// If you have variable-length batches using the same statement, using Session.ExecuteBatch is more efficient.
//
// See Example_batch for an example.
//
// # Lightweight transactions
//
// Query.ScanCAS or Query.MapScanCAS can be used to execute a single-statement lightweight transaction (an
// INSERT/UPDATE .. IF statement) and reading its result. See example for Query.MapScanCAS.
//
// Multiple-statement lightweight transactions can be executed as a logged batch that contains at least one conditional
// statement. All the conditions must return true for the batch to be applied. You can use Session.ExecuteBatchCAS and
// Session.MapExecuteBatchCAS when executing the batch to learn about the result of the LWT. See example for
// Session.MapExecuteBatchCAS.
//
// # Retries and speculative execution
//
// Queries can be marked as idempotent. Marking the query as idempotent tells the driver that the query can be executed
// multiple times without affecting its result. Non-idempotent queries are not eligible for retrying nor speculative
// execution.
//
// Idempotent queries are retried in case of errors based on the configured RetryPolicy.
//
// Queries can be retried even before they fail by setting a SpeculativeExecutionPolicy. The policy can
// cause the driver to retry on a different node if the query is taking longer than a specified delay even before the
// driver receives an error or timeout from the server. When a query is speculatively executed, the original execution
// is still executing. The two parallel executions of the query race to return a result, the first received result will
// be returned.
//
// # User-defined types
//
// UDTs can be mapped (un)marshaled from/to map[string]interface{} a Go struct (or a type implementing
// UDTUnmarshaler, UDTMarshaler, Unmarshaler or Marshaler interfaces).
//
// For structs, cql tag can be used to specify the CQL field name to be mapped to a struct field:
//
//	type MyUDT struct {
//		FieldA int32 `cql:"a"`
//		FieldB string `cql:"b"`
//	}
//
// See Example_userDefinedTypesMap, Example_userDefinedTypesStruct, ExampleUDTMarshaler, ExampleUDTUnmarshaler.
//
// # Metrics and tracing
//
// It is possible to provide observer implementations that could be used to gather metrics:
//
//   - QueryObserver for monitoring individual queries.
//   - BatchObserver for monitoring batch queries.
//   - ConnectObserver for monitoring new connections from the driver to the database.
//   - FrameHeaderObserver for monitoring individual protocol frames.
//
// CQL protocol also supports tracing of queries. When enabled, the database will write information about
// internal events that happened during execution of the query. You can use Query.Trace to request tracing and receive
// the session ID that the database used to store the trace information in system_traces.sessions and
// system_traces.events tables. NewTraceWriter returns an implementation of Tracer that writes the events to a writer.
// Gathering trace information might be essential for debugging and optimizing queries, but writing traces has overhead,
// so this feature should not be used on production systems with very high load unless you know what you are doing.
package gocql // import "github.com/gocql/gocql"
