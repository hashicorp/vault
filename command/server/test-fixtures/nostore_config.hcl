disable_cache = true
disable_mlock = true

ui = true

listener "tcp" {
    address = "127.0.0.1:1024"
    tls_disable = true
}

ha_backend "consul" {
    bar = "baz"
    advertise_addr = "snafu"
    disable_clustering = "true"
}

// No backend stanza in config!
