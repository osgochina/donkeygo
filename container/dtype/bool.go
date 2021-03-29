package dtype

import (
	"bytes"
	"donkeygo/util/dconv"
	"sync/atomic"
)

//并发安全的 bool 对象
type Bool struct {
	value int32
}

//创建一个并发安全的bool对象，并且默认值是可选
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

//并发安全的克隆它
func (that *Bool) Clone() *Bool {
	return NewBool(that.Val())
}

//并发安全的获取bool对象的值
func (that *Bool) Val() bool {
	return atomic.LoadInt32(&that.value) > 0
}

//并发安全的设置值
func (that *Bool) Set(val bool) (old bool) {
	if val {
		old = atomic.SwapInt32(&that.value, 1) == 1
	} else {
		old = atomic.SwapInt32(&that.value, 0) == 0
	}
	return
}

//并发安全的值比较与值交换
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

//json序列化
func (that *Bool) MarshalJson() ([]byte, error) {
	if that.Val() {
		return bytesTrue, nil
	}
	return bytesFalse, nil
}

//json反序列化
func (that *Bool) UnmarshalJson(b []byte) error {
	that.Set(dconv.Bool(bytes.Trim(b, `"`)))
	return nil
}

//json反序列化的时候针对值的判断
func (that *Bool) UnmarshalValue(value interface{}) error {
	that.Set(dconv.Bool(value))
	return nil
}
