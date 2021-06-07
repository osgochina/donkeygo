// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.
//

package dmap

import (
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/util/dconv"
)

type StrIntMap struct {
	mu   rwmutex.RWMutex
	data map[string]int
}

// NewStrIntMap returns an empty StrIntMap object.
// The parameter <safe> is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewStrIntMap(safe ...bool) *StrIntMap {
	return &StrIntMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[string]int),
	}
}

// NewStrIntMapFrom creates and returns a hash map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewStrIntMapFrom(data map[string]int, safe ...bool) *StrIntMap {
	return &StrIntMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Iterator iterates the hash map readonly with custom callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (that *StrIntMap) Iterator(f func(k string, v int) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for k, v := range that.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (that *StrIntMap) Clone() *StrIntMap {
	return NewStrIntMapFrom(that.MapCopy(), that.mu.IsSafe())
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (that *StrIntMap) Map() map[string]int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if !that.mu.IsSafe() {
		return that.data
	}
	data := make(map[string]int, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
func (that *StrIntMap) MapStrAny() map[string]interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	data := make(map[string]interface{}, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// MapCopy returns a copy of the underlying data of the hash map.
func (that *StrIntMap) MapCopy() map[string]int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	data := make(map[string]int, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (that *StrIntMap) FilterEmpty() {
	that.mu.Lock()
	for k, v := range that.data {
		if empty.IsEmpty(v) {
			delete(that.data, k)
		}
	}
	that.mu.Unlock()
}

// Set sets key-value to the hash map.
func (that *StrIntMap) Set(key string, val int) {
	that.mu.Lock()
	if that.data == nil {
		that.data = make(map[string]int)
	}
	that.data[key] = val
	that.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (that *StrIntMap) Sets(data map[string]int) {
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

// Search searches the map with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (that *StrIntMap) Search(key string) (value int, found bool) {
	that.mu.RLock()
	if that.data != nil {
		value, found = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (that *StrIntMap) Get(key string) (value int) {
	that.mu.RLock()
	if that.data != nil {
		value, _ = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (that *StrIntMap) Pop() (key string, value int) {
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
func (that *StrIntMap) Pops(size int) map[string]int {
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
		newMap = make(map[string]int, size)
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
// It returns value with given <key>.
func (that *StrIntMap) doSetWithLockCheck(key string, value int) int {
	that.mu.Lock()
	if that.data == nil {
		that.data = make(map[string]int)
	}
	if v, ok := that.data[key]; ok {
		that.mu.Unlock()
		return v
	}
	that.data[key] = value
	that.mu.Unlock()
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (that *StrIntMap) GetOrSet(key string, value int) int {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (that *StrIntMap) GetOrSetFunc(key string, f func() int) int {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func (that *StrIntMap) GetOrSetFuncLock(key string, f func() int) int {
	if v, ok := that.Search(key); !ok {
		that.mu.Lock()
		defer that.mu.Unlock()
		if that.data == nil {
			that.data = make(map[string]int)
		}
		if v, ok = that.data[key]; ok {
			return v
		}
		v = f()
		that.data[key] = v
		return v
	} else {
		return v
	}
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (that *StrIntMap) SetIfNotExist(key string, value int) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (that *StrIntMap) SetIfNotExistFunc(key string, f func() int) bool {
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
func (that *StrIntMap) SetIfNotExistFuncLock(key string, f func() int) bool {
	if !that.Contains(key) {
		that.mu.Lock()
		defer that.mu.Unlock()
		if that.data == nil {
			that.data = make(map[string]int)
		}
		if _, ok := that.data[key]; !ok {
			that.data[key] = f()
		}
		return true
	}
	return false
}

// Removes batch deletes values of the map by keys.
func (that *StrIntMap) Removes(keys []string) {
	that.mu.Lock()
	if that.data != nil {
		for _, key := range keys {
			delete(that.data, key)
		}
	}
	that.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (that *StrIntMap) Remove(key string) (value int) {
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
func (that *StrIntMap) Keys() []string {
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
func (that *StrIntMap) Values() []int {
	that.mu.RLock()
	var (
		values = make([]int, len(that.data))
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
func (that *StrIntMap) Contains(key string) bool {
	var ok bool
	that.mu.RLock()
	if that.data != nil {
		_, ok = that.data[key]
	}
	that.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (that *StrIntMap) Size() int {
	that.mu.RLock()
	length := len(that.data)
	that.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (that *StrIntMap) IsEmpty() bool {
	return that.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (that *StrIntMap) Clear() {
	that.mu.Lock()
	that.data = make(map[string]int)
	that.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (that *StrIntMap) Replace(data map[string]int) {
	that.mu.Lock()
	that.data = data
	that.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (that *StrIntMap) LockFunc(f func(m map[string]int)) {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (that *StrIntMap) RLockFunc(f func(m map[string]int)) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.data)
}

// Flip exchanges key-value of the map to value-key.
func (that *StrIntMap) Flip() {
	that.mu.Lock()
	defer that.mu.Unlock()
	n := make(map[string]int, len(that.data))
	for k, v := range that.data {
		n[dconv.String(v)] = dconv.Int(k)
	}
	that.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (that *StrIntMap) Merge(other *StrIntMap) {
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
func (that *StrIntMap) String() string {
	b, _ := that.MarshalJSON()
	return dconv.UnsafeBytesToStr(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (that *StrIntMap) MarshalJSON() ([]byte, error) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return json.Marshal(that.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *StrIntMap) UnmarshalJSON(b []byte) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[string]int)
	}
	if err := json.UnmarshalUseNumber(b, &that.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (that *StrIntMap) UnmarshalValue(value interface{}) (err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[string]int)
	}
	switch value.(type) {
	case string, []byte:
		return json.UnmarshalUseNumber(dconv.Bytes(value), &that.data)
	default:
		for k, v := range dconv.Map(value) {
			that.data[k] = dconv.Int(v)
		}
	}
	return
}
