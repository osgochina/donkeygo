package dpool

import (
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/os/dtimer"
	"time"
)

// Pool 并发安全的对象复用池
type Pool struct {
	list       *dlist.List
	closed     *dtype.Bool
	TTL        time.Duration
	NewFunc    func() (interface{}, error)
	ExpireFunc func(interface{})
}

type poolItem struct {
	value    interface{}
	expireAt int64
}

type NewFunc func() (interface{}, error)

type ExpireFunc func(interface{})

// New 创建对象池
// 注意：
// ttl = 0: 表示不过期。
// ttl < 0: 使用后立即过期
// ttl > 0: 设置的过期时间
func New(ttl time.Duration, newFunc NewFunc, expireFunc ...ExpireFunc) *Pool {
	p := &Pool{
		list:    dlist.New(true),
		closed:  dtype.NewBool(),
		TTL:     ttl,
		NewFunc: newFunc,
	}
	if len(expireFunc) > 0 {
		p.ExpireFunc = expireFunc[0]
	}
	dtimer.AddSingleton(time.Second, p.checkExpireItems)
	return p
}

// 检查对象池中的对象是否过期
func (that *Pool) checkExpireItems() {
	//检查到对象池已关闭，则结束当前定时器，且退出定时器，以后不会执行
	if that.closed.Val() {

		if that.ExpireFunc != nil {
			for {
				r := that.list.PopFront()
				if r == nil {
					break
				}
				that.ExpireFunc(r.(*poolItem).value)
			}
		}
		dtimer.Exit()
	}

	//全局设置了不会过期
	if that.TTL == 0 {
		return
	}
	//最后一个对象的过期时间
	var latestExpire int64 = -1

	//当前时间的毫秒数
	var timestampMilli = dtime.TimestampMilli()

	for {
		//如果最后过期时间大于当前时间，说明池子里面的对象都是没有过期的
		if latestExpire > timestampMilli {
			break
		}
		//从列表中取出一个对象
		r := that.list.PopFront()
		if r == nil {
			break
		}
		item := r.(*poolItem)
		latestExpire = item.expireAt
		// 判断对象是否过期，如果未过期，则把对象放回去，结束检查
		if item.expireAt > timestampMilli {
			that.list.PushFront(item)
			break
		}
		//如果过期回调函数存在，则执行
		if that.ExpireFunc != nil {
			that.ExpireFunc(item.value)
		}
	}
}

// Close 关闭该对象池
func (that *Pool) Close() {
	that.closed.Set(true)
}

// Put 把对象放回对象池
func (that *Pool) Put(value interface{}) error {
	if that.closed.Val() {
		return derror.New("pool is closed")
	}
	item := &poolItem{
		value: value,
	}
	if that.TTL == 0 {
		item.expireAt = 0
	} else {
		item.expireAt = dtime.TimestampMilli() + that.TTL.Nanoseconds()/1e6
	}
	that.list.PushBack(item)
	return nil
}

// Get 从对象池中获取对象，如果池中不存在对象，则调用对象生成方法生成对象
func (that *Pool) Get() (interface{}, error) {
	for !that.closed.Val() {
		r := that.list.PopFront()
		if r == nil {
			break
		}
		item := r.(*poolItem)
		if item.expireAt == 0 || item.expireAt > dtime.TimestampMilli() {
			return item.value, nil
		}
		//如果拿到的对象发现已过期，则不反回该对象，并且调用过期回调方法
		if that.ExpireFunc != nil {
			that.ExpireFunc(item.value)
		}
	}
	if that.NewFunc != nil {
		return that.NewFunc()
	}

	return nil, derror.New("pool is empty")
}

// Size 对象池中的对象数目
func (that *Pool) Size() int {
	return that.list.Len()
}

// Clear 清除对象池
func (that *Pool) Clear() {
	if that.ExpireFunc != nil {
		for {
			r := that.list.PopFront()
			if r == nil {
				break
			}
			that.ExpireFunc(r.(*poolItem).value)
		}
	} else {
		that.list.RemoveAll()
	}
}
