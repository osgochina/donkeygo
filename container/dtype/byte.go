package dtype

import (
	"github.com/osgochina/donkeygo/util/dconv"
	"strconv"
	"sync/atomic"
)

// Byte 并发安全的byte类型
type Byte struct {
	value int32
}

// NewByte 新建一个byte对象
func NewByte(value ...byte) *Byte {
	if len(value) > 0 {
		return &Byte{
			value: int32(value[0]),
		}
	}
	return &Byte{}
}

// Clone 克隆被返回一个新的并发安全byte对象
func (that *Byte) Clone() *Byte {
	return NewByte(that.Val())
}

// Set 设置值
func (that *Byte) Set(value byte) (old byte) {
	return byte(atomic.SwapInt32(&that.value, int32(value)))
}

// Val 获取值
func (that *Byte) Val() byte {
	return byte(atomic.LoadInt32(&that.value))
}

// Add 添加delta
func (that *Byte) Add(delta byte) (new byte) {
	return byte(atomic.AddInt32(&that.value, int32(delta)))
}

// Cas 原子性的替换新旧两个值
func (that *Byte) Cas(old, new byte) (swapped bool) {
	return atomic.CompareAndSwapInt32(&that.value, int32(old), int32(new))
}

// 转成字符串
func (that *Byte) String() string {
	return strconv.FormatUint(uint64(that.Val()), 10)
}

func (that *Byte) MarshalJSON() ([]byte, error) {
	return dconv.UnsafeStrToBytes(strconv.FormatUint(uint64(that.Val()), 10)), nil
}

func (that *Byte) UnmarshalJSON(b []byte) error {
	that.Set(dconv.Uint8(dconv.UnsafeBytesToStr(b)))
	return nil
}

func (that *Byte) UnmarshalValue(value interface{}) error {
	that.Set(dconv.Byte(value))
	return nil
}
