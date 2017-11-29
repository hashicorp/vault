package plugin

import (
	"context"
	"net/rpc"
	"time"
)

type ContextCancelClient struct {
	client *rpc.Client
}

func (c *ContextCancelClient) Cancel() {
	c.client.Call("Plugin.ContextCancel", struct{}{}, struct{}{})
}

type ContextCancelServer struct {
	f context.CancelFunc
}

func (c ContextCancelServer) ContextCancel(_ struct{}, _ struct{}) {
	c.f()
}

// StorageClient is an implementation of logical.Storage that communicates
// over RPC.
type ContextClient struct {
	client    *rpc.Client
	d         chan struct{}
	cachedErr error
}

func (c *ContextClient) CancelFunc() context.CancelFunc {
	return func() { close(c.d) }
}

func (c *ContextClient) Deadline() (deadline time.Time, ok bool) {
	var reply ContextDeadlineReply
	err := c.client.Call("Plugin.Deadline", struct{}{}, &reply)
	if err != nil {
		return time.Time{}, false
	}

	return reply.Deadline, reply.Ok
}
func (c *ContextClient) Done() <-chan struct{} {
	return c.d
}

func (c *ContextClient) Err() error {
	if c.cachedErr != nil {
		return c.cachedErr
	}

	var reply ContextErrReply
	err := c.client.Call("Plugin.Err", struct{}{}, &reply)
	if err != nil {
		return err
	}

	if reply.Err != nil {
		c.cachedErr = reply.Err
	}

	return reply.Err
}

func (c *ContextClient) Value(key interface{}) interface{} {
	var reply ContextValueReply
	err := c.client.Call("Plugin.Value", key, &reply)
	if err != nil {
		return nil
	}

	return reply.Value
}

// StorageServer is a net/rpc compatible structure for serving
type ContextServer struct {
	ctx context.Context
}

func (c *ContextServer) Deadline(_ struct{}, reply *ContextDeadlineReply) error {
	d, ok := c.ctx.Deadline()
	*reply = ContextDeadlineReply{
		Deadline: d,
		Ok:       ok,
	}
	return nil
}

func (c *ContextServer) Err(_ struct{}, reply *ContextErrReply) error {
	err := c.ctx.Err()
	*reply = ContextErrReply{
		Err: err,
	}
	return nil
}

func (c *ContextServer) Value(key interface{}, reply *ContextValueReply) error {
	v := c.ctx.Value(key)
	*reply = ContextValueReply{
		Value: v,
	}
	return nil
}

type ContextErrReply struct {
	Err error
}

type ContextDeadlineReply struct {
	Deadline time.Time
	Ok       bool
}

type ContextValueReply struct {
	Value interface{}
}
