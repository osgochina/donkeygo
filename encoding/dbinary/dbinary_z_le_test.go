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

func Test_LeEncodeAndLeDecode(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for k, v := range testData {
			ve := dbinary.LeEncode(v)
			ve1 := dbinary.LeEncodeByLength(len(ve), v)

			//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
			switch v.(type) {
			case int:
				t.Assert(dbinary.LeDecodeToInt(ve), v)
				t.Assert(dbinary.LeDecodeToInt(ve1), v)
			case int8:
				t.Assert(dbinary.LeDecodeToInt8(ve), v)
				t.Assert(dbinary.LeDecodeToInt8(ve1), v)
			case int16:
				t.Assert(dbinary.LeDecodeToInt16(ve), v)
				t.Assert(dbinary.LeDecodeToInt16(ve1), v)
			case int32:
				t.Assert(dbinary.LeDecodeToInt32(ve), v)
				t.Assert(dbinary.LeDecodeToInt32(ve1), v)
			case int64:
				t.Assert(dbinary.LeDecodeToInt64(ve), v)
				t.Assert(dbinary.LeDecodeToInt64(ve1), v)
			case uint:
				t.Assert(dbinary.LeDecodeToUint(ve), v)
				t.Assert(dbinary.LeDecodeToUint(ve1), v)
			case uint8:
				t.Assert(dbinary.LeDecodeToUint8(ve), v)
				t.Assert(dbinary.LeDecodeToUint8(ve1), v)
			case uint16:
				t.Assert(dbinary.LeDecodeToUint16(ve1), v)
				t.Assert(dbinary.LeDecodeToUint16(ve), v)
			case uint32:
				t.Assert(dbinary.LeDecodeToUint32(ve1), v)
				t.Assert(dbinary.LeDecodeToUint32(ve), v)
			case uint64:
				t.Assert(dbinary.LeDecodeToUint64(ve), v)
				t.Assert(dbinary.LeDecodeToUint64(ve1), v)
			case bool:
				t.Assert(dbinary.LeDecodeToBool(ve), v)
				t.Assert(dbinary.LeDecodeToBool(ve1), v)
			case string:
				t.Assert(dbinary.LeDecodeToString(ve), v)
				t.Assert(dbinary.LeDecodeToString(ve1), v)
			case float32:
				t.Assert(dbinary.LeDecodeToFloat32(ve), v)
				t.Assert(dbinary.LeDecodeToFloat32(ve1), v)
			case float64:
				t.Assert(dbinary.LeDecodeToFloat64(ve), v)
				t.Assert(dbinary.LeDecodeToFloat64(ve1), v)
			default:
				if v == nil {
					continue
				}
				res := make([]byte, len(ve))
				err := dbinary.LeDecode(ve, res)
				if err != nil {
					t.Errorf("test data: %s, %v, error:%v", k, v, err)
				}
				t.Assert(res, v)
			}
		}
	})
}

func Test_LeEncodeStruct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		user := User{"wenzi1", 999, "www.baidu.com"}
		ve := dbinary.LeEncode(user)
		s := dbinary.LeDecodeToString(ve)
		t.Assert(s, s)
	})
}
