package darray

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/util/dconv"
	"github.com/osgochina/donkeygo/util/drand"
	"math"
	"sort"
)

// IntArray int类型的数组
type IntArray struct {
	mu    rwmutex.RWMutex
	array []int
}

// NewIntArray 创建int类型的数组
func NewIntArray(safe ...bool) *IntArray {
	return NewIntArraySize(0, 0, safe...)
}

// NewIntArraySize 创建指定长度的int类型的数组
func NewIntArraySize(size int, cap int, safe ...bool) *IntArray {
	return &IntArray{
		mu:    rwmutex.Create(safe...),
		array: make([]int, size, cap),
	}
}

// NewIntArrayRange 创建一个切片，从start开始到end结束，步进未step
func NewIntArrayRange(start, end, step int, safe ...bool) *IntArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewIntArrayFrom(slice, safe...)
}

// NewIntArrayFrom 从基础数据结构切片中创建自定义数组
func NewIntArrayFrom(array []int, safe ...bool) *IntArray {
	return &IntArray{
		mu:    rwmutex.Create(safe...),
		array: array,
	}
}

// NewIntArrayFromCopy 创建一个新的自定义数组，值是复制传入的切片
func NewIntArrayFromCopy(array []int, safe ...bool) *IntArray {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return &IntArray{
		mu:    rwmutex.Create(safe...),
		array: newArray,
	}
}

// Get 获取指定index的值
func (that *IntArray) Get(index int) (value int, found bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if index < 0 || index >= len(that.array) {
		return 0, false
	}
	return that.array[index], true
}

// SetArray 赋值
func (that *IntArray) SetArray(array []int) *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	that.array = array
	return that
}

// Replace 把传入的array替换原来的数组
func (that *IntArray) Replace(array []int) *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	max := len(array)
	if max > len(that.array) {
		max = len(that.array)
	}
	for i := 0; i < max; i++ {
		that.array[i] = array[i]
	}
	return that
}

// Sum 统计int的总数
func (that *IntArray) Sum() (sum int) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, v := range that.array {
		sum += v
	}
	return
}

// Sort 数组排序
func (that *IntArray) Sort(reverse ...bool) *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	if len(reverse) > 0 && reverse[0] {
		sort.Slice(that.array, func(i, j int) bool {
			if that.array[i] < that.array[j] {
				return false
			}
			return true
		})
	} else {
		sort.Ints(that.array)
	}
	return that
}

// SortFunc 使用指定的排序函数排序数组
func (that *IntArray) SortFunc(less func(v1, v2 int) bool) *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	sort.Slice(that.array, func(i, j int) bool {
		return less(that.array[i], that.array[j])
	})
	return that
}

// InsertBefore 在指定index之前插入value
func (that *IntArray) InsertBefore(index int, value int) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if index < 0 || index >= len(that.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(that.array)))
	}
	rear := append([]int{}, that.array[index:]...)
	that.array = append(that.array[0:index], value)
	that.array = append(that.array, rear...)
	return nil
}

// InsertAfter 在数组index之后插入value
func (that *IntArray) InsertAfter(index int, value int) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if index < 0 || index >= len(that.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(that.array)))
	}
	rear := append([]int{}, that.array[index+1:]...)
	that.array = append(that.array[0:index+1], value)
	that.array = append(that.array, rear...)
	return nil
}

// Remove 移除指定的index对应的值
func (that *IntArray) Remove(index int) (value int, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.doRemoveWithoutLock(index)
}

// 从数组中找到指定indxe的值，从数组中移除它，并返回它。
func (that *IntArray) doRemoveWithoutLock(index int) (value int, found bool) {
	if index < 0 || index >= len(that.array) {
		return 0, false
	}
	// 在删除的时候确认数组边界，提升效率
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

// RemoveValue 查找指定的值，并从数组中移除它
func (that *IntArray) RemoveValue(value int) bool {
	if i := that.Search(value); i != -1 {
		_, found := that.Remove(i)
		return found
	}
	return false
}

func (that *IntArray) PushLeft(value ...int) *IntArray {
	that.mu.Lock()
	that.array = append(value, that.array...)
	that.mu.Unlock()
	return that
}

// PushRight pushes one or multiple items to the end of array.
// It equals to Append.
func (that *IntArray) PushRight(value ...int) *IntArray {
	that.mu.Lock()
	that.array = append(that.array, value...)
	that.mu.Unlock()
	return that
}

// PopLeft pops and returns an item from the beginning of array.
// Note that if the array is empty, the <found> is false.
func (that *IntArray) PopLeft() (value int, found bool) {
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
func (that *IntArray) PopRight() (value int, found bool) {
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

// PopRand randomly pops and return thatn item out of array.
// Note that if the array is empty, the <found> is false.
func (that *IntArray) PopRand() (value int, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.doRemoveWithoutLock(drand.Intn(len(that.array)))
}

// PopRands randomly pops and returns <size> items out of array.
// If the given <size> is greater than size of the array, it returns all elements of the array.
// Note that if given <size> <= 0 or the array is empty, it returns nil.
func (that *IntArray) PopRands(size int) []int {
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
func (that *IntArray) PopLefts(size int) []int {
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
func (that *IntArray) PopRights(size int) []int {
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
// else a pointer to the underlying datthat.
//
// If <end> is negative, then the offset will start from the end of array.
// If <end> is omitted, then the sequence will have everything from start up
// until the end of the array.
func (that *IntArray) Range(start int, end ...int) []int {
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
func (that *IntArray) SubSlice(offset int, length ...int) []int {
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

// Append See PushRight.
func (that *IntArray) Append(value ...int) *IntArray {
	that.mu.Lock()
	that.array = append(that.array, value...)
	that.mu.Unlock()
	return that
}

// Len returns the length of array.
func (that *IntArray) Len() int {
	that.mu.RLock()
	length := len(that.array)
	that.mu.RUnlock()
	return length
}

// Slice returns the underlying data of array.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying datthat.
func (that *IntArray) Slice() []int {
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
func (that *IntArray) Interfaces() []interface{} {
	that.mu.RLock()
	defer that.mu.RUnlock()
	array := make([]interface{}, len(that.array))
	for k, v := range that.array {
		array[k] = v
	}
	return array
}

// Clone returns a new array, which is a copy of current array.
func (that *IntArray) Clone() (newArray *IntArray) {
	that.mu.RLock()
	array := make([]int, len(that.array))
	copy(array, that.array)
	that.mu.RUnlock()
	return NewIntArrayFrom(array, that.mu.IsSafe())
}

// Clear deletes all items of current array.
func (that *IntArray) Clear() *IntArray {
	that.mu.Lock()
	if len(that.array) > 0 {
		that.array = make([]int, 0)
	}
	that.mu.Unlock()
	return that
}

// Contains checks whether a value exists in the array.
func (that *IntArray) Contains(value int) bool {
	return that.Search(value) != -1
}

// Search searches array by <value>, returns the index of <value>,
// or returns -1 if not exists.
func (that *IntArray) Search(value int) int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if len(that.array) == 0 {
		return -1
	}
	result := -1
	for index, v := range that.array {
		if v == value {
			result = index
			break
		}
	}
	return result
}

// Unique uniques the array, clear repeated items.
// Example: [1,1,2,3,2] -> [1,2,3]
func (that *IntArray) Unique() *IntArray {
	that.mu.Lock()
	for i := 0; i < len(that.array)-1; i++ {
		for j := i + 1; j < len(that.array); {
			if that.array[i] == that.array[j] {
				that.array = append(that.array[:j], that.array[j+1:]...)
			} else {
				j++
			}
		}
	}
	that.mu.Unlock()
	return that
}

// LockFunc locks writing by callback function <f>.
func (that *IntArray) LockFunc(f func(array []int)) *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.array)
	return that
}

// RLockFunc locks reading by callback function <f>.
func (that *IntArray) RLockFunc(f func(array []int)) *IntArray {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.array)
	return that
}

// Merge merges <array> into current array.
// The parameter <array> can be any garray or slice type.
// The difference between Merge and Append is Append supports only specified slice type,
// but Merge supports more parameter types.
func (that *IntArray) Merge(array interface{}) *IntArray {
	return that.Append(dconv.Ints(array)...)
}

// Fill fills an array with num entries of the value <value>,
// keys starting at the <startIndex> parameter.
func (that *IntArray) Fill(startIndex int, num int, value int) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if startIndex < 0 || startIndex > len(that.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", startIndex, len(that.array)))
	}
	for i := startIndex; i < startIndex+num; i++ {
		if i > len(that.array)-1 {
			that.array = append(that.array, value)
		} else {
			that.array[i] = value
		}
	}
	return nil
}

// Chunk splits an array into multiple arrays,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
func (that *IntArray) Chunk(size int) [][]int {
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

// Pad pads array to the specified length with <value>.
// If size is positive then the array is padded on the right, or negative on the left.
// If the absolute value of <size> is less than or equal to the length of the array
// then no padding takes place.
func (that *IntArray) Pad(size int, value int) *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	if size == 0 || (size > 0 && size < len(that.array)) || (size < 0 && size > -len(that.array)) {
		return that
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(that.array)
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		tmp[i] = value
	}
	if size > 0 {
		that.array = append(that.array, tmp...)
	} else {
		that.array = append(tmp, that.array...)
	}
	return that
}

// Rand randomly returns one item from array(no deleting).
func (that *IntArray) Rand() (value int, found bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if len(that.array) == 0 {
		return 0, false
	}
	return that.array[drand.Intn(len(that.array))], true
}

// Rands randomly returns <size> items from array(no deleting).
func (that *IntArray) Rands(size int) []int {
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

// Shuffle randomly shuffles the array.
func (that *IntArray) Shuffle() *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	for i, v := range drand.Perm(len(that.array)) {
		that.array[i], that.array[v] = that.array[v], that.array[i]
	}
	return that
}

// Reverse makes array with elements in reverse order.
func (that *IntArray) Reverse() *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	for i, j := 0, len(that.array)-1; i < j; i, j = i+1, j-1 {
		that.array[i], that.array[j] = that.array[j], that.array[i]
	}
	return that
}

// Join joins array elements with a string <glue>.
func (that *IntArray) Join(glue string) string {
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
func (that *IntArray) CountValues() map[int]int {
	m := make(map[int]int)
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, v := range that.array {
		m[v]++
	}
	return m
}

// Iterator is alias of IteratorAsc.
func (that *IntArray) Iterator(f func(k int, v int) bool) {
	that.IteratorAsc(f)
}

// IteratorAsc iterates the array readonly in ascending order with given callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (that *IntArray) IteratorAsc(f func(k int, v int) bool) {
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
func (that *IntArray) IteratorDesc(f func(k int, v int) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for i := len(that.array) - 1; i >= 0; i-- {
		if !f(i, that.array[i]) {
			break
		}
	}
}

// String returns current array as a string, which implements like json.Marshal does.
func (that *IntArray) String() string {
	return "[" + that.Join(",") + "]"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (that IntArray) MarshalJSON() ([]byte, error) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return json.Marshal(that.array)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *IntArray) UnmarshalJSON(b []byte) error {
	if that.array == nil {
		that.array = make([]int, 0)
	}
	that.mu.Lock()
	defer that.mu.Unlock()
	if err := json.Unmarshal(b, &that.array); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for array.
func (that *IntArray) UnmarshalValue(value interface{}) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	switch value.(type) {
	case string, []byte:
		return json.Unmarshal(dconv.Bytes(value), &that.array)
	default:
		that.array = dconv.SliceInt(value)
	}
	return nil
}

// FilterEmpty removes all zero value of the array.
func (that *IntArray) FilterEmpty() *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	for i := 0; i < len(that.array); {
		if that.array[i] == 0 {
			that.array = append(that.array[:i], that.array[i+1:]...)
		} else {
			i++
		}
	}
	return that
}

// Walk applies a user supplied function <f> to every item of array.
func (that *IntArray) Walk(f func(value int) int) *IntArray {
	that.mu.Lock()
	defer that.mu.Unlock()
	for i, v := range that.array {
		that.array[i] = f(v)
	}
	return that
}

// IsEmpty checks whether the array is empty.
func (that *IntArray) IsEmpty() bool {
	return that.Len() == 0
}
