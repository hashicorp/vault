---
layout: "docs"
page_title: "Telemetry"
sidebar_current: "docs-internals-telemetry"
description: |-
  Learn about the telemetry data available in Vault.
---

# Telemetry

The Vault server process collects various runtime metrics about the performance of different libraries and subsystems. These metrics are aggregated on a ten second interval and are retained for one minute.

To view the raw data, you must send a signal to the Vault process: on Unix-style operating systems, this is `USR1` while on Windows it is `BREAK`. When the Vault process receives this signal it will dump the current telemetry information to the process's `stderr`.

This telemetry information can be used for debugging or otherwise getting a better view of what Vault is doing.

Telemetry information can also be streamed directly from Vault to a range of metrics aggregation solutions as described in the [telemetry Stanza documentation][telemetry-stanza].

The following is an example telemetry dump snippet:

```text
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.expire.num_leases': 5100.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.num_goroutines': 39.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.sys_bytes': 222746880.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.malloc_count': 109189192.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.free_count': 108408240.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.heap_objects': 780953.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.total_gc_runs': 232.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.alloc_bytes': 72954392.000
[2017-12-19 20:37:50 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.total_gc_pause_ns': 150293024.000
[2017-12-19 20:37:50 +0000 UTC][S] 'vault.merkle.flushDirty': Count: 100 Min: 0.008 Mean: 0.027 Max: 0.183 Stddev: 0.024 Sum: 2.681 LastUpdated: 2017-12-19 20:37:59.848733035 +0000 UTC m=+10463.692105920
[2017-12-19 20:37:50 +0000 UTC][S] 'vault.merkle.saveCheckpoint': Count: 4 Min: 0.021 Mean: 0.054 Max: 0.110 Stddev: 0.039 Sum: 0.217 LastUpdated: 2017-12-19 20:37:57.048458148 +0000 UTC m=+10460.891835029
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.alloc_bytes': 73326136.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.sys_bytes': 222746880.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.malloc_count': 109195904.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.free_count': 108409568.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.heap_objects': 786342.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.total_gc_pause_ns': 150293024.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.expire.num_leases': 5100.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.num_goroutines': 39.000
[2017-12-19 20:38:00 +0000 UTC][G] 'vault.7f320e57f9fe.runtime.total_gc_runs': 232.000
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.route.rollback.consul-': Count: 1 Sum: 0.013 LastUpdated: 2017-12-19 20:38:01.968471579 +0000 UTC m=+10465.811842067
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.rollback.attempt.consul-': Count: 1 Sum: 0.073 LastUpdated: 2017-12-19 20:38:01.968502743 +0000 UTC m=+10465.811873131
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.rollback.attempt.pki-': Count: 1 Sum: 0.070 LastUpdated: 2017-12-19 20:38:01.96867005 +0000 UTC m=+10465.812041936
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.route.rollback.auth-app-id-': Count: 1 Sum: 0.012 LastUpdated: 2017-12-19 20:38:01.969146401 +0000 UTC m=+10465.812516689
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.rollback.attempt.identity-': Count: 1 Sum: 0.063 LastUpdated: 2017-12-19 20:38:01.968029888 +0000 UTC m=+10465.811400276
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.rollback.attempt.database-': Count: 1 Sum: 0.066 LastUpdated: 2017-12-19 20:38:01.969394215 +0000 UTC m=+10465.812764603
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.barrier.get': Count: 16 Min: 0.010 Mean: 0.015 Max: 0.031 Stddev: 0.005 Sum: 0.237 LastUpdated: 2017-12-19 20:38:01.983268118 +0000 UTC m=+10465.826637008
[2017-12-19 20:38:00 +0000 UTC][S] 'vault.merkle.flushDirty': Count: 100 Min: 0.006 Mean: 0.024 Max: 0.098 Stddev: 0.019 Sum: 2.386 LastUpdated: 2017-12-19 20:38:09.848158309 +0000 UTC m=+10473.691527099
```

You'll note that log entries are prefixed with the metric type as follows:

- **[C]** is a counter
- **[G]** is a gauge
- **[S]** is a summary


The following sections describe available Vault metrics. The metrics interval can be assumed to be 10 seconds when manually triggering metrics output using the above described signals.

## Internal Metrics

These metrics represent operational aspects of the running Vault instance.

<<<<<<< HEAD
| Metric           | Description                       | Unit | Type |
| ---------------- | ----------------------------------| ---- | ---- |
|`vault.audit.log_request`| This measures the number of audit log requests | Number of requests | Summary |
|`vault.audit.log_response`| This measures the number of audit log responses | Number of responses | Summary |
|`vault.barrier.delete`| This measures the number of delete operations at the barrier | Number of operations | Summary |
|`vault.barrier.get`| This measures the number of get operations at the barrier | Number of operations | Summary |
|`vault.barrier.put`| This measures the number of put operations at the barrier | Number of operations | Summary |
|`vault.barrier.list`| This measures the number of list operations at the barrier | Number of operations | Counter |
|`vault.core.check_token`| This measures the number of token checks | Number of checks | Summary |
|`vault.core.fetch_acl_and_token`| This measures the number of ACL and corresponding token entry fetches | Number of fetches | Summary |
|`vault.core.handle_request`| This measures the number of requests | Number of requests | Summary |
|`vault.core.handle_login_request`| This measures the number of login requests | Number of requests | Summary |
|`vault.core.leadership_setup_failed`| This measures the number of cluster leadership setup failures | Number of failures | Summary |
|`vault.core.leadership_lost`| This measures the number of cluster leadership losses | Number of losses | Summary |
|`vault.core.post_unseal` | This measures the number of post-unseal operations | Number of operations | Gauge |
|`vault.core.pre_seal`| This measures the number of pre-seal operations | Number of operations | Gauge |
|`vault.core.seal-with-request`| This measures the number of requested seal operations | Number of operations | Gauge |
|`vault.core.seal`| This measures the number of seal operations | Number of operations | Gauge |
|`vault.core.seal-internal`| This measures the number of internal seal operations | Number of operations | Gauge |
|`vault.core.step_down`| This measures the number of cluster leadership step downs | Number of stepdowns | Summary |
|`vault.core.unseal`| This measures the number of unseal operations | Number of operations | Summary |
|`vault.runtime.alloc_bytes` | This measures the number of bytes allocated by the Vault process. This may burst from time to time but should return to a steady state value.| Number of bytes | Gauge |
|`vault.runtime.free_count`| This measures the number of `free` operations | Number of operations | Gauge |
|`vault.runtime.heap_objects`| This measures the number of objects on the heap and is a good general memory pressure indicator | Number of heap objects | Gauge |
|`vault.runtime.malloc_count`| This measures the number of `malloc` operations | Number of operations | Gauge |
|`vault.runtime.num_goroutines`| This measures the number of goroutines and serves as a general load indicator | Number of goroutines| Gauge |
|`vault.runtime.sys_bytes`| This measures the number of bytes allocated to Vault and includes what is being used by the heap and what has been reclaimed but not given back| Number of bytes | Gauge |
|`vault.runtime.total_gc_pause_ns` | This measures the total garbage collector pause time since the Vault instance was last started | Nanosecond | Summary |
| `vault.runtime.total_gc_runs` | Total number of garbage collection runs since the Vault instance was last started  | Number of operations | Gauge |

### Policy and Token Metrics

These metrics relate to policies and tokens.

| Metric           | Description                       | Unit | Type |
| ---------------- | ----------------------------------| ---- | ---- |
`vault.expire.fetch-lease-times`| This measures the number of lease time fetch operations | Number of operations | Gauge |
`vault.expire.fetch-lease-times-by-token`| This measures the number of operations which compute lease times by token | Number of operations | Gauge |
`vault.expire.num_leases`| This measures the number of expired leases | Number of expired leases | Gauge |
`vault.expire.revoke`| This measures the number of revoke operations | Number of operations | Counter |
`vault.expire.revoke-force`| This measures the number of forced revoke operations | Number of operations | Counter |
`vault.expire.revoke-prefix`| This measures the number of operations used to revoke all secrets with a given prefix | Number of operations | Counter |
`vault.expire.revoke-by-token`| This measures the number of operations used to revoke all secrets issued with a given token | Number of operations | Counter |
`vault.expire.renew`| This measures the number of renew operations | Number of operations | Counter |
`vault.expire.renew-token`| This measures the number of renew token operations to renew a token which does not need to invoke a logical backend | Number of operations | Gauge |
`vault.expire.register`| This measures the number of register operations which  take a request and response with an associated lease and register a lease entry with lease ID | Number of operations | Gauge |
`vault.expire.register-auth`| This measures the number of register auth operations which create lease entries without lease ID | Number of operations | Gauge |
`vault.policy.get_policy`| This measures the number of policy get operations | Number of operations | Counter |
`vault.policy.list_policies`| This measures the number of policy list operations | Number of operations | Counter |
`vault.policy.delete_policy`| This measures the number of policy delete operations | Number of operations | Counter |
`vault.policy.set_policy`| This measures the number of policy set operations | Number of operations | Gauge |
`vault.token.create`| This measures the number of token create operations | Number of operations | Gauge |
`vault.token.createAccessor`| This measures the number of Token ID identifier operations | Number of operations | Gauge |
`vault.token.lookup`| This measures the number of token lookups | Number of lookups | Counter |
`vault.token.revoke`| This measures the number of token revocation operations | Number of operations | Gauge |
`vault.token.revoke-tree`| This measures the number of revoke tree operations | Number of operations | Gauge |
`vault.token.store`| This measures the number of operations to store an updated token entry without writing to the secondary index | Number of operations | Gauge |

### Auth Method Metrics
=======
### vault.audit.log_request

**[S]** Summary (Number of requests): Number of audit log requests

### vault.audit.log_response

**[S]** Summary (Number of responses): Number of audit log responses

### vault.audit.log_request_failure

**[C]** Counter (Number of failures): Number of audit log request failures

**NOTE**: This is a particularly important metric. Any non-zero value here indicates that there was a failure to make an audit log request to any of the configured audit log backends; **when Vault cannot log to any of the configured audit log backends it ceases all user operations**, and you should begin troubleshooting the audit log backends immediately if this metric continually increases.

### vault.audit.log_response_failure

**[C]** Counter (Number of failures): Number of audit log response failures

**NOTE**: This is a particularly important metric. Any non-zero value here indicates that there was a failure to receive a response to a request made to one of the configured audit log backends; **when Vault cannot log to any of the configured audit log backends it ceases all user operations**, and you should begin troubleshooting the audit log backends immediately if this metric continually increases.

### vault.barrier.delete

**[S]** Summary (Number of operations): Number of DELETE operations at the barrier

### vault.barrier.get

**[S]** Summary (Number of operations): Number of GET operations at the barrier

### vault.barrier.put

**[S]** Summary (Number of operations): Number of PUT operations at the barrier

### vault.barrier.list

**[S]** Summary (Number of operations): Number of LIST operations at the barrier

### vault.core.check_token

**[S]** Summary (Number of checks): Number of token checks handled by Vault core

### vault.core.fetch_acl_and_token

**[S]** Summary (Number of fetches): Number of ACL and corresponding token entry fetches handled by Vault core

### vault.core.handle_request

**[S]** Summary (Number of requests) Number of requests handled by Vault core

### vault.core.handle_login_request

**[S]** Summary (Number of requests): Number of login requests handled by Vault core

### vault.core.leadership_setup_failed

**[S]** Summary (Number of failures): Number of cluster leadership setup failures which have occurred in a highly available Vault cluster

This should be monitored and alerted on for overall cluster leadership status

### vault.core.leadership_lost

**[S]** Summary (Number of losses): Number of cluster leadership losses which have occurred in a highly available Vault cluster

This should be monitored and alerted on for overall cluster leadership status

### vault.core.post_unseal

**[G]** Gauge (Number of operations): Number of post-unseal operations handled by Vault core

### vault.core.pre_seal

**[G]** Gauge (Number of operations) Number of pre-seal operations

### vault.core.seal-with-request

**[G]** Gauge (Number of operations): Number of requested seal operations

### vault.core.seal

**[G]** Gauge (Number of operations): Number of seal operations

### vault.core.seal-internal

**[G]** Gauge (Number of operations): Number of internal seal operations

### vault.core.step_down

**[S]** Summary (Number of step downs): Number of cluster leadership step downs

This should be monitored and alerted on for overall cluster leadership status

### vault.core.unseal

**[S]** Summary (Number of operations): Number of unseal operations

### vault.runtime.alloc_bytes

**[G]** Gauge (Number of bytes): Number of bytes allocated by the Vault process.

This could burst from time to time, but should return to a steady state value.

### vault.runtime.free_count

**[G]** Gauge (Number of objects): Number of freed objects

### vault.runtime.heap_objects

**[G]** Gauge (Number of objects): Number of objects on the heap

This is a good general memory pressure indicator worth establishing a baseline and thresholds for alerting.

### vault.runtime.malloc_count

**[G]** Gauge (Number of objects): Cumulative count of allocated heap objects

### vault.runtime.num_goroutines

**[G]** Gauge (Number of goroutines): Number of goroutines

This serves as a general system load indicator worth establishing a baseline and thresholds for alerting.

### vault.runtime.sys_bytes

**[G]** Gauge (Number of bytes): Number of bytes allocated to Vault

This includes what is being used by Vault's heap and what has been reclaimed but not given back to the operating system.

### vault.runtime.total_gc_pause_ns

**[S]** Summary (Nanoseconds): The total garbage collector pause time since Vault was last started

### vault.runtime.total_gc_runs

**[G]** Gauge (Number of operations): Total number of garbage collection runs since Vault was last started

## Policy and Token Metrics

These metrics relate to policies and tokens.

### vault.expire.fetch-lease-times

**[S]** Summary (Nanoseconds): Time taken to fetch lease times

### vault.expire.fetch-lease-times-by-token

**[S]** Summary (Nanoseconds): Time taken to fetch lease times by token

### vault.expire.num_leases

**[G]** Gauge (Number of leases): Number of all leases which are eligible for eventual expiry

### vault.expire.revoke

**[S]** Summary (Nanoseconds): Time taken to revoke a token

### vault.expire.revoke-force

**[S]** Summary (Nanoseconds): Time taken to forcibly revoke a token

### vault.expire.revoke-prefix

**[S]** Summary (Nanoseconds): Time taken to revoke tokens on a prefix

### vault.expire.revoke-by-token

**[S]** Summary (Nanoseconds): Time taken to revoke all secrets issued with a given token

### vault.expire.renew

**[S]** Summary (Nanoseconds): Time taken to renew a lease

### vault.expire.renew-token

**[S]** Summary (Nanoseconds): Time taken to renew a token which does not need to invoke a logical backend

### vault.expire.register

**[S]** Summary (Nanoseconds): Time taken for register operations

Thes operations take a request and response with an associated lease and register a lease entry with lease ID

### vault.expire.register-auth

**[S]** Summary (Nanoseconds): Time taken for register authentication operations which create lease entries without lease ID

### vault.policy.get_policy

**[S]** Summary (Nanoseconds): Time taken to get a policy

### vault.policy.list_policies

**[S]** Summary (Nanoseconds): Time taken to list policies

### vault.policy.delete_policy

**[S]** Summary (Nanoseconds): Time taken to delete a policy

### vault.policy.set_policy

**[S]** Summary (Nanoseconds): Time taken to set a policy

### vault.token.create

**[S]** Summary (Nanoseconds): The time taken to create a token

### vault.token.createAccessor

**[S]** Summary (Nanoseconds): The time taken to create a token

### vault.token.lookup

**[S]** Summary (Nanoseconds): The time taken to look up a token

### vault.token.revoke

**[S]** Summary (Nanoseconds): Time taken to revoke a token

### vault.token.revoke-tree

**[S]** Summary (Nanoseconds): Time taken to revoke a token tree

### vault.token.store

**[S]** Summary (Nanoseconds): Time taken to store an updated token entry without writing to the secondary index

## Authentication Backend Metrics

These metrics relate to supported auth methods.

### vault.rollback.attempt.auth-token-

**[S]** Summary (Nanoseconds): Time taken to perform a rollback operation for the [token authentication backend][token-auth-backend]

### vault.rollback.attempt.auth-ldap-

**[S]** Summary (Nanoseconds): Time taken to perform a rollback operation for the [LDAP authentication backend][ldap-auth-backend]

### vault.rollback.attempt.cubbyhole-

**[S]** Summary (Nanoseconds): Time taken to perform a rollback operation for the [Cubbyhole secret backend][cubbyhole-secret-backend]

### vault.rollback.attempt.secret-

**[S]** Summary (Nanoseconds): Time taken to perform a rollback operation for the [K/V secret backend][kv-secret-backend]

### vault.rollback.attempt.sys-

**[S]** Summary (Nanoseconds): Time taken to perform a rollback operation for the system backend

### vault.route.rollback.auth-ldap-

**[S]** Summary (Nanoseconds): Time taken to perform a route rollback operation for the [LDAP authentication backend][ldap-auth-backend]

### vault.route.rollback.auth-token-

**[S]** Summary (Nanoseconds): Time taken to perform a route rollback operation for the [token authentication backend][token-auth-backend]

### vault.route.rollback.cubbyhole-

**[S]** Summary (Nanoseconds): Time taken to perform a route rollback operation for the [Cubbyhole secret backend][cubbyhole-secret-backend]

### vault.route.rollback.secret-

**[S]** Summary (Nanoseconds): Time taken to perform a route rollback operation for the [K/V secret backend][kv-secret-backend]

### vault.route.rollback.sys-

**[S]** Summary (Nanoseconds): Time taken to perform a route rollback operation for the system backend

## Storage Backend Metrics

These metrics relate to the supported storage backends.

### vault.azure.put

**[S]** Summary (Number of operations): Number of put operations against the [Azure storage backend][azure-storage-backend]

### vault.azure.get

**[S]** Summary (Number of operations):Number of get operations against the [
Azure storage backend][azure-storage-backend]

### vault.azure.delete

**[S]** Summary (Number of operations):Number of delete operations against the [Azure storage backend][azure-storage-backend]

### vault.azure.list

**[S]** Summary (Number of operations):Number of list operations against the [Azure storage backend][azure-storage-backend]

### vault.cassandra.put

**[S]** Summary (Number of operations): Number of PUT operations against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cassandra.get

**[S]** Summary (Number of operations): Number of GET operations against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cassandra.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cassandra.list

**[S]** Summary (Number of operations): Number of LIST operations against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cockroachdb.put

**[S]** Summary (Number of operations): Number of PUT operations against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.cockroachdb.get

**[S]** Summary (Number of operations): Number of GET operations against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.cockroachdb.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.cockroachdb.list

**[S]** Summary (Number of operations): Number of LIST operations against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.consul.put

**[S]** Summary (Number of operations): Number of PUT operations against the [Consul storage backend][consul-storage-backend]

### vault.consul.get

**[S]** Summary (Number of operations): Number of GET operations against the [Consul storage backend][consul-storage-backend]

### vault.consul.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [Consul storage backend][consul-storage-backend]

### vault.consul.list

**[S]** Summary (Number of operations): Number of LIST operations against the [Consul storage backend][consul-storage-backend]

### vault.couchdb.put

**[S]** Summary (Number of operations): Number of PUT operations against the [CouchDB storage backend][couchdb-storage-backend]

### vault.couchdb.get

**[S]** Summary (Number of operations): Number of GET operations against the [CouchDB storage backend][couchdb-storage-backend]

### vault.couchdb.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [CouchDB storage backend][couchdb-storage-backend]

### vault.couchdb.list

**[S]** Summary (Number of operations): Number of LIST operations against the [CouchDB storage backend][couchdb-storage-backend]

### vault.dynamodb.put

**[S]** Summary (Number of operations): Number of PUT operations against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.dynamodb.get

**[S]** Summary (Number of operations): Number of GET operations against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.dynamodb.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.dynamodb.list

**[S]** Summary (Number of operations): Number of LIST operations against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.etcd.put

**[S]** Summary (Number of operations): Number of PUT operations against the [etcd storage backend][etcd-storage-backend]

### vault.etcd.get

**[S]** Summary (Number of operations): Number of GET operations against the [etcd storage backend][etcd-storage-backend]

### vault.etcd.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [etcd storage backend][etcd-storage-backend]

### vault.etcd.list

**[S]** Summary (Number of operations): Number of LIST operations against the [etcd storage backend][etcd-storage-backend]

### vault.gcs.put

**[S]** Summary (Number of operations): Number of PUT operations against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.gcs.get

**[S]** Summary (Number of operations): Number of GET operations against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.gcs.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.gcs.list

**[S]** Summary (Number of operations): Number of LIST operations against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.mssql.put

**[S]** Summary (Number of operations): Number of PUT operations against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mssql.get

**[S]** Summary (Number of operations): Number of GET operations against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mssql.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mssql.list

**[S]** Summary (Number of operations): Number of LIST operations against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mysql.put

**[S]** Summary (Number of operations): Number of PUT operations against the [MySQL storage backend][mysql-storage-backend]

### vault.mysql.get

**[S]** Summary (Number of operations): Number of GET operations against the [MySQL storage backend][mysql-storage-backend]

### vault.mysql.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [MySQL storage backend][mysql-storage-backend]

### vault.mysql.list

**[S]** Summary (Number of operations): Number of LIST operations against the [MySQL storage backend][mysql-storage-backend]

### vault.postgres.put

**[S]** Summary (Number of operations): Number of PUT operations against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.postgres.get

**[S]** Summary (Number of operations): Number of GET operations against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.postgres.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.postgres.list

**[S]** Summary (Number of operations): Number of LIST operations against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.s3.put

**[S]** Summary (Number of operations): Number of PUT operations against the [Amazon S3 storage backend][s3-storage-backend]

### vault.s3.get

**[S]** Summary (Number of operations): Number of GET operations against the [Amazon S3 storage backend][s3-storage-backend]

### vault.s3.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [Amazon S3 storage backend][s3-storage-backend]

### vault.s3.list

**[S]** Summary (Number of operations): Number of LIST operations against the [Amazon S3 storage backend][s3-storage-backend]

### vault.swift.put

**[S]** Summary (Number of operations): Number of PUT operations against the [Swift storage backend][swift-storage-backend]

### vault.swift.get

**[S]** Summary (Number of operations): Number of GET operations against the [Swift storage backend][swift-storage-backend]

### vault.swift.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [Swift storage backend][swift-storage-backend]

### vault.swift.list

**[S]** Summary (Number of operations): Number of LIST operations against the [Swift storage backend][swift-storage-backend]

### vault.zookeeper.put

**[S]** Summary (Number of operations): Number of PUT operations against the [ZooKeeper storage backend][zookeeper-storage-backend]

### vault.zookeeper.get

**[S]** Summary (Number of operations): Number of GET operations against the [ZooKeeper storage backend][zookeeper-storage-backend]

### vault.zookeeper.delete

**[S]** Summary (Number of operations): Number of DELETE operations against the [ZooKeeper storage backend][zookeeper-storage-backend]

### vault.zookeeper.list

**[S]** Summary (Number of operations): Number of LIST operations against the [ZooKeeper storage backend][zookeeper-storage-backend]

[telemetry-stanza]: /docs/configuration/telemetry.html
[cubbyhole-secret-backend]: /docs/secrets/cubbyhole/index.html
[kv-secret-backend]: /docs/secrets/kv/index.html
[ldap-auth-backend]: /docs/auth/ldap.html
[token-auth-backend]: /docs/auth/token.html
[azure-storage-backend]: /docs/configuration/storage/azure.html
[cassandra-storage-backend]: /docs/configuration/storage/cassandra.html
[cockroachdb-storage-backend]: /docs/configuration/storage/cockroachdb.html
[consul-storage-backend]: /docs/configuration/storage/consul.html
[couchdb-storage-backend]: /docs/configuration/storage/couchdb.html
[dynamodb-storage-backend]: /docs/configuration/storage/dynamodb.html
[etcd-storage-backend]: /docs/configuration/storage/etcd.html
[gcs-storage-backend]: /docs/configuration/storage/google-cloud.html
[mssql-storage-backend]: /docs/configuration/storage/mssql.html
[mysql-storage-backend]: /docs/configuration/storage/mysql.html
[postgresql-storage-backend]: /docs/configuration/storage/postgresql.html
[s3-storage-backend]: /docs/configuration/storage/s3.html
[swift-storage-backend]: /docs/configuration/storage/swift.html
[zookeeper-storage-backend]: /docs/configuration/storage/zookeeper.html
