telemetry {
    statsd_address = "baz"
    statsite_address = "qux"
    disable_hostname = true
}

default_lease_ttl = "10h"
cluster_name = "testcluster"
enable_cors     = true
allowed_origins = "http://localhost:8[0-9]{3}"
