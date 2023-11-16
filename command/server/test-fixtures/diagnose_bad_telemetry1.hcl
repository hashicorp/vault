# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

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
	circonus_check_id = "bar"
}

cluster_addr = "127.0.0.1:8201"
