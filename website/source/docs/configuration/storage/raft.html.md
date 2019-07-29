---
layout: "docs"
page_title: "Raft - Storage Backends - Configuration"
sidebar_title: "Raft"
sidebar_current: "docs-configuration-storage-raft"
description: |-

 The Raft storage backend is used to persist Vault's data. Unlike all the other
 storage backends, this backend does not operate from a single source for the
 data. Instead all the nodes in a Vault cluster will have a replicated copy of
 the entire data. The data is replicated across the nodes using the Raft
 Consensus Algorithm.

---

# Raft Storage Backend

The Raft storage backend is used to persist Vault's data. Unlike other storage
backends, Raft storage does not operate from a single source of data. Instead
all the nodes in a Vault cluster will have a replicated copy of Vault's data.
Data gets replicated across the all the nodes via the [Raft Consensus
Algorithm][raft].


- **High Availability** – the Raft storage backend supports high availability.

- **HashiCorp Supported** – the Raft storage backend is officially supported
  by HashiCorp.

```hcl
storage "raft" {
  path = "/path/to/raft/data"
  node_id = "raft_node_1"
}
cluster_addr = "http://127.0.0.1:8201"
```

**Note:** When using the Raft storage backend, it is required to provide `cluster_addr` to indicate the address and port to be used for communication between the nodes in the Raft cluster.

## `raft` Parameters

- `path` `(string: "")` – The file system path where all the Vault data gets
  stored.

- `node_id` `(string: "")` - The identifier for the node in the Raft cluster.

[raft]: https://raft.github.io/ "The Raft Consensus Algorithm"
