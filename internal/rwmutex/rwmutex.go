package rwmutex

import "sync"

// RWMutex 读写锁
type RWMutex struct {
	*sync.RWMutex
}

// New 创建读写锁，传入的 safe表示该读写锁是否支持并发安全，默认是 false
func New(safe ...bool) *RWMutex {
	mu := Create(safe...)
	return &mu
}

// Create 创建读写锁，传入的 safe表示该读写锁是否支持并发安全，默认是 false
func Create(safe ...bool) RWMutex {
	mu := RWMutex{}
	if len(safe) > 0 && safe[0] {
		mu.RWMutex = new(sync.RWMutex)
	}
	return mu
}

// IsSafe 判断该锁是否是并发安全锁
func (that *RWMutex) IsSafe() bool {
	return that.RWMutex != nil
}

// Lock 加锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) Lock() {
	if that.RWMutex != nil {
		that.RWMutex.Lock()
	}
}

// Unlock 解锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) Unlock() {
	if that.RWMutex != nil {
		that.RWMutex.Unlock()
	}
}

// RLock 读加锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) RLock() {
	if that.RWMutex != nil {
		that.RWMutex.Lock()
	}
}

// RUnlock 读解锁 如果不是并发安全锁，则什么都不做
func (that *RWMutex) RUnlock() {
	if that.RWMutex != nil {
		that.RWMutex.RUnlock()
	}
}
