// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package functions

package dtimer_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestSetTimeout(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.SetTimeout(200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestSetInterval(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.SetInterval(300*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func TestAddEntry(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.AddEntry(200*time.Millisecond, func() {
			array.Append(1)
		}, false, 2, dtimer.StatusReady)
		time.Sleep(1100 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestAddSingleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.AddSingleton(200*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(1100 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestAddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.AddTimes(200*time.Millisecond, 2, func() {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestDelayAdd(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.DelayAdd(500*time.Millisecond, 500*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddEntry(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.DelayAddEntry(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
		}, false, 2, dtimer.StatusReady)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestDelayAddSingleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.DelayAddSingleton(500*time.Millisecond, 500*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddOnce(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.DelayAddOnce(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := darray.New(true)
		dtimer.DelayAddTimes(500*time.Millisecond, 500*time.Millisecond, 2, func() {
			array.Append(1)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1500 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
