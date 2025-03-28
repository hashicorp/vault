# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

disable_cache = true
disable_mlock = true
ui = true

listener "tcp" {
  address = "127.0.0.1:8200"
  # duplicate key
  address = "127.0.0.1:8201"
}

storage "raft" {
  path = "/storage/path/raft"
  node_id = "raft1"
  performance_multiplier = 1
}

telemetry {
  statsd_address = "bar"
  usage_gauge_period = "5m"
  maximum_gauge_cardinality = 100

  statsite_address = "foo"
  dogstatsd_addr = "127.0.0.1:7254"
  dogstatsd_tags = ["tag_1:val_1", "tag_2:val_2"]
  metrics_prefix = "myprefix"
}

cluster_addr = "127.0.0.1:8201"