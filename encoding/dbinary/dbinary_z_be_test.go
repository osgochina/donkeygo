// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dbinary_test

import (
	"github.com/osgochina/donkeygo/encoding/dbinary"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_BeEncodeAndBeDecode(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for k, v := range testData {
			ve := dbinary.BeEncode(v)
			ve1 := dbinary.BeEncodeByLength(len(ve), v)

			//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
			switch v.(type) {
			case int:
				t.Assert(dbinary.BeDecodeToInt(ve), v)
				t.Assert(dbinary.BeDecodeToInt(ve1), v)
			case int8:
				t.Assert(dbinary.BeDecodeToInt8(ve), v)
				t.Assert(dbinary.BeDecodeToInt8(ve1), v)
			case int16:
				t.Assert(dbinary.BeDecodeToInt16(ve), v)
				t.Assert(dbinary.BeDecodeToInt16(ve1), v)
			case int32:
				t.Assert(dbinary.BeDecodeToInt32(ve), v)
				t.Assert(dbinary.BeDecodeToInt32(ve1), v)
			case int64:
				t.Assert(dbinary.BeDecodeToInt64(ve), v)
				t.Assert(dbinary.BeDecodeToInt64(ve1), v)
			case uint:
				t.Assert(dbinary.BeDecodeToUint(ve), v)
				t.Assert(dbinary.BeDecodeToUint(ve1), v)
			case uint8:
				t.Assert(dbinary.BeDecodeToUint8(ve), v)
				t.Assert(dbinary.BeDecodeToUint8(ve1), v)
			case uint16:
				t.Assert(dbinary.BeDecodeToUint16(ve1), v)
				t.Assert(dbinary.BeDecodeToUint16(ve), v)
			case uint32:
				t.Assert(dbinary.BeDecodeToUint32(ve1), v)
				t.Assert(dbinary.BeDecodeToUint32(ve), v)
			case uint64:
				t.Assert(dbinary.BeDecodeToUint64(ve), v)
				t.Assert(dbinary.BeDecodeToUint64(ve1), v)
			case bool:
				t.Assert(dbinary.BeDecodeToBool(ve), v)
				t.Assert(dbinary.BeDecodeToBool(ve1), v)
			case string:
				t.Assert(dbinary.BeDecodeToString(ve), v)
				t.Assert(dbinary.BeDecodeToString(ve1), v)
			case float32:
				t.Assert(dbinary.BeDecodeToFloat32(ve), v)
				t.Assert(dbinary.BeDecodeToFloat32(ve1), v)
			case float64:
				t.Assert(dbinary.BeDecodeToFloat64(ve), v)
				t.Assert(dbinary.BeDecodeToFloat64(ve1), v)
			default:
				if v == nil {
					continue
				}
				res := make([]byte, len(ve))
				err := dbinary.BeDecode(ve, res)
				if err != nil {
					t.Errorf("test data: %s, %v, error:%v", k, v, err)
				}
				t.Assert(res, v)
			}
		}
	})
}

func Test_BeEncodeStruct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		user := User{"wenzi1", 999, "www.baidu.com"}
		ve := dbinary.BeEncode(user)
		s := dbinary.BeDecodeToString(ve)
		t.Assert(string(s), s)
	})
}
