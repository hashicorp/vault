disable_cache = true
disable_mlock = true
statsd_addr = "bar"
statsite_addr = "foo"

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
}

max_lease_ttl = "10h"
default_lease_ttl = "10h"
