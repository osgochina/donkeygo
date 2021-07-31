// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package dmap_test

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

func Test_TreeMap_Var(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var m dmap.TreeMap
		m.SetComparator(dutil.ComparatorString)
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
	})
}

func Test_TreeMap_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := dmap.NewTreeMap(dutil.ComparatorString)
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

		m2 := dmap.NewTreeMapFrom(dutil.ComparatorString, map[interface{}]interface{}{1: 1, "key1": "val1"})
		t.Assert(m2.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
	})
}

func Test_TreeMap_Set_Fun(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := dmap.NewTreeMap(dutil.ComparatorString)
		m.GetOrSetFunc("fun", getValue)
		m.GetOrSetFuncLock("funlock", getValue)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
		m.GetOrSetFunc("fun", getValue)
		t.Assert(m.SetIfNotExistFunc("fun", getValue), false)
		t.Assert(m.SetIfNotExistFuncLock("funlock", getValue), false)
	})
}

func Test_TreeMap_Batch(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m := dmap.NewTreeMap(dutil.ComparatorString)
		m.Sets(map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		t.Assert(m.Map(), map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		m.Removes([]interface{}{"key1", 1})
		t.Assert(m.Map(), map[interface{}]interface{}{"key2": "val2", "key3": "val3"})
	})
}
func Test_TreeMap_Iterator(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := map[interface{}]interface{}{1: 1, "key1": "val1"}
		m := dmap.NewTreeMapFrom(dutil.ComparatorString, expect)
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

	dtest.C(t, func(t *dtest.T) {
		expect := map[interface{}]interface{}{1: 1, "key1": "val1"}
		m := dmap.NewTreeMapFrom(dutil.ComparatorString, expect)
		for i := 0; i < 10; i++ {
			m.IteratorAsc(func(k interface{}, v interface{}) bool {
				t.Assert(expect[k], v)
				return true
			})
		}
		j := 0
		for i := 0; i < 10; i++ {
			m.IteratorAsc(func(k interface{}, v interface{}) bool {
				j++
				return false
			})
		}
		t.Assert(j, 10)
	})
}

func Test_TreeMap_Clone(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		//clone 方法是深克隆
		m := dmap.NewTreeMapFrom(dutil.ComparatorString, map[interface{}]interface{}{1: 1, "key1": "val1"})
		m_clone := m.Clone()
		m.Remove(1)
		//修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove("key1")
		//修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}

func Test_TreeMap_Json(t *testing.T) {
	// Marshal
	dtest.C(t, func(t *dtest.T) {
		data := d.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		m1 := dmap.NewTreeMapFrom(dutil.ComparatorString, data)
		b1, err1 := json.Marshal(m1)
		b2, err2 := json.Marshal(dconv.Map(data))
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})
	// Unmarshal
	dtest.C(t, func(t *dtest.T) {
		data := d.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(dconv.Map(data))
		t.Assert(err, nil)

		m := dmap.NewTreeMap(dutil.ComparatorString)
		err = json.UnmarshalUseNumber(b, m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
	dtest.C(t, func(t *dtest.T) {
		data := d.MapAnyAny{
			"k1": "v1",
			"k2": "v2",
		}
		b, err := json.Marshal(dconv.Map(data))
		t.Assert(err, nil)

		var m dmap.TreeMap
		err = json.UnmarshalUseNumber(b, &m)
		t.Assert(err, nil)
		t.Assert(m.Get("k1"), data["k1"])
		t.Assert(m.Get("k2"), data["k2"])
	})
}

func TestTreeMap_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Map  *dmap.TreeMap
	}
	// JSON
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"map":  []byte(`{"k1":"v1","k2":"v2"}`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), "v1")
		t.Assert(v.Map.Get("k2"), "v2")
	})
	// Map
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"map": d.Map{
				"k1": "v1",
				"k2": "v2",
			},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Map.Size(), 2)
		t.Assert(v.Map.Get("k1"), "v1")
		t.Assert(v.Map.Get("k2"), "v2")
	})
}
