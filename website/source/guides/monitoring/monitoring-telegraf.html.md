---
layout: "guides"
page_title: "Monitoring with Telegraf - Guides"
sidebar_current: "guides-monitoring-telegraf"
description: |-
  This guide discusses the concepts necessary to monitor Vault
---

# Monitoring Vault with Telegraf

Vault has a number of metrics avaliable through through http endpoints. Measuring these metrics is useful for both health checking and debugging of cluster issues.

The core points will be the same if other tools are used, but the names of the metrics may be different.

## Configuring Telegraf

# Installing Telegraf

Installing Telegraf is straightforward on most Linux distributions. We recomend following the [official Telegraf installation documentation][telegraf-install].

# Configuring Telegraf

Besides acting as a statsd agent, Telegraf can collect additional metrics of its own. Telegraf itself ships
with a wide range of [input plugins][telegraf-input-plugins] to collect data from lots of sources.
We're going to enable some of the most common ones to monitor CPU, memory, disk I/O, networking, and process status, as these are useful for debugging Consul cluster issues.

The `telegraf.conf` file starts with global options:

```ini
[agent]
  interval = "10s"
  flush_interval = "10s"
  omit_hostname = false
```

We set the default collection interval to 10 seconds and ask Telegraf to include a `host` tag in each metric.

As mentioned above, Telegraf also allows you to set additional tags on the metrics that pass through it. In this
case, we are adding tags for the server role and datacenter. We can then use these tags in Grafana to filter
queries (for example, to create a dashboard showing only servers with the `consul-server` role, or only servers
in the `us-east-1` datacenter).

```ini
[global_tags]
  role = "consul-server"
  datacenter = "us-east-1"
```

Next, we set up a statsd listener on UDP port 8125, with instructions to calculate percentile metrics and to
parse DogStatsd-compatible tags, when they're sent:

```ini
[[inputs.statsd]]
  protocol = "udp"
  service_address = ":8125"
  delete_gauges = true
  delete_counters = true
  delete_sets = true
  delete_timings = true
  percentiles = [90]
  metric_separator = "_"
  parse_data_dog_tags = true
  allowed_pending_messages = 10000
  percentile_limit = 1000
```

The full reference to all the available statsd-related options in Telegraf is [here][telegraf-statsd-input].

Now we can configure inputs for things like CPU, memory, network I/O, and disk I/O. Most of them don't require any configuration, but make sure the `interfaces` list in `inputs.net` matches the interface names you see in `ifconfig`.

```ini
[[inputs.cpu]]
  percpu = true
  totalcpu = true
  collect_cpu_time = false

[[inputs.disk]]
  # mount_points = ["/"]
  # ignore_fs = ["tmpfs", "devtmpfs"]

[[inputs.diskio]]
  # devices = ["sda", "sdb"]
  # skip_serial_number = false

[[inputs.kernel]]
  # no configuration

[[inputs.linux_sysctl_fs]]
  # no configuration

[[inputs.mem]]
  # no configuration

[[inputs.net]]
  interfaces = ["enp0s*"]

[[inputs.netstat]]
  # no configuration

[[inputs.processes]]
  # no configuration

[[inputs.swap]]
  # no configuration

[[inputs.system]]
  # no configuration
```

Another useful plugin is the [procstat][telegraf-procstat-input] plugin, which reports metrics for processes you select:

```ini
[[inputs.procstat]]
  pattern = "(vault)"
```

[telegraf-install]: https://docs.influxdata.com/telegraf/v1.6/introduction/installation/
[telegraf-consul-input]: https://github.com/influxdata/telegraf/tree/release-1.6/plugins/inputs/consul
[telegraf-statsd-input]: https://github.com/influxdata/telegraf/tree/release-1.6/plugins/inputs/statsd
[telegraf-procstat-input]: https://github.com/influxdata/telegraf/tree/release-1.6/plugins/inputs/procstat
[telegraf-input-plugins]: https://docs.influxdata.com/telegraf/v1.6/plugins/inputs/

## Telegraf Configuration for Consul

Asking Consul to send telemetry to Telegraf is as simple as adding a `telemetry` section to your agent configuration:

```json
{
  "telemetry": {
    "dogstatsd_addr": "localhost:8125",
    "disable_hostname": true
  }
}
```

As you can see, we only need to specify two options. The `dogstatsd_addr` specifies the hostname and port of the
statsd daemon.

Note that we specify DogStatsd format instead of plain statsd, which tells Consul to send [tags][tagging]
with each metric. Tags can be used by Grafana to filter data on your dashboards (for example, displaying only
the data for which `role=consul-server`. Telegraf is compatible with the DogStatsd format and allows us to add
our own tags too, as you'll see later.

The second option tells Consul not to insert the hostname in the names of the metrics it sends to statsd, since the hostnames will be sent as tags. Without this option, the single metric `consul.raft.apply` would become multiple metrics:

        vault.server1.core.handle_request
        vault.server2.core.handle_request
        vault.server3.core.handle_request

If you are using a different agent (e.g. Circonus, Statsite, or plain statsd), you can find the configuration reference [here][consul-telemetry-config].

[tagging]: https://docs.datadoghq.com/getting_started/tagging/
[telegraf-input-plugins]: https://docs.influxdata.com/telegraf/v1.6/plugins/inputs/

## Important Metric Aggregates and Alerting

# Request processing

| Metric Name | Description |
| :---------- | :---------- |
| `vault.core.handle_request` | Duration of requests handled by the Vault core. |

**Why it's important:** This is the key measure of Vault's response time.

**What to look for:** Changes to the `count` or `mean` fields that exceed 50% of baseline values, or
more than 3 standard deviations above baseline.

# Consul response time

| Metric Name | Description |
| :---------- | :---------- |
| `vault.consul.get` | Count and duration of `GET` operations against the Consul storage backend. |
| `vault.consul.put` | Count and duration of `PUT` operations against the Consul storage backend. |
| `vault.consul.list` | Count and duration of `LIST` operations against the Consul storage backend. |
| `vault.consul.delete` | Count and duration of `DELETE` operations against the Consul storage backend. |

**Why they're important:** These metrics indicate how long it takes for Consul to handle requests from Vault.

**What to look for:** Large deltas in the `count`, `upper`, or `90_percentile` fields.

# Write-ahead log processing

| Metric Name | Description |
| :---------- | :---------- |
| `vault.wal.persistWALs` | Amount of time required to persist the Vault write-ahead logs (WAL) to the Consul backend. |
| `vault.wal.flushReady` | Amount of time required to flush the Vault write-ahead logs (WAL) to the persist queue. |

**Why they're important:** The Vault write-ahead logs (WALs) are used to replicate Vault between clusters.
Surprisingly, the WAL's are kept even if replication is not currently enabled. The WAL is purged every few seconds by
a garbage collector. But if Vault is under heavy load, the WAL may start to grow, putting pressure on Consul.

**What to look for:** If `flushReady` is over 500ms, or if `persistWALs` is over 1000ms.

# Leadership changes

| Metric Name | Description |
| :---------- | :---------- |
| `vault.core.leadership_lost` | Total duration of cluster leadership losses in a highly-available cluster. |

**Why it's important:** There should not be a leadership change unless the leader crashes or becomes otherwise
unavailable. While the other servers elect a leader, Vault is unable to process any requests.

**What to monitor:** Any value greater than 0 should cause an alert condition.

# Seal status

| Metric Name | Description |
| :---------- | :---------- |
| `consul_health_checks[check_name="Vault Sealed Status"].passing` | Value of 1 indicates Vault is unsealed; 0 means sealed. |

**Why they're important:** By default, Vault is sealed on startup, so if this value changes to 0 during the day,
Vault has restarted for some reason. And until it's unsealed, it won't answer requests from clients.

**What to look for:** A value of 0 being reported by any host.

**NOTE:** This metric is actually reported by the [Consul plugin to Telegraf][telegraf_plugin].

# Memory usage

| Metric Name | Description |
| :---------- | :---------- |
| `vault.runtime.alloc_bytes` | This measures the number of bytes allocated by the Vault process.  |
| `vault.runtime.sys_bytes`   | This is the total number of bytes of memory obtained from the OS.  |
| `mem.total_bytes`           | Total amount of physical memory (RAM) available on the server.     |
| `mem.used_percent`           | Percentage of physical memory in use. |
| `swap.used_percent`          | Percentage of swap space in use. |

**Why they're important:** Vault doesn't need as much memory as Consul, but if it does run out, it too will crash.
You should also monitor total available RAM to make sure some RAM is available for other processes, and swap usage
should remain at 0% for best performance.

**What to look for:** If `sys_bytes` exceeds 90% of `total_bytes`, if `mem.used_percent` is over 90%, or if
`swap.used_percent` is greater than 0.

# Garbage collection

| Metric Name | Description |
| :---------- | :---------- |
| `vault.runtime.total_gc_pause_ns` | Number of nanoseconds consumed by stop-the-world garbage collection (GC) pauses since Vault started. |

**Why it's important:** As mentioned above, GC pause is a "stop-the-world" event, meaning that all runtime threads are
blocked until GC completes. Normally these pauses last only a few nanoseconds. But if memory usage is high, the Go
runtime may GC so frequently that it starts to slow down Vault.

**What to look for:** Warning if `total_gc_pause_ns` exceeds 2 seconds/minute, critical if it exceeds 5 seconds/minute.

**Additional notes:** `total_gc_pause_ns` is a cumulative counter, so in order to calculate rates (such as GC/minute),
you will need to apply a function such as [non_negative_difference][].

# File descriptors

| Metric Name | Description |
| :---------- | :---------- |
| `linux_sysctl_fs.file-nr` | Number of file handles being used across all processes on the host. |
| `linux_sysctl_fs.file-max` | Total number of available file handles. |

**Why it's important:** Practically anything Vault does -- receiving a connection from another host, sending data
between servers, writing snapshots to disk -- requires a file descriptor handle. If Vault runs out of handles, it
will stop accepting connections.

By default, process and kernel limits are fairly conservative. You will want to increase these beyond the defaults.

**What to look for:** If `file-nr` exceeds 80% of `file-max`.

# CPU usage

| Metric Name | Description |
| :---------- | :---------- |
| `cpu.user_cpu` | Percentage of CPU being used by user processes (such as Vault or Consul). |
| `cpu.iowait_cpu` | Percentage of CPU time spent waiting for I/O tasks to complete. |

**Why they're important:** Encryption can place a heavy demand on CPU. If the CPU is too busy, Vault may have
trouble keeping up with the incoming request load. You may also want to monitor each CPU individually to make sure
requests are evenly balanced across all CPUs.

**What to look for:** if `cpu.iowait_cpu` greater than 10%.

# Network activity

| Metric Name | Description |
| :---------- | :---------- |
| `net.bytes_recv` | Bytes received on each network interface. |
| `net.bytes_sent` | Bytes transmitted on each network interface. |

**Why they're important:** A sudden spike in network traffic to Vault might be the result of a misconfigured
client causing too many requests, or additional load you didn't plan for.

**What to look for:**
Sudden large changes to the `net` metrics (greater than 50% deviation from baseline).

**NOTE:** The `net` metrics are counters, so in order to calculate rates (such as bytes/second),
you will need to apply a function such as [non_negative_difference][].

# Disk activity

| Metric Name | Description |
| :---------- | :---------- |
| `diskio.read_bytes` | Bytes read from each block device. |
| `diskio.write_bytes` | Bytes written to each block device. |

**Why they're important:** Vault generally doesn't require too much disk I/O, so a sudden change in disk activity
could mean that debug/trace logging has accidentally been enabled in production, which can impact performance.
Too much disk I/O can cause the rest of the system to slow down or become unavailable, as the kernel spends all
its time waiting for I/O to complete.

**What to look for:** Sudden large changes to the `diskio` metrics (greater than 50% deviation from baseline,
or more than 3 standard deviations from baseline).

**NOTE:** The `diskio` metrics are counters, so in order to calculate rates (such as bytes/second),
you will need to apply a function such as [non_negative_difference][].


[non_negative_difference]: https://docs.influxdata.com/influxdb/v1.5/query_language/functions/#non-negative-difference
[consul_faq_fds]: https://www.consul.io/docs/faq.html#q-does-consul-require-certain-user-process-resource-limits-
[telegraf_plugin]: https://github.com/influxdata/telegraf/tree/master/plugins/inputs/consul
