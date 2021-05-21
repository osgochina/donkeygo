// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Timer Operations

package dtimer_test

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func New() *dtimer.Timer {
	return dtimer.NewTimer(10, 10*time.Millisecond)
}

func TestTimer_Add_Close(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		//fmt.Println("start", time.Now())
		timer.Add(200*time.Millisecond, func() {
			//fmt.Println("entry1", time.Now())
			array.Append(1)
		})
		timer.Add(200*time.Millisecond, func() {
			//fmt.Println("entry2", time.Now())
			array.Append(1)
		})
		timer.Add(400*time.Millisecond, func() {
			//fmt.Println("entry3", time.Now())
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 5)
		timer.Close()
		time.Sleep(250 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestTimer_Start_Stop_Close(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.Add(200*time.Millisecond, func() {
			//glog.Println("add...")
			array.Append(1)
		})
		t.Assert(array.Len(), 0)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.Stop()
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.Start()
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		timer.Close()
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_Reset(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		glog.Printf("start time:%d", time.Now().Unix())
		singleton := timer.AddSingleton(2*time.Second, func() {
			timestamp := time.Now().Unix()
			glog.Println(timestamp)
			array.Append(timestamp)
		})
		time.Sleep(5 * time.Second)
		glog.Printf("reset time:%d", time.Now().Unix())
		singleton.Reset()
		time.Sleep(10 * time.Second)
		t.Assert(array.Len(), 6)
	})
}

func TestTimer_AddSingleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.AddSingleton(200*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_AddOnce(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.AddOnce(200*time.Millisecond, func() {
			array.Append(1)
		})
		timer.AddOnce(200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		timer.Close()
		time.Sleep(250 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestTimer_AddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.AddTimes(200*time.Millisecond, 2, func() {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_DelayAdd(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.DelayAdd(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddEntry(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.DelayAddEntry(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
		}, false, 100, dtimer.StatusReady)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddSingleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.DelayAddSingleton(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddOnce(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.DelayAddOnce(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.DelayAddTimes(200*time.Millisecond, 500*time.Millisecond, 2, func() {
			array.Append(1)
		})
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 2)

		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_AddLessThanInterval(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := dtimer.NewTimer(10, 100*time.Millisecond)
		array := darray.New(true)
		timer.Add(20*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(110 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(110 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_AddLeveledEntry1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		//glog.Println("start")
		timer.DelayAdd(1000*time.Millisecond, 1000*time.Millisecond, func() {
			//glog.Println("add")
			array.Append(1)
		})
		time.Sleep(1500 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1300 * time.Millisecond)
		//glog.Println("check")
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_Exit(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		timer.Add(200*time.Millisecond, func() {
			array.Append(1)
			dtimer.Exit()
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}
