// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

func Test_MapToMap1(t *testing.T) {
	// map[int]int -> map[string]string
	// empty original map.
	dtest.C(t, func(t *dtest.T) {
		m1 := d.MapIntInt{}
		m2 := d.MapStrStr{}
		t.Assert(dconv.MapToMap(m1, &m2), nil)
		t.Assert(len(m1), len(m2))
	})
	// map[int]int -> map[string]string
	dtest.C(t, func(t *dtest.T) {
		m1 := d.MapIntInt{
			1: 100,
			2: 200,
		}
		m2 := d.MapStrStr{}
		t.Assert(dconv.MapToMap(m1, &m2), nil)
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
		t.Assert(dconv.MapToMap(m1, &m2), nil)
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
		t.Assert(dconv.MapToMap(m1, &m2), nil)
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
		t.Assert(dconv.MapToMap(m1, &m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
}

func Test_MapToMap2(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	params := d.Map{
		"key": d.Map{
			"id":   1,
			"name": "john",
		},
	}
	dtest.C(t, func(t *dtest.T) {
		m := make(map[string]User)
		err := dconv.MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	dtest.C(t, func(t *dtest.T) {
		m := (map[string]User)(nil)
		err := dconv.MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	dtest.C(t, func(t *dtest.T) {
		m := make(map[string]*User)
		err := dconv.MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	dtest.C(t, func(t *dtest.T) {
		m := (map[string]*User)(nil)
		err := dconv.MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
}

func Test_MapToMapDeep(t *testing.T) {
	type Ids struct {
		Id  int
		Uid int
	}
	type Base struct {
		Ids
		Time string
	}
	type User struct {
		Base
		Name string
	}
	params := d.Map{
		"key": d.Map{
			"id":   1,
			"name": "john",
		},
	}
	dtest.C(t, func(t *dtest.T) {
		m := (map[string]*User)(nil)
		err := dconv.MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
}

func Test_MapToMaps(t *testing.T) {
	params := d.Slice{
		d.Map{"id": 1, "name": "john"},
		d.Map{"id": 2, "name": "smith"},
	}
	dtest.C(t, func(t *dtest.T) {
		var s []d.Map
		err := dconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s, params)
	})
	dtest.C(t, func(t *dtest.T) {
		var s []*d.Map
		err := dconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s, params)
	})
}

func Test_MapToMaps_StructParams(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	params := d.Slice{
		User{1, "name1"},
		User{2, "name2"},
	}
	dtest.C(t, func(t *dtest.T) {
		var s []d.Map
		err := dconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
	})
	dtest.C(t, func(t *dtest.T) {
		var s []*d.Map
		err := dconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
	})
}
