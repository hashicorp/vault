---
layout: "docs"
page_title: "Telemetry"
sidebar_title: "Telemetry"
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

## Audit Metrics

These metrics relate to auditing.

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.audit.log_request`| Duration of time taken by all audit log requests across all audit log devices | ms | summary |
| `vault.audit.log_response`| Duration of time taken by audit log responses across all audit log devices | ms | summary |
| `vault.audit.log_request_failure` | Number of audit log request failures.  **NOTE**: This is a particularly important metric. Any non-zero value here indicates that there was a failure to make an audit log request to any of the configured audit log devices; **when Vault cannot log to any of the configured audit log devices it ceases all user operations**, and you should begin troubleshooting the audit log devices immediately if this metric continually increases. | failures | counter |
| `vault.audit.log_response_failure` | Number of audit log response failures. **NOTE**: This is a particularly important metric. Any non-zero value here indicates that there was a failure to receive a response to a request made to one of the configured audit log devices; **when Vault cannot log to any of the configured audit log devices it ceases all user operations**, and you should begin troubleshooting the audit log devices immediately if this metric continually increases. | failures | counter |

**NOTE:** In addition, there are audit metrics for each enabled audit device represented as `vault.audit.<type>.log_request`.  For example, if a file audit device is enabled, its metrics would be `vault.audit.file.log_request` and `vault.audit.file.log_response` .

## Core Metrics

These metrics represent operational aspects of the running Vault instance.

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.barrier.delete` | Duration of time taken by DELETE operations at the barrier | ms | summary |
| `vault.barrier.get` | Duration of time taken by GET operations at the barrier | ms | summary |
| `vault.barrier.put` | Duration of time taken by PUT operations at the barrier | ms | summary |
| `vault.barrier.list` | Duration of time taken by LIST operations at the barrier | ms | summary |
| `vault.core.check_token` | Duration of time taken by token checks handled by Vault core | ms | summary |
| `vault.core.fetch_acl_and_token` | Duration of time taken by ACL and corresponding token entry fetches handled by Vault core | ms | summary |
| `vault.core.handle_request` | Duration of time taken by requests handled by Vault core | ms | summary |
| `vault.core.handle_login_request` | Duration of time taken by login requests handled by Vault core | ms | summary |
| `vault.core.leadership_setup_failed` | Duration of time taken by cluster leadership setup failures which have occurred in a highly available Vault cluster. This should be monitored and alerted on for overall cluster leadership status. | ms | summary |
| `vault.core.leadership_lost` | Duration of time taken by cluster leadership losses which have occurred in a highly available Vault cluster.  This should be monitored and alerted on for overall cluster leadership status. | ms | summary |
| `vault.core.post_unseal` | Duration of time taken by post-unseal operations handled by Vault core | ms | gauge |
| `vault.core.pre_seal` | Duration of time taken by pre-seal operations | ms | gauge |
| `vault.core.seal-with-request` | Duration of time taken by requested seal operations | ms | gauge |
| `vault.core.seal` | Duration of time taken by seal operations | ms | gauge |
| `vault.core.seal-internal` | Duration of time taken by internal seal operations | ms | gauge |
| `vault.core.step_down` | Duration of time taken by cluster leadership step downs.  This should be monitored and alerted on for overall cluster leadership status. | ms | summary |
| `vault.core.unseal` | Duration of time taken by unseal operations | ms | summary |

## Runtime Metrics

These metrics represent runtime aspects of the running Vault instance.

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.runtime.alloc_bytes` | Number of bytes allocated by the Vault process.  This could burst from time to time, but should return to a steady state value. | bytes | gauge |
| `vault.runtime.free_count` | Number of freed objects | objects | gauge |
| `vault.runtime.heap_objects` | Number of objects on the heap.  This is a good general memory pressure indicator worth establishing a baseline and thresholds for alerting. | objects | gauge |
| `vault.runtime.malloc_count` | Cumulative count of allocated heap objects | objects | gauge |
| `vault.runtime.num_goroutines` | Number of goroutines.  This serves as a general system load indicator worth establishing a baseline and thresholds for alerting. | goroutines | gauge |
| `vault.runtime.sys_bytes` | Number of bytes allocated to Vault.  This includes what is being used by Vault's heap and what has been reclaimed but not given back to the operating system. | bytes | gauge |
| `vault.runtime.total_gc_pause_ns` | The total garbage collector pause time since Vault was last started | ms | summary |
| `vault.runtime.total_gc_runs` | Total number of garbage collection runs since Vault was last started | operations | gauge |

## Policy and Token Metrics

These metrics relate to policies and tokens.

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.expire.fetch-lease-times` | Time taken to fetch lease times | ms | summary |
| `vault.expire.fetch-lease-times-by-token` | Time taken to fetch lease times by token | ms | summary |
| `vault.expire.num_leases` | Number of all leases which are eligible for eventual expiry | leases | gauge |
| `vault.expire.revoke` | Time taken to revoke a token | ms | summary |
| `vault.expire.revoke-force` | Time taken to forcibly revoke a token | ms | summary |
| `vault.expire.revoke-prefix` | Time taken to revoke tokens on a prefix | ms | summary |
| `vault.expire.revoke-by-token` | Time taken to revoke all secrets issued with a given token | ms | summary |
| `vault.expire.renew` | Time taken to renew a lease | ms | summary |
| `vault.expire.renew-token` | Time taken to renew a token which does not need to invoke a logical backend | ms | summary |
| `vault.expire.register` | Time taken for register operations | ms | summary |

These operations take a request and response with an associated lease and register a lease entry with lease ID

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.expire.register-auth` | Time taken for register authentication operations which create lease entries without lease ID | ms | summary |
| `vault.policy.get_policy` | Time taken to get a policy | ms | summary |
| `vault.policy.list_policies` | Time taken to list policies | ms | summary |
| `vault.policy.delete_policy` | Time taken to delete a policy | ms | summary |
| `vault.policy.set_policy` | Time taken to set a policy | ms | summary |
| `vault.token.create` | The time taken to create a token | ms | summary |
| `vault.token.createAccessor` | The time taken to create a token accessor | ms | summary |
| `vault.token.lookup` | The time taken to look up a token | ms | summary |
| `vault.token.revoke` | Time taken to revoke a token | ms | summary |
| `vault.token.revoke-tree` | Time taken to revoke a token tree | ms | summary |
| `vault.token.store` | Time taken to store an updated token entry without writing to the secondary index | ms | summary |

## Auth Methods Metrics

These metrics relate to supported authentication methods.

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.rollback.attempt.auth-token` | Time taken to perform a rollback operation for the [token auth method][token-auth-backend] | ms | summary |
| `vault.rollback.attempt.auth-ldap` | Time taken to perform a rollback operation for the [LDAP auth method][ldap-auth-backend] | ms | summary |
| `vault.rollback.attempt.cubbyhole` | Time taken to perform a rollback operation for the [Cubbyhole secret backend][cubbyhole-secrets-engine] | ms | summary |
| `vault.rollback.attempt.secret` | Time taken to perform a rollback operation for the [K/V secret backend][kv-secrets-engine] | ms | summary |
| `vault.rollback.attempt.sys` | Time taken to perform a rollback operation for the system backend | ms | summary |
| `vault.route.rollback.auth-ldap` | Time taken to perform a route rollback operation for the [LDAP auth method][ldap-auth-backend] | ms | summary |
| `vault.route.rollback.auth-token` | Time taken to perform a route rollback operation for the [token auth method][token-auth-backend] | ms | summary |
| `vault.route.rollback.cubbyhole` | Time taken to perform a route rollback operation for the [Cubbyhole secret backend][cubbyhole-secrets-engine] | ms | summary |
| `vault.route.rollback.secret` | Time taken to perform a route rollback operation for the [K/V secret backend][kv-secrets-engine] | ms | summary |
| `vault.route.rollback.sys` | Time taken to perform a route rollback operation for the system backend | ms | summary |

## Merkle Tree and Write Ahead Log Metrics

These metrics relate to internal operations on Merkle Trees and Write Ahead Logs (WAL)

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.merkle_flushdirty` | Time taken to flush any dirty pages to cold storage | ms | summary |
| `vault.merkle_savecheckpoint` | Time taken to save the checkpoint | ms | summary |
| `vault.wal_deletewals` | Time taken to delete a Write Ahead Log (WAL) | ms | summary |
| `vault.wal_gc_deleted` | Number of Write Ahead Logs (WAL) deleted during each garbage collection run | WAL | counter |
| `vault.wal_gc_total` | Total Number of Write Ahead Logs (WAL) on disk | WAL | counter |
| `vault.wal_loadWAL` | Time taken to load a Write Ahead Log (WAL) | ms | summary |
| `vault.wal_persistwals` | Time taken to persist a Write Ahead Log (WAL) | ms | summary |
| `vault.wal_flushready` | Time taken to flush a ready Write Ahead Log (WAL) to storage | ms | summary |

## Replication Metrics

These metrics relate to [Vault Enterprise Replication](https://www.vaultproject.io/docs/enterprise/replication/index.html).

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `logshipper.streamWALs.missing_guard` | Number of incidences where the starting Merkle Tree index used to begin streaming WAL entries is not matched/found | missing guards | counter |
| `logshipper.streamWALs.guard_found` | Number of incidences where the starting Merkle Tree index used to begin streaming WAL entries is matched/found | found guards | counter |
| `replication.fetchRemoteKeys` | Time taken to fetch keys from a remote cluster participating in replication prior to Merkle Tree based delta generation | ms | summary |
| `replication.merkleDiff` | Time taken to perform a Merkle Tree based delta generation between the clusters participating in replication | ms | summary |
| `replication.merkleSync` | Time taken to perform a Merkle Tree based synchronization using the last delta generated between the clusters participating in replication | ms | summary |
| `replication.merkle.commit_index` | The last committed index in the Merkle Tree. | sequence number | gauge |
| `replication.wal.last_wal` | The index of the last WAL | sequence number | gauge |
| `replication.wal.last_dr_wal` | The index of the last DR WAL | sequence number | gauge |
| `replication.wal.last_performance_wal` | The index of the last Performance WAL | sequence number | gauge |
| `replication.fsm.last_remote_wal` | The index of the last remote WAL | sequence number | gauge |
| `replication.rpc.server.auth_request` | Duration of time taken by auth request | ms | summary |
| `replication.rpc.server.bootstrap_request` | Duration of time taken by bootstrap request | ms | summary |
| `replication.rpc.server.conflicting_pages_request` | Duration of time taken by conflicting pages request | ms | summary |
| `replication.rpc.server.echo` | Duration of time taken by echo | ms | summary |
| `replication.rpc.server.forwarding_request` | Duration of time taken by forwarding request | ms | summary |
| `replication.rpc.server.guard_hash_request` | Duration of time taken by guard hash request | ms | summary |
| `replication.rpc.server.persist_alias_request` | Duration of time taken by persist alias request | ms | summary |
| `replication.rpc.server.persist_persona_request` | Duration of time taken by persist persona request | ms | summary |
| `replication.rpc.server.stream_wals_request` | Duration of time taken by stream wals request | ms | summary |
| `replication.rpc.server.sub_page_hashes_request` | Duration of time taken by sub page hashes request | ms | summary |
| `replication.rpc.server.sync_counter_request` | Duration of time taken by sync counter request | ms | summary |
| `replication.rpc.server.upsert_group_request` | Duration of time taken by upsert group request | ms | summary |
| `replication.rpc.client.conflicting_pages` | Duration of time taken by client conflicting pages request | ms | summary |
| `replication.rpc.client.fetch_keys` | Duration of time taken by client fetch keys request | ms | summary |
| `replication.rpc.client.forward` | Duration of time taken by client forward request | ms | summary |
| `replication.rpc.client.guard_hash` | Duration of time taken by client guard hash request | ms | summary |
| `replication.rpc.client.persist_alias` | Duration of time taken by | ms | summary |
| `replication.rpc.client.register_auth` | Duration of time taken by client register auth request | ms | summary |
| `replication.rpc.client.register_lease` | Duration of time taken by client register lease request | ms | summary |
| `replication.rpc.client.stream_wals` | Duration of time taken by client s | ms | summary |
| `replication.rpc.client.sub_page_hashes` | Duration of time taken by client sub page hashes request | ms | summary |
| `replication.rpc.client.sync_counter` | Duration of time taken by client sync counter request | ms | summary |
| `replication.rpc.client.upsert_group` | Duration of time taken by client upstert group request | ms | summary |
| `replication.rpc.client.wrap_in_cubbyhole` | Duration of time taken by client wrap in cubbyhole request | ms | summary |
| `replication.rpc.dr.server.echo` | Duration of time taken by DR echo request | ms | summary |
| `replication.rpc.dr.server.fetch_keys_request` | Duration of time taken by DR fetch keys request | ms | summary |
| `replication.rpc.standby.server.echo` | Duration of time taken by standby echo request | ms | summary |
| `replication.rpc.standby.server.register_auth_request` | Duration of time taken by standby register auth request | ms | summary |
| `replication.rpc.standby.server.register_lease_request` | Duration of time taken by standby register lease request | ms | summary |
| `replication.rpc.standby.server.wrap_token_request` | Duration of time taken by standby wrap token request | ms | summary |

## Secrets Engines Metrics

These metrics relate to the supported [secrets engines][secrets-engines].

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `database.Initialize` | Time taken to initialize a database secret engine across all database secrets engines | ms | summary |
| `database.&lt;name&gt;.Initialize` | Time taken to initialize a database secret engine for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Initialize` | ms | summary |
| `database.Initialize.error` | Number of database secrets engine initialization operation errors across all database secrets engines | errors | counter |
| `database.&lt;name&gt;.Initialize.error` | Number of database secrets engine initialization operation errors for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Initialize.error` | errors | counter |
| `database.Close` | Time taken to close a database secret engine across all database secrets engines | ms | summary |
| `database.&lt;name&gt;.Close` | Time taken to close a database secret engine for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Close` | ms | summary |
| `database.Close.error` | Number of database secrets engine close operation errors across all database secrets engines | errors | counter |
| `database.&lt;name&gt;.Close.error` | Number of database secrets engine close operation errors for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Close.error` | errors | counter |
| `database.CreateUser` | Time taken to create a user across all database secrets engines | ms | summary |
| `database.&lt;name&gt;.CreateUser` | Time taken to create a user for the named database secrets engine `<name>` | ms | summary |
| `database.CreateUser.error` | Number of user creation operation errors across all database secrets engines | errors | counter |
| `database.&lt;name&gt;.CreateUser.error` | Number of user creation operation errors for the named database secrets engine `<name>`, for example: `database.postgresql-prod.CreateUser.error` | errors | counter |
| `database.RenewUser` | Time taken to renew a user across all database secrets engines | ms | summary |
| `database.&lt;name&gt;.RenewUser` | Time taken to renew a user for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RenewUser` | ms | summary |
| `database.RenewUser.error` | Number of user renewal operation errors across all database secrets engines | errors | counter |
| `database.&lt;name&gt;.RenewUser.error` | Number of user renewal operations for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RenewUser.error` | errors | counter |
| `database.RevokeUser` | Time taken to revoke a user across all database secrets engines | ms | summary |
| `database.&lt;name&gt;.RevokeUser` | Time taken to revoke a user for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RevokeUser` | ms | summary |
| `database.RevokeUser.error` | Number of user revocation operation errors across all database secrets engines | errors | counter |
| `database.&lt;name&gt;.RevokeUser.error` | Number of user revocation operations for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RevokeUser.error` | errors | counter |

## Storage Backend Metrics

These metrics relate to the supported [storage backends][storage-backends].

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.azure.put` | Duration of a PUT operation against the [Azure storage backend][azure-storage-backend] | ms | summary |
| `vault.azure.get` | Duration of a GET operation against the [Azure storage backend][azure-storage-backend] | ms | summary |
| `vault.azure.delete` | Duration of a DELETE operation against the [Azure storage backend][azure-storage-backend] | ms | summary |
| `vault.azure.list` | Duration of a LIST operation against the [Azure storage backend][azure-storage-backend] | ms | summary |
| `vault.cassandra.put` | Duration of a PUT operation against the [Cassandra storage backend][cassandra-storage-backend] | ms | summary |
| `vault.cassandra.get` | Duration of a GET operation against the [Cassandra storage backend][cassandra-storage-backend] | ms | summary |
| `vault.cassandra.delete` | Duration of a DELETE operation against the [Cassandra storage backend][cassandra-storage-backend] | ms | summary |
| `vault.cassandra.list` | Duration of a LIST operation against the [Cassandra storage backend][cassandra-storage-backend] | ms | summary |
| `vault.cockroachdb.put` | Duration of a PUT operation against the [CockroachDB storage backend][cockroachdb-storage-backend] | ms | summary |
| `vault.cockroachdb.get` | Duration of a GET operation against the [CockroachDB storage backend][cockroachdb-storage-backend] | ms | summary |
| `vault.cockroachdb.delete` | Duration of a DELETE operation against the [CockroachDB storage backend][cockroachdb-storage-backend] | ms | summary |
| `vault.cockroachdb.list` | Duration of a LIST operation against the [CockroachDB storage backend][cockroachdb-storage-backend] | ms | summary |
| `vault.consul.put` | Duration of a PUT operation against the [Consul storage backend][consul-storage-backend] | ms | summary |
| `vault.consul.get` | Duration of a GET operation against the [Consul storage backend][consul-storage-backend] | ms | summary |
| `vault.consul.delete` | Duration of a DELETE operation against the [Consul storage backend][consul-storage-backend] | ms | summary |
| `vault.consul.list` | Duration of a LIST operation against the [Consul storage backend][consul-storage-backend] | ms | summary |
| `vault.couchdb.put` | Duration of a PUT operation against the [CouchDB storage backend][couchdb-storage-backend] | ms | summary |
| `vault.couchdb.get` | Duration of a GET operation against the [CouchDB storage backend][couchdb-storage-backend] | ms | summary |
| `vault.couchdb.delete` | Duration of a DELETE operation against the [CouchDB storage backend][couchdb-storage-backend] | ms | summary |
| `vault.couchdb.list` | Duration of a LIST operation against the [CouchDB storage backend][couchdb-storage-backend] | ms | summary |
| `vault.dynamodb.put` | Duration of a PUT operation against the [DynamoDB storage backend][dynamodb-storage-backend] | ms | summary |
| `vault.dynamodb.get` | Duration of a GET operation against the [DynamoDB storage backend][dynamodb-storage-backend] | ms | summary |
| `vault.dynamodb.delete` | Duration of a DELETE operation against the [DynamoDB storage backend][dynamodb-storage-backend] | ms | summary |
| `vault.dynamodb.list` | Duration of a LIST operation against the [DynamoDB storage backend][dynamodb-storage-backend] | ms | summary |
| `vault.etcd.put` | Duration of a PUT operation against the [etcd storage backend][etcd-storage-backend] | ms | summary |
| `vault.etcd.get` | Duration of a GET operation against the [etcd storage backend][etcd-storage-backend] | ms | summary |
| `vault.etcd.delete` | Duration of a DELETE operation against the [etcd storage backend][etcd-storage-backend] | ms | summary |
| `vault.etcd.list` | Duration of a LIST operation against the [etcd storage backend][etcd-storage-backend] | ms | summary |
| `vault.gcs.put` | Duration of a PUT operation against the [Google Cloud Storage storage backend][gcs-storage-backend] | ms | summary |
| `vault.gcs.get` | Duration of a GET operation against the [Google Cloud Storage storage backend][gcs-storage-backend] | ms | summary |
| `vault.gcs.delete` | Duration of a DELETE operation against the [Google Cloud Storage storage backend][gcs-storage-backend] | ms | summary |
| `vault.gcs.list` | Duration of a LIST operation against the [Google Cloud Storage storage backend][gcs-storage-backend] | ms | summary |
| `vault.gcs.lock.unlock` | Duration of an UNLOCK operation against the [Google Cloud Storage storage backend][gcs-storage-backend] in HA mode | ms | summary |
| `vault.gcs.lock.lock` | Duration of a LOCK operation against the [Google Cloud Storage storage backend][gcs-storage-backend] in HA mode | ms | summary |
| `vault.gcs.lock.value` | Duration of a VALUE operation against the [Google Cloud Storage storage backend][gcs-storage-backend] in HA mode | ms | summary |
| `vault.mssql.put` | Duration of a PUT operation against the [MS-SQL storage backend][mssql-storage-backend] | ms | summary |
| `vault.mssql.get` | Duration of a GET operation against the [MS-SQL storage backend][mssql-storage-backend] | ms | summary |
| `vault.mssql.delete` | Duration of a DELETE operation against the [MS-SQL storage backend][mssql-storage-backend] | ms | summary |
| `vault.mssql.list` | Duration of a LIST operation against the [MS-SQL storage backend][mssql-storage-backend] | ms | summary |
| `vault.mysql.put` | Duration of a PUT operation against the [MySQL storage backend][mysql-storage-backend] | ms | summary |
| `vault.mysql.get` | Duration of a GET operation against the [MySQL storage backend][mysql-storage-backend] | ms | summary |
| `vault.mysql.delete` | Duration of a DELETE operation against the [MySQL storage backend][mysql-storage-backend] | ms | summary |
| `vault.mysql.list` | Duration of a LIST operation against the [MySQL storage backend][mysql-storage-backend] | ms | summary |
| `vault.postgres.put` | Duration of a PUT operation against the [PostgreSQL storage backend][postgresql-storage-backend] | ms | summary |
| `vault.postgres.get` | Duration of a GET operation against the [PostgreSQL storage backend][postgresql-storage-backend] | ms | summary |
| `vault.postgres.delete` | Duration of a DELETE operation against the [PostgreSQL storage backend][postgresql-storage-backend] | ms | summary |
| `vault.postgres.list` | Duration of a LIST operation against the [PostgreSQL storage backend][postgresql-storage-backend] | ms | summary |
| `vault.s3.put` | Duration of a PUT operation against the [Amazon S3 storage backend][s3-storage-backend] | ms | summary |
| `vault.s3.get` | Duration of a GET operation against the [Amazon S3 storage backend][s3-storage-backend] | ms | summary |
| `vault.s3.delete` | Duration of a DELETE operation against the [Amazon S3 storage backend][s3-storage-backend] | ms | summary |
| `vault.s3.list` | Duration of a LIST operation against the [Amazon S3 storage backend][s3-storage-backend] | ms | summary |
| `vault.spanner.put` | Duration of a PUT operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] | ms | summary |
| `vault.spanner.get` | Duration of a GET operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] | ms | summary |
| `vault.spanner.delete` | Duration of a DELETE operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] | ms | summary |
| `vault.spanner.list` | Duration of a LIST operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] | ms | summary |
| `vault.spanner.lock.unlock` | Duration of an UNLOCK operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] in HA mode | ms | summary |
| `vault.spanner.lock.lock` | Duration of a LOCK operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] in HA mode | ms | summary |
| `vault.spanner.lock.value` | Duration of a VALUE operation against the [Google Cloud Spanner storage backend][gcs-storage-backend] in HA mode | ms | summary |
| `vault.swift.put` | Duration of a PUT operation against the [Swift storage backend][swift-storage-backend] | ms | summary |
| `vault.swift.get` | Duration of a GET operation against the [Swift storage backend][swift-storage-backend] | ms | summary |
| `vault.swift.delete` | Duration of a DELETE operation against the [Swift storage backend][swift-storage-backend] | ms | summary |
| `vault.swift.list` | Duration of a LIST operation against the [Swift storage backend][swift-storage-backend] | ms | summary |
| `vault.zookeeper.put` | Duration of a PUT operation against the [ZooKeeper storage backend][zookeeper-storage-backend] | ms | summary |
| `vault.zookeeper.get` | Duration of a GET operation against the [ZooKeeper storage backend][zookeeper-storage-backend] | ms | summary |
| `vault.zookeeper.delete` | Duration of a DELETE operation against the [ZooKeeper storage backend][zookeeper-storage-backend] | ms | summary |
| `vault.zookeeper.list` | Duration of a LIST operation against the [ZooKeeper storage backend][zookeeper-storage-backend] | ms | summary |

[secrets-engines]: /docs/secrets/index.html
[storage-backends]: /docs/configuration/storage/index.html
[telemetry-stanza]: /docs/configuration/telemetry.html
[cubbyhole-secrets-engine]: /docs/secrets/cubbyhole/index.html
[kv-secrets-engine]: /docs/secrets/kv/index.html
[ldap-auth-backend]: /docs/auth/ldap.html
[token-auth-backend]: /docs/auth/token.html
[azure-storage-backend]: /docs/configuration/storage/azure.html
[cassandra-storage-backend]: /docs/configuration/storage/cassandra.html
[cockroachdb-storage-backend]: /docs/configuration/storage/cockroachdb.html
[consul-storage-backend]: /docs/configuration/storage/consul.html
[couchdb-storage-backend]: /docs/configuration/storage/couchdb.html
[dynamodb-storage-backend]: /docs/configuration/storage/dynamodb.html
[etcd-storage-backend]: /docs/configuration/storage/etcd.html
[gcs-storage-backend]: /docs/configuration/storage/google-cloud-storage.html
[spanner-storage-backend]: /docs/configuration/storage/google-cloud-spanner.html
[mssql-storage-backend]: /docs/configuration/storage/mssql.html
[mysql-storage-backend]: /docs/configuration/storage/mysql.html
[postgresql-storage-backend]: /docs/configuration/storage/postgresql.html
[s3-storage-backend]: /docs/configuration/storage/s3.html
[swift-storage-backend]: /docs/configuration/storage/swift.html
[zookeeper-storage-backend]: /docs/configuration/storage/zookeeper.html
