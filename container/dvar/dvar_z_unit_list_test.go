// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar_test

import (
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestVar_ListItemValues_Map(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 99},
			d.Map{"id": 3, "score": 99},
		}
		t.Assert(dvar.New(listMap).ListItemValues("id"), d.Slice{1, 2, 3})
		t.Assert(dvar.New(listMap).ListItemValues("score"), d.Slice{100, 99, 99})
	})
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": nil},
			d.Map{"id": 3, "score": 0},
		}
		t.Assert(dvar.New(listMap).ListItemValues("id"), d.Slice{1, 2, 3})
		t.Assert(dvar.New(listMap).ListItemValues("score"), d.Slice{100, nil, 0})
	})
}

func TestVar_ListItemValues_Struct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Id    int
			Score float64
		}
		listStruct := d.Slice{
			T{1, 100},
			T{2, 99},
			T{3, 0},
		}
		t.Assert(dvar.New(listStruct).ListItemValues("Id"), d.Slice{1, 2, 3})
		t.Assert(dvar.New(listStruct).ListItemValues("Score"), d.Slice{100, 99, 0})
	})
	// Pointer items.
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Id    int
			Score float64
		}
		listStruct := d.Slice{
			&T{1, 100},
			&T{2, 99},
			&T{3, 0},
		}
		t.Assert(dvar.New(listStruct).ListItemValues("Id"), d.Slice{1, 2, 3})
		t.Assert(dvar.New(listStruct).ListItemValues("Score"), d.Slice{100, 99, 0})
	})
	// Nil element value.
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Id    int
			Score interface{}
		}
		listStruct := d.Slice{
			T{1, 100},
			T{2, nil},
			T{3, 0},
		}
		t.Assert(dvar.New(listStruct).ListItemValues("Id"), d.Slice{1, 2, 3})
		t.Assert(dvar.New(listStruct).ListItemValues("Score"), d.Slice{100, nil, 0})
	})
}

func TestVar_ListItemValuesUnique(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 100},
			d.Map{"id": 3, "score": 100},
			d.Map{"id": 4, "score": 100},
			d.Map{"id": 5, "score": 100},
		}
		t.Assert(dvar.New(listMap).ListItemValuesUnique("id"), d.Slice{1, 2, 3, 4, 5})
		t.Assert(dvar.New(listMap).ListItemValuesUnique("score"), d.Slice{100})
	})
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 100},
			d.Map{"id": 3, "score": 100},
			d.Map{"id": 4, "score": 100},
			d.Map{"id": 5, "score": 99},
		}
		t.Assert(dvar.New(listMap).ListItemValuesUnique("id"), d.Slice{1, 2, 3, 4, 5})
		t.Assert(dvar.New(listMap).ListItemValuesUnique("score"), d.Slice{100, 99})
	})
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 100},
			d.Map{"id": 3, "score": 0},
			d.Map{"id": 4, "score": 100},
			d.Map{"id": 5, "score": 99},
		}
		t.Assert(dvar.New(listMap).ListItemValuesUnique("id"), d.Slice{1, 2, 3, 4, 5})
		t.Assert(dvar.New(listMap).ListItemValuesUnique("score"), d.Slice{100, 0, 99})
	})
}
