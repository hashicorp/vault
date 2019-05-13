---
layout: "guides"
page_title: "Vault Cluster Monitoring - Guides"
sidebar_title: "Vault Cluster Monitoring"
sidebar_current: "guides-operations-monitoring"
description: |-
  Learn how to set up and manage Vault Enterprise Performance Replication.
---

# Vault Cluster Monitoring

~> [Download **Vault Cluster Monitoring Guide**](https://s3-us-west-2.amazonaws.com/hashicorp-education/whitepapers/Vault/Vault-Consul-Monitoring-Guide.pdf)


This _Vault Cluster Monitoring Guide_ demonstrates monitoring of a Vault cluster
configured  with Consul as its storage backend.

The guide walks you through:

- How to set up a monitoring stack ([Telegraf](https://www.influxdata.com/time-series-platform/telegraf/), [InfluxDB](https://www.influxdata.com/time-series-platform/influxdb/), and [Grafana](https://grafana.com/)).
- How to configure Vault and Consul to send telemetry to a monitoring agent.
- Which metrics are important to monitor, and why.

![Dashboard Example](/img/vault_cluster.png)


## Reference Materials

- [Vault Deployment Reference Architecture](/guides/operations/reference-architecture.html)
- [Vault High Availability](/guides/operations/vault-ha-consul.html)
- [Production Hardening](/guides/operations/production.html)
