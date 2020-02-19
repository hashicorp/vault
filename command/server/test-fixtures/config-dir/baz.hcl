telemetry {
    statsd_address = "baz"
    statsite_address = "qux"
    disable_hostname = true
}

backend "consul" {
    path = "vault"
    foo = "baz"
}

ui=true
raw_storage_endpoint=true
default_lease_ttl = "10h"
cluster_name = "testcluster"

disable_clustering = true

disable_cache = true