disable_cache = true
disable_mlock = true

ui = true

listener "tcp" {
    address = "127.0.0.1:1029"
    tls_disable = true
}

backend "consul" {
    foo = "bar"
    advertise_addr = "foo"
    address = "127.0.0.1:8500"
}

ha_backend "consul" {
    bar = "baz"
    advertise_addr = "snafu"
    disable_clustering = "true"
    address = "127.0.0.1:8500"
}

service_registration "consul" {
    address = "https://consulserverIP:8500"
    foo = "bar"
    tls_cert_file = "./../vault/diagnose/test-fixtures/expiredcert.pem"
    tls_key_file = "./../vault/diagnose/test-fixtures/expiredprivatekey.pem"
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
