storage "raft" {
	path = "/storage/path/raft"
	node_id = "raft1"
	retry_join = [
    {
    "leader_api_addr" = "http://127.0.0.1:8200"
    },
    {
    "leader_api_addr" = "http://127.0.0.2:8200"
    },
    {
    "leader_api_addr" = "http://127.0.0.3:8200"
    }
  ]
}
listener "tcp" {
	address = "127.0.0.1:8200"
}
disable_mlock = true
