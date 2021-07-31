// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar_test

import (
	"bytes"
	"encoding/binary"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
	"time"
)

func Test_Set(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var v dvar.Var
		v.Set(123.456)
		t.Assert(v.Val(), 123.456)
	})
	dtest.C(t, func(t *dtest.T) {
		var v dvar.Var
		v.Set(123.456)
		t.Assert(v.Val(), 123.456)
	})
	dtest.C(t, func(t *dtest.T) {
		v := dvar.Create(123.456)
		t.Assert(v.Val(), 123.456)
	})
	dtest.C(t, func(t *dtest.T) {
		objOne := dvar.New("old", true)
		objOneOld, _ := objOne.Set("new").(string)
		t.Assert(objOneOld, "old")

		objTwo := dvar.New("old", false)
		objTwoOld, _ := objTwo.Set("new").(string)
		t.Assert(objTwoOld, "old")
	})
}

func Test_Val(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		objOne := dvar.New(1, true)
		objOneOld, _ := objOne.Val().(int)
		t.Assert(objOneOld, 1)

		objTwo := dvar.New(1, false)
		objTwoOld, _ := objTwo.Val().(int)
		t.Assert(objTwoOld, 1)
	})
}
func Test_Interface(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		objOne := dvar.New(1, true)
		objOneOld, _ := objOne.Interface().(int)
		t.Assert(objOneOld, 1)

		objTwo := dvar.New(1, false)
		objTwoOld, _ := objTwo.Interface().(int)
		t.Assert(objTwoOld, 1)
	})
}
func Test_IsNil(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		objOne := dvar.New(nil, true)
		t.Assert(objOne.IsNil(), true)

		objTwo := dvar.New("noNil", false)
		t.Assert(objTwo.IsNil(), false)

	})
}

func Test_Bytes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		x := int32(1)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, x)

		objOne := dvar.New(bytesBuffer.Bytes(), true)

		bBuf := bytes.NewBuffer(objOne.Bytes())
		var y int32
		binary.Read(bBuf, binary.BigEndian, &y)

		t.Assert(x, y)

	})
}

func Test_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var str string = "hello"
		objOne := dvar.New(str, true)
		t.Assert(objOne.String(), str)

	})
}
func Test_Bool(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var ok bool = true
		objOne := dvar.New(ok, true)
		t.Assert(objOne.Bool(), ok)

		ok = false
		objTwo := dvar.New(ok, true)
		t.Assert(objTwo.Bool(), ok)

	})
}

func Test_Int(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num int = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Int(), num)

	})
}

func Test_Int8(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num int8 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Int8(), num)

	})
}

func Test_Int16(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num int16 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Int16(), num)

	})
}

func Test_Int32(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num int32 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Int32(), num)

	})
}

func Test_Int64(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num int64 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Int64(), num)

	})
}

func Test_Uint(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num uint = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Uint(), num)

	})
}

func Test_Uint8(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num uint8 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Uint8(), num)

	})
}

func Test_Uint16(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num uint16 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Uint16(), num)

	})
}

func Test_Uint32(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num uint32 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Uint32(), num)

	})
}

func Test_Uint64(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num uint64 = 1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Uint64(), num)

	})
}
func Test_Float32(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num float32 = 1.1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Float32(), num)

	})
}

func Test_Float64(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var num float64 = 1.1
		objOne := dvar.New(num, true)
		t.Assert(objOne.Float64(), num)

	})
}

func Test_Time(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var timeUnix int64 = 1556242660
		objOne := dvar.New(timeUnix, true)
		t.Assert(objOne.Time().Unix(), timeUnix)
	})
}

func Test_GTime(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var timeUnix int64 = 1556242660
		objOne := dvar.New(timeUnix, true)
		t.Assert(objOne.GTime().Unix(), timeUnix)
	})
}

func Test_Duration(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var timeUnix int64 = 1556242660
		objOne := dvar.New(timeUnix, true)
		t.Assert(objOne.Duration(), time.Duration(timeUnix))
	})
}

func Test_UnmarshalJson(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type V struct {
			Name string
			Var  *dvar.Var
		}
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
	dtest.C(t, func(t *dtest.T) {
		type V struct {
			Name string
			Var  dvar.Var
		}
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
}

func Test_UnmarshalValue(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type V struct {
			Name string
			Var  *dvar.Var
		}
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
	dtest.C(t, func(t *dtest.T) {
		type V struct {
			Name string
			Var  dvar.Var
		}
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
}
