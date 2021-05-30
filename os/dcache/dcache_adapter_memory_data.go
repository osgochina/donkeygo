package dcache

import (
	"github.com/osgochina/donkeygo/os/dtime"
	"sync"
	"time"
)

type adapterMemoryData struct {
	mu   sync.RWMutex
	data map[interface{}]adapterMemoryItem
}

// 创建数据存储对象
func newAdapterMemoryData() *adapterMemoryData {
	return &adapterMemoryData{
		data: make(map[interface{}]adapterMemoryItem),
	}
}

// Update 更新key对应的值，如果该值存在，则把value写入，并返回旧的值，如果key对应的值不存在，则返回nil
func (that *adapterMemoryData) Update(key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	that.mu.Lock()
	defer that.mu.Unlock()

	if item, ok := that.data[key]; ok {
		that.data[key] = adapterMemoryItem{value: value, expire: item.expire}
		return item, true, nil
	}
	return nil, false, nil
}

// UpdateExpire 更新key对应值得过期时间，更新成功，返回旧的过期时间，如果未找到对应的值，则返回-1
func (that *adapterMemoryData) UpdateExpire(key interface{}, expireTime int64) (oldDuration time.Duration, err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if item, ok := that.data[key]; ok {
		that.data[key] = adapterMemoryItem{value: item.value, expire: expireTime}
		return time.Duration(item.expire-dtime.TimestampMilli()) * time.Millisecond, nil
	}

	return -1, nil
}

// Remove 移除一个或多个key，则返回移除的key的列表和被移除的最后一个key的值，如果是一个key，则返回改key和对应的值
func (that *adapterMemoryData) Remove(keys ...interface{}) (removedKeys []interface{}, value interface{}, err error) {
	that.mu.Lock()
	defer that.mu.Unlock()

	var removeKeys = make([]interface{}, 0)
	for _, key := range keys {
		item, ok := that.data[key]
		if ok {
			value = item.value
			delete(that.data, key)
			removeKeys = append(removedKeys, key)
		}
	}
	return removeKeys, value, nil
}

func (that *adapterMemoryData) Data() (map[interface{}]interface{}, error) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	m := make(map[interface{}]interface{}, len(that.data))
	for k, v := range that.data {
		if !v.IsExpired() {
			m[k] = v.value
		}
	}
	return m, nil
}
