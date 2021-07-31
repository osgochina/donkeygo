package dmlock_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dmlock"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func Test_Locker_RLock(t *testing.T) {
	//RLock before Lock
	dtest.C(t, func(t *dtest.T) {
		key := "testRLockBeforeLock"
		array := darray.New(true)
		go func() {
			dmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.Lock(key)
			array.Append(1)
			dmlock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLock
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeRLock"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.RLock(key)
			array.Append(1)
			dmlock.RUnlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLocks
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeRLocks"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.RUnlock(key)
		}()
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLock(t *testing.T) {
	//Lock before TryRLock
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeTryRLock"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if dmlock.TryRLock(key) {
				array.Append(1)
				dmlock.RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

	//Lock before TryRLocks
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeTryRLocks"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if dmlock.TryRLock(key) {
				array.Append(1)
				dmlock.RUnlock(key)
			}
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			if dmlock.TryRLock(key) {
				array.Append(1)
				dmlock.RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Locker_RLockFunc(t *testing.T) {
	//RLockFunc before Lock
	dtest.C(t, func(t *dtest.T) {
		key := "testRLockFuncBeforeLock"
		array := darray.New(true)
		go func() {
			dmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.Lock(key)
			array.Append(1)
			dmlock.Unlock(key)
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLockFunc
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeRLockFunc"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.RLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLockFuncs
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeRLockFuncs"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLockFunc(t *testing.T) {
	//Lock before TryRLockFunc
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeTryRLockFunc"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

	//Lock before TryRLockFuncs
	dtest.C(t, func(t *dtest.T) {
		key := "testLockBeforeTryRLockFuncs"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			dmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
