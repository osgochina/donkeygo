// go test *.go

package dset_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/util/dconv"

	"github.com/osgochina/donkeygo/container/dset"
	"github.com/osgochina/donkeygo/test/dtest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestIntSet_Var(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var s dset.IntSet
		s.Add(1, 1, 2)
		s.Add([]int{3, 4}...)
		t.Assert(s.Size(), 4)
		t.AssertIN(1, s.Slice())
		t.AssertIN(2, s.Slice())
		t.AssertIN(3, s.Slice())
		t.AssertIN(4, s.Slice())
		t.AssertNI(0, s.Slice())
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), false)
		s.Remove(1)
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestIntSet_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet()
		s.Add(1, 1, 2)
		s.Add([]int{3, 4}...)
		t.Assert(s.Size(), 4)
		t.AssertIN(1, s.Slice())
		t.AssertIN(2, s.Slice())
		t.AssertIN(3, s.Slice())
		t.AssertIN(4, s.Slice())
		t.AssertNI(0, s.Slice())
		t.Assert(s.Contains(4), true)
		t.Assert(s.Contains(5), false)
		s.Remove(1)
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestIntSet_Iterator(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)

		a1 := darray.New(true)
		a2 := darray.New(true)
		s.Iterator(func(v int) bool {
			a1.Append(1)
			return false
		})
		s.Iterator(func(v int) bool {
			a2.Append(1)
			return true
		})
		t.Assert(a1.Len(), 1)
		t.Assert(a2.Len(), 3)
	})
}

func TestIntSet_LockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[int]struct{}) {
			delete(m, 1)
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[int]struct{}) {
			t.Assert(m, map[int]struct{}{
				3: struct{}{},
				2: struct{}{},
			})
		})
	})
}

func TestIntSet_Equal(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s2 := dset.NewIntSet()
		s3 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(1, 2, 3)
		s3.Add(1, 2, 3, 4)
		t.Assert(s1.Equal(s2), true)
		t.Assert(s1.Equal(s3), false)
	})
}

func TestIntSet_IsSubsetOf(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s2 := dset.NewIntSet()
		s3 := dset.NewIntSet()
		s1.Add(1, 2)
		s2.Add(1, 2, 3)
		s3.Add(1, 2, 3, 4)
		t.Assert(s1.IsSubsetOf(s2), true)
		t.Assert(s2.IsSubsetOf(s3), true)
		t.Assert(s1.IsSubsetOf(s3), true)
		t.Assert(s2.IsSubsetOf(s1), false)
		t.Assert(s3.IsSubsetOf(s2), false)
	})
}

func TestIntSet_Union(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s2 := dset.NewIntSet()
		s1.Add(1, 2)
		s2.Add(3, 4)
		s3 := s1.Union(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), true)
	})
}

func TestIntSet_Diff(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s2 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), false)
		t.Assert(s3.Contains(4), false)
	})
}

func TestIntSet_Intersect(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s2 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), false)
	})
}

func TestIntSet_Complement(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s2 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(4), true)
		t.Assert(s3.Contains(5), true)
	})
}

func TestIntSet_Size(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet(true)
		s1.Add(1, 2, 3)
		t.Assert(s1.Size(), 3)

	})

}

func TestIntSet_Merge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s2 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Merge(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(5), true)
		t.Assert(s3.Contains(6), false)
	})
}

func TestIntSet_Join(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s3 := s1.Join(",")
		t.Assert(strings.Contains(s3, "1"), true)
		t.Assert(strings.Contains(s3, "2"), true)
		t.Assert(strings.Contains(s3, "3"), true)
	})
}

func TestIntSet_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s3 := s1.String()
		t.Assert(strings.Contains(s3, "["), true)
		t.Assert(strings.Contains(s3, "]"), true)
		t.Assert(strings.Contains(s3, "1"), true)
		t.Assert(strings.Contains(s3, "2"), true)
		t.Assert(strings.Contains(s3, "3"), true)
	})
}

func TestIntSet_Sum(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSet()
		s1.Add(1, 2, 3)
		s2 := dset.NewIntSet()
		s2.Add(5, 6, 7)
		t.Assert(s2.Sum(), 18)

	})

}

func TestIntSet_Pop(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet()
		s.Add(4, 2, 3)
		t.Assert(s.Size(), 3)
		t.AssertIN(s.Pop(), []int{4, 2, 3})
		t.AssertIN(s.Pop(), []int{4, 2, 3})
		t.Assert(s.Size(), 1)
	})
}

func TestIntSet_Pops(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet()
		s.Add(1, 4, 2, 3)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(1), []int{1, 4, 2, 3})
		t.Assert(s.Size(), 3)
		a := s.Pops(2)
		t.Assert(len(a), 2)
		t.AssertIN(a, []int{1, 4, 2, 3})
		t.Assert(s.Size(), 1)
	})

	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet(true)
		a := []int{1, 2, 3, 4}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestIntSet_AddIfNotExist(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.AddIfNotExist(1), false)
		t.Assert(s.AddIfNotExist(2), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExist(2), false)
		t.Assert(s.Contains(2), true)
	})
}

func TestIntSet_AddIfNotExistFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return false }), false)
		t.Assert(s.Contains(2), false)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return true }), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExistFunc(2, func() bool { return true }), false)
		t.Assert(s.Contains(2), true)
	})
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet(true)
		wd := sync.WaitGroup{}
		wd.Add(1)
		go func() {
			defer wd.Done()
			r := s.AddIfNotExistFunc(1, func() bool {
				time.Sleep(100 * time.Millisecond)
				return true
			})
			t.Assert(r, false)
		}()
		s.Add(1)
		wd.Wait()
	})
}

func TestIntSet_AddIfNotExistFuncLock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewIntSet(true)
		wd := sync.WaitGroup{}
		wd.Add(2)
		go func() {
			defer wd.Done()
			r := s.AddIfNotExistFuncLock(1, func() bool {
				time.Sleep(500 * time.Millisecond)
				return true
			})
			t.Assert(r, true)
		}()
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wd.Done()
			r := s.AddIfNotExistFuncLock(1, func() bool {
				return true
			})
			t.Assert(r, false)
		}()
		wd.Wait()
	})
}

func TestIntSet_Json(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []int{1, 3, 2, 4}
		a1 := dset.NewIntSetFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := dset.NewIntSet()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains(1), true)
		t.Assert(a2.Contains(2), true)
		t.Assert(a2.Contains(3), true)
		t.Assert(a2.Contains(4), true)
		t.Assert(a2.Contains(5), false)

		var a3 dset.IntSet
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a2.Contains(1), true)
		t.Assert(a2.Contains(2), true)
		t.Assert(a2.Contains(3), true)
		t.Assert(a2.Contains(4), true)
		t.Assert(a2.Contains(5), false)
	})
}

func TestIntSet_Walk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var set dset.IntSet
		set.Add(d.SliceInt{1, 2}...)
		set.Walk(func(item int) int {
			return item + 10
		})
		t.Assert(set.Size(), 2)
		t.Assert(set.Contains(11), true)
		t.Assert(set.Contains(12), true)
	})
}

func TestIntSet_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Set  *dset.IntSet
	}
	// JSON
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(d.Map{
			"name": "john",
			"set":  []byte(`[1,2,3]`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains(1), true)
		t.Assert(v.Set.Contains(2), true)
		t.Assert(v.Set.Contains(3), true)
		t.Assert(v.Set.Contains(4), false)
	})
	// Map
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(d.Map{
			"name": "john",
			"set":  d.Slice{1, 2, 3},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains(1), true)
		t.Assert(v.Set.Contains(2), true)
		t.Assert(v.Set.Contains(3), true)
		t.Assert(v.Set.Contains(4), false)
	})
}
