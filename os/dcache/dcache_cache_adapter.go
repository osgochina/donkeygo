// Package dcache 针对适配器方法的封装
package dcache

import "time"

// Set 写入缓存，duration是超时时间，设置为0表示不过期。< 0 表示立即过期， > 0 表示指定过期时间
func (that *Cache) Set(key interface{}, value interface{}, duration time.Duration) error {
	return that.adapter.Set(that.getCtx(), key, value, duration)
}

// Sets 批量写入缓存， duration是超时时间，设置为0表示不过期。< 0 表示立即过期， > 0 表示指定过期时间
func (that *Cache) Sets(data map[interface{}]interface{}, duration time.Duration) error {
	return that.adapter.Sets(that.getCtx(), data, duration)
}

// SetIfNotExist 当键名不存在时写入
func (that *Cache) SetIfNotExist(key interface{}, value interface{}, duration time.Duration) (bool, error) {
	return that.adapter.SetIfNotExist(that.getCtx(), key, value, duration)
}

// Get 获取key对应的值，如果不存在则返回nil
func (that *Cache) Get(key interface{}) (interface{}, error) {
	return that.adapter.Get(that.getCtx(), key)
}

// GetOrSet 获取指定键值，如果不存在时写入，并返回键值
func (that *Cache) GetOrSet(key interface{}, value interface{}, duration time.Duration) (interface{}, error) {
	return that.adapter.GetOrSet(that.getCtx(), key, value, duration)
}

// GetOrSetFunc 获取一个缓存值，当缓存不存在时执行指定的f func() (interface{}, error)，缓存该f方法的结果值，并返回该结果。
// GetOrSetFunc的缓存方法参数f是在缓存的锁机制外执行，因此在f内部也可以嵌套调用GetOrSetFunc。
// 但如果f的执行比较耗时，高并发的时候容易出现f被多次执行的情况(缓存设置只有第一个执行的f返回结果能够设置成功，其余的被抛弃掉)。
func (that *Cache) GetOrSetFunc(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return that.adapter.GetOrSetFunc(that.getCtx(), key, f, duration)
}

// GetOrSetFuncLock 获取一个缓存值，当缓存不存在时执行指定的f func() (interface{}, error)，缓存该f方法的结果值，并返回该结果。
// GetOrSetFuncLock的缓存方法f是在缓存的锁机制内执行，因此可以保证当缓存项不存在时只会执行一次f，但是缓存写锁的时间随着f方法的执行时间而定。
func (that *Cache) GetOrSetFuncLock(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return that.adapter.GetOrSetFuncLock(that.getCtx(), key, f, duration)
}

// Contains 判断key是否存在于缓存中
func (that *Cache) Contains(key interface{}) (bool, error) {
	return that.adapter.Contains(that.getCtx(), key)
}

// GetExpire 获取指定key的过期时间
func (that *Cache) GetExpire(key interface{}) (time.Duration, error) {
	return that.adapter.GetExpire(that.getCtx(), key)
}

// Remove 删除一个或多个key对应的值，如果一个key，则返回删除的值，如果是多个key，则返回最后那个key对应的值
func (that *Cache) Remove(keys ...interface{}) (value interface{}, err error) {
	return that.adapter.Remove(that.getCtx(), keys...)
}

// Update 更新指定key对应的值为value，并返回旧的值，如果旧的值不存在，则oldValue返回nil，exist返回false
func (that *Cache) Update(key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	return that.adapter.Update(that.getCtx(), key, value)
}

// UpdateExpire 更新指定key的过期时间，并返回旧的过期时间
func (that *Cache) UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return that.adapter.UpdateExpire(that.getCtx(), key, duration)
}

// Size 缓存中有多少个条目
func (that *Cache) Size() (size int, err error) {
	return that.adapter.Size(that.getCtx())
}

// Data 返回缓存的所有数据，以数组形式返回
// 注意，这个函数如果缓存中的条目非常多，可能会占用大量的内存和调用时间，谨慎实现
func (that *Cache) Data() (map[interface{}]interface{}, error) {
	return that.adapter.Data(that.getCtx())
}

// Keys 以切片的形式返回所有的key
func (that *Cache) Keys() ([]interface{}, error) {
	return that.adapter.Keys(that.getCtx())
}

// Values 以切片的形式返回所有的值
func (that *Cache) Values() ([]interface{}, error) {
	return that.adapter.Values(that.getCtx())
}

// Clear 清空缓存
func (that *Cache) Clear() error {
	return that.adapter.Clear(that.getCtx())
}

// Close 关闭缓存对象
func (that *Cache) Close() error {
	return that.adapter.Close(that.getCtx())
}
