storage "raft" {
	path = "/storage/path/raft"
	node_id = "raft1"
	autopilot = {
        cleanup_dead_servers = "true"
        last_contact_threshold = "500ms"
        max_trailing_logs = "250"
        min_quorum = "3"
        server_stabilization_time = "10s"
	}
}
listener "tcp" {
	address = "127.0.0.1:8200"
}
disable_mlock = true
