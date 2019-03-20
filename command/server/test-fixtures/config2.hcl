disable_cache = true
disable_mlock = true

ui = true

api_addr = "top_level_api_addr"
cluster_addr = "top_level_cluster_addr"

listener "tcp" {
    address = "127.0.0.1:443"
}

storage "consul" {
    foo = "bar"
    redirect_addr = "foo"
}

ha_storage "consul" {
    bar = "baz"
    redirect_addr = "snafu"
    disable_clustering = "true"
}

telemetry {
    statsd_address = "bar"
    statsite_address = "foo"
    dogstatsd_addr = "127.0.0.1:7254"
    dogstatsd_tags = ["tag_1:val_1", "tag_2:val_2"]
    prometheus_retention_time = "30s"
}

max_lease_ttl = "10h"
default_lease_ttl = "10h"
cluster_name = "testcluster"
pid_file = "./pidfile"
raw_storage_endpoint = true
disable_sealwrap = true
