package dcache

import (
	"context"
	"time"
)

// Adapter 缓存适配器接口，只要实现了这个接口的对象，都能使用dcache缓存
type Adapter interface {

	// Set 写入缓存，duration是超时时间，设置为0表示不过期。< 0 表示立即过期， > 0 表示指定过期时间
	Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error

	// Sets 批量写入缓存， duration是超时时间，设置为0表示不过期。< 0 表示立即过期， > 0 表示指定过期时间
	Sets(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error

	// SetIfNotExist 当键名不存在时写入
	SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (bool, error)

	// Get 获取key对应的值，如果不存在则返回nil
	Get(ctx context.Context, key interface{}) (interface{}, error)

	// GetOrSet 获取指定键值，如果不存在时写入，并返回键值
	GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (interface{}, error)

	// GetOrSetFunc 获取一个缓存值，当缓存不存在时执行指定的f func() (interface{}, error)，缓存该f方法的结果值，并返回该结果。
	// GetOrSetFunc的缓存方法参数f是在缓存的锁机制外执行，因此在f内部也可以嵌套调用GetOrSetFunc。
	// 但如果f的执行比较耗时，高并发的时候容易出现f被多次执行的情况(缓存设置只有第一个执行的f返回结果能够设置成功，其余的被抛弃掉)。
	GetOrSetFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error)

	// GetOrSetFuncLock 获取一个缓存值，当缓存不存在时执行指定的f func() (interface{}, error)，缓存该f方法的结果值，并返回该结果。
	// GetOrSetFuncLock的缓存方法f是在缓存的锁机制内执行，因此可以保证当缓存项不存在时只会执行一次f，但是缓存写锁的时间随着f方法的执行时间而定。
	GetOrSetFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error)

	// Contains 判断key是否存在于缓存中
	Contains(ctx context.Context, key interface{}) (bool, error)

	// GetExpire 获取指定key的过期时间
	GetExpire(ctx context.Context, key interface{}) (time.Duration, error)

	// Remove 删除一个或多个key对应的值，如果一个key，则返回删除的值，如果是多个key，则返回最后那个key对应的值
	Remove(ctx context.Context, keys ...interface{}) (value interface{}, err error)

	// Update 更新指定key对应的值为value，并返回旧的值，如果旧的值不存在，则oldValue返回nil，exist返回false
	Update(ctx context.Context, key interface{}, value interface{}) (oldValue interface{}, exist bool, err error)

	// UpdateExpire 更新指定key的过期时间，并返回旧的过期时间
	UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error)

	// Size 缓存中有多少个条目
	Size(ctx context.Context) (size int, err error)

	// Data 返回缓存的所有数据，以数组形式返回
	// 注意，这个函数如果缓存中的条目非常多，可能会占用大量的内存和调用时间，谨慎实现
	Data(ctx context.Context) (map[interface{}]interface{}, error)

	// Keys 以切片的形式返回所有的key
	Keys(ctx context.Context) ([]interface{}, error)

	// Values 以切片的形式返回所有的值
	Values(ctx context.Context) ([]interface{}, error)

	// Clear 清空缓存
	Clear(ctx context.Context) error

	// Close 关闭缓存对象
	Close(ctx context.Context) error
}
