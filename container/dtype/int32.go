package dtype

import (
	"github.com/osgochina/donkeygo/util/dconv"
	"strconv"
	"sync/atomic"
)

// Int32  并发安全的int对象
type Int32 struct {
	value int32
}

// NewInt32 创建Int32对象
func NewInt32(val ...int32) *Int32 {
	if len(val) > 0 {
		return &Int32{
			value: val[0],
		}
	}
	return &Int32{}
}

// Clone 并发安全克隆一个int对象
func (that *Int32) Clone() *Int32 {
	return NewInt32(that.Val())
}

// Val 并发安全获取int的值
func (that *Int32) Val() int32 {
	return atomic.LoadInt32(&that.value)
}

// Set 并发安全设置int的值
func (that *Int32) Set(val int32) (old int32) {
	return atomic.SwapInt32(&that.value, val)
}

// Add 并发安全设置int的值
func (that *Int32) Add(delta int32) (new int32) {
	return atomic.AddInt32(&that.value, delta)
}

// Cas 安全的值比较与值交换
func (that *Int32) Cas(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&that.value, old, new)
}

//转换成字符串
func (that *Int32) String() string {
	return strconv.Itoa(int(that.Val()))
}

// MarshalJSON json序列化
func (that *Int32) MarshalJSON() ([]byte, error) {
	return dconv.UnsafeStrToBytes(that.String()), nil
}

// UnmarshalJSON json反序列化
func (that *Int32) UnmarshalJSON(b []byte) error {
	that.Set(dconv.Int32(dconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue 值的反序列化
func (that *Int32) UnmarshalValue(value interface{}) error {
	that.Set(dconv.Int32(value))
	return nil
}

// Reduce 减少 delta
func (that *Int32) Reduce(delta int32) (new int32) {
	return atomic.AddInt32(&that.value, -delta)
}
