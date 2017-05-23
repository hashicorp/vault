package connutil

import (
	"errors"
	"sync"
)

var (
	ErrNotInitialized = errors.New("connection has not been initalized")
)

// ConnectionProducer can be used as an embeded interface in the Database
// definition. It implements the methods dealing with individual database
// connections and is used in all the builtin database types.
type ConnectionProducer interface {
	Close() error
	Initialize(map[string]interface{}, bool) error
	Connection() (interface{}, error)

	sync.Locker
}
