// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package darray

import (
	"bytes"
	"fmt"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/util/dconv"
	"github.com/osgochina/donkeygo/util/drand"
	"math"
	"sort"
)

// SortedIntArray is a golang sorted int array with rich features.
// It is using increasing order in default, which can be changed by
// setting it a custom comparator.
// It contains a concurrent-safe/unsafe switch, which should be set
// when its initialization and cannot be changed then.
type SortedIntArray struct {
	mu         rwmutex.RWMutex
	array      []int
	unique     bool               // Whether enable unique feature(false)
	comparator func(a, b int) int // Comparison function(it returns -1: a < b; 0: a == b; 1: a > b)
}

// NewSortedIntArray creates and returns an empty sorted array.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArray(safe ...bool) *SortedIntArray {
	return NewSortedIntArraySize(0, safe...)
}

// NewSortedIntArrayComparator creates and returns an empty sorted array with specified comparator.
// The parameter <safe> is used to specify whether using array in concurrent-safety which is false in default.
func NewSortedIntArrayComparator(comparator func(a, b int) int, safe ...bool) *SortedIntArray {
	array := NewSortedIntArray(safe...)
	array.comparator = comparator
	return array
}

// NewSortedIntArraySize create and returns an sorted array with given size and cap.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArraySize(cap int, safe ...bool) *SortedIntArray {
	return &SortedIntArray{
		mu:         rwmutex.Create(safe...),
		array:      make([]int, 0, cap),
		comparator: defaultComparatorInt,
	}
}

// NewSortedIntArrayRange creates and returns a array by a range from <start> to <end>
// with step value <step>.
func NewSortedIntArrayRange(start, end, step int, safe ...bool) *SortedIntArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewSortedIntArrayFrom(slice, safe...)
}

// NewIntArrayFrom creates and returns an sorted array with given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArrayFrom(array []int, safe ...bool) *SortedIntArray {
	a := NewSortedIntArraySize(0, safe...)
	a.array = array
	sort.Ints(a.array)
	return a
}

// NewSortedIntArrayFromCopy creates and returns an sorted array from a copy of given slice <array>.
// The parameter <safe> is used to specify whether using array in concurrent-safety,
// which is false in default.
func NewSortedIntArrayFromCopy(array []int, safe ...bool) *SortedIntArray {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return NewSortedIntArrayFrom(newArray, safe...)
}

// SetArray sets the underlying slice array with the given <array>.
func (that *SortedIntArray) SetArray(array []int) *SortedIntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	that.array = array
	quickSortInt(that.array, that.getComparator())
	return that
}

// Sort sorts the array in increasing order.
// The parameter <reverse> controls whether sort
// in increasing order(default) or decreasing order.
func (that *SortedIntArray) Sort() *SortedIntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	quickSortInt(that.array, that.getComparator())
	return that
}

// Add adds one or multiple values to sorted array, the array always keeps sorted.
// It's alias of function Append, see Append.
func (that *SortedIntArray) Add(values ...int) *SortedIntArray {
	return that.Append(values...)
}

// Append adds one or multiple values to sorted array, the array always keeps sorted.
func (that *SortedIntArray) Append(values ...int) *SortedIntArray {
	if len(values) == 0 {
		return that
	}
	that.mu.Lock()
	defer that.mu.Unlock()
	for _, value := range values {
		index, cmp := that.binSearch(value, false)
		if that.unique && cmp == 0 {
			continue
		}
		if index < 0 {
			that.array = append(that.array, value)
			continue
		}
		if cmp > 0 {
			index++
		}
		rear := append([]int{}, that.array[index:]...)
		that.array = append(that.array[0:index], value)
		that.array = append(that.array, rear...)
	}
	return that
}

// Get returns the value by the specified index.
// If the given <index> is out of range of the array, the <found> is false.
func (that *SortedIntArray) Get(index int) (value int, found bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if index < 0 || index >= len(that.array) {
		return 0, false
	}
	return that.array[index], true
}

// Remove removes an item by index.
// If the given <index> is out of range of the array, the <found> is false.
func (that *SortedIntArray) Remove(index int) (value int, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.doRemoveWithoutLock(index)
}

// doRemoveWithoutLock removes an item by index without lock.
func (that *SortedIntArray) doRemoveWithoutLock(index int) (value int, found bool) {
	if index < 0 || index >= len(that.array) {
		return 0, false
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := that.array[0]
		that.array = that.array[1:]
		return value, true
	} else if index == len(that.array)-1 {
		value := that.array[index]
		that.array = that.array[:index]
		return value, true
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	value = that.array[index]
	that.array = append(that.array[:index], that.array[index+1:]...)
	return value, true
}

// RemoveValue removes an item by value.
// It returns true if value is found in the array, or else false if not found.
func (that *SortedIntArray) RemoveValue(value int) bool {
	if i := that.Search(value); i != -1 {
		_, found := that.Remove(i)
		return found
	}
	return false
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the <found> is false.
func (that *SortedIntArray) PopLeft() (value int, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if len(that.array) == 0 {
		return 0, false
	}
	value = that.array[0]
	that.array = that.array[1:]
	return value, true
}

// PopRight pops and returns an item from the end of array.
// Note that if the array is empty, the <found> is false.
func (that *SortedIntArray) PopRight() (value int, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	index := len(that.array) - 1
	if index < 0 {
		return 0, false
	}
	value = that.array[index]
	that.array = that.array[:index]
	return value, true
}

// PopRand randomly pops and return an item out of array.
// Note that if the array is empty, the <found> is false.
func (that *SortedIntArray) PopRand() (value int, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.doRemoveWithoutLock(drand.Intn(len(that.array)))
}

// PopRands randomly pops and returns <size> items out of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (that *SortedIntArray) PopRands(size int) []int {
	that.mu.Lock()
	defer that.mu.Unlock()
	if size <= 0 || len(that.array) == 0 {
		return nil
	}
	if size >= len(that.array) {
		size = len(that.array)
	}
	array := make([]int, size)
	for i := 0; i < size; i++ {
		array[i], _ = that.doRemoveWithoutLock(drand.Intn(len(that.array)))
	}
	return array
}

// PopLefts pops and returns <size> items from the beginning of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (that *SortedIntArray) PopLefts(size int) []int {
	that.mu.Lock()
	defer that.mu.Unlock()
	if size <= 0 || len(that.array) == 0 {
		return nil
	}
	if size >= len(that.array) {
		array := that.array
		that.array = that.array[:0]
		return array
	}
	value := that.array[0:size]
	that.array = that.array[size:]
	return value
}

// PopRights pops and returns <size> items from the end of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (that *SortedIntArray) PopRights(size int) []int {
	that.mu.Lock()
	defer that.mu.Unlock()
	if size <= 0 || len(that.array) == 0 {
		return nil
	}
	index := len(that.array) - size
	if index <= 0 {
		array := that.array
		that.array = that.array[:0]
		return array
	}
	value := that.array[index:]
	that.array = that.array[:index]
	return value
}

// Range picks and returns items by range, like array[start:end].
// Notice, if in concurrent-safe usage, it returns a copy of slice;
// else a pointer to the underlying data.
//
// If <end> is negative, then the offset will start from the end of array.
// If <end> is omitted, then the sequence will have everything from start up
// until the end of the array.
func (that *SortedIntArray) Range(start int, end ...int) []int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	offsetEnd := len(that.array)
	if len(end) > 0 && end[0] < offsetEnd {
		offsetEnd = end[0]
	}
	if start > offsetEnd {
		return nil
	}
	if start < 0 {
		start = 0
	}
	array := ([]int)(nil)
	if that.mu.IsSafe() {
		array = make([]int, offsetEnd-start)
		copy(array, that.array[start:offsetEnd])
	} else {
		array = that.array[start:offsetEnd]
	}
	return array
}

// SubSlice returns a slice of elements from the array as specified
// by the <offset> and <size> parameters.
// If in concurrent safe usage, it returns a copy of the slice; else a pointer.
//
// If offset is non-negative, the sequence will start at that offset in the array.
// If offset is negative, the sequence will start that far from the end of the array.
//
// If length is given and is positive, then the sequence will have up to that many elements in it.
// If the array is shorter than the length, then only the available array elements will be present.
// If length is given and is negative then the sequence will stop that many elements from the end of the array.
// If it is omitted, then the sequence will have everything from offset up until the end of the array.
//
// Any possibility crossing the left border of array, it will fail.
func (that *SortedIntArray) SubSlice(offset int, length ...int) []int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	size := len(that.array)
	if len(length) > 0 {
		size = length[0]
	}
	if offset > len(that.array) {
		return nil
	}
	if offset < 0 {
		offset = len(that.array) + offset
		if offset < 0 {
			return nil
		}
	}
	if size < 0 {
		offset += size
		size = -size
		if offset < 0 {
			return nil
		}
	}
	end := offset + size
	if end > len(that.array) {
		end = len(that.array)
		size = len(that.array) - offset
	}
	if that.mu.IsSafe() {
		s := make([]int, size)
		copy(s, that.array[offset:])
		return s
	} else {
		return that.array[offset:end]
	}
}

// Len returns the length of array.
func (that *SortedIntArray) Len() int {
	that.mu.RLock()
	length := len(that.array)
	that.mu.RUnlock()
	return length
}

// Sum returns the sum of values in an array.
func (that *SortedIntArray) Sum() (sum int) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, v := range that.array {
		sum += v
	}
	return
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (that *SortedIntArray) Slice() []int {
	array := ([]int)(nil)
	if that.mu.IsSafe() {
		that.mu.RLock()
		defer that.mu.RUnlock()
		array = make([]int, len(that.array))
		copy(array, that.array)
	} else {
		array = that.array
	}
	return array
}

// Interfaces returns current array as []interface{}.
func (that *SortedIntArray) Interfaces() []interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	array := make([]interface{}, len(that.array))
	for k, v := range that.array {
		array[k] = v
	}
	return array
}

// Contains checks whether a value exists in the array.
func (that *SortedIntArray) Contains(value int) bool {
	return that.Search(value) != -1
}

// Search searches array by <value>, returns the index of <value>,
// or returns -1 if not exists.
func (that *SortedIntArray) Search(value int) (index int) {
	if i, r := that.binSearch(value, true); r == 0 {
		return i
	}
	return -1
}

// Binary search.
// It returns the last compared index and the result.
// If <result> equals to 0, it means the value at <index> is equals to <value>.
// If <result> lesser than 0, it means the value at <index> is lesser than <value>.
// If <result> greater than 0, it means the value at <index> is greater than <value>.
func (that *SortedIntArray) binSearch(value int, lock bool) (index int, result int) {
	if lock {
		that.mu.RLock()
		defer that.mu.RUnlock()
	}
	if len(that.array) == 0 {
		return -1, -2
	}
	min := 0
	max := len(that.array) - 1
	mid := 0
	cmp := -2
	for min <= max {
		mid = min + int((max-min)/2)
		cmp = that.getComparator()(value, that.array[mid])
		switch {
		case cmp < 0:
			max = mid - 1
		case cmp > 0:
			min = mid + 1
		default:
			return mid, cmp
		}
	}
	return mid, cmp
}

// SetUnique sets unique mark to the array,
// which means it does not contain any repeated items.
// It also do unique check, remove all repeated items.
func (that *SortedIntArray) SetUnique(unique bool) *SortedIntArray {
	oldUnique := that.unique
	that.unique = unique
	if unique && oldUnique != unique {
		that.Unique()
	}
	return that
}

// Unique uniques the array, clear repeated items.
func (that *SortedIntArray) Unique() *SortedIntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	if len(that.array) == 0 {
		return that
	}
	i := 0
	for {
		if i == len(that.array)-1 {
			break
		}
		if that.getComparator()(that.array[i], that.array[i+1]) == 0 {
			that.array = append(that.array[:i+1], that.array[i+1+1:]...)
		} else {
			i++
		}
	}
	return that
}

// Clone returns a new array, which is a copy of current array.
func (that *SortedIntArray) Clone() (newArray *SortedIntArray) {
	that.mu.RLock()
	array := make([]int, len(that.array))
	copy(array, that.array)
	that.mu.RUnlock()
	return NewSortedIntArrayFrom(array, that.mu.IsSafe())
}

// Clear deletes all items of current array.
func (that *SortedIntArray) Clear() *SortedIntArray {
	that.mu.Lock()
	if len(that.array) > 0 {
		that.array = make([]int, 0)
	}
	that.mu.Unlock()
	return that
}

// LockFunc locks writing by callback function <f>.
func (that *SortedIntArray) LockFunc(f func(array []int)) *SortedIntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.array)
	return that
}

// RLockFunc locks reading by callback function <f>.
func (that *SortedIntArray) RLockFunc(f func(array []int)) *SortedIntArray {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.array)
	return that
}

// Merge merges <array> into current array.
// The parameter <array> can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (that *SortedIntArray) Merge(array interface{}) *SortedIntArray {
	return that.Add(dconv.Ints(array)...)
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
func (that *SortedIntArray) Chunk(size int) [][]int {
	if size < 1 {
		return nil
	}
	that.mu.RLock()
	defer that.mu.RUnlock()
	length := len(that.array)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]int
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, that.array[i*size:end])
		i++
	}
	return n
}

// Rand randomly returns one item from array(no deleting).
func (that *SortedIntArray) Rand() (value int, found bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if len(that.array) == 0 {
		return 0, false
	}
	return that.array[drand.Intn(len(that.array))], true
}

// Rands randomly returns <size> items from array(no deleting).
func (that *SortedIntArray) Rands(size int) []int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if size <= 0 || len(that.array) == 0 {
		return nil
	}
	array := make([]int, size)
	for i := 0; i < size; i++ {
		array[i] = that.array[drand.Intn(len(that.array))]
	}
	return array
}

// Join joins array elements with a string <glue>.
func (that *SortedIntArray) Join(glue string) string {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if len(that.array) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range that.array {
		buffer.WriteString(dconv.String(v))
		if k != len(that.array)-1 {
			buffer.WriteString(glue)
		}
	}
	return buffer.String()
}

// CountValues counts the number of occurrences of all values in the array.
func (that *SortedIntArray) CountValues() map[int]int {
	m := make(map[int]int)
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, v := range that.array {
		m[v]++
	}
	return m
}

// Iterator is alias of IteratorAsc.
func (that *SortedIntArray) Iterator(f func(k int, v int) bool) {
	that.IteratorAsc(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (that *SortedIntArray) IteratorAsc(f func(k int, v int) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for k, v := range that.array {
		if !f(k, v) {
			break
		}
	}
}

// IteratorDesc iterates the array readonly in descending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (that *SortedIntArray) IteratorDesc(f func(k int, v int) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for i := len(that.array) - 1; i >= 0; i-- {
		if !f(i, that.array[i]) {
			break
		}
	}
}

// String returns current array as a string, which implements like json.Marshal does.
func (that *SortedIntArray) String() string {
	return "[" + that.Join(",") + "]"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (that SortedIntArray) MarshalJSON() ([]byte, error) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return json.Marshal(that.array)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *SortedIntArray) UnmarshalJSON(b []byte) error {
	if that.comparator == nil {
		that.array = make([]int, 0)
		that.comparator = defaultComparatorInt
	}
	that.mu.Lock()
	defer that.mu.Unlock()
	if err := json.UnmarshalUseNumber(b, &that.array); err != nil {
		return err
	}
	if that.array != nil {
		sort.Ints(that.array)
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (that *SortedIntArray) UnmarshalValue(value interface{}) (err error) {
	if that.comparator == nil {
		that.comparator = defaultComparatorInt
	}
	that.mu.Lock()
	defer that.mu.Unlock()
	switch value.(type) {
	case string, []byte:
		err = json.UnmarshalUseNumber(dconv.Bytes(value), &that.array)
	default:
		that.array = dconv.SliceInt(value)
	}
	if that.array != nil {
		sort.Ints(that.array)
	}
	return err
}

// FilterEmpty removes all zero value of the array.
func (that *SortedIntArray) FilterEmpty() *SortedIntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	for i := 0; i < len(that.array); {
		if that.array[i] == 0 {
			that.array = append(that.array[:i], that.array[i+1:]...)
		} else {
			break
		}
	}
	for i := len(that.array) - 1; i >= 0; {
		if that.array[i] == 0 {
			that.array = append(that.array[:i], that.array[i+1:]...)
		} else {
			break
		}
	}
	return that
}

// Walk applies a user supplied function <f> to every item of array.
func (that *SortedIntArray) Walk(f func(value int) int) *SortedIntArray {
	that.mu.Lock()
	defer that.mu.Unlock()

	// Keep the array always sorted.
	defer quickSortInt(that.array, that.getComparator())

	for i, v := range that.array {
		that.array[i] = f(v)
	}
	return that
}

// IsEmpty checks whether the array is empty.
func (that *SortedIntArray) IsEmpty() bool {
	return that.Len() == 0
}

// getComparator returns the comparator if it's previously set,
// or else it returns a default comparator.
func (that *SortedIntArray) getComparator() func(a, b int) int {
	if that.comparator == nil {
		return defaultComparatorInt
	}
	return that.comparator
}
