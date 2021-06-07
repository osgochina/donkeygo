// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package dmap

import (
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/util/dconv"
)

type IntStrMap struct {
	mu   rwmutex.RWMutex
	data map[int]string
}

// NewIntStrMap returns an empty IntStrMap object.
// The parameter <safe> is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewIntStrMap(safe ...bool) *IntStrMap {
	return &IntStrMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[int]string),
	}
}

// NewIntStrMapFrom creates and returns a hash map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewIntStrMapFrom(data map[int]string, safe ...bool) *IntStrMap {
	return &IntStrMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

// Iterator iterates the hash map readonly with custom callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (that *IntStrMap) Iterator(f func(k int, v string) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for k, v := range that.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (that *IntStrMap) Clone() *IntStrMap {
	return NewIntStrMapFrom(that.MapCopy(), that.mu.IsSafe())
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (that *IntStrMap) Map() map[int]string {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if !that.mu.IsSafe() {
		return that.data
	}
	data := make(map[int]string, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
func (that *IntStrMap) MapStrAny() map[string]interface{} {
	that.mu.RLock()
	data := make(map[string]interface{}, len(that.data))
	for k, v := range that.data {
		data[dconv.String(k)] = v
	}
	that.mu.RUnlock()
	return data
}

// MapCopy returns a copy of the underlying data of the hash map.
func (that *IntStrMap) MapCopy() map[int]string {
	that.mu.RLock()
	defer that.mu.RUnlock()
	data := make(map[int]string, len(that.data))
	for k, v := range that.data {
		data[k] = v
	}
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (that *IntStrMap) FilterEmpty() {
	that.mu.Lock()
	for k, v := range that.data {
		if empty.IsEmpty(v) {
			delete(that.data, k)
		}
	}
	that.mu.Unlock()
}

// Set sets key-value to the hash map.
func (that *IntStrMap) Set(key int, val string) {
	that.mu.Lock()
	if that.data == nil {
		that.data = make(map[int]string)
	}
	that.data[key] = val
	that.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (that *IntStrMap) Sets(data map[int]string) {
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
func (that *IntStrMap) Search(key int) (value string, found bool) {
	that.mu.RLock()
	if that.data != nil {
		value, found = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (that *IntStrMap) Get(key int) (value string) {
	that.mu.RLock()
	if that.data != nil {
		value, _ = that.data[key]
	}
	that.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (that *IntStrMap) Pop() (key int, value string) {
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
func (that *IntStrMap) Pops(size int) map[int]string {
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
		newMap = make(map[int]string, size)
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
func (that *IntStrMap) doSetWithLockCheck(key int, value string) string {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[int]string)
	}
	if v, ok := that.data[key]; ok {
		return v
	}
	that.data[key] = value
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (that *IntStrMap) GetOrSet(key int, value string) string {
	if v, ok := that.Search(key); !ok {
		return that.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist and returns this value.
func (that *IntStrMap) GetOrSetFunc(key int, f func() string) string {
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
func (that *IntStrMap) GetOrSetFuncLock(key int, f func() string) string {
	if v, ok := that.Search(key); !ok {
		that.mu.Lock()
		defer that.mu.Unlock()
		if that.data == nil {
			that.data = make(map[int]string)
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
func (that *IntStrMap) SetIfNotExist(key int, value string) bool {
	if !that.Contains(key) {
		that.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (that *IntStrMap) SetIfNotExistFunc(key int, f func() string) bool {
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
func (that *IntStrMap) SetIfNotExistFuncLock(key int, f func() string) bool {
	if !that.Contains(key) {
		that.mu.Lock()
		defer that.mu.Unlock()
		if that.data == nil {
			that.data = make(map[int]string)
		}
		if _, ok := that.data[key]; !ok {
			that.data[key] = f()
		}
		return true
	}
	return false
}

// Removes batch deletes values of the map by keys.
func (that *IntStrMap) Removes(keys []int) {
	that.mu.Lock()
	if that.data != nil {
		for _, key := range keys {
			delete(that.data, key)
		}
	}
	that.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (that *IntStrMap) Remove(key int) (value string) {
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
func (that *IntStrMap) Keys() []int {
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
func (that *IntStrMap) Values() []string {
	that.mu.RLock()
	var (
		values = make([]string, len(that.data))
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
func (that *IntStrMap) Contains(key int) bool {
	var ok bool
	that.mu.RLock()
	if that.data != nil {
		_, ok = that.data[key]
	}
	that.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (that *IntStrMap) Size() int {
	that.mu.RLock()
	length := len(that.data)
	that.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (that *IntStrMap) IsEmpty() bool {
	return that.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (that *IntStrMap) Clear() {
	that.mu.Lock()
	that.data = make(map[int]string)
	that.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (that *IntStrMap) Replace(data map[int]string) {
	that.mu.Lock()
	that.data = data
	that.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (that *IntStrMap) LockFunc(f func(m map[int]string)) {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (that *IntStrMap) RLockFunc(f func(m map[int]string)) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.data)
}

// Flip exchanges key-value of the map to value-key.
func (that *IntStrMap) Flip() {
	that.mu.Lock()
	defer that.mu.Unlock()
	n := make(map[int]string, len(that.data))
	for k, v := range that.data {
		n[dconv.Int(v)] = dconv.String(k)
	}
	that.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (that *IntStrMap) Merge(other *IntStrMap) {
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
func (that *IntStrMap) String() string {
	b, _ := that.MarshalJSON()
	return dconv.UnsafeBytesToStr(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (that *IntStrMap) MarshalJSON() ([]byte, error) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return json.Marshal(that.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *IntStrMap) UnmarshalJSON(b []byte) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[int]string)
	}
	if err := json.UnmarshalUseNumber(b, &that.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (that *IntStrMap) UnmarshalValue(value interface{}) (err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.data == nil {
		that.data = make(map[int]string)
	}
	switch value.(type) {
	case string, []byte:
		return json.UnmarshalUseNumber(dconv.Bytes(value), &that.data)
	default:
		for k, v := range dconv.Map(value) {
			that.data[dconv.Int(k)] = dconv.String(v)
		}
	}
	return
}
