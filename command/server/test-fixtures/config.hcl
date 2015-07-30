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

max_lease_duration = 10
default_lease_duration = 10