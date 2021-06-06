package dcache

import (
	"context"
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dset"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/os/dtimer"
	"math"
	"time"
)

// 默认的最大过期时间 math.MaxInt64/1000000.
const defaultMaxExpire = 9223372036854

// 内存缓存的实现
type adapterMemory struct {

	// 是否启动容量限制，cap>0，则表示启用，缓存最多存在cap个item，多余的按照过期时间淘汰
	cap         int
	data        *adapterMemoryData        //数据存放
	expireTimes *adapterMemoryExpireTimes // 到期时间的map，用于快速索引和删除
	expireSets  *adapterMemoryExpireSets  // 相同过期时间的 时间=>key值 映射
	lru         *adapterMemoryLru         // lru淘汰算法管理器
	lruGetList  *dlist.List               // 每次获取缓存，则把缓存的key丢到队列，更新lru管理器中的key排列顺序
	eventList   *dlist.List               // 内部数据同步的异步事件列表
	closed      *dtype.Bool               // 当前cache是否被关闭
}

//创建内存缓存对象
func newAdapterMemory(lruCap ...int) *adapterMemory {
	c := &adapterMemory{
		data:        newAdapterMemoryData(),
		expireTimes: newAdapterMemoryExpireTimes(),
		expireSets:  newAdapterMemoryExpireSets(),
		lruGetList:  dlist.New(true),
		eventList:   dlist.New(true),
		closed:      dtype.NewBool(),
	}
	if len(lruCap) > 0 {
		c.cap = lruCap[0]
		c.lru = newMemCacheLru(c)
	}
	return c
}

// Set 写入数据到缓存
func (that *adapterMemory) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	//计算有效期
	expireTime := that.getInternalExpire(duration)
	//写入数据
	that.data.Set(key, adapterMemoryItem{
		value:  value,
		expire: expireTime,
	})
	//推送写入事件
	that.eventList.PushBack(&adapterMemoryEvent{
		k: key,
		e: expireTime,
	})
	return nil
}

// Update 更新指定key对应的值为value，并返回旧的值，如果旧的值不存在，则oldValue返回nil，exist返回false
func (that *adapterMemory) Update(ctx context.Context, key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	return that.data.Update(key, value)
}

// UpdateExpire 更新指定key的有效期
func (that *adapterMemory) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	newExpireTime := that.getInternalExpire(duration)
	oldDuration, err = that.data.UpdateExpire(key, newExpireTime)
	if err != nil {
		return
	}
	if oldDuration != -1 {
		that.eventList.PushBack(&adapterMemoryEvent{
			k: key,
			e: newExpireTime,
		})
	}
	return
}

// GetExpire 获取指定key的有效期
func (that *adapterMemory) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	if item, ok := that.data.Get(key); ok {
		return time.Duration(item.expire-dtime.TimestampMilli()) * time.Millisecond, nil
	}
	return -1, nil
}

// SetIfNotExist 判断key是否存在，如果存在则写入失败，如果不存在则写入。
func (that *adapterMemory) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (bool, error) {
	isContained, err := that.Contains(ctx, key)
	if err != nil {
		return false, err
	}
	if !isContained {
		_, err := that.doSetWithLockCheck(key, value, duration)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// Sets 批量写入数据到缓存
func (that *adapterMemory) Sets(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	var (
		expireTime = that.getInternalExpire(duration)
		err        = that.data.Sets(data, expireTime)
	)
	if err != nil {
		return err
	}
	for k := range data {
		that.eventList.PushBack(&adapterMemoryEvent{
			k: k,
			e: expireTime,
		})
	}
	return nil
}

// Get 从缓存中获取指定key的数据
func (that *adapterMemory) Get(ctx context.Context, key interface{}) (interface{}, error) {
	item, ok := that.data.Get(key)
	if ok && !item.IsExpired() {
		// 如果启动了LRU算法，则把key加入到LRU列表中
		if that.cap > 0 {
			that.lruGetList.PushBack(key)
		}
		return item.value, nil
	}
	return nil, nil
}

// GetOrSet 从缓存中获取指定key的数据，如果存在则直接返回，如果不存在，则写入它
func (that *adapterMemory) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (interface{}, error) {
	v, err := that.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return that.doSetWithLockCheck(key, value, duration)
	} else {
		return v, nil
	}
}

// GetOrSetFunc 从缓存中获取指定key的数据，如果存在则直接返回，如果不存在则执行 f 方法，生成值，当值为nil则返回，不为nil则写入
func (that *adapterMemory) GetOrSetFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	v, err := that.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		value, err := f()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return that.doSetWithLockCheck(key, value, duration)
	} else {
		return v, nil
	}
}

// GetOrSetFuncLock 从缓存中获取指定key的数据，如果存在则直接返回，如果不存在则执行 f 方法，生成值，并写入
func (that *adapterMemory) GetOrSetFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (interface{}, error) {
	v, err := that.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return that.doSetWithLockCheck(key, f, duration)
	} else {
		return v, nil
	}
}

// Contains 判断指定的key是否存在与缓存中
func (that *adapterMemory) Contains(ctx context.Context, key interface{}) (bool, error) {
	v, err := that.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return v != nil, nil
}

// Remove 从缓存中移除指定的key，并返回它
func (that *adapterMemory) Remove(ctx context.Context, keys ...interface{}) (value interface{}, err error) {
	var removedKeys []interface{}
	removedKeys, value, err = that.data.Remove(keys...)
	if err != nil {
		return
	}
	for _, key := range removedKeys {
		that.eventList.PushBack(&adapterMemoryEvent{
			k: key,
			e: dtime.TimestampMilli() - 1000000,
		})
	}
	return
}

// Data 返回缓存中的所有数据
func (that *adapterMemory) Data(ctx context.Context) (map[interface{}]interface{}, error) {
	return that.data.Data()
}

// Keys 返回缓存中的所有key
func (that *adapterMemory) Keys(ctx context.Context) ([]interface{}, error) {
	return that.data.Keys()
}

// Values 返回缓存中的所有值
func (that *adapterMemory) Values(ctx context.Context) ([]interface{}, error) {
	return that.data.Values()
}

// Size 获取缓存item的数量
func (that *adapterMemory) Size(ctx context.Context) (size int, err error) {
	return that.data.Size()
}

// Clear 清空缓存数据
func (that *adapterMemory) Clear(ctx context.Context) error {
	return that.data.Clear()
}

// Close 关闭缓存
func (that *adapterMemory) Close(ctx context.Context) error {
	if that.cap > 0 {
		that.lru.Close()
	}
	that.closed.Set(true)
	return nil
}

// 写入数据到缓存
func (that *adapterMemory) doSetWithLockCheck(key interface{}, value interface{}, duration time.Duration) (result interface{}, err error) {
	expireTimestamp := that.getInternalExpire(duration)
	result, err = that.data.SetWithLock(key, value, expireTimestamp)
	that.eventList.PushBack(&adapterMemoryEvent{k: key, e: expireTimestamp})
	return
}

// 获取有效期
func (that *adapterMemory) getInternalExpire(duration time.Duration) int64 {
	if duration == 0 {
		return defaultMaxExpire
	} else {
		return dtime.TimestampMilli() + duration.Nanoseconds()/1000000
	}
}

// 生成有效期的key，按微秒取整
func (that *adapterMemory) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

// 同步添加事件，并且清理过期缓存
func (that *adapterMemory) syncEventAndClearExpired() {
	//如果缓存已关闭，则退出定时事件
	if that.closed.Val() {
		dtimer.Exit()
		return
	}
	var (
		event         *adapterMemoryEvent
		oldExpireTime int64
		newExpireTime int64
	)
	// ========================
	// Data Synchronization.
	// ========================

	for {
		//从事件列表中获取事件
		v := that.eventList.PopFront()
		if v == nil {
			break
		}
		event = v.(*adapterMemoryEvent)
		//获取指定key的老的过期时间
		oldExpireTime = that.expireTimes.Get(event.k)
		// 生成新的过期时间
		newExpireTime = that.makeExpireKey(event.e)
		//如果新老不一样，说明过期时间变了
		if newExpireTime != oldExpireTime {
			// 把缓存key添加到将要过期的时间对应的排序列表中
			that.expireSets.GetOrNew(newExpireTime).Add(event.k)
			// 如果旧的过期时间存在，则把缓存key从旧的过期时间排序列表删除
			if oldExpireTime != 0 {
				that.expireSets.GetOrNew(oldExpireTime).Remove(event.k)
			}
			// 把缓存key的过期时间更新
			that.expireTimes.Set(event.k, newExpireTime)
		}
		//如果启动了LRU算法
		if that.cap > 0 {
			that.lru.Push(event.k)
		}
	}
	// 处理lru
	if that.cap > 0 && that.lruGetList.Len() > 0 {
		for {
			if v := that.lruGetList.PopFront(); v != nil {
				that.lru.Push(v)
			} else {
				break
			}
		}
	}
	// ========================
	// Data Cleaning up.
	// ========================
	var (
		expireSet *dset.Set
		ek        = that.makeExpireKey(dtime.TimestampMilli())
		eks       = []int64{ek - 1000, ek - 2000, ek - 3000, ek - 4000, ek - 5000}
	)
	//最近5000微秒的key都过期
	for _, expireTime := range eks {
		expireSet = that.expireSets.Get(expireTime)
		if expireSet != nil {
			expireSet.Iterator(func(key interface{}) bool {
				that.clearByKey(key)
				return true
			})
			that.expireSets.Delete(expireTime)
		}
	}
}

// 清除指定key的缓存，
func (that *adapterMemory) clearByKey(key interface{}, force ...bool) {
	// 清理的时候进行二次检查
	that.data.DeleteWithDoubleCheck(key, force...)

	// 删除指定key的过期时间
	that.expireTimes.Delete(key)

	// 从LRU中删除
	if that.cap > 0 {
		that.lru.Remove(key)
	}
}
