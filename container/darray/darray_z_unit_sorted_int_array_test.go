// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package darray_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/util/dconv"

	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestNewSortedIntArrayFrom(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{0, 3, 2, 1, 4, 5, 6}
		array1 := darray.NewSortedIntArrayFrom(a1, true)
		t.Assert(array1.Join("."), "0.1.2.3.4.5.6")
		t.Assert(array1.Slice(), a1)
		t.Assert(array1.Interfaces(), a1)
	})
}

func TestNewSortedIntArrayFromCopy(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{0, 5, 2, 1, 4, 3, 6}
		array1 := darray.NewSortedIntArrayFromCopy(a1, false)
		t.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestSortedIntArray_SetArray(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{0, 1, 2, 3}
		a2 := []int{4, 5, 6}
		array1 := darray.NewSortedIntArrayFrom(a1)
		array2 := array1.SetArray(a2)

		t.Assert(array2.Len(), 3)
		t.Assert(array2.Search(3), -1)
		t.Assert(array2.Search(5), 1)
		t.Assert(array2.Search(6), 2)
	})
}

func TestSortedIntArray_Sort(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{0, 3, 2, 1}

		array1 := darray.NewSortedIntArrayFrom(a1)
		array2 := array1.Sort()

		t.Assert(array2.Len(), 4)
		t.Assert(array2, []int{0, 1, 2, 3})
	})
}

func TestSortedIntArray_Get(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 0}
		array1 := darray.NewSortedIntArrayFrom(a1)
		v, ok := array1.Get(0)
		t.Assert(v, 0)
		t.Assert(ok, true)

		v, ok = array1.Get(1)
		t.Assert(v, 1)
		t.Assert(ok, true)

		v, ok = array1.Get(3)
		t.Assert(v, 5)
		t.Assert(ok, true)
	})
}

func TestSortedIntArray_Remove(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 0}
		array1 := darray.NewSortedIntArrayFrom(a1)

		v, ok := array1.Remove(-1)
		t.Assert(v, 0)
		t.Assert(ok, false)

		v, ok = array1.Remove(-100000)
		t.Assert(v, 0)
		t.Assert(ok, false)

		v, ok = array1.Remove(2)
		t.Assert(v, 3)
		t.Assert(ok, true)

		t.Assert(array1.Search(5), 2)

		v, ok = array1.Remove(0)
		t.Assert(v, 0)
		t.Assert(ok, true)

		t.Assert(array1.Search(5), 1)

		a2 := []int{1, 3, 4}
		array2 := darray.NewSortedIntArrayFrom(a2)

		v, ok = array2.Remove(1)
		t.Assert(v, 3)
		t.Assert(ok, true)
		t.Assert(array2.Search(1), 0)

		v, ok = array2.Remove(1)
		t.Assert(v, 4)
		t.Assert(ok, true)

		t.Assert(array2.Search(4), -1)
	})
}

func TestSortedIntArray_PopLeft(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 2}
		array1 := darray.NewSortedIntArrayFrom(a1)
		v, ok := array1.PopLeft()
		t.Assert(v, 1)
		t.Assert(ok, true)
		t.Assert(array1.Len(), 3)
		t.Assert(array1.Search(1), -1)
	})
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewSortedIntArrayFrom(g.SliceInt{1, 2, 3})
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

func TestSortedIntArray_PopRight(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 2}
		array1 := darray.NewSortedIntArrayFrom(a1)
		v, ok := array1.PopRight()
		t.Assert(v, 5)
		t.Assert(ok, true)
		t.Assert(array1.Len(), 3)
		t.Assert(array1.Search(5), -1)
	})
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewSortedIntArrayFrom(g.SliceInt{1, 2, 3})
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

func TestSortedIntArray_PopRand(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 2}
		array1 := darray.NewSortedIntArrayFrom(a1)
		i1, ok := array1.PopRand()
		t.Assert(ok, true)
		t.Assert(array1.Len(), 3)
		t.Assert(array1.Search(i1), -1)
		t.AssertIN(i1, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_PopRands(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 2}
		array1 := darray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopRands(2)
		t.Assert(array1.Len(), 2)
		t.AssertIN(ns1, []int{1, 3, 5, 2})

		a2 := []int{1, 3, 5, 2}
		array2 := darray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopRands(5)
		t.Assert(array2.Len(), 0)
		t.Assert(len(ns2), 4)
		t.AssertIN(ns2, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_Empty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewSortedIntArray()
		v, ok := array.PopLeft()
		t.Assert(v, 0)
		t.Assert(ok, false)
		t.Assert(array.PopLefts(10), nil)

		v, ok = array.PopRight()
		t.Assert(v, 0)
		t.Assert(ok, false)
		t.Assert(array.PopRights(10), nil)

		v, ok = array.PopRand()
		t.Assert(v, 0)
		t.Assert(ok, false)
		t.Assert(array.PopRands(10), nil)
	})
}

func TestSortedIntArray_PopLefts(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 2}
		array1 := darray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopLefts(2)
		t.Assert(array1.Len(), 2)
		t.Assert(ns1, []int{1, 2})

		a2 := []int{1, 3, 5, 2}
		array2 := darray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopLefts(5)
		t.Assert(array2.Len(), 0)
		t.AssertIN(ns2, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_PopRights(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 2}
		array1 := darray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopRights(2)
		t.Assert(array1.Len(), 2)
		t.Assert(ns1, []int{3, 5})

		a2 := []int{1, 3, 5, 2}
		array2 := darray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopRights(5)
		t.Assert(array2.Len(), 0)
		t.AssertIN(ns2, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_Range(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5, 2, 6, 7}
		array1 := darray.NewSortedIntArrayFrom(a1)
		array2 := darray.NewSortedIntArrayFrom(a1, true)
		ns1 := array1.Range(1, 4)
		t.Assert(len(ns1), 3)
		t.Assert(ns1, []int{2, 3, 5})

		ns2 := array1.Range(5, 4)
		t.Assert(len(ns2), 0)

		ns3 := array1.Range(-1, 4)
		t.Assert(len(ns3), 4)

		nsl := array1.Range(5, 8)
		t.Assert(len(nsl), 1)
		t.Assert(array2.Range(1, 2), []int{2})
	})
}

func TestSortedIntArray_Sum(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		n1 := array1.Sum()
		t.Assert(n1, 9)
	})
}

func TestSortedIntArray_Join(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		t.Assert(array1.Join("."), `1.3.5`)
	})
}

func TestSortedIntArray_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		t.Assert(array1.String(), `[1,3,5]`)
	})
}

func TestSortedIntArray_Contains(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		t.Assert(array1.Contains(4), false)
	})
}

func TestSortedIntArray_Clone(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		array2 := array1.Clone()
		t.Assert(array2.Len(), 3)
		t.Assert(array2, array1)
	})
}

func TestSortedIntArray_Clear(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 3, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		array1.Clear()
		t.Assert(array1.Len(), 0)
	})
}

func TestSortedIntArray_Chunk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Chunk(2) //按每几个元素切成一个数组
		ns2 := array1.Chunk(-1)
		t.Assert(len(ns1), 3)
		t.Assert(ns1[0], []int{1, 2})
		t.Assert(ns1[2], []int{5})
		t.Assert(len(ns2), 0)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []int{1, 2, 3})
		t.Assert(chunks[1], []int{4, 5})
		t.Assert(array1.Chunk(0), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := darray.NewSortedIntArrayFrom(a1)
		chunks := array1.Chunk(2)
		t.Assert(len(chunks), 3)
		t.Assert(chunks[0], []int{1, 2})
		t.Assert(chunks[1], []int{3, 4})
		t.Assert(chunks[2], []int{5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 6}
		array1 := darray.NewSortedIntArrayFrom(a1)
		chunks := array1.Chunk(3)
		t.Assert(len(chunks), 2)
		t.Assert(chunks[0], []int{1, 2, 3})
		t.Assert(chunks[1], []int{4, 5, 6})
		t.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedIntArray_SubSlice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		array2 := darray.NewSortedIntArrayFrom(a1, true)
		ns1 := array1.SubSlice(1, 2)
		t.Assert(len(ns1), 2)
		t.Assert(ns1, []int{2, 3})

		ns2 := array1.SubSlice(7, 2)
		t.Assert(len(ns2), 0)

		ns3 := array1.SubSlice(3, 5)
		t.Assert(len(ns3), 2)
		t.Assert(ns3, []int{4, 5})

		ns4 := array1.SubSlice(3, 1)
		t.Assert(len(ns4), 1)
		t.Assert(ns4, []int{4})
		t.Assert(array1.SubSlice(-1, 1), []int{5})
		t.Assert(array1.SubSlice(-9, 1), nil)
		t.Assert(array1.SubSlice(1, -9), nil)
		t.Assert(array2.SubSlice(1, 2), []int{2, 3})
	})
}

func TestSortedIntArray_Rand(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		ns1, ok := array1.Rand()
		t.AssertIN(ns1, a1)
		t.Assert(ok, true)
	})
}

func TestSortedIntArray_Rands(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Rands(2)
		t.AssertIN(ns1, a1)
		t.Assert(len(ns1), 2)

		ns2 := array1.Rands(6)
		t.Assert(len(ns2), 6)
	})
}

func TestSortedIntArray_CountValues(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 3}
		array1 := darray.NewSortedIntArrayFrom(a1)
		ns1 := array1.CountValues() //按每几个元素切成一个数组
		t.Assert(len(ns1), 5)
		t.Assert(ns1[2], 1)
		t.Assert(ns1[3], 2)
	})
}

func TestSortedIntArray_SetUnique(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		array1.SetUnique(true)
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []int{1, 2, 3, 4, 5})
	})
}

func TestSortedIntArray_Unique(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a1 := []int{1, 2, 3, 4, 5, 3, 2, 2, 3, 5, 5}
		array1 := darray.NewSortedIntArrayFrom(a1)
		array1.Unique()
		t.Assert(array1.Len(), 5)
		t.Assert(array1, []int{1, 2, 3, 4, 5})
	})
}

func TestSortedIntArray_LockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []int{1, 2, 3, 4}
		a1 := darray.NewSortedIntArrayFrom(s1, true)
		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []int) { //读写锁
			time.Sleep(2 * time.Second) //暂停2秒
			n1[2] = 6
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
		t.Assert(a1.Contains(6), true)
	})
}

func TestSortedIntArray_RLockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []int{1, 2, 3, 4}
		a1 := darray.NewSortedIntArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []int) { //读锁
			time.Sleep(2 * time.Second) //暂停1秒
			n1[2] = 6
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
		t.Assert(a1.Contains(6), true)
	})
}

func TestSortedIntArray_Merge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		func1 := func(v1, v2 interface{}) int {
			if dconv.Int(v1) < dconv.Int(v2) {
				return 0
			}
			return 1
		}
		i0 := []int{1, 2, 3, 4}
		s2 := []string{"e", "f"}
		i1 := darray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := darray.NewArrayFrom([]interface{}{3})
		s3 := darray.NewStrArrayFrom([]string{"g", "h"})
		s4 := darray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := darray.NewSortedStrArrayFrom(s2)
		s6 := darray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := darray.NewSortedIntArrayFrom(i0)

		t.Assert(a1.Merge(s2).Len(), 6)
		t.Assert(a1.Merge(i1).Len(), 9)
		t.Assert(a1.Merge(i2).Len(), 10)
		t.Assert(a1.Merge(s3).Len(), 12)
		t.Assert(a1.Merge(s4).Len(), 14)
		t.Assert(a1.Merge(s5).Len(), 16)
		t.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestSortedIntArray_Json(t *testing.T) {
	// array pointer
	dtest.C(t, func(t *dtest.T) {
		s1 := []int{1, 4, 3, 2}
		s2 := []int{1, 2, 3, 4}
		a1 := darray.NewSortedIntArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := darray.NewSortedIntArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(a2.Slice(), s2)

		var a3 darray.SortedIntArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Slice(), s1)
	})
	// array value
	dtest.C(t, func(t *dtest.T) {
		s1 := []int{1, 4, 3, 2}
		s2 := []int{1, 2, 3, 4}
		a1 := *darray.NewSortedIntArrayFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(b1, b2)
		t.Assert(err1, err2)

		a2 := darray.NewSortedIntArray()
		err1 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(a2.Slice(), s2)

		var a3 darray.SortedIntArray
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Slice(), s1)
	})
	// array pointer
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Name   string
			Scores *darray.SortedIntArray
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
		t.Assert(user.Scores, []int{98, 99, 100})
	})
	// array value
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Name   string
			Scores darray.SortedIntArray
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
		t.Assert(user.Scores, []int{98, 99, 100})
	})
}

func TestSortedIntArray_Iterator(t *testing.T) {
	slice := g.SliceInt{10, 20, 30, 40}
	array := darray.NewSortedIntArrayFrom(slice)
	dtest.C(t, func(t *dtest.T) {
		array.Iterator(func(k int, v int) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		array.IteratorAsc(func(k int, v int) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		array.IteratorDesc(func(k int, v int) bool {
			t.Assert(v, slice[k])
			return true
		})
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.Iterator(func(k int, v int) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.IteratorAsc(func(k int, v int) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
	dtest.C(t, func(t *dtest.T) {
		index := 0
		array.IteratorDesc(func(k int, v int) bool {
			index++
			return false
		})
		t.Assert(index, 1)
	})
}

func TestSortedIntArray_RemoveValue(t *testing.T) {
	slice := g.SliceInt{10, 20, 30, 40}
	array := darray.NewSortedIntArrayFrom(slice)
	dtest.C(t, func(t *dtest.T) {
		t.Assert(array.RemoveValue(99), false)
		t.Assert(array.RemoveValue(20), true)
		t.Assert(array.RemoveValue(10), true)
		t.Assert(array.RemoveValue(20), false)
		t.Assert(array.RemoveValue(88), false)
		t.Assert(array.Len(), 2)
	})
}

func TestSortedIntArray_UnmarshalValue(t *testing.T) {
	type V struct {
		Name  string
		Array *darray.SortedIntArray
	}
	// JSON
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(g.Map{
			"name":  "john",
			"array": []byte(`[2,3,1]`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.Slice{1, 2, 3})
	})
	// Map
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(g.Map{
			"name":  "john",
			"array": g.Slice{2, 3, 1},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Array.Slice(), g.Slice{1, 2, 3})
	})
}

func TestSortedIntArray_FilterEmpty(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewSortedIntArrayFrom(g.SliceInt{0, 1, 2, 3, 4, 0})
		t.Assert(array.FilterEmpty(), g.SliceInt{1, 2, 3, 4})
	})
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewSortedIntArrayFrom(g.SliceInt{1, 2, 3, 4})
		t.Assert(array.FilterEmpty(), g.SliceInt{1, 2, 3, 4})
	})
}

func TestSortedIntArray_Walk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.NewSortedIntArrayFrom(g.SliceInt{1, 2})
		t.Assert(array.Walk(func(value int) int {
			return 10 + value
		}), g.Slice{11, 12})
	})
}
