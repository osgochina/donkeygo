package dcache

import (
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dtimer"
	"time"
)

type adapterMemoryLru struct {
	cache   *adapterMemory // 缓存对象
	data    *dmap.Map      // 缓存list中的key，方便快速索引
	list    *dlist.List    // lru算法管理器的key列表
	rawList *dlist.List    // 添加历史记录，需要处理
	closed  *dtype.Bool    // 是否关闭
}

//创建内存缓存lru淘汰算法
func newMemCacheLru(cache *adapterMemory) *adapterMemoryLru {
	lru := &adapterMemoryLru{
		cache:   cache,
		data:    dmap.New(true),
		rawList: dlist.New(true),
		list:    dlist.New(true),
		closed:  dtype.NewBool(),
	}
	//每秒执行一次淘汰
	dtimer.AddSingleton(time.Second, lru.SyncAndClear)
	return lru
}

// Close 关闭lru
func (that *adapterMemoryLru) Close() {
	that.closed.Set(true)
}

// Remove 移除
func (that *adapterMemoryLru) Remove(key interface{}) {
	if v := that.data.Get(key); v != nil {
		that.data.Remove(key)
		that.list.Remove(v.(*dlist.Element))
	}
}

// Size 管理器中的数据长度
func (that *adapterMemoryLru) Size() int {
	return that.data.Size()
}

// Push 把缓存key放入队列尾部
func (that *adapterMemoryLru) Push(key interface{}) {
	that.rawList.PushBack(key)
}

// Pop 从队列尾部取出key
func (that *adapterMemoryLru) Pop() interface{} {
	if v := that.list.PopBack(); v != nil {
		that.data.Remove(v)
		return v
	}
	return nil
}

// SyncAndClear 同步缓存key的使用记录，并清除超过缓存容量，最后为使用的key
func (that *adapterMemoryLru) SyncAndClear() {
	//如果lru管理器关闭，则结束执行该定时任务
	if that.closed.Val() {
		dtimer.Exit()
		return
	}
	for {
		// 从key的处理历史记录中获取key
		if v := that.rawList.PopFront(); v != nil {
			//如果该key已经存在与lru队列，则把该key从lru队列中删除
			v2 := that.data.Get(v)
			if v2 != nil {
				that.list.Remove(v2.(*dlist.Element))
			}
			// 把key放到lru队列头部，这样保证了，头部永远是最后使用过的key
			that.data.Set(v, that.list.PushFront(v))
		} else {
			break
		}
	}

	//判断缓存key的数目超出容量限制多少，删除多余的key
	for i := that.Size() - that.cache.cap; i > 0; i-- {
		// 把队列尾部取出最久不被使用的key，并删除
		if s := that.Pop(); s != nil {
			that.cache.clearByKey(s, true)
		}
	}
}
