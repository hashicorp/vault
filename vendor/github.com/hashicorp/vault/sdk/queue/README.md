Vault SDK - Queue
=================

The `queue` package provides Vault plugins with a Priority Queue. It can be used
as an in-memory list of `queue.Item` sorted by their `priority`, and offers
methods to find or remove items by their key. Internally it
uses `container/heap`; see [Example Priority
Queue](https://golang.org/pkg/container/heap/#example__priorityQueue)

