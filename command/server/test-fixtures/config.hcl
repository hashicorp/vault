listener "tcp" {
    address = "127.0.0.1:443"
}

backend "consul" {
    foo = "bar"
}
