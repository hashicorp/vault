disable_cache = true
disable_mlock = true

ui = true

listener "tcp" {
    address = "127.0.0.1:1024"
    tls_disable = true
}

backend "consul" {
    address = "consulserver:8500"
    foo = "bar"
    advertise_addr = "foo"
}

ha_backend "consul" {
    address = "127.0.0.1:1024"
    bar = "baz"
    advertise_addr = "https://127.0.0.1:8500"
    disable_clustering = "true"
}

service_registration "consul" {
    address = "127.0.0.1:8500"
    foo = "bar"
}