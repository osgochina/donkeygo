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

func Test_ListItemValues_Map(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 99},
			d.Map{"id": 3, "score": 99},
		}
		t.Assert(dutil.ListItemValues(listMap, "id"), d.Slice{1, 2, 3})
		t.Assert(dutil.ListItemValues(listMap, "score"), d.Slice{100, 99, 99})
	})
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": nil},
			d.Map{"id": 3, "score": 0},
		}
		t.Assert(dutil.ListItemValues(listMap, "id"), d.Slice{1, 2, 3})
		t.Assert(dutil.ListItemValues(listMap, "score"), d.Slice{100, nil, 0})
	})
}

func Test_ListItemValues_Map_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "scores": Scores{100, 60}},
			d.Map{"id": 2, "scores": Scores{0, 100}},
			d.Map{"id": 3, "scores": Scores{59, 99}},
		}
		t.Assert(dutil.ListItemValues(listMap, "scores", "Math"), d.Slice{100, 0, 59})
		t.Assert(dutil.ListItemValues(listMap, "scores", "English"), d.Slice{60, 100, 99})
		t.Assert(dutil.ListItemValues(listMap, "scores", "PE"), d.Slice{})
	})
}

func Test_ListItemValues_Map_Array_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "scores": []Scores{{1, 2}, {3, 4}}},
			d.Map{"id": 2, "scores": []Scores{{5, 6}, {7, 8}}},
			d.Map{"id": 3, "scores": []Scores{{9, 10}, {11, 12}}},
		}
		t.Assert(dutil.ListItemValues(listMap, "scores", "Math"), d.Slice{1, 3, 5, 7, 9, 11})
		t.Assert(dutil.ListItemValues(listMap, "scores", "English"), d.Slice{2, 4, 6, 8, 10, 12})
		t.Assert(dutil.ListItemValues(listMap, "scores", "PE"), d.Slice{})
	})
}

func Test_ListItemValues_Struct(t *testing.T) {
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
		t.Assert(dutil.ListItemValues(listStruct, "Id"), d.Slice{1, 2, 3})
		t.Assert(dutil.ListItemValues(listStruct, "Score"), d.Slice{100, 99, 0})
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
		t.Assert(dutil.ListItemValues(listStruct, "Id"), d.Slice{1, 2, 3})
		t.Assert(dutil.ListItemValues(listStruct, "Score"), d.Slice{100, 99, 0})
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
		t.Assert(dutil.ListItemValues(listStruct, "Id"), d.Slice{1, 2, 3})
		t.Assert(dutil.ListItemValues(listStruct, "Score"), d.Slice{100, nil, 0})
	})
}

func Test_ListItemValues_Struct_SubKey(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []Student
		}
		listStruct := d.Slice{
			Class{2, []Student{{1, 1}, {2, 2}}},
			Class{3, []Student{{3, 3}, {4, 4}, {5, 5}}},
			Class{1, []Student{{6, 6}}},
		}
		t.Assert(dutil.ListItemValues(listStruct, "Total"), d.Slice{2, 3, 1})
		t.Assert(dutil.ListItemValues(listStruct, "Students"), `[[{"Id":1,"Score":1},{"Id":2,"Score":2}],[{"Id":3,"Score":3},{"Id":4,"Score":4},{"Id":5,"Score":5}],[{"Id":6,"Score":6}]]`)
		t.Assert(dutil.ListItemValues(listStruct, "Students", "Id"), d.Slice{1, 2, 3, 4, 5, 6})
	})
	dtest.C(t, func(t *dtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []*Student
		}
		listStruct := d.Slice{
			&Class{2, []*Student{{1, 1}, {2, 2}}},
			&Class{3, []*Student{{3, 3}, {4, 4}, {5, 5}}},
			&Class{1, []*Student{{6, 6}}},
		}
		t.Assert(dutil.ListItemValues(listStruct, "Total"), d.Slice{2, 3, 1})
		t.Assert(dutil.ListItemValues(listStruct, "Students"), `[[{"Id":1,"Score":1},{"Id":2,"Score":2}],[{"Id":3,"Score":3},{"Id":4,"Score":4},{"Id":5,"Score":5}],[{"Id":6,"Score":6}]]`)
		t.Assert(dutil.ListItemValues(listStruct, "Students", "Id"), d.Slice{1, 2, 3, 4, 5, 6})
	})
}

func Test_ListItemValuesUnique(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 100},
			d.Map{"id": 3, "score": 100},
			d.Map{"id": 4, "score": 100},
			d.Map{"id": 5, "score": 100},
		}
		t.Assert(dutil.ListItemValuesUnique(listMap, "id"), d.Slice{1, 2, 3, 4, 5})
		t.Assert(dutil.ListItemValuesUnique(listMap, "score"), d.Slice{100})
	})
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 100},
			d.Map{"id": 3, "score": 100},
			d.Map{"id": 4, "score": 100},
			d.Map{"id": 5, "score": 99},
		}
		t.Assert(dutil.ListItemValuesUnique(listMap, "id"), d.Slice{1, 2, 3, 4, 5})
		t.Assert(dutil.ListItemValuesUnique(listMap, "score"), d.Slice{100, 99})
	})
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "score": 100},
			d.Map{"id": 2, "score": 100},
			d.Map{"id": 3, "score": 0},
			d.Map{"id": 4, "score": 100},
			d.Map{"id": 5, "score": 99},
		}
		t.Assert(dutil.ListItemValuesUnique(listMap, "id"), d.Slice{1, 2, 3, 4, 5})
		t.Assert(dutil.ListItemValuesUnique(listMap, "score"), d.Slice{100, 0, 99})
	})
}

func Test_ListItemValuesUnique_Struct_SubKey(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []Student
		}
		listStruct := d.Slice{
			Class{2, []Student{{1, 1}, {1, 2}}},
			Class{3, []Student{{2, 3}, {2, 4}, {5, 5}}},
			Class{1, []Student{{6, 6}}},
		}
		t.Assert(dutil.ListItemValuesUnique(listStruct, "Total"), d.Slice{2, 3, 1})
		t.Assert(dutil.ListItemValuesUnique(listStruct, "Students", "Id"), d.Slice{1, 2, 5, 6})
	})
	dtest.C(t, func(t *dtest.T) {
		type Student struct {
			Id    int
			Score float64
		}
		type Class struct {
			Total    int
			Students []*Student
		}
		listStruct := d.Slice{
			&Class{2, []*Student{{1, 1}, {1, 2}}},
			&Class{3, []*Student{{2, 3}, {2, 4}, {5, 5}}},
			&Class{1, []*Student{{6, 6}}},
		}
		t.Assert(dutil.ListItemValuesUnique(listStruct, "Total"), d.Slice{2, 3, 1})
		t.Assert(dutil.ListItemValuesUnique(listStruct, "Students", "Id"), d.Slice{1, 2, 5, 6})
	})
}

func Test_ListItemValuesUnique_Map_Array_SubKey(t *testing.T) {
	type Scores struct {
		Math    int
		English int
	}
	dtest.C(t, func(t *dtest.T) {
		listMap := d.List{
			d.Map{"id": 1, "scores": []Scores{{1, 2}, {1, 2}}},
			d.Map{"id": 2, "scores": []Scores{{5, 8}, {5, 8}}},
			d.Map{"id": 3, "scores": []Scores{{9, 10}, {11, 12}}},
		}
		t.Assert(dutil.ListItemValuesUnique(listMap, "scores", "Math"), d.Slice{1, 5, 9, 11})
		t.Assert(dutil.ListItemValuesUnique(listMap, "scores", "English"), d.Slice{2, 8, 10, 12})
		t.Assert(dutil.ListItemValuesUnique(listMap, "scores", "PE"), d.Slice{})
	})
}
