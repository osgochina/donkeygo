// go test *.go

package dset_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/container/dset"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStrSet_Var(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var s dset.StrSet
		s.Add("1", "1", "2")
		s.Add([]string{"3", "4"}...)
		t.Assert(s.Size(), 4)
		t.AssertIN("1", s.Slice())
		t.AssertIN("2", s.Slice())
		t.AssertIN("3", s.Slice())
		t.AssertIN("4", s.Slice())
		t.AssertNI("0", s.Slice())
		t.Assert(s.Contains("4"), true)
		t.Assert(s.Contains("5"), false)
		s.Remove("1")
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestStrSet_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet()
		s.Add("1", "1", "2")
		s.Add([]string{"3", "4"}...)
		t.Assert(s.Size(), 4)
		t.AssertIN("1", s.Slice())
		t.AssertIN("2", s.Slice())
		t.AssertIN("3", s.Slice())
		t.AssertIN("4", s.Slice())
		t.AssertNI("0", s.Slice())
		t.Assert(s.Contains("4"), true)
		t.Assert(s.Contains("5"), false)
		s.Remove("1")
		t.Assert(s.Size(), 3)
		s.Clear()
		t.Assert(s.Size(), 0)
	})
}

func TestStrSet_ContainsI(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet()
		s.Add("a", "b", "C")
		t.Assert(s.Contains("A"), false)
		t.Assert(s.Contains("a"), true)
		t.Assert(s.ContainsI("A"), true)
	})
}

func TestStrSet_Iterator(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet()
		s.Add("1", "2", "3")
		t.Assert(s.Size(), 3)

		a1 := darray.New(true)
		a2 := darray.New(true)
		s.Iterator(func(v string) bool {
			a1.Append("1")
			return false
		})
		s.Iterator(func(v string) bool {
			a2.Append("1")
			return true
		})
		t.Assert(a1.Len(), 1)
		t.Assert(a2.Len(), 3)
	})
}

func TestStrSet_LockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet()
		s.Add("1", "2", "3")
		t.Assert(s.Size(), 3)
		s.LockFunc(func(m map[string]struct{}) {
			delete(m, "1")
		})
		t.Assert(s.Size(), 2)
		s.RLockFunc(func(m map[string]struct{}) {
			t.Assert(m, map[string]struct{}{
				"3": struct{}{},
				"2": struct{}{},
			})
		})
	})
}

func TestStrSet_Equal(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s2 := dset.NewStrSet()
		s3 := dset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("1", "2", "3")
		s3.Add("1", "2", "3", "4")
		t.Assert(s1.Equal(s2), true)
		t.Assert(s1.Equal(s3), false)
	})
}

func TestStrSet_IsSubsetOf(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s2 := dset.NewStrSet()
		s3 := dset.NewStrSet()
		s1.Add("1", "2")
		s2.Add("1", "2", "3")
		s3.Add("1", "2", "3", "4")
		t.Assert(s1.IsSubsetOf(s2), true)
		t.Assert(s2.IsSubsetOf(s3), true)
		t.Assert(s1.IsSubsetOf(s3), true)
		t.Assert(s2.IsSubsetOf(s1), false)
		t.Assert(s3.IsSubsetOf(s2), false)
	})
}

func TestStrSet_Union(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s2 := dset.NewStrSet()
		s1.Add("1", "2")
		s2.Add("3", "4")
		s3 := s1.Union(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("2"), true)
		t.Assert(s3.Contains("3"), true)
		t.Assert(s3.Contains("4"), true)
	})
}

func TestStrSet_Diff(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s2 := dset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Diff(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("2"), true)
		t.Assert(s3.Contains("3"), false)
		t.Assert(s3.Contains("4"), false)
	})
}

func TestStrSet_Intersect(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s2 := dset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Intersect(s2)
		t.Assert(s3.Contains("1"), false)
		t.Assert(s3.Contains("2"), false)
		t.Assert(s3.Contains("3"), true)
		t.Assert(s3.Contains("4"), false)
	})
}

func TestStrSet_Complement(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s2 := dset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Complement(s2)
		t.Assert(s3.Contains("1"), false)
		t.Assert(s3.Contains("2"), false)
		t.Assert(s3.Contains("4"), true)
		t.Assert(s3.Contains("5"), true)
	})
}

func TestNewIntSetFrom(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewIntSetFrom([]int{1, 2, 3, 4})
		s2 := dset.NewIntSetFrom([]int{5, 6, 7, 8})
		t.Assert(s1.Contains(3), true)
		t.Assert(s1.Contains(5), false)
		t.Assert(s2.Contains(3), false)
		t.Assert(s2.Contains(5), true)
	})
}

func TestStrSet_Merge(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s2 := dset.NewStrSet()
		s1.Add("1", "2", "3")
		s2.Add("3", "4", "5")
		s3 := s1.Merge(s2)
		t.Assert(s3.Contains("1"), true)
		t.Assert(s3.Contains("6"), false)
		t.Assert(s3.Contains("4"), true)
		t.Assert(s3.Contains("5"), true)
	})
}

func TestNewStrSetFrom(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		t.Assert(s1.Contains("b"), true)
		t.Assert(s1.Contains("d"), false)
	})
}

func TestStrSet_Join(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, "b"), true)
		t.Assert(strings.Contains(str1, "d"), false)
	})

	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSet()
		s1.Add("a", `"b"`, `\c`)
		str1 := s1.Join(",")
		t.Assert(strings.Contains(str1, `"b"`), true)
		t.Assert(strings.Contains(str1, `\c`), true)
		t.Assert(strings.Contains(str1, `a`), true)
	})
}

func TestStrSet_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "b"), true)
		t.Assert(strings.Contains(str1, "d"), false)
	})

	dtest.C(t, func(t *dtest.T) {
		s1 := dset.New(true)
		s1.Add("a", "a2", "b", "c")
		str1 := s1.String()
		t.Assert(strings.Contains(str1, "["), true)
		t.Assert(strings.Contains(str1, "]"), true)
		t.Assert(strings.Contains(str1, "a2"), true)
	})
}

func TestStrSet_Sum(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		s2 := dset.NewIntSetFrom([]int{2, 3, 4}, true)
		t.Assert(s1.Sum(), 0)
		t.Assert(s2.Sum(), 9)
	})
}

func TestStrSet_Size(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		t.Assert(s1.Size(), 3)

	})
}

func TestStrSet_Remove(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := dset.NewStrSetFrom([]string{"a", "b", "c"}, true)
		s1.Remove("b")
		t.Assert(s1.Contains("b"), false)
		t.Assert(s1.Contains("c"), true)
	})
}

func TestStrSet_Pop(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a := []string{"a", "b", "c", "d"}
		s := dset.NewStrSetFrom(a, true)
		t.Assert(s.Size(), 4)
		t.AssertIN(s.Pop(), a)
		t.Assert(s.Size(), 3)
		t.AssertIN(s.Pop(), a)
		t.Assert(s.Size(), 2)
	})
}

func TestStrSet_Pops(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a := []string{"a", "b", "c", "d"}
		s := dset.NewStrSetFrom(a, true)
		array := s.Pops(2)
		t.Assert(len(array), 2)
		t.Assert(s.Size(), 2)
		t.AssertIN(array, a)
		t.Assert(s.Pops(0), nil)
		t.AssertIN(s.Pops(2), a)
		t.Assert(s.Size(), 0)
	})

	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet(true)
		a := []string{"1", "2", "3", "4"}
		s.Add(a...)
		t.Assert(s.Size(), 4)
		t.Assert(s.Pops(-2), nil)
		t.AssertIN(s.Pops(-1), a)
	})
}

func TestStrSet_AddIfNotExist(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet(true)
		s.Add("1")
		t.Assert(s.Contains("1"), true)
		t.Assert(s.AddIfNotExist("1"), false)
		t.Assert(s.AddIfNotExist("2"), true)
		t.Assert(s.Contains("2"), true)
		t.Assert(s.AddIfNotExist("2"), false)
		t.Assert(s.Contains("2"), true)
	})
}

func TestStrSet_AddIfNotExistFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet(true)
		s.Add("1")
		t.Assert(s.Contains("1"), true)
		t.Assert(s.Contains("2"), false)
		t.Assert(s.AddIfNotExistFunc("2", func() bool { return false }), false)
		t.Assert(s.Contains("2"), false)
		t.Assert(s.AddIfNotExistFunc("2", func() bool { return true }), true)
		t.Assert(s.Contains("2"), true)
		t.Assert(s.AddIfNotExistFunc("2", func() bool { return true }), false)
		t.Assert(s.Contains("2"), true)
	})
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet(true)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFunc("1", func() bool {
				time.Sleep(100 * time.Millisecond)
				return true
			})
			t.Assert(r, false)
		}()
		s.Add("1")
		wg.Wait()
	})
}

func TestStrSet_AddIfNotExistFuncLock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dset.NewStrSet(true)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFuncLock("1", func() bool {
				time.Sleep(500 * time.Millisecond)
				return true
			})
			t.Assert(r, true)
		}()
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			r := s.AddIfNotExistFuncLock("1", func() bool {
				return true
			})
			t.Assert(r, false)
		}()
		wg.Wait()
	})
}

func TestStrSet_Json(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := []string{"a", "b", "d", "c"}
		a1 := dset.NewStrSetFrom(s1)
		b1, err1 := json.Marshal(a1)
		b2, err2 := json.Marshal(s1)
		t.Assert(len(b1), len(b2))
		t.Assert(err1, err2)

		a2 := dset.NewStrSet()
		err2 = json.UnmarshalUseNumber(b2, &a2)
		t.Assert(err2, nil)
		t.Assert(a2.Contains("a"), true)
		t.Assert(a2.Contains("b"), true)
		t.Assert(a2.Contains("c"), true)
		t.Assert(a2.Contains("d"), true)
		t.Assert(a2.Contains("e"), false)

		var a3 dset.StrSet
		err := json.UnmarshalUseNumber(b2, &a3)
		t.Assert(err, nil)
		t.Assert(a3.Contains("a"), true)
		t.Assert(a3.Contains("b"), true)
		t.Assert(a3.Contains("c"), true)
		t.Assert(a3.Contains("d"), true)
		t.Assert(a3.Contains("e"), false)
	})
}

func TestStrSet_Walk(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			set    dset.StrSet
			names  = g.SliceStr{"user", "user_detail"}
			prefix = "gf_"
		)
		set.Add(names...)
		// Add prefix for given table names.
		set.Walk(func(item string) string {
			return prefix + item
		})
		t.Assert(set.Size(), 2)
		t.Assert(set.Contains("gf_user"), true)
		t.Assert(set.Contains("gf_user_detail"), true)
	})
}

func TestStrSet_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Set  *dset.StrSet
	}
	// JSON
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  []byte(`["1","2","3"]`),
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("1"), true)
		t.Assert(v.Set.Contains("2"), true)
		t.Assert(v.Set.Contains("3"), true)
		t.Assert(v.Set.Contains("4"), false)
	})
	// Map
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := gconv.Struct(g.Map{
			"name": "john",
			"set":  g.SliceStr{"1", "2", "3"},
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Set.Size(), 3)
		t.Assert(v.Set.Contains("1"), true)
		t.Assert(v.Set.Contains("2"), true)
		t.Assert(v.Set.Contains("3"), true)
		t.Assert(v.Set.Contains("4"), false)
	})
}
