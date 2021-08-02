// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
	"time"
)

type apiString interface {
	String() string
}
type S struct {
}

func (s S) String() string {
	return "22222"
}

type apiError interface {
	Error() string
}
type S1 struct {
}

func (s1 S1) Error() string {
	return "22222"
}

func Test_Bool_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.AssertEQ(dconv.Bool(any), false)
		t.AssertEQ(dconv.Bool(false), false)
		t.AssertEQ(dconv.Bool(nil), false)
		t.AssertEQ(dconv.Bool(0), false)
		t.AssertEQ(dconv.Bool("0"), false)
		t.AssertEQ(dconv.Bool(""), false)
		t.AssertEQ(dconv.Bool("false"), false)
		t.AssertEQ(dconv.Bool("off"), false)
		t.AssertEQ(dconv.Bool([]byte{}), false)
		t.AssertEQ(dconv.Bool([]string{}), false)
		t.AssertEQ(dconv.Bool([2]int{1, 2}), true)
		t.AssertEQ(dconv.Bool([]interface{}{}), false)
		t.AssertEQ(dconv.Bool([]map[int]int{}), false)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Bool(countryCapitalMap), true)

		t.AssertEQ(dconv.Bool("1"), true)
		t.AssertEQ(dconv.Bool("on"), true)
		t.AssertEQ(dconv.Bool(1), true)
		t.AssertEQ(dconv.Bool(123.456), true)
		t.AssertEQ(dconv.Bool(boolStruct{}), true)
		t.AssertEQ(dconv.Bool(&boolStruct{}), true)
	})
}

func Test_Int_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.AssertEQ(dconv.Int(any), 0)
		t.AssertEQ(dconv.Int(false), 0)
		t.AssertEQ(dconv.Int(nil), 0)
		t.Assert(dconv.Int(nil), 0)
		t.AssertEQ(dconv.Int(0), 0)
		t.AssertEQ(dconv.Int("0"), 0)
		t.AssertEQ(dconv.Int(""), 0)
		t.AssertEQ(dconv.Int("false"), 0)
		t.AssertEQ(dconv.Int("off"), 0)
		t.AssertEQ(dconv.Int([]byte{}), 0)
		t.AssertEQ(dconv.Int([]string{}), 0)
		t.AssertEQ(dconv.Int([2]int{1, 2}), 0)
		t.AssertEQ(dconv.Int([]interface{}{}), 0)
		t.AssertEQ(dconv.Int([]map[int]int{}), 0)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Int(countryCapitalMap), 0)

		t.AssertEQ(dconv.Int("1"), 1)
		t.AssertEQ(dconv.Int("on"), 0)
		t.AssertEQ(dconv.Int(1), 1)
		t.AssertEQ(dconv.Int(123.456), 123)
		t.AssertEQ(dconv.Int(boolStruct{}), 0)
		t.AssertEQ(dconv.Int(&boolStruct{}), 0)
	})
}

func Test_Int8_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Int8(any), int8(0))
		t.AssertEQ(dconv.Int8(false), int8(0))
		t.AssertEQ(dconv.Int8(nil), int8(0))
		t.AssertEQ(dconv.Int8(0), int8(0))
		t.AssertEQ(dconv.Int8("0"), int8(0))
		t.AssertEQ(dconv.Int8(""), int8(0))
		t.AssertEQ(dconv.Int8("false"), int8(0))
		t.AssertEQ(dconv.Int8("off"), int8(0))
		t.AssertEQ(dconv.Int8([]byte{}), int8(0))
		t.AssertEQ(dconv.Int8([]string{}), int8(0))
		t.AssertEQ(dconv.Int8([2]int{1, 2}), int8(0))
		t.AssertEQ(dconv.Int8([]interface{}{}), int8(0))
		t.AssertEQ(dconv.Int8([]map[int]int{}), int8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Int8(countryCapitalMap), int8(0))

		t.AssertEQ(dconv.Int8("1"), int8(1))
		t.AssertEQ(dconv.Int8("on"), int8(0))
		t.AssertEQ(dconv.Int8(int8(1)), int8(1))
		t.AssertEQ(dconv.Int8(123.456), int8(123))
		t.AssertEQ(dconv.Int8(boolStruct{}), int8(0))
		t.AssertEQ(dconv.Int8(&boolStruct{}), int8(0))
	})
}

func Test_Int16_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Int16(any), int16(0))
		t.AssertEQ(dconv.Int16(false), int16(0))
		t.AssertEQ(dconv.Int16(nil), int16(0))
		t.AssertEQ(dconv.Int16(0), int16(0))
		t.AssertEQ(dconv.Int16("0"), int16(0))
		t.AssertEQ(dconv.Int16(""), int16(0))
		t.AssertEQ(dconv.Int16("false"), int16(0))
		t.AssertEQ(dconv.Int16("off"), int16(0))
		t.AssertEQ(dconv.Int16([]byte{}), int16(0))
		t.AssertEQ(dconv.Int16([]string{}), int16(0))
		t.AssertEQ(dconv.Int16([2]int{1, 2}), int16(0))
		t.AssertEQ(dconv.Int16([]interface{}{}), int16(0))
		t.AssertEQ(dconv.Int16([]map[int]int{}), int16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Int16(countryCapitalMap), int16(0))

		t.AssertEQ(dconv.Int16("1"), int16(1))
		t.AssertEQ(dconv.Int16("on"), int16(0))
		t.AssertEQ(dconv.Int16(int16(1)), int16(1))
		t.AssertEQ(dconv.Int16(123.456), int16(123))
		t.AssertEQ(dconv.Int16(boolStruct{}), int16(0))
		t.AssertEQ(dconv.Int16(&boolStruct{}), int16(0))
	})
}

func Test_Int32_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Int32(any), int32(0))
		t.AssertEQ(dconv.Int32(false), int32(0))
		t.AssertEQ(dconv.Int32(nil), int32(0))
		t.AssertEQ(dconv.Int32(0), int32(0))
		t.AssertEQ(dconv.Int32("0"), int32(0))
		t.AssertEQ(dconv.Int32(""), int32(0))
		t.AssertEQ(dconv.Int32("false"), int32(0))
		t.AssertEQ(dconv.Int32("off"), int32(0))
		t.AssertEQ(dconv.Int32([]byte{}), int32(0))
		t.AssertEQ(dconv.Int32([]string{}), int32(0))
		t.AssertEQ(dconv.Int32([2]int{1, 2}), int32(0))
		t.AssertEQ(dconv.Int32([]interface{}{}), int32(0))
		t.AssertEQ(dconv.Int32([]map[int]int{}), int32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Int32(countryCapitalMap), int32(0))

		t.AssertEQ(dconv.Int32("1"), int32(1))
		t.AssertEQ(dconv.Int32("on"), int32(0))
		t.AssertEQ(dconv.Int32(int32(1)), int32(1))
		t.AssertEQ(dconv.Int32(123.456), int32(123))
		t.AssertEQ(dconv.Int32(boolStruct{}), int32(0))
		t.AssertEQ(dconv.Int32(&boolStruct{}), int32(0))
	})
}

func Test_Int64_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.AssertEQ(dconv.Int64("0x00e"), int64(14))
		t.Assert(dconv.Int64("022"), int64(18))

		t.Assert(dconv.Int64(any), int64(0))
		t.Assert(dconv.Int64(true), 1)
		t.Assert(dconv.Int64("1"), int64(1))
		t.Assert(dconv.Int64("0"), int64(0))
		t.Assert(dconv.Int64("X"), int64(0))
		t.Assert(dconv.Int64("x"), int64(0))
		t.Assert(dconv.Int64(int64(1)), int64(1))
		t.Assert(dconv.Int64(int(0)), int64(0))
		t.Assert(dconv.Int64(int8(0)), int64(0))
		t.Assert(dconv.Int64(int16(0)), int64(0))
		t.Assert(dconv.Int64(int32(0)), int64(0))
		t.Assert(dconv.Int64(uint64(0)), int64(0))
		t.Assert(dconv.Int64(uint32(0)), int64(0))
		t.Assert(dconv.Int64(uint16(0)), int64(0))
		t.Assert(dconv.Int64(uint8(0)), int64(0))
		t.Assert(dconv.Int64(uint(0)), int64(0))
		t.Assert(dconv.Int64(float32(0)), int64(0))

		t.AssertEQ(dconv.Int64(false), int64(0))
		t.AssertEQ(dconv.Int64(nil), int64(0))
		t.AssertEQ(dconv.Int64(0), int64(0))
		t.AssertEQ(dconv.Int64("0"), int64(0))
		t.AssertEQ(dconv.Int64(""), int64(0))
		t.AssertEQ(dconv.Int64("false"), int64(0))
		t.AssertEQ(dconv.Int64("off"), int64(0))
		t.AssertEQ(dconv.Int64([]byte{}), int64(0))
		t.AssertEQ(dconv.Int64([]string{}), int64(0))
		t.AssertEQ(dconv.Int64([2]int{1, 2}), int64(0))
		t.AssertEQ(dconv.Int64([]interface{}{}), int64(0))
		t.AssertEQ(dconv.Int64([]map[int]int{}), int64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Int64(countryCapitalMap), int64(0))

		t.AssertEQ(dconv.Int64("1"), int64(1))
		t.AssertEQ(dconv.Int64("on"), int64(0))
		t.AssertEQ(dconv.Int64(int64(1)), int64(1))
		t.AssertEQ(dconv.Int64(123.456), int64(123))
		t.AssertEQ(dconv.Int64(boolStruct{}), int64(0))
		t.AssertEQ(dconv.Int64(&boolStruct{}), int64(0))
	})
}

func Test_Uint_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.AssertEQ(dconv.Uint(any), uint(0))
		t.AssertEQ(dconv.Uint(false), uint(0))
		t.AssertEQ(dconv.Uint(nil), uint(0))
		t.Assert(dconv.Uint(nil), uint(0))
		t.AssertEQ(dconv.Uint(uint(0)), uint(0))
		t.AssertEQ(dconv.Uint("0"), uint(0))
		t.AssertEQ(dconv.Uint(""), uint(0))
		t.AssertEQ(dconv.Uint("false"), uint(0))
		t.AssertEQ(dconv.Uint("off"), uint(0))
		t.AssertEQ(dconv.Uint([]byte{}), uint(0))
		t.AssertEQ(dconv.Uint([]string{}), uint(0))
		t.AssertEQ(dconv.Uint([2]int{1, 2}), uint(0))
		t.AssertEQ(dconv.Uint([]interface{}{}), uint(0))
		t.AssertEQ(dconv.Uint([]map[int]int{}), uint(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Uint(countryCapitalMap), uint(0))

		t.AssertEQ(dconv.Uint("1"), uint(1))
		t.AssertEQ(dconv.Uint("on"), uint(0))
		t.AssertEQ(dconv.Uint(1), uint(1))
		t.AssertEQ(dconv.Uint(123.456), uint(123))
		t.AssertEQ(dconv.Uint(boolStruct{}), uint(0))
		t.AssertEQ(dconv.Uint(&boolStruct{}), uint(0))
	})
}

func Test_Uint8_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Uint8(any), uint8(0))
		t.AssertEQ(dconv.Uint8(uint8(1)), uint8(1))
		t.AssertEQ(dconv.Uint8(false), uint8(0))
		t.AssertEQ(dconv.Uint8(nil), uint8(0))
		t.AssertEQ(dconv.Uint8(0), uint8(0))
		t.AssertEQ(dconv.Uint8("0"), uint8(0))
		t.AssertEQ(dconv.Uint8(""), uint8(0))
		t.AssertEQ(dconv.Uint8("false"), uint8(0))
		t.AssertEQ(dconv.Uint8("off"), uint8(0))
		t.AssertEQ(dconv.Uint8([]byte{}), uint8(0))
		t.AssertEQ(dconv.Uint8([]string{}), uint8(0))
		t.AssertEQ(dconv.Uint8([2]int{1, 2}), uint8(0))
		t.AssertEQ(dconv.Uint8([]interface{}{}), uint8(0))
		t.AssertEQ(dconv.Uint8([]map[int]int{}), uint8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Uint8(countryCapitalMap), uint8(0))

		t.AssertEQ(dconv.Uint8("1"), uint8(1))
		t.AssertEQ(dconv.Uint8("on"), uint8(0))
		t.AssertEQ(dconv.Uint8(int8(1)), uint8(1))
		t.AssertEQ(dconv.Uint8(123.456), uint8(123))
		t.AssertEQ(dconv.Uint8(boolStruct{}), uint8(0))
		t.AssertEQ(dconv.Uint8(&boolStruct{}), uint8(0))
	})
}

func Test_Uint16_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Uint16(any), uint16(0))
		t.AssertEQ(dconv.Uint16(uint16(1)), uint16(1))
		t.AssertEQ(dconv.Uint16(false), uint16(0))
		t.AssertEQ(dconv.Uint16(nil), uint16(0))
		t.AssertEQ(dconv.Uint16(0), uint16(0))
		t.AssertEQ(dconv.Uint16("0"), uint16(0))
		t.AssertEQ(dconv.Uint16(""), uint16(0))
		t.AssertEQ(dconv.Uint16("false"), uint16(0))
		t.AssertEQ(dconv.Uint16("off"), uint16(0))
		t.AssertEQ(dconv.Uint16([]byte{}), uint16(0))
		t.AssertEQ(dconv.Uint16([]string{}), uint16(0))
		t.AssertEQ(dconv.Uint16([2]int{1, 2}), uint16(0))
		t.AssertEQ(dconv.Uint16([]interface{}{}), uint16(0))
		t.AssertEQ(dconv.Uint16([]map[int]int{}), uint16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Uint16(countryCapitalMap), uint16(0))

		t.AssertEQ(dconv.Uint16("1"), uint16(1))
		t.AssertEQ(dconv.Uint16("on"), uint16(0))
		t.AssertEQ(dconv.Uint16(int16(1)), uint16(1))
		t.AssertEQ(dconv.Uint16(123.456), uint16(123))
		t.AssertEQ(dconv.Uint16(boolStruct{}), uint16(0))
		t.AssertEQ(dconv.Uint16(&boolStruct{}), uint16(0))
	})
}

func Test_Uint32_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Uint32(any), uint32(0))
		t.AssertEQ(dconv.Uint32(uint32(1)), uint32(1))
		t.AssertEQ(dconv.Uint32(false), uint32(0))
		t.AssertEQ(dconv.Uint32(nil), uint32(0))
		t.AssertEQ(dconv.Uint32(0), uint32(0))
		t.AssertEQ(dconv.Uint32("0"), uint32(0))
		t.AssertEQ(dconv.Uint32(""), uint32(0))
		t.AssertEQ(dconv.Uint32("false"), uint32(0))
		t.AssertEQ(dconv.Uint32("off"), uint32(0))
		t.AssertEQ(dconv.Uint32([]byte{}), uint32(0))
		t.AssertEQ(dconv.Uint32([]string{}), uint32(0))
		t.AssertEQ(dconv.Uint32([2]int{1, 2}), uint32(0))
		t.AssertEQ(dconv.Uint32([]interface{}{}), uint32(0))
		t.AssertEQ(dconv.Uint32([]map[int]int{}), uint32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Uint32(countryCapitalMap), uint32(0))

		t.AssertEQ(dconv.Uint32("1"), uint32(1))
		t.AssertEQ(dconv.Uint32("on"), uint32(0))
		t.AssertEQ(dconv.Uint32(int32(1)), uint32(1))
		t.AssertEQ(dconv.Uint32(123.456), uint32(123))
		t.AssertEQ(dconv.Uint32(boolStruct{}), uint32(0))
		t.AssertEQ(dconv.Uint32(&boolStruct{}), uint32(0))
	})
}

func Test_Uint64_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.AssertEQ(dconv.Uint64("0x00e"), uint64(14))
		t.Assert(dconv.Uint64("022"), uint64(18))

		t.AssertEQ(dconv.Uint64(any), uint64(0))
		t.AssertEQ(dconv.Uint64(true), uint64(1))
		t.Assert(dconv.Uint64("1"), int64(1))
		t.Assert(dconv.Uint64("0"), uint64(0))
		t.Assert(dconv.Uint64("X"), uint64(0))
		t.Assert(dconv.Uint64("x"), uint64(0))
		t.Assert(dconv.Uint64(int64(1)), uint64(1))
		t.Assert(dconv.Uint64(int(0)), uint64(0))
		t.Assert(dconv.Uint64(int8(0)), uint64(0))
		t.Assert(dconv.Uint64(int16(0)), uint64(0))
		t.Assert(dconv.Uint64(int32(0)), uint64(0))
		t.Assert(dconv.Uint64(uint64(0)), uint64(0))
		t.Assert(dconv.Uint64(uint32(0)), uint64(0))
		t.Assert(dconv.Uint64(uint16(0)), uint64(0))
		t.Assert(dconv.Uint64(uint8(0)), uint64(0))
		t.Assert(dconv.Uint64(uint(0)), uint64(0))
		t.Assert(dconv.Uint64(float32(0)), uint64(0))

		t.AssertEQ(dconv.Uint64(false), uint64(0))
		t.AssertEQ(dconv.Uint64(nil), uint64(0))
		t.AssertEQ(dconv.Uint64(0), uint64(0))
		t.AssertEQ(dconv.Uint64("0"), uint64(0))
		t.AssertEQ(dconv.Uint64(""), uint64(0))
		t.AssertEQ(dconv.Uint64("false"), uint64(0))
		t.AssertEQ(dconv.Uint64("off"), uint64(0))
		t.AssertEQ(dconv.Uint64([]byte{}), uint64(0))
		t.AssertEQ(dconv.Uint64([]string{}), uint64(0))
		t.AssertEQ(dconv.Uint64([2]int{1, 2}), uint64(0))
		t.AssertEQ(dconv.Uint64([]interface{}{}), uint64(0))
		t.AssertEQ(dconv.Uint64([]map[int]int{}), uint64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Uint64(countryCapitalMap), uint64(0))

		t.AssertEQ(dconv.Uint64("1"), uint64(1))
		t.AssertEQ(dconv.Uint64("on"), uint64(0))
		t.AssertEQ(dconv.Uint64(int64(1)), uint64(1))
		t.AssertEQ(dconv.Uint64(123.456), uint64(123))
		t.AssertEQ(dconv.Uint64(boolStruct{}), uint64(0))
		t.AssertEQ(dconv.Uint64(&boolStruct{}), uint64(0))
	})
}

func Test_Float32_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Float32(any), float32(0))
		t.AssertEQ(dconv.Float32(false), float32(0))
		t.AssertEQ(dconv.Float32(nil), float32(0))
		t.AssertEQ(dconv.Float32(0), float32(0))
		t.AssertEQ(dconv.Float32("0"), float32(0))
		t.AssertEQ(dconv.Float32(""), float32(0))
		t.AssertEQ(dconv.Float32("false"), float32(0))
		t.AssertEQ(dconv.Float32("off"), float32(0))
		t.AssertEQ(dconv.Float32([]byte{}), float32(0))
		t.AssertEQ(dconv.Float32([]string{}), float32(0))
		t.AssertEQ(dconv.Float32([2]int{1, 2}), float32(0))
		t.AssertEQ(dconv.Float32([]interface{}{}), float32(0))
		t.AssertEQ(dconv.Float32([]map[int]int{}), float32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Float32(countryCapitalMap), float32(0))

		t.AssertEQ(dconv.Float32("1"), float32(1))
		t.AssertEQ(dconv.Float32("on"), float32(0))
		t.AssertEQ(dconv.Float32(float32(1)), float32(1))
		t.AssertEQ(dconv.Float32(123.456), float32(123.456))
		t.AssertEQ(dconv.Float32(boolStruct{}), float32(0))
		t.AssertEQ(dconv.Float32(&boolStruct{}), float32(0))
	})
}

func Test_Float64_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.Assert(dconv.Float64(any), float64(0))
		t.AssertEQ(dconv.Float64(false), float64(0))
		t.AssertEQ(dconv.Float64(nil), float64(0))
		t.AssertEQ(dconv.Float64(0), float64(0))
		t.AssertEQ(dconv.Float64("0"), float64(0))
		t.AssertEQ(dconv.Float64(""), float64(0))
		t.AssertEQ(dconv.Float64("false"), float64(0))
		t.AssertEQ(dconv.Float64("off"), float64(0))
		t.AssertEQ(dconv.Float64([]byte{}), float64(0))
		t.AssertEQ(dconv.Float64([]string{}), float64(0))
		t.AssertEQ(dconv.Float64([2]int{1, 2}), float64(0))
		t.AssertEQ(dconv.Float64([]interface{}{}), float64(0))
		t.AssertEQ(dconv.Float64([]map[int]int{}), float64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.Float64(countryCapitalMap), float64(0))

		t.AssertEQ(dconv.Float64("1"), float64(1))
		t.AssertEQ(dconv.Float64("on"), float64(0))
		t.AssertEQ(dconv.Float64(float64(1)), float64(1))
		t.AssertEQ(dconv.Float64(123.456), float64(123.456))
		t.AssertEQ(dconv.Float64(boolStruct{}), float64(0))
		t.AssertEQ(dconv.Float64(&boolStruct{}), float64(0))
	})
}

func Test_String_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var s []rune
		t.AssertEQ(dconv.String(s), "")
		var any interface{} = nil
		t.AssertEQ(dconv.String(any), "")
		t.AssertEQ(dconv.String("1"), "1")
		t.AssertEQ(dconv.String("0"), string("0"))
		t.Assert(dconv.String("X"), string("X"))
		t.Assert(dconv.String("x"), string("x"))
		t.Assert(dconv.String(int64(1)), uint64(1))
		t.Assert(dconv.String(int(0)), string("0"))
		t.Assert(dconv.String(int8(0)), string("0"))
		t.Assert(dconv.String(int16(0)), string("0"))
		t.Assert(dconv.String(int32(0)), string("0"))
		t.Assert(dconv.String(uint64(0)), string("0"))
		t.Assert(dconv.String(uint32(0)), string("0"))
		t.Assert(dconv.String(uint16(0)), string("0"))
		t.Assert(dconv.String(uint8(0)), string("0"))
		t.Assert(dconv.String(uint(0)), string("0"))
		t.Assert(dconv.String(float32(0)), string("0"))
		t.AssertEQ(dconv.String(true), "true")
		t.AssertEQ(dconv.String(false), "false")
		t.AssertEQ(dconv.String(nil), "")
		t.AssertEQ(dconv.String(0), string("0"))
		t.AssertEQ(dconv.String("0"), string("0"))
		t.AssertEQ(dconv.String(""), "")
		t.AssertEQ(dconv.String("false"), "false")
		t.AssertEQ(dconv.String("off"), string("off"))
		t.AssertEQ(dconv.String([]byte{}), "")
		t.AssertEQ(dconv.String([]string{}), "[]")
		t.AssertEQ(dconv.String([2]int{1, 2}), "[1,2]")
		t.AssertEQ(dconv.String([]interface{}{}), "[]")
		t.AssertEQ(dconv.String(map[int]int{}), "{}")

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(dconv.String(countryCapitalMap), `{"France":"巴黎","India ":"新德里","Italy":"罗马","Japan":"东京"}`)
		t.AssertEQ(dconv.String(int64(1)), "1")
		t.AssertEQ(dconv.String(123.456), "123.456")
		t.AssertEQ(dconv.String(boolStruct{}), "{}")
		t.AssertEQ(dconv.String(&boolStruct{}), "{}")

		var info apiString
		info = new(S)
		t.AssertEQ(dconv.String(info), "22222")
		var errinfo apiError
		errinfo = new(S1)
		t.AssertEQ(dconv.String(errinfo), "22222")
	})
}

func Test_Runes_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Runes("www"), []int32{119, 119, 119})
		var s []rune
		t.AssertEQ(dconv.Runes(s), nil)
	})
}

func Test_Rune_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Rune("www"), int32(0))
		t.AssertEQ(dconv.Rune(int32(0)), int32(0))
		var s []rune
		t.AssertEQ(dconv.Rune(s), int32(0))
	})
}

func Test_Bytes_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Bytes(nil), nil)
		t.AssertEQ(dconv.Bytes(int32(0)), []uint8{0, 0, 0, 0})
		t.AssertEQ(dconv.Bytes("s"), []uint8{115})
		t.AssertEQ(dconv.Bytes([]byte("s")), []uint8{115})
	})
}

func Test_Byte_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Byte(uint8(0)), uint8(0))
		t.AssertEQ(dconv.Byte("s"), uint8(0))
		t.AssertEQ(dconv.Byte([]byte("s")), uint8(115))
	})
}

func Test_Convert_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.AssertEQ(dconv.Convert(any, "string"), "")
		t.AssertEQ(dconv.Convert("1", "string"), "1")
		t.Assert(dconv.Convert(int64(1), "int64"), int64(1))
		t.Assert(dconv.Convert(int(0), "int"), int(0))
		t.Assert(dconv.Convert(int8(0), "int8"), int8(0))
		t.Assert(dconv.Convert(int16(0), "int16"), int16(0))
		t.Assert(dconv.Convert(int32(0), "int32"), int32(0))
		t.Assert(dconv.Convert(uint64(0), "uint64"), uint64(0))
		t.Assert(dconv.Convert(uint32(0), "uint32"), uint32(0))
		t.Assert(dconv.Convert(uint16(0), "uint16"), uint16(0))
		t.Assert(dconv.Convert(uint8(0), "uint8"), uint8(0))
		t.Assert(dconv.Convert(uint(0), "uint"), uint(0))
		t.Assert(dconv.Convert(float32(0), "float32"), float32(0))
		t.Assert(dconv.Convert(float64(0), "float64"), float64(0))
		t.AssertEQ(dconv.Convert(true, "bool"), true)
		t.AssertEQ(dconv.Convert([]byte{}, "[]byte"), []uint8{})
		t.AssertEQ(dconv.Convert([]string{}, "[]string"), []string{})
		t.AssertEQ(dconv.Convert([2]int{1, 2}, "[]int"), []int{1, 2})
		t.AssertEQ(dconv.Convert("1989-01-02", "Time", "Y-m-d"), dconv.Time("1989-01-02", "Y-m-d"))
		t.AssertEQ(dconv.Convert(1989, "Time"), dconv.Time("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(dconv.Convert(dtime.Now(), "dtime.Time", 1), *dtime.New())
		t.AssertEQ(dconv.Convert(1989, "dtime.Time"), *dconv.GTime("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(dconv.Convert(dtime.Now(), "*dtime.Time", 1), dtime.New())
		t.AssertEQ(dconv.Convert(dtime.Now(), "dtime", 1), *dtime.New())
		t.AssertEQ(dconv.Convert(1989, "*dtime.Time"), dconv.GTime(1989))
		t.AssertEQ(dconv.Convert(1989, "Duration"), time.Duration(int64(1989)))
		t.AssertEQ(dconv.Convert("1989", "Duration"), time.Duration(int64(1989)))
		t.AssertEQ(dconv.Convert("1989", ""), "1989")
	})
}

func Test_Slice_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := 123.456
		t.AssertEQ(dconv.Ints(value), []int{123})
		t.AssertEQ(dconv.Ints(nil), nil)
		t.AssertEQ(dconv.Ints([]string{"1", "2"}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]int{}), []int{})
		t.AssertEQ(dconv.Ints([]int8{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]int16{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]int32{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]int64{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]uint{1}), []int{1})
		t.AssertEQ(dconv.Ints([]uint8{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]uint16{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]uint32{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]uint64{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]bool{true}), []int{1})
		t.AssertEQ(dconv.Ints([]float32{1, 2}), []int{1, 2})
		t.AssertEQ(dconv.Ints([]float64{1, 2}), []int{1, 2})
		var inter []interface{} = make([]interface{}, 2)
		t.AssertEQ(dconv.Ints(inter), []int{0, 0})

		t.AssertEQ(dconv.Strings(value), []string{"123.456"})
		t.AssertEQ(dconv.Strings(nil), nil)
		t.AssertEQ(dconv.Strings([]string{"1", "2"}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]int{1}), []string{"1"})
		t.AssertEQ(dconv.Strings([]int8{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]int16{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]int32{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]int64{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]uint{1}), []string{"1"})
		t.AssertEQ(dconv.Strings([]uint8{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]uint16{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]uint32{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]uint64{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]bool{true}), []string{"true"})
		t.AssertEQ(dconv.Strings([]float32{1, 2}), []string{"1", "2"})
		t.AssertEQ(dconv.Strings([]float64{1, 2}), []string{"1", "2"})
		var strer = make([]interface{}, 2)
		t.AssertEQ(dconv.Strings(strer), []string{"", ""})

		t.AssertEQ(dconv.Floats(value), []float64{123.456})
		t.AssertEQ(dconv.Floats(nil), nil)
		t.AssertEQ(dconv.Floats([]string{"1", "2"}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]int{1}), []float64{1})
		t.AssertEQ(dconv.Floats([]int8{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]int16{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]int32{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]int64{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]uint{1}), []float64{1})
		t.AssertEQ(dconv.Floats([]uint8{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]uint16{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]uint32{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]uint64{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]bool{true}), []float64{0})
		t.AssertEQ(dconv.Floats([]float32{1, 2}), []float64{1, 2})
		t.AssertEQ(dconv.Floats([]float64{1, 2}), []float64{1, 2})
		var floer = make([]interface{}, 2)
		t.AssertEQ(dconv.Floats(floer), []float64{0, 0})

		t.AssertEQ(dconv.Interfaces(value), []interface{}{123.456})
		t.AssertEQ(dconv.Interfaces(nil), nil)
		t.AssertEQ(dconv.Interfaces([]interface{}{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]string{"1"}), []interface{}{"1"})
		t.AssertEQ(dconv.Interfaces([]int{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]int8{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]int16{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]int32{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]int64{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]uint{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]uint8{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]uint16{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]uint32{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]uint64{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]bool{true}), []interface{}{true})
		t.AssertEQ(dconv.Interfaces([]float32{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([]float64{1}), []interface{}{1})
		t.AssertEQ(dconv.Interfaces([1]int{1}), []interface{}{1})

		type interSlice []int
		slices := interSlice{1}
		t.AssertEQ(dconv.Interfaces(slices), []interface{}{1})

		t.AssertEQ(dconv.Maps(nil), nil)
		t.AssertEQ(dconv.Maps([]map[string]interface{}{{"a": "1"}}), []map[string]interface{}{{"a": "1"}})
		t.AssertEQ(dconv.Maps(1223), []map[string]interface{}{nil})
		t.AssertEQ(dconv.Maps([]int{}), nil)
	})
}

// 私有属性不会进行转换
func Test_Slice_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int           `json:"id"`
		name string        `json:"name"`
		Ad   []interface{} `json:"ad"`
	}
	dtest.C(t, func(t *dtest.T) {
		user := &User{1, "john", []interface{}{2}}
		array := dconv.Interfaces(user)
		t.Assert(len(array), 1)
		t.Assert(array[0].(*User).Id, 1)
		t.Assert(array[0].(*User).name, "john")
		t.Assert(array[0].(*User).Ad, []interface{}{2})
	})
}

func Test_Map_Basic_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m1 := map[string]string{
			"k": "v",
		}
		m2 := map[int]string{
			3: "v",
		}
		m3 := map[float64]float32{
			1.22: 3.1,
		}
		t.Assert(dconv.Map(m1), d.Map{
			"k": "v",
		})
		t.Assert(dconv.Map(m2), d.Map{
			"3": "v",
		})
		t.Assert(dconv.Map(m3), d.Map{
			"1.22": "3.1",
		})
		t.AssertEQ(dconv.Map(nil), nil)
		t.AssertEQ(dconv.Map(map[string]interface{}{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(dconv.Map(map[int]interface{}{1: 1}), map[string]interface{}{"1": 1})
		t.AssertEQ(dconv.Map(map[uint]interface{}{1: 1}), map[string]interface{}{"1": 1})
		t.AssertEQ(dconv.Map(map[uint]string{1: "1"}), map[string]interface{}{"1": "1"})

		t.AssertEQ(dconv.Map(map[interface{}]interface{}{"a": 1}), map[interface{}]interface{}{"a": 1})
		t.AssertEQ(dconv.Map(map[interface{}]string{"a": "1"}), map[interface{}]string{"a": "1"})
		t.AssertEQ(dconv.Map(map[interface{}]int{"a": 1}), map[interface{}]int{"a": 1})
		t.AssertEQ(dconv.Map(map[interface{}]uint{"a": 1}), map[interface{}]uint{"a": 1})
		t.AssertEQ(dconv.Map(map[interface{}]float32{"a": 1}), map[interface{}]float32{"a": 1})
		t.AssertEQ(dconv.Map(map[interface{}]float64{"a": 1}), map[interface{}]float64{"a": 1})

		t.AssertEQ(dconv.Map(map[string]bool{"a": true}), map[string]interface{}{"a": true})
		t.AssertEQ(dconv.Map(map[string]int{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(dconv.Map(map[string]uint{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(dconv.Map(map[string]float32{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(dconv.Map(map[string]float64{"a": 1}), map[string]interface{}{"a": 1})

	})
}

func Test_Map_StructWithdconvTag_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string   `dconv:"-"`
			NickName string   `dconv:"nickname,omitempty"`
			Pass1    string   `dconv:"password1"`
			Pass2    string   `dconv:"password2"`
			Ss       []string `dconv:"ss"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
			Ss:      []string{"sss", "2222"},
		}
		user2 := &user1
		map1 := dconv.Map(user1)
		map2 := dconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")
		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithJsonTag_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string   `json:"-"`
			NickName string   `json:"nickname, omitempty"`
			Pass1    string   `json:"password1,newpassword"`
			Pass2    string   `json:"password2"`
			Ss       []string `json:"omitempty"`
			ssb, ssa string
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
			Ss:      []string{"sss", "2222"},
			ssb:     "11",
			ssa:     "222",
		}
		user3 := User{
			Uid:      100,
			Name:     "john",
			NickName: "SSS",
			SiteUrl:  "https://goframe.org",
			Pass1:    "123",
			Pass2:    "456",
			Ss:       []string{"sss", "2222"},
			ssb:      "11",
			ssa:      "222",
		}
		user2 := &user1
		_ = dconv.Map(user1, "Ss")
		map1 := dconv.Map(user1, "json", "json2")
		map2 := dconv.Map(user2)
		map3 := dconv.Map(user3)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")
		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
		t.Assert(map3["NickName"], nil)
	})
}

func Test_Map_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	dtest.C(t, func(t *dtest.T) {
		user := &User{1, "john"}
		t.Assert(dconv.Map(user), d.Map{"Id": 1})
	})
}

func Test_Map_StructInherit_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string  `json:"passport"`
			Password string  `json:"password"`
			Nickname string  `json:"nickname"`
			S        *string `json:"nickname2"`
		}

		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		var s = "s"
		user.S = &s

		m := dconv.MapDeep(user)
		t.Assert(m["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], user.CreateTime)
		t.Assert(m["nickname2"], user.S)
	})
}

func Test_Struct_Basic1_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Score struct {
			Name   int
			Result string
		}

		type Score2 struct {
			Name   int
			Result string
		}

		type User struct {
			Uid      int
			Name     string
			Site_Url string
			NickName string
			Pass1    string `dconv:"password1"`
			Pass2    string `dconv:"password2"`
			As       *Score
			Ass      Score
			Assb     []interface{}
		}
		// 使用默认映射规则绑定属性值到对象
		user := new(User)
		params1 := d.Map{
			"uid":       1,
			"Name":      "john",
			"siteurl":   "https://goframe.org",
			"nick_name": "johng",
			"PASS1":     "123",
			"PASS2":     "456",
			"As":        d.Map{"Name": 1, "Result": "22222"},
			"Ass":       &Score{11, "11"},
			"Assb":      []string{"wwww"},
		}
		_ = dconv.Struct(nil, user)
		_ = dconv.Struct(params1, nil)
		_ = dconv.Struct([]interface{}{nil}, user)
		_ = dconv.Struct(user, []interface{}{nil})

		var a = []interface{}{nil}
		ab := &a
		_ = dconv.Struct(params1, *ab)
		var pi *int = nil
		_ = dconv.Struct(params1, pi)

		_ = dconv.Struct(params1, user)
		_ = dconv.Struct(params1, user, map[string]string{"uid": "Names"})
		_ = dconv.Struct(params1, user, map[string]string{"uid": "as"})

		// 使用struct tag映射绑定属性值到对象
		user = new(User)
		params2 := d.Map{
			"uid":       2,
			"name":      "smith",
			"site-url":  "https://goframe.org",
			"nick name": "johng",
			"password1": "111",
			"password2": "222",
		}
		if err := dconv.Struct(params2, user); err != nil {
			dtest.Error(err)
		}
		t.Assert(user, &User{
			Uid:      2,
			Name:     "smith",
			Site_Url: "https://goframe.org",
			NickName: "johng",
			Pass1:    "111",
			Pass2:    "222",
		})
	})
}

// 使用默认映射规则绑定属性值到对象
func Test_Struct_Basic2_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid     int
			Name    string
			SiteUrl string
			Pass1   string
			Pass2   string
		}
		user := new(User)
		params := d.Map{
			"uid":      1,
			"Name":     "john",
			"site_url": "https://goframe.org",
			"PASS1":    "123",
			"PASS2":    "456",
		}
		if err := dconv.Struct(params, user); err != nil {
			dtest.Error(err)
		}
		t.Assert(user, &User{
			Uid:     1,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		})
	})
}

// 带有指针的基础类型属性
func Test_Struct_Basic3_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid  int
			Name *string
		}
		user := new(User)
		params := d.Map{
			"uid":  1,
			"Name": "john",
		}
		if err := dconv.Struct(params, user); err != nil {
			dtest.Error(err)
		}
		t.Assert(user.Uid, 1)
		t.Assert(*user.Name, "john")
	})
}

// slice类型属性的赋值
func Test_Struct_Attr_Slice_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Scores []int
		}
		scores := []interface{}{99, 100, 60, 140}
		user := new(User)
		if err := dconv.Struct(d.Map{"Scores": scores}, user); err != nil {
			dtest.Error(err)
		} else {
			t.Assert(user, &User{
				Scores: []int{99, 100, 60, 140},
			})
		}
	})
}

// 属性为struct对象
func Test_Struct_Attr_Struct_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换
		if err := dconv.Struct(scores, user); err != nil {
			dtest.Error(err)
		} else {
			t.Assert(user, &User{
				Scores: Score{
					Name:   "john",
					Result: 100,
				},
			})
		}
	})
}

// 属性为struct对象指针
func Test_Struct_Attr_Struct_Ptr_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores *Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换
		if err := dconv.Struct(scores, user); err != nil {
			dtest.Error(err)
		} else {
			t.Assert(user.Scores, &Score{
				Name:   "john",
				Result: 100,
			})
		}
	})
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice1_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换，属性为slice类型，数值为map类型
		if err := dconv.Struct(scores, user); err != nil {
			dtest.Error(err)
		} else {
			t.Assert(user.Scores, []Score{
				{
					Name:   "john",
					Result: 100,
				},
			})
		}
	})
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice2_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": []interface{}{
				map[string]interface{}{
					"Name":   "john",
					"Result": 100,
				},
				map[string]interface{}{
					"Name":   "smith",
					"Result": 60,
				},
			},
		}

		// 嵌套struct转换，属性为slice类型，数值为slice map类型
		if err := dconv.Struct(scores, user); err != nil {
			dtest.Error(err)
		} else {
			t.Assert(user.Scores, []Score{
				{
					Name:   "john",
					Result: 100,
				},
				{
					Name:   "smith",
					Result: 60,
				},
			})
		}
	})
}

// 属性为struct对象slice ptr
func Test_Struct_Attr_Struct_Slice_Ptr_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []*Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": []interface{}{
				map[string]interface{}{
					"Name":   "john",
					"Result": 100,
				},
				map[string]interface{}{
					"Name":   "smith",
					"Result": 60,
				},
			},
		}

		// 嵌套struct转换，属性为slice类型，数值为slice map类型
		if err := dconv.Struct(scores, user); err != nil {
			dtest.Error(err)
		} else {
			t.Assert(len(user.Scores), 2)
			t.Assert(user.Scores[0], &Score{
				Name:   "john",
				Result: 100,
			})
			t.Assert(user.Scores[1], &Score{
				Name:   "smith",
				Result: 60,
			})
		}
	})
}

func Test_Struct_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		err := dconv.Struct(d.Map{"id": 1, "name": "john"}, user)
		t.Assert(err, nil)
		t.Assert(user.Id, 1)
		t.Assert(user.name, "")
	})
}

func Test_Struct_Embedded_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		data := d.Map{
			"id":          100,
			"uid":         101,
			"passport":    "t1",
			"password":    "123456",
			"nickname":    "T1",
			"create_time": "2019",
		}
		user := new(User)
		dconv.Struct(data, user)
		t.Assert(user.Id, 100)
		t.Assert(user.Uid, 101)
		t.Assert(user.Nickname, "T1")
		t.Assert(user.CreateTime, "2019")
	})
}

func Test_Struct_Time_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			CreateTime time.Time
		}
		now := time.Now()
		user := new(User)
		dconv.Struct(d.Map{
			"create_time": now,
		}, user)
		t.Assert(user.CreateTime.UTC().String(), now.UTC().String())
	})

	dtest.C(t, func(t *dtest.T) {
		type User struct {
			CreateTime *time.Time
		}
		now := time.Now()
		user := new(User)
		dconv.Struct(d.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.UTC().String(), now.UTC().String())
	})

	dtest.C(t, func(t *dtest.T) {
		type User struct {
			CreateTime *dtime.Time
		}
		now := time.Now()
		user := new(User)
		dconv.Struct(d.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})

	dtest.C(t, func(t *dtest.T) {
		type User struct {
			CreateTime dtime.Time
		}
		now := time.Now()
		user := new(User)
		dconv.Struct(d.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})

	dtest.C(t, func(t *dtest.T) {
		type User struct {
			CreateTime dtime.Time
		}
		now := time.Now()
		user := new(User)
		dconv.Struct(d.Map{
			"create_time": now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})
}
