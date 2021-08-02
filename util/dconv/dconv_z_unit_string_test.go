// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

type stringStruct1 struct {
	Name string
}

type stringStruct2 struct {
	Name string
}

func (s *stringStruct1) String() string {
	return s.Name
}

func Test_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.String(int(123)), "123")
		t.AssertEQ(dconv.String(int(-123)), "-123")
		t.AssertEQ(dconv.String(int8(123)), "123")
		t.AssertEQ(dconv.String(int8(-123)), "-123")
		t.AssertEQ(dconv.String(int16(123)), "123")
		t.AssertEQ(dconv.String(int16(-123)), "-123")
		t.AssertEQ(dconv.String(int32(123)), "123")
		t.AssertEQ(dconv.String(int32(-123)), "-123")
		t.AssertEQ(dconv.String(int64(123)), "123")
		t.AssertEQ(dconv.String(int64(-123)), "-123")
		t.AssertEQ(dconv.String(int64(1552578474888)), "1552578474888")
		t.AssertEQ(dconv.String(int64(-1552578474888)), "-1552578474888")

		t.AssertEQ(dconv.String(uint(123)), "123")
		t.AssertEQ(dconv.String(uint8(123)), "123")
		t.AssertEQ(dconv.String(uint16(123)), "123")
		t.AssertEQ(dconv.String(uint32(123)), "123")
		t.AssertEQ(dconv.String(uint64(155257847488898765)), "155257847488898765")

		t.AssertEQ(dconv.String(float32(123.456)), "123.456")
		t.AssertEQ(dconv.String(float32(-123.456)), "-123.456")
		t.AssertEQ(dconv.String(float64(1552578474888.456)), "1552578474888.456")
		t.AssertEQ(dconv.String(float64(-1552578474888.456)), "-1552578474888.456")

		t.AssertEQ(dconv.String(true), "true")
		t.AssertEQ(dconv.String(false), "false")

		t.AssertEQ(dconv.String([]byte("bytes")), "bytes")

		t.AssertEQ(dconv.String(stringStruct1{"john"}), `{"Name":"john"}`)
		t.AssertEQ(dconv.String(&stringStruct1{"john"}), "john")

		t.AssertEQ(dconv.String(stringStruct2{"john"}), `{"Name":"john"}`)
		t.AssertEQ(dconv.String(&stringStruct2{"john"}), `{"Name":"john"}`)
	})
}
