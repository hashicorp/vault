disable_cache = true
disable_mlock = true

ui = true

listener "tcp" {
    address = "127.0.0.1:1025"
    tls_cert_file = "./../api/test-fixtures/keys/cert.pem"
    tls_key_file = "./../api/test-fixtures/keys/key.pem"
}

backend "consul" {
    foo = "bar"
    advertise_addr = "foo"
    address = "127.0.0.1:8500"
}

ha_backend "consul" {
    bar = "baz"
    advertise_addr = "http://blah:8500"
    disable_clustering = "true"
    address = "127.0.0.1:8500"
}

service_registration "consul" {
    foo = "bar"
    address = "127.0.0.1:8500"
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

sentinel {
    additional_enabled_modules = []
}

max_lease_ttl = "10h"
default_lease_ttl = "10h"
cluster_name = "testcluster"
pid_file = "./pidfile"
raw_storage_endpoint = true
disable_sealwrap = true
disable_printable_check = true
