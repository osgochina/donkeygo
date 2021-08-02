// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dutil_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

func Test_Dump(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		dutil.Dump(map[int]int{
			100: 100,
		})
	})
}

func Test_Try(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := `dutil Try test`
		t.Assert(dutil.Try(func() {
			panic(s)
		}), s)
	})
}

func Test_TryCatch(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		dutil.TryCatch(func() {
			panic("dutil TryCatch test")
		})
	})

	dtest.C(t, func(t *dtest.T) {
		dutil.TryCatch(func() {
			panic("dutil TryCatch test")

		}, func(err error) {
			t.Assert(err, "dutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		defer func() {
			t.Assert(recover(), "dutil Throw test")
		}()

		dutil.Throw("dutil Throw test")
	})
}

func Test_Keys(t *testing.T) {
	// map
	dtest.C(t, func(t *dtest.T) {
		keys := dutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *map
	dtest.C(t, func(t *dtest.T) {
		keys := dutil.Keys(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *struct
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			A string
			B int
		}
		keys := dutil.Keys(new(T))
		t.Assert(keys, d.SliceStr{"A", "B"})
	})
	// *struct nil
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := dutil.Keys(pointer)
		t.Assert(keys, d.SliceStr{"A", "B"})
	})
	// **struct nil
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := dutil.Keys(&pointer)
		t.Assert(keys, d.SliceStr{"A", "B"})
	})
}

func Test_Values(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		values := dutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", values)
		t.AssertIN("2", values)
	})

	dtest.C(t, func(t *dtest.T) {
		type T struct {
			A string
			B int
		}
		keys := dutil.Values(T{
			A: "1",
			B: 2,
		})
		t.Assert(keys, d.Slice{"1", 2})
	})
}
