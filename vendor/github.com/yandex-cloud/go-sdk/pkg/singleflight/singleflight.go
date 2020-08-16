// Package singleflight based on github.com/golang/groupcache/singleflight
package singleflight

import "sync"

// call is an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup
	val interface{}
}

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex            // protects m
	m  map[interface{}]*call // lazily initialized
}

// Do executes and returns the results of the given function, making sure that
// only one execution is in-flight for a given key at a time. If a duplicate
// comes in, the duplicate caller waits for the original to complete and
// receives the same results.
func (g *Group) Do(key interface{}, fn func() interface{}) interface{} {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val
}

// DoAsync used instead of Do, when there is not need to wait for result. It
// behaves like go { Group.Do(key, fn) }(), but doesn't create goroutine when
// there is another execution for given key in-flight.
func (g *Group) DoAsync(key interface{}, fn func() interface{}) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	if _, ok := g.m[key]; ok {
		g.mu.Unlock()
		return
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	go func() {
		c.val = fn()
		c.wg.Done()
		g.mu.Lock()
		delete(g.m, key)
		g.mu.Unlock()
	}()
}

// Call represents single deduplicated function call.
type Call struct {
	mu        sync.Mutex
	callState *callState
}

type callState struct {
	wg  sync.WaitGroup
	val interface{}
}

func (c *Call) Do(fn func() interface{}) interface{} {
	c.mu.Lock()
	if c.callState != nil {
		callState := c.callState
		c.mu.Unlock()
		callState.wg.Wait()
		return callState.val
	}

	c.callState = &callState{}
	c.callState.wg.Add(1)
	c.mu.Unlock()

	res := fn()

	c.mu.Lock()
	c.callState.val = res
	c.callState.wg.Done()
	c.callState = nil
	c.mu.Unlock()

	return res
}
