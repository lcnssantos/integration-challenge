package cache

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type cacheItem struct {
	value    interface{}
	expireAt time.Time
}

type task func() (interface{}, error)

type CacheProxy struct {
	sync.Mutex
	expirationTime time.Duration
	cache          map[string]*cacheItem
}

func NewCache() *CacheProxy {
	return &CacheProxy{}
}

func (c *CacheProxy) SetDefaultExpirationTime(duration time.Duration) *CacheProxy {
	c.expirationTime = duration
	return c
}

func (c *CacheProxy) Proxy(task task, tag string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	if c.cache == nil {
		c.cache = make(map[string]*cacheItem)
	}

	item, ok := c.cache[tag]

	if ok {
		log.Debug().Str("tag", tag).Msg("get from cache")
		if item.expireAt.After(time.Now()) {
			return item.value, nil
		}
	}

	value, err := task()

	if err != nil {
		log.
			Error().
			Err(err).
			Str("tag", tag).
			Msg("error executing task")

		return nil, err
	}

	log.
		Debug().
		Str("tag", tag).
		Msg("add to cache")

	c.cache[tag] = &cacheItem{
		value:    value,
		expireAt: time.Now().Add(c.expirationTime),
	}

	return value, nil
}
