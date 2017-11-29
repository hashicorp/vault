---
layout: "docs"
page_title: "Telemetry"
sidebar_current: "docs-internals-telemetry"
description: |-
  Learn about the telemetry data available in Vault.
---

# Telemetry

The Vault agent collects various runtime metrics about the performance of
different libraries and subsystems. These metrics are aggregated on a ten
second interval and are retained for one minute.

To view this data, you must send a signal to the Vault process: on Unix,
this is `USR1` while on Windows it is `BREAK`. Once Vault receives the signal,
it will dump the current telemetry information to the agent's `stderr`.

This telemetry information can be used for debugging or otherwise
getting a better view of what Vault is doing.

Telemetry information can be streamed to both [statsite](https://github.com/armon/statsite)
as well as statsd based on providing the appropriate configuration options.

Below is sample output of a telemetry dump:

```text
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.num_goroutines': 12.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.free_count': 11882.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.total_gc_runs': 9.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.expire.num_leases': 1.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.alloc_bytes': 502992.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.sys_bytes': 3999992.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.malloc_count': 17315.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.heap_objects': 5433.000
[2015-04-20 12:24:30 -0700 PDT][G] 'vault.runtime.total_gc_pause_ns': 3794124.000
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.audit.log_response': Count: 2 Min: 0.001 Mean: 0.001 Max: 0.001 Stddev: 0.000 Sum: 0.002
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.route.read.secret-': Count: 1 Sum: 0.036
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.barrier.get': Count: 3 Min: 0.004 Mean: 0.021 Max: 0.050 Stddev: 0.025 Sum: 0.064
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.token.lookup': Count: 2 Min: 0.040 Mean: 0.074 Max: 0.108 Stddev: 0.048 Sum: 0.148
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.policy.get_policy': Count: 2 Min: 0.003 Mean: 0.004 Max: 0.005 Stddev: 0.001 Sum: 0.009
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.core.check_token': Count: 2 Min: 0.053 Mean: 0.087 Max: 0.121 Stddev: 0.048 Sum: 0.174
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.audit.log_request': Count: 2 Min: 0.001 Mean: 0.001 Max: 0.001 Stddev: 0.000 Sum: 0.002
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.barrier.put': Count: 3 Min: 0.004 Mean: 0.010 Max: 0.019 Stddev: 0.008 Sum: 0.029
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.route.write.secret-': Count: 1 Sum: 0.035
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.core.handle_request': Count: 2 Min: 0.097 Mean: 0.228 Max: 0.359 Stddev: 0.186 Sum: 0.457
[2015-04-20 12:24:30 -0700 PDT][S] 'vault.expire.register': Count: 1 Sum: 0.18
```

You'll note that log entries are prefixed with the metric type as follows:

- `[C]` is a counter
- `[G]` is a gauge
- `[S]` is a summary

## Key Metrics

The following tables described the different Vault metrics. The metrics interval can be assumed to be 10 seconds when retrieving metrics using the above described signals.

### Internal Metrics

These metrics represent operational aspects of the running Vault instance.

| Metric           | Description                       | Unit | Type |
| ---------------- | ----------------------------------| ---- | ---- |
|`vault.audit.log_request`| This measures the number of audit log requests | Number of requests | Summary |
|`vault.audit.log_response`| This measures the number of audit log responses | Number of responses | Summary |
|`vault.audit.log_request_failure` | The number of audit log request failures | Number of failures | Counter |
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
`vault.token.create`| This measures the time taken to create a token | Number of operations | Gauge |
`vault.token.createAccessor`| This measures the number of Token ID identifier operations | Number of operations | Gauge |
`vault.token.lookup`| This measures the number of token lookups | Number of lookups | Counter |
`vault.token.revoke`| This measures the number of token revocation operations | Number of operations | Gauge |
`vault.token.revoke-tree`| This measures the number of revoke tree operations | Number of operations | Gauge |
`vault.token.store`| This measures the number of operations to store an updated token entry without writing to the secondary index | Number of operations | Gauge |

### Authentication Backend Metrics

These metrics relate to supported authentication backends.

| Metric           | Description                       | Unit | Type |
| ---------------- | ----------------------------------| ---- | ---- |
| `vault.rollback.attempt.auth-token-` | This measures the number of rollback operations attempted for authentication tokens backend | Number of operations | Summary |
| `vault.rollback.attempt.auth-ldap-` | This measures the number of rollback operations attempted for the LDAP authentication backend | Number of operations | Summary |
| `vault.rollback.attempt.cubbyhole-` | This measures the number of rollback operations attempted for the cubbyhole authentication backend | Number of operations | Summary |
| `vault.rollback.attempt.secret-` | This measures the number of rollback operations attempted for the kv secret backend | Number of operations | Summary |
| `vault.rollback.attempt.sys-` | This measures the number of rollback operations attempted for the sys backend | Number of operations | Summary |
| `vault.route.rollback.auth-ldap-` | This measures the number of rollback operations for the LDAP authentication backend | Number of operations | Summary |
| `vault.route.rollback.auth-token-` | This measures the number of rollback operations for the authentication tokens backend | Number of operations | Summary |
| `vault.route.rollback.cubbyhole-` | This measures the number of rollback operations for the cubbyhole authentication backend | Number of operations | Summary |
| `vault.route.rollback.secret-` | This measures the number of rollback operations for the kv secret backend | Number of operations | Summary | 
| `vault.route.rollback.sys-` | This measures the number of rollback operations for the sys backend | Number of operations | Summary |

### Storage Backend Metrics

These metrics relate to supported storage backends.

| Metric           | Description                       | Unit | Type |
| ---------------- | ----------------------------------| ---- | ---- |
|`vault.azure.put` | This measures the number of put operations against the Azure storage backend | Number of operations | Gauge |
|`vault.azure.get` | This measures the number of get operations against the Azure storage backend | Number of operations | Gauge |
|`vault.azure.delete` | This measures the number of delete operations against the Azure storage backend | Number of operations | Gauge |
|`vault.azure.list` | This measures the number of list operations against the Azure storage backend | Number of operations | Gauge |
|`vault.consul.put` | This measures the number of put operations against the Consul storage backend | Number of operations | Gauge |
|`vault.consul.get` | This measures the number of get operations against the Consul storage backend | Number of operations | Gauge |
|`vault.consul.delete` | This measures the number of delete operations against the Consul storage backend | Number of operations | Gauge |
|`vault.consul.list` | This measures the number of list operations against the Consul storage backend | Number of operations | Gauge |
|`vault.dynamodb.put` | This measures the number of put operations against the DynamoDB storage backend | Number of operations | Gauge |
|`vault.dynamodb.get` | This measures the number of get operations against the DynamoDB storage backend | Number of operations | Gauge |
|`vault.dynamodb.delete` | This measures the number of delete operations against the DynamoDB storage backend | Number of operations | Gauge |
|`vault.dynamodb.list` | This measures the number of list operations against the DynamoDB storage backend | Number of operations | Gauge |
|`vault.etcd.put` | This measures the number of put operations against the etcd storage backend | Number of operations | Gauge |
|`vault.etcd.get` | This measures the number of get operations against the etcd storage backend | Number of operations | Gauge |
|`vault.etcd.delete` | This measures the number of delete operations against the etcd storage backend | Number of operations | Gauge |
|`vault.etcd.list` | This measures the number of list operations against the etcd storage backend | Number of operations | Gauge |
|`vault.gcs.put` | This measures the number of put operations against the Google Cloud Storage backend | Number of operations | Gauge |
|`vault.gcs.get` | This measures the number of get operations against the Google Cloud Storage backend | Number of operations | Gauge |
|`vault.gcs.delete` | This measures the number of delete operations against the Google Cloud Storage backend | Number of operations | Gauge |
|`vault.gcs.list` | This measures the number of list operations against the Google Cloud Storage backend | Number of operations | Gauge |
|`vault.mysql.put` | This measures the number of put operations against the MySQL backend | Number of operations | Gauge |
|`vault.mysql.get` | This measures the number of get operations against the MySQL backend | Number of operations | Gauge |
|`vault.mysql.delete` | This measures the number of delete operations against the MySQL backend | Number of operations | Gauge |
|`vault.mysql.list` | This measures the number of list operations against the MySQL backend | Number of operations | Gauge |
|`vault.postgres.put` | This measures the number of put operations against the PostgreSQL backend | Number of operations | Gauge |
|`vault.postgres.get` | This measures the number of get operations against the PostgreSQL backend | Number of operations | Gauge |
|`vault.postgres.delete` | This measures the number of delete operations against the PostgreSQL backend | Number of operations | Gauge |
|`vault.postgres.list` | This measures the number of list operations against the PostgreSQL backend | Number of operations | Gauge |
|`vault.s3.put` | This measures the number of put operations against the Amazon S3 backend | Number of operations | Gauge |
|`vault.s3.get` | This measures the number of get operations against the Amazon S3 backend | Number of operations | Gauge |
|`vault.s3.delete` | This measures the number of delete operations against the Amazon S3 backend | Number of operations | Gauge |
|`vault.s3.list` | This measures the number of list operations against the Amazon S3 backend | Number of operations | Gauge |
|`vault.swift.put` | This measures the number of put operations against the OpenStack Swift backend | Number of operations | Gauge |
|`vault.swift.get` | This measures the number of get operations against the OpenStack Swift backend | Number of operations | Gauge |
|`vault.swift.delete` | This measures the number of delete operations against the OpenStack Swift backend | Number of operations | Gauge |
|`vault.swift.list` | This measures the number of list operations against the OpenStack Swift backend | Number of operations | Gauge |
|`vault.zookeeper.put` | This measures the number of put operations against the ZooKeeper backend | Number of operations | Gauge |
|`vault.zookeeper.get` | This measures the number of get operations against the ZooKeeper backend | Number of operations | Gauge |
|`vault.zookeeper.delete` | This measures the number of delete operations against the ZooKeeper backend | Number of operations | Gauge |
|`vault.zookeeper.list` | This measures the number of list operations against the ZooKeeper backend | Number of operations | Gauge |
