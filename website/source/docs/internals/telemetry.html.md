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

Telemetry information can be streamed to both [statsite](http://github.com/armon/statsite)
as well as statsd based on providing the appropriate configuration options.

Below is sample output of a telemetry dump:

```text
[2014-01-29 10:56:50 -0800 PST][G] 'vault.runtime.num_goroutines': 19.000
[2014-01-29 10:56:50 -0800 PST][G] 'vault.runtime.alloc_bytes': 755960.000
[2014-01-29 10:56:50 -0800 PST][G] 'vault.runtime.malloc_count': 7550.000
[2014-01-29 10:56:50 -0800 PST][G] 'vault.runtime.free_count': 4387.000
[2014-01-29 10:56:50 -0800 PST][G] 'vault.runtime.heap_objects': 3163.000
[2014-01-29 10:56:50 -0800 PST][G] 'vault.runtime.total_gc_pause_ns': 1151002.000
[2014-01-29 10:56:50 -0800 PST][G] 'vault.runtime.total_gc_runs': 4.000
[2014-01-29 10:56:50 -0800 PST][C] 'vault.agent.ipc.accept': Count: 5 Sum: 5.000
[2014-01-29 10:56:50 -0800 PST][C] 'vault.agent.ipc.command': Count: 10 Sum: 10.000
[2014-01-29 10:56:50 -0800 PST][C] 'vault.serf.events': Count: 5 Sum: 5.000
[2014-01-29 10:56:50 -0800 PST][C] 'vault.serf.events.foo': Count: 4 Sum: 4.000
[2014-01-29 10:56:50 -0800 PST][C] 'vault.serf.events.baz': Count: 1 Sum: 1.000
[2014-01-29 10:56:50 -0800 PST][S] 'vault.memberlist.gossip': Count: 50 Min: 0.007 Mean: 0.020 Max: 0.041 Stddev: 0.007 Sum: 0.989
[2014-01-29 10:56:50 -0800 PST][S] 'vault.serf.queue.Intent': Count: 10 Sum: 0.000
[2014-01-29 10:56:50 -0800 PST][S] 'vault.serf.queue.Event': Count: 10 Min: 0.000 Mean: 2.500 Max: 5.000 Stddev: 2.121 Sum: 25.000
```

