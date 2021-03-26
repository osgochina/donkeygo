package rwmutex

import "sync"

//读写锁
type RWMutex struct {
	*sync.RWMutex
}

//创建读写锁，传入的 safe表示该读写锁是否支持并发安全，默认是 false
func New(safe ...bool) *RWMutex {
	mu := Create(safe...)
	return &mu
}

//创建读写锁，传入的 safe表示该读写锁是否支持并发安全，默认是 false
func Create(safe ...bool) RWMutex {
	mu := RWMutex{}
	if len(safe) > 0 && safe[0] {
		mu.RWMutex = new(sync.RWMutex)
	}
	return mu
}

//判断该锁是否是并发安全锁
func (that *RWMutex) IsSafe() bool {
	return that.RWMutex != nil
}

//加锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) Lock() {
	if that.RWMutex != nil {
		that.RWMutex.Lock()
	}
}

//解锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) Unlock() {
	if that.RWMutex != nil {
		that.RWMutex.Unlock()
	}
}

//读加锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) RLock() {
	if that.RWMutex != nil {
		that.RWMutex.Lock()
	}
}

//读解锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) RUnlock() {
	if that.RWMutex != nil {
		that.RWMutex.RUnlock()
	}
}
