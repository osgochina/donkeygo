package dset

import (
	"bytes"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/text/dstr"
	"github.com/osgochina/donkeygo/util/dconv"
)

type Set struct {
	mu   rwmutex.RWMutex
	data map[interface{}]struct{}
}

// New 创建并返回一个Set
func New(safe ...bool) *Set {
	return NewSet(safe...)
}

func NewSet(safe ...bool) *Set {
	return &Set{
		mu:   rwmutex.Create(safe...),
		data: make(map[interface{}]struct{}),
	}
}

// NewFrom 从items创建一个set
func NewFrom(items interface{}, safe ...bool) *Set {
	m := make(map[interface{}]struct{})
	for _, v := range dconv.Interfaces(items) {
		m[v] = struct{}{}
	}
	return &Set{
		data: m,
		mu:   rwmutex.Create(safe...),
	}
}

// Iterator 迭代集合
func (that *Set) Iterator(f func(interface{}) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()

	for v, _ := range that.data {
		if !f(v) {
			break
		}
	}
}

// Add 添加一个元素到集合
func (that *Set) Add(items ...interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[interface{}]struct{})
	}
	for _, v := range items {
		that.data[v] = struct{}{}
	}
}

// AddIfNotExist 如果元素不存在则添加
func (that *Set) AddIfNotExist(item interface{}) bool {
	if item == nil {
		return false
	}

	if !that.Contains(item) {
		that.mu.Lock()
		defer that.mu.Unlock()
		if that.data == nil {
			that.data = make(map[interface{}]struct{})
		}
		if _, ok := that.data[item]; !ok {
			that.data[item] = struct{}{}
			return true
		}
	}
	return false
}

// AddIfNotExistFunc 如果元素不存在，且执行f方法，返回成功，则添加元素到set
// 执行f方法的时候不加锁
func (that *Set) AddIfNotExistFunc(item interface{}, f func() bool) bool {
	if item == nil {
		return false
	}
	if !that.Contains(item) {
		if f() {
			that.mu.Lock()
			defer that.mu.Unlock()
			if that.data == nil {
				that.data = make(map[interface{}]struct{})
			}
			if _, ok := that.data[item]; !ok {
				that.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// AddIfNotExistFuncLock 如果元素不存在，且执行f方法，返回成功，则添加元素到set
//// 执行f方法的时候加锁
func (that *Set) AddIfNotExistFuncLock(item interface{}, f func() bool) bool {
	if item == nil {
		return false
	}
	if !that.Contains(item) {
		that.mu.Lock()
		defer that.mu.Unlock()
		if that.data == nil {
			that.data = make(map[interface{}]struct{})
		}
		if f() {
			if _, ok := that.data[item]; !ok {
				that.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// Contains 判断集合中是否存在item元素
func (that *Set) Contains(item interface{}) bool {
	var ok bool
	that.mu.RLock()
	if that.data != nil {
		_, ok = that.data[item]
	}
	that.mu.RUnlock()
	return ok
}

// Remove 从集合中移除元素
func (that *Set) Remove(item interface{}) {
	that.mu.Lock()
	if that.data != nil {
		delete(that.data, item)
	}
	that.mu.Unlock()
}

// Size 获取集合的长度
func (that *Set) Size() int {
	that.mu.RLock()
	l := len(that.data)
	that.mu.RUnlock()
	return l
}

// Clear 清空集合
func (that *Set) Clear() {
	that.mu.Lock()
	that.data = make(map[interface{}]struct{})
	that.mu.Unlock()
}

// Slice 返回集合内的元素，以切片的形式返回
func (that *Set) Slice() []interface{} {
	that.mu.RLock()
	var (
		i   = 0
		ret = make([]interface{}, len(that.data))
	)
	for item := range that.data {
		ret[i] = item
		i++
	}
	that.mu.RUnlock()
	return ret
}

// Join 把集合的元素转换成字符串，然后以glue作为分隔符链接起来，以字符的形式返回
func (that *Set) Join(glue string) string {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if len(that.data) == 0 {
		return ""
	}
	var (
		l      = len(that.data)
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	for k, _ := range that.data {
		buffer.WriteString(dconv.String(k))
		if i != l-1 {
			buffer.WriteString(glue)
		}
		i++
	}
	return buffer.String()
}

// 把集合转换成字符串
func (that *Set) String() string {
	that.mu.RLock()
	defer that.mu.RUnlock()
	var (
		s      = ""
		l      = len(that.data)
		i      = 0
		buffer = bytes.NewBuffer(nil)
	)
	buffer.WriteByte('[')
	for k, _ := range that.data {
		s = dconv.String(k)
		if dstr.IsNumeric(s) {
			buffer.WriteString(s)
		} else {
			buffer.WriteString(`"` + dstr.QuoteMeta(s, `"\`) + `"`)
		}
		if i != l-1 {
			buffer.WriteByte(',')
		}
		i++
	}
	buffer.WriteByte(']')
	return buffer.String()
}

// LockFunc 加锁执行自定义方法
func (that *Set) LockFunc(f func(m map[interface{}]struct{})) {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.data)
}

// RLockFunc 加读锁执行自定义方法
func (that *Set) RLockFunc(f func(m map[interface{}]struct{})) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.data)
}

// Equal 判断两个集合是否相等
func (that *Set) Equal(other *Set) bool {
	if that == other {
		return true
	}
	that.mu.RLock()
	defer that.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	if len(that.data) != len(other.data) {
		return false
	}
	for key := range that.data {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}
	return true
}

// IsSubsetOf  检查当前集合是否为other集合的子集
func (that *Set) IsSubsetOf(other *Set) bool {
	if that == other {
		return true
	}
	that.mu.RLock()
	defer that.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for key := range that.data {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}
	return true
}

// Union 求多个集合的并集
func (that *Set) Union(others ...*Set) (newSet *Set) {
	newSet = NewSet()
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, other := range others {
		if that != other {
			other.mu.RLock()
		}
		for k, v := range that.data {
			newSet.data[k] = v
		}
		if that != other {
			for k, v := range other.data {
				newSet.data[k] = v
			}
		}
		if that != other {
			other.mu.RUnlock()
		}
	}

	return
}

// Diff 求多个集合差集
func (that *Set) Diff(others ...*Set) (newSet *Set) {
	newSet = NewSet()
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, other := range others {
		if that == other {
			continue
		}
		other.mu.RLock()
		for k, v := range that.data {
			if _, ok := other.data[k]; !ok {
				newSet.data[k] = v
			}
		}
		other.mu.RUnlock()
	}
	return
}

// Intersect 求多个集合的交集
func (that *Set) Intersect(others ...*Set) (newSet *Set) {
	newSet = NewSet()
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, other := range others {
		if that != other {
			other.mu.RLock()
		}
		for k, v := range that.data {
			if _, ok := other.data[k]; ok {
				newSet.data[k] = v
			}
		}
		if that != other {
			other.mu.RUnlock()
		}
	}
	return
}

// Complement 求指定full集合对与当前集合的补集
func (that *Set) Complement(full *Set) (newSet *Set) {
	newSet = NewSet()
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that != full {
		full.mu.RLock()
		defer full.mu.RUnlock()
	}
	for k, v := range full.data {
		if _, ok := that.data[k]; !ok {
			newSet.data[k] = v
		}
	}
	return
}

// Merge 合并一个或多个集合
func (that *Set) Merge(others ...*Set) *Set {
	that.mu.Lock()
	defer that.mu.Unlock()
	for _, other := range others {
		if that != other {
			other.mu.RLock()
		}
		for k, v := range other.data {
			that.data[k] = v
		}
		if that != other {
			other.mu.RUnlock()
		}
	}
	return that
}

// Sum 集合所有元素累加
func (that *Set) Sum() (sum int) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for k, _ := range that.data {
		sum += dconv.Int(k)
	}
	return
}

// Pop 集合中取出一个元素
func (that *Set) Pop() interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	for k, _ := range that.data {
		delete(that.data, k)
		return k
	}
	return nil
}

// Pops 从集合中取出指定size的元素
func (that *Set) Pops(size int) []interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if size > len(that.data) || size == -1 {
		size = len(that.data)
	}
	if size <= 0 {
		return nil
	}
	index := 0
	array := make([]interface{}, size)
	for k, _ := range that.data {
		delete(that.data, k)
		array[index] = k
		index++
		if index == size {
			break
		}
	}
	return array
}

// Walk 针对集合中的每个元素执行一次k方法
func (that *Set) Walk(f func(item interface{}) interface{}) *Set {
	that.mu.Lock()
	defer that.mu.Unlock()
	m := make(map[interface{}]struct{}, len(that.data))
	for k, v := range that.data {
		m[f(k)] = v
	}
	that.data = m
	return that
}

// MarshalJSON 把集合格式化成json格式
func (that *Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(that.Slice())
}

// UnmarshalJSON 把json格式转换成集合
func (that *Set) UnmarshalJSON(b []byte) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[interface{}]struct{})
	}
	var array []interface{}
	if err := json.UnmarshalUseNumber(b, &array); err != nil {
		return err
	}
	for _, v := range array {
		that.data[v] = struct{}{}
	}
	return nil
}

// UnmarshalValue 把整个值从对象中转换成集合能识别的值
func (that *Set) UnmarshalValue(value interface{}) (err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[interface{}]struct{})
	}
	var array []interface{}
	switch value.(type) {
	case string, []byte:
		err = json.UnmarshalUseNumber(dconv.Bytes(value), &array)
	default:
		array = dconv.SliceAny(value)
	}
	for _, v := range array {
		that.data[v] = struct{}{}
	}
	return
}
