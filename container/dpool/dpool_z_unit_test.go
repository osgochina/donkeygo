// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dpool_test

import (
	"errors"
	"github.com/osgochina/donkeygo/container/dpool"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

var nf dpool.NewFunc = func() (i interface{}, e error) {
	return "hello", nil
}

type testStruct struct {
	name string
}

var newf = func() interface{} {
	return "hello"
}

var newf2 = func() interface{} {
	return &testStruct{name: "hello"}
}

var assertIndex int = 0
var ef dpool.ExpireFunc = func(i interface{}) {
	assertIndex++
	dtest.Assert(i, assertIndex)
}

func Test_Gpool(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		//
		//expire = 0
		p1 := dpool.New(0, nf)
		p1.Put(1)
		p1.Put(2)
		time.Sleep(1 * time.Second)
		//test won't be timeout
		v1, err1 := p1.Get()
		t.Assert(err1, nil)
		t.AssertIN(v1, d.Slice{1, 2})
		//test clear
		p1.Clear()
		t.Assert(p1.Size(), 0)
		//test newFunc
		v1, err1 = p1.Get()
		t.Assert(err1, nil)
		t.Assert(v1, "hello")
		//put data again
		p1.Put(3)
		p1.Put(4)
		v1, err1 = p1.Get()
		t.Assert(err1, nil)
		t.AssertIN(v1, d.Slice{3, 4})
		//test close
		p1.Close()
		v1, err1 = p1.Get()
		t.Assert(err1, nil)
		t.Assert(v1, "hello")
	})

	dtest.C(t, func(t *dtest.T) {
		//
		//expire > 0
		p2 := dpool.New(2*time.Second, nil, ef)
		for index := 0; index < 10; index++ {
			p2.Put(index)
		}
		t.Assert(p2.Size(), 10)
		v2, err2 := p2.Get()
		t.Assert(err2, nil)
		t.Assert(v2, 0)
		//test timeout expireFunc
		time.Sleep(3 * time.Second)
		v2, err2 = p2.Get()
		t.Assert(err2, errors.New("pool is empty"))
		t.Assert(v2, nil)
		//test close expireFunc
		for index := 0; index < 10; index++ {
			p2.Put(index)
		}
		t.Assert(p2.Size(), 10)
		v2, err2 = p2.Get()
		t.Assert(err2, nil)
		t.Assert(v2, 0)
		assertIndex = 0
		p2.Close()
		time.Sleep(3 * time.Second)
	})

	dtest.C(t, func(t *dtest.T) {
		//
		//expire < 0
		p3 := dpool.New(-1, nil)
		v3, err3 := p3.Get()
		t.Assert(err3, errors.New("pool is empty"))
		t.Assert(v3, nil)
	})
}

func Test_SyncPool(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		//
		//expire = 0
		p1 := dpool.NewSyncPool(newf, nil)
		p1.Put(1)
		p1.Put(2)
		time.Sleep(1 * time.Second)
		//test won't be timeout
		v1 := p1.Get()
		t.AssertIN(v1, d.Slice{1, 2})

		//test newFunc
		v1 = p1.Get()
		t.Assert(v1, "hello")
		////put data again
		p1.Put(3)
		p1.Put(4)
		v1 = p1.Get()
		v1 = p1.Get()
		t.AssertIN(v1, d.Slice{3, 4})
		v3 := p1.Get()
		t.Assert(v3, "hello")
	})

	dtest.C(t, func(t *dtest.T) {
		p3 := dpool.NewSyncPool(newf2, func(i interface{}) {
			t := i.(*testStruct)
			t.name = "hello2"
		})
		v3 := p3.Get()
		t.Assert(v3.(*testStruct).name, "hello2")
	})
}
