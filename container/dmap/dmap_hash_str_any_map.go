package dmap

import (
	"encoding/json"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/util/dconv"
)

type StrAnyMap struct {
	mu   rwmutex.RWMutex
	data map[string]interface{}
}

// NewStrAnyMap 返回一个空的StrAnyMap对象
func NewStrAnyMap(safe ...bool) *StrAnyMap {
	return &StrAnyMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[string]interface{}),
	}
}

// NewStrAnyMapFrom 从给定的map 创建并返回一个哈希映射。
//注意，param map将被设置为底层数据map(没有深层复制)，
//当改变外部映射时，可能会有一些并发安全问题。
func NewStrAnyMapFrom(data map[string]interface{}, safe ...bool) *StrAnyMap {
	return &StrAnyMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Iterator 迭代该map，如果其中有一次迭代返回false，则整个迭代过程将会终止
func (that *StrAnyMap) Iterator(f func(k string, v interface{}) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for k, v := range that.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone 返回一个新的StrAnyMap对象，其中的值是赋值的当前StrAnyMap对象
func (that *StrAnyMap) Clone() *StrAnyMap {
	return NewStrAnyMapFrom(that.MapCopy(), that.mu.IsSafe())
}

// Map 返回一个原始的map对象
func (that *StrAnyMap) Map() map[string]interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if !that.mu.IsSafe() {
		return that.data
	}
	data := make(map[string]interface{}, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// MapStrAny 返回一个map[string]interface{}类型的map对象
func (that *StrAnyMap) MapStrAny() map[string]interface{} {
	return that.Map()
}

// MapCopy 把对象的值赋值到新的地址中并返回
func (that *StrAnyMap) MapCopy() map[string]interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	data := make(map[string]interface{}, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// FilterEmpty 删除map对象中的空值
// 空值的定义: 0, nil, false, "", len(slice/map/chan) == 0
func (that *StrAnyMap) FilterEmpty() {
	that.mu.Lock()
	for k, v := range that.data {
		if empty.IsEmpty(v) {
			delete(that.data, k)
		}
	}
	that.mu.Unlock()
}

// FilterNil 删除map中值为 nil的对象
func (that *StrAnyMap) FilterNil() {
	that.mu.Lock()
	defer that.mu.Unlock()
	for k, v := range that.data {
		if empty.IsNil(v) {
			delete(that.data, k)
		}
	}
}

// Set 写入key，val到map中
func (that *StrAnyMap) Set(key string, val interface{}) {
	that.mu.Lock()
	if that.data == nil {
		that.data = make(map[string]interface{})
	}
	that.data[key] = val
	that.mu.Unlock()
}

// Sets 批量写入
func (that *StrAnyMap) Sets(data map[string]interface{}) {
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

// Search 查找key是否存在map中，知道则返回该值，found为true，没有找到value=nil，found=false
func (that *StrAnyMap) Search(key string) (value interface{}, found bool) {
	that.mu.RLock()
	if that.data != nil {
		value, found = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Get 返回key对应的值
func (that *StrAnyMap) Get(key string) (value interface{}) {
	that.mu.RLock()
	if that.data != nil {
		value, _ = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Pop 从map中返回一个元素并删除它
func (that *StrAnyMap) Pop() (key string, value interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	for key, value = range that.data {
		delete(that.data, key)
		return
	}
	return
}

// Pops 返回指定数量的元素并删除它
func (that *StrAnyMap) Pops(size int) map[string]interface{} {
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
		newMap = make(map[string]interface{}, size)
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

// doSetWithLockCheck 检查该键值是否存在,
//如果不存在，设置map的值为func的返回值，
//否则返回现有的值。
func (that *StrAnyMap) doSetWithLockCheck(key string, value interface{}) interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[string]interface{})
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

// GetOrSet 如果存在，则返回，不存在则设置并返回
func (that *StrAnyMap) GetOrSet(key string, value interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc 如果存在，则返回，不存在则设置并返回
func (that *StrAnyMap) GetOrSetFunc(key string, f func() interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock 如果存在，则返回，不存在则设置并返回
func (that *StrAnyMap) GetOrSetFuncLock(key string, f func() interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// GetVar 返回Var类型的值
func (that *StrAnyMap) GetVar(key string) *dvar.Var {
	return dvar.New(that.Get(key))
}

// GetVarOrSet returns a Var with result from GetVarOrSet.
// The returned Var is un-concurrent safe.
func (that *StrAnyMap) GetVarOrSet(key string, value interface{}) *dvar.Var {
	return dvar.New(that.GetOrSet(key, value))
}

// GetVarOrSetFunc returns a Var with result from GetOrSetFunc.
// The returned Var is un-concurrent safe.
func (that *StrAnyMap) GetVarOrSetFunc(key string, f func() interface{}) *dvar.Var {
	return dvar.New(that.GetOrSetFunc(key, f))
}

// GetVarOrSetFuncLock returns a Var with result from GetOrSetFuncLock.
// The returned Var is un-concurrent safe.
func (that *StrAnyMap) GetVarOrSetFuncLock(key string, f func() interface{}) *dvar.Var {
	return dvar.New(that.GetOrSetFuncLock(key, f))
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (that *StrAnyMap) SetIfNotExist(key string, value interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (that *StrAnyMap) SetIfNotExistFunc(key string, f func() interface{}) bool {
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
func (that *StrAnyMap) SetIfNotExistFuncLock(key string, f func() interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Removes batch deletes values of the map by keys.
func (that *StrAnyMap) Removes(keys []string) {
	that.mu.Lock()
	if that.data != nil {
		for _, key := range keys {
			delete(that.data, key)
		}
	}
	that.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (that *StrAnyMap) Remove(key string) (value interface{}) {
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
func (that *StrAnyMap) Keys() []string {
	that.mu.RLock()
	var (
		keys  = make([]string, len(that.data))
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
func (that *StrAnyMap) Values() []interface{} {
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
func (that *StrAnyMap) Contains(key string) bool {
	var ok bool
	that.mu.RLock()
	if that.data != nil {
		_, ok = that.data[key]
	}
	that.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (that *StrAnyMap) Size() int {
	that.mu.RLock()
	length := len(that.data)
	that.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (that *StrAnyMap) IsEmpty() bool {
	return that.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (that *StrAnyMap) Clear() {
	that.mu.Lock()
	that.data = make(map[string]interface{})
	that.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (that *StrAnyMap) Replace(data map[string]interface{}) {
	that.mu.Lock()
	that.data = data
	that.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (that *StrAnyMap) LockFunc(f func(m map[string]interface{})) {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (that *StrAnyMap) RLockFunc(f func(m map[string]interface{})) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.data)
}

// Flip exchanges key-value of the map to value-key.
func (that *StrAnyMap) Flip() {
	that.mu.Lock()
	defer that.mu.Unlock()
	n := make(map[string]interface{}, len(that.data))
	for k, v := range that.data {
		n[dconv.String(v)] = k
	}
	that.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (that *StrAnyMap) Merge(other *StrAnyMap) {
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
func (that *StrAnyMap) String() string {
	b, _ := that.MarshalJSON()
	return dconv.UnsafeBytesToStr(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (that *StrAnyMap) MarshalJSON() ([]byte, error) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return json.Marshal(that.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *StrAnyMap) UnmarshalJSON(b []byte) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[string]interface{})
	}
	if err := json.Unmarshal(b, &that.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (that *StrAnyMap) UnmarshalValue(value interface{}) (err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	that.data = dconv.Map(value)
	return
}
