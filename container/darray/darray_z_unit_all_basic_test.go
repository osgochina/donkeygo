// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package darray_test

import (
	"github.com/gogf/gf/util/gutil"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"strings"
	"testing"
)

func Test_Array_Var(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var array darray.Array
		expect := []int{2, 3, 1}
		array.Append(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	dtest.C(t, func(t *dtest.T) {
		var array darray.IntArray
		expect := []int{2, 3, 1}
		array.Append(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	dtest.C(t, func(t *dtest.T) {
		var array darray.StrArray
		expect := []string{"b", "a"}
		array.Append("b", "a")
		t.Assert(array.Slice(), expect)
	})
	dtest.C(t, func(t *dtest.T) {
		var array darray.SortedArray
		array.SetComparator(gutil.ComparatorInt)
		expect := []int{1, 2, 3}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	dtest.C(t, func(t *dtest.T) {
		var array darray.SortedIntArray
		expect := []int{1, 2, 3}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	dtest.C(t, func(t *dtest.T) {
		var array darray.SortedStrArray
		expect := []string{"a", "b", "c"}
		array.Add("c", "a", "b")
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedIntArray_Var(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var array darray.SortedIntArray
		expect := []int{1, 2, 3}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
}

func Test_IntArray_Unique(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []int{1, 2, 3, 4, 5, 6}
		array := darray.NewIntArray()
		array.Append(1, 1, 2, 3, 3, 4, 4, 5, 5, 6, 6)
		array.Unique()
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedIntArray1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		array := darray.NewSortedIntArray()
		for i := 10; i > -1; i-- {
			array.Add(i)
		}
		t.Assert(array.Slice(), expect)
		t.Assert(array.Add().Slice(), expect)
	})
}

func Test_SortedIntArray2(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		array := darray.NewSortedIntArray()
		for i := 0; i <= 10; i++ {
			array.Add(i)
		}
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedStrArray1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		array1 := darray.NewSortedStrArray()
		array2 := darray.NewSortedStrArray(true)
		for i := 10; i > -1; i-- {
			array1.Add(dconv.String(i))
			array2.Add(dconv.String(i))
		}
		t.Assert(array1.Slice(), expect)
		t.Assert(array2.Slice(), expect)
	})

}

func Test_SortedStrArray2(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		array := darray.NewSortedStrArray()
		for i := 0; i <= 10; i++ {
			array.Add(dconv.String(i))
		}
		t.Assert(array.Slice(), expect)
		array.Add()
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedArray1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		array := darray.NewSortedArray(func(v1, v2 interface{}) int {
			return strings.Compare(dconv.String(v1), dconv.String(v2))
		})
		for i := 10; i > -1; i-- {
			array.Add(dconv.String(i))
		}
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedArray2(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(dconv.String(v1), dconv.String(v2))
		}
		array := darray.NewSortedArray(func1)
		array2 := darray.NewSortedArray(func1, true)
		for i := 0; i <= 10; i++ {
			array.Add(dconv.String(i))
			array2.Add(dconv.String(i))
		}
		t.Assert(array.Slice(), expect)
		t.Assert(array.Add().Slice(), expect)
		t.Assert(array2.Slice(), expect)
	})
}

func TestNewFromCopy(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{"100", "200", "300", "400", "500", "600"}
		array1 := darray.NewFromCopy(a1)
		t.AssertIN(array1.PopRands(2), a1)
		t.Assert(len(array1.PopRands(1)), 1)
		t.Assert(len(array1.PopRands(9)), 3)
	})
}
