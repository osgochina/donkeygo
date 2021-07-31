// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar_test

import (
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestVar_Ints(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Ints()[0], arr[0])
	})
}

func TestVar_Uints(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Uints()[0], arr[0])
	})
}

func TestVar_Int64s(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Int64s()[0], arr[0])
	})
}

func TestVar_Uint64s(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Uint64s()[0], arr[0])
	})
}

func TestVar_Floats(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []float64{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Floats()[0], arr[0])
	})
}

func TestVar_Float32s(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []float32{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.AssertEQ(objOne.Float32s(), arr)
	})
}

func TestVar_Float64s(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []float64{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.AssertEQ(objOne.Float64s(), arr)
	})
}

func TestVar_Strings(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []string{"hello", "world"}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Strings()[0], arr[0])
	})
}

func TestVar_Interfaces(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Interfaces(), arr)
	})
}

func TestVar_Slice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, true)
		t.Assert(objOne.Slice(), arr)
	})
}

func TestVar_Array(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, false)
		t.Assert(objOne.Array(), arr)
	})
}

func TestVar_Vars(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := dvar.New(arr, false)
		t.Assert(len(objOne.Vars()), 5)
		t.Assert(objOne.Vars()[0].Int(), 1)
		t.Assert(objOne.Vars()[4].Int(), 5)
	})
}
