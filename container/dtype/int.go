package dtype

import (
	"donkeygo/util/dconv"
	"strconv"
	"sync/atomic"
)

// Int 并发安全的int对象
type Int struct {
	value int64
}

// NewInt 创建Int对象
func NewInt(val ...int) *Int {
	if len(val) > 0 {
		return &Int{
			value: int64(val[0]),
		}
	}
	return &Int{}
}

// Clone 并发安全克隆一个int对象
func (that *Int) Clone() *Int {
	return NewInt(that.Val())
}

// Val 并发安全获取int的值
func (that *Int) Val() int {
	return int(atomic.LoadInt64(&that.value))
}

// Set 并发安全设置int的值
func (that *Int) Set(val int) (old int) {
	return int(atomic.SwapInt64(&that.value, int64(val)))
}

// Add 并发安全设置int的值
func (that *Int) Add(delta int) (new int) {
	return int(atomic.AddInt64(&that.value, int64(delta)))
}

// Cas 安全的值比较与值交换
func (that *Int) Cas(old, new int) (swapped bool) {
	return atomic.CompareAndSwapInt64(&that.value, int64(old), int64(new))
}

//转换成字符串
func (that *Int) String() string {
	return strconv.Itoa(that.Val())
}

// MarshalJSON json序列化
func (that *Int) MarshalJSON() ([]byte, error) {
	return dconv.UnsafeStrToBytes(that.String()), nil
}

// UnmarshalJSON json反序列化
func (that *Int) UnmarshalJSON(b []byte) error {
	that.Set(dconv.Int(dconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue 值的反序列化
func (that *Int) UnmarshalValue(value interface{}) error {
	that.Set(dconv.Int(value))
	return nil
}

// Reduce 减少 delta
func (that *Int) Reduce(delta int) (new int) {
	return int(atomic.AddInt64(&that.value, -int64(delta)))
}
