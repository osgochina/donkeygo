package dtype

import (
	"encoding/json"
	"github.com/osgochina/donkeygo/util/dconv"
	"sync/atomic"
)

type Interface struct {
	value atomic.Value
}

// NewInterface 创建interface类型
func NewInterface(value ...interface{}) *Interface {
	t := &Interface{}
	if len(value) > 0 && value[0] != nil {
		t.value.Store(value[0])
	}
	return t
}

func (that *Interface) Clone() *Interface {
	return NewInterface(that.Val())
}

// Set 设置值
func (that *Interface) Set(value interface{}) (old interface{}) {
	old = that.Val()
	that.value.Store(value)
	return
}

// Val 获取值
func (that *Interface) Val() interface{} {
	return that.value.Load()
}

//转换成string
func (that *Interface) String() string {
	return dconv.String(that.Val())
}

// MarshalJSON 序列化json
func (that *Interface) MarshalJSON() ([]byte, error) {
	return json.Marshal(that.Val())
}

// UnmarshalJSON 反序列化json
func (that *Interface) UnmarshalJSON(b []byte) error {
	var i interface{}
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	that.Set(i)
	return nil
}

// UnmarshalValue 反序列化值
func (that *Interface) UnmarshalValue(value interface{}) error {
	that.Set(value)
	return nil
}
