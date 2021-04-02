package dtype

import (
	"bytes"
	"donkeygo/util/dconv"
	"sync/atomic"
)

type String struct {
	value atomic.Value
}

//创建并发安全的string对象
func NewString(value ...string) *String {
	s := &String{}
	if len(value) > 0 {
		s.value.Store(value[0])
	}
	return s
}

func (that *String) Clone() *String {
	return NewString(that.Val())
}

func (that *String) Val() string {
	s := that.value.Load()
	if s != nil {
		return s.(string)
	}
	return ""
}

func (that *String) String() string {
	return that.Val()
}

func (that *String) Set(value string) (old string) {
	old = that.Val()
	that.value.Store(value)
	return
}

func (that *String) MarshalJSON() ([]byte, error) {
	return dconv.UnsafeStrToBytes(`"` + that.Val() + `"`), nil
}

func (that *String) UnmarshalJSON(b []byte) error {
	that.Set(dconv.UnsafeBytesToStr(bytes.Trim(b, `"`)))
	return nil
}

func (that *String) UnmarshalValue(value interface{}) error {
	that.Set(dconv.String(value))
	return nil
}
