package mggo

import (
	"fmt"
	"sync"
	"time"
)

type cacheMethod struct {
	cacheType
	expiration int
}

type cacheType int

//CacheTypeUser is cache by user
//CacheTypeMethodParams is cache by method params
const (
	CacheTypeUser cacheType = iota
	CacheTypeMethodParams
)
const (
	cacheUserPrefix   = "-userid-"
	cacheParamsPrefix = "-params-"
)

// Cache is cache
var Cache cache

type itemCache struct {
	method     string
	key        string
	value      interface{}
	expiration int64
}

type cache struct {
	sync.RWMutex
	items        map[string]itemCache
	start        bool
	cacheMethods map[string]cacheMethod
}

func init() {
	Cache = cache{
		items:        map[string]itemCache{},
		cacheMethods: map[string]cacheMethod{},
	}
}

// AddMethod is wil add method in cache
func (c *cache) AddMethod(method string, cType cacheType, expiration int) {
	c.cacheMethods[method] = cacheMethod{cType, expiration}
}

// ClearMethodCache is clear method cache
func (c *cache) ClearCacheMethod(method string) {
	keys := []string{}
	for _, v := range c.items {
		if v.method == method {
			keys = append(keys, v.key)
		}
	}
	if len(keys) > 0 {
		c.clearItems(keys)
	}
}

// ClearMethodCacheByUserID is clear method cache by user
func (c *cache) ClearCacheMethodByUserID(method string, id int) {
	keys := []string{}
	for _, v := range c.items {
		key := method + cacheUserPrefix + string(id)
		if v.key == key {
			keys = append(keys, v.key)
		}
	}
	if len(keys) > 0 {
		c.clearItems(keys)
	}
}

func (c *cache) set(method, key string, value interface{}, expiration int) {
	c.items[key] = itemCache{
		method:     method,
		value:      value,
		key:        key,
		expiration: time.Now().Add(time.Duration(expiration) * time.Second).UnixNano(),
	}
	c.startCache()
}

func (c *cache) get(key string) (interface{}, bool) {
	if val, ok := c.items[key]; ok {
		return val.value, ok
	}

	return nil, false
}

func (c *cache) delete(key string) {
	keys := []string{key}
	c.clearItems(keys)
}

func (c *cache) getMethod(method string, params interface{}) (interface{}, bool) {
	if v, ok := c.cacheMethods[method]; ok {
		var key string
		if v.cacheType == CacheTypeMethodParams {
			key = method + cacheParamsPrefix + fmt.Sprintf("%v", params)
		} else if v.cacheType == CacheTypeUser {
			id := SAP{}.SessionUserID()
			if id == 0 {
				return nil, false
			}
			key = method + cacheUserPrefix + string(id)
		}
		return c.get(key)
	}
	return nil, false
}

func (c *cache) setMethod(method string, value interface{}, params interface{}) bool {
	if v, ok := c.cacheMethods[method]; ok {
		var key string
		if v.cacheType == CacheTypeMethodParams {
			key = method + cacheParamsPrefix + fmt.Sprintf("%v", params)
		} else if v.cacheType == CacheTypeUser {
			id := SAP{}.SessionUserID()
			if id == 0 {
				return false
			}
			key = method + cacheUserPrefix + string(id)
		}
		c.set(method, key, value, v.expiration)
	}
	return false
}

func (c *cache) isset(method string) bool {
	_, ok := c.cacheMethods[method]
	return ok
}

func (c *cache) startCache() {
	if c.start == false {
		c.start = true
		go c.gc()
	}
}

func (c *cache) gc() {
	for {
		if len(c.items) == 0 {
			c.start = false
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
		time.Sleep(1 * time.Second)
	}
}

// expiredKeys возвращает список "просроченных" ключей
func (c *cache) expiredKeys() (keys []string) {
	c.RLock()

	defer c.RUnlock()
	for k, i := range c.items {
		if time.Now().UnixNano() > i.expiration && i.expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
