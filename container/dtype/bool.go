package dtype

import (
	"bytes"
	"github.com/osgochina/donkeygo/util/dconv"
	"sync/atomic"
)

// Bool 并发安全的 bool 对象
type Bool struct {
	value int32
}

// NewBool 创建一个并发安全的bool对象，并且默认值是可选
func NewBool(value ...bool) *Bool {
	b := &Bool{}
	if len(value) > 0 {
		if value[0] {
			b.value = 1
		} else {
			b.value = 0
		}
	}
	return b
}

// Clone 并发安全的克隆它
func (that *Bool) Clone() *Bool {
	return NewBool(that.Val())
}

// Val 并发安全的获取bool对象的值
func (that *Bool) Val() bool {
	return atomic.LoadInt32(&that.value) > 0
}

// Set 并发安全的设置值
func (that *Bool) Set(val bool) (old bool) {
	if val {
		old = atomic.SwapInt32(&that.value, 1) == 1
	} else {
		old = atomic.SwapInt32(&that.value, 0) == 1
	}
	return
}

// Cas 并发安全的值比较与值交换
func (that *Bool) Cas(old, new bool) (swapped bool) {
	var oldInt32, newInt32 int32
	if old {
		oldInt32 = 1
	}
	if new {
		newInt32 = 1
	}
	return atomic.CompareAndSwapInt32(&that.value, oldInt32, newInt32)
}

//转换成字符串
func (that *Bool) String() string {
	if that.Val() {
		return "true"
	}
	return "false"
}

var (
	bytesTrue  = []byte("true")
	bytesFalse = []byte("false")
)

// MarshalJSON json序列化
func (that *Bool) MarshalJSON() ([]byte, error) {
	if that.Val() {
		return bytesTrue, nil
	}
	return bytesFalse, nil
}

// UnmarshalJSON json反序列化
func (that *Bool) UnmarshalJSON(b []byte) error {
	that.Set(dconv.Bool(bytes.Trim(b, `"`)))
	return nil
}

// UnmarshalValue json反序列化的时候针对值的判断
func (that *Bool) UnmarshalValue(value interface{}) error {
	that.Set(dconv.Bool(value))
	return nil
}
