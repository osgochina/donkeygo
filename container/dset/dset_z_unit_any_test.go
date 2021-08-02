// go test *.go

package dset_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/container/dset"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSet_Var(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var s dset.Set
		s.Add(1, 1, 2)
		s.Add([]interface{}{3, 4}...)
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

func TestSet_New(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.New()
		s.Add(1, 1, 2)
		s.Add([]interface{}{3, 4}...)
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

func TestSet_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewSet()
		s.Add(1, 1, 2)
		s.Add([]interface{}{3, 4}...)
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

func TestSet_Iterator(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewSet()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)

		a1 := darray.New(true)
		a2 := darray.New(true)
		s.Iterator(func(v interface{}) bool {
			a1.Append(1)
			return false
		})
		s.Iterator(func(v interface{}) bool {
			a2.Append(1)
			return true
		})
		t.Assert(a1.Len(), 1)
		t.Assert(a2.Len(), 3)
	})
}

func TestSet_LockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewSet()
		s.Add(1, 2, 3)
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[interface{}]struct{}) {
			delete(m, 1)
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[interface{}]struct{}) {
			t.Assert(m, map[interface{}]struct{}{
				3: struct{}{},
				2: struct{}{},
			})
		})
	})
}

func TestSet_Equal(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewSet()
		s2 := dset.NewSet()
		s3 := dset.NewSet()
		s1.Add(1, 2, 3)
		s2.Add(1, 2, 3)
		s3.Add(1, 2, 3, 4)
		t.Assert(s1.Equal(s2), true)
		t.Assert(s1.Equal(s3), false)
	})
}

func TestSet_IsSubsetOf(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewSet()
		s2 := dset.NewSet()
		s3 := dset.NewSet()
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

func TestSet_Union(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewSet()
		s2 := dset.NewSet()
		s1.Add(1, 2)
		s2.Add(3, 4)
		s3 := s1.Union(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), true)
	})
}

func TestSet_Diff(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewSet()
		s2 := dset.NewSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains(1), true)
		t.Assert(s3.Contains(2), true)
		t.Assert(s3.Contains(3), false)
		t.Assert(s3.Contains(4), false)
	})
}

func TestSet_Intersect(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewSet()
		s2 := dset.NewSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(3), true)
		t.Assert(s3.Contains(4), false)
	})
}

func TestSet_Complement(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewSet()
		s2 := dset.NewSet()
		s1.Add(1, 2, 3)
		s2.Add(3, 4, 5)
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains(1), false)
		t.Assert(s3.Contains(2), false)
		t.Assert(s3.Contains(4), true)
		t.Assert(s3.Contains(5), true)
	})
}

func TestNewFrom(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewFrom("a")
		s2 := dset.NewFrom("b", false)
		s3 := dset.NewFrom(3, true)
		s4 := dset.NewFrom([]string{"s1", "s2"}, true)
		t.Assert(s1.Contains("a"), true)
		t.Assert(s2.Contains("b"), true)
		t.Assert(s3.Contains(3), true)
		t.Assert(s4.Contains("s1"), true)
		t.Assert(s4.Contains("s3"), false)

	})
}

func TestNew(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.New()
		s1.Add("a", 2)
		s2 := dset.New(true)
		s2.Add("b", 3)
		t.Assert(s1.Contains("a"), true)

	})
}

func TestSet_Join(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.New(true)
		s1.Add("a", "a1", "b", "c")
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, "a1"), true)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.New(true)
		s1.Add("a", `"b"`, `\c`)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, `"b"`), true)
		t.Assert(strings.Contains(str1, `\c`), true)
		t.Assert(strings.Contains(str1, `a`), true)
	})
}

func TestSet_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.New(true)
		s1.Add("a", "a2", "b", "c")
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "["), true)
		t.Assert(strings.Contains(str1, "]"), true)
		t.Assert(strings.Contains(str1, "a2"), true)
	})
}

func TestSet_Merge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.New(true)
		s2 := dset.New(true)
		s1.Add("a", "a2", "b", "c")
		s2.Add("b", "b1", "e", "f")
		ss := s1.Merge(s2)
		t.Assert(ss.Contains("a2"), true)
		t.Assert(ss.Contains("b1"), true)

	})
}

func TestSet_Sum(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.New(true)
		s1.Add(1, 2, 3, 4)
		t.Assert(s1.Sum(), int(10))

	})
}

func TestSet_Pop(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.New(true)
		s.Add(1, 2, 3, 4)
		t.Assert(s.Size(), 4)
		t.AssertIN(s.Pop(), []int{1, 2, 3, 4})
		t.Assert(s.Size(), 3)
	})
}

func TestSet_Pops(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.New(true)
		s.Add(1, 2, 3, 4)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(1), []int{1, 2, 3, 4})
		t.Assert(s.Size(), 3)
		a := s.Pops(6)
		t.Assert(len(a), 3)
		t.AssertIN(a, []int{1, 2, 3, 4})
		t.Assert(s.Size(), 0)
	})

	dtest.C(t, func(t *dtest.T) {
		s := dset.New(true)
		a := []interface{}{1, 2, 3, 4}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestSet_Json(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []interface{}{"a", "b", "d", "c"}
		a1 := dset.NewFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := dset.New()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains("a"), true)
		t.Assert(a2.Contains("b"), true)
		t.Assert(a2.Contains("c"), true)
		t.Assert(a2.Contains("d"), true)
		t.Assert(a2.Contains("e"), false)

		var a3 dset.Set
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Contains("a"), true)
		t.Assert(a3.Contains("b"), true)
		t.Assert(a3.Contains("c"), true)
		t.Assert(a3.Contains("d"), true)
		t.Assert(a3.Contains("e"), false)
	})
}

func TestSet_AddIfNotExist(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.New(true)
		s.Add(1)
		t.Assert(s.Contains(1), true)
		t.Assert(s.AddIfNotExist(1), false)
		t.Assert(s.AddIfNotExist(2), true)
		t.Assert(s.Contains(2), true)
		t.Assert(s.AddIfNotExist(2), false)
		t.Assert(s.Contains(2), true)
	})
}

func TestSet_AddIfNotExistFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.New(true)
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
		s := dset.New(true)
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

func TestSet_Walk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var set dset.Set
		set.Add(d.Slice{1, 2}...)
		set.Walk(func(item interface{}) interface{} {
			return dconv.Int(item) + 10
		})
		t.Assert(set.Size(), 2)
		t.Assert(set.Contains(11), true)
		t.Assert(set.Contains(12), true)
	})
}

func TestSet_AddIfNotExistFuncLock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.New(true)
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

func TestSet_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Set  *dset.Set
	}
	// JSON
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(d.Map{
			"name": "john",
			"set":  []byte(`["k1","k2","k3"]`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("k1"), true)
		t.Assert(v.Set.Contains("k2"), true)
		t.Assert(v.Set.Contains("k3"), true)
		t.Assert(v.Set.Contains("k4"), false)
	})
	// Map
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(d.Map{
			"name": "john",
			"set":  d.Slice{"k1", "k2", "k3"},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("k1"), true)
		t.Assert(v.Set.Contains("k2"), true)
		t.Assert(v.Set.Contains("k3"), true)
		t.Assert(v.Set.Contains("k4"), false)
	})
}
