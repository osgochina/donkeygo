package dmutex_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/os/dmutex"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestMutex_RUnLock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		for index := 0; index < 1000; index++ {
			go func() {
				mu.RLockFunc(func() {
					time.Sleep(200 * time.Millisecond)
				})
			}()
		}
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsRLocked(), true)
		t.Assert(mu.IsLocked(), true)
		t.Assert(mu.IsWLocked(), false)

		for index := 0; index < 1000; index++ {
			go func() {
				mu.RUnlock()
			}()
		}
		time.Sleep(300 * time.Millisecond)
		t.Assert(mu.IsRLocked(), false)
	})

	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		mu.RLock()
		go func() {
			mu.Lock()
			time.Sleep(300 * time.Millisecond)
			mu.Unlock()
		}()
		time.Sleep(100 * time.Millisecond)
		mu.RUnlock()
		t.Assert(mu.IsRLocked(), false)
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsLocked(), true)
		time.Sleep(400 * time.Millisecond)
		t.Assert(mu.IsLocked(), false)
	})

}

func Test_Mutex_IsLocked(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		go func() {
			mu.LockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsLocked(), true)
		t.Assert(mu.IsWLocked(), true)
		t.Assert(mu.IsRLocked(), false)
		time.Sleep(300 * time.Millisecond)
		t.Assert(mu.IsLocked(), false)
		t.Assert(mu.IsWLocked(), false)

		go func() {
			mu.RLockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(mu.IsRLocked(), true)
		t.Assert(mu.IsLocked(), true)
		t.Assert(mu.IsWLocked(), false)
		time.Sleep(300 * time.Millisecond)
		t.Assert(mu.IsRLocked(), false)
	})
}

func Test_Mutex_Unlock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		array := darray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()

		go func() {
			time.Sleep(200 * time.Millisecond)
			mu.Unlock()
			mu.Unlock()
			mu.Unlock()
			mu.Unlock()
		}()

		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Mutex_LockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		array := darray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Mutex_TryLockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		array := darray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(400 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Mutex_RLockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		array := darray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	dtest.C(t, func(t *dtest.T) {
		mu := dmutex.New()
		array := darray.New(true)
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		t.Assert(array.Len(), 0)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Mutex_TryRLockFunc(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			mu    = dmutex.New()
			array = darray.New(true)
		)
		// First writing lock
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				dlog.Println("lock1 done")
				time.Sleep(2000 * time.Millisecond)
			})
		}()
		// This goroutine never gets the lock.
		go func() {
			time.Sleep(1000 * time.Millisecond)
			mu.TryRLockFunc(func() {
				array.Append(1)
			})
		}()
		for index := 0; index < 1000; index++ {
			go func() {
				time.Sleep(4000 * time.Millisecond)
				mu.TryRLockFunc(func() {
					array.Append(1)
				})
			}()
		}
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(2000 * time.Millisecond)
		t.Assert(array.Len(), 1001)
	})
}
