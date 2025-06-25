# Benchmark modules

These benchmark modules are designed to benchmark a Vault cluster using either
raft storage or an external Consul cluster for storage. It supports creating
an external telemetry collector using a combination of Prometheus and Grafana
and supports collecting Vault, Consul, and node metrics.

When using this module in your scenario there are a few things to consider:
  - There are currently assumptions that Ubuntu 22.04 is the target runner for
    the k6 runner instance.
  - The Vault and Consul clusters are assumed to be three node clusters.
