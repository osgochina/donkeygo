package dcache

import (
	"github.com/osgochina/donkeygo/os/dtime"
	"sync"
	"time"
)

// 内存缓存存放数据的结构
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

// Data 把缓存中的所有导出
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

// Keys 获取缓存的所有key
func (that *adapterMemoryData) Keys() ([]interface{}, error) {
	that.mu.RLock()
	var (
		index = 0
		keys  = make([]interface{}, len(that.data))
	)
	for k, v := range that.data {
		if !v.IsExpired() {
			keys[index] = k
			index++
		}
	}
	that.mu.RUnlock()
	return keys, nil
}

// Values 获取缓存中的所有值，导出为值组成的数组
func (that *adapterMemoryData) Values() ([]interface{}, error) {
	that.mu.RLock()
	var (
		index  = 0
		values = make([]interface{}, len(that.data))
	)
	for _, v := range that.data {
		if !v.IsExpired() {
			values[index] = v.value
			index++
		}
	}
	that.mu.RUnlock()
	return values, nil
}

// Size 获取缓存中的item数量
func (that *adapterMemoryData) Size() (size int, err error) {
	that.mu.RLock()
	size = len(that.data)
	that.mu.RUnlock()
	return size, nil
}

// Clear 清空缓存
func (that *adapterMemoryData) Clear() error {
	that.mu.Lock()
	defer that.mu.Unlock()
	that.data = make(map[interface{}]adapterMemoryItem)
	return nil
}

// Get 获取缓存中指定key的值
func (that *adapterMemoryData) Get(key interface{}) (item adapterMemoryItem, ok bool) {
	that.mu.RLock()
	item, ok = that.data[key]
	that.mu.RUnlock()
	return
}

// Set 设置值
func (that *adapterMemoryData) Set(key interface{}, value adapterMemoryItem) {
	that.mu.Lock()
	that.data[key] = value
	that.mu.Unlock()
}

// Sets 批量设置值，并且设置它的有效期
func (that *adapterMemoryData) Sets(data map[interface{}]interface{}, expireTime int64) error {
	that.mu.Lock()
	for k, v := range data {
		that.data[k] = adapterMemoryItem{
			value:  v,
			expire: expireTime,
		}
	}
	that.mu.Unlock()
	return nil
}

// SetWithLock 设置指定key的值，如果该key已存在，并且还在有效期内，则直接返回该值
// 如果要设置的key不存在或者已过期，则把该值写入，并设置为传入的有效期，如果传入的value是func，则执行该func，并把返回的值作为value写入
func (that *adapterMemoryData) SetWithLock(key interface{}, value interface{}, expireTimestamp int64) (interface{}, error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if v, ok := that.data[key]; ok && !v.IsExpired() {
		return v.value, nil
	}
	if f, ok := value.(func() (interface{}, error)); ok {
		v, err := f()
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		} else {
			value = v
		}
	}
	that.data[key] = adapterMemoryItem{value: value, expire: expireTimestamp}
	return value, nil
}

// DeleteWithDoubleCheck 删除指定的key的值，删除前要做双重检查
func (that *adapterMemoryData) DeleteWithDoubleCheck(key interface{}, force ...bool) {
	that.mu.Lock()
	// 删除之前需要进行双重检查，值必须存在，且在有效期之内，如果强制删除，则直接删除
	if item, ok := that.data[key]; (ok && item.IsExpired()) || (len(force) > 0 && force[0]) {
		delete(that.data, key)
	}
	that.mu.Unlock()
}
