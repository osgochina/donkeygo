package dmlock

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/os/dmutex"
)

// Locker 内存锁，底层使用map实现
type Locker struct {
	m *dmap.StrAnyMap
}

func New() *Locker {
	return &Locker{
		m: dmap.NewStrAnyMap(true),
	}
}

// Lock 对该key增加写锁
func (that *Locker) Lock(key string) {
	that.getOrNewMutex(key).Lock()
}

// TryLock 尝试对该key增加写锁，加锁成功返回true，加锁失败返回false
func (that *Locker) TryLock(key string) bool {
	return that.getOrNewMutex(key).TryLock()
}

// Unlock 对该key解锁
func (that *Locker) Unlock(key string) {
	if v := that.m.Get(key); v != nil {
		v.(*dmutex.Mutex).Unlock()
	}
}

// RLock 对该key增加读锁
func (that *Locker) RLock(key string) {
	that.getOrNewMutex(key).RLock()
}

// TryRLock 尝试对该key增加读锁，加锁成功返回true，加锁失败返回false
func (that *Locker) TryRLock(key string) bool {
	return that.getOrNewMutex(key).TryRLock()
}

// RUnlock 对指定的key解锁
func (that *Locker) RUnlock(key string) {
	if v := that.m.Get(key); v != nil {
		v.(*dmutex.Mutex).RUnlock()
	}
}

//LockFunc 对指定的key加锁，并执行方法
func (that *Locker) LockFunc(key string, f func()) {
	that.Lock(key)
	defer that.Unlock(key)
	f()
}

// RLockFunc 对指定的key加读锁，并执行方法
func (that *Locker) RLockFunc(key string, f func()) {
	that.RLock(key)
	defer that.RUnlock(key)
	f()
}

// TryLockFunc 尝试对指定的key加锁，加锁成功则执行f方法，如果加锁并执行成功则返回true，否则返回false
func (that *Locker) TryLockFunc(key string, f func()) bool {
	if that.TryLock(key) {
		defer that.Unlock(key)
		f()
		return true
	}
	return false
}

// TryRLockFunc 尝试对指定的key加读锁，加锁成功则执行方法f，执行成功返回true，否则返回false
func (that *Locker) TryRLockFunc(key string, f func()) bool {
	if that.TryRLock(key) {
		defer that.RUnlock(key)
		f()
		return true
	}
	return false
}

// Remove 移除指定key的锁
func (that *Locker) Remove(key string) {
	that.m.Remove(key)
}

// Clear 清空内存锁中的所有锁
func (that *Locker) Clear() {
	that.m.Clear()
}

// 每个key都是创建一个高级互斥锁，如果key存在，则返回该锁，不存在则创建
func (that *Locker) getOrNewMutex(key string) *dmutex.Mutex {
	return that.m.GetOrSetFuncLock(key, func() interface{} {
		return dmutex.New()
	}).(*dmutex.Mutex)
}
