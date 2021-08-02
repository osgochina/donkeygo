// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dbinary_test

import (
	"github.com/osgochina/donkeygo/encoding/dbinary"
	"github.com/osgochina/donkeygo/test/dtest"
	"math"
	"testing"
)

type User struct {
	Name string
	Age  int
	Url  string
}

var testData = map[string]interface{}{
	//"nil":         nil,
	"int":         int(123),
	"int8":        int8(-99),
	"int8.max":    math.MaxInt8,
	"int16":       int16(123),
	"int16.max":   math.MaxInt16,
	"int32":       int32(-199),
	"int32.max":   math.MaxInt32,
	"int64":       int64(123),
	"uint":        uint(123),
	"uint8":       uint8(123),
	"uint8.max":   math.MaxUint8,
	"uint16":      uint16(9999),
	"uint16.max":  math.MaxUint16,
	"uint32":      uint32(123),
	"uint64":      uint64(123),
	"bool.true":   true,
	"bool.false":  false,
	"string":      "hehe haha",
	"byte":        []byte("hehe haha"),
	"float32":     float32(123.456),
	"float32.max": math.MaxFloat32,
	"float64":     float64(123.456),
}

var testBitData = []int{0, 99, 122, 129, 222, 999, 22322}

func Test_EncodeAndDecode(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for k, v := range testData {
			ve := dbinary.Encode(v)
			ve1 := dbinary.EncodeByLength(len(ve), v)

			//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
			switch v.(type) {
			case int:
				t.Assert(dbinary.DecodeToInt(ve), v)
				t.Assert(dbinary.DecodeToInt(ve1), v)
			case int8:
				t.Assert(dbinary.DecodeToInt8(ve), v)
				t.Assert(dbinary.DecodeToInt8(ve1), v)
			case int16:
				t.Assert(dbinary.DecodeToInt16(ve), v)
				t.Assert(dbinary.DecodeToInt16(ve1), v)
			case int32:
				t.Assert(dbinary.DecodeToInt32(ve), v)
				t.Assert(dbinary.DecodeToInt32(ve1), v)
			case int64:
				t.Assert(dbinary.DecodeToInt64(ve), v)
				t.Assert(dbinary.DecodeToInt64(ve1), v)
			case uint:
				t.Assert(dbinary.DecodeToUint(ve), v)
				t.Assert(dbinary.DecodeToUint(ve1), v)
			case uint8:
				t.Assert(dbinary.DecodeToUint8(ve), v)
				t.Assert(dbinary.DecodeToUint8(ve1), v)
			case uint16:
				t.Assert(dbinary.DecodeToUint16(ve1), v)
				t.Assert(dbinary.DecodeToUint16(ve), v)
			case uint32:
				t.Assert(dbinary.DecodeToUint32(ve1), v)
				t.Assert(dbinary.DecodeToUint32(ve), v)
			case uint64:
				t.Assert(dbinary.DecodeToUint64(ve), v)
				t.Assert(dbinary.DecodeToUint64(ve1), v)
			case bool:
				t.Assert(dbinary.DecodeToBool(ve), v)
				t.Assert(dbinary.DecodeToBool(ve1), v)
			case string:
				t.Assert(dbinary.DecodeToString(ve), v)
				t.Assert(dbinary.DecodeToString(ve1), v)
			case float32:
				t.Assert(dbinary.DecodeToFloat32(ve), v)
				t.Assert(dbinary.DecodeToFloat32(ve1), v)
			case float64:
				t.Assert(dbinary.DecodeToFloat64(ve), v)
				t.Assert(dbinary.DecodeToFloat64(ve1), v)
			default:
				if v == nil {
					continue
				}
				res := make([]byte, len(ve))
				err := dbinary.Decode(ve, res)
				if err != nil {
					t.Errorf("test data: %s, %v, error:%v", k, v, err)
				}
				t.Assert(res, v)
			}
		}
	})
}

func Test_EncodeStruct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		user := User{"wenzi1", 999, "www.baidu.com"}
		ve := dbinary.Encode(user)
		s := dbinary.DecodeToString(ve)
		t.Assert(s, s)
	})
}

func Test_Bits(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := range testBitData {
			bits := make([]dbinary.Bit, 0)
			res := dbinary.EncodeBits(bits, testBitData[i], 64)

			t.Assert(dbinary.DecodeBits(res), testBitData[i])
			t.Assert(dbinary.DecodeBitsToUint(res), uint(testBitData[i]))

			t.Assert(dbinary.DecodeBytesToBits(dbinary.EncodeBitsToBytes(res)), res)
		}
	})
}
