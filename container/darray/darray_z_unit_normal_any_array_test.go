// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package darray_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func Test_Array_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []interface{}{0, 1, 2, 3}
		array := darray.NewArrayFrom(expect)
		array2 := darray.NewArrayFrom(expect)
		array3 := darray.NewArrayFrom([]interface{}{})
		t.Assert(array.Slice(), expect)
		t.Assert(array.Interfaces(), expect)
		array.Set(0, 100)

		v, ok := array.Get(0)
		t.Assert(v, 100)
		t.Assert(ok, true)

		v, ok = array.Get(1)
		t.Assert(v, 1)
		t.Assert(ok, true)

		t.Assert(array.Search(100), 0)
		t.Assert(array3.Search(100), -1)
		t.Assert(array.Contains(100), true)

		v, ok = array.Remove(0)
		t.Assert(v, 100)
		t.Assert(ok, true)

		v, ok = array.Remove(-1)
		t.Assert(v, nil)
		t.Assert(ok, false)

		v, ok = array.Remove(100000)
		t.Assert(v, nil)
		t.Assert(ok, false)

		v, ok = array2.Remove(3)
		t.Assert(v, 3)
		t.Assert(ok, true)

		v, ok = array2.Remove(1)
		t.Assert(v, 1)
		t.Assert(ok, true)

		t.Assert(array.Contains(100), false)
		array.Append(4)
		t.Assert(array.Len(), 4)
		array.InsertBefore(0, 100)
		array.InsertAfter(0, 200)
		t.Assert(array.Slice(), []interface{}{100, 200, 2, 2, 3, 4})
		array.InsertBefore(5, 300)
		array.InsertAfter(6, 400)
		t.Assert(array.Slice(), []interface{}{100, 200, 2, 2, 3, 300, 4, 400})
		t.Assert(array.Clear().Len(), 0)
	})
}

func TestArray_Sort(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect1 := []interface{}{0, 1, 2, 3}
		expect2 := []interface{}{3, 2, 1, 0}
		array := darray.NewArray()
		for i := 3; i >= 0; i-- {
			array.Append(i)
		}
		array.SortFunc(func(v1, v2 interface{}) bool {
			return v1.(int) < v2.(int)
		})
		t.Assert(array.Slice(), expect1)
		array.SortFunc(func(v1, v2 interface{}) bool {
			return v1.(int) > v2.(int)
		})
		t.Assert(array.Slice(), expect2)
	})
}

func TestArray_Unique(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []interface{}{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array := darray.NewArrayFrom(expect)
		t.Assert(array.Unique().Slice(), []interface{}{1, 2, 3, 4, 5})
	})
}

func TestArray_PushAndPop(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		expect := []interface{}{0, 1, 2, 3}
		array := darray.NewArrayFrom(expect)
		t.Assert(array.Slice(), expect)

		v, ok := array.PopLeft()
		t.Assert(v, 0)
		t.Assert(ok, true)

		v, ok = array.PopRight()
		t.Assert(v, 3)
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.AssertIN(v, []interface{}{1, 2})
		t.Assert(ok, true)

		v, ok = array.PopRand()
		t.AssertIN(v, []interface{}{1, 2})
		t.Assert(ok, true)

		t.Assert(array.Len(), 0)
		array.PushLeft(1).PushRight(2)
		t.Assert(array.Slice(), []interface{}{1, 2})
	})
}

func TestArray_PopRands(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{100, 200, 300, 400, 500, 600}
		array := darray.NewFromCopy(a1)
		t.AssertIN(array.PopRands(2), []interface{}{100, 200, 300, 400, 500, 600})
	})
}

func TestArray_PopLeft(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewFrom(g.Slice{1, 2, 3})
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

func TestArray_PopRight(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewFrom(g.Slice{1, 2, 3})

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

func TestArray_PopLefts(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewFrom(g.Slice{1, 2, 3})
		t.Assert(array.PopLefts(2), g.Slice{1, 2})
		t.Assert(array.Len(), 1)
		t.Assert(array.PopLefts(2), g.Slice{3})
		t.Assert(array.Len(), 0)
	})
}

func TestArray_PopRights(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewFrom(g.Slice{1, 2, 3})
		t.Assert(array.PopRights(2), g.Slice{2, 3})
		t.Assert(array.Len(), 1)
		t.Assert(array.PopLefts(2), g.Slice{1})
		t.Assert(array.Len(), 0)
	})
}

func TestArray_PopLeftsAndPopRights(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New()
		v, ok := array.PopLeft()
		t.Assert(v, nil)
		t.Assert(ok, false)
		t.Assert(array.PopLefts(10), nil)

		v, ok = array.PopRight()
		t.Assert(v, nil)
		t.Assert(ok, false)
		t.Assert(array.PopRights(10), nil)

		v, ok = array.PopRand()
		t.Assert(v, nil)
		t.Assert(ok, false)
		t.Assert(array.PopRands(10), nil)
	})

	dtest.C(t, func(t *dtest.T) {
		value1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		value2 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(value1)
		array2 := darray.NewArrayFrom(value2)
		t.Assert(array1.PopLefts(2), []interface{}{0, 1})
		t.Assert(array1.Slice(), []interface{}{2, 3, 4, 5, 6})
		t.Assert(array1.PopRights(2), []interface{}{5, 6})
		t.Assert(array1.Slice(), []interface{}{2, 3, 4})
		t.Assert(array1.PopRights(20), []interface{}{2, 3, 4})
		t.Assert(array1.Slice(), []interface{}{})
		t.Assert(array2.PopLefts(20), []interface{}{0, 1, 2, 3, 4, 5, 6})
		t.Assert(array2.Slice(), []interface{}{})
	})
}

func TestArray_Range(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(value1)
		array2 := darray.NewArrayFrom(value1, true)
		t.Assert(array1.Range(0, 1), []interface{}{0})
		t.Assert(array1.Range(1, 2), []interface{}{1})
		t.Assert(array1.Range(0, 2), []interface{}{0, 1})
		t.Assert(array1.Range(-1, 10), value1)
		t.Assert(array1.Range(10, 2), nil)
		t.Assert(array2.Range(1, 3), []interface{}{1, 2})
	})
}

func TestArray_Merge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		i1 := []interface{}{0, 1, 2, 3}
		i2 := []interface{}{4, 5, 6, 7}
		array1 := darray.NewArrayFrom(i1)
		array2 := darray.NewArrayFrom(i2)
		t.Assert(array1.Merge(array2).Slice(), []interface{}{0, 1, 2, 3, 4, 5, 6, 7})

		//s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i3 := darray.NewIntArrayFrom([]int{1, 2, 3})
		i4 := darray.NewArrayFrom([]interface{}{3})
		s3 := darray.NewStrArrayFrom([]string{"g", "h"})
		s4 := darray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := darray.NewSortedStrArrayFrom(s2)
		s6 := darray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := darray.NewArrayFrom(i1)

		t.Assert(a1.Merge(s2).Len(), 6)
		t.Assert(a1.Merge(i3).Len(), 9)
		t.Assert(a1.Merge(i4).Len(), 10)
		t.Assert(a1.Merge(s3).Len(), 12)
		t.Assert(a1.Merge(s4).Len(), 14)
		t.Assert(a1.Merge(s5).Len(), 16)
		t.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestArray_Fill(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0}
		a2 := []interface{}{0}
		array1 := darray.NewArrayFrom(a1)
		array2 := darray.NewArrayFrom(a2, true)

		t.Assert(array1.Fill(1, 2, 100), nil)
		t.Assert(array1.Slice(), []interface{}{0, 100, 100})

		t.Assert(array2.Fill(0, 2, 100), nil)
		t.Assert(array2.Slice(), []interface{}{100, 100})

		t.AssertNE(array2.Fill(-1, 2, 100), nil)
		t.Assert(array2.Slice(), []interface{}{100, 100})
	})
}

func TestArray_Chunk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5}
		array1 := darray.NewArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []interface{}{1, 2})
		t.Assert(chunks[1], []interface{}{3, 4})
		t.Assert(chunks[2], []interface{}{5})
		t.Assert(array1.Chunk(0), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5}
		array1 := darray.NewArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []interface{}{1, 2, 3})
		t.Assert(chunks[1], []interface{}{4, 5})
		t.Assert(array1.Chunk(0), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []interface{}{1, 2})
		t.Assert(chunks[1], []interface{}{3, 4})
		t.Assert(chunks[2], []interface{}{5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []interface{}{1, 2, 3})
		t.Assert(chunks[1], []interface{}{4, 5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestArray_Pad(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0}
		array1 := darray.NewArrayFrom(a1)
		t.Assert(array1.Pad(3, 1).Slice(), []interface{}{0, 1, 1})
		t.Assert(array1.Pad(-4, 1).Slice(), []interface{}{1, 0, 1, 1})
		t.Assert(array1.Pad(3, 1).Slice(), []interface{}{1, 0, 1, 1})
	})
}

func TestArray_SubSlice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		array2 := darray.NewArrayFrom(a1, true)
		t.Assert(array1.SubSlice(0, 2), []interface{}{0, 1})
		t.Assert(array1.SubSlice(2, 2), []interface{}{2, 3})
		t.Assert(array1.SubSlice(5, 8), []interface{}{5, 6})
		t.Assert(array1.SubSlice(9, 1), nil)
		t.Assert(array1.SubSlice(-2, 2), []interface{}{5, 6})
		t.Assert(array1.SubSlice(-9, 2), nil)
		t.Assert(array1.SubSlice(1, -2), nil)
		t.Assert(array2.SubSlice(0, 2), []interface{}{0, 1})
	})
}

func TestArray_Rand(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		t.Assert(len(array1.Rands(2)), 2)
		t.Assert(len(array1.Rands(10)), 10)
		t.AssertIN(array1.Rands(1)[0], a1)
	})

	dtest.C(t, func(t *dtest.T) {
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := darray.NewArrayFrom(s1)
		i1, ok := a1.Rand()
		t.Assert(ok, true)
		t.Assert(a1.Contains(i1), true)
		t.Assert(a1.Len(), 4)
	})
}

func TestArray_Shuffle(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		t.Assert(array1.Shuffle().Len(), 7)
	})
}

func TestArray_Reverse(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		t.Assert(array1.Reverse().Slice(), []interface{}{6, 5, 4, 3, 2, 1, 0})
	})
}

func TestArray_Join(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		t.Assert(array1.Join("."), `0.1.2.3.4.5.6`)
	})

	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, `"a"`, `\a`}
		array1 := darray.NewArrayFrom(a1)
		t.Assert(array1.Join("."), `0.1."a".\a`)
	})
}

func TestArray_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := darray.NewArrayFrom(a1)
		t.Assert(array1.String(), `[0,1,2,3,4,5,6]`)
	})
}

func TestArray_Replace(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		a2 := []interface{}{"a", "b", "c"}
		a3 := []interface{}{"m", "n", "p", "z", "x", "y", "d", "u"}
		array1 := darray.NewArrayFrom(a1)
		array2 := array1.Replace(a2)
		t.Assert(array2.Len(), 7)
		t.Assert(array2.Contains("b"), true)
		t.Assert(array2.Contains(4), true)
		t.Assert(array2.Contains("v"), false)
		array3 := array1.Replace(a3)
		t.Assert(array3.Len(), 7)
		t.Assert(array3.Contains(4), false)
		t.Assert(array3.Contains("p"), true)
		t.Assert(array3.Contains("u"), false)
	})
}

func TestArray_SetArray(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		a2 := []interface{}{"a", "b", "c"}

		array1 := darray.NewArrayFrom(a1)
		array1 = array1.SetArray(a2)
		t.Assert(array1.Len(), 3)
		t.Assert(array1.Contains("b"), true)
		t.Assert(array1.Contains("5"), false)
	})
}

func TestArray_Sum(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3}
		a2 := []interface{}{"a", "b", "c"}
		a3 := []interface{}{"a", "1", "2"}

		array1 := darray.NewArrayFrom(a1)
		array2 := darray.NewArrayFrom(a2)
		array3 := darray.NewArrayFrom(a3)

		t.Assert(array1.Sum(), 6)
		t.Assert(array2.Sum(), 0)
		t.Assert(array3.Sum(), 3)

	})
}

func TestArray_Clone(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{0, 1, 2, 3}
		array1 := darray.NewArrayFrom(a1)
		array2 := array1.Clone()

		t.Assert(array1.Len(), 4)
		t.Assert(array2.Sum(), 6)
		t.AssertEQ(array1, array2)

	})
}

func TestArray_CountValues(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []interface{}{"a", "b", "c", "d", "e", "d"}
		array1 := darray.NewArrayFrom(a1)
		array2 := array1.CountValues()
		t.Assert(len(array2), 5)
		t.Assert(array2["b"], 1)
		t.Assert(array2["d"], 2)
	})
}

func TestArray_LockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := darray.NewArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []interface{}) { //读写锁
			time.Sleep(2 * time.Second) //暂停2秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		t.AssertGT(t2-t1, 20) //go1加的读写互斥锁，所go2读的时候被阻塞。
		t.Assert(a1.Contains("g"), true)
	})
}

func TestArray_RLockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := darray.NewArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []interface{}) { //读锁
			time.Sleep(2 * time.Second) //暂停1秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		t.AssertLT(t2-t1, 20) //go1加的读锁，所go2读的时候，并没有阻塞。
		t.Assert(a1.Contains("g"), true)
	})
}

func TestArray_Json(t *testing.T) {
	// pointer
	dtest.C(t, func(t *dtest.T) {
		s1 := []interface{}{"a", "b", "d", "c"}
		a1 := darray.NewArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := darray.New()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Slice(), s1)

		var a3 darray.Array
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Slice(), s1)
	})
	// value.
	dtest.C(t, func(t *dtest.T) {
		s1 := []interface{}{"a", "b", "d", "c"}
		a1 := *darray.NewArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := darray.New()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Slice(), s1)

		var a3 darray.Array
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Slice(), s1)
	})
	// pointer
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Name   string
			Scores *darray.Array
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []int{99, 100, 98},
		}
		b, err := json.Marshal(data)
		t.Assert(err, nil)

		user := new(User)
		err = json.UnmarshalUseNumber(b, user)
		t.Assert(err, nil)
		t.Assert(user.Name, data["Name"])
		t.Assert(user.Scores, data["Scores"])
	})
	// value
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Name   string
			Scores darray.Array
		}
		data := g.Map{
			"Name":   "john",
			"Scores": []int{99, 100, 98},
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

func TestArray_Iterator(t *testing.T) {
	slice := g.Slice{"a", "b", "d", "c"}
	array := darray.NewArrayFrom(slice)
	dtest.C(t, func(t *dtest.T) {
		array.Iterator(func(k int, v interface{}) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		array.IteratorAsc(func(k int, v interface{}) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		array.IteratorDesc(func(k int, v interface{}) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.Iterator(func(k int, v interface{}) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.IteratorAsc(func(k int, v interface{}) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.IteratorDesc(func(k int, v interface{}) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
}

func TestArray_RemoveValue(t *testing.T) {
	slice := g.Slice{"a", "b", "d", "c"}
	array := darray.NewArrayFrom(slice)
	dtest.C(t, func(t *dtest.T) {
		t.Assert(array.RemoveValue("e"), false)
		t.Assert(array.RemoveValue("b"), true)
		t.Assert(array.RemoveValue("a"), true)
		t.Assert(array.RemoveValue("c"), true)
		t.Assert(array.RemoveValue("f"), false)
	})
}

func TestArray_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *darray.Array
	}
	// JSON
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`[1,2,3]`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.Slice{1, 2, 3})
	})
	// Map
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name":  "john",
			"array": g.Slice{1, 2, 3},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.Slice{1, 2, 3})
	})
}

func TestArray_FilterNil(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		values := g.Slice{0, 1, 2, 3, 4, "", g.Slice{}}
		array := darray.NewArrayFromCopy(values)
		t.Assert(array.FilterNil().Slice(), values)
	})
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewArrayFromCopy(g.Slice{nil, 1, 2, 3, 4, nil})
		t.Assert(array.FilterNil(), g.Slice{1, 2, 3, 4})
	})
}

func TestArray_FilterEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewArrayFrom(g.Slice{0, 1, 2, 3, 4, "", g.Slice{}})
		t.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewArrayFrom(g.Slice{1, 2, 3, 4})
		t.Assert(array.FilterEmpty(), g.Slice{1, 2, 3, 4})
	})
}

func TestArray_Walk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewArrayFrom(g.Slice{"1", "2"})
		t.Assert(array.Walk(func(value interface{}) interface{} {
			return "key-" + gconv.String(value)
		}), g.Slice{"key-1", "key-2"})
	})
}
