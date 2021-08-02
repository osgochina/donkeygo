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

func Test_SliceToMap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := d.Slice{
			"K1", "v1", "K2", "v2",
		}
		m := dutil.SliceToMap(s)
		t.Assert(len(m), 2)
		t.Assert(m, d.Map{
			"K1": "v1",
			"K2": "v2",
		})
	})
	dtest.C(t, func(t *dtest.T) {
		s := d.Slice{
			"K1", "v1", "K2",
		}
		m := dutil.SliceToMap(s)
		t.Assert(len(m), 0)
		t.Assert(m, nil)
	})
}

func Test_SliceToMapWithColumnAsKey(t *testing.T) {
	m1 := d.Map{"K1": "v1", "K2": 1}
	m2 := d.Map{"K1": "v2", "K2": 2}
	s := d.Slice{m1, m2}
	dtest.C(t, func(t *dtest.T) {
		m := dutil.SliceToMapWithColumnAsKey(s, "K1")
		t.Assert(m, d.MapAnyAny{
			"v1": m1,
			"v2": m2,
		})
	})
	dtest.C(t, func(t *dtest.T) {
		m := dutil.SliceToMapWithColumnAsKey(s, "K2")
		t.Assert(m, d.MapAnyAny{
			1: m1,
			2: m2,
		})
	})
}
