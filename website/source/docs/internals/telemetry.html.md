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

### vault.audit.log_request

**[S]** Summary (Nanoseconds): Duration of time taken by all audit log requests across all audit log backends

### vault.audit.log_response

**[S]** Summary (Nanoseconds): Duration of time taken by audit log responses across all audit log backends

Additionally, per audit log backend metrics such as those for a specific backend like `file` will be present as:

### vault.audit.file.log_request

**[S]** Summary (Nanoseconds): Duration of time taken by audit log requests for the file based audit backend mounted as `file`

### vault.audit.file.log_response

**[S]** Summary (Nanoseconds): Duration of time taken by audit log responses for the file based audit backend mounted as `file`

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

**[S]** Summary (Milliseconds): The total garbage collector pause time since Vault was last started

### vault.runtime.total_gc_runs

**[G]** Gauge (Number of operations): Total number of garbage collection runs since Vault was last started

## Policy and Token Metrics

These metrics relate to policies and tokens.

### vault.expire.fetch-lease-times

**[S]** Summary (Milliseconds): Time taken to fetch lease times

### vault.expire.fetch-lease-times-by-token

**[S]** Summary (Milliseconds): Time taken to fetch lease times by token

### vault.expire.num_leases

**[G]** Gauge (Number of leases): Number of all leases which are eligible for eventual expiry

### vault.expire.revoke

**[S]** Summary (Milliseconds): Time taken to revoke a token

### vault.expire.revoke-force

**[S]** Summary (Milliseconds): Time taken to forcibly revoke a token

### vault.expire.revoke-prefix

**[S]** Summary (Milliseconds): Time taken to revoke tokens on a prefix

### vault.expire.revoke-by-token

**[S]** Summary (Milliseconds): Time taken to revoke all secrets issued with a given token

### vault.expire.renew

**[S]** Summary (Milliseconds): Time taken to renew a lease

### vault.expire.renew-token

**[S]** Summary (Milliseconds): Time taken to renew a token which does not need to invoke a logical backend

### vault.expire.register

**[S]** Summary (Milliseconds): Time taken for register operations

Thes operations take a request and response with an associated lease and register a lease entry with lease ID

### vault.expire.register-auth

**[S]** Summary (Milliseconds): Time taken for register authentication operations which create lease entries without lease ID

### vault.policy.get_policy

**[S]** Summary (Milliseconds): Time taken to get a policy

### vault.policy.list_policies

**[S]** Summary (Milliseconds): Time taken to list policies

### vault.policy.delete_policy

**[S]** Summary (Milliseconds): Time taken to delete a policy

### vault.policy.set_policy

**[S]** Summary (Milliseconds): Time taken to set a policy

### vault.token.create

**[S]** Summary (Milliseconds): The time taken to create a token

### vault.token.createAccessor

**[S]** Summary (Milliseconds): The time taken to create a token

### vault.token.lookup

**[S]** Summary (Milliseconds): The time taken to look up a token

### vault.token.revoke

**[S]** Summary (Milliseconds): Time taken to revoke a token

### vault.token.revoke-tree

**[S]** Summary (Milliseconds): Time taken to revoke a token tree

### vault.token.store

**[S]** Summary (Milliseconds): Time taken to store an updated token entry without writing to the secondary index

## Authentication Backend Metrics

These metrics relate to supported authentication methods.

### vault.rollback.attempt.auth-token-

**[S]** Summary (Milliseconds): Time taken to perform a rollback operation for the [token authentication backend][token-auth-backend]

### vault.rollback.attempt.auth-ldap-

**[S]** Summary (Milliseconds): Time taken to perform a rollback operation for the [LDAP authentication backend][ldap-auth-backend]

### vault.rollback.attempt.cubbyhole-

**[S]** Summary (Milliseconds): Time taken to perform a rollback operation for the [Cubbyhole secret backend][cubbyhole-secret-backend]

### vault.rollback.attempt.secret-

**[S]** Summary (Milliseconds): Time taken to perform a rollback operation for the [K/V secret backend][kv-secret-backend]

### vault.rollback.attempt.sys-

**[S]** Summary (Milliseconds): Time taken to perform a rollback operation for the system backend

### vault.route.rollback.auth-ldap-

**[S]** Summary (Milliseconds): Time taken to perform a route rollback operation for the [LDAP authentication backend][ldap-auth-backend]

### vault.route.rollback.auth-token-

**[S]** Summary (Milliseconds): Time taken to perform a route rollback operation for the [token authentication backend][token-auth-backend]

### vault.route.rollback.cubbyhole-

**[S]** Summary (Milliseconds): Time taken to perform a route rollback operation for the [Cubbyhole secret backend][cubbyhole-secret-backend]

### vault.route.rollback.secret-

**[S]** Summary (Milliseconds): Time taken to perform a route rollback operation for the [K/V secret backend][kv-secret-backend]

### vault.route.rollback.sys-

**[S]** Summary (Milliseconds): Time taken to perform a route rollback operation for the system backend

## Storage Backend Metrics

These metrics relate to the supported storage backends.

### vault.azure.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [Azure storage backend][azure-storage-backend]

### vault.azure.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [Azure storage backend][azure-storage-backend]

### vault.azure.delete

**[S]** Summary (Milliseconds): Duration of a DELETE operation against the [Azure storage backend][azure-storage-backend]

### vault.azure.list

**[S]** Summary (Milliseconds): Duration of a LIST operation against the [Azure storage backend][azure-storage-backend]

### vault.cassandra.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cassandra.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cassandra.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cassandra.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [Cassandra storage backend][cassandra-storage-backend]

### vault.cockroachdb.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.cockroachdb.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.cockroachdb.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.cockroachdb.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [CockroachDB storage backend][cockroachdb-storage-backend]

### vault.consul.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [Consul storage backend][consul-storage-backend]

### vault.consul.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [Consul storage backend][consul-storage-backend]

### vault.consul.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [Consul storage backend][consul-storage-backend]

### vault.consul.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [Consul storage backend][consul-storage-backend]

### vault.couchdb.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [CouchDB storage backend][couchdb-storage-backend]

### vault.couchdb.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [CouchDB storage backend][couchdb-storage-backend]

### vault.couchdb.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [CouchDB storage backend][couchdb-storage-backend]

### vault.couchdb.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [CouchDB storage backend][couchdb-storage-backend]

### vault.dynamodb.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.dynamodb.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.dynamodb.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.dynamodb.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [DynamoDB storage backend][dynamodb-storage-backend]

### vault.etcd.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [etcd storage backend][etcd-storage-backend]

### vault.etcd.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [etcd storage backend][etcd-storage-backend]

### vault.etcd.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [etcd storage backend][etcd-storage-backend]

### vault.etcd.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [etcd storage backend][etcd-storage-backend]

### vault.gcs.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.gcs.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.gcs.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.gcs.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [Google Cloud Storage storage backend][gcs-storage-backend]

### vault.mssql.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mssql.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mssql.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mssql.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [MS-SQL storage backend][mssql-storage-backend]

### vault.mysql.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [MySQL storage backend][mysql-storage-backend]

### vault.mysql.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [MySQL storage backend][mysql-storage-backend]

### vault.mysql.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [MySQL storage backend][mysql-storage-backend]

### vault.mysql.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [MySQL storage backend][mysql-storage-backend]

### vault.postgres.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.postgres.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.postgres.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.postgres.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [PostgreSQL storage backend][postgresql-storage-backend]

### vault.s3.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [Amazon S3 storage backend][s3-storage-backend]

### vault.s3.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [Amazon S3 storage backend][s3-storage-backend]

### vault.s3.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [Amazon S3 storage backend][s3-storage-backend]

### vault.s3.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [Amazon S3 storage backend][s3-storage-backend]

### vault.swift.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [Swift storage backend][swift-storage-backend]

### vault.swift.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [Swift storage backend][swift-storage-backend]

### vault.swift.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [Swift storage backend][swift-storage-backend]

### vault.swift.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [Swift storage backend][swift-storage-backend]

### vault.zookeeper.put

**[S]** Summary (Milliseconds): Duration of a PUT operation against the [ZooKeeper storage backend][zookeeper-storage-backend]

### vault.zookeeper.get

**[S]** Summary (Milliseconds): Duration of a GET operation against the [ZooKeeper storage backend][zookeeper-storage-backend]

### vault.zookeeper.delete

**[S]** Summary (Milliseconds):  Duration of a DELETE operation against the [ZooKeeper storage backend][zookeeper-storage-backend]

### vault.zookeeper.list

**[S]** Summary (Milliseconds):  Duration of a LIST operation against the [ZooKeeper storage backend][zookeeper-storage-backend]

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
