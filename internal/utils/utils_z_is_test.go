// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestVar_IsNil(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsNil(0), false)
		t.Assert(utils.IsNil(nil), true)
		t.Assert(utils.IsNil(d.Map{}), false)
		t.Assert(utils.IsNil(d.Slice{}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsNil(1), false)
		t.Assert(utils.IsNil(0.1), false)
		t.Assert(utils.IsNil(d.Map{"k": "v"}), false)
		t.Assert(utils.IsNil(d.Slice{0}), false)
	})
}

func TestVar_IsEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsEmpty(0), true)
		t.Assert(utils.IsEmpty(nil), true)
		t.Assert(utils.IsEmpty(d.Map{}), true)
		t.Assert(utils.IsEmpty(d.Slice{}), true)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsEmpty(1), false)
		t.Assert(utils.IsEmpty(0.1), false)
		t.Assert(utils.IsEmpty(d.Map{"k": "v"}), false)
		t.Assert(utils.IsEmpty(d.Slice{0}), false)
	})
}

func TestVar_IsInt(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsInt(0), true)
		t.Assert(utils.IsInt(nil), false)
		t.Assert(utils.IsInt(d.Map{}), false)
		t.Assert(utils.IsInt(d.Slice{}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsInt(1), true)
		t.Assert(utils.IsInt(-1), true)
		t.Assert(utils.IsInt(0.1), false)
		t.Assert(utils.IsInt(d.Map{"k": "v"}), false)
		t.Assert(utils.IsInt(d.Slice{0}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsInt(int8(1)), true)
		t.Assert(utils.IsInt(uint8(1)), false)
	})
}

func TestVar_IsUint(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsUint(0), false)
		t.Assert(utils.IsUint(nil), false)
		t.Assert(utils.IsUint(d.Map{}), false)
		t.Assert(utils.IsUint(d.Slice{}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsUint(1), false)
		t.Assert(utils.IsUint(-1), false)
		t.Assert(utils.IsUint(0.1), false)
		t.Assert(utils.IsUint(d.Map{"k": "v"}), false)
		t.Assert(utils.IsUint(d.Slice{0}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsUint(int8(1)), false)
		t.Assert(utils.IsUint(uint8(1)), true)
	})
}

func TestVar_IsFloat(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsFloat(0), false)
		t.Assert(utils.IsFloat(nil), false)
		t.Assert(utils.IsFloat(d.Map{}), false)
		t.Assert(utils.IsFloat(d.Slice{}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsFloat(1), false)
		t.Assert(utils.IsFloat(-1), false)
		t.Assert(utils.IsFloat(0.1), true)
		t.Assert(utils.IsFloat(float64(1)), true)
		t.Assert(utils.IsFloat(d.Map{"k": "v"}), false)
		t.Assert(utils.IsFloat(d.Slice{0}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsFloat(int8(1)), false)
		t.Assert(utils.IsFloat(uint8(1)), false)
	})
}

func TestVar_IsSlice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsSlice(0), false)
		t.Assert(utils.IsSlice(nil), false)
		t.Assert(utils.IsSlice(d.Map{}), false)
		t.Assert(utils.IsSlice(d.Slice{}), true)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsSlice(1), false)
		t.Assert(utils.IsSlice(-1), false)
		t.Assert(utils.IsSlice(0.1), false)
		t.Assert(utils.IsSlice(float64(1)), false)
		t.Assert(utils.IsSlice(d.Map{"k": "v"}), false)
		t.Assert(utils.IsSlice(d.Slice{0}), true)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsSlice(int8(1)), false)
		t.Assert(utils.IsSlice(uint8(1)), false)
	})
}

func TestVar_IsMap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsMap(0), false)
		t.Assert(utils.IsMap(nil), false)
		t.Assert(utils.IsMap(d.Map{}), true)
		t.Assert(utils.IsMap(d.Slice{}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsMap(1), false)
		t.Assert(utils.IsMap(-1), false)
		t.Assert(utils.IsMap(0.1), false)
		t.Assert(utils.IsMap(float64(1)), false)
		t.Assert(utils.IsMap(d.Map{"k": "v"}), true)
		t.Assert(utils.IsMap(d.Slice{0}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsMap(int8(1)), false)
		t.Assert(utils.IsMap(uint8(1)), false)
	})
}

func TestVar_IsStruct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsStruct(0), false)
		t.Assert(utils.IsStruct(nil), false)
		t.Assert(utils.IsStruct(d.Map{}), false)
		t.Assert(utils.IsStruct(d.Slice{}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.IsStruct(1), false)
		t.Assert(utils.IsStruct(-1), false)
		t.Assert(utils.IsStruct(0.1), false)
		t.Assert(utils.IsStruct(float64(1)), false)
		t.Assert(utils.IsStruct(d.Map{"k": "v"}), false)
		t.Assert(utils.IsStruct(d.Slice{0}), false)
	})
	dtest.C(t, func(t *dtest.T) {
		a := &struct {
		}{}
		t.Assert(utils.IsStruct(a), true)
		t.Assert(utils.IsStruct(*a), true)
		t.Assert(utils.IsStruct(&a), true)
	})
}
