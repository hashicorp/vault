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

## Internal Metrics

These metrics represent operational aspects of the running Vault instance.

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.audit.log_request`| Duration of time taken by all audit log requests across all audit log devices | ms | summary |
| `vault.audit.log_response`| Duration of time taken by audit log responses across all audit log devices | ms | summary |

Additionally, per audit log device metrics such as those for a specific backend like `file` will be present as:

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.audit.file.log_request`| Duration of time taken by audit log requests for the file based audit device mounted as `file` | ms | summary |
| `vault.audit.file.log_response`| Duration of time taken by audit log responses for the file based audit device mounted as `file` | ms | summary |

| Metric | Description | Unit | Type |
| :----- | :---------- | :--- | :--- |
| `vault.audit.log_request_failure` | Number of audit log request failures.  **NOTE**: This is a particularly important metric. Any non-zero value here indicates that there was a failure to make an audit log request to any of the configured audit log devices; **when Vault cannot log to any of the configured audit log devices it ceases all user operations**, and you should begin troubleshooting the audit log devices immediately if this metric continually increases. | failures | counter |
| `vault.audit.log_response_failure` | Number of audit log response failures. **NOTE**: This is a particularly important metric. Any non-zero value here indicates that there was a failure to receive a response to a request made to one of the configured audit log devices; **when Vault cannot log to any of the configured audit log devices it ceases all user operations**, and you should begin troubleshooting the audit log devices immediately if this metric continually increases. | failures | counter |
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

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`vault.expire.fetch-lease-times`</td>
    <td>Time taken to fetch lease times</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.fetch-lease-times-by-token`</td>
    <td>Time taken to fetch lease times by token</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.num_leases`</td>
    <td>Number of all leases which are eligible for eventual expiry</td>
    <td>leases</td>
    <td>gauge</td>
  </tr>

  <tr>
    <td>`vault.expire.revoke`</td>
    <td>Time taken to revoke a token</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.revoke-force`</td>
    <td>Time taken to forcibly revoke a token</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.revoke-prefix`</td>
    <td>Time taken to revoke tokens on a prefix</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.revoke-by-token`</td>
    <td>Time taken to revoke all secrets issued with a given token</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.renew`</td>
    <td>Time taken to renew a lease</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.renew-token`</td>
    <td>Time taken to renew a token which does not need to invoke a logical backend</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.expire.register`</td>
    <td>Time taken for register operations</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
</table>

Thes operations take a request and response with an associated lease and register a lease entry with lease ID

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`vault.expire.register-auth`</td>
    <td>Time taken for register authentication operations which create lease entries without lease ID</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.merkle_flushdirty`</td>
    <td>Time taken to flush any dirty pages to cold storage</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.merkle_savecheckpoint`</td>
    <td>Time taken to save the checkpoint</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.policy.get_policy`</td>
    <td>Time taken to get a policy</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.policy.list_policies`</td>
    <td>Time taken to list policies</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.policy.delete_policy`</td>
    <td>Time taken to delete a policy</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.policy.set_policy`</td>
    <td>Time taken to set a policy</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.token.create`</td>
    <td>The time taken to create a token</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.token.createAccessor`</td>
    <td>The time taken to create a token accessor</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.token.lookup`</td>
    <td>The time taken to look up a token</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.token.revoke`</td>
    <td>Time taken to revoke a token</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.token.revoke-tree`</td>
    <td>Time taken to revoke a token tree</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.token.store`</td>
    <td>Time taken to store an updated token entry without writing to the secondary index</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.wal_deletewals`</td>
    <td>Time taken to delete a Write Ahead Log (WAL)</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.wal_gc_deleted`</td>
    <td>Number of Write Ahead Logs (WAL) deleted during each garbage collection run</td>
    <td>WAL</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`vault.wal_gc_total`</td>
    <td>Total Number of Write Ahead Logs (WAL) on disk</td>
    <td>WAL</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`vault.wal_persistwals`</td>
    <td>Time taken to persist a Write Ahead Log (WAL)</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.wal_flushready`</td>
    <td>Time taken to flush a ready Write Ahead Log (WAL) to storage</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
</table>  

## Auth Methods Metrics

These metrics relate to supported authentication methods.

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`vault.rollback.attempt.auth-token-`</td>
    <td>Time taken to perform a rollback operation for the [token auth method][token-auth-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.rollback.attempt.auth-ldap-`</td>
    <td>Time taken to perform a rollback operation for the [LDAP auth method][ldap-auth-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.rollback.attempt.cubbyhole-`</td>
    <td>Time taken to perform a rollback operation for the [Cubbyhole secret backend][cubbyhole-secrets-engine]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.rollback.attempt.secret-`</td>
    <td>Time taken to perform a rollback operation for the [K/V secret backend][kv-secrets-engine]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.rollback.attempt.sys-`</td>
    <td>Time taken to perform a rollback operation for the system backend</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.route.rollback.auth-ldap-`</td>
    <td>Time taken to perform a route rollback operation for the [LDAP auth method][ldap-auth-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.route.rollback.auth-token-`</td>
    <td>Time taken to perform a route rollback operation for the [token auth method][token-auth-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.route.rollback.cubbyhole-`</td>
    <td>Time taken to perform a route rollback operation for the [Cubbyhole secret backend][cubbyhole-secrets-engine]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.route.rollback.secret-`</td>
    <td>Time taken to perform a route rollback operation for the [K/V secret backend][kv-secrets-engine]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.route.rollback.sys-`</td>
    <td>Time taken to perform a route rollback operation for the system backend</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
</table>

## Replication Metrics

These metrics relate to [Vault Enterprise Replication](https://www.vaultproject.io/docs/enterprise/replication/index.html).

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`logshipper.streamWALs.missing_guard`</td>
    <td>Number of incidences where the starting Merkle Tree index used to begin streaming WAL entries is not matched/found</td>
    <td>missing guards</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`logshipper.streamWALs.guard_found`</td>
    <td>Number of incidences where the starting Merkle Tree index used to begin streaming WAL entries is matched/found</td>
    <td>found guards</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`replication.fetchRemoteKeys`</td>
    <td>Time taken to fetch keys from a remote cluster participating in replication prior to Merkle Tree based delta generation</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`replication.merkleDiff`</td>
    <td>Time taken to perform a Merkle Tree based delta generation between the
        clusters participating in replication</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`replication.merkleSync`</td>
    <td>Time taken to perform a Merkle Tree based synchronization using the
        last delta generated between the clusters participating in replication</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`replication.merkle.commit_index`</td>
    <td>The last committed index in the Merkle Tree.</td>
    <td>sequence number</td>
    <td>gauge</td>
  </tr>

  <tr>
    <td>`replication.wal.last_wal`</td>
    <td>The index of the last WAL</td>
    <td>sequence number</td>
    <td>gauge</td>
  </tr>

  <tr>
    <td>`replication.wal.last_dr_wal`</td>
    <td>The index of the last DR WAL</td>
    <td>sequence number</td>
    <td>gauge</td>
  </tr>

  <tr>
    <td>`replication.wal.last_performance_wal`</td>
    <td>The index of the last Performance WAL</td>
    <td>sequence number</td>
    <td>gauge</td>
  </tr>

  <tr>
    <td>`replication.fsm.last_remote_wal`</td>
    <td>The index of the last remote WAL</td>
    <td>sequence number</td>
    <td>gauge</td>
  </tr>

  <tr>
    <td>`replication.rpc.server.auth_request`</td>
    <td>Duration of time taken by auth request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.bootstrap_request`</td>
    <td>Duration of time taken by bootstrap request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.conflicting_pages_request`</td>
    <td>Duration of time taken by conflicting pages request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.echo`</td>
    <td>Duration of time taken by echo</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.forwarding_request`</td>
    <td>Duration of time taken by forwarding request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.guard_hash_request`</td>
    <td>Duration of time taken by guard hash request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.persist_alias_request`</td>
    <td>Duration of time taken by persist alias request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.persist_persona_request`</td>
    <td>Duration of time taken by persist persona request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.stream_wals_request`</td>
    <td>Duration of time taken by stream wals request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.sub_page_hashes_request`</td>
    <td>Duration of time taken by sub page hashes request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.sync_counter_request`</td>
    <td>Duration of time taken by sync counter request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.server.upsert_group_request`</td>
    <td>Duration of time taken by upsert group request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.conflicting_pages`</td>
    <td>Duration of time taken by client conflicting pages request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.fetch_keys`</td>
    <td>Duration of time taken by client fetch keys request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.forward`</td>
    <td>Duration of time taken by client forward request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.guard_hash`</td>
    <td>Duration of time taken by client guard hash request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.persist_alias`</td>
    <td>Duration of time taken by</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.register_auth`</td>
    <td>Duration of time taken by client register auth request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.register_lease`</td>
    <td>Duration of time taken by client register lease request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.stream_wals`</td>
    <td>Duration of time taken by client s</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.sub_page_hashes`</td>
    <td>Duration of time taken by client sub page hashes request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.sync_counter`</td>
    <td>Duration of time taken by client sync counter request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.upsert_group`</td>
    <td>Duration of time taken by client upstert group request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.client.wrap_in_cubbyhole`</td>
    <td>Duration of time taken by client wrap in cubbyhole request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.dr.server.echo`</td>
    <td>Duration of time taken by DR echo request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.dr.server.fetch_keys_request`</td>
    <td>Duration of time taken by DR fetch keys request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.standby.server.echo`</td>
    <td>Duration of time taken by standby echo request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.standby.server.register_auth_request`</td>
    <td>Duration of time taken by standby register auth request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.standby.server.register_lease_request`</td>
    <td>Duration of time taken by standby register lease request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
  <tr>
    <td>`replication.rpc.standby.server.wrap_token_request`</td>
    <td>Duration of time taken by standby wrap token request</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

</table>

## Secrets Engines Metrics

These metrics relate to the supported [secrets engines][secrets-engines].

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`database.Initialize`</td>
    <td>Time taken to initialize a database secret engine across all database secrets engines</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.Initialize`</td>
    <td>Time taken to initialize a database secret engine for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Initialize`</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.Initialize.error`</td>
    <td>Number of database secrets engine initialization operation errors across all database secrets engines</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.Initialize.error`</td>
    <td>Number of database secrets engine initialization operation errors for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Initialize.error`</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.Close`</td>
    <td>Time taken to close a database secret engine across all database secrets engines</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.Close`</td>
    <td>Time taken to close a database secret engine for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Close`</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.Close.error`</td>
    <td>Number of database secrets engine close operation errors across all database secrets engines</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.Close.error`</td>
    <td>Number of database secrets engine close operation errors for the named database secrets engine `<name>`, for example: `database.postgresql-prod.Close.error`</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.CreateUser`</td>
    <td>Time taken to create a user across all database secrets engines</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.CreateUser`</td>
    <td>Time taken to create a user for the named database secrets engine `<name>`</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.CreateUser.error`</td>
    <td>Number of user creation operation errors across all database secrets engines</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.CreateUser.error`</td>
    <td>Number of user creation operation errors for the named database secrets engine `<name>`, for example: `database.postgresql-prod.CreateUser.error`</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.RenewUser`</td>
    <td>Time taken to renew a user across all database secrets engines</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.RenewUser`</td>
    <td>Time taken to renew a user for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RenewUser`</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.RenewUser.error`</td>
    <td>Number of user renewal operation errors across all database secrets engines</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.RenewUser.error`</td>
    <td>Number of user renewal operations for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RenewUser.error`</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.RevokeUser`</td>
    <td>Time taken to revoke a user across all database secrets engines</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.RevokeUser`</td>
    <td>Time taken to revoke a user for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RevokeUser`</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`database.RevokeUser.error`</td>
    <td>Number of user revocation operation errors across all database secrets engines</td>
    <td>errors</td>
    <td>counter</td>
  </tr>

  <tr>
    <td>`database.&lt;name&gt;.RevokeUser.error`</td>
    <td>Number of user revocation operations for the named database secrets engine `<name>`, for example: `database.postgresql-prod.RevokeUser.error`</td>
    <td>errors</td>
    <td>counter</td>
  </tr>
</table>

## Storage Backend Metrics

These metrics relate to the supported [storage backends][storage-backends].

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`vault.azure.put`</td>
    <td>Duration of a PUT operation against the [Azure storage backend][azure-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.azure.get`</td>
    <td>Duration of a GET operation against the [Azure storage backend][azure-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.azure.delete`</td>
    <td>Duration of a DELETE operation against the [Azure storage backend][azure-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.azure.list`</td>
    <td>Duration of a LIST operation against the [Azure storage backend][azure-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cassandra.put`</td>
    <td>Duration of a PUT operation against the [Cassandra storage backend][cassandra-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cassandra.get`</td>
    <td>Duration of a GET operation against the [Cassandra storage backend][cassandra-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cassandra.delete`</td>
    <td>Duration of a DELETE operation against the [Cassandra storage backend][cassandra-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cassandra.list`</td>
    <td>Duration of a LIST operation against the [Cassandra storage backend][cassandra-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cockroachdb.put`</td>
    <td>Duration of a PUT operation against the [CockroachDB storage backend][cockroachdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cockroachdb.get`</td>
    <td>Duration of a GET operation against the [CockroachDB storage backend][cockroachdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cockroachdb.delete`</td>
    <td>Duration of a DELETE operation against the [CockroachDB storage backend][cockroachdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.cockroachdb.list`</td>
    <td>Duration of a LIST operation against the [CockroachDB storage backend][cockroachdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.consul.put`</td>
    <td>Duration of a PUT operation against the [Consul storage backend][consul-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.consul.get`</td>
    <td>Duration of a GET operation against the [Consul storage backend][consul-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.consul.delete`</td>
    <td>Duration of a DELETE operation against the [Consul storage backend][consul-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.consul.list`</td>
    <td>Duration of a LIST operation against the [Consul storage backend][consul-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.couchdb.put`</td>
    <td>Duration of a PUT operation against the [CouchDB storage backend][couchdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.couchdb.get`</td>
    <td>Duration of a GET operation against the [CouchDB storage backend][couchdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.couchdb.delete`</td>
    <td>Duration of a DELETE operation against the [CouchDB storage backend][couchdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.couchdb.list`</td>
    <td>Duration of a LIST operation against the [CouchDB storage backend][couchdb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.dynamodb.put`</td>
    <td>Duration of a PUT operation against the [DynamoDB storage backend][dynamodb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.dynamodb.get`</td>
    <td>Duration of a GET operation against the [DynamoDB storage backend][dynamodb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.dynamodb.delete`</td>
    <td>Duration of a DELETE operation against the [DynamoDB storage backend][dynamodb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.dynamodb.list`</td>
    <td>Duration of a LIST operation against the [DynamoDB storage backend][dynamodb-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.etcd.put`</td>
    <td>Duration of a PUT operation against the [etcd storage backend][etcd-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.etcd.get`</td>
    <td>Duration of a GET operation against the [etcd storage backend][etcd-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.etcd.delete`</td>
    <td>Duration of a DELETE operation against the [etcd storage backend][etcd-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.etcd.list`</td>
    <td>Duration of a LIST operation against the [etcd storage backend][etcd-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.gcs.put`</td>
    <td>Duration of a PUT operation against the [Google Cloud Storage storage backend][gcs-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.gcs.get`</td>
    <td>Duration of a GET operation against the [Google Cloud Storage storage backend][gcs-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.gcs.delete`</td>
    <td>Duration of a DELETE operation against the [Google Cloud Storage storage backend][gcs-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.gcs.list`</td>
    <td>Duration of a LIST operation against the [Google Cloud Storage storage backend][gcs-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.gcs.lock.unlock`</td>
    <td>Duration of an UNLOCK operation against the [Google Cloud Storage storage backend][gcs-storage-backend] in HA mode</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.gcs.lock.lock`</td>
    <td>Duration of a LOCK operation against the [Google Cloud Storage storage backend][gcs-storage-backend] in HA mode</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.gcs.lock.value`</td>
    <td>Duration of a VALUE operation against the [Google Cloud Storage storage backend][gcs-storage-backend] in HA mode</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mssql.put`</td>
    <td>Duration of a PUT operation against the [MS-SQL storage backend][mssql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mssql.get`</td>
    <td>Duration of a GET operation against the [MS-SQL storage backend][mssql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mssql.delete`</td>
    <td>Duration of a DELETE operation against the [MS-SQL storage backend][mssql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mssql.list`</td>
    <td>Duration of a LIST operation against the [MS-SQL storage backend][mssql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mysql.put`</td>
    <td>Duration of a PUT operation against the [MySQL storage backend][mysql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mysql.get`</td>
    <td>Duration of a GET operation against the [MySQL storage backend][mysql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mysql.delete`</td>
    <td>Duration of a DELETE operation against the [MySQL storage backend][mysql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.mysql.list`</td>
    <td>Duration of a LIST operation against the [MySQL storage backend][mysql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.postgres.put`</td>
    <td>Duration of a PUT operation against the [PostgreSQL storage backend][postgresql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.postgres.get`</td>
    <td>Duration of a GET operation against the [PostgreSQL storage backend][postgresql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.postgres.delete`</td>
    <td>Duration of a DELETE operation against the [PostgreSQL storage backend][postgresql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.postgres.list`</td>
    <td>Duration of a LIST operation against the [PostgreSQL storage backend][postgresql-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.s3.put`</td>
    <td>Duration of a PUT operation against the [Amazon S3 storage backend][s3-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.s3.get`</td>
    <td>Duration of a GET operation against the [Amazon S3 storage backend][s3-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.s3.delete`</td>
    <td>Duration of a DELETE operation against the [Amazon S3 storage backend][s3-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.s3.list`</td>
    <td>Duration of a LIST operation against the [Amazon S3 storage backend][s3-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.spanner.put`</td>
    <td>Duration of a PUT operation against the [Google Cloud Spanner storage backend][spanner-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.spanner.get`</td>
    <td>Duration of a GET operation against the [Google Cloud Spanner storage backend][spanner-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.spanner.delete`</td>
    <td>Duration of a DELETE operation against the [Google Cloud Spanner storage backend][spanner-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.spanner.list`</td>
    <td>Duration of a LIST operation against the [Google Cloud Spanner storage backend][spanner-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.spanner.lock.unlock`</td>
    <td>Duration of an UNLOCK operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] in HA mode</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.spanner.lock.lock`</td>
    <td>Duration of a LOCK operation against the [Google Cloud Spanner storage backend][spanner-storage-backend] in HA mode</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.spanner.lock.value`</td>
    <td>Duration of a VALUE operation against the [Google Cloud Spanner storage backend][gcs-storage-backend] in HA mode</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.swift.put`</td>
    <td>Duration of a PUT operation against the [Swift storage backend][swift-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.swift.get`</td>
    <td>Duration of a GET operation against the [Swift storage backend][swift-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.swift.delete`</td>
    <td>Duration of a DELETE operation against the [Swift storage backend][swift-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.swift.list`</td>
    <td>Duration of a LIST operation against the [Swift storage backend][swift-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.zookeeper.put`</td>
    <td>Duration of a PUT operation against the [ZooKeeper storage backend][zookeeper-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.zookeeper.get`</td>
    <td>Duration of a GET operation against the [ZooKeeper storage backend][zookeeper-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.zookeeper.delete`</td>
    <td>Duration of a DELETE operation against the [ZooKeeper storage backend][zookeeper-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>

  <tr>
    <td>`vault.zookeeper.list`</td>
    <td>Duration of a LIST operation against the [ZooKeeper storage backend][zookeeper-storage-backend]</td>
    <td>ms</td>
    <td>summary</td>
  </tr>
</table>

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
