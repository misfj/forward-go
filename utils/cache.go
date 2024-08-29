package utils

import (
	"errors"
	"sync"
	"time"
)

type single struct {
	key      string
	value    interface{}
	expireAt time.Time
}
type GCache struct {
	size int
	pool map[string]*single
	mux  sync.RWMutex
}

func New(size int) *GCache {
	return &GCache{
		size: size,
		pool: make(map[string]*single, size),
	}
}

func (gc *GCache) Add(key string, value interface{}, expire int) {
	gc.mux.Lock()
	defer gc.mux.Unlock()
	for k, _ := range gc.pool {
		if k == key {
			//覆蓋數據
			gc.pool[k].value = value
			return
		}
	}
	gc.pool[key] = &single{key: key, value: value, expireAt: time.Now().Add(time.Duration(expire) * time.Second)}
}

func (gc *GCache) Get(key string) (interface{}, error) {
	gc.mux.Lock()
	defer gc.mux.Unlock()
	v, ok := gc.pool[key]
	if !ok {
		return "xxx", errors.New("key not found")
	} else {
		if time.Now().After(v.expireAt) {
			delete(gc.pool, key) // 删除过期项
			return nil, errors.New("expired")
		}
		return v.value, nil
	}
}
