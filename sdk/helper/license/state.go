package license

import "time"

type State struct {
	State      string
	ExpiryTime time.Time
	Terminated bool
}
