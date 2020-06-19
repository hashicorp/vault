storage "inmem" {}

ha_storage "raft" {
	path = "/storage/path/raft"
	node_id = "raft1"
	retry_join = {}
}

listener "tcp" {
	address = "127.0.0.1:8200"
}
disable_mlock = true
