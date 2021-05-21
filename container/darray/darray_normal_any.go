package darray

import (
	"github.com/osgochina/donkeygo/internal/rwmutex"
)

type Array struct {
	mu    rwmutex.RWMutex
	array []interface{}
}

// New 创建一个不限制大小的数组
func New(safe ...bool) *Array {
	return NewArraySize(0, 0, safe...)
}

// NewArray 创建一个不限制大小的数组
func NewArray(safe ...bool) *Array {
	return NewArraySize(0, 0, safe...)
}

// NewArraySize 创建一个指定大小和上限的数组
// size 数组大小
// cap  数组元素上限
// safe 是否并发安全
func NewArraySize(size int, cap int, safe ...bool) *Array {
	return &Array{
		mu:    rwmutex.Create(safe...),
		array: make([]interface{}, size, cap),
	}
}

func (that *Array) PushRight(value ...interface{}) *Array {
	that.mu.Lock()
	that.array = append(that.array, value...)
	that.mu.Unlock()
	return that
}

func (that *Array) Append(value ...interface{}) *Array {
	that.PushRight(value...)
	return that
}

func (that *Array) Len() int {
	that.mu.RLock()
	length := len(that.array)
	that.mu.RUnlock()
	return length
}

func (that *Array) Slice() []interface{} {
	if that.mu.IsSafe() {
		that.mu.RLock()
		defer that.mu.RUnlock()
		array := make([]interface{}, len(that.array))
		copy(array, that.array)
		return array
	} else {
		return that.array
	}
}
