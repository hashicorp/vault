disable_cache = true
disable_mlock = true
ui = true

listener "tcp" {
	address = "127.0.0.1:8200"
}

backend "consul" {
	advertise_addr = "foo"
	token = "foo"
}

telemetry {
	dogstatsd_tags = ["bar"]
}

cluster_addr = "127.0.0.1:8201"
