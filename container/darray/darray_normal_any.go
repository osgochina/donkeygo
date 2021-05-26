package darray

import (
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/util/dconv"
	"github.com/osgochina/donkeygo/util/drand"
	"math"
	"sort"
)

// Array 自定义的的数组
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

// NewArrayRange 创建一个切片，从start开始到end结束，步进未step
func NewArrayRange(start, end, step int, safe ...bool) *Array {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]interface{}, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewArrayFrom(slice, safe...)
}

// NewFrom see NewArrayFrom
func NewFrom(array []interface{}, safe ...bool) *Array {
	return NewArrayFrom(array, safe...)
}

// NewFromCopy see NewArrayFromCopy
func NewFromCopy(array []interface{}, safe ...bool) *Array {
	return NewArrayFromCopy(array, safe...)
}

// NewArrayFrom 从基础数据结构切片中创建自定义数组
func NewArrayFrom(array []interface{}, safe ...bool) *Array {
	return &Array{
		mu:    rwmutex.Create(safe...),
		array: array,
	}
}

// NewArrayFromCopy 创建一个新的自定义数组，值是复制传入的切片
func NewArrayFromCopy(array []interface{}, safe ...bool) *Array {
	newArray := make([]interface{}, len(array))
	copy(newArray, array)
	return &Array{
		mu:    rwmutex.Create(safe...),
		array: newArray,
	}
}

// Get 获取指定index的值
func (that *Array) Get(index int) (value interface{}, found bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if index < 0 || index >= len(that.array) {
		return nil, false
	}
	return that.array[index], true
}

// Set 设置指定index对应的值
func (that *Array) Set(index int, value interface{}) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if index < 0 || index >= len(that.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(that.array)))
	}
	that.array[index] = value
	return nil
}

// SetArray 直接把切片赋值给数组
func (that *Array) SetArray(array []interface{}) *Array {
	that.mu.Lock()
	defer that.mu.Unlock()
	that.array = array
	return that
}

// Replace 把数组中的切片值替换成传入的切片
func (that *Array) Replace(array []interface{}) *Array {
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

// Sum 把切片中的值转换成 int类型，并计算他们的总数
func (that *Array) Sum() (sum int) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	for _, v := range that.array {
		sum += dconv.Int(v)
	}
	return
}

// SortFunc 使用排序方法对数组进行排序
func (that *Array) SortFunc(less func(v1, v2 interface{}) bool) *Array {
	that.mu.Lock()
	defer that.mu.Unlock()
	sort.Slice(that.array, func(i, j int) bool {
		return less(that.array[i], that.array[j])
	})
	return that
}

// InsertBefore 在指定的index之前插入值
func (that *Array) InsertBefore(index int, value interface{}) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if index < 0 || index >= len(that.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(that.array)))
	}
	rear := append([]interface{}{}, that.array[index:]...)
	that.array = append(that.array[0:index], value)
	that.array = append(that.array, rear...)
	return nil
}

// InsertAfter 在指定的index之后插入值
func (that *Array) InsertAfter(index int, value interface{}) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if index < 0 || index >= len(that.array) {
		return errors.New(fmt.Sprintf("index %d out of array range %d", index, len(that.array)))
	}
	rear := append([]interface{}{}, that.array[index+1:]...)
	that.array = append(that.array[0:index+1], value)
	that.array = append(that.array, rear...)
	return nil
}

// Remove 移除指定的index对应的值
func (that *Array) Remove(index int) (value interface{}, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.doRemoveWithoutLock(index)
}

// 从数组中找到指定indxe的值，从数组中移除它，并返回它。
func (that *Array) doRemoveWithoutLock(index int) (value interface{}, found bool) {
	if index < 0 || index >= len(that.array) {
		return nil, false
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
func (that *Array) RemoveValue(value interface{}) bool {
	if i := that.Search(value); i != -1 {
		that.Remove(i)
		return true
	}
	return false
}

// PushLeft 把值插入到数组的左边
func (that *Array) PushLeft(value ...interface{}) *Array {
	that.mu.Lock()
	that.array = append(value, that.array...)
	that.mu.Unlock()
	return that
}

// PushRight 把值插入到数组的右边
func (that *Array) PushRight(value ...interface{}) *Array {
	that.mu.Lock()
	that.array = append(that.array, value...)
	that.mu.Unlock()
	return that
}

// PopRand 随机从数组中拿出一个值
func (that *Array) PopRand() (value interface{}, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.doRemoveWithoutLock(drand.Intn(len(that.array)))
}

// PopRands 随机从数组中拿出指定的值
func (that *Array) PopRands(size int) []interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if size <= 0 || len(that.array) == 0 {
		return nil
	}
	if size >= len(that.array) {
		size = len(that.array)
	}
	array := make([]interface{}, size)
	for i := 0; i < size; i++ {
		array[i], _ = that.doRemoveWithoutLock(drand.Intn(len(that.array)))
	}
	return array
}

// PopLeft 从左边拿出一个值
func (that *Array) PopLeft() (value interface{}, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if len(that.array) == 0 {
		return nil, false
	}
	value = that.array[0]
	that.array = that.array[1:]
	return value, true
}

// PopRight 从右边拿出一个值
func (that *Array) PopRight() (value interface{}, found bool) {
	that.mu.Lock()
	defer that.mu.Unlock()
	index := len(that.array) - 1
	if index < 0 {
		return nil, false
	}
	value = that.array[index]
	that.array = that.array[:index]
	return value, true
}

// PopLefts 从左边拿出指定数量的值
func (that *Array) PopLefts(size int) []interface{} {
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

// PopRights 从右边拿出指定数量的值
func (that *Array) PopRights(size int) []interface{} {
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

// Range 获取从start到end的index的值，把它们赋值给一个新的切片，并返回
func (that *Array) Range(start int, end ...int) []interface{} {
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
	array := ([]interface{})(nil)
	if that.mu.IsSafe() {
		array = make([]interface{}, offsetEnd-start)
		copy(array, that.array[start:offsetEnd])
	} else {
		array = that.array[start:offsetEnd]
	}
	return array
}

// SubSlice 截取数组中的某一段数组，并返回
// 支持offset，length负数
func (that *Array) SubSlice(offset int, length ...int) []interface{} {
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
		s := make([]interface{}, size)
		copy(s, that.array[offset:])
		return s
	} else {
		return that.array[offset:end]
	}
}

// Append see PushRight
func (that *Array) Append(value ...interface{}) *Array {
	that.PushRight(value...)
	return that
}

// Len 数组的长度
func (that *Array) Len() int {
	that.mu.RLock()
	length := len(that.array)
	that.mu.RUnlock()
	return length
}

// Slice 把数组转换成切片返回
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

// Interfaces 返回数组的当前值
func (that *Array) Interfaces() []interface{} {
	return that.Slice()
}

// Clone 复制数组
func (that *Array) Clone() (newArray *Array) {
	that.mu.RLock()
	array := make([]interface{}, len(that.array))
	copy(array, that.array)
	that.mu.RUnlock()
	return NewArrayFrom(array, that.mu.IsSafe())
}

// Clear 清空数组数据
func (that *Array) Clear() *Array {
	that.mu.Lock()
	if len(that.array) > 0 {
		that.array = make([]interface{}, 0)
	}
	that.mu.Unlock()
	return that
}

// Contains 判断值在数组是否存在
func (that *Array) Contains(value interface{}) bool {
	return that.Search(value) != -1
}

// Search 遍历数组，查找指定的值是否在数组中存在，不存在返回-1，存在返回对应的index
func (that *Array) Search(value interface{}) int {
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

// Unique 去除数组中相同的值
func (that *Array) Unique() *Array {
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

// LockFunc 加锁执行指定方法，把数组值作为参数传入
func (that *Array) LockFunc(f func(array []interface{})) *Array {
	that.mu.Lock()
	defer that.mu.Unlock()
	f(that.array)
	return that
}

// RLockFunc 加读锁执行方法，把数组作为参数传入
func (that *Array) RLockFunc(f func(array []interface{})) *Array {
	that.mu.RLock()
	defer that.mu.RUnlock()
	f(that.array)
	return that
}

// Merge 合并两个数组
func (that *Array) Merge(array interface{}) *Array {
	return that.Append(dconv.Interfaces(array)...)
}

// Fill 使用指定的value填充数组，从startindex开始，填充num个索引
func (that *Array) Fill(startIndex int, num int, value interface{}) error {
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

// Chunk 分片，使用size个数目分片成多个数组，返回二维数组
func (that *Array) Chunk(size int) [][]interface{} {
	if size < 1 {
		return nil
	}
	that.mu.RLock()
	defer that.mu.RUnlock()
	length := len(that.array)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
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
