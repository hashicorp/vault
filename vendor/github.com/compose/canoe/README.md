#canoe

canoe was born to get the simplicity of Hashicorp's raft library with the
full set of features you get with the raft library in etcd.

On the backend, canoe is a wrapper around
[etcd-raft](https://github.com/coreos/etcd/tree/master/raft). 

canoe is currently not considered production ready. There are still
several caveats which need to be worked out.

## Configuration API
Note: Do not mess with this unless you know what you're doing.
Canoe should handle all these requests internally, but some advanced applications - such as arbiters - may wish to do manual adjustment to cluster membership

Endpoint `:<APIPort>/peers`

### POST
Request JSON Data:
  * `id` - The ID of the new canoe node to add to the cluster
  * `raft_port` - The RaftPort of the new canoe node to add to the cluster
  * `config_port` - The APIPort of the new canoe node to add to the cluster

Response JSON Data:
  * `id` - The ID of the canoe node which the request was sent to
  * `raft_port` - The RaftPort of the canoe node which the request was sent to
  * `config_port` - The APIPort of the canoe node which the request was send to
  * `remote_peers` - List of other peers in the canoe cluster with `id`, `raft_port`, and `config_port` as fields

### GET
Response JSON Data:
  * `id` - The ID of the canoe node which the request was sent to
  * `raft_port` - The RaftPort of the canoe node which the request was sent to
  * `config_port` - The APIPort of the canoe node which the request was send to
  * `remote_peers` - List of other peers in the canoe cluster with `id`, `raft_port`, and `config_port` as fields

### DELETE
Request JSON Data:
  * `id` - The ID of the new canoe node to delete from the cluster

Success - 200
