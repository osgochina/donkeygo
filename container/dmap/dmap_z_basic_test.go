// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package dmap_test

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

func getValue() interface{} {
	return 3
}

func Test_Map_Var(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var m dmap.Map
		m.Set(1, 11)
		t.Assert(m.Get(1), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.IntAnyMap
		m.Set(1, 11)
		t.Assert(m.Get(1), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.IntIntMap
		m.Set(1, 11)
		t.Assert(m.Get(1), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.IntStrMap
		m.Set(1, "11")
		t.Assert(m.Get(1), "11")
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.StrAnyMap
		m.Set("1", "11")
		t.Assert(m.Get("1"), "11")
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.StrStrMap
		m.Set("1", "11")
		t.Assert(m.Get("1"), "11")
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.StrIntMap
		m.Set("1", 11)
		t.Assert(m.Get("1"), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.ListMap
		m.Set("1", 11)
		t.Assert(m.Get("1"), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		var m dmap.TreeMap
		m.SetComparator(dutil.ComparatorString)
		m.Set("1", 11)
		t.Assert(m.Get("1"), 11)
	})
}

func Test_Map_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := dmap.New()
		m.Set("key1", "val1")
		t.Assert(m.Keys(), []interface{}{"key1"})

		t.Assert(m.Get("key1"), "val1")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet("key2", "val2"), "val2")
		t.Assert(m.SetIfNotExist("key2", "val2"), false)

		t.Assert(m.SetIfNotExist("key3", "val3"), true)

		t.Assert(m.Remove("key2"), "val2")
		t.Assert(m.Contains("key2"), false)

		t.AssertIN("key3", m.Keys())
		t.AssertIN("key1", m.Keys())
		t.AssertIN("val3", m.Values())
		t.AssertIN("val1", m.Values())

		m.Flip()
		t.Assert(m.Map(), map[interface{}]interface{}{"val3": "key3", "val1": "key1"})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := dmap.NewFrom(map[interface{}]interface{}{1: 1, "key1": "val1"})
		t.Assert(m2.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
	})
}
func Test_Map_Set_Fun(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := dmap.New()
		m.GetOrSetFunc("fun", getValue)
		m.GetOrSetFuncLock("funlock", getValue)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
		m.GetOrSetFunc("fun", getValue)
		t.Assert(m.SetIfNotExistFunc("fun", getValue), false)
		t.Assert(m.SetIfNotExistFuncLock("funlock", getValue), false)
	})
}

func Test_Map_Batch(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := dmap.New()
		m.Sets(map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		t.Assert(m.Map(), map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		m.Removes([]interface{}{"key1", 1})
		t.Assert(m.Map(), map[interface{}]interface{}{"key2": "val2", "key3": "val3"})
	})
}
func Test_Map_Iterator(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := map[interface{}]interface{}{1: 1, "key1": "val1"}

		m := dmap.NewFrom(expect)
		m.Iterator(func(k interface{}, v interface{}) bool {
			t.Assert(expect[k], v)
			return true
		})
		// 断言返回值对遍历控制
		i := 0
		j := 0
		m.Iterator(func(k interface{}, v interface{}) bool {
			i++
			return true
		})
		m.Iterator(func(k interface{}, v interface{}) bool {
			j++
			return false
		})
		t.Assert(i, 2)
		t.Assert(j, 1)
	})
}

func Test_Map_Lock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := map[interface{}]interface{}{1: 1, "key1": "val1"}
		m := dmap.NewFrom(expect)
		m.LockFunc(func(m map[interface{}]interface{}) {
			t.Assert(m, expect)
		})
		m.RLockFunc(func(m map[interface{}]interface{}) {
			t.Assert(m, expect)
		})
	})
}

func Test_Map_Clone(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		//clone 方法是深克隆
		m := dmap.NewFrom(map[interface{}]interface{}{1: 1, "key1": "val1"})
		m_clone := m.Clone()
		m.Remove(1)
		//修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove("key1")
		//修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}
func Test_Map_Basic_Merge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m1 := dmap.New()
		m2 := dmap.New()
		m1.Set("key1", "val1")
		m2.Set("key2", "val2")
		m1.Merge(m2)
		t.Assert(m1.Map(), map[interface{}]interface{}{"key1": "val1", "key2": "val2"})
	})
}
