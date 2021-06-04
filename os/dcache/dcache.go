package dcache

import (
	"context"
	"github.com/osgochina/donkeygo/container/dvar"
	"time"
)

var defaultCache = New()

// Ctx 设置上下文
func Ctx(ctx context.Context) *Cache {
	return defaultCache.Ctx(ctx)
}

// Set 写入缓存，duration是超时时间，设置为0表示不过期。< 0 表示立即过期， > 0 表示指定过期时间
func Set(key interface{}, value interface{}, duration time.Duration) {
	_ = defaultCache.Set(key, value, duration)
}

// SetIfNotExist 当键名不存在时写入
func SetIfNotExist(key interface{}, value interface{}, duration time.Duration) (bool, error) {
	return defaultCache.SetIfNotExist(key, value, duration)
}

// Sets 批量写入缓存， duration是超时时间，设置为0表示不过期。< 0 表示立即过期， > 0 表示指定过期时间
func Sets(data map[interface{}]interface{}, duration time.Duration) error {
	return defaultCache.Sets(data, duration)
}

// Get 获取key对应的值，如果不存在则返回nil
func Get(key interface{}) (interface{}, error) {
	return defaultCache.Get(key)
}

// GetVar 从缓存中获取指定key的值，并把该值转换成Var类型
func GetVar(key interface{}) (*dvar.Var, error) {
	return defaultCache.GetVar(key)
}

// GetOrSet 获取指定键值，如果不存在时写入，并返回键值
func GetOrSet(key interface{}, value interface{}, duration time.Duration) (interface{}, error) {
	return defaultCache.GetOrSet(key, value, duration)
}

// GetOrSetFunc 获取一个缓存值，当缓存不存在时执行指定的f func() (interface{}, error)，缓存该f方法的结果值，并返回该结果。
// GetOrSetFunc的缓存方法参数f是在缓存的锁机制外执行，因此在f内部也可以嵌套调用GetOrSetFunc。
// 但如果f的执行比较耗时，高并发的时候容易出现f被多次执行的情况(缓存设置只有第一个执行的f返回结果能够设置成功，其余的被抛弃掉)。
func GetOrSetFunc(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return defaultCache.GetOrSetFunc(key, f, duration)
}

// GetOrSetFuncLock 获取一个缓存值，当缓存不存在时执行指定的f func() (interface{}, error)，缓存该f方法的结果值，并返回该结果。
// GetOrSetFuncLock的缓存方法f是在缓存的锁机制内执行，因此可以保证当缓存项不存在时只会执行一次f，但是缓存写锁的时间随着f方法的执行时间而定。
func GetOrSetFuncLock(key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	return defaultCache.GetOrSetFuncLock(key, f, duration)
}

// Contains 判断key是否存在于缓存中
func Contains(key interface{}) (bool, error) {
	return defaultCache.Contains(key)
}

// Remove 删除一个或多个key对应的值，如果一个key，则返回删除的值，如果是多个key，则返回最后那个key对应的值
func Remove(keys ...interface{}) (value interface{}, err error) {
	return defaultCache.Remove(keys...)
}

// Removes 删除一个或多个key对应的值，如果一个key，则返回删除的值，如果是多个key，则返回最后那个key对应的值
func Removes(keys []interface{}) {
	_ = defaultCache.Removes(keys)
}

// Data 返回缓存的所有数据，以数组形式返回
// 注意，这个函数如果缓存中的条目非常多，可能会占用大量的内存和调用时间，谨慎实现
func Data() (map[interface{}]interface{}, error) {
	return defaultCache.Data()
}

// Keys 以切片的形式返回所有的key
func Keys() ([]interface{}, error) {
	return defaultCache.Keys()
}

// KeyStrings 把缓存中的key，以字符串切片的形式返回
func KeyStrings() ([]string, error) {
	return defaultCache.KeyStrings()
}

// Values 以切片的形式返回所有的值
func Values() ([]interface{}, error) {
	return defaultCache.Values()
}

// Size 缓存中有多少个条目
func Size() (int, error) {
	return defaultCache.Size()
}

// GetExpire 获取指定key的过期时间
func GetExpire(key interface{}) (time.Duration, error) {
	return defaultCache.GetExpire(key)
}

// Update 更新指定key对应的值为value，并返回旧的值，如果旧的值不存在，则oldValue返回nil，exist返回false
func Update(key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	return defaultCache.Update(key, value)
}

// UpdateExpire 更新指定key的过期时间，并返回旧的过期时间
func UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	return defaultCache.UpdateExpire(key, duration)
}
