package dvar

import (
	"encoding/json"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/util/dconv"
	"time"
)

type Var struct {
	value interface{} // Underlying value.
	safe  bool        // Concurrent safe or not.
}

// New 创建var的对象指针
func New(value interface{}, safe ...bool) *Var {
	v := Var{}
	if len(safe) > 0 && !safe[0] {
		v.safe = true
		v.value = dtype.NewInterface(value)
	} else {
		v.value = value
	}
	return &v
}

// Create 创建var对象
func Create(value interface{}, safe ...bool) Var {
	v := Var{}
	if len(safe) > 0 && !safe[0] {
		v.safe = true
		v.value = dtype.NewInterface(value)
	} else {
		v.value = value
	}
	return v
}

// Set 设置值
func (that *Var) Set(value interface{}) (old interface{}) {
	if that.safe {
		if t, ok := that.value.(*dtype.Interface); ok {
			old = t.Set(value)
			return
		}
	}
	old = that.value
	that.value = value
	return
}

func (that *Var) Clone() *Var {
	return New(that.Val(), that.safe)
}

// Val 获取dvar的值
func (that *Var) Val() interface{} {
	if that == nil {
		return nil
	}
	if that.safe {
		if t, ok := that.value.(*dtype.Interface); ok {
			return t.Val()
		}
	}
	return that.value
}

// Interface 返回var的值
func (that *Var) Interface() interface{} {
	return that.Val()
}

func (that *Var) Bytes() []byte {
	return dconv.Bytes(that.Val())
}

func (that *Var) String() string {
	return dconv.String(that.Val())
}

func (that *Var) Bool() bool {
	return dconv.Bool(that.Val())
}

func (that *Var) Int() int {
	return dconv.Int(that.Val())
}

// Int8 converts and returns <v> as int8.
func (that *Var) Int8() int8 {
	return dconv.Int8(that.Val())
}

// Int16 converts and returns <v> as int16.
func (that *Var) Int16() int16 {
	return dconv.Int16(that.Val())
}

// Int32 converts and returns <v> as int32.
func (that *Var) Int32() int32 {
	return dconv.Int32(that.Val())
}

// Int64 converts and returns <v> as int64.
func (that *Var) Int64() int64 {
	return dconv.Int64(that.Val())
}

// Uint converts and returns <v> as uint.
func (that *Var) Uint() uint {
	return dconv.Uint(that.Val())
}

// Uint8 converts and returns <v> as uint8.
func (that *Var) Uint8() uint8 {
	return dconv.Uint8(that.Val())
}

// Uint16 converts and returns <v> as uint16.
func (that *Var) Uint16() uint16 {
	return dconv.Uint16(that.Val())
}

// Uint32 converts and returns <v> as uint32.
func (that *Var) Uint32() uint32 {
	return dconv.Uint32(that.Val())
}

// Uint64 converts and returns <v> as uint64.
func (that *Var) Uint64() uint64 {
	return dconv.Uint64(that.Val())
}

// Float32 converts and returns <v> as float32.
func (that *Var) Float32() float32 {
	return dconv.Float32(that.Val())
}

// Float64 converts and returns <v> as float64.
func (that *Var) Float64() float64 {
	return dconv.Float64(that.Val())
}

// Time converts and returns <v> as time.Time.
// The parameter <format> specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (that *Var) Time(format ...string) time.Time {
	return dconv.Time(that.Val(), format...)
}

// Duration converts and returns <v> as time.Duration.
// If value of <v> is string, then it uses time.ParseDuration for conversion.
func (that *Var) Duration() time.Duration {
	return dconv.Duration(that.Val())
}

// GTime converts and returns <v> as *gtime.Time.
// The parameter <format> specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (that *Var) GTime(format ...string) *dtime.Time {
	return dconv.GTime(that.Val(), format...)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (that *Var) MarshalJSON() ([]byte, error) {
	return json.Marshal(that.Val())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *Var) UnmarshalJSON(b []byte) error {
	var i interface{}
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	that.Set(i)
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for Var.
func (that *Var) UnmarshalValue(value interface{}) error {
	that.Set(value)
	return nil
}
