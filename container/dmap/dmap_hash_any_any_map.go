package dmap

import (
	"donkeygo/internal/empty"
	"donkeygo/internal/rwmutex"
	"donkeygo/util/dconv"
	"encoding/json"
	"math/rand"
)

// AnyAnyMap 并发安全的hash字典表
type AnyAnyMap struct {
	mu   rwmutex.RWMutex
	data map[interface{}]interface{}
}

// NewAnyAnyMap 创建hash表
func NewAnyAnyMap(safe ...bool) *AnyAnyMap {
	return &AnyAnyMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[interface{}]interface{}),
	}
}

// NewAnyAnyMapFrom 通过map创建hash表
func NewAnyAnyMapFrom(data map[interface{}]interface{}, safe ...bool) *AnyAnyMap {
	return &AnyAnyMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Iterator 迭代hash表
func (that *AnyAnyMap) Iterator(f func(k interface{}, v interface{}) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for k, v := range that.data {
		if !f(k, v) {
			break
		}
	}
}

// Map 返回原始结构
func (that *AnyAnyMap) Map() map[interface{}]interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if !that.mu.IsSafe() {
		return that.data
	}

	data := make(map[interface{}]interface{}, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// MapCopy copy一份map数据
func (that *AnyAnyMap) MapCopy() map[interface{}]interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	data := make(map[interface{}]interface{}, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// Clone clone一份数据。并返回新的对象指针
func (that *AnyAnyMap) Clone(safe ...bool) *AnyAnyMap {
	return NewFrom(that.MapCopy(), safe...)
}

// FilterEmpty 清除空值
func (that *AnyAnyMap) FilterEmpty() {
	that.mu.Lock()
	defer that.mu.Unlock()
	for k, v := range that.data {
		if empty.IsEmpty(v) {
			delete(that.data, k)
		}
	}
}

// FilterNil 清除nil值
func (that *AnyAnyMap) FilterNil() {
	that.mu.Lock()
	defer that.mu.Unlock()
	for k, v := range that.data {
		if empty.IsNil(v) {
			delete(that.data, k)
		}
	}
}

// Set 设置key value
func (that *AnyAnyMap) Set(key interface{}, value interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[interface{}]interface{})
	}
	that.data[key] = value
}

// Sets 批量设置
func (that *AnyAnyMap) Sets(data map[interface{}]interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = data
	} else {
		for k, v := range data {
			that.data[k] = v
		}
	}
}

// Search 通过key查找value
func (that *AnyAnyMap) Search(key interface{}) (value interface{}, found bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.data != nil {
		value, found = that.data[key]
	}
	return
}

// Get 查找key
func (that *AnyAnyMap) Get(key interface{}) (value interface{}) {
	v, _ := that.Search(key)
	return v
}

// Pop 随机获取一个字典key，value
func (that *AnyAnyMap) Pop() (key, value interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	for key, value = range that.data {
		delete(that.data, key)
		return
	}
	return
}

// Pops 随机获取指定size的字典值，-1表示获取全部
func (that *AnyAnyMap) Pops(size int) map[interface{}]interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if len(that.data) < size || size == -1 {
		size = len(that.data)
	}
	if size == 0 {
		return nil
	}
	var index = 0
	var newData = make(map[interface{}]interface{}, size)

	for key, value := range that.data {
		delete(that.data, key)
		newData[key] = value
		index++
		if index == size {
			break
		}
	}
	return newData
}

func (that *AnyAnyMap) Random() (key, value interface{}, exist bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	length := len(that.data)
	if length == 0 {
		return
	}
	i := rand.Intn(length)
	for key, value = range that.data {
		if i == 0 {
			exist = true
			return
		}
		i--
	}
	return
}

//加锁设置值，支持传入方法生成值，如果传入的key存在于map中，则直接返回已存在的值，不会使用传入的值
func (that *AnyAnyMap) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()

	if that.data == nil {
		that.data = make(map[interface{}]interface{})
	}
	//如果给定的key存在，则返回该key指向的值
	if v, ok := that.data[key]; ok {
		return v
	}
	//如果传入的值是方法，则执行该方法，并把返回值作为真正要使用的值
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	//如果值不为nil，则设置该值
	if value != nil {
		that.data[key] = value
	}
	//返回刚设置的值
	return value
}

// GetOrSet 查找key对应的值是否存在，如果未找到，则写入
func (that *AnyAnyMap) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc 获取指定key的值，如果不存在，则通过方法f生成该值，并写入到map，然后返回
//生成新值的方法未使用到锁，不会造成阻塞
func (that *AnyAnyMap) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock 与上面GetOrSetFunc方法的区别，在乎生成值的时候是否使用锁，此方法生成值的时候会使用到做，传入的方法不能阻塞，
//不然会造成map加锁不可读写
func (that *AnyAnyMap) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// SetIfNotExist 如果map中不存在该key，则设置
func (that *AnyAnyMap) SetIfNotExist(key interface{}, value interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc 如果map中不存在key，则调用方法f生成，生成的时候未加锁
func (that *AnyAnyMap) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock 如果map中不存在key，则调用方法f生成，生成的时候加锁
func (that *AnyAnyMap) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Contains 判断传入的key是否存在于map中
func (that *AnyAnyMap) Contains(key interface{}) bool {
	var ok bool
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.data != nil {
		_, ok = that.data[key]
	}
	return ok
}

// Remove 移除指定的key，并返回它对应的值
func (that *AnyAnyMap) Remove(key interface{}) (value interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data != nil {
		var ok bool
		if value, ok = that.data[key]; ok {
			delete(that.data, key)
			return value
		}
	}
	return
}

// Removes 批量删除key
func (that *AnyAnyMap) Removes(keys []interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data != nil {
		for _, key := range keys {
			delete(that.data, key)
		}
	}
}

// Keys 返回map的所有key
func (that *AnyAnyMap) Keys() []interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	var (
		keys  = make([]interface{}, len(that.data))
		index = 0
	)
	for key := range that.data {
		keys[index] = key
		index++
	}
	return keys
}

// Values 返回map所有的值
func (that *AnyAnyMap) Values() []interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	var (
		values = make([]interface{}, len(that.data))
		index  = 0
	)
	for _, value := range that.data {
		values[index] = value
		index++
	}
	return values
}

// Size 返回map的长度
func (that *AnyAnyMap) Size() int {
	that.mu.RLock()
	length := len(that.data)
	that.mu.RUnlock()
	return length
}

// IsEmpty 判断map是否为空
func (that *AnyAnyMap) IsEmpty() bool {
	return that.Size() == 0
}

// Clear 清除map
func (that *AnyAnyMap) Clear() {
	that.mu.Lock()
	that.data = make(map[interface{}]interface{})
	that.mu.Unlock()
}

// Replace 替换map
func (that *AnyAnyMap) Replace(data map[interface{}]interface{}) {
	that.mu.Lock()
	that.data = data
	that.mu.Unlock()
}

// LockFunc 加锁执行方法操作数据
func (that *AnyAnyMap) LockFunc(f func(m map[interface{}]interface{})) {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.data)
}

// RLockFunc 加读锁执行方法操作数据
func (that *AnyAnyMap) RLockFunc(f func(m map[interface{}]interface{})) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.data)
}

// Flip 翻转key和value
func (that *AnyAnyMap) Flip() {
	that.mu.Lock()
	defer that.mu.Unlock()
	n := make(map[interface{}]interface{}, len(that.data))
	for k, v := range that.data {
		n[v] = k
	}
	that.data = n
}

// Merge 合并两个map
func (that *AnyAnyMap) Merge(other *AnyAnyMap) {
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

//转换成字符串
func (that *AnyAnyMap) String() string {
	b, _ := that.MarshalJSON()
	return dconv.UnsafeBytesToStr(b)
}

// MarshalJSON json序列化
func (that *AnyAnyMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(dconv.Map(that.Map()))
}

// UnmarshalJSON json反序列化
func (that *AnyAnyMap) UnmarshalJSON(b []byte) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[interface{}]interface{})
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		that.data[k] = v
	}
	return nil
}

func (that *AnyAnyMap) UnmarshalValue(value interface{}) (err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[interface{}]interface{})
	}
	for k, v := range dconv.Map(value) {
		that.data[k] = v
	}
	return
}
