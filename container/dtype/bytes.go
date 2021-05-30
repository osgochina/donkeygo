package dtype

import (
	"bytes"
	"encoding/base64"
	"github.com/osgochina/donkeygo/util/dconv"
	"sync/atomic"
)

// Bytes 并发安全的bytes数组，[]byte
type Bytes struct {
	value atomic.Value
}

// NewBytes 新建一个并发安全的bytes对象，类似于 []byte 结构
func NewBytes(value ...[]byte) *Bytes {
	t := &Bytes{}
	if len(value) > 0 {
		t.value.Store(value[0])
	}
	return t
}

// Clone 克隆并返回一个新的bytes对象
func (that *Bytes) Clone() *Bytes {
	return NewBytes(that.Val())
}

// Set 设置内容
func (that *Bytes) Set(value []byte) (old []byte) {
	old = that.Val()
	that.value.Store(value)
	return
}

// Val 获取值
func (that *Bytes) Val() []byte {
	if s := that.value.Load(); s != nil {
		return s.([]byte)
	}
	return nil
}

// 转换成字符串
func (that *Bytes) String() string {
	return string(that.Val())
}

func (that *Bytes) MarshalJSON() ([]byte, error) {
	val := that.Val()
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(val)))
	base64.StdEncoding.Encode(dst, val)
	return dconv.UnsafeStrToBytes(`"` + dconv.UnsafeBytesToStr(dst) + `"`), nil
}

func (that *Bytes) UnmarshalJSON(b []byte) error {
	src := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	n, err := base64.StdEncoding.Decode(src, bytes.Trim(b, `"`))
	if err != nil {
		return nil
	}
	that.Set(src[:n])
	return nil
}

func (that *Bytes) UnmarshalValue(value interface{}) error {
	that.Set(dconv.Bytes(value))
	return nil
}
