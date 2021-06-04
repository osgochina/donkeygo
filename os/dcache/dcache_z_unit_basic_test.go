// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package dcache_test

import (
	"context"
	"github.com/gogf/gf/util/guid"
	"github.com/osgochina/donkeygo/os/dcache"
	"github.com/osgochina/donkeygo/test/dtest"
	"math"
	"testing"
	"time"

	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/grpool"
)

func TestCache_GCache_Set(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		dcache.Set(1, 11, 0)
		defer dcache.Removes(g.Slice{1, 2, 3})
		v, _ := dcache.Get(1)
		t.Assert(v, 11)
		b, _ := dcache.Contains(1)
		t.Assert(b, true)
	})
}

func TestCache_Set(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		c := dcache.New()
		defer c.Close()
		t.Assert(c.Set(1, 11, 0), nil)
		v, _ := c.Get(1)
		t.Assert(v, 11)
		b, _ := c.Contains(1)
		t.Assert(b, true)
	})
}

func TestCache_GetVar(t *testing.T) {
	c := dcache.New()
	defer c.Close()
	dtest.C(t, func(t *dtest.T) {
		t.Assert(c.Set(1, 11, 0), nil)
		v, _ := c.Get(1)
		t.Assert(v, 11)
		b, _ := c.Contains(1)
		t.Assert(b, true)
	})
	dtest.C(t, func(t *dtest.T) {
		v, _ := c.GetVar(1)
		t.Assert(v.Int(), 11)
		v, _ = c.GetVar(2)
		t.Assert(v.Int(), 0)
		t.Assert(v.IsNil(), true)
		t.Assert(v.IsEmpty(), true)
	})
}

func TestCache_Set_Expire(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		t.Assert(cache.Set(2, 22, 100*time.Millisecond), nil)
		v, _ := cache.Get(2)
		t.Assert(v, 22)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(2)
		t.Assert(v, nil)
		time.Sleep(3 * time.Second)
		n, _ := cache.Size()
		t.Assert(n, 0)
		t.Assert(cache.Close(), nil)
	})

	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		t.Assert(cache.Set(1, 11, 100*time.Millisecond), nil)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(1)
		t.Assert(v, nil)
	})
}

func TestCache_Update_GetExpire(t *testing.T) {
	// gcache
	dtest.C(t, func(t *dtest.T) {
		key := guid.S()
		dcache.Set(key, 11, 3*time.Second)
		expire1, _ := dcache.GetExpire(key)
		dcache.Update(key, 12)
		expire2, _ := dcache.GetExpire(key)
		v, _ := dcache.GetVar(key)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
	// dcache.Cache
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.Set(1, 11, 3*time.Second)
		expire1, _ := cache.GetExpire(1)
		cache.Update(1, 12)
		expire2, _ := cache.GetExpire(1)
		v, _ := cache.GetVar(1)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
}

func TestCache_UpdateExpire(t *testing.T) {
	// gcache
	dtest.C(t, func(t *dtest.T) {
		key := guid.S()
		dcache.Set(key, 11, 3*time.Second)
		defer dcache.Remove(key)
		oldExpire, _ := dcache.GetExpire(key)
		newExpire := 10 * time.Second
		dcache.UpdateExpire(key, newExpire)
		e, _ := dcache.GetExpire(key)
		t.AssertNE(e, oldExpire)
		e, _ = dcache.GetExpire(key)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
	// dcache.Cache
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.Set(1, 11, 3*time.Second)
		oldExpire, _ := cache.GetExpire(1)
		newExpire := 10 * time.Second
		cache.UpdateExpire(1, newExpire)
		e, _ := cache.GetExpire(1)
		t.AssertNE(e, oldExpire)

		e, _ = cache.GetExpire(1)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
}

func TestCache_Keys_Values(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		c := dcache.New()
		for i := 0; i < 10; i++ {
			t.Assert(c.Set(i, i*10, 0), nil)
		}
		var (
			keys, _   = c.Keys()
			values, _ = c.Values()
		)
		t.Assert(len(keys), 10)
		t.Assert(len(values), 10)
		t.AssertIN(0, keys)
		t.AssertIN(90, values)
	})
}

func TestCache_LRU(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New(2)
		for i := 0; i < 10; i++ {
			cache.Set(i, i, 0)
		}
		n, _ := cache.Size()
		t.Assert(n, 10)
		v, _ := cache.Get(6)
		t.Assert(v, 6)
		time.Sleep(4 * time.Second)
		n, _ = cache.Size()
		t.Assert(n, 2)
		v, _ = cache.Get(6)
		t.Assert(v, 6)
		v, _ = cache.Get(1)
		t.Assert(v, nil)
		t.Assert(cache.Close(), nil)
	})
}

func TestCache_LRU_expire(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New(2)
		t.Assert(cache.Set(1, nil, 1000), nil)
		n, _ := cache.Size()
		t.Assert(n, 1)
		v, _ := cache.Get(1)

		t.Assert(v, nil)
	})
}

func TestCache_SetIfNotExist(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.SetIfNotExist(1, 11, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
		cache.SetIfNotExist(1, 22, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
		cache.SetIfNotExist(2, 22, 0)
		v, _ = cache.Get(2)
		t.Assert(v, 22)

		dcache.Removes(g.Slice{1, 2, 3})
		dcache.SetIfNotExist(1, 11, 0)
		v, _ = dcache.Get(1)
		t.Assert(v, 11)
		dcache.SetIfNotExist(1, 22, 0)
		v, _ = dcache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_Sets(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		dcache.Removes(g.Slice{1, 2, 3})
		dcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSet(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.GetOrSet(1, 11, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
		cache.GetOrSet(1, 111, 0)

		v, _ = cache.Get(1)
		t.Assert(v, 11)
		dcache.Removes(g.Slice{1, 2, 3})
		dcache.GetOrSet(1, 11, 0)

		v, _ = cache.Get(1)
		t.Assert(v, 11)

		dcache.GetOrSet(1, 111, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.GetOrSetFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		cache.GetOrSetFunc(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		dcache.Removes(g.Slice{1, 2, 3})

		dcache.GetOrSetFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		dcache.GetOrSetFunc(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFuncLock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		cache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		dcache.Removes(g.Slice{1, 2, 3})
		dcache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		dcache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_Clear(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		cache.Clear()
		n, _ := cache.Size()
		t.Assert(n, 0)
	})
}

func TestCache_SetConcurrency(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		pool := grpool.New(4)
		go func() {
			for {
				pool.Add(func() {
					cache.SetIfNotExist(1, 11, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			//t.Log("first part end")
		}

		go func() {
			for {
				pool.Add(func() {
					cache.SetIfNotExist(1, nil, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			//t.Log("second part end")
		}
	})
}

func TestCache_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		{
			cache := dcache.New()
			cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			b, _ := cache.Contains(1)
			t.Assert(b, true)
			v, _ := cache.Get(1)
			t.Assert(v, 11)
			data, _ := cache.Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := cache.Size()
			t.Assert(n, 2)
			keys, _ := cache.Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs, _ := cache.KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values, _ := cache.Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1, _ := cache.Remove(1)
			t.Assert(removeData1, 11)
			n, _ = cache.Size()
			t.Assert(n, 1)
			cache.Removes(g.Slice{2})
			n, _ = cache.Size()
			t.Assert(n, 0)
		}

		dcache.Remove(g.Slice{1, 2, 3}...)
		{
			dcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			b, _ := dcache.Contains(1)
			t.Assert(b, true)
			v, _ := dcache.Get(1)
			t.Assert(v, 11)
			data, _ := dcache.Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := dcache.Size()
			t.Assert(n, 2)
			keys, _ := dcache.Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs, _ := dcache.KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values, _ := dcache.Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1, _ := dcache.Remove(1)
			t.Assert(removeData1, 11)
			n, _ = dcache.Size()
			t.Assert(n, 1)
			dcache.Removes(g.Slice{2})
			n, _ = dcache.Size()
			t.Assert(n, 0)
		}
	})
}

func TestCache_Ctx(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cache := dcache.New()
		cache.Ctx(context.Background()).Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		b, _ := cache.Contains(1)
		t.Assert(b, true)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
	})
}
