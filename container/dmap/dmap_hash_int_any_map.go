package dmap

import (
	"encoding/json"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/util/dconv"
)

type IntAnyMap struct {
	mu   rwmutex.RWMutex
	data map[int]interface{}
}

// NewIntAnyMap 创建key值为int类型的hash表
func NewIntAnyMap(safe ...bool) *IntAnyMap {
	return &IntAnyMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[int]interface{}),
	}
}

// NewIntAnyMapFrom 通过基础数据格式创建map
func NewIntAnyMapFrom(data map[int]interface{}, safe ...bool) *IntAnyMap {
	return &IntAnyMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Iterator 迭代map
func (that *IntAnyMap) Iterator(f func(k int, v interface{}) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for k, v := range that.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone 返回一个新的map对象，内容是当前map对象的副本
func (that *IntAnyMap) Clone() *IntAnyMap {
	return NewIntAnyMapFrom(that.MapCopy(), that.mu.IsSafe())
}

// Map 返回基础map类型
func (that *IntAnyMap) Map() map[int]interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if !that.mu.IsSafe() {
		return that.data
	}
	data := make(map[int]interface{}, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// MapStrAny 把key转换成字符串返回
func (that *IntAnyMap) MapStrAny() map[string]interface{} {
	that.mu.RLock()
	data := make(map[string]interface{}, len(that.data))
	for k, v := range that.data {
		data[dconv.String(k)] = v
	}
	that.mu.RUnlock()
	return data
}

// MapCopy 复制一份数据返回
func (that *IntAnyMap) MapCopy() map[int]interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	data := make(map[int]interface{}, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// FilterEmpty 删除空值
func (that *IntAnyMap) FilterEmpty() {
	that.mu.Lock()
	for k, v := range that.data {
		if empty.IsEmpty(v) {
			delete(that.data, k)
		}
	}
	that.mu.Unlock()
}

// FilterNil 删除nil值
func (that *IntAnyMap) FilterNil() {
	that.mu.Lock()
	defer that.mu.Unlock()
	for k, v := range that.data {
		if empty.IsNil(v) {
			delete(that.data, k)
		}
	}
}

// Set 写入单个数据
func (that *IntAnyMap) Set(key int, val interface{}) {
	that.mu.Lock()
	if that.data == nil {
		that.data = make(map[int]interface{})
	}
	that.data[key] = val
	that.mu.Unlock()
}

// Sets 写入数组
func (that *IntAnyMap) Sets(data map[int]interface{}) {
	that.mu.Lock()
	if that.data == nil {
		that.data = data
	} else {
		for k, v := range data {
			that.data[k] = v
		}
	}
	that.mu.Unlock()
}

// Search 搜索数据
func (that *IntAnyMap) Search(key int) (value interface{}, found bool) {
	that.mu.RLock()
	if that.data != nil {
		value, found = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (that *IntAnyMap) Get(key int) (value interface{}) {
	that.mu.RLock()
	if that.data != nil {
		value, _ = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (that *IntAnyMap) Pop() (key int, value interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	for key, value = range that.data {
		delete(that.data, key)
		return
	}
	return
}

// Pops retrieves and deletes <size> items from the map.
// It returns all items if size == -1.
func (that *IntAnyMap) Pops(size int) map[int]interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if size > len(that.data) || size == -1 {
		size = len(that.data)
	}
	if size == 0 {
		return nil
	}
	var (
		index  = 0
		newMap = make(map[int]interface{}, size)
	)
	for k, v := range that.data {
		delete(that.data, k)
		newMap[k] = v
		index++
		if index == size {
			break
		}
	}
	return newMap
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// When setting value, if <value> is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with <key>.
//
// It returns value with given <key>.
func (that *IntAnyMap) doSetWithLockCheck(key int, value interface{}) interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[int]interface{})
	}
	if v, ok := that.data[key]; ok {
		return v
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		that.data[key] = value
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (that *IntAnyMap) GetOrSet(key int, value interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist and returns this value.
func (that *IntAnyMap) GetOrSetFunc(key int, f func() interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist and returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func (that *IntAnyMap) GetOrSetFuncLock(key int, f func() interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// GetVar returns a Var with the value by given <key>.
// The returned Var is un-concurrent safe.
func (that *IntAnyMap) GetVar(key int) *dvar.Var {
	return dvar.New(that.Get(key))
}

// GetVarOrSet returns a Var with result from GetVarOrSet.
// The returned Var is un-concurrent safe.
func (that *IntAnyMap) GetVarOrSet(key int, value interface{}) *dvar.Var {
	return dvar.New(that.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a Var with result from GetOrSetFunc.
// The returned Var is un-concurrent safe.
func (that *IntAnyMap) GetVarOrSetFunc(key int, f func() interface{}) *dvar.Var {
	return dvar.New(that.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a Var with result from GetOrSetFuncLock.
// The returned Var is un-concurrent safe.
func (that *IntAnyMap) GetVarOrSetFuncLock(key int, f func() interface{}) *dvar.Var {
	return dvar.New(that.GetOrSetFuncLock(key, f))
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (that *IntAnyMap) SetIfNotExist(key int, value interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (that *IntAnyMap) SetIfNotExistFunc(key int, f func() interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function <f> with mutex.Lock of the hash map.
func (that *IntAnyMap) SetIfNotExistFuncLock(key int, f func() interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Removes batch deletes values of the map by keys.
func (that *IntAnyMap) Removes(keys []int) {
	that.mu.Lock()
	if that.data != nil {
		for _, key := range keys {
			delete(that.data, key)
		}
	}
	that.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (that *IntAnyMap) Remove(key int) (value interface{}) {
	that.mu.Lock()
	if that.data != nil {
		var ok bool
		if value, ok = that.data[key]; ok {
			delete(that.data, key)
		}
	}
	that.mu.Unlock()
	return
}

// Keys returns all keys of the map as a slice.
func (that *IntAnyMap) Keys() []int {
	that.mu.RLock()
	var (
		keys  = make([]int, len(that.data))
		index = 0
	)
	for key := range that.data {
		keys[index] = key
		index++
	}
	that.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (that *IntAnyMap) Values() []interface{} {
	that.mu.RLock()
	var (
		values = make([]interface{}, len(that.data))
		index  = 0
	)
	for _, value := range that.data {
		values[index] = value
		index++
	}
	that.mu.RUnlock()
	return values
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (that *IntAnyMap) Contains(key int) bool {
	var ok bool
	that.mu.RLock()
	if that.data != nil {
		_, ok = that.data[key]
	}
	that.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (that *IntAnyMap) Size() int {
	that.mu.RLock()
	length := len(that.data)
	that.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (that *IntAnyMap) IsEmpty() bool {
	return that.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (that *IntAnyMap) Clear() {
	that.mu.Lock()
	that.data = make(map[int]interface{})
	that.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (that *IntAnyMap) Replace(data map[int]interface{}) {
	that.mu.Lock()
	that.data = data
	that.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (that *IntAnyMap) LockFunc(f func(m map[int]interface{})) {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (that *IntAnyMap) RLockFunc(f func(m map[int]interface{})) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.data)
}

// Flip exchanges key-value of the map to value-key.
func (that *IntAnyMap) Flip() {
	that.mu.Lock()
	defer that.mu.Unlock()
	n := make(map[int]interface{}, len(that.data))
	for k, v := range that.data {
		n[dconv.Int(v)] = k
	}
	that.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (that *IntAnyMap) Merge(other *IntAnyMap) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = other.MapCopy()
		return
	}
	if other != that {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range other.data {
		that.data[k] = v
	}
}

// String returns the map as a string.
func (that *IntAnyMap) String() string {
	b, _ := that.MarshalJSON()
	return dconv.UnsafeBytesToStr(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (that *IntAnyMap) MarshalJSON() ([]byte, error) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return json.Marshal(that.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *IntAnyMap) UnmarshalJSON(b []byte) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[int]interface{})
	}
	if err := json.Unmarshal(b, &that.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (that *IntAnyMap) UnmarshalValue(value interface{}) (err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[int]interface{})
	}
	switch value.(type) {
	case string, []byte:
		return json.Unmarshal(dconv.Bytes(value), &that.data)
	default:
		for k, v := range dconv.Map(value) {
			that.data[dconv.Int(k)] = v
		}
	}
	return
}
