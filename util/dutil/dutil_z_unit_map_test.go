// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dutil_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

func Test_MapCopy(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		m2 := dutil.MapCopy(m1)
		m2["k2"] = "v2"

		t.Assert(m1["k1"], "v1")
		t.Assert(m1["k2"], nil)
		t.Assert(m2["k1"], "v1")
		t.Assert(m2["k2"], "v2")
	})
}

func Test_MapContains(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		t.Assert(dutil.MapContains(m1, "k1"), true)
		t.Assert(dutil.MapContains(m1, "K1"), false)
		t.Assert(dutil.MapContains(m1, "k2"), false)
	})
}

func Test_MapMerge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		m2 := g.Map{
			"k2": "v2",
		}
		m3 := g.Map{
			"k3": "v3",
		}
		dutil.MapMerge(m1, m2, m3, nil)
		t.Assert(m1["k1"], "v1")
		t.Assert(m1["k2"], "v2")
		t.Assert(m1["k3"], "v3")
		t.Assert(m2["k1"], nil)
		t.Assert(m3["k1"], nil)
	})
}

func Test_MapMergeCopy(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		m2 := g.Map{
			"k2": "v2",
		}
		m3 := g.Map{
			"k3": "v3",
		}
		m := dutil.MapMergeCopy(m1, m2, m3, nil)
		t.Assert(m["k1"], "v1")
		t.Assert(m["k2"], "v2")
		t.Assert(m["k3"], "v3")
		t.Assert(m1["k1"], "v1")
		t.Assert(m1["k2"], nil)
		t.Assert(m2["k1"], nil)
		t.Assert(m3["k1"], nil)
	})
}

func Test_MapPossibleItemByKey(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := g.Map{
			"name":     "guo",
			"NickName": "john",
		}
		k, v := dutil.MapPossibleItemByKey(m, "NAME")
		t.Assert(k, "name")
		t.Assert(v, "guo")

		k, v = dutil.MapPossibleItemByKey(m, "nick name")
		t.Assert(k, "NickName")
		t.Assert(v, "john")

		k, v = dutil.MapPossibleItemByKey(m, "none")
		t.Assert(k, "")
		t.Assert(v, nil)
	})
}

func Test_MapContainsPossibleKey(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := g.Map{
			"name":     "guo",
			"NickName": "john",
		}
		t.Assert(dutil.MapContainsPossibleKey(m, "name"), true)
		t.Assert(dutil.MapContainsPossibleKey(m, "NAME"), true)
		t.Assert(dutil.MapContainsPossibleKey(m, "nickname"), true)
		t.Assert(dutil.MapContainsPossibleKey(m, "nick name"), true)
		t.Assert(dutil.MapContainsPossibleKey(m, "nick_name"), true)
		t.Assert(dutil.MapContainsPossibleKey(m, "nick-name"), true)
		t.Assert(dutil.MapContainsPossibleKey(m, "nick.name"), true)
		t.Assert(dutil.MapContainsPossibleKey(m, "none"), false)
	})
}

func Test_MapOmitEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := g.Map{
			"k1": "john",
			"e1": "",
			"e2": 0,
			"e3": nil,
			"k2": "smith",
		}
		dutil.MapOmitEmpty(m)
		t.Assert(len(m), 2)
		t.AssertNE(m["k1"], nil)
		t.AssertNE(m["k2"], nil)
	})
}

func Test_MapToSlice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		s := dutil.MapToSlice(m)
		t.Assert(len(s), 4)
		t.AssertIN(s[0], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[1], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[2], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[3], g.Slice{"k1", "k2", "v1", "v2"})
	})
	dtest.C(t, func(t *dtest.T) {
		m := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		s := dutil.MapToSlice(m)
		t.Assert(len(s), 4)
		t.AssertIN(s[0], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[1], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[2], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[3], g.Slice{"k1", "k2", "v1", "v2"})
	})
	dtest.C(t, func(t *dtest.T) {
		m := g.MapStrStr{}
		s := dutil.MapToSlice(m)
		t.Assert(len(s), 0)
	})
	dtest.C(t, func(t *dtest.T) {
		s := dutil.MapToSlice(1)
		t.Assert(s, nil)
	})
}
