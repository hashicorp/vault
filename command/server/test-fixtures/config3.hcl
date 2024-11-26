# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

disable_cache = true
disable_mlock = true
log_requests_level = "Basic"

ui = true

api_addr = "top_level_api_addr"
cluster_addr = "top_level_cluster_addr"

listener "tcp" {
  address = "127.0.0.1:443"
  chroot_namespace="admin/"
  disable_request_limiter = false
}

backend "consul" {
  advertise_addr = "foo"
  token = "foo"
}

ha_backend "consul" {
  bar = "baz"
  advertise_addr = "snafu"
  disable_clustering = "true"
  token = "foo"
}

service_registration "consul" {
  token = "foo"
}

telemetry {
  statsd_address = "bar"
  circonus_api_token = "baz"
  metrics_prefix = "pfx"
  usage_gauge_period = "5m"
  maximum_gauge_cardinality = 100
}

sentinel {
  additional_enabled_modules = ["http"]
}

seal "awskms" {
  region     = "us-east-1"
  access_key = "AKIAIOSFODNN7EXAMPLE"
  secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}

max_lease_ttl = "30d"
default_lease_ttl = "365d"
cluster_name = "testcluster"
pid_file = "./pidfile"
raw_storage_endpoint = true
disable_sealwrap = true
disable_sentinel_trace = true
administrative_namespace_path = "admin/"
enable_post_unseal_trace = true
post_unseal_trace_directory = "/tmp"
