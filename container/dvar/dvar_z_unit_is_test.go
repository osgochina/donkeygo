// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestVar_IsNil(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsNil(), false)
		t.Assert(d.NewVar(nil).IsNil(), true)
		t.Assert(d.NewVar(d.Map{}).IsNil(), false)
		t.Assert(d.NewVar(d.Slice{}).IsNil(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsNil(), false)
		t.Assert(d.NewVar(0.1).IsNil(), false)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsNil(), false)
		t.Assert(d.NewVar(d.Slice{0}).IsNil(), false)
	})
}

func TestVar_IsEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsEmpty(), true)
		t.Assert(d.NewVar(nil).IsEmpty(), true)
		t.Assert(d.NewVar(d.Map{}).IsEmpty(), true)
		t.Assert(d.NewVar(d.Slice{}).IsEmpty(), true)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsEmpty(), false)
		t.Assert(d.NewVar(0.1).IsEmpty(), false)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsEmpty(), false)
		t.Assert(d.NewVar(d.Slice{0}).IsEmpty(), false)
	})
}

func TestVar_IsInt(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsInt(), true)
		t.Assert(d.NewVar(nil).IsInt(), false)
		t.Assert(d.NewVar(d.Map{}).IsInt(), false)
		t.Assert(d.NewVar(d.Slice{}).IsInt(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsInt(), true)
		t.Assert(d.NewVar(-1).IsInt(), true)
		t.Assert(d.NewVar(0.1).IsInt(), false)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsInt(), false)
		t.Assert(d.NewVar(d.Slice{0}).IsInt(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(int8(1)).IsInt(), true)
		t.Assert(d.NewVar(uint8(1)).IsInt(), false)
	})
}

func TestVar_IsUint(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsUint(), false)
		t.Assert(d.NewVar(nil).IsUint(), false)
		t.Assert(d.NewVar(d.Map{}).IsUint(), false)
		t.Assert(d.NewVar(d.Slice{}).IsUint(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsUint(), false)
		t.Assert(d.NewVar(-1).IsUint(), false)
		t.Assert(d.NewVar(0.1).IsUint(), false)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsUint(), false)
		t.Assert(d.NewVar(d.Slice{0}).IsUint(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(int8(1)).IsUint(), false)
		t.Assert(d.NewVar(uint8(1)).IsUint(), true)
	})
}

func TestVar_IsFloat(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsFloat(), false)
		t.Assert(d.NewVar(nil).IsFloat(), false)
		t.Assert(d.NewVar(d.Map{}).IsFloat(), false)
		t.Assert(d.NewVar(d.Slice{}).IsFloat(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsFloat(), false)
		t.Assert(d.NewVar(-1).IsFloat(), false)
		t.Assert(d.NewVar(0.1).IsFloat(), true)
		t.Assert(d.NewVar(float64(1)).IsFloat(), true)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsFloat(), false)
		t.Assert(d.NewVar(d.Slice{0}).IsFloat(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(int8(1)).IsFloat(), false)
		t.Assert(d.NewVar(uint8(1)).IsFloat(), false)
	})
}

func TestVar_IsSlice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsSlice(), false)
		t.Assert(d.NewVar(nil).IsSlice(), false)
		t.Assert(d.NewVar(d.Map{}).IsSlice(), false)
		t.Assert(d.NewVar(d.Slice{}).IsSlice(), true)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsSlice(), false)
		t.Assert(d.NewVar(-1).IsSlice(), false)
		t.Assert(d.NewVar(0.1).IsSlice(), false)
		t.Assert(d.NewVar(float64(1)).IsSlice(), false)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsSlice(), false)
		t.Assert(d.NewVar(d.Slice{0}).IsSlice(), true)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(int8(1)).IsSlice(), false)
		t.Assert(d.NewVar(uint8(1)).IsSlice(), false)
	})
}

func TestVar_IsMap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsMap(), false)
		t.Assert(d.NewVar(nil).IsMap(), false)
		t.Assert(d.NewVar(d.Map{}).IsMap(), true)
		t.Assert(d.NewVar(d.Slice{}).IsMap(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsMap(), false)
		t.Assert(d.NewVar(-1).IsMap(), false)
		t.Assert(d.NewVar(0.1).IsMap(), false)
		t.Assert(d.NewVar(float64(1)).IsMap(), false)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsMap(), true)
		t.Assert(d.NewVar(d.Slice{0}).IsMap(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(int8(1)).IsMap(), false)
		t.Assert(d.NewVar(uint8(1)).IsMap(), false)
	})
}

func TestVar_IsStruct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(0).IsStruct(), false)
		t.Assert(d.NewVar(nil).IsStruct(), false)
		t.Assert(d.NewVar(d.Map{}).IsStruct(), false)
		t.Assert(d.NewVar(d.Slice{}).IsStruct(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(d.NewVar(1).IsStruct(), false)
		t.Assert(d.NewVar(-1).IsStruct(), false)
		t.Assert(d.NewVar(0.1).IsStruct(), false)
		t.Assert(d.NewVar(float64(1)).IsStruct(), false)
		t.Assert(d.NewVar(d.Map{"k": "v"}).IsStruct(), false)
		t.Assert(d.NewVar(d.Slice{0}).IsStruct(), false)
	})
	dtest.C(t, func(t *dtest.T) {
		a := &struct {
		}{}
		t.Assert(d.NewVar(a).IsStruct(), true)
		t.Assert(d.NewVar(*a).IsStruct(), true)
		t.Assert(d.NewVar(&a).IsStruct(), true)
	})
}
