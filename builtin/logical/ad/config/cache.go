package config

import (
	"sync"
)

func newCache(conf *EngineConf) *cache {
	return &cache{
		conf:     conf,
		confLock: &sync.RWMutex{},

		isValid:          true,
		invalidationLock: &sync.RWMutex{},
	}
}

// cache is thread-safe
type cache struct {
	conf     *EngineConf
	confLock *sync.RWMutex

	isValid          bool
	invalidationLock *sync.RWMutex
}

func (c *cache) Get() (*EngineConf, bool) {

	c.invalidationLock.RLock()
	defer c.invalidationLock.RLock()

	c.confLock.RLock()
	defer c.confLock.RUnlock()

	if !c.isValid {
		return nil, false
	}
	return c.conf, true
}

func (c *cache) Set(conf *EngineConf) {

	c.invalidationLock.Lock()
	defer c.invalidationLock.Unlock()

	c.confLock.Lock()
	defer c.confLock.Unlock()

	c.conf = conf
	c.isValid = true
}

func (c *cache) Invalidate() {

	c.invalidationLock.Lock()
	defer c.invalidationLock.Unlock()

	c.isValid = false
}
