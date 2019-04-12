package lib

import (
	"context"
	"time"
)

// StopChannelContext implements the context.Context interface
// You provide the channel to select on to determine whether
// the context should be canceled and other code such
// as the rate.Limiter will automatically use the channel
// appropriately
type StopChannelContext struct {
	StopCh <-chan struct{}
}

func (c *StopChannelContext) Deadline() (deadline time.Time, ok bool) {
	ok = false
	return
}

func (c *StopChannelContext) Done() <-chan struct{} {
	return c.StopCh
}

func (c *StopChannelContext) Err() error {
	select {
	case <-c.StopCh:
		return context.Canceled
	default:
		return nil
	}
}

func (c *StopChannelContext) Value(key interface{}) interface{} {
	return nil
}
