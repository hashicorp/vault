telemetry {
    statsd_address = "baz"
    statsite_address = "qux"
    disable_hostname = true
    usage_gauge_period = "5m"
    maximum_gauge_cardinality = 100
}
sentinel {
    additional_enabled_modules = ["http"]
}
ui=true
raw_storage_endpoint=true
default_lease_ttl = "10h"
cluster_name = "testcluster"
