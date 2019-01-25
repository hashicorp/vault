package queue

import "time"

// Item is something managed in the priority queue
type Item struct {
        value     interface{}
        priority  int64 // priority of item in queue
        index     int   // index is needed by update and maintained by heap package
        createdAt time.Time
}
