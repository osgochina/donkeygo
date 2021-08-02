// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package empty_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

type TestInt int

type TestString string

type TestPerson interface {
	Say() string
}
type TestWoman struct {
}

func (woman TestWoman) Say() string {
	return "nice"
}

func TestIsEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		tmpT1 := "0"
		tmpT2 := func() {}
		tmpT2 = nil
		tmpT3 := make(chan int, 0)
		var (
			tmpT4 TestPerson  = nil
			tmpT5 *TestPerson = nil
			tmpT6 TestPerson  = TestWoman{}
			tmpT7 TestInt     = 0
			tmpT8 TestString  = ""
		)
		tmpF1 := "1"
		tmpF2 := func(a string) string { return "1" }
		tmpF3 := make(chan int, 1)
		tmpF3 <- 1
		var (
			tmpF4 TestPerson = &TestWoman{}
			tmpF5 TestInt    = 1
			tmpF6 TestString = "1"
		)

		// true
		t.Assert(empty.IsEmpty(nil), true)
		t.Assert(empty.IsEmpty(0), true)
		t.Assert(empty.IsEmpty(dconv.Int(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Int8(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Int16(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Int32(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Int64(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Uint64(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Uint(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Uint16(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Uint32(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Uint64(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Float32(tmpT1)), true)
		t.Assert(empty.IsEmpty(dconv.Float64(tmpT1)), true)
		t.Assert(empty.IsEmpty(false), true)
		t.Assert(empty.IsEmpty([]byte("")), true)
		t.Assert(empty.IsEmpty(""), true)
		t.Assert(empty.IsEmpty(d.Map{}), true)
		t.Assert(empty.IsEmpty(d.Slice{}), true)
		t.Assert(empty.IsEmpty(d.Array{}), true)
		t.Assert(empty.IsEmpty(tmpT2), true)
		t.Assert(empty.IsEmpty(tmpT3), true)
		t.Assert(empty.IsEmpty(tmpT3), true)
		t.Assert(empty.IsEmpty(tmpT4), true)
		t.Assert(empty.IsEmpty(tmpT5), true)
		t.Assert(empty.IsEmpty(tmpT6), true)
		t.Assert(empty.IsEmpty(tmpT7), true)
		t.Assert(empty.IsEmpty(tmpT8), true)

		// false
		t.Assert(empty.IsEmpty(dconv.Int(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Int8(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Int16(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Int32(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Int64(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Uint(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Uint8(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Uint16(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Uint32(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Uint64(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Float32(tmpF1)), false)
		t.Assert(empty.IsEmpty(dconv.Float64(tmpF1)), false)
		t.Assert(empty.IsEmpty(true), false)
		t.Assert(empty.IsEmpty(tmpT1), false)
		t.Assert(empty.IsEmpty([]byte("1")), false)
		t.Assert(empty.IsEmpty(d.Map{"a": 1}), false)
		t.Assert(empty.IsEmpty(d.Slice{"1"}), false)
		t.Assert(empty.IsEmpty(d.Array{"1"}), false)
		t.Assert(empty.IsEmpty(tmpF2), false)
		t.Assert(empty.IsEmpty(tmpF3), false)
		t.Assert(empty.IsEmpty(tmpF4), false)
		t.Assert(empty.IsEmpty(tmpF5), false)
		t.Assert(empty.IsEmpty(tmpF6), false)
	})
}

func TestIsNil(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(empty.IsNil(nil), true)
	})
	dtest.C(t, func(t *dtest.T) {
		var i int
		t.Assert(empty.IsNil(i), false)
	})
	dtest.C(t, func(t *dtest.T) {
		var i *int
		t.Assert(empty.IsNil(i), true)
	})
	dtest.C(t, func(t *dtest.T) {
		var i *int
		t.Assert(empty.IsNil(&i), false)
		t.Assert(empty.IsNil(&i, true), true)
	})
}
