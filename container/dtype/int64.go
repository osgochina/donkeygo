package dtype

import (
	"github.com/osgochina/donkeygo/util/dconv"
	"strconv"
	"sync/atomic"
)

// Int64 并发安全的int对象
type Int64 struct {
	value int64
}

// NewInt64 创建Int对象
func NewInt64(val ...int64) *Int64 {
	if len(val) > 0 {
		return &Int64{
			value: val[0],
		}
	}
	return &Int64{}
}

// Clone 并发安全克隆一个int对象
func (that *Int64) Clone() *Int64 {
	return NewInt64(that.Val())
}

// Val 并发安全获取int的值
func (that *Int64) Val() int64 {
	return atomic.LoadInt64(&that.value)
}

// Set 并发安全设置int的值
func (that *Int64) Set(val int64) (old int64) {
	return atomic.SwapInt64(&that.value, val)
}

// Add 并发安全设置int的值
func (that *Int64) Add(delta int64) (new int64) {
	return atomic.AddInt64(&that.value, delta)
}

// Cas 安全的值比较与值交换
func (that *Int64) Cas(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&that.value, old, new)
}

//转换成字符串
func (that *Int64) String() string {
	return strconv.Itoa(int(that.Val()))
}

// MarshalJSON json序列化
func (that *Int64) MarshalJSON() ([]byte, error) {
	return dconv.UnsafeStrToBytes(that.String()), nil
}

// UnmarshalJSON json反序列化
func (that *Int64) UnmarshalJSON(b []byte) error {
	that.Set(dconv.Int64(dconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue 值的反序列化
func (that *Int64) UnmarshalValue(value interface{}) error {
	that.Set(dconv.Int64(value))
	return nil
}

// Reduce 减少 delta
func (that *Int64) Reduce(delta int64) (new int64) {
	return atomic.AddInt64(&that.value, -delta)
}
