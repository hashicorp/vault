disable_cache = true
disable_mlock = true

ui = true

listener "atlas" {
    token = "foobar"
    infrastructure = "foo/bar"
    endpoint = "https://foo.bar:1111"
    node_id = "foo_node"
}

listener "tcp" {
    address = "127.0.0.1:443"
}

backend "consul" {
    foo = "bar"
    advertise_addr = "foo"
}

ha_backend "consul" {
    bar = "baz"
    advertise_addr = "snafu"
    disable_clustering = "true"
}

telemetry {
    statsd_address = "bar"
    statsite_address = "foo"
}

max_lease_ttl = "10h"
default_lease_ttl = "10h"
cluster_name = "testcluster"
