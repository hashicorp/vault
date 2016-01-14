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
