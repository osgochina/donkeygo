package dcache

import "sync"

// 存放key的过期时间，方便快速索引
type adapterMemoryExpireTimes struct {
	mu          sync.RWMutex
	expireTimes map[interface{}]int64
}

func newAdapterMemoryExpireTimes() *adapterMemoryExpireTimes {
	return &adapterMemoryExpireTimes{
		expireTimes: make(map[interface{}]int64),
	}
}

func (that *adapterMemoryExpireTimes) Get(key interface{}) (value int64) {
	that.mu.RLock()
	value = that.expireTimes[key]
	that.mu.RUnlock()
	return
}

func (that *adapterMemoryExpireTimes) Set(key interface{}, value int64) {
	that.mu.Lock()
	that.expireTimes[key] = value
	that.mu.Unlock()
}

func (that *adapterMemoryExpireTimes) Delete(key interface{}) {
	that.mu.Lock()
	delete(that.expireTimes, key)
	that.mu.Unlock()
}
