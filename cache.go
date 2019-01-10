package authentic

import (
	"fmt"
	"sync"
	"time"
)

type (
	cache struct {
		data   map[string]*cacheValue
		maxAge time.Duration
		mtx    *sync.RWMutex
	}

	cacheValue struct {
		Created time.Time
		Value   interface{}
	}
)

func newCache(maxAge time.Duration) *cache {
	return &cache{
		data:   map[string]*cacheValue{},
		maxAge: maxAge,
		// Using mutual exclusion lock to prevent conflicts
		mtx: &sync.RWMutex{},
	}
}

func (c *cache) cacheKey(iss, kid string) string {
	return fmt.Sprintf("%s/%s", iss, kid)
}

// GetKey from cache
func (c *cache) GetKey(iss, kid string) interface{} {
	return c.Get(c.cacheKey(iss, kid))
}

// SetKey in cache
func (c *cache) SetKey(iss, kid string, value interface{}) interface{} {
	return c.Set(c.cacheKey(iss, kid), value)
}

// Get value from cache
func (c *cache) Get(key string) interface{} {
	if v := c.get(key); v != nil {
		return v.Value
	}
	return nil
}

func (c *cache) get(key string) *cacheValue {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.data[key]
}

// Set value in cache
func (c *cache) Set(key string, value interface{}) interface{} {
	c.mtx.Lock()
	c.data[key] = c.newCacheValue(value)
	c.mtx.Unlock()
	return value
}

// KeyIsExpired stale cache
func (c *cache) KeyIsExpired(iss, kid string) bool {
	return c.IsExpired(c.cacheKey(iss, kid))
}

// IsExpired stale cache
func (c *cache) IsExpired(key string) bool {
	// Created time + max age is before now
	if val := c.get(key); val != nil {
		return val.Created.Add(c.maxAge).Before(time.Now())
	}
	return true
}

func (c *cache) newCacheValue(value interface{}) *cacheValue {
	return &cacheValue{time.Now(), value}
}
