# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

disable_cache = true
disable_mlock = true

ui_config {
  enabled = true
}

listener "tcp" {
	address = "127.0.0.1:8200"
}

storage "raft" {
	path = "/storage/path/raft"
	node_id = "raft1"
	performance_multiplier = 1
	foo = "bar"
	baz = true
}

cluster_addr = "127.0.0.1:8201"
