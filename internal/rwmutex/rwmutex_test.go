// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package rwmutex_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/internal/rwmutex"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestRwmutexIsSafe(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		lock := rwmutex.New()
		t.Assert(lock.IsSafe(), false)

		lock = rwmutex.New(false)
		t.Assert(lock.IsSafe(), false)

		lock = rwmutex.New(false, false)
		t.Assert(lock.IsSafe(), false)

		lock = rwmutex.New(true, false)
		t.Assert(lock.IsSafe(), true)

		lock = rwmutex.New(true, true)
		t.Assert(lock.IsSafe(), true)

		lock = rwmutex.New(true)
		t.Assert(lock.IsSafe(), true)
	})
}

func TestSafeRwmutex(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		safeLock := rwmutex.New(true)
		array := darray.New(true)

		go func() {
			safeLock.Lock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			safeLock.Unlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			safeLock.Lock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			safeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(80 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 4)
	})
}

func TestSafeReaderRwmutex(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		safeLock := rwmutex.New(true)
		array := darray.New(true)

		go func() {
			safeLock.RLock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			safeLock.RUnlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			safeLock.RLock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			safeLock.RUnlock()
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			safeLock.Lock()
			array.Append(1)
			safeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 4)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 6)
	})
}

func TestUnsafeRwmutex(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		unsafeLock := rwmutex.New()
		array := darray.New(true)

		go func() {
			unsafeLock.Lock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			unsafeLock.Unlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			unsafeLock.Lock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			unsafeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 4)
	})
}
