// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package darray_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/util/dconv"

	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/test/dtest"
	"strings"
	"testing"
	"time"
)

func Test_StrArray_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []string{"0", "1", "2", "3"}
		array := darray.NewStrArrayFrom(expect)
		array2 := darray.NewStrArrayFrom(expect, true)
		array3 := darray.NewStrArrayFrom([]string{})
		t.Assert(array.Slice(), expect)
		t.Assert(array.Interfaces(), expect)
		array.Set(0, "100")

		v, ok := array.Get(0)
		t.Assert(v, 100)
		t.Assert(ok, true)

		t.Assert(array.Search("100"), 0)
		t.Assert(array.Contains("100"), true)

		v, ok = array.Remove(0)
		t.Assert(v, 100)
		t.Assert(ok, true)

		v, ok = array.Remove(-1)
		t.Assert(v, "")
		t.Assert(ok, false)

		v, ok = array.Remove(100000)
		t.Assert(v, "")
		t.Assert(ok, false)

		t.Assert(array.Contains("100"), false)
		array.Append("4")
		t.Assert(array.Len(), 4)
		array.InsertBefore(0, "100")
		array.InsertAfter(0, "200")
		t.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "4"})
		array.InsertBefore(5, "300")
		array.InsertAfter(6, "400")
		t.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "300", "4", "400"})
		t.Assert(array.Clear().Len(), 0)
		t.Assert(array2.Slice(), expect)
		t.Assert(array3.Search("100"), -1)
	})
}

func TestStrArray_ContainsI(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := darray.NewStrArray()
		s.Append("a", "b", "C")
		t.Assert(s.Contains("A"), false)
		t.Assert(s.Contains("a"), true)
		t.Assert(s.ContainsI("A"), true)
	})
}

func TestStrArray_Sort(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect1 := []string{"0", "1", "2", "3"}
		expect2 := []string{"3", "2", "1", "0"}
		array := darray.NewStrArray()
		for i := 3; i >= 0; i-- {
			array.Append(dconv.String(i))
		}
		array.Sort()
		t.Assert(array.Slice(), expect1)
		array.Sort(true)
		t.Assert(array.Slice(), expect2)
	})
}

func TestStrArray_Unique(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []string{"1", "1", "2", "2", "3", "3", "2", "2"}
		array := darray.NewStrArrayFrom(expect)
		t.Assert(array.Unique().Slice(), []string{"1", "2", "3"})
	})
}

func TestStrArray_PushAndPop(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []string{"0", "1", "2", "3"}
		array := darray.NewStrArrayFrom(expect)
		t.Assert(array.Slice(), expect)

		v, ok := array.PopLeft()
		t.Assert(v, "0")
		t.Assert(ok, true)

		v, ok = array.PopRight()
		t.Assert(v, "3")
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.AssertIN(v, []string{"1", "2"})
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.AssertIN(v, []string{"1", "2"})
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.Assert(v, "")
		t.Assert(ok, false)

		t.Assert(array.Len(), 0)
		array.PushLeft("1").PushRight("2")
		t.Assert(array.Slice(), []string{"1", "2"})
	})
}

func TestStrArray_PopLeft(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArrayFrom(d.SliceStr{"1", "2", "3"})
		v, ok := array.PopLeft()
		t.Assert(v, 1)
		t.Assert(ok, true)
		t.Assert(array.Len(), 2)
		v, ok = array.PopLeft()
		t.Assert(v, 2)
		t.Assert(ok, true)
		t.Assert(array.Len(), 1)
		v, ok = array.PopLeft()
		t.Assert(v, 3)
		t.Assert(ok, true)
		t.Assert(array.Len(), 0)
	})
}

func TestStrArray_PopRight(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArrayFrom(d.SliceStr{"1", "2", "3"})

		v, ok := array.PopRight()
		t.Assert(v, 3)
		t.Assert(ok, true)
		t.Assert(array.Len(), 2)

		v, ok = array.PopRight()
		t.Assert(v, 2)
		t.Assert(ok, true)
		t.Assert(array.Len(), 1)

		v, ok = array.PopRight()
		t.Assert(v, 1)
		t.Assert(ok, true)
		t.Assert(array.Len(), 0)
	})
}

func TestStrArray_PopLefts(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArrayFrom(d.SliceStr{"1", "2", "3"})
		t.Assert(array.PopLefts(2), d.Slice{"1", "2"})
		t.Assert(array.Len(), 1)
		t.Assert(array.PopLefts(2), d.Slice{"3"})
		t.Assert(array.Len(), 0)
	})
}

func TestStrArray_PopRights(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArrayFrom(d.SliceStr{"1", "2", "3"})
		t.Assert(array.PopRights(2), d.Slice{"2", "3"})
		t.Assert(array.Len(), 1)
		t.Assert(array.PopLefts(2), d.Slice{"1"})
		t.Assert(array.Len(), 0)
	})
}

func TestStrArray_PopLeftsAndPopRights(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArray()
		v, ok := array.PopLeft()
		t.Assert(v, "")
		t.Assert(ok, false)
		t.Assert(array.PopLefts(10), nil)

		v, ok = array.PopRight()
		t.Assert(v, "")
		t.Assert(ok, false)
		t.Assert(array.PopRights(10), nil)

		v, ok = array.PopRand()
		t.Assert(v, "")
		t.Assert(ok, false)
		t.Assert(array.PopRands(10), nil)
	})

	dtest.C(t, func(t *dtest.T) {
		value1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		value2 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(value1)
		array2 := darray.NewStrArrayFrom(value2)
		t.Assert(array1.PopLefts(2), []interface{}{"0", "1"})
		t.Assert(array1.Slice(), []interface{}{"2", "3", "4", "5", "6"})
		t.Assert(array1.PopRights(2), []interface{}{"5", "6"})
		t.Assert(array1.Slice(), []interface{}{"2", "3", "4"})
		t.Assert(array1.PopRights(20), []interface{}{"2", "3", "4"})
		t.Assert(array1.Slice(), []interface{}{})
		t.Assert(array2.PopLefts(20), []interface{}{"0", "1", "2", "3", "4", "5", "6"})
		t.Assert(array2.Slice(), []interface{}{})
	})
}

func TestString_Range(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(value1)
		array2 := darray.NewStrArrayFrom(value1, true)
		t.Assert(array1.Range(0, 1), []interface{}{"0"})
		t.Assert(array1.Range(1, 2), []interface{}{"1"})
		t.Assert(array1.Range(0, 2), []interface{}{"0", "1"})
		t.Assert(array1.Range(-1, 10), value1)
		t.Assert(array1.Range(10, 1), nil)
		t.Assert(array2.Range(0, 1), []interface{}{"0"})
	})
}

func TestStrArray_Merge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a11 := []string{"0", "1", "2", "3"}
		a21 := []string{"4", "5", "6", "7"}
		array1 := darray.NewStrArrayFrom(a11)
		array2 := darray.NewStrArrayFrom(a21)
		t.Assert(array1.Merge(array2).Slice(), []string{"0", "1", "2", "3", "4", "5", "6", "7"})

		func1 := func(v1, v2 interface{}) int {
			if dconv.Int(v1) < dconv.Int(v2) {
				return 0
			}
			return 1
		}

		s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i1 := darray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := darray.NewArrayFrom([]interface{}{3})
		s3 := darray.NewStrArrayFrom([]string{"g", "h"})
		s4 := darray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := darray.NewSortedStrArrayFrom(s2)
		s6 := darray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := darray.NewStrArrayFrom(s1)

		t.Assert(a1.Merge(s2).Len(), 6)
		t.Assert(a1.Merge(i1).Len(), 9)
		t.Assert(a1.Merge(i2).Len(), 10)
		t.Assert(a1.Merge(s3).Len(), 12)
		t.Assert(a1.Merge(s4).Len(), 14)
		t.Assert(a1.Merge(s5).Len(), 16)
		t.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestStrArray_Fill(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0"}
		a2 := []string{"0"}
		array1 := darray.NewStrArrayFrom(a1)
		array2 := darray.NewStrArrayFrom(a2)
		t.Assert(array1.Fill(1, 2, "100"), nil)
		t.Assert(array1.Slice(), []string{"0", "100", "100"})
		t.Assert(array2.Fill(0, 2, "100"), nil)
		t.Assert(array2.Slice(), []string{"100", "100"})
		t.AssertNE(array2.Fill(-1, 2, "100"), nil)
		t.Assert(array2.Len(), 2)
	})
}

func TestStrArray_Chunk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"1", "2", "3", "4", "5"}
		array1 := darray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []string{"1", "2"})
		t.Assert(chunks[1], []string{"3", "4"})
		t.Assert(chunks[2], []string{"5"})
		t.Assert(len(array1.Chunk(0)), 0)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"1", "2", "3", "4", "5"}
		array1 := darray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []string{"1", "2", "3"})
		t.Assert(chunks[1], []string{"4", "5"})
		t.Assert(array1.Chunk(0), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []string{"1", "2"})
		t.Assert(chunks[1], []string{"3", "4"})
		t.Assert(chunks[2], []string{"5", "6"})
		t.Assert(array1.Chunk(0), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []string{"1", "2", "3"})
		t.Assert(chunks[1], []string{"4", "5", "6"})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestStrArray_Pad(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.Pad(3, "1").Slice(), []string{"0", "1", "1"})
		t.Assert(array1.Pad(-4, "1").Slice(), []string{"1", "0", "1", "1"})
		t.Assert(array1.Pad(3, "1").Slice(), []string{"1", "0", "1", "1"})
	})
}

func TestStrArray_SubSlice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		array2 := darray.NewStrArrayFrom(a1, true)
		t.Assert(array1.SubSlice(0, 2), []string{"0", "1"})
		t.Assert(array1.SubSlice(2, 2), []string{"2", "3"})
		t.Assert(array1.SubSlice(5, 8), []string{"5", "6"})
		t.Assert(array1.SubSlice(8, 2), nil)
		t.Assert(array1.SubSlice(1, -2), nil)
		t.Assert(array1.SubSlice(-5, 2), []string{"2", "3"})
		t.Assert(array1.SubSlice(-10, 1), nil)
		t.Assert(array2.SubSlice(0, 2), []string{"0", "1"})
	})
}

func TestStrArray_Rand(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(len(array1.Rands(2)), "2")
		t.Assert(len(array1.Rands(10)), 10)
		t.AssertIN(array1.Rands(1)[0], a1)
		v, ok := array1.Rand()
		t.Assert(ok, true)
		t.AssertIN(v, a1)
	})
}

func TestStrArray_PopRands(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"a", "b", "c", "d", "e", "f", "g"}
		array1 := darray.NewStrArrayFrom(a1)
		t.AssertIN(array1.PopRands(1), []string{"a", "b", "c", "d", "e", "f", "g"})
		t.AssertIN(array1.PopRands(1), []string{"a", "b", "c", "d", "e", "f", "g"})
		t.AssertNI(array1.PopRands(1), array1.Slice())
		t.AssertNI(array1.PopRands(1), array1.Slice())
		t.Assert(len(array1.PopRands(10)), 3)
	})
}

func TestStrArray_Shuffle(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.Shuffle().Len(), 7)
	})
}

func TestStrArray_Reverse(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.Reverse().Slice(), []string{"6", "5", "4", "3", "2", "1", "0"})
	})
}

func TestStrArray_Join(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.Join("."), `0.1.2.3.4.5.6`)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", `"a"`, `\a`}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.Join("."), `0.1."a".\a`)
	})
}

func TestStrArray_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.String(), `["0","1","2","3","4","5","6"]`)
	})
}

func TestNewStrArrayFromCopy(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := darray.NewStrArrayFromCopy(a1)
		a3 := darray.NewStrArrayFromCopy(a1, true)
		t.Assert(a2.Contains("1"), true)
		t.Assert(a2.Len(), 7)
		t.Assert(a2, a3)
	})
}

func TestStrArray_SetArray(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"a", "b", "c", "d"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.Contains("2"), true)
		t.Assert(array1.Len(), 7)

		array1 = array1.SetArray(a2)
		t.Assert(array1.Contains("2"), false)
		t.Assert(array1.Contains("c"), true)
		t.Assert(array1.Len(), 4)
	})
}

func TestStrArray_Replace(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"a", "b", "c", "d"}
		a3 := []string{"o", "p", "q", "x", "y", "z", "w", "r", "v"}
		array1 := darray.NewStrArrayFrom(a1)
		t.Assert(array1.Contains("2"), true)
		t.Assert(array1.Len(), 7)

		array1 = array1.Replace(a2)
		t.Assert(array1.Contains("2"), false)
		t.Assert(array1.Contains("c"), true)
		t.Assert(array1.Contains("5"), true)
		t.Assert(array1.Len(), 7)

		array1 = array1.Replace(a3)
		t.Assert(array1.Contains("2"), false)
		t.Assert(array1.Contains("c"), false)
		t.Assert(array1.Contains("5"), false)
		t.Assert(array1.Contains("p"), true)
		t.Assert(array1.Contains("r"), false)
		t.Assert(array1.Len(), 7)

	})
}

func TestStrArray_Sum(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"0", "a", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		array2 := darray.NewStrArrayFrom(a2)
		t.Assert(array1.Sum(), 21)
		t.Assert(array2.Sum(), 18)
	})
}

func TestStrArray_PopRand(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		str1, ok := array1.PopRand()
		t.Assert(strings.Contains("0,1,2,3,4,5,6", str1), true)
		t.Assert(array1.Len(), 6)
		t.Assert(ok, true)
	})
}

func TestStrArray_Clone(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := darray.NewStrArrayFrom(a1)
		array2 := array1.Clone()
		t.Assert(array2, array1)
		t.Assert(array2.Len(), 7)
	})
}

func TestStrArray_CountValues(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"0", "1", "2", "3", "4", "4", "6"}
		array1 := darray.NewStrArrayFrom(a1)

		m1 := array1.CountValues()
		t.Assert(len(m1), 6)
		t.Assert(m1["2"], 1)
		t.Assert(m1["4"], 2)
	})
}

func TestStrArray_Remove(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := darray.NewStrArrayFrom(a1)
		s1, ok := array1.Remove(1)
		t.Assert(s1, "a")
		t.Assert(ok, true)
		t.Assert(array1.Len(), 4)
		s1, ok = array1.Remove(3)
		t.Assert(s1, "c")
		t.Assert(ok, true)
		t.Assert(array1.Len(), 3)
	})
}

func TestStrArray_RLockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := darray.NewStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []string) { //读锁
			time.Sleep(2 * time.Second) //暂停1秒
			n1[2] = "g"
			ch2 <- dconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- dconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- dconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		t.AssertLT(t2-t1, 20) //go1加的读锁，所go2读的时候，并没有阻塞。
		t.Assert(a1.Contains("g"), true)
	})
}

func TestStrArray_SortFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []string{"a", "d", "c", "b"}
		a1 := darray.NewStrArrayFrom(s1)
		func1 := func(v1, v2 string) bool {
			return v1 < v2
		}
		a11 := a1.SortFunc(func1)
		t.Assert(a11, []string{"a", "b", "c", "d"})
	})
}

func TestStrArray_LockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []string{"a", "b", "c", "d"}
		a1 := darray.NewStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []string) { //读写锁
			time.Sleep(2 * time.Second) //暂停2秒
			n1[2] = "g"
			ch2 <- dconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- dconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- dconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		t.AssertGT(t2-t1, 20) //go1加的读写互斥锁，所go2读的时候被阻塞。
		t.Assert(a1.Contains("g"), true)
	})
}

func TestStrArray_Json(t *testing.T) {
	// array pointer
	dtest.C(t, func(t *dtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := darray.NewStrArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := darray.NewStrArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(a2.Slice(), s1)

		var a3 darray.StrArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Slice(), s1)
	})
	// array value
	dtest.C(t, func(t *dtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := *darray.NewStrArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := darray.NewStrArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(a2.Slice(), s1)

		var a3 darray.StrArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Slice(), s1)
	})
	// array pointer
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Name   string
			Scores *darray.StrArray
		}
		data := d.Map{
			"Name":   "john",
			"Scores": []string{"A+", "A", "A"},
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		user := new(User)
		err = json.UnmarshalUseNumber(b, user)
		t.Assert(err, nil)
		t.Assert(user.Name, data["Name"])
		t.Assert(user.Scores, data["Scores"])
	})
	// array value
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Name   string
			Scores darray.StrArray
		}
		data := d.Map{
			"Name":   "john",
			"Scores": []string{"A+", "A", "A"},
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		user := new(User)
		err = json.UnmarshalUseNumber(b, user)
		t.Assert(err, nil)
		t.Assert(user.Name, data["Name"])
		t.Assert(user.Scores, data["Scores"])
	})
}

func TestStrArray_Iterator(t *testing.T) {
	slice := d.SliceStr{"a", "b", "d", "c"}
	array := darray.NewStrArrayFrom(slice)
	dtest.C(t, func(t *dtest.T) {
		array.Iterator(func(k int, v string) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		array.IteratorAsc(func(k int, v string) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		array.IteratorDesc(func(k int, v string) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.Iterator(func(k int, v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.IteratorAsc(func(k int, v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.IteratorDesc(func(k int, v string) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
}

func TestStrArray_RemoveValue(t *testing.T) {
	slice := d.SliceStr{"a", "b", "d", "c"}
	array := darray.NewStrArrayFrom(slice)
	dtest.C(t, func(t *dtest.T) {
		t.Assert(array.RemoveValue("e"), false)
		t.Assert(array.RemoveValue("b"), true)
		t.Assert(array.RemoveValue("a"), true)
		t.Assert(array.RemoveValue("c"), true)
		t.Assert(array.RemoveValue("f"), false)
	})
}

func TestStrArray_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *darray.StrArray
	}
	// JSON
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(d.Map{
			"name":  "john",
			"array": []byte(`["1","2","3"]`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), d.SliceStr{"1", "2", "3"})
	})
	// Map
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(d.Map{
			"name":  "john",
			"array": d.SliceStr{"1", "2", "3"},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), d.SliceStr{"1", "2", "3"})
	})
}

func TestStrArray_FilterEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArrayFrom(d.SliceStr{"", "1", "2", "0"})
		t.Assert(array.FilterEmpty(), d.SliceStr{"1", "2", "0"})
	})
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArrayFrom(d.SliceStr{"1", "2"})
		t.Assert(array.FilterEmpty(), d.SliceStr{"1", "2"})
	})
}

func TestStrArray_Walk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewStrArrayFrom(d.SliceStr{"1", "2"})
		t.Assert(array.Walk(func(value string) string {
			return "key-" + value
		}), d.Slice{"key-1", "key-2"})
	})
}
