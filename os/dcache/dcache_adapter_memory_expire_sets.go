package dcache

import (
	"github.com/gogf/gf/container/gset"
	"sync"
)

// 每微秒到期的key列表
type adapterMemoryExpireSets struct {
	mu         sync.RWMutex
	expireSets map[int64]*gset.Set
}

func newAdapterMemoryExpireSets() *adapterMemoryExpireSets {
	return &adapterMemoryExpireSets{
		expireSets: make(map[int64]*gset.Set),
	}
}

// Get 获取指定到期时间的key排序列表
func (that *adapterMemoryExpireSets) Get(key int64) (result *gset.Set) {
	that.mu.RLock()
	result = that.expireSets[key]
	that.mu.RUnlock()
	return
}

// GetOrNew 获取指定到期时间的key列表，如果不存在该时间，则新建一个空的排序列表
func (that *adapterMemoryExpireSets) GetOrNew(key int64) (result *gset.Set) {
	if result = that.Get(key); result != nil {
		return
	}
	that.mu.Lock()
	if es, ok := that.expireSets[key]; ok {
		result = es
	} else {
		result = gset.New(true)
		that.expireSets[key] = result
	}
	that.mu.Unlock()
	return
}

// Delete 删除指定时间的排序列表
func (that *adapterMemoryExpireSets) Delete(key int64) {
	that.mu.Lock()
	delete(that.expireSets, key)
	that.mu.Unlock()
}
