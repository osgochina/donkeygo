package dcache

import (
	"context"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/util/dconv"
	"time"
)

// Cache 缓存
type Cache struct {
	adapter Adapter         // 缓存适配器
	ctx     context.Context // 上下文
}

// New 创建一个缓存对象，默认使用内存缓存
func New(lruCap ...int) *Cache {
	memAdapter := newAdapterMemory(lruCap...)
	c := &Cache{
		adapter: memAdapter,
	}
	// 启动一个定时任务，每秒过期一次已过期的item
	dtimer.AddSingleton(time.Second, memAdapter.syncEventAndClearExpired)
	return c
}

// Clone clone一个新的缓存管理器
func (that *Cache) Clone() *Cache {
	return &Cache{
		adapter: that.adapter,
		ctx:     that.ctx,
	}
}

// Ctx 设置上下文
func (that *Cache) Ctx(ctx context.Context) *Cache {
	newCache := that.Clone()
	newCache.ctx = ctx
	return newCache
}

// SetAdapter 设置一个全新的适配器
func (that *Cache) SetAdapter(adapter Adapter) {
	that.adapter = adapter
}

// GetVar 从缓存中获取指定key的值，并把该值转换成Var类型
func (that *Cache) GetVar(key interface{}) (*dvar.Var, error) {
	v, err := that.Get(key)
	return dvar.New(v), err
}

// Removes 删除传入key列表对应的值
func (that *Cache) Removes(keys []interface{}) error {
	_, err := that.Remove(keys...)
	return err
}

// KeyStrings 把缓存中的key，以字符串切片的形式返回
func (that *Cache) KeyStrings() ([]string, error) {
	keys, err := that.Keys()
	if err != nil {
		return nil, err
	}
	return dconv.Strings(keys), nil
}

// 获取上下文
func (that *Cache) getCtx() context.Context {
	if that.ctx == nil {
		return context.Background()
	}
	return that.ctx
}
