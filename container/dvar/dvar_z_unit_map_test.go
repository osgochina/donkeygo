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

func TestVar_Map(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := d.Map{
			"k1": "v1",
			"k2": "v2",
		}
		objOne := dvar.New(m, true)
		t.Assert(objOne.Map()["k1"], m["k1"])
		t.Assert(objOne.Map()["k2"], m["k2"])
	})
}

func TestVar_MapToMap(t *testing.T) {
	// map[int]int -> map[string]string
	// empty original map.
	dtest.C(t, func(t *dtest.T) {
		m1 := d.MapIntInt{}
		m2 := d.MapStrStr{}
		t.Assert(dvar.New(m1).MapToMap(&m2), nil)
		t.Assert(len(m1), len(m2))
	})
	// map[int]int -> map[string]string
	dtest.C(t, func(t *dtest.T) {
		m1 := d.MapIntInt{
			1: 100,
			2: 200,
		}
		m2 := d.MapStrStr{}
		t.Assert(dvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["1"], m1[1])
		t.Assert(m2["2"], m1[2])
	})
	// map[string]interface{} -> map[string]string
	dtest.C(t, func(t *dtest.T) {
		m1 := d.Map{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := d.MapStrStr{}
		t.Assert(dvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]string -> map[string]interface{}
	dtest.C(t, func(t *dtest.T) {
		m1 := d.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := d.Map{}
		t.Assert(dvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]interface{} -> map[interface{}]interface{}
	dtest.C(t, func(t *dtest.T) {
		m1 := d.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := d.MapAnyAny{}
		t.Assert(dvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
}
