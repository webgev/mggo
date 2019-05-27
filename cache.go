package mggo

import (
    "sync"
    "time"
)

// Cache is cache
var Cache cache

type itemCache struct {
    key        string
    value      interface{}
    expiration int64
}

type cache struct {
    sync.RWMutex
    items map[string]itemCache
    start bool
}

func init() {
    Cache = cache{
        items: map[string]itemCache{},
    }
}
func (c *cache) Set(key string, value interface{}, expiration int) {
    c.items[key] = itemCache{
        value:      value,
        key:        key,
        expiration: time.Now().Add(time.Duration(expiration) * time.Second).UnixNano(),
    }
    c.startCache()
}
func (c *cache) Get(key string) (interface{}, bool) {
    if val, ok := c.items[key]; ok {
        return val.value, ok
    }

    return nil, false
}
func (c *cache) Delete(key string) {
    keys := []string{key}
    c.clearItems(keys)
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
