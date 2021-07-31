package dmlock_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dmlock"
	"github.com/osgochina/donkeygo/test/dtest"
	"sync"
	"testing"
	"time"
)

func Test_Locker_Lock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		key := "testLock"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.Lock(key)
			array.Append(1)
			dmlock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		dmlock.Remove(key)
	})

	dtest.C(t, func(t *dtest.T) {
		key := "testLock"
		array := darray.New(true)
		lock := dmlock.New()
		go func() {
			lock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			lock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			lock.Lock(key)
			array.Append(1)
			lock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		lock.Clear()
	})

}

func Test_Locker_TryLock(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		key := "testTryLock"
		array := darray.New(true)
		go func() {
			dmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			dmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(150 * time.Millisecond)
			if dmlock.TryLock(key) {
				array.Append(1)
				dmlock.Unlock(key)
			}
		}()
		go func() {
			time.Sleep(400 * time.Millisecond)
			if dmlock.TryLock(key) {
				array.Append(1)
				dmlock.Unlock(key)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

}

func Test_Locker_LockFunc(t *testing.T) {
	//no expire
	dtest.C(t, func(t *dtest.T) {
		key := "testLockFunc"
		array := darray.New(true)
		go func() {
			dmlock.LockFunc(key, func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			}) //
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.LockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1) //
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
func Test_Locker_TryLockFunc(t *testing.T) {
	//no expire
	dtest.C(t, func(t *dtest.T) {
		key := "testTryLockFunc"
		array := darray.New(true)
		go func() {
			dmlock.TryLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			dmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			dmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Multiple_Goroutine(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		ch := make(chan struct{}, 0)
		num := 1000
		wait := sync.WaitGroup{}
		wait.Add(num)
		for i := 0; i < num; i++ {
			go func() {
				defer wait.Done()
				<-ch
				dmlock.Lock("test")
				defer dmlock.Unlock("test")
				time.Sleep(time.Millisecond)
			}()
		}
		close(ch)
		wait.Wait()
	})

	dtest.C(t, func(t *dtest.T) {
		ch := make(chan struct{}, 0)
		num := 100
		wait := sync.WaitGroup{}
		wait.Add(num * 2)
		for i := 0; i < num; i++ {
			go func() {
				defer wait.Done()
				<-ch
				dmlock.Lock("test")
				defer dmlock.Unlock("test")
				time.Sleep(time.Millisecond)
			}()
		}
		for i := 0; i < num; i++ {
			go func() {
				defer wait.Done()
				<-ch
				dmlock.RLock("test")
				defer dmlock.RUnlock("test")
				time.Sleep(time.Millisecond)
			}()
		}
		close(ch)
		wait.Wait()
	})
}
