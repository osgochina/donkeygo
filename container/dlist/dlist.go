package dlist

import (
	"bytes"
	"container/list"
	"donkeygo/internal/rwmutex"
	"donkeygo/util/dconv"
	"encoding/json"
)

type (
	// List 是一个包含并发安全开关的双向链表，开关在对象创建的时候确认，并且不能更改。
	List struct {
		list *list.List
		mu   rwmutex.RWMutex
	}
	// Element 链表中的节点类型
	Element = list.Element
)

// New 创建
func New(safe ...bool) *List {
	return &List{
		list: list.New(),
		mu:   rwmutex.Create(safe...),
	}
}

// NewFrom 创建一个链表，并把切片中的对象填充进该链表
func NewFrom(array []interface{}, safe ...bool) *List {
	l := list.New()
	for _, v := range array {
		l.PushBack(v)
	}

	return &List{
		mu:   rwmutex.Create(safe...),
		list: l,
	}
}

// PushFront 往链表的头部插入一个值为v的节点，并返回该节点
func (that *List) PushFront(v interface{}) *Element {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	return that.list.PushFront(v)
}

// PushBack 往链表的尾部插入一个值为v的节点，并返回该节点
func (that *List) PushBack(v interface{}) *Element {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	return that.list.PushBack(v)
}

// PushFronts 往链表的头部插入多个值
func (that *List) PushFronts(values []interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	for _, v := range values {
		that.list.PushFront(v)
	}
}

// PushBacks 往链表的尾部插入多个值
func (that *List) PushBacks(values []interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	for _, v := range values {
		that.list.PushBack(v)
	}
}

// PopFront 获取链表的头节点的值，并且移除该节点
func (that *List) PopFront() interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
		return nil
	}
	if e := that.list.Front(); e != nil {
		return that.list.Remove(e)
	}
	return nil
}

// PopBack 获取链表的尾节点的值，并且移除该节点
func (that *List) PopBack() interface{} {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
		return nil
	}
	if e := that.list.Back(); e != nil {
		return that.list.Remove(e)
	}
	return nil
}

// PopBacks 取出链表尾部最多max个节点的值
func (that *List) PopBacks(max int) (values []interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
		return
	}
	length := that.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]interface{}, length)
		for i := 0; i < length; i++ {
			values[i] = that.list.Remove(that.list.Back())
		}
	}
	return
}

// PopFronts 取出链表头部最多max个节点的值
func (that *List) PopFronts(max int) (values []interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
		return
	}
	length := that.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]interface{}, length)
		for i := 0; i < length; i++ {
			values[i] = that.list.Remove(that.list.Front())
		}
	}
	return
}

// PopBackAll 从链表尾部开始取，取出全部的值
func (that *List) PopBackAll() []interface{} {
	return that.PopBacks(-1)
}

// PopFrontAll 从链表头部开始取，取出全部的值
func (that *List) PopFrontAll() []interface{} {
	return that.PopFronts(-1)
}

// FrontAll 查看链表从头开始的所有值
func (that *List) FrontAll() (values []interface{}) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return nil
	}
	length := that.list.Len()
	if length > 0 {
		values = make([]interface{}, length)
		for i, e := 0, that.list.Front(); i < length; i, e = i+1, e.Next() {
			values[i] = e.Value
		}
	}
	return
}

// BackAll 查看链表从尾开始的所有值
func (that *List) BackAll() (values []interface{}) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return nil
	}
	length := that.list.Len()
	if length > 0 {
		values = make([]interface{}, length)
		for i, e := 0, that.list.Back(); i < length; i, e = i+1, e.Prev() {
			values[i] = e.Value
		}
	}
	return
}

// FrontValue 查看链表头部的值
func (that *List) FrontValue() (value interface{}) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return
	}
	if e := that.list.Front(); e != nil {
		value = e.Value
	}
	return
}

// BackValue 查看链表尾部的值
func (that *List) BackValue() (value interface{}) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return
	}
	if e := that.list.Back(); e != nil {
		value = e.Value
	}
	return
}

// Front 获取链表的头部节点
func (that *List) Front() (e *Element) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return
	}
	e = that.list.Front()
	return
}

// Back 获取链表的尾部节点
func (that *List) Back() (e *Element) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return
	}
	e = that.list.Back()
	return
}

// Len 获取链表的长度
func (that *List) Len() (length int) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return
	}
	length = that.list.Len()
	return
}

// Size 获取链表的长度
func (that *List) Size() int {
	return that.Len()
}

// RemoveAll 移除所有节点
func (that *List) RemoveAll() {
	that.mu.Lock()
	that.list = list.New()
	that.mu.Unlock()
}

// Clear 清空链表
func (that *List) Clear() {
	that.RemoveAll()
}

// MoveBefore 把e节点从链表中移动到p节点的前面
func (that *List) MoveBefore(e, p *Element) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	that.list.MoveBefore(e, p)
}

// MoveAfter 把e节点移动到p节点的后面
func (that *List) MoveAfter(e, p *Element) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	that.list.MoveAfter(e, p)
}

// MoveToFront 把节点e移动到链表的头部
func (that *List) MoveToFront(e *Element) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	that.list.MoveToFront(e)
}

// MoveToBack 把节点e移动到链表的尾部
func (that *List) MoveToBack(e *Element) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	that.list.MoveToBack(e)
}

// PushBackList 把链表other链接到链表的尾部
func (that *List) PushBackList(other *List) {
	if that != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	that.list.PushBackList(other.list)
}

// PushFrontList 把链表other链接到链表的头部
func (that *List) PushFrontList(other *List) {
	if that != other {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	that.list.PushFrontList(other.list)
}

// InsertAfter 插入值v到节点p的后面，并返回该节点
func (that *List) InsertAfter(p *Element, v interface{}) (e *Element) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	e = that.list.InsertAfter(v, p)
	return
}

// InsertBefore 插入值v到节点p的前面，并返回该节点
func (that *List) InsertBefore(p *Element, v interface{}) (e *Element) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	e = that.list.InsertBefore(v, p)
	return
}

// Remove 移除节点e，并返回该节点的值
func (that *List) Remove(e *Element) (value interface{}) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	value = that.list.Remove(e)
	return
}

// Removes 移除指定的一批节点
func (that *List) Removes(es []*Element) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	for _, e := range es {
		that.list.Remove(e)
	}
	return
}

// RLockFunc 加读锁执行函数
func (that *List) RLockFunc(f func(list *list.List)) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list != nil {
		f(that.list)
	}
}

// LockFunc 加读写锁执行函数
func (that *List) LockFunc(f func(list *list.List)) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	f(that.list)
}

// Iterator 迭代
func (that *List) Iterator(f func(e *Element) bool) {
	that.IteratorAsc(f)
}

// IteratorAsc 从头开始迭代
func (that *List) IteratorAsc(f func(e *Element) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return
	}
	length := that.list.Len()
	if length > 0 {
		for i, e := 0, that.list.Front(); i < length; i, e = i+1, e.Next() {
			if !f(e) {
				break
			}
		}
	}
}

// IteratorDesc 从尾部开始迭代
func (that *List) IteratorDesc(f func(e *Element) bool) {
	that.mu.RLock()
	defer that.mu.RUnlock()
	if that.list == nil {
		return
	}
	length := that.list.Len()
	if length > 0 {
		for i, e := 0, that.list.Back(); i < length; i, e = i+1, e.Prev() {
			if !f(e) {
				break
			}
		}
	}
}

// Join 把链表中的值，按glue分隔符链接起来
func (that *List) Join(glue string) string {
	that.mu.RLock()
	defer that.mu.RUnlock()

	if that.list == nil {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	length := that.list.Len()
	if length > 0 {
		for i, e := 0, that.list.Front(); i < length; i, e = i+1, e.Next() {
			buffer.WriteString(dconv.String(e.Value))
			if i != length-1 {
				buffer.WriteString(glue)
			}
		}
	}

	return buffer.String()
}

//把链表转成字符串
func (that *List) String() string {
	return "[" + that.Join(",") + "]"
}

// MarshalJSON 把链表转成json串
func (that *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(that.FrontAll())
}

// UnmarshalJSON 把json串反序列化成链表
func (that *List) UnmarshalJSON(b []byte) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	var array []interface{}
	if err := json.Unmarshal(b, &array); err != nil {
		return err
	}
	that.PushBacks(array)
	return nil
}

// UnmarshalValue 解码内容
func (that *List) UnmarshalValue(value interface{}) (err error) {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.list == nil {
		that.list = list.New()
	}
	var array []interface{}
	switch value.(type) {
	case string, []byte:
		err = json.Unmarshal(dconv.Bytes(value), &array)
	default:
		array = dconv.SliceAny(value)
	}
	that.PushBacks(array)
	return err
}
