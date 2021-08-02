// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Job Operations

package dtimer_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestJob_Start_Stop_Close(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		job := timer.Add(200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
		job.Stop()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
		job.Start()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		job.Close()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)

		t.Assert(job.Status(), dtimer.StatusClosed)
	})
}

func TestJob_Singleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		job := timer.Add(200*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		t.Assert(job.IsSingleton(), false)
		job.SetSingleton(true)
		t.Assert(job.IsSingleton(), true)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestJob_SetTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		job := timer.Add(200*time.Millisecond, func() {
			array.Append(1)
		})
		job.SetTimes(2)
		//job.IsSingleton()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestJob_Run(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timer := New()
		array := darray.New(true)
		job := timer.Add(1000*time.Millisecond, func() {
			array.Append(1)
		})
		job.Job()()
		t.Assert(array.Len(), 1)
	})
}
