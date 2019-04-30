package queue

// Item is something managed in the priority queue
type Item struct {
	Key      string // key is used as an index in the internal map of a Queue
	Value    interface{}
	Priority int64 // priority of item in queue
	index    int   // index is needed by update and maintained by heap package
}
